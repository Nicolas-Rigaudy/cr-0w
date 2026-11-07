package main

// The checks source: management gaps from bin/crow-checks (sprint rituals, 1:1s, mail,
// rituals, overdue notes), shown in Home "You decide".

import "time"

const srcChecks = "checks"

type checkRow struct {
	ID     string `json:"id"`
	Check  string `json:"check"`
	What   string `json:"what"`
	Why    string `json:"why"`
	Action string `json:"action"`
	Due    string `json:"due"`
}

type checksSnap struct {
	Rows    []checkRow        `json:"rows"`
	Skipped map[string]string `json:"skipped"`
}

func fetchChecks() (*checksSnap, error) {
	var s checksSnap
	return &s, runJSON(30*time.Second, nil, &s, "crow-checks", "json")
}

// checkDecide maps the check rows to You decide rows: what and why on the left, the
// suggested action in the why column, the age from the due date.
func checkDecide(s *checksSnap, now time.Time) []nextRow {
	if s == nil {
		return nil
	}
	var out []nextRow
	for _, r := range s.Rows {
		what := r.What
		if r.Why != "" {
			what += ", " + r.Why
		}
		out = append(out, nextRow{Project: "checks", What: what, Actor: "me", Why: r.Action, Age: threadAge(r.Due, now)})
	}
	return out
}
