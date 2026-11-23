package api

// Login on the TCP listener (item_237). Behind tailscale serve the owner's tailnet login
// in Tailscale-User-Login lets a request in, as before. Without Tailscale the core token
// does: opening /?token=<token> once sets an HttpOnly cookie and redirects to the same
// page without the token; every later request carries the cookie.

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// TokenCookie holds the core token in the browser.
const TokenCookie = "crow_token"

// TokenFile is the core token, in the machine state dir.
const TokenFile = "core.token"

type loginResult int

const (
	loginRefused loginResult = iota
	loginOK
	loginRedirected
)

func (s *Server) token() string {
	if s.o.Token == nil {
		return ""
	}
	return s.o.Token()
}

func (s *Server) tokenOn() bool { return s.token() != "" }

func sameToken(a, b string) bool {
	return a != "" && b != "" && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// login checks the tailnet identity, then the token cookie, then a ?token= link.
func (s *Server) login(w http.ResponseWriter, r *http.Request) loginResult {
	if id := r.Header.Get(IdentityHeader); id != "" {
		if owner := s.o.Owner(); owner != "" && id == owner {
			return loginOK
		}
		return loginRefused
	}
	tok := s.token()
	if tok == "" {
		return loginRefused
	}
	if c, err := r.Cookie(TokenCookie); err == nil && sameToken(c.Value, tok) {
		return loginOK
	}
	q := r.URL.Query()
	if (r.Method == http.MethodGet || r.Method == http.MethodHead) && sameToken(q.Get("token"), tok) {
		http.SetCookie(w, &http.Cookie{
			Name: TokenCookie, Value: tok, Path: "/", MaxAge: 400 * 24 * 3600,
			HttpOnly: true, SameSite: http.SameSiteStrictMode,
			Secure: r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https"),
		})
		q.Del("token")
		u := *r.URL
		u.RawQuery = q.Encode()
		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, u.RequestURI(), http.StatusSeeOther)
		return loginRedirected
	}
	return loginRefused
}

// EnsureToken reads the token file in dir, or writes a new random one (0600).
func EnsureToken(dir string) (string, error) {
	p := filepath.Join(dir, TokenFile)
	b, err := os.ReadFile(p)
	if err == nil {
		if t := strings.TrimSpace(string(b)); len(t) >= 32 {
			return t, nil
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	t := hex.EncodeToString(raw)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(p, []byte(t+"\n"), 0o600); err != nil {
		return "", err
	}
	return t, nil
}
