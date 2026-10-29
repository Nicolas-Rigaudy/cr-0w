package main

// Live data: the bin/ scripts' JSON (crow-runs, crow-sessions, crow-usage, crow-accounts,
// crow-registry) is fetched by exec and mapped to fixtures. Nothing here calls a model:
// crow-usage runs with CROW_TITLES_PER_DAY=0 so it never starts its title worker.

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func crowHome() string {
	if h := os.Getenv("CROW_HOME"); h != "" {
		return h
	}
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "..", "..")
}

func crowBin(name string) string { return filepath.Join(crowHome(), "bin", name) }

// ----- source types (mirror the scripts' JSON) -----

type event struct {
	TS     string `json:"ts"`
	Action string `json:"action"`
}

type liveAgent struct {
	Kind          string  `json:"kind"`
	PID           int     `json:"pid"`
	Account       string  `json:"account"`
	Project       string  `json:"project"`
	Cwd           string  `json:"cwd"`
	SessionID     string  `json:"session_id"`
	State         string  `json:"state"`
	StartedAt     string  `json:"started_at"`
	LastActivityS int     `json:"last_activity_s"`
	LastAction    string  `json:"last_action"`
	LastPrompt    string  `json:"last_prompt"`
	ContextTokens int     `json:"context_tokens"`
	Recent        []event `json:"recent"`
	RunID         string  `json:"run_id"`
	Label         string  `json:"label"`
	Title         string  `json:"title"`
}

type finishedRun struct {
	RunID     string   `json:"run_id"`
	Project   string   `json:"project"`
	Label     string   `json:"label"`
	Status    string   `json:"status"`
	DurationS int      `json:"duration_s"`
	EndedAt   string   `json:"ended_at"`
	Commits   []string `json:"commits"`
}

type runsSnap struct {
	Slots struct {
		Busy int `json:"busy"`
		Max  int `json:"max"`
	} `json:"slots"`
	Agents   []liveAgent   `json:"agents"`
	Finished []finishedRun `json:"finished"`
}

// intrSession: a session cut off by a usage limit (crow-sessions list --json).
type intrSession struct {
	ID          string `json:"conversation_id"`
	Account     string `json:"account"`
	Project     string `json:"project"`
	Cwd         string `json:"cwd"`
	LastEvent   string `json:"last_event"`
	Limit       string `json:"limit"`
	ResetsAt    int64  `json:"resets_at"`
	FirstPrompt string `json:"first_prompt"`
	Relaunched  any    `json:"relaunched"`
}

type limitRow struct {
	Account          string   `json:"account"`
	Remaining5hPct   *int     `json:"remaining_5h_pct"`
	Reset5hAt        string   `json:"reset_5h_at"`
	RemainingWeekPct *int     `json:"remaining_week_pct"`
	ResetWeekAt      string   `json:"reset_week_at"`
	Burn5h           *float64 `json:"burn_5h_pct_per_h"`
	EmptyAt          string   `json:"empty_at"`
}

type usageSession struct {
	ID       string           `json:"id"`
	Account  string           `json:"account"`
	Project  string           `json:"project"`
	Title    string           `json:"title"`
	Start    string           `json:"start"`
	End      string           `json:"end"`
	Activity map[string][]int `json:"activity"`
	Kind     string           `json:"kind"`
	Branch   string           `json:"branch"`
	Models   []string         `json:"models"`
	Tokens   int              `json:"tokens"`
	CostUSD  float64          `json:"cost_usd"`
	Messages int              `json:"user_messages"`
	Today    struct {
		CostUSD float64 `json:"cost_usd"`
	} `json:"today"`
}

type usageName struct {
	Name    string  `json:"name"`
	CostUSD float64 `json:"cost_usd"`
}

type usageTotals struct {
	CostUSD         float64     `json:"cost_usd"`
	CacheSavedUSD   float64     `json:"cache_saved_usd"`
	RoutingSavedUSD float64     `json:"routing_saved_usd"`
	CachePct        int         `json:"cache_pct"`
	ByProject       []usageName `json:"by_project"`
	ByModel         []usageName `json:"by_model"`
}

type usageSnap struct {
	GeneratedAt string         `json:"generated_at"`
	Limits      []limitRow     `json:"limits"`
	Sessions    []usageSession `json:"sessions"`
	Today       usageTotals    `json:"today"`
	Week        usageTotals    `json:"week"`
	Days        []struct {
		Date    string  `json:"date"`
		CostUSD float64 `json:"cost_usd"`
	} `json:"days"`
	History []struct {
		Date    string  `json:"date"`
		CostUSD float64 `json:"cost_usd"`
	} `json:"history"`
	RTK struct {
		TodaySavedTokens int     `json:"today_saved_tokens"`
		TodaySavedPct    float64 `json:"today_saved_pct"`
		SavedTokens      int     `json:"saved_tokens"`
		SavedPct         float64 `json:"saved_pct"`
	} `json:"rtk"`
}

type acctPolicy struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`
}

type regGit struct {
	Email      string `json:"email"`
	Dirty      int    `json:"dirty"`
	Ahead      int    `json:"ahead"`
	StaleTrees int    `json:"stale_trees"`
}

type regProject struct {
	Name    string `json:"name"`
	Short   string `json:"short"`
	Scope   string `json:"scope"`
	Forge   string `json:"forge"`
	Remote  string `json:"remote"`
	Path    string `json:"path"`
	Status  string `json:"status"`
	Logics  bool   `json:"logics"`
	Release struct {
		Target  string `json:"target"`
		State   string `json:"state"`
		Roadmap string `json:"roadmap"`
	} `json:"release"`
	Road *regRoad `json:"road"`
	Git  regGit   `json:"git"`
}

// regRoad is the milestone progress crow-registry reads from the roadmap doc.
type regRoad struct {
	Done    int    `json:"done"`
	Total   int    `json:"total"`
	Current string `json:"current"`
	Exit    string `json:"exit"`
}

type registry struct {
	Identities map[string]string `json:"identities"`
	StateDir   string            `json:"state_dir"`
	Projects   []regProject      `json:"projects"`
}

// ----- fetching -----

func runJSON(timeout time.Duration, env []string, out any, name string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, crowBin(name), args...)
	cmd.Env = append(os.Environ(), env...)
	b, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return json.Unmarshal(b, out)
}

func fetchRuns() (*runsSnap, error) {
	var s runsSnap
	return &s, runJSON(20*time.Second, nil, &s, "crow-runs", "json")
}

func fetchInterrupted() ([]intrSession, error) {
	var r struct {
		Sessions []intrSession `json:"sessions"`
	}
	err := runJSON(20*time.Second, nil, &r, "crow-sessions", "list", "--json")
	return r.Sessions, err
}

func fetchUsage() (*usageSnap, error) {
	var s usageSnap
	return &s, runJSON(30*time.Second, []string{"CROW_TITLES_PER_DAY=0"}, &s, "crow-usage", "--json")
}

func fetchAccounts() ([]acctPolicy, error) {
	var r struct {
		Accounts []acctPolicy `json:"accounts"`
	}
	err := runJSON(10*time.Second, nil, &r, "crow-accounts", "json")
	return r.Accounts, err
}

func fetchRegistry() (*registry, error) {
	var r registry
	return &r, runJSON(30*time.Second, nil, &r, "crow-registry", "json")
}

// ----- live state -----

const (
	srcRuns     = "runs"
	srcSessions = "sessions"
	srcUsage    = "usage"
	srcAccounts = "accounts"
	srcRegistry = "registry"
	srcMRs      = "mrs"
)

// liveState holds the last good payload of every source with its age and last error.
type liveState struct {
	runs  *runsSnap
	intr  []intrSession
	usage *usageSnap
	acct  []acctPolicy
	reg   *registry
	mrs   *mrsSnap
	at    map[string]time.Time
	err   map[string]string
	// MR diffs and threads, fetched when an MR is selected.
	details map[string]*detailEntry
}

type detailEntry struct {
	d    *mrDetail
	at   time.Time
	err  string
	busy bool
}

func newLiveState() *liveState {
	return &liveState{at: map[string]time.Time{}, err: map[string]string{}, details: map[string]*detailEntry{}}
}

// srcMsg is one source's result.
type srcMsg struct {
	src string
	val any
	err error
	at  time.Time
}

func (ls *liveState) apply(msg srcMsg) {
	if msg.err != nil {
		ls.err[msg.src] = msg.err.Error()
		return
	}
	delete(ls.err, msg.src)
	ls.at[msg.src] = msg.at
	switch v := msg.val.(type) {
	case *runsSnap:
		ls.runs = v
	case []intrSession:
		ls.intr = v
	case *usageSnap:
		ls.usage = v
	case []acctPolicy:
		ls.acct = v
	case *registry:
		ls.reg = v
	case *mrsSnap:
		ls.mrs = v
	}
}

// applyDetail stores the diff of one MR; a failure is kept with its time so it is retried
// after the TTL, not on every tick.
func (ls *liveState) applyDetail(msg detailMsg) {
	e := ls.details[msg.key]
	if e == nil {
		e = &detailEntry{}
		ls.details[msg.key] = e
	}
	e.busy, e.at = false, msg.at
	if msg.err != nil {
		e.err = firstLine(msg.err.Error())
		return
	}
	e.d, e.err = msg.d, ""
}

// syncState summarises the sources for the header dot: loading until runs arrived, red
// naming the failing source, yellow when the oldest source is older than staleAfter.
const staleAfter = 90 * time.Second

func (ls *liveState) syncState(now time.Time) (state, note string) {
	if ls.runs == nil && len(ls.err) == 0 {
		return stLoading, ""
	}
	for _, src := range []string{srcRuns, srcSessions, srcUsage, srcAccounts, srcRegistry, srcMRs} {
		if _, bad := ls.err[src]; bad {
			return stError, map[string]string{
				srcRuns: "cdx not running", srcSessions: "sessions scan failed", srcUsage: "usage failed",
				srcAccounts: "accounts failed", srcRegistry: "registry failed", srcMRs: "forge sync failed",
			}[src]
		}
	}
	if ls.mrs != nil {
		for _, name := range []string{"gitlab", "github"} {
			if st := ls.mrs.Forges[name]; st != "" && st != "ok" {
				return stError, name + " sync failed"
			}
		}
	}
	for _, t := range ls.at {
		if now.Sub(t) > staleAfter {
			return stStale, "stale " + shortAge(now.Sub(t))
		}
	}
	return stFresh, ""
}

func shortAge(d time.Duration) string {
	s := int(d.Seconds())
	switch {
	case s < 0:
		return "0s"
	case s < 60:
		return fmt.Sprintf("%ds", s)
	case s < 3600:
		return fmt.Sprintf("%dm", s/60)
	case s < 86400:
		return fmt.Sprintf("%dh%02d", s/3600, (s%3600)/60)
	}
	return fmt.Sprintf("%dd", s/86400)
}

// ----- mapping to fixtures -----

func parseTime(s string) (time.Time, bool) {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05.000Z"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func clock(t time.Time) string { return t.Local().Format("15:04") }

// shortName maps a project name from the scripts to the registry's short name.
func (r *registry) shortName(name string) string {
	name = strings.TrimSuffix(name, " (hub)")
	if r != nil {
		for _, p := range r.Projects {
			if p.Name == name || p.Short == name {
				return p.Short
			}
		}
	}
	return name
}

func (r *registry) scopeOf(name string) string {
	name = strings.TrimSuffix(name, " (hub)")
	if r != nil {
		for _, p := range r.Projects {
			if p.Name == name || p.Short == name {
				return p.Scope
			}
		}
	}
	return ""
}

func firstLine(s string) string {
	s, _, _ = strings.Cut(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(s)
}

func tokensK(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return fmt.Sprintf("%d", n)
}

// agentRows maps the runs snapshot and the interrupted registry to rows, newest first
// inside each group. The group is what the Agents tab sorts on.
func agentRows(now time.Time, runs *runsSnap, intr []intrSession, reg *registry) []agentRow {
	var rows []agentRow
	if runs == nil {
		return nil
	}
	for _, a := range runs.Agents {
		r := agentRow{
			Account: a.Account, PID: a.PID, Session: a.SessionID, Cwd: a.Cwd, Project: reg.shortName(a.Project),
			Title: a.Title, Ctx: a.ContextTokens,
		}
		started, _ := parseTime(a.StartedAt)
		idle := time.Duration(a.LastActivityS) * time.Second
		last := firstLine(strings.TrimPrefix(a.LastAction, "says: "))
		switch {
		case a.Kind == "run":
			r.Group, r.State = gRuns, "working"
			if a.State == "quiet" {
				r.State = "idle"
			}
			r.Session = a.RunID
			r.Name = cmpOr(a.Label, "run "+short4(a.RunID))
			r.Short = r.Name
			r.Doing = cmpOr(last, "running")
			r.Since = shortAge(now.Sub(started))
		case a.State == "waiting":
			r.Group, r.State = gNeeds, "needs"
			r.Name = r.Project
			r.Short = cmpOr(a.Title, r.Project)
			r.Doing = "asks you: " + cmpOr(last, "waiting for your answer")
			r.Ask = strings.TrimSpace(strings.TrimPrefix(a.LastAction, "says: "))
			r.Since = shortAge(idle)
		default:
			r.Group, r.State = gRunning, "working"
			if a.State == "idle" {
				r.State = "idle"
			}
			r.Name = r.Project
			r.Short = cmpOr(a.Title, r.Project)
			r.Doing = cmpOr(a.Title, last)
			r.Since = shortAge(now.Sub(started))
		}
		r.Meta = a.Account + " · " + tokensK(a.ContextTokens) + " tokens · " + shortAge(now.Sub(started))
		if !started.IsZero() {
			r.Started = started
		}
		for _, e := range a.Recent {
			if t, ok := parseTime(e.TS); ok {
				r.Tail = append(r.Tail, tailLine{At: clock(t), Text: e.Action})
			}
		}
		// recent is newest first in crow-runs; the tail reads oldest first.
		for i, j := 0, len(r.Tail)-1; i < j; i, j = i+1, j-1 {
			r.Tail[i], r.Tail[j] = r.Tail[j], r.Tail[i]
		}
		rows = append(rows, r)
	}
	for _, s := range intr {
		if s.Relaunched != nil {
			continue
		}
		t, ok := parseTime(s.LastEvent)
		if !ok || !sameDay(t, now) {
			continue
		}
		resets := ""
		if s.ResetsAt > 0 {
			resets = "back " + clock(time.Unix(s.ResetsAt, 0))
		}
		name := cmpOr(reg.shortName(s.Project), "session")
		rows = append(rows, agentRow{
			Group: gRuns, State: "waiting", Name: name, Short: name + " " + short4(s.ID), Account: s.Account,
			Doing: "paused on " + limitName(s.Limit) + ", " + resets + ", R relaunch", Since: shortAge(now.Sub(t)),
			Session: s.ID, Cwd: s.Cwd, Interrupted: true, Meta: s.Account + " · " + limitName(s.Limit),
			Tail: []tailLine{{At: clock(t), Text: "stopped on the " + limitName(s.Limit)}, {At: "", Text: resets + " (R relaunches it on another account)"}},
		})
	}
	for _, f := range runs.Finished {
		t, ok := parseTime(f.EndedAt)
		if !ok || !sameDay(t, now) {
			continue
		}
		name := reg.shortName(f.Project)
		doing := "finished: " + cmpOr(f.Label, "run")
		if n := len(f.Commits); n > 0 {
			doing += fmt.Sprintf(", %d %s", n, plural(n, "commit", "commits"))
		}
		rows = append(rows, agentRow{
			Group: gDone, State: "done", Name: name, Short: name + " " + cmpOr(f.Label, short4(f.RunID)), Account: "",
			Doing: doing, Since: clock(t), Session: f.RunID, Meta: f.Status + " · " + shortAge(time.Duration(f.DurationS)*time.Second),
			Tail: []tailLine{{At: clock(t), Text: doing}},
		})
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Group < rows[j].Group })
	return rows
}

func limitName(kind string) string {
	switch kind {
	case "five_hour":
		return "5h limit"
	case "seven_day", "week", "weekly":
		return "weekly limit"
	}
	return "usage limit"
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func cmpOr(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func short4(id string) string {
	if len(id) > 4 {
		return id[:4]
	}
	return id
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Local().Date()
	by, bm, bd := b.Local().Date()
	return ay == by && am == bm && ad == bd
}

// gaugeRows maps cdx limits and the account policy to the Usage gauges; accounts with
// policy never are left out (crow does not use them).
func gaugeRows(runs *runsSnap, usage *usageSnap, acct []acctPolicy) []gauge {
	if usage == nil {
		return nil
	}
	policy := map[string]string{}
	for _, a := range acct {
		policy[a.Name] = a.Policy
	}
	used := map[string]int{}
	slots := 0
	if runs != nil {
		slots = runs.Slots.Max
		for _, a := range runs.Agents {
			if a.Kind == "run" {
				used[a.Account]++
			}
		}
	}
	var out []gauge
	for _, l := range usage.Limits {
		pol := policy[l.Account]
		if pol == "never" || l.Remaining5hPct == nil {
			continue
		}
		g := gauge{Account: l.Account, FiveH: *l.Remaining5hPct, Week: -1, Policy: pol, RunsUsed: used[l.Account], RunsMax: slots}
		g.Reset, g.WeekReset = resetTime(l.Reset5hAt), l.ResetWeekAt
		g.Empty = emptyClock(l.EmptyAt)
		if l.RemainingWeekPct != nil {
			g.Week = *l.RemainingWeekPct
		}
		if l.Burn5h != nil {
			g.Burn, g.BurnPct = fmt.Sprintf("burn %.0f%%/h", *l.Burn5h), *l.Burn5h
		}
		out = append(out, g)
	}
	return out
}

// resetTime keeps the clock of a "Oct 29 21:10" reset, whatever the day.
func resetTime(s string) string {
	if i := strings.LastIndex(s, " "); i >= 0 {
		return s[i+1:]
	}
	return s
}

// whoWorked builds the timeline rows from today's 10-minute activity buckets.
func whoWorked(now time.Time, usage *usageSnap, intr []intrSession, reg *registry) []whoRow {
	if usage == nil {
		return nil
	}
	day := now.Local().Format("2006-01-02")
	paused := map[string]intrSession{}
	for _, s := range intr {
		if s.Relaunched == nil {
			paused[s.ID] = s
		}
	}
	var rows []whoRow
	for _, s := range usage.Sessions {
		b := s.Activity[day]
		if len(b) == 0 {
			continue
		}
		row := whoRow{Name: whoName(s, reg), Account: s.Account}
		sort.Ints(b)
		from, prev := b[0], b[0]
		for _, x := range append(b[1:], 1<<30) {
			if x > prev+1 {
				row.On = append(row.On, hspan{float64(from) / 6, float64(prev+1) / 6})
				from = x
			}
			prev = x
		}
		if p, ok := paused[s.ID]; ok && p.ResetsAt > 0 {
			end := time.Unix(p.ResetsAt, 0).Local()
			row.Off = append(row.Off, hspan{float64(b[len(b)-1]+1) / 6, float64(end.Hour()) + float64(end.Minute())/60})
		}
		rows = append(rows, row)
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].On[0].From < rows[j].On[0].From })
	return rows
}

// whoName is the session title without its "project: " prefix; a title that is only a slash
// command or empty falls back to the project name.
func whoName(s usageSession, reg *registry) string {
	proj := reg.shortName(s.Project)
	title := strings.TrimSpace(s.Title)
	if i := strings.Index(title, ": "); i >= 0 && i < 24 {
		title = strings.TrimSpace(title[i+2:])
	}
	if title == "" || strings.HasPrefix(title, "/") {
		return proj
	}
	return title
}

// openThreads reads the unchecked "- [ ] date [tag] text" lines of open-threads.md as
// the Next list: project column from the tag, text up to the first arrow.
func openThreads(path string, reg *registry) []nextRow {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []nextRow
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		line := sc.Text()
		rest, ok := strings.CutPrefix(line, "- [ ] ")
		if !ok {
			continue
		}
		fields := strings.SplitN(rest, " ", 3)
		if len(fields) < 3 || !strings.HasPrefix(fields[1], "[") {
			continue
		}
		tag := strings.Trim(fields[1], "[]")
		what, _, _ := strings.Cut(fields[2], " → ")
		name := "crow"
		if tag != "hub" {
			name = reg.shortName(tag)
		}
		scope := reg.scopeOf(tag)
		if tag == "hub" {
			scope = "personal"
		}
		out = append(out, nextRow{Project: name, What: strings.TrimSpace(what), Perso: scope == "personal"})
	}
	return out
}

// health: identity matches the registry for every repo, cdx answers, worktrees merged.
func healthOf(ls *liveState) health {
	h := health{Identity: true, CDX: ls.runs != nil && ls.err[srcRuns] == ""}
	if ls.reg == nil {
		return h
	}
	for _, p := range ls.reg.Projects {
		if p.Git.Email != "" && ls.reg.Identities[p.Scope] != "" && p.Git.Email != ls.reg.Identities[p.Scope] {
			h.Identity = false
		}
		h.StaleTrees += p.Git.StaleTrees
	}
	return h
}

type liveInput struct {
	Now time.Time
	LS  *liveState
}

// buildFixtures maps the live payloads to what the pages draw. Sources not wired yet keep
// an honest note instead of invented rows.
func buildFixtures(in liveInput) *fixtures {
	ls, now := in.LS, in.Now
	fx := &fixtures{Now: now.Local().Format("Mon 2 Jan 15:04"), Live: true}
	fx.Synced = clock(now)
	state, note := ls.syncState(now)
	fx.SyncNote = note
	fx.Health = healthOf(ls)
	if t, ok := ls.at[srcRuns]; ok {
		fx.Synced = clock(t)
	}
	fx.Agents = agentRows(now, ls.runs, ls.intr, ls.reg)
	if ls.runs != nil {
		fx.Slots = [2]int{ls.runs.Slots.Busy, ls.runs.Slots.Max}
	}
	fx.Gauges = gaugeRows(ls.runs, ls.usage, ls.acct)
	fx.Usage = usageViewOf(ls.usage, ls.reg, now)
	fx.Who = whoWorked(now, ls.usage, ls.intr, ls.reg)
	fx.Hist = histOf(ls.usage, ls.reg, now)
	if ls.reg != nil {
		fx.Next = openThreads(filepath.Join(ls.reg.StateDir, "open-threads.md"), ls.reg)
		for _, p := range ls.reg.Projects {
			fx.Repos++
			if p.Git.Dirty > 0 {
				fx.Dirty++
			}
			if p.Git.Ahead > 0 {
				fx.Ahead++
			}
		}
	}
	if ls.mrs != nil {
		fx.MRs, fx.Branches, fx.MRsReady = mrOrder(ls.mrs.MRs), ls.mrs.Branches, true
		fx.MRsAwaiting, fx.MRsMine, fx.MRsThreads = mrCounts(fx.MRs)
		fx.Blockers = mrBlockers(now, fx.MRs)
		for _, name := range []string{"gitlab", "github"} {
			if st := ls.mrs.Forges[name]; st != "" && st != "ok" {
				fx.ForgeErrs = append(fx.ForgeErrs, st)
			}
		}
	}
	fx.Cards, fx.Quiet = projectCardsOf(now, ls.reg, fx.Agents, fx.Next, fx.MRs)
	fx.Projects = projectRows(fx.Cards)
	fx.MRDetails, fx.MRErrs = map[string]*mrDetail{}, map[string]string{}
	for k, e := range ls.details {
		if e.d != nil {
			fx.MRDetails[k] = e.d
		}
		if e.err != "" {
			fx.MRErrs[k] = e.err
		}
	}
	fx.Sync = state
	return fx
}
