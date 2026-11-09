package mission

import (
	"strings"
	"testing"
	"time"

	"crow-core/contract"
)

// cycle logs one resume that parks on a question again and gets answered: the mission goes
// back to queued with one more run_start carrying resume_of.
func cycle(t *testing.T, env Env, id, run, prev string) {
	t.Helper()
	ts := env.Now().Format(time.RFC3339)
	parked, answered := "parked: questions", "answered"
	for _, e := range []struct {
		typ  string
		data any
	}{
		{"state", contract.StateData{From: "queued", To: "running"}},
		{"run_start", contract.RunStartData{Run: run, AccountClass: "personal", Version: 1, ResumeOf: &prev}},
		{"run_end", contract.RunEndData{Run: run, Cause: "exit"}},
		{"state", contract.StateData{From: "running", To: "waiting", Reason: &parked}},
		{"state", contract.StateData{From: "waiting", To: "queued", Reason: &answered}},
	} {
		if err := env.append(id, "core", e.typ, ts, e.data); err != nil {
			t.Fatal(err)
		}
	}
}

// AC1 (tick side): an answered mission resumes on the next tick by continuing its last run.
func TestTickResumesAnAnsweredMission(t *testing.T) {
	env, _, id, _, _ := answered(t)
	r := &seqRunner{}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Started) != 1 || len(r.dispatched) != 1 || r.dispatched[0].ContinueRun != "r1" {
		t.Fatalf("tick = %+v, dispatched = %+v", rep, r.dispatched)
	}
	if got := show(t, env, id); got.State != "running" || len(got.Runs) != 2 {
		t.Errorf("mission = %s with %d runs", got.State, len(got.Runs))
	}
}

// A run cut by the usage limit resumes cold: its session stays on the old account.
func TestTickResumesColdAfterUsageLimit(t *testing.T) {
	env, root := newEnv(t)
	id, _ := startOne(t, env)
	ts := env.Now().Format(time.RFC3339)
	reason, back := "parked: usage_limit", "reset"
	for _, e := range []struct {
		typ  string
		data any
	}{
		{"run_end", contract.RunEndData{Run: "r1", Cause: "usage_limit"}},
		{"state", contract.StateData{From: "running", To: "waiting", Reason: &reason}},
		{"state", contract.StateData{From: "waiting", To: "queued", Reason: &back}},
	} {
		if err := env.append(id, "core", e.typ, ts, e.data); err != nil {
			t.Fatal(err)
		}
	}
	r := &seqRunner{}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	if len(r.dispatched) != 1 || r.dispatched[0].ContinueRun != "" || r.dispatched[0].Resumes != "r1" {
		t.Fatalf("dispatched = %+v", r.dispatched)
	}
	if sd := lastRunStart(t, root, id); !sd.Cold {
		t.Errorf("run_start = %+v", sd)
	}
}

// AC3: after budget.resumes resumes in a row without a gate, the tick asks the owner
// through the inbox instead of starting a run.
func TestTickAsksAtTheResumeCap(t *testing.T) {
	env, root, id, m, _ := answered(t)
	if m.Budget.Resumes != 2 {
		t.Fatalf("research resumes = %d, want 2", m.Budget.Resumes)
	}
	cycle(t, env, id, "c1", "r1")
	cycle(t, env, id, "c2", "c1")
	r := &seqRunner{}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.dispatched) != 0 || len(rep.Started) != 0 || len(rep.Finished) != 1 || rep.Finished[0].State != "waiting" {
		t.Fatalf("tick = %+v, dispatched = %+v", rep, r.dispatched)
	}
	if got := show(t, env, id); got.State != "waiting" {
		t.Errorf("state = %s", got.State)
	}
	if e := lastEvent(t, root, id); e.Type != "question" || e.Actor != "core" {
		t.Errorf("last event = %s by %s", e.Type, e.Actor)
	}
	items, err := env.Inbox(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != id+"/v1/"+CapQuestion || !strings.Contains(items[0].Ask, "cap 2") || !items[0].Blocking {
		t.Fatalf("inbox = %+v", items)
	}

	// The owner answers continue: one more resume runs, then the next one asks again.
	res, err := env.InboxReply(items[0].ID, Reply{Verb: "answer", Words: []string{CapContinue}, Surface: "cli"})
	if err != nil || res.State != "queued" {
		t.Fatalf("reply = %+v, %v", res, err)
	}
	if rep, err = env.Tick(r, 2); err != nil || len(rep.Started) != 1 || len(r.dispatched) != 1 {
		t.Fatalf("tick after continue = %+v, %v", rep, err)
	}
	if s, err := env.resumeStreak(id); err != nil || s.resumes != 1 || s.capped(show(t, env, id)) {
		t.Errorf("streak = %+v, %v", s, err)
	}
}

// The owner answers stop at the cap: the next tick cancels the mission.
func TestTickCancelsWhenTheOwnerStopsAtTheCap(t *testing.T) {
	env, _, id, _, _ := answered(t)
	cycle(t, env, id, "c1", "r1")
	cycle(t, env, id, "c2", "c1")
	r := &seqRunner{}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	if _, err := env.InboxReply(id+"/v1/"+CapQuestion, Reply{Verb: "answer", Words: []string{"stop"}, Surface: "cli"}); err != nil {
		t.Fatal(err)
	}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.dispatched) != 0 || len(rep.Finished) != 1 || rep.Finished[0].State != "cancelled" {
		t.Fatalf("tick = %+v, dispatched = %+v", rep, r.dispatched)
	}
	if got := show(t, env, id); got.State != "cancelled" {
		t.Errorf("state = %s", got.State)
	}
}

// A gate resets the streak: a mission that reached a version starts counting again.
func TestResumeStreakResetsOnGate(t *testing.T) {
	env, _, id, m, _ := answered(t)
	cycle(t, env, id, "c1", "r1")
	if err := env.append(id, "core", "gate", env.Now().Format(time.RFC3339),
		contract.GateData{File: "report.v1.md", Verdict: "blocked"}); err != nil {
		t.Fatal(err)
	}
	cycle(t, env, id, "c2", "c1")
	s, err := env.resumeStreak(id)
	if err != nil || s.resumes != 1 || s.capped(m) {
		t.Errorf("streak = %+v, %v", s, err)
	}
}
