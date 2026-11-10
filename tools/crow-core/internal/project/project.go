// Package project folds mission events into the rows of the store: the mission record,
// its runs and its inbox rows. It is pure: the same events give the same rows, whether
// they are applied one by one or replayed by a rebuild.
//
// Inbox rows here carry what the events say. The ask text, artifacts and blocking flag
// come from the version files and questions.json, and are filled by the inbox item
// (item_134).
package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"

	"crow-core/contract"
)

// View is everything the store keeps for one mission.
type View struct {
	Mission contract.Mission
	Inbox   []contract.InboxItem
}

// Kinds resolves the defaults of a kind file for a created event that leaves them out.
// A nil Kinds applies no defaults.
type Kinds func(name string) (*contract.Kind, error)

// KindDirs reads <name>.yaml from overrides first, then from dir.
func KindDirs(dir, overrides string) Kinds {
	return func(name string) (*contract.Kind, error) {
		for _, d := range []string{overrides, dir} {
			if d == "" {
				continue
			}
			b, err := os.ReadFile(filepath.Join(d, name+".yaml"))
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, err
			}
			return contract.ParseKind(b)
		}
		return nil, nil
	}
}

var versionFile = regexp.MustCompile(`^(plan|report|mockup|answers)\.v([0-9]+)\.md$`)

func parseVersionFile(f string) (string, int, bool) {
	m := versionFile.FindStringSubmatch(f)
	if m == nil {
		return "", 0, false
	}
	n, _ := strconv.Atoi(m[2])
	return m[1], n, true
}

// Gate shapes of req_037: the inbox type and verbs of each version file.
func gateRow(source string, straight bool) (typ string, verbs []string) {
	switch source {
	case "plan":
		return "approval", []string{"go", "dig", "no"}
	case "mockup":
		return "choice", []string{"pick", "change", "approve", "abandon"}
	default: // report
		if straight {
			return "fyi", []string{"dig"}
		}
		return "review", []string{"ok", "msg", "abandon"}
	}
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// Apply folds one event into v. The first event of a mission must be created.
func Apply(v *View, e contract.Event, kinds Kinds) error {
	d, err := contract.DecodeEventData(e)
	if err != nil {
		return err
	}
	m := &v.Mission
	if m.ID == "" && e.Type != "created" {
		return fmt.Errorf("%s: first event is %s, want created", e.Mission, e.Type)
	}
	if m.ID != "" && m.ID != e.Mission {
		return fmt.Errorf("event of %s applied to %s", e.Mission, m.ID)
	}
	switch d := d.(type) {
	case *contract.CreatedData:
		if m.ID != "" {
			return fmt.Errorf("%s: created twice", e.Mission)
		}
		if err := created(v, e, d, kinds); err != nil {
			return err
		}
	case *contract.StateData:
		if d.From != m.State {
			return fmt.Errorf("%s: state %s -> %s, but the mission is %s", e.Mission, d.From, d.To, m.State)
		}
		if !contract.CanMove(d.From, d.To) {
			return fmt.Errorf("%s: illegal move %s -> %s (from %s: %v)", e.Mission, d.From, d.To, d.From, contract.Moves[d.From])
		}
		m.State = d.To
	case *contract.RunStartData:
		for _, r := range m.Runs {
			if r.Run == d.Run {
				return fmt.Errorf("%s: run %s started twice", e.Mission, d.Run)
			}
		}
		m.Runs = append(m.Runs, contract.Run{
			Run: d.Run, Mission: m.ID, Account: deref(d.Account),
			Version: d.Version, Cold: d.Cold, Started: e.TS,
		})
	case *contract.RunEndData:
		r := findRun(m, d.Run)
		if r == nil {
			return fmt.Errorf("%s: run %s ended but never started", e.Mission, d.Run)
		}
		ts, cause := e.TS, d.Cause
		r.Ended, r.Cause, r.Exit = &ts, &cause, d.Exit
	case *contract.QuestionData:
		bump(m, d.Version)
		for _, q := range d.IDs {
			v.upsert(contract.InboxItem{
				ID: fmt.Sprintf("%s/v%d/%s", m.ID, d.Version, q), Type: "choice", Source: "questions",
				Version: d.Version, File: d.File, Verbs: []string{"answer"},
			}, m, e.TS)
		}
	case *contract.GateData:
		source, n, ok := parseVersionFile(d.File)
		if !ok || source == "answers" {
			return fmt.Errorf("%s: gate on %q", e.Mission, d.File)
		}
		bump(m, n)
		typ, verbs := gateRow(source, len(m.Gates) == 0)
		if d.Verdict != "ready" {
			// needs-input or blocked: the owner answers in words, nothing validates yet.
			typ, verbs = d.Verdict, slices.DeleteFunc(verbs, func(vb string) bool {
				return slices.Contains([]string{"go", "ok", "approve", "pick"}, vb)
			})
		}
		v.upsert(contract.InboxItem{
			ID: fmt.Sprintf("%s/v%d", m.ID, n), Type: typ, Source: source,
			Version: n, File: d.File, Verbs: verbs,
		}, m, e.TS)
	case *contract.AnswerData:
		id := deref(d.Item)
		if id == "" {
			_, n, ok := parseVersionFile(d.File)
			if !ok {
				return fmt.Errorf("%s: answer in %q", e.Mission, d.File)
			}
			id = fmt.Sprintf("%s/v%d", m.ID, n)
		}
		if it := v.find(id); it != nil {
			ts := e.TS
			it.State, it.Answered = "answered", &ts
		}
	case *contract.ValidationData:
		// The record chain is checked before each publish (mission.Verify, item_108).
		if e.Type == "validation" {
			ts := e.TS
			m.Validation = &ts
		} else {
			m.Validation = nil
		}
	}
	m.Updated = e.TS
	return nil
}

func created(v *View, e contract.Event, d *contract.CreatedData, kinds Kinds) error {
	m := &v.Mission
	*m = contract.Mission{
		ID: e.Mission, Kind: d.Kind, Batch: d.Batch,
		Scope: deref(d.Scope), DataClass: deref(d.DataClass), Project: deref(d.Project),
		Goal: deref(d.Goal), EndState: deref(d.EndState), Autonomy: deref(d.Autonomy),
		Gates: []string{}, State: "preflight", Runs: []contract.Run{}, FileScope: []string{},
		Created: e.TS, Updated: e.TS,
	}
	// A created event written by core carries the resolved values. A shorter one (an old
	// or hand-written log) takes the kind file defaults.
	var k *contract.Kind
	if kinds != nil && (d.Gates == nil || d.Autonomy == nil || d.Budget == nil) {
		var err error
		if k, err = kinds(d.Kind); err != nil {
			return fmt.Errorf("%s: kind %s: %w", e.Mission, d.Kind, err)
		}
	}
	switch {
	case d.Gates != nil:
		m.Gates = append(m.Gates, *d.Gates...)
	case k != nil:
		m.Gates = append(m.Gates, k.Gates...)
	}
	if d.Autonomy == nil && k != nil {
		m.Autonomy = k.Autonomy
	}
	switch {
	case d.Budget != nil:
		m.Budget = *d.Budget
	case k != nil:
		m.Budget = contract.Budget{TimeMin: k.Caps.TimeMin, BudgetUSD: k.Caps.BudgetUSD, Resumes: k.Caps.Resumes}
	}
	if d.FileScope != nil {
		m.FileScope = append(m.FileScope, *d.FileScope...)
	}
	return nil
}

func findRun(m *contract.Mission, id string) *contract.Run {
	for i := range m.Runs {
		if m.Runs[i].Run == id {
			return &m.Runs[i]
		}
	}
	return nil
}

func bump(m *contract.Mission, n int) {
	if n > m.Version {
		m.Version = n
	}
}

func (v *View) find(id string) *contract.InboxItem {
	for i := range v.Inbox {
		if v.Inbox[i].ID == id {
			return &v.Inbox[i]
		}
	}
	return nil
}

// upsert adds a row, or reopens an existing one with the same id.
func (v *View) upsert(it contract.InboxItem, m *contract.Mission, ts string) {
	it.Mission, it.Kind, it.Project = m.ID, m.Kind, m.Project
	it.Artifacts, it.State, it.Created = []string{}, "open", ts
	if old := v.find(it.ID); old != nil {
		*old = it
		return
	}
	v.Inbox = append(v.Inbox, it)
}
