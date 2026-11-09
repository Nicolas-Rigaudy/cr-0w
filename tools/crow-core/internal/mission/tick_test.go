package mission

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"crow-core/contract"
)

// fakeRunner starts runs r1, r2... and reports the end states the test sets.
type fakeRunner struct {
	dispatched []Dispatch
	ends       map[string]RunState
	refuse     *DispatchError
}

func (f *fakeRunner) Dispatch(d Dispatch) (Started, error) {
	if f.refuse != nil {
		return Started{}, f.refuse
	}
	f.dispatched = append(f.dispatched, d)
	return Started{RunID: fmt.Sprintf("r%d", len(f.dispatched)), Account: "main"}, nil
}

func (f *fakeRunner) State(run string) (RunState, error) {
	if st, ok := f.ends[run]; ok {
		return st, nil
	}
	return RunState{RunID: run}, nil
}

func ended(state string, exit int) RunState {
	return RunState{Ended: true, State: state, ExitCode: &exit}
}

func writeVersion(t *testing.T, root, id, file, verdict string, version int) {
	t.Helper()
	md := fmt.Sprintf("# Report\n\n```crow\ncontract: 1\nmission: %s\nkind: research\nphase: report\nversion: %d\n"+
		"summary: one provider\nconfidence: 80\nverdict: %s\n```\n", id, version, verdict)
	if err := os.WriteFile(filepath.Join(root, id, file), []byte(md), 0o600); err != nil {
		t.Fatal(err)
	}
}

func show(t *testing.T, env Env, id string) contract.Mission {
	t.Helper()
	m, err := env.Show(id)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func events(t *testing.T, root, id string) []string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, id, "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var types []string
	for _, l := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var e contract.Event
		if err := contract.Decode([]byte(l), &e); err != nil {
			t.Fatal(err)
		}
		types = append(types, e.Type)
	}
	return types
}

// AC1: tick starts a research mission, a later tick finishes it straight through.
func TestTickStartsThenFinishesStraightThrough(t *testing.T) {
	env, root := newEnv(t)
	m, err := env.New(research("Compare push providers for the Android app."))
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Started) != 1 || rep.Started[0].Run != "r1" {
		t.Fatalf("report = %+v", rep)
	}
	d := r.dispatched[0]
	if d.Mission != m.ID || d.Project != "demo" || d.Kind != "review" || d.Permission != "review" ||
		d.Model != "opus" || d.Power != "medium" || d.BudgetUSD != 5 || d.TimeoutMin != 60 {
		t.Errorf("dispatch = %+v", d)
	}
	if d.PromptFile != filepath.Join(root, m.ID, "brief.md") {
		t.Errorf("prompt = %s", d.PromptFile)
	}
	if _, err := os.Stat(d.PromptFile); err != nil {
		t.Errorf("brief not written: %v", err)
	}
	got := show(t, env, m.ID)
	if got.State != "running" || len(got.Runs) != 1 || got.Runs[0].Run != "r1" || got.Runs[0].Version != 1 {
		t.Fatalf("after start = %+v", got)
	}

	// The run is live: the next tick leaves it.
	if rep, _ = env.Tick(r, 2); len(rep.Finished) != 0 || show(t, env, m.ID).State != "running" {
		t.Fatalf("live run finished: %+v", rep)
	}

	writeVersion(t, root, m.ID, "report.v1.md", "ready", 1)
	r.ends["r1"] = ended("done", 0)
	if rep, err = env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	got = show(t, env, m.ID)
	if got.State != "done" || got.Runs[0].Ended == nil || *got.Runs[0].Cause != "done" || got.Version != 1 {
		t.Fatalf("after end = %+v", got)
	}
	want := []string{"created", "state", "state", "run_start", "run_end", "gate", "state"}
	if e := events(t, root, m.ID); strings.Join(e, ",") != strings.Join(want, ",") {
		t.Errorf("events = %v, want %v", e, want)
	}
	if len(r.dispatched) != 1 {
		t.Errorf("a done mission was dispatched again: %d", len(r.dispatched))
	}
}

// AC1: a gated mission goes to review.
func TestTickGatedMissionGoesToReview(t *testing.T) {
	env, root := newEnv(t)
	spec := research("Compare push providers for the Android app.")
	spec.Gates = &[]string{"report"}
	m, err := env.New(spec)
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	writeVersion(t, root, m.ID, "report.v1.md", "ready", 1)
	r.ends["r1"] = ended("done", 0)
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	if got := show(t, env, m.ID); got.State != "review" {
		t.Errorf("state = %s", got.State)
	}
}

// AC2: the third mission stays queued while two run.
func TestTickCapsAdmission(t *testing.T) {
	env, _ := newEnv(t)
	var ids []string
	for _, g := range []string{"First push study.", "Second push study.", "Third push study."} {
		m, err := env.New(research(g))
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, m.ID)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Started) != 2 || len(rep.Held) != 1 || rep.Held[0].Mission != ids[2] {
		t.Fatalf("report = %+v", rep)
	}
	if !strings.Contains(rep.Held[0].Detail, "2 of 2 runs busy") {
		t.Errorf("held detail = %s", rep.Held[0].Detail)
	}
	if show(t, env, ids[2]).State != "queued" {
		t.Error("third mission left queued state")
	}
	// A refused dispatch leaves the mission queued with no event.
	r.ends["r1"] = ended("stopped", 1)
	r.refuse = &DispatchError{Code: "too_many_runs", Message: "2 runs already"}
	rep, err = env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if show(t, env, ids[2]).State != "queued" || len(rep.Held) != 1 || !strings.Contains(rep.Held[0].Detail, "too_many_runs") {
		t.Errorf("report = %+v", rep)
	}
}

// AC3: a run that leaves no version file fails the mission with an incident.
func TestTickNoVersionFileFails(t *testing.T) {
	env, root := newEnv(t)
	m, err := env.New(research("Compare push providers for the Android app."))
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	r.ends["r1"] = ended("failed", 2)
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	got := show(t, env, m.ID)
	if got.State != "failed" || *got.Runs[0].Cause != "error" {
		t.Fatalf("mission = %+v", got)
	}
	if len(rep.Finished) != 1 || !strings.Contains(rep.Finished[0].Detail, "no plan, report or mockup v1") {
		t.Errorf("report = %+v", rep)
	}
	b, _ := os.ReadFile(filepath.Join(root, m.ID, "events.jsonl"))
	if !strings.Contains(string(b), `"cause":"no_version_file"`) || !strings.Contains(string(b), "ended error (exit 2)") {
		t.Errorf("no incident naming the cause:\n%s", b)
	}
}

// A version file for another mission does not count.
func TestTickRejectsForeignCrowBlock(t *testing.T) {
	env, root := newEnv(t)
	m, err := env.New(research("Compare push providers for the Android app."))
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	writeVersion(t, root, m.ID, "report.v1.md", "ready", 1)
	b, _ := os.ReadFile(filepath.Join(root, m.ID, "report.v1.md"))
	os.WriteFile(filepath.Join(root, m.ID, "report.v1.md"), []byte(strings.Replace(string(b), "mission: "+m.ID, "mission: other", 1)), 0o600)
	r.ends["r1"] = ended("done", 0)
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if show(t, env, m.ID).State != "failed" || !strings.Contains(rep.Finished[0].Detail, "names mission other") {
		t.Errorf("report = %+v", rep)
	}
}

// A usage limit parks the mission for the resume item.
func TestTickUsageLimitParks(t *testing.T) {
	env, _ := newEnv(t)
	m, err := env.New(research("Compare push providers for the Android app."))
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	st := ended("limit", 1)
	st.Policy = &Policy{Choice: "move", Reason: "reset in 180 min"}
	r.ends["r1"] = st
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	got := show(t, env, m.ID)
	if got.State != "waiting" || *got.Runs[0].Cause != "usage_limit" {
		t.Fatalf("mission = %+v", got)
	}
	if !strings.Contains(rep.Finished[0].Detail, "policy move") {
		t.Errorf("detail = %s", rep.Finished[0].Detail)
	}
}

// Kinds past the live phase stay queued.
func TestTickHoldsLaterPhases(t *testing.T) {
	env, _ := newEnv(t)
	m, err := env.New(Spec{Kind: "build", Project: "demo", Goal: "Build the push client.", EndState: "Tests pass."})
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.dispatched) != 0 || len(rep.Held) != 1 || !strings.Contains(rep.Held[0].Detail, "phase 3") {
		t.Errorf("report = %+v", rep)
	}
	if show(t, env, m.ID).State != "queued" {
		t.Error("build mission moved")
	}
}

func TestDispatchForWritableKind(t *testing.T) {
	env, _ := newEnv(t)
	k, err := env.Kinds("plan")
	if err != nil {
		t.Fatal(err)
	}
	d := dispatchFor(contract.Mission{ID: "x", Project: "demo"}, k, "x", "/b.md")
	if d.Kind != "task" || d.Permission != "" {
		t.Errorf("dispatch = %+v", d)
	}
}
