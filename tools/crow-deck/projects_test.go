package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func projModel(w, h int) model {
	m := newModel(defaultConfig())
	m.page, m.w, m.h, m.mode = 3, w, h, "all"
	return m
}

func TestProjectsShowsCardsAndTimeline(t *testing.T) {
	out := projModel(170, 46).View()
	for _, want := range []string{"NOW", "NEXT", "BLOCK", "ROAD", "MILE", "OPEN", "6/8 · target Dec 1 · on time", "Timeline", "▼ today", "◆", "no roadmap"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestProjectsActionsOnlyOnSelectedCard(t *testing.T) {
	out := projModel(170, 46).View()
	if n := strings.Count(out, "continue"); n != 1 {
		t.Errorf("%d action rows, want 1", n)
	}
}

func TestProjectsKeys(t *testing.T) {
	m := projModel(170, 46)
	press := func(k string) {
		nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)})
		m = nm.(model)
	}
	press("j")
	if m.prjSel != 1 {
		t.Errorf("j: card %d", m.prjSel)
	}
	for i := 0; i < 9; i++ {
		press("j")
	}
	if m.prjSel != len(m.fx.Cards)-1 {
		t.Errorf("selection ran past the cards: %d", m.prjSel)
	}
	press("c")
	if !strings.Contains(m.flash, "continue on crow") {
		t.Errorf("flash %q", m.flash)
	}
	press("z")
	if !m.prjZoom {
		t.Error("z did not zoom")
	}
	if !strings.Contains(m.View(), "3 months · z zoom") {
		t.Error("zoomed timeline lacks its label")
	}
}

func TestProjectsModeFilter(t *testing.T) {
	m := projModel(170, 46)
	m.mode = "work"
	out := m.View()
	if !strings.Contains(out, "atlas") || strings.Contains(out, "lint · release") {
		t.Errorf("work mode shows personal cards:\n%s", out)
	}
}

func TestProjectsLiveWithoutDates(t *testing.T) {
	reg := &registry{Projects: []regProject{
		{Name: "alpha", Short: "alpha", Scope: "work", Road: &regRoad{Done: 2, Total: 5, Current: "0.3 - Next", Exit: "tests pass"}},
		{Name: "beta", Short: "beta", Scope: "personal"},
	}}
	now := testNow()
	cards, quiet := projectCardsOf(now, reg,
		[]agentRow{{Name: "alpha", Group: gRunning, Short: "alpha session", Since: "3m"}},
		[]nextRow{{Project: "alpha", What: "ship it"}},
		[]mrRow{{Project: "alpha", Ref: "!7", Role: "REVIEWER", Author: "sam", Updated: "2026-10-27T09:00:00Z"}})
	if len(cards) != 1 || len(quiet) != 1 || quiet[0] != "beta" {
		t.Fatalf("cards %+v quiet %v", cards, quiet)
	}
	c := cards[0]
	if c.State != "running" || c.Now != "alpha session" || c.Next != "ship it" || c.Done != 2 || c.Of != 5 || c.Open != "1 MR" {
		t.Errorf("card %+v", c)
	}
	if c.roadText() != "2/5 · no date" || !strings.Contains(c.Block, "!7 sam waits for your review") {
		t.Errorf("road %q block %q", c.roadText(), c.Block)
	}
	if rows := projectRows(cards); len(rows) != 1 || rows[0].Done != 2 || rows[0].Date != "" {
		t.Errorf("home rows %+v", rows)
	}
	m := projModel(170, 46)
	m.fx = &fixtures{Cards: cards, Quiet: quiet, Live: true}
	if out := m.View(); !strings.Contains(out, "No dated roadmap") {
		t.Errorf("live view lacks the no-date note:\n%s", out)
	}
}

func TestProjectsTimelineScale(t *testing.T) {
	m := projModel(170, 46)
	now := time.Date(2026, 10, 29, 14, 0, 0, 0, time.Local)
	pv := m.projectsView()
	pv.now = now
	start, months := pv.window()
	if start.Month() != time.October || start.Day() != 1 || months != 5 {
		t.Errorf("window %v %d", start, months)
	}
	pv.zoom = true
	if _, months = pv.window(); months != 3 {
		t.Errorf("zoom %d", months)
	}
}

func TestRegistryRoadDecodes(t *testing.T) {
	var reg registry
	raw := `{"projects":[{"name":"a","short":"a","scope":"work","road":{"done":2,"total":5,"current":"0.3 - Next","exit":"tests pass"}},{"name":"b","short":"b","road":null}]}`
	if err := json.Unmarshal([]byte(raw), &reg); err != nil {
		t.Fatal(err)
	}
	r := reg.Projects[0].Road
	if r == nil || r.Done != 2 || r.Total != 5 || r.Current != "0.3 - Next" || reg.Projects[1].Road != nil {
		t.Errorf("road %+v", reg.Projects)
	}
}
