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

func TestProjectsShowsListDetailAndTimeline(t *testing.T) {
	for _, sz := range [][2]int{{235, 62}, {118, 60}} {
		out := strings.ToLower(projModel(sz[0], sz[1]).View())
		for _, want := range []string{"needs you", "running or ready", "health", "left", "detail: lint", "next action", "milestone",
			"roadmap and work", "today", "worked", "time split", "▼", "◆"} {
			if !strings.Contains(out, want) {
				t.Errorf("%dx%d missing %q", sz[0], sz[1], want)
			}
		}
		if strings.Contains(out, "not wired") {
			t.Errorf("%dx%d shows a not wired box", sz[0], sz[1])
		}
	}
	out := strings.ToLower(projModel(235, 62).View())
	for _, want := range []string{"roadmap", "signals", "quiet", "→ ", "c continue", "h handoff", "sort health"} {
		if !strings.Contains(out, want) {
			t.Errorf("board missing %q", want)
		}
	}
}

func TestProjectsActionsOnlyInDetail(t *testing.T) {
	out := projModel(235, 62).View()
	if n := strings.Count(out, "continue"); n != 1 {
		t.Errorf("%d action rows, want 1", n)
	}
	m := projModel(235, 62)
	m.prjSel = 1
	if out := strings.ToLower(m.View()); !strings.Contains(out, "open mrs") || !strings.Contains(out, "acts: ") {
		t.Error("detail of a project with MRs lacks who must act")
	}
}

func TestProjectsKeys(t *testing.T) {
	m := projModel(235, 62)
	press := func(k string) {
		nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)})
		m = nm.(model)
	}
	press("j")
	if m.prjSel != 1 {
		t.Errorf("j: card %d", m.prjSel)
	}
	for i := 0; i < 20; i++ {
		press("j")
	}
	if m.prjSel != len(m.fx.Cards)-1 {
		t.Errorf("selection ran past the cards: %d", m.prjSel)
	}
	press("c")
	if !strings.Contains(m.flash, "continue on") || strings.Contains(m.flash, "not wired") {
		t.Errorf("flash %q", m.flash)
	}
	_, n0, _ := m.projectsView().window(false)
	press("z")
	if !m.prjZoom {
		t.Error("z did not zoom")
	}
	if _, n1, _ := m.projectsView().window(false); n1 != 2*n0 {
		t.Errorf("zoom window %d days, want %d", n1, 2*n0)
	}
}

func TestProjectsModeFilter(t *testing.T) {
	m := projModel(235, 62)
	m.mode = "work"
	out := m.View()
	if !strings.Contains(out, "atlas-platform") || strings.Contains(out, "vault-plugin") {
		t.Errorf("work mode shows personal cards:\n%s", out)
	}
}

func TestProjectsLiveFromRegistry(t *testing.T) {
	reg := &registry{Projects: []regProject{
		{Name: "alpha", Short: "alpha", Scope: "work", Road: &regRoad{Done: 2, Total: 5, Current: "0.3 - Next", Exit: "tests pass"},
			Work: map[string]int{"2026-10-27": 3}},
		{Name: "beta", Short: "beta", Scope: "personal"},
		{Name: "gamma", Short: "gamma", Scope: "work"},
	}}
	reg.Projects[2].Release.Date, reg.Projects[2].Release.Label = "2026-12-01", "Prod launch"
	now := testNow()
	cards := projectCardsOf(now, reg,
		[]agentRow{{Name: "alpha", Group: gRunning, Short: "alpha session", Since: "3m"}},
		[]nextRow{{Project: "alpha", What: "ship it", Actor: "sam"}, {Project: "alpha", What: "write the notes", Actor: "me"}},
		[]mrRow{{Project: "alpha", Ref: "!7", Role: "REVIEWER", Author: "sam", Updated: "2026-10-27T09:00:00Z"}})
	if len(cards) != 3 {
		t.Fatalf("cards %+v", cards)
	}
	by := map[string]projCard{}
	for _, c := range cards {
		by[c.Name] = c
	}
	a, b, g := by["alpha"], by["beta"], by["gamma"]
	if a.Group != pgNeeds || a.Glyph != "needs" || a.Done != 2 || a.Of != 5 || a.Health != "no date" || len(a.MRs) != 1 || a.Work["2026-10-27"] != 3 {
		t.Errorf("alpha %+v", a)
	}
	if n := a.lead(); n == nil || n.What != "write the notes" {
		t.Errorf("lead %+v", n)
	}
	if b.Group != pgQuiet || b.Health != "idle" {
		t.Errorf("beta %+v", b)
	}
	if g.Target.Format("2006-01-02") != "2026-12-01" || g.Mile != "Prod launch" || g.Health != "on time" || !strings.HasPrefix(g.left(now), "D-") {
		t.Errorf("gamma %+v left %q", g, g.left(now))
	}
	rows := homeProjects(now, reg, cards)
	if len(rows) != 3 || rows[0].Done != 2 || rows[0].Now != "0.3 - Next" || rows[2].Date != "Dec 1" {
		t.Errorf("home rows %+v", rows)
	}
	m := projModel(235, 62)
	m.now = now
	m.fx = &fixtures{Cards: cards, Live: true}
	out := m.View()
	for _, want := range []string{"◆ prod launch Dec 1", "no date", "worked  1d", "alpha 100%"} {
		if !strings.Contains(out, want) {
			t.Errorf("live view lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "waiting on sam: ship it") {
		t.Error("an other-actor row led while the owner has one")
	}
}

func TestProjectsHealth(t *testing.T) {
	now := testNow()
	d := func(n int) time.Time { return day(now).AddDate(0, 0, n) }
	cases := []struct {
		c    projCard
		want string
	}{
		{projCard{}, "no date"},
		{projCard{Target: d(-1)}, "late"},
		{projCard{Target: d(30), Start: d(-90), Done: 1, Of: 5}, "at risk"},
		{projCard{Target: d(30), Start: d(-30), Done: 3, Of: 5}, "on time"},
		{projCard{Target: d(20), Glyph: "blocked"}, "at risk"},
	}
	for _, tc := range cases {
		if got := cardHealth(now, tc.c); got != tc.want {
			t.Errorf("%+v: %q, want %q", tc.c, got, tc.want)
		}
	}
}

func TestProjectsTimelineWindow(t *testing.T) {
	pv := projModel(235, 62).projectsView()
	pv.now = time.Date(2026, 10, 29, 14, 0, 0, 0, time.Local)
	start, days, past := pv.window(false)
	if start.Weekday() != time.Monday || days != 98 || past != 5 || start.Format("01-02") != "09-21" {
		t.Errorf("window %v %d %d", start, days, past)
	}
	if _, days, _ = pv.window(true); days != 63 {
		t.Errorf("portrait %d", days)
	}
	if heatGlyph(1) != "▃" || heatGlyph(5) != "▅" || heatGlyph(10) != "▇" || heatGlyph(40) != "█" {
		t.Error("heat glyphs")
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

func TestProjectsTimelineUsesSpareRows(t *testing.T) {
	if got := heatCells(16, 2); got[0] != "█" || got[1] != "█" {
		t.Errorf("16 commits on 2 rows: %q", got)
	}
	if got := heatCells(1, 3); got[0] != " " || got[2] == " " {
		t.Errorf("1 commit on 3 rows: %q", got)
	}
	if hh, _ := heatRows(1, 20); hh != 4 {
		t.Errorf("one project, 20 rows: %d-row bars", hh)
	}
	if hh, gap := heatRows(6, 0); hh != 1 || gap != 0 {
		t.Errorf("no room: %d %d", hh, gap)
	}
	// work mode has one active project: the timeline stays compact and the detail lists
	// every next action instead of "+N more"
	m := projModel(235, 62)
	m.mode = "work"
	if out := m.View(); strings.Contains(out, "more in the open threads") {
		t.Error("work mode still truncates the next actions")
	}
}
