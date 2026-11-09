package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"crow-core/contract"
	"crow-core/internal/eventlog"
	"crow-core/internal/project"
	"crow-core/internal/testfix"
)

func open(t *testing.T, path, root string) *Store {
	t.Helper()
	s, err := Open(path, root, project.KindDirs(testfix.KindsDir, ""))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func dump(t *testing.T, s *Store) string {
	t.Helper()
	d, err := s.Dump()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(d, "", "  ")
	return string(b)
}

// rebuilt is the dump of a fresh store rebuilt from the fixture logs.
func rebuilt(t *testing.T) string {
	t.Helper()
	root := testfix.MissionsRoot(t)
	s := open(t, filepath.Join(t.TempDir(), "core.db"), root)
	if _, err := s.Rebuild(); err != nil {
		t.Fatal(err)
	}
	return dump(t, s)
}

// AC1: appending the events one by one, syncing after each, gives the rows of a rebuild.
func TestOneByOneEqualsRebuild(t *testing.T) {
	src := testfix.MissionsRoot(t)
	root := t.TempDir()
	s := open(t, filepath.Join(root, "core.db"), root)
	ids, err := eventlog.Missions(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		recs, err := eventlog.Read(eventlog.Path(src, id), 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range recs {
			end, err := eventlog.Append(root, r.Event)
			if err != nil {
				t.Fatal(err)
			}
			if n, err := s.Sync(id); err != nil || n != 1 {
				t.Fatalf("sync after one append: %d events, err %v", n, err)
			}
			if got, _ := s.Applied(id); got != end {
				t.Fatalf("applied %d, log ends at %d", got, end)
			}
		}
	}
	if n, err := s.SyncAll(); err != nil || n != 0 {
		t.Errorf("sync with nothing new applied %d, err %v", n, err)
	}
	if got, want := dump(t, s), rebuilt(t); got != want {
		t.Errorf("one by one\n%s\nrebuild\n%s", got, want)
	}
}

// AC3: deleting the store and rebuilding brings every row back.
func TestDeletedStoreRebuilds(t *testing.T) {
	root := testfix.MissionsRoot(t)
	db := filepath.Join(root, "core.db")
	s := open(t, db, root)
	if _, err := s.SyncAll(); err != nil {
		t.Fatal(err)
	}
	before := dump(t, s)
	s.Close()
	for _, p := range []string{db, db + "-wal", db + "-shm"} {
		if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	}
	s = open(t, db, root)
	if _, err := s.Rebuild(); err != nil {
		t.Fatal(err)
	}
	if after := dump(t, s); after != before {
		t.Errorf("before\n%s\nafter rebuild\n%s", before, after)
	}
	var d Dump
	json.Unmarshal([]byte(before), &d)
	if len(d.Missions) != 2 || len(d.Runs) != 4 || len(d.Inbox) != 4 {
		t.Errorf("rows: %d missions, %d runs, %d inbox", len(d.Missions), len(d.Runs), len(d.Inbox))
	}
}

func TestRebuildTwiceIsStable(t *testing.T) {
	root := testfix.MissionsRoot(t)
	s := open(t, filepath.Join(root, "core.db"), root)
	n1, err := s.Rebuild()
	if err != nil {
		t.Fatal(err)
	}
	first := dump(t, s)
	n2, err := s.Rebuild()
	if err != nil || n1 != n2 || dump(t, s) != first {
		t.Errorf("second rebuild differs: %d vs %d events, err %v", n1, n2, err)
	}
}

func TestShrunkLogNeedsARebuild(t *testing.T) {
	root := testfix.MissionsRoot(t)
	s := open(t, filepath.Join(root, "core.db"), root)
	if _, err := s.SyncAll(); err != nil {
		t.Fatal(err)
	}
	id := "20261109-research-android-push-options"
	if err := os.Truncate(eventlog.Path(root, id), 10); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Sync(id); !errors.Is(err, eventlog.ErrShrunk) {
		t.Errorf("want ErrShrunk, got %v", err)
	}
}

func TestSchemaMismatchIsRefused(t *testing.T) {
	root := t.TempDir()
	db := filepath.Join(root, "core.db")
	s := open(t, db, root)
	if _, err := s.db.Exec(`UPDATE meta SET value = '0' WHERE key = 'schema'`); err != nil {
		t.Fatal(err)
	}
	s.Close()
	if _, err := Open(db, root, nil); err == nil {
		t.Error("a store of another schema opened")
	}
}

func TestAppendRefusesAnEventThatDoesNotApply(t *testing.T) {
	root := t.TempDir()
	s := open(t, filepath.Join(root, "core.db"), root)
	ev := func(typ, data string) contract.Event {
		return contract.Event{TS: "2026-11-09T14:00:00+01:00", Mission: "20261109-research-x", Type: typ, Actor: "core", Data: json.RawMessage(data)}
	}
	if _, err := s.Append(ev("state", `{"from":"preflight","to":"queued"}`)); err == nil {
		t.Fatal("state before created accepted")
	}
	if _, err := os.Stat(eventlog.Path(root, "20261109-research-x")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("a refused first event created the log")
	}
	queued := ev("state", `{"from":"preflight","to":"queued"}`)
	for _, e := range []contract.Event{ev("created", `{"kind":"research","batch":null}`), queued} {
		if _, err := s.Append(e); err != nil {
			t.Fatal(err)
		}
	}
	size := func() int64 { st, _ := os.Stat(eventlog.Path(root, "20261109-research-x")); return st.Size() }
	before := size()
	if _, err := s.Append(queued); err == nil {
		t.Fatal("stale state move accepted")
	}
	if size() != before {
		t.Error("a refused event reached the log")
	}
	if _, err := s.Rebuild(); err != nil {
		t.Errorf("log no longer replays: %v", err)
	}
}
