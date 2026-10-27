## item_011_background_runs_from_cdx - Count background runs from the cdx CLI (#2)
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Small
> Theme: Hub portability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: crow-runs counts background runs from `cdx runs --json` (no ended_at = running); /proc lists interactive sessions only and matches claude and codex, `-p` and `--print` (#2).
- Keywords: count, background, runs, cdx, cli
- Use when: working on this slice of req_005.
- Skip when: the slice is done.

# Problem
- crow-runs counts background runs from `cdx runs --json` (no ended_at = running); /proc lists interactive sessions only and matches claude and codex, `-p` and `--print` (#2).

# Scope
- In:
  - this slice only.
- Out:
  - other slices of req_005.

# Acceptance criteria
- AC1: the run cap counts running cdx runs of any provider.
- AC2: /proc matching covers `claude`, `codex`, `-p`, `--print`.
- AC3: tests use a fake cdx and a fake /proc.

# AC Traceability
- request-AC2 -> this slice. Proof: its acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_005_codex_support_and_cleanups`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: order in req_005 (profiles helper first).

# Notes
- Generated locally by logics-manager.
