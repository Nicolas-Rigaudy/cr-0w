package main

// The Today rail: calendar and mail counts from bin/crow-today (Microsoft Graph).

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const srcToday = "today"

type todayEvent struct {
	Start   string `json:"start"` // local "2006-01-02T15:04"
	End     string `json:"end"`
	Subject string `json:"subject"`
	AllDay  bool   `json:"all_day"`
}

type todaySnap struct {
	Events    []todayEvent `json:"events"`
	Tomorrow  []todayEvent `json:"tomorrow"`
	Unread    int          `json:"unread"`
	Flagged   []todayFlag  `json:"flagged"`
	Drafts    int          `json:"drafts"`
	Important []todayMail  `json:"important"`
}

type todayFlag struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	From    string `json:"from"`
	Crow    bool   `json:"crow"`
	WebLink string `json:"webLink"`
	Until   string `json:"until"` // a snoozed flag returns at this time
}

type todayMail struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	From    string `json:"from"`
	Unread  bool   `json:"unread"`
	High    bool   `json:"high"`
	WebLink string `json:"webLink"`
}

// awake drops the flags snoozed past now.
func awake(fs []todayFlag, now time.Time) []todayFlag {
	var out []todayFlag
	for _, f := range fs {
		if t, err := time.Parse(time.RFC3339, f.Until); err == nil && t.After(now) {
			continue
		}
		out = append(out, f)
	}
	return out
}

func fetchToday() (*todaySnap, error) {
	var s todaySnap
	return &s, runJSON(60*time.Second, nil, &s, "crow-today", "json")
}

// todayRows lays out the rail: today's events (the next one marked with its countdown),
// then the mail crow picked as worth a look, then tomorrow. mail is the pane's right label.
func todayRows(s *todaySnap, now time.Time) (rows []todayRow, mail []string) {
	if s == nil {
		return nil, nil
	}
	const layout = "2006-01-02T15:04"
	marked := false
	flagged := awake(s.Flagged, now)
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
	if len(flagged) > 0 {
		// The owner and crow share Outlook flags; crow's carry the crow category.
		rows = append(rows, todayRow{Kind: "gap"}, todayRow{Kind: "head", What: "Flagged"})
		for _, f := range flagged {
			st := "warn"
			if f.Crow {
				st = "crow"
			}
			rows = append(rows, todayRow{Kind: "flag", What: strings.Join(strings.Fields(f.Subject), " "), Tag: f.From, TagState: st,
				ID: f.ID, Link: f.WebLink, Subject: f.Subject})
		}
	} else if len(s.Important) > 0 {
		rows = append(rows, todayRow{Kind: "gap"}, todayRow{Kind: "head", What: "Mail to look at"})
		for _, m := range s.Important {
			st := ""
			if m.High {
				st = "warn"
			} else if !m.Unread {
				st = "past"
			}
			rows = append(rows, todayRow{Kind: "flag", What: strings.Join(strings.Fields(m.Subject), " "), Tag: m.From, TagState: st,
				ID: m.ID, Link: m.WebLink, Subject: m.Subject})
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
	if n := len(flagged); n > 0 {
		mail = append(mail, fmt.Sprintf("%d flagged", n))
	} else if n := len(s.Important); n > 0 {
		mail = append(mail, fmt.Sprintf("%d to look at", n))
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

// mailActMsg ends a done or snooze: the result flashes and Today is fetched again.
type mailActMsg struct{ text string }

// mailCmd runs crow-today on one mail (done, snooze).
func mailCmd(verb string, r todayRow) tea.Cmd {
	return func() tea.Msg {
		out, err := runBin("crow-today", verb, r.ID)
		if err != nil {
			return mailActMsg{verb + " failed: " + lastLine(cmpOr(out, err.Error()))}
		}
		return mailActMsg{fit(r.What, 40) + ": " + firstLine(cmpOr(out, "done"))}
	}
}

// mailTask is the hand-off for enter: a hub session reads the thread and drafts the reply.
func mailTask(r todayRow) (task, end string) {
	task = fmt.Sprintf("Mail from %s: %q (Graph message id %s). Read the thread with `bin/crow-today thread %s`, "+
		"say in a few lines what it needs, then save a reply draft in French with `bin/crow-today draft %s [--all] < text` "+
		"(sign-off per local/owner-prefs.md). Never send.", r.Tag, r.Subject, r.ID, r.ID, r.ID)
	return task, "A reply draft in Outlook Drafts and a short summary of the thread; nothing sent."
}

// todayKey runs a mail action on the selected Today row.
func (m model) todayKey(k string) (tea.Model, tea.Cmd) {
	r := m.homeView().curMail()
	if r == nil {
		m.flash = "no mail selected"
		return m, nil
	}
	if r.ID == "" {
		m.flash = k + ": this mail has no Graph id (mock data)"
		return m, nil
	}
	switch k {
	case "o":
		if r.Link == "" {
			m.flash = "no Outlook link for this mail"
			return m, nil
		}
		m.flash = "opening in Outlook…"
		return m, actionCmd("open", "crow-open", r.Link)
	case "x":
		m.flash = "marking done…"
		return m, mailCmd("done", *r)
	case "s":
		m.flash = "snoozing to tomorrow…"
		return m, mailCmd("snooze", *r)
	case "enter":
		task, end := mailTask(*r)
		m.flash = "opening a hub session on the mail…"
		return m, actionCmd("mail session", "crow", "handoff", "cr-0w", "--task", task, "--end-state", end)
	}
	m.flash = k + ": not a Today action (o open, enter session, x done, s snooze)"
	return m, nil
}
