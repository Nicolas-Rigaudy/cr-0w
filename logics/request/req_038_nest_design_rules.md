## req_038_nest_design_rules - nest design rules
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 75
> Confidence: 65
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Write the nest design rules before the first screen: tokens shared with the deck terminal kit, page block order, breakpoints, touch targets, icons, motion, status colours, phone checks, and the data contract each screen reads.
- Keywords: design tokens, dark theme, breakpoints, touch targets, icons, block order, phone checks, data contract, ETag, sections
- Use when: building or reviewing any nest screen.
- Skip when: picking the framework (req_031).

# Needs
1. **Tokens.** One token set for colour, space, type, radius, motion. Dark first; light optional. Derive from the deck ANSI theme so the terminal and the web look like one product (owner preference: lazytf look). Status colours fixed: ok, running, waiting on owner, waiting on others, failed.
2. **Block order per page.** Asks first, then live state, then history. Every page has one primary action. Long lists page server side.
3. **Layout.** 7 screens and 3 overlays from prod_005 (Home, Missions, Mission detail, Work, Personal, System, Logics; Butler chat, New mission sheet, Command palette). Desktop three columns (nav or list, content, chat). Phone bar: Home, Missions, Work, Personal, More; "More" shows the worst hidden status. Breakpoints named once.
4. **Galaxy Z Fold 6 (the owner's phone).** Cover screen: single column, bottom bar, chat as a full sheet. Inner screen (near square, tablet width): two panes, list and detail or content and chat. Folding and unfolding keeps the open screen, scroll position and any draft answer. Half-folded (Flex mode): content on top, composer or actions below, when the browser exposes the fold (Viewport Segments API; support checked in req_031). No layout keyed on device names: breakpoints and segments only.
5. **Touch and input.** Targets at least 44 px; heavy buttons ask a second tap; keyboard shortcuts and a command palette on desktop; every notification deep links.
6. **Icons and motion.** One SVG icon set, no emoji in chrome. Motion only to show state change; respects reduced motion.
7. **Phone checks.** Each screen checked on the Fold 6 cover screen, the inner screen in both orientations, and across a fold and unfold, with a screenshot pair in the mockup loop before merge (Jarvis lesson 3).
8. **Data contract.** Each screen names the API sections it reads. Sections carry `updated_at` and an ETag; a stale section shows its age. No offline cache: a stale screen never looks live.
9. **Visual choices.** Options shown side by side with `bin/crow-compare` before any choice (owner preference).

# Context
- Jarvis lesson 3: write design rules before the first screen. Ultron used OKLCH tokens, dark only; Jarvis tokens with dark and light and a per-user accent.
- crow deck already has a theme and the terminal kit (req_020); deck UX2 boards are the approved look.
- Owner: design choices are asked with mockups first; information stays visible.
- Owner 2026-11-06: the phone is a Galaxy Z Fold 6; nest must handle both screens well. Screen list cut from 22 to 7 plus 3 overlays.

# Acceptance criteria
- AC1: A tokens file usable by both the deck kit and nest, with a contrast check.
- AC2: A rules page with block order, layout, breakpoints, touch, icons, motion, phone checks.
- AC3: Two reference screens (Home and Mission detail) mocked on the cover screen, the inner screen and desktop, approved through the compare view.
- AC4: Folding and unfolding on the Fold 6 keeps the open screen and a half-written answer.

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
- `logics/request/req_020_terminal_design_system_and_design_review_workflow.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
