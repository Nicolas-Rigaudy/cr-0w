package main

// Relaunch picker and no-model mode (item_009). The account table comes from
// `crow-relaunch status --json`: every allowed account with its policy, 5h and week
// headroom, reset times and the reason crow treats it as blocked. The Relaunch overlay
// lets the owner pick the target (the recommended one preselected); choosing an account
// there is the owner's yes, so the deck passes --picked. When every account is blocked,
// Home Next opens with a "No tokens" group of work that needs no model.

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const srcRelay = "relay"

// relAcct is one account of crow-relaunch status.
type relAcct struct {
	Name         string `json:"name"`
	Policy       string `json:"policy"`
	Current      bool   `json:"current"`
	Blocked      string `json:"blocked"`
	Left5h       int    `json:"remaining_5h_pct"`
	LeftWeek     int    `json:"remaining_week_pct"`
	Reset5h      string `json:"reset_5h_at"`
	Reset5hISO   string `json:"reset_5h_at_iso"`
	ResetWeek    string `json:"reset_week_at"`
	ResetWeekISO string `json:"reset_week_at_iso"`
}

type relStatus struct {
	Accounts    []relAcct `json:"accounts"`
	Recommended string    `json:"recommended"`
}

var fetchRelStatus = func() (*relStatus, error) {
	out, err := runBin("crow-relaunch", "status", "--json")
	if err != nil {
		return nil, fmt.Errorf("%s", firstLine(cmpOr(out, err.Error())))
	}
	var st relStatus
	if err := json.Unmarshal([]byte(out), &st); err != nil {
		return nil, err
	}
	return &st, nil
}

// pickable: the account can take this session (not blocked, not never, not its own).
func (a relAcct) pickable(source string) bool {
	return a.Blocked == "" && a.Policy != "never" && a.Name != source
}

// relWhy is why an account cannot be picked, "" when it can.
func (a relAcct) why(source string) string {
	switch {
	case a.Name == source:
		return "the session's own account"
	case a.Policy == "never":
		return "policy never"
	}
	return a.Blocked
}

// relPreselect is the recommended account when it can take the session, else the first that can.
func relPreselect(st *relStatus, source string) int {
	if st == nil {
		return -1
	}
	first := -1
	for i, a := range st.Accounts {
		if !a.pickable(source) {
			continue
		}
		if a.Name == st.Recommended {
			return i
		}
		if first < 0 {
			first = i
		}
	}
	return first
}

// relStep moves the selection to the next pickable account in direction d.
func relStep(st *relStatus, source string, cur, d int) int {
	if st == nil {
		return cur
	}
	for i := cur + d; i >= 0 && i < len(st.Accounts); i += d {
		if st.Accounts[i].pickable(source) {
			return i
		}
	}
	return cur
}

func relaunchCmd(a agentRow, target string) tea.Cmd {
	return func() tea.Msg {
		out, err := runGated(kSession, "", "crow-relaunch", a.Session, target, "--picked")
		if err != nil {
			return actionMsg{"relaunch failed: " + lastLine(cmpOr(out, err.Error()))}
		}
		return actionMsg{"relaunched on " + target}
	}
}

// openRelaunch shows the picker and refreshes the account table behind it.
func (m *model) openRelaunch() tea.Cmd {
	m.overlay = ovConfirm
	src := ""
	if a := m.curAgent(); a != nil {
		src = a.Account
	}
	m.relSel = relPreselect(m.fx.Rel, src)
	if m.live == nil {
		return nil
	}
	return fetchCmd(srcRelay)
}

func (m model) relaunchKey(k string) (tea.Model, tea.Cmd) {
	a := m.curAgent()
	src := ""
	if a != nil {
		src = a.Account
	}
	switch k {
	case "j", "down":
		m.relSel = relStep(m.fx.Rel, src, m.relSel, 1)
	case "k", "up":
		m.relSel = relStep(m.fx.Rel, src, m.relSel, -1)
	case "enter", "y":
		st := m.fx.Rel
		if a == nil || st == nil || m.relSel < 0 || m.relSel >= len(st.Accounts) || !st.Accounts[m.relSel].pickable(src) {
			m.flash = "no account can take this session"
			return m, nil
		}
		m.overlay = ovNone
		target := st.Accounts[m.relSel].Name
		m.flash = "relaunching on " + target + "…"
		return m, relaunchCmd(*a, target)
	case "esc", "n", "q":
		m.overlay = ovNone
	}
	return m, nil
}

func (lk *look) confirmOverlay(w, h int, a *agentRow, st *relStatus, sel int) []string {
	sp := boxSpec{title: "Relaunch", sel: -1}
	name, src := "this session", ""
	if a != nil {
		name, src = cmpOr(a.Short, a.Name), a.Account
	}
	sp.rows = []line{{{"Continue ", lk.plain}, {name, lk.bold}, {" on:", lk.plain}}, nil}
	foot := line{{"j k", lk.key}, spc(1), {"account", lk.sub}, spc(2), {"enter", lk.key}, spc(1), {"relaunch", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"cancel", lk.sub}}
	if st == nil {
		sp.rows = append(sp.rows, line{{"Reading the accounts…", lk.sub}}, nil, foot)
		return lk.relBlock(sp, w, h)
	}
	gw := 10
	sp.rows = append(sp.rows, line{{"  " + fit("account", 9), lk.head}, {fit("policy", 11), lk.head}, {fit("5h left", gw+12), lk.head}, {"week left", lk.head}})
	var notes []line
	for i, x := range st.Accounts {
		st, g := lk.plain, "  "
		if i == sel {
			st, g = lk.acc, "● "
		}
		why := x.why(src)
		if why != "" {
			st = lk.dim
			notes = append(notes, line{{x.Name + ": " + why, lk.bad}})
		}
		pol := seg{fit(x.Policy, 11), lk.sub}
		if x.Policy == "ask-first" {
			pol.style = lk.warn
		}
		g5 := lk.ok
		if x.Left5h < 20 {
			g5 = lk.bad
		}
		l := line{{g, st}, {fit(x.Name, 9), st}, pol}
		l = append(l, gaugeBar(lk, x.Left5h, gw, g5)...)
		l = append(l, seg{fit(" "+strconv.Itoa(x.Left5h)+"% "+clockOf(x.Reset5h), 12), lk.sub})
		l = append(l, gaugeBar(lk, x.LeftWeek, gw, lk.ok)...)
		l = append(l, seg{" " + strconv.Itoa(x.LeftWeek) + "%", lk.sub})
		sp.rows = append(sp.rows, l)
	}
	if len(notes) > 0 {
		sp.rows = append(append(sp.rows, nil), notes...)
	}
	if sel >= 0 && sel < len(st.Accounts) && st.Accounts[sel].Policy == "ask-first" {
		sp.rows = append(sp.rows, nil, line{{st.Accounts[sel].Name + " is ask-first: enter counts as your yes.", lk.warn}})
	}
	if sel < 0 {
		sp.rows = append(sp.rows, nil, line{{"No account can take this session now.", lk.bad}})
	}
	sp.rows = append(sp.rows, nil, foot)
	return lk.relBlock(sp, w, h)
}

// relBlock is overlayBlock at the table's width (the gauges need 68 columns).
func (lk *look) relBlock(sp boxSpec, w, h int) []string {
	ow := min(w, 72)
	sp.focused, sp.index = true, 1
	out := lk.box(sp, ow, min(h, lk.boxH(len(sp.rows))))
	for i := range out {
		out[i] += blank(w - ow)
	}
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out
}

// clockOf keeps the time of "Nov 5 15:20".
func clockOf(s string) string {
	if i := strings.LastIndex(s, " "); i >= 0 {
		return s[i+1:]
	}
	return s
}

// noTokens reports whether every allowed account is blocked, and the earliest time one
// frees up ("" when unknown).
func noTokens(st *relStatus, now time.Time) (bool, string) {
	if st == nil || len(st.Accounts) == 0 {
		return false, ""
	}
	var first time.Time
	for _, a := range st.Accounts {
		if a.Policy == "never" {
			continue
		}
		if a.Blocked == "" {
			return false, ""
		}
		iso := a.Reset5hISO
		if strings.HasPrefix(a.Blocked, "week") {
			iso = a.ResetWeekISO
		}
		if t, err := time.Parse(time.RFC3339, iso); err == nil && t.After(now) && (first.IsZero() || t.Before(first)) {
			first = t
		}
	}
	if first.IsZero() {
		return true, ""
	}
	if first.Sub(now) >= 24*time.Hour {
		return true, first.Local().Format("Mon 15:04")
	}
	return true, first.Local().Format("15:04")
}

// noModelRows is the work that needs no model: the MRs of You decide (red pipelines, reviews
// to read, threads), repos to tidy, and the runs paused on a limit.
func noModelRows(decide []nextRow, reg *registry, agents []agentRow) []nextRow {
	out := append([]nextRow(nil), decide...)
	if reg != nil {
		for _, p := range reg.Projects {
			name := reg.shortName(p.Name)
			if p.Git.Dirty > 0 {
				out = append(out, nextRow{Project: name, What: strconv.Itoa(p.Git.Dirty) + " " + plural(p.Git.Dirty, "file", "files") + " uncommitted", Perso: p.Scope == "personal", Actor: "me", Why: "repo"})
			}
			if p.Git.StaleTrees > 0 {
				out = append(out, nextRow{Project: name, What: strconv.Itoa(p.Git.StaleTrees) + " stale " + plural(p.Git.StaleTrees, "worktree", "worktrees"), Perso: p.Scope == "personal", Actor: "me", Why: "repo"})
			}
		}
	}
	for _, a := range agents {
		if a.Interrupted {
			out = append(out, nextRow{Project: a.Project, What: a.Doing, Actor: "me", Why: "paused run"})
		}
	}
	return out
}
