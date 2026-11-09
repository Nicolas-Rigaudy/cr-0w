// Package testfix gives the core tests the shared spec fixtures as mission dirs.
package testfix

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Repo is the hub checkout, found from this file's path.
func Repo() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..", "..", "..")
}

// Spec paths.
var (
	CoreSpec = filepath.Join(Repo(), "logics", "specs", "core_contract")
	RunSpec  = filepath.Join(Repo(), "logics", "specs", "req_037_run_contract")
	KindsDir = filepath.Join(RunSpec, "kinds")
	Records  = filepath.Join(CoreSpec, "fixtures", "records")
)

// Logs maps each fixture mission to its event log in the specs.
func Logs() map[string]string {
	return map[string]string{
		"20261109-research-android-push-options": filepath.Join(CoreSpec, "fixtures", "events", "20261109-research-android-push-options.jsonl"),
		"20261109-build-item-095":                filepath.Join(RunSpec, "examples", "20261109-build-item-095", "events.jsonl"),
	}
}

// Read reads a file or fails the test.
func Read(t testing.TB, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// MissionsRoot copies every fixture log into a fresh missions root.
func MissionsRoot(t testing.TB) string {
	t.Helper()
	root := t.TempDir()
	for id, src := range Logs() {
		dir := filepath.Join(root, id)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), Read(t, src), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
