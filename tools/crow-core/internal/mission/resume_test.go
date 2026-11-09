package mission

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"crow-core/contract"
)

// seqRunner refuses the dispatches listed in refuse (by call number, from 1) and starts
// the others as s1, s2...
type seqRunner struct {
	dispatched []Dispatch
	refuse     map[int]*DispatchError
	calls      int
}

func (f *seqRunner) Dispatch(d Dispatch) (Started, error) {
	f.calls++
	if e := f.refuse[f.calls]; e != nil {
		return Started{}, e
	}
	f.dispatched = append(f.dispatched, d)
	return Started{RunID: fmt.Sprintf("s%d", len(f.dispatched)), Account: "main"}, nil
}

func (f *seqRunner) State(string) (RunState, error) { return RunState{}, nil }

// answered parks a research mission on its version 1 question and answers it.
func answered(t *testing.T) (Env, string, string, contract.Mission, *contract.Kind) {
	t.Helper()
	env, root, id := parked(t, false)
	if _, err := env.InboxReply(id+"/v1", Reply{Verb: "answer", Words: []string{"skip it"}, Surface: "cli"}); err != nil {
		t.Fatal(err)
	}
	m := show(t, env, id)
	k, err := env.Kinds(m.Kind)
	if err != nil {
		t.Fatal(err)
	}
	return env, root, id, m, k
}

func lastRunStart(t *testing.T, root, id string) contract.RunStartData {
	t.Helper()
	e := lastEvent(t, root, id)
	if e.Type != "run_start" {
		t.Fatalf("last event = %s", e.Type)
	}
	d, err := contract.DecodeEventData(e)
	if err != nil {
		t.Fatal(err)
	}
	return *d.(*contract.RunStartData)
}

// The default resume continues the last run with a note that quotes the answer.
func TestResumeContinuesTheLastRun(t *testing.T) {
	env, root, id, m, k := answered(t)
	r := &seqRunner{}
	line, err := env.resume(r, m, k, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.dispatched) != 1 {
		t.Fatalf("dispatched = %+v", r.dispatched)
	}
	d := r.dispatched[0]
	if d.ContinueRun != "r1" || d.Resumes != "" || d.PromptFile != filepath.Join(root, id, ResumeNote) {
		t.Errorf("dispatch = %+v", d)
	}
	note, err := os.ReadFile(d.PromptFile)
	if err != nil || !strings.Contains(string(note), "skip it") || !strings.Contains(string(note), "Write version 2 next") {
		t.Errorf("note = %q, %v", note, err)
	}
	if b, err := os.ReadFile(filepath.Join(root, id, "brief.md")); err != nil || !strings.Contains(string(b), "skip it") {
		t.Errorf("brief.md misses the answer: %v", err)
	}
	sd := lastRunStart(t, root, id)
	if sd.Run != "s1" || sd.Cold || sd.Version != 2 || sd.ResumeOf == nil || *sd.ResumeOf != "r1" {
		t.Errorf("run_start = %+v", sd)
	}
	if got := show(t, env, id); got.State != "running" || len(got.Runs) != 2 {
		t.Errorf("mission = %s with %d runs", got.State, len(got.Runs))
	}
	if line.State != "running" || !strings.HasPrefix(line.Detail, "continue resume of r1") {
		t.Errorf("line = %+v", line)
	}
}

// AC1: a continuation cdx refuses (isolated workspace gone) resumes cold from brief.md,
// logs why, and the version chain goes on at version 2.
func TestResumeFallsBackColdWhenContinueIsRefused(t *testing.T) {
	for _, code := range []string{"continue_run_workspace_gone", "continue_run_settings_mismatch", "continue_run_unavailable"} {
		t.Run(code, func(t *testing.T) {
			env, root, id, m, k := answered(t)
			r := &seqRunner{refuse: map[int]*DispatchError{1: {Code: code, Message: "start a normal run"}}}
			line, err := env.resume(r, m, k, false)
			if err != nil {
				t.Fatal(err)
			}
			if r.calls != 2 || len(r.dispatched) != 1 {
				t.Fatalf("calls = %d, dispatched = %+v", r.calls, r.dispatched)
			}
			d := r.dispatched[0]
			if d.ContinueRun != "" || d.Resumes != "r1" || d.PromptFile != filepath.Join(root, id, "brief.md") {
				t.Errorf("dispatch = %+v", d)
			}
			sd := lastRunStart(t, root, id)
			if !sd.Cold || sd.Version != 2 || sd.ResumeOf == nil || *sd.ResumeOf != "r1" {
				t.Errorf("run_start = %+v", sd)
			}
			ev := events(t, root, id)
			if got := strings.Join(ev[len(ev)-3:], ","); got != "incident,state,run_start" {
				t.Errorf("events tail = %s", got)
			}
			if !strings.HasPrefix(line.Detail, "cold resume of r1") {
				t.Errorf("line = %+v", line)
			}
		})
	}
}

// AC3: a mission that moves account resumes cold without trying to continue.
func TestResumeColdWhenTheMissionMovesAccount(t *testing.T) {
	env, root, id, m, k := answered(t)
	r := &seqRunner{}
	if _, err := env.resume(r, m, k, true); err != nil {
		t.Fatal(err)
	}
	if r.calls != 1 || r.dispatched[0].ContinueRun != "" || r.dispatched[0].Resumes != "r1" {
		t.Errorf("dispatched = %+v", r.dispatched)
	}
	if sd := lastRunStart(t, root, id); !sd.Cold {
		t.Errorf("run_start = %+v", sd)
	}
	if _, err := os.Stat(filepath.Join(root, id, ResumeNote)); !os.IsNotExist(err) {
		t.Errorf("a cold resume wrote %s: %v", ResumeNote, err)
	}
}

// Any other refusal starts nothing and leaves the mission queued.
func TestResumeOtherRefusalStartsNothing(t *testing.T) {
	env, root, id, m, k := answered(t)
	before := len(events(t, root, id))
	r := &seqRunner{refuse: map[int]*DispatchError{1: {Code: "too_many_runs", Message: "2 runs already running"}}}
	line, err := env.resume(r, m, k, false)
	if err == nil || line != nil || r.calls != 1 {
		t.Fatalf("line = %+v, err = %v, calls = %d", line, err, r.calls)
	}
	if after := len(events(t, root, id)); after != before {
		t.Errorf("%d events appended", after-before)
	}
	if got := show(t, env, id); got.State != "queued" {
		t.Errorf("state = %s", got.State)
	}
}
