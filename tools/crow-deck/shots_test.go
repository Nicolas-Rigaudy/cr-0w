package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestShots writes ANSI renders for bin/ansi-shot:
//
//	SHOTS=/abs/dir go test -run TestShots .
//
// Files: <theme>-<page>-<size>[-<state>].ans. Sizes: wide 170x46, laptop 160x45,
// portrait 118x60, narrow 120x40, small 90x30.
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
	sizes := map[string][2]int{"wide": {170, 46}, "portrait": {118, 60}, "narrow": {120, 40}, "small": {90, 30}}
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
	}
	m := newModel(defaultConfig())
	m.overlay = ovHelp
	write("overlay-help", m)
	m.overlay, m.query = ovPalette, "th"
	write("overlay-palette", m)
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
	m := base("wide")
	m.selectAgent(2)
	write(name+"-Agents-wide-running", m)
	m.selectAgent(4)
	write(name+"-Agents-wide-paused", m)
	m = base("wide")
	m.input, m.typing = "yes, but keep the scope", true
	write(name+"-Agents-wide-typing", m)
	for pane, pn := range []string{"list", "tail", "who"} {
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
	m := base("wide")
	m.mrSel = 2
	write(name+"-MRs-wide-review", m)
	m.mrItem = 1
	write(name+"-MRs-wide-review-2", m)
	m = base("wide")
	m.mrSel = 1
	write(name+"-MRs-wide-threads", m)
	m = base("wide")
	m.mrSel, m.mrToggle = 1, true
	write(name+"-MRs-wide-threads-flipped", m)
	m = base("wide")
	m.mode = "all"
	write(name+"-MRs-wide-all", m)
	m = base("portrait")
	m.mrSel = 2
	write(name+"-MRs-portrait-review", m)
	for pane, pn := range []string{"list", "diff", "card"} {
		m = base("small")
		m.mrSel, m.mrPane = 2, pane
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
