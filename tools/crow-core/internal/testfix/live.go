package testfix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Missions in flight for the action tests (item_147), next to the finished fixtures.
const (
	Waiting = "20261110-research-waiting-on-owner" // parked on q1 of the example questions
	Running = "20261110-research-running-now"      // run f6a7b8c9 live
	LiveRun = "f6a7b8c9"
)

// Registry is a project registry with the demo project the fixtures use.
func Registry(t testing.TB) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "projects.yaml")
	if err := os.WriteFile(p, []byte("projects:\n  - name: demo\n    scope: personal\n    path: /x/demo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// Live adds the waiting and the running mission to a missions root.
func Live(t testing.TB, root string) {
	t.Helper()
	write(t, filepath.Join(root, Waiting, "events.jsonl"), log(Waiting, "waiting-on-owner", "e5f6a7b8",
		`{"run":"e5f6a7b8","cause":"done","exit":0}`, "run_end", "08:20:00",
		`{"file":"outbox/questions.json","version":1,"ids":["q1"]}`, "question", "08:20:01",
		`{"from":"running","to":"waiting"}`, "state", "08:20:02"))
	write(t, filepath.Join(root, Waiting, "outbox", "questions.json"),
		Read(t, filepath.Join(RunSpec, "examples", "gates", "questions.json")))
	write(t, filepath.Join(root, Running, "events.jsonl"), log(Running, "running-now", LiveRun))
}

// log is created, queued, running and run_start, then (data, type, time) triples.
func log(id, key, run string, more ...string) []byte {
	line := func(ts, typ, actor, data string) string {
		return fmt.Sprintf(`{"ts":"2026-11-10T%s+01:00","mission":"%s","type":"%s","actor":"%s","data":%s}`, ts, id, typ, actor, data)
	}
	ls := []string{
		line("08:00:00", "created", "owner", `{"kind":"research","key":"`+key+`","surface":"cli","project":"demo","scope":"personal","data_class":"personal","goal":"Compare push providers.","end_state":"A report with one recommendation.","autonomy":"auto","gates":[],"budget":{"time_min":60,"budget_usd":5,"resumes":2},"file_scope":[]}`),
		line("08:00:01", "state", "core", `{"from":"preflight","to":"queued"}`),
		line("08:05:00", "state", "core", `{"from":"queued","to":"running"}`),
		line("08:05:01", "run_start", "core", `{"run":"`+run+`","account":"main","account_class":"personal","cold":true,"version":1}`),
	}
	for i := 0; i+2 < len(more); i += 3 {
		actor := "core"
		if more[i+1] == "question" {
			actor = "run"
		}
		ls = append(ls, line(more[i+2], more[i+1], actor, more[i]))
	}
	return []byte(strings.Join(ls, "\n") + "\n")
}

func write(t testing.TB, p string, b []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, b, 0o600); err != nil {
		t.Fatal(err)
	}
}
