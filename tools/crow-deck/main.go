// crow-deck: the chief of staff's terminal control room (skeleton: chrome, tabs, themes,
// config and a mock Home). Run `crow deck`. Spec: local/state/research/crow-deck-design.md.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	w, h    int
	page    int
	th      int
	accent  int // slot, -1 = theme default
	mode    string
	sync    string
	fx      *fixtures
	focus   int
	sel     [homePanes]int
	overlay int
	query   string
	palSel  int
	flash   string
	cfgPath string
}

func newModel(cfg config) model {
	th := max(0, themeIndex(cfg.Theme))
	return model{th: th, accent: cfg.Accent, mode: cfg.Mode, sync: stFresh, fx: mockFixtures(), focus: pBlockers, w: 160, h: 45}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) look() *look { return newLook(themes[m.th], m.accent) }

func (m *model) setSync(s string) {
	m.sync = s
	m.fx = mockFixtures()
	if s == stEmpty {
		m.fx = emptyFixtures()
		m.sync = stFresh
	}
	m.sel = [homePanes]int{}
}

func (m *model) cycleTheme(d int) {
	m.th = (m.th + d + len(themes)) % len(themes)
	if err := saveKey(m.cfgPath, "theme", themes[m.th].name); err != nil {
		m.flash = "config not saved: " + err.Error()
	}
}

func (m model) homeView() homeView {
	return homeView{lk: m.look(), fx: m.fx, mode: m.mode, sync: m.sync, focus: m.focus, sel: m.sel}
}

func (m *model) moveSel(d int) {
	if m.page != 0 {
		return
	}
	hv := m.homeView()
	if !hv.selectable(m.focus) {
		return
	}
	n := len(hv.spec(m.focus, 80).rows)
	m.sel[m.focus] = max(0, min(n-1, m.sel[m.focus]+d))
}

func (m *model) moveFocus(d int) {
	if m.page == 0 {
		m.focus = (m.focus + d + homePanes) % homePanes
	}
}

// exec runs a palette command.
func (m *model) exec(c string) tea.Cmd {
	f := strings.Fields(c)
	if len(f) == 0 {
		return nil
	}
	arg := ""
	if len(f) > 1 {
		arg = f[1]
	}
	switch f[0] {
	case "quit":
		return tea.Quit
	case "tab":
		for i, n := range pageNames {
			if strings.EqualFold(n, arg) {
				m.page = i
			}
		}
	case "theme":
		if i := themeIndex(arg); i >= 0 {
			m.th = i
			_ = saveKey(m.cfgPath, "theme", arg)
		}
	case "mode":
		if validMode(arg) {
			m.mode = arg
		}
	case "mock":
		m.setSync(arg)
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyMsg:
		return m.key(msg)
	}
	return m, nil
}

func (m model) key(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()
	m.flash = ""
	if k == "ctrl+c" {
		return m, tea.Quit
	}
	switch m.overlay {
	case ovPalette:
		return m.paletteKey(msg)
	case ovTheme:
		switch k {
		case "t", "j", "down":
			m.cycleTheme(1)
		case "k", "up":
			m.cycleTheme(-1)
		case "esc", "enter":
			m.overlay = ovNone
		case "q":
			return m, tea.Quit
		}
		return m, nil
	case ovHelp:
		switch k {
		case "esc", "enter", "?":
			m.overlay = ovNone
		case "q":
			return m, tea.Quit
		}
		return m, nil
	}
	switch k {
	case "q":
		return m, tea.Quit
	case "1", "2", "3", "4", "5":
		m.page = int(k[0] - '1')
	case "t":
		m.cycleTheme(1)
		m.overlay = ovTheme
	case "m":
		m.mode = map[string]string{"work": "perso", "perso": "all", "all": "work"}[m.mode]
	case "?":
		m.overlay = ovHelp
	case ":":
		m.overlay, m.query, m.palSel = ovPalette, "", 0
	case "tab":
		m.moveFocus(1)
	case "shift+tab":
		m.moveFocus(-1)
	case "j", "down":
		m.moveSel(1)
	case "k", "up":
		m.moveSel(-1)
	case "/":
		m.flash = "/ search and filter come with the list pages"
	case "enter", "h", "r", "x", "R", "a", "e", "c":
		m.flash = k + ": actions are not wired in the skeleton"
	}
	return m, nil
}

func (m model) paletteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.overlay = ovNone
	case tea.KeyEnter:
		cmds := filterCommands(m.query)
		m.overlay = ovNone
		if m.palSel < len(cmds) {
			return m, m.exec(cmds[m.palSel].name)
		}
	case tea.KeyDown, tea.KeyTab:
		m.palSel = min(m.palSel+1, max(0, len(filterCommands(m.query))-1))
	case tea.KeyUp, tea.KeyShiftTab:
		m.palSel = max(0, m.palSel-1)
	case tea.KeyBackspace:
		if r := []rune(m.query); len(r) > 0 {
			m.query = string(r[:len(r)-1])
		}
		m.palSel = 0
	case tea.KeyRunes, tea.KeySpace:
		m.query += string(msg.Runes)
		if msg.Type == tea.KeySpace {
			m.query += " "
		}
		m.palSel = 0
	}
	return m, nil
}

func (m model) context() []keyHint {
	if m.page != 0 {
		return []keyHint{{"1-5", "tab"}, {"j k", "move"}}
	}
	ctx := []keyHint{{"1-5", "tab"}}
	if layoutFor(m.w) == laySingle {
		ctx = append(ctx, keyHint{"tab", "pane"})
	}
	return append(ctx, keyHint{"enter", "open"}, keyHint{"m", "work / perso / all"}, keyHint{"h", "handoff"}, keyHint{"r", "review"}, keyHint{"x", "snooze"})
}

func (m model) View() string {
	if m.w < 40 || m.h < 8 {
		return "terminal too small" + strings.Repeat("\n", max(0, m.h-1))
	}
	lk := m.look()
	cs := chromeState{page: m.page, mode: m.mode, sync: m.sync, fx: m.fx, flash: m.flash, context: m.context()}
	bh := m.h - 4
	var body []string
	switch {
	case m.overlay == ovHelp:
		body = lk.helpOverlay(m.w, bh, m.page)
	case m.overlay == ovPalette:
		body = lk.paletteOverlay(m.w, bh, m.query, m.palSel)
	case m.overlay == ovTheme:
		body = lk.themeOverlay(m.w, bh, m.th)
	case m.page == 0:
		body = m.homeView().body(m.w, bh)
	default:
		body = lk.placeholder(m.page, m.w, bh, m.sync)
	}
	out := []string{lk.header(m.w, cs), lk.tabs(m.w, cs), lk.rule(m.w)}
	out = append(out, body...)
	out = append(out, lk.status(m.w, cs, themes[m.th].name))
	return strings.Join(out, "\n")
}

func main() {
	cfgPath := configPath()
	cfg := loadConfig(cfgPath)
	themeFlag := flag.String("theme", "", "theme: "+strings.Join(themeNames(), ", ")+" (default from config)")
	accentFlag := flag.String("accent", "", "accent slot: name or 0-15 (default from config)")
	modeFlag := flag.String("mode", "", "mode: all, work, perso")
	stateFlag := flag.String("state", "", "mock state: "+strings.Join(syncStates, ", "))
	flag.Parse()
	if *themeFlag != "" {
		if themeIndex(*themeFlag) < 0 {
			fmt.Fprintf(os.Stderr, "crow-deck: unknown theme %q (%s)\n", *themeFlag, strings.Join(themeNames(), ", "))
			os.Exit(2)
		}
		cfg.Theme = strings.ToLower(*themeFlag)
	}
	if *accentFlag != "" {
		n, ok := parseAccent(*accentFlag)
		if !ok {
			fmt.Fprintf(os.Stderr, "crow-deck: unknown accent %q (name or 0-15)\n", *accentFlag)
			os.Exit(2)
		}
		cfg.Accent = n
	}
	if *modeFlag != "" {
		if !validMode(*modeFlag) {
			fmt.Fprintf(os.Stderr, "crow-deck: unknown mode %q (all, work, perso)\n", *modeFlag)
			os.Exit(2)
		}
		cfg.Mode = *modeFlag
	}
	m := newModel(cfg)
	m.cfgPath = cfgPath
	if *stateFlag != "" {
		m.setSync(*stateFlag)
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "crow-deck:", err)
		os.Exit(1)
	}
}
