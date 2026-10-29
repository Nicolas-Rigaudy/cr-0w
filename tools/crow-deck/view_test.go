package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func plain(s string) string { return ansiRe.ReplaceAllString(s, "") }

// hasText looks for text in the rendered view, ignoring colour and case (some themes capitalise).
func hasText(out, text string) bool {
	return strings.Contains(strings.ToLower(plain(out)), strings.ToLower(text))
}

func withANSI(t *testing.T) {
	t.Helper()
	lipgloss.SetColorProfile(termenv.ANSI)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
}

func send(m model, msgs ...tea.Msg) model {
	for _, msg := range msgs {
		next, _ := m.Update(msg)
		m = next.(model)
	}
	return m
}

func keyMsg(s string) tea.KeyMsg {
	switch s {
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func resize(w, h int) tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} }

// checkFit fails when a view does not fill the terminal exactly.
func checkFit(t *testing.T, label, out string, w, h int) {
	t.Helper()
	lines := strings.Split(out, "\n")
	if len(lines) != h {
		t.Errorf("%s: %d lines, want %d", label, len(lines), h)
	}
	for i, l := range lines {
		if got := lipgloss.Width(l); got > w {
			t.Errorf("%s: line %d is %d wide, max %d: %q", label, i, got, w, plain(l))
			return
		}
	}
}

var sizes = [][2]int{{170, 46}, {160, 45}, {150, 40}, {149, 40}, {118, 60}, {120, 40}, {100, 30}, {99, 30}, {90, 30}, {60, 20}}

// Every page, theme, state and size must fill the terminal and never overflow a line.
func TestViewsFitTerminal(t *testing.T) {
	withANSI(t)
	for ti, th := range themes {
		for _, sz := range sizes {
			for page := range pageNames {
				for _, st := range syncStates {
					m := newModel(defaultConfig())
					m.th, m.page, m.w, m.h = ti, page, sz[0], sz[1]
					m.setSync(st)
					checkFit(t, fmt.Sprintf("%s %s %dx%d %s", th.name, pageNames[page], sz[0], sz[1], st), m.View(), sz[0], sz[1])
				}
			}
		}
	}
}

// AC1: a resize between landscape and portrait re-lays the page out at both sizes.
func TestResizeBothOrientations(t *testing.T) {
	withANSI(t)
	landscape, portrait := [2]int{170, 46}, [2]int{118, 60}
	for ti, th := range themes {
		m := newModel(defaultConfig())
		m.th = ti
		for _, sz := range [][2]int{landscape, portrait, landscape, {90, 30}, portrait} {
			m = send(m, resize(sz[0], sz[1]))
			if m.w != sz[0] || m.h != sz[1] {
				t.Fatalf("%s: size not applied: %dx%d", th.name, m.w, m.h)
			}
			out := m.View()
			checkFit(t, fmt.Sprintf("%s resize %dx%d", th.name, sz[0], sz[1]), out, sz[0], sz[1])
			if layoutFor(sz[0]) != laySingle && !hasText(out, "Blockers") {
				t.Errorf("%s %dx%d: Blockers pane missing", th.name, sz[0], sz[1])
			}
		}
	}
}

func TestLayoutBreakpoints(t *testing.T) {
	for w, want := range map[int]int{170: layLandscape, 150: layLandscape, 149: layPortrait, 120: layPortrait, 100: layPortrait, 99: laySingle, 90: laySingle, 60: laySingle} {
		if got := layoutFor(w); got != want {
			t.Errorf("layoutFor(%d) = %d, want %d", w, got, want)
		}
	}
}

// AC2: `t` cycles the four themes, lazytf is the default, and the choice is saved.
func TestThemeCycle(t *testing.T) {
	m := newModel(defaultConfig())
	m.cfgPath = filepath.Join(t.TempDir(), "deck.conf")
	if themes[m.th].name != "lazytf" {
		t.Fatalf("default theme %q, want lazytf", themes[m.th].name)
	}
	want := []string{"hud", "cyber", "soft", "lazytf"}
	for _, name := range want {
		m = send(m, keyMsg("t"), keyMsg("esc"))
		if got := themes[m.th].name; got != name {
			t.Fatalf("after t: theme %q, want %q", got, name)
		}
	}
	if len(themes) != 4 {
		t.Errorf("%d themes, want 4", len(themes))
	}
	b, err := os.ReadFile(m.cfgPath)
	if err != nil || !strings.Contains(string(b), "theme = lazytf") {
		t.Errorf("theme not saved: %q %v", b, err)
	}
	m = send(m, keyMsg("t"))
	if m.overlay != ovTheme || !hasText(m.View(), "hud") {
		t.Error("t should open the theme overlay listing the themes")
	}
}

func TestThemesDifferAndStatusNamesTheme(t *testing.T) {
	withANSI(t)
	seen := map[string]string{}
	for ti, th := range themes {
		m := newModel(defaultConfig())
		m.th = ti
		out := m.View()
		if !hasText(out, "t theme: "+th.name) {
			t.Errorf("%s: status bar does not name the theme", th.name)
		}
		if other, ok := seen[out]; ok {
			t.Errorf("%s renders like %s", th.name, other)
		}
		seen[out] = th.name
	}
}

// ANSI 16 only: even on a truecolour terminal every colour is an index from 0 to 15.
func TestNoTruecolor(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	defer lipgloss.SetColorProfile(termenv.Ascii)
	indexed := regexp.MustCompile(`[34]8;5;(\d+)`)
	for ti, th := range themes {
		for _, page := range []int{0, 1} {
			m := newModel(defaultConfig())
			m.th, m.page = ti, page
			out := m.View()
			if strings.Contains(out, "38;2;") || strings.Contains(out, "48;2;") {
				t.Errorf("%s: truecolour sequence", th.name)
			}
			for _, sgr := range ansiRe.FindAllString(out, -1) {
				for _, mm := range indexed.FindAllStringSubmatch(sgr, -1) {
					var n int
					fmt.Sscan(mm[1], &n)
					if n > 15 {
						t.Errorf("%s: colour %d outside ANSI 16", th.name, n)
					}
				}
			}
		}
	}
}

// AC3: Home at 170, 120 and 90 columns, per the breakpoints.
func TestHomeBreakpoints(t *testing.T) {
	withANSI(t)
	titles := []string{"Usage", "Right now", "Blockers", "Today", "Next", "Projects and roadmap"}
	for ti, th := range themes {
		for _, tc := range []struct {
			w, h   int
			panes  int // panes drawn at once
			single bool
		}{{170, 46, 6, false}, {120, 50, 6, false}, {90, 30, 1, true}} {
			m := newModel(defaultConfig())
			m.th = ti
			m = send(m, resize(tc.w, tc.h))
			out := plain(m.View())
			checkFit(t, fmt.Sprintf("%s home %d", th.name, tc.w), m.View(), tc.w, tc.h)
			var got []string
			for _, ti := range titles {
				if hasText(out, ti) {
					got = append(got, ti)
				}
			}
			if len(got) < tc.panes || (tc.single && len(got) > 2) {
				t.Errorf("%s %d cols: panes shown %v, want %d", th.name, tc.w, got, tc.panes)
			}
		}
	}
}

// Under 100 columns one pane shows at a time and tab walks through all six.
func TestSinglePaneTabWalk(t *testing.T) {
	m := newModel(defaultConfig())
	m = send(m, resize(90, 30))
	seen := map[int]bool{m.focus: true}
	first := m.focus
	for i := 0; i < homePanes-1; i++ {
		before := plain(m.View())
		m = send(m, keyMsg("tab"))
		seen[m.focus] = true
		if plain(m.View()) == before {
			t.Errorf("tab %d did not change the view", i)
		}
	}
	if len(seen) != homePanes {
		t.Errorf("tab reached %d panes, want %d", len(seen), homePanes)
	}
	m = send(m, keyMsg("tab"))
	if m.focus != first {
		t.Errorf("tab does not wrap: focus %d, want %d", m.focus, first)
	}
}

func TestHomeStates(t *testing.T) {
	withANSI(t)
	for _, sz := range [][2]int{{170, 46}, {120, 50}, {90, 30}} {
		for ti, th := range themes {
			base := newModel(defaultConfig())
			base.th, base.w, base.h = ti, sz[0], sz[1]
			base.focus = pBlockers
			norm := plain(base.View())
			for st, want := range map[string]string{
				stEmpty:   "Nothing blocks you.",
				stError:   "forge sync failed",
				stLoading: "syncing",
				stStale:   "stale",
			} {
				m := base
				m.setSync(st)
				out := plain(m.View())
				checkFit(t, fmt.Sprintf("%s %dx%d %s", th.name, sz[0], sz[1], st), m.View(), sz[0], sz[1])
				if !hasText(out, want) {
					t.Errorf("%s %dx%d %s: %q not shown", th.name, sz[0], sz[1], st, want)
				}
				if out == norm {
					t.Errorf("%s %dx%d %s: same as the fresh view", th.name, sz[0], sz[1], st)
				}
			}
		}
	}
}

func TestTabsSwitchPages(t *testing.T) {
	m := newModel(defaultConfig())
	for i, name := range pageNames {
		m = send(m, keyMsg(fmt.Sprint(i+1)))
		if m.page != i {
			t.Fatalf("key %d: page %d", i+1, m.page)
		}
		if !hasText(m.View(), fmt.Sprintf("%d %s", i+1, name)) {
			t.Errorf("tab strip does not show %d %s", i+1, name)
		}
	}
}

func TestModeCycle(t *testing.T) {
	m := newModel(defaultConfig())
	for _, want := range []string{"perso", "all", "work"} {
		m = send(m, keyMsg("m"))
		if m.mode != want {
			t.Fatalf("mode %q, want %q", m.mode, want)
		}
	}
}

func TestTooSmall(t *testing.T) {
	m := newModel(defaultConfig())
	m = send(m, resize(30, 5))
	if !strings.Contains(m.View(), "too small") {
		t.Error("expected a too-small notice")
	}
}

func TestParseConfig(t *testing.T) {
	c := parseConfig("# c\ntheme = Cyber\naccent = 5\nmode = perso\nbogus = 1\n")
	if c.Theme != "cyber" || c.Accent != 5 || c.Mode != "perso" {
		t.Errorf("parsed %+v", c)
	}
	c = parseConfig("theme = nope\nmode = x\naccent = 99\n")
	if c != defaultConfig() {
		t.Errorf("bad values should fall back to defaults, got %+v", c)
	}
	if defaultConfig().Theme != "lazytf" {
		t.Error("default theme must be lazytf")
	}
	if c := loadConfig(filepath.Join(t.TempDir(), "missing.conf")); c != defaultConfig() {
		t.Errorf("missing file should give defaults, got %+v", c)
	}
}

func TestSaveKeyKeepsOtherLines(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sub", "deck.conf")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("# mine\ntheme = hud\naccent = 3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := saveKey(p, "theme", "soft"); err != nil {
		t.Fatal(err)
	}
	if err := saveKey(p, "mode", "all"); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if got, want := string(b), "# mine\ntheme = soft\naccent = 3\nmode = all\n"; got != want {
		t.Errorf("file %q, want %q", got, want)
	}
	if err := saveKey("", "theme", "hud"); err != nil {
		t.Errorf("empty path should be a no-op: %v", err)
	}
}

func TestAccentOverride(t *testing.T) {
	withANSI(t)
	cfg := defaultConfig()
	a := newModel(cfg)
	cfg.Accent = 5
	b := newModel(cfg)
	if a.View() == b.View() {
		t.Error("accent override changes nothing")
	}
	if n, ok := parseAccent("magenta"); !ok || n != 5 {
		t.Errorf("parseAccent(magenta) = %d, %v", n, ok)
	}
	if _, ok := parseAccent("16"); ok {
		t.Error("accent above 15 must be rejected")
	}
}

func TestOverlaysFit(t *testing.T) {
	withANSI(t)
	for _, sz := range [][2]int{{170, 46}, {90, 30}} {
		for ov := range []int{ovHelp, ovPalette, ovTheme} {
			m := newModel(defaultConfig())
			m.w, m.h = sz[0], sz[1]
			m.overlay = []int{ovHelp, ovPalette, ovTheme}[ov]
			checkFit(t, fmt.Sprintf("overlay %d %dx%d", m.overlay, sz[0], sz[1]), m.View(), sz[0], sz[1])
		}
	}
}
