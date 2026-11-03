package main

// Usage tab data: what crow-usage reports, shaped for the boxes. Costs are API-equivalent
// dollars from cdx's price table, not what is billed.

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Periods of the d / w / m switch.
const (
	perToday = iota
	perWeek
	perMonth
)

var periodNames = []string{"today", "week", "month"}

type shareRow struct {
	Name     string
	USD      float64
	Tokens   int
	Sessions int
}

type toolRow struct {
	Name  string
	Calls int
}

// periodUse is one column of crow-usage's today / week / month totals.
type periodUse struct {
	Cost                float64
	Tokens              int
	CachePct            int
	CacheUSD, RouteUSD  float64
	Runs, Routed        int
	Compactions, Retros int
	RTKTokens, RTKCmds  int
	RTKPct              float64
	Projects, Models    []shareRow
	Tools               []toolRow
}

func (p periodUse) savedUSD() float64 { return p.CacheUSD + p.RouteUSD }

type daySpend struct {
	Label  string // "Tue 27"
	Date   string
	USD    float64
	Tokens int
}

// ledgerEntry is one session on one day, so a session that spans days shows each day's numbers.
type ledgerEntry struct {
	ID, Project, Title, Account, Model, Kind string
	Day                                      string // YYYY-MM-DD
	Start, End                               time.Time
	Cost                                     float64
	Tokens, Out, Messages                    int
	Week                                     float64 // % of the account's week, -1 unknown
	Ctx                                      int
	CtxLevel                                 string
	Compactions                              int
	Active                                   bool
}

// usageView is everything on the Usage tab except the account gauges (fixtures.Gauges).
type usageView struct {
	Updated  string
	Periods  [3]periodUse
	Days     []daySpend // oldest first, up to 14
	Entries  []ledgerEntry
	Activity map[string]map[int64]int // account -> 10-min bucket start (unix) -> tokens
	Since    string                   // first day the ledger covers
}

func shareRows(in []usageName, name func(string) string) []shareRow {
	var out []shareRow
	for _, r := range in {
		if r.Name == "?" && r.CostUSD == 0 {
			continue
		}
		out = append(out, shareRow{name(r.Name), r.CostUSD, r.Tokens, r.Sessions})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].USD > out[j].USD })
	return out
}

func periodOf(t usageTotals, name func(string) string) periodUse {
	p := periodUse{Cost: t.CostUSD, Tokens: t.Tokens, CachePct: t.CachePct, CacheUSD: t.CacheSavedUSD, RouteUSD: t.RoutingSavedUSD,
		Runs: t.Runs, Routed: t.RoutedRuns, Compactions: t.Compactions,
		Projects: shareRows(t.ByProject, name), Models: shareRows(t.ByModel, func(s string) string { return s })}
	for _, x := range t.Tools {
		p.Tools = append(p.Tools, toolRow{x.Name, x.Calls})
	}
	return p
}

// isoDay turns 2026-10-27 into "Tue 27".
func isoDay(date string) string {
	t, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return date
	}
	return t.Format("Mon 2")
}

// usageViewOf maps crow-usage's payload; nil until the first answer.
func usageViewOf(u *usageSnap, reg *registry, now time.Time) *usageView {
	if u == nil {
		return nil
	}
	name := func(s string) string { return reg.shortName(s) }
	v := &usageView{Updated: clock(now), Activity: map[string]map[int64]int{}}
	for i, t := range []usageTotals{u.Today, u.Week, u.Month} {
		v.Periods[i] = periodOf(t, name)
	}
	v.Periods[perToday].RTKTokens, v.Periods[perToday].RTKPct, v.Periods[perToday].RTKCmds = u.RTK.TodaySavedTokens, u.RTK.TodaySavedPct, u.RTK.TodayCommands
	for _, i := range []int{perWeek, perMonth} {
		v.Periods[i].RTKTokens, v.Periods[i].RTKPct, v.Periods[i].RTKCmds = u.RTK.SavedTokens, u.RTK.SavedPct, u.RTK.Commands
	}
	v.Periods[perToday].Retros, v.Periods[perWeek].Retros = u.Retros.Today, u.Retros.Week
	for _, d := range u.History {
		v.Days = append(v.Days, daySpend{Label: isoDay(d.Date), Date: d.Date, USD: d.CostUSD, Tokens: d.Tokens})
	}
	if n := len(v.Days); n > 14 {
		v.Days = v.Days[n-14:]
	}
	for _, s := range u.Sessions {
		start, _ := time.Parse(time.RFC3339, s.Start)
		end, err := time.Parse(time.RFC3339, s.End)
		if err != nil {
			end = start
		}
		week := -1.0
		if s.LimitShare.Week != nil {
			week = *s.LimitShare.Week
		}
		for day, d := range s.ByDay {
			if d.Tokens == 0 && d.CostUSD == 0 {
				continue
			}
			v.Entries = append(v.Entries, ledgerEntry{
				ID: s.ID, Project: reg.shortName(s.Project), Title: cmpOr(strings.TrimSpace(s.Title), reg.shortName(s.Project)),
				Account: s.Account, Model: modelFamily(s.Models), Kind: s.Kind, Day: day, Start: start.Local(), End: end.Local(),
				Cost: d.CostUSD, Tokens: d.Tokens, Out: s.OutputTokens, Messages: s.Messages, Week: week,
				Ctx: s.ContextTokens, CtxLevel: cmpOr(s.ContextLevel, "ok"), Compactions: s.Compactions, Active: s.Active,
			})
			if v.Since == "" || day < v.Since {
				v.Since = day
			}
		}
		for day, buckets := range s.ActivityTokens {
			t0, err := time.ParseInLocation("2006-01-02", day, time.Local)
			if err != nil {
				continue
			}
			for b, tok := range buckets {
				n, err := strconv.Atoi(b)
				if err != nil {
					continue
				}
				if v.Activity[s.Account] == nil {
					v.Activity[s.Account] = map[int64]int{}
				}
				v.Activity[s.Account][t0.Add(time.Duration(n)*10*time.Minute).Unix()] += tok
			}
		}
	}
	sort.SliceStable(v.Entries, func(i, j int) bool { return v.Entries[i].Cost > v.Entries[j].Cost })
	return v
}

// entries keeps the ledger rows of a period: today, the last 7 days, or all the sessions
// crow-usage lists (it keeps 7 days of sessions, so month shows the same rows as week).
func (u *usageView) entries(period int, now time.Time) []ledgerEntry {
	day := now.Local().Format("2006-01-02")
	from := now.Local().AddDate(0, 0, -6).Format("2006-01-02")
	var out []ledgerEntry
	for _, e := range u.Entries {
		switch {
		case period == perToday && e.Day != day, period == perWeek && e.Day < from:
			continue
		}
		out = append(out, e)
	}
	return out
}

func mockUsage(now time.Time) *usageView {
	day := func(back int) string { return now.AddDate(0, 0, -back).Format("2006-01-02") }
	at := func(back, h, m int) time.Time {
		d := now.AddDate(0, 0, -back)
		return time.Date(d.Year(), d.Month(), d.Day(), h, m, 0, 0, now.Location())
	}
	v := &usageView{Updated: "14:05", Activity: map[string]map[int64]int{}}
	costs := []float64{104, 116, 62, 20, 0, 0, 9.9, 84, 0, 63, 0, 0, 0, 39}
	for i, c := range costs {
		d := day(len(costs) - 1 - i)
		v.Days = append(v.Days, daySpend{Label: isoDay(d), Date: d, USD: c, Tokens: int(c * 2_100_000)})
	}
	proj := []shareRow{{"crow", 118, 310_500_000, 61}, {"atlas", 49.8, 129_200_000, 5}, {"vault", 6.6, 14_800_000, 3}, {"shell", 2.7, 6_800_000, 2},
		{"lint", 2.4, 4_900_000, 1}, {"tools", 2.2, 7_000_000, 2}, {"web", 1.8, 6_000_000, 1}, {"notes", 1.0, 3_100_000, 1}}
	mod := []shareRow{{"opus-5-5", 137, 332_700_000, 30}, {"sonnet-5-5", 49.7, 153_400_000, 51}, {"haiku-5-5", 0.01, 27_000, 1}}
	tools := []toolRow{{"Bash", 3589}, {"Read", 277}, {"Edit", 169}, {"Grep", 10}}
	scale := func(rows []shareRow, f float64) []shareRow {
		out := make([]shareRow, len(rows))
		for i, r := range rows {
			out[i] = shareRow{r.Name, r.USD * f, int(float64(r.Tokens) * f), max(1, int(float64(r.Sessions)*f))}
		}
		return out
	}
	v.Periods[perWeek] = periodUse{Cost: 187, Tokens: 486_000_000, CachePct: 97, CacheUSD: 1512, RouteUSD: 50.7, Runs: 57, Routed: 43,
		Compactions: 40, Retros: 28, RTKTokens: 2_700_000, RTKPct: 22, RTKCmds: 10353, Projects: proj, Models: mod, Tools: tools}
	v.Periods[perToday] = periodUse{Cost: 39, Tokens: 81_000_000, CachePct: 96, CacheUSD: 267, RouteUSD: 6.6, Runs: 10, Routed: 9,
		Compactions: 14, Retros: 9, RTKTokens: 513_000, RTKPct: 43, RTKCmds: 1728, Projects: scale(proj[:4], .2), Models: scale(mod[:2], .2), Tools: tools[:3]}
	v.Periods[perMonth] = periodUse{Cost: 655, Tokens: 1_770_000_000, CachePct: 96, CacheUSD: 5210, RouteUSD: 160, Runs: 210, Routed: 150,
		Compactions: 131, RTKTokens: 2_700_000, RTKPct: 22, RTKCmds: 10353, Projects: scale(proj, 3.5), Models: scale(mod, 3.5), Tools: tools}
	type ses struct {
		proj, title, acct, model, kind, lvl string
		back, h0, m0, h1, m1                int
		cost                                float64
		tok, ctx, comp                      int
		week                                float64
		active                              bool
	}
	for i, s := range []ses{
		{"crow", "Deck review page", "main", "sonnet", "interactive", "warn", 0, 11, 2, 13, 40, 12.9, 42_600_000, 145_900, 8, 6, false},
		{"crow", "Review threads disposition column fix", "main", "opus", "interactive", "warn", 0, 9, 29, 10, 27, 8.2, 16_800_000, 122_200, 2, 3, false},
		{"crow", "Run item_045 findings", "main", "opus", "run", "warn", 1, 12, 23, 13, 8, 6.0, 12_300_000, 138_700, 1, 2, false},
		{"crow", "Deck critique and redesign planning", "main", "opus", "interactive", "ok", 1, 17, 2, 17, 44, 5.6, 7_300_000, 109_400, 2, 2, false},
		{"crow", "Handoff audit follow-through", "main", "opus+sonnet", "interactive", "ok", 2, 16, 59, 17, 17, 4.8, 9_400_000, 107_700, 1, 2, false},
		{"crow", "deck-tabs-2", "spare", "sonnet", "run", "warn", 2, 21, 11, 21, 36, 4.7, 15_800_000, 144_300, 2, 2, false},
		{"crow", "Driver terms reread, comment", "main", "opus", "interactive", "warn", 3, 10, 45, 12, 23, 4.6, 8_700_000, 124_900, 1, 2, false},
		{"atlas", "ADR 003 typed IDs", "main", "opus", "interactive", "warn", 2, 10, 58, 18, 2, 37.4, 99_700_000, 139_900, 2, 16, false},
		{"atlas", "Review atlas MRs 213 214", "main", "opus", "interactive", "over", 3, 11, 14, 13, 24, 7.5, 21_700_000, 270_900, 1, 3, false},
		{"atlas", "Add the demo car to the seed data", "main", "opus", "interactive", "ok", 0, 13, 41, 14, 6, 1.8, 3_800_000, 117_700, 0, 0.4, true},
		{"vault", "Vault layout pass", "spare", "sonnet", "interactive", "ok", 1, 15, 0, 16, 10, 4.0, 5_700_000, 90_000, 0, 1, false},
		{"shell", "Bluetooth error surfacing", "main", "sonnet", "run", "ok", 4, 9, 0, 9, 40, 2.7, 6_800_000, 60_000, 0, 1, false},
		{"lint", "Release notes", "main", "haiku", "run", "ok", 5, 8, 0, 8, 20, 0.6, 1_200_000, 20_000, 0, 0.2, false},
	} {
		st, en := at(s.back, s.h0, s.m0), at(s.back, s.h1, s.m1)
		v.Entries = append(v.Entries, ledgerEntry{ID: fmt.Sprintf("m%02d", i), Project: s.proj, Title: s.title, Account: s.acct, Model: s.model,
			Kind: s.kind, Day: day(s.back), Start: st, End: en, Cost: s.cost, Tokens: s.tok, Out: s.tok / 100, Messages: 3 + i%5,
			Week: s.week, Ctx: s.ctx, CtxLevel: s.lvl, Compactions: s.comp, Active: s.active})
		if s.back <= 1 {
			for t := st; t.Before(en); t = t.Add(10 * time.Minute) {
				if v.Activity[s.acct] == nil {
					v.Activity[s.acct] = map[int64]int{}
				}
				v.Activity[s.acct][t.Truncate(10*time.Minute).Unix()] += s.tok / max(1, int(en.Sub(st)/(10*time.Minute)))
			}
		}
	}
	sort.SliceStable(v.Entries, func(i, j int) bool { return v.Entries[i].Cost > v.Entries[j].Cost })
	v.Since = day(5)
	return v
}

// emptyClock turns "2026-10-29T22:45+01:00" into "22:45" local; empty when absent.
func emptyClock(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04Z07:00", s)
	}
	if err != nil {
		return ""
	}
	return t.Local().Format("15:04")
}
