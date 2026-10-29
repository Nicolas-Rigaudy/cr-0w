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
			m := newModel(defaultConfig())
			m.th, m.page = ti, p
			write(fmt.Sprintf("%s-%s-wide", th.name, pageNames[p]), m)
		}
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
