package main

// MRs page: per project list of every MR or PR the owner authors, reviews or is assigned,
// the diff of the selected one in the centre, and one right box with the findings (review)
// or threads (own MR) list, the card of the selected item and the end-of-review footer.
// Landscape puts the three side by side; portrait stacks them; under 100 columns one pane
// shows at a time. Read only: the action keys explain what arrives with a later slice.

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	mpList = iota
	mpDiff
	mpCard
	mrPanes
)

var mrPaneNames = [mrPanes]string{"MRs and PRs", "Diff", "Card"}

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
	toggle bool // T: the other kind of card (threads on a review, findings on an own MR)
}

func (mv mrsView) loading() bool {
	return mv.sync == stLoading || !mv.fx.MRsReady && len(mv.fx.ForgeErrs) == 0 && mv.fx.Live
}

// inMode reports whether a scope shows in the current mode.
func inMode(mode, scope string) bool {
	return mode == "all" || (mode == "perso") == (scope == "personal")
}

func (mv mrsView) visible() []mrRow { return visibleMRs(mv.fx, mv.mode) }

func visibleMRs(fx *fixtures, mode string) []mrRow {
	var out []mrRow
	for _, r := range fx.MRs {
		if inMode(mode, r.Scope) {
			out = append(out, r)
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

// threadsMode: own MRs show their threads, reviews their findings; T flips it.
func (mv mrsView) threadsMode() bool {
	r := mv.cur()
	if r == nil {
		return false
	}
	return (r.Role == "AUTHOR") != mv.toggle
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

func (mv mrsView) statusParts(r mrRow) line {
	lk := mv.lk
	var parts []line
	if r.Draft {
		parts = append(parts, line{{"draft", lk.sub}})
	}
	switch r.Pipeline {
	case "running":
		parts = append(parts, line{{"CI running", lk.warn}})
	case "failed":
		parts = append(parts, line{{"CI failed", lk.bad}})
	case "passed":
		if len(parts) == 0 && r.ThreadsOpen == 0 && r.Role != "REVIEWER" {
			parts = append(parts, line{{"CI passed", lk.ok}})
		}
	}
	if r.ThreadsOpen > 0 {
		parts = append(parts, line{{fmt.Sprintf("%d %s", r.ThreadsOpen, plural(r.ThreadsOpen, "thread", "threads")), lk.warn}})
	}
	if r.Role == "REVIEWER" {
		age := "now"
		if t, ok := parseTime(r.Updated); ok {
			age = shortDays(mv.now.Sub(t))
		}
		parts = append(parts, line{{"awaiting me " + age, lk.plain}})
	}
	var out line
	for i, p := range parts {
		if i > 0 {
			out = append(out, seg{" · ", lk.sub})
		}
		out = append(out, p...)
	}
	return out
}

// listRows builds the grouped list. start[i] is the first display row of visible MR i.
func (mv mrsView) listRows() (rows []line, start []int) {
	lk, fx := mv.lk, mv.fx
	vis := mv.visible()
	start = make([]int, len(vis))
	var branches []mrBranch
	for _, b := range fx.Branches {
		if inMode(mv.mode, b.Scope) {
			branches = append(branches, b)
		}
	}
	seen := map[string]bool{}
	group := func(project string) {
		if seen[project] {
			return
		}
		seen[project] = true
		rows = append(rows, mv.groupHdr(project))
		for i, r := range vis {
			if r.Project != project {
				continue
			}
			start[i] = len(rows)
			rows = append(rows, line{{r.Ref, lk.bold}, spc(1), {r.Title, lk.plain}})
			rows = append(rows, append(append(mv.chip(r.Role), spc(1)), mv.statusParts(r)...))
		}
		for _, b := range branches {
			if b.Project == project {
				rows = append(rows, line{{"no MR", lk.warn}, spc(1), {b.Branch, lk.plain}},
					line{{fmt.Sprintf("%d %s ahead", b.Ahead, plural(b.Ahead, "commit", "commits")), lk.sub}, spc(1), {"M", lk.key}, spc(1), {"open MR", lk.sub}})
			}
		}
	}
	for _, r := range vis {
		group(r.Project)
	}
	for _, b := range branches {
		group(b.Project)
	}
	if len(rows) == 0 {
		rows = append(rows, line{{"Nothing open: no MR or PR is yours.", lk.sub}})
	}
	return rows, start
}

func (mv mrsView) listSpec() boxSpec {
	lk := mv.lk
	sp := boxSpec{title: "MRs and PRs", index: 1, sel: -1, selN: 2, focused: mv.pane == mpList}
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
		sp.title += " " + strconv.Itoa(len(vis))
	}
	var start []int
	sp.rows, start = mv.listRows()
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

// ----- diff -----

func (mv mrsView) fileIdx(d *mrDetail) int {
	if d == nil || len(d.Files) == 0 {
		return 0
	}
	return max(0, min(mv.file, len(d.Files)-1))
}

// markers lists the cards of the current file by new line number.
func (mv mrsView) markers(d *mrDetail, path string) map[int]int {
	out := map[int]int{}
	if d == nil {
		return out
	}
	if mv.threadsMode() {
		for i, t := range d.Threads {
			if t.Path == path && t.Line > 0 {
				out[t.Line] = i + 1
			}
		}
		return out
	}
	for i, f := range d.Findings {
		if f.Path == path && f.Line > 0 {
			out[f.Line] = i + 1
		}
	}
	return out
}

func (mv mrsView) diffRow(d diffLine, mark string, cur bool) line {
	lk := mv.lk
	if d.Kind == '@' {
		return line{{d.Text, lk.dim}}
	}
	n := d.New
	if n == 0 {
		n = d.Old
	}
	num := rep(" ", 4)
	if n > 0 {
		num = fmt.Sprintf("%3d ", n)
	}
	st, gut := lk.plain, lk.dim
	switch d.Kind {
	case '+':
		st, gut = lk.ok, lk.ok
	case '-':
		st, gut = lk.bad, lk.bad
	}
	ms := lk.warn
	if cur {
		ms = lk.warn.Bold(true)
	}
	l := line{{num, lk.dim}, {fit(mark, 3), ms}}
	if lk.th.frame == frameHUD || lk.th.frame == frameCyber {
		g := " "
		if d.Kind != ' ' {
			g = "▏"
		}
		l = append(l, seg{g, gut})
	}
	return append(l, seg{string(d.Kind) + " " + strings.ReplaceAll(d.Text, "\t", "    "), st})
}

// window returns n diff lines around the line of the current card.
func diffWindow(lines []diffLine, focus, n int) []diffLine {
	if len(lines) <= n {
		return lines
	}
	at := 0
	for i, l := range lines {
		if focus > 0 && l.New == focus && l.Kind != '-' {
			at = i
			break
		}
	}
	from := max(0, min(at-3, len(lines)-n))
	return lines[from : from+n]
}

func (mv mrsView) totals(d *mrDetail) (add, del int) {
	if d != nil {
		for _, f := range d.Files {
			add, del = add+f.Add, del+f.Del
		}
	}
	return
}

func (mv mrsView) ciSeg(r mrRow) seg {
	lk := mv.lk
	switch r.Pipeline {
	case "passed":
		return seg{"CI passed", lk.ok}
	case "failed":
		return seg{"CI failed", lk.bad}
	case "running":
		return seg{"CI running", lk.warn}
	}
	return seg{"", lk.sub}
}

func (mv mrsView) diffSpec(iw, h int) boxSpec {
	lk := mv.lk
	sp := boxSpec{title: "Diff", index: 2, sel: -1, focused: mv.pane == mpDiff}
	r := mv.cur()
	if mv.loading() {
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	}
	if r == nil {
		sp.rows = []line{{{"Nothing to show.", lk.sub}}}
		return sp
	}
	sp.title = r.Ref + " " + r.Title
	d := mv.detail()
	if mv.threadsMode() {
		sp.title += " · threads"
	}
	add, del := mv.totals(d)
	meta := r.Author
	if r.Branch != "" {
		meta += " · " + r.Branch + " → " + r.Target
	}
	if d != nil {
		meta += fmt.Sprintf(" · +%d -%d", add, del)
	}
	if mv.threadsMode() {
		meta = strings.TrimPrefix(strings.Join([]string{ifs(r.Draft, "draft"), fmt.Sprintf("%d open, %d resolved", r.ThreadsOpen, r.ThreadsResolved)}, " · "), " · ")
	}
	sp.right = line{{meta, lk.sub}}
	if d == nil {
		msg := "loading the diff…"
		if e := mv.fx.MRErrs[r.key()]; e != "" {
			msg = "◆ diff unavailable: " + e
			sp.rows = []line{{{msg, lk.bad}}}
			return sp
		}
		sp.rows = []line{{{msg, lk.sub}}}
		return sp
	}
	slots := max(1, h-lk.boxH(0))
	if mv.threadsMode() {
		sp.rows = mv.threadRows(d, r, iw, slots)
		return sp
	}
	sp.rows = mv.reviewRows(d, r, iw, slots)
	return sp
}

func ifs(ok bool, s string) string {
	if ok {
		return s
	}
	return ""
}

// reviewRows: the file bar, then the diff of the current file around the current finding.
func (mv mrsView) reviewRows(d *mrDetail, r *mrRow, iw, slots int) []line {
	lk := mv.lk
	if len(d.Files) == 0 {
		return []line{{{"No file changed.", lk.sub}}}
	}
	fi := mv.fileIdx(d)
	bar := line{{lk.cap("files"), lk.sub}, spc(1)}
	for i, f := range d.Files {
		st := lk.sub
		if i == fi {
			st = lk.underlineAcc
		}
		bar = append(bar, seg{f.Path, st})
		if n := len(mv.markers(d, f.Path)); n > 0 {
			bar = append(bar, seg{fmt.Sprintf(" ◆%d", n), lk.warn})
		}
		bar = append(bar, spc(2))
	}
	rows := []line{lr(bar, line{mv.ciSeg(*r)}, iw), nil}
	// Findings with no file sit in the commit row at the bottom.
	var commit line
	for i, f := range d.Findings {
		if f.Path == "" && !mv.threadsMode() {
			commit = append(commit, seg{fmt.Sprintf("◆%d %s", i+1, firstSentence(f.Body)), lk.warn})
		}
	}
	reserve := 0
	if commit != nil {
		reserve = 2
	}
	focus := 0
	mk := mv.markers(d, d.Files[fi].Path)
	cardLine := -1
	if f := mv.curFinding(d); f != nil && f.Path == d.Files[fi].Path {
		focus, cardLine = f.Line, f.Line
	}
	for _, dl := range diffWindow(d.Files[fi].Lines, focus, max(1, slots-len(rows)-reserve)) {
		m := ""
		if n, ok := mk[dl.New]; ok && dl.Kind != '-' {
			m = fmt.Sprintf("◆%d", n)
		}
		rows = append(rows, mv.diffRow(dl, m, dl.New == cardLine && dl.Kind != '-'))
	}
	if commit != nil {
		for len(rows) < slots-2 {
			rows = append(rows, nil)
		}
		rows = append(rows, mv.divider(iw), lr(line{{lk.cap("commit"), lk.sub}, spc(2), {r.Title, lk.plain}}, commit, iw))
	}
	return rows
}

func firstSentence(s string) string {
	s, _, _ = strings.Cut(s, ". ")
	return strings.TrimSuffix(s, ".")
}

// threadRows: the current thread, its saved fix and test or the code it points at, and the
// other threads in one line.
func (mv mrsView) threadRows(d *mrDetail, r *mrRow, iw, slots int) []line {
	lk := mv.lk
	t := mv.curThread(d)
	if t == nil {
		return []line{{{"No review thread on this MR.", lk.sub}}}
	}
	where := t.Path
	if t.Line > 0 {
		where += ":" + strconv.Itoa(t.Line)
	}
	rows := []line{{{t.Author, lk.bold}, spc(2), {mv.when(t.At) + " · " + where, lk.sub}}}
	for _, l := range wrapLines(t.Body, iw) {
		rows = append(rows, line{{l, lk.plain}})
	}
	sections := []*mrSection{t.Fix, t.Test}
	any := false
	for _, s := range sections {
		if s == nil {
			continue
		}
		any = true
		rows = append(rows, nil, lr(line{{s.Title, lk.head}}, line{{s.Meta, lk.sub}}, iw))
		for _, dl := range s.Lines {
			rows = append(rows, mv.diffRow(dl, "", false))
		}
	}
	if !any {
		rows = append(rows, nil, line{{lk.cap("code at " + where), lk.head}})
		for _, f := range d.Files {
			if f.Path != t.Path {
				continue
			}
			room := max(3, slots-len(rows)-3)
			for _, dl := range diffWindow(f.Lines, t.Line, room) {
				rows = append(rows, mv.diffRow(dl, "", dl.New == t.Line && dl.Kind != '-'))
			}
		}
		rows = append(rows, nil, line{{"No fix prepared yet. h hands this MR to a repo session.", lk.sub}})
	}
	var others []string
	for i, o := range d.Threads {
		if i != mv.itemIdx(len(d.Threads)) {
			others = append(others, fmt.Sprintf("%s %s", o.Path, ifs(o.Resolved, "resolved")+ifs(!o.Resolved, "open")))
		}
	}
	if len(others) > 0 && len(rows) < slots-2 {
		for len(rows) < slots-2 {
			rows = append(rows, nil)
		}
		rows = append(rows, mv.divider(iw), line{{lk.cap("other threads"), lk.sub}, spc(2), {strings.Join(others, " · "), lk.plain}})
	}
	return rows
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

// ----- right box -----

func stateStyle(lk *look, state string) lipgloss.Style {
	switch state {
	case "accepted", "resolved":
		return lk.ok
	case "refused":
		return lk.sub
	}
	return lk.warn
}

func (mv mrsView) cardSpec(iw, h int) boxSpec {
	lk := mv.lk
	sp := boxSpec{title: "Card", index: 3, sel: -1, focused: mv.pane == mpCard}
	r, d := mv.cur(), mv.detail()
	switch {
	case mv.loading():
		sp.rows = []line{{{"syncing…", lk.sub}}}
		return sp
	case r == nil:
		sp.rows = []line{{{"Nothing to show.", lk.sub}}}
		return sp
	case d == nil:
		sp.rows = []line{{{"loading…", lk.sub}}}
		return sp
	}
	threads := mv.threadsMode()
	n := mv.items()
	sp.title = ifs(threads, "Threads") + ifs(!threads, "Findings")
	if n > 0 {
		sp.title += " " + strconv.Itoa(n)
	}
	if hits := patternsHit(d, threads); hits > 0 {
		sp.right = line{{fmt.Sprintf("patterns hit: %d", hits), lk.sub}}
	}
	if n == 0 {
		sp.rows = []line{{{ifs(threads, "No review thread on this MR.") + ifs(!threads, "No findings saved for this MR."), lk.sub}}}
		if !threads {
			sp.rows = append(sp.rows, line{{"s self-review writes them; h hands the MR to a repo session.", lk.sub}})
		}
		return sp
	}
	idx := mv.itemIdx(n)
	var list []line
	for i := 0; i < n; i++ {
		var title, state string
		if threads {
			t := d.Threads[i]
			title = t.Path + ifs(t.Line > 0, ":"+strconv.Itoa(t.Line)) + " " + cmpOr(t.Pattern, firstWords(t.Body, 4))
			state = ifs(t.Resolved, "resolved") + ifs(!t.Resolved, "pending")
		} else {
			f := d.Findings[i]
			title, state = f.Title, f.State
		}
		g := seg{fmt.Sprintf("◆%d", i+1), stateStyle(lk, state)}
		list = append(list, lr(line{g, spc(1), {title, lk.bold}}, line{{state, stateStyle(lk, state)}}, iw))
	}
	card := mv.cardRows(d, iw)
	foot := mv.footRows(d, r, iw)
	slots := max(1, h-lk.boxH(0))
	if len(list)+1+len(card)+len(foot) > slots && len(foot) > 1 {
		foot = foot[len(foot)-1:]
	}
	sp.rows = append(sp.rows, list...)
	sp.sel = idx
	sp.rows = append(sp.rows, mv.divider(iw))
	room := max(0, slots-len(sp.rows)-len(foot)-1)
	if len(card) > room {
		card = card[:room]
	}
	sp.rows = append(sp.rows, card...)
	for len(sp.rows) < slots-len(foot)-1 && len(foot) > 0 {
		sp.rows = append(sp.rows, nil)
	}
	if len(foot) > 0 {
		sp.rows = append(sp.rows, mv.divider(iw))
		sp.rows = append(sp.rows, foot...)
	}
	sp.bottom = fmt.Sprintf("%d of %d", idx+1, n)
	return sp
}

func firstWords(s string, n int) string {
	f := strings.Fields(s)
	if len(f) > n {
		f = f[:n]
	}
	return strings.Join(f, " ")
}

func patternsHit(d *mrDetail, threads bool) int {
	seen := map[string]bool{}
	if threads {
		for _, t := range d.Threads {
			if t.Pattern != "" {
				seen[t.Pattern] = true
			}
		}
	} else {
		for _, f := range d.Findings {
			if f.Pattern != "" {
				seen[f.Pattern] = true
			}
		}
	}
	return len(seen)
}

func (mv mrsView) cardRows(d *mrDetail, iw int) []line {
	lk := mv.lk
	var rows []line
	pattern := func(name string, seen int) {
		if name == "" {
			return
		}
		rows = append(rows, nil, line{{lk.cap("pattern"), lk.sub}, spc(1), {name, lk.warn}, spc(1), {fmt.Sprintf("· seen %dx", seen), lk.sub}})
	}
	if mv.threadsMode() {
		t := mv.curThread(d)
		if t == nil {
			return nil
		}
		rows = append(rows, line{{fmt.Sprintf("Thread %d · %s", mv.itemIdx(len(d.Threads))+1, t.Path+ifs(t.Line > 0, ":"+strconv.Itoa(t.Line))), lk.head}})
		for _, l := range wrapLines(t.Author+": "+t.Body, iw) {
			rows = append(rows, line{{l, lk.sub}})
		}
		pattern(t.Pattern, t.Seen)
		if t.Reply != "" {
			rows = append(rows, nil, line{{lk.cap("reply"), lk.sub}, spc(1), {"0", lk.key.Bold(true)}, spc(1), {"none, fix only", lk.sub}, spc(2), {"1", lk.key.Bold(true)}, spc(1), {"short", lk.bold}})
			rows = append(rows, line{{"reply (FR), skipped with 0", lk.sub}})
			rows = append(rows, mv.replyBlock(t.Reply, iw)...)
			rows = append(rows, nil, mv.buttons(0, [2]string{"a", "accept"}, [2]string{"x", "refuse"}, [2]string{"e", "edit"}, [2]string{"c", "chat"}))
		}
		return rows
	}
	f := mv.curFinding(d)
	if f == nil {
		return nil
	}
	where := ifs(f.Path != "", f.Path+ifs(f.Line > 0, ":"+strconv.Itoa(f.Line))) + ifs(f.Path == "", "commit")
	rows = append(rows, line{{fmt.Sprintf("Finding %d · %s", mv.itemIdx(len(d.Findings))+1, where), lk.head}})
	for _, l := range wrapLines(f.Body, iw) {
		rows = append(rows, line{{l, lk.sub}})
	}
	pattern(f.Pattern, f.Seen)
	if f.Comment != "" {
		rows = append(rows, nil, line{{"comment (FR)", lk.sub}})
		rows = append(rows, mv.replyBlock(f.Comment, iw)...)
		rows = append(rows, nil, mv.buttons(0, [2]string{"a", "accept"}, [2]string{"x", "refuse"}, [2]string{"e", "edit"}, [2]string{"c", "chat"}))
	}
	return rows
}

// footRows is the end of the review: counts and the next step. Posting and pushing wait for
// the owner's OK and for a later slice.
func (mv mrsView) footRows(d *mrDetail, r *mrRow, iw int) []line {
	lk := mv.lk
	if mv.threadsMode() {
		fixes, resolved := 0, 0
		for _, t := range d.Threads {
			if t.Fix != nil {
				fixes++
			}
			if t.Resolved {
				resolved++
			}
		}
		if fixes == 0 {
			return []line{line{{"h", lk.key}, spc(1), {"hand this MR to a repo session", lk.sub}}}
		}
		return []line{
			lr(line{{"End of review · " + r.Ref, lk.bold}}, line{{"nothing pushed yet", lk.sub}}, iw),
			line{{fmt.Sprintf("%d fix + reply", fixes), lk.plain}, spc(2), {fmt.Sprintf("%d resolved", resolved), lk.ok}},
			line{{"next", lk.sub}, spc(2), {"push, reply, resolve the thread", lk.plain}},
			lr(nil, line{{"P", lk.key.Bold(true)}, spc(1), {"push, after your OK", lk.sub}}, iw),
		}
	}
	acc, pend, ref := 0, 0, 0
	for _, f := range d.Findings {
		switch f.State {
		case "accepted":
			acc++
		case "refused":
			ref++
		default:
			pend++
		}
	}
	return []line{
		lr(line{{"End of review · " + r.Ref, lk.bold}}, line{{"nothing posted yet", lk.sub}}, iw),
		line{{fmt.Sprintf("%d accepted", acc), lk.ok}, spc(2), {fmt.Sprintf("%d pending", pend), lk.warn}, spc(2), {fmt.Sprintf("%d refused", ref), lk.sub}},
		line{{"next", lk.sub}, spc(2), {fmt.Sprintf("post %d %s, mark reviewed", acc, plural(acc, "comment", "comments")), lk.plain}},
		lr(nil, line{{"P", lk.key.Bold(true)}, spc(1), {"post, after your OK", lk.sub}}, iw),
	}
}

// ----- composition -----

func (mv mrsView) body(w, h int) []string {
	lk := mv.lk
	lay := layoutFor(w)
	if lay == laySingle {
		return mv.single(w, h)
	}
	gapH, gapV := 0, 0
	if lk.th.frame == frameSoft {
		gapH, gapV = 3, 1
	}
	var out block
	if lay == layLandscape {
		w1 := min(36, w*22/100)
		w3 := min(48, w*29/100)
		w2 := w - w1 - w3 - 2*gapH
		out = hcat([]int{w1, w2, w3}, gapH, mv.listBox(w1, h), mv.diffBox(w2, h), mv.cardBox(w3, h))
	} else {
		list := mv.listSpec()
		listH := max(lk.boxH(4), min(lk.boxH(len(list.rows)), h*30/100))
		cardH := max(lk.boxH(8), h*36/100)
		diffH := max(lk.boxH(3), h-listH-cardH-2*gapV)
		out = vcat(mv.listBox(w, listH), blankLines(w, gapV), mv.diffBox(w, diffH), blankLines(w, gapV), mv.cardBox(w, h-listH-diffH-2*gapV))
	}
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out[:h]
}

func (mv mrsView) listBox(w, h int) block { return mv.lk.box(mv.listSpec(), w, h) }
func (mv mrsView) diffBox(w, h int) block {
	return mv.lk.box(mv.diffSpec(mv.lk.innerW(w), h), w, h)
}
func (mv mrsView) cardBox(w, h int) block {
	return mv.lk.box(mv.cardSpec(mv.lk.innerW(w), h), w, h)
}

func (mv mrsView) single(w, h int) []string {
	lk := mv.lk
	pos := fmt.Sprintf("pane %d of %d · tab", mv.pane+1, mrPanes)
	var sp boxSpec
	switch mv.pane {
	case mpDiff:
		sp = mv.diffSpec(lk.innerW(w), h)
	case mpCard:
		sp = mv.cardSpec(lk.innerW(w), h)
	default:
		sp = mv.listSpec()
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
