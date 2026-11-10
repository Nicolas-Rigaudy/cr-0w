package api

import (
	"encoding/json"
	"strings"
	"time"
)

// Job is one crow timer of the user systemd instance.
type Job struct {
	Unit      string  `json:"unit"`
	Activates string  `json:"activates"`
	Next      *string `json:"next"`
	Last      *string `json:"last"`
}

// ParseTimers reads `systemctl --user list-timers -o json` and keeps the crow-* timers.
// systemd gives times in microseconds since the epoch, 0 when unknown.
func ParseTimers(b []byte) ([]Job, error) {
	var rows []struct {
		Unit, Activates string
		Next, Last      int64
	}
	if err := json.Unmarshal(b, &rows); err != nil {
		return nil, err
	}
	jobs := []Job{}
	for _, r := range rows {
		if !strings.HasPrefix(r.Unit, "crow-") {
			continue
		}
		jobs = append(jobs, Job{Unit: r.Unit, Activates: r.Activates, Next: usTime(r.Next), Last: usTime(r.Last)})
	}
	return jobs, nil
}

func usTime(us int64) *string {
	if us <= 0 {
		return nil
	}
	s := time.UnixMicro(us).UTC().Format(time.RFC3339)
	return &s
}
