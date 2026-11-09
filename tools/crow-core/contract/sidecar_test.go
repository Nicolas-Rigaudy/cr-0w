package contract

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestPythonSidecar runs the Python base sidecar and checks it answers the lifecycle
// fixtures exactly, so both sides of the contract read the same files.
func TestPythonSidecar(t *testing.T) {
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not found")
	}
	types := []string{"hello", "health", "shutdown"}
	var in bytes.Buffer
	for _, typ := range types {
		var req Request
		if err := Decode(read(t, filepath.Join(fix, "messages", typ+".request.json")), &req); err != nil {
			t.Fatal(err)
		}
		line, _ := json.Marshal(req)
		in.Write(append(line, '\n'))
	}
	cmd := exec.Command(py, filepath.Join(repo, "bin", "crow_sidecar.py"))
	cmd.Stdin = &in
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	done := make(chan struct{})
	var out []byte
	go func() { out, err = cmd.Output(); close(done) }()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("sidecar did not exit after shutdown")
	}
	if err != nil {
		t.Fatalf("sidecar: %v\n%s", err, stderr.String())
	}
	sc := bufio.NewScanner(bytes.NewReader(out))
	i := 0
	for sc.Scan() {
		if i >= len(types) {
			t.Fatalf("extra reply after shutdown: %s", sc.Text())
		}
		var rep Reply
		if err := Decode(sc.Bytes(), &rep); err != nil {
			t.Fatalf("reply %d: %v", i, err)
		}
		sameJSON(t, types[i]+" reply", read(t, filepath.Join(fix, "messages", types[i]+".reply.json")), sc.Bytes())
		i++
	}
	if i != len(types) {
		t.Fatalf("got %d replies, want %d", i, len(types))
	}
}
