## req_033_spike_android_extras_for_nest - Spike Android extras for nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 50
> Confidence: 45
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. List what Android adds beyond a PWA for a tinkerer: widgets, tiles, automations, NFC, watch, always-on display, assistant.
- Keywords: android, widgets, quick settings, tasker, kwgt, nfc, wear os, always-on, twa, assistant
- Use when: planning phone features for nest.
- Skip when: building core or nest screens.

# Needs
1. Home screen widgets (inbox count, live missions) through KWGT or a native app.
2. Quick Settings tiles (focus mode, approve last ask, dispatch a preset mission).
3. Tasker or Automate flows over the core API (arrive at the office → work focus).
4. NFC tags (desk tag starts the morning brief).
5. Wear OS: approvals from the watch.
6. Old phone as an always-on cyberdeck display (screensaver from req_036).
7. Android assistant entry to the butler.
8. Which of these need a native or TWA app versus a PWA.

# Context
- Owner decision 2026-11-06: the owner is on Android and a tinkerer; colleagues are on iPhone, so this is new ground.

# Acceptance criteria
- AC1: A table of features with effort, need for a native app, and fun value.
- AC2: A short list of the first three to build.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_031_spike_nest_stack_and_app_shape.md`
- `logics/request/req_015_wishlist_ambient_mobile_voice_colleague_ideas.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
