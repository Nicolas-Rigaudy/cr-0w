// Package brief renders brief.md, the file a mission run reads first (item_107). It
// follows logics/specs/req_037_run_contract/brief.template.md and is built from the kind
// file, the mission record, the owner prefs and every answers file. core rewrites it
// before each run, so a cold run on any account restarts from files alone.
//
// The brief names no account and no local path: the repo and worktree are placeholders
// the run resolves from its working directory.
package brief

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"crow-core/contract"
)

// Pref is one owner preference. A pref with no kinds applies to every kind.
type Pref struct {
	Rule  string   `yaml:"rule"`
	Kinds []string `yaml:"kinds"`
}

// LoadPrefs reads the prefs file (a YAML list of prefs). A missing file gives none.
func LoadPrefs(path string) ([]Pref, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var ps []Pref
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(&ps); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	for i, p := range ps {
		if strings.TrimSpace(p.Rule) == "" {
			return nil, fmt.Errorf("%s: pref %d has no rule", path, i+1)
		}
	}
	return ps, nil
}

// ForKind keeps the prefs that apply to a kind, in file order.
func ForKind(ps []Pref, kind string) []string {
	var out []string
	for _, p := range ps {
		if len(p.Kinds) == 0 || slices.Contains(p.Kinds, kind) {
			out = append(out, strings.TrimSpace(p.Rule))
		}
	}
	return out
}

// Answer is one answers.vN.md: the owner's verb and words for version N.
type Answer struct {
	Version int
	Verb    string
	Words   string // the body after the front matter, verbatim (trimmed)
}

var answersFile = regexp.MustCompile(`^answers\.v([0-9]+)\.md$`)

// ReadAnswers reads every answers.vN.md of a mission dir, in version order.
func ReadAnswers(dir string) ([]Answer, error) {
	es, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []Answer
	for _, e := range es {
		m := answersFile.FindStringSubmatch(e.Name())
		if m == nil || e.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		a, err := parseAnswer(b)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		if n, _ := strconv.Atoi(m[1]); a.Version != n {
			return nil, fmt.Errorf("%s: front matter says version %d", e.Name(), a.Version)
		}
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	return out, nil
}

func parseAnswer(b []byte) (Answer, error) {
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return Answer{}, errors.New("no front matter")
	}
	head, body, ok := strings.Cut(s[4:], "\n---")
	if !ok {
		return Answer{}, errors.New("front matter not closed")
	}
	var fm struct {
		Version int    `yaml:"version"`
		Verb    string `yaml:"verb"`
	}
	if err := yaml.Unmarshal([]byte(head), &fm); err != nil {
		return Answer{}, err
	}
	if fm.Version < 1 || fm.Verb == "" {
		return Answer{}, errors.New("front matter needs version and verb")
	}
	return Answer{Version: fm.Version, Verb: fm.Verb, Words: strings.TrimSpace(body)}, nil
}

// Input is everything one brief is built from.
type Input struct {
	Mission   contract.Mission
	Key       string // from the created event; derived from the id when empty
	Dir       string // absolute mission folder; the brief names it so the run writes there
	Kind      *contract.Kind
	Prefs     []string          // already filtered by kind
	Answers   []Answer          // in version order
	Preflight map[string]string // preflight answers by question id; nil when skipped
}

// Version is the version the next run writes: one past the last answer.
func Version(answers []Answer) int {
	n := 0
	for _, a := range answers {
		n = max(n, a.Version)
	}
	return n + 1
}

// Phase is the file the next run writes: the first gate before the report the owner has
// not passed with `go`, or the report.
func Phase(gates []string, answers []Answer) string {
	passed := 0
	for _, a := range answers {
		if a.Verb == "go" {
			passed++
		}
	}
	for _, g := range gates {
		if g != "plan" && g != "mockup" {
			continue
		}
		if passed == 0 {
			return g
		}
		passed--
	}
	return "report"
}

// Key is the short part of a mission id (YYYYMMDD-kind-key). The id lowercases the key,
// so callers pass the created event's key when they have it.
func Key(m contract.Mission) string {
	_, rest, ok := strings.Cut(m.ID, "-")
	if !ok {
		return m.ID
	}
	return strings.TrimPrefix(rest, m.Kind+"-")
}

// Render writes the brief.
func Render(in Input) string {
	m, k := in.Mission, in.Kind
	n := Version(in.Answers)
	phase := Phase(m.Gates, in.Answers)
	file := fmt.Sprintf("%s.v%d.md", phase, n)
	worktree := k.Worktree != nil && *k.Worktree

	var b strings.Builder
	line := func(f string, a ...any) { fmt.Fprintf(&b, f+"\n", a...) }

	line("# Brief: %s", m.ID)
	line("")
	line("<!-- Written by core. The run reads it first and never edits it. -->")
	line("")
	line("## Goal")
	line("%s", orNone(m.Goal, "(no goal recorded)"))
	line("")

	line("## Mission")
	line("- kind: %s", m.Kind)
	if len(m.Gates) == 0 {
		line("- gates: none (straight through)")
	} else {
		line("- gates: %s", strings.Join(m.Gates, ", "))
	}
	key := in.Key
	if key == "" {
		key = Key(m)
	}
	line("- key: %s", key)
	line("- scope: %s", orNone(m.Scope, "none"))
	if worktree {
		line("- repo: `<repo>` on branch `mission/%s` (worktree `<repo>.wt/%s`)", m.ID, m.ID)
	}
	caps := fmt.Sprintf("%d min, %s USD, %d resumes", m.Budget.TimeMin, num(m.Budget.BudgetUSD), m.Budget.Resumes)
	if k.Caps.Rounds != nil {
		caps += fmt.Sprintf(", %d rounds", *k.Caps.Rounds)
	}
	line("- caps: %s", caps)
	line("- current version: %d (write `%s`)", n, file)
	if in.Dir != "" {
		line("- mission folder: `%s`", in.Dir)
	}
	line("")

	line("## Owner preferences")
	if len(in.Prefs) == 0 {
		line("None.")
	}
	for _, p := range in.Prefs {
		line("- %s", p)
	}
	line("")

	line("## Preflight")
	switch {
	case len(k.Preflight) == 0:
		line("None for this kind.")
	case in.Preflight == nil:
		line("Skipped: use the safest defaults.")
		for _, q := range k.Preflight {
			line("- %s %s.", q.Text, q.Safest)
		}
	default:
		for _, q := range k.Preflight {
			a, ok := in.Preflight[q.ID]
			if !ok {
				a = q.Safest + " (safest default)"
			}
			line("- %s %s.", q.Text, strings.TrimSuffix(a, "."))
		}
	}
	line("")

	line("## Answers so far")
	if len(in.Answers) == 0 {
		line("None yet.")
	}
	for _, a := range in.Answers {
		if a.Words == "" {
			line("- v%d: %s.", a.Version, a.Verb)
			continue
		}
		line("- v%d: %s.", a.Version, a.Verb)
		for _, w := range strings.Split(a.Words, "\n") {
			if w = strings.TrimRight(w, " \t"); w == "" {
				line("")
			} else {
				line("  %s", w)
			}
		}
	}
	line("")

	line("## Constraints")
	inDir := ""
	if in.Dir != "" {
		inDir = " in the mission folder"
	}
	if worktree {
		line("- Allowed: read and edit inside the worktree, run the repo checks, commit on the mission branch.")
		line("- Not allowed: push, merge, post, install packages, touch files outside the worktree.")
	} else {
		line("- Allowed: %s", readOnlyAllowed(k))
		line("- Not allowed: push, merge, post, install packages, write anywhere but `outbox/` and the version file%s.", inDir)
	}
	dc := "none"
	if k.DataClasses != nil && len(*k.DataClasses) > 0 {
		dc = strings.Join(*k.DataClasses, ", ")
	}
	line("- Data classes allowed: %s.", dc)
	line("")

	line("## Reply contract")
	at := func(p string) string { return p }
	if in.Dir != "" {
		at = func(p string) string { return filepath.Join(in.Dir, p) }
		line("- Paths below are in the mission folder `%s`, not in the repo. List `artifacts` relative to it.", in.Dir)
	}
	line("- Write exactly one file for this turn: `%s`, with one fenced `crow` block (schema `crow-block.schema.json`, contract %d).", at(file), contract.Version)
	line("- Block shape (YAML or JSON; drop the keys you do not use; list entries are objects, never bare strings):")
	line("%s", blockShape)
	outbox := fmt.Sprintf("- Put any file you want shown under `%s/` and list it in `artifacts`.", at("outbox"))
	if k.Caps.OutboxFileMB != nil && k.Caps.OutboxTotalMB != nil {
		outbox += fmt.Sprintf(" Caps: %s MB per file, %s MB in total.", num(*k.Caps.OutboxFileMB), num(*k.Caps.OutboxTotalMB))
	}
	line("%s", outbox)
	line("- A choice you cannot make: add it to `questions` with `blocking: true`, set `verdict: needs-input`, and end the turn. Outside a plan or report, write `%s`.", at("outbox/questions.json"))
	line("- `questions.json` shape (JSON; `choices`, `default` and `recommended` are optional; ids are `q<N>`, choice ids lowercase):")
	line("%s", questionsShape(m.ID, n))
	line("- Never wait for an answer inside the run. End the turn. core resumes you with `answers.v%d.md`.", n)
	line("- From version 2, start the file with a `## Responses` section and fill `responses`.")
	line("- A report with `verdict: ready` carries `proof`.")
	line("- List every outward action in `publish` with its exact text. core publishes only what the owner validated.")
	return b.String()
}

// blockShape shows every list of the crow block with its object form: runs that only
// see the schema name guess bare strings (e2e 2026-11-10).
const blockShape = "```yaml\n" +
	"contract: 1\nmission: <mission id>\nkind: <kind>\nphase: <plan|report|mockup>\nversion: <N>\n" +
	"summary: <one line>\nconfidence: <0-100>\nverdict: <ready|needs-input|blocked>\n" +
	"questions:\n  - {id: q1, text: <question>, blocking: true, choices: [{id: a, label: <choice>}, {id: b, label: <choice>}]}\n" +
	"responses:\n  - {ref: a1, outcome: done, note: <quote of the answer>}\n" +
	"artifacts:\n  - {path: outbox/<file>}\n" +
	"publish:\n  - {action: <action>, target: <where>, body: <exact text>}\n" +
	"proof:\n  - {what: command, cmd: <command>, exit: 0}\n" +
	"```"

// questionsShape shows outbox/questions.json with the mission and version filled in: runs
// that only see the file name guess crow-block fields and core rejects the file (unknown
// field kind).
func questionsShape(id string, n int) string {
	return "```json\n" +
		fmt.Sprintf(`{"contract": %d, "mission": %q, "version": %d, "questions": [`, contract.Version, id, n) + "\n" +
		`  {"id": "q1", "text": "<question>", "blocking": true, "choices": [{"id": "a", "label": "<choice>"}, {"id": "b", "label": "<choice>"}], "recommended": "a"}` + "\n" +
		"]}\n" +
		"```"
}

func readOnlyAllowed(k *contract.Kind) string {
	if k.Permission != nil && *k.Permission == "review" {
		return "read the repo and the diff under review, run read-only checks."
	}
	return "read the repo and the web, run read-only commands."
}

func orNone(s, none string) string {
	if s = strings.TrimSpace(s); s == "" {
		return none
	}
	return s
}

// num prints a float without a trailing .0.
func num(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Write renders the brief into <dir>/brief.md, replacing the previous one.
func Write(dir string, in Input) error {
	tmp := filepath.Join(dir, ".brief.md.tmp")
	if err := os.WriteFile(tmp, []byte(Render(in)), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, "brief.md"))
}
