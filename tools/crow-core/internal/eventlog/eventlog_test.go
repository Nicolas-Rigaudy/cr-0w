package eventlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"

	"crow-core/contract"
	"crow-core/internal/testfix"
)

const mission = "20261109-research-lock-test"

func incident(i int, pad int) contract.Event {
	detail := fmt.Sprintf("writer %d %s", i, strings.Repeat("x", pad))
	data, _ := json.Marshal(contract.IncidentData{Cause: "test", Detail: &detail})
	return contract.Event{TS: "2026-11-09T14:00:00+01:00", Mission: mission, Type: "incident", Actor: "core", Data: data}
}

func TestAppendThenRead(t *testing.T) {
	root := t.TempDir()
	var ends []int64
	src, err := Read(Path(testfix.MissionsRoot(t), "20261109-build-item-095"), 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range src {
		e := r.Event
		e.Mission = mission
		end, err := Append(root, e)
		if err != nil {
			t.Fatal(err)
		}
		ends = append(ends, end)
	}
	recs, err := Read(Path(root, mission), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != len(ends) {
		t.Fatalf("read %d events, appended %d", len(recs), len(ends))
	}
	for i, r := range recs {
		if r.End != ends[i] {
			t.Errorf("event %d ends at %d, Append said %d", i, r.End, ends[i])
		}
	}
	tail, err := Read(Path(root, mission), recs[2].End)
	if err != nil || len(tail) != len(recs)-3 {
		t.Fatalf("tail from event 3: %d events, err %v", len(tail), err)
	}
}

func TestAppendRefusesBadEvents(t *testing.T) {
	root := t.TempDir()
	good := incident(0, 0)
	for name, mut := range map[string]func(*contract.Event){
		"mission": func(e *contract.Event) { e.Mission = "../escape" },
		"actor":   func(e *contract.Event) { e.Actor = "someone" },
		"ts":      func(e *contract.Event) { e.TS = "yesterday" },
		"type":    func(e *contract.Event) { e.Type = "nope" },
		"data":    func(e *contract.Event) { e.Data = json.RawMessage(`{"cause":"x","extra":1}`) },
	} {
		e := good
		mut(&e)
		if _, err := Append(root, e); err == nil {
			t.Errorf("%s: bad event accepted", name)
		}
	}
	if _, err := os.Stat(Path(root, mission)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a refused event created the log")
	}
}

func TestReadLeavesAPartialLastLine(t *testing.T) {
	root := t.TempDir()
	if _, err := Append(root, incident(1, 0)); err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(Path(root, mission), os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`{"ts":"2026-11-09T14:00:00+01:00","mission":`)
	f.Close()
	recs, err := Read(Path(root, mission), 0)
	if err != nil || len(recs) != 1 {
		t.Fatalf("got %d events, err %v; want the one complete line", len(recs), err)
	}
}

func TestReadRefusesAShrunkLog(t *testing.T) {
	root := t.TempDir()
	end, err := Append(root, incident(1, 0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Read(Path(root, mission), end+10); !errors.Is(err, ErrShrunk) {
		t.Fatalf("want ErrShrunk, got %v", err)
	}
}

// checkLines fails unless the log holds exactly want lines, each one whole event.
func checkLines(t *testing.T, root string, want int) {
	t.Helper()
	recs, err := Read(Path(root, mission), 0)
	if err != nil {
		t.Fatalf("log has a broken line: %v", err)
	}
	if len(recs) != want {
		t.Fatalf("log has %d events, want %d", len(recs), want)
	}
}

// AC2: lines bigger than a pipe buffer, written from many goroutines, never interleave.
func TestParallelAppendsInOneProcess(t *testing.T) {
	root := t.TempDir()
	const writers, each = 8, 25
	var wg sync.WaitGroup
	for w := range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range each {
				if _, err := Append(root, incident(w*each+i, 70_000)); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()
	checkLines(t, root, writers*each)
}

// AC2 across processes: the test binary re-runs itself as writer processes.
func TestParallelAppendsAcrossProcesses(t *testing.T) {
	if os.Getenv("CROW_CORE_WRITER") != "" {
		return
	}
	root := t.TempDir()
	const writers, each = 4, 20
	var cmds []*exec.Cmd
	for w := range writers {
		c := exec.Command(os.Args[0], "-test.run=^TestWriterProcess$")
		c.Env = append(os.Environ(), "CROW_CORE_WRITER="+strconv.Itoa(w), "CROW_CORE_ROOT="+root, "CROW_CORE_EACH="+strconv.Itoa(each))
		if err := c.Start(); err != nil {
			t.Fatal(err)
		}
		cmds = append(cmds, c)
	}
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			t.Fatalf("writer: %v", err)
		}
	}
	checkLines(t, root, writers*each)
}

func TestWriterProcess(t *testing.T) {
	w, err := strconv.Atoi(os.Getenv("CROW_CORE_WRITER"))
	if err != nil {
		t.Skip("run by TestParallelAppendsAcrossProcesses")
	}
	each, _ := strconv.Atoi(os.Getenv("CROW_CORE_EACH"))
	for i := range each {
		if _, err := Append(os.Getenv("CROW_CORE_ROOT"), incident(w*each+i, 70_000)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestMissionsListsLogsOnly(t *testing.T) {
	root := testfix.MissionsRoot(t)
	os.MkdirAll(root+"/20261109-build-no-log", 0o700)
	os.WriteFile(root+"/core.db", nil, 0o600)
	ids, err := Missions(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != "20261109-build-item-095,20261109-research-android-push-options" {
		t.Errorf("missions: %v", ids)
	}
}
