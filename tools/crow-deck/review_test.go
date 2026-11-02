package main

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// reviewDetail is a review of !5 with two findings and the given run state.
func reviewDetail(state, reason string) *mrDetail {
	raw := rawDetail{}
	raw.Files = append(raw.Files, struct {
		Path string `json:"path"`
		Add  int    `json:"add"`
		Del  int    `json:"del"`
		Diff string `json:"diff"`
	}{"a.py", 2, 1, testDiff})
	d := mapDetail(raw)
	d.Findings = []mrFinding{
		{ID: "f1", Title: "missing pagination", Severity: "blocking", Path: "a.py", Line: 11, Body: "First page only.",
			Comment: "Il manque la pagination.", Fix: "-items = client.list()\n+items = list(paginate(client.list))", State: "pending"},
		{ID: "f2", Title: "commit scope", Severity: "nit", Body: "No file scope.", Comment: "Pas de scope.", State: "edited", Note: "Ajoute le scope PY."},
	}
	d.Review = &mrReview{ID: "rv1", State: state, Reason: reason, Model: "opus", StartedAt: "2026-10-29T13:00:00Z"}
	return d
}

// reviewModel is the live MRs tab on !5 (a review) with the given detail loaded.
func reviewModel(t *testing.T, d *mrDetail) model {
	t.Helper()
	m, _ := liveMRModel(t)
	m.mrSel = 2 // the list is sorted: !6, !12, !5 in "all" mode
	for i, r := range visibleMRs(m.fx, m.mode) {
		if r.Ref == "!5" {
			m.selectMR(i)
		}
	}
	m = send(m, detailMsg{key: "grp/atlas!5", d: d, at: testNow()})
	return m
}

func TestReviewKeyConfirmsThenStartsReadOnlyRun(t *testing.T) {
	calls := recordRuns(t, "review rv1 started", nil)
	m := reviewModel(t, reviewDetail("done", ""))
	m = send(m, keyMsg("r"))
	if m.overlay != ovReview || !hasText(m.View(), "read-only review of !5") {
		t.Fatalf("r must open the review confirm, overlay %d", m.overlay)
	}
	if len(*calls) != 0 {
		t.Fatal("nothing may run before enter")
	}
	m2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m2.(model)
	if m.overlay != ovNone || cmd == nil {
		t.Fatal("enter must close the overlay and return the command")
	}
	var msg tea.Msg = cmd()
	if bm, ok := msg.(tea.BatchMsg); ok {
		for _, c := range bm {
			if c != nil {
				if rm, ok := c().(reviewStartedMsg); ok {
					msg = rm
				}
			}
		}
	}
	want := []string{"start", "grp/atlas", "5", "--mode", "peer"}
	if len(*calls) != 1 || (*calls)[0].name != "crow-review-run" || fmt.Sprint((*calls)[0].args) != fmt.Sprint(want) {
		t.Fatalf("want crow-review-run %v, got %v", want, *calls)
	}
	m = send(m, msg)
	if !m.live.details["grp/atlas!5"].at.IsZero() {
		t.Error("a started review must refetch the detail")
	}
}

func TestReviewKeyRefusesWhileRunning(t *testing.T) {
	m := reviewModel(t, reviewDetail("running", ""))
	m = send(m, keyMsg("r"))
	if m.overlay != ovNone || m.flash == "" {
		t.Errorf("a running review must not start again: overlay %d flash %q", m.overlay, m.flash)
	}
}

func TestDecisionsAreWrittenToTheFindingsJSON(t *testing.T) {
	calls := recordRuns(t, "f1 accepted", nil)
	m := reviewModel(t, reviewDetail("done", ""))
	m2, cmd := m.Update(keyMsg("a"))
	m = m2.(model)
	runAll(cmd)
	if f := m.fx.MRDetails["grp/atlas!5"].Findings[0]; f.State != "accepted" {
		t.Errorf("a: state %q", f.State)
	}
	// e opens the input with the comment, enter saves the edited text.
	m = send(m, keyMsg("e"))
	if m.mrTyping != "e" || m.mrInput != "Il manque la pagination." {
		t.Fatalf("e: typing %q input %q", m.mrTyping, m.mrInput)
	}
	m = send(m, keyMsg(" ici"))
	m2, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m2.(model)
	runAll(cmd)
	// c records a chat note; x refuses.
	m = send(m, keyMsg("c"))
	m = send(m, keyMsg("pourquoi"))
	m2, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m2.(model)
	runAll(cmd)
	m2, cmd = m.Update(keyMsg("x"))
	m = m2.(model)
	runAll(cmd)
	got := []string{}
	for _, c := range *calls {
		if c.name != "crow-review-run" {
			t.Fatalf("unexpected call %v", c)
		}
		got = append(got, fmt.Sprint(c.args))
	}
	want := []string{
		"[decide grp/atlas 5 f1 accepted]",
		"[decide grp/atlas 5 f1 edited --note Il manque la pagination. ici]",
		"[decide grp/atlas 5 f1 chat --note pourquoi]",
		"[decide grp/atlas 5 f1 refused]",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("calls\n got %v\nwant %v", got, want)
	}
}

func TestEscDropsTheEditAndSavesNothing(t *testing.T) {
	calls := recordRuns(t, "", nil)
	m := reviewModel(t, reviewDetail("done", ""))
	m = send(m, keyMsg("e"), tea.KeyMsg{Type: tea.KeyEsc})
	if m.mrTyping != "" || len(*calls) != 0 || m.fx.MRDetails["grp/atlas!5"].Findings[0].State != "pending" {
		t.Errorf("esc: typing %q calls %v", m.mrTyping, *calls)
	}
}

func TestSavedDecisionsShowAfterRestart(t *testing.T) {
	withANSI(t)
	m := reviewModel(t, reviewDetail("done", ""))
	m.mrItem = 1
	out := m.View()
	for _, want := range []string{"Findings 2", "edited", "your edit", "Ajoute le scope PY.", "1 accepted"} {
		if !hasText(out, want) {
			t.Errorf("view lacks %q", want)
		}
	}
}

func TestReviewRunStateShows(t *testing.T) {
	withANSI(t)
	d := reviewDetail("running", "")
	d.Findings = nil
	m := reviewModel(t, d)
	if out := m.View(); !hasText(out, "review running") || !hasText(out, "opus") {
		t.Error("a running review must show as running")
	}
	d = reviewDetail("failed", "no findings JSON in the run's final message")
	d.Findings = nil
	m = reviewModel(t, d)
	if out := m.View(); !hasText(out, "review failed: no findings JSON") || !hasText(out, "runs the review again") {
		t.Error("a failed review must show its reason")
	}
}

func TestRunningReviewRefetchesSooner(t *testing.T) {
	m := reviewModel(t, reviewDetail("running", ""))
	m.now = testNow().Add(reviewTTL + 1)
	if m.wantDetail() == nil {
		t.Error("a running review must be refetched after reviewTTL")
	}
	m = reviewModel(t, reviewDetail("done", ""))
	m.now = testNow().Add(reviewTTL + 1)
	if m.wantDetail() != nil {
		t.Error("a finished review keeps the normal TTL")
	}
}

func TestReviewCardFitsEveryThemeAndSize(t *testing.T) {
	withANSI(t)
	for ti, th := range themes {
		for _, sz := range [][2]int{{80, 24}, {170, 46}, {118, 60}} {
			for _, typing := range []string{"", "e"} {
				m := reviewModel(t, reviewDetail("failed", "run exited 1: usage limit"))
				m.th, m.w, m.h, m.mrPane = ti, sz[0], sz[1], mpCard
				if typing != "" {
					m = send(m, keyMsg("e"))
				}
				out := m.View()
				name := fmt.Sprintf("%s %dx%d typing=%q", th.name, sz[0], sz[1], typing)
				checkFit(t, name, out, sz[0], sz[1])
				if !hasText(out, "missing pagination") {
					t.Errorf("%s: finding not shown", name)
				}
			}
		}
	}
}

func runAll(cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	if bm, ok := cmd().(tea.BatchMsg); ok {
		for _, c := range bm {
			runAll(c)
		}
	}
}
