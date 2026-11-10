package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"crow-core/internal/api"
	"crow-core/internal/mission"
	"crow-core/internal/project"
	"crow-core/internal/store"
	"crow-core/internal/testfix"
)

func liveEnv(t *testing.T) (mission.Env, string) {
	t.Helper()
	root := t.TempDir()
	testfix.Live(t, root)
	kinds := project.KindDirs(testfix.KindsDir, "")
	s, err := store.Open(filepath.Join(t.TempDir(), "core.db"), root, kinds)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	now := time.Date(2026, 11, 10, 9, 0, 0, 0, time.FixedZone("CET", 3600))
	return mission.Env{Store: s, Kinds: kinds, Projects: testfix.Registry(t), Now: func() time.Time { return now }}, root
}

// item_147 AC4: `crow inbox answer` and POST /api/missions/<id>/answer write the same
// events and the same answers file; only the surface differs.
func TestInboxAnswerMatchesTheEndpoint(t *testing.T) {
	cliEnv, cliRoot := liveEnv(t)
	if err := inboxCmd(cliEnv, []string{"answer", testfix.Waiting + "/v1/q1", "skip", "it"}, io.Discard); err != nil {
		t.Fatal(err)
	}
	apiEnv, apiRoot := liveEnv(t)
	r := httptest.NewRequest(http.MethodPost, "/api/missions/"+testfix.Waiting+"/answer",
		strings.NewReader(`{"item":"v1/q1","text":"skip it","version":1}`))
	r.Header.Set(api.CSRFHeader, "1")
	w := httptest.NewRecorder()
	api.New(api.Options{Env: apiEnv}).Handler(true).ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("endpoint: %d %s", w.Code, w.Body)
	}
	for _, f := range []string{"events.jsonl", "answers.v1.md"} {
		a := read(t, filepath.Join(cliRoot, testfix.Waiting, f))
		b := read(t, filepath.Join(apiRoot, testfix.Waiting, f))
		if a = strings.ReplaceAll(a, `"surface":"cli","verb"`, `"surface":"nest","verb"`); f == "answers.v1.md" {
			a = strings.Replace(a, "surface: cli", "surface: nest", 1)
		}
		if a != b {
			t.Errorf("%s differs\ncli:\n%s\napi:\n%s", f, a, b)
		}
	}
}

func read(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
