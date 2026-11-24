package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"crow-core/contract"
	"crow-core/internal/work"
)

// ChangesMax caps the files and commits GET /api/changes lists.
const ChangesMax = 50

// changesTimeout bounds the git calls of one GET /api/changes.
const changesTimeout = 5 * time.Second

var errNoGit = errors.New("not a git checkout with a main or master branch")

// changes serves GET /api/changes?cwd=<path> (item_273): the git state of a live lane's
// checkout. The cwd must sit in a registered project or its .wt worktrees; finished runs
// read their stored diff stat from GET /api/runs/{id} instead.
func (s *Server) changes(w http.ResponseWriter, r *http.Request) {
	cwd := r.URL.Query().Get("cwd")
	if cwd == "" || !filepath.IsAbs(cwd) || filepath.Clean(cwd) != cwd {
		fail(w, http.StatusBadRequest, "cwd: an absolute clean path is required")
		return
	}
	var reg []work.Project
	if e := s.ext["projects"]; e != nil {
		m, _ := s.extSection(e, true).Data.(map[string]any)
		reg = work.Registry(m)
	}
	if work.ProjectOfCwd(cwd, reg) == "" {
		fail(w, http.StatusNotFound, "cwd is not in a registered project")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), changesTimeout)
	defer cancel()
	c, err := readChanges(ctx, cwd)
	if errors.Is(err, errNoGit) {
		fail(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		fail(w, http.StatusInternalServerError, err.Error())
		return
	}
	reply(w, r, hash(c), c)
}

// readChanges runs git in cwd: branch, ahead/behind the main branch, numstat of
// <base>...HEAD, the uncommitted count and the commits of <base>..HEAD.
func readChanges(ctx context.Context, cwd string) (contract.Changes, error) {
	if st, err := os.Stat(cwd); err != nil || !st.IsDir() {
		return contract.Changes{}, errNoGit
	}
	git := func(args ...string) (string, error) {
		cmd := exec.CommandContext(ctx, "git", append([]string{"-C", cwd}, args...)...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		return strings.TrimRight(out.String(), "\n"), err
	}
	c := contract.Changes{Cwd: cwd, Commits: []contract.Commit{}}
	for _, b := range []string{"main", "master"} {
		if _, err := git("rev-parse", "--verify", "-q", "refs/heads/"+b); err == nil {
			c.Base = b
			break
		}
	}
	if c.Base == "" {
		return c, errNoGit
	}
	c.Branch, _ = git("symbolic-ref", "-q", "--short", "HEAD")
	lr, err := git("rev-list", "--left-right", "--count", c.Base+"...HEAD")
	if err != nil {
		return c, err
	}
	if f := strings.Fields(lr); len(f) == 2 {
		c.Behind, _ = strconv.Atoi(f[0])
		c.Ahead, _ = strconv.Atoi(f[1])
	}
	st, err := git("status", "--porcelain")
	if err != nil {
		return c, err
	}
	if st != "" {
		c.Uncommitted = strings.Count(st, "\n") + 1
	}
	if c.DiffStat.Base, err = git("merge-base", c.Base, "HEAD"); err != nil {
		return c, err
	}
	ns, err := git("diff", "--numstat", c.Base+"...HEAD")
	if err != nil {
		return c, err
	}
	c.DiffStat.PerFile = []contract.FileStat{}
	for _, l := range strings.Split(ns, "\n") {
		f := strings.SplitN(l, "\t", 3)
		if len(f) < 3 {
			continue
		}
		ins, _ := strconv.Atoi(f[0]) // "-" for binary files counts 0
		del, _ := strconv.Atoi(f[1])
		c.DiffStat.Files++
		c.DiffStat.Insertions += ins
		c.DiffStat.Deletions += del
		if len(c.DiffStat.PerFile) < ChangesMax {
			c.DiffStat.PerFile = append(c.DiffStat.PerFile, contract.FileStat{Path: f[2], Insertions: ins, Deletions: del})
		}
	}
	c.DiffStat.Truncated = c.DiffStat.Files > len(c.DiffStat.PerFile)
	lg, err := git("log", "--format=%h%x09%cI%x09%s", "-n", strconv.Itoa(ChangesMax+1), c.Base+"..HEAD")
	if err != nil {
		return c, err
	}
	for _, l := range strings.Split(lg, "\n") {
		f := strings.SplitN(l, "\t", 3)
		if len(f) < 3 {
			continue
		}
		if len(c.Commits) == ChangesMax {
			c.More = true
			break
		}
		c.Commits = append(c.Commits, contract.Commit{SHA: f[0], Time: f[1], Subject: f[2]})
	}
	return c, nil
}
