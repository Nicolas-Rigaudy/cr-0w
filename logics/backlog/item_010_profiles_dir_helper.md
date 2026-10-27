## item_010_profiles_dir_helper - One shared cdx profiles helper (#4)
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Small
> Theme: Hub portability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Add `bin/crow_lib.py` with `profiles_dir()` (CDX_PROFILES, then CDX_HOME/profiles, then ~/.cdx/profiles); crow-memory, crow-relaunch, crow-sessions, crow-usage use it (#4).
- Keywords: shared, cdx, profiles, helper
- Use when: working on this slice of req_005.
- Skip when: the slice is done.

# Problem
- Add `bin/crow_lib.py` with `profiles_dir()` (CDX_PROFILES, then CDX_HOME/profiles, then ~/.cdx/profiles); crow-memory, crow-relaunch, crow-sessions, crow-usage use it (#4).

# Scope
- In:
  - this slice only.
- Out:
  - other slices of req_005.

# Acceptance criteria
- AC1: `profiles_dir()` honours the three sources in that order.
- AC2: no script reads CDX_PROFILES or CDX_HOME directly.
- AC3: tests cover each source.

# AC Traceability
- request-AC1 -> this slice. Proof: its acceptance criteria.

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
