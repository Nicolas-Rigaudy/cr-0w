package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func usageModel(w, h int) model {
	m := newModel(defaultConfig())
	m.page, m.w, m.h = 4, w, h
	return m
}

func pressUsage(m model, keys ...string) model {
	for _, k := range keys {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
		switch k {
		case "tab":
			msg = tea.KeyMsg{Type: tea.KeyTab}
		case "enter":
			msg = tea.KeyMsg{Type: tea.KeyEnter}
		case "esc":
			msg = tea.KeyMsg{Type: tea.KeyEsc}
		}
		nm, _ := m.Update(msg)
		m = nm.(model)
	}
	return m
}

func TestUsageShowsEveryBox(t *testing.T) {
	for _, sz := range [][2]int{{235, 62}, {118, 60}} {
		out := usageModel(sz[0], sz[1]).View()
		for _, want := range []string{"Runway", "Daily spend", "Saved", "Ledger", "By project", "By model", "now", "week resets", "prompt cache", "model routing", "rtk"} {
			if !strings.Contains(out, want) {
				t.Errorf("%dx%d: missing %q", sz[0], sz[1], want)
			}
		}
		if strings.Contains(out, "not wired") {
			t.Errorf("%dx%d: a not wired box", sz[0], sz[1])
		}
		if n := len(strings.Split(out, "\n")); n > sz[1] {
			t.Errorf("%dx%d: %d lines", sz[0], sz[1], n)
		}
	}
}

func TestUsageChartAxes(t *testing.T) {
	out := usageModel(235, 62).View()
	for _, want := range []string{"$0 ┤", "└", "Today"} {
		if !strings.Contains(out, want) && !(want == "Today" && strings.Contains(out, mockNow().Format("Mon 2"))) {
			t.Errorf("daily chart lacks %q", want)
		}
	}
}

func TestUsagePeriodKeys(t *testing.T) {
	m := usageModel(235, 62)
	if m.usePeriod != perToday {
		t.Fatalf("starts on %d", m.usePeriod)
	}
	m = pressUsage(m, "p")
	if m.usePeriod != perWeek || !strings.Contains(m.View(), "By project · week") || !strings.Contains(m.View(), "period: week") {
		t.Error("p did not move to the week")
	}
	mode := m.mode
	m = pressUsage(m, "p", "m")
	if m.usePeriod != perMonth || m.mode == mode {
		t.Errorf("p: period %d; m kept mode %q", m.usePeriod, m.mode)
	}
	m = pressUsage(m, "p")
	if m.usePeriod != perToday {
		t.Error("p did not wrap to today")
	}
}

func TestUsageSinglePaneWalk(t *testing.T) {
	m := usageModel(90, 30)
	for i, want := range []string{"Runway", "Daily spend", "Saved", "Ledger", "By project", "By model"} {
		if !strings.Contains(m.View(), want) {
			t.Errorf("pane %d lacks %q", i, want)
		}
		m = pressUsage(m, "tab")
	}
	if m.usePane != 0 {
		t.Errorf("tab did not wrap, pane %d", m.usePane)
	}
}

func TestUsageProjectFilter(t *testing.T) {
	m := usageModel(235, 62)
	m.usePane = upProjects
	m = pressUsage(m, "p", "enter")
	top := m.usageView().per().Projects[0].Name
	if m.useFilter != top {
		t.Fatalf("filter %q, want %q", m.useFilter, top)
	}
	for _, e := range m.usageView().ledgerEntries() {
		if e.Project != top {
			t.Errorf("entry of %s under filter %s", e.Project, top)
		}
	}
	m = pressUsage(m, "esc")
	if m.useFilter != "" {
		t.Error("esc kept the filter")
	}
}

func TestUsageLedgerFoldAndOpen(t *testing.T) {
	m := usageModel(235, 62)
	m.usePane = upLedger
	m = pressUsage(m, "p")
	rows := func() []ledRow { _, r := m.usageView().ledgerSpec(200, 1000, false); return r }
	before := len(rows())
	m = pressUsage(m, "enter") // first group: open -> all
	m = pressUsage(m, "enter") // all -> closed
	if len(rows()) >= before {
		t.Errorf("closing the first group kept %d rows (was %d)", len(rows()), before)
	}
	m = pressUsage(m, "j", "j", "o") // row 2: the first entry of the second group
	if m.page != 1 || !m.hist || m.histF.query == "" {
		t.Errorf("o did not open the history: page %d hist %v query %q", m.page, m.hist, m.histF.query)
	}
}

func TestUsageFlatAndTokens(t *testing.T) {
	m := pressUsage(usageModel(235, 62), "g", "c")
	if !m.useFlat || !m.useTok {
		t.Fatal("g or c did not toggle")
	}
	for _, r := range m.usageView().ledRows(100, 1) {
		if r.kind != 1 {
			t.Error("flat ledger has a group row")
		}
	}
}

func TestUsageRunwayPan(t *testing.T) {
	m := pressUsage(usageModel(235, 62), "h")
	if m.useShift != 6 {
		t.Errorf("h: shift %d", m.useShift)
	}
	m = pressUsage(m, "l", "l")
	if m.useShift != -6 {
		t.Errorf("l: shift %d", m.useShift)
	}
}

func TestUsageRefreshFetchesUsage(t *testing.T) {
	m := usageModel(235, 62)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	if cmd == nil {
		t.Error("r on Usage starts no fetch")
	}
}

func TestUsageStates(t *testing.T) {
	m := usageModel(235, 62)
	m.setSync(stLoading)
	if !strings.Contains(m.View(), "syncing") {
		t.Error("loading state")
	}
	m = usageModel(235, 62)
	m.setSync(stEmpty)
	out := m.View()
	if !strings.Contains(out, "No spend recorded") || !strings.Contains(out, "No session in this period") {
		t.Errorf("empty state:\n%s", out)
	}
}

func TestUsageViewOfRealPayload(t *testing.T) {
	wk := 12.5
	u := &usageSnap{
		Today:   usageTotals{CostUSD: 12, ByProject: []usageName{{Name: "?", CostUSD: 0}, {Name: "a", CostUSD: 12}}},
		History: []dayUse{{Date: "2026-10-27", CostUSD: 10}, {Date: "2026-10-28", CostUSD: 20}},
		Sessions: []usageSession{{ID: "s1", Account: "main", Project: "a", Title: "fix it", Start: "2026-10-28T09:00:00Z", End: "2026-10-28T10:00:00Z",
			ByDay: map[string]dayUse{"2026-10-28": {CostUSD: 3, Tokens: 1000}}, ActivityTokens: map[string]map[string]int{"2026-10-28": {"54": 500}}}},
	}
	u.Sessions[0].LimitShare.Week = &wk
	v := usageViewOf(u, nil, testNow())
	if v == nil || len(v.Days) != 2 || v.Days[1].USD != 20 || v.Days[0].Label != "Tue 27" {
		t.Fatalf("days %+v", v)
	}
	if len(v.Periods[perToday].Projects) != 1 {
		t.Errorf("the empty ? row was kept: %+v", v.Periods[perToday].Projects)
	}
	if len(v.Entries) != 1 || v.Entries[0].Week != 12.5 || v.Entries[0].Title != "fix it" {
		t.Errorf("entries %+v", v.Entries)
	}
	if len(v.Activity["main"]) != 1 || v.Since != "2026-10-28" {
		t.Errorf("activity %+v since %q", v.Activity, v.Since)
	}
}
