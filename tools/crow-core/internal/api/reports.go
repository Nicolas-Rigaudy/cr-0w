package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ReportsSection lists what is settled (item_258): AFK reports whose decisions are all
// taken, the weekly rollups and the retros. Read only.
const ReportsSection = "reports"

// reportCap bounds the rows and the body of each row.
const (
	reportRows = 60
	reportBody = 64 << 10
)

// ReportDirs are the folders read, resolved like crow_lib.crow_path. Vault is the vault
// root ("" when the notes sit outside a vault): notes under it get an obsidian:// link.
type ReportDirs struct {
	AFK     string `json:"afk"`
	Rollups string `json:"rollups"`
	Retros  string `json:"retros"`
	Vault   string `json:"-"`
}

// Report is one settled output.
type Report struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"` // afk | weekly | retro
	Title   string `json:"title"`
	At      string `json:"at"`
	Summary string `json:"summary,omitempty"`
	Path    string `json:"path"`
	Link    string `json:"link,omitempty"`
	Body    string `json:"body"`
}

// Reports is the section body. Held counts AFK reports left out for an open decision
// (they wait in Inbox).
type Reports struct {
	Items []Report   `json:"items"`
	Held  int        `json:"held"`
	Dirs  ReportDirs `json:"dirs"`
}

// BuildReports reads the three folders, newest first. A missing folder is empty.
func BuildReports(d ReportDirs) (Reports, error) {
	out := Reports{Items: []Report{}, Dirs: d}
	afk, held, err := afkReports(d.AFK)
	if err != nil {
		return out, err
	}
	out.Held = held
	out.Items = append(out.Items, afk...)
	for _, k := range []struct{ kind, dir string }{{"weekly", d.Rollups}, {"retro", d.Retros}} {
		rs, err := noteReports(k.kind, k.dir, d.Vault)
		if err != nil {
			return out, err
		}
		out.Items = append(out.Items, rs...)
	}
	sort.SliceStable(out.Items, func(i, j int) bool { return out.Items[i].At > out.Items[j].At })
	if len(out.Items) > reportRows {
		out.Items = out.Items[:reportRows]
	}
	return out, nil
}

type afkDecision struct {
	Q string `json:"q"`
}

type afkPoint struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Where string   `json:"where"`
	Runs  []string `json:"runs"`
}

type afkProject struct {
	Name      string        `json:"name"`
	TLDR      string        `json:"tldr"`
	Points    []afkPoint    `json:"points"`
	Decisions []afkDecision `json:"decisions"`
	Decided   []string      `json:"decided"`
}

type afkSummary struct {
	Start    string          `json:"start"`
	End      string          `json:"end"`
	Projects json.RawMessage `json:"projects"`
}

// projects reads projects as a list or as a map keyed by name, sorted by name.
func (s afkSummary) projects() ([]afkProject, error) {
	raw := strings.TrimSpace(string(s.Projects))
	if raw == "" || raw == "null" {
		return nil, nil
	}
	var ps []afkProject
	if strings.HasPrefix(raw, "[") {
		if err := json.Unmarshal(s.Projects, &ps); err != nil {
			return nil, err
		}
		return ps, nil
	}
	var m map[string]afkProject
	if err := json.Unmarshal(s.Projects, &m); err != nil {
		return nil, err
	}
	for name, p := range m {
		if p.Name == "" {
			p.Name = name
		}
		ps = append(ps, p)
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Name < ps[j].Name })
	return ps, nil
}

// settled: every decision of the project is named in its decided[].
func settled(p afkProject) bool {
	for _, d := range p.Decisions {
		q := strings.TrimSpace(d.Q)
		found := false
		for _, t := range p.Decided {
			if q != "" && strings.Contains(t, q) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// afkReports reads summary-<date>.json files. A report with an open decision is held.
func afkReports(dir string) ([]Report, int, error) {
	files, err := filepath.Glob(filepath.Join(dir, "summary-*.json"))
	if err != nil || dir == "" {
		return nil, 0, err
	}
	var out []Report
	held := 0
	for _, f := range files {
		date := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(f), "summary-"), ".json")
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, 0, err
		}
		var s afkSummary
		if err := json.Unmarshal(b, &s); err != nil {
			continue
		}
		ps, err := s.projects()
		if err != nil {
			continue
		}
		open := false
		for _, p := range ps {
			if !settled(p) {
				open = true
			}
		}
		if open {
			held++
			continue
		}
		out = append(out, afkReport(f, date, s, ps))
	}
	return out, held, nil
}

func afkReport(path, date string, s afkSummary, ps []afkProject) Report {
	at := s.End
	if at == "" {
		at = s.Start
	}
	if at == "" {
		at = date
	}
	points, runs, decided := 0, 0, 0
	var body strings.Builder
	for _, p := range ps {
		fmt.Fprintf(&body, "## %s\n\n", p.Name)
		if p.TLDR != "" {
			fmt.Fprintf(&body, "%s\n\n", p.TLDR)
		}
		for _, pt := range p.Points {
			points++
			runs += len(pt.Runs)
			fmt.Fprintf(&body, "- **%s**", pt.Title)
			if pt.Text != "" {
				fmt.Fprintf(&body, ": %s", pt.Text)
			}
			body.WriteString("\n")
		}
		if len(p.Decided) > 0 {
			body.WriteString("\n### Decided\n\n")
			for _, t := range p.Decided {
				fmt.Fprintf(&body, "- %s\n", t)
			}
		}
		decided += len(p.Decisions)
		body.WriteString("\n")
	}
	if len(ps) == 0 {
		body.WriteString("No project ran.\n")
	}
	sum := fmt.Sprintf("%d points · %d runs · ", points, runs)
	if decided == 0 {
		sum += "no decision"
	} else {
		sum += fmt.Sprintf("%d decisions answered", decided)
	}
	return Report{
		ID: "afk:" + date, Kind: "afk", Title: "AFK night → " + date, At: at,
		Summary: sum, Path: path, Body: body.String(),
	}
}

// noteReports reads the Markdown notes of one ritual folder.
func noteReports(kind, dir, vault string) ([]Report, error) {
	if dir == "" {
		return nil, nil
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}
	var out []Report
	for _, f := range files {
		st, err := os.Stat(f)
		if err != nil || st.IsDir() {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if len(b) > reportBody {
			b = b[:reportBody]
		}
		name := strings.TrimSuffix(filepath.Base(f), ".md")
		title, sum := headOf(string(b))
		if title == "" {
			title = name
		}
		out = append(out, Report{
			ID: kind + ":" + name, Kind: kind, Title: title, At: st.ModTime().Format(time.RFC3339),
			Summary: sum, Path: f, Link: obsidianLink(vault, f), Body: string(b),
		})
	}
	return out, nil
}

// headOf returns the first heading and the first plain line after it. Front matter is skipped.
func headOf(s string) (title, sum string) {
	lines := strings.Split(s, "\n")
	i := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i = 1; i < len(lines) && strings.TrimSpace(lines[i]) != "---"; i++ {
		}
		i++
	}
	for ; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "#") {
			if title == "" {
				title = strings.TrimSpace(strings.TrimLeft(l, "#"))
			}
			continue
		}
		sum = strings.TrimSpace(strings.TrimLeft(l, "-*> "))
		break
	}
	if len(sum) > 120 {
		sum = sum[:120]
	}
	return title, sum
}

// obsidianLink opens a note under the vault root in Obsidian, "" outside it.
func obsidianLink(vault, path string) string {
	if vault == "" {
		return ""
	}
	rel, err := filepath.Rel(vault, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return ""
	}
	rel = strings.TrimSuffix(filepath.ToSlash(rel), ".md")
	return "obsidian://open?vault=" + url.PathEscape(filepath.Base(vault)) + "&file=" + url.PathEscape(rel)
}
