package main

import (
	"strings"
	"testing"
)

// recordViewer swaps startViewer for a recorder that answers state.
func recordViewer(t *testing.T, state, out string) *[]viewerLaunch {
	old := startViewer
	calls := &[]viewerLaunch{}
	startViewer = func(l viewerLaunch) (string, string) {
		*calls = append(*calls, l)
		return state, out
	}
	t.Cleanup(func() { startViewer = old })
	return calls
}

func logicsCards(m model, on bool) model {
	for i := range m.fx.Cards {
		m.fx.Cards[i].Logics, m.fx.Cards[i].Path = on, "/p/"+m.fx.Cards[i].Name
	}
	return m
}

func TestViewerKeyOpensProjectCorpus(t *testing.T) {
	calls := recordViewer(t, "started", "")
	m := logicsCards(projModel(235, 62), true)
	name := m.projectsView().cur().Name
	nm, cmd := m.Update(keyMsg("v"))
	if !strings.Contains(nm.(model).flash, "opening the Logics viewer for "+name) || cmd == nil {
		t.Fatalf("flash %q cmd %v", nm.(model).flash, cmd)
	}
	msg := cmd()
	if len(*calls) != 1 {
		t.Fatalf("calls %d", len(*calls))
	}
	l := (*calls)[0]
	if l.Dir != "/p/"+name || strings.Join(l.Args, " ") != "view --open" {
		t.Errorf("launch %+v", l)
	}
	if !strings.Contains(msg.(actionMsg).text, "server started") {
		t.Errorf("msg %q", msg.(actionMsg).text)
	}
}

func TestViewerKeyFleetReusesServer(t *testing.T) {
	calls := recordViewer(t, "reused", "")
	m := projModel(235, 62)
	_, cmd := m.Update(keyMsg("V"))
	msg := cmd()
	l := (*calls)[0]
	if strings.Join(l.Args, " ") != "view --fleet --open" || l.Dir != crowHome() {
		t.Errorf("launch %+v", l)
	}
	if !strings.Contains(msg.(actionMsg).text, "running server") {
		t.Errorf("msg %q", msg.(actionMsg).text)
	}
}

func TestViewerKeyFailureFlashesLastLine(t *testing.T) {
	recordViewer(t, "failed", "no Logics corpus here")
	_, cmd := logicsCards(projModel(235, 62), true).Update(keyMsg("v"))
	if msg := cmd().(actionMsg).text; !strings.Contains(msg, "failed: no Logics corpus here") {
		t.Errorf("msg %q", msg)
	}
}

func TestViewerKeyWithoutCorpus(t *testing.T) {
	calls := recordViewer(t, "started", "")
	m := logicsCards(projModel(235, 62), false)
	name := m.projectsView().cur().Name
	nm, cmd := m.Update(keyMsg("v"))
	if cmd != nil || len(*calls) != 0 || nm.(model).flash != "v: "+name+" has no Logics corpus" {
		t.Fatalf("flash %q cmd %v calls %d", nm.(model).flash, cmd, len(*calls))
	}
	m.fx.Cards = nil
	nm, _ = m.Update(keyMsg("v"))
	if nm.(model).flash != "v: no project selected" {
		t.Errorf("empty flash %q", nm.(model).flash)
	}
}

func TestViewerKeysFreeOnOtherPages(t *testing.T) {
	calls := recordViewer(t, "started", "")
	for _, p := range []int{0, 1, 2, 4} {
		m := logicsCards(projModel(235, 62), true)
		m.page = p
		for _, k := range []string{"v", "V"} {
			send(m, keyMsg(k))
		}
	}
	if len(*calls) != 0 {
		t.Errorf("v/V ran on another page: %d calls", len(*calls))
	}
}

func TestViewerEnvUsesRealHome(t *testing.T) {
	t.Setenv("HOME", "/sandbox/home")
	t.Setenv("CLAUDE_CODE_OAUTH_TOKEN", "x")
	env := strings.Join(viewerEnv(), "\n")
	if strings.Contains(env, "/sandbox/home") || strings.Contains(env, "CLAUDE_CODE_OAUTH_TOKEN") {
		t.Errorf("sandbox leaked: %s", env)
	}
	if !strings.Contains(env, "BROWSER="+crowBin("crow-open")) {
		t.Errorf("no crow-open browser")
	}
}

func TestViewerKeysInFooterAndHelp(t *testing.T) {
	m := projModel(235, 62)
	if !strings.Contains(m.View(), "logics / fleet") {
		t.Errorf("footer lacks v V")
	}
	m.overlay = ovHelp
	if !strings.Contains(m.View(), "Logics viewer") {
		t.Errorf("help lacks v V")
	}
}
