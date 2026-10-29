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
}

// hspan is a stretch of the day in hours (8.5 = 08:30).
type hspan struct{ From, To float64 }

// whoRow is one line of "Who worked when": active stretches, then paused ones.
type whoRow struct {
	Name, Account string
	On, Off       []hspan
}

type blocker struct {
	Glyph string // needs, blocked, waiting
	Item  string
	Why   string
	Age   string
}

type todayRow struct {
	When, What, Tag string
	TagState        string // ok, warn, ""
}

type nextRow struct {
	Project string
	What    string
	Perso   bool
}

type projRow struct {
	Name       string
	Done, Of   int
	Date       string
	DateAtRisk bool
}

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
	Projects            []projRow
	Cards               []projCard // Projects tab
	Quiet               []string   // projects without a roadmap or activity
	Repos, Dirty, Ahead int
	Slots               [2]int // background runs busy, cap
	Who                 []whoRow
	Live                bool   // from the scripts, not invented
	Sync                string // sync state of live data
	SyncNote            string // stale age or failing source for the header dot
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
		Who:    mockWho(),
		Sync:   stFresh,
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
			{"11:30", "daily", "sent", "ok"},
			{"14:00", "design doc to PDF", "", ""},
			{"Fri", "terms reply", "draft", "warn"},
			{"Fri", "weekly slides", "", ""},
			{"Mon", "sprint planning", "", ""},
		},
		Mail: []string{"mail 3 unread", "1 flagged"},
		Next: []nextRow{
			{"atlas", "review !214, merge 214 then 213", false},
			{"atlas", "review !220 for Sam", false},
			{"atlas", "purge job for retention", false},
			{"vault", "design doc review, export PDF", false},
			{"vault", "terms reply, draft ready", false},
			{"crow", "audit: retros out of sessions", false},
			{"crow", "deck: real data on Home", false},
			{"atlas", "retention policy tests", false},
			{"vault", "weekly slides draft", false},
			{"crow", "command map and doctor", false},
			{"atlas", "alarm review after launch", false},
			{"vault", "decision log cleanup", false},
			{"lint", "panel layout", true},
			{"tools", "AC1 to AC4", true},
			{"web", "0.40 deck pill", true},
			{"lint", "release notes", true},
		},
		Projects: []projRow{
			{"atlas", 6, 8, "Dec 1", true},
			{"web", 4, 13, "", false},
			{"crow", 3, 9, "", false},
			{"vault", 5, 6, "", false},
		},
		Repos: 7, Dirty: 3, Ahead: 2,
		Usage: mockUsage(),
		MRs:   mockMRs(), Branches: mockBranches(), MRDetails: mockDetails(), MRsReady: true,
	}
	f.Cards, f.Quiet = mockCards(mockNow())
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

func mockWho() []whoRow {
	return []whoRow{
		{Name: "crow designs", Account: "main", On: []hspan{{9.5, 14.1}}},
		{Name: "atlas session", Account: "main", On: []hspan{{10, 11.4}, {13.2, 14.1}}},
		{Name: "vault", Account: "main", On: []hspan{{10.8, 12}, {13.6, 14.1}}},
		{Name: "notes", Account: "main", On: []hspan{{12, 12.5}}},
		{Name: "run 99a8", Account: "spare", On: []hspan{{11.4, 13.4}}, Off: []hspan{{13.4, 16.2}}},
	}
}

// emptyFixtures: no agent, no MR, no blocker.
func emptyFixtures() *fixtures {
	f := mockFixtures()
	f.Agents, f.Blockers, f.Today, f.Mail, f.Next, f.Projects = nil, nil, nil, nil, nil, nil
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
