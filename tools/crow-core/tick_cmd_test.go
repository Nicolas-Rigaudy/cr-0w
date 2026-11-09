package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"crow-core/internal/mission"
)

// fakeCrowRuns writes a crow-runs stand-in that logs its args and prints body.
func fakeCrowRuns(t *testing.T, body string) (crowRuns, string) {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "args")
	bin := filepath.Join(dir, "crow-runs")
	script := "#!/bin/sh\necho \"$@\" >> " + log + "\n" + body + "\n"
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return crowRuns{bin: bin}, log
}

func TestCrowRunsDispatchFlags(t *testing.T) {
	c, log := fakeCrowRuns(t, `echo '{"ok": true, "run_id": "ab12", "account": "main"}'`)
	s, err := c.Dispatch(mission.Dispatch{
		Mission: "20261109-research-x", Project: "demo", Label: "x", PromptFile: "/m/brief.md",
		Kind: "review", Model: "opus", Power: "medium", Permission: "review", BudgetUSD: 5, TimeoutMin: 60,
	})
	if err != nil || s != (mission.Started{RunID: "ab12", Account: "main"}) {
		t.Fatalf("started = %+v, %v", s, err)
	}
	b, _ := os.ReadFile(log)
	want := "dispatch demo x --prompt-file /m/brief.md --mission 20261109-research-x --kind review --json " +
		"--model opus --power medium --permission review --budget 5 --timeout-min 60"
	if got := strings.TrimSpace(string(b)); got != want {
		t.Errorf("args = %s\nwant   %s", got, want)
	}
}

func TestCrowRunsDispatchRefusals(t *testing.T) {
	c, _ := fakeCrowRuns(t, `echo '{"ok": false, "error": {"code": "too_many_runs", "message": "2 runs"}}'; exit 1`)
	_, err := c.Dispatch(mission.Dispatch{Project: "demo", Label: "x"})
	var de *mission.DispatchError
	if !errors.As(err, &de) || de.Code != "too_many_runs" {
		t.Errorf("err = %v", err)
	}
	c, _ = fakeCrowRuns(t, `echo 'crow: no allowed account' >&2; exit 1`)
	_, err = c.Dispatch(mission.Dispatch{Project: "demo", Label: "x"})
	if !errors.As(err, &de) || de.Code != "dispatch_failed" || !strings.Contains(de.Message, "no allowed account") {
		t.Errorf("err = %v", err)
	}
}

func TestCrowRunsState(t *testing.T) {
	c, log := fakeCrowRuns(t, `echo '{"run_id": "ab12", "ended": true, "state": "limit", "exit_code": 1, `+
		`"policy": {"choice": "move", "reason": "reset in 3h", "account": "main"}}'`)
	st, err := c.State("ab12")
	if err != nil || !st.Ended || st.State != "limit" || *st.ExitCode != 1 || st.Policy.Choice != "move" {
		t.Fatalf("state = %+v, %v", st, err)
	}
	if b, _ := os.ReadFile(log); strings.TrimSpace(string(b)) != "state ab12" {
		t.Errorf("args = %s", b)
	}
	c, _ = fakeCrowRuns(t, `echo 'boom' >&2; exit 3`)
	if _, err := c.State("x"); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Errorf("err = %v", err)
	}
}

func TestMaxRuns(t *testing.T) {
	env := map[string]string{}
	get := func(k string) string { return env[k] }
	if maxRuns(get) != 2 {
		t.Error("default cap")
	}
	env["CROW_MAX_RUNS"] = "1"
	if maxRuns(get) != 1 {
		t.Error("env cap")
	}
}
