package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// settingsYAML has comments and other bricks around modules.nest: a save must leave them.
const settingsYAML = `# crow.yaml of a test home
version: 1
modules:
  mail: false  # no mail
  forge: false
  nest:
    theme: modern   # first paint
    pages:
      - home
      - desk
  notes: true

# trailing note
`

// settingsServer is a module server whose home also holds the garden brick: a theme
// (moss), a hero (garden) and a page (garden).
func settingsServer(t *testing.T) (*Server, string) {
	t.Helper()
	s, home := moduleServer(t, settingsYAML)
	garden := filepath.Join(home, "modules", "garden")
	copyTree(t, filepath.Join(s.o.Engine, "tools", "nest", "fixtures", "bricks", "garden"), garden)
	writeFile(t, filepath.Join(garden, "module.yaml"), []byte("summary: Garden pieces for Nest.\n"))
	return s, home
}

var writeHeader = map[string]string{CSRFHeader: "1", "Origin": "https://core.example"}

func readYAML(t *testing.T, home string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(home, "crow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestNestSettingsGet(t *testing.T) {
	s, home := settingsServer(t)
	w := do(t, s, call{path: "/api/settings/nest"})
	if w.Code != 200 || w.Header().Get("ETag") == "" {
		t.Fatalf("GET: %d %s", w.Code, w.Body)
	}
	var v NestSettings
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if v.Theme != "modern" || v.Hero != "holotable" || v.File != filepath.Join(home, "crow.yaml") {
		t.Errorf("values: %+v", v)
	}
	if !slices.Contains(v.Known.Themes, "moss") || !slices.Contains(v.Known.Heroes, "garden") ||
		!slices.Contains(v.Known.Pages, "garden") || !slices.Contains(v.Known.Pages, "home") {
		t.Errorf("known: %+v", v.Known)
	}
}

func TestNestSettingsSave(t *testing.T) {
	s, home := settingsServer(t)
	w := do(t, s, call{method: "POST", path: "/api/settings/nest", header: writeHeader,
		body: `{"theme":"moss","hero":"garden","pages":["home","garden"]}`})
	if w.Code != 200 {
		t.Fatalf("POST: %d %s", w.Code, w.Body)
	}
	var v NestSaved
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if v.Theme != "moss" || v.Hero != "garden" || !v.NeedsBuild || len(v.Changed) != 3 {
		t.Errorf("answer: %+v", v)
	}
	want := strings.Replace(settingsYAML, "    theme: modern   # first paint\n    pages:\n      - home\n      - desk\n",
		"    theme: moss  # first paint\n    pages: [home, garden]\n    hero: garden\n", 1)
	if got := readYAML(t, home); got != want {
		t.Errorf("crow.yaml:\n%s\nwant:\n%s", got, want)
	}
	log, err := os.ReadFile(s.o.ActionLog)
	if err != nil || !strings.Contains(string(log), `"action":"settings"`) {
		t.Errorf("action log: %v %s", err, log)
	}

	// a theme-only save needs no build
	w = do(t, s, call{method: "POST", path: "/api/settings/nest", header: writeHeader, body: `{"theme":"cyber"}`})
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil || w.Code != 200 || v.NeedsBuild || !slices.Equal(v.Changed, []string{"theme"}) {
		t.Errorf("theme only: %d %s", w.Code, w.Body)
	}
}

func TestNestSettingsRefuses(t *testing.T) {
	s, home := settingsServer(t)
	for _, c := range []struct {
		body   string
		header map[string]string
		code   int
		want   string
	}{
		{`{"theme":"neon"}`, writeHeader, 400, "neon"},
		{`{"hero":"ghost"}`, writeHeader, 400, "ghost"},
		{`{"pages":["home","attic"]}`, writeHeader, 400, "attic"},
		{`{"pages":["home,desk"]}`, writeHeader, 400, "home,desk"},
		{`{"theme":"cyber","hero":"nope"}`, writeHeader, 400, "nope"},
		{`{}`, writeHeader, 400, "want theme"},
		{`{"colour":"red"}`, writeHeader, 400, "colour"},
		{`{"theme":"cyber"}`, map[string]string{"Origin": "https://core.example"}, 403, "header"},
		{`{"theme":"cyber"}`, map[string]string{CSRFHeader: "1", "Origin": "https://evil.example"}, 403, "origin"},
	} {
		w := do(t, s, call{method: "POST", path: "/api/settings/nest", header: c.header, body: c.body})
		if w.Code != c.code || !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("%s: %d %s, want %d %q", c.body, w.Code, w.Body, c.code, c.want)
		}
		if got := readYAML(t, home); got != settingsYAML {
			t.Fatalf("%s changed crow.yaml:\n%s", c.body, got)
		}
	}
}
