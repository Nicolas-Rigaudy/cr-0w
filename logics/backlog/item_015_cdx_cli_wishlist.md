## item_015_cdx_cli_wishlist - List what crow still needs from the cdx CLI (#6)
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Small
> Theme: Hub portability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Write the list of what crow still needs from the cdx CLI as the reply text for #6 (already drafted in the review note).
- Keywords: list, crow, still, needs, cdx, cli
- Use when: working on this slice of req_005.
- Skip when: the slice is done.

# Problem
- Write the list of what crow still needs from the cdx CLI as the reply text for #6 (already drafted in the review note).

# Scope
- In:
  - this slice only.
- Out:
  - other slices of req_005.

# Acceptance criteria
- AC1: list is in the review note and matches `cdx --help` of the installed version.
- AC2: owner approves the exact text before posting.

# AC Traceability
- request-AC6 -> this slice. Proof: its acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_005_codex_support_and_cleanups`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: order in req_005 (profiles helper first).

# Notes
- Generated locally by logics-manager.
- Wishlist text (checked against cdx 0.22.0 `--help`): `runs --running` filter; `sessions --json` for live interactive sessions; `final_message` in `run-report --json`; documented `profile_dir` per account; `can-resume` taking a conversation id. Reply draft: `local/state/drafts/codex-support-fork-review.md`.
