package main

import "time"

// Mock data for the skeleton. All of it is invented and kept consistent across pages
// (counts, badges, run slots). Real sources arrive with the later slices.

type gauge struct {
	Account  string
	FiveH    int // % remaining
	Reset    string
	Week     int // % remaining, -1 none
	Policy   string
	Burn     string
	RunsUsed int
	RunsMax  int
	// Usage tab.
	WeekReset string // "Nov 2 17:00" or "Mon"
	Empty     string // clock the 5h window runs empty at the current burn
	BurnPct   float64
}

// Agent groups on the Agents tab, in display order.
const (
	gNeeds   = iota // waits for the owner
	gRunning        // interactive sessions at work
	gRuns           // background runs and paused runs
	gDone           // finished today
)

var groupNames = []string{"needs you", "running", "background runs", "finished today"}

type tailLine struct{ At, Text string }

type agentRow struct {
	State   string // working, idle, needs, waiting (paused), done
	Name    string // project, the Home column
	Account string
	Doing   string
	Since   string
	// Agents tab.
	Group       int
	Short       string // list label: session title, run label
	Meta        string // tail header: account, tokens, age
	Ask         string // the agent's last words when it waits for the owner
	Tail        []tailLine
	Project     string
	Title       string
	Ctx         int
	PID         int
	Session     string // session id, or run id for runs
	Cwd         string
	Interrupted bool
	Started     time.Time
	Quiet       time.Time // last activity
	Model       string    // short model name, from the usage session
}

type blocker struct {
	Glyph string // needs, blocked, waiting
	Item  string
	Why   string
	Age   string
}

type todayRow struct {
	When, What, Tag string
	TagState        string // ok, warn, next, past, ""
	Kind            string // "" event, head, flag, gap
	ID, Link        string // flag rows: the Graph message id and its Outlook web link
	Subject         string // flag rows: the subject as sent
}

type nextRow struct {
	Project string
	What    string
	Perso   bool
	Actor   string // me, crow, or the person a Waiting line names
	Since   string // YYYY-MM-DD the thread was opened
	Why     string // first ref in the line (item_047, !214, CD-580), else the next step
	Repo    string // full repo name for the Home table, "" uses Project
	Age     string // "today", "3d"; set when the fixtures are built
}

// projRow is one line of the Home Projects table.
type projRow struct {
	Name         string
	Scope        string
	Done, Of     int    // roadmap milestones, Of 0: no roadmap
	Now          string // current milestone, else what runs on it
	Date, Left   string // target "Dec 1" and "D-29", "" without a date
	DateAtRisk   bool
	Dirty, Ahead int
}

// doneRow is one line of today's activity log.
type doneRow struct{ When, Tag, What string }

type health struct {
	Identity, CDX bool
	StaleTrees    int
}

type fixtures struct {
	Now                 string // header date
	Synced              string
	Health              health
	Gauges              []gauge
	Agents              []agentRow
	Blockers            []blocker
	Today               []todayRow
	Mail                []string
	Next                []nextRow
	Decide              []nextRow // MRs that need the owner, on top of You decide
	Done                []doneRow // today's activity log, newest first
	Team                []gauge   // tier team accounts, one compact Pulse line
	Projects            []projRow
	Cards               []projCard // Projects tab
	Quiet               []string   // projects without a roadmap or activity
	Repos, Dirty, Ahead int
	Slots               [2]int              // background runs busy, cap
	Lanes               []laneRow           // Agents day lanes
	LaneGit             map[string]*laneGit // git state by checkout, fetched for the selected lane
	Hist                []histSession       // History, newest first
	HistTotal           int                 // sessions in the archive, 0: len(Hist)
	Live                bool                // from the scripts, not invented
	Sync                string              // sync state of live data
	SyncNote            string              // stale age or failing source for the header dot
	// Counts behind the tab badges and the Blockers title.
	MRsAwaiting, MRsMine, MRsThreads int
	// MRs tab.
	MRs       []mrRow
	Branches  []mrBranch
	MRDetails map[string]*mrDetail
	MRErrs    map[string]string // detail fetch errors by MR key
	MRsReady  bool              // the forges answered at least once
	ForgeErrs []string
	// Usage tab; nil until crow-usage answered.
	Usage *usageView
}

func mockFixtures() *fixtures {
	f := &fixtures{
		Now: "Thu 29 Oct 14:02", Synced: "14:01",
		Health: health{Identity: true, CDX: true, StaleTrees: 1},
		Gauges: []gauge{
			{Account: "main", FiveH: 60, Reset: "16:10", Week: 63, WeekReset: "Mon", Empty: "19:40", BurnPct: 8, Policy: "use", Burn: "burn 8%/h", RunsUsed: 0, RunsMax: 2},
			{Account: "spare", FiveH: 40, Reset: "19:30", Week: 31, WeekReset: "Sun", Empty: "21:10", BurnPct: 4, Policy: "ask-first", Burn: "burn 4%/h", RunsUsed: 1, RunsMax: 2},
		},
		Agents: mockAgents(),
		Slots:  [2]int{1, 2},
		Hist:   mockHist(mockNow()), HistTotal: 212,
		Sync: stFresh,
		Blockers: []blocker{
			{"needs", "!214 atlas", "PRJ-565 typed ids, CI running, blocks !213", "2h"},
			{"blocked", "!213 atlas", "PRJ-503 init prod, waits for !214", "prod"},
			{"needs", "!220 atlas", "Sam waits for your review", "1d"},
			{"waiting", "PRJ-571", "your decision on backups", "2d"},
			{"waiting", "Jo", "terms answers expected", "3d"},
			{"waiting", "PRJ-559", "sprint scope for the retention work", "3d"},
			{"waiting", "vault", "design doc review before PDF export", "4d"},
		},
		Today: []todayRow{
			{When: "09:30–10:00", What: "Sprint sync", TagState: "past"},
			{When: "11:30–11:45", What: "Team daily", TagState: "past"},
			{When: "16:00–16:30", What: "1:1 product", Tag: "in 8 min", TagState: "next"},
			{When: "17:30–18:00", What: "Demo dry run"},
			{Kind: "gap"}, {Kind: "head", What: "Mail to look at"},
			{Kind: "flag", What: "Terms v3 feedback", Tag: "Jo", TagState: "warn"},
			{Kind: "flag", What: "Shipping date", Tag: "Lee", TagState: "past"},
			{Kind: "gap"}, {Kind: "head", What: "Tomorrow"},
			{When: "09:00", What: "Travel to the demo site"},
		},
		Mail: []string{"mail 4 unread", "2 to look at"},
		Next: []nextRow{
			{Project: "atlas", What: "review !214, merge 214 then 213", Perso: false, Actor: "me", Since: "2026-10-29", Why: "!214"},
			{Project: "atlas", What: "purge job for retention", Perso: false, Actor: "me", Since: "2026-10-27", Why: "decide"},
			{Project: "vault", What: "design doc review, export PDF", Perso: false, Actor: "me", Since: "2026-10-29", Why: "doc"},
			{Project: "vault", What: "terms reply, draft ready", Perso: false, Actor: "me", Since: "2026-10-28", Why: "draft"},
			{Project: "lint", What: "release notes", Perso: true, Actor: "me", Since: "2026-10-20", Why: "release"},
			{Project: "tools", What: "AC1 to AC4", Perso: true, Actor: "me", Since: "2026-10-22", Why: "manual"},
			{Project: "crow", What: "audit: retros out of sessions", Perso: true, Actor: "crow", Since: "2026-10-29", Why: "item_047"},
			{Project: "crow", What: "deck: real data on Home", Perso: true, Actor: "crow", Since: "2026-10-28", Why: "req_011"},
			{Project: "atlas", What: "retention policy tests", Perso: false, Actor: "crow", Since: "2026-10-26", Why: "session"},
			{Project: "web", What: "0.40 deck pill", Perso: true, Actor: "crow", Since: "2026-10-19", Why: "item"},
			{Project: "atlas", What: "review !220 after the rebase", Perso: false, Actor: "Sam", Since: "2026-10-29", Why: "!220"},
			{Project: "vault", What: "answers on legal notices", Perso: false, Actor: "Lee", Since: "2026-10-18", Why: "mail"},
		},
		Decide: []nextRow{
			{Project: "atlas", What: "!214 pipeline failed", Actor: "me", Why: "MR", Age: "2h"},
			{Project: "atlas", What: "!220 Sam waits for your review", Actor: "me", Why: "MR", Age: "1d"},
		},
		Team: []gauge{
			{Account: "team-a", FiveH: 20, Week: 68, Policy: "never"},
			{Account: "team-b", FiveH: 100, Week: 2, Policy: "never"},
		},
		Projects: []projRow{
			{Name: "atlas", Scope: "work", Done: 6, Of: 8, Now: "init prod seed", Date: "Dec 1", Left: "D-33", DateAtRisk: true, Dirty: 2, Ahead: 1},
			{Name: "web", Scope: "personal", Done: 4, Of: 13, Now: "design system v2"},
			{Name: "crow", Scope: "personal", Done: 3, Of: 9, Now: "deck rework", Dirty: 4, Ahead: 3},
			{Name: "vault", Scope: "work"},
			{Name: "lint", Scope: "personal", Done: 1, Of: 8, Now: "current flow correctness", Date: "Nov 20", Left: "D-22"},
			{Name: "tools", Scope: "personal", Done: 5, Of: 6, Now: "acceptance checks"},
		},
		Done: []doneRow{
			{"13:58", "atlas", "!214 rebased on main, pipeline restarted"},
			{"13:40", "crow", "done: deck Next rows carry actor and age"},
			{"12:15", "vault", "design doc exported to PDF for review"},
			{"11:30", "hub", "daily posted"},
			{"10:05", "web", "0.40 tokens merged into the deck pill"},
		},
		Repos: 7, Dirty: 3, Ahead: 2,
		Usage: mockUsage(),
		MRs:   mockMRs(), Branches: mockBranches(), MRDetails: mockDetails(), MRsReady: true,
	}
	f.Cards, f.Quiet = mockCards(mockNow())
	mockLaneSetup(f, mockNow())
	for i := range f.Next {
		f.Next[i].Age = threadAge(f.Next[i].Since, time.Date(2026, 10, 29, 14, 2, 0, 0, time.Local))
	}
	mockMRCounts(f)
	return f
}

func mockAgents() []agentRow {
	ask := "The commit subject has no file scope. Keep this finding?"
	rows := []agentRow{
		{State: "needs", Name: "atlas", Account: "main", Doing: "asks you: finding 2, the commit subject has no file scope", Since: "2m", Group: gNeeds,
			Short: "review !220", Meta: "main · medium · 31k tokens · 2m", Ask: ask,
			Tail: []tailLine{
				{"14:04", "$ crow review !220 --patterns atlas"}, {"14:04", "loaded 14 ledger patterns, 3 touch retry.py"},
				{"14:05", "reading retry.py (84 lines changed)"}, {"14:05", "◆ BACKOFF comes from the environment, not in config.md"},
				{"14:06", "checking config.md"}}},
		{State: "needs", Name: "atlas", Account: "main", Doing: "asks you: undraft the MR?", Since: "1m", Group: gNeeds,
			Short: "atlas undraft?", Meta: "main · 22k tokens · 1m", Ask: "Undraft !213 now that CI is green?",
			Tail: []tailLine{{"14:05", "CI passed on !213"}, {"14:05", "waiting for your answer"}}},
		{State: "working", Name: "crow", Account: "main", Doing: "writing designs, Home landscape in the canvas", Since: "4m", Group: gRunning,
			Short: "crow designs", Meta: "main · 64k tokens · 4m", Tail: []tailLine{{"14:02", "Edit: HomeL_merge.dc.html"}, {"14:05", "Bash: crow-compare"}}},
		{State: "working", Name: "vault", Account: "main", Doing: "designing the CR-0W data layout", Since: "12m", Group: gRunning,
			Short: "vault data layout", Meta: "main · 41k tokens · 12m", Tail: []tailLine{{"13:58", "Read: layout.md"}, {"14:03", "Write: layout.md"}}},
		{State: "waiting", Name: "run 99a8", Account: "spare", Doing: "paused on the 5h limit, back 16:10, R relaunch", Since: "9m", Group: gRuns,
			Short: "99a8 paused", Meta: "spare · 5h limit", Interrupted: true,
			Tail: []tailLine{{"13:53", "stopped on the 5h limit"}, {"", "back 16:10 (R relaunches it on another account)"}}},
		{State: "done", Name: "notes", Account: "main", Doing: "finished: 0.40 pill, 3 commits, unpushed", Since: "12:30", Group: gDone,
			Short: "notes 3 commits", Meta: "done · 18m", Tail: []tailLine{{"12:30", "finished: 0.40 pill, 3 commits, unpushed"}}},
		{State: "done", Name: "tools", Account: "main", Doing: "finished: AC2 tests green, 1 commit", Since: "11:05", Group: gDone,
			Short: "tools AC2", Meta: "done · 9m", Tail: []tailLine{{"11:05", "finished: AC2 tests green, 1 commit"}}},
		{State: "done", Name: "lint", Account: "main", Doing: "finished: panel layout notes", Since: "09:40", Group: gDone,
			Short: "lint notes", Meta: "done · 6m", Tail: []tailLine{{"09:40", "finished: panel layout notes"}}},
	}
	return rows
}

// emptyFixtures: no agent, no MR, no blocker.
func emptyFixtures() *fixtures {
	f := mockFixtures()
	f.Agents, f.Blockers, f.Today, f.Mail, f.Next, f.Projects = nil, nil, nil, nil, nil, nil
	f.Lanes, f.LaneGit = nil, nil
	f.Decide, f.Done = nil, nil
	f.Hist, f.HistTotal = nil, 0
	f.MRs, f.Branches, f.MRDetails = nil, nil, nil
	f.MRsAwaiting, f.MRsMine, f.MRsThreads = 0, 0, 0
	f.Usage = nil
	f.Cards, f.Quiet = nil, nil
	f.Gauges[1].RunsUsed = 0
	return f
}

// agentsNeedYou is the badge on the Agents tab: agents waiting for the owner.
func (f *fixtures) agentsNeedYou() int {
	n := 0
	for _, a := range f.Agents {
		if a.Group == gNeeds {
			n++
		}
	}
	return n
}

// mrsNeedYou is the badge on the MRs tab: MRs awaiting a review plus threads to fix.
func (f *fixtures) mrsNeedYou() int { return f.MRsAwaiting + f.MRsThreads }
