package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func histModel(w, h int) model {
	m := newModel(defaultConfig())
	m.page, m.w, m.h, m.mode, m.hist = 1, w, h, "all", true
	return m
}

func TestSlashOpensHistoryAndEscLeaves(t *testing.T) {
	m := newModel(defaultConfig())
	m.page = 1
	m = send(m, keyMsg("/"))
	if !m.hist || !m.histTyping {
		t.Fatalf("/ must open History in typing mode: hist=%v typing=%v", m.hist, m.histTyping)
	}
	m = send(m, keyMsg("r"), keyMsg("e"), keyMsg("t"))
	if m.histF.query != "ret" {
		t.Errorf("query %q", m.histF.query)
	}
	m = send(m, keyMsg("esc"))
	if m.histTyping || !m.hist {
		t.Errorf("esc while typing stops typing only: typing=%v hist=%v", m.histTyping, m.hist)
	}
	m = send(m, keyMsg("esc"))
	if m.hist {
		t.Error("second esc must go back to live")
	}
}

func TestHistorySearchFiltersRows(t *testing.T) {
	m := histModel(170, 46)
	all := len(m.histView().rows())
	m.histF.query = "retry"
	got := m.histView().rows()
	if len(got) == 0 || len(got) >= all {
		t.Fatalf("search retry: %d of %d rows", len(got), all)
	}
	for _, s := range got {
		if !strings.Contains(strings.ToLower(s.Task), "retry") {
			t.Errorf("row %q does not match", s.Task)
		}
	}
	m.histF.query = "zzz"
	if !hasText(m.View(), "No session matches") {
		t.Error("no-match state missing")
	}
}

func TestHistoryChipsCycleAndClear(t *testing.T) {
	m := histModel(170, 46)
	m = send(m, keyMsg("p"))
	if m.histF.project == "" {
		t.Fatal("p must pick a project")
	}
	for _, s := range m.histView().rows() {
		if s.Project != m.histF.project {
			t.Errorf("row %q outside the project filter", s.Project)
		}
	}
	m = send(m, keyMsg("d"), keyMsg("s"), keyMsg("M"))
	if m.histF.days != 7 || m.histF.status != "done" || m.histF.model == "" {
		t.Errorf("chips not set: %+v", m.histF)
	}
	m = send(m, keyMsg("x"))
	if m.histF.active() {
		t.Errorf("x must clear the filters: %+v", m.histF)
	}
	// every value cycles back to any
	vals := []string{"a", "b"}
	if nextValue("b", vals) != "" || nextValue("", vals) != "a" || nextValue("a", vals) != "b" {
		t.Error("nextValue does not cycle through any")
	}
}

func TestHistoryDateRange(t *testing.T) {
	m := histModel(170, 46)
	m.histF.days = 1
	for _, s := range m.histView().rows() {
		if dayLabel(s.Start, m.now) != "Today" {
			t.Errorf("%s is not today", s.Start)
		}
	}
}

func TestHistoryModeHidesPerso(t *testing.T) {
	m := histModel(170, 46)
	m.mode = "work"
	for _, s := range m.histView().rows() {
		if s.personal() {
			t.Errorf("perso row %q in work mode", s.Task)
		}
	}
}

func TestHistoryActionsAreReadOnly(t *testing.T) {
	m := histModel(170, 46)
	for _, k := range []string{"enter", "o", "r", "y", "h"} {
		m2 := send(m, keyMsg(k))
		if !strings.Contains(m2.flash, "actions slice") {
			t.Errorf("%s: flash %q", k, m2.flash)
		}
		if !m2.hist {
			t.Errorf("%s left History", k)
		}
	}
}

func TestHistorySelectionAndPane(t *testing.T) {
	m := histModel(170, 46)
	m = send(m, keyMsg("j"), keyMsg("j"))
	if m.histSel != 2 {
		t.Errorf("sel %d", m.histSel)
	}
	for i := 0; i < 40; i++ {
		m = send(m, keyMsg("j"))
	}
	if n := len(m.histView().rows()); m.histSel != n-1 {
		t.Errorf("sel %d, want last %d", m.histSel, n-1)
	}
	small := histModel(90, 30)
	small = send(small, keyMsg("tab"))
	if small.histPane != hpPreview || !hasText(small.View(), "end state") {
		t.Error("tab must show the preview in single-pane mode")
	}
}

func TestHistoryPreviewWithoutDetail(t *testing.T) {
	m := histModel(170, 46)
	for i := range m.fx.Hist {
		m.fx.Hist[i].Detail = false
	}
	if !hasText(m.View(), "not recorded by the sources yet") {
		t.Error("a row without detail must say so")
	}
}

func TestLiveSwitchOnAgentsTab(t *testing.T) {
	m := histModel(170, 46)
	if !hasText(m.View(), "live /") || !hasText(m.View(), "history") {
		t.Error("switch missing on Agents")
	}
	other := newModel(defaultConfig())
	other.page = 2
	if hasText(other.View(), "live /") {
		t.Error("switch must show on Agents only")
	}
}

func TestHistFromUsage(t *testing.T) {
	now := time.Date(2026, 10, 29, 14, 0, 0, 0, time.UTC)
	u := &usageSnap{Sessions: []usageSession{
		{ID: "a", Account: "x", Project: "p", Title: "p: older task", Start: "2026-10-28T09:00:00Z", End: "2026-10-29T09:00:00Z", Models: []string{"claude-opus-5"}, CostUSD: 3, Activity: map[string][]int{"2026-10-28": {1, 2, 3}}},
		{ID: "b", Account: "x", Project: "p", Title: "p: now task", Start: "2026-10-29T13:00:00Z", End: "2026-10-29T13:59:00Z", Models: []string{"claude-sonnet-5", "claude-haiku-5"}, CostUSD: 1},
	}}
	h := histOf(u, nil, now)
	if len(h) != 2 || h[0].ID != "b" {
		t.Fatalf("order: %+v", h)
	}
	if !h[0].Active || h[1].Active {
		t.Error("only the session that ended a minute ago is active")
	}
	if h[0].Model != "mixed" || h[1].Model != "opus" {
		t.Errorf("models %q %q", h[0].Model, h[1].Model)
	}
	if h[1].Took != 30*time.Minute {
		t.Errorf("took %v, want the 3 active buckets", h[1].Took)
	}
	if h[1].Task != "older task" {
		t.Errorf("task %q", h[1].Task)
	}
	if histOf(nil, nil, now) != nil {
		t.Error("nil usage gives nil history")
	}
}

func TestHistoryStatesFit(t *testing.T) {
	withANSI(t)
	for ti, th := range themes {
		for _, sz := range sizes {
			for _, st := range syncStates {
				m := histModel(sz[0], sz[1])
				m.th = ti
				m.setSync(st)
				checkFit(t, fmt.Sprintf("%s History %dx%d %s", th.name, sz[0], sz[1], st), m.View(), sz[0], sz[1])
			}
			m := histModel(sz[0], sz[1])
			m.th, m.histTyping, m.histF.query = ti, true, "retry"
			checkFit(t, fmt.Sprintf("%s History search %dx%d", th.name, sz[0], sz[1]), m.View(), sz[0], sz[1])
		}
	}
}

func TestHistoryEmptyAndSourceError(t *testing.T) {
	m := histModel(170, 46)
	m.setSync(stEmpty)
	m.hist = true
	if !hasText(m.View(), "No session yet") {
		t.Error("empty state missing")
	}
	m = histModel(170, 46)
	m.fx = emptyFixtures()
	m.fx.Live, m.sync = true, stError
	if !hasText(m.View(), "usage source did not answer") {
		t.Error("source error missing")
	}
}
