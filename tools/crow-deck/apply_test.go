package main

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// selfReviewModel is the live MRs tab on !6 (the owner's MR, findings view) with the given detail.
func selfReviewModel(t *testing.T, d *mrDetail) model {
	t.Helper()
	m, _ := liveMRModel(t)
	for i, r := range visibleMRs(m.fx, m.mode) {
		if r.Ref == "!6" {
			m.selectMR(i)
		}
	}
	m.mrToggle = true // own MR: T shows the findings
	return send(m, detailMsg{key: "grp/atlas!6", d: d, at: testNow()})
}

func applyDetail() *mrDetail {
	d := reviewDetail("done", "")
	d.Findings[0].State = "accepted"
	d.Findings = append(d.Findings, mrFinding{ID: "f3", Title: "naming", Severity: "question", Path: "a.py", Line: 12,
		Body: "Why this name?", State: "chat", Note: "Pourquoi ce nom ?"})
	return d
}

func TestApplyKeyConfirmsThenResumesTheReview(t *testing.T) {
	calls := recordRuns(t, "apply ap1 started", nil)
	m := selfReviewModel(t, applyDetail())
	m = send(m, keyMsg("A"))
	if m.overlay != ovApply || !hasText(m.View(), "fix 2 findings") || !hasText(m.View(), "answer 1") || !hasText(m.View(), "nothing pushed") {
		t.Fatalf("A must open the apply confirm with the counts, overlay %d flash %q", m.overlay, m.flash)
	}
	if len(*calls) != 0 {
		t.Fatal("nothing may run before enter")
	}
	m2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m2.(model)
	runAll(cmd)
	if m.overlay != ovNone || len(*calls) != 1 || fmt.Sprint((*calls)[0].args) != "[apply grp/atlas 6]" {
		t.Fatalf("enter must run crow-review-run apply grp/atlas 6, got %v", *calls)
	}
}

func TestApplyKeyEscRunsNothing(t *testing.T) {
	calls := recordRuns(t, "", nil)
	m := selfReviewModel(t, applyDetail())
	m = send(m, keyMsg("A"), tea.KeyMsg{Type: tea.KeyEsc})
	if m.overlay != ovNone || len(*calls) != 0 {
		t.Errorf("esc: overlay %d calls %v", m.overlay, *calls)
	}
}

func TestApplyKeyRefusals(t *testing.T) {
	running := applyDetail()
	running.Apply = &mrApply{ID: "ap1", State: "running"}
	handled := applyDetail()
	for i := range handled.Findings {
		f := &handled.Findings[i]
		f.Applied, f.AppliedFor = "fixed", [2]string{f.State, f.Note}
	}
	cases := []struct {
		name string
		m    model
		want string
	}{
		{"review running", selfReviewModel(t, reviewDetail("running", "")), "finished review"},
		{"apply running", selfReviewModel(t, running), "already running"},
		{"all handled", selfReviewModel(t, handled), "nothing to apply"},
		// A peer review applies chat notes only: f1 pending and f2 edited give nothing.
		{"peer without chat", reviewModel(t, reviewDetail("done", "")), "only chat notes"},
	}
	for _, c := range cases {
		m := send(c.m, keyMsg("A"))
		if m.overlay != ovNone || !hasText(m.flash, c.want) {
			t.Errorf("%s: overlay %d flash %q", c.name, m.overlay, m.flash)
		}
	}
}

func TestApplyResultShows(t *testing.T) {
	withANSI(t)
	d := applyDetail()
	d.Apply = &mrApply{ID: "ap1", State: "done", Commit: "0123456789abcdef", Subject: "fix[PY]: paginate the list call",
		Warnings: []string{"files outside the accepted findings: b.py"}}
	d.Findings[2].Applied, d.Findings[2].AppliedDetail = "answered", "Le nom suit le module voisin."
	d.Findings[2].AppliedFor = [2]string{"chat", "Pourquoi ce nom ?"}
	m := selfReviewModel(t, d)
	m.mrItem = 2
	out := m.View()
	for _, want := range []string{"committed", "01234567 fix[PY]: paginate", "not pushed", "! files outside", "answer", "Le nom suit le module voisin."} {
		if !hasText(out, want) {
			t.Errorf("view lacks %q", want)
		}
	}
	if hasText(out, "before your last decision") {
		t.Error("an answer for the current note is not stale")
	}
	// A new decision on the finding marks the answer as older than it.
	m.fx.MRDetails["grp/atlas!6"].Findings[2].Note = "Et pour les tests ?"
	if !hasText(m.View(), "before your last decision") {
		t.Error("an answer to an older note must say so")
	}
}

func TestApplyFailureSaysWhy(t *testing.T) {
	withANSI(t)
	d := applyDetail()
	d.Apply = &mrApply{ID: "ap1", State: "failed", Reason: "the resumed session lost the review context; nothing was applied. r runs the review again"}
	if out := selfReviewModel(t, d).View(); !hasText(out, "apply failed: the resumed session lost") {
		t.Error("a lost context must show on the card")
	}
}

func TestRunningApplyRefetchesSooner(t *testing.T) {
	d := applyDetail()
	d.Apply = &mrApply{ID: "ap1", State: "running"}
	m := selfReviewModel(t, d)
	m.now = testNow().Add(reviewTTL + 1)
	if m.wantDetail() == nil {
		t.Error("a running apply must be refetched after reviewTTL")
	}
}

func TestReviewStartPassesTheBranch(t *testing.T) {
	calls := recordRuns(t, "review rv1 started", nil)
	r := mrRow{Repo: "grp/atlas", Number: 6, Ref: "!6", Role: "AUTHOR", Branch: "feat/seed"}
	reviewCmd(r)()
	if fmt.Sprint((*calls)[0].args) != "[start grp/atlas 6 --mode self --branch feat/seed]" {
		t.Errorf("got %v", *calls)
	}
}

func TestApplyCardFitsEveryThemeAndSize(t *testing.T) {
	withANSI(t)
	for ti, th := range themes {
		for _, sz := range [][2]int{{80, 24}, {170, 46}, {118, 60}} {
			d := applyDetail()
			d.Apply = &mrApply{ID: "ap1", State: "done", Commit: "0123456789abcdef", Subject: "fix[PY]: paginate", Warnings: []string{"a long warning about files outside the accepted findings: b.py, c.py, d.py"}}
			d.Findings[2].Applied, d.Findings[2].AppliedDetail = "answered", "Une réponse assez longue pour passer sur plusieurs lignes dans la carte."
			for _, ov := range []int{ovNone, ovApply} {
				m := selfReviewModel(t, d)
				m.th, m.w, m.h, m.mrPane, m.mrItem, m.overlay = ti, sz[0], sz[1], mpCard, 2, ov
				checkFit(t, fmt.Sprintf("%s %dx%d overlay=%d", th.name, sz[0], sz[1], ov), m.View(), sz[0], sz[1])
			}
		}
	}
}
