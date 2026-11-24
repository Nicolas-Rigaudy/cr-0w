package contract

// RunReport is the data of GET /api/runs/{id}: one crow-runs run from its ledger, or a
// mission run from the core store.
type RunReport struct {
	Run      string    `json:"run"`
	Source   string    `json:"source"` // dispatch or mission
	Project  string    `json:"project"`
	Label    string    `json:"label"`
	Account  string    `json:"account"`
	Status   string    `json:"status"`
	Started  string    `json:"started"`
	Ended    *string   `json:"ended"`
	CostUSD  *float64  `json:"cost_usd"`
	Commits  []string  `json:"commits"`
	Branch   *string   `json:"branch"`
	Worktree *string   `json:"worktree"`
	Item     *string   `json:"item"`
	Mission  *string   `json:"mission"`
	Report   string    `json:"report"`
	DiffStat *DiffStat `json:"diffstat"`
}

// DiffStat is `git diff --numstat <start>...HEAD` in a run's cwd, stored by crow-runs when
// the run ends so it outlives the worktree.
type DiffStat struct {
	Base       string     `json:"base"`
	Files      int        `json:"files"`
	Insertions int        `json:"insertions"`
	Deletions  int        `json:"deletions"`
	PerFile    []FileStat `json:"per_file"`
	Truncated  bool       `json:"truncated"`
}

// FileStat is one file of a diff stat; binary files count 0.
type FileStat struct {
	Path       string `json:"path"`
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
}
