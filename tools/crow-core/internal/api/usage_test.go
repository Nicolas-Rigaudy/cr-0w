package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

func usageData() map[string]any {
	return map[string]any{
		"generated_at": "2026-11-10T09:00:00+00:00",
		"limits":       []any{map[string]any{"account": "acct1", "remaining_5h_pct": 40.0, "remaining_week_pct": 70.0}},
		"today":        map[string]any{"cost_usd": 12.0},
		"sessions": []any{
			map[string]any{"id": "s1", "account": "acct1", "project": "api", "kind": "run", "cost_usd": 3.5,
				"tokens": 1000.0, "path": "/private/transcript.jsonl", "cwd": "/private", "activity": []any{1.0}},
			"not a row",
		},
	}
}

func costData() map[string]any {
	return map[string]any{"days": 35.0, "cost_usd": 42.0,
		"by_account": map[string]any{"acct1": map[string]any{"cost_usd": 42.0}},
		"by_day":     []any{map[string]any{"date": "2026-11-10", "cost_usd": 42.0, "limits": nil}}}
}

func usageServer(clk *fakeClock, costFail *atomic.Bool) *Server {
	return New(Options{Now: clk.Now,
		Sources: []Source{{Name: "usage", TTL: 30 * time.Second, Fetch: func() (any, error) { return usageData(), nil }}},
		Costs: func() (any, error) {
			if costFail != nil && costFail.Load() {
				return nil, errors.New("crow-cost-audit: exit status 1")
			}
			return costData(), nil
		}})
}

func getUsage(t *testing.T, s *Server, header map[string]string) (Usage, *http.Response) {
	t.Helper()
	w := do(t, s, call{path: "/api/usage", trusted: true, header: header})
	var u Usage
	if w.Code == http.StatusOK {
		if err := json.Unmarshal(w.Body.Bytes(), &u); err != nil {
			t.Fatal(err)
		}
	}
	return u, w.Result()
}

func TestUsageEndpoint(t *testing.T) {
	clk := &fakeClock{t: clock}
	s := usageServer(clk, nil)

	cold, res := getUsage(t, s, nil)
	if res.StatusCode != 200 || !cold.Stale || !cold.Cost.Stale || len(cold.Sessions) != 0 || len(cold.Limits.([]any)) != 0 {
		t.Fatalf("cold = %d %+v, want empty and stale", res.StatusCode, cold)
	}
	idle(t, s, "usage")
	idle(t, s, UsageCostSection)

	u, res := getUsage(t, s, nil)
	if u.Stale || u.Cost.Stale || u.GeneratedAt != "2026-11-10T09:00:00+00:00" {
		t.Fatalf("warm = %+v", u)
	}
	if got := u.Cost.Data.(map[string]any)["cost_usd"]; got != 42.0 {
		t.Fatalf("cost = %v", got)
	}
	if l := u.Limits.([]any); len(l) != 1 || l[0].(map[string]any)["account"] != "acct1" {
		t.Fatalf("limits = %v", u.Limits)
	}
	if len(u.Sessions) != 1 {
		t.Fatalf("sessions = %v, want the one object row", u.Sessions)
	}
	row := u.Sessions[0].(map[string]any)
	keys := make([]string, 0, len(row))
	for k := range row {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	if want := []string{"account", "cost_usd", "id", "kind", "project", "tokens"}; !reflect.DeepEqual(keys, want) {
		t.Fatalf("session keys = %v, want %v (path, cwd, activity dropped)", keys, want)
	}

	etag := res.Header.Get("ETag")
	if _, res := getUsage(t, s, map[string]string{"If-None-Match": etag}); res.StatusCode != http.StatusNotModified {
		t.Fatalf("unchanged: %d, want 304", res.StatusCode)
	}
}

// A failed audit keeps the last good figures, marked stale, and moves the ETag.
func TestUsageCostStale(t *testing.T) {
	clk := &fakeClock{t: clock}
	var bad atomic.Bool
	s := usageServer(clk, &bad)
	getUsage(t, s, nil)
	idle(t, s, UsageCostSection)
	good, res := getUsage(t, s, nil)
	etag := res.Header.Get("ETag")

	bad.Store(true)
	clk.Add(UsageCostTTL + time.Second)
	getUsage(t, s, nil)
	idle(t, s, UsageCostSection)
	stale, res := getUsage(t, s, map[string]string{"If-None-Match": etag})
	if res.StatusCode != 200 || !stale.Cost.Stale || stale.Cost.Error == "" || !reflect.DeepEqual(stale.Cost.Data, good.Cost.Data) {
		t.Fatalf("after a failed audit = %d %+v", res.StatusCode, stale.Cost)
	}
}

// The audit is no state section: /api/state neither lists nor runs it.
func TestUsageCostNotInState(t *testing.T) {
	clk := &fakeClock{t: clock}
	var n atomic.Int32
	s := New(Options{Now: clk.Now, Costs: func() (any, error) { n.Add(1); return costData(), nil }})
	if slices.Contains(s.Sections(), UsageCostSection) {
		t.Fatalf("sections = %v", s.Sections())
	}
	if w := do(t, s, call{path: "/api/state/" + UsageCostSection, trusted: true}); w.Code != http.StatusNotFound {
		t.Fatalf("state/%s = %d, want 404", UsageCostSection, w.Code)
	}
	s.tick()
	time.Sleep(10 * time.Millisecond)
	if n.Load() != 0 {
		t.Fatalf("tick ran the audit %d times", n.Load())
	}
	// Without the usage source the route answers 404.
	if w := do(t, s, call{path: "/api/usage", trusted: true}); w.Code != http.StatusNotFound {
		t.Fatalf("usage without source = %d", w.Code)
	}
}
