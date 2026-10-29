package main

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func testMRs() *mrsSnap {
	return &mrsSnap{
		Forges: map[string]string{"gitlab": "ok", "github": "ok"},
		MRs: []mrRow{
			{Project: "web", Name: "web", Scope: "personal", Forge: "github", Repo: "me/web", Number: 12, Ref: "#12", Title: "cache", Author: "me", Role: "AUTHOR", Pipeline: "passed", Updated: "2026-10-29T09:00:00Z"},
			{Project: "atlas", Name: "atlas", Scope: "work", Forge: "gitlab", Repo: "grp/atlas", Number: 5, Ref: "!5", Title: "retry", Author: "sam", Role: "REVIEWER", Updated: "2026-10-27T09:00:00Z"},
			{Project: "atlas", Name: "atlas", Scope: "work", Forge: "gitlab", Repo: "grp/atlas", Number: 6, Ref: "!6", Title: "seed", Author: "me", Role: "AUTHOR", ThreadsOpen: 2, Pipeline: "failed", Updated: "2026-10-28T09:00:00Z"},
		},
	}
}

const testDiff = "@@ -10,3 +10,4 @@ def f():\n a\n-b\n+c\n+d\n e\n"

func TestParseDiffNumbersLines(t *testing.T) {
	got := parseDiff(testDiff)
	want := []diffLine{
		{Kind: '@', Text: "@@ -10,3 +10,4 @@ def f():"},
		{Old: 10, New: 10, Kind: ' ', Text: "a"},
		{Old: 11, Kind: '-', Text: "b"},
		{New: 11, Kind: '+', Text: "c"},
		{New: 12, Kind: '+', Text: "d"},
		{Old: 12, New: 13, Kind: ' ', Text: "e"},
	}
	if len(got) != len(want) {
		t.Fatalf("%d lines, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestMROrderAndCounts(t *testing.T) {
	rows := mrOrder(testMRs().MRs)
	var refs []string
	for _, r := range rows {
		refs = append(refs, r.Ref)
	}
	// Projects alphabetical, own MRs before reviews.
	if strings.Join(refs, " ") != "!6 !5 #12" {
		t.Errorf("order %v", refs)
	}
	a, mine, th := mrCounts(rows)
	if a != 1 || mine != 2 || th != 2 {
		t.Errorf("counts %d %d %d", a, mine, th)
	}
}

func TestMRBlockers(t *testing.T) {
	b := mrBlockers(testNow(), mrOrder(testMRs().MRs))
	if len(b) != 2 || b[0].Item != "!6 atlas" || b[0].Why != "2 threads to answer" || b[1].Why != "sam waits for your review" {
		t.Errorf("blockers %+v", b)
	}
}

func liveMRModel(t *testing.T) (model, *[]string) {
	t.Helper()
	m := liveModel(t)
	m.page, m.mode = 0, "all"
	m.w, m.h = 170, 46
	var asked []string
	old := detailFetcher
	detailFetcher = func(r mrRow) (*mrDetail, error) {
		asked = append(asked, r.key())
		if r.Ref == "!6" {
			return nil, errors.New("forge down")
		}
		return mapDetail(rawDetail{}), nil
	}
	t.Cleanup(func() { detailFetcher = old })
	m = send(m, srcMsg{src: srcMRs, val: testMRs(), at: testNow()})
	m.page = 2
	return m, &asked
}

func TestMRsPageLive(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	out := m.View()
	for _, want := range []string{"atlas", "web", "!6", "retry", "awaiting me", "2 threads", "CI failed"} {
		if !hasText(out, want) {
			t.Errorf("view lacks %q", want)
		}
	}
	if m.fx.MRsAwaiting != 1 || m.fx.MRsThreads != 2 {
		t.Errorf("counts %d %d", m.fx.MRsAwaiting, m.fx.MRsThreads)
	}
	if len(m.fx.Blockers) != 2 {
		t.Errorf("home blockers %d, want 2 from the MRs", len(m.fx.Blockers))
	}
}

func TestMRsModeFiltersByScope(t *testing.T) {
	m, _ := liveMRModel(t)
	m.mode = "work"
	if v := visibleMRs(m.fx, m.mode); len(v) != 2 || v[0].Project != "atlas" {
		t.Errorf("work: %+v", v)
	}
	m.mode = "perso"
	if v := visibleMRs(m.fx, m.mode); len(v) != 1 || v[0].Ref != "#12" {
		t.Errorf("perso: %+v", v)
	}
}

func TestMRSelectionFetchesDetailOnce(t *testing.T) {
	m, asked := liveMRModel(t)
	m.page = 0
	next, cmd := m.Update(keyMsg("3"))
	m = next.(model)
	if cmd == nil {
		t.Fatal("opening the tab did not fetch the diff")
	}
	msg := cmd()
	if bm, ok := msg.(tea.BatchMsg); ok {
		for _, c := range bm {
			if c != nil {
				if dm, ok := c().(detailMsg); ok {
					msg = dm
				}
			}
		}
	}
	if len(*asked) != 1 || (*asked)[0] != "grp/atlas!6" {
		t.Fatalf("asked %v", *asked)
	}
	// A failing fetch shows its error and is not retried on every key.
	m = send(m, msg)
	if m.fx.MRErrs["grp/atlas!6"] != "forge down" {
		t.Errorf("errs %v", m.fx.MRErrs)
	}
	if _, c := m.Update(keyMsg("j")); c == nil {
		// moving to !5 starts its fetch
		t.Error("j did not fetch the next MR")
	}
	m.live.details["grp/atlas!6"].at = testNow().Add(-time.Hour)
	if cmd := m.wantDetail(); cmd == nil {
		t.Error("stale entry not refetched")
	}
}

func TestMRCursorFollowsItsMRAcrossRefresh(t *testing.T) {
	m, _ := liveMRModel(t)
	m = send(m, keyMsg("j"))
	if m.mrKey != "grp/atlas!5" {
		t.Fatalf("on %q", m.mrKey)
	}
	s := testMRs()
	s.MRs = append([]mrRow{{Project: "aaa", Name: "aaa", Scope: "work", Repo: "g/aaa", Number: 1, Ref: "!1", Title: "new", Role: "AUTHOR"}}, s.MRs...)
	m = send(m, srcMsg{src: srcMRs, val: s, at: testNow()})
	if v := visibleMRs(m.fx, m.mode); v[m.mrSel].Ref != "!5" {
		t.Errorf("cursor moved to %q", v[m.mrSel].Ref)
	}
}

func TestMRsForgeFailureShowsOnlyThatForge(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	s := testMRs()
	s.Forges["github"] = "github sync failed: auth"
	m = send(m, srcMsg{src: srcMRs, val: s, at: testNow()})
	if !hasText(m.View(), "github sync failed: auth") || !hasText(m.View(), "!6") {
		t.Error("the failing forge is not named next to the other forge's MRs")
	}
}

func TestMRsMockCards(t *testing.T) {
	withANSI(t)
	m := newModel(defaultConfig())
	m.page, m.w, m.h = 2, 170, 46
	m.mrSel = 2 // !220, a review
	out := m.View()
	for _, want := range []string{"Findings 3", "env var undocumented", "End of review", "nothing posted yet", "patterns hit"} {
		if !hasText(out, want) {
			t.Errorf("review view lacks %q", want)
		}
	}
	m.mrSel = 1 // !213, own MR with a thread and a saved fix
	out = m.View()
	for _, want := range []string{"Threads 2", "the fix", "the test", "End of review"} {
		if !hasText(out, want) {
			t.Errorf("thread view lacks %q", want)
		}
	}
	m = send(m, keyMsg("T"))
	if !hasText(m.View(), "Findings") {
		t.Error("T did not flip the card kind")
	}
}

func TestMRsActionKeysDoNotPretend(t *testing.T) {
	m := newModel(defaultConfig())
	m.page = 2
	for _, k := range []string{"a", "M"} {
		if n := send(m, keyMsg(k)); n.flash == "" {
			t.Errorf("%s: no flash", k)
		}
	}
}

func TestPostConfirmIsDisabled(t *testing.T) {
	m := newModel(defaultConfig())
	m.page = 2
	m = send(m, keyMsg("P"))
	if m.overlay != ovPost || !hasText(m.View(), "post (not wired)") {
		t.Fatalf("P must open the post confirm, overlay %d", m.overlay)
	}
	m = send(m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.overlay != ovNone || !strings.Contains(m.flash, "not wired") {
		t.Errorf("enter must close and say it is not wired: overlay %d, flash %q", m.overlay, m.flash)
	}
}

func TestHomeBlockersComeFromMRs(t *testing.T) {
	withANSI(t)
	m, _ := liveMRModel(t)
	m.page = 0
	out := m.View()
	for _, want := range []string{"waits for your review", "2 threads to answer", "1 awaits you"} {
		if !hasText(out, want) {
			t.Errorf("home lacks %q", want)
		}
	}
	m = send(m, srcMsg{src: srcMRs, val: &mrsSnap{Forges: map[string]string{"gitlab": "ok", "github": "ok"}}, at: testNow()})
	if !hasText(m.View(), "No MR blocks you") {
		t.Error("an MR list with nothing to do should say so")
	}
}
