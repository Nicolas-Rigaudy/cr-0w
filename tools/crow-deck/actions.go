package main

// The deck's actions as deckActions: each one shows its exact text and command in the
// confirm overlay before it runs (autonomy.go).

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// scopeOf is a project's scope from the registry, else from the cards; unknown reads as work.
func (m model) scopeOf(project string) string {
	if m.live != nil && m.live.reg != nil {
		if p, ok := m.live.reg.find(project); ok && p.Scope != "" {
			return p.Scope
		}
	}
	for _, c := range m.fx.Cards {
		if c.Name == project || (c.Short != "" && c.Short == project) {
			return c.Scope
		}
	}
	for _, p := range m.fx.Projects {
		if p.Name == project {
			return p.Scope
		}
	}
	return "work"
}

// proposeSend asks before typing text into the agent's window.
func (m model) proposeSend(a agentRow, text string) (model, tea.Cmd) {
	if a.PID == 0 {
		m.flash = "no running session to send to"
		return m, nil
	}
	return m.propose(deckAction{
		Kind: kAnswer, Scope: m.scopeOf(a.Project), Title: "Send",
		What:  "Type this answer into " + cmpOr(a.Short, a.Name) + " (" + a.Account + ") and press enter:",
		Text:  text,
		Label: "sent",
		Cmd:   []string{"crow-focus", "--pid", strconv.Itoa(a.PID), "--send", text},
	})
}

// proposeAgentCtl asks before pausing (Escape in the window) or stopping (SIGTERM) a session.
func (m model) proposeAgentCtl(a agentRow, act string) (model, tea.Cmd) {
	flag := map[string]string{"pause": "--interrupt", "stop": "--stop"}[act]
	if a.PID == 0 || flag == "" {
		m.flash = "no running session to " + act
		return m, nil
	}
	what := map[string]string{
		"pause": "Interrupt the current turn of " + cmpOr(a.Short, a.Name) + " (Escape in its window).",
		"stop":  "End the claude session of " + cmpOr(a.Short, a.Name) + " (SIGTERM).",
	}[act]
	return m.propose(deckAction{
		Kind: kAgentCtl, Scope: m.scopeOf(a.Project), Title: map[string]string{"pause": "Pause", "stop": "Stop"}[act],
		What: what, Label: act,
		Cmd: []string{"crow-focus", "--pid", strconv.Itoa(a.PID), flag},
	})
}

// handoffAction opens a repo session through crow handoff with a literal task and end state.
func (m model) handoffAction(title, project, what, task, end string) deckAction {
	return deckAction{
		Kind: kSession, Scope: m.scopeOf(project), Title: title, What: what,
		Text: "task: " + task, Extra: []string{"end: " + end},
		Note:  "The session starts on a free allowed account and asks before any push.",
		Label: "handoff",
		Cmd:   []string{"crow", "handoff", project, "--task", task, "--end-state", end},
	}
}

// mrHandoff hands the selected MR to a repo session.
func mrHandoff(r mrRow, findings string) deckAction {
	task, end := handoffTask(&r, findings)
	return deckAction{
		Kind: kSession, Scope: r.Scope, Title: "Handoff",
		What: "Open a repo session on " + r.Project + " for " + r.Ref + ".",
		Text: "task: " + task, Extra: []string{"end: " + end},
		Label: "handoff",
		Cmd:   []string{"crow", "handoff", r.Name, "--task", task, "--end-state", end},
	}
}

// mailAction marks the mail done or snoozes it to tomorrow (Outlook flags, nothing sent).
func mailAction(verb string, r todayRow) deckAction {
	what := map[string]string{"done": "Mark this mail done (flag complete):", "snooze": "Snooze this mail to tomorrow:"}[verb]
	return deckAction{
		Kind: kMailFlag, Scope: "work", Title: map[string]string{"done": "Done", "snooze": "Snooze"}[verb],
		What: what, Text: r.Tag + ": " + cmpOr(r.Subject, r.What),
		Note:  "Changes the Outlook flag only; nothing is sent.",
		Label: verb,
		Cmd:   []string{"crow-today", verb, r.ID},
		done: func(out string, err error) tea.Msg {
			if err != nil {
				return mailActMsg{verb + " failed: " + lastLine(cmpOr(out, err.Error()))}
			}
			return mailActMsg{fit(r.What, 40) + ": " + firstLine(cmpOr(out, "done"))}
		},
	}
}
