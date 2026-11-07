## item_126_plan_the_workshop_nest_build_from_the_approved_concept_art - Plan the workshop nest build from the approved concept art
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Turn the approved workshop concept art into a build plan for the nest view: one illustrated room with droids moving along the mission loop, driven by real crow data, with click targets that open panels. Planning only.
- Keywords: nest, workshop, build plan, mission loop, droids, panels, three.js, 2.5D, crow runs
- Use when: planning how the nest workshop gets built after item_125.
- Skip when: still choosing the look (item_125).

# Problem
- The owner wants a living board for crow's mission flow (req_040 purpose section). It shows something needs me, moving or stuck, quota and while-away at a glance, and it is mostly a viewer with a few actions. Nothing says yet how the picked workshop gets built.

# Scope
- In:
  - Data map: every station and droid state mapped to a real crow source (`crow-runs`, sessions, usage, decisions, MRs, AFK queue).
  - Render approach: 2D image of a 3D world is enough (no free camera). Compare a pre-rendered room with animated sprite droids against three.js with glTF, and include the phone version.
  - Actions: click a droid to open its mission, click the desk to open decisions, click the counter to open MRs, click the racks to open the away report, plus a dispatch button.
  - Background items (quota bar, retro and others) placed as small details.
  - Backlog items for the build with ACs, dependencies and order, then stop.
- Out:
  - Building anything in this item.

# Acceptance criteria
- AC1: Data map table, one row per visual element and its crow source.
- AC2: Render approach chosen by the owner from at most 3 options, shown side by side.
- AC3: Build backlog items written and linted, in execution order.

# AC Traceability
- request-AC1 -> AC1, AC2, AC3. Proof: this item's notes and the new backlog items.

# Decision framing
- Product framing: Not needed
- Architecture framing: Consider

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_040_nest_deck_two_themes_3d_world`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: R&D. Depends on item_125.

# Notes
- Depends on item_125. Nest stack itself: req_031.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
