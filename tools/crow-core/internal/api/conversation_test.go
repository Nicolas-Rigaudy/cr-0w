package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

const convPage = `{"schema":1,"account":"spare","conversation_id":"c-9","ok":true,"status":"ok","entries":[{"role":"user","text":"hi"},{"role":"assistant","text":"hello"}],"next_cursor":"n1","before_cursor":"b1","has_more":true}`

func TestSessionConversation(t *testing.T) {
	s, calls, _ := liveServer(t, func(scriptCall) ([]byte, error) { return []byte(convPage + "\n"), nil })
	for _, tc := range []struct {
		path string
		args string
	}{
		{"/api/sessions/4242/conversation", "conversation --account main --conversation s-1 --json"},
		{"/api/sessions/5151/conversation?tail=10", "conversation --run a1b2c3d4 --tail 10 --json"},
		{"/api/sessions/5151/conversation?cursor=n1&tail=5", "conversation --run a1b2c3d4 --cursor n1 --json"},
		{"/api/sessions/5151/conversation?before=b1", "conversation --run a1b2c3d4 --before b1 --json"},
	} {
		*calls = nil
		w := do(t, s, call{path: tc.path})
		if w.Code != http.StatusOK {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body)
		}
		if len(*calls) != 1 || (*calls)[0].Name != "crow-runs" || strings.Join((*calls)[0].Args, " ") != tc.args {
			t.Fatalf("%s: calls %+v", tc.path, *calls)
		}
		var c Conversation
		if err := json.Unmarshal(w.Body.Bytes(), &c); err != nil || len(c.Entries) != 2 || c.NextCursor != "n1" || c.BeforeCursor != "b1" || !c.HasMore || c.ConversationID != "c-9" {
			t.Fatalf("%s: body %s", tc.path, w.Body)
		}
	}
}

func TestConversationErrors(t *testing.T) {
	var out []byte
	var err error
	s, _, _ := liveServer(t, func(scriptCall) ([]byte, error) { return out, err })
	for _, tc := range []struct {
		path string
		out  string
		err  error
		code int
	}{
		{"/api/sessions/9999/conversation", convPage, nil, http.StatusNotFound},
		{"/api/sessions/5151/conversation?tail=0", convPage, nil, http.StatusBadRequest},
		{"/api/sessions/5151/conversation", `{"ok":false,"error":{"code":"no_conversation","message":"run has no conversation yet"}}`, errors.New("exit 1"), http.StatusConflict},
		{"/api/sessions/5151/conversation", `boom`, errors.New("exit 2"), http.StatusBadGateway},
		{"/api/missions/nope/conversation", convPage, nil, http.StatusNotFound},
	} {
		out, err = []byte(tc.out), tc.err
		if w := do(t, s, call{path: tc.path}); w.Code != tc.code {
			t.Fatalf("%s %q: %d %s", tc.path, tc.out, w.Code, w.Body)
		}
	}
}

func TestMissionConversation(t *testing.T) {
	s, calls, _ := liveServer(t, func(scriptCall) ([]byte, error) { return []byte(convPage), nil })
	s.mu.Lock()
	d, err := s.o.Env.Detail("20261109-build-item-095")
	s.mu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	w := do(t, s, call{path: "/api/missions/20261109-build-item-095/conversation"})
	if len(d.Runs) == 0 {
		if w.Code != http.StatusConflict {
			t.Fatalf("no runs: %d %s", w.Code, w.Body)
		}
		return
	}
	if w.Code != http.StatusOK {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	want := "conversation --run " + d.Runs[len(d.Runs)-1].Run + " --json"
	if got := strings.Join((*calls)[0].Args, " "); got != want {
		t.Fatalf("args %q want %q", got, want)
	}
}
