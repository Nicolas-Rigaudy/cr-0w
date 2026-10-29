package main

// Usage tab data: what crow-usage reports, shaped for the boxes. Costs are API-equivalent
// dollars from cdx's price table, not what is billed.

import (
	"sort"
	"time"
)

type costRow struct {
	Name string
	USD  float64
}

type daySpend struct {
	Label string // weekday
	Date  string
	USD   float64
}

// usageView is everything on the Usage tab except the account gauges (fixtures.Gauges).
type usageView struct {
	Updated   string
	Days      []daySpend // oldest first, up to 14
	RTKTokens int        // tokens rtk filtered today
	RTKPct    float64
	CachePct  int
	CacheUSD  float64 // saved by cache reads this week
	RouteUSD  float64 // saved by routing runs to a smaller model this week
	Sessions  []costRow
	Projects  []costRow
	Models    []costRow
}

func (u *usageView) savedTotal() float64 { return u.CacheUSD + u.RouteUSD }

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

func topCosts(rows []costRow, n int) []costRow {
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].USD > rows[j].USD })
	if len(rows) > n {
		rows = rows[:n]
	}
	var out []costRow
	for _, r := range rows {
		if r.USD > 0 {
			out = append(out, r)
		}
	}
	return out
}

// usageViewOf maps crow-usage's payload; nil until the first answer.
func usageViewOf(u *usageSnap, reg *registry, now time.Time) *usageView {
	if u == nil {
		return nil
	}
	short := map[string]string{}
	if reg != nil {
		for _, p := range reg.Projects {
			short[p.Name] = p.Short
		}
	}
	v := &usageView{Updated: clock(now), CachePct: u.Week.CachePct, CacheUSD: u.Week.CacheSavedUSD, RouteUSD: u.Week.RoutingSavedUSD,
		RTKTokens: u.RTK.TodaySavedTokens, RTKPct: u.RTK.TodaySavedPct}
	// History holds 30 days, Days the last 7 with detail; both carry cost.
	for _, d := range u.History {
		if t, err := time.ParseInLocation("2006-01-02", d.Date, time.Local); err == nil {
			v.Days = append(v.Days, daySpend{Label: t.Format("Mon"), Date: d.Date, USD: d.CostUSD})
		}
	}
	if n := len(v.Days); n > 14 {
		v.Days = v.Days[n-14:]
	}
	var sess []costRow
	for _, s := range u.Sessions {
		if s.Today.CostUSD > 0 {
			sess = append(sess, costRow{cmpOr(s.Title, s.Project), s.Today.CostUSD})
		}
	}
	v.Sessions = topCosts(sess, 5)
	var proj []costRow
	for _, p := range u.Today.ByProject {
		proj = append(proj, costRow{cmpOr(short[p.Name], p.Name), p.CostUSD})
	}
	v.Projects = topCosts(proj, 5)
	var mod []costRow
	for _, m := range u.Today.ByModel {
		mod = append(mod, costRow{m.Name, m.CostUSD})
	}
	v.Models = topCosts(mod, 5)
	return v
}

func mockUsage() *usageView {
	return &usageView{
		Updated:   "14:05",
		Days:      []daySpend{{"Fri", "", 18}, {"Sat", "", 4}, {"Sun", "", 6}, {"Mon", "", 31}, {"Tue", "", 38}, {"Wed", "", 44}, {"Thu", "", 27}},
		RTKTokens: 1_200_000, RTKPct: 74, CachePct: 91, CacheUSD: 11, RouteUSD: 14,
		Sessions: []costRow{{"review !220", 6.2}, {"crow designs", 5.1}, {"atlas MR", 4.4}, {"vault layout", 3.0}, {"retro batch", 1.1}},
		Projects: []costRow{{"atlas", 14.0}, {"crow", 8.1}, {"vault", 3.9}, {"tools", 1.0}},
		Models:   []costRow{{"opus", 16.5}, {"sonnet", 9.8}, {"haiku", 0.7}},
	}
}
