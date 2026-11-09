package contract

import (
	"encoding/json"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// CrowBlock is the fenced crow block of plan.vN.md, report.vN.md or mockup.vN.md.
type CrowBlock struct {
	Contract    int         `json:"contract"`
	Mission     string      `json:"mission"`
	Kind        string      `json:"kind"`
	Phase       string      `json:"phase"`
	Version     int         `json:"version"`
	Summary     string      `json:"summary"`
	Confidence  int         `json:"confidence"`
	Verdict     string      `json:"verdict"`
	Assumptions *[]string   `json:"assumptions,omitempty"`
	Questions   *[]Question `json:"questions,omitempty"`
	Responses   *[]Response `json:"responses,omitempty"`
	Artifacts   *[]Artifact `json:"artifacts,omitempty"`
	Publish     *[]Publish  `json:"publish,omitempty"`
	Proof       *[]Proof    `json:"proof,omitempty"`
	Next        *string     `json:"next,omitempty"`
}

type Response struct {
	Ref     string  `json:"ref"`
	Outcome string  `json:"outcome"`
	Note    *string `json:"note,omitempty"`
}

type Artifact struct {
	Path     string  `json:"path"`
	Kind     *string `json:"kind,omitempty"`
	Option   *string `json:"option,omitempty"`
	Scenario *string `json:"scenario,omitempty"`
	Caption  *string `json:"caption,omitempty"`
}

type Publish struct {
	Action  string    `json:"action"`
	Target  string    `json:"target"`
	Title   *string   `json:"title,omitempty"`
	Body    *string   `json:"body,omitempty"`
	Commits *[]string `json:"commits,omitempty"`
}

type Proof struct {
	What    string  `json:"what"`
	Cmd     *string `json:"cmd,omitempty"`
	Exit    *int    `json:"exit,omitempty"`
	URL     *string `json:"url,omitempty"`
	Checked *string `json:"checked,omitempty"`
	Path    *string `json:"path,omitempty"`
	Note    *string `json:"note,omitempty"`
}

// Questions is outbox/questions.json.
type Questions struct {
	Contract  int        `json:"contract"`
	Mission   string     `json:"mission"`
	Version   int        `json:"version"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	Blocking    bool      `json:"blocking"`
	Choices     *[]Choice `json:"choices,omitempty"`
	Default     *string   `json:"default,omitempty"`
	Recommended *string   `json:"recommended,omitempty"`
}

type Choice struct {
	ID      string  `json:"id"`
	Label   string  `json:"label"`
	Note    *string `json:"note,omitempty"`
	Preview *string `json:"preview,omitempty"`
}

// Kind is one kinds/<name>.yaml file.
type Kind struct {
	Name          string                    `json:"name"`
	SchemaVersion int                       `json:"schema_version"`
	Phase         int                       `json:"phase"`
	Purpose       string                    `json:"purpose"`
	Absorbs       *[]string                 `json:"absorbs,omitempty"`
	Skill         KindSkill                 `json:"skill"`
	Key           KindKey                   `json:"key"`
	Scopes        []string                  `json:"scopes"`
	DataClasses   *[]string                 `json:"data_classes,omitempty"`
	Gates         []string                  `json:"gates"`
	Autonomy      string                    `json:"autonomy"`
	Account       KindAccount               `json:"account"`
	Effort        string                    `json:"effort"`
	Caps          KindCaps                  `json:"caps"`
	Permission    *string                   `json:"permission,omitempty"`
	Worktree      *bool                     `json:"worktree,omitempty"`
	Outputs       []string                  `json:"outputs"`
	Publish       KindPublish               `json:"publish"`
	Preflight     []Preflight               `json:"preflight"`
	ByScope       *map[string]ScopeOverride `json:"by_scope,omitempty"`
}

type KindSkill struct {
	Name      *string   `json:"name"`
	Mode      *string   `json:"mode,omitempty"`
	Fallbacks *[]string `json:"fallbacks,omitempty"`
	Status    *string   `json:"status,omitempty"`
}

type KindKey struct {
	Format  string  `json:"format"`
	Pattern string  `json:"pattern"`
	Example *string `json:"example,omitempty"`
}

type KindAccount struct {
	Route         string    `json:"route"`
	Model         string    `json:"model"`
	SubagentModel *string   `json:"subagent_model,omitempty"`
	Overflow      *[]string `json:"overflow,omitempty"`
}

type KindCaps struct {
	TimeMin       int      `json:"time_min"`
	BudgetUSD     float64  `json:"budget_usd"`
	Resumes       int      `json:"resumes"`
	Rounds        *int     `json:"rounds,omitempty"`
	OutboxFileMB  *float64 `json:"outbox_file_mb,omitempty"`
	OutboxTotalMB *float64 `json:"outbox_total_mb,omitempty"`
}

type KindPublish struct {
	Actions         []string  `json:"actions"`
	NeedsValidation *bool     `json:"needs_validation,omitempty"`
	AskActions      *[]string `json:"ask_actions,omitempty"`
}

type Preflight struct {
	ID      string    `json:"id"`
	Text    string    `json:"text"`
	Choices *[]string `json:"choices,omitempty"`
	Safest  string    `json:"safest"`
}

type ScopeOverride struct {
	Gates             *[]string    `json:"gates,omitempty"`
	Autonomy          *string      `json:"autonomy,omitempty"`
	Publish           *KindPublish `json:"publish,omitempty"`
	Effort            *string      `json:"effort,omitempty"`
	SensitivePathsAsk *bool        `json:"sensitive_paths_ask,omitempty"`
}

var crowFence = regexp.MustCompile("(?ms)^```crow[ \t]*\n(.*?)^```[ \t]*$")

// ExtractCrowBlock returns the body of the one crow block of a version file.
func ExtractCrowBlock(md []byte) ([]byte, error) {
	all := crowFence.FindAllSubmatch(md, -1)
	if len(all) != 1 {
		return nil, fmt.Errorf("want one crow block, found %d", len(all))
	}
	return all[0][1], nil
}

// YAMLToJSON converts a YAML (or JSON) document to JSON, so YAML files decode with the
// same strict JSON rules.
func YAMLToJSON(b []byte) ([]byte, error) {
	var v any
	if err := yaml.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

// ParseCrowBlock extracts and strictly decodes the crow block of a version file.
func ParseCrowBlock(md []byte) (*CrowBlock, error) {
	body, err := ExtractCrowBlock(md)
	if err != nil {
		return nil, err
	}
	j, err := YAMLToJSON(body)
	if err != nil {
		return nil, err
	}
	var cb CrowBlock
	if err := Decode(j, &cb); err != nil {
		return nil, err
	}
	return &cb, nil
}

// ParseKind strictly decodes a kinds/<name>.yaml file.
func ParseKind(b []byte) (*Kind, error) {
	j, err := YAMLToJSON(b)
	if err != nil {
		return nil, err
	}
	var k Kind
	if err := Decode(j, &k); err != nil {
		return nil, err
	}
	return &k, nil
}
