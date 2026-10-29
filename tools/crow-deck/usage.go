package main

// Usage tab: accounts and policy with the 5h and week gauges, the forecast, daily spend,
// savings and the cost of today by session, project and model. Costs are API-equivalent
// dollars from crow-usage, not what is billed. Read only: policy edits arrive with the
// actions slice.

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	upAccounts = iota
	upForecast
	upSpend
	upSavings
	upCost
	usagePanes
)

type usageViewState struct {
	lk        *look
	fx        *fixtures
	sync      string
	sel       int
	pane      int
	fortnight bool
	now       time.Time
}

func (uv usageViewState) loading() bool {
	return uv.sync == stLoading || uv.fx.Live && uv.fx.Usage == nil && len(uv.fx.Gauges) == 0
}

func money(v float64) string {
	if v >= 10 {
		return fmt.Sprintf("$%.0f", v)
	}
	return fmt.Sprintf("$%.1f", v)
}

func tokensShort(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 1000:
		return fmt.Sprintf("%dk", n/1000)
	}
	return fmt.Sprint(n)
}

// ----- accounts -----

func (uv usageViewState) accountsSpec(iw int) boxSpec {
	lk, fx := uv.lk, uv.fx
	sp := boxSpec{title: "Accounts and policy", index: 1, sel: -1, focused: uv.pane == upAccounts}
	if uv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if len(fx.Gauges) == 0 {
		sp.rows = []line{{{"No account reports limits.", lk.sub}}}
		return sp
	}
	cells := max(6, min(14, (iw-78)/2))
	sp.rows = append(sp.rows, line{{"#  ", lk.sub}, {fit("account", 9), lk.sub}, {fit("mode", 10), lk.sub}, {fit("5h window", cells+13), lk.sub}, {fit("week", cells+21), lk.sub}, {"burn · runs", lk.sub}})
	for i, g := range fx.Gauges {
		l := line{{fmt.Sprintf("%-3d", i+1), lk.sub}, {fit(g.Account, 9), lk.bold}, {fit(cmpOr(g.Policy, "-"), 10), lk.plain}}
		l = append(l, gaugeBar(lk, g.FiveH, cells, lk.level(g.FiveH))...)
		l = append(l, seg{fit(fmt.Sprintf(" %d%% %s", g.FiveH, g.Reset), 13), lk.plain})
		if g.Week >= 0 {
			l = append(l, gaugeBar(lk, g.Week, cells, lk.level(g.Week))...)
			l = append(l, seg{fit(fmt.Sprintf(" %d%% %s", g.Week, weekWord(g.WeekReset)), 21), lk.plain})
		} else {
			l = append(l, seg{fit("no week limit", cells+21), lk.sub})
		}
		l = append(l, seg{fmt.Sprintf("%s · %d/%d", strings.TrimPrefix(cmpOr(g.Burn, "no burn"), "burn "), g.RunsUsed, g.RunsMax), lk.sub})
		sp.rows = append(sp.rows, l)
	}
	sp.rows = append(sp.rows, nil, line{{"p mode  l cap  [ ] priority", lk.key}, spc(2), {"editing the policy arrives with the actions slice", lk.sub}})
	sp.sel = 1 + max(0, min(uv.sel, len(fx.Gauges)-1))
	if fx.Usage != nil {
		sp.right = line{{"r refresh · " + fx.Usage.Updated, lk.sub}}
	}
	return sp
}

// weekWord turns "Nov 2 17:00" into "Nov 2"; a bare weekday stays.
func weekWord(s string) string {
	if f := strings.Fields(s); len(f) == 3 {
		return "resets " + f[0] + " " + f[1]
	}
	return ifs(s != "", "resets "+s)
}

// ----- forecast -----

func (uv usageViewState) forecastSpec(iw int) boxSpec {
	lk, fx := uv.lk, uv.fx
	sp := boxSpec{title: "Forecast", index: 2, sel: -1, focused: uv.pane == upForecast}
	if uv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if len(fx.Gauges) == 0 {
		sp.rows = []line{{{"Nothing to forecast.", lk.sub}}}
		return sp
	}
	for _, g := range fx.Gauges {
		l := line{{fit(g.Account, 8), lk.bold}}
		switch {
		case g.BurnPct <= 0:
			l = append(l, seg{"5h idle, nothing burning", lk.sub})
		case g.Empty != "":
			l = append(l, seg{"5h empty ", lk.sub}, seg{"~" + g.Empty, lk.level(g.FiveH)}, seg{ifs(g.Reset != "", " · resets "+g.Reset), lk.sub})
		default:
			l = append(l, seg{"5h lasts until the reset", lk.sub})
		}
		sp.rows = append(sp.rows, l)
		if g.BurnPct > 0 {
			sp.rows = append(sp.rows, line{{fit("", 8), lk.sub}, {fmt.Sprintf("burn %.0f%%/h", g.BurnPct), lk.sub}, {ifs(g.Week >= 0, fmt.Sprintf(" · week %d%% left", g.Week)), lk.sub}})
		}
	}
	if u := fx.Usage; u != nil && len(u.Days) > 0 {
		sp.rows = append(sp.rows, nil, line{{fit("spend", 8), lk.sub}, {spark(u.Days), lk.cyan}, {" " + fmt.Sprint(len(u.Days)) + " days", lk.sub}})
	}
	return sp
}

func spark(days []daySpend) string {
	peak := 0.0
	for _, d := range days {
		peak = max(peak, d.USD)
	}
	var b strings.Builder
	for _, d := range days {
		i := 0
		if peak > 0 {
			i = int(d.USD / peak * 7)
		}
		b.WriteRune([]rune("▁▂▃▄▅▆▇█")[min(7, i)])
	}
	return b.String()
}

// ----- daily spend -----

func (uv usageViewState) spendSpec(iw, h int) boxSpec {
	lk, u := uv.lk, uv.fx.Usage
	sp := boxSpec{title: "Daily spend", index: 3, sel: -1, focused: uv.pane == upSpend}
	n := 7
	if uv.fortnight {
		n = 14
	}
	sp.title += fmt.Sprintf(" · %d days", n)
	sp.right = line{{"d " + map[bool]string{false: "14", true: "7"}[uv.fortnight] + " days", lk.sub}}
	if uv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if u == nil || len(u.Days) == 0 {
		sp.rows = []line{{{"No spend recorded.", lk.sub}}}
		return sp
	}
	days := u.Days
	if len(days) > n {
		days = days[len(days)-n:]
	}
	bars := max(2, h-lk.boxH(0)-2)
	cw := max(4, min(8, iw/len(days)))
	peak := 0.0
	for _, d := range days {
		peak = max(peak, d.USD)
	}
	height := func(d daySpend) int {
		if peak == 0 || d.USD <= 0 {
			return 0
		}
		return max(1, int(d.USD/peak*float64(bars)+0.5))
	}
	center := func(s string, st lipgloss.Style) seg {
		pad := max(0, cw-len([]rune(s)))
		return seg{rep(" ", pad/2) + s + rep(" ", pad-pad/2), st}
	}
	for r := 0; r <= bars; r++ {
		var l line
		for i, d := range days {
			hh := height(d)
			st := lk.cyan
			if i == len(days)-1 {
				st = lk.acc
			}
			switch top := bars - hh; {
			case r == top:
				l = append(l, center(fmt.Sprintf("$%.0f", d.USD), lk.plain))
			case r > top:
				l = append(l, center("███", st))
			default:
				l = append(l, spc(cw))
			}
		}
		sp.rows = append(sp.rows, l)
	}
	var lab line
	for _, d := range days {
		lab = append(lab, center(d.Label, lk.sub))
	}
	sp.rows = append(sp.rows, lab)
	return sp
}

// ----- savings -----

func (uv usageViewState) savingsSpec() boxSpec {
	lk, u := uv.lk, uv.fx.Usage
	sp := boxSpec{title: "Savings", index: 4, sel: -1, focused: uv.pane == upSavings}
	if uv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if u == nil {
		sp.rows = []line{{{"No usage data.", lk.sub}}}
		return sp
	}
	row := func(k, v, note string) line {
		return line{{fit(k, 15), lk.sub}, {fit(v, 13), lk.ok}, {note, lk.sub}}
	}
	sp.rows = []line{
		{{"rtk today and this week", lk.sub}},
		row("rtk filtering", tokensShort(u.RTKTokens)+" tokens", fmt.Sprintf("%.0f%% of input", u.RTKPct)),
		row("cache reads", fmt.Sprintf("%d%%", u.CachePct), "~"+money(u.CacheUSD)+" saved"),
		row("model choice", "~"+money(u.RouteUSD), "smaller model"),
		nil,
		{{fit("total saved", 15), lk.bold}, {"~" + money(u.savedTotal()) + " this week", lk.ok.Bold(true)}},
	}
	return sp
}

// ----- cost today -----

func costCol(lk *look, title string, rows []costRow, cw int) []line {
	out := []line{{{title, lk.sub}}}
	peak := 0.0
	for _, r := range rows {
		peak = max(peak, r.USD)
	}
	name := min(14, max(8, cw/3))
	cells := max(3, cw-name-8)
	for _, r := range rows {
		pct := 0
		if peak > 0 {
			pct = int(r.USD / peak * 100)
		}
		l := line{{fit(r.Name, name), lk.plain}, spc(1)}
		l = append(l, gaugeBar(lk, pct, cells, lk.cyan)...)
		l = append(l, seg{fmt.Sprintf(" %6s", money(r.USD)), lk.plain})
		out = append(out, l)
	}
	if len(rows) == 0 {
		out = append(out, line{{"nothing today", lk.sub}})
	}
	return out
}

func (uv usageViewState) costSpec(iw int) boxSpec {
	lk, u := uv.lk, uv.fx.Usage
	sp := boxSpec{title: "Cost today", index: 5, sel: -1, focused: uv.pane == upCost}
	sp.right = line{{"API-equivalent $, not billed", lk.sub}}
	if uv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if u == nil {
		sp.rows = []line{{{"No usage data.", lk.sub}}}
		return sp
	}
	if iw < 90 {
		for _, c := range [][]line{costCol(lk, "per session", u.Sessions, iw), costCol(lk, "per project", u.Projects, iw), costCol(lk, "per model", u.Models, iw)} {
			sp.rows = append(sp.rows, c...)
			sp.rows = append(sp.rows, nil)
		}
		return sp
	}
	gap := 3
	cw := (iw - 2*gap) / 3
	cols := [][]line{costCol(lk, "per session", u.Sessions, cw), costCol(lk, "per project", u.Projects, cw), costCol(lk, "per model", u.Models, cw)}
	n := 0
	for _, c := range cols {
		n = max(n, len(c))
	}
	for i := 0; i < n; i++ {
		var l line
		for j, c := range cols {
			var cell line
			if i < len(c) {
				cell = c[i].trunc(cw)
			}
			l = append(l, cell...)
			l = append(l, spc(cw-cell.width()))
			if j < 2 {
				l = append(l, spc(gap))
			}
		}
		sp.rows = append(sp.rows, l)
	}
	return sp
}

// ----- composition -----

func (uv usageViewState) body(w, h int) []string {
	lk := uv.lk
	lay := layoutFor(w)
	if lay == laySingle {
		return uv.single(w, h)
	}
	gapH, gapV := 0, 0
	if lk.th.frame == frameSoft {
		gapH, gapV = 3, 1
	}
	accH := lk.boxH(len(uv.fx.Gauges) + 3)
	if len(uv.fx.Gauges) == 0 {
		accH = lk.boxH(1)
	}
	var out block
	acc := lk.box(uv.accountsSpec(lk.innerW(w)), w, accH)
	if lay == layLandscape {
		costH := lk.boxH(7)
		midH := h - accH - costH - 2*gapV
		if midH < lk.boxH(9) {
			midH, costH = lk.boxH(9), max(lk.boxH(4), h-accH-lk.boxH(9)-2*gapV)
		}
		w1 := (w - 2*gapH) * 32 / 100
		w3 := (w - 2*gapH) * 30 / 100
		w2 := w - 2*gapH - w1 - w3
		mid := hcat([]int{w1, w2, w3}, gapH,
			lk.box(uv.forecastSpec(lk.innerW(w1)), w1, midH),
			lk.box(uv.spendSpec(lk.innerW(w2), midH), w2, midH),
			lk.box(uv.savingsSpec(), w3, midH))
		out = vcat(acc, blankLines(w, gapV), mid, blankLines(w, gapV), lk.box(uv.costSpec(lk.innerW(w)), w, costH))
	} else {
		fh := lk.boxH(6)
		ch := lk.boxH(8)
		sh := max(lk.boxH(8), h-accH-fh-ch-3*gapV)
		w1 := (w - gapH) / 2
		w2 := w - gapH - w1
		out = vcat(acc, blankLines(w, gapV),
			hcat([]int{w1, w2}, gapH, lk.box(uv.forecastSpec(lk.innerW(w1)), w1, fh), lk.box(uv.savingsSpec(), w2, fh)), blankLines(w, gapV),
			lk.box(uv.spendSpec(lk.innerW(w), sh), w, sh), blankLines(w, gapV),
			lk.box(uv.costSpec(lk.innerW(w)), w, ch))
	}
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func (uv usageViewState) single(w, h int) []string {
	lk := uv.lk
	iw := lk.innerW(w)
	var sp boxSpec
	switch uv.pane {
	case upForecast:
		sp = uv.forecastSpec(iw)
	case upSpend:
		sp = uv.spendSpec(iw, h)
	case upSavings:
		sp = uv.savingsSpec()
	case upCost:
		sp = uv.costSpec(iw)
	default:
		sp = uv.accountsSpec(iw)
	}
	sp.bottom = fmt.Sprintf("pane %d of %d · tab", uv.pane+1, usagePanes)
	b := lk.box(sp, w, h)
	for len(b) < h {
		b = append(b, blank(w))
	}
	return b[:h]
}
