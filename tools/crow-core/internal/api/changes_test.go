package api

import (
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"crow-core/contract"
)

func gitT(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimRight(string(out), "\n")
}

// changesServer serves a registry holding one project at repo.
func changesServer(t *testing.T, repo string) *Server {
	s := New(Options{
		Owner: func() string { return owner },
		Sources: []Source{{Name: "projects", Fetch: func() (any, error) {
			return map[string]any{"projects": []any{map[string]any{"name": "demo", "short": "demo", "path": repo}}}, nil
		}}},
	})
	s.Refresh("projects")
	return s
}

// TestChanges matches git diff --numstat main... and git log main.. in a worktree (AC1).
func TestChanges(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("needs git")
	}
	root := t.TempDir()
	repo := filepath.Join(root, "demo")
	writeFile(t, filepath.Join(repo, "a.txt"), []byte("one\n"))
	gitT(t, repo, "init", "-q", "-b", "main")
	gitT(t, repo, "add", ".")
	gitT(t, repo, "commit", "-qm", "base")
	wt := filepath.Join(root, "demo.wt", "item-1")
	gitT(t, repo, "worktree", "add", "-q", "-b", "item-1", wt)
	writeFile(t, filepath.Join(wt, "a.txt"), []byte("one\ntwo\n"))
	writeFile(t, filepath.Join(wt, "dir", "b.txt"), []byte("x\ny\n"))
	gitT(t, wt, "add", ".")
	gitT(t, wt, "commit", "-qm", "first")
	writeFile(t, filepath.Join(wt, "a.txt"), []byte("two\n"))
	gitT(t, wt, "commit", "-qam", "second")
	writeFile(t, filepath.Join(repo, "c.txt"), []byte("main moved\n"))
	gitT(t, repo, "add", ".")
	gitT(t, repo, "commit", "-qm", "on main")
	writeFile(t, filepath.Join(wt, "dirty.txt"), []byte("?\n"))

	s := changesServer(t, repo)
	w := do(t, s, call{path: "/api/changes?cwd=" + url.QueryEscape(wt)})
	if w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	var c contract.Changes
	if err := json.Unmarshal(w.Body.Bytes(), &c); err != nil {
		t.Fatal(err)
	}
	if c.Branch != "item-1" || c.Base != "main" || c.Ahead != 2 || c.Behind != 1 || c.Uncommitted != 1 {
		t.Errorf("branch card: %+v", c)
	}
	var files []string
	ins, del := 0, 0
	for _, l := range strings.Split(gitT(t, wt, "diff", "--numstat", "main..."), "\n") {
		f := strings.SplitN(l, "\t", 3)
		i, _ := strconv.Atoi(f[0])
		d, _ := strconv.Atoi(f[1])
		ins, del = ins+i, del+d
		files = append(files, f[2])
	}
	if c.DiffStat.Files != len(files) || c.DiffStat.Insertions != ins || c.DiffStat.Deletions != del || len(c.DiffStat.PerFile) != len(files) {
		t.Errorf("diffstat %+v, git numstat %v +%d -%d", c.DiffStat, files, ins, del)
	}
	for i, f := range files {
		if c.DiffStat.PerFile[i].Path != f {
			t.Errorf("file %d: %s, want %s", i, c.DiffStat.PerFile[i].Path, f)
		}
	}
	log := strings.Split(gitT(t, wt, "log", "--format=%h %s", "main.."), "\n")
	if len(c.Commits) != len(log) || c.More {
		t.Fatalf("commits %+v, git log %v", c.Commits, log)
	}
	for i, l := range log {
		if got := c.Commits[i].SHA + " " + c.Commits[i].Subject; got != l {
			t.Errorf("commit %d: %q, want %q", i, got, l)
		}
	}
}

// TestChangesRefuses keeps git to registered checkouts.
func TestChangesRefuses(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "demo")
	writeFile(t, filepath.Join(repo, "a.txt"), []byte("x"))
	s := changesServer(t, repo)
	for _, c := range []struct {
		cwd  string
		code int
	}{
		{"", 400},
		{"relative/path", 400},
		{repo + "/../demo", 400},
		{root, 404},                              // not in the registry
		{filepath.Join(root, "demo-other"), 404}, // a sibling with the same prefix
		{repo, 404},                              // registered but not a git checkout
		{filepath.Join(root, "demo.wt", "gone"), 404}, // worktree removed
	} {
		if w := do(t, s, call{path: "/api/changes?cwd=" + url.QueryEscape(c.cwd)}); w.Code != c.code {
			t.Errorf("%q: %d, want %d (%s)", c.cwd, w.Code, c.code, w.Body)
		}
	}
}
