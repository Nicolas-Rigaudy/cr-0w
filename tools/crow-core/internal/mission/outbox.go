package mission

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"crow-core/contract"
)

// Outbox caps when a kind file sets none (req_037: most kinds take 20 MB and 100 MB).
const (
	defaultFileMB  = 20.0
	defaultTotalMB = 100.0
)

// QuestionsFile is where a run asks questions outside a plan or report.
const QuestionsFile = "outbox/questions.json"

// Drop is one listed outbox file core does not keep. Cause is the incident cause.
type Drop struct {
	Path, Cause, Detail string
}

// Caps are the outbox limits of one mission, in bytes.
type Caps struct {
	File, Total int64
}

// capsOf reads the outbox caps of a kind. A missing kind or cap takes the default.
func capsOf(k *contract.Kind) Caps {
	file, total := defaultFileMB, defaultTotalMB
	if k != nil && k.Caps.OutboxFileMB != nil {
		file = *k.Caps.OutboxFileMB
	}
	if k != nil && k.Caps.OutboxTotalMB != nil {
		total = *k.Caps.OutboxTotalMB
	}
	return Caps{File: int64(file * (1 << 20)), Total: int64(total * (1 << 20))}
}

// outboxFile resolves a listed path to a regular file inside <dir>/outbox. It refuses
// absolute paths, paths outside outbox/, any ".." element, and links that lead out.
func outboxFile(dir, p string) (string, error) {
	if filepath.IsAbs(p) || !strings.HasPrefix(p, "outbox/") {
		return "", errors.New("not under outbox/")
	}
	for _, el := range strings.Split(p, "/") {
		if el == ".." {
			return "", errors.New(`contains ".."`)
		}
	}
	box, err := filepath.EvalSymlinks(filepath.Join(dir, "outbox"))
	if err != nil {
		return "", err
	}
	full, err := filepath.EvalSymlinks(filepath.Join(dir, filepath.FromSlash(p)))
	if err != nil {
		return "", err
	}
	if rel, err := filepath.Rel(box, full); err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return "", errors.New("leads outside outbox/")
	}
	return full, nil
}

// Ingest checks the listed outbox files of a turn against the caps, in list order. It
// keeps each file once, and drops what is missing, refused, over the per-file cap, or
// over the total once earlier files filled it. Files not listed are ignored. Nothing on
// disk changes.
func Ingest(dir string, paths []string, c Caps) (kept []string, drops []Drop) {
	kept = []string{}
	seen := map[string]bool{}
	var total int64
	for _, p := range paths {
		if seen[p] {
			continue
		}
		seen[p] = true
		full, err := outboxFile(dir, p)
		switch {
		case errors.Is(err, os.ErrNotExist):
			drops = append(drops, Drop{p, "artifact_missing", p + ": no such file"})
			continue
		case err != nil:
			drops = append(drops, Drop{p, "artifact_refused", p + ": " + err.Error()})
			continue
		}
		fi, err := os.Stat(full)
		if err != nil {
			drops = append(drops, Drop{p, "artifact_missing", p + ": " + err.Error()})
			continue
		}
		if !fi.Mode().IsRegular() {
			drops = append(drops, Drop{p, "artifact_refused", p + ": not a regular file"})
			continue
		}
		size := fi.Size()
		if size > c.File {
			drops = append(drops, Drop{p, "artifact_over_cap",
				fmt.Sprintf("%s: %s over the %s per-file cap", p, mb(size), mb(c.File))})
			continue
		}
		if total+size > c.Total {
			drops = append(drops, Drop{p, "artifact_over_cap",
				fmt.Sprintf("%s: %s would pass the %s outbox cap (%s kept)", p, mb(size), mb(c.Total), mb(total))})
			continue
		}
		total += size
		kept = append(kept, p)
	}
	return kept, drops
}

func mb(n int64) string { return fmt.Sprintf("%.1f MB", float64(n)/(1<<20)) }

var (
	questionID = regexp.MustCompile(`^q[0-9]+$`)
	choiceID   = regexp.MustCompile(`^[a-z0-9-]+$`)
)

// readQuestions reads outbox/questions.json for the run of version n. ok is false when
// there is no file or the file belongs to an earlier version (a question already
// parked). A file that is present but invalid gives why.
func readQuestions(dir, id string, n int) (q *contract.Questions, ok bool, why string) {
	full, err := outboxFile(dir, QuestionsFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, ""
	}
	if err != nil {
		return nil, false, QuestionsFile + ": " + err.Error()
	}
	b, err := os.ReadFile(full)
	if err != nil {
		return nil, false, QuestionsFile + ": " + err.Error()
	}
	q = &contract.Questions{}
	if err := contract.Decode(b, q); err != nil {
		return nil, false, QuestionsFile + ": " + err.Error()
	}
	if q.Mission == id && q.Version < n {
		return nil, false, ""
	}
	if err := checkQuestions(q, id, n); err != nil {
		return nil, false, QuestionsFile + ": " + err.Error()
	}
	return q, true, ""
}

// checkQuestions applies questions.schema.json and the mission and version of the run.
func checkQuestions(q *contract.Questions, id string, n int) error {
	switch {
	case q.Contract != contract.Version:
		return fmt.Errorf("contract %d, want %d", q.Contract, contract.Version)
	case q.Mission != id:
		return fmt.Errorf("names mission %s", q.Mission)
	case q.Version != n:
		return fmt.Errorf("version %d, the run writes version %d", q.Version, n)
	case len(q.Questions) == 0:
		return errors.New("no questions")
	}
	ids := map[string]bool{}
	for _, x := range q.Questions {
		if !questionID.MatchString(x.ID) || ids[x.ID] {
			return fmt.Errorf("question id %q is not q<N> or repeats", x.ID)
		}
		ids[x.ID] = true
		if strings.TrimSpace(x.Text) == "" {
			return fmt.Errorf("%s: empty text", x.ID)
		}
		if x.Choices == nil {
			continue
		}
		for _, c := range *x.Choices {
			if !choiceID.MatchString(c.ID) {
				return fmt.Errorf("%s: choice id %q", x.ID, c.ID)
			}
			if c.Preview != nil && !strings.HasPrefix(*c.Preview, "outbox/") {
				return fmt.Errorf("%s/%s: preview %q not under outbox/", x.ID, c.ID, *c.Preview)
			}
		}
	}
	return nil
}

// previews lists the preview files of the questions, in order.
func previews(q *contract.Questions) []string {
	var ps []string
	for _, x := range q.Questions {
		if x.Choices == nil {
			continue
		}
		for _, c := range *x.Choices {
			if c.Preview != nil {
				ps = append(ps, *c.Preview)
			}
		}
	}
	return ps
}

func artifactPaths(cb *contract.CrowBlock) []string {
	var ps []string
	if cb.Artifacts != nil {
		for _, a := range *cb.Artifacts {
			ps = append(ps, a.Path)
		}
	}
	return ps
}
