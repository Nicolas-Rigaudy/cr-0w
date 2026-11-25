package api

// Request log (item_279): one line per write request under /api, refused ones included,
// so an action from Nest always shows in the service journal.

import (
	"net/http"
	"strings"
	"time"
)

// statusWriter keeps the status code a handler writes.
type statusWriter struct {
	http.ResponseWriter
	code int
}

func (w *statusWriter) WriteHeader(code int) {
	if w.code == 0 {
		w.code = code
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.code == 0 {
		w.code = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Unwrap lets http.ResponseController reach the real writer (flush, deadlines).
func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// logRequests logs method, path, status and duration of every non-GET /api request.
// A nil logf logs nothing.
func (s *Server) logRequests(trusted bool, next http.Handler) http.Handler {
	logf := s.o.Logf
	if logf == nil {
		return next
	}
	via := "tcp"
	if trusted {
		via = "unix"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		start := s.o.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		code := sw.code
		if code == 0 {
			code = http.StatusOK
		}
		logf("api %s %s %d %dms (%s)", r.Method, r.URL.EscapedPath(), code,
			s.o.Now().Sub(start)/time.Millisecond, via)
	})
}
