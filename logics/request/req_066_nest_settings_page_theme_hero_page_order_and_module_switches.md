## req_066_nest_settings_page_theme_hero_page_order_and_module_switches - Nest settings page: theme, hero, page order and module switches
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 16:12:19

# AI Context
- Summary: A Settings page in Nest to pick the theme, hero, page order and module switches without editing crow.yaml, then rebuild Nest.
- Keywords: nest, settings, page, theme, hero, order, module, switches, crow.yaml, rebuild
- Use when: building or scoping a Nest UI that writes Nest or module config.
- Skip when: the change is only build-time brick loading (item_267).

# Needs
- The owner changes Nest look and layout from Nest itself: theme, hero, page order and which pages show.
- Later, the same page turns bricks on and off (`crow modules`), so a phone or laptop can manage modules.
- Room for other general settings as they come.

# Context
- item_267 made theme, hero and pages plain config (`modules.nest.{theme,hero,pages}`) with known value lists from `virtual:nest-modules`.
- Brick themes, heroes and pages are compiled in, so a change needs `crow nest build`. Theme alone already switches live (data-theme).
- Core is the only writer: Nest posts to a core endpoint, core edits crow.yaml and runs the rebuild, the same way section actions run scripts.
- Mockups: options A, B and C in `tools/nest/scripts/mockups/settings.html` (render with `node tools/nest/scripts/settings-mockups.mjs <dir>`). Items 282 to 285 stay Draft until the owner picks one.
- Out of scope: installing bricks from git, editing brick config other than nest, any runtime code loading.

# Acceptance criteria
- AC1: A Settings page lists theme, hero and pages with their known values and shows the current crow.yaml values.
- AC2: Saving writes `modules.nest.*` in crow.yaml through core, keeps the rest of the file and its comments, and refuses unknown values.
- AC3: After a save that needs a build, core runs `crow nest build` and Nest reloads on success or shows the build error with the brick name.
- AC4: Module switches (second slice) turn a brick on or off through `crow modules`, show `requires` blocks, and relink skills.
- AC5: Tests and crow-check pass.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed: crow.yaml write that keeps comments, build time on the owner's machine, page layout (mockup first).

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `tools/nest/nest-modules.ts`
- `bin/crow-modules`
- `docs/writing-a-brick.md`
- `logics/backlog/item_267_nest_themes_hero_and_module_pages_built_from_the_user_s_modules.md`

# Backlog
- `item_282_core_settings_endpoint_writes_modules_nest_in_crow_yaml`
- `item_283_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors`
- `item_284_nest_settings_page_for_theme_hero_and_page_order_from_the_picked_board`
- `item_285_module_switches_on_the_nest_settings_page_through_crow_modules`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
