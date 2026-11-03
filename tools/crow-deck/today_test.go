package main

import (
	"strings"
	"testing"
	"time"
)

func TestTodayRows(t *testing.T) {
	now := time.Date(2026, 11, 3, 15, 52, 0, 0, time.Local)
	s := &todaySnap{
		Events: []todayEvent{
			{Start: "2026-11-03T00:00", End: "2026-11-04T00:00", Subject: "Office", AllDay: true},
			{Start: "2026-11-03T11:30", End: "2026-11-03T12:00", Subject: "Daily"},
			{Start: "2026-11-03T16:00", End: "2026-11-03T16:30", Subject: "1:1"},
			{Start: "2026-11-03T17:30", End: "2026-11-03T18:00", Subject: "Dry run"},
		},
		Tomorrow: []todayEvent{{Start: "2026-11-04T09:00", End: "2026-11-04T10:00", Subject: "Travel"}},
		Unread:   4,
		Drafts:   1,
	}
	s.Important = append(s.Important, todayMail{Subject: "Tr :  Terms   v3", From: "Jo", Unread: true})
	rows, mail := todayRows(s, now)
	defer func() {
		s.Flagged = append(s.Flagged, todayFlag{Subject: "Budget", From: "Al", Crow: true})
		rows, mail := todayRows(s, now)
		if rows[5].What != "Flagged" || rows[6].TagState != "crow" || mail[1] != "1 flagged" {
			t.Fatalf("flagged wins over the picker: %+v %v", rows[5:7], mail)
		}
	}()
	if rows[0].When != "all day" || rows[1].TagState != "past" {
		t.Fatalf("head rows: %+v", rows[:2])
	}
	if rows[2].TagState != "next" || rows[2].Tag != "in 8 min" || rows[3].TagState != "" {
		t.Fatalf("next marker: %+v", rows[2:4])
	}
	var kinds []string
	for _, r := range rows[4:] {
		kinds = append(kinds, r.Kind+":"+r.What)
	}
	if got := strings.Join(kinds, ","); got != "gap:,head:Mail to look at,flag:Tr : Terms v3,gap:,head:Tomorrow,:Travel" {
		t.Fatalf("sections: %s", got)
	}
	if strings.Join(mail, " · ") != "mail 4 unread · 1 to look at · 1 drafts" {
		t.Fatalf("mail: %v", mail)
	}
	if r, m := todayRows(nil, now); r != nil || m != nil {
		t.Fatal("nil snap must hide Today")
	}
}

func TestSnoozedFlagsHideUntilTheyReturn(t *testing.T) {
	now := time.Date(2026, 11, 3, 15, 0, 0, 0, time.UTC)
	s := &todaySnap{Flagged: []todayFlag{
		{ID: "a", Subject: "Snoozed", Until: "2026-11-04T07:00:00+00:00"},
		{ID: "b", Subject: "Back", Until: "2026-11-03T07:00:00+00:00"},
		{ID: "c", Subject: "Plain"},
	}}
	rows, mail := todayRows(s, now)
	var got []string
	for _, r := range rows {
		if r.Kind == "flag" {
			got = append(got, r.ID)
		}
	}
	if strings.Join(got, ",") != "b,c" || mail[1] != "2 flagged" {
		t.Fatalf("snooze filter: %v %v", got, mail)
	}
}

// todayModel is the Home page with the Today pane focused on two mails between events.
func todayModel(t *testing.T) model {
	t.Helper()
	m := newModel(defaultConfig())
	m.w, m.h, m.page, m.focus = 170, 46, 0, pToday
	m.fx.Today = []todayRow{
		{When: "11:30–12:00", What: "Daily"},
		{Kind: "gap"}, {Kind: "head", What: "Flagged"},
		{Kind: "flag", What: "Budget", Subject: "Budget", Tag: "Al", ID: "m1", Link: "https://outlook/m1"},
		{Kind: "flag", What: "Terms v3", Subject: "Terms v3", Tag: "Jo", ID: "m2", Link: "https://outlook/m2"},
		{Kind: "gap"}, {Kind: "head", What: "Tomorrow"}, {When: "09:00", What: "Travel"},
	}
	return m
}

func TestTodayCursorWalksMailRowsOnly(t *testing.T) {
	m := todayModel(t)
	if n := m.homeView().items(pToday); n != 2 {
		t.Fatalf("cursor must walk the 2 mail rows, got %d", n)
	}
	if r := m.homeView().curMail(); r == nil || r.ID != "m1" {
		t.Fatalf("first mail: %+v", r)
	}
	m = send(m, keyMsg("j"), keyMsg("j"), keyMsg("j"))
	if r := m.homeView().curMail(); r == nil || r.ID != "m2" {
		t.Fatalf("cursor must stop on the last mail: %+v", r)
	}
	if !hasText(m.View(), "Open in Outlook") && !hasText(m.View(), "open in Outlook") {
		t.Fatal("footer must show the mail keys when Today has focus")
	}
}

func TestTodayKeysRunTheirCommands(t *testing.T) {
	cases := map[string][]string{
		"o":     {"crow-open", "https://outlook/m1"},
		"x":     {"crow-today", "done", "m1"},
		"s":     {"crow-today", "snooze", "m1"},
		"enter": {"crow", "handoff", "cr-0w", "--task"},
	}
	for k, want := range cases {
		calls := recordRuns(t, "ok", nil)
		_, cmd := todayModel(t).Update(keyMsg(k))
		if cmd == nil {
			t.Fatalf("%s: no command", k)
		}
		msg := cmd()
		if len(*calls) != 1 {
			t.Fatalf("%s: %d runs", k, len(*calls))
		}
		got := append([]string{(*calls)[0].name}, (*calls)[0].args...)
		if strings.Join(got[:len(want)], " ") != strings.Join(want, " ") {
			t.Fatalf("%s ran %v", k, got)
		}
		if k == "enter" && (!strings.Contains(got[4], "crow-today thread m1") || !strings.Contains(got[4], "Never send")) {
			t.Fatalf("enter task: %s", got[4])
		}
		if _, ok := msg.(mailActMsg); ok != (k == "x" || k == "s") {
			t.Fatalf("%s: only done and snooze refetch Today, got %T", k, msg)
		}
	}
}

func TestTodayKeysNeedAMailID(t *testing.T) {
	calls := recordRuns(t, "ok", nil)
	m := todayModel(t)
	m.fx.Today[3].ID = ""
	m = send(m, keyMsg("x"))
	if len(*calls) != 0 || !strings.Contains(m.flash, "no Graph id") {
		t.Fatalf("mock rows must not run anything: %v %q", *calls, m.flash)
	}
}
