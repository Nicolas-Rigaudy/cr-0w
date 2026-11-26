## item_300_mrs_review_tab_mockup_options - MRs review tab mockup options
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Board options for the Review tab of the Nest MR detail (findings, decisions, apply, post).
- Keywords: mrs, review, tab, mockup, options
- Use when: designing the Nest review tab.
- Skip when: building it (item_301).

# Problem
- No approved board covers findings in Nest. Owner rule: show mockup options before building a UI.

# Scope
- In:
  - 2 or 3 options in tools/nest/scripts/mockups (laptop and phone) built on board C of mrs-v2: a Review tab in the MR detail with run state, finding list, the current finding (file:line, comment, proposed fix), accept/refuse/edit/chat, apply, post one action at a time.
  - Shots side by side in the vault for the owner's pick.
- Out:
  - Code view of the diff (the deck has one; Nest links to the forge line).

# Acceptance criteria
- AC1: Options rendered with real 245-like data, laptop and phone.
- AC2: Owner pick recorded in this item.

# AC Traceability
- req_071 -> this slice. Proof: its acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_071_nest_mrs_review_flow_and_visible_session_starts
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Reviews rank second in the owner's priorities; no-terminal days depend on Nest.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

