package mission

// Owner queue actions (item_272): move a queued mission to a place in the queue, and
// start one now past the run cap. The queue is every mission in preflight or queued,
// placed missions first by position, then the rest oldest first (contract.QueueCompare).

import (
	"slices"
	"time"

	"crow-core/contract"
)

func inQueue(m contract.Mission) bool { return m.State == "queued" || m.State == "preflight" }

// Queue lists the missions waiting for a run, in the order the tick admits them.
func (env Env) Queue() ([]contract.Mission, error) {
	ms, err := env.List()
	if err != nil {
		return nil, err
	}
	q := slices.DeleteFunc(ms, func(m contract.Mission) bool { return !inQueue(m) })
	slices.SortStableFunc(q, contract.QueueCompare)
	return q, nil
}

// Move puts mission id at place to of the queue (1 = next; past the end = last). Every
// mission whose place changes gets a queue action, so the order holds as written.
func (env Env) Move(id string, to int, surface string) ([]contract.Mission, error) {
	if err := checkSurface(surface); err != nil {
		return nil, err
	}
	q, err := env.Queue()
	if err != nil {
		return nil, err
	}
	at := slices.IndexFunc(q, func(m contract.Mission) bool { return m.ID == id })
	if at < 0 {
		if _, err := env.exact(id); err != nil {
			return nil, err
		}
		return nil, failf(ErrConflict, "%s is not queued", id)
	}
	if to < 1 {
		return nil, failf(ErrInvalid, "place %d: want 1 or more", to)
	}
	to = min(to, len(q))
	m := q[at]
	q = slices.Insert(slices.Delete(q, at, at+1), to-1, m)
	ts := env.Now().Format(time.RFC3339)
	for i, m := range q {
		p := i + 1
		if m.Position != nil && *m.Position == p {
			continue
		}
		ad := contract.ActionData{Action: "queue", Surface: surface, Position: &p}
		if err := env.append(m.ID, "owner", "action", ts, ad); err != nil {
			return nil, err
		}
	}
	return env.Queue()
}

// Step moves mission id one place up (dir -1) or down (dir +1). At either end it stays.
func (env Env) Step(id string, dir int, surface string) ([]contract.Mission, error) {
	q, err := env.Queue()
	if err != nil {
		return nil, err
	}
	at := slices.IndexFunc(q, func(m contract.Mission) bool { return m.ID == id })
	if at < 0 {
		return env.Move(id, 1, surface) // reports why
	}
	return env.Move(id, max(1, at+1+dir), surface)
}

// StartNow marks a queued mission to start on the next tick, past the run cap when the
// cap is reached. It also moves the mission to the head of the queue.
func (env Env) StartNow(id, surface string) (contract.Mission, error) {
	m, err := env.exact(id)
	if err != nil {
		return m, err
	}
	if err := checkSurface(surface); err != nil {
		return m, err
	}
	if m.State != "queued" {
		return m, failf(ErrConflict, "%s is %s; only a queued mission starts now", m.ID, m.State)
	}
	if !m.OverCap {
		ad := contract.ActionData{Action: "start_now", Surface: surface}
		if err := env.append(m.ID, "owner", "action", env.Now().Format(time.RFC3339), ad); err != nil {
			return m, err
		}
	}
	if _, err := env.Move(m.ID, 1, surface); err != nil {
		return m, err
	}
	return env.find(m.ID)
}

// exact finds the mission with exactly this id: queue actions never take a prefix.
func (env Env) exact(id string) (contract.Mission, error) {
	m, err := env.find(id)
	if err == nil && m.ID != id {
		err = failf(ErrNotFound, "no mission %q", id)
	}
	return m, err
}
