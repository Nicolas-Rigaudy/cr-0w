package main

// Autonomy: what the deck may do, per action kind and project scope. The table mirrors the
// Autonomy table of CLAUDE.md (a test compares them). Every action that changes something
// goes through propose: a confirm overlay with the exact text and command, then one run on
// one OK. runGated is the only way an action reaches a bin/ script; it refuses what the
// table forbids even when a caller skips the overlay.

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type rule int

const (
	ruleAllow rule = iota // done after the overlay's OK
	ruleAsk               // approval-gated: exact text shown, one OK per action
	ruleNever             // refused, no overlay
)

type actKind string

// Kinds from the CLAUDE.md table.
const (
	kCommit    actKind = "commit"
	kPush      actKind = "push"
	kMerge     actKind = "merge"
	kMR        actKind = "mr"
	kTicket    actKind = "ticket"
	kWiki      actKind = "wiki"
	kVault     actKind = "vault"
	kMail      actKind = "mail"
	kInfraDev  actKind = "infra-dev"
	kInfraProd actKind = "infra-prod"
	kDispatch  actKind = "dispatch"
)

// Deck kinds: actions that stay on this machine.
const (
	kSession  actKind = "session"  // open, hand off or resume a session
	kAnswer   actKind = "answer"   // type into an agent's window
	kAgentCtl actKind = "agentctl" // pause or stop an agent
	kMailFlag actKind = "mailflag" // mark a mail done or snooze it
	kLocal    actKind = "local"    // deck and review state files
)

type autonomyRow struct {
	row            string // the Action cell of the CLAUDE.md table; "" for deck kinds
	work, personal rule
}

var autonomy = map[actKind]autonomyRow{
	kCommit:    {"Branch / commit", ruleAllow, ruleAllow},
	kPush:      {"Push", ruleAsk, ruleAllow},
	kMerge:     {"Merge", ruleNever, ruleAllow},
	kMR:        {"Open MR/PR, write descriptions, review, verify review fixes", ruleAllow, ruleAllow},
	kTicket:    {"Create Jira tickets", ruleAllow, ruleNever},
	kWiki:      {"Edit Confluence", ruleAllow, ruleNever},
	kVault:     {"Write / reorganize the vault", ruleAllow, ruleAllow},
	kMail:      {"Send email, post on Teams", ruleAsk, ruleAsk},
	kInfraDev:  {"Infra on dev", ruleAllow, ruleAllow},
	kInfraProd: {"Infra on int / prod", ruleAsk, ruleNever},
	kDispatch:  {"Dispatch cdx runs", ruleAllow, ruleAllow},

	kSession:  {"", ruleAllow, ruleAllow},
	kAnswer:   {"", ruleAllow, ruleAllow},
	kAgentCtl: {"", ruleAllow, ruleAllow},
	kMailFlag: {"", ruleAllow, ruleAllow},
	kLocal:    {"", ruleAllow, ruleAllow},
}

// gate is the rule for kind in scope. Unknown kinds are refused; any scope but personal
// reads as work, the stricter column.
func gate(kind actKind, scope string) rule {
	r, ok := autonomy[kind]
	if !ok {
		return ruleNever
	}
	if scope == "personal" {
		return r.personal
	}
	return r.work
}

// scopeWord names the column gate used.
func scopeWord(scope string) string {
	if scope == "personal" {
		return "personal"
	}
	return "work"
}

// refusal is the flash for a refused action.
func refusal(kind actKind, scope string) string {
	what := string(kind)
	if r := autonomy[kind].row; r != "" {
		what = strings.ToLower(r)
	}
	return "refused: " + what + " on " + scopeWord(scope) + " projects is never done from the deck (autonomy table)"
}

var errRefused = errors.New("refused by the autonomy table")

// runGated runs a bin/ script for an action of kind in scope, unless the table forbids it.
func runGated(kind actKind, scope, name string, args ...string) (string, error) {
	if gate(kind, scope) == ruleNever {
		return refusal(kind, scope), errRefused
	}
	return runBin(name, args...)
}

// postKind maps a crow-review-run action to its table row: a push is a push, a merge a
// merge, comments, replies and resolves are review work.
func postKind(kind string) actKind {
	switch kind {
	case "push":
		return kPush
	case "merge":
		return kMerge
	case "reply", "resolve", "comment", "note":
		return kMR
	}
	return actKind("post:" + kind)
}

// deckAction is one action waiting for its OK in the confirm overlay.
type deckAction struct {
	Kind  actKind
	Scope string
	Title string   // overlay title: "Send", "Handoff"
	What  string   // one line: what happens, on what
	Text  string   // the exact text sent or handed over, "" when none
	Extra []string // more exact fields (end state), shown under the text
	Note  string   // what it does not do
	Label string   // flash prefix of the result
	Cmd   []string // bin/ script and its arguments
	// done turns the script's answer into the message the deck reacts to; nil gives an actionMsg.
	done func(out string, err error) tea.Msg
}

// propose opens the confirm overlay for a, or refuses it at once.
func (m model) propose(a deckAction) (model, tea.Cmd) {
	if gate(a.Kind, a.Scope) == ruleNever {
		m.flash = refusal(a.Kind, a.Scope)
		return m, nil
	}
	m.act, m.overlay = &a, ovAct
	return m, nil
}

// actKey answers the confirm overlay: enter or y runs the action once, esc, n or q drops it.
func (m model) actKey(k string) (tea.Model, tea.Cmd) {
	a := m.act
	switch k {
	case "enter", "y":
		m.act, m.overlay = nil, ovNone
		if a == nil {
			return m, nil
		}
		m.flash = cmpOr(a.Label, a.Title) + "…"
		return m, runAction(*a)
	case "esc", "n", "q":
		m.act, m.overlay = nil, ovNone
		if a != nil {
			m.flash = cmpOr(a.Label, a.Title) + " cancelled"
		}
	}
	return m, nil
}

// runAction runs one confirmed action through the gate.
func runAction(a deckAction) tea.Cmd {
	return func() tea.Msg {
		if len(a.Cmd) == 0 {
			return actionMsg{a.Label + ": nothing to run"}
		}
		out, err := runGated(a.Kind, a.Scope, a.Cmd[0], a.Cmd[1:]...)
		if a.done != nil {
			return a.done(out, err)
		}
		if err != nil {
			return actionMsg{a.Label + " failed: " + lastLine(cmpOr(out, err.Error()))}
		}
		return actionMsg{a.Label + ": " + firstLine(cmpOr(out, "done"))}
	}
}

// actOverlay shows the action exactly: what, the text, the command, and the rule that applies.
func (lk *look) actOverlay(w, h int, a *deckAction) []string {
	if a == nil {
		return lk.overlayBlock(boxSpec{title: "Confirm", sel: -1, rows: []line{{{"Nothing to confirm.", lk.sub}}}}, w, h)
	}
	sp := boxSpec{title: a.Title, sel: -1}
	iw := overlayWidth(w) - 6
	for _, l := range wrap(a.What, iw) {
		sp.rows = append(sp.rows, line{{l, lk.plain}})
	}
	if a.Text != "" {
		sp.rows = append(sp.rows, nil)
		for _, l := range wrapLines(a.Text, iw) {
			sp.rows = append(sp.rows, line{{l, lk.bold}})
		}
	}
	for _, x := range a.Extra {
		for _, l := range wrap(x, iw) {
			sp.rows = append(sp.rows, line{{l, lk.sub}})
		}
	}
	sp.rows = append(sp.rows, nil, line{{"command", lk.sub}})
	for _, l := range wrapLines(shellLine(a.Cmd), iw) {
		sp.rows = append(sp.rows, line{{l, lk.sub}})
	}
	if a.Note != "" {
		sp.rows = append(sp.rows, nil)
		for _, l := range wrap(a.Note, iw) {
			sp.rows = append(sp.rows, line{{l, lk.sub}})
		}
	}
	if gate(a.Kind, a.Scope) == ruleAsk {
		sp.rows = append(sp.rows, line{{"! needs your OK each time on " + scopeWord(a.Scope) + " projects", lk.warn}})
	}
	sp.rows = append(sp.rows, nil, line{{"enter", lk.key}, spc(1), {"run this once", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"cancel", lk.sub}})
	return lk.overlayBlock(sp, w, h)
}
