## item_278_desk_hides_the_butler_from_its_session_list - Desk hides the butler from its session list
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The butler's own Claude process shows up in Desk's session and run lists. Picking it shows the butler conversation in the main panel, next to the same conversation in the side panel. Hide it from Desk.
- Keywords: nest, desk, butler, sessions, duplicate
- Use when: which agents Desk lists.
- Skip when: butler chat content.

# Problem
- `crow runs` lists the butler as `crow-butler (not dispatched by crow)` under background runs: the session scan picks up its headless `claude -p` (cwd `local/state/crow-butler`, `--system-prompt-file local/state/butler/persona-*.md`).
- Desk lists it like any other agent, so the butler conversation shows twice on screen: side panel and main panel.

# Scope
- In:
  - Tag the butler process where agents are gathered (crow-sessions / core dispatch section), from its cwd or persona file, as `kind: butler`.
  - Desk, Ops and the hero skip `kind: butler`. `crow runs` and `crow doctor` show it on one line as the butler, not as a run.
- Out:
  - The side panel chat itself.

# Acceptance criteria
- AC1: With the butler live, Desk lists no butler row at 1440x900 and 900x700.
- AC2: `crow runs` shows the butler once, labelled butler, not under background runs.
- AC3: Tests cover the tag (pytest or go test); tsc and vitest pass.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`. Related item_176 (butler chat).

# Priority
- Priority: Medium
- Rationale: confusing, not blocking.

# Notes
- Owner 2026-11-25: the butler chat shows both in the side panel and in the Desk session list. Desk's own ask box stays.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
