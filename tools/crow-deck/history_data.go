package main

// History data: the sessions crow-usage knows, newest first. Costs are API-equivalent
// dollars, not what is billed. The sources record no end state, touched files, links or last
// messages, so live rows carry only what crow-usage reports; mock rows carry invented detail.

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

type histSession struct {
	ID, Project, Account, Model, Task, Scope string
	Start                                    time.Time
	Took                                     time.Duration
	Cost                                     float64
	Tokens, Messages                         int
	Active                                   bool
	Branch                                   string
	// Detail: the sources record these; false on live rows.
	Detail   bool
	EndState string
	Touched  []string
	Commits  int
	Linked   []string
	Last     []tailLine
}

func (s histSession) personal() bool { return s.Scope == "personal" }

// modelFamily: "claude-sonnet-5-5" is "sonnet"; mixed families read "mixed".
func modelFamily(models []string) string {
	fam := ""
	for _, m := range models {
		if strings.HasPrefix(m, "<") {
			continue // <synthetic> turns carry no model
		}
		f := m
		for _, n := range []string{"opus", "sonnet", "haiku", "fable"} {
			if strings.Contains(m, n) {
				f = n
			}
		}
		switch {
		case fam == "":
			fam = f
		case fam != f:
			return "mixed"
		}
	}
	return fam
}

// histOf maps crow-usage's sessions; nil until the first answer.
func histOf(u *usageSnap, reg *registry, now time.Time) []histSession {
	if u == nil {
		return nil
	}
	var out []histSession
	for _, s := range u.Sessions {
		start, err := time.Parse(time.RFC3339, s.Start)
		if err != nil {
			continue
		}
		end, err := time.Parse(time.RFC3339, s.End)
		if err != nil {
			end = start
		}
		took := end.Sub(start)
		buckets := 0
		for _, b := range s.Activity {
			buckets += len(b)
		}
		if active := time.Duration(buckets) * 10 * time.Minute; active > 0 && active < took {
			took = active
		}
		out = append(out, histSession{
			ID: s.ID, Project: reg.shortName(s.Project), Account: s.Account, Model: modelFamily(s.Models),
			Task: whoName(s, reg), Scope: reg.scopeOf(s.Project), Start: start.Local(), Took: took,
			Cost: s.CostUSD, Tokens: s.Tokens, Messages: s.Messages, Branch: s.Branch,
			Active: now.Sub(end) < 3*time.Minute,
		})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Start.After(out[j].Start) })
	return out
}

// histFilter is the search text and the chips of the History page.
type histFilter struct {
	query, project, account, model, status string
	days                                   int // 30, 7 or 1
}

func newHistFilter() histFilter { return histFilter{days: 30} }

var histDays = []int{30, 7, 1}

func (f histFilter) dayLabel() string {
	switch f.days {
	case 1:
		return "today"
	case 7:
		return "last 7 days"
	}
	return "last 30 days"
}

func (f histFilter) active() bool {
	return f.query != "" || f.project != "" || f.account != "" || f.model != "" || f.status != "" || f.days != 30
}

// visibleHist applies the mode, the date range, the chips and the search text.
func visibleHist(fx *fixtures, f histFilter, mode string, now time.Time) []histSession {
	cut := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1-max(f.days, 1))
	words := strings.Fields(strings.ToLower(f.query))
	var out []histSession
	for _, s := range fx.Hist {
		if !inMode(mode, s.Scope) || s.Start.Before(cut) {
			continue
		}
		if f.project != "" && s.Project != f.project || f.account != "" && s.Account != f.account || f.model != "" && s.Model != f.model {
			continue
		}
		if f.status == "active" && !s.Active || f.status == "done" && s.Active {
			continue
		}
		hay := strings.ToLower(s.Task + " " + s.Project + " " + s.Account + " " + s.Model)
		ok := true
		for _, w := range words {
			ok = ok && strings.Contains(hay, w)
		}
		if ok {
			out = append(out, s)
		}
	}
	return out
}

// histValues lists the distinct values of a field in the mode, most used first, for the chips.
func histValues(fx *fixtures, mode string, pick func(histSession) string) []string {
	n := map[string]int{}
	for _, s := range fx.Hist {
		if inMode(mode, s.Scope) && pick(s) != "" {
			n[pick(s)]++
		}
	}
	var vals []string
	for v := range n {
		vals = append(vals, v)
	}
	sort.Slice(vals, func(i, j int) bool {
		if n[vals[i]] != n[vals[j]] {
			return n[vals[i]] > n[vals[j]]
		}
		return vals[i] < vals[j]
	})
	return vals
}

// nextValue steps through "" (any), then each value, then back to "".
func nextValue(cur string, vals []string) string {
	if cur == "" {
		if len(vals) > 0 {
			return vals[0]
		}
		return ""
	}
	for i, v := range vals {
		if v == cur && i+1 < len(vals) {
			return vals[i+1]
		}
	}
	return ""
}

func (f *histFilter) cycle(key string, fx *fixtures, mode string) {
	switch key {
	case "p":
		f.project = nextValue(f.project, histValues(fx, mode, func(s histSession) string { return s.Project }))
	case "a":
		f.account = nextValue(f.account, histValues(fx, mode, func(s histSession) string { return s.Account }))
	case "m":
		f.model = nextValue(f.model, histValues(fx, mode, func(s histSession) string { return s.Model }))
	case "s":
		f.status = nextValue(f.status, []string{"done", "active"})
	case "d":
		for i, d := range histDays {
			if d == f.days {
				f.days = histDays[(i+1)%len(histDays)]
				return
			}
		}
		f.days = histDays[0]
	}
}

func histTotals(rows []histSession) (cost float64, took time.Duration) {
	for _, s := range rows {
		cost += s.Cost
		took += s.Took
	}
	return cost, took
}

func fmtTook(d time.Duration) string {
	m := int(d.Round(time.Minute) / time.Minute)
	switch {
	case m < 1:
		return "<1m"
	case m < 60:
		return strconv.Itoa(m) + "m"
	}
	return strconv.Itoa(m/60) + "h" + twoDigits(m%60)
}

func dayLabel(t, now time.Time) string {
	if t.Format("2006-01-02") == now.Format("2006-01-02") {
		return "Today"
	}
	return t.Format("Mon 2")
}

func mockHist(now time.Time) []histSession {
	at := func(back, h, m int) time.Time {
		d := now.AddDate(0, 0, -back)
		return time.Date(d.Year(), d.Month(), d.Day(), h, m, 0, 0, d.Location())
	}
	mk := func(t time.Time, project, account, model, task string, cost float64, took time.Duration, tokens int) histSession {
		scope := "work"
		if project == "lint" || project == "web" || project == "tools" {
			scope = "personal"
		}
		return histSession{ID: "mock-" + t.Format("0102-1504"), Project: project, Account: account, Model: model, Task: task, Scope: scope,
			Start: t, Took: took, Cost: cost, Tokens: tokens, Messages: 6, Branch: "main", Detail: true}
	}
	mn := time.Minute
	rows := []histSession{
		mk(at(0, 14, 4), "atlas", "main", "sonnet", "review !220", 6.2, 2*mn, 40_000),
		mk(at(0, 11, 0), "atlas", "main", "sonnet", "retry policy, first pass", 4.4, 90*mn, 31_000),
		mk(at(0, 9, 20), "crow", "main", "opus", "deck designs, Home landscape", 5.1, 190*mn, 64_000),
		mk(at(1, 16, 40), "atlas", "main", "opus", "typed ids, CI fix", 9.1, 130*mn, 88_000),
		mk(at(1, 10, 15), "atlas", "spare", "sonnet", "retry backoff design", 3.3, 48*mn, 27_000),
		mk(at(1, 8, 30), "vault", "main", "sonnet", "CR-0W data layout", 3.0, 65*mn, 41_000),
		mk(at(3, 20, 5), "lint", "main", "sonnet", "panel layout notes", 1.9, 40*mn, 18_000),
		mk(at(3, 15, 20), "atlas", "main", "sonnet", "webhook retry research", 2.8, 35*mn, 22_000),
		mk(at(3, 11, 10), "crow", "main", "haiku", "retro batch", 0.4, 12*mn, 9_000),
		mk(at(6, 9, 50), "atlas", "main", "haiku", "env var inventory", 0.4, 12*mn, 8_000),
		mk(at(6, 9, 5), "vault", "main", "sonnet", "terms reply draft", 2.1, 44*mn, 19_000),
	}
	rows[0].Active, rows[0].Detail = true, true
	rows[0].EndState = "findings in review, none posted"
	rows[0].Touched, rows[0].Linked = []string{"retry.py", "config.md"}, []string{"!220", "PRJ-571"}
	rows[0].Last = []tailLine{{"14:05", "reading retry.py (84 lines changed)"}, {"14:06", "checking config.md"}}
	rows[1].EndState = "comments drafted, none posted"
	rows[1].Touched, rows[1].Linked = []string{"retry.py", "config.md"}, []string{"!220", "PRJ-571"}
	rows[1].Last = []tailLine{{"11:52", "3 comments drafted, waiting for your OK"}, {"11:58", "› good, keep the first two"}, {"12:01", "saved to local/state/handoffs"}}
	rows[2].EndState = "HomeL merged board saved"
	rows[2].Touched, rows[2].Commits = []string{"HomeL_merge.dc.html"}, 0
	rows[2].Last = []tailLine{{"12:20", "Edit: HomeL_merge.dc.html"}, {"12:31", "boards compared, merge approved"}}
	rows[3].EndState = "CI green, MR updated"
	rows[3].Touched, rows[3].Commits, rows[3].Linked = []string{"ids.py", "tests/test_ids.py"}, 2, []string{"!214", "PRJ-565"}
	rows[3].Last = []tailLine{{"18:45", "pytest: 212 passed"}, {"18:50", "pushed 2 commits"}}
	for _, i := range []int{4, 5, 6, 7, 8, 9, 10} {
		rows[i].EndState = "notes saved, nothing pending"
		rows[i].Last = []tailLine{{rows[i].Start.Add(rows[i].Took).Format("15:04"), "saved to local/state/handoffs"}}
	}
	return rows
}
