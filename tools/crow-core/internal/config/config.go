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
	Serve     Serve  `json:"serve" yaml:"serve"`         // crow-core serve
	File      string `json:"file" yaml:"-"`              // config file read, empty when none
}

// Serve configures the HTTP service. Listen is the loopback TCP port tailscale serve
// proxies to: every request on it must carry the owner's tailnet login. Socket is the
// unix socket local clients use without that header. Owner defaults to the login of
// this node in `tailscale status` when empty.
type Serve struct {
	Listen string `json:"listen" yaml:"listen"`
	Socket string `json:"socket" yaml:"socket"`
	Nest   string `json:"nest" yaml:"nest"` // built Nest app served at /
	Owner  string `json:"owner_login" yaml:"owner_login"`
}

// DefaultListen is the loopback address of the service.
const DefaultListen = "127.0.0.1:8787"

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
		Serve:     Serve{Listen: DefaultListen, Socket: filepath.Join(home, "local", "state", "core.sock")},
	}
	if engine := getenv("CROW_ENGINE"); engine != "" {
		c.Engine = engine
		c.Kinds = filepath.Join(engine, "logics", "specs", "req_037_run_contract", "kinds")
		c.Serve.Nest = filepath.Join(engine, "tools", "nest", "dist")
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
		{&c.Serve.Socket, f.Serve.Socket}, {&c.Serve.Nest, f.Serve.Nest},
	} {
		if p.v == "" {
			continue
		}
		if !filepath.IsAbs(p.v) {
			p.v = filepath.Join(home, p.v)
		}
		*p.dst = p.v
	}
	if f.Serve.Listen != "" {
		c.Serve.Listen = f.Serve.Listen
	}
	c.Serve.Owner = f.Serve.Owner
	c.File = file
	return c, nil
}
