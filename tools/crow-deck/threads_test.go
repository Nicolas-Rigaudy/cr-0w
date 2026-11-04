package main

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// threadsDetail is a threads review of !6 (the owner's MR): f1 fixed with a test and no reply,
// f2 a reply only, f3 a thread the run disagrees with.
func threadsDetail() *mrDetail {
	d := reviewDetail("done", "")
	d.Review.Mode = "threads"
	d.Findings = []mrFinding{
		{ID: "f1", Title: "Liste non paginée", Severity: "blocking", Path: "a.py", Line: 11, Body: "list() rend une page.", Agree: "yes",
			ThreadAuthor: "sam", ThreadBody: "Pagine cet appel.", Fix: "-items = client.list()\n+items = list(paginate(client.list))",
			Test: "+assert len(items) == 3", State: "pending"},
		{ID: "f2", Title: "Nom", Severity: "question", Body: "Question sur le nom.", Agree: "yes", ThreadAuthor: "sam",
			ThreadBody: "Pourquoi ce nom ?", Comment: "Le nom suit le module voisin.", State: "pending"},
		{ID: "f3", Title: "Renommage", Severity: "nit", Path: "a.py", Line: 11, Body: "Le nom suit la convention.", Agree: "no",
			ThreadAuthor: "sam", ThreadBody: "Renomme.", Comment: "Je garde ce nom.", State: "pending"},
	}
	d.Threads = []mrThread{{Author: "sam", Path: "a.py", Line: 11, Body: "Pagine cet appel."}}
	return d
}

// threadsModel is the live MRs tab on !6 with a threads review loaded; no T needed.
func threadsModel(t *testing.T, d *mrDetail) model {
	t.Helper()
	m, _ := liveMRModel(t)
	for i, r := range visibleMRs(m.fx, m.mode) {
		if r.Ref == "!6" {
			m.selectMR(i)
		}
	}
	return send(m, detailMsg{key: "grp/atlas!6", d: d, at: testNow()})
}

func TestReviewModeOnOwnMRWithThreads(t *testing.T) {
	cases := map[string]mrRow{
		"threads": {Role: "AUTHOR", ThreadsOpen: 2},
		"self":    {Role: "AUTHOR"},
		"peer":    {Role: "REVIEWER", ThreadsOpen: 2},
	}
	for want, r := range cases {
		if got := reviewMode(&r); got != want {
			t.Errorf("%+v: %s, want %s", r, got, want)
		}
	}
}

func TestThreadsReviewShowsEachAnswer(t *testing.T) {
	withANSI(t)
	m := threadsModel(t, threadsDetail())
	out := m.View()
	for _, want := range []string{"finding 1 of 3", "sam", "Pagine cet appel.", "no reply: the fix answers it",
		"Proposed fix", "Test", "assert len(items) == 3", "push + reply"} {
		if !hasText(out, want) {
			t.Errorf("thread 1 lacks %q", want)
		}
	}
	m.mrItem = 1
	if out := m.View(); !hasText(out, "Reply") || !hasText(out, "Le nom suit le module voisin.") {
		t.Error("thread 2 must show its reply")
	}
	m.mrItem = 2
	if out := m.View(); !hasText(out, "disagree") || !hasText(out, "stays open") {
		t.Error("thread 3 must say the run disagrees and the thread stays open")
	}
	if out := send(m, keyMsg("T")).View(); !hasText(out, "Pagine cet appel") || !hasText(out, "thread 1 of 1") {
		t.Error("T shows the raw threads")
	}
}

func TestThreadsDecisionsAndApply(t *testing.T) {
	calls := recordRuns(t, "", nil)
	m := threadsModel(t, threadsDetail())
	step := func(k tea.Msg) {
		next, cmd := m.Update(k)
		m = next.(model)
		runAll(cmd)
	}
	step(keyMsg("a"))
	m.mrItem = 1
	step(keyMsg("e"))
	if m.mrInput != "Le nom suit le module voisin." {
		t.Fatalf("e must prefill the reply, got %q", m.mrInput)
	}
	step(tea.KeyMsg{Type: tea.KeyEnter})
	if len(*calls) < 2 || fmt.Sprint((*calls)[0].args) != "[decide grp/atlas 6 f1 accepted]" {
		t.Fatalf("decisions must be saved, got %v", *calls)
	}
	// Only f1 has a fix: the edited reply of f2 needs no apply run.
	d := m.fx.MRDetails["grp/atlas!6"]
	n := 0
	for _, f := range d.Findings {
		if f.toApply("threads") {
			n++
		}
	}
	if n != 1 {
		t.Errorf("threads apply acts on fixes only, got %d", n)
	}
	m = send(m, keyMsg("A"))
	if m.overlay != ovApply || !hasText(m.View(), "fix 1 finding") {
		t.Errorf("A must offer the one fix, overlay %d flash %q", m.overlay, m.flash)
	}
}

func TestThreadsReplyAndResolveOverlays(t *testing.T) {
	fakeBin(t, []postAction{
		{ID: "reply-f2", Kind: "reply", Title: "Reply to sam on the MR", Text: "Le nom suit le module voisin.", Hash: "h1",
			Cmd: []string{"glab", "api", "-X", "POST", "projects/grp%2Fatlas/merge_requests/6/discussions/d2/notes", "-f", "body=Le nom suit le module voisin."}},
		{ID: "resolve-f2", Kind: "resolve", Title: "Resolve the thread of sam", Text: "Pourquoi ce nom ?", Hash: "h2",
			Cmd: []string{"glab", "api", "-X", "PUT", "projects/grp%2Fatlas/merge_requests/6/discussions/d2", "-f", "resolved=true"}},
	}, nil)
	m := press(threadsModel(t, threadsDetail()), keyMsg("P"))
	out := m.View()
	for _, want := range []string{"Reply to sam", "Le nom suit le module voisin.", "body=<the text above>", "enter post this reply"} {
		if !hasText(out, want) {
			t.Errorf("reply overlay lacks %q", want)
		}
	}
	m = send(m, keyMsg("s"))
	if out := m.View(); !hasText(out, "Resolve the thread of sam") || !hasText(out, "resolved=true") || !hasText(out, "enter resolve this thread") {
		t.Error("the resolve overlay must show its command")
	}
}

func TestThreadsStatesInTheList(t *testing.T) {
	withANSI(t)
	d := threadsDetail()
	d.Findings[0].State, d.Findings[0].Resolved = "accepted", true
	d.Findings[1].State, d.Findings[1].Posted = "accepted", true
	m := threadsModel(t, d)
	if !hasText(m.View(), "resolved") {
		t.Error("the card must show a resolved thread")
	}
	m.mrItem = 1
	if !hasText(m.View(), "replied") {
		t.Error("the card must show a replied thread")
	}
}

func TestThreadsCardKeepsButtonsAndStatusInSight(t *testing.T) {
	withANSI(t)
	d := threadsDetail()
	d.Findings[0].Body = strings.Repeat("A long assessment of the thread. ", 40)
	d.Apply = &mrApply{ID: "ap1", State: "running", StartedAt: "2026-10-29T13:59:00Z"}
	m := threadsModel(t, d)
	m.w, m.h = 170, 30
	out := m.View()
	for _, want := range []string{"accept the fix", "chat", "apply running", "more lines"} {
		if !hasText(out, want) {
			t.Errorf("a short window must still show %q", want)
		}
	}
}

func TestApplyStartShowsAtOnceAndResultsInTheList(t *testing.T) {
	withANSI(t)
	recordRuns(t, "apply ap1 started", nil)
	d := threadsDetail()
	d.Findings[0].State = "accepted"
	m := threadsModel(t, d)
	m = send(m, keyMsg("A"), tea.KeyMsg{Type: tea.KeyEnter})
	if !hasText(m.View(), "apply starting") {
		t.Error("enter on the apply confirm must show the run at once")
	}
	m = send(m, reviewStartedMsg{key: "grp/atlas!6", text: "apply failed to start: no slot", apply: true, failed: true})
	m.mrPane = mpCard
	m = send(m, keyMsg("j"))
	if !hasText(m.View(), "apply failed: apply failed to start") {
		t.Error("a failed start must stay on the card after the next key")
	}
	d.Findings[0].Applied, d.Findings[0].AppliedFor = "fixed", [2]string{"accepted", ""}
	if !hasText(threadsModel(t, d).View(), "fixed") {
		t.Error("the list must show the apply result")
	}
}

func TestThreadsDiffOpensAtTheReviewerLine(t *testing.T) {
	d := threadsDetail()
	d.Files = append(d.Files, mrFile{Path: "b.py", Lines: parseDiff("@@ -1,1 +1,2 @@\n x\n+y\n")})
	d.Files[0], d.Files[1] = d.Files[1], d.Files[0] // b.py first: the diff must still open a.py
	d.Findings[0].ThreadPath, d.Findings[0].ThreadLine = "a.py", 11
	m := threadsModel(t, d)
	if mv := m.mrsView(); mv.detail().Files[mv.fileIdx(mv.detail())].Path != "a.py" {
		t.Errorf("the diff opened %q", mv.detail().Files[mv.fileIdx(mv.detail())].Path)
	}
}

func TestThreadsShowPushedAtTheTop(t *testing.T) {
	withANSI(t)
	d := threadsDetail()
	d.Apply = &mrApply{State: "done", Commit: "abc1234def", Pushed: true}
	d.Findings[0].Applied, d.Findings[0].AppliedFor = "fixed", [2]string{"accepted", ""}
	d.Findings[0].State = "accepted"
	d.Findings[1].State, d.Findings[1].Resolved = "accepted", true
	d.Findings[2].State = "refused"
	out := threadsModel(t, d).View()
	for _, want := range []string{"committed abc1234d", "pushed to the MR branch", "1 thread to resolve", "resolves 1 thread", "pushed abc1234d"} {
		if !hasText(out, want) {
			t.Errorf("missing %q", want)
		}
	}
	d.Apply.Pushed = false
	out = threadsModel(t, d).View()
	for _, want := range []string{"committed abc1234d", "local only, P pushes it", "push abc1234d", "fixed"} {
		if !hasText(out, want) {
			t.Errorf("unpushed: missing %q", want)
		}
	}
}
