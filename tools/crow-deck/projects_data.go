package main

// Projects tab data: one card per project with a roadmap or some activity, built from the
// registry (crow-registry, incl. milestone progress), the agents, the open-threads rows and
// the MRs. Roadmap docs carry milestones but no target dates, so live cards have no date.

import (
	"strconv"
	"strings"
	"time"
)

type projCard struct {
	Name, Release string // Release: "release 1.0" or ""
	Scope         string
	State         string // needs you, blocked, running, ready
	Now           string // what runs on it now, empty when idle
	NowSince      string
	Next, Block   string
	Done, Of      int // milestones
	Target        time.Time
	AtRisk        bool
	Current, Exit string // first milestone not done and its exit signal
	Open          string // "2 MRs · 1 thread"
}

// roadText is the ROAD line: "6/8 · target Dec 1 · at risk".
func (c projCard) roadText() string {
	if c.Of == 0 {
		return ""
	}
	s := itoa(c.Done) + "/" + itoa(c.Of)
	switch {
	case c.Target.IsZero():
		return s + " · no date"
	case c.AtRisk:
		return s + " · target " + c.Target.Format("Jan 2") + " · at risk"
	}
	return s + " · target " + c.Target.Format("Jan 2") + " · on time"
}

func itoa(n int) string { return strconv.Itoa(n) }

// projectCardsOf returns the cards, then the names of the projects without a roadmap and
// without activity (listed on one line).
func projectCardsOf(now time.Time, reg *registry, agents []agentRow, next []nextRow, mrs []mrRow) (cards []projCard, quiet []string) {
	if reg == nil {
		return nil, nil
	}
	for _, p := range reg.Projects {
		c := projCard{Name: p.Short, Scope: p.Scope}
		if p.Release.Target != "" {
			c.Release = "release " + p.Release.Target
		}
		if r := p.Road; r != nil {
			c.Done, c.Of, c.Current, c.Exit = r.Done, r.Total, r.Current, r.Exit
		}
		running := false
		for _, a := range agents {
			if a.Name != p.Short && a.Project != p.Short {
				continue
			}
			switch a.Group {
			case gNeeds:
				c.State = "needs you"
				c.Now = cmpOr(a.Short, a.Doing)
			case gRunning:
				running = true
				c.Now = cmpOr(a.Short, a.Doing)
			}
			c.NowSince = a.Since
		}
		for _, n := range next {
			if n.Project == p.Short {
				c.Next = n.What
				break
			}
		}
		var mine []mrRow
		for _, r := range mrs {
			if r.Project == p.Short {
				mine = append(mine, r)
			}
		}
		var blocks []string
		for _, b := range mrBlockers(now, mine) {
			blocks = append(blocks, strings.TrimSuffix(b.Item, " "+p.Short)+" "+b.Why)
			if b.Glyph == "blocked" && c.State == "" {
				c.State = "blocked"
			}
		}
		c.Block = strings.Join(blocks, " · ")
		_, _, th := mrCounts(mine)
		c.Open = openText(len(mine), th)
		if c.State == "" {
			c.State = "ready"
			if running {
				c.State = "running"
			}
		}
		if c.Of == 0 && c.Now == "" && c.Next == "" && c.Block == "" && len(mine) == 0 {
			quiet = append(quiet, p.Short)
			continue
		}
		cards = append(cards, c)
	}
	return cards, quiet
}

func openText(mrs, threads int) string {
	var parts []string
	if mrs > 0 {
		parts = append(parts, itoa(mrs)+" "+plural(mrs, "MR", "MRs"))
	}
	if threads > 0 {
		parts = append(parts, itoa(threads)+" "+plural(threads, "thread", "threads"))
	}
	return strings.Join(parts, " · ")
}

// projectRows feeds the Home Projects box: the cards that have a roadmap.
func projectRows(cards []projCard) []projRow {
	var out []projRow
	for _, c := range cards {
		if c.Of > 0 {
			d := ""
			if !c.Target.IsZero() {
				d = c.Target.Format("Jan 2")
			}
			out = append(out, projRow{c.Name, c.Done, c.Of, d, c.AtRisk})
		}
	}
	return out
}

func mockCards(now time.Time) ([]projCard, []string) {
	y := now.Year()
	day := func(m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, now.Location()) }
	return []projCard{
		{Name: "atlas", Release: "release 1.0", Scope: "work", State: "blocked", Now: "review session on !214, CI running", NowSince: "1m",
			Next: "review !214, then merge 214 then 213", Block: "!213 waits for !214 · !220 waits for your review",
			Done: 6, Of: 8, Target: day(time.December, 1), Current: "init prod seed", Exit: "seed on a clean prod account", Open: "3 MRs · 1 thread"},
		{Name: "web", Release: "release 1.0", Scope: "personal", State: "ready", Now: "", Next: "push, then 0.40 deck pill",
			Done: 4, Of: 13, Current: "design system v2", Exit: "every page uses the tokens"},
		{Name: "lint", Release: "release 1.0", Scope: "personal", State: "needs you", Next: "panel layout, waits for your decision",
			Done: 1, Of: 8, Target: day(time.November, 20), AtRisk: true, Current: "current flow correctness", Exit: "no known correctness bug is open"},
		{Name: "vault", Release: "release 1.0", Scope: "personal", State: "ready", Next: "acceptance checks AC1 to AC4",
			Done: 5, Of: 6, Current: "books flow", Exit: "the checks in the brief pass"},
		{Name: "crow", Scope: "personal", State: "running", Now: "writing designs", NowSince: "4m", Next: "P1 actions menu",
			Done: 3, Of: 9, Current: "P1 actions", Exit: "every key does what its hint says"},
	}, []string{"notes", "dots", "cdx", "tools"}
}
