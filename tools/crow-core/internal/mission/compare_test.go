package mission

import (
	"os"
	"path/filepath"
	"testing"
)

// item_211: a crow-compare spec in the outbox gives the review its labels and scenarios.
func TestCompareOfReadsTheFirstDecision(t *testing.T) {
	dir := t.TempDir()
	box := filepath.Join(dir, "outbox", "boot")
	if err := os.MkdirAll(box, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"a-frame.png", "a-dark.png", "b-frame.png"} {
		if err := os.WriteFile(filepath.Join(box, f), []byte("png"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	spec := `{"title": "Boot", "decisions": [{"id": "d1", "question": "Which retry banner?", "recommend": "Quiet",
	  "rows": ["boot frame", "dark room"],
	  "options": [{"label": "Quiet", "note": "one line", "images": ["a-frame.png", "a-dark.png"]},
	              {"label": "Loud", "images": ["b-frame.png", "missing.png"]},
	              {"label": "Escape", "images": ["../../events.jsonl", "/etc/passwd"]}]},
	  {"id": "d2", "question": "ignored", "options": []}]}`
	if err := os.WriteFile(filepath.Join(box, "spec.json"), []byte(spec), 0o600); err != nil {
		t.Fatal(err)
	}
	c := compareOf(dir, []string{"outbox/notes.md", "outbox/boot/spec.json"})
	if c == nil {
		t.Fatal("no compare")
	}
	if c.Question != "Which retry banner?" || c.Recommend != "Quiet" || len(c.Rows) != 2 || len(c.Options) != 3 {
		t.Fatalf("compare = %+v", c)
	}
	got := func(o, r int) string {
		if p := c.Options[o].Images[r]; p != nil {
			return *p
		}
		return ""
	}
	if got(0, 0) != "outbox/boot/a-frame.png" || got(0, 1) != "outbox/boot/a-dark.png" || got(1, 0) != "outbox/boot/b-frame.png" {
		t.Errorf("images = %+v", c.Options)
	}
	if got(1, 1) != "" || got(2, 0) != "" || got(2, 1) != "" {
		t.Errorf("missing or escaping images kept: %+v", c.Options)
	}
	if c.Options[0].Note != "one line" || c.Options[1].Note != "" {
		t.Errorf("notes = %+v", c.Options)
	}
	if n := len(compareImages(c)); n != 3 {
		t.Errorf("compareImages = %d, want 3", n)
	}
	if compareOf(dir, []string{"outbox/boot/a-frame.png"}) != nil {
		t.Error("a non-spec artifact gave a compare")
	}
}
