package main

// Whose move on an MR: the walkthrough steps (review, decide, apply, commit, push, resolve,
// done), the current one, who acts next and the key that moves it on. One rule feeds the
// stage line, the Your move panel, the list groups and the tab badge.

import (
	"fmt"
	"strconv"
	"time"
)

// Who acts next on an MR.
const (
	moveYou    = "you"    // the owner
	moveCrow   = "crow"   // a background run
	moveOthers = "others" // a reviewer or the author
	moveDone   = "done"   // nothing left on the owner's side
)

// mrStep is one step of the walkthrough: Name on the stage line, Text in the panel, Sub on
// the right of the current one.
type mrStep struct {
	Name, Text, Sub string
}

// mrStage is where an MR stands. Cur is the index of the current step; len(Steps) means done.
type mrStage struct {
	Steps   []mrStep
	Cur     int
	Who     string
	Reason  string // the list line: "2 decisions pending", "review asked of you"
	Keys    string // the key that moves it on, "" when none
	Problem bool   // a failed run or pipeline: the badge dot
}

func (s mrStage) done() bool { return s.Cur >= len(s.Steps) }

// mrGroup is the list group of a stage: 0 your move, 1 waiting on others, 2 done.
func (s mrStage) group() int {
	switch s.Who {
	case moveYou:
		return 0
	case moveDone:
		return 2
	}
	return 1
}

// stageSteps lists the steps of a review mode.
func stageSteps(mode string) []mrStep {
	switch mode {
	case "peer":
		return []mrStep{{Name: "review"}, {Name: "decide"}, {Name: "post"}, {Name: "done"}}
	case "self":
		return []mrStep{{Name: "review"}, {Name: "decide"}, {Name: "apply"}, {Name: "commit"}, {Name: "push"}, {Name: "done"}}
	}
	return []mrStep{{Name: "review"}, {Name: "decide"}, {Name: "apply"}, {Name: "commit"}, {Name: "push"}, {Name: "resolve"}, {Name: "done"}}
}

// rowMove is the rule when no review detail is known: a review asked of the owner, threads
// to answer or a red pipeline are the owner's move; a draft or an MR in review waits.
func rowMove(r mrRow) (who, reason string, problem bool) {
	switch {
	case r.Role == "REVIEWER":
		return moveYou, "review asked of you", false
	case r.Pipeline == "failed":
		return moveYou, "pipeline failed", true
	case r.ThreadsOpen > 0:
		return moveYou, strconv.Itoa(r.ThreadsOpen) + " " + plural(r.ThreadsOpen, "thread", "threads") + " to answer", false
	case r.Draft:
		return moveOthers, "draft", false
	}
	return moveOthers, "waiting for review", false
}

// mrStageOf computes the stage of r from its review detail (nil when not loaded).
func mrStageOf(r mrRow, d *mrDetail, now time.Time) mrStage {
	mode := detailMode(&r, d)
	st := mrStage{Steps: stageSteps(mode)}
	idx := func(name string) int {
		for i, s := range st.Steps {
			if s.Name == name {
				return i
			}
		}
		return len(st.Steps)
	}
	set := func(name, text, sub string) {
		i := idx(name)
		if i < len(st.Steps) {
			st.Steps[i].Text, st.Steps[i].Sub = text, sub
		}
	}
	// Default texts: what each step does, before it happens.
	threads := mode == "threads"
	set("review", "review", "")
	set("decide", "decide each finding", "")
	set("apply", "apply: background run", "")
	set("commit", "commit (local only)", "")
	set("push", ifs(threads, "push + reply to threads")+ifs(!threads, "push the fix"), "")
	set("post", "post the comments", "")
	set("resolve", "resolve the threads", "")
	set("done", "done", "")

	who, reason, problem := rowMove(r)
	st.Who, st.Reason, st.Problem = who, reason, problem
	if d == nil || d.Review == nil {
		set("review", ifs(threads, "review the threads")+ifs(!threads, "review the diff"), "")
		st.Keys = "r"
		return st
	}
	switch d.Review.State {
	case "running":
		since := ""
		if t, err := time.Parse(time.RFC3339, d.Review.StartedAt); err == nil {
			since = " " + shortAge(now.Sub(t))
		}
		set("review", "review running"+since+ifs(d.Review.Model != "", " · "+d.Review.Model), "")
		st.Who, st.Reason = moveCrow, "review running"
		return st
	case "failed":
		set("review", "review failed", "")
		st.Who, st.Reason, st.Keys, st.Problem = moveYou, "review failed", "r", true
		return st
	}
	n := len(d.Findings)
	set("review", fmt.Sprintf("review: %d %s", n, plural(n, "finding", "findings")), "")
	pend, unposted, open := 0, 0, 0
	for _, f := range d.Findings {
		if f.State == "pending" || f.State == "" {
			pend++
		}
		if (f.State == "accepted" || f.State == "edited") && !f.Posted {
			unposted++
		}
		if !f.Resolved && f.Agree != "no" {
			open++
		}
	}
	todo := 0
	for _, f := range d.Findings {
		if f.toApply(mode) {
			todo++
		}
	}
	set("resolve", fmt.Sprintf("resolve %d %s", open, plural(open, "thread", "threads")), "")
	reviewer := ""
	for _, f := range d.Findings {
		if f.ThreadAuthor != "" {
			reviewer = f.ThreadAuthor
			break
		}
	}
	if mode == "peer" {
		set("done", "done: "+cmpOr(r.Author, "the author")+" next", "")
	} else {
		set("done", "done: ask "+cmpOr(reviewer, "a reviewer"), "")
	}
	if pend > 0 {
		st.Cur, st.Who, st.Keys = idx("decide"), moveYou, "a x e c"
		set("decide", "decide each finding", fmt.Sprintf("%d of %d decided", n-pend, n))
		st.Reason = fmt.Sprintf("%d %s pending", pend, plural(pend, "decision", "decisions"))
		return st
	}
	set("decide", fmt.Sprintf("decide: %d of %d decided", n, n), "")
	if mode == "peer" {
		if unposted > 0 {
			st.Cur, st.Who, st.Keys = idx("post"), moveYou, "P"
			st.Reason = fmt.Sprintf("%d %s to post", unposted, plural(unposted, "comment", "comments"))
			set("post", st.Reason, "")
			return st
		}
		st.Cur, st.Who, st.Reason = len(st.Steps), moveDone, "done, reviewed"
		return st
	}
	ap := d.Apply
	switch {
	case ap != nil && (ap.State == "running" || ap.State == "starting"):
		since := ""
		if t, err := time.Parse(time.RFC3339, ap.StartedAt); err == nil && ap.State == "running" {
			since = " " + shortAge(now.Sub(t))
		}
		set("apply", "apply "+ap.State+since, "background run")
		st.Cur, st.Who, st.Reason = idx("apply"), moveCrow, "apply "+ap.State
		return st
	case todo > 0:
		set("apply", fmt.Sprintf("apply %d %s", todo, plural(todo, "decision", "decisions")), "")
		st.Cur, st.Who, st.Keys = idx("apply"), moveYou, "A"
		st.Reason = fmt.Sprintf("%d to apply", todo)
		if ap != nil && ap.State == "failed" {
			set("apply", "apply failed: "+cmpOr(ap.Reason, "no reason given"), "")
			st.Reason, st.Problem = "apply failed", true
		}
		return st
	}
	committed := ap != nil && ap.Commit != ""
	if committed {
		set("apply", "apply: done", "")
		set("commit", "commit "+shortSHA(ap.Commit), "")
	} else {
		set("apply", "apply: no code change", "")
		set("commit", "commit: nothing to commit", "")
	}
	if committed && !ap.Pushed {
		st.Cur, st.Who, st.Keys = idx("push"), moveYou, "P"
		st.Reason = "push " + shortSHA(ap.Commit)
		return st
	}
	if committed {
		set("push", "pushed "+shortSHA(ap.Commit), "")
	} else {
		set("push", "push: nothing to push", "")
	}
	if threads && open > 0 {
		st.Cur, st.Who, st.Keys = idx("resolve"), moveYou, "P"
		st.Reason = fmt.Sprintf("%d %s to resolve", open, plural(open, "thread", "threads"))
		return st
	}
	if threads {
		set("resolve", "threads resolved", "")
	}
	st.Cur, st.Who = len(st.Steps), moveDone
	st.Reason = "done"
	if committed {
		st.Reason += ", pushed " + shortSHA(ap.Commit)
	}
	if r.Pipeline == "failed" {
		// Nothing left in the walkthrough, but a red pipeline is still the owner's move.
		st.Who, st.Reason = moveYou, "pipeline failed"
	}
	return st
}

// mrStages returns the stage of every MR by key.
func mrStages(fx *fixtures, now time.Time) map[string]mrStage {
	out := map[string]mrStage{}
	for _, r := range fx.MRs {
		out[r.key()] = mrStageOf(r, fx.MRDetails[r.key()], now)
	}
	return out
}

// mrMoveCounts feeds the tab badge: MRs on the owner's move, and whether one has a real
// problem (failed pipeline or run).
func mrMoveCounts(fx *fixtures, now time.Time) (yours int, problem bool) {
	for _, r := range fx.MRs {
		st := mrStageOf(r, fx.MRDetails[r.key()], now)
		if st.Who == moveYou {
			yours++
		}
		problem = problem || st.Problem
	}
	return
}
