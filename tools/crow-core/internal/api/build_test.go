package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// fakeBuild is a stand-in for `crow nest build`. It reads its mode from <dir>/mode:
// ok writes <outDir>/index.html, fail prints a brick error and exits 1, wait blocks
// until <dir>/go exists and then builds. Each run adds a line to <dir>/runs.
const fakeBuild = `#!/bin/sh
dir=%DIR%
out=""
while [ $# -gt 0 ]; do [ "$1" = --outDir ] && out=$2; shift; done
echo run >> "$dir/runs"
echo "vite building for production..."
case "$(cat "$dir/mode")" in
fail) printf 'error during build:\n\033[31m[nest-modules] brick garden: page garden.tsx has no default export\033[0m\n' >&2; exit 1 ;;
wait) while [ ! -f "$dir/go" ]; do sleep 0.02; done ;;
esac
mkdir -p "$out"
echo "<html>build $(wc -l < "$dir/runs")</html>" > "$out/index.html"
`

// buildServer is a settings server that serves a temp dist (an "old" build) and
// rebuilds it with fakeBuild in the given mode.
func buildServer(t *testing.T, mode string) (s *Server, home, fake, dist string) {
	t.Helper()
	s, home = settingsServer(t)
	fake = t.TempDir()
	dist = filepath.Join(t.TempDir(), "dist")
	writeFile(t, filepath.Join(dist, "index.html"), []byte("<html>old</html>"))
	script := filepath.Join(fake, "build")
	writeFile(t, script, []byte(strings.ReplaceAll(fakeBuild, "%DIR%", fake)))
	if err := os.Chmod(script, 0o755); err != nil {
		t.Fatal(err)
	}
	setMode(t, fake, mode)
	s.o.Nest, s.o.NestBuild = dist, []string{script, "nest", "build"}
	return s, home, fake, dist
}

func setMode(t *testing.T, fake, mode string) {
	t.Helper()
	writeFile(t, filepath.Join(fake, "mode"), []byte(mode))
}

func save(t *testing.T, s *Server, body string) NestSaved {
	t.Helper()
	w := do(t, s, call{method: "POST", path: "/api/settings/nest", header: writeHeader, body: body})
	var v NestSaved
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil || w.Code != 200 {
		t.Fatalf("POST %s: %d %s", body, w.Code, w.Body)
	}
	return v
}

// waitBuild polls the build route until no build runs.
func waitBuild(t *testing.T, s *Server) NestBuild {
	t.Helper()
	for end := time.Now().Add(10 * time.Second); time.Now().Before(end); time.Sleep(20 * time.Millisecond) {
		w := do(t, s, call{path: "/api/settings/nest/build"})
		var st NestBuild
		if err := json.Unmarshal(w.Body.Bytes(), &st); err != nil || w.Code != 200 {
			t.Fatalf("GET build: %d %s", w.Code, w.Body)
		}
		if st.State != BuildBuilding {
			return st
		}
	}
	t.Fatal("build still running after 10s")
	return NestBuild{}
}

func runs(t *testing.T, fake string) int {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(fake, "runs"))
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(b), "\n")
}

// liveModules is the /api/modules answer Nest polls.
func liveModules(t *testing.T, s *Server) Modules {
	t.Helper()
	w := do(t, s, call{path: "/api/modules"})
	var m Modules
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil || w.Code != 200 || w.Header().Get("ETag") != `"`+m.Version+`"` {
		t.Fatalf("GET modules: %d %s %s", w.Code, w.Header().Get("ETag"), w.Body)
	}
	return m
}

func TestNestBuildOK(t *testing.T) {
	s, _, fake, dist := buildServer(t, "ok")
	before := liveModules(t, s)
	if before.NestBuild == "" || before.NestBuild != distID(dist) {
		t.Fatalf("served build: %q", before.NestBuild)
	}

	// a theme-only save starts no build
	if v := save(t, s, `{"theme":"cyber"}`); v.NeedsBuild || v.Build.State != BuildIdle {
		t.Errorf("theme only: %+v", v.Build)
	}
	if n := runs(t, fake); n != 0 {
		t.Fatalf("theme only ran %d builds", n)
	}

	v := save(t, s, `{"hero":"garden"}`)
	if !v.NeedsBuild || (v.Build.State != BuildBuilding && v.Build.State != BuildOK) {
		t.Errorf("save answer: needs_build %v, build %+v", v.NeedsBuild, v.Build)
	}
	st := waitBuild(t, s)
	if st.State != BuildOK || st.Error != "" || st.StartedAt == "" || !slices.Contains(st.Log, "vite building for production...") {
		t.Errorf("state: %+v", st)
	}
	if b, _ := os.ReadFile(filepath.Join(dist, "index.html")); string(b) != "<html>build 1</html>\n" {
		t.Errorf("dist not replaced: %q", b)
	}
	for _, left := range []string{dist + ".next", dist + ".prev"} {
		if _, err := os.Stat(left); !os.IsNotExist(err) {
			t.Errorf("%s left behind: %v", left, err)
		}
	}
	after := liveModules(t, s)
	if after.NestBuild == before.NestBuild || after.NestBuild != st.Served || after.Version == before.Version {
		t.Errorf("live feed did not move: before %s/%s, after %s/%s", before.NestBuild, before.Version, after.NestBuild, after.Version)
	}
}

func TestNestBuildFailedKeepsDist(t *testing.T) {
	s, home, fake, dist := buildServer(t, "fail")
	before := distID(dist)
	save(t, s, `{"hero":"garden","pages":["home","garden"]}`)
	st := waitBuild(t, s)
	if st.State != BuildFailed || st.Brick != "garden" ||
		st.Error != "brick garden: page garden.tsx has no default export" {
		t.Errorf("state: %+v", st)
	}
	if !slices.Contains(st.Log, "[nest-modules] brick garden: page garden.tsx has no default export") {
		t.Errorf("log keeps no plain error line: %q", st.Log)
	}
	if distID(dist) != before || st.Served != before {
		t.Errorf("last good dist not kept: %s, want %s", distID(dist), before)
	}
	if _, err := os.Stat(dist + ".next"); !os.IsNotExist(err) {
		t.Errorf("staging left behind: %v", err)
	}
	if y := readYAML(t, home); !strings.Contains(y, "hero: garden") || !strings.Contains(y, "pages: [home, garden]") {
		t.Errorf("crow.yaml lost the new values:\n%s", y)
	}

	// the owner fixes the brick: the next save builds again
	setMode(t, fake, "ok")
	save(t, s, `{"pages":["home"]}`)
	if st := waitBuild(t, s); st.State != BuildOK || st.Brick != "" || st.Error != "" {
		t.Errorf("after fix: %+v", st)
	}
}

func TestNestBuildQueuesOneMore(t *testing.T) {
	s, _, fake, dist := buildServer(t, "wait")
	if v := save(t, s, `{"hero":"garden"}`); v.Build.State != BuildBuilding || v.Build.Queued {
		t.Fatalf("first save: %+v", v.Build)
	}
	for _, body := range []string{`{"pages":["home","garden"]}`, `{"pages":["garden"]}`} {
		if v := save(t, s, body); v.Build.State != BuildBuilding || !v.Build.Queued {
			t.Fatalf("save during a build: %+v", v.Build)
		}
	}
	writeFile(t, filepath.Join(fake, "go"), nil)
	st := waitBuild(t, s)
	if st.State != BuildOK || st.Queued {
		t.Errorf("state: %+v", st)
	}
	if n := runs(t, fake); n != 2 {
		t.Errorf("%d builds ran, want 2 (one queued for both saves)", n)
	}
	if b, _ := os.ReadFile(filepath.Join(dist, "index.html")); string(b) != "<html>build 2</html>\n" {
		t.Errorf("dist is not the last build: %q", b)
	}
}

func TestNestBuildOff(t *testing.T) {
	s, _ := settingsServer(t) // no build command: the save still answers needs_build
	if v := save(t, s, `{"hero":"garden"}`); !v.NeedsBuild || v.Build.State != BuildIdle {
		t.Errorf("answer: %+v", v)
	}
}
