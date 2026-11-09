package mission

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crow-core/contract"
	"crow-core/internal/brief"
)

// Resume modes recorded in the tick line and the run_start event.
const (
	ResumeContinue = "continue" // cdx --continue-run of the last run, same account and sandbox
	ResumeCold     = "cold"     // a fresh run from brief.md, which holds every answers file
)

// ResumeNote is the prompt file of a continued run: the answers since the last run.
const ResumeNote = "resume.md"

// coldFallback lists the cdx refusals of --continue-run that a cold run still serves
// (item_116): the isolated workspace or its conversation is gone, the sandbox settings
// changed, the transcript expired or the run is unknown to cdx.
var coldFallback = map[string]bool{
	"continue_run_workspace_gone":    true,
	"continue_run_settings_mismatch": true,
	"continue_run_unavailable":       true,
	"continue_run_not_found":         true,
}

// lastRun is the most recent ended run of a mission, or nil.
func lastRun(m contract.Mission) *contract.Run {
	for i := len(m.Runs) - 1; i >= 0; i-- {
		if m.Runs[i].Ended != nil {
			return &m.Runs[i]
		}
	}
	return nil
}

// resume starts the next run of an answered mission. It continues the last run natively
// unless cold is asked (the mission moves account, item_115 "move"); a continuation cdx
// refuses with a coldFallback code becomes a cold run, logged as an incident. Either way
// brief.md is rewritten first, so a cold run on any account restarts from files alone.
// A nil line with an error means nothing started.
func (env Env) resume(r Runner, m contract.Mission, k *contract.Kind, cold bool) (*TickLine, error) {
	prev := lastRun(m)
	if prev == nil {
		return env.start(r, m, k)
	}
	dir, in, err := env.Brief(m.ID)
	if err != nil {
		return nil, err
	}
	if err := brief.Write(dir, in); err != nil {
		return nil, err
	}
	label := in.Key
	if label == "" {
		label = brief.Key(m)
	}
	next := brief.Version(in.Answers) // the version the brief asks for
	mode, fellBack := ResumeCold, ""
	var s Started
	if !cold {
		note, err := writeResumeNote(dir, prev.Version, next, in.Answers)
		if err != nil {
			return nil, err
		}
		d := dispatchFor(m, k, dir, label, note)
		d.ContinueRun = prev.Run
		s, err = r.Dispatch(d)
		var de *DispatchError
		switch {
		case err == nil:
			mode = ResumeContinue
		case errors.As(err, &de) && coldFallback[de.Code]:
			fellBack = de.Error()
		default:
			return nil, err
		}
	}
	if mode == ResumeCold {
		d := dispatchFor(m, k, dir, label, filepath.Join(dir, "brief.md"))
		d.Resumes = prev.Run
		if s, err = r.Dispatch(d); err != nil {
			return nil, err
		}
	}
	ts := env.Now().Format(time.RFC3339)
	out := &TickLine{Mission: m.ID, Run: s.RunID}
	if fellBack != "" {
		detail := "continue of " + prev.Run + " refused, resumed cold: " + fellBack
		if err := env.append(m.ID, "core", "incident", ts, contract.IncidentData{Cause: "continue_refused", Detail: &detail}); err != nil {
			return out, err
		}
	}
	if err := env.append(m.ID, "core", "state", ts, contract.StateData{From: "queued", To: "running"}); err != nil {
		return out, err
	}
	of := prev.Run
	sd := contract.RunStartData{Run: s.RunID, AccountClass: m.Scope, Cold: mode == ResumeCold, Version: next, ResumeOf: &of}
	if s.Account != "" {
		sd.Account = &s.Account
	}
	if err := env.append(m.ID, "core", "run_start", ts, sd); err != nil {
		return out, err
	}
	out.State, out.Detail = "running", fmt.Sprintf("%s resume of %s on %s", mode, prev.Run, s.Account)
	return out, nil
}

// writeResumeNote writes the prompt of a continued run: every answer given since version
// last, quoted, and the version to write next. It names no path: the answers are inline.
func writeResumeNote(dir string, last, next int, answers []brief.Answer) (string, error) {
	var b strings.Builder
	b.WriteString("# Resume\n\nThe owner answered. Read the answers below, then continue the mission from the brief you already have.\n")
	n := 0
	for _, a := range answers {
		if a.Version < last {
			continue
		}
		n++
		fmt.Fprintf(&b, "\n## Answer to version %d (%s)\n\n", a.Version, a.Verb)
		if a.Words != "" {
			b.WriteString(a.Words + "\n")
		} else {
			b.WriteString("(no words)\n")
		}
	}
	if n == 0 {
		return "", fmt.Errorf("no answers file since version %d to resume from", last)
	}
	fmt.Fprintf(&b, "\nWrite version %d next, with the same rules as the brief. Quote each answer in its Responses.\n", next)
	path := filepath.Join(dir, ResumeNote)
	return path, os.WriteFile(path, []byte(b.String()), 0o600)
}
