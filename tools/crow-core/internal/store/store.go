// Package store keeps the SQLite view of the mission event logs: one row per mission,
// run and inbox item. The view is disposable: Rebuild drops it and replays every log,
// and Sync applies only the tail of a log past the offset the store last applied.
package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // pure Go, no cgo

	"crow-core/contract"
	"crow-core/internal/eventlog"
	"crow-core/internal/project"
)

// SchemaVersion changes when the tables change. A store of another version is rebuilt.
const SchemaVersion = 1

const schema = `
CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS missions (
  id TEXT PRIMARY KEY, state TEXT NOT NULL, kind TEXT NOT NULL, project TEXT NOT NULL,
  updated TEXT NOT NULL, applied INTEGER NOT NULL, record TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS runs (
  run TEXT PRIMARY KEY, mission TEXT NOT NULL, started TEXT NOT NULL, ended TEXT, record TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS runs_mission ON runs(mission);
CREATE TABLE IF NOT EXISTS inbox (
  id TEXT PRIMARY KEY, mission TEXT NOT NULL, state TEXT NOT NULL, created TEXT NOT NULL, record TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS inbox_mission ON inbox(mission);
`

// Store is an open store over one missions root.
type Store struct {
	db    *sql.DB
	root  string
	kinds project.Kinds
}

// Open opens (or creates) the store at path for the logs under root. Writes take the
// database lock at BEGIN, so two processes syncing the same mission run one after the other.
func Open(path, root string, kinds project.Kinds) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	dsn := "file:" + path + "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_txlock=immediate"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, root: root, kinds: kinds}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the database.
func (s *Store) Close() error { return s.db.Close() }

// Root is the missions root: one dir per mission, holding its log and files.
func (s *Store) Root() string { return s.root }

func (s *Store) migrate() error {
	if _, err := s.db.Exec(schema); err != nil {
		return err
	}
	var v string
	err := s.db.QueryRow(`SELECT value FROM meta WHERE key = 'schema'`).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = s.db.Exec(`INSERT INTO meta(key, value) VALUES ('schema', ?)`, fmt.Sprint(SchemaVersion))
		return err
	}
	if err != nil {
		return err
	}
	if v != fmt.Sprint(SchemaVersion) {
		return fmt.Errorf("store schema %s, core wants %d: delete the store file and run crow-core rebuild", v, SchemaVersion)
	}
	return nil
}

// Sync applies the events of one mission past its applied offset. It returns how many
// events it applied.
func (s *Store) Sync(mission string) (int, error) {
	if err := eventlog.ValidMission(mission); err != nil {
		return 0, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // no-op after commit
	n, err := s.syncTx(tx, mission)
	if err != nil {
		return 0, err
	}
	return n, tx.Commit()
}

// SyncAll syncs every mission under the root and returns the events applied.
func (s *Store) SyncAll() (int, error) {
	ids, err := eventlog.Missions(s.root)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, id := range ids {
		n, err := s.Sync(id)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

// Rebuild empties every table and replays every log, in one transaction.
func (s *Store) Rebuild() (int, error) {
	ids, err := eventlog.Missions(s.root)
	if err != nil {
		return 0, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // no-op after commit
	for _, t := range []string{"missions", "runs", "inbox"} {
		if _, err := tx.Exec("DELETE FROM " + t); err != nil {
			return 0, err
		}
	}
	total := 0
	for _, id := range ids {
		n, err := s.syncTx(tx, id)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, tx.Commit()
}

func (s *Store) syncTx(tx *sql.Tx, mission string) (int, error) {
	v, applied, err := current(tx, mission)
	if err != nil {
		return 0, err
	}
	recs, err := eventlog.Read(eventlog.Path(s.root, mission), applied)
	if err != nil || len(recs) == 0 {
		return 0, err
	}
	for _, r := range recs {
		if err := project.Apply(&v, r.Event, s.kinds); err != nil {
			return 0, fmt.Errorf("%s at %d: %w", eventlog.Path(s.root, mission), r.End, err)
		}
	}
	return len(recs), save(tx, v, recs[len(recs)-1].End)
}

// current is the stored view of a mission and the log offset it reflects.
func current(tx *sql.Tx, mission string) (project.View, int64, error) {
	var applied int64
	var record string
	err := tx.QueryRow(`SELECT applied, record FROM missions WHERE id = ?`, mission).Scan(&applied, &record)
	if errors.Is(err, sql.ErrNoRows) {
		return project.View{}, 0, nil
	}
	if err != nil {
		return project.View{}, 0, err
	}
	v, err := load(tx, mission, record)
	return v, applied, err
}

// Append writes e to its mission log only if it applies to the mission as it stands, then
// syncs. The store write lock is held throughout, so two appends to one mission never
// both pass the check. An event that does not apply is refused and the log stays clean.
func (s *Store) Append(e contract.Event) (int64, error) {
	if err := eventlog.Validate(e); err != nil {
		return 0, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // no-op after commit
	if _, err := s.syncTx(tx, e.Mission); err != nil && !errors.Is(err, os.ErrNotExist) {
		return 0, err
	}
	v, _, err := current(tx, e.Mission)
	if err != nil {
		return 0, err
	}
	if err := project.Apply(&v, e, s.kinds); err != nil {
		return 0, fmt.Errorf("refused: %w", err)
	}
	end, err := eventlog.Append(s.root, e)
	if err != nil {
		return 0, err
	}
	if _, err := s.syncTx(tx, e.Mission); err != nil {
		return 0, err
	}
	return end, tx.Commit()
}

func load(tx *sql.Tx, mission, record string) (project.View, error) {
	var v project.View
	if err := json.Unmarshal([]byte(record), &v.Mission); err != nil {
		return v, err
	}
	rows, err := tx.Query(`SELECT record FROM inbox WHERE mission = ? ORDER BY created, id`, mission)
	if err != nil {
		return v, err
	}
	defer rows.Close()
	for rows.Next() {
		var rec string
		var it contract.InboxItem
		if err := rows.Scan(&rec); err != nil {
			return v, err
		}
		if err := json.Unmarshal([]byte(rec), &it); err != nil {
			return v, err
		}
		v.Inbox = append(v.Inbox, it)
	}
	return v, rows.Err()
}

func save(tx *sql.Tx, v project.View, applied int64) error {
	m := v.Mission
	rec, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO missions(id, state, kind, project, updated, applied, record)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET state=excluded.state, kind=excluded.kind, project=excluded.project,
		  updated=excluded.updated, applied=excluded.applied, record=excluded.record`,
		m.ID, m.State, m.Kind, m.Project, m.Updated, applied, string(rec)); err != nil {
		return err
	}
	for _, q := range []string{`DELETE FROM runs WHERE mission = ?`, `DELETE FROM inbox WHERE mission = ?`} {
		if _, err := tx.Exec(q, m.ID); err != nil {
			return err
		}
	}
	for _, r := range m.Runs {
		b, err := json.Marshal(r)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO runs(run, mission, started, ended, record) VALUES (?, ?, ?, ?, ?)`,
			r.Run, r.Mission, r.Started, r.Ended, string(b)); err != nil {
			return fmt.Errorf("run %s: %w", r.Run, err)
		}
	}
	for _, it := range v.Inbox {
		b, err := json.Marshal(it)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO inbox(id, mission, state, created, record) VALUES (?, ?, ?, ?, ?)`,
			it.ID, it.Mission, it.State, it.Created, string(b)); err != nil {
			return err
		}
	}
	return nil
}

// Dump is every row of the store, each list sorted by id.
type Dump struct {
	Missions []contract.Mission   `json:"missions"`
	Runs     []contract.Run       `json:"runs"`
	Inbox    []contract.InboxItem `json:"inbox"`
}

// Dump reads every row.
func (s *Store) Dump() (Dump, error) {
	d := Dump{Missions: []contract.Mission{}, Runs: []contract.Run{}, Inbox: []contract.InboxItem{}}
	if err := scan(s.db, `SELECT record FROM missions ORDER BY id`, &d.Missions); err != nil {
		return d, err
	}
	if err := scan(s.db, `SELECT record FROM runs ORDER BY run`, &d.Runs); err != nil {
		return d, err
	}
	return d, scan(s.db, `SELECT record FROM inbox ORDER BY id`, &d.Inbox)
}

func scan[T any](db *sql.DB, q string, out *[]T) error {
	rows, err := db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var rec string
		var v T
		if err := rows.Scan(&rec); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(rec), &v); err != nil {
			return err
		}
		*out = append(*out, v)
	}
	return rows.Err()
}

// Applied is the log offset the store has applied for a mission, 0 when unknown.
func (s *Store) Applied(mission string) (int64, error) {
	var n int64
	err := s.db.QueryRow(`SELECT applied FROM missions WHERE id = ?`, mission).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return n, err
}

// Version changes whenever the store applies an event: it is the mission count and the
// sum of their applied log offsets. Readers key caches and ETags on it.
func (s *Store) Version() (string, error) {
	var n, sum int64
	err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(applied), 0) FROM missions`).Scan(&n, &sum)
	return fmt.Sprintf("%d.%d", n, sum), err
}
