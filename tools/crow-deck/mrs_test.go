package main

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func testMRs() *mrsSnap {
	return &mrsSnap{
		Forges: map[string]string{"gitlab": "ok", "github": "ok"},
		MRs: []mrRow{
			{Project: "web", Name: "web", Scope: "personal", Forge: "github", Repo: "me/web", Number: 12, Ref: "#12", Title: "cache", Author: "me", Role: "AUTHOR", Pipeline: "passed", Updated: "2026-10-29T09:00:00Z"},
			{Project: "atlas", Name: "atlas", Scope: "work", Forge: "gitlab", Repo: "grp/atlas", Number: 5, Ref: "!5", Title: "retry", Author: "sam", Role: "REVIEWER", Updated: "2026-10-27T09:00:00Z"},
			{Project: "atlas", Name: "atlas", Scope: "work", Forge: "gitlab", Repo: "grp/atlas", Number: 6, Ref: "!6", Title: "seed", Author: "me", Role: "AUTHOR", ThreadsOpen: 2, Pipeline: "failed", Updated: "2026-10-28T09:00:00Z"},
		},
	}
}

const testDiff = "@@ -10,3 +10,4 @@ def f():\n a\n-b\n+c\n+d\n e\n"

func TestParseDiffNumbersLines(t *testing.T) {
	got := parseDiff(testDiff)
	want := []diffLine{
		{Kind: '@', Text: "@@ -10,3 +10,4 @@ def f():"},
		{Old: 10, New: 10, Kind: ' ', Text: "a"},
		{Old: 11, Kind: '-', Text: "b"},
		{New: 11, Kind: '+', Text: "c"},
		{New: 12, Kind: '+', Text: "d"},
		{Old: 12, New: 13, Kind: ' ', Text: "e"},
	}
	if len(got) != len(want) {
		t.Fatalf("%d lines, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestMROrderAndCounts(t *testing.T) {
	rows := mrOrder(testMRs().MRs)
	var refs []string
	for _, r := range rows {
		refs = append(refs, r.Ref)
	}
	// Projects alphabetical, own MRs before reviews.
	if strings.Join(refs, " ") != "!6 !5 #12" {
		t.Errorf("order %v", refs)
	}
	a, mine, th := mrCounts(rows)
	if a != 1 || mine != 2 || th != 2 {
		t.Errorf("counts %d %d %d", a, mine, th)
	}
}

func TestMRBlockers(t *testing.T) {
	b := mrBlockers(testNow(), mrOrder(testMRs().MRs))
	if len(b) != 2 || b[0].Item != "!6 atlas" || b[0].Why != "2 threads to answer" || b[1].Why != "sam waits for your review" {
		t.Errorf("blockers %+v", b)
	}
}

func liveMRModel(t *testing.T) (model, *[]string) {
	t.Helper()
	m := liveModel(t)
	m.page, m.mode = 0, "all"
	m.w, m.h = 170, 46
	var asked []string
	old := detailFetcher
	detailFetcher = func(r mrRow) (*mrDetail, error) {
		asked = append(asked, r.key())
		if r.Ref == "!6" {
			return nil, errors.New("forge down")
		}
		return mapDetail(rawDetail{}), nil
	}
	t.Cleanup(func() { detailFetcher = old })
	m = send(m, srcMsg{src: srcMRs, val: testMRs(), at: testNow()})
	m.page, m.mrPane = 2, mpList
	return m, &asked
}

func TestMRsPageLive(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	out := m.View()
	for _, want := range []string{"atlas", "web", "!6", "retry", "Your move", "review asked of you", "pipeline failed", "Waiting on others"} {
		if !hasText(out, want) {
			t.Errorf("view lacks %q", want)
		}
	}
	if m.fx.MRsAwaiting != 1 || m.fx.MRsThreads != 2 {
		t.Errorf("counts %d %d", m.fx.MRsAwaiting, m.fx.MRsThreads)
	}
	if len(m.fx.Blockers) != 2 {
		t.Errorf("home blockers %d, want 2 from the MRs", len(m.fx.Blockers))
	}
}

func TestMRsModeFiltersByScope(t *testing.T) {
	m, _ := liveMRModel(t)
	m.mode = "work"
	if v := visibleMRs(m.fx, m.mode); len(v) != 2 || v[0].Project != "atlas" {
		t.Errorf("work: %+v", v)
	}
	m.mode = "perso"
	if v := visibleMRs(m.fx, m.mode); len(v) != 1 || v[0].Ref != "#12" {
		t.Errorf("perso: %+v", v)
	}
}

func TestMRSelectionFetchesDetailOnce(t *testing.T) {
	m, asked := liveMRModel(t)
	m.page = 0
	next, cmd := m.Update(keyMsg("3"))
	m = next.(model)
	if cmd == nil {
		t.Fatal("opening the tab did not fetch the diff")
	}
	msg := cmd()
	if bm, ok := msg.(tea.BatchMsg); ok {
		for _, c := range bm {
			if c != nil {
				if dm, ok := c().(detailMsg); ok {
					msg = dm
				}
			}
		}
	}
	if len(*asked) != 1 || (*asked)[0] != "grp/atlas!6" {
		t.Fatalf("asked %v", *asked)
	}
	// A failing fetch shows its error and is not retried on every key.
	m = send(m, msg)
	if m.fx.MRErrs["grp/atlas!6"] != "forge down" {
		t.Errorf("errs %v", m.fx.MRErrs)
	}
	if _, c := m.Update(keyMsg("j")); c == nil {
		// moving to !5 starts its fetch
		t.Error("j did not fetch the next MR")
	}
	m.live.details["grp/atlas!6"].at = testNow().Add(-time.Hour)
	if cmd := m.wantDetail(); cmd == nil {
		t.Error("stale entry not refetched")
	}
}

func TestMRCursorFollowsItsMRAcrossRefresh(t *testing.T) {
	m, _ := liveMRModel(t)
	m = send(m, keyMsg("j"))
	if m.mrKey != "grp/atlas!5" {
		t.Fatalf("on %q", m.mrKey)
	}
	s := testMRs()
	s.MRs = append([]mrRow{{Project: "aaa", Name: "aaa", Scope: "work", Repo: "g/aaa", Number: 1, Ref: "!1", Title: "new", Role: "AUTHOR"}}, s.MRs...)
	m = send(m, srcMsg{src: srcMRs, val: s, at: testNow()})
	if v := visibleMRs(m.fx, m.mode); v[m.mrSel].Ref != "!5" {
		t.Errorf("cursor moved to %q", v[m.mrSel].Ref)
	}
}

func TestMRsForgeFailureShowsOnlyThatForge(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	s := testMRs()
	s.Forges["github"] = "github sync failed: auth"
	m = send(m, srcMsg{src: srcMRs, val: s, at: testNow()})
	if !hasText(m.View(), "github sync failed: auth") || !hasText(m.View(), "!6") {
		t.Error("the failing forge is not named next to the other forge's MRs")
	}
}

func TestMRsMockCards(t *testing.T) {
	withANSI(t)
	m := newModel(defaultConfig())
	m.page, m.w, m.h = 2, 170, 46
	m.mrSel = 1 // !220, a review (the list puts your move first: !213, !220, then !214)
	out := m.View()
	for _, want := range []string{"finding 1 of 3", "env var undocumented", "Decide F1", "a", "accept", "Then", "next finding", "post the comments", "pattern"} {
		if !hasText(out, want) {
			t.Errorf("review view lacks %q", want)
		}
	}
	m.mrSel = 0 // !213, own MR with a thread and a saved fix
	out = m.View()
	for _, want := range []string{"thread 1 of 2", "Proposed fix", "Test", "review the threads"} {
		if !hasText(out, want) {
			t.Errorf("thread view lacks %q", want)
		}
	}
	m = send(m, keyMsg("T"))
	if !hasText(m.View(), "No findings saved") {
		t.Error("T did not flip the card kind")
	}
}

func TestMRsActionKeysDoNotPretend(t *testing.T) {
	m := newModel(defaultConfig())
	m.page = 2
	for _, k := range []string{"a", "M"} {
		if n := send(m, keyMsg(k)); n.flash == "" {
			t.Errorf("%s: no flash", k)
		}
	}
}

func TestHomeBlockersComeFromMRs(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	m.page = 0
	out := m.View()
	for _, want := range []string{"You decide", "waits for your review", "2 threads to answer"} {
		if !hasText(out, want) {
			t.Errorf("home lacks %q", want)
		}
	}
	m = send(m, srcMsg{src: srcMRs, val: &mrsSnap{Forges: map[string]string{"gitlab": "ok", "github": "ok"}}, at: testNow()})
	if hasText(m.View(), "waits for your review") {
		t.Error("MR rows stay after the MRs are gone")
	}
}

func TestHandoffOverlayThenRun(t *testing.T) {
	calls := recordRuns(t, "ok", nil)
	m := newModel(defaultConfig())
	m.page = 2
	m = send(m, keyMsg("h"))
	if m.overlay != ovAct || !hasText(m.View(), "Open a repo session") {
		t.Fatalf("h must open the handoff overlay, overlay %d", m.overlay)
	}
	if len(*calls) != 0 {
		t.Fatal("nothing may run before enter")
	}
	m2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m2.(model).overlay != ovNone || cmd == nil {
		t.Fatal("enter must close the overlay and return the command")
	}
	cmd()
	if len(*calls) != 1 || (*calls)[0].args[0] != "handoff" {
		t.Errorf("want one crow handoff call, got %v", *calls)
	}
}

func TestHandoffTaskRolesAndFindings(t *testing.T) {
	rev := &mrRow{Role: "REVIEWER", Ref: "!7", Title: "retry"}
	own := &mrRow{Role: "AUTHOR", Ref: "!8", Title: "cache", ThreadsOpen: 2}
	if task, _ := handoffTask(rev, ""); !strings.HasPrefix(task, "Review !7") || strings.Contains(task, "findings") {
		t.Errorf("reviewer task: %q", task)
	}
	if task, _ := handoffTask(own, ""); !strings.HasPrefix(task, "Address the 2 open review threads") {
		t.Errorf("own task: %q", task)
	}
	if task, _ := handoffTask(rev, "/s/reviews/g_a-7.json"); !strings.Contains(task, "/s/reviews/g_a-7.json") {
		t.Errorf("findings path missing: %q", task)
	}
}

func TestHandoffOverlayShowsFindingsFile(t *testing.T) {
	m := newModel(defaultConfig())
	m.page = 2
	r := m.mrsView().cur()
	m.fx.MRDetails[r.key()] = &mrDetail{FindingsFile: "/s/reviews/x-1.json"}
	m = send(m, keyMsg("h"))
	if !hasText(m.View(), "/s/reviews/x-1.json") {
		t.Error("overlay must show the findings file before enter")
	}
}
