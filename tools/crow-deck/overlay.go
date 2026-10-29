package main

// Overlays: ? help, : palette, t theme. They replace the page body.

import (
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
		cs = append(cs, command{"mode " + m, "filter Next and Blockers by scope"})
	}
	for _, s := range syncStates {
		cs = append(cs, command{"mock " + s, "show the " + s + " state (invented data)"})
	}
	cs = append(cs, command{"history", "search the session history (Agents)"}, command{"live", "back to the live agents"}, command{"attach", "jump to the selected agent's window"}, command{"relaunch", "continue the selected paused session on another account"})
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
		{"1-5", "tabs"}, {"tab", "next pane (one pane at a time under 100 columns)"}, {"j k", "move in the list"},
		{"m", "mode: work, perso, all"}, {"t", "switch theme"}, {":", "palette"}, {"?", "this help"},
		{"/", "Agents: history and search; esc back to live"},
		{"p a d M s x", "History: filter project, account, date, model, status; clear"}, {"a x e c", "Agents: answer yes, no, edit, chat"},
		{"enter", "Agents: send the typed answer"}, {"g", "Agents: jump to the agent's window"},
		{"R", "Agents: relaunch a paused run"}, {"h r x", "handoff, review, snooze (later slices)"}, {"q", "quit"},
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

// accentFor is the accent a theme would show: the config accent, else its default.
func (lk *look) accentFor(t theme) int {
	if lk.custom >= 0 {
		return lk.custom
	}
	return t.accent
}

func (lk *look) confirmOverlay(w, h int, a *agentRow) []string {
	sp := boxSpec{title: "Relaunch", sel: -1}
	name := "this session"
	if a != nil {
		name = cmpOr(a.Short, a.Name)
	}
	sp.rows = []line{
		{{"Continue " + name + " on the account crow-relaunch recommends?", lk.plain}},
		{{"It opens a new window and uses that account's limits.", lk.sub}},
		nil,
		{{"enter", lk.key}, spc(1), {"relaunch", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"cancel", lk.sub}},
	}
	return lk.overlayBlock(sp, w, h)
}

// postOverlay is the confirm drawn before a post or push. The deck never sends: enter is
// shown disabled and the text says where it goes through.
func (lk *look) postOverlay(w, h int, r *mrRow) []string {
	sp := boxSpec{title: "Post", sel: -1}
	ref := "this MR"
	if r != nil {
		ref = r.Ref + " " + r.Name
	}
	sp.rows = []line{
		{{"Post the accepted findings and replies on " + ref + "?", lk.plain}},
		{{"Nothing leaves the deck yet.", lk.sub}},
		nil,
		{{"enter", lk.sub}, spc(1), {"post (not wired)", lk.sub}, spc(2), {"esc", lk.key}, spc(1), {"close", lk.sub}},
	}
	return lk.overlayBlock(sp, w, h)
}
