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
		return reviewMode(r)
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
		if p, l := f.at(); p == path && l > 0 {
			out[l] = i + 1
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
	// The current file and the card's line first; the bar used to list every file and got cut.
	cur := d.Files[fi].Path
	bar := line{{lk.cap("file"), lk.sub}, spc(1), {cur, lk.underlineAcc}}
	if f := mv.curFinding(d); f != nil {
		if p, l := f.at(); p == cur && l > 0 {
			bar = append(bar, seg{fmt.Sprintf(":%d", l), lk.acc})
		}
	}
	bar = append(bar, seg{fmt.Sprintf("  %d of %d", fi+1, len(d.Files)), lk.sub})
	cards := 0
	for _, f := range d.Files {
		cards += len(mv.markers(d, f.Path))
	}
	if n := len(mv.markers(d, cur)); n > 0 || cards > 0 {
		bar = append(bar, seg{fmt.Sprintf("  ◆%d of %d", n, cards), lk.warn})
	}
	rows := []line{lr(bar, line{mv.ciSeg(*r)}, iw), nil}
	// Findings with no file sit in the commit row at the bottom.
	var commit line
	for i, f := range d.Findings {
		if p, _ := f.at(); p == "" && !mv.threadsMode() {
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
	if f := mv.curFinding(d); f != nil {
		if p, l := f.at(); p == d.Files[fi].Path {
			focus, cardLine = l, l
		}
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
	case "accepted", "resolved", "edited", "posted", "replied", "fixed":
		return lk.ok
	case "answered":
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

// severityStyle colours a finding's severity: blocking stands out, nits fade.
func severityStyle(lk *look, sev string) lipgloss.Style {
	switch sev {
	case "blocking":
		return lk.bad.Bold(true)
	case "nit":
		return lk.dim
	case "question":
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
		return line{{"● apply running" + since, lk.warn}, {" · background run, commits locally, pushes nothing; the card refreshes every 20 s", lk.sub}}
	case "failed":
		return line{{"✗ apply failed: ", lk.bad}, {cmpOr(ap.Reason, "no reason given"), lk.sub}}
	}
	if ap.Commit == "" {
		return line{{"✓ applied", lk.ok}, {" · no code change", lk.sub}}
	}
	return line{{"✓ committed ", lk.ok}, {shortSHA(ap.Commit) + " " + ap.Subject + ifs(ap.Pushed, " · pushed") + ifs(!ap.Pushed, " · local only, P pushes it"), lk.sub}}
}

func shortSHA(s string) string {
	if len(s) > 8 {
		return s[:8]
	}
	return s
}

// reviewLine is the state of the review run, nil when there is none or it is done.
func (mv mrsView) reviewLine(rv *mrReview) line {
	lk := mv.lk
	if rv == nil {
		return nil
	}
	switch rv.State {
	case "running":
		since := ""
		if t, err := time.Parse(time.RFC3339, rv.StartedAt); err == nil {
			since = " · " + shortAge(mv.now.Sub(t))
		}
		return line{{"● review running", lk.warn}, {since + ifs(rv.Model != "", " · "+rv.Model), lk.sub}}
	case "failed":
		return line{{"✗ review failed: ", lk.bad}, {cmpOr(rv.Reason, "no reason given"), lk.sub}}
	}
	return nil
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
	sp.title = ifs(threads || mv.answering(), "Threads") + ifs(!threads && !mv.answering(), "Findings")
	if n > 0 {
		sp.title += " " + strconv.Itoa(n)
	}
	if hits := patternsHit(d, threads); hits > 0 {
		sp.right = line{{fmt.Sprintf("patterns hit: %d", hits), lk.sub}}
	}
	var run line
	if !threads {
		run = mv.reviewLine(d.Review)
	}
	if n == 0 {
		if run != nil {
			sp.rows = []line{run}
			if d.Review.State == "running" {
				sp.rows = append(sp.rows, line{{"Read only. Findings appear here when the run ends.", lk.sub}})
			} else {
				sp.rows = append(sp.rows, line{{"r", lk.key}, spc(1), {"runs the review again", lk.sub}})
			}
			return sp
		}
		sp.rows = []line{{{ifs(threads, "No review thread on this MR.") + ifs(!threads, "No findings saved for this MR."), lk.sub}}}
		if !threads {
			sp.rows = append(sp.rows, line{{"r", lk.key}, spc(1), {"starts a read-only review; h hands the MR to a repo session.", lk.sub}})
		}
		return sp
	}
	// Run state sits in the pinned bottom block, so it stays in sight however long the card is.
	var status []line
	if run != nil {
		status = append(status, run)
	}
	if ap := mv.applyLine(d.Apply); ap != nil && !threads {
		status = append(status, wrapLine(ap, iw)...)
		for _, w := range d.Apply.Warnings {
			for _, l := range wrapLines("! "+w, iw) {
				status = append(status, line{{l, lk.warn}})
			}
		}
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
			switch {
			case f.Resolved:
				state = "resolved"
			case f.Posted && mv.answering():
				state = "replied"
			case f.Posted:
				state = "posted"
			case f.Applied != "" && f.AppliedFor == [2]string{f.State, f.Note}:
				state = f.Applied
			}
		}
		g := seg{fmt.Sprintf("◆%d", i+1), stateStyle(lk, state)}
		list = append(list, lr(line{g, spc(1), {title, lk.bold}}, line{{state, stateStyle(lk, state)}}, iw))
	}
	card := mv.cardRows(d, iw)
	tail := mv.cardTail(d, iw) // the decision buttons or the input line: always shown
	foot := mv.footRows(d, r, iw)
	slots := max(1, h-lk.boxH(0))
	// Pinned bottom block: buttons, then run state, then the end of review (cut to its last line
	// when room is short). The card body above gets what is left and says what it hides.
	bottom := func() []line {
		var b []line
		b = append(b, tail...)
		if len(status)+len(foot) > 0 {
			b = append(b, mv.divider(iw))
		}
		return append(append(b, status...), foot...)
	}
	if len(list)+1+4+len(bottom()) > slots && len(foot) > 1 {
		foot = foot[len(foot)-1:]
	}
	pinned := bottom()
	off := len(sp.rows)
	sp.rows = append(sp.rows, list...)
	sp.sel = off + idx
	sp.rows = append(sp.rows, mv.divider(iw))
	room := max(1, slots-len(sp.rows)-len(pinned))
	if len(card) > room {
		hidden := len(card) - room + 1
		card = append(card[:room-1], line{{fmt.Sprintf("… %d more %s, the window is too short", hidden, plural(hidden, "line", "lines")), lk.sub}})
	}
	sp.rows = append(sp.rows, card...)
	for len(sp.rows) < slots-len(pinned) {
		sp.rows = append(sp.rows, nil)
	}
	sp.rows = append(sp.rows, pinned...)
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
	head := line{{fmt.Sprintf(ifs(mv.answering(), "Thread")+ifs(!mv.answering(), "Finding")+" %d · ", mv.itemIdx(len(d.Findings))+1), lk.head}}
	if f.Severity != "" {
		head = append(head, seg{f.Severity, severityStyle(lk, f.Severity)}, seg{" · ", lk.head})
	}
	rows = append(rows, append(head, seg{where, lk.head}))
	if mv.answering() {
		// What matters to decide comes first: comment, verdict, reply, fix, test. The run's
		// reasoning and the pattern go last, where a short window cuts them.
		rows = append(rows, mv.answerRows(f, iw)...)
		return append(rows, mv.answerTail(f, iw)...)
	}
	for _, l := range wrapLines(f.Body, iw) {
		rows = append(rows, line{{l, lk.sub}})
	}
	pattern(f.Pattern, f.Seen)
	if f.Comment != "" {
		rows = append(rows, nil, line{{ifs(mv.answering(), "reply (FR)") + ifs(!mv.answering(), "comment (FR)"), lk.sub}})
		rows = append(rows, mv.replyBlock(f.Comment, iw)...)
	}
	if f.Fix != "" {
		rows = append(rows, nil, line{{"proposed fix", lk.sub}})
		for _, l := range strings.Split(strings.TrimRight(f.Fix, "\n"), "\n") {
			st := lk.plain
			switch {
			case strings.HasPrefix(l, "+"):
				st = lk.ok
			case strings.HasPrefix(l, "-"):
				st = lk.bad
			}
			rows = append(rows, line{{fit(l, iw), st}})
		}
	}
	if f.Test != "" {
		rows = append(rows, nil, line{{"test", lk.sub}})
		for _, l := range strings.Split(strings.TrimRight(f.Test, "\n"), "\n") {
			st := lk.plain
			if strings.HasPrefix(l, "+") {
				st = lk.ok
			}
			rows = append(rows, line{{fit(l, iw), st}})
		}
	}
	if f.Note != "" && mv.typing == "" {
		rows = append(rows, nil, line{{ifs(f.State == "chat", "chat note"), lk.cyan}, {ifs(f.State != "chat", ifs(mv.answering(), "your reply")+ifs(!mv.answering(), "your edit")), lk.ok}})
		rows = append(rows, mv.replyBlock(f.Note, iw)...)
	}
	if f.Applied != "" && mv.typing == "" {
		label, st := "apply: "+f.Applied, lk.ok
		switch f.Applied {
		case "answered":
			label, st = "answer", lk.cyan
		case "skipped", "failed":
			st = lk.bad
		}
		stale := f.AppliedFor != [2]string{f.State, f.Note}
		rows = append(rows, nil, line{{label, st}, {ifs(stale, " · before your last decision"), lk.sub}})
		rows = append(rows, mv.replyBlock(f.AppliedDetail, iw)...)
	}
	return rows
}

// cardTail is the decision line of the current finding: the input while typing, else the buttons.
func (mv mrsView) cardTail(d *mrDetail, iw int) []line {
	lk := mv.lk
	f := mv.curFinding(d)
	if f == nil {
		return nil
	}
	if mv.typing != "" {
		label := ifs(mv.typing == "e", ifs(mv.answering(), "reply › ")+ifs(!mv.answering(), "edit › ")) + ifs(mv.typing == "c", "chat › ")
		var rows []line
		for i, l := range wrapLines(mv.input+"_", max(10, iw-len(label))) {
			rows = append(rows, line{{ifs(i == 0, label) + ifs(i > 0, rep(" ", len(label))), lk.key}, {l, lk.plain}})
		}
		return rows
	}
	active := map[string]int{"accepted": 1, "refused": 2, "edited": 3, "chat": 4}[f.State] - 1
	return []line{mv.buttons(active, [2]string{"a", "accept"}, [2]string{"x", "refuse"}, [2]string{"e", ifs(mv.answering(), "reply") + ifs(!mv.answering(), "edit")}, [2]string{"c", "chat"})}
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
	acc, pend, ref, chat := 0, 0, 0, 0
	for _, f := range d.Findings {
		switch f.State {
		case "accepted", "edited":
			acc++
		case "refused":
			ref++
		case "chat":
			chat++
		default:
			pend++
		}
	}
	mode, todo := detailMode(r, d), 0
	self := mode != "peer"
	for _, f := range d.Findings {
		if f.toApply(mode) {
			todo++
		}
	}
	unposted := 0
	for _, f := range d.Findings {
		if (f.State == "accepted" || f.State == "edited") && !f.Posted {
			unposted++
		}
	}
	next := fmt.Sprintf("post %d %s, mark reviewed", unposted, plural(unposted, "comment", "comments"))
	switch ap := d.Apply; {
	case !self:
	case todo == 0 && ap != nil && ap.Commit != "" && !ap.Pushed:
		next = "push " + shortSHA(ap.Commit)
	case todo == 0 && ap != nil && ap.Pushed && mode == "threads":
		next = "pushed; P replies and resolves"
	case todo == 0 && ap != nil && ap.Pushed:
		next = "pushed; ask a reviewer"
	case todo == 0 && mode == "threads":
		next = "P replies and resolves the threads"
	default:
		next = "apply the decisions, then push"
	}
	verb := map[string]string{"peer": "post", "self": "push", "threads": "push, reply, resolve"}[mode]
	end := lr(nil, line{{"P", lk.key.Bold(true)}, spc(1), {verb + ", after your OK", lk.sub}}, iw)
	if todo > 0 && (d.Apply == nil || d.Apply.State != "running") {
		end = lr(nil, line{{"A", lk.key.Bold(true)}, spc(1), {fmt.Sprintf("apply %d %s, commit locally", todo, plural(todo, "decision", "decisions")), lk.sub}}, iw)
	}
	sent := ifs(self, "nothing pushed yet") + ifs(!self, "nothing posted yet")
	if self && d.Apply != nil && d.Apply.Pushed {
		sent = "pushed"
	} else if posted := acc - unposted; !self && posted > 0 {
		sent = fmt.Sprintf("%d posted", posted)
	}
	return []line{
		lr(line{{"End of review · " + r.Ref, lk.bold}}, line{{sent, lk.sub}}, iw),
		line{{fmt.Sprintf("%d accepted", acc), lk.ok}, spc(2), {fmt.Sprintf("%d pending", pend), lk.warn}, spc(2), {fmt.Sprintf("%d refused", ref), lk.sub}, {ifs(chat > 0, fmt.Sprintf("  %d chat", chat)), lk.cyan}},
		line{{"next", lk.sub}, spc(2), {next, lk.plain}},
		end,
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
		if mv.answering() {
			// Answering threads: the card holds the comment, the fix and the reply, so it takes
			// almost half the width and the list shrinks.
			w1, w3 = min(28, w*15/100), w*45/100
		}
		w2 := w - w1 - w3 - 2*gapH
		out = hcat([]int{w1, w2, w3}, gapH, mv.listBox(w1, h), mv.diffBox(w2, h), mv.cardBox(w3, h))
	} else {
		list := mv.listSpec()
		listH := max(lk.boxH(4), min(lk.boxH(len(list.rows)), h*30/100))
		cardH := max(lk.boxH(8), h*36/100)
		if mv.answering() {
			// Same reason as landscape: the card gets the room, the diff keeps the thread's lines.
			listH = max(lk.boxH(4), min(lk.boxH(len(list.rows)), h*15/100))
			cardH = max(lk.boxH(8), h*55/100)
		}
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

// answerRows: the reviewer's comment and the run's verdict on it.
func (mv mrsView) answerRows(f *mrFinding, iw int) []line {
	lk := mv.lk
	var rows []line
	if f.ThreadBody != "" {
		for _, l := range wrapLines(cmpOr(f.ThreadAuthor, "reviewer")+": "+f.ThreadBody, iw) {
			rows = append(rows, line{{l, lk.plain}})
		}
	}
	verdict := map[string]string{"yes": "agree", "no": "disagree", "partly": "partly agree"}[f.Agree]
	st := lk.ok
	if f.Agree == "no" {
		st = lk.warn
	}
	if verdict != "" {
		rows = append(rows, nil, line{{verdict, st}, {ifs(f.Agree == "no", " · the thread stays open for the reviewer"), lk.sub}})
	}
	return rows
}

// answerTail: reply, fix, test, the owner's note, the apply result, then the reasoning and pattern.
func (mv mrsView) answerTail(f *mrFinding, iw int) []line {
	lk := mv.lk
	var rows []line
	reply := cmpOr(ifs(f.State == "edited", f.Note), f.Comment)
	if reply == "" {
		rows = append(rows, line{{"no reply: the fix answers it, resolved after the push", lk.sub}})
	} else {
		rows = append(rows, line{{ifs(f.State == "edited", "your reply (FR)") + ifs(f.State != "edited", "reply (FR)"), lk.sub}})
		rows = append(rows, mv.replyBlock(reply, iw)...)
	}
	diff := func(title, text string) {
		if text == "" {
			return
		}
		rows = append(rows, nil, line{{title, lk.sub}})
		for _, l := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
			if strings.HasPrefix(l, "+++") || strings.HasPrefix(l, "---") {
				continue
			}
			st := lk.plain
			switch {
			case strings.HasPrefix(l, "@@"):
				st = lk.dim
			case strings.HasPrefix(l, "+"):
				st = lk.ok
			case strings.HasPrefix(l, "-"):
				st = lk.bad
			}
			rows = append(rows, line{{fit(l, iw), st}})
		}
	}
	diff("proposed fix", f.Fix)
	diff("test", f.Test)
	if f.State == "chat" && f.Note != "" && mv.typing == "" {
		rows = append(rows, nil, line{{"chat note", lk.cyan}})
		rows = append(rows, mv.replyBlock(f.Note, iw)...)
	}
	if f.Applied != "" && mv.typing == "" {
		st := lk.ok
		if f.Applied == "skipped" || f.Applied == "failed" {
			st = lk.bad
		}
		stale := f.AppliedFor != [2]string{f.State, f.Note}
		rows = append(rows, nil, line{{"apply: " + f.Applied, st}, {ifs(stale, " · before your last decision"), lk.sub}})
		rows = append(rows, mv.replyBlock(f.AppliedDetail, iw)...)
	}
	rows = append(rows, nil, line{{"why", lk.sub}})
	for _, l := range wrapLines(f.Body, iw) {
		rows = append(rows, line{{l, lk.sub}})
	}
	if f.Pattern != "" {
		rows = append(rows, nil, line{{lk.cap("pattern"), lk.sub}, spc(1), {f.Pattern, lk.warn}, spc(1), {fmt.Sprintf("· seen %dx", f.Seen), lk.sub}})
	}
	return rows
}
