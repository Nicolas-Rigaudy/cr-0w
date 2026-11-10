## prod_007_nest_home_design_pass - Nest Home design pass
> Date: 2026-11-10
> Status: Proposed
> Related request: `req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome`
> Related backlog: `item_161_netrunner_palette_for_the_nest_cyber_theme`, `item_162_holographic_hero_on_the_dot_matrix_background`, `item_163_design_skills_audit_of_the_laptop_home`, `item_164_jarvis_chrome_corner_brackets_and_a_segmented_ring_around_the_hero`, `item_165_telemetry_strip_on_the_laptop_home`, `item_166_boot_sequence_on_first_load_of_the_nest`
> Related task: `task_019_orchestrate_the_nest_home_design_pass`
> Related architecture: (none yet)
> Reminder: Update status, linked refs, scope, decisions, success signals, and open questions when you edit this doc.

# Overview
A calmer cyber palette, a hero that reads as a hologram, and JARVIS-style chrome on the approved laptop Home.

# Goals
- One palette shared with the archeotech shell
- A hero that belongs to the background
- Chrome that serves real data, inside the glow budget

# Non-goals
- Other screens (Missions etc. follow task_018)
- Voice mode, memory graph, head tracking (saved for later in the vault note)
- New UI libraries such as Arwes

# Scope and guardrails
- In: scaffolded request, product, backlog, orchestration task, validation, and handoff context.
- Out: unrelated workflow docs and implementation of generated tasks.

# Key product decisions
- Use structured input as the source of truth for generated docs.
- Keep generated write paths local and repo-bounded.

# Success signals
- Generated docs pass lint and audit without broad manual rewrites.
- Context-pack output can be handed to an implementation agent directly.

# References
- Product back-reference: `req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome`
- Task back-reference: `task_019_orchestrate_the_nest_home_design_pass`
