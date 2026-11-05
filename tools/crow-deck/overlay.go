package main

// Overlays: ? help, : palette, t theme. They replace the page body.

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ovNone = iota
	ovHelp
	ovPalette
	ovTheme
	ovConfirm
	ovPost
	ovHandoff
	ovReview
	ovApply
)

type command struct{ name, desc string }

func paletteCommands() []command {
	var cs []command
	for i, n := range pageNames {
		cs = append(cs, command{"tab " + strings.ToLower(n), "open " + strconv.Itoa(i+1) + " " + n})
	}
	for _, n := range themeNames() {
		cs = append(cs, command{"theme " + n, "switch the theme"})
	}
	for _, m := range []string{"work", "perso", "all"} {
		cs = append(cs, command{"mode " + m, "filter Next and Projects by scope"})
	}
	for _, s := range syncStates {
		cs = append(cs, command{"mock " + s, "show the " + s + " state (invented data)"})
	}
	cs = append(cs, command{"history", "search the session history (Agents)"}, command{"live", "back to the live agents"}, command{"attach", "jump to the selected agent's window"}, command{"pause", "interrupt the selected agent's turn (Escape, asks first)"}, command{"stop", "end the selected claude session (SIGTERM, asks first)"}, command{"relaunch", "continue the selected paused session on another account"})
	return append(cs, command{"quit", "leave the deck"})
}

func filterCommands(q string) []command {
	var out []command
	for _, c := range paletteCommands() {
		ok := true
		for _, word := range strings.Fields(strings.ToLower(q)) {
			if !strings.Contains(c.name, word) {
				ok = false
			}
		}
		if ok {
			out = append(out, c)
		}
	}
	return out
}

func overlayWidth(w int) int { return min(w, 60) }

func (lk *look) overlayBlock(sp boxSpec, w, h int) []string {
	ow := overlayWidth(w)
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

func (lk *look) helpOverlay(w, h int, page int) []string {
	keys := [][2]string{
		{"1-5 & é\" ' (", "tabs (AZERTY digit row works unshifted)"}, {"tab", "next pane (one pane at a time under 100 columns)"}, {"j k", "move in the list"},
		{"m", "mode: work, perso, all"}, {"t", "switch theme"}, {":", "palette"}, {"?", "this help"},
		{"/", "Agents: history and search; esc back to live"},
		{"p a d M s x", "History: filter project, account, date, model, status; clear"}, {"a x e c", "Agents: answer yes, no, edit, chat"},
		{"enter", "Agents: open the lane's session; send a typed answer"},
		{"R", "Agents: relaunch a paused run on an account you pick"}, {"< > w", "Agents: move the hour range; full day"},
		{"d", "Agents: Detail full screen"},
		{"o enter x s", "Home Today: open the mail in Outlook, crow session drafts a reply, done, snooze to tomorrow"},
		{"q", "quit"},
	}
	sp := boxSpec{title: "Keys", sel: -1}
	for _, k := range keys {
		sp.rows = append(sp.rows, line{{fit(k[0], 14), lk.key}, {k[1], lk.plain}})
	}
	sp.rows = append(sp.rows, nil, line{{"esc closes", lk.sub}})
	return lk.overlayBlock(sp, w, h)
}

func (lk *look) paletteOverlay(w, h int, q string, sel int) []string {
	sp := boxSpec{title: "Palette", sel: -1}
	sp.rows = append(sp.rows, line{{": ", lk.key}, {q + "▏", lk.plain}}, nil)
	cmds := filterCommands(q)
	for i, c := range cmds {
		l := line{{fit(c.name, 18), lk.plain}, {c.desc, lk.sub}}
		if i == sel {
			sp.sel = len(sp.rows)
		}
		sp.rows = append(sp.rows, l)
	}
	if len(cmds) == 0 {
		sp.rows = append(sp.rows, line{{"No command matches.", lk.sub}})
	}
	return lk.overlayBlock(sp, w, h)
}

func (lk *look) themeOverlay(w, h, cur int) []string {
	sp := boxSpec{title: "Theme", sel: -1}
	for i, t := range themes {
		g, st := "○", lk.sub
		if i == cur {
			g, st = "◆", lk.acc
		}
		sp.rows = append(sp.rows, line{{g + " ", st}, {fit(t.name, 8), st}, {"accent " + slotName(lk.accentFor(t)), lk.sub}})
	}
	sp.rows = append(sp.rows, nil, line{{"t next · k previous · enter or esc close", lk.sub}})
	return lk.overlayBlock(sp, w, h)
}

// accentFor is the accent a theme would show: the config accent, else the default.
func (lk *look) accentFor(theme) int {
	if lk.custom >= 0 {
		return lk.custom
	}
	return defaultAccent
}

// postOverlay shows one push or comment: its exact text and command. enter runs that one
// action only, s skips it, esc stops. The text is never cut silently: a long one says so.
func (lk *look) postOverlay(w, h int, acts []postAction, i int, running bool) []string {
	sp := boxSpec{title: "Push and post", sel: -1}
	if i >= len(acts) {
		sp.rows = []line{{{"Nothing left to push or post.", lk.sub}}}
		return lk.overlayBlock(sp, w, h)
	}
	a, iw := acts[i], overlayWidth(w)-6
	head := []line{{{fmt.Sprintf("%d/%d ", i+1, len(acts)), lk.sub}, {a.Title, lk.bold}}, nil}
	cmd := []line{nil, {{"command", lk.sub}}}
	for _, l := range wrapLines(shellLine(a.Cmd), iw) {
		cmd = append(cmd, line{{l, lk.sub}})
	}
	var warn []line
	for _, x := range a.Warnings {
		for _, l := range wrapLines("! "+x, iw) {
			warn = append(warn, line{{l, lk.warn}})
		}
	}
	foot := line{{"enter", lk.key}, spc(1), {postVerb(a.Kind), lk.sub}, spc(2),
		{"s", lk.key}, spc(1), {"skip", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"stop", lk.sub}}
	if running {
		foot = line{{"● running, wait for the answer", lk.warn}}
	}
	text := wrapLines(a.Text, iw)
	room := max(1, h-2-len(head)-len(warn)-2) // box border, then the blank line and the footer
	if len(text)+len(cmd) > room {
		cmd = cmd[:min(len(cmd), max(2, room-len(text)))]
	}
	if keep := room - len(cmd); len(text) > keep {
		text = append(text[:max(0, keep-1)], fmt.Sprintf("… %d more lines", len(text)-max(0, keep-1)))
	}
	sp.rows = head
	for _, l := range text {
		sp.rows = append(sp.rows, line{{l, lk.plain}})
	}
	sp.rows = append(append(append(sp.rows, cmd...), warn...), nil, foot)
	return lk.overlayBlock(sp, w, h)
}

// shellLine prints a command the way a shell would read it; the body is the text shown above.
func shellLine(cmd []string) string {
	var out []string
	for i, x := range cmd {
		switch {
		case strings.HasPrefix(x, "body="):
			x = "body=<the text above>"
		case i > 0 && cmd[i-1] == "--body":
			x = "<the text above>"
		}
		if strings.ContainsAny(x, " '\"$[]{}<>") {
			x = "'" + strings.ReplaceAll(x, "'", `'\''`) + "'"
		}
		out = append(out, x)
	}
	return strings.Join(out, " ")
}

// reviewOverlay confirms a read-only review run on the selected MR.
func (lk *look) reviewOverlay(w, h int, r *mrRow, d *mrDetail) []string {
	sp := boxSpec{title: "Review", sel: -1}
	if r == nil {
		sp.rows = []line{{{"No MR selected.", lk.sub}}}
		return lk.overlayBlock(sp, w, h)
	}
	sp.rows = []line{{{"Start a read-only review of " + r.Ref + " " + r.Name + " (" + reviewMode(r, d) + ", opus)?", lk.plain}}, nil}
	for _, l := range wrap("The run cannot edit the repo. Findings land here; nothing is posted.", overlayWidth(w)-6) {
		sp.rows = append(sp.rows, line{{l, lk.sub}})
	}
	sp.rows = append(sp.rows, nil, line{{"enter", lk.key}, spc(1), {"start the run", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"cancel", lk.sub}})
	return lk.overlayBlock(sp, w, h)
}

// applyOverlay confirms an apply run: what it will do, where, and that nothing leaves the machine.
func (lk *look) applyOverlay(w, h int, r *mrRow, d *mrDetail) []string {
	sp := boxSpec{title: "Apply", sel: -1}
	if r == nil || d == nil {
		sp.rows = []line{{{"No review loaded.", lk.sub}}}
		return lk.overlayBlock(sp, w, h)
	}
	fix, chat := 0, 0
	for _, f := range d.Findings {
		if f.toApply(detailMode(r, d)) {
			if f.State == "chat" {
				chat++
			} else {
				fix++
			}
		}
	}
	q := fmt.Sprintf("Resume the review of %s to fix %d %s and answer %d chat %s?", r.Ref, fix, plural(fix, "finding", "findings"), chat, plural(chat, "note", "notes"))
	for _, l := range wrap(q, overlayWidth(w)-6) {
		sp.rows = append(sp.rows, line{{l, lk.plain}})
	}
	sp.rows = append(sp.rows, nil)
	for _, l := range wrap("Write rights in the review worktree only. One local commit, nothing pushed or posted.", overlayWidth(w)-6) {
		sp.rows = append(sp.rows, line{{l, lk.sub}})
	}
	sp.rows = append(sp.rows, nil, line{{"enter", lk.key}, spc(1), {"start the apply run", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"cancel", lk.sub}})
	return lk.overlayBlock(sp, w, h)
}

// reviewMode is threads on the owner's own MR with open threads, self on another own MR,
// peer otherwise. d is the loaded detail or nil; its threads count over the list row, which can predate them.
func reviewMode(r *mrRow, d *mrDetail) string {
	switch {
	case r.Role == "AUTHOR" && (r.ThreadsOpen > 0 || d.openThreads() > 0):
		return "threads"
	case r.Role == "AUTHOR":
		return "self"
	}
	return "peer"
}

// handoffTask is the literal task and end state `crow handoff` receives for an MR.
// findings is the findings file path of the reviewed MR, empty when it has none.
func handoffTask(r *mrRow, findings string) (task, end string) {
	note := ""
	if findings != "" {
		note = " Start from the review findings, decisions and notes in " + findings + "."
	}
	if r.Role == "REVIEWER" {
		return "Review " + r.Ref + " (" + r.Title + ") with the review skill and write the comments to post." + note,
			"Comments drafted and shown to me, nothing posted without my OK."
	}
	return "Address the " + strconv.Itoa(r.ThreadsOpen) + " open review threads on " + r.Ref + " (" + r.Title + ") with the review-threads skill." + note,
		"Fixes committed locally with the decision table shown to me. No push, no reply, no resolve without my OK."
}

// handoffOverlay shows the exact text before the repo session opens.
func (lk *look) handoffOverlay(w, h int, r *mrRow, findings string) []string {
	sp := boxSpec{title: "Handoff", sel: -1}
	if r == nil {
		sp.rows = []line{{{"No MR selected.", lk.sub}}}
		return lk.overlayBlock(sp, w, h)
	}
	task, end := handoffTask(r, findings)
	sp.rows = []line{{{"Open a repo session on " + r.Project + " for " + r.Ref + ".", lk.plain}}, nil}
	for _, l := range wrap("task: "+task, overlayWidth(w)-6) {
		sp.rows = append(sp.rows, line{{l, lk.plain}})
	}
	for _, l := range wrap("end: "+end, overlayWidth(w)-6) {
		sp.rows = append(sp.rows, line{{l, lk.sub}})
	}
	sp.rows = append(sp.rows, nil, line{{"enter", lk.key}, spc(1), {"open the session", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"cancel", lk.sub}})
	return lk.overlayBlock(sp, w, h)
}

// postVerb is what enter does on one pending action.
func postVerb(kind string) string {
	switch kind {
	case "push":
		return "push"
	case "reply":
		return "post this reply"
	case "resolve":
		return "resolve this thread"
	}
	return "post this comment"
}
