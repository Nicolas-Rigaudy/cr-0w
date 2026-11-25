## item_294_core_snooze_route_for_threads_and_checks_agenda_snooze_button - Core snooze route for threads and checks, Agenda snooze button
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: crow-core serves a snooze for Open Threads entries and crow-checks; Agenda shows SNOOZE on due chips.
- Keywords: core, snooze, route, threads, checks, agenda, button
- Use when: building snooze for Agenda due items.
- Skip when: mail snooze (already served).

# Problem
- Agenda (item_226) lists due threads and checks but cannot snooze them. Core only serves a mail snooze; `crow-checks snooze` is CLI only.

# Scope
- In:
  - a core POST route that snoozes a check or a thread until a date, reusing `crow-checks snooze`.
  - a SNOOZE button on Agenda due chips and rows (1d, 3d, next week).
- Out:
  - mail snooze, other pages.

# Acceptance criteria
- AC1: a snoozed check or thread leaves the Agenda due list until its date, and comes back after it.
- AC2: the route has Go tests; the button has a page test.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: bounded delivery slice.
- request-AC2 -> This backlog slice. Proof: promotable backlog item.
- request-AC3 -> This backlog slice. Proof: delivery chain includes a task-ready backlog item.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Default until groomed.

# Notes
- Owner decision 2026-11-25: later, not in the Agenda build.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

