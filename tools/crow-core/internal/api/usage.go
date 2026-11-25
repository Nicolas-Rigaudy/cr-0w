package api

import (
	"net/http"
	"time"
)

// UsageCostSection names the cost audit source. It is not a state section: only
// GET /api/usage reads it, so /api/state never pays for the audit.
const UsageCostSection = "usage-cost"

// UsageCostTTL is how long one cost audit answers before a request starts the next.
const UsageCostTTL = 5 * time.Minute

// sessionKeys are the session fields the Usage page reads (PER SESSION, TOP RUNS).
var sessionKeys = []string{
	"id", "account", "project", "title", "kind", "models", "start", "end", "active",
	"tokens", "input_tokens", "output_tokens", "cache_read_tokens", "cache_write_tokens",
	"cache_pct", "cost_usd", "cache_saved_usd", "by_day",
}

// Usage is GET /api/usage: the cost audit (totals, by_day, by_account...) and, from the
// usage section, the account limits and the slimmed session list.
type Usage struct {
	Version  string  `json:"version"`
	Cost     Section `json:"cost"`
	Limits   any     `json:"limits"`
	Sessions []any   `json:"sessions"`
	// GeneratedAt is when crow-usage built the limits and sessions, "" before its first run.
	GeneratedAt string `json:"generated_at"`
	Stale       bool   `json:"stale,omitempty"`
	Error       string `json:"error,omitempty"`
}

func (s *Server) routeUsage(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/usage", s.usage)
}

func (s *Server) usage(w http.ResponseWriter, r *http.Request) {
	if s.mods().sectionOff("usage") {
		fail(w, http.StatusNotFound, "module off")
		return
	}
	src := s.ext["usage"]
	if src == nil {
		fail(w, http.StatusNotFound, "no usage source")
		return
	}
	out := Usage{Limits: []any{}, Sessions: []any{}}
	if e := s.ext[UsageCostSection]; e != nil {
		out.Cost = s.readExt(e)
	} else {
		out.Cost = Section{Section: UsageCostSection, Data: map[string]any{}, Stale: true, Error: "no cost audit"}
	}
	u := s.readExt(src)
	out.Stale, out.Error = u.Stale, u.Error
	if d, ok := u.Data.(map[string]any); ok {
		if l, ok := d["limits"].([]any); ok {
			out.Limits = l
		}
		if g, ok := d["generated_at"].(string); ok {
			out.GeneratedAt = g
		}
		if l, ok := d["sessions"].([]any); ok {
			out.Sessions = slimSessions(l)
		}
	}
	out.Version = hash([]any{out.Cost.Version, u.Version})
	reply(w, r, out.Version, out)
}

// slimSessions keeps the sessionKeys of each session row; rows that are not objects drop.
func slimSessions(rows []any) []any {
	out := make([]any, 0, len(rows))
	for _, row := range rows {
		m, ok := row.(map[string]any)
		if !ok {
			continue
		}
		slim := map[string]any{}
		for _, k := range sessionKeys {
			if v, ok := m[k]; ok {
				slim[k] = v
			}
		}
		out = append(out, slim)
	}
	return out
}
