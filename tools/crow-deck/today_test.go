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
	s.Important = append(s.Important, struct {
		Subject string `json:"subject"`
		From    string `json:"from"`
		Unread  bool   `json:"unread"`
		High    bool   `json:"high"`
	}{"Tr :  Terms   v3", "Jo", true, false})
	rows, mail := todayRows(s, now)
	defer func() {
		s.Flagged = append(s.Flagged, struct {
			Subject string `json:"subject"`
			From    string `json:"from"`
			Crow    bool   `json:"crow"`
		}{"Budget", "Al", true})
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
