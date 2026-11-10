package mission

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"crow-core/contract"
	"crow-core/internal/eventlog"
)

// VersionEntry is one version of a mission: the file the run wrote, its block summary,
// the questions it asked and the answer points the owner gave.
type VersionEntry struct {
	Version   int                 `json:"version"`
	File      *string             `json:"file"`
	Summary   string              `json:"summary"`
	Questions []contract.Question `json:"questions"`
	Answers   []string            `json:"answers"`
}

// MissionDetail is one mission with its versions and runs, for a client screen.
type MissionDetail struct {
	Mission  contract.Mission `json:"mission"`
	Versions []VersionEntry   `json:"versions"`
	Runs     []contract.Run   `json:"runs"`
}

// ErrNotFound is a mission or file a reader may not see. Callers answer 404.
var ErrNotFound = errors.New("not found")

// Detail returns one mission by exact id with every version up to its current one.
func (env Env) Detail(id string) (MissionDetail, error) {
	if eventlog.ValidMission(id) != nil {
		return MissionDetail{}, ErrNotFound
	}
	ms, err := env.List()
	if err != nil {
		return MissionDetail{}, err
	}
	i := slices.IndexFunc(ms, func(m contract.Mission) bool { return m.ID == id })
	if i < 0 {
		return MissionDetail{}, ErrNotFound
	}
	m := ms[i]
	dir := filepath.Join(env.Store.Root(), id)
	d := MissionDetail{Mission: m, Versions: []VersionEntry{}, Runs: m.Runs}
	if d.Runs == nil {
		d.Runs = []contract.Run{}
	}
	for n := 1; n <= m.Version; n++ {
		d.Versions = append(d.Versions, versionOf(dir, id, n))
	}
	return d, nil
}

// versionOf reads what version n left in the mission dir. Missing files leave the
// fields empty: a run may end a version with questions only.
func versionOf(dir, id string, n int) VersionEntry {
	v := VersionEntry{Version: n, Questions: []contract.Question{}, Answers: []string{}}
	if f, cb, _ := versionFile(dir, id, n); cb != nil {
		v.File, v.Summary = &f, cb.Summary
		if cb.Questions != nil {
			v.Questions = append(v.Questions, *cb.Questions...)
		}
	}
	if q := questionsOf(dir, id, n); q != nil {
		v.Questions = append(v.Questions, q.Questions...)
	}
	if pts, err := readPoints(filepath.Join(dir, AnswersFile(n))); err == nil && pts != nil {
		v.Answers = pts
	}
	return v
}

// Artifact resolves an outbox file of a mission for reading. It answers ErrNotFound for
// an unknown mission, a path outside outbox/ and a file over the kind's per-file cap.
func (env Env) Artifact(id, p string) (string, error) {
	if eventlog.ValidMission(id) != nil {
		return "", ErrNotFound
	}
	m, err := env.Show(id)
	if err != nil || m.ID != id {
		return "", ErrNotFound
	}
	dir := filepath.Join(env.Store.Root(), id)
	full, err := outboxFile(dir, strings.TrimPrefix(p, "/"))
	if err != nil {
		return "", ErrNotFound
	}
	fi, err := os.Stat(full)
	if err != nil || !fi.Mode().IsRegular() {
		return "", ErrNotFound
	}
	caps := capsOf(nil)
	if env.Kinds != nil {
		if k, err := env.Kinds(m.Kind); err == nil {
			caps = capsOf(k)
		}
	}
	if fi.Size() > caps.File {
		return "", ErrNotFound
	}
	return full, nil
}
