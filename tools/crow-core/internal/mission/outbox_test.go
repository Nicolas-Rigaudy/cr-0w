package mission

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"crow-core/contract"
	"crow-core/internal/testfix"
)

func writeFile(t *testing.T, path string, size int64) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := f.Truncate(size); err != nil {
		t.Fatal(err)
	}
}

func dropCauses(drops []Drop) string {
	var cs []string
	for _, d := range drops {
		cs = append(cs, d.Path+"="+d.Cause)
	}
	return strings.Join(cs, ",")
}

// AC1: over the per-file cap, or over the total once earlier files filled it, is dropped.
func TestIngestCaps(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "outbox", "a.html"), 40)
	writeFile(t, filepath.Join(dir, "outbox", "big.png"), 200)
	writeFile(t, filepath.Join(dir, "outbox", "b.png"), 50)
	writeFile(t, filepath.Join(dir, "outbox", "c.png"), 30)
	paths := []string{"outbox/a.html", "outbox/big.png", "outbox/b.png", "outbox/a.html", "outbox/c.png", "outbox/none.md"}
	kept, drops := Ingest(dir, paths, Caps{File: 100, Total: 80})
	if strings.Join(kept, ",") != "outbox/a.html,outbox/c.png" {
		t.Errorf("kept = %v", kept)
	}
	want := "outbox/big.png=artifact_over_cap,outbox/b.png=artifact_over_cap,outbox/none.md=artifact_missing"
	if got := dropCauses(drops); got != want {
		t.Errorf("drops = %s, want %s", got, want)
	}
}

// AC2: a path with "..", outside outbox/, absolute, or a link that leads out is refused.
func TestIngestRefusesPathsOutsideOutbox(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "outbox", "ok.md"), 1)
	writeFile(t, filepath.Join(dir, "secret.md"), 1)
	writeFile(t, filepath.Join(dir, "notoutbox", "x.md"), 1)
	if err := os.Symlink(filepath.Join(dir, "secret.md"), filepath.Join(dir, "outbox", "link.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "outbox", "sub"), 0o700); err != nil {
		t.Fatal(err)
	}
	paths := []string{
		"outbox/../secret.md", "outbox/sub/../../secret.md", filepath.Join(dir, "secret.md"),
		"notoutbox/x.md", "outbox/link.md", "outbox/sub", "outbox/ok.md",
	}
	kept, drops := Ingest(dir, paths, capsOf(nil))
	if strings.Join(kept, ",") != "outbox/ok.md" {
		t.Errorf("kept = %v", kept)
	}
	if len(drops) != 6 {
		t.Fatalf("drops = %+v", drops)
	}
	for _, d := range drops {
		if d.Cause != "artifact_refused" {
			t.Errorf("%s: cause %s", d.Path, d.Cause)
		}
	}
}

func TestCapsOfKind(t *testing.T) {
	file, total := 2.5, 10.0
	k := &contract.Kind{}
	k.Caps.OutboxFileMB, k.Caps.OutboxTotalMB = &file, &total
	if c := capsOf(k); c.File != 5<<19 || c.Total != 10<<20 {
		t.Errorf("caps = %+v", c)
	}
	if c := capsOf(nil); c.File != 20<<20 || c.Total != 100<<20 {
		t.Errorf("default caps = %+v", c)
	}
}

// questions copies the spec example for mission id and version n.
func questions(t *testing.T, root, id string, n int) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(testfix.RunSpec, "examples", "gates", "questions.json"))
	if err != nil {
		t.Fatal(err)
	}
	s := strings.Replace(string(b), "20261109-research-android-push-options", id, 1)
	s = strings.Replace(s, `"version": 1`, fmt.Sprintf(`"version": %d`, n), 1)
	writeFile(t, filepath.Join(root, id, QuestionsFile), 0)
	if err := os.WriteFile(filepath.Join(root, id, QuestionsFile), []byte(s), 0o600); err != nil {
		t.Fatal(err)
	}
}

func startOne(t *testing.T, env Env) (string, *fakeRunner) {
	t.Helper()
	m, err := env.New(research("Compare push providers for the Android app."))
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{}}
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	return m.ID, r
}

// AC3: a questions file parks the mission and its questions reach the inbox.
func TestTickParksOnQuestions(t *testing.T) {
	env, root := newEnv(t)
	id, r := startOne(t, env)
	questions(t, root, id, 1)
	r.ends["r1"] = ended("done", 0)
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Finished) != 1 || rep.Finished[0].State != "waiting" || rep.Finished[0].Detail != "1 question(s): q1" {
		t.Fatalf("report = %+v", rep.Finished)
	}
	if got := show(t, env, id); got.State != "waiting" {
		t.Fatalf("state = %s", got.State)
	}
	want := "created,state,state,run_start,run_end,question,state"
	if e := strings.Join(events(t, root, id), ","); e != want {
		t.Errorf("events = %s, want %s", e, want)
	}
	d, err := env.Store.Dump()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, it := range d.Inbox {
		ids = append(ids, it.ID)
	}
	if strings.Join(ids, ",") != id+"/v1/q1" {
		t.Errorf("inbox = %v", ids)
	}
}

// A questions file left from an earlier version is not asked again.
func TestTickIgnoresStaleQuestions(t *testing.T) {
	env, root := newEnv(t)
	id, r := startOne(t, env)
	questions(t, root, id, 0)
	writeVersion(t, root, id, "report.v1.md", "ready", 1)
	r.ends["r1"] = ended("done", 0)
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	if got := show(t, env, id); got.State != "done" {
		t.Errorf("state = %s", got.State)
	}
}

// A bad questions file is noted and the turn ends on its version file.
func TestTickBadQuestionsFile(t *testing.T) {
	env, root := newEnv(t)
	id, r := startOne(t, env)
	writeFile(t, filepath.Join(root, id, QuestionsFile), 0)
	if err := os.WriteFile(filepath.Join(root, id, QuestionsFile), []byte(`{"contract":1}`), 0o600); err != nil {
		t.Fatal(err)
	}
	writeVersion(t, root, id, "report.v1.md", "ready", 1)
	r.ends["r1"] = ended("done", 0)
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	if got := show(t, env, id); got.State != "done" {
		t.Errorf("state = %s", got.State)
	}
	want := "created,state,state,run_start,run_end,incident,gate,state"
	if e := strings.Join(events(t, root, id), ","); e != want {
		t.Errorf("events = %s, want %s", e, want)
	}
}

// AC1/AC2 through a tick: dropped artifacts are logged and the mission still reaches its gate.
func TestTickLogsDroppedArtifacts(t *testing.T) {
	env, root := newEnv(t)
	id, r := startOne(t, env)
	writeFile(t, filepath.Join(root, id, "outbox", "small.html"), 10)
	writeFile(t, filepath.Join(root, id, "outbox", "huge.mp4"), 21<<20)
	md := fmt.Sprintf("# Report\n\n```crow\ncontract: 1\nmission: %s\nkind: research\nphase: report\nversion: 1\n"+
		"summary: one provider\nconfidence: 80\nverdict: ready\nartifacts:\n"+
		"  - path: outbox/small.html\n  - path: outbox/huge.mp4\n  - path: outbox/../brief.md\n```\n", id)
	if err := os.WriteFile(filepath.Join(root, id, "report.v1.md"), []byte(md), 0o600); err != nil {
		t.Fatal(err)
	}
	r.ends["r1"] = ended("done", 0)
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	if got := show(t, env, id); got.State != "done" {
		t.Errorf("state = %s", got.State)
	}
	want := "created,state,state,run_start,run_end,incident,incident,gate,state"
	if e := strings.Join(events(t, root, id), ","); e != want {
		t.Errorf("events = %s, want %s", e, want)
	}
	b, _ := os.ReadFile(filepath.Join(root, id, "events.jsonl"))
	for _, c := range []string{"artifact_over_cap", "artifact_refused"} {
		if !strings.Contains(string(b), c) {
			t.Errorf("no %s incident", c)
		}
	}
}
