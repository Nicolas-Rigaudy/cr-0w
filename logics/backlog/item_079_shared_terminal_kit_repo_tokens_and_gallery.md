## item_079_shared_terminal_kit_repo_tokens_and_gallery - Shared terminal kit repo: tokens and gallery
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: design
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Create the public Go module: ANSI-slot token roles, glyphs, borders, monochrome mode, gallery command.
- Keywords: shared, terminal, kit, repo, tokens, gallery
- Use when: Creating the kit repo or changing tokens.
- Skip when: Component extraction details.

# Problem
- Tokens, glyphs and borders are redrawn per app and drift.

# Scope
- In:
  - New public Go module on the personal GitHub (name proposed to the owner before creation; check git identity first).
  - tokens: roles (accent, muted, border, focus, selection, ok, warn, error) mapped to ANSI slots; glyph sets with fallbacks (sparkline vs block); border styles; monochrome shape-only mode.
  - gallery command: prints every registered component x state x size as ANSI, plus a manifest consumed by ansi-shot and crow-compare.
  - Design-system round 1 through the new workflow: tokens, glyphs, borders, monochrome.
- Out:
  - Components beyond what round 1 needs to demonstrate tokens.

# Acceptance criteria
- AC1: Repo builds, README explains roles and the gallery.
- AC2: Gallery output renders through ansi-shot into a crow-compare page.
- AC3: Round 1 decisions recorded in the ledger.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: AC1: Repo builds, README explains roles and the gallery.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_004_terminal_design_system`
- Architecture decision(s): (none yet)
- Request: `req_020_terminal_design_system_and_design_review_workflow`
- Primary task(s): `task_010_orchestrate_terminal_design_system`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

