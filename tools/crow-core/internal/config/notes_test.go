package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNotePath(t *testing.T) {
	home := t.TempDir()
	c := Config{Home: home, Engine: repoRoot,Projects: filepath.Join(home, "local", "projects.yaml")}
	if got := NotePath(c, "retros"); got != filepath.Join(home, "local", "state", "retros") {
		t.Errorf("no registry: %s", got)
	}
	if err := os.MkdirAll(filepath.Dir(c.Projects), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(c.Projects, []byte("vault:\n  path: /v/Vault\n  crow_dir: CR-0W\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := NotePath(c, "rollups"); got != "/v/Vault/CR-0W/Rollups" {
		t.Errorf("relative crow_dir: %s", got)
	}
	if _, vault := NotesDir(c); vault != "/v/Vault" {
		t.Errorf("vault = %s", vault)
	}
	if err := os.WriteFile(filepath.Join(home, "local", "crow.yaml"), []byte("version: 1\nmodules:\n  vault: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := NotePath(c, "retros"); got != filepath.Join(home, "notes", "Retros") {
		t.Errorf("vault off: %s", got)
	}
}
