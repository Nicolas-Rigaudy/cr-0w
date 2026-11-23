package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"crow-core/contract"
	"crow-core/internal/config"
	"crow-core/internal/testfix"
)

// engineRoot is the checkout: its modules/ holds the engine manifests.
var engineRoot = filepath.Join("..", "..", "..", "..")

const notesBrick = `summary: Notes kept by a script.
scripts: [crow-notes]
nest:
  sections:
    - id: notes
      page: desk
      title: Notes
      view: list
      run: [crow-notes, json]
      ttl: 1m
      actions:
        - {id: done, label: Done, script: crow-notes, args: [done, "{item}"], confirm: true}
        - {id: broken, label: Break, script: crow-notes, args: [refuse]}
`

const notesScript = `#!/bin/sh
case "$1" in
  json) echo '{"items":[]}' ;;
  done) echo "closed $2" ;;
  refuse) echo "nothing to break" >&2; exit 2 ;;
esac
`

// brickHome is a CROW_HOME with mail and forge off and one user brick, notes.
func brickHome(t *testing.T, crowYAML string) string {
	t.Helper()
	home := t.TempDir()
	writeFile(t, filepath.Join(home, "crow.yaml"), []byte(crowYAML))
	writeFile(t, filepath.Join(home, "modules", "notes", "module.yaml"), []byte(notesBrick))
	bin := filepath.Join(home, "modules", "notes", "bin", "crow-notes")
	writeFile(t, bin, []byte(notesScript))
	if err := os.Chmod(bin, 0o700); err != nil {
		t.Fatal(err)
	}
	return home
}

func moduleServer(t *testing.T, crowYAML string) (*Server, string) {
	t.Helper()
	s, _ := newServer(t)
	home := brickHome(t, crowYAML)
	engine, err := filepath.Abs(engineRoot)
	if err != nil {
		t.Fatal(err)
	}
	s.o.Home, s.o.Engine = home, engine
	s.o.ActionLog = filepath.Join(t.TempDir(), ActionLogName)
	s.o.Modules = LoadModules(home, engine)
	return s, home
}

const mailForgeOff = "version: 1\nmodules:\n  mail: false\n  forge: false\n"

func TestModulesList(t *testing.T) {
	s, _ := moduleServer(t, mailForgeOff)
	w := do(t, s, call{path: "/api/modules"})
	if w.Code != 200 || w.Header().Get("ETag") == "" {
		t.Fatalf("GET /api/modules: %d %v", w.Code, w.Header())
	}
	var m Modules
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	if m.ContractVersion != contract.Version || m.Version == "" {
		t.Errorf("contract_version %d version %q", m.ContractVersion, m.Version)
	}
	got := map[string]ModuleBrick{}
	for _, b := range m.Modules {
		got[b.Name] = b
	}
	if got["mail"].On || got["forge"].On || !got["team"].On {
		t.Errorf("switches: mail %v forge %v team %v", got["mail"].On, got["forge"].On, got["team"].On)
	}
	n := got["notes"]
	if !n.On || !n.User || len(n.Nest.Sections) != 1 || strings.Join(n.Sections, ",") != "notes" {
		t.Fatalf("notes brick: %+v", n)
	}
	sec := n.Nest.Sections[0]
	if sec.Page != "desk" || sec.View != "list" || len(sec.Actions) != 2 || !sec.Actions[0].Row || sec.Actions[1].Row {
		t.Errorf("notes section: %+v", sec)
	}
	if strings.Contains(w.Body.String(), "crow-notes") {
		t.Error("/api/modules leaks script names")
	}
	if strings.Join(got["forge"].Deck.Pages, ",") != "MRs" {
		t.Errorf("forge deck pages: %v", got["forge"].Deck.Pages)
	}
}

func TestModuleOffGates(t *testing.T) {
	s, _ := moduleServer(t, mailForgeOff)
	w := do(t, s, call{path: "/api/state"})
	var body struct {
		Sections map[string]json.RawMessage `json:"sections"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	st := body.Sections
	for _, off := range []string{"today", "mrs"} {
		if _, ok := st[off]; ok {
			t.Errorf("/api/state still serves %s", off)
		}
	}
	if _, ok := st["checks"]; !ok {
		t.Error("/api/state lost checks (team is on)")
	}
	post := map[string]string{CSRFHeader: "1", "Origin": "https://core.example"}
	for _, c := range []call{
		{path: "/api/state/today"},
		{path: "/api/state/mrs"},
		{method: "POST", path: "/api/today/abc/done", header: post},
		{method: "POST", path: "/api/modules/mail/actions/x", header: post},
	} {
		w := do(t, s, c)
		if w.Code != 404 || !strings.Contains(w.Body.String(), "module off") {
			t.Errorf("%s %s: %d %s, want 404 module off", c.method, c.path, w.Code, w.Body.String())
		}
	}
	if w := do(t, s, call{path: "/api/state/sessions"}); w.Code != 200 {
		t.Errorf("core-api section: %d", w.Code)
	}
}

// TestModulesReload sees a crow.yaml change without a restart.
func TestModulesReload(t *testing.T) {
	s, home := moduleServer(t, mailForgeOff)
	before := do(t, s, call{path: "/api/modules"}).Header().Get("ETag")
	writeFile(t, filepath.Join(home, "crow.yaml"), []byte("version: 1\nmodules:\n  forge: false\n"))
	later := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(filepath.Join(home, "crow.yaml"), later, later); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond)
	w := do(t, s, call{path: "/api/modules"})
	if w.Header().Get("ETag") == before {
		t.Fatal("modules version did not move after crow.yaml changed")
	}
	if w := do(t, s, call{path: "/api/state/today"}); w.Code != 200 {
		t.Errorf("mail back on: /api/state/today %d", w.Code)
	}
}

func TestModuleActions(t *testing.T) {
	s, _ := moduleServer(t, mailForgeOff)
	post := map[string]string{CSRFHeader: "1", "Origin": "https://core.example"}
	for _, c := range []struct {
		path, body string
		code       int
		want       string
	}{
		{"/api/modules/notes/actions/done", `{"item":"n-42"}`, 200, "closed n-42"},
		{"/api/modules/notes/actions/done", `{"item":"../x;rm"}`, 400, "row id"},
		{"/api/modules/notes/actions/done", ``, 400, "row id"},
		{"/api/modules/notes/actions/broken", ``, 409, "nothing to break"},
		{"/api/modules/notes/actions/nope", ``, 404, "no action"},
		{"/api/modules/ghost/actions/done", ``, 404, "no brick"},
	} {
		w := do(t, s, call{method: "POST", path: c.path, body: c.body, header: post})
		if w.Code != c.code || !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("%s %s: %d %s, want %d %q", c.path, c.body, w.Code, w.Body.String(), c.code, c.want)
		}
	}
	log, err := os.ReadFile(s.o.ActionLog)
	if err != nil || !strings.Contains(string(log), `"target":"notes/done"`) {
		t.Errorf("action log: %v %s", err, log)
	}
}

func TestTokenLogin(t *testing.T) {
	s, _ := newServer(t)
	tok, err := EnsureToken(t.TempDir())
	if err != nil || len(tok) != 64 {
		t.Fatalf("token %q %v", tok, err)
	}
	s.o.Token = func() string { return tok }
	anon := map[string]string{IdentityHeader: ""}
	if w := do(t, s, call{path: "/api/state/inbox", header: anon}); w.Code != 403 || !strings.Contains(w.Body.String(), "crow core url") {
		t.Errorf("no cookie: %d %s", w.Code, w.Body.String())
	}
	if w := do(t, s, call{path: "/?token=wrong", header: anon}); w.Code != 403 {
		t.Errorf("wrong token: %d", w.Code)
	}
	w := do(t, s, call{path: "/?token=" + tok + "&tab=desk", header: anon})
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/?tab=desk" {
		t.Fatalf("login link: %d %v", w.Code, w.Header())
	}
	ck := w.Result().Cookies()
	if len(ck) != 1 || ck[0].Name != TokenCookie || !ck[0].HttpOnly || ck[0].SameSite != http.SameSiteStrictMode {
		t.Fatalf("cookie: %+v", ck)
	}
	if w := do(t, s, call{path: "/api/state/inbox", header: map[string]string{IdentityHeader: "", "Cookie": TokenCookie + "=" + tok}}); w.Code != 200 {
		t.Errorf("with cookie: %d", w.Code)
	}
	if w := do(t, s, call{path: "/api/state/inbox", header: map[string]string{IdentityHeader: "", "Cookie": TokenCookie + "=nope"}}); w.Code != 403 {
		t.Errorf("bad cookie: %d", w.Code)
	}
	if w := do(t, s, call{path: "/api/state/inbox", header: map[string]string{IdentityHeader: "someone@example.com", "Cookie": TokenCookie + "=" + tok}}); w.Code != 403 {
		t.Errorf("foreign tailnet login with cookie: %d, want 403", w.Code)
	}
	again, _ := EnsureToken(filepath.Dir(filepath.Join(t.TempDir(), "x")))
	if again == tok {
		t.Error("a new state dir reused the token")
	}
}

// TestModulesFixture writes GET /api/modules for a fixed brick set to the contract
// fixtures (modules.json), which the Python test checks against modules.schema.json.
func TestModulesFixture(t *testing.T) {
	notes := config.NestSection{ID: "notes", Page: "desk", Title: "Notes", View: "list", Run: []string{"crow-notes", "json"},
		Actions: []config.NestAction{
			{ID: "done", Label: "Done", Script: "crow-notes", Args: []string{"done", "{item}"}, Confirm: true, Row: true},
			{ID: "sweep", Label: "Sweep", Script: "crow-notes", Args: []string{"sweep"}},
		}}
	bricks := []config.Brick{
		{Manifest: config.Manifest{Name: "core-api", CoreSections: []string{"sessions"}, CoreRoutes: []string{"GET /api/state"}}, On: true},
		{Manifest: config.Manifest{Name: "forge", CoreSections: []string{"mrs"},
			Nest: config.NestManifest{Sections: []config.NestSection{{ID: "mrs", Page: "mrs", Source: "mrs"}}},
			Deck: config.DeckManifest{Pages: []string{"MRs"}}}, On: true},
		{Manifest: config.Manifest{Name: "mail", CoreSections: []string{"today"}, CoreRoutes: []string{"POST /api/today/{id}/done"},
			Nest: config.NestManifest{Sections: []config.NestSection{{ID: "today", Page: "home"}}},
			Deck: config.DeckManifest{Boxes: []string{"today"}}}, On: false},
		{Manifest: config.Manifest{Name: "notes", User: true, Nest: config.NestManifest{Sections: []config.NestSection{notes}}}, On: true},
	}
	m := BuildModules(bricks, config.ModuleState{Errors: []config.ModuleError{{File: "crow.yaml", Key: "modules.ghost", Error: "no brick named ghost"}}})
	if !m.sectionOff("today") || m.sectionOff("mrs") || m.sectionOff("notes") {
		t.Errorf("off sections: %v", m.offSections)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	got := indent(t, b)
	p := filepath.Join(Fixtures, "modules.json")
	if *update {
		writeFile(t, p, got)
		return
	}
	if want := testfix.Read(t, p); !bytes.Equal(got, want) {
		t.Errorf("modules differ from %s (go test -update to rewrite)\n%s", p, got)
	}
}

// TestModuleOffNoRefresh stops a section's script once its brick is switched off (item_262).
func TestModuleOffNoRefresh(t *testing.T) {
	home := brickHome(t, "version: 1\nmodules:\n  forge: false\n")
	engine, err := filepath.Abs(engineRoot)
	if err != nil {
		t.Fatal(err)
	}
	clk := &fakeClock{t: clock}
	var n atomic.Int32
	s := New(Options{Now: clk.Now, Modules: LoadModules(home, engine),
		Sources: []Source{{Name: "today", TTL: 5 * time.Second, Fetch: func() (any, error) {
			n.Add(1)
			return []any{}, nil
		}}}})
	sectionOf(t, s, "today", nil)
	idle(t, s, "today")
	clk.Add(5 * time.Second)
	s.tick()
	idle(t, s, "today")
	if n.Load() != 2 {
		t.Fatalf("runs with mail on = %d, want 2", n.Load())
	}
	writeFile(t, filepath.Join(home, "crow.yaml"), []byte("version: 1\nmodules:\n  mail: false\n  forge: false\n"))
	later := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(filepath.Join(home, "crow.yaml"), later, later); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond)
	for i := 0; i < 3; i++ {
		clk.Add(5 * time.Second)
		s.tick()
		idle(t, s, "today")
	}
	if n.Load() != 2 {
		t.Fatalf("runs with mail off = %d, want 2", n.Load())
	}
}
