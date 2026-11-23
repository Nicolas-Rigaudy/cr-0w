## item_236_first_module_split_forge_mrs_optional_and_folder_projects - First module split: forge/MRs optional and folder projects
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 95
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-23 15:45:36

# AI Context
- Summary: Apply the ADR to the first brick: forge/MR features (crow-mrs, review skills, MR hooks, Nest MRs tab, deck MRs tab) turn off with one manifest switch, and folder projects skip git steps.
- Keywords: modules, forge, mrs, optional, folder projects, manifest
- Use when: making the MR brick optional.
- Skip when: other bricks.

# Problem
- Paul works alone in a folder with no forge; MR features are noise for him.

# Scope
- In:
  - Manifest flag `modules.forge`; when off: MR scripts not run by brief/gather, skills not installed, tabs hidden, doctor quiet (generic switch from item_261).
  - `layout: folder` projects (adr_002; `kind` stays the Nest planet type): crow-activity reads file changes, commit/push rules, identity check and worktrees skip; default `repo` when the path holds `.git`.
  - Doctor fails on `forge` set for a folder project.
  - Test both settings in crow-check.
- Out:
  - Other bricks.

# Acceptance criteria
- AC1: With forge off, brief, gather, doctor and Nest run without any forge call or error.
- AC2: A folder project shows activity in brief.
- AC3: With forge on, nothing changes (owner's setup).

# AC Traceability
- req_064_paul_s_feedback_and_a_modular_public_crow-AC -> This backlog slice. Proof: With forge off, brief, gather, doctor and Nest run without any forge call or error.
- AC1 proof: a throwaway hub with `forge: false` and `glab`/`gh` replaced by traps: `crow-gather brief`, `crow-doctor`, `crow-activity` and `crow-mrs json` ran with zero forge calls; `crow-mrs json` returns `off: true` with empty lists, `crow-mrs show` exits 2 with the off message. The Nest (`hiddenTabs`) and the deck (`deck_pages` from the manifests, via `crow-modules json`) hide the MRs tab; tests in `tests/test_folder_projects.py`, `tools/nest/src/ui/tabs.test.ts`, `tools/crow-deck/modules_test.go`.
- AC2 proof: a folder project lists `changed <file>` lines in `crow-activity` and the brief gather; `crow-registry json` serves `layout` and per-day file counts as `work`.
- AC3 proof: the owner's setup (forge on) keeps every tab and forge call; `bin/crow-check` green (823 passed), go test and vitest pass. `crow doctor` reports `layout 12 repo and 3 folder project(s)`: the vaults and godot are now folder projects and show their changed notes in the brief.

# Dependencies
- `adr_002_crow_bricks_module_manifest_and_project_kinds` accepted; `item_261_skills_hooks_and_gather_follow_the_module_switch` done.

# Links
- Request: req_064_paul_s_feedback_and_a_modular_public_crow.
- Architecture decision(s): `adr_002_crow_bricks_module_manifest_and_project_kinds` (migration step 3).
- Next: `item_237_nest_and_crow_core_as_separate_front_and_back_bricks`

# Priority
- Priority: Medium
- Rationale: Smallest split with the most value for a solo user.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-23T15:45:00+01:00 | hub | hub-item236 | forge off and folder projects | succeeded |  | `layout` key, `project_forge`, `changed_files`; doctor `layout` check; dispatch refuses --cwd/--isolate on a folder project; schemas: registry `layout`, mrs `off` |
