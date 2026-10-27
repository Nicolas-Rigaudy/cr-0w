## item_012_linux_docs_and_portable_tests - Document Linux requirements and make tests portable (#3)
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
- Summary: README lists Linux requirements (/proc, kitty, setsid); tests create a temp accounts.yaml; `pwd` replaces /home/<user> (#3).
- Keywords: document, linux, requirements, tests, portable
- Use when: working on this slice of req_005.
- Skip when: the slice is done.

# Problem
- README lists Linux requirements (/proc, kitty, setsid); tests create a temp accounts.yaml; `pwd` replaces /home/<user> (#3).

# Scope
- In:
  - this slice only.
- Out:
  - other slices of req_005.

# Acceptance criteria
- AC1: README section.
- AC2: `pytest` passes with no local/accounts.yaml.
- AC3: no hard-coded /home path in bin/ or tests/.

# AC Traceability
- request-AC3 -> this slice. Proof: its acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_005_codex_support_and_cleanups`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: order in req_005 (profiles helper first).

# Notes
- Generated locally by logics-manager.
