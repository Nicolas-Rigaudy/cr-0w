package api

// Bricks in core (item_237, adr_002). GET /api/modules lists every brick with its switch,
// the core sections it serves and what it adds to Nest and the deck. An off brick's
// sections leave /api/state and its routes answer 404 "module off". A section a manifest
// declares with actions gets POST /api/modules/<brick>/actions/<id>, which runs the
// brick's bin/ script under the script contract. crow.yaml is read again when it changes.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"crow-core/contract"
	"crow-core/internal/config"
)

// ModuleBrick is one row of GET /api/modules.
type ModuleBrick struct {
	Name     string              `json:"name"`
	On       bool                `json:"on"`
	User     bool                `json:"user"`
	Sections []string            `json:"sections"` // core sections it serves
	Nest     config.NestManifest `json:"nest"`
	Deck     config.DeckManifest `json:"deck"`
	routes   []string
	dir      string
}

// Modules is GET /api/modules.
type Modules struct {
	ContractVersion int                  `json:"contract_version"`
	Version         string               `json:"version"` // changes when a switch or manifest changes
	Modules         []ModuleBrick        `json:"modules"`
	Errors          []config.ModuleError `json:"errors"`

	offSections map[string]string // section -> its off brick
	offRoutes   *http.ServeMux    // routes of off bricks, nil when none
	bricks      map[string]*ModuleBrick
}

// ServerBrick serves the API itself: its sections and routes are never gated.
const ServerBrick = "core-api"

// BuildModules turns the bricks into the API answer and the gate tables.
func BuildModules(bricks []config.Brick, st config.ModuleState) *Modules {
	m := &Modules{ContractVersion: contract.Version, Modules: []ModuleBrick{}, Errors: st.Errors,
		offSections: map[string]string{}, bricks: map[string]*ModuleBrick{}}
	if m.Errors == nil {
		m.Errors = []config.ModuleError{}
	}
	var off []string
	for _, b := range bricks {
		row := ModuleBrick{Name: b.Name, On: b.On, User: b.User, Sections: []string{},
			Nest: b.Nest, Deck: b.Deck, routes: b.CoreRoutes, dir: b.Dir}
		row.Nest.Sections = make([]config.NestSection, len(b.Nest.Sections))
		for i, sec := range b.Nest.Sections {
			sec.Source = SectionOf(sec) // Nest reads /api/state/<source>
			row.Nest.Sections[i] = sec
		}
		row.Deck.Pages, row.Deck.Boxes = orEmpty(row.Deck.Pages), orEmpty(row.Deck.Boxes)
		row.Sections = append(row.Sections, b.CoreSections...)
		for _, sec := range b.Nest.Sections {
			if len(sec.Run) > 0 {
				row.Sections = append(row.Sections, SectionOf(sec))
			}
		}
		if !b.On && b.Name != ServerBrick {
			for _, sec := range row.Sections {
				m.offSections[sec] = b.Name
			}
			off = append(off, b.CoreRoutes...)
		}
		m.Modules = append(m.Modules, row)
	}
	for i := range m.Modules {
		m.bricks[m.Modules[i].Name] = &m.Modules[i]
	}
	if len(off) > 0 {
		m.offRoutes = http.NewServeMux()
		for _, p := range off {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("modules: route %q not gated: %v", p, r)
					}
				}()
				m.offRoutes.HandleFunc(p, func(http.ResponseWriter, *http.Request) {})
			}()
		}
	}
	m.Version = hash(m)
	return m
}

func orEmpty(l []string) []string {
	if l == nil {
		return []string{}
	}
	return l
}

// SectionOf is the core section a Nest section reads: its source, else its id.
func SectionOf(sec config.NestSection) string {
	if sec.Source != "" {
		return sec.Source
	}
	return sec.ID
}

// routeOff names the off brick that owns the request's route, "" when none.
func (m *Modules) routeOff(r *http.Request) bool {
	if m == nil || m.offRoutes == nil {
		return false
	}
	_, pat := m.offRoutes.Handler(r)
	return pat != ""
}

// sectionOff reports whether the section belongs to an off brick.
func (m *Modules) sectionOff(name string) bool {
	if m == nil {
		return false
	}
	_, off := m.offSections[name]
	return off
}

// LoadModules reads the bricks once and again whenever crow.yaml or a manifest changes,
// checked at most once a second.
func LoadModules(home, engine string) func() (*Modules, error) {
	var mu sync.Mutex
	var last *Modules
	var key string
	var checked time.Time
	return func() (*Modules, error) {
		mu.Lock()
		defer mu.Unlock()
		if last != nil && time.Since(checked) < time.Second {
			return last, nil
		}
		checked = time.Now()
		k := filesKey(home, engine)
		if last != nil && k == key {
			return last, nil
		}
		bricks, st, err := config.Bricks(home, engine)
		if err != nil {
			return last, err
		}
		last, key = BuildModules(bricks, st), k
		return last, nil
	}
}

// filesKey is the size and mtime of crow.yaml and every manifest.
func filesKey(home, engine string) string {
	paths := []string{filepath.Join(home, "crow.yaml"), filepath.Join(home, "local", "crow.yaml")}
	eng, _ := filepath.Glob(filepath.Join(engine, "modules", "*.yaml"))
	usr, _ := filepath.Glob(filepath.Join(config.UserModulesDir(home, engine), "*", config.UserManifest))
	var b strings.Builder
	for _, p := range append(append(paths, eng...), usr...) {
		if fi, err := os.Stat(p); err == nil {
			fmt.Fprintf(&b, "%s %d %d;", p, fi.Size(), fi.ModTime().UnixNano())
		}
	}
	return b.String()
}

// mods is the current brick state, nil without a loader or after a failed first read.
func (s *Server) mods() *Modules {
	if s.o.Modules == nil {
		return nil
	}
	m, err := s.o.Modules()
	if err != nil {
		log.Printf("modules: %v", err)
	}
	return m
}

func (s *Server) routeModules(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/modules", s.modules)
	mux.HandleFunc("POST /api/modules/{brick}/actions/{id}", s.moduleAction)
}

func (s *Server) modules(w http.ResponseWriter, r *http.Request) {
	m := s.mods()
	if m == nil {
		fail(w, http.StatusServiceUnavailable, "modules not loaded")
		return
	}
	reply(w, r, m.Version, m)
}

// ModuleActionBody is the optional body of POST /api/modules/{brick}/actions/{id}.
type ModuleActionBody struct {
	Item string `json:"item,omitempty"`
}

var (
	itemRe   = regexp.MustCompile(`^[A-Za-z0-9._:@/+=~][A-Za-z0-9._:@/+=~ -]{0,199}$`)
	scriptRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
)

// ActionTimeout bounds one manifest action.
const ActionTimeout = 90 * time.Second

func (s *Server) moduleAction(w http.ResponseWriter, r *http.Request) {
	m := s.mods()
	if m == nil {
		fail(w, http.StatusServiceUnavailable, "modules not loaded")
		return
	}
	b := m.bricks[r.PathValue("brick")]
	if b == nil {
		fail(w, http.StatusNotFound, "no brick "+r.PathValue("brick"))
		return
	}
	if !b.On {
		fail(w, http.StatusNotFound, "module off")
		return
	}
	id := r.PathValue("id")
	var act config.NestAction
	var section string
	for _, sec := range b.Nest.Sections {
		for _, a := range sec.Actions {
			if a.ID == id && act.ID == "" {
				act, section = a, SectionOf(sec)
			}
		}
	}
	if act.ID == "" {
		fail(w, http.StatusNotFound, "no action "+id)
		return
	}
	var body ModuleActionBody
	if raw, _ := io.ReadAll(io.LimitReader(r.Body, 4096)); len(strings.TrimSpace(string(raw))) > 0 {
		if err := json.Unmarshal(raw, &body); err != nil {
			fail(w, http.StatusBadRequest, "body: "+err.Error())
			return
		}
	}
	if act.Row && !itemRe.MatchString(body.Item) {
		fail(w, http.StatusBadRequest, "item: expects a row id")
		return
	}
	if !scriptRe.MatchString(act.Script) {
		fail(w, http.StatusInternalServerError, "action "+id+": script must be a bin/ name")
		return
	}
	bin := filepath.Join(s.o.Engine, "bin", act.Script)
	if b.User {
		bin = filepath.Join(b.dir, "bin", act.Script)
	}
	args := make([]string, len(act.Args))
	for i, a := range act.Args {
		args[i] = strings.ReplaceAll(a, "{item}", body.Item)
	}
	ev := LiveAction{TS: s.o.Now().UTC().Format(time.RFC3339), Actor: "owner", Surface: "nest",
		Action: "module_action", Target: b.Name + "/" + id, Text: body.Item}
	out, code, err := s.runBrick(bin, args)
	switch {
	case err == nil:
		ev.OK, ev.Message = true, strings.TrimSpace(string(tailLine(out)))
		if ev.Message == "" {
			ev.Message = act.Label + " done"
		}
	case code == 2:
		ev.Code, ev.Message = 1, err.Error()
	default:
		ev.Code, ev.Message = 2, err.Error()
	}
	if lerr := s.logAction(ev); lerr != nil {
		log.Printf("module action %s: not logged: %v", ev.Target, lerr)
	}
	if section != "" {
		go s.Refresh(section)
	}
	status := http.StatusOK
	switch ev.Code {
	case 1:
		status = http.StatusConflict
	case 2:
		status = http.StatusBadGateway
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(ev)
}

// runBrick runs a brick script with CROW_HOME and CROW_ENGINE set. It returns stdout, the
// exit code (-1 when it did not run) and an error carrying the end of stderr.
func (s *Server) runBrick(bin string, args []string) ([]byte, int, error) {
	if s.o.Engine == "" && s.o.Home == "" {
		return nil, -1, errors.New("no CROW_ENGINE: brick actions are off")
	}
	ctx, cancel := context.WithTimeout(context.Background(), ActionTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = append(os.Environ(), "CROW_HOME="+s.o.Home, "CROW_ENGINE="+s.o.Engine)
	out, err := cmd.Output()
	if err == nil {
		return out, 0, nil
	}
	name := filepath.Base(bin)
	if ctx.Err() != nil {
		return out, -1, fmt.Errorf("%s: timed out", name)
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		msg := strings.TrimSpace(string(tailLine(ee.Stderr)))
		if msg == "" {
			msg = err.Error()
		}
		return out, ee.ExitCode(), fmt.Errorf("%s: %s", name, msg)
	}
	return out, -1, fmt.Errorf("%s: %w", name, err)
}
