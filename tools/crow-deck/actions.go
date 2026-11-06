package main

// The deck's actions as deckActions: each one shows its exact text and command in the
// confirm overlay before it runs (autonomy.go).

import (
	"strconv"
	"strings"

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

// sameProject: a and b name the same project (short or full registry name, or the same text).
func (m model) sameProject(a, b string) bool {
	if strings.EqualFold(a, b) {
		return true
	}
	if m.live == nil || m.live.reg == nil {
		return false
	}
	pa, oka := m.live.reg.find(a)
	pb, okb := m.live.reg.find(b)
	return oka && okb && pa.Name == pb.Name
}

// reviewJump opens the MRs page on the project's MR (ref when given) and its review overlay,
// which asks before starting the review.
// project may list several names for the same project (card name and short name).
func (m model) reviewJump(ref string, project ...string) (tea.Model, tea.Cmd) {
	same := func(p string) bool {
		for _, n := range project {
			if n != "" && m.sameProject(p, n) {
				return true
			}
		}
		return false
	}
	vis := visibleMRs(m.fx, m.mode)
	at := -1
	for i, r := range vis {
		if ref != "" && r.Ref == ref && same(r.Project) {
			at = i
			break
		}
		if at < 0 && ref == "" && same(r.Project) {
			at = i
		}
	}
	if at < 0 {
		m.flash = "r: no open MR on " + cmpOr(ref, project[0])
		return m, nil
	}
	m.page = 2
	m.selectMR(at)
	cmd := m.wantDetail()
	nm, c2 := m.mrKeyPress("r")
	return nm, tea.Batch(cmd, c2)
}

// mrRef is the MR ref a Next line starts with ("!214", "#31"), else "".
func mrRef(what string) string {
	f := strings.Fields(what)
	if len(f) > 0 && len(f[0]) > 1 && (f[0][0] == '!' || f[0][0] == '#') {
		if _, err := strconv.Atoi(f[0][1:]); err == nil {
			return f[0]
		}
	}
	return ""
}

// cardFor is the Projects card of a project, nil when it has none.
func (m model) cardFor(project string) *projCard {
	for i, c := range m.fx.Cards {
		if m.sameProject(c.Name, project) || (c.Short != "" && m.sameProject(c.Short, project)) {
			return &m.fx.Cards[i]
		}
	}
	return nil
}

// cardTask is what a card's session starts on: its first Next line, else the milestone.
func cardTask(c *projCard) (task, end string) {
	task = "Continue milestone " + cmpOr(c.Mile, "the current milestone") + " of " + c.Name + "."
	if len(c.Next) > 0 && c.Next[0].What != "" {
		task = "Continue " + c.Name + ": " + c.Next[0].What + "."
	}
	end = "milestone exit: " + c.Exit
	if c.Exit == "" {
		end = "that step done, or what is left listed as open threads"
	}
	return task, end
}

// termAction opens a session on project in a new window, on a first message when task is set.
func (m model) termAction(title, project, what, task, end string) deckAction {
	cmd := []string{"crow", "term", "--on-busy", "fail", project}
	a := deckAction{Kind: kSession, Scope: m.scopeOf(project), Title: title, What: what, Label: strings.ToLower(title), Cmd: cmd,
		Note: "Opens a new window on a free allowed account. Nothing is pushed."}
	if task != "" {
		a.Text, a.Extra = "task: "+task, []string{"end: " + end}
		a.Cmd = append(cmd, "--task", task, "--end-state", end)
	}
	return a
}

// projectAction runs a card action key on the selected Projects card.
func (m model) projectAction(k string) (tea.Model, tea.Cmd) {
	c := m.projectsView().cur()
	if c == nil {
		m.flash = "no project selected"
		return m, nil
	}
	task, end := cardTask(c)
	switch k {
	case "c":
		return m.propose(m.termAction("Continue", c.Name, "Open a session on "+c.Name+" that continues its roadmap:", task, end))
	case "s":
		rt := "Research what milestone " + cmpOr(c.Mile, "the next milestone") + " of " + c.Name + " needs next and write the findings into its Logics docs. Change no code."
		return m.propose(m.termAction("Research", c.Name, "Open a research session on "+c.Name+":", rt, "findings in Logics, open questions listed"))
	case "o":
		return m.propose(m.termAction("Session", c.Name, "Open a session on "+c.Name+" with no first message.", "", ""))
	case "h":
		return m.propose(m.handoffAction("Handoff", c.Name, "Open a repo session on "+c.Name+" with a hand-off note:", task, end))
	case "r":
		return m.reviewJump("", c.Name, c.Short)
	}
	m.flash = k + ": no action on this key"
	return m, nil
}

// homeAction runs h (handoff) and r (review) on the Home row in focus.
func (m model) homeAction(k string) (tea.Model, tea.Cmd) {
	hv := m.homeView()
	project, what, ref := "", "", ""
	switch m.focus {
	case pNext:
		rows := hv.nextItems()
		if i := m.sel[pNext]; i >= 0 && i < len(rows) {
			r := rows[i]
			project, what, ref = cmpOr(r.Repo, r.Project), r.What, mrRef(r.What)
		}
	case pProjects:
		rows := hv.projects()
		if i := m.sel[pProjects]; i >= 0 && i < len(rows) {
			project = rows[i].Name
		}
	default:
		m.flash = k + ": pick a row in Next or Projects"
		return m, nil
	}
	if project == "" {
		m.flash = "no row selected"
		return m, nil
	}
	switch k {
	case "r":
		return m.reviewJump(ref, project)
	case "h":
		task, end := what, "done, or what is left listed as open threads"
		if task == "" {
			if c := m.cardFor(project); c != nil {
				task, end = cardTask(c)
			} else {
				task = "Continue the open work on " + project + "."
			}
		}
		return m.propose(m.handoffAction("Handoff", project, "Open a repo session on "+project+" with a hand-off note:", task, end))
	case "x":
		m.flash = "x: snooze is not built yet (where a snooze is kept is the owner's call)"
	}
	return m, nil
}
