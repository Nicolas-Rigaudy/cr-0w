package main

// Home page: Usage, Right now, Blockers, Today, Next, Projects and roadmap.
// Landscape and portrait share the board grid with different ratios; under 100 columns
// one pane is shown at a time and `tab` switches.

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	layLandscape = iota // 150+ columns
	layPortrait         // 100 to 149
	laySingle           // under 100
)

func layoutFor(w int) int {
	switch {
	case w >= 150:
		return layLandscape
	case w >= 100:
		return layPortrait
	}
	return laySingle
}

const (
	pUsage = iota
	pNow
	pBlockers
	pToday
	pNext
	pProjects
	homePanes
)

var homePaneNames = [homePanes]string{"Usage", "Right now", "Blockers", "Today", "Next", "Projects and roadmap"}

// homeView is the data and focus a Home render needs.
type homeView struct {
	lk    *look
	fx    *fixtures
	mode  string
	sync  string
	focus int
	sel   [homePanes]int
}

func (hv homeView) loading() bool { return hv.sync == stLoading }

// nextRows filters the Next list by mode.
func (hv homeView) nextRows() []nextRow {
	var out []nextRow
	for _, r := range hv.fx.Next {
		if hv.mode == "all" || (hv.mode == "perso") == r.Perso {
			out = append(out, r)
		}
	}
	return out
}

func (hv homeView) paneLen(i int) int {
	if hv.loading() {
		return 1
	}
	switch i {
	case pNow:
		return max(1, len(hv.fx.Agents))
	case pBlockers:
		return max(1, len(hv.fx.Blockers))
	case pToday:
		return max(1, len(hv.fx.Today)+len(hv.fx.Mail))
	case pNext:
		return max(1, len(hv.nextRows()))
	case pProjects:
		return max(1, len(hv.fx.Projects)+1)
	}
	return len(hv.fx.Gauges)
}

// hasHdr: panes with a table header row above the list.
func (hv homeView) hasHdr(i int) bool {
	if hv.loading() {
		return false
	}
	switch i {
	case pNow:
		return len(hv.fx.Agents) > 0
	case pBlockers:
		return len(hv.fx.Blockers) > 0 && hv.sync != stError
	}
	return false
}

func (hv homeView) spec(i, iw int) boxSpec {
	lk, fx := hv.lk, hv.fx
	sp := boxSpec{title: homePaneNames[i], index: i + 1, sel: -1, focused: hv.focus == i}
	if hv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	link := func(s string) line { return line{{s, lk.label}} }
	switch i {
	case pUsage:
		sp.right = link("→ 5 Usage")
		sp.rows = hv.usageRows(iw)
	case pNow:
		sp.right = link("→ 2 Agents")
		if len(fx.Agents) == 0 {
			sp.rows = []line{{{"No agent running.", lk.sub}}}
			break
		}
		sp.hdr = lk.headRow(iw, []col{{"agent", 12, 2}, {"account", 8, 1}, {"doing", 0, 0}}, "since")
		for _, a := range fx.Agents {
			sp.rows = append(sp.rows, hv.agentRow(a, iw))
		}
	case pBlockers:
		if n := len(fx.Blockers); n > 0 {
			sp.title += " " + strconv.Itoa(n)
		}
		if !fx.Live || fx.MRsReady {
			sp.right = line{{"MRs: ", lk.sub}, {fmt.Sprintf("%d awaits you", fx.MRsAwaiting), lk.warn}, {fmt.Sprintf(" · %d mine · %d thread", fx.MRsMine, fx.MRsThreads), lk.sub}}
		}
		switch {
		case hv.sync == stError:
			sp.rows = []line{{{"◆ forge sync failed: login expired, log in again", lk.bad}}}
		case len(fx.Blockers) == 0 && fx.Live && !fx.MRsReady:
			sp.rows = []line{{{"Waiting for the forges (MRs tab).", lk.sub}}}
		case len(fx.Blockers) == 0 && fx.Live:
			sp.rows = []line{{{"No MR blocks you. Ticket blockers are not wired.", lk.sub}}}
		case len(fx.Blockers) == 0:
			sp.rows = []line{{{"Nothing blocks you.", lk.sub}}}
		default:
			sp.hdr = lk.headRow(iw, []col{{"item", 14, 1}, {"why it blocks", 0, 0}}, "age")
			for _, b := range fx.Blockers {
				sp.rows = append(sp.rows, hv.blockerRow(b, iw))
			}
		}
	case pToday:
		for _, t := range fx.Today {
			tag := line{}
			if t.Tag != "" {
				st := lk.ok
				if t.TagState == "warn" {
					st = lk.warn
				}
				tag = line{{t.Tag, st}}
			}
			sp.rows = append(sp.rows, lr(line{{fit(t.When, 5), lk.cyan}, spc(1), {t.What, lk.plain}}, tag, iw))
		}
		for _, m := range fx.Mail {
			sp.rows = append(sp.rows, line{{m, lk.sub}})
		}
		if len(sp.rows) == 0 && fx.Live {
			sp.rows = []line{{{"Calendar and mail are not connected.", lk.sub}}}
		} else if len(sp.rows) == 0 {
			sp.rows = []line{{{"Nothing scheduled.", lk.sub}}}
		}
	case pNext:
		rows := hv.nextRows()
		if len(rows) > 0 {
			sp.title += " " + strconv.Itoa(len(rows))
		}
		sp.right = line{{map[string]string{"work": "work · m shows perso", "perso": "perso · m shows all", "all": "all · m shows work"}[hv.mode], lk.sub}}
		for _, r := range rows {
			l := line{{fit(r.Project, 8), lk.key}, {r.What, lk.plain}}
			var tag line
			if r.Perso && hv.mode == "all" {
				tag = line{{"perso", lk.sub}}
			}
			sp.rows = append(sp.rows, lr(l, tag, iw))
		}
		if len(rows) == 0 {
			sp.rows = []line{{{"Nothing queued.", lk.sub}}}
		}
	case pProjects:
		sp.right = link("→ 4 Projects")
		for _, p := range fx.Projects {
			date, ds := "no date", lk.sub
			if p.Date != "" {
				date, ds = p.Date, lk.sub
				if p.DateAtRisk {
					ds = lk.bad
				}
			}
			l := line{{fit(p.Name, 8), lk.bold}}
			l = append(l, gaugeBar(lk, p.Done*100/p.Of, 8, lk.ok)...)
			l = append(l, sp1(), seg{fmt.Sprintf("%d/%d", p.Done, p.Of), lk.plain}, sp1(), seg{date, ds})
			sp.rows = append(sp.rows, l)
		}
		if len(fx.Projects) == 0 {
			sp.rows = []line{{{"No roadmap yet.", lk.sub}}}
			break
		}
		sp.rows = append(sp.rows, line{{fmt.Sprintf("%d repos", fx.Repos), lk.sub}, spc(1), {fmt.Sprintf("%d dirty", fx.Dirty), lk.warn}, spc(1), {fmt.Sprintf("%d ahead", fx.Ahead), lk.warn}})
	}
	return sp
}

func sp1() seg { return spc(1) }

type col struct {
	name  string
	width int // 0: flexible
	gap   int
}

// headRow: table header above a list; the first column starts after the glyph cell.
func (lk *look) headRow(iw int, cols []col, right string) line {
	prefix := "  "
	if lk.th.frame == frameCyber {
		prefix = "//"
	}
	l := line{{prefix, lk.head}}
	for _, c := range cols {
		text := lk.cap(c.name)
		if c.width > 0 {
			text = fit(text, c.width)
		}
		l = append(l, seg{text + rep(" ", c.gap), lk.head})
	}
	return lr(l, line{{lk.cap(right), lk.head}}, iw)
}

func (hv homeView) usageRows(iw int) []line {
	lk := hv.lk
	cells := max(6, min(16, (iw-37)/2))
	var out []line
	for _, g := range hv.fx.Gauges {
		l := line{{fit(g.Account, 7), lk.level(g.FiveH)}, spc(1), {"5h ", lk.sub}}
		l = append(l, gaugeBar(lk, g.FiveH, cells, lk.level(g.FiveH))...)
		txt := fmt.Sprintf(" %d%%", g.FiveH)
		if g.Reset != "" {
			txt += " resets " + g.Reset
		}
		l = append(l, seg{txt, lk.plain})
		if g.Week >= 0 {
			l = append(l, spc(2), seg{"wk ", lk.sub})
			l = append(l, gaugeBar(lk, g.Week, cells, lk.level(g.Week))...)
			l = append(l, seg{fmt.Sprintf(" %d%%", g.Week), lk.plain})
		}
		if g.Policy != "" {
			l = append(l, spc(2), seg{fmt.Sprintf("%s · %s · runs %d/%d", g.Policy, g.Burn, g.RunsUsed, g.RunsMax), lk.sub})
		}
		out = append(out, l)
	}
	return out
}

func (hv homeView) agentRow(a agentRow, iw int) line {
	lk := hv.lk
	var g seg
	switch a.State {
	case "working":
		g = seg{"●", lk.ok}
	case "idle":
		g = seg{"○", lk.sub}
	case "needs":
		g = seg{"◆", lk.mag}
	case "waiting":
		g = seg{"‖", lk.warn}
	default:
		g = seg{"✓", lk.ok}
	}
	left := line{g, spc(1), {fit(a.Name, 12), lk.bold}, spc(2), {fit(a.Account, 8), lk.sub}, spc(1), {a.Doing, lk.plain}}
	return lr(left, line{{a.Since, lk.sub}}, iw)
}

func (hv homeView) blockerRow(b blocker, iw int) line {
	lk := hv.lk
	g := seg{"◆", lk.warn}
	switch b.Glyph {
	case "blocked":
		g = seg{"◆", lk.bad}
	case "waiting":
		g = seg{"‖", lk.warn}
	}
	left := line{g, spc(1), {fit(b.Item, 14), lk.bold}, spc(1), {b.Why, lk.plain}}
	return lr(left, line{{b.Age, lk.sub}}, iw)
}

// ----- composition -----

type band struct {
	panes    []int
	min, max int // body lines (header included)
	ratio    [2]int
	prio     int // lower grows first
}

func (hv homeView) bands(lay int) []band {
	blockRatio, nextRatio := [2]int{16, 10}, [2]int{13, 10}
	if lay == layPortrait {
		blockRatio = [2]int{17, 10}
	}
	bs := []band{
		{panes: []int{pUsage}, min: 2, prio: 9},
		{panes: []int{pNow}, min: 6, prio: 0},
		{panes: []int{pBlockers, pToday}, min: 6, ratio: blockRatio, prio: 1},
		{panes: []int{pNext, pProjects}, min: 6, ratio: nextRatio, prio: 2},
	}
	for i := range bs {
		for _, p := range bs[i].panes {
			l := hv.paneLen(p) + btoi(hv.hasHdr(p))
			bs[i].max = max(bs[i].max, l)
		}
		bs[i].min = min(bs[i].min, bs[i].max)
	}
	return bs
}

func (hv homeView) gapH() int {
	if hv.lk.th.frame == frameSoft {
		return 3
	}
	return 0
}

func (hv homeView) gapV() int {
	if hv.lk.th.frame == frameSoft {
		return 1
	}
	return 0
}

// body draws the page at w x h.
func (hv homeView) body(w, h int) []string {
	lay := layoutFor(w)
	lk := hv.lk
	if lay == laySingle {
		return hv.single(w, h)
	}
	bs := hv.bands(lay)
	lines := make([]int, len(bs))
	total := func() int {
		t := 0
		for i := range bs {
			t += lk.boxH(lines[i])
		}
		return t + hv.gapV()*(len(bs)-1)
	}
	for i, b := range bs {
		lines[i] = b.min
	}
	// Leftover height goes to more rows, in priority order, up to the data available.
	for grew := true; grew && total() < h; {
		grew = false
		for p := 0; p <= 9 && total() < h; p++ {
			for i, b := range bs {
				if b.prio == p && lines[i] < b.max && total() < h {
					lines[i]++
					grew = true
				}
			}
		}
	}
	// Too tall: shrink the lowest priority bands first, down to header plus one row.
	for total() > h {
		shrunk := false
		for p := 9; p >= 0 && total() > h; p-- {
			for i, b := range bs {
				if b.prio == p && lines[i] > 2 && total() > h {
					lines[i]--
					shrunk = true
				}
			}
		}
		if !shrunk {
			break
		}
	}
	var out block
	for i, b := range bs {
		if i > 0 && hv.gapV() > 0 {
			out = append(out, blank(w))
		}
		bh := lk.boxH(lines[i])
		if len(b.panes) == 1 {
			out = append(out, hv.pane(b.panes[0], w, bh, "")...)
			continue
		}
		gap := hv.gapH()
		w1 := (w - gap) * b.ratio[0] / (b.ratio[0] + b.ratio[1])
		w2 := w - gap - w1
		out = append(out, hcat([]int{w1, w2}, gap, hv.pane(b.panes[0], w1, bh, ""), hv.pane(b.panes[1], w2, bh, ""))...)
	}
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

// selectable: the pane holds data rows (not a note line) the cursor can sit on.
func (hv homeView) selectable(i int) bool {
	if hv.loading() {
		return false
	}
	switch i {
	case pNow:
		return len(hv.fx.Agents) > 0
	case pBlockers:
		return len(hv.fx.Blockers) > 0 && hv.sync != stError
	case pToday:
		return len(hv.fx.Today)+len(hv.fx.Mail) > 0
	case pNext:
		return len(hv.nextRows()) > 0
	case pProjects:
		return len(hv.fx.Projects) > 0
	}
	return false
}

// pane draws one pane at w x h with its selection and counters; pos is the single-pane
// position label.
func (hv homeView) pane(i, w, h int, pos string) block {
	lk := hv.lk
	sp := hv.spec(i, lk.innerW(w))
	var bottom []string
	if hv.focus == i && hv.selectable(i) {
		sp.sel = min(hv.sel[i], len(sp.rows)-1)
		if i == pBlockers && len(sp.rows) >= 5 {
			bottom = append(bottom, fmt.Sprintf("%d of %d", sp.sel+1, len(sp.rows)))
		}
	}
	if i == pNext && hv.selectable(i) {
		bottom = append(bottom, fmt.Sprintf("top %d", min(len(sp.rows), max(0, h-lk.boxH(0)))))
	}
	if pos != "" {
		bottom = append(bottom, pos)
	}
	sp.bottom = strings.Join(bottom, " · ")
	return lk.box(sp, w, h)
}

func (hv homeView) single(w, h int) []string {
	return hv.pane(hv.focus, w, h, fmt.Sprintf("pane %d of %d · tab", hv.focus+1, homePanes))
}
