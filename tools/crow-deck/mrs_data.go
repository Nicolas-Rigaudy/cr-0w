package main

// MRs data: crow-mrs json (the list) and crow-mrs show (diff, threads, saved findings of one
// MR) mapped to rows. Read only: nothing here posts, resolves or pushes.

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

type mrRow struct {
	Project         string `json:"project"`
	Name            string `json:"name"`
	Scope           string `json:"scope"`
	Forge           string `json:"forge"`
	Repo            string `json:"repo"`
	Number          int    `json:"number"`
	Ref             string `json:"ref"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	Role            string `json:"role"`
	Draft           bool   `json:"draft"`
	Pipeline        string `json:"pipeline"`
	ThreadsOpen     int    `json:"threads_open"`
	ThreadsResolved int    `json:"threads_resolved"`
	Updated         string `json:"updated"`
	Branch          string `json:"branch"`
	Target          string `json:"target"`
	URL             string `json:"url"`
}

func (r mrRow) key() string { return r.Repo + r.Ref }

type mrBranch struct {
	Project string `json:"project"`
	Scope   string `json:"scope"`
	Branch  string `json:"branch"`
	Ahead   int    `json:"ahead"`
}

type mrsSnap struct {
	Forges   map[string]string `json:"forges"`
	MRs      []mrRow           `json:"mrs"`
	Branches []mrBranch        `json:"branches"`
}

// diffLine is one line of a unified diff; Kind is '+', '-' or ' '.
type diffLine struct {
	Old, New int
	Kind     byte
	Text     string
}

type mrFile struct {
	Path     string
	Add, Del int
	Lines    []diffLine
}

// mrThread is a review thread on the owner's MR. Pattern, Fix, Test and Reply stay empty
// until the review-threads skill saves them.
type mrThread struct {
	Author, At, Path string
	Line             int
	Body             string
	Resolved         bool
	Pattern          string
	Seen             int
	Fix, Test        *mrSection
	Reply            string
}

// mrSection is a titled diff: the fix or the test of a thread.
type mrSection struct {
	Title string
	Meta  string
	Lines []diffLine
}

// mrFinding is one finding of a review run (runbooks/review-findings.md).
type mrFinding struct {
	ID, Title, Path string
	Line            int
	Severity        string // blocking, suggestion, nit, question
	Body            string
	Pattern         string
	Seen            int
	Comment         string
	Fix             string
	State           string    // pending, accepted, refused, edited, chat
	Note            string    // the owner's edited comment or chat note
	Applied         string    // apply run result: fixed, answered, skipped, failed ("" before any)
	AppliedDetail   string    // what changed, or the answer to the chat note
	AppliedFor      [2]string // the decision and note the apply run acted on
	Posted          bool      // the comment went to the forge (peer review) or the reply to the thread
	// Threads review only: the reviewer's thread, the run's assessment, the test, and whether
	// the thread is resolved on the forge.
	ThreadAuthor, ThreadBody string
	ThreadPath               string // where the reviewer commented
	ThreadLine               int
	Agree                    string // yes, no, partly
	Test                     string
	Resolved                 bool
}

// at is where the card points in the diff: the reviewer's line for a thread, else the finding's.
func (f mrFinding) at() (string, int) {
	if f.ThreadPath != "" {
		return f.ThreadPath, f.ThreadLine
	}
	return f.Path, f.Line
}

// toApply is true when the next apply run acts on this finding: a chat note in any mode,
// an accepted or edited finding on a self review, or on a threads review when it has a fix,
// not already handled for the same decision.
func (f mrFinding) toApply(mode string) bool {
	acc := f.State == "accepted" || f.State == "edited"
	ok := f.State == "chat" || (mode == "self" && acc) || (mode == "threads" && acc && f.Fix != "")
	return ok && f.AppliedFor != [2]string{f.State, f.Note}
}

// mrApply is the apply run that resumes the review session (commit local, never pushed).
type mrApply struct {
	ID, State, Reason, Commit, Subject string
	StartedAt                          string
	Warnings                           []string
	Pushed                             bool // the commit went to the MR branch
}

// mrReview is the state of the review run behind the findings.
type mrReview struct {
	ID, State, Reason, Model, StartedAt string
	Mode                                string // peer, self, threads
}

// openThreads counts the unresolved threads of a loaded detail; 0 for nil.
func (d *mrDetail) openThreads() int {
	n := 0
	if d != nil {
		for _, t := range d.Threads {
			if !t.Resolved {
				n++
			}
		}
	}
	return n
}

type mrDetail struct {
	Files    []mrFile
	Threads  []mrThread
	Findings []mrFinding
	Review   *mrReview
	Apply    *mrApply
	// FindingsFile is the path of the review findings file, empty when the MR has none.
	FindingsFile string
}

// ----- fetching -----

func fetchMRs() (*mrsSnap, error) {
	var s mrsSnap
	return &s, runJSON(90*time.Second, nil, &s, "crow-mrs", "json")
}

type rawDetail struct {
	Files []struct {
		Path string `json:"path"`
		Add  int    `json:"add"`
		Del  int    `json:"del"`
		Diff string `json:"diff"`
	} `json:"files"`
	Threads []struct {
		Author   string `json:"author"`
		At       string `json:"at"`
		Path     string `json:"path"`
		Line     int    `json:"line"`
		Body     string `json:"body"`
		Resolved bool   `json:"resolved"`
	} `json:"threads"`
	Findings []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Severity string `json:"severity"`
		Path     string `json:"path"`
		Line     int    `json:"line"`
		Body     string `json:"body"`
		Pattern  string `json:"pattern"`
		Seen     int    `json:"seen"`
		Comment  string `json:"comment"`
		Fix      string `json:"fix"`
		State    string `json:"state"`
		Note     string `json:"note"`
		Apply    *struct {
			Status string   `json:"status"`
			Detail string   `json:"detail"`
			For    []string `json:"for"`
		} `json:"apply"`
		Posted *struct {
			At string `json:"at"`
		} `json:"posted"`
		Thread *struct {
			Author string `json:"author"`
			Body   string `json:"body"`
			Path   string `json:"path"`
			Line   int    `json:"line"`
		} `json:"thread"`
		Agree    string `json:"agree"`
		Test     string `json:"test"`
		Resolved *struct {
			At string `json:"at"`
		} `json:"resolved"`
	} `json:"findings"`
	FindingsFile string `json:"findings_file"`
	Review       *struct {
		ID        string `json:"id"`
		State     string `json:"state"`
		Reason    string `json:"reason"`
		Model     string `json:"model"`
		StartedAt string `json:"started_at"`
		Mode      string `json:"mode"`
	} `json:"review"`
	Apply *struct {
		ID       string   `json:"id"`
		State    string   `json:"state"`
		Reason   string   `json:"reason"`
		Commit   string   `json:"commit"`
		Subject  string   `json:"subject"`
		Warnings []string `json:"warnings"`
		Started  string   `json:"started_at"`
		Pushed   *struct {
			Commit string `json:"commit"`
		} `json:"pushed"`
	} `json:"apply"`
}

func fetchMRDetail(r mrRow) (*mrDetail, error) {
	var raw rawDetail
	if err := runJSON(60*time.Second, nil, &raw, "crow-mrs", "show", r.Repo, strconv.Itoa(r.Number)); err != nil {
		return nil, err
	}
	return mapDetail(raw), nil
}

func mapDetail(raw rawDetail) *mrDetail {
	d := &mrDetail{FindingsFile: raw.FindingsFile}
	for _, f := range raw.Files {
		d.Files = append(d.Files, mrFile{Path: f.Path, Add: f.Add, Del: f.Del, Lines: parseDiff(f.Diff)})
	}
	for _, t := range raw.Threads {
		d.Threads = append(d.Threads, mrThread{Author: t.Author, At: t.At, Path: t.Path, Line: t.Line, Body: strings.TrimSpace(t.Body), Resolved: t.Resolved})
	}
	for _, f := range raw.Findings {
		mf := mrFinding{ID: f.ID, Title: f.Title, Severity: f.Severity, Path: f.Path, Line: f.Line, Body: f.Body,
			Pattern: f.Pattern, Seen: f.Seen, Comment: f.Comment, Fix: f.Fix, State: cmpOr(f.State, "pending"), Note: f.Note}
		if f.Apply != nil {
			mf.Applied, mf.AppliedDetail = f.Apply.Status, f.Apply.Detail
			if len(f.Apply.For) == 2 {
				mf.AppliedFor = [2]string{f.Apply.For[0], f.Apply.For[1]}
			}
		}
		mf.Posted, mf.Resolved = f.Posted != nil, f.Resolved != nil
		mf.Agree, mf.Test = f.Agree, f.Test
		if f.Thread != nil {
			mf.ThreadAuthor, mf.ThreadBody = f.Thread.Author, strings.TrimSpace(f.Thread.Body)
			mf.ThreadPath, mf.ThreadLine = f.Thread.Path, f.Thread.Line
		}
		d.Findings = append(d.Findings, mf)
	}
	if ap := raw.Apply; ap != nil {
		d.Apply = &mrApply{ID: ap.ID, State: ap.State, Reason: ap.Reason, Commit: ap.Commit, Subject: ap.Subject, Warnings: ap.Warnings, StartedAt: ap.Started,
			Pushed: ap.Pushed != nil && ap.Pushed.Commit == ap.Commit && ap.Commit != ""}
	}
	if rv := raw.Review; rv != nil {
		d.Review = &mrReview{ID: rv.ID, State: rv.State, Reason: rv.Reason, Model: rv.Model, StartedAt: rv.StartedAt, Mode: rv.Mode}
	}
	return d
}

// parseDiff reads the hunks of a unified diff into numbered lines. Hunk headers become
// a ' ' line with Old and New 0 and the header text, so gaps between hunks stay visible.
func parseDiff(diff string) []diffLine {
	var out []diffLine
	oldN, newN := 0, 0
	for _, l := range strings.Split(strings.TrimRight(diff, "\n"), "\n") {
		switch {
		case strings.HasPrefix(l, "@@"):
			oldN, newN = hunkStart(l)
			out = append(out, diffLine{Kind: '@', Text: l})
		case strings.HasPrefix(l, "+++"), strings.HasPrefix(l, "---"), strings.HasPrefix(l, "\\"):
		case strings.HasPrefix(l, "+"):
			out = append(out, diffLine{New: newN, Kind: '+', Text: l[1:]})
			newN++
		case strings.HasPrefix(l, "-"):
			out = append(out, diffLine{Old: oldN, Kind: '-', Text: l[1:]})
			oldN++
		default:
			out = append(out, diffLine{Old: oldN, New: newN, Kind: ' ', Text: strings.TrimPrefix(l, " ")})
			oldN++
			newN++
		}
	}
	return out
}

// hunkStart returns the first old and new line of "@@ -a,b +c,d @@".
func hunkStart(h string) (int, int) {
	f := strings.Fields(h)
	num := func(s string) int {
		s = strings.TrimLeft(s, "-+")
		s, _, _ = strings.Cut(s, ",")
		n, _ := strconv.Atoi(s)
		return n
	}
	if len(f) < 3 {
		return 1, 1
	}
	return num(f[1]), num(f[2])
}

// ----- mapping -----

// mrOrder: by project, the owner's own MRs before the ones to review, newest first.
func mrOrder(rows []mrRow) []mrRow {
	out := append([]mrRow{}, rows...)
	rank := func(r mrRow) int {
		if r.Role == "AUTHOR" {
			return 0
		}
		return 1
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Project != b.Project {
			return a.Project < b.Project
		}
		if rank(a) != rank(b) {
			return rank(a) < rank(b)
		}
		return a.Updated > b.Updated
	})
	return out
}

// mrCounts feeds the tab badge and the Blockers title: MRs awaiting the owner's review,
// the owner's own MRs, and the threads waiting for a fix.
func mrCounts(rows []mrRow) (awaiting, mine, threads int) {
	for _, r := range rows {
		switch r.Role {
		case "REVIEWER":
			awaiting++
		case "AUTHOR":
			mine++
			threads += r.ThreadsOpen
		}
	}
	return
}

// mrBlockers turns MRs into Home Blockers rows: a review to do, threads to answer, a red pipeline.
func mrBlockers(now time.Time, rows []mrRow) []blocker {
	var out []blocker
	for _, r := range rows {
		age := ""
		if t, ok := parseTime(r.Updated); ok {
			age = shortDays(now.Sub(t))
		}
		item := r.Ref + " " + r.Project
		switch {
		case r.Role == "REVIEWER":
			out = append(out, blocker{"needs", item, cmpOr(r.Author, "someone") + " waits for your review", age})
		case r.ThreadsOpen > 0:
			out = append(out, blocker{"needs", item, strconv.Itoa(r.ThreadsOpen) + " " + plural(r.ThreadsOpen, "thread", "threads") + " to answer", age})
		case r.Pipeline == "failed":
			out = append(out, blocker{"blocked", item, "pipeline failed", age})
		}
	}
	return out
}

// shortDays: minutes under an hour, hours under a day, then days.
func shortDays(d time.Duration) string {
	if d >= 24*time.Hour {
		return strconv.Itoa(int(d/(24*time.Hour))) + "d"
	}
	return shortAge(d)
}

// mrDecide is the MR part of Home "You decide": the same rows as mrBlockers, with the
// project's full name and scope.
func mrDecide(now time.Time, rows []mrRow, reg *registry) []nextRow {
	var out []nextRow
	for _, r := range rows {
		why := ""
		switch {
		case r.Role == "REVIEWER":
			why = cmpOr(r.Author, "someone") + " waits for your review"
		case r.ThreadsOpen > 0:
			why = strconv.Itoa(r.ThreadsOpen) + " " + plural(r.ThreadsOpen, "thread", "threads") + " to answer"
		case r.Pipeline == "failed":
			why = "pipeline failed"
		default:
			continue
		}
		age := ""
		if t, ok := parseTime(r.Updated); ok {
			age = shortDays(now.Sub(t))
		}
		out = append(out, nextRow{Project: reg.shortName(r.Project), Repo: reg.fullName(r.Project), What: r.Ref + " " + why,
			Perso: r.Scope == "personal", Actor: "me", Why: "MR", Age: age})
	}
	return out
}
