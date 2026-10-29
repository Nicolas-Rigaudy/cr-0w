package main

// Lines, gauges and the per-theme box frames. Every block is a slice of lines of an
// exact width, so panes can be joined side by side without measuring ANSI again.

import (
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

// gaugeBar draws `━` filled in colour on a `━` track in slot 8.
func gaugeBar(lk *look, pct, cells int, style lipgloss.Style) line {
	pct = max(0, min(100, pct))
	on := (pct*cells + 50) / 100
	if cells > 1 {
		if pct > 0 {
			on = max(on, 1)
		}
		if pct < 100 {
			on = min(on, cells-1)
		}
	}
	return line{{rep("━", on), style}, {rep("━", cells-on), lk.dim}}
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
}

// has reports whether display row i is inside the selection starting at sel.
func (sp boxSpec) has(i, sel int) bool { return sel >= 0 && i >= sel && i < sel+max(1, sp.selN) }

// innerW is the width available to a row inside a box of width w.
func (lk *look) innerW(w int) int {
	switch lk.th.frame {
	case frameHUD:
		return w - 3
	case frameSoft:
		return w - 4
	}
	return w - 4
}

// boxH is the height of a box with n body lines (header included).
func (lk *look) boxH(n int) int {
	if lk.th.frame == frameSoft {
		return n + 1
	}
	return n + 2
}

// box draws a pane at width w and height h (h < 0: fit the content). Rows scroll so the
// selected one stays visible; the header row stays on top.
func (lk *look) box(sp boxSpec, w, h int) block {
	hdr := btoi(sp.hdr != nil)
	if h < 0 {
		h = lk.boxH(hdr + len(sp.rows))
	}
	n := max(0, h-lk.boxH(0)-hdr) // row slots
	off := 0
	if end := sp.sel + max(1, sp.selN); sp.sel >= 0 && end > n && n > 0 {
		off = min(sp.sel, end-n)
	}
	var lines []line
	if hdr == 1 {
		lines = append(lines, sp.hdr)
	}
	rows := sp.rows[min(off, len(sp.rows)):]
	lines = append(lines, rows[:min(n, len(rows))]...)
	for len(lines) < h-lk.boxH(0) {
		lines = append(lines, nil)
	}
	sel := -1
	if sp.sel >= 0 {
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

func (lk *look) borderStyle(focused bool) lipgloss.Style {
	if focused {
		return lk.acc
	}
	return lk.dim
}

// ----- lazytf: rounded border, title on the border -----

func (lk *look) boxRound(sp boxSpec, w, inner int, lines []line, sel int) block {
	bs := lk.borderStyle(sp.focused)
	title := line{{sp.title, lk.title}}
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
		out = append(out, bs.Render("│")+" "+l.render(inner, sp.has(i, sel), lk.sel)+" "+bs.Render("│"))
	}
	bot := bs.Render("╰" + rep("─", w-2) + "╯")
	if sp.bottom != "" && len(sp.bottom)+6 <= w {
		lab := " " + sp.bottom + " "
		bot = bs.Render("╰"+rep("─", w-4-runewidth.StringWidth(lab))) + lk.sub.Render(lab) + bs.Render("─╯")
	}
	return append(out, bot)
}

// ----- HUD: brackets, tag title, panel index on the bottom edge -----

func (lk *look) boxHUD(sp boxSpec, w, inner int, lines []line, sel int) block {
	bs := lk.borderStyle(sp.focused)
	tagStyle := lk.bold
	if sp.focused {
		tagStyle = lk.acc
	}
	tag := "◢ " + lk.cap(sp.title) + " ◣"
	top := lk.acc.Render("┏━ ") + tagStyle.Render(tag) + " "
	used := 3 + runewidth.StringWidth(tag) + 1
	lab := ""
	if sp.right != nil {
		lab = lk.cap(plainText(sp.right))
	}
	if lab != "" && used+runewidth.StringWidth(lab)+5 <= w {
		top += lk.dim.Render(rep("─", w-used-runewidth.StringWidth(lab)-5)) + " " + lk.dim.Render("╌╌") + " " + lk.label.Render(lab) + " "
	} else {
		top += lk.dim.Render(rep("─", w-used))
	}
	out := block{top}
	for i, l := range lines {
		out = append(out, bs.Render("▏")+" "+l.render(inner, sp.has(i, sel), lk.sel)+" ")
	}
	idx := " " + twoDigits(sp.index) + " "
	left := "▏▏▏▏" + idx + "▏▏▏▏ "
	bot := bs.Render(left)
	used = runewidth.StringWidth(left)
	if sp.bottom != "" {
		lab := " " + lk.cap(sp.bottom) + " "
		bot += lk.dim.Render(rep("─", w-used-runewidth.StringWidth(lab)-3)) + lk.sub.Render(lab) + lk.acc.Render("━━┛")
	} else {
		bot += lk.dim.Render(rep("─", w-used-3)) + lk.acc.Render("━━┛")
	}
	return append(out, bot)
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

// ----- Cyber: single rules, ▌// TITLE //, ◢ corner -----

func (lk *look) boxCyber(sp boxSpec, w, inner int, lines []line, sel int) block {
	bs := lk.borderStyle(sp.focused)
	ts := lk.bold
	if sp.focused {
		ts = lk.acc
	}
	title := lk.acc.Render("▌//") + " " + ts.Render(lk.cap(sp.title)) + " " + lk.dim.Render("//")
	tw := 3 + 1 + runewidth.StringWidth(lk.cap(sp.title)) + 1 + 2
	top := bs.Render("┌─ ") + title + " "
	used := 3 + tw + 1
	end := bs.Render("┐")
	if sp.right != nil && used+sp.right.width()+7 <= w {
		lab := plainText(sp.right)
		end = " " + lk.dim.Render("░▒") + " " + lk.label.Render(lab) + bs.Render(" ─┐")
		used += runewidth.StringWidth(lab) + 7
	} else {
		used++
	}
	out := block{top + bs.Render(rep("─", w-used)) + end}
	for i, l := range lines {
		out = append(out, bs.Render("│")+" "+l.render(inner, sp.has(i, sel), lk.sel)+" "+bs.Render("│"))
	}
	corner := lk.dim.Render("◢")
	if sp.focused {
		corner = lk.acc.Render("◢")
	}
	bot := bs.Render("└"+rep("─", w-2)) + corner
	if sp.bottom != "" && len(sp.bottom)+8 <= w {
		lab := " " + sp.bottom + " "
		bot = bs.Render("└"+rep("─", w-4-runewidth.StringWidth(lab))) + lk.sub.Render(lab) + bs.Render("──") + corner
	}
	return append(out, bot)
}

// ----- Soft: no lines, sentence-case title row, spacing only -----

func (lk *look) boxSoft(sp boxSpec, w, inner int, lines []line, sel int) block {
	ts := lk.bold
	if sp.focused {
		ts = lk.acc
	}
	right := sp.right
	if sp.bottom != "" {
		right = append(append(line{}, right...), seg{"  " + sp.bottom, lk.sub})
		if sp.right == nil {
			right = line{{sp.bottom, lk.sub}}
		}
	}
	titleRow := lr(line{spc(2), {sp.title, ts}}, right, w)
	out := block{titleRow.render(w, false, lk.sel)}
	for i, l := range lines {
		out = append(out, sp2(2)+l.render(inner, sp.has(i, sel), lk.sel)+sp2(2))
	}
	return out
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
