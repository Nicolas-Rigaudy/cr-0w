package mission

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"crow-core/contract"
)

func queueOf(t *testing.T, env Env, goals ...string) []string {
	t.Helper()
	var ids []string
	for _, g := range goals {
		m, err := env.New(research(g))
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, m.ID)
	}
	return ids
}

func order(t *testing.T, env Env) []string {
	t.Helper()
	q, err := env.Queue()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, m := range q {
		ids = append(ids, m.ID)
	}
	return ids
}

// AC1 (item_272): placed missions come first by position, then the rest by age.
func TestQueueCompare(t *testing.T) {
	one, two := 1, 2
	ms := []contract.Mission{
		{ID: "d", Created: "2026-11-01T10:00:00Z"},
		{ID: "c", Created: "2026-11-01T09:00:00Z"},
		{ID: "b", Created: "2026-11-02T09:00:00Z", Position: &two},
		{ID: "a", Created: "2026-11-03T09:00:00Z", Position: &one},
	}
	slices.SortStableFunc(ms, contract.QueueCompare)
	var got []string
	for _, m := range ms {
		got = append(got, m.ID)
	}
	if strings.Join(got, " ") != "a b c d" {
		t.Errorf("order = %v", got)
	}
}

// Move, Step and next rewrite only the places that change; a new mission goes last.
func TestQueueMove(t *testing.T) {
	env, root := newEnv(t)
	ids := queueOf(t, env, "First push study.", "Second push study.", "Third push study.")
	if _, err := env.Move(ids[2], 1, "nest"); err != nil {
		t.Fatal(err)
	}
	want := []string{ids[2], ids[0], ids[1]}
	if got := order(t, env); !slices.Equal(got, want) {
		t.Fatalf("after next = %v", got)
	}
	if _, err := env.Step(ids[1], -1, "nest"); err != nil {
		t.Fatal(err)
	}
	want = []string{ids[2], ids[1], ids[0]}
	if got := order(t, env); !slices.Equal(got, want) {
		t.Fatalf("after up = %v", got)
	}
	// Down at the end stays; a place past the end is last.
	if _, err := env.Step(ids[0], 1, "nest"); err != nil {
		t.Fatal(err)
	}
	if _, err := env.Move(ids[2], 9, "nest"); err != nil {
		t.Fatal(err)
	}
	want = []string{ids[1], ids[0], ids[2]}
	if got := order(t, env); !slices.Equal(got, want) {
		t.Fatalf("after move to 9 = %v", got)
	}
	n := strings.Count(strings.Join(events(t, root, ids[0]), " "), "action")
	if n != 3 { // 2 (next), 3 (up), 2 (move to 9)
		t.Errorf("%s has %d queue actions", ids[0], n)
	}
	late := queueOf(t, env, "A later push study.")
	if got := order(t, env); got[len(got)-1] != late[0] {
		t.Errorf("new mission not last: %v", got)
	}
}

func TestQueueMoveRefuses(t *testing.T) {
	env, _ := newEnv(t)
	ids := queueOf(t, env, "First push study.")
	if _, err := env.Move("nope", 1, "nest"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown: %v", err)
	}
	if _, err := env.Move(ids[0][:12], 1, "nest"); !errors.Is(err, ErrNotFound) {
		t.Errorf("prefix: %v", err)
	}
	if _, err := env.Move(ids[0], 0, "nest"); !errors.Is(err, ErrInvalid) {
		t.Errorf("place 0: %v", err)
	}
	if _, err := env.Cancel(ids[0], "", "cli", func(string) {}); err != nil {
		t.Fatal(err)
	}
	if _, err := env.Move(ids[0], 1, "nest"); !errors.Is(err, ErrConflict) {
		t.Errorf("cancelled: %v", err)
	}
	if _, err := env.StartNow(ids[0], "nest"); !errors.Is(err, ErrConflict) {
		t.Errorf("start cancelled: %v", err)
	}
}

// AC1: the tick admits by position, then by age.
func TestTickFollowsQueueOrder(t *testing.T) {
	env, _ := newEnv(t)
	ids := queueOf(t, env, "First push study.", "Second push study.", "Third push study.")
	if _, err := env.Move(ids[2], 1, "nest"); err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	rep, err := env.Tick(r, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Started) != 1 || rep.Started[0].Mission != ids[2] {
		t.Fatalf("started = %+v", rep.Started)
	}
	if m := show(t, env, ids[2]); m.Position != nil {
		t.Errorf("running mission keeps position %d", *m.Position)
	}
	if rep.Held[0].Mission != ids[0] || rep.Held[1].Mission != ids[1] {
		t.Errorf("held = %+v", rep.Held)
	}
}

// AC3: start now runs a mission past the cap, recorded as over cap.
func TestTickStartNowOverCap(t *testing.T) {
	env, _ := newEnv(t)
	ids := queueOf(t, env, "First push study.", "Second push study.", "Third push study.")
	r := &fakeRunner{ends: map[string]RunState{}}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	m, err := env.StartNow(ids[2], "nest")
	if err != nil {
		t.Fatal(err)
	}
	if !m.OverCap || m.Position == nil || *m.Position != 1 {
		t.Fatalf("after start now = %+v", m)
	}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Started) != 1 || rep.Started[0].Mission != ids[2] || !strings.Contains(rep.Started[0].Detail, "over cap") {
		t.Fatalf("report = %+v", rep)
	}
	if !r.dispatched[2].OverCap || r.dispatched[0].OverCap {
		t.Errorf("dispatches = %+v", r.dispatched)
	}
	m = show(t, env, ids[2])
	if m.State != "running" || m.OverCap {
		t.Errorf("mission = %+v", m)
	}
}

// Start now with a free slot starts in the normal way, not over cap.
func TestTickStartNowUnderCap(t *testing.T) {
	env, _ := newEnv(t)
	ids := queueOf(t, env, "First push study.", "Second push study.")
	if _, err := env.StartNow(ids[1], "nest"); err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	rep, err := env.Tick(r, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Started) != 1 || rep.Started[0].Mission != ids[1] || r.dispatched[0].OverCap {
		t.Fatalf("report = %+v, dispatched %+v", rep, r.dispatched)
	}
}
