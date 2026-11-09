package mission

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crow-core/contract"
	"crow-core/internal/eventlog"
)

// CapQuestion is the id of the question core asks when a mission hits its resume cap.
const CapQuestion = "resume-cap"

// Choices of the resume cap question.
const (
	CapContinue = "continue"
	CapStop     = "stop"
)

// streak is what the event log says about the resumes of a mission since the last reset:
// a gate (the run reached a version the owner reviews) or a resume cap question core
// asked. capAsked is true when that reset is a cap question.
type streak struct {
	resumes  int
	capAsked bool
}

func (env Env) resumeStreak(id string) (streak, error) {
	ls, err := rawLines(eventlog.Path(env.Store.Root(), id))
	if err != nil {
		return streak{}, err
	}
	var s streak
	for _, l := range ls {
		switch l.Event.Type {
		case "gate":
			s = streak{}
		case "question":
			if l.Event.Actor == "core" {
				s = streak{capAsked: true}
			}
		case "run_start":
			var d contract.RunStartData
			if err := json.Unmarshal(l.Event.Data, &d); err != nil {
				return s, fmt.Errorf("%s line %d: %w", id, l.No, err)
			}
			if d.ResumeOf != nil {
				s.resumes++
			}
		}
	}
	return s, nil
}

// capped reports whether the next resume of m goes over its resume cap: budget.resumes
// resumes in a row without a gate. The first resume after the owner answered the cap
// question always runs, so a cap of 0 asks before every resume.
func (s streak) capped(m contract.Mission) bool {
	if s.capAsked && s.resumes == 0 {
		return false
	}
	return s.resumes >= m.Budget.Resumes
}

// askResumeCap parks a queued mission on a core question instead of resuming it. The
// question goes to outbox/questions.json under the current version, so the answer lands
// in the same answers file the run will read.
func (env Env) askResumeCap(m contract.Mission, s streak) (*TickLine, error) {
	dir := filepath.Join(env.Store.Root(), m.ID)
	text := fmt.Sprintf("%s resumed %d times in a row without reaching a gate (cap %d). Resume it again?",
		m.ID, s.resumes, m.Budget.Resumes)
	rec := CapContinue
	q := contract.Questions{Contract: contract.Version, Mission: m.ID, Version: m.Version, Questions: []contract.Question{{
		ID: CapQuestion, Text: text, Blocking: true, Recommended: &rec,
		Choices: &[]contract.Choice{{ID: CapContinue, Label: "Resume once more"}, {ID: CapStop, Label: "Cancel the mission"}},
	}}}
	b, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, QuestionsFile)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o600); err != nil {
		return nil, err
	}
	ts := env.Now().Format(time.RFC3339)
	reason := "parked: resume cap"
	if err := env.append(m.ID, "core", "state", ts, contract.StateData{From: "queued", To: "waiting", Reason: &reason}); err != nil {
		return nil, err
	}
	qd := contract.QuestionData{File: QuestionsFile, Version: m.Version, IDs: []string{CapQuestion}}
	if err := env.append(m.ID, "core", "question", ts, qd); err != nil {
		return nil, err
	}
	return &TickLine{Mission: m.ID, State: "waiting", Detail: text}, nil
}

// capAnswer is the owner's words on the last cap question, read from the answers file of
// the current version, or "".
func (env Env) capAnswer(m contract.Mission) (string, error) {
	points, err := readPoints(filepath.Join(env.Store.Root(), m.ID, AnswersFile(m.Version)))
	if err != nil {
		return "", err
	}
	words := ""
	for _, p := range points {
		if w, ok := strings.CutPrefix(p, CapQuestion+": "); ok {
			words = w
		}
	}
	return words, nil
}

// stopAfterCap cancels a mission whose owner answered stop to the cap question.
func (env Env) stopAfterCap(m contract.Mission) (*TickLine, error) {
	ts := env.Now().Format(time.RFC3339)
	reason := "owner stopped it at the resume cap"
	if err := env.append(m.ID, "core", "state", ts, contract.StateData{From: "queued", To: "cancelled", Reason: &reason}); err != nil {
		return nil, err
	}
	return &TickLine{Mission: m.ID, State: "cancelled", Detail: reason}, nil
}

// beforeResume settles what the tick does with an answered mission before it resumes:
// cancel it when the owner answered stop at the resume cap, park it on a cap question
// when it resumed budget.resumes times in a row, or nil to resume it.
func (env Env) beforeResume(m contract.Mission) (*TickLine, error) {
	s, err := env.resumeStreak(m.ID)
	if err != nil {
		return nil, err
	}
	if s.capAsked && s.resumes == 0 {
		words, err := env.capAnswer(m)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(words)), CapStop) {
			return env.stopAfterCap(m)
		}
	}
	if s.capped(m) {
		return env.askResumeCap(m, s)
	}
	return nil, nil
}
