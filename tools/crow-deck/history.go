package main

// History page, behind `/` in Agents: a search box with filter chips, the sessions grouped
// by day, and a preview of the selected one. Read only: resume, transcript, retro, copy id
// and handoff note arrive with the actions slice.

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	hpList = iota
	hpPreview
	histPanes
)

var histPaneNames = [histPanes]string{"Sessions", "Preview"}

type histView struct {
	lk     *look
	fx     *fixtures
	now    time.Time
	sync   string
	mode   string
	f      histFilter
	typing bool
	sel    int
	pane   int
}

func (hv histView) loading() bool { return hv.sync == stLoading }

func (hv histView) rows() []histSession { return visibleHist(hv.fx, hv.f, hv.mode, hv.now) }

func (hv histView) cur(rows []histSession) *histSession {
	if len(rows) == 0 {
		return nil
	}
	return &rows[max(0, min(len(rows)-1, hv.sel))]
}

// noSource: the live usage source never answered, so there is nothing to list.
func (hv histView) noSource() bool { return hv.fx.Live && hv.fx.Hist == nil && hv.sync == stError }

func (hv histView) total() int { return max(hv.fx.HistTotal, len(hv.fx.Hist)) }

func (hv histView) chip(label, val, any string) line {
	st := hv.lk.sub
	text := label + " " + any
	if val != "" {
		st, text = hv.lk.label.Bold(true), label+" "+val
	}
	return line{{text, st}}
}

func (hv histView) searchSpec(rows []histSession, iw int) boxSpec {
	lk := hv.lk
	sp := boxSpec{title: "Search", index: 1, sel: -1, focused: hv.typing}
	sp.right = line{{"/ text · p project · a account · d date · M model · s status", lk.sub}}
	q := line{{"› ", lk.ok}}
	switch {
	case hv.f.query != "":
		q = append(q, seg{hv.f.query, lk.plain})
	case !hv.typing:
		q = append(q, seg{"press / to search the session titles", lk.sub})
	}
	if hv.typing {
		q = append(q, seg{"▍", lk.warn})
	}
	count := fmt.Sprintf("%d of %d sessions", len(rows), hv.total())
	if hv.loading() {
		count = "syncing…"
	}
	chips := line{}
	for i, c := range []line{
		hv.chip("project", hv.f.project, "any"), hv.chip("account", hv.f.account, "any"),
		{{hv.f.dayLabel(), lk.label.Bold(true)}}, hv.chip("model", hv.f.model, "any"), hv.chip("status", hv.f.status, "any"),
	} {
		if i > 0 {
			chips = append(chips, spc(2))
		}
		chips = append(chips, c...)
	}
	clear := line{{"x", lk.key}, spc(1), {"clear filters", lk.sub}}
	if !hv.f.active() {
		clear = line{{"x clear filters", lk.sub}}
	}
	sp.rows = []line{lr(q, line{{count, lk.sub}}, iw), lr(chips, clear, iw)}
	return sp
}

// listCols: the table columns that fit; portrait drops account and model (the preview has them).
type histCols struct{ account, model bool }

func histColsFor(iw int) histCols { return histCols{account: iw >= 70, model: iw >= 82} }

const (
	hTime, hProj, hAcct, hModel, hCost, hTook = 5, 9, 7, 7, 6, 6
)

func (hv histView) listSpec(iw, h int, bottomNote string) boxSpec {
	lk := hv.lk
	sp := boxSpec{title: "Sessions", index: 2, sel: -1, focused: hv.pane == hpList && !hv.typing}
	if hv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	rows := hv.rows()
	switch {
	case hv.noSource():
		sp.rows = []line{{{"No session history: the usage source did not answer.", lk.bad}}}
		return sp
	case len(hv.fx.Hist) == 0:
		sp.rows = []line{{{"No session yet.", lk.sub}}}
		return sp
	case len(rows) == 0:
		sp.rows = []line{{{"No session matches.", lk.sub}}, {{"x clears the filters", lk.sub}}}
		return sp
	}
	sp.title += fmt.Sprintf(" %d", len(rows))
	sp.right = line{{"newest first · j k", lk.sub}}
	cols := histColsFor(iw)
	hc := []col{{"start", hTime, 1}, {"project", hProj, 1}}
	if cols.account {
		hc = append(hc, col{"account", hAcct, 1})
	}
	if cols.model {
		hc = append(hc, col{"model", hModel, 1})
	}
	hc = append(hc, col{"task", 0, 0})
	sp.hdr = lk.headRow(iw, hc, fmt.Sprintf("%*s %*s", hCost, "cost", hTook, "took"))

	sel := max(0, min(len(rows)-1, hv.sel))
	var all []line
	selRow, day := 0, ""
	for i, s := range rows {
		if d := dayLabel(s.Start, hv.now); d != day {
			day = d
			all = append(all, line{{lk.cap(d), lk.head}})
		}
		if i == sel {
			selRow = len(all)
		}
		g := seg{"✓", lk.ok}
		if s.Active {
			g = seg{"●", lk.ok}
		}
		l := line{g, spc(1), {fit(s.Start.Format("15:04"), hTime), lk.sub}, spc(1), {fit(s.Project, hProj), lk.bold}, spc(1)}
		if cols.account {
			l = append(l, seg{fit(s.Account, hAcct), lk.sub}, spc(1))
		}
		if cols.model {
			l = append(l, seg{fit(s.Model, hModel), lk.sub}, spc(1))
		}
		l = append(l, seg{s.Task, lk.plain})
		right := line{{fmt.Sprintf("%*s", hCost, money(s.Cost)), lk.plain}, spc(1), {fmt.Sprintf("%*s", hTook, fmtTook(s.Took)), lk.sub}}
		all = append(all, lr(l, right, iw))
	}
	cost, took := histTotals(rows)
	foot := line{{fmt.Sprintf("%d sessions · %s · %s active", len(rows), money(cost), fmtTook(took)), lk.sub}}
	note := line{{"api-equivalent, not billed", lk.sub}}

	// Window the rows by hand so the totals stay at the bottom.
	slots := max(1, h-lk.boxH(0)-1-2)
	off := 0
	if selRow >= slots {
		off = selRow - slots + 1
	}
	vis := all[min(off, len(all)):]
	if len(vis) > slots {
		vis = vis[:slots]
	}
	sp.rows = append(sp.rows, vis...)
	sp.rows = append(sp.rows, foot, note)
	sp.sel = selRow - off
	if len(rows) >= 5 {
		sp.bottom = fmt.Sprintf("%d of %d", sel+1, len(rows))
	}
	if bottomNote != "" {
		sp.bottom = bottomNote
	}
	return sp
}

// previewSpec: the selected session. Live rows have only what crow-usage records.
func (hv histView) previewSpec(iw, h int) boxSpec {
	lk := hv.lk
	sp := boxSpec{title: "Preview", index: 3, sel: -1, focused: hv.pane == hpPreview && !hv.typing}
	if hv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	s := hv.cur(hv.rows())
	if s == nil {
		sp.rows = []line{{{"Nothing to show.", lk.sub}}}
		return sp
	}
	sp.title = s.Task
	sp.right = line{{strings.Join([]string{s.Project, s.Account, s.Model, s.Start.Format("15:04")}, " · "), lk.sub}}
	lab := func(name string) seg { return seg{fit(lk.cap(name), 10), lk.head} }
	val := func(parts ...seg) line { return line(parts) }
	var rows []line
	for i, l := range wrap(s.Task, max(10, iw-10)) {
		name := ""
		if i == 0 {
			name = "task"
		}
		rows = append(rows, line{lab(name), {l, lk.plain}})
	}
	glyph, result := seg{"✓", lk.ok}, "done"
	if s.Active {
		glyph, result = seg{"●", lk.ok}, "running"
	}
	meta := fmt.Sprintf(" · %s · %s tokens · %s", fmtTook(s.Took), tokensShort(s.Tokens), money(s.Cost))
	if s.Detail {
		rows = append(rows, line{lab("end state"), {s.EndState, lk.plain}})
	}
	rows = append(rows, val(lab("result"), glyph, seg{" " + result + meta, lk.plain}))
	if s.Branch != "" {
		rows = append(rows, val(lab("branch"), seg{fmt.Sprintf("%s · %d messages", s.Branch, s.Messages), lk.plain}))
	}
	if s.Detail {
		touched := line{lab("touched")}
		for _, f := range s.Touched {
			touched = append(touched, seg{f, lk.plain}, spc(2))
		}
		touched = append(touched, seg{fmt.Sprintf("· %d commits", s.Commits), lk.sub})
		rows = append(rows, touched)
		linked := line{lab("linked")}
		for _, l := range s.Linked {
			linked = append(linked, seg{l, lk.cyan}, spc(2))
		}
		if len(s.Linked) == 0 {
			linked = append(linked, seg{"none", lk.sub})
		}
		rows = append(rows, linked, nil, line{{lk.cap("last messages"), lk.head}})
		for _, t := range s.Last {
			rows = append(rows, line{{fit(t.At, 5), lk.cyan}, spc(1), {t.Text, lk.plain}})
		}
	} else {
		rows = append(rows, nil, line{{"End state, touched files, links and last messages", lk.sub}}, line{{"are not recorded by the sources yet.", lk.sub}})
	}
	hints := line{{"enter resume  o transcript  r retro  y copy id  h handoff note", lk.sub}}
	note := line{{"read only: these arrive with the actions slice", lk.sub}}
	slots := max(0, h-lk.boxH(0)-2)
	if len(rows) > slots {
		rows = rows[:slots]
	}
	for len(rows) < slots {
		rows = append(rows, nil)
	}
	sp.rows = append(rows, hints, note)
	return sp
}

// rowsAll counts display rows of the list: sessions plus one header per day.
func (hv histView) rowsAll() []int {
	var out []int
	day := ""
	for _, s := range hv.rows() {
		if d := dayLabel(s.Start, hv.now); d != day {
			day = d
			out = append(out, -1)
		}
		out = append(out, 0)
	}
	return out
}

func (hv histView) body(w, h int) []string {
	lk := hv.lk
	rows := hv.rows()
	gapV := 0
	if lk.th.frame == frameSoft {
		gapV = 1
	}
	searchH := lk.boxH(2)
	search := lk.box(hv.searchSpec(rows, lk.innerW(w)), w, searchH)
	restH := h - searchH - gapV
	var rest block
	switch layoutFor(w) {
	case layLandscape:
		gap := 0
		if lk.th.frame == frameSoft {
			gap = 3
		}
		w1 := (w - gap) * 60 / 100
		w2 := w - gap - w1
		rest = hcat([]int{w1, w2}, gap,
			lk.box(hv.listSpec(lk.innerW(w1), restH, ""), w1, restH),
			lk.box(hv.previewSpec(lk.innerW(w2), restH), w2, restH))
	case layPortrait:
		listH := min(restH*3/5, max(lk.boxH(4), lk.boxH(1+len(hv.rowsAll())+2)))
		rest = vcat(lk.box(hv.listSpec(lk.innerW(w), listH, ""), w, listH), blankLines(w, gapV),
			lk.box(hv.previewSpec(lk.innerW(w), restH-listH-gapV), w, restH-listH-gapV))
	default:
		pos := fmt.Sprintf("pane %d of %d · tab", hv.pane+1, histPanes)
		if hv.pane == hpPreview {
			sp := hv.previewSpec(lk.innerW(w), restH)
			sp.bottom = pos
			rest = lk.box(sp, w, restH)
		} else {
			rest = lk.box(hv.listSpec(lk.innerW(w), restH, pos), w, restH)
		}
	}
	out := vcat(search, blankLines(w, gapV), rest)
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

// histKey handles the History keys; ok is false for keys the page shares with the rest (tabs,
// mode, theme, help, palette, j k, tab, quit).
func (m model) histKey(msg tea.KeyMsg) (model, bool) {
	k := msg.String()
	if m.histTyping {
		switch msg.Type {
		case tea.KeyEsc, tea.KeyEnter:
			m.histTyping = false
		case tea.KeyBackspace:
			if r := []rune(m.histF.query); len(r) > 0 {
				m.histF.query = string(r[:len(r)-1])
			}
		case tea.KeyRunes:
			m.histF.query += string(msg.Runes)
		case tea.KeySpace:
			m.histF.query += " "
		case tea.KeyCtrlC:
			return m, false
		}
		m.histSel = 0
		return m, true
	}
	switch k {
	case "esc":
		m.hist = false
	case "/":
		m.histTyping, m.histPane = true, hpList
	case "p", "a", "d", "s":
		m.histF.cycle(k, m.fx, m.mode)
		m.histSel = 0
	case "m":
		// m stays the global mode key; the model chip is on M.
		return m, false
	case "M":
		m.histF.cycle("m", m.fx, m.mode)
		m.histSel = 0
	case "x":
		m.histF = newHistFilter()
		m.histSel = 0
	case "pgdown", "ctrl+d":
		m.moveSel(10)
	case "pgup", "ctrl+u":
		m.moveSel(-10)
	case "enter", "o", "r", "y", "h":
		m.flash = k + ": " + map[string]string{"enter": "resume", "o": "transcript", "r": "retro", "y": "copy id", "h": "handoff note"}[k] + " arrives with the actions slice"
	default:
		return m, false
	}
	return m, true
}
