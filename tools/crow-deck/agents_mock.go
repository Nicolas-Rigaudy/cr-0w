package main

// Mock day lanes: crow-usage sessions for the mock day, joined to mockAgents by session id.

import (
	"strconv"
	"time"
)

type mockSpan struct{ from, to int } // buckets, inclusive

func mockSession(now time.Time, id, project, title, account, kind, model string, cost float64, spans ...mockSpan) usageSession {
	day := now.Local().Format("2006-01-02")
	s := usageSession{ID: id, Project: project, Title: title, Account: account, Kind: kind, Models: []string{"claude-" + model},
		Activity: map[string][]int{}, ActivityTokens: map[string]map[string]int{day: {}}, CachePct: 88, Compactions: 0,
		Branch: "main", FilesEdited: map[string]int{}}
	s.Today.CostUSD = cost
	first, last := -1, -1
	for _, sp := range spans {
		for b := sp.from; b <= sp.to; b++ {
			s.Activity[day] = append(s.Activity[day], b)
			tok := 4000 + (b*7919)%9000
			s.ActivityTokens[day][strconv.Itoa(b)] = tok
			s.Tokens += tok
			if first < 0 {
				first = b
			}
			last = b
		}
	}
	s.Start = bucketTime(now, first).Format(time.RFC3339)
	s.End = bucketTime(now, last).Add(8 * time.Minute).Format(time.RFC3339)
	return s
}

func mockSessions(now time.Time) []usageSession {
	ss := []usageSession{
		mockSession(now, "s-review", "atlas", "review !220", "main", "interactive", "opus-5-5", 1.84, mockSpan{78, 80}, mockSpan{82, 83}),
		mockSession(now, "s-undraft", "atlas", "undraft !213", "main", "interactive", "sonnet-5-5", 0.42, mockSpan{74, 76}, mockSpan{81, 83}),
		mockSession(now, "s-retention", "atlas", "retention policy tests", "main", "interactive", "opus-5-5", 2.10, mockSpan{60, 66}, mockSpan{68, 70}),
		mockSession(now, "s-designs", "crow", "deck designs, Home landscape", "main", "interactive", "fable-5-1", 3.95, mockSpan{70, 76}, mockSpan{79, 84}),
		mockSession(now, "s-home", "crow", "deck: real data on Home", "main", "interactive", "opus-5-5", 2.60, mockSpan{55, 62}),
		mockSession(now, "r-audit-1", "crow", "audit retros", "spare", "run", "sonnet-5-5", 0.31, mockSpan{57, 58}),
		mockSession(now, "r-audit-2", "crow", "audit retros", "spare", "run", "sonnet-5-5", 0.28, mockSpan{66, 67}),
		mockSession(now, "r-audit-3", "crow", "audit retros", "spare", "run", "sonnet-5-5", 0.30, mockSpan{75, 76}),
		mockSession(now, "s-vault", "vault", "CR-0W data layout", "main", "interactive", "opus-5-5", 1.12, mockSpan{77, 84}),
		mockSession(now, "r-99a8", "vault", "design doc PDF export", "spare", "run", "sonnet-5-5", 0.95, mockSpan{70, 79}),
		mockSession(now, "s-notes", "web", "0.40 deck pill", "main", "interactive", "opus-5-5", 0.80, mockSpan{72, 74}),
		mockSession(now, "s-tools", "tools", "AC2 acceptance tests", "main", "interactive", "sonnet-5-5", 0.35, mockSpan{63, 65}),
		mockSession(now, "s-lint", "lint", "panel layout notes", "main", "interactive", "haiku-5-5", 0.06, mockSpan{57, 57}),
	}
	ss[0].Branch, ss[0].Compactions, ss[0].CachePct = "fix/retry-backoff", 1, 91
	ss[0].FilesEdited = map[string]int{"retry.py": 3, "config.md": 1}
	ss[3].Branch, ss[3].FilesEdited = "chg/deck-agents", map[string]int{"agents.go": 9, "agents_data.go": 4, "main.go": 3}
	ss[8].FilesEdited = map[string]int{"layout.md": 2}
	return ss
}

// mockLaneSetup gives the mock agents their sessions and builds the lanes.
func mockLaneSetup(f *fixtures, now time.Time) {
	ids := []string{"s-review", "s-undraft", "s-designs", "s-vault", "r-99a8"}
	projects := []string{"atlas", "atlas", "crow", "vault", "vault"}
	titles := []string{"review !220", "undraft !213", "deck designs, Home landscape", "CR-0W data layout", "design doc PDF export"}
	ctx := []int{31_000, 22_000, 164_000, 41_000, 0}
	quiet := []int{2, 1, 0, 0, 9}
	for i := range ids {
		a := &f.Agents[i]
		a.Session, a.Project, a.Title, a.Ctx = ids[i], projects[i], titles[i], ctx[i]
		a.Quiet = now.Add(-time.Duration(quiet[i]) * time.Minute)
		a.Cwd = "/w/" + projects[i]
	}
	f.Lanes = lanesOf(now, mockSessions(now), f.Agents, f.Slots, nil)
	f.LaneGit = map[string]*laneGit{
		"/w/atlas": {At: now, Last: "4c1e9d2 fix[PY]: retry reads BACKOFF from config", Files: []gitFile{
			{Path: "src/retry.py", Add: 42, Del: 17}, {Path: "docs/config.md", Add: 6, Del: 1}},
			Diff: []string{"diff --git a/src/retry.py b/src/retry.py", "@@ -12,3 +12,4 @@ def backoff():",
				"-    base = float(os.environ.get(\"BACKOFF\", \"1.5\"))", "+    base = config.get(\"retry.backoff\", 1.5)", "+    cap = config.get(\"retry.cap\", 30)",
				"     return min(cap, base ** attempt)"}},
	}
}
