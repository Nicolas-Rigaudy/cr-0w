package main

// Projects tab: a strip (scope, sort, the nearest dated target), the project list in three
// groups with two rows per project (facts, then the next action), the selected project's
// detail, and the roadmap and work timeline: per project its roadmap bar and the owner's
// commits per day under it. Landscape puts list and detail side by side, portrait stacks them.

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	jpList = iota
	jpDetail
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
	zoom bool // timeline over twice the weeks
}

var cardActions = [][2]string{{"c", "continue"}, {"r", "review"}, {"s", "research"}, {"e", "roadmap"}, {"o", "session"}, {"h", "handoff"}}

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

func (pv projectsView) healthStyle(h string) lipgloss.Style {
	switch h {
	case "late":
		return pv.lk.bad.Bold(true)
	case "at risk":
		return pv.lk.warn.Bold(true)
	case "on time":
		return pv.lk.ok
	}
	return pv.lk.dim
}

func (pv projectsView) barStyle(h string) lipgloss.Style {
	switch h {
	case "late":
		return pv.lk.bad
	case "at risk":
		return pv.lk.warn
	case "on time":
		return pv.lk.ok
	}
	return pv.lk.sub
}

func (pv projectsView) glyph(g string) seg {
	lk := pv.lk
	switch g {
	case "blocked":
		return seg{"■", lk.bad.Bold(true)}
	case "needs":
		return seg{"◐", lk.warn.Bold(true)}
	case "running":
		return seg{"●", lk.cyan.Bold(true)}
	case "quiet":
		return seg{"·", lk.dim}
	}
	return seg{"○", lk.sub}
}

var glyphWords = map[string]string{"blocked": "blocked", "needs": "needs you", "running": "running", "waiting": "waiting", "ready": "ready", "quiet": "quiet"}

// actor is the next action's prefix: nothing when the owner acts.
func (pv projectsView) actor(n *nextRow, long bool) line {
	lk := pv.lk
	switch {
	case n.Actor == "me" || n.Actor == "":
		return nil
	case n.Actor == "crow":
		return line{{"crow: ", lk.cyan.Bold(true)}}
	case long:
		return line{{"waiting on " + n.Actor + ": ", lk.sub.Bold(true)}}
	}
	return line{{"waiting " + n.Actor + ": ", lk.sub.Bold(true)}}
}

func (pv projectsView) whoLabel(actor string) seg {
	lk := pv.lk
	switch actor {
	case "me", "", "you":
		return seg{fit("you", 8), lk.warn.Bold(true)}
	case "crow":
		return seg{fit("crow", 8), lk.cyan.Bold(true)}
	}
	return seg{fit(actor, 8), lk.sub.Bold(true)}
}

func dateLabel(t time.Time) string { return t.Format("Jan 2") }

// launch is the nearest dated target ahead, work first: the strip and the ┊ line.
func (pv projectsView) launch() *projCard {
	var best *projCard
	for i, c := range pv.fx.Cards {
		if c.Target.IsZero() || c.Target.Before(day(pv.now)) {
			continue
		}
		if best == nil || (c.Scope == "work") != (best.Scope == "work") && c.Scope == "work" ||
			(c.Scope == "work") == (best.Scope == "work") && c.Target.Before(best.Target) {
			best = &pv.fx.Cards[i]
		}
	}
	return best
}

func (pv projectsView) strip(w int) string {
	lk := pv.lk
	left := line{{"m", lk.key}, {" scope ", lk.sub}}
	for i, md := range []string{"work", "perso", "all"} {
		if i > 0 {
			left = append(left, seg{" · ", lk.dim})
		}
		st := lk.sub
		if md == pv.mode {
			st = lk.acc
		}
		left = append(left, seg{md, st})
	}
	left = append(left, seg{"   sort ", lk.sub}, seg{"health", lk.bold})
	var right line
	if l := pv.launch(); l != nil {
		right = append(right, seg{"◆ " + strings.ToLower(l.Mile) + " " + dateLabel(l.Target) + " · " + l.left(pv.now), lk.bad.Bold(true)}, spc(3))
	}
	late, risk := 0, 0
	v := pv.visible()
	for _, c := range v {
		late += btoi(c.Health == "late")
		risk += btoi(c.Health == "at risk")
	}
	if late > 0 {
		right = append(right, seg{fmt.Sprintf("%d late", late), lk.bad}, spc(3))
	}
	if risk > 0 {
		right = append(right, seg{fmt.Sprintf("%d at risk", risk), lk.warn}, spc(3))
	}
	right = append(right, seg{fmt.Sprintf("%d projects", len(v)), lk.sub})
	return lr(line{spc(1)}.join(left), right.join(line{spc(1)}), w).render(w, false, lk.sel)
}

func (l line) join(r line) line { return append(append(line{}, l...), r...) }

// ----- list -----

type listCols struct{ repo, health, road, left, sig, next int }

func (pv projectsView) cols(port bool) listCols {
	if port {
		return listCols{repo: 4, health: 30, left: 41, next: 52}
	}
	return listCols{repo: 4, health: 33, road: 46, left: 67, sig: 80}
}

func (pv projectsView) listSpec(iw int, port bool) boxSpec {
	lk := pv.lk
	cards := pv.visible()
	sp := boxSpec{title: "Projects " + itoa(len(cards)), index: 1, sel: -1, focused: pv.pane == jpList,
		right: line{{"sorted by health, then target", lk.sub}}}
	if pv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if len(cards) == 0 {
		sp.rows = []line{{{"No active project in this mode.", lk.sub}}}
		return sp
	}
	c := pv.cols(port)
	at := func(l line, x int) line { return append(l, spc(x-l.width())) }
	hd := at(nil, c.repo).join(line{{lk.cap("repo"), lk.head}})
	hd = at(hd, c.health).join(line{{lk.cap("health"), lk.head}})
	if port {
		hd = at(hd, c.left).join(line{{lk.cap("left"), lk.head}})
		hd = at(hd, c.next).join(line{{lk.cap("next"), lk.head}})
	} else {
		hd = at(hd, c.road).join(line{{lk.cap("roadmap"), lk.head}})
		hd = at(hd, c.left).join(line{{lk.cap("left"), lk.head}})
		hd = at(hd, c.sig).join(line{{lk.cap("signals"), lk.head}})
	}
	sp.rows = append(sp.rows, hd)
	if !port {
		sp.rows = append(sp.rows, nil)
	}
	group := -1
	for i, cd := range cards {
		if cd.Group != group {
			if group >= 0 && !port {
				sp.rows = append(sp.rows, nil)
			}
			group = cd.Group
			n := 0
			for _, x := range cards {
				n += btoi(x.Group == group)
			}
			gst := []lipgloss.Style{lk.warn.Bold(true), lk.cyan.Bold(true), lk.dim.Bold(true)}[group]
			name := lk.cap(projGroupNames[group])
			head := line{{name, gst}, {fmt.Sprintf("  %d ", n), lk.sub}}
			sp.rows = append(sp.rows, head.join(line{{rep("╌", iw-head.width()), lk.dim}}))
		}
		if i == pv.sel {
			sp.sel = len(sp.rows)
			sp.selN = 1 + btoi(!port && cd.Group != pgQuiet)
		}
		sp.rows = append(sp.rows, pv.listRows(cd, i == pv.sel, iw, port)...)
	}
	return sp
}

func (pv projectsView) listRows(cd projCard, selected bool, iw int, port bool) []line {
	lk := pv.lk
	c := pv.cols(port)
	quiet := cd.Group == pgQuiet
	at := func(l line, x int) line { return append(l, spc(x-l.width())) }
	mark := seg{"  ", lk.plain}
	if selected {
		mark = seg{"› ", lk.acc}
	}
	nameSt := lk.bold
	if quiet {
		nameSt = lk.dim
	}
	row := line{mark, pv.glyph(cd.Glyph)}
	row = at(row, c.repo).join(line{{fit(cd.Name, c.health-c.repo-2), nameSt}})
	row = at(row, c.health).join(line{{cd.Health, pv.healthStyle(cd.Health)}})
	if port {
		row = at(row, c.left).join(line{{cd.left(pv.now), pv.healthStyle(cd.Health)}})
		if n := cd.lead(); n != nil && !quiet {
			row = at(row, c.next).join(pv.actor(n, false)).join(line{{n.What, lk.sub}})
		}
		return []line{row.trunc(iw)}
	}
	if cd.Of > 0 {
		row = at(row, c.road).join(gaugeBar(lk, cd.Done*100/cd.Of, 10, pv.barStyle(cd.Health))).join(line{{fmt.Sprintf(" %d/%d", cd.Done, cd.Of), lk.sub}})
	} else if !quiet && cd.Target.IsZero() {
		row = at(row, c.road).join(line{{"no roadmap", lk.dim}})
	}
	row = at(row, c.left).join(line{{cd.left(pv.now), pv.healthStyle(cd.Health)}})
	row = at(row, c.sig).join(pv.signals(cd))
	rows := []line{row.trunc(iw)}
	if quiet {
		return rows
	}
	two := at(nil, c.repo).join(line{{"→ ", lk.dim}})
	switch n := cd.lead(); {
	case n != nil:
		two = two.join(pv.actor(n, true)).join(line{{n.What, lk.plain}})
	case cd.Now != "":
		two = two.join(line{{"crow: ", lk.cyan.Bold(true)}, {cd.Now, lk.plain}})
	default:
		two = two.join(line{{"nobody assigned, ready to pick", lk.dim}})
	}
	two = two.trunc(c.sig - 2)
	if cd.Mile != "" {
		mile := line{{"◆ ", lk.dim}, {cd.Mile, lk.sub}}
		if !cd.Target.IsZero() {
			mile = append(mile, seg{"  " + dateLabel(cd.Target), lk.dim})
		}
		two = at(two, c.sig).join(mile)
	}
	return append(rows, two.trunc(iw))
}

func (pv projectsView) signals(cd projCard) line {
	lk := pv.lk
	var parts []line
	if n := len(cd.MRs); n > 0 {
		bad := 0
		for _, m := range cd.MRs {
			bad += btoi(m.Bad)
		}
		s := fmt.Sprintf("%d MR", n)
		st := lk.sub
		if bad > 0 {
			s += fmt.Sprintf(" %d red", bad)
			st = lk.bad
		}
		parts = append(parts, line{{s, st}})
	}
	var g []string
	if cd.Dirty > 0 {
		g = append(g, fmt.Sprintf("dirty %d", cd.Dirty))
	}
	if cd.Ahead > 0 {
		g = append(g, fmt.Sprintf("+%d", cd.Ahead))
	}
	if cd.Stale > 0 {
		g = append(g, fmt.Sprintf("%d stale", cd.Stale))
	}
	if len(g) > 0 {
		parts = append(parts, line{{strings.Join(g, " "), lk.warn}})
	}
	var out line
	for i, p := range parts {
		if i > 0 {
			out = append(out, seg{"  ·  ", lk.dim})
		}
		out = append(out, p...)
	}
	return out
}

// ----- detail -----

func (pv projectsView) secHead(title string, right string, iw int) line {
	lk := pv.lk
	l := line{{lk.cap(title), lk.head.Bold(true)}, spc(1)}
	r := line{}
	if right != "" {
		r = line{spc(1), {right, lk.sub}}
	}
	return l.join(line{{rep("╌", iw-l.width()-r.width()), lk.dim}}).join(r)
}

// detailSpec fits the detail into n rows: every next action when they fit, else as many as the
// room allows (6 in landscape, 2 in portrait at least) and a "+N more" line.
func (pv projectsView) detailSpec(iw, n int, port bool) boxSpec {
	base := 6
	if port {
		base = 2
	}
	all := 0
	if cd := pv.cur(); cd != nil {
		all = len(cd.Next)
	}
	keep := all
	if over := pv.detailNeedK(iw, port, all) - n; over > 0 {
		keep = max(base, all-over-1)
	}
	return pv.detailWith(iw, n, port, keep)
}

func (pv projectsView) detailWith(iw, n int, port bool, keep int) boxSpec {
	lk := pv.lk
	cd := pv.cur()
	sp := boxSpec{title: "Detail", index: 2, sel: -1, focused: pv.pane == jpDetail}
	if cd == nil {
		sp.rows = []line{{{"No project selected.", lk.sub}}}
		return sp
	}
	sp.title = "Detail: " + cd.Name
	right := cd.Scope + " · " + cd.Forge
	if l := cd.left(pv.now); l != "" {
		right += " · " + l
	}
	sp.right = line{{right, lk.sub}}
	var rows []line
	g := pv.glyph(cd.Glyph)
	state := line{g, spc(1), {glyphWords[cd.Glyph], g.style}, {"   health ", lk.sub}, {cd.Health, pv.healthStyle(cd.Health)}}
	if !port {
		state = append(state, seg{"   " + cd.Scope + " · " + cd.Forge, lk.sub})
	}
	rows = append(rows, state)
	if cd.Now != "" && !port {
		rows = append(rows, line{{"● ", lk.cyan}, {"crow session ", lk.bold}, {cd.Now, lk.sub}, {"  " + cd.NowSince, lk.dim}})
	}
	if cd.Block != "" {
		rows = append(rows, line{{"■ ", lk.bad}, {"blocked by " + cd.Block, lk.bad}})
	}
	gap := func() {
		if !port {
			rows = append(rows, nil)
		}
	}
	gap()
	rows = append(rows, pv.secHead("Next action", "who acts", iw))
	next := cd.Next
	if l := cd.lead(); l != nil && len(next) > 0 && next[0] != *l {
		next = append([]nextRow{*l}, next...)
		for i := 1; i < len(next); i++ {
			if next[i] == *l {
				next = append(next[:i], next[i+1:]...)
				break
			}
		}
	}
	more := len(next) - keep
	if more > 0 {
		next = next[:keep]
	}
	for _, x := range next {
		meta := strings.TrimSpace(strings.Trim(x.Why+" · "+x.Age, " ·"))
		rows = append(rows, line{pv.whoLabel(x.Actor), {x.What, lk.plain}, {"  " + meta, lk.dim}})
	}
	if len(next) == 0 {
		rows = append(rows, line{{"nothing open in the threads", lk.dim}})
	}
	if more > 0 {
		rows = append(rows, line{{fmt.Sprintf("        +%d more in the open threads", more), lk.dim}})
	}
	gap()
	if cd.Mile != "" || cd.Of > 0 {
		r := ""
		if !cd.Target.IsZero() {
			r = dateLabel(cd.Target) + " · " + cd.left(pv.now)
		}
		rows = append(rows, pv.secHead("Milestone", r, iw))
		ml := line{{"◆ ", lk.acc}, {cmpOr(cd.Mile, "all milestones done"), lk.bold}}
		if cd.Of > 0 {
			ml = append(ml, seg{fmt.Sprintf("   %d of %d done", cd.Done, cd.Of), lk.sub})
		}
		rows = append(rows, ml)
		if cd.Exit != "" {
			rows = append(rows, line{{"exit: ", lk.sub}, {cd.Exit, lk.plain}})
		}
		if !port && cd.Of > 0 {
			for _, m := range cd.Miles {
				mark := seg{"○ ", lk.sub}
				if m.Done {
					mark = seg{"✓ ", lk.ok}
				}
				rows = append(rows, line{mark, {m.Title, lk.plain}, {"  " + dateLabel(m.Date), lk.dim}})
			}
		}
		gap()
	}
	if len(cd.MRs) > 0 {
		rows = append(rows, pv.secHead("Open MRs", "who must act", iw))
		mrs := cd.MRs
		if port && len(mrs) > 2 {
			mrs = mrs[:2]
		}
		for _, m := range mrs {
			st := lk.sub
			if m.Bad {
				st = lk.bad
			}
			who := lk.sub.Bold(true)
			if m.Acts == "you" {
				who = lk.warn.Bold(true)
			}
			state := m.State
			if !port {
				state = fit(state, 26)
			}
			rows = append(rows, line{{m.Ref + " ", lk.acc}, {state, st}, {"  acts: ", lk.sub}, {fit(m.Acts, 10), who}, {m.Todo, lk.sub}})
		}
		gap()
	}
	if !port {
		rows = append(rows, pv.secHead("Git", "", iw))
		if s := pv.signals(projCard{Dirty: cd.Dirty, Ahead: cd.Ahead, Stale: cd.Stale}); len(s) > 0 {
			rows = append(rows, line{{"dirty, ahead, stale trees: ", lk.sub}}.join(s))
		} else {
			rows = append(rows, line{{"clean", lk.sub}})
		}
	}
	// the action keys sit on the last row of the box
	act := line{}
	for i, a := range cardActions {
		if i > 0 {
			act = append(act, spc(3))
		}
		act = append(act, seg{a[0], lk.key}, seg{" " + a[1], lk.sub})
	}
	if len(rows) > n-1 {
		rows = rows[:max(0, n-1)]
	}
	for len(rows) < n-1 {
		rows = append(rows, nil)
	}
	sp.rows = append(rows, act)
	return sp
}

// timelineMax is the timeline's height in rows with a gap between projects.
func (pv projectsView) timelineMax() int {
	active := 0
	for _, cd := range pv.visible() {
		active += btoi(cd.Group != pgQuiet)
	}
	return 5 + active*(2+chartMax)
}

// detailNeed is the detail pane's natural height in rows: its content, a blank, the actions.
func (pv projectsView) detailNeed(iw int, port bool) int {
	base := 6
	if port {
		base = 2
	}
	return pv.detailNeedK(iw, port, base)
}

func (pv projectsView) detailNeedK(iw int, port bool, keep int) int {
	rows := pv.detailWith(iw, 200, port, keep).rows
	n := len(rows) - 1
	for n > 0 && len(rows[n-1]) == 0 {
		n--
	}
	return n + 2
}

// ----- timeline -----

// window: the Monday some weeks back, and the days shown (past weeks then planned weeks).
func (pv projectsView) window(port bool) (time.Time, int, int) {
	past, ahead := 5, 9
	if port {
		past, ahead = 3, 6
	}
	if pv.zoom {
		past, ahead = past*2, ahead*2
	}
	t := day(pv.now)
	t = t.AddDate(0, 0, -((int(t.Weekday())+6)%7)-7*past)
	return t, 7 * (past + ahead), past
}

// grid is one timeline row: a rune and a style per cell.
type grid struct {
	r  []rune
	st []lipgloss.Style
}

func newGrid(w int, st lipgloss.Style) *grid {
	g := &grid{r: make([]rune, w), st: make([]lipgloss.Style, w)}
	for i := range g.r {
		g.r[i], g.st[i] = ' ', st
	}
	return g
}

func (g *grid) put(x int, s string, st lipgloss.Style) {
	for _, r := range s {
		if x >= 0 && x < len(g.r) {
			g.r[x], g.st[x] = r, st
		}
		x++
	}
}

// free says whether n cells from x are blank, so labels never overwrite each other.
func (g *grid) free(x, n int) bool {
	for i := x; i < x+n; i++ {
		if i >= 0 && i < len(g.r) && g.r[i] != ' ' {
			return false
		}
	}
	return true
}

func (g *grid) line() line {
	var out line
	for i := range g.r {
		out = append(out, seg{string(g.r[i]), g.st[i]})
	}
	return out
}

// barCells draws v against top as a bar of hh rows, top row first, in eighths of a cell.
// Any value above zero shows at least one eighth.
func barCells(v, top, hh int) []string {
	out := make([]string, hh)
	e := 0
	if top > 0 && v > 0 {
		e = max(1, (min(v, top)*hh*8+top/2)/top)
	}
	for r := range out {
		fill := e - (hh-1-r)*8
		switch {
		case fill >= 8:
			out[r] = "█"
		case fill <= 0:
			out[r] = " "
		default:
			out[r] = string([]rune("▁▂▃▄▅▆▇")[fill-1])
		}
	}
	return out
}

// commitTop: the most commits any shown project made in a day since start, so every
// project's chart shares one y scale.
func (pv projectsView) commitTop(start time.Time) int {
	top := 0
	for _, cd := range pv.visible() {
		if cd.Group == pgQuiet {
			continue
		}
		perDay := map[string]int{}
		for ds, n := range cd.Work {
			d := parseDay(ds, pv.now.Location())
			if d.IsZero() || d.Before(start) || d.After(day(pv.now)) {
				continue
			}
			perDay[ds] += n
			top = max(top, perDay[ds])
		}
	}
	return top
}

func weekend(d time.Time) bool { return d.Weekday() == time.Saturday || d.Weekday() == time.Sunday }

// chartMax is the tallest commit chart, in rows.
const chartMax = 3

// timelineFit: the commit chart height and the gap between projects that fit in avail
// inner rows (avail 0: no limit). The chart shrinks before the gap goes.
func timelineFit(projects, avail int) (hh, gap int) {
	if avail <= 0 {
		return chartMax, 1
	}
	for _, f := range [][2]int{{3, 1}, {3, 0}, {2, 0}} {
		if 5+projects*(1+f[0])+(projects-1)*f[1] <= avail {
			return f[0], f[1]
		}
	}
	return 1, 0
}

// yLabels: the commits-a-day value at the top of each chart row, the bottom row 0.
func yLabels(top, hh int) []string {
	out := make([]string, hh)
	for r := range out {
		out[r] = itoa(top * (hh - 1 - r) / max(1, hh-1))
	}
	if hh == 1 {
		out[0] = itoa(top)
	}
	for r := hh - 2; r >= 0; r-- { // a value repeated from the row below says nothing
		if out[r] == itoa(top*(hh-2-r)/max(1, hh-1)) {
			out[r] = ""
		}
	}
	return out
}

func (pv projectsView) timelineSpec(iw int, port bool, avail int) boxSpec {
	lk := pv.lk
	sp := boxSpec{title: "Roadmap and work", index: 3, sel: -1, focused: pv.pane == jpTime,
		right: line{{"█ done  ░ left  ◆ milestone  │ today  ┊ target", lk.sub}}}
	if pv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	start, days, past := pv.window(port)
	gut := 28
	if port {
		gut = 26
	}
	mw := max(10, iw-gut)
	col := func(t time.Time) int { return int(day(t).Sub(start).Hours() / 24 * float64(mw) / float64(days)) }
	end := start.AddDate(0, 0, days)
	today := col(pv.now)
	launch := -1
	l := pv.launch()
	if l != nil && l.Target.Before(end) {
		launch = col(l.Target)
	}
	// axis: months and the launch label, then Monday dates and today, then the rule
	months, dates, wdays, rule := newGrid(mw, lk.sub), newGrid(mw, lk.dim), newGrid(mw, lk.dim), newGrid(mw, lk.dim)
	if launch >= 0 && l != nil {
		lbl := "◆ " + dateLabel(l.Target) + " " + strings.ToLower(l.Mile) + " " + l.left(pv.now)
		months.put(max(0, min(launch, mw-len([]rune(lbl)))), lbl, lk.bad.Bold(true))
	}
	dates.put(today, "today", lk.acc)
	rule.put(today, "▼", lk.acc)
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		x := col(d)
		if (d.Day() == 1 || d.Equal(start)) && months.free(x-1, 5) {
			months.put(x, d.Format("Jan"), lk.sub)
		}
		// one weekday initial per day, weekends dim
		wst := lk.sub
		if weekend(d) {
			wst = lk.dim
		}
		narrow := mw < 2*days // under 2 cells a day: Mondays and weekends only
		if wdays.r[x] == ' ' && (!narrow || weekend(d) || d.Weekday() == time.Monday) {
			wdays.put(x, d.Weekday().String()[:1], wst)
		}
		if d.Weekday() == time.Monday {
			if dates.free(x-1, 4) {
				dates.put(x, d.Format("02"), lk.dim)
			}
			if rule.r[x] == ' ' {
				rule.put(x, "┬", lk.dim)
			}
		}
	}
	for x := range rule.r {
		if rule.r[x] == ' ' {
			rule.r[x] = '─'
		}
	}
	pad := func(s string, st lipgloss.Style) line { return line{{fit(s, gut), st}} }
	sp.rows = append(sp.rows, pad("", lk.sub).join(months.line()),
		pad(fmt.Sprintf("last %d weeks worked", past), lk.sub).join(dates.line()),
		pad(fmt.Sprintf("next %d weeks planned", days/7-past), lk.sub).join(wdays.line()),
		pad("", lk.sub).join(rule.line()))
	guides := func(g *grid) {
		if today >= 0 && today < mw && g.r[today] == ' ' {
			g.put(today, "│", lk.dim)
		}
		if launch >= 0 && launch < mw && g.r[launch] == ' ' {
			g.put(launch, "┊", lk.bad)
		}
	}
	tot := map[string]int{}
	active := 0
	for _, cd := range pv.visible() {
		active += btoi(cd.Group != pgQuiet)
	}
	hh, gapRows := timelineFit(active, avail)
	top := pv.commitTop(start)
	first := true
	for _, cd := range pv.visible() {
		if cd.Group == pgQuiet {
			continue
		}
		if !first {
			for i := 0; i < gapRows; i++ {
				sp.rows = append(sp.rows, nil)
			}
		}
		first = false
		bar := newGrid(mw, lk.plain)
		bs := pv.barStyle(cd.Health)
		switch {
		case cd.Start.IsZero() && cd.Target.IsZero():
			bar.put(0, "· · ·  no date", lk.dim)
		default:
			s, e := cd.Start, cd.Target
			openEnd := e.IsZero() || !e.Before(end)
			x0, x1 := 0, mw-1
			if !s.IsZero() && s.After(start) {
				x0 = col(s)
			}
			if !openEnd {
				x1 = col(e)
			}
			doneTo := x0
			if cd.Of > 0 {
				doneTo = x0 + (x1-x0)*cd.Done/cd.Of
			}
			for x := x0; x < x1; x++ {
				if x < doneTo {
					bar.put(x, "█", bs)
				} else {
					bar.put(x, "░", bs)
				}
			}
			for _, m := range cd.Miles {
				if m.Date.After(start) && m.Date.Before(end) {
					bar.put(col(m.Date), "◆", lk.plain)
				}
			}
			if s.IsZero() || !s.After(start) {
				bar.put(0, "<", lk.sub)
			}
			if openEnd {
				bar.put(mw-1, ">", lk.sub)
			} else if !port {
				bar.put(x1+2, dateLabel(e), lk.sub)
			}
		}
		guides(bar)
		sp.rows = append(sp.rows, line{pv.glyph(cd.Glyph), spc(1), {fit(cd.Name, gut-2), lk.bold}}.join(bar.line()))
		// commits a day as a bar chart: y labelled in the gutter, x on the shared date axis
		perDay := map[int]int{}
		sum, last := 0, ""
		for ds, n := range cd.Work {
			d := parseDay(ds, pv.now.Location())
			if d.IsZero() || d.After(day(pv.now)) {
				continue
			}
			if ds > last {
				last = ds
			}
			if d.Before(start) {
				continue
			}
			k := int(d.Sub(start).Hours() / 24)
			perDay[k] += n
			sum += n
			tot[cd.Name] += n
		}
		chart := make([]*grid, hh)
		for r := range chart {
			chart[r] = newGrid(mw, lk.cyan)
		}
		for k, n := range perDay {
			d := start.AddDate(0, 0, k)
			x0, x1 := col(d), max(col(d)+1, col(d.AddDate(0, 0, 1))-1)
			for r, g := range barCells(n, top, hh) {
				if g != " " {
					chart[r].put(x0, rep(g, x1-x0), lk.cyan)
				}
			}
		}
		note := "not touched"
		if last != "" {
			if d := parseDay(last, pv.now.Location()); d.Equal(day(pv.now)) {
				note = "last today"
			} else {
				note = "last " + dateLabel(d)
			}
		}
		left := []string{"commits a day", fmt.Sprintf("%d in %d weeks", sum, past), note}
		if hh < 3 {
			when := strings.TrimPrefix(note, "last ")
			left = []string{"commits a day", fmt.Sprintf("%d in %dw · %s", sum, past, when)}
		}
		yl := yLabels(top, hh)
		const lab = 6 // "  12 ┤"
		for d := start; !d.After(day(pv.now)); d = d.AddDate(0, 0, 1) {
			if weekend(d) {
				x0, x1 := col(d), max(col(d)+1, col(d.AddDate(0, 0, 1))-1)
				for x := x0; x < x1; x++ {
					if chart[hh-1].r[x] == ' ' {
						chart[hh-1].put(x, "┈", lk.dim)
					}
				}
			}
		}
		for r, g := range chart {
			guides(g)
			txt := ""
			if r < len(left) {
				txt = left[r]
			}
			sty := lk.dim
			if r == 0 {
				sty = lk.sub
			}
			sp.rows = append(sp.rows, line{{fit("  "+txt, gut-lab), sty}, {fmt.Sprintf("%4s ┤", yl[r]), lk.dim}}.join(g.line()))
		}
	}
	sp.rows = append(sp.rows, pv.timeSplit(tot, port))
	return sp
}

// timeSplit is where the owner's commits went over the window.
func (pv projectsView) timeSplit(tot map[string]int, port bool) line {
	lk := pv.lk
	type kv struct {
		k string
		v int
	}
	var all []kv
	sum := 0
	for k, v := range tot {
		all = append(all, kv{k, v})
		sum += v
	}
	out := line{{lk.cap("time split") + "  ", lk.head.Bold(true)}}
	if sum == 0 {
		return append(out, seg{"no commits in the window", lk.dim})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].v > all[j].v || all[i].v == all[j].v && all[i].k < all[j].k })
	top := 6
	if port {
		top = 4
	}
	for i, x := range all[:min(top, len(all))] {
		if i > 0 {
			out = append(out, seg{"  ·  ", lk.dim})
		}
		out = append(out, seg{x.k, lk.plain}, seg{fmt.Sprintf(" %d%%", (100*x.v+sum/2)/sum), lk.cyan})
	}
	return out
}

// ----- composition -----

func (pv projectsView) body(w, h int) []string {
	lk := pv.lk
	var out block
	switch layoutFor(w) {
	case laySingle:
		out = pv.single(w, h)
	case layLandscape:
		out = pv.landscape(w, h-1)
	default:
		out = pv.portrait(w, h-1)
	}
	if layoutFor(w) != laySingle {
		out = append(block{pv.strip(w)}, out...)
	}
	_ = lk
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func (pv projectsView) gapV() int {
	if pv.lk.th.frame == frameSoft {
		return 1
	}
	return 0
}

func (pv projectsView) landscape(w, h int) block {
	lk := pv.lk
	gv, gh := pv.gapV(), 0
	if lk.th.frame == frameSoft {
		gh = 3
	}
	lw := w * 128 / 235
	dw := w - lw - gh
	ls := pv.listSpec(lk.innerW(lw), false)
	// the list shows every row, the timeline then takes up to 4-row commit bars, and the
	// top boxes keep whatever is left
	need := max(len(ls.rows), pv.detailNeed(lk.innerW(dw), false), 14)
	natural := lk.boxH(len(pv.timelineSpec(lk.innerW(w), false, 0).rows))
	topH := max(lk.boxH(need), h-gv-lk.boxH(pv.timelineMax()))
	topH = max(min(topH, h-gv-natural), lk.boxH(14))
	th := h - topH - gv
	tl := pv.timelineSpec(lk.innerW(w), false, th-lk.boxH(0))
	list := lk.box(ls, lw, topH)
	det := lk.box(pv.detailSpec(lk.innerW(dw), topH-lk.boxH(0), false), dw, topH)
	return vcat(hcat([]int{lw, dw}, gh, list, det), blankLines(w, gv), lk.box(tl, w, th))
}

func (pv projectsView) portrait(w, h int) block {
	lk := pv.lk
	gv := pv.gapV()
	iw := lk.innerW(w)
	ls := pv.listSpec(iw, true)
	lh := min(lk.boxH(len(ls.rows)), h/3)
	natural := lk.boxH(len(pv.timelineSpec(iw, true, 0).rows))
	th := max(min(lk.boxH(pv.timelineMax()), h-lh-2*gv-lk.boxH(pv.detailNeed(iw, true))), natural)
	th = min(th, h-lh-2*gv-lk.boxH(8))
	dh := h - lh - th - 2*gv
	tl := pv.timelineSpec(iw, true, th-lk.boxH(0))
	return vcat(lk.box(ls, w, lh), blankLines(w, gv), lk.box(pv.detailSpec(iw, dh-lk.boxH(0), true), w, dh),
		blankLines(w, gv), lk.box(tl, w, th))
}

func (pv projectsView) single(w, h int) []string {
	lk := pv.lk
	iw := lk.innerW(w)
	var sp boxSpec
	switch pv.pane {
	case jpDetail:
		sp = pv.detailSpec(iw, h-lk.boxH(0), true)
	case jpTime:
		sp = pv.timelineSpec(iw, true, h-lk.boxH(0))
	default:
		sp = pv.listSpec(iw, true)
	}
	sp.bottom = fmt.Sprintf("pane %d of %d · tab", pv.pane+1, projPanes)
	return lk.box(sp, w, h)
}

func scopeTag(scope string) string {
	if scope == "personal" {
		return "perso"
	}
	return scope
}
