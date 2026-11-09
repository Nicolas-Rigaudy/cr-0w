package project

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"crow-core/contract"
	"crow-core/internal/eventlog"
	"crow-core/internal/testfix"
)

func fold(t *testing.T, mission string, kinds Kinds) View {
	t.Helper()
	recs, err := eventlog.Read(eventlog.Path(testfix.MissionsRoot(t), mission), 0)
	if err != nil {
		t.Fatal(err)
	}
	var v View
	for _, r := range recs {
		if err := Apply(&v, r.Event, kinds); err != nil {
			t.Fatal(err)
		}
	}
	return v
}

func records[T any](t *testing.T, name string) []T {
	t.Helper()
	var out []T
	if err := json.Unmarshal(testfix.Read(t, filepath.Join(testfix.Records, name)), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func asJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// eventFields is what an inbox row takes from the events; item_134 fills the rest.
type eventFields struct {
	ID, Mission, Kind, Project, Type, Source, File, State, Created string
	Version                                                        int
	Verbs                                                          []string
	Answered                                                       *string
}

func fields(it contract.InboxItem) eventFields {
	return eventFields{it.ID, it.Mission, it.Kind, it.Project, it.Type, it.Source, it.File, it.State, it.Created, it.Version, it.Verbs, it.Answered}
}

func TestResearchFixtureGivesTheRecords(t *testing.T) {
	v := fold(t, "20261109-research-android-push-options", nil)
	want := records[contract.Mission](t, "missions.json")[0]
	if got, w := asJSON(t, v.Mission), asJSON(t, want); got != w {
		t.Errorf("mission\n got %s\nwant %s", got, w)
	}
	if got, w := asJSON(t, v.Mission.Runs[0]), asJSON(t, records[contract.Run](t, "runs.json")[0]); got != w {
		t.Errorf("run\n got %s\nwant %s", got, w)
	}
	inbox := records[contract.InboxItem](t, "inbox.json")[:2]
	if len(v.Inbox) != len(inbox) {
		t.Fatalf("%d inbox rows, want %d", len(v.Inbox), len(inbox))
	}
	for i := range inbox {
		if got, w := fields(v.Inbox[i]), fields(inbox[i]); !reflect.DeepEqual(got, w) {
			t.Errorf("inbox row %d\n got %+v\nwant %+v", i, got, w)
		}
	}
}

func TestBuildFixtureTakesKindDefaults(t *testing.T) {
	v := fold(t, "20261109-build-item-095", KindDirs(testfix.KindsDir, ""))
	m := v.Mission
	if !reflect.DeepEqual(m.Gates, []string{"plan", "report"}) {
		t.Errorf("gates %v, want the build kind's [plan report]", m.Gates)
	}
	if m.State != "done" || m.Version != 3 || len(m.Runs) != 3 {
		t.Errorf("state %s version %d runs %d, want done 3 3", m.State, m.Version, len(m.Runs))
	}
	if m.Validation == nil || *m.Validation != "2026-11-09T11:02:02+01:00" {
		t.Errorf("validation %v", m.Validation)
	}
	rows := map[string]contract.InboxItem{}
	for _, it := range v.Inbox {
		rows[it.ID] = it
	}
	id := "20261109-build-item-095"
	if it, ok := rows[id+"/v2"]; !ok || it.Type != "approval" || it.State != "answered" {
		t.Errorf("v2: %+v", it)
	}
	if it, ok := rows[id+"/v3"]; !ok || it.Type != "review" || it.State != "answered" {
		t.Errorf("v3: %+v", it)
	}
	if _, ok := rows[id+"/v1"]; ok {
		t.Errorf("v1 needs input and must not reach the inbox")
	}
}

func TestApplyRefusesBrokenSequences(t *testing.T) {
	ev := func(typ, data string) contract.Event {
		return contract.Event{TS: "2026-11-09T14:00:00+01:00", Mission: "20261109-research-x", Type: typ, Actor: "core", Data: json.RawMessage(data)}
	}
	created := ev("created", `{"kind":"research","batch":null}`)
	for name, seq := range map[string][]contract.Event{
		"no created first": {ev("state", `{"from":"preflight","to":"queued"}`)},
		"created twice":    {created, created},
		"wrong from":       {created, ev("state", `{"from":"running","to":"done"}`)},
		"unknown run end":  {created, ev("run_end", `{"run":"aaaa1111","cause":"done","exit":0}`)},
	} {
		var v View
		var err error
		for _, e := range seq {
			if err = Apply(&v, e, nil); err != nil {
				break
			}
		}
		if err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}
