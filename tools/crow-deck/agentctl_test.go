package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Pause and stop from the palette ask first; the overlay runs crow-focus once on y or enter, n cancels.
func TestPauseAndStopAskFirst(t *testing.T) {
	var calls []call
	old := runBin
	runBin = func(name string, args ...string) (string, error) {
		calls = append(calls, call{name, args})
		return "done", nil
	}
	t.Cleanup(func() { runBin = old })

	for _, tc := range []struct{ act, flag string }{{"pause", "--interrupt"}, {"stop", "--stop"}} {
		m := send(newModel(defaultConfig()), resize(120, 40), keyMsg("2"))
		for i := range m.fx.Agents {
			m.fx.Agents[i].PID = 4242
		}
		if a := m.curAgent(); a == nil || a.PID != 4242 {
			t.Fatalf("%s: no selected agent with a pid", tc.act)
		}
		calls = nil
		_ = m.exec(tc.act)
		if m.overlay != ovAct || m.act == nil || m.act.Label != tc.act || len(calls) != 0 {
			t.Fatalf("%s: overlay %d act %+v calls %d", tc.act, m.overlay, m.act, len(calls))
		}
		c := send(m, keyMsg("n"))
		if c.act != nil || c.overlay != ovNone || !strings.Contains(c.flash, "cancelled") {
			t.Fatalf("%s: n did not cancel: %q", tc.act, c.flash)
		}
		_, cmd := m.Update(keyMsg("y"))
		if cmd == nil {
			t.Fatalf("%s: y returned no command", tc.act)
		}
		runAll(cmd)
		if len(calls) != 1 || calls[0].name != "crow-focus" || strings.Join(calls[0].args, " ") != "--pid 4242 "+tc.flag {
			t.Fatalf("%s: calls %+v", tc.act, calls)
		}
	}
}

func TestFailedActionShowsTheLastLine(t *testing.T) {
	old := runBin
	defer func() { runBin = old }()
	runBin = func(name string, args ...string) (string, error) {
		return "crow: using 'main' (30% left)\ncrow: unknown project 'x'", errors.New("exit 1")
	}
	msg := actionCmd("mail session", "crow", "handoff", "x")().(actionMsg)
	if !strings.Contains(fmt.Sprint(msg), "unknown project") {
		t.Fatalf("flash hides the error: %v", msg)
	}
}
