## prod_006_crow_data_in_the_vault - Crow data in the vault
> Date: 2026-11-07
> Status: Settled
> Related request: `req_041_crow_data_in_the_vault`
> Related backlog: `item_127_path_helper_separating_vault_notes_from_machine_state`
> Related task: `task_015_build_crow_data_in_the_vault`
> Related architecture: (none yet)
> Reminder: Update status, linked refs, scope, decisions, success signals, and open questions when you edit this doc.
> Indicators reviewed: 2026-11-09 10:46:09

# Overview
Split crow's data in two: owner-facing notes in the vault CR-0W folder, machine state in local/state, with one path helper every script and skill uses.

# Goals
- The owner finds every crow note in Obsidian
- No machine file ever lands in the vault
- Activity logging stays model-free

# Non-goals
- Vault dashboard, Bases and template edits (vault Claude)
- Retention and archiving rules
- Promotion of crow notes into 30 Knowledge

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
- Product back-reference: `item_127_path_helper_separating_vault_notes_from_machine_state`
- Task back-reference: `task_015_build_crow_data_in_the_vault`
