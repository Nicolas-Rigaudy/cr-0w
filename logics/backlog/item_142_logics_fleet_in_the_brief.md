## item_142_logics_fleet_in_the_brief - Logics fleet view in the brief across registered repos
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80
> Confidence: 75
> Progress: 0%
> Complexity: Small
> Theme: Hub upkeep
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 17:28:44

# AI Context
- Summary: the brief adds `logics-manager fleet brief --root <projects dir> --format json` (next moves and blocked items across every corpus) next to its per-repo reads.
- Keywords: logics fleet, brief, next moves, health, corpora
- Use when: building the brief's Logics section.
- Skip when: one repo's state is enough (`logics-manager status` in that repo).

# Problem
- The brief reads each repo's Logics state one by one. logics-manager 2.24.0 added `fleet status|health|roadmap|brief` that scans every corpus under a root in one call.

# Scope
- In:
  - The brief skill calls `fleet brief --format json` once and shows the top moves and blocked items per registered repo. Repos missing from local/projects.yaml are dropped.
  - Keep the call within the brief's token budget (`--limit`).
- Out:
  - Writing to any corpus.

# Acceptance criteria
- AC1: one `fleet` call replaces the per-repo Logics reads in the brief, and the output names the repo of each move.
- AC2: a repo with no corpus or an unregistered one does not appear.

# AC Traceability
- request-AC4 -> This backlog slice. Proof: pending.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_045_periodic_tool_update_watch`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Found in the 2026-11-09 manual tool pass. Small gain.

# Notes
- 2026-11-09: owner picked it from the manual pass.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
