package main

// MRs page, a walkthrough per MR: the list grouped by whose move, the stage line, the code
// of the current file that scrolls in place with the commented line marked, and the Your
// move panel (steps, the current finding, its fix, the decision keys, then apply and push).
// Landscape puts list | stage + code | panel side by side; portrait stacks stage, code and a
// compact panel; under 100 columns one pane shows at a time.

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const (
	mpList = iota
	mpDiff
	mpCard
	mrPanes
)

var mrPaneNames = [mrPanes]string{"MRs and PRs", "Code", "Your move"}

// mrsView is the data and state an MRs render needs.
type mrsView struct {
	lk     *look
	fx     *fixtures
	now    time.Time
	sync   string
	mode   string
	sel    int // index in the visible list
	pane   int
	file   int
	item   int
	scroll int    // first code row the owner scrolled to, -1: follow the card
	toggle bool   // T: the other kind of card (threads on a review, findings on an own MR)
	typing string // "e" or "c" while the owner types on the current finding
	input  string
}

func (mv mrsView) loading() bool {
	return mv.sync == stLoading || !mv.fx.MRsReady && len(mv.fx.ForgeErrs) == 0 && mv.fx.Live
}

// inMode reports whether a scope shows in the current mode.
func inMode(mode, scope string) bool {
	return mode == "all" || (mode == "perso") == (scope == "personal")
}

func (mv mrsView) visible() []mrRow { return visibleMRs(mv.fx, mv.mode) }

// visibleMRs lists the MRs of the mode grouped by whose move: yours, waiting, done.
func visibleMRs(fx *fixtures, mode string) []mrRow {
	var out []mrRow
	for g := 0; g < 3; g++ {
		for _, r := range fx.MRs {
			if inMode(mode, r.Scope) && mrStageOf(r, fx.MRDetails[r.key()], time.Time{}).group() == g {
				out = append(out, r)
			}
		}
	}
	return out
}

func (mv mrsView) cur() *mrRow {
	v := mv.visible()
	if mv.sel >= 0 && mv.sel < len(v) {
		return &v[mv.sel]
	}
	return nil
}

func (mv mrsView) detail() *mrDetail {
	if r := mv.cur(); r != nil {
		return mv.fx.MRDetails[r.key()]
	}
	return nil
}

// threadsMode: own MRs show their threads, reviews their findings; T flips it. Once a threads
// review exists, its findings (one per thread, with fix and reply) replace the raw threads.
func (mv mrsView) threadsMode() bool {
	r := mv.cur()
	if r == nil {
		return false
	}
	if d := mv.detail(); d != nil && d.Review != nil && d.Review.Mode == "threads" {
		return mv.toggle
	}
	return (r.Role == "AUTHOR") != mv.toggle
}

// mode2 is the mode of the saved review, else the one r would start (peer, self, threads).
func (mv mrsView) mode2() string {
	r := mv.cur()
	if r == nil {
		return ""
	}
	return detailMode(r, mv.detail())
}

// answering: the findings shown are the answers to the reviewers' threads.
func (mv mrsView) answering() bool {
	d := mv.detail()
	return !mv.threadsMode() && d != nil && d.Review != nil && d.Review.Mode == "threads"
}

// detailMode: the saved review's mode (self or peer by role for a file older than modes),
// else the mode r would start.
func detailMode(r *mrRow, d *mrDetail) string {
	switch {
	case d == nil || d.Review == nil:
		return reviewMode(r, d)
	case d.Review.Mode != "":
		return d.Review.Mode
	case r.Role == "AUTHOR":
		return "self"
	}
	return "peer"
}

// items counts the cards of the right box.
func (mv mrsView) items() int {
	d := mv.detail()
	if d == nil {
		return 0
	}
	if mv.threadsMode() {
		return len(d.Threads)
	}
	return len(d.Findings)
}

// ----- components per theme -----

// chip is the role label: bold coloured text, with the theme's mark.
func (mv mrsView) chip(role string) line {
	lk := mv.lk
	st := lk.bold
	if role == "AUTHOR" {
		st = lk.ok.Bold(true)
	}
	switch lk.th.frame {
	case frameHUD:
		return line{{"▏", lk.acc}, {role, st}}
	case frameCyber:
		return line{{"[", lk.dim}, {role, st}, {"]", lk.dim}}
	}
	return line{{role, st}}
}

// button is `a accept` with the key in the key colour; the active one is bold.
func (mv mrsView) button(key, word string, on bool) line {
	lk := mv.lk
	ks, ws := lk.key.Bold(true), lk.sub
	if on {
		ks, ws = lk.acc, lk.bold
	}
	return line{{key, ks}, spc(1), {lk.cap(word), ws}}
}

func (mv mrsView) buttons(active int, labels ...[2]string) line {
	var l line
	for i, b := range labels {
		if i > 0 {
			l = append(l, spc(2))
		}
		l = append(l, mv.button(b[0], b[1], i == active)...)
	}
	return l
}

func (mv mrsView) divider(iw int) line {
	lk := mv.lk
	if lk.th.frame == frameSoft {
		return nil
	}
	return line{{rep("─", iw), lk.dim}}
}

func (mv mrsView) groupHdr(name string) line {
	return line{{mv.lk.cap(name), mv.lk.head}}
}

// replyBlock draws saved reply text with the theme's bar; no fill in any theme.
func (mv mrsView) replyBlock(text string, iw int) []line {
	lk := mv.lk
	bar := seg{"  ", lk.plain}
	switch lk.th.frame {
	case frameHUD:
		bar = seg{"▏ ", lk.acc}
	case frameCyber:
		bar = seg{"▌ ", lk.acc}
	case frameRound:
		bar = seg{"│ ", lk.dim}
	}
	var out []line
	for _, l := range wrapLines(text, max(10, iw-2)) {
		out = append(out, line{bar, {l, lk.plain}})
	}
	return out
}

// ----- list -----

// stage is the walkthrough stage of an MR (mrs_move.go).
func (mv mrsView) stage(r mrRow) mrStage {
	return mrStageOf(r, mv.fx.MRDetails[r.key()], mv.now)
}

var mrGroupNames = [3]string{"Your move", "Waiting on others", "Done"}

// listRows builds the list grouped by whose move. start[i] is the first display row of visible MR i.
func (mv mrsView) listRows(iw int) (rows []line, start []int) {
	lk, fx := mv.lk, mv.fx
	vis := mv.visible()
	start = make([]int, len(vis))
	for g := 0; g < 3; g++ {
		var idx []int
		for i, r := range vis {
			if mv.stage(r).group() == g {
				idx = append(idx, i)
			}
		}
		if len(idx) == 0 {
			continue
		}
		gs := []lipgloss.Style{lk.warn.Bold(true), lk.sub.Bold(true), lk.ok.Bold(true)}[g]
		rows = append(rows, lk.groupHead(mrGroupNames[g], gs, len(idx), iw))
		for _, i := range idx {
			r := vis[i]
			st := mv.stage(r)
			mark, ms, ts, rs := "● ", lk.warn, lk.bold, lk.warn.Bold(true)
			switch g {
			case 1:
				mark, ms, ts, rs = "○ ", lk.dim, lk.sub, lk.sub
			case 2:
				mark, ms, ts, rs = "✓ ", lk.ok, lk.plain, lk.ok
			}
			if st.Problem {
				ms, rs = lk.bad, lk.bad.Bold(true)
			}
			cur := ifs(i == mv.sel, "› ") + ifs(i != mv.sel, "  ")
			title := r.Title
			if mv.mode == "all" {
				title = r.Project + " · " + title
			}
			start[i] = len(rows)
			rows = append(rows, line{{cur, lk.acc}, {mark, ms}, {r.Ref + " ", ts.Bold(true)}, {title, ts.Bold(i == mv.sel)}})
			age := ""
			if t, ok := parseTime(r.Updated); ok {
				age = shortDays(mv.now.Sub(t))
			}
			rows = append(rows, lr(line{spc(4), {st.Reason, rs}}, line{{age, lk.sub}}, iw))
		}
	}
	var branches []mrBranch
	for _, b := range fx.Branches {
		if inMode(mv.mode, b.Scope) {
			branches = append(branches, b)
		}
	}
	if len(branches) > 0 {
		hdr := line{{lk.cap("No MR yet"), lk.sub.Bold(true)}, {"  " + strconv.Itoa(len(branches)) + " ", lk.sub}}
		rows = append(rows, append(hdr, seg{rep("╌", iw-hdr.width()), lk.dim}))
		for _, b := range branches {
			rows = append(rows, line{spc(2), {"○ ", lk.dim}, {b.Project + " ", lk.sub.Bold(true)}, {b.Branch, lk.sub}},
				line{spc(4), {fmt.Sprintf("%d %s ahead", b.Ahead, plural(b.Ahead, "commit", "commits")), lk.sub}, spc(1), {"M", lk.key}, spc(1), {"open MR", lk.sub}})
		}
	}
	if len(rows) == 0 {
		rows = append(rows, line{{"Nothing open: no MR or PR is yours.", lk.sub}})
	}
	return rows, start
}

func (mv mrsView) yourCount() int {
	n := 0
	for _, r := range mv.visible() {
		if mv.stage(r).Who == moveYou {
			n++
		}
	}
	return n
}

func (mv mrsView) listSpec(iw int) boxSpec {
	lk := mv.lk
	sp := boxSpec{title: "MRs", index: 1, sel: -1, selN: 2, quiet: true, focused: mv.pane == mpList}
	switch {
	case mv.loading():
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	case len(mv.fx.ForgeErrs) > 0 && len(mv.fx.MRs) == 0:
		for _, e := range mv.fx.ForgeErrs {
			sp.rows = append(sp.rows, line{{"◆ " + e, lk.bad}})
		}
		return sp
	}
	vis := mv.visible()
	if len(vis) > 0 {
		sp.right = line{{fmt.Sprintf("%d/%d your move", mv.yourCount(), len(vis)), lk.sub}}
	}
	var start []int
	sp.rows, start = mv.listRows(iw)
	for _, e := range mv.fx.ForgeErrs {
		sp.rows = append(sp.rows, line{{"◆ " + e, lk.bad}})
	}
	if mv.sel >= 0 && mv.sel < len(start) {
		sp.sel = start[mv.sel]
	}
	if len(vis) >= 4 {
		sp.bottom = fmt.Sprintf("%d of %d", mv.sel+1, len(vis))
	}
	return sp
}

// ----- cards: the findings of a review, or the raw threads of an own MR -----

// mrCard is what the code pane marks and the Your move panel shows.
type mrCard struct {
	Label     string // F1, T1
	Path      string
	Line      int
	Who, When string
	Text      string
	State     string // the word shown: pending, accepted, fixed, pushed, replied, resolved…
	Fix, Test string // diff text
	Reply     string // what goes to the forge with accept
	Note      string // the owner's edited reply or chat note
	Applied   string
	AppDetail string
	AppStale  bool
	Disagree  bool
	Agree     string // threads review: yes, partly, no
	Why       string
	Pattern   string
	Seen      int
	Severity  string
	decided   string // raw decision state: pending, accepted, refused, edited, chat
}

func sectionText(s *mrSection) string {
	if s == nil {
		return ""
	}
	var b strings.Builder
	for _, l := range s.Lines {
		if l.Kind == '@' {
			continue
		}
		b.WriteString(string(l.Kind) + l.Text + "\n")
	}
	return b.String()
}

func (mv mrsView) cards() []mrCard {
	d := mv.detail()
	if d == nil {
		return nil
	}
	var out []mrCard
	if mv.threadsMode() {
		for i, t := range d.Threads {
			out = append(out, mrCard{Label: fmt.Sprintf("T%d", i+1), Path: t.Path, Line: t.Line, Who: t.Author, When: mv.when(t.At), Text: t.Body,
				State: ifs(t.Resolved, "resolved") + ifs(!t.Resolved, "pending"), Fix: sectionText(t.Fix), Test: sectionText(t.Test),
				Reply: t.Reply, Pattern: t.Pattern, Seen: t.Seen, decided: "pending"})
		}
		return out
	}
	for i, f := range d.Findings {
		p, l := f.at()
		c := mrCard{Label: fmt.Sprintf("F%d", i+1), Path: p, Line: l, Text: f.Body, Fix: f.Fix, Test: f.Test, Reply: f.Comment,
			Note: f.Note, Pattern: f.Pattern, Seen: f.Seen, Severity: f.Severity, decided: f.State,
			Applied: f.Applied, AppDetail: f.AppliedDetail, AppStale: f.AppliedFor != [2]string{f.State, f.Note}}
		if mv.answering() {
			c.Who, c.Text, c.Why, c.Disagree, c.Agree = cmpOr(f.ThreadAuthor, "reviewer"), cmpOr(f.ThreadBody, f.Body), f.Body, f.Agree == "no", f.Agree
		} else {
			c.Who = cmpOr(f.Severity, "finding")
			if f.Title != "" {
				c.Text = f.Title + ". " + f.Body
			}
		}
		c.State = f.State
		switch {
		case f.Resolved:
			c.State = "resolved"
		case f.Posted && mv.answering():
			c.State = "replied"
		case f.Posted:
			c.State = "posted"
		case f.Applied == "fixed" && d.Apply != nil && d.Apply.Pushed:
			c.State = "pushed"
		case f.Applied != "" && !c.AppStale:
			c.State = f.Applied
		}
		out = append(out, c)
	}
	return out
}

func (mv mrsView) curCard() *mrCard {
	cs := mv.cards()
	if len(cs) == 0 {
		return nil
	}
	c := cs[mv.itemIdx(len(cs))]
	return &c
}

// ----- stage line -----

func (mv mrsView) rail(st mrStage) line {
	lk := mv.lk
	var l line
	for i, s := range st.Steps {
		if i > 0 {
			l = append(l, seg{" ─ ", lk.dim})
		}
		switch {
		case i < st.Cur:
			l = append(l, seg{"✓ " + s.Name, lk.ok})
		case i == st.Cur:
			l = append(l, seg{"● " + s.Name, lk.warn.Bold(true)})
		case s.Name == "done" && st.done():
			l = append(l, seg{"✓ done", lk.ok.Bold(true)})
		default:
			l = append(l, seg{"○ " + s.Name, lk.dim})
		}
	}
	return l
}

func (mv mrsView) stageSpec(iw int, portrait bool) boxSpec {
	lk := mv.lk
	sp := boxSpec{title: "MR", sel: -1}
	r := mv.cur()
	if r == nil || mv.loading() {
		sp.rows = []line{{{ifs(mv.loading(), "syncing…") + ifs(!mv.loading(), "Nothing to show."), lk.sub}}}
		return sp
	}
	sp.title = r.Ref + " " + r.Title
	st := mv.stage(*r)
	if portrait {
		sp.right = line{{"tab", lk.key.Bold(true)}, {fmt.Sprintf(" MR list  %d your move", mv.yourCount()), lk.sub}}
	} else {
		meta := nonEmpty(ifs(r.Draft, "draft"), mv.ciSeg(*r).text, mv.mode2()+" review")
		sp.right = line{{strings.Join(meta, " · "), lk.sub}}
	}
	l := mv.rail(st)
	if st.Keys != "" && l.width()+4+len(st.Keys)+10 <= iw {
		l = lr(l, line{{"next ", lk.sub}, {st.Keys, lk.key.Bold(true)}}, iw)
	}
	sp.rows = []line{l}
	return sp
}

// ----- code pane -----

func (mv mrsView) fileIdx(d *mrDetail) int {
	if d == nil || len(d.Files) == 0 {
		return 0
	}
	return max(0, min(mv.file, len(d.Files)-1))
}

// codeRow is one display row of the code pane; hunk is true on hunk headers.
type codeRow struct {
	l    line
	hunk bool
	num  int // new (else old) line number, 0 on headers and boxes
}

// codeRows lays out the whole diff of the current file: numbered lines, the current card's
// line reversed with its marker, and the card's thread drawn under it. at is the row of that line.
func (mv mrsView) codeRows(d *mrDetail, iw int) (rows []codeRow, at int) {
	lk := mv.lk
	at = -1
	f := d.Files[mv.fileIdx(d)]
	c := mv.curCard()
	marks := map[int]string{}
	for _, o := range mv.cards() {
		if o.Path == f.Path && o.Line > 0 {
			marks[o.Line] = o.Label
		}
	}
	tw := iw - 2 // the scrollbar takes the last column
	for _, dl := range f.Lines {
		if dl.Kind == '@' {
			rows = append(rows, codeRow{l: line{{"┄┄ " + dl.Text, lk.dim}}, hunk: true})
			continue
		}
		num := func(n int) string {
			if n == 0 {
				return "    "
			}
			return fmt.Sprintf("%4d", n)
		}
		st, ks := lk.plain, lk.plain
		switch dl.Kind {
		case '+':
			st, ks = lk.ok, lk.ok.Bold(true)
		case '-':
			st, ks = lk.bad, lk.bad.Bold(true)
		}
		text := strings.ReplaceAll(dl.Text, "\t", "    ")
		lab, marked := marks[dl.New], dl.Kind != '-' && dl.New > 0
		hi := marked && c != nil && c.Path == f.Path && c.Line == dl.New
		l := line{{num(dl.Old) + " " + num(dl.New) + " ", lk.dim}}
		switch {
		case hi:
			rv := st.Bold(true).Reverse(true)
			body := string(dl.Kind) + " " + text
			tail := "  ◀ " + c.Label + ifs(c.Who != "", " "+c.Who) + " "
			room := max(4, tw-l.width()-2-runewidth.StringWidth(tail))
			body = runewidth.Truncate(body, room, "…")
			l = append(l, seg{"▶ ", lk.warn.Bold(true)}, seg{body + tail, rv})
			at = len(rows)
		case marked && lab != "":
			l = append(l, seg{"◆ ", lk.warn}, seg{string(dl.Kind) + " " + text, st})
		default:
			l = append(l, spc(2), seg{string(dl.Kind), ks}, seg{" " + text, st})
		}
		n := dl.New
		if n == 0 {
			n = dl.Old
		}
		rows = append(rows, codeRow{l: l, num: n})
		if hi && c.Text != "" {
			for _, b := range mv.threadBox(c, tw-13) {
				rows = append(rows, codeRow{l: append(line{spc(13)}, b...)})
			}
		}
	}
	return rows, at
}

// threadBox is the card's comment drawn under its line, in a thin box.
func (mv mrsView) threadBox(c *mrCard, w int) []line {
	lk := mv.lk
	w = min(84, w)
	if w < 20 {
		return nil
	}
	head := line{{"╭─ ", lk.warn}, {c.Who, lk.title}, {ifs(c.When != "", " "+c.When), lk.sub}, {" ", lk.warn}}
	kind := ifs(mv.answering() || mv.threadsMode(), "thread") + ifs(!mv.answering() && !mv.threadsMode(), "finding")
	head = append(head, seg{rep("─", w-head.width()-len(kind)-3), lk.warn}, seg{" " + kind + " ", lk.warn}, seg{"╮", lk.warn})
	out := []line{head}
	for _, t := range wrapLines(c.Text, w-4) {
		out = append(out, line{{"│ ", lk.warn}, {fit(t, w-4), lk.plain}, {" │", lk.warn}})
	}
	return append(out, line{{"╰" + rep("─", w-2) + "╯", lk.warn}})
}

// codeTop is the first row shown: the owner's scroll, else the card's line a third down.
func (mv mrsView) codeTop(n, at, slots int) int {
	top := mv.scroll
	if top < 0 {
		top = 0
		if at >= 0 {
			top = at - slots/3
		}
	}
	return max(0, min(top, n-slots))
}

// fileStrip lists the files that carry a card with their cards: decided ✓, open ○; the
// current file is marked ▶.
func (mv mrsView) fileStrip(d *mrDetail, iw int) line {
	lk := mv.lk
	cur := d.Files[mv.fileIdx(d)].Path
	var paths []string
	cards := map[string][]mrCard{}
	for _, c := range mv.cards() {
		if c.Path == "" {
			continue
		}
		if cards[c.Path] == nil {
			paths = append(paths, c.Path)
		}
		cards[c.Path] = append(cards[c.Path], c)
	}
	if cards[cur] == nil {
		paths = append([]string{cur}, paths...)
	}
	var l line
	for i, p := range paths {
		if i > 0 {
			l = append(l, spc(4))
		}
		ps := lk.sub
		if p == cur {
			l, ps = append(l, seg{"▶ ", lk.acc}), lk.bold
		}
		l = append(l, seg{p, ps})
		for _, c := range cards[p] {
			mk, ms := "○", lk.warn
			if c.decided != "pending" || c.State == "resolved" {
				mk, ms = "✓", lk.ok
			}
			l = append(l, seg{" " + mk + " " + c.Label, ms})
		}
	}
	return l.trunc(iw)
}

func (mv mrsView) codeSpec(iw, h int) boxSpec {
	lk := mv.lk
	sp := boxSpec{title: "Code", index: 2, sel: -1, focused: mv.pane == mpDiff}
	r := mv.cur()
	switch {
	case mv.loading():
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	case r == nil:
		sp.rows = []line{{{"Nothing to show.", lk.sub}}}
		return sp
	}
	d := mv.detail()
	if d == nil {
		if e := mv.fx.MRErrs[r.key()]; e != "" {
			sp.rows = []line{{{"◆ diff unavailable: " + e, lk.bad}}}
			return sp
		}
		sp.rows = []line{{{"loading the diff…", lk.sub}}}
		return sp
	}
	if len(d.Files) == 0 {
		sp.rows = []line{{{"No file changed.", lk.sub}}}
		return sp
	}
	f := d.Files[mv.fileIdx(d)]
	sp.title = fmt.Sprintf("%s +%d -%d", f.Path, f.Add, f.Del)
	sp.right = line{{"j k", lk.key.Bold(true)}, {" scroll  ", lk.sub}, {"J K", lk.key.Bold(true)}, {" hunk  ", lk.sub}, {"[ ]", lk.key.Bold(true)}, {fmt.Sprintf(" file %d of %d", mv.fileIdx(d)+1, len(d.Files)), lk.sub}}
	slots := max(1, h-lk.boxH(0))
	sp.rows = []line{mv.fileStrip(d, iw), {{rep("╌", iw), lk.dim}}}
	body := max(1, slots-len(sp.rows))
	rows, at := mv.codeRows(d, iw)
	top := mv.codeTop(len(rows), at, body)
	end := min(len(rows), top+body)
	// Scrollbar on the last column: the thumb is the share of the diff in sight.
	t0, t1 := 0, body
	if len(rows) > body {
		t0 = top * body / len(rows)
		t1 = max(t0+1, (end*body+len(rows)-1)/len(rows))
	}
	for k := 0; k < body; k++ {
		var l line
		if top+k < end {
			l = rows[top+k].l.trunc(iw - 2)
		}
		l = append(l, spc(iw-1-l.width()))
		if t0 <= k && k < t1 {
			l = append(l, seg{"█", lk.acc})
		} else {
			l = append(l, seg{"│", lk.dim})
		}
		sp.rows = append(sp.rows, l)
	}
	first, last := 0, 0
	for _, cr := range rows[top:end] {
		if cr.num > 0 {
			if first == 0 {
				first = cr.num
			}
			last = cr.num
		}
	}
	pos := []string{fmt.Sprintf("↑ %d above · ↓ %d below", top, len(rows)-end)}
	if first > 0 {
		pos = append(pos, fmt.Sprintf("lines %d–%d", first, last))
	}
	sp.bottom = strings.Join(pos, " · ")
	return sp
}

// hunkStarts lists the code rows that open a hunk.
func (mv mrsView) hunkStarts(iw int) []int {
	d := mv.detail()
	if d == nil || len(d.Files) == 0 {
		return nil
	}
	rows, _ := mv.codeRows(d, iw)
	var out []int
	for i, r := range rows {
		if r.hunk {
			out = append(out, i)
		}
	}
	return out
}

func (mv mrsView) ciSeg(r mrRow) seg {
	lk := mv.lk
	switch r.Pipeline {
	case "passed":
		return seg{"pipeline passed", lk.ok}
	case "failed":
		return seg{"pipeline failed", lk.bad}
	case "running":
		return seg{"pipeline running", lk.warn}
	}
	return seg{"", lk.sub}
}

func ifs(ok bool, s string) string {
	if ok {
		return s
	}
	return ""
}

// when formats a thread time as today, yesterday or a date; text that is not a time stays.
func (mv mrsView) when(at string) string {
	t, ok := parseTime(at)
	if !ok {
		return at
	}
	switch days := int(mv.now.Local().Sub(t.Local()).Hours() / 24); {
	case sameDay(t, mv.now):
		return "today " + clock(t)
	case days <= 1:
		return "yesterday " + clock(t)
	}
	return t.Local().Format("2 Jan 15:04")
}

func (mv mrsView) itemIdx(n int) int {
	if n == 0 {
		return 0
	}
	return max(0, min(mv.item, n-1))
}

func (mv mrsView) curFinding(d *mrDetail) *mrFinding {
	if d == nil || len(d.Findings) == 0 || mv.threadsMode() {
		return nil
	}
	return &d.Findings[mv.itemIdx(len(d.Findings))]
}

func (mv mrsView) curThread(d *mrDetail) *mrThread {
	if d == nil || len(d.Threads) == 0 {
		return nil
	}
	return &d.Threads[mv.itemIdx(len(d.Threads))]
}

// ----- Your move panel -----

func stateStyle(lk *look, state string) lipgloss.Style {
	switch state {
	case "accepted", "resolved", "edited", "posted", "replied", "fixed":
		return lk.ok
	case "answered", "pushed":
		return lk.cyan
	case "failed", "skipped":
		return lk.bad
	case "refused":
		return lk.sub
	case "chat":
		return lk.cyan
	}
	return lk.warn
}

// applyLine is the state of the apply run, nil when there is none.
func (mv mrsView) applyLine(ap *mrApply) line {
	lk := mv.lk
	if ap == nil {
		return nil
	}
	switch ap.State {
	case "starting":
		return line{{"● apply starting", lk.warn}, {" · asking cdx for a background run", lk.sub}}
	case "running":
		since := ""
		if t, err := time.Parse(time.RFC3339, ap.StartedAt); err == nil {
			since = " " + shortAge(mv.now.Sub(t))
		}
		return line{{"● apply running" + since, lk.warn}, {" · background run, commits locally, pushes nothing", lk.sub}}
	case "failed":
		return line{{"✗ apply failed: ", lk.bad}, {cmpOr(ap.Reason, "no reason given"), lk.sub}}
	}
	if ap.Commit == "" {
		return line{{"✓ applied", lk.ok}, {" · no code change", lk.sub}}
	}
	return line{{"✓ committed ", lk.ok}, {shortSHA(ap.Commit) + ifs(ap.Subject != "", " "+ap.Subject), lk.sub}}
}

// applyRows is the apply state under the steps; a commit says whether it is pushed.
func (mv mrsView) applyRows(ap *mrApply, iw int) []line {
	l := mv.applyLine(ap)
	if l == nil {
		return nil
	}
	if ap.State == "running" || ap.State == "starting" || ap.State == "failed" {
		return wrapLine(l, iw)
	}
	rows := []line{l.trunc(iw)}
	if ap.Commit != "" {
		rows = append(rows, line{spc(2), {ifs(ap.Pushed, "pushed to the MR branch") + ifs(!ap.Pushed, "local only, P pushes it"), mv.lk.sub}})
	}
	return rows
}

func shortSHA(s string) string {
	if len(s) > 8 {
		return s[:8]
	}
	return s
}

// keycap draws a key as a reversed cap.
func (mv mrsView) keycap(k string) seg {
	return seg{" " + k + " ", mv.lk.key.Bold(true).Reverse(true)}
}

// action is one big key row: cap, bold label, muted detail. on marks the decision taken.
func (mv mrsView) action(k, label, detail string, on bool, iw int) line {
	lk := mv.lk
	ls := lk.bold
	if on {
		ls = lk.acc
	}
	l := line{spc(1), mv.keycap(k), spc(2), {label, ls}}
	if detail != "" {
		l = append(l, spc(2), seg{detail, lk.sub})
	}
	return l.trunc(iw)
}

// actionRows is an action on one line, or in roomy mode the detail on a second line under
// the label.
func (mv mrsView) actionRows(k, label, detail string, on, two bool, iw int) []line {
	if !two {
		return []line{mv.action(k, label, detail, on, iw)}
	}
	return []line{mv.action(k, label, "", on, iw), line{spc(6), {detail, mv.lk.sub}}.trunc(iw)}
}

func (mv mrsView) steps(st mrStage, iw int) []line {
	lk := mv.lk
	var out []line
	for i, s := range st.Steps {
		text := cmpOr(s.Text, s.Name)
		var l line
		switch {
		case i < st.Cur || st.done():
			l = line{spc(2), {fmt.Sprintf("✓ %d  ", i+1), lk.ok}, {text, lk.sub}}
		case i == st.Cur:
			l = line{{"▶ ", lk.warn.Bold(true)}, {fmt.Sprintf("● %d  ", i+1), lk.warn.Bold(true)}, {text, lk.bold}}
			if s.Sub != "" {
				l = lr(l, line{{s.Sub, lk.warn.Bold(true)}}, iw)
			}
		default:
			l = line{spc(2), {fmt.Sprintf("○ %d  ", i+1), lk.dim}, {text, lk.dim}}
		}
		out = append(out, l)
	}
	return out
}

// diffText colours a fix or test diff.
func (mv mrsView) diffText(text string, iw int) []line {
	lk := mv.lk
	var rows []line
	for _, l := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		if strings.HasPrefix(l, "+++") || strings.HasPrefix(l, "---") {
			continue
		}
		st, mark, body := lk.plain, "  ", l
		switch {
		case strings.HasPrefix(l, "@@"):
			st = lk.dim
		case strings.HasPrefix(l, "+"):
			st, mark, body = lk.ok, "+ ", l[1:]
		case strings.HasPrefix(l, "-"):
			st, mark, body = lk.bad, "- ", l[1:]
		}
		rows = append(rows, line{spc(2), {mark, st.Bold(true)}, {fit(strings.TrimLeft(body, " \t"), iw-4), st}})
	}
	return rows
}

// cardBody is the current card: who and where, the comment, the fix, the test, the owner's
// note, the apply result, then the run's reasoning and pattern (cut first when room is short).
func (mv mrsView) cardBody(c *mrCard, n int, iw int) []line {
	lk := mv.lk
	where := c.Path + ifs(c.Line > 0, ":"+strconv.Itoa(c.Line))
	if c.Path == "" {
		where = "commit"
	}
	head := line{{c.Label + "  ", lk.warn.Bold(true)}, {where, lk.bold}, {ifs(c.Who != "", "   "+c.Who) + ifs(c.When != "", " · "+c.When), lk.sub}}
	rows := []line{lr(head, line{{c.State, stateStyle(lk, c.State)}}, iw)}
	for _, l := range wrapLines("« "+c.Text+" »", iw-2) {
		rows = append(rows, line{spc(2), {l, lk.plain}})
	}
	if c.Disagree {
		rows = append(rows, line{spc(2), {"disagree", lk.warn}, {" · the thread stays open for the reviewer", lk.sub}})
	}
	if c.Fix != "" {
		rows = append(rows, nil, line{{"Proposed fix", lk.head.Bold(true)}})
		rows = append(rows, mv.diffText(c.Fix, iw)...)
	}
	if c.Test != "" {
		rows = append(rows, nil, line{{"Test", lk.head.Bold(true)}})
		rows = append(rows, mv.diffText(c.Test, iw)...)
	}
	if mv.answering() && c.Note == "" {
		rows = append(rows, nil, line{{"Reply", lk.head.Bold(true)}, {ifs(c.Agree == "partly", "  agree in part"), lk.sub}})
		if c.Reply == "" {
			rows = append(rows, line{spc(2), {"no reply: the fix answers it", lk.sub}})
		}
		for _, l := range wrapLines(c.Reply, iw-2) {
			rows = append(rows, line{spc(2), {l, lk.plain}})
		}
	}
	if c.Note != "" && mv.typing == "" {
		lab := ifs(c.decided == "chat", "Chat note") + ifs(c.decided != "chat", ifs(mv.answering(), "Your reply")+ifs(!mv.answering(), "Your edit"))
		rows = append(rows, nil, line{{lab, lk.head.Bold(true)}})
		for _, l := range wrapLines(c.Note, iw-2) {
			rows = append(rows, line{spc(2), {l, lk.plain}})
		}
	}
	if c.Applied != "" && mv.typing == "" {
		st := lk.ok
		if c.Applied == "skipped" || c.Applied == "failed" {
			st = lk.bad
		}
		rows = append(rows, nil, line{{"apply: " + c.Applied, st.Bold(true)}, {ifs(c.AppStale, " · before your last decision"), lk.sub}})
		for _, l := range wrapLines(c.AppDetail, iw-2) {
			rows = append(rows, line{spc(2), {l, lk.sub}})
		}
	}
	if c.Why != "" && c.Why != c.Text {
		rows = append(rows, nil, line{{"Why", lk.head.Bold(true)}})
		for _, l := range wrapLines(c.Why, iw-2) {
			rows = append(rows, line{spc(2), {l, lk.sub}})
		}
	}
	if c.Pattern != "" {
		rows = append(rows, nil, line{{lk.cap("pattern"), lk.sub}, spc(1), {c.Pattern, lk.warn}, spc(1), {fmt.Sprintf("· seen %dx", c.Seen), lk.sub}})
	}
	return rows
}

// decideRows are the decision keys of the current card, or the input line while typing.
func (mv mrsView) decideRows(c *mrCard, gap bool, iw int) []line {
	lk := mv.lk
	if mv.typing != "" {
		label := ifs(mv.typing == "e", ifs(mv.answering(), "reply › ")+ifs(!mv.answering(), "edit › ")) + ifs(mv.typing == "c", "chat › ")
		var rows []line
		for i, l := range wrapLines(mv.input+"_", max(10, iw-len(label))) {
			rows = append(rows, line{{ifs(i == 0, label) + ifs(i > 0, rep(" ", len(label))), lk.key}, {l, lk.plain}})
		}
		return rows
	}
	// Once the fix left the machine, the card is settled: one line instead of the keycaps.
	switch c.State {
	case "pushed", "replied", "resolved", "posted":
		return []line{{{"Decided " + c.Label + " ", lk.head.Bold(true)}, {c.decided, stateStyle(lk, c.decided)}, {"   a x e c", lk.key}, {" change it", lk.sub}}}
	}
	mode := mv.mode2()
	quote := func(s string) string {
		return "« " + firstWords(s, 8) + ifs(len(strings.Fields(s)) > 8, "…") + " »"
	}
	var acc, ref, edit string
	switch {
	case mv.answering():
		acc = ifs(c.Reply != "", "reply "+quote(c.Reply)) + ifs(c.Reply == "", "fix only, no reply")
		ref, edit = "the thread stays open", "type your own reply"
	case mode == "peer":
		acc = ifs(c.Reply != "", "post "+quote(c.Reply))
		ref, edit = "nothing posted", "type your own comment"
	default:
		acc, ref, edit = "fix it in the apply run", "keep the code", "type the fix you want"
	}
	accLabel := ifs(c.Fix != "", "accept the fix") + ifs(c.Fix == "", "accept")
	rows := []line{{{"Decide " + c.Label, lk.head.Bold(true)}}}
	for _, a := range [][4]string{{"a", accLabel, acc, "accepted"}, {"x", "refuse", ref, "refused"}, {"e", ifs(mv.answering(), "reply") + ifs(!mv.answering(), "edit"), edit, "edited"}, {"c", "chat", "ask the agent about it", "chat"}} {
		rows = append(rows, mv.actionRows(a[0], a[1], a[2], c.decided == a[3], gap, iw)...)
	}
	return rows
}

// thenRows: what comes after this card. n next, A apply, P push or post.
func (mv mrsView) thenRows(d *mrDetail, st mrStage, gap bool, iw int) []line {
	lk := mv.lk
	r := mv.cur()
	mode := mv.mode2()
	cs := mv.cards()
	rows := []line{{{"Then", lk.head.Bold(true)}}}
	add := func(l []line) { rows = append(rows, l...) }
	if i := mv.itemIdx(len(cs)); i+1 < len(cs) {
		nx := cs[i+1]
		add(mv.actionRows("n", "next finding", nx.Label+" "+nx.Path+ifs(nx.Line > 0, ":"+strconv.Itoa(nx.Line))+ifs(nx.Who != "", " ("+nx.Who+")"), false, gap, iw))
	} else if len(cs) > 0 {
		add(mv.actionRows("n", "next MR", "last finding of "+r.Ref, false, gap, iw))
	}
	todo := 0
	for _, f := range d.Findings {
		if f.toApply(mode) {
			todo++
		}
	}
	applied := false
	for i, sp := range st.Steps {
		if sp.Name == "apply" {
			applied = st.Cur > i
		}
	}
	if mode != "peer" && !applied {
		det := "background run, local commit"
		switch {
		case d.Apply != nil && (d.Apply.State == "running" || d.Apply.State == "starting"):
			det = "running in the background"
		case todo == 0:
			det = "nothing to apply yet"
		}
		add(mv.actionRows("A", fmt.Sprintf("apply %d %s", todo, plural(todo, "decision", "decisions")), det, st.Keys == "A", gap, iw))
	}
	open := 0
	for _, f := range d.Findings {
		if !f.Resolved && f.Agree != "no" {
			open++
		}
	}
	switch mode {
	case "peer":
		add(mv.actionRows("P", "post the comments", "after your OK, then mark reviewed", st.Keys == "P", gap, iw))
	case "threads":
		add(mv.actionRows("P", "push + reply", fmt.Sprintf("resolves %d %s", open, plural(open, "thread", "threads")), st.Keys == "P", gap, iw))
	default:
		add(mv.actionRows("P", "push", "after your OK", st.Keys == "P", gap, iw))
	}
	return rows
}

// doneRows: the finished walkthrough.
func (mv mrsView) doneRows(d *mrDetail, st mrStage, iw int) []line {
	lk := mv.lk
	r := mv.cur()
	sent := "nothing pushed"
	if d.Apply != nil && d.Apply.Commit != "" {
		sent = shortSHA(d.Apply.Commit) + ifs(d.Apply.Pushed, " pushed") + ifs(!d.Apply.Pushed, " local")
	}
	if mv.mode2() == "peer" {
		sent = "comments posted"
	}
	n, res := len(d.Findings), 0
	for _, f := range d.Findings {
		if f.Resolved {
			res++
		}
	}
	rows := []line{
		{{"✓ done", lk.ok.Bold(true)}, spc(3), {sent, lk.plain}},
	}
	if mv.mode2() == "threads" {
		rows = append(rows, line{{fmt.Sprintf("%d/%d threads replied and resolved", res, n), lk.sub}, {ifs(r.Pipeline != "", " · "+mv.ciSeg(*r).text), lk.sub}})
	}
	last := st.Steps[len(st.Steps)-1]
	rows = append(rows, line{{"next: ", lk.sub}, {strings.TrimPrefix(last.Text, "done: "), lk.bold}})
	rows = append(rows, nil, mv.action("n", "next MR", "your move first", false, iw))
	return rows
}

func (mv mrsView) moveSpec(iw, h int, compact bool) boxSpec {
	lk := mv.lk
	sp := boxSpec{title: "Your move", index: 3, sel: -1, focused: mv.pane == mpCard}
	r, d := mv.cur(), mv.detail()
	switch {
	case mv.loading():
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	case r == nil:
		sp.rows = []line{{{"Nothing to show.", lk.sub}}}
		return sp
	}
	st := mv.stage(*r)
	sp.title = map[string]string{moveYou: "Your move", moveCrow: "Crow is on it", moveOthers: "Waiting on others", moveDone: "Done"}[st.Who] + " · " + r.Ref
	slots := max(1, h-lk.boxH(0))
	var steps, top []line
	if !compact {
		steps = append(mv.steps(st, iw), line{{rep("╌", iw), lk.dim}})
	}
	if d == nil {
		sp.rows = append(steps, line{{"loading…", lk.sub}})
		return sp
	}
	// Run state: a running or failed review or apply stays in sight under the steps.
	if mv.mode2() != "peer" && !mv.threadsMode() {
		if ap := mv.applyRows(d.Apply, iw); ap != nil {
			top = append(top, ap...)
			for _, w := range d.Apply.Warnings {
				for _, l := range wrapLines("! "+w, iw) {
					top = append(top, line{{l, lk.warn}})
				}
			}
		}
	}
	if d.Review != nil && d.Review.State == "failed" {
		top = append(top, line{{"✗ review failed: ", lk.bad}, {cmpOr(d.Review.Reason, "no reason given"), lk.sub}})
	}
	if st.done() && !mv.threadsMode() {
		sp.rows = append(append(append(steps, top...), nil), mv.doneRows(d, st, iw)...)
		return sp
	}
	cs := mv.cards()
	if len(cs) == 0 {
		var msg []line
		switch {
		case d.Review != nil && d.Review.State == "running":
			msg = []line{{{"Read only. Findings appear here when the run ends.", lk.sub}}}
		case mv.threadsMode():
			msg = []line{{{"No review thread on this MR.", lk.sub}}}
		default:
			msg = []line{{{"No findings saved for this MR.", lk.sub}}}
		}
		again := d.Review != nil && d.Review.State == "failed"
		msg = append(msg, nil, mv.action("r", "review", ifs(again, "runs the review again")+ifs(!again, "a read-only run prepares the findings"), st.Keys == "r", iw),
			mv.action("h", "hand off", "opens a repo session on this MR", false, iw))
		sp.rows = append(append(steps, top...), msg...)
		return sp
	}
	idx := mv.itemIdx(len(cs))
	c := cs[idx]
	sp.right = line{{fmt.Sprintf("%s %d of %d", ifs(mv.threadsMode(), "thread")+ifs(!mv.threadsMode(), "finding"), idx+1, len(cs)), lk.warn.Bold(true)}}
	body := mv.cardBody(&c, len(cs), iw)
	var pinned []line
	if mv.threadsMode() {
		pinned = []line{{{rep("╌", iw), lk.dim}}, {{"Next", lk.head.Bold(true)}},
			mv.action("r", "review the threads", "prepares a fix and a reply for each", false, iw),
			mv.action("h", "hand off", "opens a repo session on this MR", false, iw),
			mv.action("T", "findings", "shows the review findings", false, iw)}
	} else {
		build := func(gap bool) []line {
			p := []line{nil}
			p = append(p, mv.decideRows(&c, gap, iw)...)
			p = append(p, line{{rep("╌", iw), lk.dim}})
			return append(p, mv.thenRows(d, st, gap, iw)...)
		}
		pinned = build(!compact)
		if len(steps)+len(top)+len(body)+len(pinned)+1 > slots {
			pinned = build(false)
		}
	}
	foot := line{{"p", lk.key.Bold(true)}, {" previous finding   ", lk.sub}, {"M", lk.key.Bold(true)}, {" open MR", lk.sub}}
	if len(steps)+len(top)+4+len(pinned)+1 <= slots {
		top = append(steps, top...)
	} // else a short window: the steps give way to the card; the stage line keeps the stage.
	room := max(1, slots-len(top)-len(pinned)-1)
	if len(body) > room {
		hidden := len(body) - room + 1
		body = append(body[:room-1], line{{fmt.Sprintf("… %d more %s, the window is too short", hidden, plural(hidden, "line", "lines")), lk.sub}})
	}
	sp.rows = append(top, body...)
	sp.rows = append(sp.rows, pinned...)
	for len(sp.rows) < slots-1 {
		sp.rows = append(sp.rows, nil)
	}
	if len(sp.rows) < slots {
		sp.rows = append(sp.rows, foot)
	}
	return sp
}

func firstWords(s string, n int) string {
	f := strings.Fields(s)
	if len(f) > n {
		f = f[:n]
	}
	return strings.Join(f, " ")
}

// ----- composition -----

func (mv mrsView) body(w, h int) []string {
	lk := mv.lk
	lay := layoutFor(w)
	if lay == laySingle {
		return mv.single(w, h)
	}
	gapH, gapV := lk.gaps()
	stageH := lk.boxH(1)
	var out block
	if lay == layLandscape {
		w1 := min(44, w*19/100)
		w3 := min(66, w*28/100)
		w2 := w - w1 - w3 - 2*gapH
		mid := vcat(lk.box(mv.stageSpec(lk.innerW(w2), false), w2, stageH), blankLines(w2, gapV), mv.codeBox(w2, h-stageH-gapV))
		out = hcat([]int{w1, w2, w3}, gapH, mv.listBox(w1, h), mid, mv.moveBox(w3, h, false))
	} else {
		stage := lk.box(mv.stageSpec(lk.innerW(w), true), w, stageH)
		rest := h - stageH - 2*gapV
		moveH := max(lk.boxH(10), rest*45/100)
		codeH := rest - moveH
		var mid block
		if mv.pane == mpList {
			mid = mv.listBox(w, codeH)
		} else {
			mid = mv.codeBox(w, codeH)
		}
		out = vcat(stage, blankLines(w, gapV), mid, blankLines(w, gapV), mv.moveBox(w, moveH, true))
	}
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func (mv mrsView) listBox(w, h int) block { return mv.lk.box(mv.listSpec(mv.lk.innerW(w)), w, h) }
func (mv mrsView) codeBox(w, h int) block {
	return mv.lk.box(mv.codeSpec(mv.lk.innerW(w), h), w, h)
}
func (mv mrsView) moveBox(w, h int, compact bool) block {
	return mv.lk.box(mv.moveSpec(mv.lk.innerW(w), h, compact), w, h)
}

// codeSlots is the number of code rows the current layout shows (for scrolling by keys).
func (mv mrsView) codeSlots(w, h int) (iw, body int) {
	lk := mv.lk
	switch layoutFor(w) {
	case laySingle:
		return lk.innerW(w), h - lk.boxH(0) - 2
	case layLandscape:
		w2 := w - min(44, w*19/100) - min(66, w*28/100)
		return lk.innerW(w2), h - lk.boxH(1) - lk.boxH(0) - 2
	}
	rest := h - lk.boxH(1)
	return lk.innerW(w), rest - max(lk.boxH(10), rest*45/100) - lk.boxH(0) - 2
}

func (mv mrsView) single(w, h int) []string {
	lk := mv.lk
	pos := fmt.Sprintf("pane %d of %d · tab", mv.pane+1, mrPanes)
	var sp boxSpec
	switch mv.pane {
	case mpDiff:
		sp = mv.codeSpec(lk.innerW(w), h)
	case mpCard:
		sp = mv.moveSpec(lk.innerW(w), h, true)
	default:
		sp = mv.listSpec(lk.innerW(w))
	}
	sp.bottom = strings.Join(nonEmpty(sp.bottom, pos), " · ")
	b := lk.box(sp, w, h)
	for len(b) < h {
		b = append(b, blank(w))
	}
	return b[:h]
}

func nonEmpty(s ...string) []string {
	var out []string
	for _, x := range s {
		if x != "" {
			out = append(out, x)
		}
	}
	return out
}

// wrapLine wraps a head segment and its tail over the width: the head keeps its style on the
// first row, the tail style carries over.
func wrapLine(l line, iw int) []line {
	if len(l) < 2 {
		return []line{l}
	}
	head, tail := l[0], l[1]
	var rows []line
	for i, t := range wrapLines(head.text+tail.text, iw) {
		if i == 0 && strings.HasPrefix(t, head.text) {
			rows = append(rows, line{head, {t[len(head.text):], tail.style}})
		} else {
			rows = append(rows, line{{t, tail.style}})
		}
	}
	return rows
}
