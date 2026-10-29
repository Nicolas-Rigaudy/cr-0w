// crow-deck: the chief of staff's terminal control room (skeleton: chrome, tabs, themes,
// config and a mock Home). Run `crow deck`. Spec: local/state/research/crow-deck-design.md.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

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

	live  *liveState // nil in mock mode
	now   time.Time
	ticks int

	agentSel  int
	agentKey  string // session or run id of the selected agent, keeps the cursor on it across refreshes
	agentPane int
	input     string
	typing    bool

	mrSel, mrPane, mrFile, mrItem int
	mrKey                         string // key of the selected MR, keeps the cursor on it across refreshes
	mrToggle                      bool

	prjSel, prjPane int
	prjZoom         bool

	usePane, useSel int
	useFortnight    bool
}

func newModel(cfg config) model {
	th := max(0, themeIndex(cfg.Theme))
	return model{th: th, accent: cfg.Accent, mode: cfg.Mode, sync: stFresh, fx: mockFixtures(), focus: pBlockers, w: 160, h: 45, now: mockNow()}
}

// mockNow is a fixed time of day (14:07) so mock renders are stable.
func mockNow() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 14, 7, 0, 0, n.Location())
}

// goLive switches the model to real data; Init then starts the fetches.
func (m *model) goLive() {
	m.now = time.Now()
	m.live = newLiveState()
	m.sync = stLoading
	m.fx = buildFixtures(liveInput{Now: m.now, LS: m.live})
}

func (m model) Init() tea.Cmd {
	if m.live == nil {
		return nil
	}
	var cmds []tea.Cmd
	for _, src := range allSources() {
		cmds = append(cmds, fetchCmd(src))
	}
	return tea.Batch(append(cmds, tickCmd())...)
}

// rebuild maps the live sources to fixtures again and keeps the agent cursor on its agent.
func (m *model) rebuild() {
	if m.live == nil {
		return
	}
	m.fx = buildFixtures(liveInput{Now: m.now, LS: m.live})
	m.sync = m.fx.Sync
	if m.sync == stLoading {
		return
	}
	m.keepMR()
	for i, a := range m.fx.Agents {
		if a.Session != "" && a.Session == m.agentKey {
			m.agentSel = i
			return
		}
	}
	m.selectAgent(m.agentSel)
}

// keepMR keeps the MR cursor on its MR when the list changes.
func (m *model) keepMR() {
	vis := visibleMRs(m.fx, m.mode)
	for i, r := range vis {
		if r.key() == m.mrKey {
			m.mrSel = i
			return
		}
	}
	m.selectMR(m.mrSel)
}

func (m *model) selectMR(i int) {
	vis := visibleMRs(m.fx, m.mode)
	i = max(0, min(len(vis)-1, i))
	m.mrSel = i
	key := ""
	if i < len(vis) {
		key = vis[i].key()
	}
	if key != m.mrKey {
		m.mrKey, m.mrFile, m.mrItem, m.mrToggle = key, 0, 0, false
	}
}

// wantDetail starts the diff fetch of the selected MR when none is fresh or running.
func (m *model) wantDetail() tea.Cmd {
	if m.live == nil || m.page != 2 {
		return nil
	}
	vis := visibleMRs(m.fx, m.mode)
	if m.mrSel >= len(vis) {
		return nil
	}
	r := vis[m.mrSel]
	e := m.live.details[r.key()]
	if e != nil && (e.busy || m.now.Sub(e.at) < detailTTL) {
		return nil
	}
	if e == nil {
		e = &detailEntry{}
		m.live.details[r.key()] = e
	}
	e.busy = true
	return detailCmd(r)
}

func (m model) mrsView() mrsView {
	return mrsView{lk: m.look(), fx: m.fx, now: m.now, sync: m.sync, mode: m.mode, sel: m.mrSel, pane: m.mrPane,
		file: m.mrFile, item: m.mrItem, toggle: m.mrToggle}
}

func (m model) projectsView() projectsView {
	return projectsView{lk: m.look(), fx: m.fx, now: m.now, sync: m.sync, mode: m.mode, sel: m.prjSel, pane: m.prjPane, zoom: m.prjZoom}
}

func (m model) usageView() usageViewState {
	return usageViewState{lk: m.look(), fx: m.fx, sync: m.sync, sel: m.useSel, pane: m.usePane, fortnight: m.useFortnight, now: m.now}
}

func (m *model) selectAgent(i int) {
	m.agentSel = max(0, min(len(m.fx.Agents)-1, i))
	if m.agentSel < len(m.fx.Agents) {
		m.agentKey = m.fx.Agents[m.agentSel].Session
	}
}

func (m model) agentsView() agentsView {
	return agentsView{lk: m.look(), fx: m.fx, now: m.now, sync: m.sync, pane: m.agentPane, sel: m.agentSel,
		input: m.input, typing: m.typing, accents: agentAccounts(m.fx)}
}

func (m model) look() *look { return newLook(themes[m.th], m.accent) }

func (m *model) setSync(s string) {
	m.live = nil
	m.sync = s
	m.fx = mockFixtures()
	if s == stEmpty {
		m.fx = emptyFixtures()
		m.sync = stFresh
	}
	m.sel = [homePanes]int{}
	m.agentSel, m.agentKey, m.input, m.typing = 0, "", "", false
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
	if m.page == 1 {
		m.selectAgent(m.agentSel + d)
		m.input, m.typing = "", false
		return
	}
	if m.page == 2 {
		mv := m.mrsView()
		switch m.mrPane {
		case mpList:
			m.selectMR(m.mrSel + d)
		case mpDiff:
			if det := mv.detail(); det != nil {
				m.mrFile = max(0, min(len(det.Files)-1, mv.fileIdx(det)+d))
			}
		case mpCard:
			m.mrItem = max(0, min(mv.items()-1, mv.itemIdx(mv.items())+d))
			m.followItem()
		}
		return
	}
	if m.page == 3 {
		m.prjSel = max(0, min(len(m.projectsView().visible())-1, m.prjSel+d))
		return
	}
	if m.page == 4 {
		if m.usePane == upAccounts {
			m.useSel = max(0, min(len(m.fx.Gauges)-1, m.useSel+d))
		}
		return
	}
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
	switch m.page {
	case 0:
		m.focus = (m.focus + d + homePanes) % homePanes
	case 1:
		m.agentPane = (m.agentPane + d + agentPanes) % agentPanes
	case 2:
		m.mrPane = (m.mrPane + d + mrPanes) % mrPanes
	case 3:
		m.prjPane = (m.prjPane + d + projPanes) % projPanes
	case 4:
		m.usePane = (m.usePane + d + usagePanes) % usagePanes
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
			m.keepMR()
		}
	case "mock":
		m.setSync(arg)
	case "attach":
		if a := m.curAgent(); a != nil {
			return attachCmd(*a)
		}
	case "relaunch":
		if a := m.curAgent(); a != nil && a.Interrupted {
			m.overlay = ovConfirm
		} else {
			m.flash = "relaunch needs a session paused on a limit"
		}
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nm, cmd := m.update(msg)
	mm := nm.(model)
	return mm, tea.Batch(cmd, mm.wantDetail())
}

func (m model) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyMsg:
		return m.key(msg)
	case tickMsg:
		if m.live == nil {
			return m, nil
		}
		m.now = time.Time(msg)
		m.ticks++
		m.rebuild()
		return m, tea.Batch(append(refreshCmds(m.ticks), tickCmd())...)
	case srcMsg:
		if m.live != nil {
			m.live.apply(msg)
			m.rebuild()
		}
	case detailMsg:
		if m.live != nil {
			m.live.applyDetail(msg)
			m.rebuild()
		}
	case actionMsg:
		m.flash = msg.text
	}
	return m, nil
}

func (m model) key(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()
	m.flash = ""
	if k == "ctrl+c" {
		return m, tea.Quit
	}
	if m.overlay == ovNone && m.page == 1 && m.typing {
		return m.typingKey(msg)
	}
	switch m.overlay {
	case ovConfirm:
		switch k {
		case "enter", "y":
			m.overlay = ovNone
			if a := m.curAgent(); a != nil {
				m.flash = "relaunching…"
				return m, relaunchCmd(*a)
			}
		case "esc", "n":
			m.overlay = ovNone
		}
		return m, nil
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
		m.keepMR()
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
	case "z":
		if m.page == 3 {
			m.prjZoom = !m.prjZoom
		}
	case "o":
		if m.page == 3 {
			m.flash = m.projectKey(k)
		}
	case "d":
		if m.page == 4 {
			m.useFortnight = !m.useFortnight
		}
	case "p", "l", "[", "]":
		if m.page == 4 {
			m.flash = k + ": editing the account policy arrives with the actions slice"
		}
	case "/":
		m.flash = "/ history and search arrive with a later slice"
	case "enter", "h", "r", "x", "R", "a", "e", "c", "g", "T", "M", "P", "s":
		if m.page == 1 {
			return m.agentKeyPress(k)
		}
		if m.page == 2 {
			return m.mrKeyPress(k)
		}
		if m.page == 3 {
			m.flash = m.projectKey(k)
			return m, nil
		}
		if m.page == 4 && k == "r" {
			m.flash = "refreshing usage…"
			return m, fetchCmd(srcUsage)
		}
		m.flash = k + ": actions are not wired yet"
	}
	return m, nil
}

// projectKey explains an action key on the Projects page: it needs a selected project, and
// the actions themselves arrive with the actions slice.
func (m model) projectKey(k string) string {
	c := m.projectsView().cur()
	if c == nil {
		return "no project selected"
	}
	for _, a := range cardActions {
		if a[0] == k {
			return k + ": " + a[1] + " on " + c.Name + " is not wired yet"
		}
	}
	return k + ": not wired"
}

func (m model) curAgent() *agentRow {
	if m.agentSel >= 0 && m.agentSel < len(m.fx.Agents) {
		return &m.fx.Agents[m.agentSel]
	}
	return nil
}

// followItem moves the diff to the file of the card the cursor is on.
func (m *model) followItem() {
	mv := m.mrsView()
	d := mv.detail()
	if d == nil {
		return
	}
	path := ""
	if f := mv.curFinding(d); f != nil {
		path = f.Path
	} else if t := mv.curThread(d); t != nil && mv.threadsMode() {
		path = t.Path
	}
	for i, f := range d.Files {
		if path != "" && f.Path == path {
			m.mrFile = i
		}
	}
}

// mrKeyPress handles the MRs page keys. The page is read only: T flips the card kind, the
// other action keys say what arrives later.
func (m model) mrKeyPress(k string) (tea.Model, tea.Cmd) {
	mv := m.mrsView()
	switch k {
	case "T":
		m.mrToggle, m.mrItem = !m.mrToggle, 0
		m.followItem()
	case "enter":
		m.mrPane = min(mpCard, m.mrPane+1)
	case "a", "x", "e", "c":
		m.flash = k + ": saving a decision on a finding arrives with the review slice"
	case "M", "P":
		m.flash = k + ": opening and posting go through the mr and review skills, not wired here"
	case "g", "h", "r", "R", "s":
		if mv.cur() == nil {
			m.flash = "no MR selected"
		} else {
			m.flash = k + ": not wired on the MRs page yet"
		}
	}
	return m, nil
}

// agentKeyPress handles the Agents page action keys.
func (m model) agentKeyPress(k string) (tea.Model, tea.Cmd) {
	a := m.curAgent()
	if a == nil {
		return m, nil
	}
	switch k {
	case "a", "x", "e", "c":
		if a.Group != gNeeds {
			m.flash = "this agent does not wait for an answer"
			return m, nil
		}
		m.input = map[string]string{"a": "yes", "x": "no", "e": "yes, but ", "c": ""}[k]
		m.typing = true
		m.agentPane = apTail
	case "g":
		return m, attachCmd(*a)
	case "R":
		if !a.Interrupted {
			m.flash = "R relaunches a session paused on a limit"
			return m, nil
		}
		m.overlay = ovConfirm
	case "enter":
		m.typing = true
		m.agentPane = apTail
	default:
		m.flash = k + ": not an Agents action"
	}
	return m, nil
}

// typingKey edits the answer line; enter sends it, esc leaves it.
func (m model) typingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.typing, m.input = false, ""
	case tea.KeyEnter:
		a := m.curAgent()
		text := strings.TrimSpace(m.input)
		if a == nil || text == "" {
			return m, nil
		}
		m.typing, m.input = false, ""
		return m, sendCmd(*a, text)
	case tea.KeyBackspace:
		if r := []rune(m.input); len(r) > 0 {
			m.input = string(r[:len(r)-1])
		}
	case tea.KeyRunes:
		m.input += string(msg.Runes)
	case tea.KeySpace:
		m.input += " "
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
	if m.page == 1 {
		if m.typing {
			return []keyHint{{"enter", "send"}, {"esc", "cancel"}}
		}
		return []keyHint{{"j k", "agent"}, {"a x e c", "answer"}, {"enter", "send"}, {"g", "jump"}, {"/", "history"}, {"tab", "pane"}}
	}
	if m.page == 2 {
		return []keyHint{{"j k", "move"}, {"tab", "pane"}, {"T", "threads / findings"}, {"a x e c", "decide"}, {"M", "open MR"}, {"P", "post"}, {"m", "work / perso / all"}}
	}
	if m.page == 3 {
		return []keyHint{{"j k", "project"}, {"tab", "pane"}, {"c r s e o", "action"}, {"z", "zoom"}, {"/", "filter"}, {"m", "work / perso / all"}}
	}
	if m.page == 4 {
		return []keyHint{{"j k", "account"}, {"tab", "pane"}, {"d", "7 / 14 days"}, {"r", "refresh"}, {"p l [ ]", "policy"}}
	}
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
	case m.overlay == ovConfirm:
		body = lk.confirmOverlay(m.w, bh, m.curAgent())
	case m.page == 0:
		body = m.homeView().body(m.w, bh)
	case m.page == 1:
		body = m.agentsView().body(m.w, bh)
	case m.page == 2:
		body = m.mrsView().body(m.w, bh)
	case m.page == 3:
		body = m.projectsView().body(m.w, bh)
	default:
		body = m.usageView().body(m.w, bh)
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
	mockFlag := flag.Bool("mock", false, "show invented data instead of the live sources")
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
	switch {
	case *stateFlag != "":
		m.setSync(*stateFlag)
	case !*mockFlag:
		m.goLive()
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "crow-deck:", err)
		os.Exit(1)
	}
}
