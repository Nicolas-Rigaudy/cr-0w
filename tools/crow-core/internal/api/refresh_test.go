package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"crow-core/internal/testfix"
)

// fakeClock is a settable clock for the refresher tests.
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *fakeClock) Add(d time.Duration) {
	c.mu.Lock()
	c.t = c.t.Add(d)
	c.mu.Unlock()
}

// idle waits until no refresh of name is running.
func idle(t *testing.T, s *Server, name string) {
	t.Helper()
	e := s.ext[name]
	for i := 0; i < 500; i++ {
		e.mu.Lock()
		r := e.running
		e.mu.Unlock()
		if !r {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("%s: refresh still running", name)
}

func sectionOf(t *testing.T, s *Server, name string, header map[string]string) (Section, *http.Response) {
	t.Helper()
	w := do(t, s, call{path: "/api/state/" + name, header: header, trusted: true})
	if w.Code != http.StatusOK && w.Code != http.StatusNotModified {
		t.Fatalf("%s: %d %s", name, w.Code, w.Body)
	}
	var sec Section
	if w.Code == http.StatusOK {
		if err := json.Unmarshal(w.Body.Bytes(), &sec); err != nil {
			t.Fatal(err)
		}
	}
	return sec, w.Result()
}

// TestScriptSections serves each script's JSON unchanged (AC1).
func TestScriptSections(t *testing.T) {
	s, _ := newServer(t)
	for _, p := range scriptPayloads {
		var want any
		if err := json.Unmarshal(testfix.Read(t, filepath.Join(testfix.ScriptSpec, "fixtures", p.fixture+".json")), &want); err != nil {
			t.Fatal(err)
		}
		sec, res := sectionOf(t, s, p.section, nil)
		if res.StatusCode != 200 || sec.Stale || !reflect.DeepEqual(sec.Data, want) {
			t.Errorf("%s: %d stale=%v, data differs from %s.json", p.section, res.StatusCode, sec.Stale, p.fixture)
		}
	}
}

// TestSlowScript keeps every other section fast while one script hangs (AC2).
func TestSlowScript(t *testing.T) {
	s, _ := newServer(t)
	release := make(chan struct{})
	defer close(release)
	s.ext["mrs"].src.Fetch = func() (any, error) { <-release; return []any{}, nil }
	s.ext["mrs"].ran = time.Time{} // due now
	do(t, s, call{path: "/api/state/mrs"})
	if !s.ext["mrs"].running {
		t.Fatal("mrs refresh did not start")
	}
	for _, p := range []string{"/api/state/missions", "/api/state/sessions", "/api/state/mrs", "/api/state"} {
		start := time.Now()
		if w := do(t, s, call{path: p}); w.Code != 200 {
			t.Fatalf("%s: %d", p, w.Code)
		}
		if d := time.Since(start); d > 100*time.Millisecond {
			t.Errorf("%s took %s while mrs hangs", p, d)
		}
	}
}

// TestFailingScript keeps the last good data, flags it stale with the error and the old
// updated_at, and clears on the next good run (AC3). It answers 304 while unchanged (AC4).
func TestFailingScript(t *testing.T) {
	clk := &fakeClock{t: clock}
	var fail atomic.Bool
	var n atomic.Int32
	s := New(Options{Now: clk.Now, Sources: []Source{{Name: "dispatch", TTL: 10 * time.Second, Fetch: func() (any, error) {
		n.Add(1)
		if fail.Load() {
			return nil, errors.New("crow-runs: exit status 2: cdx not found")
		}
		return map[string]any{"schema": 1.0, "live": []any{}}, nil
	}}}})

	cold, _ := sectionOf(t, s, "dispatch", nil)
	if !cold.Stale || cold.UpdatedAt != nil || !reflect.DeepEqual(cold.Data, []any{}) {
		t.Fatalf("first load = %+v, want empty and stale", cold)
	}
	idle(t, s, "dispatch")
	good, res := sectionOf(t, s, "dispatch", nil)
	if good.Stale || good.Error != "" || good.UpdatedAt == nil {
		t.Fatalf("after a run = %+v", good)
	}
	etag := res.Header.Get("ETag")
	if _, res := sectionOf(t, s, "dispatch", map[string]string{"If-None-Match": etag}); res.StatusCode != http.StatusNotModified {
		t.Fatalf("unchanged: %d, want 304", res.StatusCode)
	}

	fail.Store(true)
	clk.Add(11 * time.Second)
	sectionOf(t, s, "dispatch", nil) // starts the failing run
	idle(t, s, "dispatch")
	bad, res := sectionOf(t, s, "dispatch", map[string]string{"If-None-Match": etag})
	if res.StatusCode != 200 || !bad.Stale || bad.Error == "" || *bad.UpdatedAt != *good.UpdatedAt || !reflect.DeepEqual(bad.Data, good.Data) {
		t.Fatalf("after a failure = %d %+v", res.StatusCode, bad)
	}

	fail.Store(false)
	clk.Add(11 * time.Second)
	sectionOf(t, s, "dispatch", nil)
	idle(t, s, "dispatch")
	back, res := sectionOf(t, s, "dispatch", nil)
	if back.Stale || back.Error != "" || res.Header.Get("ETag") != etag {
		t.Fatalf("after recovery = %+v %s, want the first ETag", back, res.Header.Get("ETag"))
	}
	if n.Load() != 3 {
		t.Errorf("runs = %d, want 3", n.Load())
	}
}

// TestIdleSection refreshes only sections requested in the last IdleAfter (AC5).
func TestIdleSection(t *testing.T) {
	clk := &fakeClock{t: clock}
	var n atomic.Int32
	s := New(Options{Now: clk.Now, Sources: []Source{{Name: "sessions", TTL: 5 * time.Second, Fetch: func() (any, error) {
		n.Add(1)
		return []any{}, nil
	}}}})
	run := func() { s.tick(); idle(t, s, "sessions") }

	run()
	if n.Load() != 0 {
		t.Fatalf("refreshed before any request: %d", n.Load())
	}
	sectionOf(t, s, "sessions", nil)
	idle(t, s, "sessions")
	for i := 0; i < 23; i++ { // 115 s of ticks, one run per TTL
		clk.Add(5 * time.Second)
		run()
	}
	if got := n.Load(); got != 24 {
		t.Fatalf("runs while requested = %d, want 24", got)
	}
	for i := 0; i < 20; i++ { // past IdleAfter
		clk.Add(5 * time.Second)
		run()
	}
	if got := n.Load(); got != 24 {
		t.Fatalf("runs after 2 min idle = %d, want 24", got)
	}
	sectionOf(t, s, "sessions", nil)
	idle(t, s, "sessions")
	if got := n.Load(); got != 25 {
		t.Fatalf("a new request did not refresh: %d", got)
	}
}
