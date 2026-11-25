package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// notes maps a note kind to its vault name and its legacy local/state name (crow_lib.NOTES).
var notes = map[string][2]string{
	"rollups": {"Rollups", "weekly"},
	"retros":  {"Retros", "retros"},
}

// NotesDir is the folder for owner-facing notes, as crow_lib.crow_dir: <home>/notes when
// the vault brick is off, vault.crow_dir (under vault.path when relative) otherwise, ""
// when unset. vault is the vault root when known.
func NotesDir(c Config) (dir, vault string) {
	if st, err := LoadModules(c.Home, c.Engine); err == nil {
		if m, ok := st.Modules["vault"]; ok && !m.On {
			return filepath.Join(c.Home, "notes"), ""
		}
	}
	var reg struct {
		Vault struct {
			Path    string `yaml:"path"`
			CrowDir string `yaml:"crow_dir"`
		} `yaml:"vault"`
	}
	b, err := os.ReadFile(c.Projects)
	if err != nil || yaml.Unmarshal(b, &reg) != nil || reg.Vault.CrowDir == "" {
		return "", ""
	}
	vault = expand(reg.Vault.Path)
	dir = expand(reg.Vault.CrowDir)
	if !filepath.IsAbs(dir) && vault != "" {
		dir = filepath.Join(vault, dir)
	}
	return dir, vault
}

// NotePath is crow_lib.crow_path for the kinds core reads.
func NotePath(c Config, kind string) string {
	n := notes[kind]
	if dir, _ := NotesDir(c); dir != "" {
		return filepath.Join(dir, n[0])
	}
	return filepath.Join(c.Home, "local", "state", n[1])
}

func expand(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if h, err := os.UserHomeDir(); err == nil {
			return filepath.Join(h, strings.TrimPrefix(p, "~"))
		}
	}
	return p
}
