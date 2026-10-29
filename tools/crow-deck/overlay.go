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
		cs = append(cs, command{"mock " + s, "show the " + s + " state (skeleton data)"})
	}
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
		{"/", "search or filter (lists, later slices)"}, {"a x e c", "answer: accept, refuse, edit, chat (later slices)"},
		{"enter h r x R", "open, handoff, review, snooze, relaunch (later slices)"}, {"q", "quit"},
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
