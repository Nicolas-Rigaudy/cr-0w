package main

import (
	"fmt"
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
	for _, want := range []string{"Threads 3", "Thread 1", "sam: Pagine cet appel.", "agree", "no reply: the fix answers it",
		"proposed fix", "test", "assert len(items) == 3", "push, reply, resolve, after your OK"} {
		if !hasText(out, want) {
			t.Errorf("thread 1 lacks %q", want)
		}
	}
	m.mrItem = 1
	if out := m.View(); !hasText(out, "reply (FR)") || !hasText(out, "Le nom suit le module voisin.") {
		t.Error("thread 2 must show its reply")
	}
	m.mrItem = 2
	if out := m.View(); !hasText(out, "disagree") || !hasText(out, "stays open") {
		t.Error("thread 3 must say the run disagrees and the thread stays open")
	}
	if out := send(m, keyMsg("T")).View(); !hasText(out, "Pagine cet appel") || hasText(out, "Threads 3") {
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
	out := threadsModel(t, d).View()
	if !hasText(out, "resolved") || !hasText(out, "replied") {
		t.Error("the list must show resolved and replied threads")
	}
}
