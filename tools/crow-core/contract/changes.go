package contract

// Changes is the data of GET /api/changes?cwd=<path> (item_273): the git state of a live
// lane's checkout against the repo's main branch, read on demand when the lane sheet opens.
type Changes struct {
	Cwd         string   `json:"cwd"`
	Branch      string   `json:"branch"` // empty on a detached HEAD
	Base        string   `json:"base"`   // main or master
	Ahead       int      `json:"ahead"`  // commits on HEAD not on base
	Behind      int      `json:"behind"` // commits on base not on HEAD
	Uncommitted int      `json:"uncommitted"`
	DiffStat    DiffStat `json:"diffstat"` // git diff --numstat <base>...HEAD
	Commits     []Commit `json:"commits"`  // git log <base>..HEAD, newest first
	More        bool     `json:"more"`     // more commits than Commits lists
}

// Commit is one commit of a lane's branch.
type Commit struct {
	SHA     string `json:"sha"`
	Subject string `json:"subject"`
	Time    string `json:"time"` // committer date, RFC 3339
}
