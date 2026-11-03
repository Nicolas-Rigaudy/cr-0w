package main

// Usage tab: the runway of each account (5h window, activity, forecast, week reset), daily
// spend, what crow saved, the ledger of sessions, and the cost by project and by model.
// Costs are API-equivalent dollars from crow-usage, not what is billed.

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	upRunway = iota
	upSpend
	upSaved
	upLedger
	upProjects
	upModels
	usagePanes
)

// Ledger group fold states; absent means the default (top three groups open).
const (
	foldClosed = iota + 1
	foldOpen
	foldAll
)

type usageViewState struct {
	lk     *look
	fx     *fixtures
	sync   string
	now    time.Time
	pane   int
	sel    int // account on the runway
	period int
	led    int // ledger row
	prj    int // By project row
	fold   map[string]int
	flat   bool   // g: one list by cost instead of project groups
	filter string // ledger limited to one project
	tokens bool   // c: daily chart in tokens
	shift  int    // runway pan, hours
}

func (uv usageViewState) loading() bool {
	return uv.sync == stLoading || uv.fx.Live && uv.fx.Usage == nil && len(uv.fx.Gauges) == 0
}

func (uv usageViewState) per() periodUse {
	if u := uv.fx.Usage; u != nil {
		return u.Periods[uv.period]
	}
	return periodUse{}
}

// money is a chart label: $104, $9.9.
func money(v float64) string {
	if v >= 10 {
		return fmt.Sprintf("$%.0f", v)
	}
	if v == 0 {
		return "$0"
	}
	return fmt.Sprintf("$%.1f", v)
}

// usd is a table amount: $1,512, $118, $49.76.
func usd(v float64) string {
	switch {
	case v >= 1000:
		n := int(v + .5)
		return fmt.Sprintf("$%d,%03d", n/1000, n%1000)
	case v >= 100:
		return fmt.Sprintf("$%.0f", v)
	}
	return fmt.Sprintf("$%.2f", v)
}

func tokensShort(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 1000:
		return fmt.Sprintf("%.1fk", float64(n)/1e3)
	}
	return fmt.Sprint(n)
}

func pctText(v float64) string {
	if v < 1 {
		return "<1%"
	}
	return fmt.Sprintf("%.0f%%", v)
}

// hbar is a horizontal bar of frac on w cells, in eighths; zero draws nothing.
func hbar(frac float64, w int, st lipgloss.Style) line {
	e := int(math.Round(max(0, min(1, frac)) * float64(w) * 8))
	if frac > 0 && e == 0 {
		e = 1
	}
	s := rep("█", e/8)
	if e%8 > 0 {
		s += string([]rune("▏▎▍▌▋▊▉")[e%8-1])
	}
	return line{{s, st}, spc(w - (e+7)/8)}
}

func right(s string, w int) string { return fmt.Sprintf("%*s", w, s) }

func (uv usageViewState) empty(sp boxSpec, msg string) boxSpec {
	if uv.loading() {
		sp.rows = []line{{{"syncing…", uv.lk.sub}}}
	} else {
		sp.rows = []line{{{msg, uv.lk.sub}}}
	}
	return sp
}

// ----- runway -----

// runwayCols: gutter width, plot width, columns per hour and hours shown before now.
func (uv usageViewState) runwayCols(iw int, port bool) (g, p, cph, past int) {
	g = min(100, max(64, iw*42/100))
	if port {
		g = 17
	}
	p = max(10, iw-g-1)
	cph = 4
	hours := p / cph
	past = hours*14/27 + uv.shift
	return
}

func runwayRows(mine, team int) int {
	n := 2 + 3*mine + max(0, mine-1)
	if team > 0 {
		n += 1 + team
	}
	return n
}

func (uv usageViewState) runwaySpec(iw int, port bool) boxSpec {
	lk, fx := uv.lk, uv.fx
	g, p, cph, past := uv.runwayCols(iw, port)
	hours := p / cph
	sp := boxSpec{title: fmt.Sprintf("Runway · last %dh, next %dh", past, hours-past), index: 1, sel: -1, focused: uv.pane == upRunway,
		right: line{{"h l pan 6h · P policy", lk.sub}}, bottom: "━ used  ─ left  ╌ forecast  ┃ week reset"}
	if uv.loading() || len(fx.Gauges)+len(fx.Team) == 0 {
		return uv.empty(sp, "No account reports limits.")
	}
	now := uv.now
	t0 := now.Truncate(time.Hour).Add(-time.Duration(past) * time.Hour)
	col := func(t time.Time) int { return int(math.Floor(t.Sub(t0).Hours() * float64(cph))) }
	nowC := col(now)
	// ruler: hour labels every 3h, ticks under them, the now marker.
	labs, ticks := newGrid(p, lk.sub), newGrid(p, lk.dim)
	ticks.put(0, rep("─", p), lk.dim)
	for h := 0; h <= hours; h++ {
		t := t0.Add(time.Duration(h) * time.Hour)
		c := h * cph
		if t.Hour()%3 != 0 || c >= p {
			continue
		}
		lab, st := t.Format("15"), lk.sub
		if t.Hour() == 0 {
			lab, st = t.Format("Mon 2"), lk.acc
		}
		if c-len(lab)/2 >= 0 && labs.free(c-len(lab)/2-1, len(lab)+2) && abs(c-nowC) > 3 {
			labs.put(c-len(lab)/2, lab, st)
		}
		ticks.put(c, "┴", lk.dim)
	}
	if nowC >= 0 && nowC < p {
		labs.put(nowC-1, "now", lk.acc)
		ticks.put(nowC, "▼", lk.acc)
	}
	pad := line{spc(g + 1)}
	sp.rows = append(sp.rows, pad.join(labs.line()), pad.join(ticks.line()))
	// peak activity over every lane, so lanes share one scale
	act := map[string]map[int64]int{}
	if fx.Usage != nil {
		act = fx.Usage.Activity
	}
	colTok := func(acct string) []int {
		out := make([]int, p)
		for b, tok := range act[acct] {
			if c := col(time.Unix(b, 0)); c >= 0 && c < p && c <= nowC {
				out[c] += tok
			}
		}
		return out
	}
	peak := 0
	lanes := map[string][]int{}
	for _, gg := range fx.Gauges {
		lanes[gg.Account] = colTok(gg.Account)
		for _, v := range lanes[gg.Account] {
			peak = max(peak, v)
		}
	}
	for i, gg := range fx.Gauges {
		if i > 0 {
			sp.rows = append(sp.rows, line{{rep("╌", iw), lk.dim}})
		}
		gut := uv.gutter(gg, g, port)
		if i == uv.sel && uv.pane == upRunway {
			gut[0] = append(line{{"› ", lk.acc}}, gut[0][1:]...)
		}
		w5, wa, wk := newGrid(p, lk.dim), newGrid(p, lk.dim), newGrid(p, lk.dim)
		// row 1: the current 5h window, used then left
		if !gg.ResetAt.IsZero() {
			a, b := col(gg.ResetAt.Add(-5*time.Hour)), col(gg.ResetAt)
			cut := a + int(math.Round(float64(b-a)*float64(100-gg.FiveH)/100))
			for c := a; c <= b; c++ {
				switch {
				case c == a:
					w5.put(c, "├", lk.dim)
				case c == b:
					w5.put(c, "┤", lk.dim)
				case c <= cut:
					w5.put(c, "━", lk.level(gg.FiveH))
				default:
					w5.put(c, "─", lk.dim)
				}
			}
		}
		// row 2: tokens per column so far, then the forecast to the reset
		for c, v := range lanes[gg.Account] {
			if v > 0 {
				wa.put(c, barCells(v, peak, 1)[0], lk.acc)
			}
		}
		switch {
		case uv.emptiesFirst(gg) && !gg.ResetAt.IsZero():
			e := col(clockOn(now, gg.Empty))
			wa.put(nowC+1, rep("╌", max(1, e-nowC-1)), lk.bad)
			wa.put(max(e, nowC+2)+1, "empty "+gg.Empty+", locked until "+gg.Reset, lk.bad.Bold(true))
		case gg.BurnPct > 0 && !gg.ResetAt.IsZero():
			e := col(gg.ResetAt)
			left := gg.ResetAt.Sub(now).Hours()
			used := min(100, int(math.Round(float64(100-gg.FiveH)+gg.BurnPct*left)))
			wa.put(nowC+1, rep("╌", max(1, e-nowC-1)), lk.sub)
			wa.put(max(e, nowC+2)+1, fmt.Sprintf("%d%% used at reset", used), lk.sub)
		default:
			wa.put(nowC+2, "idle", lk.dim)
		}
		// row 3: the week reset
		if c := col(gg.WeekAt); !gg.WeekAt.IsZero() && c >= 0 && c < p {
			lab := "week resets " + clock(gg.WeekAt) + " "
			if c-len(lab) < 0 {
				wk.put(c+1, " week resets "+clock(gg.WeekAt), lk.sub)
			} else {
				wk.put(c-len(lab), lab, lk.sub)
			}
			wk.put(c, "┃", lk.warn.Bold(true))
		}
		for r, gr := range []*grid{w5, wa, wk} {
			sp.rows = append(sp.rows, gut[r].join(line{spc(1)}).join(gr.line()))
		}
	}
	if len(fx.Team) > 0 {
		head := line{{"Team ", lk.title}, {"teammates' accounts, crow never uses them ", lk.dim}}
		sp.rows = append(sp.rows, head.join(line{{rep("╌", iw-head.width()), lk.dim}}))
		gw := 16
		if port {
			gw = 10
		}
		for _, t := range fx.Team {
			l := line{spc(2), {fit(t.Account, 8), lk.sub}, {"5h ", lk.dim}}
			l = append(l, gaugeBar(lk, t.FiveH, gw, lk.level(t.FiveH))...)
			l = append(l, seg{right(fmt.Sprintf("%d%%", t.FiveH), 5) + "   ", lk.sub}, seg{"wk ", lk.dim})
			if t.Week >= 0 {
				l = append(l, gaugeBar(lk, t.Week, gw, lk.level(t.Week))...)
				l = append(l, seg{right(fmt.Sprintf("%d%%", t.Week), 5), lk.sub})
			} else {
				l = append(l, seg{fit("no week limit", gw+5), lk.dim})
			}
			if !port {
				l = append(l, seg{"   " + ifs(t.Reset != "", "5h resets "+t.Reset) + ifs(t.WeekReset != "", "  ·  week resets "+t.WeekReset), lk.dim})
			}
			sp.rows = append(sp.rows, l)
		}
	}
	if fx.Usage != nil {
		sp.right = line{{"h l pan 6h · P policy · r refresh " + fx.Usage.Updated, lk.sub}}
	}
	return sp
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// clockOn puts an "HH:MM" on now's day, or the next day when it is already past.
func clockOn(now time.Time, hhmm string) time.Time {
	t, err := time.ParseInLocation("15:04", hhmm, now.Location())
	if err != nil {
		return now
	}
	at := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	if at.Before(now.Add(-time.Hour)) {
		at = at.AddDate(0, 0, 1)
	}
	return at
}

// gutter: the three text rows left of an account's lane, g cells wide.
func (uv usageViewState) gutter(gg gauge, g int, port bool) [3]line {
	lk := uv.lk
	pol := map[string]lipgloss.Style{"use": lk.ok, "ask-first": lk.warn, "never": lk.bad}[gg.Policy]
	if gg.Policy == "" {
		pol = lk.sub
	}
	dot := seg{"  ", lk.plain}
	name := lk.bold
	if gg.Active {
		dot, name = seg{"● ", lk.ok}, lk.acc
	}
	pctLeft := func(v int) seg { return seg{right(fmt.Sprintf("%d%%", v), 4) + " left", lk.level(v)} }
	var out [3]line
	if port {
		out[0] = line{dot, {gg.Account, name}}
		out[1] = line{spc(2), {fmt.Sprintf("5h %d%%", gg.FiveH), lk.level(gg.FiveH)}}
		if gg.Week >= 0 {
			out[1] = append(out[1], seg{fmt.Sprintf(" wk %d%%", gg.Week), lk.level(gg.Week)})
		}
		out[2] = line{spc(2), {cmpOr(gg.Policy, "-"), pol}}
	} else {
		gw := max(8, min(20, g-56))
		out[0] = line{dot, {fit(gg.Account, 8), name}, {fit("["+cmpOr(gg.Policy, "-")+"]", 12), pol}, {"5h ", lk.sub}}
		out[0] = append(out[0], gaugeBar(lk, gg.FiveH, gw, lk.level(gg.FiveH))...)
		out[0] = append(out[0], spc(1), pctLeft(gg.FiveH), seg{ifs(gg.Reset != "", " resets "+gg.Reset), lk.sub})
		out[1] = line{spc(22), {"wk ", lk.sub}}
		if gg.Week >= 0 {
			out[1] = append(out[1], gaugeBar(lk, gg.Week, gw, lk.level(gg.Week))...)
			out[1] = append(out[1], spc(1), pctLeft(gg.Week), seg{ifs(gg.WeekReset != "", " resets "+gg.WeekReset), lk.sub})
		} else {
			out[1] = append(out[1], seg{"no week limit", lk.sub})
		}
		burn := seg{fit(fmt.Sprintf("burn %.1f%%/h", gg.BurnPct), 20), lk.sub}
		if gg.BurnPct > 20 {
			burn.style = lk.bad
		}
		out[2] = line{spc(2), burn, {uv.status(gg), uv.statusStyle(gg)}}
		if gg.RunsMax > 0 {
			out[2] = append(out[2], seg{fmt.Sprintf("  ·  runs %d/%d", gg.RunsUsed, gg.RunsMax), lk.sub})
		}
	}
	for i := range out {
		out[i] = out[i].trunc(g)
		out[i] = append(out[i], spc(g-out[i].width()))
	}
	return out
}

// emptiesFirst: the 5h window runs dry before it resets.
func (uv usageViewState) emptiesFirst(gg gauge) bool {
	return gg.Empty != "" && (gg.ResetAt.IsZero() || clockOn(uv.now, gg.Empty).Before(gg.ResetAt))
}

// status is the one-line state of an account: empty soon, week spent, active or idle.
func (uv usageViewState) status(gg gauge) string {
	var s string
	switch {
	case uv.emptiesFirst(gg):
		s = "5h empty at " + gg.Empty
	case gg.Week >= 0 && gg.Week < 5:
		s = "week spent"
		if gg.Policy == "ask-first" {
			s += ", asks first"
		}
	case gg.Active:
		s = "active"
	case gg.BurnPct > 0:
		s = "burning"
	default:
		s = "idle"
	}
	if d := gg.WeekAt.Sub(uv.now); !gg.WeekAt.IsZero() && d > 0 && d < 24*time.Hour {
		s += fmt.Sprintf(", week resets in %dh%02d", int(d.Hours()), int(d.Minutes())%60)
	}
	return s
}

func (uv usageViewState) statusStyle(gg gauge) lipgloss.Style {
	switch {
	case uv.emptiesFirst(gg):
		return uv.lk.bad
	case gg.Week >= 0 && gg.Week < 5:
		return uv.lk.warn
	}
	return uv.lk.sub
}

// ----- daily spend -----

func (uv usageViewState) totalsLine() line {
	lk, u := uv.lk, uv.fx.Usage
	var l line
	for i, n := range periodNames {
		st := lk.sub
		if i == uv.period {
			st = lk.warn.Bold(true)
		}
		if i > 0 {
			l = append(l, spc(2))
		}
		l = append(l, seg{n + " ", lk.sub}, seg{usd(u.Periods[i].Cost), st})
	}
	return l
}

func (uv usageViewState) spendSpec(iw, ih, n int) boxSpec {
	lk, u := uv.lk, uv.fx.Usage
	unit := "cost"
	if uv.tokens {
		unit = "tokens"
	}
	sp := boxSpec{title: fmt.Sprintf("Daily spend · %d days", n), index: 2, sel: -1, focused: uv.pane == upSpend, bottom: "c " + map[bool]string{false: "tokens", true: "cost"}[uv.tokens]}
	if u == nil || len(u.Days) == 0 {
		return uv.empty(sp, "No spend recorded.")
	}
	if iw >= 100 {
		sp.right = uv.totalsLine()
	} else {
		sp.right = line{{unit, lk.sub}}
	}
	days := u.Days
	if len(days) > n {
		days = days[len(days)-n:]
	}
	val := func(d daySpend) float64 {
		if uv.tokens {
			return float64(d.Tokens)
		}
		return d.USD
	}
	label := func(v float64) string {
		if uv.tokens {
			return tokensShort(int(v))
		}
		return money(v)
	}
	top := 0.0
	for _, d := range days {
		top = max(top, val(d))
	}
	const lab = 7      // "$116 ┤"
	ph := max(2, ih-2) // plot rows; one row for the value labels, one for the dates
	slot := max(3, (iw-lab)/len(days))
	bw := max(1, min(8, slot-2))
	plot := make([]*grid, ph+1)
	for i := range plot {
		plot[i] = newGrid(slot*len(days), lk.plain)
	}
	scaled := func(v float64) int {
		if top <= 0 {
			return 0
		}
		return int(math.Round(v / top * 1000))
	}
	dates := newGrid(slot*len(days), lk.sub)
	short := false // labels too wide for a slot: day numbers only
	for _, d := range days {
		short = short || len(d.Label) >= slot
	}
	for i, d := range days {
		x := i*slot + (slot-bw)/2
		v := val(d)
		st := lk.acc
		last := i == len(days)-1
		if last {
			st = lk.warn
		}
		cells := barCells(scaled(v), 1000, ph)
		h := 0
		for r, c := range cells {
			if c != " " {
				plot[r+1].put(x, rep(c, bw), st)
				h = max(h, ph-r)
			}
		}
		if v <= 0 {
			plot[ph].put(x, rep("▁", bw), lk.dim)
		}
		s := label(v)
		lst := lk.sub
		if last {
			lst = lk.warn.Bold(true)
		}
		plot[ph-h].put(x+bw/2-len([]rune(s))/2, s, lst)
		dl := d.Label
		if short {
			dl = strings.TrimPrefix(d.Date[len(d.Date)-2:], "0")
		}
		dst := lk.sub
		if last {
			dst = lk.warn.Bold(true)
		}
		dates.put(i*slot+(slot-len(dl))/2, dl, dst)
	}
	for r, gr := range plot {
		ax := "       "
		switch r {
		case 1:
			ax = right(label(top), 5) + " ┤"
		case 1 + ph/2:
			ax = right(label(top*float64(ph-ph/2)/float64(ph)), 5) + " ┤"
		case ph:
			ax = right(label(0), 5) + " ┤"
		default:
			if r > 0 {
				ax = "      │"
			}
		}
		sp.rows = append(sp.rows, line{{ax, lk.dim}}.join(gr.line()))
	}
	sp.rows = append(sp.rows, line{{"      └", lk.dim}}.join(dates.line()))
	return sp
}

// ----- saved -----

func (uv usageViewState) savedSpec(iw int, port bool) boxSpec {
	lk, u := uv.lk, uv.fx.Usage
	sp := boxSpec{title: "Saved · " + periodNames[uv.period], index: 3, sel: -1, focused: uv.pane == upSaved}
	if u == nil {
		return uv.empty(sp, "No usage data.")
	}
	p := u.Periods[uv.period]
	peak := max(p.CacheUSD, p.RouteUSD, 1)
	routed := ""
	if p.Runs > 0 {
		routed = fmt.Sprintf("%d of %d runs on a cheaper model", p.Routed, p.Runs)
	}
	rows := []struct {
		k, v, det string
		f         float64
	}{
		{"prompt cache", usd(p.CacheUSD), fmt.Sprintf("%d%% of input read from cache", p.CachePct), p.CacheUSD / peak},
		{"model routing", usd(p.RouteUSD), routed, p.RouteUSD / peak},
		{"rtk", tokensShort(p.RTKTokens) + " tok", fmt.Sprintf("%.0f%% of command output cut, %d cmds", p.RTKPct, p.RTKCmds), p.RTKPct / 100},
	}
	bw := max(6, iw-14-11)
	for _, r := range rows {
		l := line{{fit(r.k, 14), lk.head}}.join(hbar(r.f, bw, lk.ok)).join(line{{right(r.v, 11), lk.ok.Bold(true)}})
		sp.rows = append(sp.rows, l)
		if !port && r.det != "" {
			sp.rows = append(sp.rows, line{spc(14), {r.det, lk.sub}})
		} else if !port {
			sp.rows = append(sp.rows, nil)
		}
	}
	sp.rows = append(sp.rows, line{{rep("╌", iw), lk.dim}})
	tot := line{{"total ", lk.head}, {usd(p.savedUSD()), lk.ok.Bold(true)}, {" + " + tokensShort(p.RTKTokens) + " tok", lk.sub},
		{"   compactions ", lk.head}, {fmt.Sprint(p.Compactions), lk.warn.Bold(true)}}
	if uv.period != perMonth {
		tot = append(tot, seg{"   retros ", lk.head}, seg{fmt.Sprint(p.Retros), lk.ok.Bold(true)})
	}
	sp.rows = append(sp.rows, tot)
	return sp
}

// ----- by project -----

func (uv usageViewState) projectsSpec(iw, ih int) boxSpec {
	lk := uv.lk
	sp := boxSpec{title: "By project · " + periodNames[uv.period], index: 5, sel: -1, focused: uv.pane == upProjects, right: line{{"enter filter ledger", lk.sub}}}
	p := uv.per()
	if uv.fx.Usage == nil || len(p.Projects) == 0 {
		return uv.empty(sp, "Nothing spent.")
	}
	nw := min(22, max(10, iw/4))
	rw := 9 + 5 + 8
	bw := max(4, iw-nw-1-rw-1)
	sp.hdr = line{{fit("project", nw+1+bw+1), lk.head}, {right("cost", 9) + right("%", 5) + right("tokens", 8), lk.head}}
	tot := max(p.Cost, 0.01)
	peak := max(p.Projects[0].USD, 0.01)
	rows := p.Projects
	room := max(1, ih-1)
	cut := len(rows) > room && uv.pane != upProjects
	if cut {
		rows = rows[:room-1]
	}
	for i, r := range rows {
		st, bst := lk.plain, lk.cyan
		if i == 0 {
			st, bst = lk.bold, lk.acc
		}
		if r.Name == uv.filter {
			st = lk.warn.Bold(true)
		}
		l := line{{fit(r.Name, nw), st}, spc(1)}.join(hbar(r.USD/peak, bw, bst)).join(line{spc(1),
			{right(usd(r.USD), 9), ifStyle(i == 0, lk.warn.Bold(true), lk.plain)}, {right(fmt.Sprintf("%.0f%%", r.USD/tot*100), 5), lk.sub},
			{right(tokensShort(r.Tokens), 8), lk.sub}})
		sp.rows = append(sp.rows, l)
	}
	if cut {
		rest, sum := p.Projects[len(rows):], 0.0
		for _, r := range rest {
			sum += r.USD
		}
		sp.rows = append(sp.rows, line{{fmt.Sprintf("… %d more, %s", len(rest), usd(sum)), lk.dim}})
	}
	if uv.pane == upProjects {
		sp.sel = max(0, min(uv.prj, len(rows)-1))
	}
	return sp
}

func ifStyle(ok bool, a, b lipgloss.Style) lipgloss.Style {
	if ok {
		return a
	}
	return b
}

// ----- by model -----

func (uv usageViewState) modelsSpec(iw int, port bool) boxSpec {
	lk := uv.lk
	sp := boxSpec{title: "By model · " + periodNames[uv.period], index: 6, sel: -1, focused: uv.pane == upModels}
	p := uv.per()
	if uv.fx.Usage == nil || len(p.Models) == 0 {
		return uv.empty(sp, "Nothing spent.")
	}
	tot := 0.0
	for _, m := range p.Models {
		tot += m.USD
	}
	tot = max(tot, 0.01)
	const nw = 12
	bw := max(6, iw-nw-9-6)
	for i, m := range p.Models {
		share := m.USD / tot
		st := lk.cyan
		if i == 0 {
			st = lk.acc
		}
		bar := hbar(share, bw, st)
		if share < .005 {
			bar = line{{"·", lk.dim}, spc(bw - 1)}
		}
		sp.rows = append(sp.rows, line{{fit(m.Name, nw), lk.bold}}.join(bar).join(line{{right(usd(m.USD), 9), lk.warn}, {right(fmt.Sprintf("%.0f%%", share*100), 6), lk.sub}}))
		if !port {
			sp.rows = append(sp.rows, line{spc(nw), {fmt.Sprintf("%s tok · %d sessions", tokensShort(m.Tokens), m.Sessions), lk.sub}})
		}
	}
	cache := line{{fit("cache read", nw), lk.head}}.join(gaugeBar(lk, p.CachePct, bw, lk.ok)).join(line{{right(fmt.Sprintf("%d%%", p.CachePct), 15), lk.ok.Bold(true)}})
	sp.rows = append(sp.rows, cache)
	if len(p.Tools) == 0 {
		return sp
	}
	tools := p.Tools[:min(4, len(p.Tools))]
	if port {
		l := line{{fit("tools", nw), lk.head}}
		var parts []string
		for _, t := range tools {
			parts = append(parts, fmt.Sprintf("%s %d", t.Name, t.Calls))
		}
		sp.rows = append(sp.rows, append(l, seg{strings.Join(parts, "  "), lk.sub}))
		return sp
	}
	tw := max(4, bw-16)
	for i, t := range tools {
		k := ""
		if i == 0 {
			k = "tools"
		}
		st := lk.sub
		if i > 0 {
			st = lk.dim
		}
		sp.rows = append(sp.rows, line{{fit(k, nw), lk.head}, {fit(t.Name, 16), lk.plain}}.join(hbar(float64(t.Calls)/float64(max(1, tools[0].Calls)), tw, st)).join(line{{right(fmt.Sprint(t.Calls), 15), lk.sub}}))
	}
	return sp
}

// ----- ledger -----

type ledGroup struct {
	Name   string
	Rows   []ledgerEntry
	Cost   float64
	Tokens int
}

// ledRow is one display row of the ledger: a group head, an entry, or the "… more" line.
type ledRow struct {
	kind  int // 0 group, 1 entry, 2 more
	group string
	e     ledgerEntry
	rest  []ledgerEntry
}

func (uv usageViewState) ledgerEntries() []ledgerEntry {
	u := uv.fx.Usage
	if u == nil {
		return nil
	}
	var out []ledgerEntry
	for _, e := range u.entries(uv.period, uv.now) {
		if uv.filter == "" || e.Project == uv.filter {
			out = append(out, e)
		}
	}
	return out
}

func groupsOf(es []ledgerEntry) []ledGroup {
	idx := map[string]int{}
	var gs []ledGroup
	for _, e := range es {
		i, ok := idx[e.Project]
		if !ok {
			i = len(gs)
			idx[e.Project] = i
			gs = append(gs, ledGroup{Name: e.Project})
		}
		gs[i].Rows = append(gs[i].Rows, e)
		gs[i].Cost += e.Cost
		gs[i].Tokens += e.Tokens
	}
	sort.SliceStable(gs, func(i, j int) bool { return gs[i].Cost > gs[j].Cost })
	return gs
}

// foldOf: the fold state of group i, default open for the three costliest groups.
func (uv usageViewState) foldOf(name string, i int) int {
	if f, ok := uv.fold[name]; ok {
		return f
	}
	if i < 3 {
		return foldOpen
	}
	return foldClosed
}

// ledRows lays the ledger out for room display lines (rh lines per entry): open groups share
// the room, the costliest first; a group shows "… N more" when it does not fit.
func (uv usageViewState) ledRows(room, rh int) []ledRow {
	es := uv.ledgerEntries()
	if uv.flat {
		var out []ledRow
		for _, e := range es {
			out = append(out, ledRow{kind: 1, e: e})
		}
		return out
	}
	gs := groupsOf(es)
	shown := make([]int, len(gs))
	budget := room - len(gs)
	for i, g := range gs {
		if uv.foldOf(g.Name, i) == foldAll {
			shown[i] = len(g.Rows)
			budget -= len(g.Rows) * rh
		}
	}
	weights := []int{6, 3, 1}
	for pass := 0; pass < 2; pass++ {
		for i, g := range gs {
			if uv.foldOf(g.Name, i) != foldOpen {
				continue
			}
			want := len(g.Rows)
			if pass == 0 {
				w := 1
				if i < len(weights) {
					w = weights[i]
				}
				want = min(want, max(1, (budget-1)*w/10/rh))
			}
			for shown[i] < want {
				more := btoi(shown[i]+1 < len(g.Rows))
				had := btoi(shown[i] > 0 && shown[i] < len(g.Rows))
				cost := rh + more - had
				if shown[i] == 0 {
					cost = rh + more
				}
				if cost > budget && shown[i] > 0 {
					break
				}
				budget -= cost
				shown[i]++
			}
		}
	}
	var out []ledRow
	for i, g := range gs {
		out = append(out, ledRow{kind: 0, group: g.Name})
		f := uv.foldOf(g.Name, i)
		if f == foldClosed {
			continue
		}
		n := max(1, shown[i])
		for _, e := range g.Rows[:min(n, len(g.Rows))] {
			out = append(out, ledRow{kind: 1, group: g.Name, e: e})
		}
		if n < len(g.Rows) {
			out = append(out, ledRow{kind: 2, group: g.Name, rest: g.Rows[n:]})
		}
	}
	return out
}

func (uv usageViewState) ledgerSpec(iw, ih int, port bool) (boxSpec, []ledRow) {
	lk := uv.lk
	es := uv.ledgerEntries()
	tot, tok := 0.0, 0
	for _, e := range es {
		tot += e.Cost
		tok += e.Tokens
	}
	gs := groupsOf(es)
	sp := boxSpec{title: fmt.Sprintf("Ledger · %d entries · %d projects · %s", len(es), len(gs), usd(tot)), index: 4, sel: -1, focused: uv.pane == upLedger,
		right: line{{ifs(uv.filter != "", "project "+uv.filter+" · ") + "g group · enter fold · o open in history", lk.sub}}}
	if u := uv.fx.Usage; u != nil && uv.period == perMonth && u.Since != "" {
		sp.title += " · since " + isoDay(u.Since)
	}
	if len(es) == 0 {
		return uv.empty(sp, "No session in this period."), nil
	}
	sp.bottom = fmt.Sprintf("%d entries · %s tok · %s · %s", len(es), tokensShort(tok), usd(tot), periodNames[uv.period])
	rh := 1
	if port {
		rh = 2
	}
	cols := ledCols(iw, port)
	sp.hdr = cols.head(lk)
	rows := uv.ledRows(max(1, ih-1), rh)
	gcost := map[string]ledGroup{}
	for _, g := range gs {
		gcost[g.Name] = g
	}
	sel := max(0, min(uv.led, len(rows)-1))
	for i, r := range rows {
		if i == sel && uv.pane == upLedger {
			sp.sel = len(sp.rows)
			if r.kind == 1 {
				sp.selN = rh
			}
		}
		switch r.kind {
		case 0:
			g := gcost[r.group]
			mark := "+"
			if j := indexOfGroup(gs, r.group); uv.foldOf(r.group, j) != foldClosed {
				mark = "-"
			}
			share := g.Cost / max(tot, 0.01)
			name := line{{g.Name, lk.bold}, {fmt.Sprintf("  %d %s", len(g.Rows), plural(len(g.Rows), "entry", "entries")), lk.sub}}
			cells := map[string]line{"_g": {{mark, lk.acc}}, "session": name, "tokens": {{tokensShort(g.Tokens), lk.head}}, "cost": {{usd(g.Cost), lk.head}},
				"context": line{{"", lk.plain}}.join(hbar(share, 5, lk.acc)).join(line{{right(fmt.Sprintf("%.0f%%", share*100), 5), lk.sub}})}
			sp.rows = append(sp.rows, cols.row(cells))
		case 1:
			sp.rows = append(sp.rows, uv.entryRows(cols, r.e, port)...)
		case 2:
			c, t := 0.0, 0
			for _, e := range r.rest {
				c += e.Cost
				t += e.Tokens
			}
			sp.rows = append(sp.rows, cols.row(map[string]line{"session": {{fmt.Sprintf("… %d more, %s tok, %s · enter shows all", len(r.rest), tokensShort(t), usd(c)), lk.sub}}}))
		}
	}
	return sp, rows
}

func indexOfGroup(gs []ledGroup, name string) int {
	for i, g := range gs {
		if g.Name == name {
			return i
		}
	}
	return len(gs)
}

func (uv usageViewState) entryRows(cols ledColumns, e ledgerEntry, port bool) []line {
	lk := uv.lk
	g := seg{"✓", lk.sub}
	switch {
	case e.Active:
		g = seg{"●", lk.ok}
	case e.Kind == "run":
		g = seg{"◆", lk.cyan}
	}
	ctxSt := map[string]lipgloss.Style{"warn": lk.warn.Bold(true), "over": lk.bad.Bold(true)}[e.CtxLevel]
	if e.CtxLevel == "ok" || e.CtxLevel == "" {
		ctxSt = lk.dim
	}
	costSt := lk.plain
	if e.Cost >= 10 {
		costSt = lk.warn.Bold(true)
	}
	ctx := line{{tokensShort(e.Ctx) + " " + cmpOr(e.CtxLevel, "ok"), ctxSt}}
	modelSt := lk.sub
	if strings.Contains(e.Model, "opus") {
		modelSt = lk.cyan
	}
	week := ""
	if e.Week >= 0 {
		week = pctText(e.Week)
	}
	cells := map[string]line{"_g": {g}, "day": {{isoDay(e.Day), lk.sub}}, "time": {{entryTime(e), lk.sub}}, "acct": {{e.Account, lk.sub}},
		"session": {{e.Title, lk.plain}}, "model": {{e.Model, modelSt}}, "tokens": {{tokensShort(e.Tokens), lk.plain}},
		"cost": {{usd(e.Cost), costSt}}, "week": {{week, lk.sub}}, "context": ctx, "cmp": {{ifs(e.Compactions > 0, fmt.Sprint(e.Compactions)), lk.warn}}}
	out := []line{cols.row(cells)}
	if port {
		d := fmt.Sprintf("%s   %s %s   %s   out %s", e.Model, isoDay(e.Day), entryTime(e), e.Account, tokensShort(e.Out))
		if week != "" {
			d += "   week " + week
		}
		if e.Compactions > 0 {
			d += fmt.Sprintf("   %d %s", e.Compactions, plural(e.Compactions, "compaction", "compactions"))
		}
		out = append(out, line{spc(3), {d, lk.sub}})
	}
	return out
}

// entryTime: "11:02–13:40" within the day, "…13:40" when it started before, "11:02…" when it ran on.
func entryTime(e ledgerEntry) string {
	sd, ed := e.Start.Format("2006-01-02"), e.End.Format("2006-01-02")
	switch {
	case sd == e.Day && ed == e.Day:
		return e.Start.Format("15:04") + "–" + e.End.Format("15:04")
	case ed == e.Day:
		return "…" + e.End.Format("15:04")
	}
	return e.Start.Format("15:04") + "…"
}

type ledCol struct {
	key, name string
	w         int // 0: takes the rest
	right     bool
}

type ledColumns struct {
	cols []ledCol
	gap  int
}

func ledCols(iw int, port bool) ledColumns {
	var cs []ledCol
	gap := 3
	if port {
		gap = 2
		cs = []ledCol{{"_g", "", 1, false}, {"session", "session", 0, false}, {"tokens", "tokens", 7, true}, {"cost", "cost", 8, true}, {"context", "context", 12, false}}
	} else {
		cs = []ledCol{{"_g", "", 1, false}, {"day", "day", 6, false}, {"time", "time", 11, false}, {"acct", "acct", 7, false}, {"session", "session", 0, false},
			{"model", "model", 11, false}, {"tokens", "tokens", 7, true}, {"cost", "cost", 8, true}, {"week", "week", 4, true}, {"context", "context", 12, false}, {"cmp", "cmp", 3, true}}
	}
	fixed := gap * (len(cs) - 1)
	for _, c := range cs {
		fixed += c.w
	}
	for i := range cs {
		if cs[i].w == 0 {
			cs[i].w = max(12, iw-fixed)
		}
	}
	return ledColumns{cs, gap}
}

func (lc ledColumns) head(lk *look) line {
	cells := map[string]line{}
	for _, c := range lc.cols {
		cells[c.key] = line{{c.name, lk.head}}
	}
	return lc.row(cells)
}

func (lc ledColumns) row(cells map[string]line) line {
	var out line
	for i, c := range lc.cols {
		if i > 0 {
			out = append(out, spc(lc.gap))
		}
		cell := cells[c.key].trunc(c.w)
		pad := spc(c.w - cell.width())
		if c.right {
			out = append(append(out, pad), cell...)
		} else {
			out = append(append(out, cell...), pad)
		}
	}
	return out
}

// ----- composition -----

func (uv usageViewState) body(w, h int) []string {
	lk := uv.lk
	lay := layoutFor(w)
	if lay == laySingle {
		return uv.single(w, h)
	}
	gapH, gapV := lk.gaps()
	port := lay != layLandscape
	iw := lk.innerW(w)
	rwH := lk.boxH(runwayRows(len(uv.fx.Gauges), len(uv.fx.Team)))
	if uv.loading() || len(uv.fx.Gauges)+len(uv.fx.Team) == 0 {
		rwH = lk.boxH(3)
	}
	out := vcat(lk.box(uv.runwaySpec(iw, port), w, rwH), blankLines(w, gapV))
	rest := h - len(out)
	if !port {
		bh := min(13, max(lk.boxH(6), rest/3))
		sw := max(52, (w-gapH)*28/100)
		dw := w - gapH - sw
		out = vcat(out, hcat([]int{dw, sw}, gapH,
			lk.box(uv.spendSpec(lk.innerW(dw), bh-lk.boxH(0), 14), dw, bh),
			lk.box(uv.savedSpec(lk.innerW(sw), false), sw, bh)), blankLines(w, gapV))
		lh := h - len(out)
		lw := (w - gapH) * 64 / 100
		rw := w - gapH - lw
		ph := lh / 2
		led, _ := uv.ledgerSpec(lk.innerW(lw), lh-lk.boxH(0), false)
		right := vcat(lk.box(uv.projectsSpec(lk.innerW(rw), ph-lk.boxH(0)-1), rw, ph), blankLines(rw, gapV),
			lk.box(uv.modelsSpec(lk.innerW(rw), false), rw, lh-ph-gapV))
		out = vcat(out, hcat([]int{lw, rw}, gapH, lk.box(led, lw, lh), right))
	} else {
		bh := lk.boxH(8)
		half := (w - gapH) / 2
		w2 := w - gapH - half
		out = vcat(out,
			hcat([]int{half, w2}, gapH, lk.box(uv.spendSpec(lk.innerW(half), bh-lk.boxH(0), 7), half, bh), lk.box(uv.savedSpec(lk.innerW(w2), true), w2, bh)),
			blankLines(w, gapV),
			hcat([]int{half, w2}, gapH, lk.box(uv.modelsSpec(lk.innerW(half), true), half, bh), lk.box(uv.projectsSpec(lk.innerW(w2), bh-lk.boxH(0)-1), w2, bh)),
			blankLines(w, gapV))
		lh := max(lk.boxH(3), h-len(out))
		led, _ := uv.ledgerSpec(iw, lh-lk.boxH(0), true)
		out = vcat(out, lk.box(led, w, lh))
	}
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func (uv usageViewState) single(w, h int) []string {
	lk := uv.lk
	iw := lk.innerW(w)
	ih := h - lk.boxH(0)
	var sp boxSpec
	switch uv.pane {
	case upSpend:
		sp = uv.spendSpec(iw, ih, 7)
	case upSaved:
		sp = uv.savedSpec(iw, true)
	case upLedger:
		sp, _ = uv.ledgerSpec(iw, ih, true)
	case upProjects:
		sp = uv.projectsSpec(iw, ih-1)
	case upModels:
		sp = uv.modelsSpec(iw, true)
	default:
		sp = uv.runwaySpec(iw, true)
	}
	sp.bottom = fmt.Sprintf("pane %d of %d · tab", uv.pane+1, usagePanes)
	b := lk.box(sp, w, h)
	for len(b) < h {
		b = append(b, blank(w))
	}
	return b[:h]
}

// ----- keys -----

func (m *model) usageMove(d int) {
	uv := m.usageView()
	switch m.usePane {
	case upRunway:
		m.useSel = max(0, min(len(m.fx.Gauges)-1, m.useSel+d))
	case upLedger:
		_, rows := uv.ledgerSpec(200, 1000, false)
		m.useLed = max(0, min(len(rows)-1, m.useLed+d))
	case upProjects:
		m.usePrj = max(0, min(len(uv.per().Projects)-1, m.usePrj+d))
	}
}

// usageKey handles the Usage tab keys; ok false hands the key to the global switch.
func (m model) usageKey(k string) (tea.Model, tea.Cmd, bool) {
	uv := m.usageView()
	switch k {
	case "p":
		m.usePeriod = (m.usePeriod + 1) % len(periodNames)
		m.useLed, m.usePrj = 0, 0
	case "h", "l", "[", "]", "<", ">":
		if k == "h" || k == "[" || k == "<" {
			m.useShift = min(m.useShift+6, 48)
		} else {
			m.useShift = max(m.useShift-6, -12)
		}
	case "c":
		m.useTok = !m.useTok
	case "g":
		m.useFlat, m.useLed = !m.useFlat, 0
	case "P":
		m.flash = "P: editing the account policy arrives with the actions slice"
	case "r":
		m.flash = "refreshing usage…"
		return m, fetchCmd(srcUsage), true
	case "esc":
		if m.useFilter == "" {
			return m, nil, false
		}
		m.useFilter, m.useLed = "", 0
	case "enter":
		switch m.usePane {
		case upProjects:
			if ps := uv.per().Projects; m.usePrj < len(ps) {
				name := ps[m.usePrj].Name
				m.useFilter = ifs(m.useFilter != name, name)
				m.useLed = 0
			}
		case upLedger:
			_, rows := uv.ledgerSpec(200, 1000, false)
			if m.useLed >= len(rows) || rows[m.useLed].kind == 1 {
				return m, nil, true
			}
			r := rows[m.useLed]
			gs := groupsOf(uv.ledgerEntries())
			cur := uv.foldOf(r.group, indexOfGroup(gs, r.group))
			next := map[int]int{foldClosed: foldOpen, foldOpen: foldAll, foldAll: foldClosed}[cur]
			if r.kind == 2 {
				next = foldAll
			}
			fold := map[string]int{}
			for k, v := range m.useFold {
				fold[k] = v
			}
			fold[r.group] = next
			m.useFold = fold
		}
	case "o":
		_, rows := uv.ledgerSpec(200, 1000, false)
		if m.usePane != upLedger || m.useLed >= len(rows) || rows[m.useLed].kind != 1 {
			m.flash = "o: select a session in the Ledger"
			return m, nil, true
		}
		m.page, m.hist, m.histTyping, m.histPane = 1, true, false, hpList
		m.histF = newHistFilter()
		m.histF.query = rows[m.useLed].e.Title
	default:
		return m, nil, false
	}
	return m, nil, true
}
