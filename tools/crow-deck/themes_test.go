package main

import (
	"regexp"
	"strings"
	"testing"
)

var sgrRe = regexp.MustCompile(`\x1b\[([0-9;]*)m`)

// mono drops colour SGR parameters and keeps bold, underline and reverse, like a
// monochrome terminal would show the screen.
func mono(s string) string {
	return sgrRe.ReplaceAllStringFunc(s, func(m string) string {
		var keep []string
		for _, p := range strings.Split(sgrRe.FindStringSubmatch(m)[1], ";") {
			switch p {
			case "1", "4", "7":
				keep = append(keep, p)
			}
		}
		if len(keep) == 0 {
			return ""
		}
		return "\x1b[" + strings.Join(keep, ";") + "m"
	})
}

// AC: themes differ by shape: Home and MRs renders differ between every pair of
// themes once colour is stripped, and each gauge has its own glyphs.
func TestThemesDifferInMonochrome(t *testing.T) {
	withANSI(t)
	for _, page := range []int{0, 2} {
		seen := map[string]string{}
		for ti, th := range themes {
			m := newModel(defaultConfig())
			m.th, m.page, m.w, m.h = ti, page, 170, 46
			out := mono(m.View())
			if other, ok := seen[out]; ok {
				t.Errorf("page %d: %s and %s look the same in monochrome", page, th.name, other)
			}
			seen[out] = th.name
		}
	}
	shapes := map[string][]string{"lazytf": {"━", "─"}, "soft": {"●", "○"}, "cyber": {"[", "/", ".", "]"}, "hud": {"▕", "▮", "▯", "│", "▏"}}
	for _, th := range themes {
		lk := newLook(th, -1)
		g := plainText(gaugeBar(lk, 40, 12, lk.ok))
		for _, glyph := range shapes[th.name] {
			if !strings.Contains(g, glyph) {
				t.Errorf("%s gauge %q has no %q", th.name, g, glyph)
			}
		}
		if n := len([]rune(g)); n != 12 {
			t.Errorf("%s gauge %q is %d cells, want 12", th.name, g, n)
		}
	}
}

// HUD gauge: ▕ start, ▏ end, │ tick at 80 % of the track while the fill is below it.
func TestHUDGaugeTick(t *testing.T) {
	lk := newLook(themes[themeIndex("hud")], -1)
	if got := plainText(gaugeBar(lk, 50, 12, lk.ok)); got != "▕▮▮▮▮▮▯▯▯│▯▏" {
		t.Errorf("50%%: %q", got)
	}
	if got := plainText(gaugeBar(lk, 95, 12, lk.ok)); strings.Contains(got, "│") {
		t.Errorf("95%%: tick drawn over the fill: %q", got)
	}
}

// AC2: only the 16 ANSI foreground slots; no background fills, no 256 or true colour.
func TestThemesUseANSI16Only(t *testing.T) {
	withANSI(t)
	for ti, th := range themes {
		for page := range pageNames {
			m := newModel(defaultConfig())
			m.th, m.page, m.w, m.h = ti, page, 170, 46
			for _, sm := range sgrRe.FindAllStringSubmatch(m.View(), -1) {
				ps := strings.Split(sm[1], ";")
				for _, p := range ps {
					switch {
					case p == "38" || p == "48":
						t.Fatalf("%s %s: extended colour %q", th.name, pageNames[page], sm[0])
					case len(p) == 2 && (p[0] == '4' && p[1] != '9') || len(p) == 3 && p[:2] == "10":
						t.Fatalf("%s %s: background fill %q", th.name, pageNames[page], sm[0])
					}
				}
			}
		}
	}
}

// Every theme shares the one accent: the config accent, else the default slot.
func TestThemesShareAccent(t *testing.T) {
	for _, th := range themes {
		if lk := newLook(th, -1); lk.accent != defaultAccent {
			t.Errorf("%s accent %d, want %d", th.name, lk.accent, defaultAccent)
		}
		if lk := newLook(th, 5); lk.accent != 5 {
			t.Errorf("%s ignores the config accent", th.name)
		}
	}
}

// Group headers stay in sentence case in the caps themes.
func TestGroupHeadersSentenceCase(t *testing.T) {
	withANSI(t)
	for _, name := range []string{"cyber", "hud"} {
		m := newModel(defaultConfig())
		m.th, m.w, m.h = themeIndex(name), 170, 46
		if out := plain(m.View()); !strings.Contains(out, "You decide") || strings.Contains(out, "YOU DECIDE") {
			t.Errorf("%s: group header not in sentence case", name)
		}
	}
}
