## req_031_spike_nest_stack_and_app_shape - Spike nest stack and app shape
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 60
> Confidence: 55
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. Pick the nest stack and app shape: React front end, Go or Python backend, and PWA, web app or phone app.
- Keywords: nest, react, vite, go, python, pwa, twa, web push, sse, etag, design tokens
- Use when: choosing technologies for core API and nest.
- Skip when: building nest screens.

# Needs
1. Front end: React (the work team's stack) with a strong UI ecosystem. Compare with the colleague choices (React 19 + Vite; vanilla ES modules with no build).
2. Backend: Go (shares code and types with deck) or Python (shares `crow_lib.py` and scripts). Compare effort, deploy, live updates (SSE vs polling with ETags), Web Push libraries.
3. App shape: installable PWA, plain web app, or Android app (TWA wrapper or native) for widgets, tiles and assistant hooks (see req_033).
4. Screens to support: dashboard home (inbox, today), chat on every screen, mission detail, mockup loop, live terminal (xterm), groups nav.
5. Design: dark tokens matching the terminal kit (req_020) so deck and nest look like one product.

# Context
- Owner decisions 2026-11-06: React front end likely; Go or Python backend; Web Push; dashboard home with chat everywhere.
- adr_001 fixes one core API; this spike picks how it is built.

# Acceptance criteria
- AC1: A decision table with a pick for front end, backend, live updates and app shape.
- AC2: A thin prototype: one page reading the inbox from core and receiving one Web Push on the Android phone.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control.md`
- `logics/request/req_020_terminal_design_system_and_design_review_workflow.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
