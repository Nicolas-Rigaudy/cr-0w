package api

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"crow-core/internal/testfix"
)

var csrf = map[string]string{CSRFHeader: "1", "Origin": "https://core.example"}

// actionServer is newServer plus a registry, a waiting and a running mission, and a
// StopRun that records the runs it stops.
func actionServer(t *testing.T) (*Server, string, *[]string) {
	t.Helper()
	s, root := newServer(t)
	testfix.Live(t, root)
	var stopped []string
	s.o.Env.Projects = testfix.Registry(t)
	s.o.Env.KindDirs = []string{testfix.KindsDir}
	s.o.Env.StopRun = func(run string) error { stopped = append(stopped, run); return nil }
	return s, root, &stopped
}

func request(t *testing.T, name string) string {
	return string(testfix.Read(t, filepath.Join(Fixtures, "requests", name+".json")))
}

func post(path, body string) call {
	return call{method: http.MethodPost, path: path, header: csrf, body: body}
}

// TestActionFixtures checks each action reply against its fixture, one fresh server per
// case. Requests come from fixtures/api/requests.
func TestActionFixtures(t *testing.T) {
	cases := map[string]struct {
		c    call
		code int
	}{
		"action-dispatch":         {post("/api/missions", request(t, "dispatch")), 201},
		"action-answer":           {post("/api/missions/"+testfix.Waiting+"/answer", request(t, "answer")), 200},
		"action-stop":             {post("/api/missions/"+testfix.Running+"/stop", ""), 200},
		"action-cancel":           {post("/api/missions/"+testfix.Waiting+"/cancel", request(t, "cancel")), 200},
		"action-hide":             {post("/api/inbox/20261109-research-android-push-options/v1/hide", ""), 200},
		"error-answer-stale":      {post("/api/missions/"+testfix.Waiting+"/answer", request(t, "answer-stale")), 409},
		"error-answer-bad-body":   {post("/api/missions/"+testfix.Waiting+"/answer", `{"verb":"answer","version":1,"files":[]}`), 400},
		"error-stop-not-running":  {post("/api/missions/"+testfix.Waiting+"/stop", ""), 409},
		"error-cancel-prefix":     {post("/api/missions/20261110-research-run/cancel", ""), 404},
		"error-hide-unknown":      {post("/api/inbox/20261109-build-nothing/v1/hide", ""), 404},
		"error-dispatch-bad-kind": {post("/api/missions", `{"kind":"nope","project":"demo","goal":"g","end_state":"e"}`), 400},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			s, _, _ := actionServer(t)
			w := do(t, s, tc.c)
			if w.Code != tc.code {
				t.Errorf("%s: status %d, want %d: %s", tc.c.path, w.Code, tc.code, w.Body)
			}
			got := indent(t, w.Body.Bytes())
			p := filepath.Join(Fixtures, name+".json")
			if *update {
				writeFile(t, p, got)
				return
			}
			if want := testfix.Read(t, p); !bytes.Equal(got, want) {
				t.Errorf("%s: reply differs from %s (go test -update to rewrite)\n%s", tc.c.path, p, got)
			}
		})
	}
}

func eventTypes(t *testing.T, root, id string) []string {
	t.Helper()
	var out []string
	for _, l := range strings.Split(strings.TrimSpace(string(testfix.Read(t, filepath.Join(root, id, "events.jsonl")))), "\n") {
		i := strings.Index(l, `"type":"`) + len(`"type":"`)
		out = append(out, l[i:i+strings.IndexByte(l[i:], '"')])
	}
	return out
}

// AC1: an answer from the phone moves the mission out of waiting and writes the text.
func TestAnswerLeavesWaiting(t *testing.T) {
	s, root, _ := actionServer(t)
	if w := do(t, s, post("/api/missions/"+testfix.Waiting+"/answer", `{"text":"skip it","version":1}`)); w.Code != 200 {
		t.Fatalf("answer: %d %s", w.Code, w.Body)
	}
	m, err := s.o.Env.Show(testfix.Waiting)
	if err != nil || m.State == "waiting" {
		t.Errorf("state = %s, %v", m.State, err)
	}
	b, err := os.ReadFile(filepath.Join(root, testfix.Waiting, "answers.v1.md"))
	if err != nil || !strings.Contains(string(b), "skip it") || !strings.Contains(string(b), "surface: nest") {
		t.Errorf("answers.v1.md = %q, %v", b, err)
	}
	// The store sections already show the move.
	w := do(t, s, call{path: "/api/state/inbox"})
	if strings.Contains(w.Body.String(), testfix.Waiting) {
		t.Errorf("inbox still lists the answered row: %s", w.Body)
	}
}

// AC2 and AC3: a stale version gets 409, a POST without X-Crow gets 403; neither writes.
func TestRefusedActionsWriteNothing(t *testing.T) {
	s, root, stopped := actionServer(t)
	log := filepath.Join(root, testfix.Waiting, "events.jsonl")
	before := testfix.Read(t, log)
	for _, c := range []struct {
		c    call
		code int
	}{
		{post("/api/missions/"+testfix.Waiting+"/answer", `{"text":"skip it","version":2}`), 409},
		{call{method: "POST", path: "/api/missions/" + testfix.Waiting + "/answer", body: `{"text":"skip it","version":1}`}, 403},
		{call{method: "POST", path: "/api/missions/" + testfix.Waiting + "/cancel", header: map[string]string{CSRFHeader: "1", "Origin": "https://evil.example"}}, 403},
		{call{method: "POST", path: "/api/missions/" + testfix.Running + "/stop"}, 403},
		{post("/api/missions/"+testfix.Waiting+"/answer", `{"text":"skip it"}`), 400},
	} {
		if w := do(t, s, c.c); w.Code != c.code {
			t.Errorf("%s %s: %d, want %d", c.c.path, c.c.body, w.Code, c.code)
		}
	}
	if !bytes.Equal(testfix.Read(t, log), before) || len(*stopped) != 0 {
		t.Errorf("a refused action wrote: stopped %v", *stopped)
	}
	if _, err := os.Stat(filepath.Join(root, testfix.Waiting, "answers.v1.md")); err == nil {
		t.Error("a refused answer wrote answers.v1.md")
	}
}

// Stop stops the live run and records it; cancel and hide record their action first.
func TestActionEvents(t *testing.T) {
	s, root, stopped := actionServer(t)
	for _, c := range []call{
		post("/api/missions/"+testfix.Running+"/stop", ""),
		post("/api/missions/"+testfix.Waiting+"/cancel", `{"reason":"changed my mind"}`),
		post("/api/inbox/20261109-research-android-push-options/v1/hide", "{}"),
	} {
		if w := do(t, s, c); w.Code != 200 {
			t.Fatalf("%s: %d %s", c.path, w.Code, w.Body)
		}
	}
	if !slices.Equal(*stopped, []string{testfix.LiveRun}) {
		t.Errorf("stopped = %v", *stopped)
	}
	for id, tail := range map[string]string{
		testfix.Running:                          "run_start,action",
		testfix.Waiting:                          "state,action,state",
		"20261109-research-android-push-options": "state,action",
	} {
		if got := strings.Join(eventTypes(t, root, id), ","); !strings.HasSuffix(got, tail) {
			t.Errorf("%s events = %s, want ...%s", id, got, tail)
		}
	}
	// Hidden rows leave the open inbox.
	if w := do(t, s, call{path: "/api/state/inbox"}); strings.Contains(w.Body.String(), "android-push-options/v1") {
		t.Errorf("hidden row still open: %s", w.Body)
	}
}
