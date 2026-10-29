package main

// Agents page: one grouped list, the selected agent's live tail with its question and the
// answer keys, and "Who worked when" underneath. Landscape puts list and tail side by
// side; portrait stacks them; under 100 columns one pane shows at a time.

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const (
	apList = iota
	apTail
	apWho
	agentPanes
)

var agentPaneNames = [agentPanes]string{"Agents", "Tail", "Who worked when"}

// agentsView is the data and state an Agents render needs.
type agentsView struct {
	lk      *look
	fx      *fixtures
	now     time.Time
	sync    string
	pane    int
	sel     int
	input   string
	typing  bool
	accents []string // account names in colour order
}

func (av agentsView) loading() bool { return av.sync == stLoading }

func (av agentsView) cur() *agentRow {
	if av.sel >= 0 && av.sel < len(av.fx.Agents) {
		return &av.fx.Agents[av.sel]
	}
	return nil
}

// acctStyle: the first account is green, then yellow, cyan, magenta (the board's legend).
func (av agentsView) acctStyle(account string) lipgloss.Style {
	lk := av.lk
	styles := []lipgloss.Style{lk.ok, lk.warn, lk.cyan, lk.mag}
	for i, a := range av.accents {
		if a == account {
			return styles[i%len(styles)]
		}
	}
	return lk.sub
}

func agentAccounts(fx *fixtures) []string {
	var out []string
	seen := map[string]bool{}
	add := func(a string) {
		if a != "" && !seen[a] {
			seen[a], out = true, append(out, a)
		}
	}
	for _, g := range fx.Gauges {
		add(g.Account)
	}
	for _, a := range fx.Agents {
		add(a.Account)
	}
	for _, w := range fx.Who {
		add(w.Account)
	}
	return out
}

func glyphOf(lk *look, state string) seg {
	switch state {
	case "working":
		return seg{"●", lk.ok}
	case "idle":
		return seg{"○", lk.sub}
	case "needs":
		return seg{"◆", lk.mag}
	case "waiting":
		return seg{"‖", lk.warn}
	}
	return seg{"✓", lk.ok}
}

// listRows builds the grouped list. idx maps a display row to its agent (-1: group header).
const labelW, acctW = 16, 8

func (av agentsView) listRows(iw int) (rows []line, idx []int) {
	lk, fx := av.lk, av.fx
	for g := range groupNames {
		var members []int
		for i, a := range fx.Agents {
			if a.Group == g {
				members = append(members, i)
			}
		}
		if len(members) == 0 && g != gRuns {
			continue
		}
		count := fmt.Sprintf("%d", len(members))
		if g == gRuns {
			count = fmt.Sprintf("runs %d/%d", fx.Slots[0], max(fx.Slots[1], 1))
			if len(members) > 0 && fx.Slots[1] == 0 {
				count = fmt.Sprintf("%d", len(members))
			}
		}
		rows = append(rows, line{{lk.cap(groupNames[g]) + " ", lk.head}, {count, lk.sub}})
		idx = append(idx, -1)
		for _, i := range members {
			a := fx.Agents[i]
			label := cmpOr(a.Short, a.Name)
			left := line{glyphOf(lk, a.State), spc(1), {fit(label, labelW), lk.bold}, spc(1), {fit(a.Account, acctW), lk.sub}, spc(1), {a.Doing, lk.plain}}
			right := line{{a.Since, lk.sub}}
			rows = append(rows, lr(left, right, iw))
			idx = append(idx, i)
		}
	}
	return rows, idx
}

func (av agentsView) listSpec(iw int) boxSpec {
	lk := av.lk
	sp := boxSpec{title: "Agents", index: 1, sel: -1, focused: av.pane == apList}
	if av.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if len(av.fx.Agents) == 0 {
		sp.rows = []line{{{"No agent running.", lk.sub}}}
		return sp
	}
	sp.title += fmt.Sprintf(" %d", len(av.fx.Agents))
	sp.hdr = lk.headRow(iw, []col{{"agent", labelW, 1}, {"account", acctW, 1}, {"doing", 0, 0}}, "since")
	var idx []int
	sp.rows, idx = av.listRows(iw)
	for d, i := range idx {
		if i == av.sel {
			sp.sel = d
		}
	}
	if len(av.fx.Agents) >= 5 {
		sp.bottom = fmt.Sprintf("%d of %d", av.sel+1, len(av.fx.Agents))
	}
	return sp
}

func wrap(s string, w int) []string {
	var out []string
	cur := ""
	for _, word := range strings.Fields(s) {
		if cur != "" && runewidth.StringWidth(cur)+1+runewidth.StringWidth(word) > w {
			out, cur = append(out, cur), word
			continue
		}
		cur = strings.TrimSpace(cur + " " + word)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

// tailSpec: the selected agent's tail, then its question card with the answer keys when
// it waits for the owner.
func (av agentsView) tailSpec(iw, h int) boxSpec {
	lk := av.lk
	sp := boxSpec{title: "Tail", index: 2, sel: -1, focused: av.pane == apTail}
	if av.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	a := av.cur()
	if a == nil {
		sp.rows = []line{{{"Nothing to show.", lk.sub}}}
		return sp
	}
	sp.title = cmpOr(a.Short, a.Name)
	sp.right = line{glyphOf(lk, a.State), spc(1), {a.Meta, lk.sub}}
	var fixed []line
	if a.Group == gNeeds {
		fixed = append(fixed, nil, line{{lk.cap("agent asks"), lk.warn}})
		ask := wrap(a.Ask, iw)
		if len(ask) > 3 {
			ask = append(ask[:2:2], ask[2]+" …")
		}
		for _, l := range ask {
			fixed = append(fixed, line{{l, lk.plain}})
		}
		keys := line{{"a", lk.key}, spc(1), {"keep", lk.sub}, spc(2), {"x", lk.key}, spc(1), {"drop", lk.sub}, spc(2),
			{"e", lk.key}, spc(1), {"edit text", lk.sub}, spc(2), {"c", lk.key}, spc(1), {"chat", lk.sub}, spc(2),
			{"g", lk.key}, spc(1), {"jump to line", lk.sub}}
		fixed = append(fixed, keys)
	}
	if a.Group == gNeeds || a.Group == gRunning {
		prompt := line{{"› ", lk.key}}
		if av.input != "" {
			prompt = append(prompt, seg{av.input, lk.plain}, seg{"▍", lk.acc})
		} else {
			prompt = append(prompt, seg{"type to the agent, enter to send", lk.sub}, seg{" ▍", lk.acc})
		}
		fixed = append(fixed, prompt)
	}
	hints := line{{"enter", lk.key}, spc(1), {"send", lk.sub}}
	if a.Interrupted {
		hints = append(hints, spc(2), seg{"R", lk.key}, spc(1), seg{"relaunch", lk.sub})
	}
	hints = append(hints, spc(2), seg{"attach, pause, stop in", lk.sub}, spc(1), seg{":", lk.key})
	slots := max(0, h-lk.boxH(0)-len(fixed)-1)
	tail := a.Tail
	if len(tail) > slots {
		tail = tail[len(tail)-slots:]
	}
	for i, t := range tail {
		l := line{{fit(t.At, 5), lk.cyan}, spc(1), {t.Text, lk.plain}}
		if i == len(tail)-1 && a.State == "working" {
			l = append(l, spc(1), seg{"▍", lk.acc})
		}
		sp.rows = append(sp.rows, l)
	}
	sp.rows = append(sp.rows, fixed...)
	for len(sp.rows) < h-lk.boxH(0)-1 {
		sp.rows = append(sp.rows, nil)
	}
	sp.rows = append(sp.rows, hints)
	return sp
}

// hourScale returns the hour range the ruler covers: 08 to 18, widened to fit the data.
func (av agentsView) hourScale() (start, end int) {
	start, end = 8, 18
	h := float64(av.now.Hour()) + float64(av.now.Minute())/60
	end = max(end, int(math.Ceil(h))+1)
	for _, w := range av.fx.Who {
		for _, s := range append(append([]hspan{}, w.On...), w.Off...) {
			start = min(start, int(math.Floor(s.From)))
			end = max(end, int(math.Ceil(s.To)))
		}
	}
	return max(0, start), min(24, end)
}

func fmtHours(h float64) string {
	m := int(math.Round(h * 60))
	return fmt.Sprintf("%dh%02d", m/60, m%60)
}

// unionHours is the total length of the union of spans.
func unionHours(spans []hspan) float64 {
	if len(spans) == 0 {
		return 0
	}
	cp := append([]hspan{}, spans...)
	for i := range cp {
		for j := i + 1; j < len(cp); j++ {
			if cp[j].From < cp[i].From {
				cp[i], cp[j] = cp[j], cp[i]
			}
		}
	}
	total, cur := 0.0, cp[0]
	for _, s := range cp[1:] {
		if s.From <= cur.To {
			cur.To = math.Max(cur.To, s.To)
			continue
		}
		total += cur.To - cur.From
		cur = s
	}
	return total + cur.To - cur.From
}

const whoNameW, whoAcctW = 18, 8

func (av agentsView) whoSpec(iw, maxRows int) boxSpec {
	lk := av.lk
	sp := boxSpec{title: "Who worked when", index: 3, sel: -1, focused: av.pane == apWho}
	sp.right = line{{"today · ", lk.sub}, {"▼ now", lk.label}}
	if av.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if len(av.fx.Who) == 0 {
		sp.rows = []line{{{"Nobody worked yet today.", lk.sub}}}
		return sp
	}
	start, end := av.hourScale()
	barW := max(10, iw-whoNameW-whoAcctW-2)
	span := float64(end - start)
	x := func(h float64) int { return int(math.Round((h - float64(start)) / span * float64(barW))) }
	pad := whoNameW + whoAcctW + 2
	ruler := []rune(rep(" ", barW))
	ticks := []rune(rep("·", barW))
	for hr := start; hr <= end; hr++ {
		p := x(float64(hr))
		if p < barW {
			ticks[p] = '│'
			label := fmt.Sprintf("%02d", hr)
			if p+2 <= barW && hr < end {
				copy(ruler[p:], []rune(label))
			}
		}
	}
	sp.rows = append(sp.rows, line{spc(pad), {string(ruler), lk.sub}}, line{spc(pad), {string(ticks), lk.dim}})
	nowH := float64(av.now.Hour()) + float64(av.now.Minute())/60
	nowMark := fmt.Sprintf("▼ %02d:%02d", av.now.Hour(), av.now.Minute())
	np := min(max(0, x(nowH)), max(0, barW-len(nowMark)))
	sp.rows = append(sp.rows, line{spc(pad + np), {nowMark, lk.mag}})
	rows := av.fx.Who
	if len(rows) > maxRows {
		rows = rows[len(rows)-maxRows:]
	}
	perAcct := map[string][]hspan{}
	for _, w := range av.fx.Who {
		perAcct[w.Account] = append(perAcct[w.Account], w.On...)
	}
	for _, w := range rows {
		kind := make([]int, barW) // 0 idle, 1 paused, 2 active
		paint := func(spans []hspan, k int) {
			for _, s := range spans {
				from, to := max(0, x(s.From)), min(barW, max(x(s.To), x(s.From)+1))
				for i := from; i < to; i++ {
					kind[i] = k
				}
			}
		}
		paint(w.Off, 1)
		paint(w.On, 2)
		l := line{{fit(w.Name, whoNameW), lk.bold}, spc(1), {fit(w.Account, whoAcctW), lk.sub}, spc(1)}
		for i := 0; i < barW; {
			j := i
			for j < barW && kind[j] == kind[i] {
				j++
			}
			switch kind[i] {
			case 2:
				l = append(l, seg{rep("━", j-i), av.acctStyle(w.Account)})
			case 1:
				l = append(l, seg{rep("╌", j-i), lk.dim})
			default:
				l = append(l, spc(j-i))
			}
			i = j
		}
		sp.rows = append(sp.rows, l)
	}
	legend := line{}
	for _, a := range av.accents {
		if len(perAcct[a]) > 0 {
			legend = append(legend, seg{"━", av.acctStyle(a)}, spc(1), seg{a, lk.sub}, spc(2))
		}
	}
	legend = append(legend, seg{"╌", lk.dim}, spc(1), seg{"paused", lk.sub})
	var tot []string
	for _, a := range av.accents {
		if hrs := unionHours(perAcct[a]); hrs > 0 {
			tot = append(tot, fmt.Sprintf("%s %s", a, fmtHours(hrs)))
		}
	}
	sp.rows = append(sp.rows, lr(legend, line{{strings.Join(tot, " · "), lk.sub}}, iw))
	return sp
}

func (av agentsView) whoHeight() int {
	rows := min(len(av.fx.Who), 8)
	if av.loading() || len(av.fx.Who) == 0 {
		return av.lk.boxH(1)
	}
	return av.lk.boxH(3 + rows + 1)
}

// body draws the page at w x h.
func (av agentsView) body(w, h int) []string {
	lk := av.lk
	lay := layoutFor(w)
	if lay == laySingle {
		return av.single(w, h)
	}
	gapV := 0
	if lk.th.frame == frameSoft {
		gapV = 1
	}
	whoH := min(av.whoHeight(), max(lk.boxH(1), h/3))
	topH := h - whoH - gapV
	var top block
	if lay == layLandscape {
		gap := 0
		if lk.th.frame == frameSoft {
			gap = 3
		}
		w1 := (w - gap) * 36 / 100
		w2 := w - gap - w1
		top = hcat([]int{w1, w2}, gap, av.listBox(w1, topH), av.tailBox(w2, topH))
	} else {
		listH := min(topH*2/5, lk.boxH(len(av.listRowsOnly())+1))
		listH = max(listH, lk.boxH(3))
		top = vcat(av.listBox(w, listH), blankLines(w, gapV), av.tailBox(w, topH-listH-gapV))
	}
	out := top
	if gapV > 0 {
		out = append(out, blank(w))
	}
	out = append(out, av.whoBox(w, whoH)...)
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func blankLines(w, n int) block {
	var b block
	for i := 0; i < n; i++ {
		b = append(b, blank(w))
	}
	return b
}

func (av agentsView) listRowsOnly() []line {
	r, _ := av.listRows(80)
	return r
}

func (av agentsView) listBox(w, h int) block {
	return av.lk.box(av.listSpec(av.lk.innerW(w)), w, h)
}

func (av agentsView) tailBox(w, h int) block {
	return av.lk.box(av.tailSpec(av.lk.innerW(w), h), w, h)
}

func (av agentsView) whoBox(w, h int) block {
	lk := av.lk
	rows := max(1, h-lk.boxH(0)-4)
	return lk.box(av.whoSpec(lk.innerW(w), rows), w, h)
}

func (av agentsView) single(w, h int) []string {
	lk := av.lk
	pos := fmt.Sprintf("pane %d of %d · tab", av.pane+1, agentPanes)
	var b block
	switch av.pane {
	case apTail:
		sp := av.tailSpec(lk.innerW(w), h)
		sp.bottom = pos
		b = lk.box(sp, w, h)
	case apWho:
		sp := av.whoSpec(lk.innerW(w), max(1, h-lk.boxH(0)-4))
		sp.bottom = pos
		b = lk.box(sp, w, h)
	default:
		sp := av.listSpec(lk.innerW(w))
		sp.bottom = pos
		b = lk.box(sp, w, h)
	}
	for len(b) < h {
		b = append(b, blank(w))
	}
	return b[:h]
}
