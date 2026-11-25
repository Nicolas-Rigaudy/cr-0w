package mission

// Design reviews from a crow-compare spec (item_211).
//
// A mockup or report may list a crow-compare spec (bin/crow-compare) among its artifacts.
// The inbox reads its first decision: the option labels and notes, the scenario rows and
// one image per option and row. Image paths in the spec are relative to the spec file and
// must stay under outbox/; the others are dropped.

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"crow-core/contract"
)

// compareMax is the largest spec read, in bytes.
const compareMax = 1 << 20

type compareSpec struct {
	Decisions []struct {
		Question  string   `json:"question"`
		Recommend string   `json:"recommend"`
		Rows      []string `json:"rows"`
		Options   []struct {
			Label  string    `json:"label"`
			Note   string    `json:"note"`
			Images []*string `json:"images"`
		} `json:"options"`
	} `json:"decisions"`
}

// compareOf reads the first crow-compare spec among the artifacts. It returns nil when there
// is none or it holds no option with an image.
func compareOf(dir string, artifacts []string) *contract.Compare {
	for _, a := range artifacts {
		if !strings.HasSuffix(strings.ToLower(a), ".json") {
			continue
		}
		if c := readCompare(dir, a); c != nil {
			return c
		}
	}
	return nil
}

func readCompare(dir, spec string) *contract.Compare {
	full, err := outboxFile(dir, spec)
	if err != nil {
		return nil
	}
	if fi, err := os.Stat(full); err != nil || fi.Size() > compareMax {
		return nil
	}
	b, err := os.ReadFile(full)
	if err != nil {
		return nil
	}
	var s compareSpec
	if json.Unmarshal(b, &s) != nil || len(s.Decisions) == 0 {
		return nil
	}
	d := s.Decisions[0]
	rows := d.Rows
	if len(rows) == 0 {
		rows = []string{""}
	}
	c := &contract.Compare{Question: d.Question, Recommend: d.Recommend, Rows: rows, Options: []contract.CompareOption{}}
	seen := false
	for _, o := range d.Options {
		if strings.TrimSpace(o.Label) == "" {
			continue
		}
		imgs := make([]*string, len(rows))
		for i := range rows {
			if i < len(o.Images) && o.Images[i] != nil {
				imgs[i] = compareImage(dir, spec, *o.Images[i])
				seen = seen || imgs[i] != nil
			}
		}
		c.Options = append(c.Options, contract.CompareOption{Label: o.Label, Note: o.Note, Images: imgs})
	}
	if !seen {
		return nil
	}
	return c
}

// compareImage resolves an image path of the spec to an outbox path that exists, or nil.
func compareImage(dir, spec, p string) *string {
	if p == "" || path.IsAbs(p) {
		return nil
	}
	q := path.Clean(path.Join(path.Dir(spec), p))
	if _, err := outboxFile(dir, q); err != nil {
		return nil
	}
	return &q
}

// compareImages lists the images of c, in option then row order.
func compareImages(c *contract.Compare) []string {
	var ps []string
	for _, o := range c.Options {
		for _, p := range o.Images {
			if p != nil {
				ps = append(ps, *p)
			}
		}
	}
	return ps
}
