package brief

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"crow-core/contract"
	"crow-core/internal/eventlog"
	"crow-core/internal/project"
	"crow-core/internal/testfix"
)

const example = "20261109-build-item-095"

var exampleDir = filepath.Join(testfix.RunSpec, "examples", example)

// exampleInput folds the example mission from its event log and loads the build kind
// and the example prefs file.
func exampleInput(t *testing.T) Input {
	t.Helper()
	kinds := project.KindDirs(testfix.KindsDir, "")
	recs, err := eventlog.Read(eventlog.Path(testfix.MissionsRoot(t), example), 0)
	if err != nil {
		t.Fatal(err)
	}
	var v project.View
	var key string
	for _, r := range recs {
		if r.Event.Type == "created" {
			d, err := contract.DecodeEventData(r.Event)
			if err != nil {
				t.Fatal(err)
			}
			key = d.(*contract.CreatedData).Key
		}
		if err := project.Apply(&v, r.Event, kinds); err != nil {
			t.Fatal(err)
		}
	}
	k, err := kinds(v.Mission.Kind)
	if err != nil {
		t.Fatal(err)
	}
	ps, err := LoadPrefs(filepath.Join(testfix.RunSpec, "prefs.example.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return Input{Mission: v.Mission, Key: key, Prefs: ForKind(ps, k.Name), Kind: k}
}

// copyAnswers copies the example answers files up to version n into a fresh dir.
func copyAnswers(t *testing.T, n int) string {
	t.Helper()
	dir := t.TempDir()
	for v := 1; v <= n; v++ {
		name := "answers.v" + string(rune('0'+v)) + ".md"
		if err := os.WriteFile(filepath.Join(dir, name), testfix.Read(t, filepath.Join(exampleDir, name)), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// AC1: the first brief of the example mission is the example brief.
func TestFirstBriefIsTheExample(t *testing.T) {
	got := Render(exampleInput(t))
	want := string(testfix.Read(t, filepath.Join(exampleDir, "brief.md")))
	if got != want {
		t.Errorf("brief differs from the example.\n--- got\n%s\n--- want\n%s", got, want)
	}
	var sections []string
	for _, l := range strings.Split(got, "\n") {
		if h, ok := strings.CutPrefix(l, "## "); ok {
			sections = append(sections, h)
		}
	}
	order := "Goal,Mission,Owner preferences,Preflight,Answers so far,Constraints,Reply contract"
	if strings.Join(sections, ",") != order {
		t.Errorf("sections %v, want %s", sections, order)
	}
}

// AC2: after answers v1 and v2 the brief lists both in the owner's words and asks for v3.
func TestResumeBriefCarriesTheAnswers(t *testing.T) {
	in := exampleInput(t)
	dir := copyAnswers(t, 2)
	var err error
	if in.Answers, err = ReadAnswers(dir); err != nil {
		t.Fatal(err)
	}
	if err := Write(dir, in); err != nil {
		t.Fatal(err)
	}
	got := string(testfix.Read(t, filepath.Join(dir, "brief.md")))
	for _, want := range []string{
		"- v1: dig.\n  1. q1: refuse.\n  2. Also list the repos with no scope in `crow doctor`.\n",
		"- v2: go.\n",
		"- current version: 3 (write `report.v3.md`)",
		"core resumes you with `answers.v3.md`.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("brief misses %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "None yet.") {
		t.Error("resume brief still says None yet.")
	}
	if _, err := os.Stat(filepath.Join(dir, ".brief.md.tmp")); !os.IsNotExist(err) {
		t.Error("temp file left behind")
	}
}

// AC3: the brief names no account and no local path, even when the record has both.
func TestBriefNamesNoAccountOrPath(t *testing.T) {
	in := exampleInput(t)
	home, _ := os.UserHomeDir()
	wt := filepath.Join(home, "Projects", "secret-repo.wt", example)
	in.Mission.Worktree = &wt
	in.Mission.Runs = append(in.Mission.Runs, contract.Run{Run: "zz99", Account: "owner-work-account"})
	got := Render(in)
	for _, bad := range []string{"owner-work-account", "secret-repo", home} {
		if bad != "" && strings.Contains(got, bad) {
			t.Errorf("brief leaks %q", bad)
		}
	}
	if regexp.MustCompile(`(^|[\s(])/[a-z]`).MatchString(got) {
		t.Errorf("brief carries an absolute path:\n%s", got)
	}
}

// item_143 AC2: with the mission folder set, the brief names it and every path the run writes,
// and still no repo path.
func TestBriefNamesTheMissionFolder(t *testing.T) {
	in := exampleInput(t)
	in.Dir = "/state/core/missions/" + example
	home, _ := os.UserHomeDir()
	wt := filepath.Join(home, "Projects", "secret-repo.wt", example)
	in.Mission.Worktree = &wt
	got := Render(in)
	for _, want := range []string{
		"- mission folder: `" + in.Dir + "`",
		"`" + in.Dir + "/" + Phase(in.Mission.Gates, nil) + ".v1.md`",
		"`" + in.Dir + "/outbox/`",
		"`" + in.Dir + "/outbox/questions.json`",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("brief misses %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "secret-repo") {
		t.Error("brief leaks the repo path")
	}
}

func TestPhaseWalksTheGates(t *testing.T) {
	go1 := []Answer{{Version: 1, Verb: "go"}}
	dig := []Answer{{Version: 1, Verb: "dig"}}
	for _, c := range []struct {
		gates   []string
		answers []Answer
		want    string
	}{
		{nil, nil, "report"},
		{[]string{"plan", "report"}, nil, "plan"},
		{[]string{"plan", "report"}, dig, "plan"},
		{[]string{"plan", "report"}, go1, "report"},
		{[]string{"mockup"}, nil, "mockup"},
		{[]string{"mockup"}, go1, "report"},
		{[]string{"report"}, nil, "report"},
	} {
		if got := Phase(c.gates, c.answers); got != c.want {
			t.Errorf("Phase(%v, %v) = %s, want %s", c.gates, c.answers, got, c.want)
		}
	}
}

func TestPrefsFilterByKind(t *testing.T) {
	ps, err := LoadPrefs(filepath.Join(testfix.RunSpec, "prefs.example.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	w := ForKind(ps, "write")
	if len(w) != 2 || !strings.HasPrefix(w[1], "Drafts for others") {
		t.Errorf("write prefs %v", w)
	}
	if r := ForKind(ps, "research"); len(r) != 1 {
		t.Errorf("research prefs %v, want the one shared rule", r)
	}
	if none, err := LoadPrefs(filepath.Join(t.TempDir(), "missing.yaml")); err != nil || none != nil {
		t.Errorf("missing file: %v %v", none, err)
	}
	bad := filepath.Join(t.TempDir(), "p.yaml")
	os.WriteFile(bad, []byte("- rul: x\n"), 0o600)
	if _, err := LoadPrefs(bad); err == nil {
		t.Error("unknown key accepted")
	}
}

func TestReadOnlyKindHasNoWorktreeLine(t *testing.T) {
	k, err := project.KindDirs(testfix.KindsDir, "")("research")
	if err != nil {
		t.Fatal(err)
	}
	m := contract.Mission{ID: "20261109-research-android-push-options", Kind: "research", Goal: "g", Gates: []string{}}
	got := Render(Input{Mission: m, Kind: k})
	for _, want := range []string{"- gates: none (straight through)", "- key: android-push-options", "write `report.v1.md`", "anywhere but `outbox/`", "Owner preferences\nNone."} {
		if !strings.Contains(got, want) {
			t.Errorf("research brief misses %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "- repo:") || strings.Contains(got, "worktree") {
		t.Errorf("research brief names a worktree:\n%s", got)
	}
}

// The brief shows the questions.json shape with the mission and version filled in, and the
// shape decodes strictly as a questions file (runs used to copy crow-block fields).
func TestBriefShowsTheQuestionsShape(t *testing.T) {
	in := exampleInput(t)
	got := Render(in)
	_, rest, ok := strings.Cut(got, "`questions.json` shape")
	if !ok {
		t.Fatalf("brief has no questions.json shape:\n%s", got)
	}
	_, rest, ok = strings.Cut(rest, "```json\n")
	body, _, ok2 := strings.Cut(rest, "```")
	if !ok || !ok2 {
		t.Fatalf("questions.json shape is not a json fence:\n%s", got)
	}
	var q contract.Questions
	if err := contract.Decode([]byte(body), &q); err != nil {
		t.Fatalf("shape does not decode as questions.json: %v\n%s", err, body)
	}
	if q.Contract != contract.Version || q.Mission != in.Mission.ID || q.Version != Version(in.Answers) {
		t.Errorf("shape names contract %d, mission %s, version %d", q.Contract, q.Mission, q.Version)
	}
	if len(q.Questions) != 1 || !q.Questions[0].Blocking || q.Questions[0].Choices == nil {
		t.Errorf("shape question %+v", q.Questions)
	}
}
