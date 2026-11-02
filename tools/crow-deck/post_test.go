package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// fakeBin answers crow-review-run like the real one: actions lists acts, act answers per id.
// It records every call; nothing real runs.
func fakeBin(t *testing.T, acts []postAction, fail map[string]string) *[]call {
	t.Helper()
	var calls []call
	old := runBin
	runBin = func(name string, args ...string) (string, error) {
		calls = append(calls, call{name, args})
		switch args[0] {
		case "actions":
			b, _ := json.Marshal(acts)
			return string(b), nil
		case "act":
			if msg, ok := fail[args[3]]; ok {
				return "crow-review-run: " + msg, errors.New("exit status 1")
			}
			return "posted " + args[3], nil
		}
		return "", nil
	}
	t.Cleanup(func() { runBin = old })
	return &calls
}

func comments() []postAction {
	return []postAction{
		{ID: "post-f1", Kind: "comment", Title: "Comment f1 on a.py:11", Text: "Il manque la pagination.", Hash: "aaa",
			Cmd: []string{"glab", "api", "-X", "POST", "projects/grp%2Fatlas/merge_requests/5/discussions", "-f", "body=Il manque la pagination.", "-F", "position[new_line]=11"}},
		{ID: "post-f2", Kind: "comment", Title: "Comment f2 on grp/atlas !5", Text: "Ajoute le scope PY.", Hash: "bbb",
			Cmd: []string{"glab", "api", "-X", "POST", "projects/grp%2Fatlas/merge_requests/5/notes", "-f", "body=Ajoute le scope PY."}},
	}
}

// press sends a key and feeds back every message its command returns.
func press(m model, k tea.Msg) model {
	m2, cmd := m.Update(k)
	m = m2.(model)
	for _, msg := range collect(cmd) {
		m = send(m, msg)
	}
	return m
}

func collect(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if bm, ok := msg.(tea.BatchMsg); ok {
		var out []tea.Msg
		for _, c := range bm {
			out = append(out, collect(c)...)
		}
		return out
	}
	return []tea.Msg{msg}
}

func acted(calls []call) []string {
	var out []string
	for _, c := range calls {
		if c.args[0] == "act" {
			out = append(out, fmt.Sprint(c.args))
		}
	}
	return out
}

func TestPostShowsEachActionAndRunsOnlyThatOne(t *testing.T) {
	calls := fakeBin(t, comments(), nil)
	m := press(reviewModel(t, reviewDetail("done", "")), keyMsg("P"))
	out := m.View()
	if m.overlay != ovPost || !hasText(out, "1/2 Comment f1 on a.py:11") || !hasText(out, "Il manque la pagination.") ||
		!hasText(out, "glab api -X POST") || !hasText(out, "body=<the text above>") {
		t.Fatalf("P must show the first action with its text and command, overlay %d flash %q", m.overlay, m.flash)
	}
	if len(acted(*calls)) != 0 {
		t.Fatal("nothing may be posted before enter")
	}
	m = press(m, tea.KeyMsg{Type: tea.KeyEnter})
	if got := acted(*calls); fmt.Sprint(got) != "[[act grp/atlas 5 post-f1 --hash aaa]]" {
		t.Fatalf("enter must run that one action, got %v", got)
	}
	if m.overlay != ovPost || !hasText(m.View(), "2/2") || !hasText(m.View(), "Ajoute le scope PY.") || m.flash != "posted post-f1" {
		t.Fatalf("after one OK the next action waits for its own: overlay %d flash %q", m.overlay, m.flash)
	}
	if !m.live.details["grp/atlas!5"].at.IsZero() {
		t.Error("a posted action must refetch the detail")
	}
	m = press(m, tea.KeyMsg{Type: tea.KeyEsc})
	if m.overlay != ovNone || len(acted(*calls)) != 1 {
		t.Errorf("esc stops: overlay %d acts %v", m.overlay, acted(*calls))
	}
}

func TestPostSkipAndWaitWhileRunning(t *testing.T) {
	calls := fakeBin(t, comments(), nil)
	m := press(reviewModel(t, reviewDetail("done", "")), keyMsg("P"))
	m = press(m, keyMsg("s"))
	if !hasText(m.View(), "2/2") || len(acted(*calls)) != 0 {
		t.Fatal("s skips without running")
	}
	// enter starts the action; keys wait until its answer arrives.
	m2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m2.(model)
	m = press(m, tea.KeyMsg{Type: tea.KeyEnter})
	m = press(m, tea.KeyMsg{Type: tea.KeyEsc})
	if !m.posting || m.overlay != ovPost || !hasText(m.View(), "running") {
		t.Fatal("while an action runs the overlay stays and ignores keys")
	}
	for _, msg := range collect(cmd) {
		m = send(m, msg)
	}
	if got := acted(*calls); fmt.Sprint(got) != "[[act grp/atlas 5 post-f2 --hash bbb]]" || m.overlay != ovNone {
		t.Errorf("only f2 runs, then the overlay closes: %v overlay %d", got, m.overlay)
	}
}

func TestPostFailureStopsAndSaysWhy(t *testing.T) {
	calls := fakeBin(t, comments(), map[string]string{"post-f1": "post-f1 changed since it was shown; P shows it again"})
	m := press(reviewModel(t, reviewDetail("done", "")), keyMsg("P"))
	m = press(m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.overlay != ovNone || !strings.Contains(m.flash, "changed since it was shown") || len(acted(*calls)) != 1 {
		t.Errorf("a failed action closes the overlay and names the reason: overlay %d flash %q", m.overlay, m.flash)
	}
}

func TestPostNothingPendingAndMock(t *testing.T) {
	fakeBin(t, nil, nil)
	m := press(reviewModel(t, reviewDetail("done", "")), keyMsg("P"))
	if m.overlay != ovNone || !strings.Contains(m.flash, "nothing to push or post on !5") {
		t.Errorf("empty list: overlay %d flash %q", m.overlay, m.flash)
	}
	calls := fakeBin(t, comments(), nil)
	mock := newModel(defaultConfig())
	mock.page = 2
	mock = press(mock, keyMsg("P"))
	if mock.overlay != ovNone || len(*calls) != 0 || !strings.Contains(mock.flash, "mock deck") {
		t.Errorf("the mock deck runs nothing: %v %q", *calls, mock.flash)
	}
}

func TestPushShowsTheCommitAndWarnings(t *testing.T) {
	fakeBin(t, []postAction{{ID: "push", Kind: "push", Title: "Push 1a2b3c4d to origin/feat/ids", Text: "fix[PY]: paginate",
		Warnings: []string{"files outside the accepted findings: b.py"}, Hash: "ccc",
		Cmd: []string{"git", "-C", "/w/atlas.wt/review-6", "push", "origin", "1a2b3c4d:refs/heads/feat/ids"}}}, nil)
	m := press(selfReviewModel(t, applyDetail()), keyMsg("P"))
	out := m.View()
	for _, want := range []string{"Push 1a2b3c4d to origin/feat/ids", "fix[PY]: paginate", "git -C /w/atlas.wt/review-6 push origin", "! files outside", "enter push"} {
		if !hasText(out, want) {
			t.Errorf("push overlay lacks %q", want)
		}
	}
}

func TestPostedAndPushedShowOnTheCard(t *testing.T) {
	withANSI(t)
	d := reviewDetail("done", "")
	d.Findings[0].State, d.Findings[0].Posted = "accepted", true
	if out := reviewModel(t, d).View(); !hasText(out, "posted") || !hasText(out, "1 posted") || !hasText(out, "post 1 comment,") {
		t.Error("a posted finding must say so")
	}
	d = applyDetail()
	for i := range d.Findings {
		d.Findings[i].AppliedFor = [2]string{d.Findings[i].State, d.Findings[i].Note}
	}
	d.Apply = &mrApply{ID: "ap1", State: "done", Commit: "1a2b3c4d5e", Subject: "fix[PY]: paginate"}
	if out := selfReviewModel(t, d).View(); !hasText(out, "nothing pushed yet") || !hasText(out, "push 1a2b3c4d") {
		t.Errorf("a committed apply must offer the push\n%s", out)
	}
	d.Apply.Pushed = true
	if out := selfReviewModel(t, d).View(); hasText(out, "nothing pushed yet") || !hasText(out, "pushed; ask a reviewer") {
		t.Error("a pushed commit must say so")
	}
}

func TestPostOverlayFitsEveryThemeAndSize(t *testing.T) {
	withANSI(t)
	long := comments()
	long[0].Text = strings.Repeat("Une ligne de commentaire assez longue pour se replier. ", 90)
	for ti, th := range themes {
		for _, sz := range [][2]int{{80, 24}, {170, 46}, {118, 60}, {60, 16}} {
			fakeBin(t, long, nil)
			m := reviewModel(t, reviewDetail("done", ""))
			m.th, m.w, m.h = ti, sz[0], sz[1]
			m = press(m, keyMsg("P"))
			out := m.View()
			name := fmt.Sprintf("%s %dx%d", th.name, sz[0], sz[1])
			checkFit(t, name, out, sz[0], sz[1])
			if !hasText(out, "more lines") || !hasText(out, "esc") {
				t.Errorf("%s: a long text must say it is cut and keep the keys", name)
			}
		}
	}
}
