package mission

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"crow-core/contract"
	"crow-core/internal/testfix"
)

// parked starts a research mission and parks it on the example question of version 1.
func parked(t *testing.T, blocking bool) (Env, string, string) {
	t.Helper()
	env, root := newEnv(t)
	id, r := startOne(t, env)
	questions(t, root, id, 1)
	if blocking {
		p := filepath.Join(root, id, QuestionsFile)
		b, _ := os.ReadFile(p)
		if err := os.WriteFile(p, []byte(strings.Replace(string(b), `"blocking": false`, `"blocking": true`, 1)), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	r.ends["r1"] = ended("done", 0)
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	return env, root, id
}

// readyPlan replays the example build mission up to its review on plan.v1.md, with the
// gate verdict set to ready so the plan opens an inbox row. plan.v1.md has blocking q1.
func readyPlan(t *testing.T) (Env, string) {
	t.Helper()
	env, root := newEnv(t)
	lines := exampleLines(t)
	lines[5] = bytes.Replace(lines[5], []byte(`"needs-input"`), []byte(`"ready"`), 1)
	replayCore(t, env, lines, 0, 7)
	copyExample(t, filepath.Join(root, example), "plan.v1.md")
	return env, root
}

func lastEvent(t *testing.T, root, id string) contract.Event {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, id, "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	ls := strings.Split(strings.TrimSpace(string(b)), "\n")
	var e contract.Event
	if err := contract.Decode([]byte(ls[len(ls)-1]), &e); err != nil {
		t.Fatal(err)
	}
	return e
}

// AC1: a parked mission shows one item; the answer closes it, writes answers.v1.md and an
// answer event, and queues the mission. Tick holds it until the resume of item_135.
func TestInboxAnswersAParkedQuestion(t *testing.T) {
	env, root, id := parked(t, false)
	items, err := env.Inbox(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("inbox = %+v", items)
	}
	it := items[0]
	if it.ID != id+"/v1/q1" || it.Type != "choice" || !strings.HasPrefix(it.Ask, "The second provider") ||
		it.Blocking || len(it.Artifacts) != 0 {
		t.Errorf("item = %+v", it)
	}
	res, err := env.InboxReply(id+"/v1", Reply{Verb: "answer", Words: []string{"skip it"}, Surface: "cli"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Item != it.ID || res.File != "answers.v1.md" || res.State != "queued" {
		t.Errorf("result = %+v", res)
	}
	b, err := os.ReadFile(filepath.Join(root, id, "answers.v1.md"))
	if err != nil || !strings.Contains(string(b), "verb: answer\nsurface: cli\n") || !strings.HasSuffix(string(b), "1. q1: skip it\n") {
		t.Errorf("answers.v1.md = %q, %v", b, err)
	}
	want := "created,state,state,run_start,run_end,question,state,answer,state"
	if e := strings.Join(events(t, root, id), ","); e != want {
		t.Errorf("events = %s, want %s", e, want)
	}
	if left, _ := env.Inbox(false); len(left) != 0 {
		t.Errorf("open items left: %+v", left)
	}
	all, _ := env.Inbox(true)
	if len(all) != 1 || all[0].State != "answered" || all[0].Answered == nil {
		t.Errorf("answered item = %+v", all)
	}
	rep, err := env.Tick(&fakeRunner{ends: map[string]RunState{}}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Started) != 0 || len(rep.Held) != 1 || !strings.Contains(rep.Held[0].Detail, "item_135") {
		t.Errorf("tick = %+v", rep)
	}
}

// AC2: a verb outside the gate shape is refused with the allowed list, and an open blocking
// question refuses go.
func TestInboxRefusesWrongVerbAndOpenBlockingQuestion(t *testing.T) {
	env, root := readyPlan(t)
	items, err := env.Inbox(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Type != "approval" || !items[0].Blocking {
		t.Fatalf("inbox = %+v", items)
	}
	ref := items[0].ID
	if _, err := env.InboxReply(ref, Reply{Verb: "ok", Surface: "cli"}); err == nil ||
		!strings.Contains(err.Error(), "want one of [go dig no]") {
		t.Errorf("ok on a plan: err = %v", err)
	}
	if _, err := env.InboxReply(ref, Reply{Verb: "go", Surface: "cli"}); err == nil ||
		!strings.Contains(err.Error(), "blocking question q1 is open") {
		t.Errorf("go with q1 open: err = %v", err)
	}
	// The answer to q1 lets go through; build has a report gate after the plan, so the
	// mission queues for its next run.
	if _, err := env.InboxReply(ref, Reply{Verb: "dig", Words: []string{"q1: refuse."}, Surface: "cli"}); err != nil {
		t.Fatal(err)
	}
	if got := show(t, env, example); got.State != "queued" {
		t.Errorf("state after dig = %s", got.State)
	}
	e := lastEvent(t, root, example)
	if e.Type != "state" || !strings.Contains(string(e.Data), `"to":"queued"`) {
		t.Errorf("last event = %s %s", e.Type, e.Data)
	}
}

func TestInboxQuestionBlocksGoOnTheSameMission(t *testing.T) {
	env, _, _ := parked(t, true)
	items, _ := env.Inbox(false)
	if len(items) != 1 || !items[0].Blocking {
		t.Fatalf("inbox = %+v", items)
	}
	if _, err := env.InboxReply(items[0].ID, Reply{Verb: "go", Surface: "cli"}); err == nil ||
		!strings.Contains(err.Error(), "want one of [answer]") {
		t.Errorf("go on a question: err = %v", err)
	}
	if _, err := env.InboxReply(items[0].ID, Reply{Verb: "answer", Surface: "cli"}); err == nil ||
		!strings.Contains(err.Error(), "needs words") {
		t.Errorf("empty answer: err = %v", err)
	}
}

// Gate verbs end or reopen the mission: no cancels, ok without a publish list ends it,
// dig on an FYI report reopens a done mission.
func TestInboxGateVerbsMoveTheMission(t *testing.T) {
	for _, c := range []struct{ verb, state string }{{"no", "cancelled"}, {"dig", "queued"}} {
		env, _ := readyPlan(t)
		words := []string{"q1: refuse."}
		if _, err := env.InboxReply(example+"/v1", Reply{Verb: c.verb, Words: words, Surface: "cli"}); err != nil {
			t.Fatal(err)
		}
		if got := show(t, env, example); got.State != c.state {
			t.Errorf("%s: state = %s, want %s", c.verb, got.State, c.state)
		}
	}

	env, root := newEnv(t)
	id, r := startOne(t, env)
	writeVersion(t, root, id, "report.v1.md", "ready", 1)
	r.ends["r1"] = ended("done", 0)
	if _, err := env.Tick(r, 2); err != nil {
		t.Fatal(err)
	}
	items, _ := env.Inbox(false)
	if len(items) != 1 || items[0].Type != "fyi" || items[0].Ask != "one provider" {
		t.Fatalf("inbox = %+v", items)
	}
	res, err := env.InboxReply(items[0].ID, Reply{Verb: "dig", Words: []string{"compare prices too"}, Surface: "cli"})
	if err != nil || res.State != "queued" {
		t.Fatalf("dig on fyi = %+v, %v", res, err)
	}
	e := lastEvent(t, root, id)
	var prev contract.Event
	b, _ := os.ReadFile(filepath.Join(root, id, "events.jsonl"))
	ls := strings.Split(strings.TrimSpace(string(b)), "\n")
	if err := contract.Decode([]byte(ls[len(ls)-2]), &prev); err != nil {
		t.Fatal(err)
	}
	if e.Type != "state" || prev.Type != "answer" ||
		!strings.Contains(string(prev.Data), `"item":"`+id+`/v1","surface":"cli"`) {
		t.Errorf("events = %s %s / %s", prev.Type, prev.Data, e.Type)
	}
}

// AC3: inbox rows read from the research fixture log have the fixture shape: the same
// keys, and the same id, type, source, verbs and state per row. The fixture rows of the
// build mission have no log; the project test compares the same two rows.
func TestInboxMatchesTheFixtureShape(t *testing.T) {
	env, root := newEnv(t)
	const research = "20261109-research-android-push-options"
	for id, src := range testfix.Logs() {
		if id != research {
			continue
		}
		dir := filepath.Join(root, id)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), testfix.Read(t, src), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	items, err := env.Inbox(true)
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}
	var have, want []map[string]any
	if err := json.Unmarshal(got, &have); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(testfix.Read(t, filepath.Join(testfix.Records, "inbox.json")), &want); err != nil {
		t.Fatal(err)
	}
	index := func(rows []map[string]any) map[string]map[string]any {
		out := map[string]map[string]any{}
		for _, r := range rows {
			out[r["id"].(string)] = r
		}
		return out
	}
	var mine []map[string]any
	for _, r := range want {
		if r["mission"] == research {
			mine = append(mine, r)
		}
	}
	h, w := index(have), index(mine)
	if len(h) != len(w) {
		t.Fatalf("%d rows, want %d", len(h), len(w))
	}
	keys := func(r map[string]any) []string {
		var ks []string
		for k := range r {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return ks
	}
	for id, wr := range w {
		hr, ok := h[id]
		if !ok {
			t.Errorf("no row %s", id)
			continue
		}
		if !reflect.DeepEqual(keys(hr), keys(wr)) {
			t.Errorf("%s keys = %v, want %v", id, keys(hr), keys(wr))
		}
		for _, k := range []string{"type", "source", "verbs", "state", "version", "file", "created", "answered"} {
			if !reflect.DeepEqual(hr[k], wr[k]) {
				t.Errorf("%s %s = %v, want %v", id, k, hr[k], wr[k])
			}
		}
		var it contract.InboxItem
		raw, _ := json.Marshal(hr)
		if err := contract.Decode(raw, &it); err != nil {
			t.Errorf("%s: %v", id, err)
		}
	}
}
