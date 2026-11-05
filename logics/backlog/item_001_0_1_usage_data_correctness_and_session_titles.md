## item_001_0_1_usage_data_correctness_and_session_titles - 0.1 usage data correctness and session titles
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Fix crow-usage pricing of synthetic messages, Haiku session titles cached by id, rtk read across all history databases, active-state agreement, token split and week totals.
- Keywords: usage, data, correctness, session, titles
- Use when: changing bin/crow-usage or the title field in crow-runs json.
- Skip when: layout-only changes.

# Problem
- Costs undercount because one `<synthetic>` message voids a whole session's price.
- Sessions are all titled with the account name, so nothing tells them apart.
- rtk numbers depend on which HOME reads them, and 'active' means different things in the two views.

# Scope
- In:
  - bin/crow-usage pricing, title cache and Haiku titler, rtk multi-database read, token split and week totals.
  - bin/crow-runs json: title field for session agents from the same cache.
  - Fixture updates for tools/crow-watch/testdata and the QML test generators.
- Out:
  - Any layout change beyond showing the new fields.

# Acceptance criteria
- AC1: no session with a priced real model has cost null; synthetic messages are skipped.
- AC2: every session has a non-account title; runs use their dispatch label; generation is cached, async and bounded.
- AC3: rtk totals equal the sum over all history databases.
- AC4: usage active state matches crow-runs for known sessions.
- AC5: per-session and per-day token split plus week totals in the JSON payload.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: no session with a priced real model has cost null; synthetic messages are skipped.
- request-AC2 -> This backlog slice. Proof: AC2: every session has a non-account title; runs use their dispatch label; generation is cached, async and bounded.
- request-AC3 -> This backlog slice. Proof: AC3: rtk totals equal the sum over all history databases.
- request-AC4 -> This backlog slice. Proof: AC4: usage active state matches crow-runs for known sessions.
- request-AC5 -> This backlog slice. Proof: AC5: per-session and per-day token split plus week totals in the JSON payload.

# Validation
- 2026-11-05 live check against `crow-usage --json` (136 sessions) and `crow-runs json`:
  - AC1: no session with a non-synthetic model has `cost_usd` null; `week.unpriced` is 0.
  - AC2: no title equals the account name; 222 Haiku, 70 dispatch, 30 placeholder titles cached; the Haiku titler runs live (`claude -p --model claude-haiku-5-5`) under the 40-call daily cap.
  - AC3: `rtk.databases` is 3 and the totals sum them.
  - AC4: both live agents in `crow-runs` are `active` in crow-usage, and no other session is.
  - AC5: per-session in/out/cache split and `by_day`, per-day split in `days`, `week` totals.
- Fix found during the check: one-prompt `/retro` runs spent a quarter of the daily title budget on near-identical titles. They now take `retro: <reviewed session title>` with no model call (`tests/test_crow_usage_titles.py::test_retro_run_titled_after_its_target_without_model`). Full suite: 273 passed.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_001_crow_observability_surfaces`
- Architecture decision(s): (none yet)
- Request: `req_000_crow_watch_and_bar_modules_ui_ux_pass_session_titles_savings_stats`
- Primary task(s): `task_001_deliver_the_crow_ui_ux_pass_in_two_phases`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.
