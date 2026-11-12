package api

// Conversation pages of a live session, a dispatch run or a mission's last run
// (item_148 read path, item_177 chat). Each read runs `crow-runs conversation`, which
// asks cdx by the exact conversation id. Paging follows cdx: tail (default 50), cursor
// for what came after, before for older turns.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"crow-core/internal/mission"
)

func (s *Server) routeConversation(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sessions/{pid}/conversation", s.sessionConversation)
	mux.HandleFunc("GET /api/missions/{id}/conversation", s.missionConversation)
}

// Conversation is the reply: the cdx page plus where it was read from.
type Conversation struct {
	Account        string           `json:"account"`
	ConversationID string           `json:"conversation_id"`
	Entries        []map[string]any `json:"entries"`
	NextCursor     string           `json:"next_cursor,omitempty"`
	BeforeCursor   string           `json:"before_cursor,omitempty"`
	HasMore        bool             `json:"has_more"`
	ResetRequired  bool             `json:"reset_required,omitempty"`
}

func (s *Server) sessionConversation(w http.ResponseWriter, r *http.Request) {
	a, err := s.agent(r.PathValue("pid"))
	if err != nil {
		liveFail(w, err, http.StatusNotFound)
		return
	}
	if run := str(a, "run_id"); run != "" {
		s.conversation(w, r, "--run", run)
		return
	}
	acct, id := str(a, "account"), str(a, "session_id")
	if acct == "" || id == "" {
		fail(w, http.StatusConflict, "session has no conversation id yet")
		return
	}
	s.conversation(w, r, "--account", acct, "--conversation", id)
}

func (s *Server) missionConversation(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	d, err := s.o.Env.Detail(r.PathValue("id"))
	s.mu.Unlock()
	if errors.Is(err, mission.ErrNotFound) {
		fail(w, http.StatusNotFound, "no mission")
		return
	}
	if err != nil {
		fail(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(d.Runs) == 0 {
		fail(w, http.StatusConflict, "mission has no run yet")
		return
	}
	s.conversation(w, r, "--run", d.Runs[len(d.Runs)-1].Run)
}

// conversation runs the script with the target args and the page asked in the query.
func (s *Server) conversation(w http.ResponseWriter, r *http.Request, target ...string) {
	if s.o.Script == nil {
		fail(w, http.StatusServiceUnavailable, "no script runner (CROW_ENGINE unset)")
		return
	}
	q := r.URL.Query()
	args := append([]string{"conversation"}, target...)
	switch {
	case q.Get("cursor") != "":
		args = append(args, "--cursor", q.Get("cursor"))
	case q.Get("before") != "":
		args = append(args, "--before", q.Get("before"))
	case q.Get("tail") != "":
		n, err := strconv.Atoi(q.Get("tail"))
		if err != nil || n < 1 || n > 200 {
			fail(w, http.StatusBadRequest, "tail: 1 to 200")
			return
		}
		args = append(args, "--tail", strconv.Itoa(n))
	}
	out, err := s.o.Script(nil, "crow-runs", append(args, "--json")...)
	var page struct {
		OK    bool `json:"ok"`
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		Reason        *string          `json:"reason"`
		Account       string           `json:"account"`
		ID            string           `json:"conversation_id"`
		Entries       []map[string]any `json:"entries"`
		Next          *string          `json:"next_cursor"`
		Before        *string          `json:"before_cursor"`
		HasMore       bool             `json:"has_more"`
		ResetRequired bool             `json:"reset_required"`
	}
	if jerr := json.Unmarshal(lastLineOf(out), &page); jerr != nil {
		msg := "crow-runs conversation: no reply"
		if err != nil {
			msg = err.Error()
		}
		fail(w, http.StatusBadGateway, msg)
		return
	}
	if !page.OK {
		msg := "conversation not readable"
		switch {
		case page.Error != nil:
			msg = page.Error.Message
		case page.Reason != nil:
			msg = *page.Reason
		}
		fail(w, http.StatusConflict, msg)
		return
	}
	c := Conversation{Account: page.Account, ConversationID: page.ID, Entries: page.Entries,
		HasMore: page.HasMore, ResetRequired: page.ResetRequired}
	if c.Entries == nil {
		c.Entries = []map[string]any{}
	}
	if page.Next != nil {
		c.NextCursor = *page.Next
	}
	if page.Before != nil {
		c.BeforeCursor = *page.Before
	}
	b, _ := json.Marshal(c)
	reply(w, r, hash(b), c)
}

// lastLineOf is the last non-empty line of a script's stdout.
func lastLineOf(b []byte) []byte {
	end := len(b)
	for end > 0 && (b[end-1] == '\n' || b[end-1] == ' ' || b[end-1] == '\r') {
		end--
	}
	b = b[:end]
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] == '\n' {
			return b[i+1:]
		}
	}
	return b
}
