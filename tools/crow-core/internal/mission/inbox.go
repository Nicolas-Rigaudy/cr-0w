package mission

// The owner inbox (req_028 need 3, item_134).
//
// The store folds the rows from events: one per parked question and one per version file
// at a gate. The ask, the blocking flag and the artifacts come from the files on disk when
// the inbox is read. An answer goes through the answers file and the validation record of
// item_108, closes the row and moves the mission: back to queued for a resume (item_135)
// or to its end state.

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"crow-core/contract"
)

// askMax is the longest ask of inbox.schema.json, in runes.
const askMax = 200

// Inbox lists the inbox rows, oldest first. Only open rows unless all is set.
func (env Env) Inbox(all bool) ([]contract.InboxItem, error) {
	if _, err := env.Store.SyncAll(); err != nil {
		return nil, err
	}
	d, err := env.Store.Dump()
	if err != nil {
		return nil, err
	}
	out := []contract.InboxItem{}
	for _, it := range d.Inbox {
		if all || it.State == "open" {
			out = append(out, env.enrich(it))
		}
	}
	slices.SortStableFunc(out, func(a, b contract.InboxItem) int {
		if c := strings.Compare(a.Created, b.Created); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return out, nil
}

// enrich fills the ask, the blocking flag and the artifacts of a row from its file.
func (env Env) enrich(it contract.InboxItem) contract.InboxItem {
	dir := filepath.Join(env.Store.Root(), it.Mission)
	it.Ask, it.Artifacts = it.File, []string{}
	if it.Source == "questions" {
		if q := questionOf(dir, it); q != nil {
			it.Ask, it.Blocking = q.Text, q.Blocking
			it.Artifacts = previews(&contract.Questions{Questions: []contract.Question{*q}})
		}
	} else if cb, err := blockOf(filepath.Join(dir, it.File)); err == nil {
		if s := strings.Join(strings.Fields(cb.Summary), " "); s != "" {
			it.Ask = s
		}
		points, _ := readPoints(filepath.Join(dir, AnswersFile(it.Version)))
		it.Blocking = len(openBlocking(cb, points)) > 0
		it.Artifacts = artifactPaths(cb)
	}
	if it.Artifacts == nil {
		it.Artifacts = []string{}
	}
	if r := []rune(it.Ask); len(r) > askMax {
		it.Ask = string(r[:askMax-1]) + "…"
	}
	return it
}

// questionsOf reads outbox/questions.json when it belongs to version n of the mission.
func questionsOf(dir, id string, n int) *contract.Questions {
	b, err := os.ReadFile(filepath.Join(dir, QuestionsFile))
	if err != nil {
		return nil
	}
	q := &contract.Questions{}
	if contract.Decode(b, q) != nil || q.Mission != id || q.Version != n {
		return nil
	}
	return q
}

// questionOf is the question a questions row asks.
func questionOf(dir string, it contract.InboxItem) *contract.Question {
	q := questionsOf(dir, it.Mission, it.Version)
	if q == nil {
		return nil
	}
	qid := it.ID[strings.LastIndex(it.ID, "/")+1:]
	for i := range q.Questions {
		if q.Questions[i].ID == qid {
			return &q.Questions[i]
		}
	}
	return nil
}

// findItem finds a row by id or unique id prefix; open rows only unless all is set.
func (env Env) findItem(ref string, all bool) (contract.InboxItem, error) {
	items, err := env.Inbox(true)
	if err != nil {
		return contract.InboxItem{}, err
	}
	var hits []contract.InboxItem
	for _, it := range items {
		if it.ID == ref {
			return it, nil
		}
		if strings.HasPrefix(it.ID, ref) && (all || it.State == "open") {
			hits = append(hits, it)
		}
	}
	switch len(hits) {
	case 1:
		return hits[0], nil
	case 0:
		return contract.InboxItem{}, fmt.Errorf("no inbox item %q", ref)
	}
	ids := make([]string, len(hits))
	for i, it := range hits {
		ids[i] = it.ID
	}
	return contract.InboxItem{}, fmt.Errorf("%q matches %d items: %s", ref, len(hits), strings.Join(ids, ", "))
}

// ItemDetail is what `crow inbox show` prints.
type ItemDetail struct {
	Item      contract.InboxItem  `json:"item"`
	Path      string              `json:"path"`
	Summary   string              `json:"summary"`
	Questions []contract.Question `json:"questions"`
	Answers   []string            `json:"answers"`
	Responses []contract.Response `json:"responses"`
	Publish   []contract.Publish  `json:"publish"`
	Artifacts []contract.Artifact `json:"artifact_entries"`
}

// InboxShow returns one row with the block summary, questions, options and artifacts.
func (env Env) InboxShow(ref string) (ItemDetail, error) {
	it, err := env.findItem(ref, true)
	if err != nil {
		return ItemDetail{}, err
	}
	dir := filepath.Join(env.Store.Root(), it.Mission)
	d := ItemDetail{Item: it, Path: filepath.Join(dir, it.File), Questions: []contract.Question{},
		Responses: []contract.Response{}, Publish: []contract.Publish{}, Artifacts: []contract.Artifact{}}
	d.Answers, _ = readPoints(filepath.Join(dir, AnswersFile(it.Version)))
	if d.Answers == nil {
		d.Answers = []string{}
	}
	if it.Source == "questions" {
		if q := questionOf(dir, it); q != nil {
			d.Questions = append(d.Questions, *q)
		}
		return d, nil
	}
	cb, err := blockOf(d.Path)
	if err != nil {
		return d, nil
	}
	d.Summary = cb.Summary
	if cb.Questions != nil {
		d.Questions = *cb.Questions
	}
	if cb.Responses != nil {
		d.Responses = *cb.Responses
	}
	if cb.Publish != nil {
		d.Publish = *cb.Publish
	}
	if cb.Artifacts != nil {
		d.Artifacts = *cb.Artifacts
	}
	return d, nil
}

// InboxAnswer is what one inbox answer did.
type InboxAnswer struct {
	Item        string   `json:"item"`
	Verb        string   `json:"verb"`
	File        string   `json:"file"`
	Points      int      `json:"points"`
	State       string   `json:"state"` // mission state after the answer
	Validations []string `json:"validations"`
	Dropped     []string `json:"dropped"`
}

// InboxReply answers one open row. A questions row takes the verb answer and words; a gate
// row takes one verb of its gate shape.
func (env Env) InboxReply(ref string, r Reply) (InboxAnswer, error) {
	var res InboxAnswer
	it, err := env.findItem(ref, false)
	if err != nil {
		return res, err
	}
	if it.State != "open" {
		return res, fmt.Errorf("%s is %s", it.ID, it.State)
	}
	if !slices.Contains(it.Verbs, r.Verb) {
		return res, fmt.Errorf("%s (%s): verb %q, want one of %v", it.ID, it.Type, r.Verb, it.Verbs)
	}
	if !slices.Contains(surfaces, r.Surface) {
		return res, fmt.Errorf("surface %q: want one of %v", r.Surface, surfaces)
	}
	if slices.Contains(validating, r.Verb) {
		open, err := env.Inbox(false)
		if err != nil {
			return res, err
		}
		for _, o := range open {
			if o.Mission == it.Mission && o.Source == "questions" && o.Blocking {
				q := o.ID[strings.LastIndex(o.ID, "/")+1:]
				return res, fmt.Errorf("%s: blocking question %s is open; answer %s first", it.ID, q, o.ID)
			}
		}
	}
	m, err := env.find(it.Mission)
	if err != nil {
		return res, err
	}
	if it.Source == "questions" {
		return env.answerQuestion(m, it, r)
	}
	return env.answerGate(m, it, r)
}

// answerQuestion writes "qN: words" into answers.vN.md. The mission goes back to queued
// once no question of its version is open.
func (env Env) answerQuestion(m contract.Mission, it contract.InboxItem, r Reply) (InboxAnswer, error) {
	res := InboxAnswer{Item: it.ID, Verb: r.Verb, File: AnswersFile(it.Version), Validations: []string{}, Dropped: []string{}}
	if m.State != "waiting" && m.State != "running" {
		return res, fmt.Errorf("%s is %s; a question is answered while the mission waits or runs", m.ID, m.State)
	}
	words := strings.Join(strings.Fields(strings.Join(r.Words, " ")), " ")
	if words == "" {
		return res, fmt.Errorf("%s: an answer needs words", it.ID)
	}
	qid := it.ID[strings.LastIndex(it.ID, "/")+1:]
	dir := filepath.Join(env.Store.Root(), m.ID)
	points, err := readPoints(filepath.Join(dir, res.File))
	if err != nil {
		return res, err
	}
	points = append(points, qid+": "+words)
	ts := env.Now().Format(time.RFC3339)
	if err := writeAnswers(dir, m.ID, it.Version, r.Verb, r.Surface, ts, points); err != nil {
		return res, err
	}
	res.Points = len(points)
	item, surface := it.ID, r.Surface
	if err := env.append(m.ID, "owner", "answer", ts,
		contract.AnswerData{File: res.File, Verb: r.Verb, Surface: &surface, Item: &item}); err != nil {
		return res, err
	}
	res.State = m.State
	if m.State != "waiting" {
		return res, nil
	}
	open, err := env.Inbox(false)
	if err != nil {
		return res, err
	}
	if slices.ContainsFunc(open, func(o contract.InboxItem) bool { return o.Mission == m.ID }) {
		return res, nil
	}
	return env.moveAfter(m, res, "queued", "answered")
}

// answerGate records the verb on the version file through item_108, then moves the
// mission: no and abandon cancel it; ok, approve and a last go end it, or keep it in
// review for the publish step (item_099) when they validated a publish list; every other
// verb queues it for a resume (item_135).
func (env Env) answerGate(m contract.Mission, it contract.InboxItem, r Reply) (InboxAnswer, error) {
	res := InboxAnswer{Item: it.ID, Verb: r.Verb}
	fyi := it.Type == "fyi"
	if fyi && m.State != "done" || !fyi && m.State != "review" {
		return res, fmt.Errorf("%s is %s; %s waits on a mission in %s", m.ID, m.State, it.ID,
			map[bool]string{true: "done", false: "review"}[fyi])
	}
	if m.Version != it.Version {
		return res, fmt.Errorf("%s: the mission is at v%d", it.ID, m.Version)
	}
	item := it.ID
	a, err := env.reply(m, r, &item)
	res.File, res.Points, res.Validations, res.Dropped = a.File, a.Points, a.Validations, a.Dropped
	if err != nil {
		return res, err
	}
	to, reason := "queued", "answered: "+r.Verb
	switch r.Verb {
	case "no", "abandon":
		to = "cancelled"
	case "ok", "approve", "go":
		if r.Verb == "go" && env.laterGate(m, it.Source) {
			break
		}
		if len(res.Validations) > 0 {
			res.State = m.State
			return res, nil
		}
		to = "done"
	}
	return env.moveAfter(m, res, to, reason)
}

// laterGate reports whether the kind of m has a gate after source.
func (env Env) laterGate(m contract.Mission, source string) bool {
	k, err := env.Kinds(m.Kind)
	if err != nil || k == nil {
		return false
	}
	i := slices.Index(k.Gates, source)
	return i >= 0 && i < len(k.Gates)-1
}

func (env Env) moveAfter(m contract.Mission, res InboxAnswer, to, reason string) (InboxAnswer, error) {
	if !contract.CanMove(m.State, to) {
		return res, fmt.Errorf("%s: no move %s -> %s", m.ID, m.State, to)
	}
	ts := env.Now().Format(time.RFC3339)
	if err := env.append(m.ID, "owner", "state", ts, contract.StateData{From: m.State, To: to, Reason: &reason}); err != nil {
		return res, err
	}
	res.State = to
	return res, nil
}
