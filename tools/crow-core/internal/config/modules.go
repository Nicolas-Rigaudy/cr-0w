package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// CoreBrick is always on (adr_002).
const CoreBrick = "core"

// yaml11Bools are the YAML 1.1 switch words. PyYAML reads them as booleans and yaml.v3
// keeps them as strings, so both loaders map them the same way.
var yaml11Bools = map[string]bool{
	"y": true, "yes": true, "on": true, "true": true,
	"n": false, "no": false, "off": false, "false": false,
}

// Manifest is the part of an engine module manifest (modules/<brick>.yaml) the loader reads.
type Manifest struct {
	Name         string         `yaml:"name"`
	Summary      string         `yaml:"summary"`
	Default      bool           `yaml:"default"`
	Requires     []string       `yaml:"requires"`
	Config       map[string]any `yaml:"config"`
	CoreSections []string       `yaml:"core_sections"`
	CoreRoutes   []string       `yaml:"core_routes"`
	Nest         NestManifest   `yaml:"nest"`
	Deck         DeckManifest   `yaml:"deck"`
	User         bool           `yaml:"-"`
	Dir          string         `yaml:"-"` // folder of a user brick (its bin/ holds the scripts)
}

// NestManifest lists the Nest sections a brick adds (item_237).
type NestManifest struct {
	Sections []NestSection `yaml:"sections" json:"sections"`
}

// NestSection is one Nest panel. Page is home, desk, ops or reports (or a brick's built-in
// page). Source is the core section it reads. View is "list" for the generic list view;
// empty means a view built into Nest. Run makes core serve the section from the brick's
// bin/ script (user bricks): [script, args...], its stdout JSON is the section data.
type NestSection struct {
	ID      string       `yaml:"id" json:"id"`
	Page    string       `yaml:"page" json:"page"`
	Source  string       `yaml:"source" json:"source"`
	Title   string       `yaml:"title" json:"title,omitempty"`
	View    string       `yaml:"view" json:"view,omitempty"`
	Run     []string     `yaml:"run" json:"-"`
	TTL     string       `yaml:"ttl" json:"-"`
	Actions []NestAction `yaml:"actions" json:"actions,omitempty"`
}

// NestAction is a button of a section. Core runs Script (bin/ of the brick) with Args
// through POST /api/modules/<brick>/actions/<id>; "{item}" in an arg is the row id.
type NestAction struct {
	ID      string   `yaml:"id" json:"id"`
	Label   string   `yaml:"label" json:"label"`
	Script  string   `yaml:"script" json:"-"`
	Args    []string `yaml:"args" json:"-"`
	Confirm bool     `yaml:"confirm" json:"confirm"`
	Row     bool     `yaml:"-" json:"row"` // an arg takes {item}: the button sits on each row
}

// DeckManifest lists the deck pages and Home boxes a brick owns.
type DeckManifest struct {
	Pages []string `yaml:"pages" json:"pages"`
	Boxes []string `yaml:"boxes" json:"boxes"`
}

// Module is the state of one brick after crow.yaml.
type Module struct {
	On      bool           `json:"on"`
	Default bool           `json:"default"`
	Config  map[string]any `json:"config"`
}

// ModuleError is a crow.yaml problem: the file, the dotted key ("" for the whole file) and the message.
type ModuleError struct {
	File  string `json:"file"`
	Key   string `json:"key"`
	Error string `json:"error"`
}

// ModuleState is the merge of the engine manifests and crow.yaml.
type ModuleState struct {
	File    string            `json:"file"`
	Exists  bool              `json:"exists"`
	Modules map[string]Module `json:"modules"`
	Errors  []ModuleError     `json:"errors"`
}

// EngineManifests reads <engine>/modules/*.yaml, keyed by brick name.
func EngineManifests(engine string) (map[string]Manifest, error) {
	paths, err := filepath.Glob(filepath.Join(engine, "modules", "*.yaml"))
	if err != nil {
		return nil, err
	}
	out := map[string]Manifest{}
	for _, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var m Manifest
		if err := yaml.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if m.Name == "" {
			m.Name = strings.TrimSuffix(filepath.Base(p), ".yaml")
		}
		out[m.Name] = m
	}
	return out, nil
}

// UserManifest is the file name of a user brick manifest: <user modules dir>/<name>/module.yaml.
const UserManifest = "module.yaml"

// UserModulesDir is where user bricks live: <home>/modules, or <home>/local/modules when home
// is the engine checkout itself (its modules/ holds the engine manifests).
func UserModulesDir(home, engine string) string {
	if sameDir(home, engine) {
		return filepath.Join(home, "local", "modules")
	}
	return filepath.Join(home, "modules")
}

func sameDir(a, b string) bool {
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	if errA != nil || errB != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return ra == rb
}

// UserManifests reads the user bricks (item_265), as crow_lib.user_manifests does. A folder
// with no manifest, a broken manifest, a name that differs from its folder or a clash with
// an engine brick is listed in the errors and left out. Default is on when the manifest
// does not say.
func UserManifests(home, engine string, engineNames map[string]Manifest) (map[string]Manifest, []ModuleError) {
	root := UserModulesDir(home, engine)
	out, errs := map[string]Manifest{}, []ModuleError{}
	entries, err := os.ReadDir(root)
	if err != nil {
		return out, errs
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") || !e.IsDir() {
			continue
		}
		path := filepath.Join(root, name, UserManifest)
		fail := func(msg string) { errs = append(errs, ModuleError{path, "modules." + name, msg}) }
		raw, err := os.ReadFile(path)
		if err != nil {
			fail("user brick has no " + UserManifest + ", ignored")
			continue
		}
		var top map[string]any
		if err := yaml.Unmarshal(raw, &top); err != nil {
			fail("manifest does not parse: " + strings.SplitN(err.Error(), "\n", 2)[0])
			continue
		}
		var m Manifest
		_ = yaml.Unmarshal(raw, &m)
		if m.Name != "" && m.Name != name {
			fail(fmt.Sprintf("manifest name '%s' differs from its folder, ignored", m.Name))
			continue
		}
		if _, clash := engineNames[name]; clash {
			fail(fmt.Sprintf("user brick '%s' clashes with the engine brick of the same name, ignored (rename its folder)", name))
			continue
		}
		if _, has := top["default"]; !has {
			m.Default = true
		}
		m.Name, m.User, m.Dir = name, true, filepath.Join(root, name)
		out[name] = m
	}
	return out, errs
}

// CrowYAMLPath is crow.yaml at the CROW_HOME root when present, else local/crow.yaml.
func CrowYAMLPath(home string) string {
	root := filepath.Join(home, "crow.yaml")
	if _, err := os.Stat(root); err == nil {
		return root
	}
	return filepath.Join(home, "local", "crow.yaml")
}

// LoadModules merges the adr_002 defaults, crow.yaml and the engine manifests, as
// crow_lib.load_modules does. A broken file or key is listed in Errors and keeps the
// default; core stays on. The error is for unreadable engine manifests only.
func LoadModules(home, engine string) (ModuleState, error) {
	st, _, err := loadAll(home, engine)
	return st, err
}

// Brick is one brick with its switch and manifest.
type Brick struct {
	Manifest
	On bool
}

// Bricks returns every brick sorted by name, with the crow.yaml state and its errors.
func Bricks(home, engine string) ([]Brick, ModuleState, error) {
	st, manifests, err := loadAll(home, engine)
	if err != nil {
		return nil, st, err
	}
	out := make([]Brick, 0, len(manifests))
	for _, n := range sortedKeys(manifests) {
		m := manifests[n]
		for i, sec := range m.Nest.Sections {
			for j, a := range sec.Actions {
				for _, arg := range a.Args {
					if strings.Contains(arg, "{item}") {
						m.Nest.Sections[i].Actions[j].Row = true
					}
				}
			}
		}
		out = append(out, Brick{Manifest: m, On: st.Modules[n].On})
	}
	return out, st, nil
}

func loadAll(home, engine string) (ModuleState, map[string]Manifest, error) {
	manifests, err := EngineManifests(engine)
	if err != nil {
		return ModuleState{}, nil, err
	}
	users, userErrs := UserManifests(home, engine, manifests)
	for n, m := range users {
		manifests[n] = m
	}
	path := CrowYAMLPath(home)
	st := ModuleState{File: path, Modules: map[string]Module{}, Errors: userErrs}
	for n, m := range manifests {
		conf := map[string]any{}
		for k, v := range m.Config {
			conf[k] = plain(v)
		}
		st.Modules[n] = Module{On: m.Default || n == CoreBrick, Default: m.Default, Config: conf}
	}
	fail := func(key, msg string) { st.Errors = append(st.Errors, ModuleError{path, key, msg}) }

	var data any
	raw, err := os.ReadFile(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err != nil:
		st.Exists = true
		fail("", "does not read: "+err.Error())
	default:
		st.Exists = true
		if err := yaml.Unmarshal(raw, &data); err != nil {
			fail("", "does not parse: "+strings.SplitN(err.Error(), "\n", 2)[0])
			data = nil
		}
	}
	top, ok := asMap(data)
	if data != nil && !ok {
		fail("", "expects a mapping with version and modules")
	}
	for _, k := range sortedKeys(top) {
		if k != "version" && k != "modules" {
			fail(k, "unknown key (expected version, modules)")
		}
	}
	if v, has := top["version"]; has && v != 1 {
		fail("version", fmt.Sprintf("unsupported version %v (expected 1)", v))
	}
	given, ok := asMap(top["modules"])
	if top["modules"] != nil && !ok {
		fail("modules", "expects a mapping of brick to true, false or a mapping")
	}
	for _, name := range sortedKeys(given) {
		val, key := given[name], "modules."+name
		mod, known := st.Modules[name]
		if !known {
			fail(key, "unknown brick (known: "+strings.Join(sortedKeys(st.Modules), ", ")+")")
			continue
		}
		on, isSwitch := switchOf(val)
		conf := map[string]any{}
		if body, isMap := asMap(val); !isSwitch && isMap {
			for _, k := range sortedKeys(body) {
				switch v := body[k]; {
				case k == "on":
					if on, isSwitch = switchOf(v); !isSwitch {
						fail(key+".on", "expects true or false")
					}
				case !hasKey(mod.Config, k):
					allowed := append([]string{"on"}, sortedKeys(mod.Config)...)
					fail(key+"."+k, "unknown key (allowed: "+strings.Join(allowed, ", ")+")")
				default:
					conf[k] = plain(v)
				}
			}
			if !isSwitch {
				on = mod.On
			}
		} else if !isSwitch {
			fail(key, "expects true, false or a mapping with on and config keys")
			continue
		}
		if name == CoreBrick && !on {
			fail(key, "core is always on")
			on = true
		}
		mod.On = on
		for k, v := range conf {
			mod.Config[k] = v
		}
		st.Modules[name] = mod
	}
	for _, name := range sortedKeys(st.Modules) {
		for _, req := range manifests[name].Requires {
			if st.Modules[name].On && !st.Modules[req].On {
				fail("modules."+name, "requires "+req+", which is off")
				if manifests[name].User {
					m := st.Modules[name]
					m.On = false // a user brick never runs without what it needs
					st.Modules[name] = m
				}
			}
		}
	}
	sort.Slice(st.Errors, func(i, j int) bool {
		a, b := st.Errors[i], st.Errors[j]
		return a.Key < b.Key || a.Key == b.Key && a.Error < b.Error
	})
	return st, manifests, nil
}

// Modules is the on set with config: brick -> config for every brick that is on.
func Modules(home, engine string) (map[string]map[string]any, error) {
	st, err := LoadModules(home, engine)
	if err != nil {
		return nil, err
	}
	out := map[string]map[string]any{}
	for n, m := range st.Modules {
		if m.On {
			out[n] = m.Config
		}
	}
	return out, nil
}

// switchOf reads a brick switch: a bool or a YAML 1.1 word, case-insensitive.
func switchOf(v any) (bool, bool) {
	switch s := v.(type) {
	case bool:
		return s, true
	case string:
		b, ok := yaml11Bools[strings.ToLower(s)]
		return b, ok
	}
	return false, false
}

// asMap reads a YAML mapping with keys as strings (a bool key becomes on/off, as PyYAML reads them).
func asMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, val := range m {
			out[keyString(k)] = val
		}
		return out, true
	}
	return nil, false
}

func keyString(k any) string {
	if b, ok := k.(bool); ok {
		if b {
			return "on"
		}
		return "off"
	}
	return fmt.Sprint(k)
}

// plain turns nested YAML mappings into map[string]any so config encodes as JSON.
func plain(v any) any {
	if m, ok := asMap(v); ok {
		out := make(map[string]any, len(m))
		for k, val := range m {
			out[k] = plain(val)
		}
		return out
	}
	if l, ok := v.([]any); ok {
		out := make([]any, len(l))
		for i, val := range l {
			out[i] = plain(val)
		}
		return out
	}
	return v
}

func hasKey[V any](m map[string]V, k string) bool {
	_, ok := m[k]
	return ok
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
