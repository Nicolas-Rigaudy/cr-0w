## item_122_crow_home_is_the_private_config_repo_itself - CROW_HOME is the private config repo itself
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 17:14:58

# AI Context
- Summary: An installed crow uses the private config repo clone as CROW_HOME, files at its root, with no extra local/ level.
- Keywords: crow, home, private, config, repo, itself
- Use when: Working on install, restore, data lookup or publishing crow publicly.
- Skip when: Engine-only features that never read owner data.

# Problem
- item_034 shipped ~/.config/crow/local/... In a checkout, local/ is already the private repo (cr-0w-private), so once installed the parent folder holds nothing else. Owner decision 2026-11-07: CROW_HOME = the private repo clone itself.

# Scope
- In:
  - Data lookup in bin/crow_lib.py and every Python script: on an installed engine, data files sit at the CROW_HOME root.
  - A checkout keeps reading local/ next to the code (fallback).
  - install.sh --restore clones the private repo straight into CROW_HOME.
  - Shell hooks (commit-msg, pre-push, pre-commit, afk-preflight) resolve data through CROW_HOME too (gap left by item_034).
  - Migration note in runbooks/install-and-update.md.
- Out:
  - Release workflow (item_035) and package recipes (item_036).

# Acceptance criteria
- AC1: With CROW_HOME set to a clone of the private repo, every crow command finds its data at the clone root.
- AC2: A checkout with local/ and no CROW_HOME behaves as today.
- AC3: install.sh --restore <git-url> produces a working CROW_HOME in one step (docker smoke test).
- AC4: The shell hooks honour CROW_HOME.
- AC5: Tests cover both layouts; nothing outside local/ names the owner.

# AC Traceability
- AC1 -> this item. Proof: to record at closeout.
- AC2 -> this item. Proof: to record at closeout.
- AC3 -> this item. Proof: to record at closeout.
- AC4 -> this item. Proof: to record at closeout.
- AC5 -> this item. Proof: to record at closeout.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Blocks publishing crow publicly with private config kept apart.

# Notes
- Owner decision 2026-11-07 on return from the AFK run (local/state/afk/decisions-2026-11-07-afk.md).
- Before each push of the private config repo: whitelist of paths plus a secret scanner, refuse the push on a hit (pattern from the harness-maxxing backup, 2026-11-07).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
