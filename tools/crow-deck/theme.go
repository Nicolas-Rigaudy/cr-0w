package main

// Theme engine: a theme is a role-to-slot table, a frame kind and style flags.
//
// Colours are the terminal's 16 ANSI slots only, so the terminal theme decides the real
// colours. Secondary text is slot 7; slot 8 is for borders, leaders and gauge tracks.
// The only filled cells are the selected row and the mode chip (accent fill, dark text,
// done with reverse video so the text takes the terminal background).

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type frameKind int

const (
	frameRound frameKind = iota // lazytf: rounded border, title on the border
	frameHUD                    // brackets, tags, panel index
	frameCyber                  // single rules, //, corner notch
	frameSoft                   // no lines, spacing only
)

const follow = -1 // role follows the accent slot

type theme struct {
	name   string
	frame  frameKind
	accent int // default accent slot (config can override)
	title  int // box title slot, follow = accent
	key    int // key hint slot
	label  int // right labels and links
	head   int // table header slot
	caps   bool
}

// Slots: 1 red, 2 green, 3 yellow, 4 blue, 5 magenta, 6 cyan, 7 white.
var themes = []theme{
	{"lazytf", frameRound, 4, 5, 4, 4, 7, false},
	{"hud", frameHUD, 6, follow, follow, follow, follow, true},
	{"cyber", frameCyber, 5, follow, 4, 6, 6, true},
	{"soft", frameSoft, 6, follow, 4, 4, 7, false},
}

func themeIndex(name string) int {
	for i, t := range themes {
		if t.name == strings.ToLower(strings.TrimSpace(name)) {
			return i
		}
	}
	return -1
}

func themeNames() []string {
	out := make([]string, len(themes))
	for i, t := range themes {
		out[i] = t.name
	}
	return out
}

var slotNames = []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}

// parseAccent reads a slot as a name (blue, bright-blue) or a number 0-15.
func parseAccent(s string) (int, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	if n, err := strconv.Atoi(s); err == nil {
		return n, n >= 0 && n <= 15
	}
	bright := 0
	if rest, ok := strings.CutPrefix(s, "bright-"); ok {
		s, bright = rest, 8
	}
	for i, n := range slotNames {
		if n == s {
			return i + bright, true
		}
	}
	return 0, false
}

func slotName(n int) string {
	if n < 8 {
		return slotNames[n]
	}
	return "bright-" + slotNames[n-8]
}

func slot(n int) lipgloss.Color { return lipgloss.Color(strconv.Itoa(n)) }

// look is a theme resolved to styles, for one accent.
type look struct {
	th     theme
	accent int
	custom int // accent set in config or by flag, -1 none

	acc, title, key, label, head       lipgloss.Style
	ok, warn, bad, mag, cyan, sub, dim lipgloss.Style
	plain, bold, sel, underlineAcc     lipgloss.Style
}

func newLook(th theme, accent int) *look {
	custom := accent
	if accent < 0 {
		accent = th.accent
	}
	pick := func(s int) int {
		if s == follow {
			return accent
		}
		return s
	}
	fg := func(n int) lipgloss.Style { return lipgloss.NewStyle().Foreground(slot(n)) }
	return &look{
		th: th, accent: accent, custom: custom,
		acc:   fg(accent).Bold(true),
		title: fg(pick(th.title)).Bold(true),
		key:   fg(pick(th.key)),
		label: fg(pick(th.label)),
		head:  fg(pick(th.head)),
		ok:    fg(2), warn: fg(3), bad: fg(1), mag: fg(5), cyan: fg(6),
		sub: fg(7), dim: fg(8),
		plain: lipgloss.NewStyle(), bold: lipgloss.NewStyle().Bold(true),
		sel:          fg(accent).Reverse(true),
		underlineAcc: fg(accent).Bold(true).Underline(true),
	}
}

// cap uppercases text for themes that draw captions in capitals.
func (lk *look) cap(s string) string {
	if lk.th.caps {
		return strings.ToUpper(s)
	}
	return s
}
