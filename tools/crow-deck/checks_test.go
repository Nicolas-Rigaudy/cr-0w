package main

import (
	"testing"
	"time"
)

func testChecks() *checksSnap {
	return &checksSnap{Rows: []checkRow{
		{ID: "sprint:Sprint 33:2026-10-29", Check: "sprint", What: "Sprint 33 ends Thu 29 Oct", Why: "no planning in the calendar", Action: "propose planning", Due: "2026-10-29"},
		{ID: "1on1:Alex:2026-10-26", Check: "one_on_one", What: "1:1 with Alex", Why: "due Mon 26 Oct, none booked", Action: "propose a 1:1", Due: "2026-10-26"},
	}}
}

func TestCheckDecideRows(t *testing.T) {
	now := time.Date(2026, 10, 29, 15, 0, 0, 0, time.Local)
	rows := checkDecide(testChecks(), now)
	if len(rows) != 2 {
		t.Fatalf("rows %+v", rows)
	}
	r := rows[0]
	if r.Project != "checks" || r.What != "Sprint 33 ends Thu 29 Oct, no planning in the calendar" || r.Why != "propose planning" || r.Age != "today" || r.Actor != "me" || r.Perso {
		t.Errorf("sprint row %+v", r)
	}
	if rows[1].Age != "3d" {
		t.Errorf("1:1 age %q", rows[1].Age)
	}
	if checkDecide(nil, now) != nil {
		t.Error("no snapshot gives rows")
	}
}

func TestHomeShowsChecksInYouDecide(t *testing.T) {
	m := liveModel(t)
	m.page = 0
	m = send(m, srcMsg{src: srcChecks, at: testNow(), val: testChecks()})
	gs := m.homeView().groups()
	if len(gs) == 0 || gs[0].name != "You decide" {
		t.Fatalf("groups %+v", gs)
	}
	found := false
	for _, r := range gs[0].rows {
		if r.Project == "checks" && r.Why == "propose planning" {
			found = true
		}
	}
	if !found {
		t.Errorf("no sprint check in You decide: %+v", gs[0].rows)
	}
	if !hasText(m.View(), "Sprint 33 ends Thu 29 Oct") {
		t.Error("Home does not draw the check row")
	}
}
