// Package config resolves the paths crow-core works on. Every path comes from here, so
// core can move to another host by changing the config file, not the code.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the paths core reads and writes.
type Config struct {
	Home      string `json:"home" yaml:"-"`              // CROW_HOME, holds local/
	Missions  string `json:"missions" yaml:"missions"`   // root of missions/<id>/
	Store     string `json:"store" yaml:"store"`         // SQLite store, a view rebuilt from the event logs
	Kinds     string `json:"kinds" yaml:"kinds"`         // shipped kind files
	Overrides string `json:"overrides" yaml:"overrides"` // private kind overrides, same file names
	Projects  string `json:"projects" yaml:"projects"`   // project registry (local/projects.yaml)
	Prefs     string `json:"prefs" yaml:"prefs"`         // owner prefs for briefs, filtered by kind
	Engine    string `json:"engine" yaml:"-"`            // CROW_ENGINE, holds bin/ (crow-runs stop)
	File      string `json:"file" yaml:"-"`              // config file read, empty when none
}

// Load builds the config from the environment: defaults from CROW_HOME and CROW_ENGINE,
// then the YAML file at CROW_CORE_CONFIG (default $CROW_HOME/local/core.yaml) when it
// exists. Relative paths in the file are relative to CROW_HOME.
func Load(getenv func(string) string) (Config, error) {
	home := getenv("CROW_HOME")
	if home == "" {
		return Config{}, errors.New("CROW_HOME is not set (run crow-core through `crow core`)")
	}
	state := filepath.Join(home, "local", "state", "missions")
	c := Config{
		Home:      home,
		Missions:  state,
		Store:     filepath.Join(state, "core.db"),
		Overrides: filepath.Join(home, "local", "kinds"),
		Projects:  filepath.Join(home, "local", "projects.yaml"),
		Prefs:     filepath.Join(home, "local", "mission-prefs.yaml"),
	}
	if engine := getenv("CROW_ENGINE"); engine != "" {
		c.Engine = engine
		c.Kinds = filepath.Join(engine, "logics", "specs", "req_037_run_contract", "kinds")
	}
	file := getenv("CROW_CORE_CONFIG")
	if file == "" {
		file = filepath.Join(home, "local", "core.yaml")
	}
	b, err := os.ReadFile(file)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return c, nil
	case err != nil:
		return Config{}, err
	}
	var f Config
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(&f); err != nil && !errors.Is(err, io.EOF) {
		return Config{}, fmt.Errorf("%s: %w", file, err)
	}
	for _, p := range []struct {
		dst *string
		v   string
	}{
		{&c.Missions, f.Missions}, {&c.Store, f.Store}, {&c.Kinds, f.Kinds}, {&c.Overrides, f.Overrides}, {&c.Projects, f.Projects}, {&c.Prefs, f.Prefs},
	} {
		if p.v == "" {
			continue
		}
		if !filepath.IsAbs(p.v) {
			p.v = filepath.Join(home, p.v)
		}
		*p.dst = p.v
	}
	c.File = file
	return c, nil
}
