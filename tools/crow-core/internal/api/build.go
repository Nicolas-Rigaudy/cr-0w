package api

// Nest rebuild after a settings save (item_283). A hero or pages change is compiled into
// Nest, so a save that needs a build starts `crow nest build` here, and so does a module
// switch that turns a brick with Nest code on or off (buildOnSwitch). The build writes a
// staging dir next to the served dist; only a good build replaces it, so a failed one
// leaves the last good Nest served. One build runs at a time: a save during a build
// queues one more. GET /api/settings/nest/build answers the state; /api/state carries
// the served build id, so its version moves when a new Nest is in place.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// BuildTimeout bounds one Nest build.
const BuildTimeout = 10 * time.Minute

// buildLogLines is how many log lines the build state keeps.
const buildLogLines = 20

// Build states.
const (
	BuildIdle     = "idle"
	BuildBuilding = "building"
	BuildOK       = "ok"
	BuildFailed   = "failed"
)

// NestBuild is GET /api/settings/nest/build.
type NestBuild struct {
	State     string   `json:"state"`                // idle, building, ok or failed
	StartedAt string   `json:"started_at,omitempty"` // RFC 3339, the build running or the last one
	Elapsed   float64  `json:"elapsed_s"`            // seconds, so far or in total
	Queued    bool     `json:"queued"`               // a save came in during this build
	Brick     string   `json:"brick,omitempty"`      // the brick a failed build names
	Error     string   `json:"error,omitempty"`      // why the last build failed
	Log       []string `json:"log"`                  // last lines of the build output
	Served    string   `json:"served"`               // id of the Nest build being served
}

// builder runs the Nest builds of one server.
type builder struct {
	mu      sync.Mutex
	st      NestBuild
	started time.Time
	ended   time.Time
	running bool
	mods    *Modules // brick state last seen by buildOnSwitch
	nestKey string   // its on bricks that ship Nest code
}

var (
	brickErr = regexp.MustCompile(`brick ([a-z][a-z0-9-]*):`)
	ansi     = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)
)

// buildOn says whether this server can rebuild Nest: it serves a dist and has a command.
func (s *Server) buildOn() bool { return s.o.Nest != "" && len(s.o.NestBuild) > 0 }

// nestBuild answers the build state.
func (s *Server) nestBuild(w http.ResponseWriter, r *http.Request) {
	st := s.buildState()
	w.Header().Set("Cache-Control", "no-store")
	reply(w, r, hash(st), st)
}

// buildState is a copy of the state with the elapsed time and served id filled in.
func (s *Server) buildState() NestBuild {
	b := &s.build
	b.mu.Lock()
	defer b.mu.Unlock()
	st := b.st
	st.Log = append([]string{}, b.st.Log...)
	if st.State == "" {
		st.State = BuildIdle
	}
	switch {
	case b.running:
		st.Elapsed = s.o.Now().Sub(b.started).Seconds()
	case !b.ended.IsZero():
		st.Elapsed = b.ended.Sub(b.started).Seconds()
	}
	st.Elapsed = float64(int(st.Elapsed*10)) / 10
	st.Served = distID(s.o.Nest)
	return st
}

// servedBuild is the id of the Nest build being served, "" without one.
func (s *Server) servedBuild() string {
	if s.o.Nest == "" {
		return ""
	}
	return distID(s.o.Nest)
}

// buildOnSwitch starts a build when a module switch changed the on bricks that ship Nest
// code: switches go through the same rebuild and reload path as hero and pages. The
// first brick state core reads is the one its dist was built from.
func (s *Server) buildOnSwitch(m *Modules) {
	if m == nil {
		return
	}
	b := &s.build
	b.mu.Lock()
	if m == b.mods {
		b.mu.Unlock()
		return
	}
	key := m.nestBricks()
	moved := b.mods != nil && key != b.nestKey
	b.mods, b.nestKey = m, key
	b.mu.Unlock()
	if moved && s.buildOn() {
		s.startBuild()
	}
}

// startBuild starts a build, or queues one more when one runs.
func (s *Server) startBuild() {
	b := &s.build
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.running {
		b.st.Queued = true
		return
	}
	b.running = true
	b.begin(s.o.Now())
	go s.runBuilds()
}

// begin resets the state for a new build. The caller holds mu.
func (b *builder) begin(now time.Time) {
	b.started, b.ended = now, time.Time{}
	b.st = NestBuild{State: BuildBuilding, StartedAt: now.UTC().Format(time.RFC3339), Log: []string{}}
}

// runBuilds runs builds until none is queued.
func (s *Server) runBuilds() {
	b := &s.build
	for {
		brick, err := s.buildOnce()
		b.mu.Lock()
		b.ended = s.o.Now()
		if err != nil {
			b.st.State, b.st.Brick, b.st.Error = BuildFailed, brick, err.Error()
		} else {
			b.st.State = BuildOK
		}
		if b.st.Queued {
			b.begin(s.o.Now())
			b.mu.Unlock()
			continue
		}
		b.running = false
		b.mu.Unlock()
		return
	}
}

// buildOnce builds into <dist>.next and swaps it in. On failure it answers the brick the
// error names, if any, and the error.
func (s *Server) buildOnce() (string, error) {
	dist := filepath.Clean(s.o.Nest)
	next, prev := dist+".next", dist+".prev"
	if err := os.RemoveAll(next); err != nil {
		return "", err
	}
	defer os.RemoveAll(next)
	ctx, cancel := context.WithTimeout(context.Background(), BuildTimeout)
	defer cancel()
	argv := append(append([]string{}, s.o.NestBuild...), "--outDir", next, "--emptyOutDir")
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Env = append(os.Environ(), "CROW_HOME="+s.o.Home, "CROW_ENGINE="+s.o.Engine, "NO_COLOR=1", "FORCE_COLOR=0")
	out := &buildLog{b: &s.build}
	cmd.Stdout, cmd.Stderr = out, out
	err := cmd.Run()
	out.flush()
	switch {
	case err == nil:
	case ctx.Err() != nil:
		return out.brick, errors.New("nest build timed out")
	case out.brickLine != "":
		return out.brick, errors.New(out.brickLine)
	case out.last != "":
		return "", fmt.Errorf("nest build: %s", out.last)
	default:
		return "", fmt.Errorf("nest build: %w", err)
	}
	if _, err := os.Stat(filepath.Join(next, "index.html")); err != nil {
		return "", errors.New("nest build wrote no index.html")
	}
	if err := os.RemoveAll(prev); err != nil {
		return "", err
	}
	if err := os.Rename(dist, prev); err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.Rename(next, dist); err != nil {
		_ = os.Rename(prev, dist)
		return "", err
	}
	return "", os.RemoveAll(prev)
}

// distID is a short hash of a dist's index.html, which names the hashed assets of its
// build; "" when there is none.
func distID(dist string) string {
	b, err := os.ReadFile(filepath.Join(dist, "index.html"))
	if err != nil {
		return ""
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])[:12]
}

// buildLog keeps the last lines of the build output in the state and the first line
// that names a brick.
type buildLog struct {
	b         *builder
	part      []byte
	last      string
	brick     string
	brickLine string
}

func (l *buildLog) Write(p []byte) (int, error) {
	l.part = append(l.part, p...)
	for {
		i := bytes.IndexByte(l.part, '\n')
		if i < 0 {
			break
		}
		l.line(string(l.part[:i]))
		l.part = l.part[i+1:]
	}
	return len(p), nil
}

func (l *buildLog) flush() {
	if len(l.part) > 0 {
		l.line(string(l.part))
		l.part = nil
	}
}

func (l *buildLog) line(s string) {
	s = strings.TrimRight(ansi.ReplaceAllString(s, ""), " \r\t")
	if strings.TrimSpace(s) == "" {
		return
	}
	l.last = strings.TrimSpace(s)
	if m := brickErr.FindStringSubmatchIndex(s); m != nil && l.brick == "" {
		l.brick, l.brickLine = s[m[2]:m[3]], strings.TrimSpace(s[m[0]:])
	}
	l.b.mu.Lock()
	l.b.st.Log = append(l.b.st.Log, s)
	if n := len(l.b.st.Log); n > buildLogLines {
		l.b.st.Log = append([]string{}, l.b.st.Log[n-buildLogLines:]...)
	}
	l.b.mu.Unlock()
}
