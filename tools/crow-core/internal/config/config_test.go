package config

import (
	"os"
	"path/filepath"
	"testing"
)

func env(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestDefaults(t *testing.T) {
	home := t.TempDir()
	c, err := Load(env(map[string]string{"CROW_HOME": home, "CROW_ENGINE": "/engine"}))
	if err != nil {
		t.Fatal(err)
	}
	if c.Missions != filepath.Join(home, "local/state/missions") || c.Store != filepath.Join(c.Missions, "core.db") {
		t.Errorf("paths: %+v", c)
	}
	if c.Prefs != filepath.Join(home, "local/mission-prefs.yaml") {
		t.Errorf("prefs: %s", c.Prefs)
	}
	if c.Kinds != "/engine/logics/specs/req_037_run_contract/kinds" {
		t.Errorf("kinds: %s", c.Kinds)
	}
}

func TestFileOverridesAndRefusesUnknownKeys(t *testing.T) {
	home := t.TempDir()
	os.MkdirAll(filepath.Join(home, "local"), 0o700)
	f := filepath.Join(home, "local/core.yaml")
	os.WriteFile(f, []byte("store: db/core.db\n"), 0o600)
	c, err := Load(env(map[string]string{"CROW_HOME": home}))
	if err != nil || c.Store != filepath.Join(home, "db/core.db") {
		t.Errorf("store %s, err %v", c.Store, err)
	}
	os.WriteFile(f, []byte("stor: x\n"), 0o600)
	if _, err := Load(env(map[string]string{"CROW_HOME": home})); err == nil {
		t.Error("unknown key accepted")
	}
}

func TestNeedsHome(t *testing.T) {
	if _, err := Load(env(nil)); err == nil {
		t.Error("loaded without CROW_HOME")
	}
}
