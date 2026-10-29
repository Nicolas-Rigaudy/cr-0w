package main

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
}

type agentRow struct {
	State   string // working, needs, waiting, done
	Name    string
	Account string
	Doing   string
	Since   string
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
	Repos, Dirty, Ahead int
	// Counts behind the tab badges and the Blockers title.
	MRsAwaiting, MRsMine, MRsThreads int
}

func mockFixtures() *fixtures {
	return &fixtures{
		Now: "Thu 29 Oct 14:02", Synced: "14:01",
		Health: health{Identity: true, CDX: true, StaleTrees: 1},
		Gauges: []gauge{
			{"main", 60, "16:10", 63, "", "", 0, 0},
			{"spare", 40, "", -1, "ask-first", "burn 4%/h", 1, 2},
		},
		Agents: []agentRow{
			{"working", "crow", "main", "writing designs, Home landscape in the canvas", "4m"},
			{"needs", "atlas", "main", "MR review session, asks you to undraft", "1m"},
			{"working", "vault", "main", "designing the CR-0W data layout", "12m"},
			{"waiting", "run 99a8", "spare", "asks before using spare (ask-first), R relaunch", "9m"},
			{"done", "notes", "main", "finished: 0.40 pill, 3 commits, unpushed", "12:30"},
			{"done", "tools", "main", "finished: AC2 tests green, 1 commit", "11:05"},
			{"done", "lint", "main", "finished: panel layout notes", "09:40"},
			{"done", "run 77c1", "spare", "finished: retro batch, 4 sessions", "08:15"},
		},
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
		MRsAwaiting: 1, MRsMine: 2, MRsThreads: 1,
	}
}

// emptyFixtures: no agent, no MR, no blocker.
func emptyFixtures() *fixtures {
	f := mockFixtures()
	f.Agents, f.Blockers, f.Today, f.Mail, f.Next, f.Projects = nil, nil, nil, nil, nil, nil
	f.MRsAwaiting, f.MRsMine, f.MRsThreads = 0, 0, 0
	f.Gauges[1].RunsUsed = 0
	return f
}

// needsYou is the badge on the Agents tab: agents waiting for the owner.
func (f *fixtures) agentsNeedYou() int {
	n := 0
	for _, a := range f.Agents {
		if a.State == "needs" || a.State == "waiting" {
			n++
		}
	}
	return n
}

// mrsNeedYou is the badge on the MRs tab: MRs awaiting a review plus threads to fix.
func (f *fixtures) mrsNeedYou() int { return f.MRsAwaiting + f.MRsThreads }
