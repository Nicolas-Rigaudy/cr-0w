## item_129_move_owner_notes_into_the_vault_and_switch_on_crow_dir - Move owner notes into the vault and switch on crow_dir
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub data
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: move, owner, notes, vault, switch, crow, dir
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- Open threads, status, decisions and retro reports still sit in local/state.

# Scope
- In:
  - One-shot script: open-threads.md to CR-0W/Open Threads.md, status.md to Status.md, decisions.md split into Decisions/<date> <title>.md (type: decision), retro reports to Retros/, past activity lines into the matching daily notes' CR-0W section. Frontmatter on create, originals kept in the local repo history.
  - Set vault.crow_dir: CR-0W in local/projects.yaml.
  - Update skills, runbooks and CLAUDE.md paths (wrap, eod, brief, retro, afk, open-threads runbook).
  - Run a brief gather, a crow-log line and a thread edit end to end.
- Out:
  - Retention windows
  - Vault dashboard and Bases

# Acceptance criteria
- Vault CR-0W holds Open Threads.md, Status.md, Decisions/, Retros/; local/state holds no owner markdown.
- crow-gather brief and crow-deck show the same threads as before.
- Tests pass.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: Vault CR-0W holds Open Threads.md, Status.md, Decisions/, Retros/; local/state holds no owner markdown.
- request-AC4 -> This backlog slice. Proof: crow-gather brief and crow-deck show the same threads as before.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_006_crow_data_in_the_vault`
- Architecture decision(s): (none yet)
- Request: `req_041_crow_data_in_the_vault`
- Primary task(s): `task_015_build_crow_data_in_the_vault`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

