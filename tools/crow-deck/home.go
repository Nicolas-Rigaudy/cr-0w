package main

// Home page, direction A: Pulse across the top, then one Next list grouped by who acts
// (You decide, Crow works, Waiting on others), the Projects table and Done today on the
// left, with Right now and Today on a rail at 150+ columns. Under 100 columns one pane is
// shown at a time and `tab` switches. Panes without data are left out.

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
	pPulse = iota
	pNext
	pProjects
	pDone
	pNow
	pToday
	homePanes
)

var homePaneNames = [homePanes]string{"Pulse", "Next", "Projects", "Done today", "Right now", "Today"}

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

func (hv homeView) inMode(perso bool) bool {
	return hv.mode == "all" || (hv.mode == "perso") == perso
}

// nextRows filters the open threads by mode.
func (hv homeView) nextRows() []nextRow {
	var out []nextRow
	for _, r := range hv.fx.Next {
		if hv.inMode(r.Perso) {
			out = append(out, r)
		}
	}
	return out
}

type nextGroup struct {
	name  string
	glyph seg
	rows  []nextRow
}

// groups: MRs and the owner's threads, then crow's, then the ones waiting on someone.
// Empty groups are left out.
func (hv homeView) groups() []nextGroup {
	lk := hv.lk
	g := []nextGroup{{"You decide", seg{"●", lk.warn}, nil}, {"Crow works", seg{"◆", lk.cyan}, nil}, {"Waiting on others", seg{"○", lk.sub}, nil}}
	var nm nextGroup
	if len(hv.fx.NoModel) > 0 {
		// Every account is blocked: work that needs no model comes first, crow's own waits.
		nm = nextGroup{"No tokens, work without a model", seg{"●", lk.warn}, nil}
		if hv.fx.NoModelUntil != "" {
			nm.name = "No tokens until " + hv.fx.NoModelUntil + ", work without a model"
		}
		for _, r := range hv.fx.NoModel {
			if hv.inMode(r.Perso) {
				nm.rows = append(nm.rows, r)
			}
		}
		g[1] = nextGroup{"Crow works, paused", seg{"◆", lk.dim}, nil}
	}
	for _, r := range hv.fx.Decide {
		if hv.inMode(r.Perso) {
			g[0].rows = append(g[0].rows, r)
		}
	}
	for _, r := range hv.nextRows() {
		switch r.Actor {
		case "", "me":
			g[0].rows = append(g[0].rows, r)
		case "crow":
			g[1].rows = append(g[1].rows, r)
		default:
			g[2].rows = append(g[2].rows, r)
		}
	}
	var out []nextGroup
	if len(nm.rows) > 0 {
		out = append(out, nm)
	}
	for _, x := range g {
		if len(x.rows) > 0 {
			out = append(out, x)
		}
	}
	return out
}

func (hv homeView) nextItems() []nextRow {
	var out []nextRow
	for _, g := range hv.groups() {
		out = append(out, g.rows...)
	}
	return out
}

func (hv homeView) projects() []projRow {
	var out []projRow
	for _, p := range hv.fx.Projects {
		if hv.inMode(p.Scope == "personal") {
			out = append(out, p)
		}
	}
	return out
}

func (hv homeView) hasToday() bool { return len(hv.fx.Today)+len(hv.fx.Mail) > 0 }

// items is the number of rows the cursor walks in pane i.
func (hv homeView) items(i int) int {
	if hv.loading() {
		return 0
	}
	switch i {
	case pNext:
		return len(hv.nextItems())
	case pProjects:
		return len(hv.projects())
	case pDone:
		return len(hv.fx.Done)
	case pNow:
		return len(hv.fx.Agents)
	case pToday:
		return len(hv.mails())
	}
	return 0
}

func (hv homeView) selectable(i int) bool { return hv.items(i) > 0 }

// panes lists the panes the layout shows, in tab order.
func (hv homeView) panes(w int) []int {
	out := []int{pNext, pPulse, pProjects}
	if hv.loading() || len(hv.fx.Done) > 0 {
		out = append(out, pDone)
	}
	if layoutFor(w) != layPortrait {
		out = append(out, pNow)
		if hv.hasToday() {
			out = append(out, pToday)
		}
	}
	return out
}

// nextFocus steps the focus through the shown panes.
func (hv homeView) nextFocus(w, d int) int {
	ps := hv.panes(w)
	at := 0
	for i, p := range ps {
		if p == hv.focus {
			at = i
		}
	}
	return ps[(at+d+len(ps))%len(ps)]
}

// ----- pane content -----

// display: the rows of pane i at inner width iw, with the display row each item starts on.
func (hv homeView) display(i, iw int) (sp boxSpec, at []int) {
	lk, fx := hv.lk, hv.fx
	sp = boxSpec{title: homePaneNames[i], index: i + 1, sel: -1, focused: hv.focus == i}
	if hv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp, nil
	}
	link := func(s string) line { return line{{s, lk.label}} }
	add := func(l line) {
		at = append(at, len(sp.rows))
		sp.rows = append(sp.rows, l)
	}
	switch i {
	case pPulse:
		sp.right = link("→ 5 Usage")
		sp.rows = hv.pulseRows(iw)
		if len(sp.rows) == 0 {
			sp.rows = []line{{{"No usage data yet.", lk.sub}}}
		}
	case pNext:
		sp.right = line{{"me · crow · waiting", lk.sub}}
		if hv.sync == stError {
			sp.rows = append(sp.rows, line{{"◆ forge sync failed: login expired, log in again", lk.bad}})
		}
		gs := hv.groups()
		n := 0
		for k, g := range gs {
			n += len(g.rows)
			sp.rows = append(sp.rows, lk.groupHead(g.name, g.glyph.style.Bold(true), len(g.rows), iw))
			if k == 0 {
				sp.rows = append(sp.rows, hv.nextHead(iw))
				if r := lk.headRule(iw); r != nil {
					sp.rows = append(sp.rows, r)
				}
			}
			for _, r := range g.rows {
				add(hv.nextRow(g, r, iw))
			}
		}
		if n > 0 {
			sp.title += " " + strconv.Itoa(n)
		} else {
			sp.rows = append(sp.rows, line{{"Nothing queued.", lk.sub}})
		}
	case pProjects:
		sp.right = link("→ 4 Projects")
		ps := hv.projects()
		if len(ps) == 0 {
			sp.rows = []line{{{"No project in this mode.", lk.sub}}}
			break
		}
		sp.hdr = hv.projHead(iw)
		for _, p := range ps {
			add(hv.projRow(p, iw))
		}
	case pDone:
		sp.title += " " + strconv.Itoa(len(fx.Done))
		sp.right = link("activity log")
		tw := 8
		if layoutFor(iw+4) == layLandscape {
			tw = 22
		}
		for _, d := range fx.Done {
			add(line{{"✓ ", lk.ok}, {d.When + "  ", lk.sub}, {fit(d.Tag, tw), lk.sub}, spc(1), {d.What, lk.plain}})
		}
	case pNow:
		sp.right = link("→ 2 Agents")
		sp.selN = 2
		if len(fx.Agents) == 0 {
			sp.rows = []line{{{"No agent running.", lk.sub}}}
			sp.selN = 0
			break
		}
		for _, a := range fx.Agents {
			l1, l2 := hv.agentLines(a, iw)
			add(l1)
			sp.rows = append(sp.rows, l2)
		}
	case pToday:
		for _, t := range fx.Today {
			if t.Kind == "flag" {
				add(hv.todayLine(t, iw)) // the cursor walks the mail rows only
			} else {
				sp.rows = append(sp.rows, hv.todayLine(t, iw))
			}
		}
		if len(fx.Mail) > 0 {
			sp.right = line{{strings.Join(fx.Mail, " · "), lk.sub}}
		}
	}
	return sp, at
}

// mails are the Today rows the cursor walks.
func (hv homeView) mails() []todayRow {
	var out []todayRow
	for _, t := range hv.fx.Today {
		if t.Kind == "flag" {
			out = append(out, t)
		}
	}
	return out
}

// curMail is the selected mail when the Today pane has focus.
func (hv homeView) curMail() *todayRow {
	ms := hv.mails()
	if hv.focus != pToday || len(ms) == 0 {
		return nil
	}
	r := ms[max(0, min(len(ms)-1, hv.sel[pToday]))]
	return &r
}

// todayLine draws one rail row: an event, a section heading, a flagged mail or a gap.
func (hv homeView) todayLine(t todayRow, iw int) line {
	lk := hv.lk
	switch t.Kind {
	case "gap":
		return line{}
	case "head":
		return line{{t.What, lk.label}}
	case "flag":
		sym, mark, what := "● ", lk.cyan, lk.bold
		switch t.TagState {
		case "warn":
			mark = lk.bad
		case "crow":
			sym = "◆ "
		case "past":
			mark, what = lk.sub, lk.plain
		}
		return lr(line{{sym, mark}, {t.What, what}}, line{{t.Tag, lk.sub}}, iw)
	}
	when, what, tag := lk.cyan, lk.plain, lk.ok
	switch t.TagState {
	case "next":
		when, what, tag = lk.acc.Bold(true), lk.bold, lk.warn.Bold(true)
	case "past":
		when, what = lk.sub, lk.sub
	case "warn":
		tag = lk.warn
	}
	return lr(line{{fit(t.When, 11), when}, spc(2), {t.What, what}}, line{{t.Tag, tag}}, iw)
}

// pulseRows: the owner's accounts in full, the team accounts on one compact line.
func (hv homeView) pulseRows(iw int) []line {
	lk := hv.lk
	wide := iw >= 180
	gw := 14
	if wide {
		gw = 22
	}
	var out []line
	for _, g := range hv.fx.Gauges {
		name := lk.plain
		if g.Policy == "use" {
			name = lk.acc.Bold(true)
		}
		l := line{{fit(g.Account, 8), name}, {"5h ", lk.sub}}
		l = append(l, gaugeBar(lk, g.FiveH, gw, lk.level(g.FiveH))...)
		l = append(l, seg{fmt.Sprintf(" %3d%% left", g.FiveH), lk.level(g.FiveH)})
		if g.Reset != "" {
			l = append(l, seg{" resets " + g.Reset, lk.sub})
		}
		if g.Week >= 0 {
			l = append(l, spc(3), seg{"week ", lk.sub})
			l = append(l, gaugeBar(lk, g.Week, gw, lk.level(g.Week))...)
			l = append(l, seg{fmt.Sprintf(" %3d%% left", g.Week), lk.level(g.Week)})
		}
		if wide {
			if g.WeekReset != "" {
				l = append(l, seg{" resets " + fit(g.WeekReset, 12), lk.sub})
			}
			l = append(l, spc(2), seg{fit(g.Policy, 10), lk.cyan})
			if g.Burn != "" {
				l = append(l, seg{g.Burn, lk.sub})
			}
		}
		switch {
		case g.Empty != "":
			l = append(l, spc(3), seg{"empty at " + g.Empty + " at this pace", lk.bad.Bold(true)})
		case g.Week >= 0 && g.Week < 5:
			l = append(l, spc(3), seg{"week nearly spent: crow skips it", lk.warn})
		}
		out = append(out, l)
	}
	if len(hv.fx.Team) > 0 {
		tw := 6
		if wide {
			tw = 10
		}
		l := line{{fit("team", 8), lk.head}}
		for k, g := range hv.fx.Team {
			if k > 0 {
				l = append(l, spc(4))
			}
			l = append(l, seg{fit(g.Account, 8), lk.sub}, seg{"5h ", lk.dim})
			l = append(l, gaugeBar(lk, g.FiveH, tw, lk.level(g.FiveH))...)
			l = append(l, seg{fmt.Sprintf(" %3d%% left  ", g.FiveH), lk.sub}, seg{"wk ", lk.dim})
			l = append(l, gaugeBar(lk, g.Week, tw, lk.level(max(0, g.Week)))...)
			l = append(l, seg{fmt.Sprintf(" %3d%% left", max(0, g.Week)), lk.sub})
		}
		if wide {
			l = append(l, spc(6), seg{"crow never uses these", lk.dim})
		}
		out = append(out, l)
	}
	return out
}

func projW(iw int) int {
	if iw < 90 {
		return 12
	}
	return 22
}

func (hv homeView) nextHead(iw int) line {
	lk := hv.lk
	left := line{{"  " + fit(lk.cap("project"), projW(iw)) + " ", lk.head}, {lk.cap("what"), lk.head}}
	return lr(left, line{{fit(lk.cap("why / link"), 15), lk.head}, {fmt.Sprintf("%6s", lk.cap("since")), lk.head}}, iw)
}

func (hv homeView) nextRow(g nextGroup, r nextRow, iw int) line {
	lk := hv.lk
	what := r.What
	if r.Actor != "" && r.Actor != "me" && r.Actor != "crow" && r.Actor != "others" {
		what = r.Actor + ": " + what
	}
	age := lk.sub
	if n, err := strconv.Atoi(strings.TrimSuffix(r.Age, "d")); err == nil && strings.HasSuffix(r.Age, "d") && n >= 7 {
		age = lk.bad
	}
	left := line{g.glyph, spc(1), {fit(cmpOr(r.Repo, r.Project), projW(iw)), lk.sub}, spc(1), {what, lk.plain}}
	return lr(left, line{spc(2), {fit(r.Why, 13), lk.sub}, {fmt.Sprintf("%6s", r.Age), age}}, iw)
}

// projCols: name, scope, roadmap bar and count, then target, left and git on the right.
func projCols(iw int) (name, bar int, dates bool, git int) {
	name, bar, git = 24, 10, 19
	if iw < 90 {
		name, bar = 14, 6
	}
	return name, bar, iw >= 140, git
}

func (hv homeView) projHead(iw int) line {
	lk := hv.lk
	nw, bw, dates, gw := projCols(iw)
	left := line{{fit(lk.cap("repo"), nw) + " ", lk.head}, {fit(lk.cap("scope"), 9), lk.head}, {fit(lk.cap("roadmap"), bw+7), lk.head}, {lk.cap("now"), lk.head}}
	var right line
	if dates {
		right = line{{fit(lk.cap("target"), 8), lk.head}, {fmt.Sprintf("%5s", lk.cap("left")), lk.head}}
	}
	return lr(left, append(right, seg{fmt.Sprintf("%*s", gw, lk.cap("git")), lk.head}), iw)
}

func (hv homeView) projRow(p projRow, iw int) line {
	lk := hv.lk
	nw, bw, dates, gw := projCols(iw)
	l := line{{fit(p.Name, nw), lk.bold}, spc(1), {fit(p.Scope, 9), lk.sub}}
	if p.Of > 0 {
		l = append(l, gaugeBar(lk, p.Done*100/p.Of, bw, lk.ok)...)
		l = append(l, seg{fmt.Sprintf(" %5s ", fmt.Sprintf("%d/%d", p.Done, p.Of)), lk.sub})
	} else {
		l = append(l, spc(bw+7))
	}
	l = append(l, seg{p.Now, lk.plain})
	var right line
	if dates {
		ds := lk.warn
		if p.DateAtRisk {
			ds = lk.bad
		}
		right = line{spc(2), {fit(p.Date, 8), ds}, {fmt.Sprintf("%5s", p.Left), lk.warn}}
	}
	var parts []string
	if p.Dirty > 0 {
		parts = append(parts, fmt.Sprintf("%d dirty", p.Dirty))
	}
	if p.Ahead > 0 {
		parts = append(parts, fmt.Sprintf("%d ahead", p.Ahead))
	}
	git, gs := "clean", lk.sub
	if len(parts) > 0 {
		git, gs = strings.Join(parts, " · "), lk.warn
	}
	if len(git) > gw {
		git = git[:gw-1] + "…"
	}
	return lr(l, append(right, seg{fmt.Sprintf("%*s", gw, git), gs}), iw)
}

func (hv homeView) agentLines(a agentRow, iw int) (line, line) {
	lk := hv.lk
	g, state := seg{"✓", lk.ok}, seg{a.Since, lk.sub}
	switch a.State {
	case "working":
		g, state = seg{"●", lk.ok}, seg{"working", lk.ok}
	case "idle":
		g = seg{"○", lk.sub}
	case "needs":
		g, state = seg{"◆", lk.mag}, seg{"needs you", lk.mag}
	case "waiting":
		g, state = seg{"‖", lk.warn}, seg{"paused", lk.warn}
	}
	l1 := lr(line{g, spc(1), {a.Name, lk.bold}, spc(2), {a.Account, lk.sub}}, line{state}, iw)
	return l1, line{spc(2), {a.Doing, lk.sub}}
}

func (hv homeView) spec(i, iw int) boxSpec {
	sp, _ := hv.display(i, iw)
	return sp
}

// ----- composition -----

func (hv homeView) gapH() int { h, _ := hv.lk.gaps(); return h }

func (hv homeView) gapV() int { _, v := hv.lk.gaps(); return v }

// want is the box height that shows all of pane i at width w.
func (hv homeView) want(i, w int) int {
	sp := hv.spec(i, hv.lk.innerW(w))
	return hv.lk.boxH(len(sp.rows) + hv.lk.hdrH(sp))
}

// stack draws panes top to bottom in w x h. Each pane takes what it wants in order;
// a pane that would not fit at its floor is left out, and the space under the last one
// stays blank. floor gives the smallest useful height of a pane.
func (hv homeView) stack(ps []int, w, h int, floor func(int) int) block {
	lk := hv.lk
	hs := make([]int, len(ps))
	left := h
	for k, p := range ps {
		gap := btoi(k > 0) * hv.gapV()
		rest := 0
		for _, q := range ps[k+1:] {
			rest += floor(q) + hv.gapV()
		}
		hs[k] = min(hv.want(p, w), max(floor(p), left-gap-rest))
		if hs[k] > left-gap || hs[k] < lk.boxH(1) {
			hs[k] = 0
			continue
		}
		left -= hs[k] + gap
	}
	var out block
	for k, p := range ps {
		if hs[k] == 0 {
			continue
		}
		if len(out) > 0 && hv.gapV() > 0 {
			out = append(out, blank(w))
		}
		out = append(out, hv.pane(p, w, hs[k], "")...)
	}
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func (hv homeView) floor(p int) int {
	lk := hv.lk
	switch p {
	case pNext:
		return lk.boxH(min(6, hv.items(p)+3))
	case pProjects:
		return lk.boxH(min(12, hv.items(p)+1))
	case pDone:
		return lk.boxH(min(5, max(1, hv.items(p))))
	case pToday:
		return lk.boxH(min(3, max(1, hv.items(p))))
	case pNow:
		return lk.boxH(min(4, max(1, 2*hv.items(p))))
	}
	return 0
}

// body draws the page at w x h.
func (hv homeView) body(w, h int) []string {
	lay := layoutFor(w)
	if lay == laySingle {
		return hv.single(w, h)
	}
	pulse := hv.pane(pPulse, w, min(h, hv.want(pPulse, w)), "")
	out := append(block{}, pulse...)
	if hv.gapV() > 0 {
		out = append(out, blank(w))
	}
	rest := max(0, h-len(out))
	left := []int{pNext, pProjects}
	if hv.loading() || len(hv.fx.Done) > 0 {
		left = append(left, pDone)
	}
	if lay == layPortrait {
		return append(out, hv.stack(left, w, rest, hv.floor)...)[:h]
	}
	gap := hv.gapH()
	lw := (w - gap) * 68 / 100
	rw := w - gap - lw
	rail := []int{pNow}
	if hv.hasToday() {
		rail = append(rail, pToday)
	}
	railFloor := func(p int) int {
		if p == pNow && len(rail) == 1 {
			return hv.want(pNow, rw)
		}
		return hv.floor(p)
	}
	out = append(out, hcat([]int{lw, rw}, gap, hv.stack(left, lw, rest, hv.floor), hv.stack(rail, rw, rest, railFloor))...)
	return out[:h]
}

// pane draws one pane at w x h with its selection; pos is the single-pane position label.
func (hv homeView) pane(i, w, h int, pos string) block {
	lk := hv.lk
	sp, at := hv.display(i, lk.innerW(w))
	var bottom []string
	slots := h - lk.boxH(0) - lk.hdrH(sp)
	if hv.focus == i && len(at) > 0 {
		k := max(0, min(hv.sel[i], len(at)-1))
		sp.sel = at[k]
		if len(sp.rows) > slots {
			bottom = append(bottom, fmt.Sprintf("%d of %d", k+1, len(at)))
		}
	} else {
		sp.selN = 0
		if len(at) > 0 && len(sp.rows) > slots {
			bottom = append(bottom, fmt.Sprintf("%d of %d", min(len(at), countBelow(at, slots)), len(at)))
		}
	}
	if pos != "" {
		bottom = append(bottom, pos)
	}
	sp.bottom = strings.Join(bottom, " · ")
	return lk.box(sp, w, h)
}

// countBelow: items whose first row is inside the first n display rows.
func countBelow(at []int, n int) int {
	c := 0
	for _, a := range at {
		if a < n {
			c++
		}
	}
	return c
}

func (hv homeView) single(w, h int) []string {
	ps := hv.panes(w)
	pos := 0
	for k, p := range ps {
		if p == hv.focus {
			pos = k
		}
	}
	return hv.pane(ps[pos], w, h, fmt.Sprintf("pane %d of %d · tab", pos+1, len(ps)))
}
