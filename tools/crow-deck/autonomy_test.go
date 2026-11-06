package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// The Go table and the Autonomy table of CLAUDE.md say the same thing, row by row.
func TestAutonomyMatchesClaudeMd(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	sec := string(b)
	i := strings.Index(sec, "## Autonomy")
	if i < 0 {
		t.Fatal("CLAUDE.md has no Autonomy section")
	}
	sec = sec[i:]
	if j := strings.Index(sec[3:], "\n## "); j >= 0 {
		sec = sec[:j+3]
	}
	byRow := map[string]actKind{}
	for k, r := range autonomy {
		if r.row != "" {
			byRow[r.row] = k
		}
	}
	cell := func(c string) rule {
		c = strings.TrimSpace(c)
		switch {
		case strings.HasPrefix(c, "✅"):
			return ruleAllow
		case strings.HasPrefix(c, "📝"):
			return ruleAsk
		}
		return ruleNever // 🚫 and n/a
	}
	seen := 0
	for _, l := range strings.Split(sec, "\n") {
		f := strings.Split(strings.Trim(strings.TrimSpace(l), "|"), "|")
		if len(f) != 3 || strings.TrimSpace(f[0]) == "Action" || strings.HasPrefix(strings.TrimSpace(f[0]), "---") {
			continue
		}
		name := strings.TrimSpace(f[0])
		k, ok := byRow[name]
		if !ok {
			t.Errorf("CLAUDE.md row %q has no kind in autonomy.go", name)
			continue
		}
		seen++
		work, pers := cell(f[1]), cell(f[2])
		if strings.TrimSpace(f[2]) == "same" {
			pers = work
		}
		if got := autonomy[k]; got.work != work || got.personal != pers {
			t.Errorf("%s: code work %d personal %d, CLAUDE.md work %d personal %d", name, got.work, got.personal, work, pers)
		}
	}
	if seen != len(byRow) {
		t.Errorf("CLAUDE.md has %d rows, autonomy.go names %d", seen, len(byRow))
	}
}

func TestGateRules(t *testing.T) {
	for _, tc := range []struct {
		kind  actKind
		scope string
		want  rule
	}{
		{kPush, "work", ruleAsk}, {kPush, "personal", ruleAllow},
		{kMerge, "work", ruleNever}, {kMerge, "personal", ruleAllow},
		{kMail, "work", ruleAsk}, {kMail, "personal", ruleAsk},
		{kInfraProd, "work", ruleAsk},
		{kMerge, "", ruleNever}, {kPush, "client", ruleAsk}, // unknown scope reads as work
		{actKind("deploy"), "personal", ruleNever}, // unknown kind
		{postKind("merge"), "work", ruleNever}, {postKind("push"), "work", ruleAsk}, {postKind("reply"), "work", ruleAllow},
		{postKind("label"), "personal", ruleNever},
	} {
		if got := gate(tc.kind, tc.scope); got != tc.want {
			t.Errorf("gate(%s, %q) = %d, want %d", tc.kind, tc.scope, got, tc.want)
		}
	}
}

// runGated refuses a forbidden action even when no overlay ran.
func TestRunGatedRefusesWithoutRunning(t *testing.T) {
	calls := recordRuns(t, "ok", nil)
	if _, err := runGated(kMerge, "work", "crow-review-run", "act", "g/r", "1", "m1"); err == nil {
		t.Fatal("a work merge ran")
	}
	if len(*calls) != 0 {
		t.Fatalf("calls %+v", *calls)
	}
	if _, err := runGated(kMerge, "personal", "crow-review-run", "act", "g/r", "1", "m1"); err != nil || len(*calls) != 1 {
		t.Fatalf("personal merge: err %v calls %d", err, len(*calls))
	}
}

// A post action of kind merge on a work MR never reaches crow-review-run.
func TestPostMergeOnWorkIsRefused(t *testing.T) {
	calls := recordRuns(t, "ok", nil)
	msg := actCmd(mrRow{Repo: "g/r", Number: 3, Scope: "work"}, postAction{ID: "m", Kind: "merge"})().(postDoneMsg)
	if msg.ok || !strings.Contains(msg.text, "refused") || len(*calls) != 0 {
		t.Fatalf("msg %+v calls %+v", msg, *calls)
	}
}

// Every place that starts a process is either the gate or a read-only or navigation path
// listed here. A new caller fails until it goes through runGated or is listed with a reason.
func TestEveryProcessGoesThroughTheGate(t *testing.T) {
	allowed := map[string]string{
		"runBin":         "the runner itself",
		"runGated":       "the gate",
		"actionCmd":      "navigation helper, callers listed below",
		"attachCmd":      "navigation: focus or resume a window",
		"todayKey":       "navigation: open a mail in Outlook (actionCmd crow-open)",
		"postActionsCmd": "read only: lists pending pushes and posts",
		"fetchRelStatus": "read only: relaunch status",
		"gitRunner":      "read only: git status of a lane",
		"runJSON":        "read only: data sources",
		"editRoadmap":    "the owner edits the roadmap doc in $EDITOR",
		"editorCmd":      "the owner's editor",
	}
	starters := map[string]bool{"runBin": true, "actionCmd": true, "Command": true, "CommandContext": true, "ExecProcess": true}
	fset := token.NewFileSet()
	files, _ := filepath.Glob("*.go")
	var bad []string
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		af, err := parser.ParseFile(fset, f, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range af.Decls {
			owner, body := "", ast.Node(d)
			switch d := d.(type) {
			case *ast.FuncDecl:
				owner = d.Name.Name
			case *ast.GenDecl:
				if vs, ok := d.Specs[0].(*ast.ValueSpec); ok && len(vs.Names) > 0 {
					owner = vs.Names[0].Name
				}
			}
			ast.Inspect(body, func(n ast.Node) bool {
				c, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				name := ""
				switch fn := c.Fun.(type) {
				case *ast.Ident:
					name = fn.Name
				case *ast.SelectorExpr:
					name = fn.Sel.Name
				}
				if starters[name] {
					if _, ok := allowed[owner]; !ok {
						bad = append(bad, fset.Position(c.Pos()).String()+" "+owner+" calls "+name)
					}
				}
				return true
			})
		}
	}
	sort.Strings(bad)
	if len(bad) > 0 {
		t.Fatalf("processes started outside the gate:\n%s", strings.Join(bad, "\n"))
	}
}

// The overlay shows the exact text and command; enter runs it once, esc runs nothing.
func TestProposeShowsExactTextAndRunsOnce(t *testing.T) {
	calls := recordRuns(t, "sent", nil)
	m := newModel(defaultConfig())
	m = send(m, resize(160, 45))
	a := deckAction{Kind: kAnswer, Scope: "personal", Title: "Send", What: "Type this answer:", Text: "yes, ship it", Label: "sent",
		Cmd: []string{"crow-focus", "--pid", "7", "--send", "yes, ship it"}}
	m, _ = m.propose(a)
	if m.overlay != ovAct || !hasText(m.View(), "yes, ship it") || !hasText(m.View(), "crow-focus --pid 7 --send") {
		t.Fatalf("overlay %d view lacks the text or command", m.overlay)
	}
	c := send(m, keyMsg("esc"))
	if c.overlay != ovNone || c.act != nil || !strings.Contains(c.flash, "cancelled") {
		t.Fatalf("esc: overlay %d flash %q", c.overlay, c.flash)
	}
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = next.(model)
	runAll(cmd)
	if len(*calls) != 1 || m.act != nil || m.overlay != ovNone {
		t.Fatalf("enter: calls %d act %v overlay %d", len(*calls), m.act, m.overlay)
	}
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	runAll(cmd)
	if len(*calls) != 1 {
		t.Fatalf("a second enter ran the action again: %d calls", len(*calls))
	}
}

func TestProposeRefusesWithoutOverlay(t *testing.T) {
	m := newModel(defaultConfig())
	m, _ = m.propose(deckAction{Kind: kMerge, Scope: "work", Title: "Merge", Cmd: []string{"x"}})
	if m.overlay != ovNone || m.act != nil || !strings.Contains(m.flash, "refused: merge on work") {
		t.Fatalf("overlay %d flash %q", m.overlay, m.flash)
	}
}

// An approval-gated action says so in its overlay.
func TestAskKindsSayTheyNeedTheOK(t *testing.T) {
	m := newModel(defaultConfig())
	m = send(m, resize(160, 45))
	m, _ = m.propose(deckAction{Kind: kPush, Scope: "work", Title: "Push", What: "Push the branch.", Cmd: []string{"git", "push"}})
	if !hasText(m.View(), "needs your OK each time on work projects") {
		t.Fatal("work push overlay lacks the approval line")
	}
}
