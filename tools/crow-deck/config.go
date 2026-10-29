package main

// Config: a small key = value file. Keys: theme, accent, mode.
// Lookup: $CROW_DECK_CONFIG, then $CROW_HOME/local/deck.conf, then <user config dir>/crow/deck.conf.
// Flags override the file for the session; `t` writes the theme back.

import (
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Theme  string
	Accent int // slot, -1 = theme default
	Mode   string
}

func defaultConfig() config { return config{Theme: "lazytf", Accent: -1, Mode: "work"} }

func configPath() string {
	if p := os.Getenv("CROW_DECK_CONFIG"); p != "" {
		return p
	}
	if h := os.Getenv("CROW_HOME"); h != "" {
		return filepath.Join(h, "local", "deck.conf")
	}
	if d, err := os.UserConfigDir(); err == nil {
		return filepath.Join(d, "crow", "deck.conf")
	}
	return ""
}

func parseConfig(text string) config {
	c := defaultConfig()
	for _, raw := range strings.Split(text, "\n") {
		raw = strings.TrimSpace(raw)
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}
		k, v, ok := strings.Cut(raw, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		switch k {
		case "theme":
			if themeIndex(v) >= 0 {
				c.Theme = strings.ToLower(v)
			}
		case "accent":
			if n, ok := parseAccent(v); ok {
				c.Accent = n
			}
		case "mode":
			if validMode(v) {
				c.Mode = v
			}
		}
	}
	return c
}

func validMode(s string) bool { return s == "all" || s == "work" || s == "perso" }

func loadConfig(path string) config {
	if path == "" {
		return defaultConfig()
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig()
	}
	return parseConfig(string(b))
}

// saveKey sets one key in the file, keeping every other line.
func saveKey(path, key, value string) error {
	if path == "" {
		return nil
	}
	var lines []string
	if b, err := os.ReadFile(path); err == nil {
		lines = strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	}
	done := false
	for i, l := range lines {
		if k, _, ok := strings.Cut(l, "="); ok && strings.TrimSpace(k) == key {
			lines[i], done = key+" = "+value, true
		}
	}
	if !done {
		lines = append(lines, key+" = "+value)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
