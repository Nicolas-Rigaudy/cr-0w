package main

// Chrome drawn on every page and theme: header, tabs, rule, status bar.

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var pageNames = []string{"Home", "Agents", "MRs", "Projects", "Usage"}

// Sync states shown by the dot in the header.
const (
	stFresh   = "fresh"
	stLoading = "loading"
	stStale   = "stale"
	stError   = "error"
	stEmpty   = "empty" // data state: nothing to show, sync is fine
)

var syncStates = []string{stFresh, stLoading, stStale, stError, stEmpty}

// chromeState is what the chrome needs from the model.
type chromeState struct {
	page    int
	mode    string
	sync    string
	fx      *fixtures
	flash   string
	context []keyHint
	hist    bool // Agents shows the history
}

type keyHint struct{ key, desc string }

func (lk *look) dot(style lipgloss.Style) seg { return seg{"●", style} }

func (lk *look) syncSegs(cs chromeState) line {
	switch cs.sync {
	case stLoading:
		return line{lk.dot(lk.warn), spc(1), {"syncing…", lk.sub}}
	case stStale:
		return line{lk.dot(lk.warn), spc(1), {cmpOr(cs.fx.SyncNote, "stale 2h"), lk.warn}}
	case stError:
		return line{lk.dot(lk.bad), spc(1), {cmpOr(cs.fx.SyncNote, "forge sync failed"), lk.bad}}
	}
	return line{lk.dot(lk.ok), spc(1), {lk.cap("synced ") + cs.fx.Synced, lk.sub}}
}

func (lk *look) healthSegs(cs chromeState, labels bool) line {
	h := cs.fx.Health
	items := []struct {
		ok    bool
		label string
	}{{h.Identity, "identity"}, {h.CDX, "cdx"}}
	var out line
	for _, it := range items {
		st := lk.ok
		if !it.ok {
			st = lk.bad
		}
		out = append(out, spc(2), lk.dot(st))
		if labels {
			out = append(out, spc(1), seg{it.label, lk.sub})
		}
	}
	trees := lk.ok
	label := "trees clean"
	if h.StaleTrees > 0 {
		trees, label = lk.warn, strconv.Itoa(h.StaleTrees)+" stale tree"
		if h.StaleTrees > 1 {
			label += "s"
		}
	}
	out = append(out, spc(2), lk.dot(trees))
	if labels {
		out = append(out, spc(1), seg{label, lk.sub})
	}
	return out
}

func (lk *look) brand(page int) line {
	switch lk.th.frame {
	case frameHUD:
		return line{{"◆ CROW//DECK", lk.acc}, spc(1), {"░▒▓", lk.dim}}
	case frameCyber:
		return line{{"▌CROW_DECK", lk.acc}, spc(1), {"//", lk.label}, spc(1), {strings.ToLower(pageNames[page]), lk.sub}, spc(1), {"░▒▓", lk.dim}}
	case frameSoft:
		return line{{"Crow deck", lk.acc}}
	}
	return line{{"◆ crow deck", lk.title}}
}

func (lk *look) modeSegs(mode string, word bool) line {
	var out line
	out = append(out, spc(2))
	if word {
		out = append(out, seg{"mode", lk.sub}, spc(1))
	}
	for i, m := range []string{"all", "work", "perso"} {
		if i > 0 {
			out = append(out, spc(1))
		}
		if m == mode {
			out = append(out, seg{" " + m + " ", lk.sel.Bold(true)})
		} else {
			out = append(out, seg{m, lk.sub})
		}
	}
	return append(out, spc(1), seg{"m", lk.key})
}

// header: brand, mode switch, then sync dot, health dots and date. Parts drop in order
// (date, health labels, mode word) until the line fits.
func (lk *look) header(w int, cs chromeState) string {
	date := line{{cs.fx.Now, lk.sub}}
	right := func(labels, withDate bool) line {
		r := append(line{}, lk.syncSegs(cs)...)
		r = append(r, lk.healthSegs(cs, labels)...)
		if withDate {
			r = append(r, spc(2))
			r = append(r, date...)
		}
		return r
	}
	type cand struct{ labels, date, word bool }
	for _, c := range []cand{{true, true, true}, {true, false, true}, {false, false, true}, {false, false, false}} {
		left := append(lk.brand(cs.page), lk.modeSegs(cs.mode, c.word)...)
		r := right(c.labels, c.date)
		if left.width()+2+r.width() <= w {
			return lr(left, r, w).render(w, false, lk.sel)
		}
	}
	left := append(lk.brand(cs.page), lk.modeSegs(cs.mode, false)...)
	return lr(left, right(false, false), w).render(w, false, lk.sel)
}

func (lk *look) badge(page int, fx *fixtures) int {
	switch page {
	case 1:
		return fx.agentsNeedYou()
	case 2:
		return fx.mrsNeedYou()
	}
	return 0
}

func (lk *look) tabs(w int, cs chromeState) string {
	var l line
	for i, name := range pageNames {
		on := i == cs.page
		label := strconv.Itoa(i+1) + " " + lk.cap(name)
		textStyle := lk.sub
		if on {
			textStyle = lk.acc
		}
		switch lk.th.frame {
		case frameRound:
			if on {
				textStyle = lk.underlineAcc
			}
		case frameHUD:
			marker := "  "
			if on {
				marker = "▶ "
			}
			l = append(l, seg{marker, lk.acc})
		case frameCyber:
			edge := lk.dim
			if on {
				edge = lk.acc
			}
			l = append(l, seg{"▌", edge})
		case frameSoft:
			l = append(l, spc(1))
			if on {
				textStyle = lk.underlineAcc
			}
		}
		l = append(l, seg{label, textStyle})
		if b := lk.badge(i, cs.fx); b > 0 && cs.sync != stLoading {
			l = append(l, spc(1), seg{strconv.Itoa(b), lk.warn.Bold(true)})
		}
		switch lk.th.frame {
		case frameCyber:
			edge := lk.dim
			if on {
				edge = lk.acc
			}
			l = append(l, seg{"╱", edge})
		case frameSoft:
			l = append(l, spc(1))
		}
		l = append(l, spc(2))
	}
	if cs.page == 1 {
		if sw := lk.liveSwitch(cs.hist); l.width()+2+sw.width() <= w {
			return lr(l, sw, w).render(w, false, lk.sel)
		}
	}
	return l.render(w, false, lk.sel)
}

// liveSwitch is the live / history switch on the Agents tab row.
func (lk *look) liveSwitch(hist bool) line {
	live, his := seg{"live", lk.acc}, seg{" history ", lk.sub}
	if hist {
		live, his = seg{"live", lk.sub}, seg{" history ", lk.sel.Bold(true)}
	}
	return line{live, spc(1), {"/", lk.key}, his}
}

// rule is the thin line under the tabs.
func (lk *look) rule(w int) string {
	switch lk.th.frame {
	case frameHUD:
		stripe := "━━━━ ╌╌╌ ░░"
		return lk.acc.Render(stripe) + lk.dim.Render(rep("╌", w-lipgloss.Width(stripe)))
	case frameCyber:
		return lk.acc.Render("//////////") + lk.dim.Render(rep("─", w-10))
	case frameSoft:
		return blank(w)
	}
	return lk.dim.Render(rep("─", w))
}

// status is the unified bar: context keys, then ? help, : palette, t theme, q quit.
func (lk *look) status(w int, cs chromeState, themeName string) string {
	global := []keyHint{{"?", "help"}, {":", "palette"}, {"t", "theme: " + themeName}, {"q", "quit"}}
	render := func(hints []keyHint) line {
		var l line
		for i, h := range hints {
			if i > 0 {
				l = append(l, spc(2))
				if lk.th.frame == frameCyber {
					l = append(l, seg{"// ", lk.dim})
				}
			}
			l = append(l, seg{h.key, lk.key}, spc(1), seg{lk.cap(h.desc), lk.sub})
		}
		return l
	}
	right := render(global)
	if cs.flash != "" {
		return lr(line{{cs.flash, lk.warn}}, right, w).render(w, false, lk.sel)
	}
	ctx := cs.context
	for n := len(ctx); n >= 0; n-- {
		left := render(ctx[:n])
		if left.width()+3+right.width() <= w {
			return lr(left, right, w).render(w, false, lk.sel)
		}
	}
	return lr(nil, right, w).render(w, false, lk.sel)
}
