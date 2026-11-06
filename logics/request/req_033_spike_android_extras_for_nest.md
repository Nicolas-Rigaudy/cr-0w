## req_033_spike_android_extras_for_nest - Spike Android extras for nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 65
> Confidence: 55
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. List what Android allows beyond the iOS-bound colleague designs and pick the extras worth building: rich push, home widgets, quick settings tiles, Tasker hooks, NFC tags, Wear OS, an old phone as an always-on display, assistant entry.
- Keywords: Android, PWA, TWA, widgets, quick settings tile, Tasker, NFC, Wear OS, always-on display, assistant, notification actions
- Use when: adding phone features to nest.
- Skip when: building the base PWA (req_031).

# Needs
1. **Push.** Action buttons (go, ok, open), inline reply, grouping, channels per scope, badges. What a PWA gets vs a TWA or native shell.
2. **Widgets.** Inbox count and next ask on the home screen; needs a native or TWA shell.
3. **Quick settings tiles.** Focus mode switch, quiet hours, "what waits on me".
4. **Tasker and intents.** Trigger missions or focus from Tasker profiles (arriving at work, leaving home).
5. **NFC tags.** Tap a tag on the desk to start the day, at home to switch to personal focus.
6. **Wear OS.** Asks on the watch with one-tap answers.
7. **Old phone display.** An always-on cyberdeck screen running the screensaver and the inbox (with req_036).
8. **Assistant entry.** Long-press or voice opens the butler.
9. For each: value, effort, shell needed (PWA, TWA, native), security impact.

# Context
- Owner decision 2026-11-06: explore everything a tinkerer can do on Android.
- Ultron chose Declarative Web Push because iOS lacks more; Android has no such limit.

# Acceptance criteria
- AC1: A table of extras with value, effort, shell and security notes.
- AC2: A recommendation of the first three to build and the shell they need.
- AC3: One working prototype (likely rich push or a tile) on the owner's phone.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
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
