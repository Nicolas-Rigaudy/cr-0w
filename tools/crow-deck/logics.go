package main

// v and V on the Projects page open the Logics viewer. logics-manager view serves it and
// reuses a running server through its own registry (~/.cache/logics-manager/viewers.json),
// so the deck only starts it detached, with the real HOME and crow-open as the browser.

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// viewerLaunch is one logics-manager view start: the directory it runs in and its arguments.
type viewerLaunch struct {
	Label string
	Dir   string
	Args  []string
}

// projectViewer is the launch for one project's corpus, or a flash text when there is none.
func (m model) projectViewer() (viewerLaunch, string) {
	c := m.projectsView().cur()
	if c == nil {
		return viewerLaunch{}, "v: no project selected"
	}
	if !c.Logics || c.Path == "" {
		return viewerLaunch{}, "v: " + c.Name + " has no Logics corpus"
	}
	return viewerLaunch{Label: c.Name, Dir: c.Path, Args: []string{"view", "--open"}}, ""
}

// fleetViewer is the launch for the fleet home, which reads its own fleet roots.
func fleetViewer() viewerLaunch {
	return viewerLaunch{Label: "fleet", Dir: crowHome(), Args: []string{"view", "--fleet", "--open"}}
}

// viewerKey runs v (selected project) or V (fleet) on the Projects page.
func (m model) viewerKey(k string) (tea.Model, tea.Cmd) {
	l := fleetViewer()
	if k == "v" {
		var why string
		if l, why = m.projectViewer(); why != "" {
			m.flash = why
			return m, nil
		}
	}
	m.flash = "opening the Logics viewer for " + l.Label + "..."
	return m, viewerCmd(l)
}

// viewerWait is how long the deck watches a start before calling the server up.
var viewerWait = 2 * time.Second

// startViewer starts logics-manager detached and reports how it went within viewerWait:
// still running (a new server), exited clean (it reused one) or failed (its last line).
var startViewer = func(l viewerLaunch) (state, out string) {
	log := filepath.Join(crowHome(), "local", "state", "logics-viewer.log")
	f, err := os.OpenFile(log, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return "failed", err.Error()
	}
	defer f.Close()
	cmd := exec.Command("logics-manager", l.Args...)
	cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = l.Dir, viewerEnv(), f, f
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return "failed", err.Error()
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		b, _ := os.ReadFile(log)
		if err != nil {
			return "failed", lastLine(cmpOr(strings.TrimSpace(string(b)), err.Error()))
		}
		return "reused", ""
	case <-time.After(viewerWait):
		return "started", ""
	}
}

func viewerCmd(l viewerLaunch) tea.Cmd {
	return func() tea.Msg {
		state, out := startViewer(l)
		switch state {
		case "failed":
			return actionMsg{"Logics viewer for " + l.Label + " failed: " + out}
		case "reused":
			return actionMsg{"Logics viewer for " + l.Label + ": opened on the running server"}
		}
		return actionMsg{"Logics viewer for " + l.Label + ": server started, opening the browser"}
	}
}

// viewerEnv is the environment with the real HOME (the viewer registry and the browser
// profile live there) and crow-open as BROWSER, without the cdx sandbox variables.
func viewerEnv() []string {
	home := os.Getenv("HOME")
	if u, err := user.Current(); err == nil && u.HomeDir != "" {
		home = u.HomeDir
	}
	env := []string{}
	for _, kv := range os.Environ() {
		k, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(k, "CLAUDE") || strings.HasPrefix(k, "ANTHROPIC_") || strings.HasPrefix(k, "CDX_") ||
			strings.HasPrefix(k, "XDG_CONFIG_HOME") || strings.HasPrefix(k, "XDG_DATA_HOME") || strings.HasPrefix(k, "XDG_CACHE_HOME") ||
			k == "HOME" || k == "BROWSER" {
			continue
		}
		env = append(env, kv)
	}
	return append(env, "HOME="+home, "BROWSER="+crowBin("crow-open"))
}
