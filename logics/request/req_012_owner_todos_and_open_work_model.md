## req_012_owner_todos_and_open_work_model - Rework how crow tracks the owner's to-dos and unfinished work (open threads)
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Replace the single open-threads log with a model that separates the owner's to-dos, work an agent still owns, things waiting on someone, and status notes, so briefs and the deck can show what needs the owner.
- Keywords: open threads, todo, next actions, wrap, brief, deck Next
- Use when: changing how /wrap, /brief, /eod or the deck record or read unfinished work.
- Skip when: Logics per-repo task state.

# Needs
- `local/state/open-threads.md` mixes done items ("MR A merged"), status notes, owner to-dos and agent follow-ups in one free-text list. The owner cannot tell what needs them; the deck Next list and project cards inherit the mess.
- Owner 2026-11-02: "why do open threads list stuff that's done? ... needs a deeper rework on how we handle your understanding of my TO DO and how you handle open items that we haven't finished working on, maybe could even be different sections".

# Context
- Writers today: /wrap, /eod, runbooks, sessions. Readers: /brief, deck (`data.go:openThreads`), crow-activity.
- Vault (`vault.crow_dir`) is the long-term home for state.

# Scope
- In: the model (sections or fields: owner to-do, agent-owned unfinished work, waiting on someone, done/status log), who writes what, how items close, migration of the current file, readers (brief, deck Next, project NEXT).
- Out: the deck layout (req_011).

# Acceptance criteria
- AC1: A design doc proposes the model with 2-3 options and the owner picks one.
- AC2: Every item says who acts and since when; done items leave the active list when they are done.
- AC3: /wrap, /eod and /brief follow the new model; the current file is migrated.
- AC4: The deck reads the new source and shows actor and age.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed: owner choice of model (AC1).

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `local/state/open-threads.md`
- `.claude/skills/brief`, `.claude/skills/eod`, `skills/wrap`
- `tools/crow-deck/data.go` (openThreads)

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
