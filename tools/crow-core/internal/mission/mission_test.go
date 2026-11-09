package mission

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"crow-core/contract"
	"crow-core/internal/project"
	"crow-core/internal/store"
	"crow-core/internal/testfix"
)

const registry = `projects:
  - name: demo
    aliases: [dm]
    scope: personal
    path: /x/demo
  - name: shop
    scope: work
    path: /x/shop
`

func newEnv(t *testing.T) (Env, string) {
	t.Helper()
	dir := t.TempDir()
	reg := filepath.Join(dir, "projects.yaml")
	if err := os.WriteFile(reg, []byte(registry), 0o600); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(dir, "missions")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	kinds := project.KindDirs(testfix.KindsDir, "")
	s, err := store.Open(filepath.Join(dir, "store.db"), root, kinds)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	now := time.Date(2026, 11, 9, 14, 0, 0, 0, time.FixedZone("CET", 3600))
	return Env{
		Store: s, Kinds: kinds, KindDirs: []string{testfix.KindsDir},
		Projects: reg, Now: func() time.Time { return now },
	}, root
}

func research(goal string) Spec {
	return Spec{Kind: "research", Project: "demo", Goal: goal, EndState: "A report with one recommendation."}
}

// AC1: new fills the kind defaults, writes created then queued.
func TestNewFillsKindDefaults(t *testing.T) {
	env, root := newEnv(t)
	m, err := env.New(research("Compare push providers for the Android app."))
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "20261109-research-compare-push-providers-android-app" {
		t.Errorf("id = %s", m.ID)
	}
	if m.State != "queued" || m.Scope != "personal" || m.DataClass != "personal" || m.Autonomy != "auto" {
		t.Errorf("mission = %+v", m)
	}
	if len(m.Gates) != 0 || m.Gates == nil {
		t.Errorf("gates = %#v, want []", m.Gates)
	}
	if m.Budget != (contract.Budget{TimeMin: 60, BudgetUSD: 5, Resumes: 2}) {
		t.Errorf("budget = %+v", m.Budget)
	}
	b, err := os.ReadFile(filepath.Join(root, m.ID, "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var types []string
	for _, l := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var e contract.Event
		if err := json.Unmarshal([]byte(l), &e); err != nil {
			t.Fatal(err)
		}
		types = append(types, e.Type+"/"+e.Actor)
	}
	if want := []string{"created/owner", "state/core"}; !slices.Equal(types, want) {
		t.Errorf("events = %v, want %v", types, want)
	}
	if !strings.Contains(string(b), `"budget":{"time_min":60`) {
		t.Errorf("created event does not carry the resolved budget: %s", b)
	}
}

func TestNewAppliesScopeOverridesAndFlags(t *testing.T) {
	env, _ := newEnv(t)
	k, err := env.Kinds("build")
	if err != nil || k == nil || k.ByScope == nil {
		t.Fatalf("build kind: %v", err)
	}
	o := (*k.ByScope)["work"]
	m, err := env.New(Spec{Kind: "build", Project: "shop", Key: "item-095", Goal: "Add login.", EndState: "An MR."})
	if err != nil {
		t.Fatal(err)
	}
	if m.Scope != "work" || m.DataClass != "work" {
		t.Errorf("scope/class = %s/%s", m.Scope, m.DataClass)
	}
	if o.Gates != nil && !slices.Equal(m.Gates, *o.Gates) {
		t.Errorf("gates = %v, want by_scope %v", m.Gates, *o.Gates)
	}
	if o.Autonomy != nil && m.Autonomy != *o.Autonomy {
		t.Errorf("autonomy = %s, want by_scope %s", m.Autonomy, *o.Autonomy)
	}
	none := []string{}
	m, err = env.New(Spec{Kind: "research", Project: "dm", Key: "flags", Goal: "g", EndState: "e", Gates: &none, Autonomy: "ask"})
	if err != nil {
		t.Fatal(err)
	}
	if m.Project != "demo" || m.Autonomy != "ask" || len(m.Gates) != 0 {
		t.Errorf("alias/flags not applied: %+v", m)
	}
	if _, err := env.New(Spec{Kind: "research", Project: "dm", Key: "flags", Goal: "g", EndState: "e"}); err == nil || !strings.Contains(err.Error(), "exists") {
		t.Errorf("duplicate id: %v", err)
	}
}

// AC4: unknown names are refused with the valid ones listed.
func TestUnknownKindOrProjectListsValidNames(t *testing.T) {
	env, _ := newEnv(t)
	s := research("g")
	s.Kind = "nope"
	_, err := env.New(s)
	if err == nil || !strings.Contains(err.Error(), "research") || !strings.Contains(err.Error(), "build") {
		t.Errorf("unknown kind: %v", err)
	}
	s = research("g")
	s.Project = "nope"
	_, err = env.New(s)
	if err == nil || !strings.Contains(err.Error(), "valid: demo, shop") {
		t.Errorf("unknown project: %v", err)
	}
	s = research("g")
	s.Scope = "elsewhere"
	if _, err = env.New(s); err == nil || !strings.Contains(err.Error(), "valid:") {
		t.Errorf("bad scope: %v", err)
	}
}

func keys(t *testing.T, b []byte) []string {
	t.Helper()
	var ms []map[string]json.RawMessage
	if err := json.Unmarshal(b, &ms); err != nil {
		t.Fatal(err)
	}
	var ks []string
	for k := range ms[0] {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// AC2: list and show print the fixture record shape.
func TestListAndShowMatchTheFixtureShape(t *testing.T) {
	env, _ := newEnv(t)
	if _, err := env.New(research("Compare push providers.")); err != nil {
		t.Fatal(err)
	}
	ms, err := env.List()
	if err != nil {
		t.Fatal(err)
	}
	got, _ := json.Marshal(ms)
	want := testfix.Read(t, filepath.Join(testfix.Records, "missions.json"))
	if !reflect.DeepEqual(keys(t, got), keys(t, want)) {
		t.Errorf("list keys = %v, want %v", keys(t, got), keys(t, want))
	}
	var fixture []contract.Mission
	if err := contract.Decode(want, &fixture); err != nil {
		t.Fatal(err)
	}
	if err := contract.Decode(got, &fixture); err != nil {
		t.Errorf("list output does not decode strictly: %v", err)
	}
	m, err := env.Show("20261109-research")
	if err != nil || m.ID != ms[0].ID {
		t.Errorf("show by prefix: %v %v", m.ID, err)
	}
	if _, err := env.Show("2026"); err != nil {
		t.Errorf("unique prefix: %v", err)
	}
	if _, err := env.Show("nope"); err == nil {
		t.Error("show nope: want an error")
	}
}

func appendEv(t *testing.T, env Env, id, typ string, data any) {
	t.Helper()
	if err := env.append(id, "core", typ, env.Now().Format(time.RFC3339), data); err != nil {
		t.Fatal(err)
	}
}

func TestCancelStopsTheLiveRun(t *testing.T) {
	env, _ := newEnv(t)
	m, err := env.New(research("Run then cancel."))
	if err != nil {
		t.Fatal(err)
	}
	appendEv(t, env, m.ID, "state", contract.StateData{From: "queued", To: "running"})
	appendEv(t, env, m.ID, "run_start", contract.RunStartData{Run: "a1b2c3d4", AccountClass: "personal", Cold: true, Version: 1})
	var stopped []string
	env.StopRun = func(r string) error { stopped = append(stopped, r); return errors.New("no live process") }
	var warned []string
	m, err = env.Cancel(m.ID, "changed my mind", func(s string) { warned = append(warned, s) })
	if err != nil {
		t.Fatal(err)
	}
	if m.State != "cancelled" || !slices.Equal(stopped, []string{"a1b2c3d4"}) || len(warned) != 1 {
		t.Errorf("state %s, stopped %v, warned %v", m.State, stopped, warned)
	}
	_, err = env.Cancel(m.ID, "", func(string) {})
	if err == nil || !strings.Contains(err.Error(), "is cancelled") {
		t.Errorf("second cancel: %v", err)
	}
}

func TestCancelQueuedNeedsNoStop(t *testing.T) {
	env, _ := newEnv(t)
	m, err := env.New(research("Cancel while queued."))
	if err != nil {
		t.Fatal(err)
	}
	m, err = env.Cancel(m.ID, "", func(string) {})
	if err != nil || m.State != "cancelled" {
		t.Errorf("cancel queued: %s %v", m.State, err)
	}
}
