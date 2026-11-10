package mission

// Owner actions shared by the CLI and the HTTP API (item_147): stop a running mission,
// hide an inbox row, and find the open row of a mission an answer is for. Each writes an
// action event with actor owner and the surface the owner acted from.

import (
	"slices"
	"strings"
	"time"

	"crow-core/contract"
)

func checkSurface(s string) error {
	if !slices.Contains(surfaces, s) {
		return failf(ErrInvalid, "surface %q: want one of %v", s, surfaces)
	}
	return nil
}

// Stop stops the live runs of a running mission. The next tick sees each run end with
// cause stopped and moves the mission on (waiting, review, done or failed).
func (env Env) Stop(ref, surface string) (contract.Mission, error) {
	m, err := env.find(ref)
	if err != nil {
		return m, err
	}
	if err := checkSurface(surface); err != nil {
		return m, err
	}
	if m.State != "running" {
		return m, failf(ErrConflict, "%s is %s; only a running mission stops", m.ID, m.State)
	}
	var live []string
	for _, r := range m.Runs {
		if r.Ended == nil {
			live = append(live, r.Run)
		}
	}
	if len(live) == 0 {
		return m, failf(ErrConflict, "%s has no live run", m.ID)
	}
	if env.StopRun == nil {
		return m, failf(ErrConflict, "%s: no way to stop run %s (CROW_ENGINE unset)", m.ID, live[0])
	}
	for _, run := range live {
		if err := env.StopRun(run); err != nil {
			return m, err
		}
		ad := contract.ActionData{Action: "stop", Surface: surface, Run: &run}
		if err := env.append(m.ID, "owner", "action", env.Now().Format(time.RFC3339), ad); err != nil {
			return m, err
		}
	}
	return env.find(m.ID)
}

// Hide drops an open inbox row without answering it. The mission does not move.
func (env Env) Hide(ref, surface string) (contract.InboxItem, error) {
	it, err := env.findItem(ref, false)
	if err != nil {
		return it, err
	}
	if err := checkSurface(surface); err != nil {
		return it, err
	}
	if it.State != "open" {
		return it, failf(ErrConflict, "%s is %s", it.ID, it.State)
	}
	id := it.ID
	ad := contract.ActionData{Action: "hide", Surface: surface, Item: &id}
	if err := env.append(it.Mission, "owner", "action", env.Now().Format(time.RFC3339), ad); err != nil {
		return it, err
	}
	return env.findItem(id, true)
}

// MissionItem finds the inbox row of mission id an answer is for. item is a row id, or
// its part after the mission id (v2, v1/q1). An empty item takes the only open row.
func (env Env) MissionItem(id, item string) (contract.InboxItem, error) {
	m, err := env.find(id)
	if err != nil {
		return contract.InboxItem{}, err
	}
	if m.ID != id {
		return contract.InboxItem{}, failf(ErrNotFound, "no mission %q", id)
	}
	if item != "" {
		if !strings.HasPrefix(item, id+"/") {
			item = id + "/" + item
		}
		it, err := env.findItem(item, true)
		if err == nil && it.Mission != id {
			err = failf(ErrNotFound, "no inbox item %q in %s", item, id)
		}
		return it, err
	}
	open, err := env.Inbox(false)
	if err != nil {
		return contract.InboxItem{}, err
	}
	var ids []string
	var hit contract.InboxItem
	for _, it := range open {
		if it.Mission == id {
			ids, hit = append(ids, it.ID), it
		}
	}
	switch len(ids) {
	case 0:
		return hit, failf(ErrConflict, "%s has no open inbox item", id)
	case 1:
		return hit, nil
	}
	return hit, failf(ErrInvalid, "%s has %d open items; name one: %s", id, len(ids), strings.Join(ids, ", "))
}
