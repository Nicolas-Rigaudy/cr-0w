// Package mission creates, reads and cancels missions (item_132). It resolves a new
// mission from its kind file and the project registry, then writes it as events: a
// created event carrying the resolved gates, autonomy and budget, and the first state
// move. The store stays a view of those events.
package mission

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"crow-core/contract"
	"crow-core/internal/brief"
	"crow-core/internal/eventlog"
	"crow-core/internal/project"
	"crow-core/internal/store"
)

// Spec is what the owner gives to `mission new`. Empty optional fields take the defaults
// of the kind file and the registry.
type Spec struct {
	Kind, Project, Goal, EndState string
	Key                           string    // short part of the id; derived from the goal when empty
	Scope, DataClass, Autonomy    string    // overrides
	Gates                         *[]string // override; an empty list means no gate
	Surface                       string    // cli by default
}

// Env is what the commands need around the store.
type Env struct {
	Store    *store.Store
	Kinds    project.Kinds
	KindDirs []string // where kind files live, to list the valid names
	Projects string   // project registry file
	Prefs    string   // owner prefs for briefs (brief.LoadPrefs)
	Now      func() time.Time
	StopRun  func(run string) error // stops a live cdx run (crow runs stop)
}

// Project is the part of a registry entry missions use.
type Project struct {
	Name    string   `yaml:"name"`
	Aliases []string `yaml:"aliases"`
	Scope   string   `yaml:"scope"`
}

// LoadProjects reads the projects list of the registry.
func LoadProjects(path string) ([]Project, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("project registry: %w", err)
	}
	var reg struct {
		Projects []Project `yaml:"projects"`
	}
	if err := yaml.Unmarshal(b, &reg); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return reg.Projects, nil
}

func findProject(ps []Project, name string) (Project, error) {
	var names []string
	for _, p := range ps {
		if p.Name == name || slices.Contains(p.Aliases, name) {
			return p, nil
		}
		names = append(names, p.Name)
	}
	slices.Sort(names)
	return Project{}, fmt.Errorf("unknown project %q (valid: %s)", name, strings.Join(names, ", "))
}

// KindNames lists the kind files found in dirs.
func KindNames(dirs []string) []string {
	var names []string
	for _, d := range dirs {
		fs, _ := filepath.Glob(filepath.Join(d, "*.yaml"))
		for _, f := range fs {
			n := strings.TrimSuffix(filepath.Base(f), ".yaml")
			if !slices.Contains(names, n) {
				names = append(names, n)
			}
		}
	}
	slices.Sort(names)
	return names
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

var filler = []string{"a", "an", "the", "for", "of", "to", "and", "in", "on", "with"}

// slug makes a key from the goal: its first five words, filler dropped, lowercase, at
// most 40 characters.
func slug(goal string) string {
	words := strings.Fields(nonSlug.ReplaceAllString(strings.ToLower(goal), " "))
	var out string
	for _, w := range words {
		if slices.Contains(filler, w) {
			continue
		}
		next := w
		if out != "" {
			next = out + "-" + w
		}
		if len(next) > 40 {
			break
		}
		out = next
		if strings.Count(out, "-") >= 4 {
			break
		}
	}
	return out
}

// Resolve checks a spec against the kind file and the registry and returns the id and
// the created event data, with every default filled.
func (env Env) Resolve(s Spec) (string, contract.CreatedData, error) {
	var d contract.CreatedData
	if s.Goal == "" || s.EndState == "" {
		return "", d, errors.New("a mission needs a goal and an end state")
	}
	k, err := env.Kinds(s.Kind)
	if err == nil && k == nil {
		err = fmt.Errorf("unknown kind %q (valid: %s)", s.Kind, strings.Join(KindNames(env.KindDirs), ", "))
	}
	if err != nil {
		return "", d, err
	}
	ps, err := LoadProjects(env.Projects)
	if err != nil {
		return "", d, err
	}
	p, err := findProject(ps, s.Project)
	if err != nil {
		return "", d, err
	}

	scope := p.Scope
	if s.Scope != "" {
		scope = s.Scope
	}
	if !slices.Contains(k.Scopes, scope) {
		return "", d, fmt.Errorf("kind %s does not run in scope %q (valid: %s)", k.Name, scope, strings.Join(k.Scopes, ", "))
	}
	dataClass := s.DataClass
	if dataClass == "" {
		dataClass = scope
		if scope != "work" {
			dataClass = "personal"
		}
	}
	classes := []string{"work", "personal", "private"}
	if k.DataClasses != nil {
		classes = *k.DataClasses
	}
	if !slices.Contains(classes, dataClass) {
		return "", d, fmt.Errorf("kind %s does not take data class %q (valid: %s)", k.Name, dataClass, strings.Join(classes, ", "))
	}

	gates, autonomy := slices.Clone(k.Gates), k.Autonomy
	if k.ByScope != nil {
		if o, ok := (*k.ByScope)[scope]; ok {
			if o.Gates != nil {
				gates = slices.Clone(*o.Gates)
			}
			if o.Autonomy != nil {
				autonomy = *o.Autonomy
			}
		}
	}
	if s.Gates != nil {
		gates = slices.Clone(*s.Gates)
	}
	for _, g := range gates {
		if !slices.Contains([]string{"plan", "report", "mockup"}, g) {
			return "", d, fmt.Errorf("unknown gate %q (valid: plan, report, mockup)", g)
		}
	}
	if s.Autonomy != "" {
		autonomy = s.Autonomy
	}
	if !slices.Contains([]string{"ask", "notify", "auto"}, autonomy) {
		return "", d, fmt.Errorf("unknown autonomy %q (valid: ask, notify, auto)", autonomy)
	}
	if gates == nil {
		gates = []string{}
	}

	key := s.Key
	if key == "" {
		key = slug(s.Goal)
	}
	if ok, _ := regexp.MatchString(k.Key.Pattern, key); !ok {
		return "", d, fmt.Errorf("key %q does not match the %s key pattern %s; give one with --key", key, k.Name, k.Key.Pattern)
	}
	id := env.Now().Format("20060102") + "-" + k.Name + "-" + strings.ToLower(key)
	if err := eventlog.ValidMission(id); err != nil {
		return "", d, err
	}

	surface := s.Surface
	if surface == "" {
		surface = "cli"
	}
	budget := contract.Budget{TimeMin: k.Caps.TimeMin, BudgetUSD: k.Caps.BudgetUSD, Resumes: k.Caps.Resumes}
	fileScope := []string{}
	name, goal, end := p.Name, s.Goal, s.EndState
	d = contract.CreatedData{
		Kind: k.Name, Key: key, Surface: surface,
		Project: &name, Scope: &scope, DataClass: &dataClass, Goal: &goal, EndState: &end,
		Autonomy: &autonomy, Gates: &gates, Budget: &budget, FileScope: &fileScope,
	}
	return id, d, nil
}

// New creates a mission: its dir, a created event with the resolved values, then the
// move to queued (preflight questions come with item_113).
func (env Env) New(s Spec) (contract.Mission, error) {
	id, d, err := env.Resolve(s)
	if err != nil {
		return contract.Mission{}, err
	}
	if _, err := env.find(id); err == nil {
		return contract.Mission{}, fmt.Errorf("mission %s exists; give another --key", id)
	}
	ts := env.Now().Format(time.RFC3339)
	if err := env.append(id, "owner", "created", ts, d); err != nil {
		return contract.Mission{}, err
	}
	if err := env.append(id, "core", "state", ts, contract.StateData{From: "preflight", To: "queued"}); err != nil {
		return contract.Mission{}, err
	}
	return env.find(id)
}

func (env Env) append(id, actor, typ, ts string, data any) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = env.Store.Append(contract.Event{TS: ts, Mission: id, Type: typ, Actor: actor, Data: raw})
	return err
}

// List returns every mission, synced from the logs first.
func (env Env) List() ([]contract.Mission, error) {
	if _, err := env.Store.SyncAll(); err != nil {
		return nil, err
	}
	d, err := env.Store.Dump()
	return d.Missions, err
}

// Show returns one mission by id or unique id prefix.
func (env Env) Show(ref string) (contract.Mission, error) {
	return env.find(ref)
}

func (env Env) find(ref string) (contract.Mission, error) {
	ms, err := env.List()
	if err != nil {
		return contract.Mission{}, err
	}
	var hits []contract.Mission
	for _, m := range ms {
		if m.ID == ref {
			return m, nil
		}
		if strings.HasPrefix(m.ID, ref) {
			hits = append(hits, m)
		}
	}
	if len(hits) == 1 {
		return hits[0], nil
	}
	if len(hits) == 0 {
		return contract.Mission{}, fmt.Errorf("no mission %q", ref)
	}
	var ids []string
	for _, m := range hits {
		ids = append(ids, m.ID)
	}
	return contract.Mission{}, fmt.Errorf("%q matches %d missions: %s", ref, len(hits), strings.Join(ids, ", "))
}

// Cancel moves a mission to cancelled. A running mission has its live run stopped first;
// a run that is already gone is reported in warn and the cancel goes on.
func (env Env) Cancel(ref, reason string, warn func(string)) (contract.Mission, error) {
	m, err := env.find(ref)
	if err != nil {
		return m, err
	}
	if !contract.CanMove(m.State, "cancelled") {
		return m, fmt.Errorf("%s is %s and cannot be cancelled", m.ID, m.State)
	}
	if m.State == "running" {
		for _, r := range m.Runs {
			if r.Ended != nil {
				continue
			}
			if env.StopRun == nil {
				return m, fmt.Errorf("%s: no way to stop run %s (CROW_ENGINE unset)", m.ID, r.Run)
			}
			if err := env.StopRun(r.Run); err != nil {
				warn(fmt.Sprintf("run %s: %v", r.Run, err))
			}
		}
	}
	sd := contract.StateData{From: m.State, To: "cancelled"}
	if reason != "" {
		sd.Reason = &reason
	}
	if err := env.append(m.ID, "owner", "state", env.Now().Format(time.RFC3339), sd); err != nil {
		return m, err
	}
	return env.find(m.ID)
}

// Brief gathers what the brief of a mission is built from: its record, kind file, the key
// of its created event, the prefs for its kind and every answers file. It returns the
// mission dir the brief belongs in.
func (env Env) Brief(ref string) (string, brief.Input, error) {
	m, err := env.find(ref)
	if err != nil {
		return "", brief.Input{}, err
	}
	k, err := env.Kinds(m.Kind)
	if err != nil {
		return "", brief.Input{}, fmt.Errorf("%s: kind %s: %w", m.ID, m.Kind, err)
	}
	in := brief.Input{Mission: m, Kind: k}
	dir, err := filepath.Abs(filepath.Join(env.Store.Root(), m.ID))
	if err != nil {
		return "", in, err
	}
	in.Dir = dir
	recs, err := eventlog.Read(eventlog.Path(env.Store.Root(), m.ID), 0)
	if err != nil {
		return "", in, err
	}
	for _, r := range recs {
		if r.Event.Type != "created" {
			continue
		}
		d, err := contract.DecodeEventData(r.Event)
		if err != nil {
			return "", in, err
		}
		in.Key = d.(*contract.CreatedData).Key
		break
	}
	ps, err := brief.LoadPrefs(env.Prefs)
	if err != nil {
		return "", in, err
	}
	in.Prefs = brief.ForKind(ps, k.Name)
	if in.Answers, err = brief.ReadAnswers(dir); err != nil {
		return "", in, err
	}
	return dir, in, nil
}
