package main

// The Today rail: calendar and mail counts from bin/crow-today (Microsoft Graph).

import (
	"fmt"
	"time"
)

const srcToday = "today"

type todayEvent struct {
	Start   string `json:"start"` // local "2006-01-02T15:04"
	End     string `json:"end"`
	Subject string `json:"subject"`
	AllDay  bool   `json:"all_day"`
}

type todaySnap struct {
	Events   []todayEvent `json:"events"`
	Tomorrow []todayEvent `json:"tomorrow"`
	Unread   int          `json:"unread"`
	Flagged  []struct {
		Subject string `json:"subject"`
		From    string `json:"from"`
	} `json:"flagged"`
	Drafts int `json:"drafts"`
}

func fetchToday() (*todaySnap, error) {
	var s todaySnap
	return &s, runJSON(60*time.Second, nil, &s, "crow-today", "json")
}

// todayRows lays out the rail: today's events (the next one marked with its countdown),
// then flagged mail, then tomorrow. mail is the pane's right label.
func todayRows(s *todaySnap, now time.Time) (rows []todayRow, mail []string) {
	if s == nil {
		return nil, nil
	}
	const layout = "2006-01-02T15:04"
	marked := false
	for _, e := range s.Events {
		if e.AllDay {
			rows = append(rows, todayRow{When: "all day", What: e.Subject})
			continue
		}
		st, err1 := time.ParseInLocation(layout, e.Start, now.Location())
		en, err2 := time.ParseInLocation(layout, e.End, now.Location())
		if err1 != nil || err2 != nil {
			continue
		}
		r := todayRow{When: st.Format("15:04") + "–" + en.Format("15:04"), What: e.Subject}
		switch {
		case en.Before(now):
			r.TagState = "past"
		case !st.After(now):
			r.TagState, r.Tag = "next", "now"
			marked = true
		case !marked:
			r.TagState, r.Tag = "next", untilLabel(st.Sub(now))
			marked = true
		}
		rows = append(rows, r)
	}
	if len(s.Flagged) > 0 {
		rows = append(rows, todayRow{Kind: "gap"}, todayRow{Kind: "head", What: "Flagged"})
		for _, f := range s.Flagged {
			rows = append(rows, todayRow{Kind: "flag", What: f.Subject, Tag: f.From})
		}
	}
	if len(s.Tomorrow) > 0 {
		rows = append(rows, todayRow{Kind: "gap"}, todayRow{Kind: "head", What: "Tomorrow"})
		for _, e := range s.Tomorrow {
			when := "all day"
			if !e.AllDay && len(e.Start) >= 16 {
				when = e.Start[11:16]
			}
			rows = append(rows, todayRow{When: when, What: e.Subject})
		}
	}
	mail = append(mail, fmt.Sprintf("mail %d unread", s.Unread))
	if n := len(s.Flagged); n > 0 {
		mail = append(mail, fmt.Sprintf("%d flagged", n))
	}
	if s.Drafts > 0 {
		mail = append(mail, fmt.Sprintf("%d drafts", s.Drafts))
	}
	return rows, mail
}

func untilLabel(d time.Duration) string {
	m := int(d.Round(time.Minute).Minutes())
	if m < 60 {
		return fmt.Sprintf("in %d min", m)
	}
	return fmt.Sprintf("in %dh%02d", m/60, m%60)
}
