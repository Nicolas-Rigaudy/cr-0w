// Package eventlog appends to and reads missions/<id>/events.jsonl, the source of truth
// of a mission (adr_001). Lines are never rewritten. Every line is checked against the
// core contract before it is written and when it is read back.
package eventlog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"syscall"
	"time"

	"crow-core/contract"
)

// FileName is the event log of a mission dir.
const FileName = "events.jsonl"

var missionID = regexp.MustCompile(`^[0-9]{8}-[a-z-]+-[a-z0-9_!#@.-]+$`)

// ValidMission checks a mission id: <YYYYMMDD>-<kind>-<key>, safe as a dir name.
func ValidMission(id string) error {
	if !missionID.MatchString(id) || id == "." || id == ".." {
		return fmt.Errorf("bad mission id %q", id)
	}
	return nil
}

// Path is the event log of a mission under the missions root.
func Path(root, mission string) string {
	return filepath.Join(root, mission, FileName)
}

// Validate checks an event against the contract: mission id, actor, ts and the data of its type.
func Validate(e contract.Event) error {
	if err := ValidMission(e.Mission); err != nil {
		return err
	}
	switch e.Actor {
	case "owner", "core", "run":
	default:
		return fmt.Errorf("bad actor %q", e.Actor)
	}
	if _, err := time.Parse(time.RFC3339, e.TS); err != nil {
		return fmt.Errorf("bad ts %q", e.TS)
	}
	_, err := contract.DecodeEventData(e)
	return err
}

// Append validates e and appends it as one line to its mission log, under an exclusive
// file lock, then fsyncs. It returns the log offset after the line.
func Append(root string, e contract.Event) (int64, error) {
	if err := Validate(e); err != nil {
		return 0, err
	}
	line, err := json.Marshal(e)
	if err != nil {
		return 0, err
	}
	line = append(line, '\n')

	dir := filepath.Join(root, e.Mission)
	newDir := false
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		newDir = true
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(filepath.Join(dir, FileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return 0, err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN) //nolint:errcheck // close releases it too
	if _, err := f.Write(line); err != nil {
		return 0, err
	}
	if err := f.Sync(); err != nil {
		return 0, err
	}
	end, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	if newDir {
		syncDir(root)
	}
	return end, nil
}

func syncDir(p string) {
	if d, err := os.Open(p); err == nil {
		_ = d.Sync()
		d.Close()
	}
}

// Record is one event read back from a log, with the offset just after its line.
type Record struct {
	Event contract.Event
	End   int64
}

// ErrShrunk means the log is shorter than an offset already applied: it was cut or
// replaced, and only a rebuild can bring the store back in line.
var ErrShrunk = errors.New("event log is shorter than the applied offset; run crow-core rebuild")

// Read returns the complete lines of a log from offset on. A last line without its
// newline is a write in progress and is left for the next read.
func Read(path string, offset int64) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if st.Size() < offset {
		return nil, fmt.Errorf("%s: %w", path, ErrShrunk)
	}
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	var out []Record
	pos := offset
	for {
		line, err := r.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
		start := pos
		pos += int64(len(line))
		var e contract.Event
		if err := contract.Decode(bytes.TrimSuffix(line, []byte("\n")), &e); err != nil {
			return nil, fmt.Errorf("%s at %d: %w", path, start, err)
		}
		if err := Validate(e); err != nil {
			return nil, fmt.Errorf("%s at %d: %w", path, start, err)
		}
		if want := filepath.Base(filepath.Dir(path)); e.Mission != want {
			return nil, fmt.Errorf("%s at %d: event of mission %s in the log of %s", path, start, e.Mission, want)
		}
		out = append(out, Record{Event: e, End: pos})
	}
}

// Missions lists the mission ids under root that have an event log, sorted.
func Missions(root string) ([]string, error) {
	ents, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, d := range ents {
		if !d.IsDir() || ValidMission(d.Name()) != nil {
			continue
		}
		if _, err := os.Stat(Path(root, d.Name())); err == nil {
			ids = append(ids, d.Name())
		}
	}
	sort.Strings(ids)
	return ids, nil
}
