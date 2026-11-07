// crow-deck: the chief of staff's terminal control room (skeleton: chrome, tabs, themes,
// config and a mock Home). Run `crow deck`. Spec: vault CR-0W/Research/crow-deck-design.md.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	relSel  int // Relaunch picker: index into fx.Rel.Accounts
	flash   string
	act     *deckAction // the action the confirm overlay shows
	cfgPath string

	live  *liveState // nil in mock mode
	now   time.Time
	ticks int

	agentSel   int    // position in the lane display order
	agentKey   string // key of the selected lane, keeps the cursor on it across refreshes
	agentPane  int
	agentShift int  // < > (or [ ]): hours the lane window moved back
	agentFull  bool // w: the whole day
	agentZoom  bool // d: Detail full screen
	input      string
	typing     bool

	mrSel, mrPane, mrFile, mrItem int
	mrScroll                      int    // first code row shown once the owner scrolled
	mrManual                      bool   // the owner scrolled the code: it no longer follows the card
	mrKey                         string // key of the selected MR, keeps the cursor on it across refreshes
	mrToggle                      bool
	mrTyping                      string // "e" edit or "c" chat on the current finding, "" when not typing
	mrInput                       string
	posts                         []postAction // pending push and post actions of the selected MR, one overlay each
	postI                         int
	posting                       bool // the shown action runs; keys wait for its answer

	prjSel, prjPane int
	prjZoom         bool

	usePane, useSel int
	usePeriod       int            // p cycles today, week, month
	useLed, usePrj  int            // Ledger and By project rows
	useFold         map[string]int // ledger groups folded by hand
	useFlat         bool           // g: ledger as one list
	useFilter       string         // ledger limited to one project
	useTok          bool           // c: daily chart in tokens
	useShift        int            // h l: runway pan, hours

	hist       bool // Agents shows the session history instead of the live view
	histF      histFilter
	histTyping bool
	histSel    int
	histPane   int
}

func newModel(cfg config) model {
	th := max(0, themeIndex(cfg.Theme))
	m := model{th: th, accent: cfg.Accent, mode: cfg.Mode, sync: stFresh, fx: mockFixtures(), focus: pNext, w: 160, h: 45, now: mockNow(), histF: newHistFilter(), mrPane: mpDiff}
	m.selectDefaultLane()
	return m
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
	m.agentSel, m.agentKey = 0, ""
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
	for i, li := range laneOrder(m.fx.Lanes) {
		if k := m.fx.Lanes[li].Key; k != "" && k == m.agentKey {
			m.agentSel = i
			return
		}
	}
	if m.agentKey == "" {
		m.selectDefaultLane()
		return
	}
	m.selectAgent(m.agentSel)
}

// selectDefaultLane puts the cursor on the first lane that waits for an answer.
func (m *model) selectDefaultLane() {
	for i, li := range laneOrder(m.fx.Lanes) {
		if m.fx.Lanes[li].State == lsNeeds {
			m.selectAgent(i)
			return
		}
	}
	m.selectAgent(0)
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
		m.mrKey, m.mrFile, m.mrItem, m.mrToggle, m.mrManual = key, 0, 0, false, false
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
	ttl := detailTTL
	if e != nil && e.d != nil && ((e.d.Review != nil && e.d.Review.State == "running") || (e.d.Apply != nil && (e.d.Apply.State == "running" || e.d.Apply.State == "starting"))) {
		ttl = reviewTTL
	}
	if e != nil && (e.busy || m.now.Sub(e.at) < ttl) {
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
	scroll := -1
	if m.mrManual {
		scroll = m.mrScroll
	}
	return mrsView{lk: m.look(), fx: m.fx, now: m.now, sync: m.sync, mode: m.mode, sel: m.mrSel, pane: m.mrPane,
		file: m.mrFile, item: m.mrItem, scroll: scroll, toggle: m.mrToggle, typing: m.mrTyping, input: m.mrInput}
}

func (m model) projectsView() projectsView {
	return projectsView{lk: m.look(), fx: m.fx, now: m.now, sync: m.sync, mode: m.mode, sel: m.prjSel, pane: m.prjPane, zoom: m.prjZoom}
}

func (m model) histView() histView {
	return histView{lk: m.look(), fx: m.fx, now: m.now, sync: m.sync, mode: m.mode, f: m.histF, typing: m.histTyping, sel: m.histSel, pane: m.histPane}
}

func (m model) usageView() usageViewState {
	return usageViewState{lk: m.look(), fx: m.fx, sync: m.sync, now: m.now, pane: m.usePane, sel: m.useSel, period: m.usePeriod,
		led: m.useLed, prj: m.usePrj, fold: m.useFold, flat: m.useFlat, filter: m.useFilter, tokens: m.useTok, shift: m.useShift}
}

func (m *model) selectAgent(i int) {
	o := laneOrder(m.fx.Lanes)
	m.agentSel = max(0, min(len(o)-1, i))
	if m.agentSel < len(o) {
		m.agentKey = m.fx.Lanes[o[m.agentSel]].Key
	}
}

func (m model) agentsView() agentsView {
	return agentsView{lk: m.look(), fx: m.fx, now: m.now, sync: m.sync, pane: m.agentPane, sel: m.agentSel,
		input: m.input, typing: m.typing, shift: m.agentShift, full: m.agentFull, zoom: m.agentZoom}
}

// wantLaneGit reads the selected lane's working tree when it is older than the TTL.
func (m *model) wantLaneGit() tea.Cmd {
	if m.live == nil || m.page != 1 || m.hist {
		return nil
	}
	l := m.agentsView().cur()
	if l == nil || l.Cwd == "" {
		return nil
	}
	g := m.live.git[l.Cwd]
	if g != nil && (g.Busy || m.now.Sub(g.At) < laneGitTTL) {
		return nil
	}
	if g == nil {
		g = &laneGit{}
		m.live.git[l.Cwd] = g
	}
	g.Busy = true
	cwd := l.Cwd
	return func() tea.Msg { return laneGitMsg{cwd: cwd, g: fetchLaneGit(cwd)} }
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
	m.agentShift, m.agentFull, m.agentZoom = 0, false, false
	m.selectDefaultLane()
	m.histSel = 0
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
	if m.page == 1 && m.hist {
		hv := m.histView()
		m.histSel = max(0, min(len(hv.rows())-1, hv.sel+d))
		return
	}
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
			m.scrollCode(d)
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
		m.usageMove(d)
		return
	}
	if m.page != 0 {
		return
	}
	hv := m.homeView()
	if !hv.selectable(m.focus) {
		return
	}
	m.sel[m.focus] = max(0, min(hv.items(m.focus)-1, m.sel[m.focus]+d))
}

func (m *model) moveFocus(d int) {
	switch m.page {
	case 0:
		m.focus = m.homeView().nextFocus(m.w, d)
	case 1:
		if m.hist {
			m.histPane = (m.histPane + d + histPanes) % histPanes
			return
		}
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
	case "history":
		m.page, m.hist = 1, true
	case "live":
		m.page, m.hist, m.histTyping = 1, false, false
	case "attach":
		if a := m.curAgent(); a != nil {
			return attachCmd(*a)
		}
	case "pause", "stop":
		a := m.curAgent()
		if a == nil || a.PID == 0 {
			m.flash = f[0] + " needs a running session"
			break
		}
		nm, cmd := m.proposeAgentCtl(*a, f[0])
		*m = nm
		return cmd
	case "relaunch":
		if a := m.curAgent(); a != nil && a.Interrupted {
			return m.openRelaunch()
		} else {
			m.flash = "relaunch needs a session paused on a limit"
		}
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nm, cmd := m.update(msg)
	mm := nm.(model)
	return mm, tea.Batch(cmd, mm.wantDetail(), mm.wantLaneGit())
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
	case laneGitMsg:
		if m.live != nil {
			g := msg.g
			m.live.git[msg.cwd] = &g
			m.rebuild()
		}
	case detailMsg:
		if m.live != nil {
			first := m.live.details[msg.key] == nil || m.live.details[msg.key].d == nil
			m.live.applyDetail(msg)
			m.rebuild()
			// The first load of the selected MR opens the diff at its first card.
			if r := m.mrsView().cur(); first && r != nil && r.key() == msg.key {
				m.followItem()
			}
		}
	case actionMsg:
		m.flash = msg.text
	case roadmapEditedMsg:
		m.flash = "roadmap of " + msg.name + " saved"
		if m.live != nil {
			return m, fetchCmd(srcRegistry)
		}
	case mailActMsg:
		m.flash = msg.text
		if m.live != nil {
			return m, fetchCmd(srcToday)
		}
	case postActionsMsg:
		r := m.mrsView().cur()
		switch {
		case r == nil || r.key() != msg.key:
		case msg.err != "":
			m.flash = msg.err
		case len(msg.acts) == 0:
			m.flash = "nothing to push or post on " + r.Ref + ifs(m.mrsView().mode2() == "self", ": A commits the decisions first") +
				ifs(m.mrsView().mode2() != "self", ": accept or edit a finding first")
		default:
			m.posts, m.postI, m.overlay = msg.acts, 0, ovPost
		}
	case postDoneMsg:
		m.posting = false
		m.flash = msg.text
		if m.live != nil {
			if e := m.live.details[msg.key]; e != nil {
				e.at = time.Time{}
			}
		}
		if msg.ok {
			m.nextPost()
		} else {
			m.overlay, m.posts = ovNone, nil
		}
		// The forge changed: read the MR again now, not at the next tick.
		return m, m.wantDetail()
	case reviewStartedMsg:
		m.flash = msg.text
		if d := m.fx.MRDetails[msg.key]; d != nil && msg.apply && msg.failed {
			d.Apply = &mrApply{State: "failed", Reason: msg.text}
		}
		if m.live != nil {
			if e := m.live.details[msg.key]; e != nil {
				e.at = time.Time{}
			}
		}
		// A run started: show its state now, not at the next tick.
		return m, m.wantDetail()
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
	if m.overlay == ovNone && m.page == 2 && m.mrTyping != "" {
		return m.mrTypingKey(msg)
	}
	if m.overlay == ovNone && m.page == 1 && m.hist {
		if nm, cmd, ok := m.histKey(msg); ok {
			return nm, cmd
		}
	}
	switch m.overlay {
	case ovConfirm:
		return m.relaunchKey(k)
	case ovAct:
		return m.actKey(k)
	case ovReview:
		switch k {
		case "enter":
			m.overlay = ovNone
			if r := m.mrsView().cur(); r != nil {
				if m.live == nil {
					m.flash = "mock deck: no review run started"
					return m, nil
				}
				m.flash = "starting a review of " + r.Ref + "…"
				return m, reviewCmd(*r, m.mrsView().detail())
			}
		case "esc", "n", "q":
			m.overlay = ovNone
		}
		return m, nil
	case ovApply:
		switch k {
		case "enter":
			m.overlay = ovNone
			if r := m.mrsView().cur(); r != nil {
				if m.live == nil {
					m.flash = "mock deck: no apply run started"
					return m, nil
				}
				m.flash = "starting the apply run on " + r.Ref + "…"
				// Shown at once on the card; the refetch replaces it with the run's state.
				if d := m.mrsView().detail(); d != nil {
					d.Apply = &mrApply{State: "starting", StartedAt: m.now.UTC().Format(time.RFC3339)}
				}
				return m, applyCmd(*r)
			}
		case "esc", "n", "q":
			m.overlay = ovNone
		}
		return m, nil
	case ovPost:
		if m.posting {
			return m, nil
		}
		r := m.mrsView().cur()
		switch k {
		case "enter":
			if r == nil || m.postI >= len(m.posts) {
				m.overlay = ovNone
				return m, nil
			}
			m.posting = true
			return m, actCmd(*r, m.posts[m.postI])
		case "s":
			m.nextPost()
		case "esc", "n", "q":
			m.overlay, m.posts = ovNone, nil
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
	if m.page == 4 {
		if mm, cmd, ok := m.usageKey(k); ok {
			return mm, cmd
		}
	}
	switch k {
	case "q":
		return m, tea.Quit
	case "1", "2", "3", "4", "5":
		m.page = int(k[0] - '1')
	case "&", "é", "\"", "'", "(":
		m.page = azertyTabs[k]
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
			return m.projectAction(k)
		}
		if m.page == 0 && m.focus == pToday {
			return m.todayKey(k)
		}
	case "d":
		if m.page == 1 {
			m.agentZoom = !m.agentZoom
		}
	case "w":
		if m.page == 1 {
			m.agentFull, m.agentShift = !m.agentFull, 0
		}
	case "[", "]", "<", ">":
		if m.page == 2 && (k == "[" || k == "]") {
			m.codeFile(map[string]int{"[": -1, "]": 1}[k])
			return m, nil
		}
		if m.page == 1 {
			if k == "[" || k == "<" {
				m.agentShift = min(m.agentShift+1, max(0, m.now.Hour()-3))
			} else {
				m.agentShift = max(0, m.agentShift-1)
			}
			m.agentFull = false
			return m, nil
		}
	case "n", "p", "J", "K":
		if m.page == 2 {
			switch k {
			case "n", "p":
				m.stepItem(map[string]int{"n": 1, "p": -1}[k])
				return m, m.wantDetail()
			default:
				m.codeHunk(map[string]int{"J": 1, "K": -1}[k])
			}
			return m, nil
		}
	case "/":
		if m.page == 1 {
			m.hist, m.histTyping, m.histPane = true, true, hpList
		} else {
			m.flash = "/ filters the Projects cards with the actions slice"
		}
	case "A":
		if m.page == 2 {
			return m.mrKeyPress(k)
		}
	case "enter", "h", "r", "x", "R", "a", "e", "c", "g", "T", "M", "P", "s":
		if m.page == 0 && m.focus == pToday {
			return m.todayKey(k)
		}
		if m.page == 1 {
			return m.agentKeyPress(k)
		}
		if m.page == 2 {
			return m.mrKeyPress(k)
		}
		if m.page == 3 {
			if k == "e" {
				return m, m.editRoadmap()
			}
			return m.projectAction(k)
		}
		if k == "h" || k == "r" || k == "x" {
			return m.homeAction(k)
		}
		m.flash = k + ": no action on this key here"
	}
	return m, nil
}

// azertyTabs maps the unshifted AZERTY digit row to tabs 1-5.
var azertyTabs = map[string]int{"&": 0, "é": 1, "\"": 2, "'": 3, "(": 4}

// editRoadmap opens the selected project's roadmap doc in $EDITOR (vi without one). The
// deck reloads the registry when the editor exits, so the card shows the edit.
func (m model) editRoadmap() tea.Cmd {
	c := m.projectsView().cur()
	if c == nil {
		return func() tea.Msg { return actionMsg{"no project selected"} }
	}
	if c.Doc == "" {
		return func() tea.Msg { return actionMsg{c.Name + " has no roadmap doc (release.roadmap in the registry)"} }
	}
	return tea.ExecProcess(editorCmd(c.Doc), func(err error) tea.Msg {
		if err != nil {
			return actionMsg{"editor failed: " + err.Error()}
		}
		return roadmapEditedMsg{c.Name}
	})
}

// editorCmd runs $EDITOR on path; $EDITOR may carry flags ("code -w").
func editorCmd(path string) *exec.Cmd {
	ed := strings.Fields(os.Getenv("EDITOR"))
	if len(ed) == 0 {
		ed = []string{"vi"}
	}
	return exec.Command(ed[0], append(ed[1:], path)...)
}

// pagerCmd runs $PAGER on path, less without one.
func pagerCmd(path string) *exec.Cmd {
	pg := strings.Fields(os.Getenv("PAGER"))
	if len(pg) == 0 {
		pg = []string{"less"}
	}
	return exec.Command(pg[0], append(pg[1:], path)...)
}

type roadmapEditedMsg struct{ name string }

// curAgent is the selected lane's agent. A lane that ended today gets a stand-in row
// so enter can resume it.
func (m model) curAgent() *agentRow {
	av := m.agentsView()
	l := av.cur()
	if l == nil {
		return nil
	}
	if a := av.agent(l); a != nil {
		return a
	}
	return &agentRow{Name: l.Project, Project: l.Project, Title: l.Title, Account: l.Account, Cwd: l.Cwd, Session: l.Key, Group: gDone}
}

// followItem moves the diff to the file of the card the cursor is on.
func (m *model) followItem() {
	m.mrManual = false
	mv := m.mrsView()
	d := mv.detail()
	if d == nil {
		return
	}
	path := ""
	if f := mv.curFinding(d); f != nil {
		path, _ = f.at()
	} else if t := mv.curThread(d); t != nil && mv.threadsMode() {
		path = t.Path
	}
	for i, f := range d.Files {
		if path != "" && f.Path == path {
			m.mrFile = i
		}
	}
}

// codeView is the code pane as the current window shows it: rows, the first one shown and
// the rows in sight.
func (m *model) codeView() (rows []codeRow, top, body int) {
	mv := m.mrsView()
	d := mv.detail()
	if d == nil || len(d.Files) == 0 {
		return nil, 0, 0
	}
	iw, body := mv.codeSlots(m.w, m.h-4)
	body = max(1, body)
	rows, at := mv.codeRows(d, iw)
	return rows, mv.codeTop(len(rows), at, body), body
}

// scrollCode moves the code by d rows; the code stops following the card.
func (m *model) scrollCode(d int) {
	rows, top, body := m.codeView()
	m.mrScroll, m.mrManual = max(0, min(len(rows)-body, top+d)), true
}

// codeHunk puts the next (d > 0) or previous hunk at the top of the code.
func (m *model) codeHunk(d int) {
	rows, top, body := m.codeView()
	to := -1
	for i, r := range rows {
		if !r.hunk {
			continue
		}
		if d > 0 && i > top && to < 0 {
			to = i
		}
		if d < 0 && i < top {
			to = i
		}
	}
	if to < 0 {
		m.flash = ifs(d > 0, "last hunk of the file") + ifs(d < 0, "first hunk of the file")
		return
	}
	m.mrScroll, m.mrManual = max(0, min(len(rows)-body, to)), true
}

// codeFile shows the next or previous changed file from its top.
func (m *model) codeFile(d int) {
	det := m.mrsView().detail()
	if det == nil || len(det.Files) == 0 {
		return
	}
	m.mrFile = max(0, min(len(det.Files)-1, m.mrsView().fileIdx(det)+d))
	m.mrScroll, m.mrManual = 0, false
}

// stepItem moves to the next or previous card; past the last card, n goes to the next MR.
func (m *model) stepItem(d int) {
	mv := m.mrsView()
	n := mv.items()
	i := mv.itemIdx(n) + d
	switch {
	case n > 0 && i >= 0 && i < n && !(d > 0 && mv.cur() != nil && mv.stage(*mv.cur()).done()):
		m.mrItem = i
		m.followItem()
	case d > 0:
		if m.mrSel+1 < len(mv.visible()) {
			m.selectMR(m.mrSel + 1)
			m.followItem()
		} else {
			m.flash = "last MR of the list"
		}
	default:
		m.flash = "first finding"
	}
}

// mrKeyPress handles the MRs page keys. Review, apply, post and handoff confirm in an
// overlay first; T flips the card kind.
func (m model) mrKeyPress(k string) (tea.Model, tea.Cmd) {
	mv := m.mrsView()
	switch k {
	case "T":
		m.mrToggle, m.mrItem = !m.mrToggle, 0
		m.followItem()
	case "enter":
		m.mrPane = min(mpCard, m.mrPane+1)
	case "a", "x", "e", "c":
		return m.decideKey(k)
	case "r":
		if mv.cur() == nil {
			m.flash = "no MR selected"
		} else if d := mv.detail(); d != nil && d.Review != nil && d.Review.State == "running" {
			m.flash = "a review of this MR is already running"
		} else {
			m.overlay = ovReview
		}
	case "A":
		m.applyKey()
	case "P":
		switch r := mv.cur(); {
		case r == nil:
			m.flash = "no MR selected"
		case m.live == nil:
			m.flash = "mock deck: nothing is pushed or posted"
		default:
			m.flash = "listing what " + r.Ref + " has to push or post…"
			return m, postActionsCmd(*r)
		}
	case "M":
		m.flash = k + ": opening goes through the mr skill, not wired here"
	case "h":
		if r := mv.cur(); r == nil {
			m.flash = "no MR selected"
		} else {
			return m.propose(mrHandoff(*r, m.findingsFile()))
		}
	case "g", "R", "s":
		if mv.cur() == nil {
			m.flash = "no MR selected"
		} else {
			m.flash = k + ": not wired on the MRs page yet"
		}
	}
	return m, nil
}

// nextPost shows the next pending action, or closes the overlay after the last one.
func (m *model) nextPost() {
	m.postI++
	if m.postI >= len(m.posts) {
		m.overlay, m.posts = ovNone, nil
	}
}

// applyKey opens the apply confirm when the review is done, nothing runs, and a decision waits.
func (m *model) applyKey() {
	mv := m.mrsView()
	r, d := mv.cur(), mv.detail()
	switch {
	case r == nil:
		m.flash = "no MR selected"
	case d == nil || d.Review == nil || d.Review.State != "done":
		m.flash = "A applies a finished review; r starts one"
	case d.Apply != nil && (d.Apply.State == "running" || d.Apply.State == "starting"):
		m.flash = "an apply run on this MR is already running"
	default:
		n, mode := 0, detailMode(r, d)
		for _, f := range d.Findings {
			if f.toApply(mode) {
				n++
			}
		}
		if n == 0 {
			m.flash = map[string]string{
				"self":    "nothing to apply: accept, edit or chat a finding first",
				"peer":    "nothing to apply: on a peer review only chat notes get an answer; P posts the rest",
				"threads": "nothing to apply: no accepted thread needs a code change; P replies and resolves",
			}[mode]
			return
		}
		m.overlay = ovApply
	}
}

// decideKey records a/x on the current finding at once; e and c open the input line first.
func (m model) decideKey(k string) (tea.Model, tea.Cmd) {
	mv := m.mrsView()
	if mv.threadsMode() {
		m.flash = k + ": r prepares an answer for each thread first" + ifs(mv.mode2() == "threads", "; T shows them")
		return m, nil
	}
	f := mv.curFinding(mv.detail())
	if f == nil {
		m.flash = "no finding selected; r starts a review"
		return m, nil
	}
	switch k {
	case "a":
		return m.decide("accepted", "")
	case "x":
		return m.decide("refused", "")
	case "e":
		m.mrTyping, m.mrInput = "e", cmpOr(f.Note, f.Comment)
	case "c":
		m.mrTyping, m.mrInput = "c", ""
	}
	m.mrPane = mpCard
	return m, nil
}

// decide sets the state of the current finding and saves it in the findings JSON (live).
func (m model) decide(state, note string) (tea.Model, tea.Cmd) {
	mv := m.mrsView()
	r, f := mv.cur(), mv.curFinding(mv.detail())
	if r == nil || f == nil {
		return m, nil
	}
	if m.live != nil && f.ID == "" {
		m.flash = "this finding has no id (old file); r runs the review again"
		return m, nil
	}
	f.State, f.Note = state, note
	if m.live == nil {
		m.flash = "mock deck: " + state + " kept in memory only"
		return m, nil
	}
	args := []string{"decide", r.Repo, strconv.Itoa(r.Number), f.ID, state}
	if note != "" {
		args = append(args, "--note", note)
	}
	// A decision is a line in the local findings file: it runs without the overlay, through the gate.
	return m, runAction(deckAction{Kind: kLocal, Scope: r.Scope, Label: "finding " + f.ID, Cmd: append([]string{"crow-review-run"}, args...)})
}

// mrTypingKey edits the edit or chat text of the current finding; enter saves it, esc drops it.
func (m model) mrTypingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mrTyping, m.mrInput = "", ""
	case tea.KeyEnter:
		text := strings.TrimSpace(m.mrInput)
		if text == "" {
			m.flash = "empty text: esc cancels"
			return m, nil
		}
		state := map[string]string{"e": "edited", "c": "chat"}[m.mrTyping]
		m.mrTyping, m.mrInput = "", ""
		return m.decide(state, text)
	case tea.KeyBackspace:
		if r := []rune(m.mrInput); len(r) > 0 {
			m.mrInput = string(r[:len(r)-1])
		}
	case tea.KeyRunes:
		m.mrInput += string(msg.Runes)
	case tea.KeySpace:
		m.mrInput += " "
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
	case "enter", "g":
		return m, attachCmd(*a)
	case "r", "R":
		if !a.Interrupted {
			m.flash = "r relaunches a session paused on a limit"
			return m, nil
		}
		return m, m.openRelaunch()
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
		return m.proposeSend(*a, text)
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
	if m.page == 1 && m.hist {
		if m.histTyping {
			return []keyHint{{"enter", "keep text"}, {"esc", "stop typing"}}
		}
		return []keyHint{{"/", "search"}, {"p a d M s", "filter"}, {"x", "clear"}, {"j k", "session"}, {"enter", "resume"}, {"o", "transcript"}, {"r", "retro"}, {"y", "copy id"}, {"h", "handoff"}, {"esc", "back to live"}}
	}
	if m.page == 1 {
		if m.typing {
			return []keyHint{{"enter", "send"}, {"esc", "cancel"}}
		}
		return []keyHint{{"enter", "open"}, {"j k", "lane"}, {"a x e c", "answer"}, {"r", "relaunch"}, {"< >", "hour range"}, {"w", "full day"}, {"d", "detail"}, {"tab", "pane"}, {"/", "history"}}
	}
	if m.page == 2 {
		if m.mrTyping != "" {
			return []keyHint{{"enter", "save"}, {"esc", "cancel"}}
		}
		return []keyHint{{"j k", "scroll"}, {"J K", "hunk"}, {"[ ]", "file"}, {"n p", "finding"}, {"tab", "pane"}, {"a x e c", "decide"}, {"A", "apply"}, {"P", "push / post"}, {"r", "review"}, {"T", "threads / findings"}, {"m", "work / perso / all"}}
	}
	if m.page == 3 {
		return []keyHint{{"j k", "project"}, {"tab", "pane"}, {"c r s e o h", "action"}, {"z", "zoom"}, {"m", "work / perso / all"}}
	}
	if m.page == 4 {
		return []keyHint{{"j k", "move"}, {"tab", "pane"}, {"p", "period: " + periodNames[m.usePeriod]}, {"enter", "fold / filter"}, {"g", "group"}, {"o", "open in history"}, {"h l", "pan"}, {"c", "cost / tokens"}, {"r", "refresh"}}
	}
	if m.page != 0 {
		return []keyHint{{"1-5 & é\" ' (", "tab"}, {"j k", "move"}}
	}
	if m.focus == pToday {
		return []keyHint{{"j k", "mail"}, {"o", "open in Outlook"}, {"r", "refresh"}, {"enter", "crow session"}, {"x", "done"}, {"s", "snooze"}, {"tab", "pane"}}
	}
	ctx := []keyHint{{"1-5 & é\" ' (", "tab"}}
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
	cs := chromeState{page: m.page, mode: m.mode, sync: m.sync, fx: m.fx, flash: m.flash, context: m.context(), hist: m.page == 1 && m.hist}
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
		body = lk.confirmOverlay(m.w, bh, m.curAgent(), m.fx.Rel, m.relSel)
	case m.overlay == ovAct:
		body = lk.actOverlay(m.w, bh, m.act)
	case m.overlay == ovReview:
		body = lk.reviewOverlay(m.w, bh, m.mrsView().cur(), m.mrsView().detail())
	case m.overlay == ovApply:
		body = lk.applyOverlay(m.w, bh, m.mrsView().cur(), m.mrsView().detail())
	case m.overlay == ovPost:
		body = lk.postOverlay(m.w, bh, m.posts, m.postI, m.posting)
	case m.page == 0:
		body = m.homeView().body(m.w, bh)
	case m.page == 1 && m.hist:
		body = m.histView().body(m.w, bh)
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

// findingsFile is the findings file of the selected MR, empty when it has none.
func (m model) findingsFile() string {
	if d := m.mrsView().detail(); d != nil {
		return d.FindingsFile
	}
	return ""
}
