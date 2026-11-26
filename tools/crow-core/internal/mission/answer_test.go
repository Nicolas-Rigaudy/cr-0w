package mission

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"crow-core/contract"
	"crow-core/internal/brief"
	"crow-core/internal/testfix"
)

const example = "20261109-build-item-095"

func exampleDir() string { return filepath.Join(testfix.RunSpec, "examples", example) }

// exampleLines are the example log lines as stored.
func exampleLines(t *testing.T) [][]byte {
	t.Helper()
	return bytes.Split(bytes.TrimSpace(testfix.Read(t, filepath.Join(exampleDir(), "events.jsonl"))), []byte("\n"))
}

// replayCore appends example lines [from, to) as events, the way core wrote them.
func replayCore(t *testing.T, env Env, lines [][]byte, from, to int) {
	t.Helper()
	for _, l := range lines[from:to] {
		var e contract.Event
		if err := json.Unmarshal(l, &e); err != nil {
			t.Fatal(err)
		}
		if _, err := env.Store.Append(e); err != nil {
			t.Fatal(err)
		}
	}
}

func copyExample(t *testing.T, dir string, names ...string) {
	t.Helper()
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), testfix.Read(t, filepath.Join(exampleDir(), n)), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func at(env *Env, ts string) {
	tm, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		panic(err)
	}
	env.Now = func() time.Time { return tm }
}

// replayExample replays the example mission up to the ok on report v3 and returns the env,
// the mission dir and the example lines.
func replayExample(t *testing.T) (Env, string, [][]byte) {
	t.Helper()
	env, root := newEnv(t)
	lines := exampleLines(t)
	dir := filepath.Join(root, example)
	replayCore(t, env, lines, 0, 7) // created .. review on plan v1
	copyExample(t, dir, "plan.v1.md")

	at(&env, "2026-11-09T09:12:00+01:00")
	if _, err := env.Answer(example, Reply{Verb: "go", Surface: "deck"}); err == nil ||
		!strings.Contains(err.Error(), "blocking question q1") {
		t.Fatalf("go with q1 open: err = %v", err)
	}
	if _, err := env.Answer(example, Reply{Verb: "dig", Surface: "deck", Words: []string{
		"q1: refuse.", "Also list the repos with no scope in `crow doctor`."}}); err != nil {
		t.Fatal(err)
	}

	replayCore(t, env, lines, 8, 13) // v2 run, gate ready, review
	copyExample(t, dir, "plan.v2.md")
	at(&env, "2026-11-09T09:40:00+01:00")
	if _, err := env.Answer(example, Reply{Verb: "go", Surface: "nest"}); err != nil {
		t.Fatal(err)
	}
	replayCore(t, env, lines, 14, 19) // v3 run, gate ready, review
	copyExample(t, dir, "report.v3.md")
	at(&env, "2026-11-09T11:02:00+01:00")
	res, err := env.Answer(example, Reply{Verb: "ok", Surface: "nest"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res.Validations, " ") != "push mr" {
		t.Fatalf("validations = %v", res.Validations)
	}
	return env, dir, lines
}

// AC1: the example mission replays to the same answers files and events, both
// validation hashes included.
func TestAnswerReplaysTheExampleMission(t *testing.T) {
	_, dir, lines := replayExample(t)
	got := bytes.Split(bytes.TrimSpace(testfix.Read(t, filepath.Join(dir, "events.jsonl"))), []byte("\n"))
	if len(got) != 22 {
		t.Fatalf("%d lines, want 22", len(got))
	}
	for i := range got {
		if !bytes.Equal(got[i], lines[i]) {
			t.Errorf("line %d:\n got %s\nwant %s", i+1, got[i], lines[i])
		}
	}
	for _, n := range []string{"answers.v1.md", "answers.v2.md", "answers.v3.md"} {
		if g, w := testfix.Read(t, filepath.Join(dir, n)), testfix.Read(t, filepath.Join(exampleDir(), n)); !bytes.Equal(g, w) {
			t.Errorf("%s:\n%s\nwant\n%s", n, g, w)
		}
	}
	if err := CheckChain(filepath.Join(exampleDir(), "events.jsonl")); err != nil {
		t.Errorf("example chain: %v", err)
	}
}

func TestVerifyPassesThenRefusesAChangedEntry(t *testing.T) {
	env, dir, _ := replayExample(t)
	vs, err := env.Verify(example)
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) != 2 || vs[0].Data.Action != "push" || vs[1].TS != "2026-11-09T11:02:02+01:00" {
		t.Fatalf("verified = %+v", vs)
	}
	// AC2: a publish entry changed after ok fails the hash check.
	p := filepath.Join(dir, "report.v3.md")
	b := bytes.Replace(testfix.Read(t, p), []byte("-> main"), []byte("-> release"), 1)
	if err := os.WriteFile(p, b, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := env.Verify(example); err == nil || !strings.Contains(err.Error(), "entry 2 (mr) changed") {
		t.Fatalf("err = %v", err)
	}
	if ev := events(t, filepath.Dir(dir), example); ev[len(ev)-1] != "incident" {
		t.Errorf("last event %s, want incident", ev[len(ev)-1])
	}
}

// AC2: a message after ok drops the validation.
func TestMessageAfterOkDropsTheValidation(t *testing.T) {
	env, dir, _ := replayExample(t)
	at(&env, "2026-11-09T11:05:00+01:00")
	res, err := env.Answer(example, Reply{Verb: "msg", Surface: "deck", Words: []string{"Squash the commits first."}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res.Dropped, " ") != "push mr" || len(res.Validations) != 0 {
		t.Fatalf("result = %+v", res)
	}
	if m := show(t, env, example); m.Validation != nil {
		t.Errorf("mission validation = %s, want none", *m.Validation)
	}
	if _, err := env.Verify(example); err == nil || !strings.Contains(err.Error(), "no live validation") {
		t.Fatalf("err = %v", err)
	}
	// A new ok validates again, chained after the first two.
	at(&env, "2026-11-09T11:06:00+01:00")
	if _, err := env.Answer(example, Reply{Verb: "ok", Surface: "nest"}); err != nil {
		t.Fatal(err)
	}
	if _, err := env.Verify(example); err != nil {
		t.Fatal(err)
	}
	if b := testfix.Read(t, filepath.Join(dir, "answers.v3.md")); !bytes.HasSuffix(b, []byte("verb: ok\nsurface: nest\nts: 2026-11-09T11:06:00+01:00\n---\n1. Squash the commits first.\n")) {
		t.Errorf("answers.v3.md = %s", b)
	}
}

// AC3: an edited or replayed validation line breaks the chain check.
func TestChainBreaksOnEditOrReplay(t *testing.T) {
	src := exampleLines(t)
	write := func(ls [][]byte) string {
		p := filepath.Join(t.TempDir(), "events.jsonl")
		if err := os.WriteFile(p, append(bytes.Join(ls, []byte("\n")), '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	edited := make([][]byte, len(src))
	copy(edited, src)
	edited[20] = bytes.Replace(src[20], []byte("mission/20261109"), []byte("other/20261109"), 1)
	if err := CheckChain(write(edited)); err == nil || !strings.Contains(err.Error(), "line 22") {
		t.Errorf("edited: err = %v", err)
	}
	replayed := append(append(append([][]byte{}, src[:22]...), src[20]), src[22:]...)
	if err := CheckChain(write(replayed)); err == nil || !strings.Contains(err.Error(), "line 23") {
		t.Errorf("replayed: err = %v", err)
	}
	forged := make([][]byte, len(src))
	copy(forged, src)
	forged[21] = bytes.Replace(src[21], []byte(`"actor":"owner"`), []byte(`"actor":"core"`), 1)
	if err := CheckChain(write(forged)); err == nil || !strings.Contains(err.Error(), "only the owner") {
		t.Errorf("forged: err = %v", err)
	}
}

// AC4: a version 2 that skips an answer point is refused, the missing point named.
func TestCoverageNamesTheMissingPoint(t *testing.T) {
	dir := t.TempDir()
	copyExample(t, dir, "plan.v1.md", "answers.v1.md", "plan.v2.md")
	_, cb, why := versionFile(dir, example, 2)
	if cb == nil {
		t.Fatal(why)
	}
	if err := Coverage(dir, example, 2, cb); err != nil {
		t.Fatalf("example v2: %v", err)
	}
	rs := (*cb.Responses)[:1]
	cb.Responses = &rs
	if err := Coverage(dir, example, 2, cb); err == nil || err.Error() != "v2 responses miss a2" {
		t.Errorf("err = %v", err)
	}
	cb.Responses = &[]contract.Response{}
	if err := Coverage(dir, example, 2, cb); err == nil || err.Error() != "v2 responses miss a1, a2, q1" {
		t.Errorf("err = %v", err)
	}
}

func TestTickRefusesAVersionThatSkipsAPoint(t *testing.T) {
	env, root := newEnv(t)
	lines := exampleLines(t)
	dir := filepath.Join(root, example)
	replayCore(t, env, lines, 0, 7)
	copyExample(t, dir, "plan.v1.md", "answers.v1.md")
	replayCore(t, env, lines, 7, 10) // answer, review -> running, run_start v2
	b := bytes.Replace(testfix.Read(t, filepath.Join(exampleDir(), "plan.v2.md")),
		[]byte("  - {ref: a2, outcome: done, note: Step 4.}\n"), nil, 1)
	if err := os.WriteFile(filepath.Join(dir, "plan.v2.md"), b, 0o600); err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{ends: map[string]RunState{"b2c3d4e5": ended("done", 0)}}
	rep, err := env.Tick(r, 2)
	if err != nil {
		t.Fatal(err)
	}
	// AC3: the first miss sends the version back once with the error.
	if len(rep.Finished) != 1 || rep.Finished[0].State != "queued" ||
		!strings.Contains(rep.Finished[0].Detail, "plan.v2.md: v2 responses miss a2") {
		t.Fatalf("first miss: report = %+v", rep)
	}
	m := show(t, env, example)
	k, err := env.Kinds(m.Kind)
	if err != nil {
		t.Fatal(err)
	}
	started, err := env.resume(r, m, k, false, false) // the tick holds build kinds (phase 3)
	if err != nil {
		t.Fatal(err)
	}
	note, err := os.ReadFile(filepath.Join(dir, ResumeNote))
	if err != nil || !strings.Contains(string(note), "core rejected your last file: plan.v2.md: v2 responses miss a2") {
		t.Errorf("note = %q, %v", note, err)
	}
	// The resumed run writes the same file: the second miss fails the mission.
	r.ends[started.Run] = ended("done", 0)
	if rep, err = env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	if len(rep.Finished) != 1 || rep.Finished[0].State != "failed" ||
		!strings.Contains(rep.Finished[0].Detail, "plan.v2.md: v2 responses miss a2") {
		t.Fatalf("second miss: report = %+v", rep)
	}
}

// deskBugs copies the desk-bugs plan mission (2026-11-25): v1 asked q1 to q3, the owner
// answered "i don't understand the question", v2 rewrote q1 as q4 and responded to a1 only.
const deskBugs = "20261125-plan-nest-desk-bugs"

func deskBugsDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, n := range []string{"plan.v1.md", "answers.v1.md", "plan.v2.md"} {
		if err := os.WriteFile(filepath.Join(dir, n), testfix.Read(t, filepath.Join("testdata", "desk-bugs", n)), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// AC1: the live v2 misses q1 to q3; the same v2 with q1 replaced, q2 withdrawn and q3
// done passes.
func TestCoverageAcceptsWithdrawnAndReplaced(t *testing.T) {
	dir := deskBugsDir(t)
	_, cb, why := versionFile(dir, deskBugs, 2)
	if cb == nil {
		t.Fatal(why)
	}
	if err := Coverage(dir, deskBugs, 2, cb); err == nil || err.Error() != "v2 responses miss q1, q2, q3" {
		t.Fatalf("live v2: err = %v", err)
	}
	note := func(s string) *string { return &s }
	rs := append(*cb.Responses,
		contract.Response{Ref: "q1", Outcome: "replaced", Note: note("replaced by q4")},
		contract.Response{Ref: "q2", Outcome: "withdrawn", Note: note("bug 1 is built on item_210")},
		contract.Response{Ref: "q3", Outcome: "done", Note: note("item_210 reused")})
	cb.Responses = &rs
	if err := Coverage(dir, deskBugs, 2, cb); err != nil {
		t.Errorf("fixed v2: %v", err)
	}
	for _, r := range rs {
		if !slices.Contains(contract.ResponseOutcomes, r.Outcome) {
			t.Errorf("outcome %s is not a contract outcome", r.Outcome)
		}
	}
}

// AC2: the resume note of version 2 names every question of version 1, and the error of
// a version core sent back.
func TestResumeNoteNamesThePreviousQuestions(t *testing.T) {
	dir := deskBugsDir(t)
	answers, err := brief.ReadAnswers(dir)
	if err != nil {
		t.Fatal(err)
	}
	path, err := writeResumeNote(dir, deskBugs, 1, 2, answers, "plan.v2.md: v2 responses miss q1, q2, q3")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Version 1 asked q1, q2, q3.", "`withdrawn` when dropped", "`replaced` when rewritten",
		"core rejected your last file: plan.v2.md: v2 responses miss q1, q2, q3", "i don't understand the question",
	} {
		if !strings.Contains(string(b), want) {
			t.Errorf("note misses %q:\n%s", want, b)
		}
	}
}

func TestAnswerRefusesWrongVerbAndState(t *testing.T) {
	env, root := newEnv(t)
	lines := exampleLines(t)
	replayCore(t, env, lines, 0, 7)
	copyExample(t, filepath.Join(root, example), "plan.v1.md")
	if _, err := env.Answer(example, Reply{Verb: "ok", Surface: "nest"}); err == nil || !strings.Contains(err.Error(), "want one of [go dig no]") {
		t.Errorf("ok on a plan: err = %v", err)
	}
	if _, err := env.Answer(example, Reply{Verb: "no", Surface: "chat"}); err == nil {
		t.Error("chat surface accepted")
	}
	replayCore(t, env, lines, 7, 9) // answer, review -> running
	if _, err := env.Answer(example, Reply{Verb: "dig", Surface: "nest"}); err == nil || !strings.Contains(err.Error(), "is running") {
		t.Errorf("answer while running: err = %v", err)
	}
}
