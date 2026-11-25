package api

// Nest settings (item_282). GET /api/settings/nest answers modules.nest.{theme,hero,pages}
// from crow.yaml with the names Nest can build; POST writes them. Both run
// `crow-modules get|set nest`: the script is the one crow.yaml writer, core parses no YAML.
// guard has already checked the owner login, the CSRF header and the origin.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// NestKnown lists the names a setting may take: engine names plus those of the on bricks.
type NestKnown struct {
	Themes []string `json:"themes"`
	Heroes []string `json:"heroes"`
	Pages  []string `json:"pages"`
}

// NestSettings is GET /api/settings/nest and the answer of a save. Theme, hero and pages
// are what crow.yaml holds, kept as-is so a broken value still shows (errors names it).
type NestSettings struct {
	File   string            `json:"file"`
	Theme  any               `json:"theme"`
	Hero   any               `json:"hero"`
	Pages  any               `json:"pages"`
	Known  NestKnown         `json:"known"`
	Errors []json.RawMessage `json:"errors"`
}

// NestSaved answers a save: the new settings, the keys whose value changed, and whether
// Nest must be rebuilt (hero or pages changed; the theme switches live through data-theme).
type NestSaved struct {
	NestSettings
	Changed    []string `json:"changed"`
	NeedsBuild bool     `json:"needs_build"`
}

// NestSave is the body of POST /api/settings/nest. A missing key keeps its value; pages
// [] shows every page.
type NestSave struct {
	Theme *string   `json:"theme,omitempty"`
	Hero  *string   `json:"hero,omitempty"`
	Pages *[]string `json:"pages,omitempty"`
}

// settingsMu keeps one crow.yaml write at a time.
var settingsMu sync.Mutex

func (s *Server) routeSettings(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/settings/nest", s.nestSettings)
	mux.HandleFunc("POST /api/settings/nest", s.saveNestSettings)
}

func (s *Server) modulesBin() string { return filepath.Join(s.o.Engine, "bin", "crow-modules") }

func (s *Server) nestSettings(w http.ResponseWriter, r *http.Request) {
	out, _, err := s.runBrick(s.modulesBin(), []string{"get", "nest"})
	if err != nil {
		fail(w, http.StatusBadGateway, err.Error())
		return
	}
	var v NestSettings
	if err := json.Unmarshal(out, &v); err != nil {
		fail(w, http.StatusBadGateway, "crow-modules get nest: "+err.Error())
		return
	}
	reply(w, r, hash(v), v)
}

func (s *Server) saveNestSettings(w http.ResponseWriter, r *http.Request) {
	var b NestSave
	if !decode(w, r, &b, false) {
		return
	}
	args := []string{"set", "nest"}
	if b.Theme != nil {
		args = append(args, "theme", *b.Theme)
	}
	if b.Hero != nil {
		args = append(args, "hero", *b.Hero)
	}
	if b.Pages != nil {
		for _, p := range *b.Pages {
			if p == "" || strings.ContainsAny(p, ", ") {
				fail(w, http.StatusBadRequest, fmt.Sprintf("unknown page %q", p))
				return
			}
		}
		args = append(args, "pages", strings.Join(*b.Pages, ","))
	}
	if len(args) == 2 {
		fail(w, http.StatusBadRequest, "body: want theme, hero or pages")
		return
	}
	settingsMu.Lock()
	out, code, err := s.runBrick(s.modulesBin(), append(args, "--json"))
	settingsMu.Unlock()
	ev := LiveAction{TS: s.o.Now().UTC().Format(time.RFC3339), Actor: "owner", Surface: Surface,
		Action: "settings", Target: "nest", Text: strings.Join(args[2:], " ")}
	if err != nil { // die() already names the script
		err = fmt.Errorf("%s", strings.Replace(err.Error(), "crow-modules: crow-modules: ", "crow-modules: ", 1))
	}
	status := http.StatusOK
	var v NestSaved
	switch {
	case err == nil:
		if jerr := json.Unmarshal(out, &v); jerr != nil {
			status, ev.Code, ev.Message = http.StatusBadGateway, 2, "crow-modules set nest: "+jerr.Error()
		} else {
			ev.OK, ev.Message = true, "saved "+strings.Join(v.Changed, ", ")
		}
	case code == 2: // a value is unknown: nothing was written
		status, ev.Code, ev.Message = http.StatusBadRequest, 1, err.Error()
	case code == 1: // crow.yaml cannot be edited safely (flow style, read-back failed)
		status, ev.Code, ev.Message = http.StatusConflict, 1, err.Error()
	default:
		status, ev.Code, ev.Message = http.StatusBadGateway, 2, err.Error()
	}
	if lerr := s.logAction(ev); lerr != nil {
		log.Printf("settings: not logged: %v", lerr)
	}
	if status != http.StatusOK {
		fail(w, status, ev.Message)
		return
	}
	if v.Changed == nil {
		v.Changed = []string{}
	}
	if v.Errors == nil {
		v.Errors = []json.RawMessage{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
