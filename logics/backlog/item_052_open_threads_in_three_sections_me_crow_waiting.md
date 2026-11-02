## item_052_open_threads_in_three_sections_me_crow_waiting - Open threads in three sections: me, crow, waiting
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Small
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Restructure open-threads.md into ## Me / ## Crow / ## Waiting sections (owner pick, option A) and make every writer follow it.
- Keywords: open threads, sections, wrap, eod, brief, afk, actor, stale
- Use when: changing how a skill writes or reads open threads.
- Skip when: the migration of existing lines (item_053) or the deck parser (item_054).

# Problem
- One free-text list mixes owner to-dos, agent work, waiting-on items and status notes; the owner cannot see what needs them.
- Model chosen 2026-11-02: option A in `local/state/research/owner-todo-model.md`.

# Scope
- In:
  - Format spec in a runbook (`runbooks/open-threads.md`, generic): three sections, line shape `- [ ] YYYY-MM-DD [project] text → next`, Waiting lines start with the person, Crow lines point to a Logics item when one exists.
  - Rules: done means gone (ticked line leaves the file in the same write and goes to `activity/YYYY-MM-DD.md`); status notes go to status.md, constraints to decisions.md; untouched for 14 days = stale.
  - Writers: skills/wrap, .claude/skills/eod, runbooks that append threads.
  - Readers: .claude/skills/brief (Me first, then Waiting with age, stale list asking keep/re-date/drop), .claude/skills/afk (Crow section only).
  - Keep it a plain checkbox list Obsidian Tasks reads (vault move planned).
- Out:
  - Migrating current lines (item_053).
  - Deck Go parser (item_054).

# Acceptance criteria
- AC1: The runbook defines the three sections, the line shape, the done and stale rules, and is linked from runbooks/README.md.
- AC2: /wrap and /eod write new items into the right section and remove ticked ones into the activity log.
- AC3: /brief shows Me and Waiting with age and lists stale items; /afk reads Crow only.
- AC4: No skill writes a status note into the file.

# AC Traceability
- request-AC2 -> AC1, AC2. Proof: `runbooks/open-threads.md` (sections, line shape, done and stale rules), linked in `runbooks/README.md`; `skills/wrap/SKILL.md` and `.claude/skills/eod/SKILL.md` write by section and move done lines to `activity/`.
- request-AC3 -> AC2, AC3. Proof: `.claude/skills/brief/SKILL.md` (Me, Waiting with age, stale list), `.claude/skills/afk/SKILL.md` (Crow only); AC4: no skill writes status notes, stated in wrap, eod, brief and the runbook.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_012_owner_todos_and_open_work_model
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Every brief and the deck Next list depend on it; small change.

# Notes
- Execution order: 1 of 3. Sonnet is enough.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

