package main

// History actions: enter resumes the session in a new window, o opens its transcript, r queues
// a retro, y copies its id, h hands it to a repo session. Resume, retro and handoff confirm
// first; transcript and copy only show something.

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

// histAction answers an action key on the selected session.
func (m model) histAction(k string) (model, tea.Cmd) {
	s := m.histView().cur(m.histView().rows())
	if s == nil {
		m.flash = "no session selected"
		return m, nil
	}
	if strings.HasPrefix(s.ID, "mock-") {
		m.flash = k + ": sample row, nothing to act on"
		return m, nil
	}
	switch k {
	case "enter":
		if s.Active {
			m.flash = "this session is still running: jump to it from Agents"
			return m, nil
		}
		if s.Cwd == "" || s.Account == "" {
			m.flash = "cannot resume: no workspace or account recorded"
			return m, nil
		}
		return m.propose(deckAction{
			Kind: kSession, Scope: s.Scope, Title: "Resume",
			What:  "Resume this session in a new window on " + s.Account + ", in " + s.Cwd + ":",
			Text:  s.Task,
			Note:  "Same account, same transcript. Nothing is sent to the session.",
			Label: "resume",
			Cmd:   []string{"crow-focus", "--pid", "0", "--account", s.Account, "--cwd", s.Cwd, "--session", s.ID},
		})
	case "o":
		if s.Path == "" {
			m.flash = "no transcript recorded for this session"
			return m, nil
		}
		return m, openTranscript(s.Path)
	case "r":
		if s.Path == "" {
			m.flash = "no transcript recorded for this session"
			return m, nil
		}
		return m.propose(deckAction{
			Kind: kDispatch, Scope: s.Scope, Title: "Retro",
			What:  "Queue a retro of this session (one worker at a time, report in the retro folder):",
			Text:  s.Task,
			Label: "retro",
			Cmd:   []string{"crow-retro", "add", s.Path},
		})
	case "y":
		copyText(s.ID)
		m.flash = "copied " + s.ID
		return m, nil
	case "h":
		task := "Pick up session " + s.ID + " (" + s.Task + "). Read its transcript first: " + cmpOr(s.Path, "resume it with cdx --resume "+s.ID) + "."
		end := "its open work done, or listed as open threads"
		a := m.handoffAction("Handoff note", s.Project, "Open a repo session on "+s.Project+" with a note on this session.", task, end)
		a.Label = "handoff note"
		return m.propose(a)
	}
	return m, nil
}

// copyText puts text on the clipboard through the terminal (OSC 52); tests replace it.
var copyText = termenv.Copy

// openTranscript shows the transcript in the owner's pager, less without one.
func openTranscript(path string) tea.Cmd {
	return tea.ExecProcess(pagerCmd(path), func(err error) tea.Msg {
		if err != nil {
			return actionMsg{"pager failed: " + err.Error()}
		}
		return nil
	})
}

// histDetail fills the live preview from the transcript: touched files from crow-usage and
// the last messages read from the file's tail.
func histDetail(s *histSession) (touched []string, last []tailLine) {
	touched = topFiles(s.Files, 4)
	if s.Path != "" {
		last = transcriptTail(s.Path, 3)
	}
	return touched, last
}

// topFiles lists the most edited files first, base names only.
func topFiles(files map[string]int, n int) []string {
	var out []string
	for f := range files {
		out = append(out, f)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if files[out[i]] != files[out[j]] {
			return files[out[i]] > files[out[j]]
		}
		return out[i] < out[j]
	})
	if len(out) > n {
		out = out[:n]
	}
	for i, f := range out {
		out[i] = baseName(f)
	}
	return out
}

func baseName(p string) string {
	if i := strings.LastIndex(p, "/"); i >= 0 {
		return p[i+1:]
	}
	return p
}

var tailCache = struct {
	sync.Mutex
	m map[string]tailEntry
}{m: map[string]tailEntry{}}

type tailEntry struct {
	mod   time.Time
	lines []tailLine
}

// transcriptTail is the last n user or assistant texts of a transcript, newest last. It reads
// the file's last 256 KiB once per change.
func transcriptTail(path string, n int) []tailLine {
	st, err := os.Stat(path)
	if err != nil {
		return nil
	}
	tailCache.Lock()
	defer tailCache.Unlock()
	if e, ok := tailCache.m[path]; ok && e.mod.Equal(st.ModTime()) {
		return e.lines
	}
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	const span = 256 << 10
	if st.Size() > span {
		if _, err := f.Seek(st.Size()-span, io.SeekStart); err != nil {
			return nil
		}
	}
	var out []tailLine
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 8<<20)
	for sc.Scan() {
		if t, ok := transcriptLine(sc.Bytes()); ok {
			out = append(out, t)
		}
	}
	if len(out) > n {
		out = out[len(out)-n:]
	}
	tailCache.m[path] = tailEntry{st.ModTime(), out}
	return out
}

// transcriptLine is one user or assistant text turn; tool calls, tool results and meta lines
// are skipped.
func transcriptLine(b []byte) (tailLine, bool) {
	var e struct {
		Type      string `json:"type"`
		Timestamp string `json:"timestamp"`
		IsMeta    bool   `json:"isMeta"`
		Message   struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal(b, &e) != nil || e.IsMeta || (e.Type != "user" && e.Type != "assistant") {
		return tailLine{}, false
	}
	text := ""
	var s string
	if json.Unmarshal(e.Message.Content, &s) == nil {
		text = s
	} else {
		var blocks []struct{ Type, Text string }
		if json.Unmarshal(e.Message.Content, &blocks) != nil {
			return tailLine{}, false
		}
		for _, bl := range blocks {
			if bl.Type == "text" && bl.Text != "" {
				text = bl.Text
				break
			}
		}
	}
	text = strings.Join(strings.Fields(text), " ")
	if text == "" || strings.HasPrefix(text, "<") {
		return tailLine{}, false
	}
	at := ""
	if t, err := time.Parse(time.RFC3339, e.Timestamp); err == nil {
		at = t.Local().Format("15:04")
	}
	who := "you "
	if e.Type == "assistant" {
		who = "› "
	}
	return tailLine{At: at, Text: who + text}, true
}
