package main

// The live loop and the actions the Agents tab triggers. External effects go through
// package variables so tests run without the scripts.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

// actionMsg is the result of a send, attach or relaunch.
type actionMsg struct{ text string }

const tickEvery = 2 * time.Second

// Refresh cadence in ticks: runs and sessions every tick, usage and accounts every 15
// (30s), the registry and the forges every 30 (60s).
const (
	usageTicks    = 15
	registryTicks = 30
)

var fetchers = map[string]func() (any, error){
	srcRuns:     func() (any, error) { return fetchRuns() },
	srcSessions: func() (any, error) { return fetchInterrupted() },
	srcUsage:    func() (any, error) { return fetchUsage() },
	srcAccounts: func() (any, error) { return fetchAccounts() },
	srcRegistry: func() (any, error) { return fetchRegistry() },
	srcMRs:      func() (any, error) { return fetchMRs() },
	srcToday:    func() (any, error) { return fetchToday() },
}

func fetchCmd(src string) tea.Cmd {
	return func() tea.Msg {
		v, err := fetchers[src]()
		return srcMsg{src: src, val: v, err: err, at: time.Now()}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(tickEvery, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func allSources() []string {
	return []string{srcRuns, srcSessions, srcUsage, srcAccounts, srcRegistry, srcMRs, srcToday}
}

// refreshCmds picks the sources due at this tick.
func refreshCmds(tick int) []tea.Cmd {
	cmds := []tea.Cmd{fetchCmd(srcRuns), fetchCmd(srcSessions)}
	if tick%usageTicks == 0 {
		cmds = append(cmds, fetchCmd(srcUsage), fetchCmd(srcAccounts))
	}
	if tick%registryTicks == 0 {
		cmds = append(cmds, fetchCmd(srcRegistry), fetchCmd(srcMRs), fetchCmd(srcToday))
	}
	return cmds
}

// runBin runs one of the bin/ scripts and returns its trimmed output.
var runBin = func(name string, args ...string) (string, error) {
	cmd := exec.Command(crowBin(name), args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func actionCmd(label string, name string, args ...string) tea.Cmd {
	return func() tea.Msg {
		out, err := runBin(name, args...)
		if err != nil {
			return actionMsg{label + " failed: " + firstLine(cmpOr(out, err.Error()))}
		}
		return actionMsg{label + ": " + firstLine(cmpOr(out, "done"))}
	}
}

// sendCmd types text into the agent's window (crow-focus terminal backend) and presses enter.
func sendCmd(a agentRow, text string) tea.Cmd {
	if a.PID == 0 {
		return func() tea.Msg { return actionMsg{"no running session to send to"} }
	}
	return actionCmd("sent", "crow-focus", "--pid", strconv.Itoa(a.PID), "--send", text)
}

// attachCmd focuses the agent's window, or resumes the session when no window holds it.
func attachCmd(a agentRow) tea.Cmd {
	if a.Cwd == "" || a.Account == "" {
		return func() tea.Msg { return actionMsg{"cannot jump: no workspace or account known"} }
	}
	return actionCmd("jump", "crow-focus", "--pid", strconv.Itoa(a.PID), "--account", a.Account, "--cwd", a.Cwd)
}

// relaunchTarget is the account crow-relaunch recommends.
var relaunchTarget = func() (string, error) {
	out, err := runBin("crow-relaunch", "status", "--json")
	if err != nil {
		return "", err
	}
	var st struct {
		Recommended string `json:"recommended"`
	}
	if err := json.Unmarshal([]byte(out), &st); err != nil {
		return "", err
	}
	if st.Recommended == "" {
		return "", fmt.Errorf("no account is free")
	}
	return st.Recommended, nil
}

func relaunchCmd(a agentRow) tea.Cmd {
	return func() tea.Msg {
		target, err := relaunchTarget()
		if err != nil {
			return actionMsg{"relaunch failed: " + err.Error()}
		}
		out, err := runBin("crow-relaunch", a.Session, target, "--yes")
		if err != nil {
			return actionMsg{"relaunch failed: " + firstLine(cmpOr(out, err.Error()))}
		}
		return actionMsg{"relaunched on " + target}
	}
}

// reviewStartedMsg ends a review start: the MR's detail is fetched again to show the run.
type reviewStartedMsg struct {
	key, text string
	apply     bool // the message ends an apply start
	failed    bool
}

func reviewCmd(r mrRow) tea.Cmd {
	return func() tea.Msg {
		args := []string{"start", r.Repo, strconv.Itoa(r.Number), "--mode", reviewMode(&r)}
		if r.Branch != "" {
			args = append(args, "--branch", r.Branch)
		}
		out, err := runBin("crow-review-run", args...)
		if err != nil {
			return reviewStartedMsg{r.key(), "review failed to start: " + firstLine(cmpOr(out, err.Error())), false, true}
		}
		return reviewStartedMsg{r.key(), "review started on " + r.Ref + ", findings appear here when it ends", false, false}
	}
}

// applyCmd resumes the review session to act on the decisions; the answer refetches the detail.
func applyCmd(r mrRow) tea.Cmd {
	return func() tea.Msg {
		out, err := runBin("crow-review-run", "apply", r.Repo, strconv.Itoa(r.Number))
		if err != nil {
			return reviewStartedMsg{r.key(), "apply failed to start: " + firstLine(cmpOr(out, err.Error())), true, true}
		}
		return reviewStartedMsg{r.key(), "apply started on " + r.Ref + ", it commits locally and pushes nothing", true, false}
	}
}

// postAction is one push or comment that leaves the machine, as crow-review-run actions lists it.
type postAction struct {
	ID       string   `json:"id"`
	Kind     string   `json:"kind"`
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	Warnings []string `json:"warnings"`
	Cmd      []string `json:"cmd"`
	Hash     string   `json:"hash"`
}

// postActionsMsg carries the pending actions of one MR; postDoneMsg ends one of them.
type postActionsMsg struct {
	key  string
	acts []postAction
	err  string
}

type postDoneMsg struct {
	key, text string
	ok        bool
}

func postActionsCmd(r mrRow) tea.Cmd {
	return func() tea.Msg {
		out, err := runBin("crow-review-run", "actions", r.Repo, strconv.Itoa(r.Number))
		if err != nil {
			return postActionsMsg{key: r.key(), err: "cannot list the actions: " + firstLine(cmpOr(out, err.Error()))}
		}
		var acts []postAction
		if err := json.Unmarshal([]byte(out), &acts); err != nil {
			return postActionsMsg{key: r.key(), err: "cannot read the actions: " + err.Error()}
		}
		return postActionsMsg{key: r.key(), acts: acts}
	}
}

// actCmd runs one action; the hash makes crow-review-run refuse it when it changed since shown.
func actCmd(r mrRow, a postAction) tea.Cmd {
	return func() tea.Msg {
		out, err := runBin("crow-review-run", "act", r.Repo, strconv.Itoa(r.Number), a.ID, "--hash", a.Hash)
		if err != nil {
			return postDoneMsg{r.key(), a.ID + " failed: " + firstLine(cmpOr(out, err.Error())), false}
		}
		return postDoneMsg{r.key(), firstLine(cmpOr(out, a.ID+" done")), true}
	}
}

// detailMsg is the diff and threads of one MR.
type detailMsg struct {
	key string
	d   *mrDetail
	err error
	at  time.Time
}

var detailFetcher = fetchMRDetail

func detailCmd(r mrRow) tea.Cmd {
	return func() tea.Msg {
		d, err := detailFetcher(r)
		return detailMsg{key: r.key(), d: d, err: err, at: time.Now()}
	}
}

const detailTTL = 2 * time.Minute

// reviewTTL is the detail refresh while a review run is going.
const reviewTTL = 20 * time.Second
