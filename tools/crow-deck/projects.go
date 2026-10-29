package main

// Projects tab: one card per project (NOW, NEXT, BLOCK, ROAD, MILE, OPEN), the action keys
// on the selected card, a timeline of the projects that have a target date, and the repos
// with no roadmap on one line. Read only: the action keys say what arrives later.

import (
	"fmt"
	"strings"
	"time"
)

const (
	jpCards = iota
	jpTime
	projPanes
)

type projectsView struct {
	lk   *look
	fx   *fixtures
	now  time.Time
	sync string
	mode string
	sel  int
	pane int
	zoom bool // timeline over three months instead of five
}

func (pv projectsView) loading() bool { return pv.sync == stLoading }

func (pv projectsView) visible() []projCard {
	var out []projCard
	for _, c := range pv.fx.Cards {
		if inMode(pv.mode, c.Scope) {
			out = append(out, c)
		}
	}
	return out
}

func (pv projectsView) cur() *projCard {
	v := pv.visible()
	if pv.sel < 0 || pv.sel >= len(v) {
		return nil
	}
	return &v[pv.sel]
}

func (pv projectsView) stateStyle(s string) seg {
	lk := pv.lk
	switch s {
	case "blocked":
		return seg{s, lk.bad}
	case "needs you":
		return seg{s, lk.warn}
	case "running":
		return seg{s, lk.ok}
	}
	return seg{s, lk.sub}
}

var cardActions = [][2]string{{"c", "continue"}, {"r", "review"}, {"s", "research"}, {"e", "roadmap"}, {"o", "session"}}

func (pv projectsView) cardRows(c projCard, selected bool, iw int) []line {
	lk := pv.lk
	lab := func(s string) seg { return seg{fit(s, 6), lk.label} }
	title := line{{c.Name, lk.bold}}
	if c.Release != "" {
		title = append(title, seg{" · " + c.Release, lk.sub})
	}
	tag := line{pv.stateStyle(c.State), {" · " + scopeTag(c.Scope), lk.sub}}
	rows := []line{lr(title, tag, iw)}
	if c.Now != "" {
		rows = append(rows, lr(line{lab("NOW"), {"● ", lk.ok}, {c.Now, lk.plain}}, line{{c.NowSince, lk.sub}}, iw))
	} else {
		rows = append(rows, line{lab("NOW"), {"idle", lk.sub}})
	}
	if c.Next != "" {
		rows = append(rows, line{lab("NEXT"), {c.Next, lk.plain}})
	}
	if c.Block != "" {
		rows = append(rows, line{lab("BLOCK"), {c.Block, lk.warn}})
	}
	if c.Of > 0 {
		st := lk.sub
		switch {
		case c.AtRisk:
			st = lk.bad
		case !c.Target.IsZero():
			st = lk.ok
		}
		road := line{lab("ROAD")}
		road = append(road, gaugeBar(lk, c.Done*100/c.Of, 12, lk.ok)...)
		road = append(road, seg{" " + c.roadText(), st})
		rows = append(rows, road)
		mile := line{lab("MILE"), {fmt.Sprintf("✓ %d done", c.Done), lk.ok}}
		if c.Current != "" {
			mile = append(mile, seg{" · ● " + c.Current, lk.plain})
		}
		if c.Exit != "" {
			mile = append(mile, seg{" · exit: " + c.Exit, lk.sub})
		}
		rows = append(rows, mile)
	}
	if c.Open != "" {
		rows = append(rows, line{lab("OPEN"), {c.Open, lk.plain}})
	}
	if selected {
		act := line{lab("")}
		for _, a := range cardActions {
			act = append(act, seg{a[0], lk.key}, seg{" " + a[1] + "  ", lk.sub})
		}
		rows = append(rows, act)
	}
	return rows
}

func (pv projectsView) cardsSpec(iw int) boxSpec {
	lk := pv.lk
	sp := boxSpec{title: "Projects", index: 1, sel: -1, focused: pv.pane == jpCards}
	if pv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	cards := pv.visible()
	if len(cards) == 0 {
		sp.rows = []line{{{"No project has a roadmap or activity in this mode.", lk.sub}}}
	}
	for i, c := range cards {
		if i > 0 {
			sp.rows = append(sp.rows, nil)
		}
		if i == pv.sel {
			sp.sel = len(sp.rows)
		}
		sp.rows = append(sp.rows, pv.cardRows(c, i == pv.sel, iw)...)
	}
	sp.right = line{{fmt.Sprintf("%d with a roadmap", pv.roadmaps()), lk.sub}}
	return sp
}

func (pv projectsView) roadmaps() int {
	n := 0
	for _, c := range pv.visible() {
		if c.Of > 0 {
			n++
		}
	}
	return n
}

// quietRows is the one line of projects with no roadmap and no activity, under the timeline.
func (pv projectsView) quietRows() []line {
	if len(pv.fx.Quiet) == 0 {
		return nil
	}
	return []line{nil, {{"no roadmap  ", pv.lk.sub}, {strings.Join(pv.fx.Quiet, " · "), pv.lk.dim}}}
}

// ----- timeline -----

// window is the first day of the month today is in, and the months shown.
func (pv projectsView) window() (time.Time, int) {
	n := time.Date(pv.now.Year(), pv.now.Month(), 1, 0, 0, 0, 0, pv.now.Location())
	return n, map[bool]int{false: 5, true: 3}[pv.zoom]
}

func (pv projectsView) timelineSpec(iw int) boxSpec {
	lk := pv.lk
	sp := boxSpec{title: "Timeline", index: 2, sel: -1, focused: pv.pane == jpTime}
	sp.right = line{{map[bool]string{false: "5", true: "3"}[pv.zoom] + " months · z zoom", lk.sub}}
	if pv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	var dated []projCard
	for _, c := range pv.visible() {
		if !c.Target.IsZero() {
			dated = append(dated, c)
		}
	}
	if len(dated) == 0 {
		sp.rows = append([]line{{{"No dated roadmap. Roadmap docs carry milestones, not target dates.", lk.sub}}}, pv.quietRows()...)
		return sp
	}
	start, months := pv.window()
	end := start.AddDate(0, months, 0)
	span := end.Sub(start).Hours() / 24
	nameW, tailW := 12, 9
	tw := max(20, iw-nameW-tailW-2)
	col := func(t time.Time) int {
		d := t.Sub(start).Hours() / 24
		return max(0, min(tw-1, int(d/span*float64(tw))))
	}
	ruler := make([]rune, tw)
	lab := make([]rune, tw)
	for i := range ruler {
		ruler[i], lab[i] = '·', ' '
	}
	for m := 0; m < months; m++ {
		ms := start.AddDate(0, m, 0)
		x := col(ms)
		ruler[x] = '│'
		for j, r := range []rune(ms.Format("Jan")) {
			if x+j < tw {
				lab[x+j] = r
			}
		}
	}
	today := col(pv.now)
	pad := spc(nameW)
	sp.rows = append(sp.rows,
		line{pad, {string(lab), lk.sub}},
		line{pad, {string(ruler), lk.dim}},
		line{pad, {rep(" ", today) + "▼ today", lk.acc}})
	for _, c := range dated {
		x := col(c.Target)
		row := line{{fit(c.Name, nameW), lk.bold}, {rep("━", today+1), lk.ok}}
		gap := max(0, x-today-1)
		st := lk.sub
		if c.AtRisk {
			st = lk.bad
		}
		if x >= today {
			row = append(row, seg{rep("╌", gap), st}, seg{"◆", st})
		}
		row = append(row, seg{" " + c.Target.Format("Jan 2"), st})
		sp.rows = append(sp.rows, row)
	}
	sp.rows = append(sp.rows, pv.quietRows()...)
	return sp
}

// ----- composition -----

func (pv projectsView) body(w, h int) []string {
	lk := pv.lk
	lay := layoutFor(w)
	if lay == laySingle {
		return pv.single(w, h)
	}
	gapV := 0
	if lk.th.frame == frameSoft {
		gapV = 1
	}
	dated := 0
	for _, c := range pv.visible() {
		if !c.Target.IsZero() {
			dated++
		}
	}
	th := lk.boxH(max(1, dated) + 3 + 2*btoi(len(pv.fx.Quiet) > 0))
	cardsH := max(lk.boxH(4), h-th-gapV)
	out := vcat(lk.box(pv.cardsSpec(lk.innerW(w)), w, cardsH), blankLines(w, gapV), lk.box(pv.timelineSpec(lk.innerW(w)), w, th))
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func (pv projectsView) single(w, h int) []string {
	lk := pv.lk
	iw := lk.innerW(w)
	sp := pv.cardsSpec(iw)
	if pv.pane == jpTime {
		sp = pv.timelineSpec(iw)
	}
	sp.bottom = fmt.Sprintf("pane %d of %d · tab", pv.pane+1, projPanes)
	b := lk.box(sp, w, h)
	for len(b) < h {
		b = append(b, blank(w))
	}
	return b[:h]
}

func scopeTag(scope string) string {
	if scope == "personal" {
		return "perso"
	}
	return scope
}
