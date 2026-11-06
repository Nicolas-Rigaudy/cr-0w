package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// liveHist is a History page whose first row is a real, finished session.
func liveHist(t *testing.T) model {
	old := copyText
	copyText = func(string) {}
	t.Cleanup(func() { copyText = old })
	m := histModel(170, 46)
	rows := m.histView().rows()
	for i := range m.fx.Hist {
		if m.fx.Hist[i].ID == rows[0].ID {
			s := &m.fx.Hist[i]
			s.ID, s.Account, s.Cwd, s.Path, s.Active = "abc-123", "main", "/w/atlas", "/p/abc-123.jsonl", false
		}
	}
	return m
}

func TestHistoryActionsAskThenRun(t *testing.T) {
	cases := map[string][]string{
		"enter": {"crow-focus", "--pid", "0", "--account", "main", "--cwd", "/w/atlas", "--session", "abc-123"},
		"r":     {"crow-retro", "add", "/p/abc-123.jsonl"},
		"h":     {"crow", "handoff"},
	}
	for k, want := range cases {
		calls := recordRuns(t, "ok", nil)
		m := send(liveHist(t), keyMsg(k))
		if m.overlay != ovAct || len(*calls) != 0 {
			t.Fatalf("%s: overlay %d calls %d flash %q", k, m.overlay, len(*calls), m.flash)
		}
		_, cmd := m.Update(keyMsg("y"))
		runAll(cmd)
		if len(*calls) != 1 {
			t.Fatalf("%s: %d runs", k, len(*calls))
		}
		got := append([]string{(*calls)[0].name}, (*calls)[0].args...)
		if strings.Join(got[:len(want)], " ") != strings.Join(want, " ") {
			t.Errorf("%s ran %v", k, got)
		}
		if k == "h" && !strings.Contains(strings.Join(got, " "), "Pick up session abc-123") {
			t.Errorf("handoff note lacks the session: %v", got)
		}
	}
	m := send(liveHist(t), keyMsg("y"))
	if m.flash != "copied abc-123" || m.overlay != ovNone {
		t.Errorf("y: flash %q", m.flash)
	}
	if _, cmd := liveHist(t).Update(keyMsg("o")); cmd == nil {
		t.Error("o: no transcript command")
	}
}

func TestHistoryRunningSessionIsNotResumed(t *testing.T) {
	m := liveHist(t)
	m.fx.Hist[0].Active, m.fx.Hist[0].ID = true, "x-1"
	m.histSel = 0
	for i, s := range m.histView().rows() {
		if s.ID == "x-1" {
			m.histSel = i
		}
	}
	m = send(m, keyMsg("enter"))
	if m.overlay != ovNone || !strings.Contains(m.flash, "still running") {
		t.Errorf("overlay %d flash %q", m.overlay, m.flash)
	}
}

func TestProjectsCardActionsAskFirst(t *testing.T) {
	for k, want := range map[string]string{"c": "crow term --on-busy fail", "s": "crow term --on-busy fail", "o": "crow term --on-busy fail", "h": "crow handoff"} {
		calls := recordRuns(t, "ok", nil)
		m := send(projModel(235, 62), keyMsg(k))
		if m.overlay != ovAct || m.act == nil || len(*calls) != 0 {
			t.Fatalf("%s: overlay %d flash %q", k, m.overlay, m.flash)
		}
		if !strings.HasPrefix(strings.Join(m.act.Cmd, " "), want) {
			t.Errorf("%s: cmd %v", k, m.act.Cmd)
		}
		if k != "o" && !strings.Contains(strings.Join(m.act.Cmd, " "), "--task") {
			t.Errorf("%s: no literal task: %v", k, m.act.Cmd)
		}
		if k == "o" && strings.Contains(strings.Join(m.act.Cmd, " "), "--task") {
			t.Errorf("o opens a bare session: %v", m.act.Cmd)
		}
	}
}

func TestProjectsReviewJumpsToTheMR(t *testing.T) {
	base := projModel(235, 62)
	jumped, missed := false, false
	for i := range base.projectsView().visible() {
		m := base
		m.prjSel = i
		c := m.projectsView().cur()
		hasMR := false
		for _, r := range visibleMRs(m.fx, m.mode) {
			hasMR = hasMR || m.sameProject(r.Project, c.Name) || m.sameProject(r.Project, c.Short)
		}
		m = send(m, keyMsg("r"))
		switch {
		case hasMR && m.page == 2 && m.overlay == ovReview:
			jumped = true
		case !hasMR && m.page == 3 && strings.Contains(m.flash, "no open MR"):
			missed = true
		default:
			t.Errorf("%s: page %d overlay %d flash %q", c.Name, m.page, m.overlay, m.flash)
		}
	}
	if !jumped || !missed {
		t.Errorf("fixtures must cover both cases: jumped %v missed %v", jumped, missed)
	}
}

func TestHomeHandoffAndReview(t *testing.T) {
	m := newModel(defaultConfig())
	m = send(m, resize(172, 48))
	m.focus = pNext
	rows := m.homeView().nextItems()
	if len(rows) == 0 {
		t.Fatal("no Next rows in the fixtures")
	}
	h := send(m, keyMsg("h"))
	if h.overlay != ovAct || h.act == nil || h.act.Cmd[1] != "handoff" || !strings.Contains(strings.Join(h.act.Cmd, " "), rows[0].What) {
		t.Fatalf("h: overlay %d act %+v", h.overlay, h.act)
	}
	// r on an MR line opens that MR's review overlay.
	for i, r := range rows {
		if mrRef(r.What) == "" {
			continue
		}
		m.sel[pNext] = i
		got := send(m, keyMsg("r"))
		if got.page != 2 || got.overlay != ovReview || got.mrsView().cur().Ref != mrRef(r.What) {
			t.Errorf("r on %q: page %d overlay %d flash %q", r.What, got.page, got.overlay, got.flash)
		}
		break
	}
	m.focus = pPulse
	if p := send(m, keyMsg("h")); p.overlay != ovNone || !strings.Contains(p.flash, "pick a row") {
		t.Errorf("pulse h: %q", p.flash)
	}
}

func TestTranscriptTailKeepsTextTurns(t *testing.T) {
	p := filepath.Join(t.TempDir(), "s.jsonl")
	lines := []string{
		`{"type":"user","timestamp":"2026-11-06T08:00:00Z","message":{"content":"first ask"}}`,
		`{"type":"assistant","timestamp":"2026-11-06T08:01:00Z","message":{"content":[{"type":"tool_use","name":"Read"}]}}`,
		`{"type":"user","timestamp":"2026-11-06T08:02:00Z","message":{"content":[{"type":"tool_result","content":"x"}]}}`,
		`{"type":"user","isMeta":true,"message":{"content":"meta"}}`,
		`{"type":"assistant","timestamp":"2026-11-06T08:03:00Z","message":{"content":[{"type":"text","text":"done:\n tests   green"}]}}`,
		`not json`,
	}
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got := transcriptTail(p, 3)
	if len(got) != 2 || got[0].Text != "you first ask" || got[1].Text != "› done: tests green" {
		t.Fatalf("tail %+v", got)
	}
	if transcriptTail(filepath.Join(t.TempDir(), "none"), 3) != nil {
		t.Error("a missing file must give nothing")
	}
}
