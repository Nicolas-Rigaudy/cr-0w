package api

import (
	"fmt"
	"strings"
	"testing"

	"crow-core/internal/testfix"
)

// TestRequestLog logs every write under /api with its status, refused ones too, and
// skips reads.
func TestRequestLog(t *testing.T) {
	s, _, _ := actionServer(t)
	var lines []string
	s.o.Logf = func(f string, a ...any) { lines = append(lines, fmt.Sprintf(f, a...)) }
	do(t, s, call{path: "/api/state"})
	do(t, s, post("/api/missions/"+testfix.Running+"/stop", ""))
	do(t, s, call{method: "POST", path: "/api/missions/" + testfix.Waiting + "/stop"})
	want := []string{
		"api POST /api/missions/" + testfix.Running + "/stop 200 ",
		"api POST /api/missions/" + testfix.Waiting + "/stop 403 ",
	}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if !strings.HasPrefix(lines[i], w) || !strings.HasSuffix(lines[i], "(tcp)") {
			t.Errorf("line %d = %q, want prefix %q", i, lines[i], w)
		}
	}
}
