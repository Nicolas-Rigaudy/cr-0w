package main

// Lines, gauges and the per-theme box frames. Every block is a slice of lines of an
// exact width, so panes can be joined side by side without measuring ANSI again.

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func init() { runewidth.DefaultCondition.EastAsianWidth = false }

type seg struct {
	text  string
	style lipgloss.Style
}
type line []seg

func spc(n int) seg { return seg{strings.Repeat(" ", max(0, n)), lipgloss.NewStyle()} }

func (l line) width() int {
	n := 0
	for _, s := range l {
		n += runewidth.StringWidth(s.text)
	}
	return n
}

// trunc cuts the line to w cells, ending in … when it was cut.
func (l line) trunc(w int) line {
	var out line
	used := 0
	for _, s := range l {
		sw := runewidth.StringWidth(s.text)
		if used+sw <= w {
			out, used = append(out, s), used+sw
			continue
		}
		if w-used > 0 {
			out = append(out, seg{runewidth.Truncate(s.text, w-used, "…"), s.style})
		}
		break
	}
	return out
}

// render truncates to width and pads with spaces to exactly width. A selected line is
// drawn in one style, so glyph colours are lost on it by design.
func (l line) render(width int, selected bool, selStyle lipgloss.Style) string {
	l = l.trunc(width)
	pad := strings.Repeat(" ", max(0, width-l.width()))
	var b strings.Builder
	if selected {
		for _, s := range l {
			b.WriteString(s.text)
		}
		return selStyle.Render(b.String() + pad)
	}
	for _, s := range l {
		b.WriteString(s.style.Render(s.text))
	}
	return b.String() + pad
}

// lr puts left and right content on one line of width w; left is cut to make room.
func lr(left, right line, w int) line {
	rw := right.width()
	if rw == 0 {
		return left.trunc(w)
	}
	if left.width()+rw+1 > w {
		left = left.trunc(max(0, w-rw-1))
	}
	out := append(line{}, left...)
	out = append(out, spc(w-left.width()-rw))
	return append(out, right...)
}

func fit(s string, w int) string {
	s = runewidth.Truncate(s, w, "…")
	return s + strings.Repeat(" ", max(0, w-runewidth.StringWidth(s)))
}

func rep(s string, n int) string { return strings.Repeat(s, max(0, n)) }

// gaugeBar draws a gauge in the theme's shape, readable without colour: lazytf `━` on
// `─`, soft `●` on `○`, cyber `[//..]`, HUD `▕▮▮▯│▯▏` with a tick at 80 %.
func gaugeBar(lk *look, pct, cells int, style lipgloss.Style) line {
	pct = max(0, min(100, pct))
	brackets := (lk.th.frame == frameCyber || lk.th.frame == frameHUD) && cells >= 5
	inner := cells
	if brackets {
		inner -= 2
	}
	on := (pct*inner + 50) / 100
	if inner > 1 {
		if pct > 0 {
			on = max(on, 1)
		}
		if pct < 100 {
			on = min(on, inner-1)
		}
	}
	off := inner - on
	switch lk.th.frame {
	case frameSoft:
		return line{{rep("●", on), style}, {rep("○", off), lk.dim}}
	case frameCyber:
		l := line{{rep("/", on), style}, {rep(".", off), lk.sub}}
		if brackets {
			l = append(line{{"[", lk.sub}}, append(l, seg{"]", lk.sub})...)
		}
		return l
	case frameHUD:
		l := line{{rep("▮", on), style}}
		if tick := (inner*80 + 50) / 100; tick >= on && tick < inner {
			l = append(l, seg{rep("▯", tick-on), lk.dim}, seg{"│", lk.label}, seg{rep("▯", inner-tick-1), lk.dim})
		} else {
			l = append(l, seg{rep("▯", off), lk.dim})
		}
		if brackets {
			l = append(line{{"▕", lk.label}}, append(l, seg{"▏", lk.label})...)
		}
		return l
	}
	return line{{rep("━", on), style}, {rep("─", off), lk.dim}}
}

// levelStyle: green above 50 % remaining, yellow 20 to 50, red below 20.
func (lk *look) level(pct int) lipgloss.Style {
	switch {
	case pct < 20:
		return lk.bad.Bold(true)
	case pct <= 50:
		return lk.warn.Bold(true)
	}
	return lk.ok.Bold(true)
}

type block []string

func blank(w int) string { return strings.Repeat(" ", w) }

// boxSpec is what a pane gives the frame; the theme decides how to draw it.
type boxSpec struct {
	title   string
	index   int  // 1-based panel number (HUD edge)
	right   line // top right: where enter goes, or meta
	bottom  string
	hdr     line // table header, above the rows
	rows    []line
	sel     int // selected row, -1 none
	selN    int // rows the selection spans (0 or 1: one)
	focused bool
	quiet   bool // the selection only scrolls into view, it is not drawn reversed
}

// has reports whether display row i is inside the selection starting at sel.
func (sp boxSpec) has(i, sel int) bool { return sel >= 0 && i >= sel && i < sel+max(1, sp.selN) }

// innerW is the width available to a row inside a box of width w.
func (lk *look) innerW(w int) int { return w - 4 }

// boxH is the height of a box with n body lines (header included).
func (lk *look) boxH(n int) int { return n + 2 }

// gaps is the blank space between panes: soft leaves columns, the boxed themes touch.
func (lk *look) gaps() (h, v int) {
	if lk.th.frame == frameSoft {
		return 2, 0
	}
	return 0, 0
}

// hdrH is the number of lines a box header takes: the header row and the theme's rule.
func (lk *look) hdrH(sp boxSpec) int {
	if sp.hdr == nil {
		return 0
	}
	if lk.headRule(1) != nil {
		return 2
	}
	return 1
}

// headRule is the rule under a table header; lazytf has none.
func (lk *look) headRule(iw int) line {
	switch lk.th.frame {
	case frameSoft, frameHUD:
		return line{{rep("╌", iw), lk.dim}}
	case frameCyber:
		return line{{rep("─", iw), lk.dim}}
	}
	return nil
}

// groupHead is a list group header: sentence case name, count, then the theme leader.
func (lk *look) groupHead(name string, st lipgloss.Style, n, iw int) line {
	l := line{{name, st}, {"  " + strconv.Itoa(n) + " ", lk.sub}}
	lead := map[frameKind]string{frameRound: "─", frameSoft: "╌", frameCyber: "╱", frameHUD: "╌"}[lk.th.frame]
	return append(l, seg{rep(lead, iw-l.width()), lk.dim})
}

// box draws a pane at width w and height h (h < 0: fit the content). Rows scroll so the
// selected one stays visible; the header row stays on top.
func (lk *look) box(sp boxSpec, w, h int) block {
	hdr := lk.hdrH(sp)
	if h < 0 {
		h = lk.boxH(hdr + len(sp.rows))
	}
	n := max(0, h-lk.boxH(0)-hdr) // row slots
	off := 0
	if end := sp.sel + max(1, sp.selN); sp.sel >= 0 && end > n && n > 0 {
		off = min(sp.sel, end-n)
	}
	var lines []line
	if hdr > 0 {
		lines = append(lines, sp.hdr)
	}
	if hdr > 1 {
		lines = append(lines, lk.headRule(lk.innerW(w)))
	}
	rows := sp.rows[min(off, len(sp.rows)):]
	lines = append(lines, rows[:min(n, len(rows))]...)
	for len(lines) < h-lk.boxH(0) {
		lines = append(lines, nil)
	}
	sel := -1
	if sp.sel >= 0 && !sp.quiet {
		sel = sp.sel - off + hdr
	}
	inner := lk.innerW(w)
	switch lk.th.frame {
	case frameHUD:
		return lk.boxHUD(sp, w, inner, lines, sel)
	case frameCyber:
		return lk.boxCyber(sp, w, inner, lines, sel)
	case frameSoft:
		return lk.boxSoft(sp, w, inner, lines, sel)
	}
	return lk.boxRound(sp, w, inner, lines, sel)
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// borderStyle: accent bold when focused, slot 7 otherwise (slot 8 fades on a
// translucent terminal).
func (lk *look) borderStyle(focused bool) lipgloss.Style {
	if focused {
		return lk.acc
	}
	return lk.sub
}

// selMark is the gutter glyph of a selected row; lazytf reverses the row instead.
func (lk *look) selMark() string {
	return map[frameKind]string{frameSoft: "▍", frameCyber: "▌", frameHUD: "▶"}[lk.th.frame]
}

// row draws a body row with its one-cell gutter: width inner+1. A selected row is
// reversed in lazytf; the other themes put a mark in the gutter and set the row bold.
func (lk *look) row(l line, inner int, selected, focused bool) string {
	mark := lk.selMark()
	if !selected || mark == "" {
		return " " + l.render(inner, selected, lk.sel)
	}
	ms := lk.acc
	if !focused {
		ms = lk.sub
	}
	return ms.Render(mark) + lk.emph(l, focused).render(inner, false, lk.sel)
}

// emph restyles a selected row: soft bold, cyber bold with the first word underlined,
// HUD all in the accent. Uncoloured text goes to slot 15.
func (lk *look) emph(l line, focused bool) line {
	out := make(line, 0, len(l)+1)
	under := lk.th.frame == frameCyber
	for _, s := range l {
		st := s.style.Bold(true)
		if _, none := st.GetForeground().(lipgloss.NoColor); none {
			st = st.Foreground(slot(15))
		}
		if lk.th.frame == frameHUD && focused {
			st = lk.acc
		}
		if under && strings.ContainsFunc(s.text, isWordRune) {
			word := strings.TrimRight(s.text, " ")
			out = append(out, seg{word, st.Underline(true)}, seg{s.text[len(word):], st})
			under = false
			continue
		}
		out = append(out, seg{s.text, st})
	}
	return out
}

func isWordRune(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '!'
}

// titleSegs splits a trailing count off the title and draws both the theme's way:
// `Next 9`, `NEXT [09]`, `NEXT ×09`.
func (lk *look) titleSegs(title string, ts lipgloss.Style) line {
	name, n := title, ""
	if i := strings.LastIndex(title, " "); i > 0 && isCount(title[i+1:]) {
		name, n = title[:i], title[i+1:]
	}
	l := line{{lk.cap(name), ts}}
	if n == "" {
		return l
	}
	if len(n) < 2 {
		n = "0" + n
	}
	switch lk.th.frame {
	case frameCyber:
		return append(l, seg{" [" + n + "]", lk.label})
	case frameHUD:
		return append(l, seg{" ×" + n, lk.label})
	}
	return append(l, seg{" " + strings.TrimLeft(n, "0"), lk.sub})
}

func isCount(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// ----- lazytf: rounded border, title on the border -----

func (lk *look) boxRound(sp boxSpec, w, inner int, lines []line, sel int) block {
	bs := lk.borderStyle(sp.focused)
	title := lk.titleSegs(sp.title, lk.title)
	top := bs.Render("╭─ ") + title.render(title.width(), false, lk.sel) + " "
	used := 3 + title.width() + 1
	end := bs.Render("╮")
	if sp.right != nil && used+sp.right.width()+4 <= w {
		end = " " + sp.right.render(sp.right.width(), false, lk.sel) + bs.Render(" ─╮")
		used += sp.right.width() + 4
	} else {
		used++
	}
	out := block{top + bs.Render(rep("─", w-used)) + end}
	for i, l := range lines {
		out = append(out, bs.Render("│")+lk.row(l, inner, sp.has(i, sel), sp.focused)+" "+bs.Render("│"))
	}
	bot := bs.Render("╰" + rep("─", w-2) + "╯")
	if sp.bottom != "" && len(sp.bottom)+6 <= w {
		lab := " " + sp.bottom + " "
		bot = bs.Render("╰"+rep("─", w-4-runewidth.StringWidth(lab))) + lk.sub.Render(lab) + bs.Render("─╯")
	}
	return append(out, bot)
}

// ----- HUD: corner brackets, ◢ TITLE ◣ tag, panel index on the bottom edge -----

func (lk *look) boxHUD(sp boxSpec, w, inner int, lines []line, sel int) block {
	bs := lk.borderStyle(sp.focused)
	edge := " "
	if sp.focused {
		edge = "╌"
	}
	// top: ┏━ ◢ TITLE ×09 ◣ ╌╌╌ ╌╌ LABEL ━┓
	tag := line{{"◢ ", bs}}
	tag = append(tag, lk.titleSegs(sp.title, map[bool]lipgloss.Style{true: lk.acc, false: lk.bold}[sp.focused])...)
	tag = append(tag, seg{" ◣", bs})
	left := append(line{{"┏━", bs}, spc(1)}, tag...)
	left = append(left, spc(1))
	right := line{{"━┓", bs}}
	if sp.right != nil {
		if lab := (line{{"╌╌ ", lk.dim}, {lk.cap(plainText(sp.right)), lk.label}, spc(1), {"━┓", bs}}); left.width()+lab.width()+2 <= w {
			right = lab
		}
	}
	top := append(append(line{}, left...), seg{rep(edge, w-left.width()-right.width()), lk.dim})
	top = append(top, right...)
	out := block{top.render(w, false, lk.sel)}
	for i, l := range lines {
		side := " "
		if len(lines) >= 3 && (i == 0 || i == len(lines)-1) {
			side = "┃"
		}
		out = append(out, bs.Render(side)+lk.row(l, inner, sp.has(i, sel), sp.focused)+" "+bs.Render(side))
	}
	// bottom: ┗━ 01 ╌╌╌ 1/5 ━┛
	bl := line{{"┗━", bs}, spc(1), {twoDigits(sp.index), lk.label}, spc(1)}
	if sp.index == 0 {
		bl = line{{"┗━", bs}, spc(1)}
	}
	br := line{{"━┛", bs}}
	if lab := lk.cap(sp.bottom); lab != "" && bl.width()+runewidth.StringWidth(lab)+5 <= w {
		br = line{spc(1), {lab, lk.sub}, spc(1), {"━┛", bs}}
	}
	bot := append(append(line{}, bl...), seg{rep(edge, w-bl.width()-br.width()), lk.dim})
	bot = append(bot, br...)
	return append(out, bot.render(w, false, lk.sel))
}

func twoDigits(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

func plainText(l line) string {
	var b strings.Builder
	for _, s := range l {
		b.WriteString(s.text)
	}
	return b.String()
}

// ----- Cyber: square single rules, double when focused, ▌// TITLE //, ◢ corner -----

func (lk *look) boxCyber(sp boxSpec, w, inner int, lines []line, sel int) block {
	bs := lk.borderStyle(sp.focused)
	set := []string{"┌", "─", "┐", "│", "└"}
	ts := lk.bold
	if sp.focused {
		set, ts = []string{"╔", "═", "╗", "║", "╚"}, lk.acc
	}
	title := append(line{{"▌//", lk.acc}, spc(1)}, lk.titleSegs(sp.title, ts)...)
	title = append(title, spc(1), seg{"//", lk.sub})
	top := bs.Render(set[0]+set[1]+" ") + title.render(title.width(), false, lk.sel) + " "
	used := 3 + title.width() + 1
	end := bs.Render(set[2])
	if sp.right != nil && used+sp.right.width()+7 <= w {
		lab := lk.cap(plainText(sp.right))
		end = " " + lk.dim.Render("░▒") + " " + lk.label.Render(lab) + bs.Render(" "+set[1]+set[2])
		used += runewidth.StringWidth(lab) + 7
	} else {
		used++
	}
	out := block{top + bs.Render(rep(set[1], w-used)) + end}
	for i, l := range lines {
		out = append(out, bs.Render(set[3])+lk.row(l, inner, sp.has(i, sel), sp.focused)+" "+bs.Render(set[3]))
	}
	corner := lk.sub.Render("◢")
	if sp.focused {
		corner = lk.acc.Render("◢")
	}
	bot := bs.Render(set[4]+rep(set[1], w-2)) + corner
	if sp.bottom != "" && len(sp.bottom)+8 <= w {
		lab := " " + lk.cap(sp.bottom) + " "
		bot = bs.Render(set[4]+rep(set[1], w-4-runewidth.StringWidth(lab))) + lk.sub.Render(lab) + bs.Render(set[1]+set[1]) + corner
	}
	return append(out, bot)
}

// ----- Soft: a top rule with the title in it and a left spine, no right or bottom edge -----

func (lk *look) boxSoft(sp boxSpec, w, inner int, lines []line, sel int) block {
	rule, spine, rs := "─", "▎", lk.sub
	ts := lk.bold
	if sp.focused {
		rule, spine, rs, ts = "━", "▌", lk.acc, lk.acc
	}
	left := append(line{{"╶" + rule, rs}, spc(1)}, lk.titleSegs(sp.title, ts)...)
	left = append(left, spc(1))
	var right line
	if sp.right != nil && left.width()+sp.right.width()+3 <= w {
		right = append(line{spc(1)}, sp.right...)
	}
	top := append(append(line{}, left...), seg{rep(rule, w-left.width()-right.width()), rs})
	out := block{append(top, right...).render(w, false, lk.sel)}
	for i, l := range lines {
		out = append(out, rs.Render(spine)+lk.row(l, inner, sp.has(i, sel), sp.focused)+"  ")
	}
	foot := line{{spine, rs}}
	if sp.bottom != "" {
		foot = lr(foot, line{{sp.bottom, lk.sub}, spc(2)}, w)
	}
	return append(out, foot.render(w, false, lk.sel))
}

func sp2(n int) string { return strings.Repeat(" ", n) }

// ----- joining blocks -----

// hcat puts blocks side by side with gap columns between; shorter blocks get blank lines.
func hcat(widths []int, gap int, blocks ...block) block {
	h := 0
	for _, b := range blocks {
		h = max(h, len(b))
	}
	out := make(block, h)
	for i := range out {
		var parts []string
		for j, b := range blocks {
			if i < len(b) {
				parts = append(parts, b[i])
			} else {
				parts = append(parts, blank(widths[j]))
			}
		}
		out[i] = strings.Join(parts, blank(gap))
	}
	return out
}

func vcat(blocks ...block) block {
	var out block
	for _, b := range blocks {
		out = append(out, b...)
	}
	return out
}
