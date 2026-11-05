package main

import (
	"strings"
	"testing"
)

// Pause and stop from the palette ask first; only y runs crow-focus, any other key cancels.
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
		if m.pending != tc.act || !strings.Contains(m.flash, "y confirms") {
			t.Fatalf("%s: pending %q flash %q", tc.act, m.pending, m.flash)
		}
		c := send(m, keyMsg("n"))
		if c.pending != "" || !strings.Contains(c.flash, "cancelled") {
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
