package main

// Day lanes: one row per session that worked today, from crow-usage's 10-minute activity
// buckets, joined to the live agents from crow-runs. Background runs with the same title
// in the same project share one lane.

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Lane states, in the glyph legend's terms.
const (
	lsNeeds   = "needs"   // waits for the owner's answer
	lsWorking = "working" // interactive, busy
	lsIdle    = "idle"    // interactive, waits for the next prompt
	lsRun     = "run"     // background run at work
	lsPaused  = "paused"  // stopped on a usage limit
	lsDone    = "done"    // ended today
)

// laneRow is one session (or a set of same-titled runs) across today.
type laneRow struct {
	Key         string // session id, or the agent's session/run id
	Project     string
	Title       string
	State       string
	Account     string
	Model       string // "opus-5-5"
	Branch      string
	Cwd         string
	Agent       int // index in fixtures.Agents, -1 none
	Run         bool
	Runs        int // runs folded into this lane
	Ctx         int // context tokens, live sessions only
	Cost        float64
	Priced      bool
	Cells       map[int]bool // 10-minute buckets with activity
	Tok         map[int]int  // tokens per bucket
	Tokens      int
	CachePct    int
	Compactions int
	Files       map[string]int // edited file basename → edits
	Start, End  time.Time
	WaitFrom    time.Time // needs: when the agent started waiting
	Next        string
	NextYou     bool
}

func (l laneRow) first() int {
	f := -1
	for b := range l.Cells {
		if f < 0 || b < f {
			f = b
		}
	}
	return f
}

func (l laneRow) last() int {
	f := -1
	for b := range l.Cells {
		if b > f {
			f = b
		}
	}
	return f
}

// mins is the active time today: ten minutes per bucket.
func (l laneRow) mins() int { return 10 * len(l.Cells) }

func (l laneRow) live() bool {
	switch l.State {
	case lsNeeds, lsWorking, lsIdle, lsRun:
		return true
	}
	return false
}

// shortModel picks the session's main model and drops the family prefix.
func shortModel(models []string) string {
	best, rank := "", -1
	for _, m := range models {
		r := 0
		for i, fam := range []string{"haiku", "sonnet", "opus", "fable"} {
			if strings.Contains(m, fam) {
				r = i + 1
			}
		}
		if r > rank {
			best, rank = m, r
		}
	}
	best = strings.TrimPrefix(best, "claude-")
	if i := strings.Index(best, "["); i > 0 {
		best = best[:i]
	}
	return best
}

func bucketTime(day time.Time, b int) time.Time {
	y, mo, d := day.Local().Date()
	return time.Date(y, mo, d, 0, b*10, 0, 0, day.Location())
}

// lanesOf builds today's lanes. Live agents without activity yet still get a lane.
func lanesOf(now time.Time, sessions []usageSession, agents []agentRow, slots [2]int, reg *registry) []laneRow {
	day := now.Local().Format("2006-01-02")
	bySess := map[string]int{}
	for i, a := range agents {
		if a.Session != "" && a.Group != gDone {
			bySess[a.Session] = i
		}
	}
	used := map[int]bool{}
	var lanes []laneRow
	runLane := map[string]int{}
	for _, s := range sessions {
		b := s.Activity[day]
		ai, isLive := bySess[s.ID]
		if len(b) == 0 && !isLive {
			continue
		}
		l := laneRow{Key: s.ID, Project: reg.shortName(s.Project), Title: whoName(s, reg), Account: s.Account,
			Model: shortModel(s.Models), Branch: s.Branch, Agent: -1, Run: s.Kind == "run", Runs: 1,
			Cost: s.Today.CostUSD, Priced: true, Cells: map[int]bool{}, Tok: map[int]int{},
			Tokens: s.Tokens, CachePct: s.CachePct, Compactions: s.Compactions, Files: s.FilesEdited}
		l.Start, _ = parseTime(s.Start)
		l.End, _ = parseTime(s.End)
		for _, x := range b {
			l.Cells[x] = true
		}
		for k, v := range s.ActivityTokens[day] {
			if n, err := strconv.Atoi(k); err == nil {
				l.Tok[n] += v
			}
		}
		if isLive {
			l.Agent, used[ai] = ai, true
		}
		if l.Run {
			key := l.Project + "\x00" + l.Title
			if j, ok := runLane[key]; ok {
				lanes[j] = mergeRun(lanes[j], l)
				continue
			}
			runLane[key] = len(lanes)
		}
		lanes = append(lanes, l)
	}
	for i, a := range agents {
		if used[i] || a.Group == gDone {
			continue
		}
		lanes = append(lanes, laneRow{Key: a.Session, Project: cmpOr(a.Project, a.Name), Title: cmpOr(a.Title, a.Short),
			Account: a.Account, Agent: i, Run: a.Group == gRuns && !a.Interrupted, Runs: 1, Ctx: a.Ctx,
			Cells: map[int]bool{}, Tok: map[int]int{}, Start: a.Started})
	}
	for i := range lanes {
		l := &lanes[i]
		if l.Agent >= 0 {
			a := agents[l.Agent]
			l.Ctx = a.Ctx
			l.Cwd = a.Cwd
			l.Project = cmpOr(a.Project, l.Project)
			if a.Model == "" {
				agents[l.Agent].Model = l.Model
			}
		}
		if l.Cwd == "" {
			if p, ok := reg.find(l.Project); ok {
				l.Cwd = p.Path
			}
		}
		l.State, l.Next, l.NextYou = laneState(*l, agents, slots)
		if l.State == lsNeeds && l.Agent >= 0 {
			l.WaitFrom = agents[l.Agent].Quiet
		}
	}
	return lanes
}

func mergeRun(a, b laneRow) laneRow {
	for k := range b.Cells {
		a.Cells[k] = true
	}
	for k, v := range b.Tok {
		a.Tok[k] += v
	}
	a.Runs++
	a.Cost += b.Cost
	a.Tokens += b.Tokens
	if b.Agent >= 0 {
		a.Agent, a.Key = b.Agent, b.Key
	}
	if b.End.After(a.End) {
		a.End, a.Branch = b.End, cmpOr(b.Branch, a.Branch)
	}
	if a.Start.IsZero() || (!b.Start.IsZero() && b.Start.Before(a.Start)) {
		a.Start = b.Start
	}
	return a
}

// laneState is the lane's glyph state and its "next" column: who acts next.
func laneState(l laneRow, agents []agentRow, slots [2]int) (state, next string, you bool) {
	if l.Agent < 0 {
		return lsDone, "nobody", false
	}
	a := agents[l.Agent]
	switch {
	case a.Interrupted:
		return lsPaused, cmpOr(strings.TrimPrefix(backAt(a.Doing), "back "), "paused"), false
	case a.Group == gNeeds:
		return lsNeeds, "answer", true
	case a.Group == gRuns:
		return lsRun, fmt.Sprintf("background %d/%d", slots[0], max(slots[1], 1)), false
	case a.State == "idle":
		return lsIdle, "next prompt", true
	}
	return lsWorking, "working", false
}

// backAt finds "back HH:MM" in a paused row's text.
func backAt(s string) string {
	if i := strings.Index(s, "back "); i >= 0 {
		rest := s[i:]
		if j := strings.Index(rest, ","); j > 0 {
			rest = rest[:j]
		}
		return rest
	}
	return ""
}

// laneGroup is one project's lanes, in display order.
type laneGroup struct {
	Project       string
	Lanes         []int // indexes in fixtures.Lanes
	Need, Running int
	LastAt        int
}

// groupLanes orders projects: most lanes that need the owner first, then running ones,
// then the latest activity. Inside a project, lanes run by start.
func groupLanes(lanes []laneRow) []laneGroup {
	idx := map[string]int{}
	var gs []laneGroup
	for i, l := range lanes {
		j, ok := idx[l.Project]
		if !ok {
			j = len(gs)
			idx[l.Project] = j
			gs = append(gs, laneGroup{Project: l.Project, LastAt: -1})
		}
		g := &gs[j]
		g.Lanes = append(g.Lanes, i)
		if l.NextYou {
			g.Need++
		}
		if l.State == lsWorking || l.State == lsRun {
			g.Running++
		}
		g.LastAt = max(g.LastAt, l.last())
	}
	for _, g := range gs {
		sort.SliceStable(g.Lanes, func(a, b int) bool {
			la, lb := lanes[g.Lanes[a]], lanes[g.Lanes[b]]
			fa, fb := la.first(), lb.first()
			if fa < 0 {
				fa = 1 << 20
			}
			if fb < 0 {
				fb = 1 << 20
			}
			return fa < fb
		})
	}
	sort.SliceStable(gs, func(a, b int) bool {
		if gs[a].Need != gs[b].Need {
			return gs[a].Need > gs[b].Need
		}
		if gs[a].Running != gs[b].Running {
			return gs[a].Running > gs[b].Running
		}
		return gs[a].LastAt > gs[b].LastAt
	})
	return gs
}

// laneOrder is the lanes in display order.
func laneOrder(lanes []laneRow) []int {
	var out []int
	for _, g := range groupLanes(lanes) {
		out = append(out, g.Lanes...)
	}
	return out
}

// ----- git facts of the selected lane's checkout -----

type gitFile struct {
	Path     string
	Add, Del int
	Binary   bool
}

// laneGit is the working tree state the Detail and diff tabs show.
type laneGit struct {
	Files []gitFile
	Diff  []string // `git diff` lines, capped
	Last  string   // last commit: "hash subject"
	Err   string
	At    time.Time
	Busy  bool
}

type laneGitMsg struct {
	cwd string
	g   laneGit
}

const laneGitTTL = 30 * time.Second
const diffCap = 400

var gitRunner = func(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	return string(out), err
}

func fetchLaneGit(cwd string) laneGit {
	g := laneGit{At: time.Now()}
	out, err := gitRunner(cwd, "diff", "--numstat", "HEAD")
	if err != nil {
		g.Err = "not a git checkout"
		return g
	}
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		f := strings.SplitN(ln, "\t", 3)
		if len(f) != 3 {
			continue
		}
		a, e1 := strconv.Atoi(f[0])
		d, e2 := strconv.Atoi(f[1])
		g.Files = append(g.Files, gitFile{Path: f[2], Add: a, Del: d, Binary: e1 != nil || e2 != nil})
	}
	if out, err := gitRunner(cwd, "diff", "--no-color", "-U1", "HEAD"); err == nil {
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) > diffCap {
			lines = lines[:diffCap]
		}
		if len(lines) == 1 && lines[0] == "" {
			lines = nil
		}
		g.Diff = lines
	}
	if out, err := gitRunner(cwd, "log", "-1", "--format=%h %s"); err == nil {
		g.Last = strings.TrimSpace(out)
	}
	return g
}
