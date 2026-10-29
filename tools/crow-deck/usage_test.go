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

func TestUsageShowsEverySection(t *testing.T) {
	out := usageModel(170, 46).View()
	for _, want := range []string{"Accounts and policy", "Forecast", "Daily spend", "Savings", "Cost today", "API-equivalent $, not billed", "per session", "per project", "per model", "total saved"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestUsageSinglePaneWalk(t *testing.T) {
	m := usageModel(90, 30)
	for i, want := range []string{"Accounts and policy", "Forecast", "Daily spend", "Savings", "Cost today"} {
		if !strings.Contains(m.View(), want) {
			t.Errorf("pane %d lacks %q", i, want)
		}
		nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
		m = nm.(model)
	}
	if m.usePane != 0 {
		t.Errorf("tab did not wrap, pane %d", m.usePane)
	}
}

func TestUsageKeys(t *testing.T) {
	m := usageModel(170, 46)
	press := func(k string) {
		nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)})
		m = nm.(model)
	}
	press("d")
	if !m.useFortnight || !strings.Contains(m.View(), "14 days") {
		t.Error("d did not switch to 14 days")
	}
	press("p")
	if !strings.Contains(m.flash, "not wired") && !strings.Contains(m.flash, "arrives") {
		t.Errorf("policy key flash %q", m.flash)
	}
	press("j")
	if m.useSel != 1 {
		t.Errorf("j: account %d", m.useSel)
	}
	for i := 0; i < 5; i++ {
		press("j")
	}
	if m.useSel != len(m.fx.Gauges)-1 {
		t.Errorf("selection ran past the accounts: %d", m.useSel)
	}
}

func TestUsageRefreshFetchesUsage(t *testing.T) {
	m := usageModel(170, 46)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	if cmd == nil {
		t.Error("r on Usage starts no fetch")
	}
}

func TestUsageStates(t *testing.T) {
	m := usageModel(170, 46)
	m.setSync(stLoading)
	if !strings.Contains(m.View(), "syncing") {
		t.Error("loading state")
	}
	m = usageModel(170, 46)
	m.setSync(stEmpty)
	out := m.View()
	if !strings.Contains(out, "No usage data") || !strings.Contains(out, "No spend recorded") {
		t.Errorf("empty state:\n%s", out)
	}
}

func TestUsageViewOfRealPayload(t *testing.T) {
	ls := newLiveState()
	ls.usage = &usageSnap{
		Today: usageTotals{CostUSD: 12},
		History: []struct {
			Date    string  `json:"date"`
			CostUSD float64 `json:"cost_usd"`
		}{{Date: "2026-10-27", CostUSD: 10}, {Date: "2026-10-28", CostUSD: 20}},
	}
	v := usageViewOf(ls.usage, nil, testNow())
	if v == nil || len(v.Days) != 2 || v.Days[1].USD != 20 {
		t.Fatalf("days %+v", v)
	}
	if spark(v.Days) != "▄█" {
		t.Errorf("spark %q", spark(v.Days))
	}
}
