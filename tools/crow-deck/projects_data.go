package main

// Projects tab data: one card per active project, built from the registry (roadmap, dated
// release, git state, the owner's commits per day), the agents, the open-threads rows and
// the MRs. Dates come only from the roadmap docs ("- Target:" lines) and release.date.

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// Project groups, in display order.
const (
	pgNeeds  = iota // the owner acts
	pgActive        // running, waiting on someone else, or ready to pick
	pgQuiet         // nothing open
)

var projGroupNames = []string{"needs you", "running or ready", "quiet"}

type projMile struct {
	Title string
	Date  time.Time
	Done  bool
}

// projMR is one open MR on the detail pane: its state and who must act.
type projMR struct {
	Ref, State, Acts, Todo string
	Bad                    bool // pipeline failed
}

type projCard struct {
	Name, Short   string // repo name, and the alias the threads and agents use
	Scope, Forge  string
	Group         int
	Glyph         string // blocked, needs, running, waiting, ready, quiet
	Health        string // late, at risk, on time, no date, idle
	Done, Of      int    // roadmap milestones
	Start, Target time.Time
	Mile, Exit    string // current milestone (or the release label) and its exit signal
	Miles         []projMile
	Now, NowSince string // the agent at work on it
	Block         string // MR blockers
	Next          []nextRow
	MRs           []projMR
	Dirty, Ahead  int
	Stale         int
	Work          map[string]int // owner commits per day, YYYY-MM-DD
	Doc           string         // the roadmap doc, absolute
	Path          string         // the repo, absolute
	Logics        bool           // it has a Logics corpus
}

func itoa(n int) string { return strconv.Itoa(n) }

func day(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func parseDay(s string, loc *time.Location) time.Time {
	t, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		return time.Time{}
	}
	return t
}

// left is "D-29", "2d late" or "" without a target.
func (c projCard) left(now time.Time) string {
	if c.Target.IsZero() {
		return ""
	}
	n := int(c.Target.Sub(day(now)).Hours() / 24)
	if n < 0 {
		return itoa(-n) + "d late"
	}
	return "D-" + itoa(n)
}

// lead is the next action the list shows: the owner's first, else the first open one.
func (c projCard) lead() *nextRow {
	for i := range c.Next {
		if c.Next[i].Actor == "me" {
			return &c.Next[i]
		}
	}
	if len(c.Next) > 0 {
		return &c.Next[0]
	}
	return nil
}

// health: late past the target; at risk when the milestones done trail the time spent
// since the roadmap start, or when a dated release without a roadmap has a blocked MR
// within 30 days; on time otherwise; no date without a target.
func cardHealth(now time.Time, c projCard) string {
	if c.Target.IsZero() {
		return "no date"
	}
	today := day(now)
	if c.Target.Before(today) {
		return "late"
	}
	if c.Of > 0 && !c.Start.IsZero() && c.Target.After(c.Start) {
		spent := today.Sub(c.Start).Hours() / c.Target.Sub(c.Start).Hours()
		if float64(c.Done)/float64(c.Of)+0.1 < spent {
			return "at risk"
		}
		return "on time"
	}
	if c.Glyph == "blocked" && c.Target.Sub(today) <= 30*24*time.Hour {
		return "at risk"
	}
	return "on time"
}

// mrOf says what an open MR waits for and who must act.
func mrOf(r mrRow) projMR {
	var st []string
	if r.Draft {
		st = append(st, "draft")
	}
	switch r.Pipeline {
	case "failed":
		st = append(st, "pipeline failed")
	case "success":
		st = append(st, "pipeline passed")
	}
	if r.ThreadsOpen > 0 {
		st = append(st, itoa(r.ThreadsOpen)+" "+plural(r.ThreadsOpen, "thread", "threads"))
	}
	m := projMR{Ref: r.Ref, State: cmpOr(strings.Join(st, ", "), "open"), Bad: r.Pipeline == "failed", Acts: "you"}
	switch {
	case r.Role == "REVIEWER":
		m.Todo = "review " + cmpOr(r.Author, "the MR")
	case r.Pipeline == "failed":
		m.Todo = "fix the pipeline"
	case r.ThreadsOpen > 0:
		m.Todo = "answer the threads"
	case r.Draft:
		m.Todo = "undraft when ready"
	default:
		m.Acts, m.Todo = "reviewers", "review and merge"
	}
	return m
}

// projectCardsOf returns a card per active project, sorted by group, health and target.
func projectCardsOf(now time.Time, reg *registry, agents []agentRow, next []nextRow, mrs []mrRow) []projCard {
	if reg == nil {
		return nil
	}
	loc := now.Location()
	var cards []projCard
	for _, p := range reg.Projects {
		if p.Status != "" && p.Status != "active" {
			continue
		}
		c := projCard{Name: p.Name, Short: p.Short, Scope: p.Scope, Forge: p.Forge,
			Dirty: p.Git.Dirty, Ahead: p.Git.Ahead, Stale: p.Git.StaleTrees, Work: p.Work, Path: p.Path, Logics: p.Logics}
		if r := p.Road; r != nil {
			c.Done, c.Of, c.Mile, c.Exit, c.Doc = r.Done, r.Total, r.Current, r.Exit, r.Doc
			c.Start, c.Target = parseDay(r.Start, loc), parseDay(r.Target, loc)
			for _, m := range r.Miles {
				if d := parseDay(m.Date, loc); !d.IsZero() {
					c.Miles = append(c.Miles, projMile{m.Title, d, m.Done})
				}
			}
		}
		if d := parseDay(p.Release.Date, loc); !d.IsZero() {
			if c.Target.IsZero() {
				c.Target, c.Mile = d, cmpOr(p.Release.Label, "Release")
			}
			c.Miles = append(c.Miles, projMile{Title: cmpOr(p.Release.Label, "Release"), Date: d})
		}
		mine := func(n string) bool { return n != "" && (n == p.Short || n == p.Name) }
		needs, running := false, false
		for _, a := range agents {
			if !mine(a.Name) && !mine(a.Project) {
				continue
			}
			switch a.Group {
			case gNeeds:
				needs = true
				c.Now, c.NowSince = cmpOr(a.Short, a.Doing), a.Since
			case gRunning:
				running = true
				c.Now, c.NowSince = cmpOr(a.Short, a.Doing), a.Since
			}
		}
		for _, n := range next {
			if mine(n.Project) {
				c.Next = append(c.Next, n)
				needs = needs || n.Actor == "me"
			}
		}
		var rows []mrRow
		for _, r := range mrs {
			if mine(r.Project) {
				rows = append(rows, r)
				m := mrOf(r)
				c.MRs = append(c.MRs, m)
				needs = needs || m.Acts == "you"
			}
		}
		var blocks []string
		blocked := false
		for _, b := range mrBlockers(now, rows) {
			blocks = append(blocks, strings.TrimSuffix(b.Item, " "+r0(rows).Project)+" "+b.Why)
			blocked = blocked || b.Glyph == "blocked"
		}
		c.Block = strings.Join(blocks, " · ")
		switch {
		case blocked:
			c.Group, c.Glyph = pgNeeds, "blocked"
		case needs:
			c.Group, c.Glyph = pgNeeds, "needs"
		case running:
			c.Group, c.Glyph = pgActive, "running"
		case len(c.Next) > 0 || len(c.MRs) > 0:
			c.Group, c.Glyph = pgActive, "waiting"
		case c.Of > 0 || !c.Target.IsZero():
			c.Group, c.Glyph = pgActive, "ready"
		default:
			c.Group, c.Glyph = pgQuiet, "quiet"
		}
		if c.Group == pgQuiet {
			c.Health = "idle"
		} else {
			c.Health = cardHealth(now, c)
		}
		cards = append(cards, c)
	}
	sortCards(cards)
	return cards
}

func r0(rows []mrRow) mrRow {
	if len(rows) == 0 {
		return mrRow{}
	}
	return rows[0]
}

var healthRank = map[string]int{"late": 0, "at risk": 1, "on time": 2, "no date": 3, "idle": 4}

func sortCards(cards []projCard) {
	sort.SliceStable(cards, func(i, j int) bool {
		a, b := cards[i], cards[j]
		if a.Group != b.Group {
			return a.Group < b.Group
		}
		if healthRank[a.Health] != healthRank[b.Health] {
			return healthRank[a.Health] < healthRank[b.Health]
		}
		if a.Target.IsZero() != b.Target.IsZero() {
			return !a.Target.IsZero()
		}
		return a.Target.Before(b.Target)
	})
}

// mockCards: invented projects consistent with the mock agents, threads and MRs.
func mockCards(now time.Time) []projCard {
	d := func(n int) time.Time { return day(now).AddDate(0, 0, n) }
	work := func(seed int, days ...int) map[string]int {
		out := map[string]int{}
		for i, n := range days {
			out[d(-n).Format("2006-01-02")] = 1 + (seed*7+i*5)%17
		}
		return out
	}
	me := func(p, what, why string) nextRow {
		return nextRow{Project: p, What: what, Actor: "me", Why: why, Age: "today"}
	}
	cards := []projCard{
		{Name: "atlas-platform", Short: "atlas", Scope: "work", Forge: "gitlab", Glyph: "blocked", Group: pgNeeds,
			Mile: "Prod launch", Target: d(29), Miles: []projMile{{"Prod launch", d(29), false}},
			Now: "review session on !214, CI running", NowSince: "1m", Block: "!214 pipeline failed · !220 sam waits for your review",
			Next: []nextRow{me("atlas", "review !214, then merge 214 then 213", "!214"), {Project: "atlas", What: "seed the demo data", Actor: "crow", Why: "item_012", Age: "2d"}},
			MRs: []projMR{{Ref: "!214", State: "pipeline failed", Acts: "you", Todo: "fix the pipeline", Bad: true},
				{Ref: "!213", State: "draft, pipeline passed", Acts: "you", Todo: "undraft when ready"},
				{Ref: "!220", State: "open", Acts: "you", Todo: "review sam"}},
			Ahead: 1, Stale: 2, Work: work(1, 0, 1, 2, 5, 6, 7, 8, 9, 12, 14, 15, 16, 20, 21, 22, 27, 29, 34)},
		{Name: "lint", Short: "lint", Scope: "personal", Forge: "github", Glyph: "needs", Group: pgNeeds,
			Done: 1, Of: 8, Start: d(-20), Target: d(-2), Mile: "0.2 - Current flow correctness", Exit: "no known correctness bug is open",
			Miles: []projMile{{"0.2 - Current flow correctness", d(-2), false}, {"1.0 - Public release", d(60), false}},
			Next:  []nextRow{me("lint", "panel layout, waits for your decision", "item_031")}, Work: work(2, 3, 4, 6, 10, 11)},
		{Name: "web", Short: "web", Scope: "personal", Forge: "github", Glyph: "running", Group: pgActive,
			Done: 4, Of: 13, Start: d(-70), Target: d(34), Mile: "0.40 - Design system v2", Exit: "every page uses the tokens",
			Miles: []projMile{{"0.40 - Design system v2", d(34), false}},
			Now:   "writing designs", NowSince: "4m", Next: []nextRow{{Project: "web", What: "push, then 0.40 deck pill", Actor: "crow", Why: "item_040", Age: "1d"}},
			Dirty: 3, Work: work(3, 0, 2, 3, 7, 8, 9, 13, 14, 22, 23, 28, 29, 30)},
		{Name: "vault-plugin", Short: "vault", Scope: "personal", Forge: "github", Glyph: "waiting", Group: pgActive,
			Done: 5, Of: 6, Start: d(-30), Mile: "Books flow", Exit: "the checks in the brief pass",
			Next: []nextRow{{Project: "vault", What: "release workflow for the v0.9 tag", Actor: "CI", Why: "CI run", Age: "1d"}},
			Work: work(4, 1, 8, 9)},
		{Name: "crow", Short: "crow", Scope: "personal", Forge: "none", Glyph: "ready", Group: pgActive,
			Done: 3, Of: 9, Start: d(-21), Mile: "P1 actions", Exit: "every key does what its hint says",
			Work: work(5, 0, 1, 2, 3, 6, 7, 13, 14)},
		{Name: "notes", Short: "notes", Scope: "personal", Glyph: "quiet", Group: pgQuiet},
		{Name: "dots", Short: "dots", Scope: "personal", Glyph: "quiet", Group: pgQuiet, Work: work(6, 40)},
	}
	for i := range cards {
		if cards[i].Group == pgQuiet {
			cards[i].Health = "idle"
		} else {
			cards[i].Health = cardHealth(now, cards[i])
		}
	}
	sortCards(cards)
	return cards
}

// homeProjects is the Home Projects table: every active project in registry order, with the
// roadmap progress and target from its card and the git state from the registry.
func homeProjects(now time.Time, reg *registry, cards []projCard) []projRow {
	if reg == nil {
		return nil
	}
	byName := map[string]projCard{}
	for _, c := range cards {
		byName[c.Name] = c
	}
	var out []projRow
	for _, p := range reg.Projects {
		if p.Status != "" && p.Status != "active" {
			continue
		}
		r := projRow{Name: p.Name, Scope: p.Scope, Dirty: p.Git.Dirty, Ahead: p.Git.Ahead}
		if c, ok := byName[p.Name]; ok {
			r.Done, r.Of = c.Done, c.Of
			r.DateAtRisk = c.Health == "late" || c.Health == "at risk"
			r.Now = cmpOr(c.Mile, c.Now)
			if !c.Target.IsZero() {
				r.Date = c.Target.Format("Jan 2")
				r.Left = c.left(now)
			}
		}
		out = append(out, r)
	}
	return out
}
