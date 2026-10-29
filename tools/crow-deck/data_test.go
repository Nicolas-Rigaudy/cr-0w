package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// All payloads below are invented. They follow the shape of the scripts' JSON.

func testNow() time.Time {
	return time.Date(2026, 10, 29, 14, 7, 0, 0, time.Local)
}

func iso(h, m int) string {
	return time.Date(2026, 10, 29, h, m, 0, 0, time.Local).UTC().Format(time.RFC3339)
}

func testRuns() *runsSnap {
	r := &runsSnap{}
	r.Slots.Busy, r.Slots.Max = 1, 2
	r.Agents = []liveAgent{
		{PID: 11, Account: "main", Project: "atlas", Cwd: "/w/atlas", SessionID: "s-wait", State: "waiting", StartedAt: iso(13, 55),
			LastActivityS: 60, LastAction: "says: Keep this finding?\nmore", ContextTokens: 31000, Title: "review !1",
			Recent: []event{{TS: iso(14, 6), Action: "checking config"}, {TS: iso(14, 4), Action: "reading file"}}},
		{PID: 12, Account: "main", Project: "orbit (hub)", Cwd: "/w/orbit", SessionID: "s-work", State: "working", StartedAt: iso(14, 0),
			LastAction: "edit x.go", ContextTokens: 5000},
		{PID: 13, Account: "main", Project: "orbit", SessionID: "s-idle", State: "idle", StartedAt: iso(9, 0), Title: "notes"},
		{Kind: "run", PID: 14, Account: "spare", Project: "atlas", RunID: "abcd1234", Label: "lint", State: "working", StartedAt: iso(14, 1)},
	}
	r.Finished = []finishedRun{
		{RunID: "ffff0001", Project: "atlas", Label: "tests", Status: "ok", DurationS: 600, EndedAt: iso(11, 5), Commits: []string{"a", "b"}},
		{RunID: "ffff0002", Project: "atlas", Label: "old", Status: "ok", EndedAt: time.Date(2026, 10, 28, 11, 0, 0, 0, time.Local).UTC().Format(time.RFC3339)},
	}
	return r
}

func testReg() *registry {
	return &registry{
		Identities: map[string]string{"work": "w@example.test", "personal": "p@example.test"},
		StateDir:   "/nonexistent",
		Projects: []regProject{
			{Name: "atlas", Short: "atl", Scope: "work", Git: regGit{Email: "w@example.test", StaleTrees: 1}},
			{Name: "orbit", Short: "orb", Scope: "personal", Git: regGit{Email: "p@example.test", Dirty: 2}},
		},
	}
}

func TestAgentRowsGroupsAndSort(t *testing.T) {
	resets := time.Date(2026, 10, 29, 16, 30, 0, 0, time.Local).Unix()
	intr := []intrSession{
		{ID: "paused-session-1", Account: "main", Project: "atlas", Cwd: "/w/atlas", LastEvent: iso(13, 0), Limit: "five_hour", ResetsAt: resets},
		{ID: "gone", Project: "atlas", LastEvent: iso(13, 0), Relaunched: "yes"},
		{ID: "yesterday", Project: "atlas", LastEvent: time.Date(2026, 10, 28, 9, 0, 0, 0, time.Local).UTC().Format(time.RFC3339)},
	}
	rows := agentRows(testNow(), testRuns(), intr, testReg())
	var got []string
	for _, r := range rows {
		got = append(got, fmt.Sprintf("%d/%s/%s", r.Group, r.State, r.Short))
	}
	want := []string{
		"0/needs/review !1", "1/working/orb", "1/idle/notes", "2/working/lint",
		"2/waiting/atl paus", "3/done/atl tests",
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("rows\n got %v\nwant %v", got, want)
	}
	needs := rows[0]
	if needs.Ask != "Keep this finding?\nmore" || needs.Doing != "asks you: Keep this finding?" {
		t.Errorf("needs row: ask %q doing %q", needs.Ask, needs.Doing)
	}
	if len(needs.Tail) != 2 || needs.Tail[0].Text != "reading file" {
		t.Errorf("tail not oldest first: %+v", needs.Tail)
	}
	if needs.Meta != "main · 31k tokens · 12m" {
		t.Errorf("meta %q", needs.Meta)
	}
	var paused agentRow
	for _, r := range rows {
		if r.Interrupted {
			paused = r
		}
	}
	if !strings.Contains(paused.Doing, "paused on 5h limit, back 16:30, R relaunch") {
		t.Errorf("paused doing %q", paused.Doing)
	}
	if n := (&fixtures{Agents: rows}).agentsNeedYou(); n != 1 {
		t.Errorf("needs you = %d, want 1", n)
	}
}

func TestSyncState(t *testing.T) {
	now := testNow()
	ls := newLiveState()
	if st, _ := ls.syncState(now); st != stLoading {
		t.Errorf("empty: %s", st)
	}
	ls.apply(srcMsg{src: srcRuns, val: testRuns(), at: now})
	if st, _ := ls.syncState(now); st != stFresh {
		t.Errorf("fresh: %s", st)
	}
	if st, note := ls.syncState(now.Add(5 * time.Minute)); st != stStale || note != "stale 5m" {
		t.Errorf("stale: %s %q", st, note)
	}
	ls.apply(srcMsg{src: srcUsage, err: errors.New("boom")})
	if st, note := ls.syncState(now); st != stError || note != "usage failed" {
		t.Errorf("error: %s %q", st, note)
	}
	ls.apply(srcMsg{src: srcUsage, val: &usageSnap{}, at: now})
	if st, _ := ls.syncState(now); st != stFresh {
		t.Errorf("recovered: %s", st)
	}
	bad := newLiveState()
	bad.apply(srcMsg{src: srcRuns, err: errors.New("no cdx")})
	if st, note := bad.syncState(now); st != stError || note != "cdx not running" {
		t.Errorf("first fetch failed: %s %q", st, note)
	}
}

func TestGaugeRowsSkipsNever(t *testing.T) {
	five, week, burn := 40, 70, 12.4
	usage := &usageSnap{Limits: []limitRow{
		{Account: "main", Remaining5hPct: &five, RemainingWeekPct: &week, Burn5h: &burn, Reset5hAt: "Oct 29 21:10"},
		{Account: "locked", Remaining5hPct: &five},
		{Account: "nodata"},
	}}
	g := gaugeRows(testRuns(), usage, []acctPolicy{{"main", "use"}, {"locked", "never"}})
	if len(g) != 1 || g[0].Account != "main" || g[0].Reset != "21:10" || g[0].Burn != "burn 12%/h" || g[0].RunsMax != 2 || g[0].Week != 70 {
		t.Fatalf("gauges %+v", g)
	}
}

func TestWhoWorkedSpansAndPause(t *testing.T) {
	resets := time.Date(2026, 10, 29, 16, 0, 0, 0, time.Local).Unix()
	usage := &usageSnap{Sessions: []usageSession{
		{ID: "s1", Account: "main", Project: "atlas", Title: "design", Activity: map[string][]int{"2026-10-29": {54, 55, 56, 60}}},
		{ID: "s2", Account: "spare", Project: "orbit", Activity: map[string][]int{"2026-10-29": {40, 41}}},
		{ID: "s3", Account: "main", Project: "orbit", Activity: map[string][]int{"2026-10-28": {10}}},
	}}
	rows := whoWorked(testNow(), usage, []intrSession{{ID: "s2", ResetsAt: resets}}, testReg())
	if len(rows) != 2 || rows[0].Name != "orb" || rows[1].Name != "design" {
		t.Fatalf("rows %+v", rows)
	}
	if len(rows[1].On) != 2 || rows[1].On[0] != (hspan{9, 9.5}) || rows[1].On[1] != (hspan{10, 10 + 1.0/6}) {
		t.Errorf("spans %+v", rows[1].On)
	}
	if len(rows[0].Off) != 1 || rows[0].Off[0].To != 16 {
		t.Errorf("paused span %+v", rows[0].Off)
	}
	if n := whoName(usageSession{Project: "atlas", Title: "atlas: fix retry"}, testReg()); n != "fix retry" {
		t.Errorf("name %q", n)
	}
	if n := whoName(usageSession{Project: "atlas", Title: "atlas: /retro run"}, testReg()); n != "atl" {
		t.Errorf("slash title %q", n)
	}
	if u := unionHours([]hspan{{9, 11}, {10, 12}, {13, 14}}); u != 4 {
		t.Errorf("union %v", u)
	}
}

func TestOpenThreads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "open-threads.md")
	body := "# Open threads\n- [ ] 2026-10-28 [hub] write the runbook → next\n- [x] 2026-10-27 [atlas] done thing\n" +
		"- [ ] 2026-10-29 [atlas] fix retry\n- [ ] 2026-10-29 [orbit] polish\nnot a thread\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	got := openThreads(path, testReg())
	want := []nextRow{{"crow", "write the runbook", true}, {"atl", "fix retry", false}, {"orb", "polish", true}}
	if len(got) != len(want) {
		t.Fatalf("rows %+v", got)
	}
	for i := range want {
		if got[i].Project != want[i].Project || got[i].What != want[i].What || got[i].Perso != want[i].Perso {
			t.Errorf("row %d: %+v, want %+v", i, got[i], want[i])
		}
	}
	if openThreads(filepath.Join(t.TempDir(), "none.md"), testReg()) != nil {
		t.Error("missing file should give no rows")
	}
}

func TestBuildFixturesLive(t *testing.T) {
	ls := newLiveState()
	now := testNow()
	ls.apply(srcMsg{src: srcRuns, val: testRuns(), at: now})
	ls.apply(srcMsg{src: srcRegistry, val: testReg(), at: now})
	fx := buildFixtures(liveInput{Now: now, LS: ls})
	if !fx.Live || fx.Sync != stFresh || fx.Slots != [2]int{1, 2} || len(fx.Agents) != 5 {
		t.Fatalf("fixtures %+v", fx)
	}
	if fx.Repos != 2 || fx.Dirty != 1 || fx.Health.StaleTrees != 1 || !fx.Health.Identity {
		t.Errorf("repos %d dirty %d health %+v", fx.Repos, fx.Dirty, fx.Health)
	}
	reg := testReg()
	reg.Projects[0].Git.Email = "other@example.test"
	ls.apply(srcMsg{src: srcRegistry, val: reg, at: now})
	if buildFixtures(liveInput{Now: now, LS: ls}).Health.Identity {
		t.Error("identity mismatch not flagged")
	}
}

// fakeHome writes stub scripts so the exec path of the fetchers runs for real.
func fakeHome(t *testing.T, scripts map[string]string) {
	t.Helper()
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, out := range scripts {
		body := "#!/bin/sh\ncat <<'EOF'\n" + out + "\nEOF\n"
		if err := os.WriteFile(filepath.Join(home, "bin", name), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("CROW_HOME", home)
}

func TestFetchThroughScripts(t *testing.T) {
	fakeHome(t, map[string]string{
		"crow-runs":     `{"slots":{"busy":1,"max":2},"agents":[{"pid":5,"account":"main","project":"atlas","state":"working"}],"finished":[]}`,
		"crow-sessions": `{"sessions":[{"conversation_id":"c1","account":"main","limit":"five_hour"}]}`,
		"crow-accounts": `{"accounts":[{"name":"main","policy":"use"}]}`,
		"crow-usage":    `not json`,
	})
	runs, err := fetchRuns()
	if err != nil || runs.Slots.Max != 2 || len(runs.Agents) != 1 {
		t.Fatalf("runs %+v %v", runs, err)
	}
	if s, err := fetchInterrupted(); err != nil || len(s) != 1 || s[0].ID != "c1" {
		t.Errorf("sessions %+v %v", s, err)
	}
	if a, err := fetchAccounts(); err != nil || a[0].Policy != "use" {
		t.Errorf("accounts %+v %v", a, err)
	}
	if _, err := fetchUsage(); err == nil {
		t.Error("bad JSON should fail")
	}
	if _, err := fetchRegistry(); err == nil {
		t.Error("missing script should fail")
	}
}

// ----- live model and actions -----

type call struct {
	name string
	args []string
}

func recordRuns(t *testing.T, out string, err error) *[]call {
	t.Helper()
	var calls []call
	old := runBin
	runBin = func(name string, args ...string) (string, error) {
		calls = append(calls, call{name, args})
		return out, err
	}
	t.Cleanup(func() { runBin = old })
	return &calls
}

func liveModel(t *testing.T) model {
	t.Helper()
	m := newModel(defaultConfig())
	m.page = 1
	m.goLive()
	m.now = testNow()
	m = send(m, srcMsg{src: srcRuns, val: testRuns(), at: testNow()}, srcMsg{src: srcRegistry, val: testReg(), at: testNow()})
	return m
}

func TestLiveModelLoadsAndKeepsCursor(t *testing.T) {
	m := newModel(defaultConfig())
	m.goLive()
	if m.sync != stLoading {
		t.Fatalf("sync %s before data", m.sync)
	}
	m = liveModel(t)
	if m.sync != stFresh || len(m.fx.Agents) != 5 {
		t.Fatalf("sync %s agents %d", m.sync, len(m.fx.Agents))
	}
	m = send(m, keyMsg("j"), keyMsg("j"))
	if m.agentKey != "s-idle" {
		t.Fatalf("cursor on %q", m.agentKey)
	}
	// A refresh that adds a row above must not move the cursor off its agent.
	r := testRuns()
	r.Agents = append([]liveAgent{{PID: 99, Account: "main", Project: "atlas", SessionID: "s-new", State: "waiting", StartedAt: iso(14, 6)}}, r.Agents...)
	m = send(m, srcMsg{src: srcRuns, val: r, at: testNow()})
	if m.fx.Agents[m.agentSel].Session != "s-idle" {
		t.Errorf("cursor moved to %q", m.fx.Agents[m.agentSel].Session)
	}
}

func TestAgentAnswerKeysSendThroughFocus(t *testing.T) {
	m := liveModel(t)
	calls := recordRuns(t, "sent", nil)
	m = send(m, keyMsg("a"))
	if !m.typing || m.input != "yes" || m.agentPane != apTail {
		t.Fatalf("a: typing %v input %q", m.typing, m.input)
	}
	m = send(m, keyMsg("esc"))
	m = send(m, keyMsg("e"))
	if m.input != "yes, but " {
		t.Errorf("e input %q", m.input)
	}
	m = send(m, keyMsg("esc"), keyMsg("x"))
	if m.input != "no" {
		t.Errorf("x input %q", m.input)
	}
	m = send(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("!")})
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = next.(model)
	if m.typing || m.input != "" || cmd == nil {
		t.Fatalf("enter: typing %v input %q cmd %v", m.typing, m.input, cmd != nil)
	}
	msg := cmd().(actionMsg)
	if len(*calls) != 1 || (*calls)[0].name != "crow-focus" || strings.Join((*calls)[0].args, " ") != "--pid 11 --send no!" {
		t.Errorf("calls %+v", *calls)
	}
	if msg.text != "sent: sent" {
		t.Errorf("flash %q", msg.text)
	}
	// An agent that does not wait for an answer ignores a x e c.
	m = send(m, keyMsg("j"), keyMsg("a"))
	if m.typing || !strings.Contains(m.flash, "does not wait") {
		t.Errorf("running agent: typing %v flash %q", m.typing, m.flash)
	}
}

func TestAttachAndFailureFlash(t *testing.T) {
	m := liveModel(t)
	calls := recordRuns(t, "no window", errors.New("exit 1"))
	_, cmd := m.Update(keyMsg("g"))
	msg := cmd().(actionMsg)
	got := strings.Join((*calls)[0].args, " ")
	if got != "--pid 11 --account main --cwd /w/atlas" || msg.text != "jump failed: no window" {
		t.Errorf("args %q flash %q", got, msg.text)
	}
}

func TestRelaunchNeedsConfirm(t *testing.T) {
	m := liveModel(t)
	resets := time.Date(2026, 10, 29, 16, 30, 0, 0, time.Local).Unix()
	m = send(m, srcMsg{src: srcSessions, at: testNow(), val: []intrSession{
		{ID: "paused-1", Account: "main", Project: "atlas", LastEvent: iso(13, 0), Limit: "five_hour", ResetsAt: resets}}})
	m = send(m, keyMsg("R"))
	if m.overlay == ovConfirm {
		t.Fatal("R on a running agent must not open the confirm")
	}
	for i, a := range m.fx.Agents {
		if a.Interrupted {
			m.selectAgent(i)
		}
	}
	calls := recordRuns(t, "", nil)
	old := relaunchTarget
	relaunchTarget = func() (string, error) { return "spare", nil }
	defer func() { relaunchTarget = old }()
	m = send(m, keyMsg("R"))
	if m.overlay != ovConfirm || len(*calls) != 0 {
		t.Fatalf("overlay %d calls %d: nothing may run before the confirm", m.overlay, len(*calls))
	}
	if !hasText(m.View(), "Relaunch") {
		t.Error("confirm overlay not drawn")
	}
	m = send(m, keyMsg("esc"))
	if m.overlay != ovNone || len(*calls) != 0 {
		t.Fatalf("esc: overlay %d calls %d", m.overlay, len(*calls))
	}
	m = send(m, keyMsg("R"))
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = next.(model)
	if cmd == nil {
		t.Fatal("enter did not start the relaunch")
	}
	if msg := cmd().(actionMsg); msg.text != "relaunched on spare" {
		t.Errorf("flash %q", msg.text)
	}
	if c := (*calls)[0]; c.name != "crow-relaunch" || strings.Join(c.args, " ") != "paused-1 spare --yes" {
		t.Errorf("call %+v", c)
	}
}

func TestRefreshSchedule(t *testing.T) {
	if n := len(refreshCmds(1)); n != 2 {
		t.Errorf("tick 1: %d fetches, want runs and sessions", n)
	}
	if n := len(refreshCmds(usageTicks)); n != 4 {
		t.Errorf("usage tick: %d, want 4", n)
	}
	if n := len(refreshCmds(registryTicks)); n != 5 {
		t.Errorf("registry tick: %d, want 5", n)
	}
}

func TestStaleAndErrorShowInHeader(t *testing.T) {
	withANSI(t)
	m := liveModel(t)
	m.w, m.h = 170, 46
	m = send(m, tickMsg(testNow().Add(10*time.Minute)))
	if m.sync != stStale || !hasText(m.View(), "stale 10m") {
		t.Errorf("stale: sync %s", m.sync)
	}
	m = send(m, srcMsg{src: srcRuns, err: errors.New("x")})
	if m.sync != stError || !hasText(m.View(), "cdx not running") {
		t.Errorf("error: sync %s", m.sync)
	}
	checkFit(t, "live error", m.View(), 170, 46)
}

func TestLiveAgentsFitEverySize(t *testing.T) {
	withANSI(t)
	for ti, th := range themes {
		for _, sz := range sizes {
			for pane := 0; pane < agentPanes; pane++ {
				m := liveModel(t)
				m.th, m.w, m.h, m.agentPane = ti, sz[0], sz[1], pane
				checkFit(t, fmt.Sprintf("%s %dx%d pane %d", th.name, sz[0], sz[1], pane), m.View(), sz[0], sz[1])
			}
		}
	}
}
