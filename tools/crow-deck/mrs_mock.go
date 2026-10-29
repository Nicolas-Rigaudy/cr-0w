package main

// Invented MR data for the mock deck and the tests. Names, ids, diffs and findings are
// made up and kept consistent with the other mock pages (atlas !213 !214 !220, 1 awaiting,
// 2 mine, 1 thread).

import "strings"

func mockMRs() []mrRow {
	return []mrRow{
		{Project: "atlas", Name: "atlas", Scope: "work", Forge: "gitlab", Repo: "grp/atlas", Number: 214, Ref: "!214", Title: "typed ids", Author: "you", Role: "AUTHOR", Pipeline: "running", Updated: "2026-10-29T12:05:00Z", Branch: "feat/typed-ids", Target: "main"},
		{Project: "atlas", Name: "atlas", Scope: "work", Forge: "gitlab", Repo: "grp/atlas", Number: 213, Ref: "!213", Title: "init prod", Author: "you", Role: "AUTHOR", Draft: true, Pipeline: "passed", ThreadsOpen: 1, ThreadsResolved: 1, Updated: "2026-10-29T11:40:00Z", Branch: "feat/init-prod", Target: "main"},
		{Project: "atlas", Name: "atlas", Scope: "work", Forge: "gitlab", Repo: "grp/atlas", Number: 220, Ref: "!220", Title: "retry policy on webhook", Author: "Sam", Role: "REVIEWER", Pipeline: "passed", Updated: "2026-10-28T13:50:00Z", Branch: "feat/retry", Target: "main"},
	}
}

func mockBranches() []mrBranch {
	return []mrBranch{{Project: "vault", Scope: "personal", Branch: "feat/export-pdf", Ahead: 3}}
}

func mockMRCounts(f *fixtures) {
	f.MRsAwaiting, f.MRsMine, f.MRsThreads = mrCounts(f.MRs)
}

func dl(old, new int, kind byte, text string) diffLine {
	return diffLine{Old: old, New: new, Kind: kind, Text: text}
}

func mockDetails() map[string]*mrDetail {
	retry := mrFile{Path: "retry.py", Add: 6, Del: 1, Lines: []diffLine{
		dl(40, 40, ' ', ""), dl(41, 41, ' ', "def deliver(event):"), dl(42, 0, '-', "    return send(event)"),
		dl(0, 42, '+', "    for attempt in range(MAX_RETRIES):"), dl(0, 43, '+', "        if send(event): return True"),
		dl(0, 44, '+', "        time.sleep(BACKOFF * 2 ** attempt)"), dl(0, 45, '+', "    return False"), dl(43, 46, ' ', ""),
		dl(0, 47, '+', `BACKOFF = float(os.environ.get("RETRY_BACKOFF", "0.5"))`), dl(0, 48, '+', `MAX_RETRIES = int(os.environ.get("RETRY_MAX", "5"))`),
	}}
	cfg := mrFile{Path: "config.md", Add: 2, Lines: []diffLine{dl(10, 10, ' ', "## Webhook"), dl(0, 11, '+', "Delivery retries with a doubling backoff."), dl(0, 12, '+', "See retry.py.")}}
	test := mrFile{Path: "test_retry.py", Add: 76, Lines: []diffLine{dl(0, 1, '+', "def test_retries_then_gives_up(monkeypatch):"), dl(0, 2, '+', "    monkeypatch.setattr(retry, \"send\", lambda e: False)"), dl(0, 3, '+', "    assert retry.deliver({}) is False")}}
	seed := mrFile{Path: "seed.py", Add: 3, Del: 1, Lines: []diffLine{
		dl(56, 56, ' ', "def seed_prod(account):"), dl(57, 57, ' ', "    tenant = get_tenant(account)"), dl(58, 0, '-', "    create_all(tenant, SEED)"),
		dl(0, 58, '+', "    if tenant_has_data(tenant):"), dl(0, 59, '+', "        raise SeedRefused(account)"), dl(0, 60, '+', "    create_all(tenant, SEED)"),
	}}
	return map[string]*mrDetail{
		"grp/atlas!220": {
			Files: []mrFile{retry, cfg, test},
			Findings: []mrFinding{
				{Title: "env var undocumented", Path: "retry.py", Line: 47, Body: "BACKOFF and MAX_RETRIES come from the environment and are not in the config doc.", Pattern: "env var undocumented", Seen: 3, Comment: "Pense à documenter RETRY_BACKOFF et MAX_RETRIES dans config.md.", State: "pending"},
				{Title: "default value not tested", Path: "retry.py", Line: 48, Body: "No test covers the default MAX_RETRIES.", Pattern: "default not tested", Seen: 2, Comment: "Ajoute un test sur la valeur par défaut de MAX_RETRIES.", State: "accepted"},
				{Title: "commit scope", Path: "", Line: 0, Body: "The commit subject has no file scope.", Pattern: "commit scope", Seen: 4, Comment: "Le sujet du commit n'a pas de scope.", State: "refused"},
			},
		},
		"grp/atlas!213": {
			Files: []mrFile{seed},
			Threads: []mrThread{
				{Author: "Kim", At: "yesterday 16:40", Path: "seed.py", Line: 58, Body: "Le seed ne vérifie pas que le compte prod est vide. On peut écraser des données si on le relance.",
					Pattern: "destructive seed, no guard", Seen: 2,
					Fix: &mrSection{Title: "the fix", Meta: "seed.py:58 · +3 -1", Lines: seed.Lines},
					Test: &mrSection{Title: "the test", Meta: "test_seed.py (new) · +9", Lines: []diffLine{
						dl(0, 1, '+', "def test_seed_refused_on_non_empty_tenant(tenant):"), dl(0, 2, '+', "    tenant.add(Item())"),
						dl(0, 3, '+', "    with pytest.raises(SeedRefused):"), dl(0, 4, '+', "        seed_prod(tenant.account)")}},
					Reply: "Bien vu. J'ajoute une garde qui refuse le seed si le compte contient déjà des données, avec un test."},
				{Author: "Kim", At: "Tue 10:12", Path: "README.md", Line: 4, Body: "README wording.", Resolved: true},
			},
		},
	}
}

// wrapLines wraps text to w and returns one string per line.
func wrapLines(s string, w int) []string {
	var out []string
	for _, para := range strings.Split(s, "\n") {
		if strings.TrimSpace(para) == "" {
			out = append(out, "")
			continue
		}
		out = append(out, wrap(para, w)...)
	}
	return out
}
