package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestShots writes ANSI renders for bin/ansi-shot:
//
//	SHOTS=/abs/dir go test -run TestShots .
//
// Files: <theme>-<page>-<size>[-<state>].ans. Sizes: wide 170x46, laptop 160x45,
// portrait 118x60, narrow 120x40, small 90x30, board 235x62 (the round 2 landscape boards).
func TestShots(t *testing.T) {
	dir := os.Getenv("SHOTS")
	if dir == "" {
		t.Skip("set SHOTS=<dir>")
	}
	lipgloss.SetColorProfile(termenv.ANSI)
	defer lipgloss.SetColorProfile(termenv.Ascii)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	sizes := map[string][2]int{"wide": {170, 46}, "portrait": {118, 60}, "narrow": {120, 40}, "small": {90, 30}, "board": {235, 62}}
	write := func(name string, m model) {
		f := filepath.Join(dir, name+".ans")
		if err := os.WriteFile(f, []byte(m.View()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for ti, th := range themes {
		for sn, sz := range sizes {
			m := newModel(defaultConfig())
			m.th, m.w, m.h = ti, sz[0], sz[1]
			write(fmt.Sprintf("%s-home-%s", th.name, sn), m)
		}
		for p := 1; p < len(pageNames); p++ {
			for _, sn := range []string{"wide", "portrait", "small"} {
				m := newModel(defaultConfig())
				m.th, m.page, m.w, m.h = ti, p, sizes[sn][0], sizes[sn][1]
				write(fmt.Sprintf("%s-%s-%s", th.name, pageNames[p], sn), m)
			}
		}
		agentShots(write, ti, sizes)
		mrsShots(write, ti, sizes)
		usageShots(write, ti, sizes)
		histShots(write, ti, sizes)
		projectShots(write, ti, sizes)
		for _, st := range []string{stEmpty, stError, stLoading, stStale} {
			for _, sn := range []string{"wide", "small"} {
				m := newModel(defaultConfig())
				m.th, m.w, m.h = ti, sizes[sn][0], sizes[sn][1]
				m.setSync(st)
				write(fmt.Sprintf("%s-home-%s-%s", th.name, sn, st), m)
			}
		}
		m := newModel(defaultConfig())
		m.th, m.mode = ti, "all"
		write(fmt.Sprintf("%s-home-wide-all", th.name), m)
		m.overlay = ovTheme
		write(fmt.Sprintf("%s-overlay-theme", th.name), m)
		m.overlay = ovHelp
		write(fmt.Sprintf("%s-overlay-help", th.name), m)
		m.overlay, m.query = ovPalette, "th"
		write(fmt.Sprintf("%s-overlay-palette", th.name), m)
		m.page, m.overlay = 2, ovPost
		write(fmt.Sprintf("%s-overlay-post", th.name), m)
		m.page, m.overlay = 1, ovConfirm
		selectLane(&m, func(l laneRow) bool { return l.State == lsPaused })
		write(fmt.Sprintf("%s-overlay-relaunch", th.name), m)
		nm := newModel(defaultConfig())
		nm.th, nm.w, nm.h = ti, sizes["wide"][0], sizes["wide"][1]
		nm.fx.NoModel, nm.fx.NoModelUntil = noModelRows(nm.fx.Decide, nil, nm.fx.Agents), "15:20"
		nm.fx.Decide = nil
		write(fmt.Sprintf("%s-home-wide-nomodel", th.name), nm)
	}
}

// agentShots: the Agents tab beyond the default view: a running agent, the paused run, the
// answer line in use, every pane in single-pane mode, and the data states.
func agentShots(write func(string, model), ti int, sizes map[string][2]int) {
	name := themes[ti].name
	base := func(sn string) model {
		m := newModel(defaultConfig())
		m.th, m.page, m.w, m.h = ti, 1, sizes[sn][0], sizes[sn][1]
		return m
	}
	for _, sn := range []string{"board", "portrait"} {
		m := base(sn)
		write(name+"-Agents-"+sn, m)
		selectLane(&m, func(l laneRow) bool { return l.State == lsWorking })
		write(name+"-Agents-"+sn+"-running", m)
		selectLane(&m, func(l laneRow) bool { return l.State == lsPaused })
		write(name+"-Agents-"+sn+"-paused", m)
		selectLane(&m, func(l laneRow) bool { return l.State == lsDone })
		write(name+"-Agents-"+sn+"-done", m)
		m = base(sn)
		m.agentPane = apDiff
		write(name+"-Agents-"+sn+"-diff", m)
		m.agentPane, m.agentFull = apDetail, true
		write(name+"-Agents-"+sn+"-detail-fullday", m)
		m.agentZoom = true
		write(name+"-Agents-"+sn+"-zoom", m)
	}
	m := base("board")
	m.input, m.typing = "yes, but keep the scope", true
	write(name+"-Agents-board-typing", m)
	for pane, pn := range agentPaneNames {
		m = base("small")
		m.agentPane = pane
		write(fmt.Sprintf("%s-Agents-small-%s", name, pn), m)
	}
	for _, st := range []string{stEmpty, stError, stLoading, stStale} {
		m = base("wide")
		m.setSync(st)
		write(fmt.Sprintf("%s-Agents-wide-%s", name, st), m)
	}
}

// mrsShots: the MRs tab beyond the default view: a review with findings, an own MR with a
// thread and its fix, the flipped card, every pane in single-pane mode, and the data states.
func mrsShots(write func(string, model), ti int, sizes map[string][2]int) {
	name := themes[ti].name
	base := func(sn string) model {
		m := newModel(defaultConfig())
		m.th, m.page, m.w, m.h = ti, 2, sizes[sn][0], sizes[sn][1]
		return m
	}
	// The list puts your move first: !213 (threads), !220 (review), then !214 (waiting).
	m := base("wide")
	m.mrSel = 1
	write(name+"-MRs-wide-review", m)
	m.mrItem = 1
	write(name+"-MRs-wide-review-2", m)
	m = base("board")
	m.mrSel = 1
	write(name+"-MRs-board-review", m)
	m = base("wide")
	m.mrSel = 0
	write(name+"-MRs-wide-threads", m)
	m = base("wide")
	m.mrSel, m.mrToggle = 0, true
	write(name+"-MRs-wide-threads-flipped", m)
	m = base("wide")
	m.mrSel = 2
	write(name+"-MRs-wide-waiting", m)
	m = base("wide")
	m.mode = "all"
	write(name+"-MRs-wide-all", m)
	m = base("portrait")
	m.mrSel = 1
	write(name+"-MRs-portrait-review", m)
	m.mrPane = mpList
	write(name+"-MRs-portrait-list", m)
	for pane, pn := range []string{"list", "code", "move"} {
		m = base("small")
		m.mrSel, m.mrPane = 1, pane
		write(fmt.Sprintf("%s-MRs-small-%s", name, pn), m)
	}
	for _, st := range []string{stEmpty, stError, stLoading, stStale} {
		m = base("wide")
		m.setSync(st)
		write(fmt.Sprintf("%s-MRs-wide-%s", name, st), m)
	}
	m = base("wide")
	m.fx.MRs, m.fx.MRDetails = m.fx.MRs[:0:0], nil
	m.fx.ForgeErrs = []string{"github sync failed: gh auth expired"}
	write(name+"-MRs-wide-forge-error", m)
}

// usageShots: the Usage tab beyond the default view: board and portrait sizes for each
// period, the ledger flat and filtered, every pane in single-pane mode, and the data states.
func usageShots(write func(string, model), ti int, sizes map[string][2]int) {
	name := themes[ti].name
	base := func(sn string) model {
		m := newModel(defaultConfig())
		m.th, m.page, m.w, m.h = ti, 4, sizes[sn][0], sizes[sn][1]
		return m
	}
	for _, sn := range []string{"board", "portrait"} {
		for p, pn := range periodNames {
			m := base(sn)
			m.usePeriod = p
			write(fmt.Sprintf("%s-Usage-%s-%s", name, sn, pn), m)
		}
	}
	m := base("board")
	m.usePane, m.useFlat, m.useTok = upLedger, true, true
	write(name+"-Usage-board-flat", m)
	m = base("board")
	m.usePane, m.useFilter = upProjects, "atlas"
	write(name+"-Usage-board-filter", m)
	for pane, pn := range []string{"runway", "spend", "saved", "ledger", "projects", "models"} {
		m = base("small")
		m.usePane = pane
		write(fmt.Sprintf("%s-Usage-small-%s", name, pn), m)
	}
	for _, st := range []string{stEmpty, stError, stLoading, stStale} {
		m = base("board")
		m.setSync(st)
		write(fmt.Sprintf("%s-Usage-board-%s", name, st), m)
	}
}

// projectShots: the Projects tab beyond the default view: the zoomed timeline, a card
// selected, every pane in single-pane mode, and the data states.
func projectShots(write func(string, model), ti int, sizes map[string][2]int) {
	name := themes[ti].name
	base := func(sn string) model {
		m := newModel(defaultConfig())
		m.th, m.page, m.w, m.h, m.mode = ti, 3, sizes[sn][0], sizes[sn][1], "all"
		return m
	}
	m := base("wide")
	write(name+"-Projects-wide-all", m)
	m.prjSel = 2
	write(name+"-Projects-wide-lint", m)
	m.prjZoom = true
	write(name+"-Projects-wide-zoom", m)
	m = base("wide")
	m.mode = "work"
	write(name+"-Projects-wide-work", m)
	for _, sn := range []string{"board", "portrait"} {
		m = base(sn)
		write(fmt.Sprintf("%s-Projects-%s-all", name, sn), m)
	}
	for pane, pn := range []string{"list", "detail", "timeline"} {
		m = base("small")
		m.prjPane = pane
		write(fmt.Sprintf("%s-Projects-small-%s", name, pn), m)
	}
	for _, st := range []string{stEmpty, stError, stLoading, stStale} {
		m = base("wide")
		m.setSync(st)
		write(fmt.Sprintf("%s-Projects-wide-%s", name, st), m)
	}
}

// histShots: the History view behind / in Agents: the list, a search, a filter, a live-style
// row without detail, every pane in single-pane mode, and the data states.
func histShots(write func(string, model), ti int, sizes map[string][2]int) {
	name := themes[ti].name
	base := func(sn string) model {
		m := newModel(defaultConfig())
		m.th, m.page, m.w, m.h, m.mode, m.hist = ti, 1, sizes[sn][0], sizes[sn][1], "all", true
		return m
	}
	for _, sn := range []string{"wide", "portrait", "small"} {
		m := base(sn)
		m.histSel = 1
		write(fmt.Sprintf("%s-History-%s", name, sn), m)
	}
	m := base("wide")
	m.histF.query, m.histTyping = "retry", true
	write(name+"-History-wide-search", m)
	m = base("wide")
	m.histF.account, m.histF.model, m.histF.days = "main", "sonnet", 7
	write(name+"-History-wide-filtered", m)
	m = base("wide")
	m.histF.query = "zzz"
	write(name+"-History-wide-nomatch", m)
	m = base("wide")
	m.histSel = 3
	for i := range m.fx.Hist {
		m.fx.Hist[i].Detail = false
	}
	write(name+"-History-wide-nodetail", m)
	for pane, pn := range []string{"list", "preview"} {
		m = base("small")
		m.histPane = pane
		write(fmt.Sprintf("%s-History-small-%s", name, pn), m)
	}
	for _, st := range []string{stEmpty, stError, stLoading, stStale} {
		m = base("wide")
		m.setSync(st)
		write(fmt.Sprintf("%s-History-wide-%s", name, st), m)
	}
}

// TestLiveShots renders Home from the real sources (crow scripts on PATH or in bin/):
//
//	LIVESHOTS=/abs/dir go test -run TestLiveShots .
//
// Files: home-, agents- and projects-live-<w>x<h>.ans at the board sizes, landscape 235x62 and 170x46, portrait 118x60.
func TestLiveShots(t *testing.T) {
	dir := os.Getenv("LIVESHOTS")
	if dir == "" {
		t.Skip("set LIVESHOTS=<dir>")
	}
	lipgloss.SetColorProfile(termenv.ANSI)
	defer lipgloss.SetColorProfile(termenv.Ascii)
	m := newModel(defaultConfig())
	m.now = time.Now()
	m.goLive()
	for _, src := range allSources() {
		if msg, ok := fetchCmd(src)().(srcMsg); ok {
			m.live.apply(msg)
		}
	}
	m.fx = buildFixtures(liveInput{Now: m.now, LS: m.live})
	m.sync = stFresh
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, sz := range [][2]int{{235, 62}, {170, 46}, {118, 60}} {
		for _, mode := range []string{"all", "work"} {
			m.w, m.h, m.mode, m.page = sz[0], sz[1], mode, 0
			f := filepath.Join(dir, fmt.Sprintf("home-live-%dx%d-%s.ans", sz[0], sz[1], mode))
			if err := os.WriteFile(f, []byte(m.View()), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		if sz[0] >= 150 {
			m.mode, m.focus, m.sel[pToday] = "all", pToday, 1
			f := filepath.Join(dir, fmt.Sprintf("home-today-%dx%d.ans", sz[0], sz[1]))
			if err := os.WriteFile(f, []byte(m.View()), 0o644); err != nil {
				t.Fatal(err)
			}
			m.focus = pNext
		}
		m.page, m.mode = 1, "all"
		f := filepath.Join(dir, fmt.Sprintf("agents-live-%dx%d.ans", sz[0], sz[1]))
		if err := os.WriteFile(f, []byte(m.View()), 0o644); err != nil {
			t.Fatal(err)
		}
		m.page = 3
		for _, mode := range []string{"all", "work"} {
			m.mode = mode
			f = filepath.Join(dir, fmt.Sprintf("projects-live-%dx%d-%s.ans", sz[0], sz[1], mode))
			if err := os.WriteFile(f, []byte(m.View()), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		m.mode = "all"
	}
}

// TestShotsReal renders a saved MR detail (crow-mrs show JSON) at the four walkthrough stages:
//
//	SHOTS=/abs/dir REAL_MR=/abs/detail.json go test -run TestShotsReal .
//
// The detail stands in for the mock !214. Work data: keep the JSON and the renders out of git.
func TestShotsReal(t *testing.T) {
	dir, src := os.Getenv("SHOTS"), os.Getenv("REAL_MR")
	if dir == "" || src == "" {
		t.Skip("set SHOTS=<dir> REAL_MR=<detail.json>")
	}
	lipgloss.SetColorProfile(termenv.ANSI)
	defer lipgloss.SetColorProfile(termenv.Ascii)
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	stages := map[string]func(d *mrDetail, r *mrRow){
		"1-deciding": func(d *mrDetail, r *mrRow) {
			for i := range d.Findings {
				d.Findings[i].State, d.Findings[i].Applied = "pending", ""
			}
			d.Findings[0].State = "accepted"
			d.Apply = nil
		},
		"2-applying": func(d *mrDetail, r *mrRow) {
			for i := range d.Findings {
				d.Findings[i].State, d.Findings[i].Applied = "accepted", ""
			}
			d.Apply = &mrApply{State: "running", StartedAt: mockNow().Add(-2 * time.Minute).UTC().Format(time.RFC3339)}
		},
		"3-pushed": func(d *mrDetail, r *mrRow) {
			for i := range d.Findings {
				d.Findings[i].State, d.Findings[i].Applied, d.Findings[i].AppliedFor = "accepted", "fixed", [2]string{"accepted", ""}
			}
			d.Apply.State, d.Apply.Pushed = "done", true
		},
		"4-done": func(d *mrDetail, r *mrRow) {
			for i := range d.Findings {
				d.Findings[i].State, d.Findings[i].Applied, d.Findings[i].AppliedFor = "accepted", "fixed", [2]string{"accepted", ""}
				d.Findings[i].Resolved, d.Findings[i].Posted = true, true
			}
			d.Apply.State, d.Apply.Pushed = "done", true
			r.Pipeline, r.ThreadsOpen = "passed", 0
		},
	}
	sizes := map[string][2]int{"board": {235, 62}, "wide": {170, 46}, "portrait": {118, 60}}
	for stage, tweak := range stages {
		for sn, sz := range sizes {
			var raw rawDetail
			if err := json.Unmarshal(b, &raw); err != nil {
				t.Fatal(err)
			}
			m := newModel(defaultConfig())
			m.page, m.w, m.h, m.mode = 2, sz[0], sz[1], "all"
			d := mapDetail(raw)
			for i, r := range m.fx.MRs {
				if r.Ref == "!214" {
					m.fx.MRs[i].ThreadsOpen = len(d.Threads)
					tweak(d, &m.fx.MRs[i])
					m.fx.MRDetails[r.key()] = d
				}
			}
			for i, r := range visibleMRs(m.fx, m.mode) {
				if r.Ref == "!214" {
					m.selectMR(i)
				}
			}
			m.followItem()
			if stage == "1-deciding" {
				m.mrItem = 1
				m.followItem()
			}
			f := filepath.Join(dir, "real214-"+stage+"-"+sn+".ans")
			if err := os.WriteFile(f, []byte(m.View()), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
}
