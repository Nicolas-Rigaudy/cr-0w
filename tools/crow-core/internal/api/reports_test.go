package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildReportsSettledAFKOnly(t *testing.T) {
	root := t.TempDir()
	afk := filepath.Join(root, "afk")
	// No project: settled, nothing ran.
	write(t, filepath.Join(afk, "summary-2026-11-23.json"), `{"start":"2026-11-22T18:00","projects":{}}`)
	// One open decision: held for Inbox.
	write(t, filepath.Join(afk, "summary-2026-11-24.json"), `{"start":"2026-11-23T18:00","end":"2026-11-24T06:00",
		"projects":[{"name":"nest","decisions":[{"q":"Pick a layout"}],"decided":[]}]}`)
	// Every decision is in decided[]: settled.
	write(t, filepath.Join(afk, "summary-2026-11-25.json"), `{"start":"2026-11-24T18:08","end":"2026-11-25T06:30",
		"projects":{"cr-0w":{"tldr":"Inbox shipped","points":[{"title":"Inbox list","text":"Rows group by kind.","runs":["a1","b2"]}],
		"decisions":[{"q":"Pick a layout"}],"decided":["Pick a layout: B"]}}}`)
	write(t, filepath.Join(afk, "summary-bad.json"), `not json`)
	r, err := BuildReports(ReportDirs{AFK: afk})
	if err != nil {
		t.Fatal(err)
	}
	if r.Held != 1 || len(r.Items) != 2 {
		t.Fatalf("held %d items %d, want 1 and 2", r.Held, len(r.Items))
	}
	got := r.Items[0]
	if got.ID != "afk:2026-11-25" || got.At != "2026-11-25T06:30" || got.Summary != "1 points · 2 runs · 1 decisions answered" {
		t.Errorf("newest = %+v", got)
	}
	if !strings.Contains(got.Body, "## cr-0w") || !strings.Contains(got.Body, "- **Inbox list**: Rows group by kind.") || !strings.Contains(got.Body, "- Pick a layout: B") {
		t.Errorf("body = %q", got.Body)
	}
	if r.Items[1].Summary != "0 points · 0 runs · no decision" || !strings.Contains(r.Items[1].Body, "No project ran.") {
		t.Errorf("empty night = %+v", r.Items[1])
	}
}

func TestBuildReportsNotes(t *testing.T) {
	vault := filepath.Join(t.TempDir(), "My Vault")
	crow := filepath.Join(vault, "CR-0W")
	roll := filepath.Join(crow, "Rollups", "2026-W47 rollup.md")
	retro := filepath.Join(crow, "Retros", "2026-11-24 batch.md")
	write(t, roll, "# Week 47 rollup\n\n- Slides drafted\n")
	write(t, retro, "---\ndate: 2026-11-24\n---\n# Retro batch\n\nFour lessons applied.\n")
	old := time.Date(2026, 11, 20, 9, 0, 0, 0, time.UTC)
	if err := os.Chtimes(roll, old, old); err != nil {
		t.Fatal(err)
	}
	r, err := BuildReports(ReportDirs{AFK: filepath.Join(vault, "none"), Rollups: filepath.Dir(roll), Retros: filepath.Dir(retro), Vault: vault})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Items) != 2 || r.Items[0].Kind != "retro" || r.Items[1].Kind != "weekly" {
		t.Fatalf("items = %+v", r.Items)
	}
	if r.Items[0].Title != "Retro batch" || r.Items[0].Summary != "Four lessons applied." {
		t.Errorf("retro = %+v", r.Items[0])
	}
	if r.Items[1].Title != "Week 47 rollup" || r.Items[1].Summary != "Slides drafted" || r.Items[1].ID != "weekly:2026-W47 rollup" {
		t.Errorf("weekly = %+v", r.Items[1])
	}
	want := "obsidian://open?vault=My%20Vault&file=CR-0W%2FRollups%2F2026-W47%20rollup"
	if r.Items[1].Link != want {
		t.Errorf("link = %q, want %q", r.Items[1].Link, want)
	}
}

func TestBuildReportsMissingDirs(t *testing.T) {
	r, err := BuildReports(ReportDirs{AFK: "/nonexistent/afk", Rollups: "/nonexistent/r"})
	if err != nil || len(r.Items) != 0 || r.Items == nil {
		t.Fatalf("got %+v, %v", r, err)
	}
}
