package mission

// Owner answers, response coverage and the validation record (req_037, item_108).
//
// An answer writes answers.vN.md for the version the mission waits on. Further words on
// the same version are appended to that file as new points, and the latest verb wins.
// ok (and go or approve where the block publishes) writes one validation event per
// publish entry, chained by hash. Any later owner message drops the live validations.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"crow-core/contract"
	"crow-core/internal/eventlog"
)

// Verbs are the owner verbs of each gate file.
var Verbs = map[string][]string{
	"plan":   {"go", "dig", "no"},
	"report": {"ok", "msg", "abandon"},
	"mockup": {"pick", "change", "approve", "abandon"},
}

// validating verbs write validations when the block has a publish list. A blocking open
// question refuses them.
var validating = []string{"ok", "go", "approve"}

var surfaces = []string{"nest", "deck", "cli"}

// Reply is one owner message on the version a mission waits on.
type Reply struct {
	Verb    string
	Words   []string // new points, one per entry
	Surface string   // nest, deck or cli
	Version int      // mission version the owner answered from; 0 skips the check
}

// AnswerResult is what one answer wrote.
type AnswerResult struct {
	File        string   `json:"file"`
	Version     int      `json:"version"`
	Verb        string   `json:"verb"`
	Points      int      `json:"points"`
	Validations []string `json:"validations"` // actions validated
	Dropped     []string `json:"dropped"`     // actions whose validation was dropped
}

// AnswersFile is the name of the answers file of version n.
func AnswersFile(n int) string { return fmt.Sprintf("answers.v%d.md", n) }

// Answer records an owner reply on the version a mission in review waits on.
func (env Env) Answer(ref string, r Reply) (AnswerResult, error) {
	var res AnswerResult
	m, err := env.find(ref)
	if err != nil {
		return res, err
	}
	if m.State != "review" {
		return res, failf(ErrConflict, "%s is %s; an answer needs a mission in review", m.ID, m.State)
	}
	return env.reply(m, r, nil)
}

// reply writes the answer of m on its current version. An inbox answer passes its row id;
// the event then carries the row and the surface.
func (env Env) reply(m contract.Mission, r Reply, item *string) (AnswerResult, error) {
	var res AnswerResult
	if err := checkSurface(r.Surface); err != nil {
		return res, err
	}
	n := m.Version
	dir := filepath.Join(env.Store.Root(), m.ID)
	file, cb, why := versionFile(dir, m.ID, n)
	if cb == nil {
		return res, fmt.Errorf("%s v%d: %s", m.ID, n, why)
	}
	// An inbox answer was checked against its row, which knows the FYI shape.
	phase, _, _ := strings.Cut(file, ".")
	if item == nil && !slices.Contains(Verbs[phase], r.Verb) {
		return res, failf(ErrInvalid, "%s: verb %q, want one of %v", file, r.Verb, Verbs[phase])
	}

	words := make([]string, 0, len(r.Words))
	for _, w := range r.Words {
		if w = strings.Join(strings.Fields(w), " "); w != "" {
			words = append(words, w)
		}
	}
	points, err := readPoints(filepath.Join(dir, AnswersFile(n)))
	if err != nil {
		return res, err
	}
	points = append(points, words...)

	publish := []contract.Publish{}
	if cb.Publish != nil {
		publish = *cb.Publish
	}
	validates := slices.Contains(validating, r.Verb) && len(publish) > 0
	if slices.Contains(validating, r.Verb) {
		if open := openBlocking(cb, points); len(open) > 0 {
			return res, fmt.Errorf("%s: blocking question %s is open; answer it as \"%s: ...\"",
				file, strings.Join(open, ", "), open[0])
		}
	}

	at := env.Now()
	ts := at.Format(time.RFC3339)
	if err := writeAnswers(dir, m.ID, n, r.Verb, r.Surface, ts, points); err != nil {
		return res, err
	}
	res = AnswerResult{File: AnswersFile(n), Version: n, Verb: r.Verb, Points: len(points),
		Validations: []string{}, Dropped: []string{}}

	log := eventlog.Path(env.Store.Root(), m.ID)
	live, err := liveValidations(log)
	if err != nil {
		return res, err
	}
	// The surface is in the answers file front matter, as in the spec example. An inbox
	// answer also names its row.
	ad := contract.AnswerData{File: res.File, Verb: r.Verb}
	if item != nil {
		surface := r.Surface
		ad.Item, ad.Surface = item, &surface
	}
	if err := env.append(m.ID, "owner", "answer", ts, ad); err != nil {
		return res, err
	}
	for _, v := range live {
		d := v.Data
		d.DroppedBy = &ts
		if err := env.append(m.ID, "owner", "validation_dropped", ts, d); err != nil {
			return res, err
		}
		res.Dropped = append(res.Dropped, d.Action)
	}
	if !validates {
		return res, nil
	}
	for i, p := range publish {
		sum, err := contract.PayloadSHA256(p)
		if err != nil {
			return res, err
		}
		prev, err := lastValidationSHA(log)
		if err != nil {
			return res, err
		}
		vts := at.Add(time.Duration(i+1) * time.Second).Format(time.RFC3339)
		d := contract.ValidationData{File: file, Version: n, Verb: r.Verb, Action: p.Action, Target: p.Target,
			PayloadSHA256: sum, OwnerEvent: ts, PrevSHA256: prev, Surface: r.Surface}
		if err := env.append(m.ID, "owner", "validation", vts, d); err != nil {
			return res, err
		}
		res.Validations = append(res.Validations, p.Action)
	}
	return res, nil
}

// openBlocking lists the blocking questions of the block no answer point starts with.
func openBlocking(cb *contract.CrowBlock, points []string) []string {
	if cb.Questions == nil {
		return nil
	}
	var open []string
	for _, q := range *cb.Questions {
		if q.Blocking && !slices.ContainsFunc(points, answersQuestion(q.ID)) {
			open = append(open, q.ID)
		}
	}
	return open
}

// answersQuestion matches a point that answers question id ("q1: refuse.").
func answersQuestion(id string) func(string) bool {
	return func(p string) bool { return strings.HasPrefix(p, id+":") }
}

var pointLine = regexp.MustCompile(`^([0-9]+)\.\s+(.*)$`)

// Points parses the numbered list of an answers body: point i is a<i+1>. A line that is
// not a numbered item continues the point above.
func Points(body string) []string {
	var out []string
	for _, l := range strings.Split(body, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if m := pointLine.FindStringSubmatch(l); m != nil {
			out = append(out, m[2])
		} else if len(out) > 0 {
			out[len(out)-1] += " " + l
		}
	}
	return out
}

// readPoints reads the points of an answers file; none when it does not exist.
func readPoints(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return nil, fmt.Errorf("%s: no front matter", filepath.Base(path))
	}
	_, body, ok := strings.Cut(s[4:], "\n---")
	if !ok {
		return nil, fmt.Errorf("%s: front matter not closed", filepath.Base(path))
	}
	return Points(strings.TrimPrefix(body, "\n")), nil
}

// writeAnswers writes answers.vN.md through a temp file and a rename.
func writeAnswers(dir, id string, n int, verb, surface, ts string, points []string) error {
	var b bytes.Buffer
	fmt.Fprintf(&b, "---\nmission: %s\nversion: %d\nverb: %s\nsurface: %s\nts: %s\n", id, n, verb, surface, ts)
	b.WriteString("---\n")
	for i, p := range points {
		fmt.Fprintf(&b, "%d. %s\n", i+1, p)
	}
	path := filepath.Join(dir, AnswersFile(n))
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b.Bytes(), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Coverage checks that the responses of version n (n >= 2) cover every answer point and
// every question of version n-1. A question also counts as covered through the answer
// point that answers it ("q1: refuse." covered as a1).
func Coverage(dir, id string, n int, cb *contract.CrowBlock) error {
	if n < 2 {
		return nil
	}
	refs := map[string]bool{}
	if cb.Responses != nil {
		for _, r := range *cb.Responses {
			refs[r.Ref] = true
		}
	}
	points, err := readPoints(filepath.Join(dir, AnswersFile(n-1)))
	if err != nil {
		return err
	}
	var missing []string
	for i := range points {
		if a := "a" + strconv.Itoa(i+1); !refs[a] {
			missing = append(missing, a)
		}
	}
	if _, prev, _ := versionFile(dir, id, n-1); prev != nil && prev.Questions != nil {
		for _, q := range *prev.Questions {
			if refs[q.ID] {
				continue
			}
			covered := false
			for i, p := range points {
				if answersQuestion(q.ID)(p) && refs["a"+strconv.Itoa(i+1)] {
					covered = true
				}
			}
			if !covered {
				missing = append(missing, q.ID)
			}
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("v%d responses miss %s", n, strings.Join(missing, ", "))
	}
	return nil
}

// rawLine is one event log line kept as stored, for the hash chain.
type rawLine struct {
	No    int
	Bytes []byte
	Event contract.Event
	Data  contract.ValidationData
}

// rawLines reads every line of a mission log as stored, without the newline.
func rawLines(path string) ([]rawLine, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []rawLine
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64<<10), 16<<20)
	for no := 1; sc.Scan(); no++ {
		b := slices.Clone(sc.Bytes())
		var e contract.Event
		if err := json.Unmarshal(b, &e); err != nil {
			return nil, fmt.Errorf("%s line %d: %w", filepath.Base(path), no, err)
		}
		l := rawLine{No: no, Bytes: b, Event: e}
		if e.Type == "validation" || e.Type == "validation_dropped" {
			if err := json.Unmarshal(e.Data, &l.Data); err != nil {
				return nil, fmt.Errorf("%s line %d: %w", filepath.Base(path), no, err)
			}
		}
		out = append(out, l)
	}
	return out, sc.Err()
}

// lastValidationSHA is the prev_sha256 of the next validation of a mission.
func lastValidationSHA(path string) (string, error) {
	ls, err := rawLines(path)
	if err != nil {
		return "", err
	}
	prev := contract.ZeroSHA256
	for _, l := range ls {
		if l.Event.Type == "validation" {
			prev = contract.SHA256Hex(l.Bytes)
		}
	}
	return prev, nil
}

// liveValidations are the validations not yet dropped or used by a publish.
func liveValidations(path string) ([]rawLine, error) {
	ls, err := rawLines(path)
	if err != nil {
		return nil, err
	}
	var live []rawLine
	for _, l := range ls {
		switch l.Event.Type {
		case "validation":
			live = append(live, l)
		case "validation_dropped":
			live = slices.DeleteFunc(live, func(v rawLine) bool {
				return v.Data.OwnerEvent == l.Data.OwnerEvent && v.Data.Action == l.Data.Action &&
					v.Data.PayloadSHA256 == l.Data.PayloadSHA256
			})
		case "publish":
			var p contract.PublishData
			if err := json.Unmarshal(l.Event.Data, &p); err != nil {
				return nil, fmt.Errorf("%s line %d: %w", filepath.Base(path), l.No, err)
			}
			live = slices.DeleteFunc(live, func(v rawLine) bool { return v.Event.TS == p.ValidationTS })
		}
	}
	return live, nil
}

// CheckChain checks the validation record of a mission log: each validation is written by
// the owner, names the owner answer that carried its verb, and its prev_sha256 is the
// SHA-256 of the validation line before it as stored. An edited or replayed line breaks it.
func CheckChain(path string) error {
	ls, err := rawLines(path)
	if err != nil {
		return err
	}
	answers := map[string]string{} // owner answer ts -> verb
	prev := contract.ZeroSHA256
	for _, l := range ls {
		e := l.Event
		switch e.Type {
		case "answer":
			var a contract.AnswerData
			if err := json.Unmarshal(e.Data, &a); err != nil {
				return fmt.Errorf("line %d: %w", l.No, err)
			}
			if e.Actor == "owner" {
				answers[e.TS] = a.Verb
			}
		case "validation":
			d := l.Data
			switch {
			case e.Actor != "owner":
				return fmt.Errorf("line %d: validation by %s; only the owner validates", l.No, e.Actor)
			case d.PrevSHA256 != prev:
				return fmt.Errorf("line %d: validation %s prev_sha256 %s, want %s", l.No, d.Action, d.PrevSHA256, prev)
			case answers[d.OwnerEvent] != d.Verb:
				return fmt.Errorf("line %d: validation %s names owner event %s with no owner %s",
					l.No, d.Action, d.OwnerEvent, d.Verb)
			}
			prev = contract.SHA256Hex(l.Bytes)
		}
	}
	return nil
}

// Validated is one live validation checked against the version file on disk.
type Validated struct {
	TS   string                  `json:"ts"`
	Data contract.ValidationData `json:"data"`
}

// Verify checks the live validations of a mission before core publishes: the chain holds,
// every publish entry of the validated file has one, and each payload hash matches the
// entry on disk. A mismatch writes a validation_mismatch incident and refuses; the mission
// stays parked at its gate until a new ok.
func (env Env) Verify(ref string) ([]Validated, error) {
	m, err := env.find(ref)
	if err != nil {
		return nil, err
	}
	log := eventlog.Path(env.Store.Root(), m.ID)
	refuse := func(detail string) ([]Validated, error) {
		ts := env.Now().Format(time.RFC3339)
		if err := env.append(m.ID, "core", "incident", ts,
			contract.IncidentData{Cause: "validation_mismatch", Detail: &detail}); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", m.ID, detail)
	}
	if err := CheckChain(log); err != nil {
		return refuse(err.Error())
	}
	live, err := liveValidations(log)
	if err != nil {
		return nil, err
	}
	if len(live) == 0 {
		return nil, fmt.Errorf("%s: no live validation", m.ID)
	}
	file, version := live[0].Data.File, live[0].Data.Version
	cb, err := blockOf(filepath.Join(env.Store.Root(), m.ID, file))
	if err != nil {
		return refuse(err.Error())
	}
	var entries []contract.Publish
	if cb.Publish != nil {
		entries = *cb.Publish
	}
	if len(entries) != len(live) {
		return refuse(fmt.Sprintf("%s lists %d publish entries, %d validated", file, len(entries), len(live)))
	}
	out := make([]Validated, len(live))
	for i, v := range live {
		d := v.Data
		if d.File != file || d.Version != version {
			return refuse(fmt.Sprintf("live validations name %s and %s", file, d.File))
		}
		sum, err := contract.PayloadSHA256(entries[i])
		if err != nil {
			return nil, err
		}
		if entries[i].Action != d.Action || sum != d.PayloadSHA256 {
			return refuse(fmt.Sprintf("%s publish entry %d (%s) changed since the validation of %s",
				file, i+1, entries[i].Action, v.Event.TS))
		}
		out[i] = Validated{TS: v.Event.TS, Data: d}
	}
	return out, nil
}

func blockOf(path string) (*contract.CrowBlock, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return contract.ParseCrowBlock(b)
}
