package main

import (
	"strings"
	"testing"
)

func lineText(l line) string {
	var b strings.Builder
	for _, s := range l {
		b.WriteString(s.text)
	}
	return b.String()
}

// A long reply wraps over several rows on a phone-width box instead of ending in "…".
func TestTailLineWrapsALongMessage(t *testing.T) {
	av := agentsView{lk: newLook(themes[0], -1)}
	msg := "Found it: two cuts. crow-runs keeps only the first line of each message, then the deck draws each entry as one row clipped at the box edge."
	rows := av.tailLine(tailLine{Text: "says: " + msg}, false, false, 40)
	if len(rows) < 3 {
		t.Fatalf("want the message over several rows, got %d", len(rows))
	}
	var got []string
	for _, r := range rows {
		if r.width() > 40 {
			t.Fatalf("row wider than the box: %q", lineText(r))
		}
		got = append(got, strings.TrimSpace(lineText(r)))
	}
	if all := strings.Join(got, " "); !strings.HasSuffix(all, "clipped at the box edge.") {
		t.Fatalf("message cut: %q", all)
	}
}

func TestWrapFirstUsesTheNarrowerFirstRow(t *testing.T) {
	got := wrapFirst("aa bb cc dd", 2, 5)
	want := []string{"aa", "bb cc", "dd"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("got %q want %q", got, want)
	}
}
