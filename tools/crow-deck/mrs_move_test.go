package main

import (
	"strconv"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// stageName is the current step of a stage, or "done".
func stageName(s mrStage) string {
	if s.done() {
		return "done"
	}
	return s.Steps[s.Cur].Name
}

func TestStageWalksTheSelfReview(t *testing.T) {
	own := mrRow{Role: "AUTHOR", Pipeline: "passed"}
	if s := mrStageOf(own, nil, testNow()); stageName(s) != "review" || s.Keys != "r" || s.Who != moveOthers {
		t.Errorf("no review: %+v", s)
	}
	d := reviewDetail("running", "")
	d.Review.Mode = "self"
	if s := mrStageOf(own, d, testNow()); s.Who != moveCrow || s.Reason != "review running" {
		t.Errorf("running review: %+v", s)
	}
	d = reviewDetail("done", "")
	d.Review.Mode = "self"
	s := mrStageOf(own, d, testNow())
	if stageName(s) != "decide" || s.Who != moveYou || s.Keys != "a x e c" || s.Reason != "1 decision pending" || s.Steps[s.Cur].Sub != "1 of 2 decided" {
		t.Errorf("decide: %+v", s)
	}
	d.Findings[0].State = "accepted"
	if s := mrStageOf(own, d, testNow()); stageName(s) != "apply" || s.Keys != "A" || s.Who != moveYou {
		t.Errorf("apply: %+v", s)
	}
	d.Apply = &mrApply{State: "running", StartedAt: testNow().Add(-3 * time.Minute).Format(time.RFC3339)}
	if s := mrStageOf(own, d, testNow()); stageName(s) != "apply" || s.Who != moveCrow {
		t.Errorf("apply running: %+v", s)
	}
	for i := range d.Findings {
		d.Findings[i].Applied, d.Findings[i].AppliedFor = "fixed", [2]string{d.Findings[i].State, d.Findings[i].Note}
	}
	d.Apply = &mrApply{State: "done", Commit: "abc1234def"}
	if s := mrStageOf(own, d, testNow()); stageName(s) != "push" || s.Keys != "P" || s.Reason != "push abc1234d" {
		t.Errorf("push: %+v", s)
	}
	d.Apply.Pushed = true
	s = mrStageOf(own, d, testNow())
	if !s.done() || s.Who != moveDone || s.group() != 2 || s.Reason != "done, pushed abc1234d" {
		t.Errorf("done: %+v", s)
	}
	own.Pipeline = "failed"
	if s := mrStageOf(own, d, testNow()); s.Who != moveYou || !s.Problem {
		t.Errorf("a red pipeline stays the owner's move: %+v", s)
	}
}

func TestStageThreadsResolveAndPeerPost(t *testing.T) {
	d := threadsDetail()
	for i := range d.Findings {
		d.Findings[i].State = "accepted"
	}
	d.Findings[0].Applied, d.Findings[0].AppliedFor = "fixed", [2]string{"accepted", ""}
	d.Apply = &mrApply{State: "done", Commit: "abc1234def", Pushed: true}
	s := mrStageOf(mrRow{Role: "AUTHOR", ThreadsOpen: 2}, d, testNow())
	if stageName(s) != "resolve" || s.Reason != "2 threads to resolve" || s.Steps[len(s.Steps)-1].Text != "done: ask sam" {
		t.Errorf("threads: %+v", s)
	}
	p := reviewDetail("done", "")
	p.Findings[0].State = "accepted"
	r := mrRow{Role: "REVIEWER", Author: "sam"}
	if s := mrStageOf(r, p, testNow()); stageName(s) != "post" || s.Reason != "2 comments to post" {
		t.Errorf("peer post: %+v", s)
	}
	p.Findings[0].Posted, p.Findings[1].Posted = true, true
	if s := mrStageOf(r, p, testNow()); !s.done() || s.Reason != "done, reviewed" {
		t.Errorf("peer done: %+v", s)
	}
	if s := mrStageOf(r, reviewDetail("failed", "boom"), testNow()); !s.Problem || s.Keys != "r" {
		t.Errorf("a failed review is a problem: %+v", s)
	}
}

func TestStageLineMarksTheCurrentStep(t *testing.T) {
	withANSI(t)
	d := reviewDetail("done", "")
	d.Review.Mode = "self"
	m := selfReviewModel(t, d)
	out := m.View()
	for _, want := range []string{"✓ review ─ ● decide ─ ○ apply ─ ○ commit ─ ○ push ─ ○ done", "next a x e c", "▶ ● 2  decide each finding", "1 of 2 decided"} {
		if !hasText(out, want) {
			t.Errorf("stage line lacks %q", want)
		}
	}
}

func TestListGroupsByWhoseMove(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	var refs []string
	for _, r := range visibleMRs(m.fx, m.mode) {
		refs = append(refs, r.Ref)
	}
	if strings.Join(refs, " ") != "!6 !5 #12" {
		t.Errorf("order %v, want your move first", refs)
	}
	out := m.View()
	for _, want := range []string{"Your move  2", "Waiting on others  1", "● !6", "○ #12", "2/3 your move"} {
		if !hasText(out, want) {
			t.Errorf("list lacks %q", want)
		}
	}
	// A pushed self review moves its MR to Done, out of the badge.
	d := applyDetail()
	for i := range d.Findings {
		d.Findings[i].Applied, d.Findings[i].AppliedFor = "fixed", [2]string{d.Findings[i].State, d.Findings[i].Note}
	}
	d.Apply = &mrApply{State: "done", Commit: "abc1234def", Pushed: true}
	s := testMRs()
	s.MRs[2].Pipeline = "passed"
	m = send(m, srcMsg{src: srcMRs, val: s, at: testNow()})
	m = send(m, detailMsg{key: "grp/atlas!6", d: d, at: testNow()})
	m.mrToggle = true
	out = m.View()
	for _, want := range []string{"Done  1", "✓ !6", "done, pushed abc1234d"} {
		if !hasText(out, want) {
			t.Errorf("done list lacks %q", want)
		}
	}
	if m.mrKey != "grp/atlas!6" || visibleMRs(m.fx, m.mode)[m.mrSel].Ref != "!6" {
		t.Errorf("the cursor must follow !6 into Done, on %q", m.mrKey)
	}
}

func TestBadgeCountsYourMoveAndDotsProblems(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	if n, bad := m.fx.mrsNeedYou(); n != 2 || !bad {
		t.Errorf("badge %d %v, want 2 and a dot for the red pipeline", n, bad)
	}
	if !hasText(m.View(), "MRs 2 ●") {
		t.Error("the tab must show the count and the dot")
	}
	s := testMRs()
	s.MRs[2].Pipeline = "passed"
	m = send(m, srcMsg{src: srcMRs, val: s, at: testNow()})
	if n, bad := m.fx.mrsNeedYou(); n != 2 || bad {
		t.Errorf("badge %d %v, want 2 without a dot", n, bad)
	}
	if hasText(m.View(), "MRs 2 ●") {
		t.Error("no dot without a problem")
	}
}

func TestDoneStateOnThePanel(t *testing.T) {
	withANSI(t)
	d := threadsDetail()
	for i := range d.Findings {
		d.Findings[i].State, d.Findings[i].Resolved, d.Findings[i].Posted = "accepted", true, true
	}
	d.Findings[0].Applied, d.Findings[0].AppliedFor = "fixed", [2]string{"accepted", ""}
	d.Apply = &mrApply{State: "done", Commit: "abc1234def", Pushed: true}
	m := threadsModel(t, d)
	s := testMRs()
	s.MRs[2].Pipeline = "passed"
	m = send(m, srcMsg{src: srcMRs, val: s, at: testNow()})
	out := m.View()
	for _, want := range []string{"Done · !6", "✓ done", "abc1234d pushed", "3/3 threads replied and resolved", "next: ask sam", "next MR"} {
		if !hasText(out, want) {
			t.Errorf("done panel lacks %q", want)
		}
	}
}

func TestPostDoneAndApplyStartRefetchAtOnce(t *testing.T) {
	m := reviewModel(t, reviewDetail("done", ""))
	m.live.details["grp/atlas!5"].busy = false
	next, cmd := m.Update(postDoneMsg{key: "grp/atlas!5", text: "posted", ok: true})
	m = next.(model)
	if cmd == nil || !m.live.details["grp/atlas!5"].busy {
		t.Fatal("a done post must refetch the MR at once")
	}
	m.live.details["grp/atlas!5"].busy = false
	if _, cmd := m.Update(reviewStartedMsg{key: "grp/atlas!5", text: "apply ap1 started", apply: true}); cmd == nil {
		t.Error("an apply start must refetch the MR at once")
	}
}

func TestCodeScrollsHunksFilesAndFindings(t *testing.T) {
	withANSI(t)
	d := reviewDetail("done", "")
	var diff strings.Builder
	for h := 0; h < 4; h++ {
		diff.WriteString("@@ -" + strconv.Itoa(h*100+1) + ",20 +" + strconv.Itoa(h*100+1) + ",20 @@\n")
		for i := 0; i < 20; i++ {
			diff.WriteString(" line\n")
		}
	}
	d.Files = append(d.Files, mrFile{Path: "b.py", Lines: parseDiff(diff.String())})
	d.Files[0], d.Files[1] = d.Files[1], d.Files[0]
	m := reviewModel(t, d)
	m.w, m.h, m.mrPane = 170, 30, mpDiff
	m.mrFile = 0 // b.py: no card, opens at the top
	if !hasText(m.View(), "↑ 0 above") {
		t.Fatalf("b.py must open at the top")
	}
	m = send(m, keyMsg("j"), keyMsg("j"))
	if !hasText(m.View(), "↑ 2 above") {
		t.Error("j must scroll the code by one row")
	}
	m = send(m, keyMsg("J"))
	if !hasText(m.View(), "↑ 21 above") {
		t.Error("J must put the next hunk at the top")
	}
	m = send(m, keyMsg("K"))
	if !hasText(m.View(), "↑ 0 above") {
		t.Error("K must go back to the previous hunk")
	}
	m = send(m, keyMsg("]"))
	if m.mrFile != 1 || !hasText(m.View(), "a.py +2 -1") {
		t.Error("] must open the next file")
	}
	// n goes to the next finding; past the last one, to the next MR.
	m = send(m, keyMsg("n"))
	if m.mrItem != 1 {
		t.Errorf("n: item %d", m.mrItem)
	}
	m = send(m, keyMsg("n"))
	if vis := visibleMRs(m.fx, m.mode); vis[m.mrSel].Ref == "!5" {
		t.Error("n after the last finding must go to the next MR")
	}
	m = send(m, tea.KeyMsg{Type: tea.KeyTab})
	if m.mrPane != mpCard {
		t.Errorf("tab: pane %d", m.mrPane)
	}
}

func TestPushedCardDropsTheKeycaps(t *testing.T) {
	withANSI(t)
	d := threadsDetail()
	for i := range d.Findings {
		d.Findings[i].State = "accepted"
	}
	d.Findings[0].Applied, d.Findings[0].AppliedFor = "fixed", [2]string{"accepted", ""}
	d.Apply = &mrApply{State: "done", Commit: "abc1234def", Pushed: true}
	m := threadsModel(t, d)
	out := m.View()
	if !hasText(out, "Decided F1 accepted") || hasText(out, "accept the fix") {
		t.Error("a pushed card must show one decided line, not the keycaps")
	}
	if hasText(out, "apply 0 decisions") {
		t.Error("the A row must go once apply is behind")
	}
}
