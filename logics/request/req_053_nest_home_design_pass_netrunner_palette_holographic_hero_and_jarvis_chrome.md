## req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome - Nest Home design pass: Netrunner palette, holographic hero and JARVIS chrome
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: nest, home, design, pass, netrunner, palette, holographic, hero, jarvis, chrome
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Needs
- The cyber theme is too punchy, especially the background. The owner chose variant B Netrunner (archeotech cyberdeck tokens) on 2026-11-10.
- The hero sits in a box on top of the background and hides the dot matrix. It should read as a hologram projected onto the dot-matrix background.
- Borrow from the JARVIS references: corner brackets, a segmented ring with tick rulers, a telemetry strip, a boot sequence.
- Run the installed design skills on Home and fix what they rank highest.

# Context
- Laptop Home hero, rails and gauges are approved and committed (2bd9e5e, c507221, dc02d78).
- The hero canvas reads its colours from CSS vars at mount (src/hero/Hero.tsx).
- Glow budget rule: at most 3 or 4 glowing elements per screen. Cyan is the one light hue. Other hues mean status only.
- Build one item at a time. Show each against the previous render and wait for the owner's OK before the next.
- The modern theme must not regress.

# Acceptance criteria
- AC1: The cyber theme uses the Netrunner tokens. The owner approves the result side by side with B-netrunner.png.
- AC2: The hero has no visible box. The dot matrix shows through and around the head. The owner approves it.
- AC3: Each later slice (audit, chrome, telemetry, boot) is approved on its own render before the next starts.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_007_nest_home_design_pass`
- Architecture decision(s): (none yet)

# References
- Vault: CR-0W/Research/nest-home-design-pass/Nest Home Design Pass.md (palette sheet, reference and skills breakdown)
- tools/nest/src/theme.css
- cyberdeck pack tokens.json in the shell repo (path in local/projects.yaml) (Netrunner register)
- logics/tasks/task_018_nest_app_shell_on_core_fixtures.md

# Backlog
- `item_161_netrunner_palette_for_the_nest_cyber_theme`
- `item_162_holographic_hero_on_the_dot_matrix_background`
- `item_163_design_skills_audit_of_the_laptop_home`
- `item_164_jarvis_chrome_corner_brackets_and_a_segmented_ring_around_the_hero`
- `item_165_telemetry_strip_on_the_laptop_home`
- `item_166_boot_sequence_on_first_load_of_the_nest`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

