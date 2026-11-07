## req_041_crow_data_in_the_vault - Crow data in the vault
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: High
> Theme: Hub data
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: crow, data, vault
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Needs
- Everything the owner reads, edits or pastes from lives in the Obsidian vault folder CR-0W, so the vault is the knowledge base for the owner and for crow.
- Machine state (JSON, JSONL, caches, locks, run queues, handoffs, renders) stays in local/state.
- Crow's activity lines land in the owner's daily note under a CR-0W heading, written by a script with no model call.

# Context
- Agreed design: vault note CR-0W in the Vault - Proposal.md. The owner answered its questions on 2026-10-29; the migration stalled as an open thread with no Logics item.
- vault.crow_dir is null and 8 scripts plus crow-deck treat it as the single state dir, so setting it alone would move run logs and caches into the vault.
- 2026-11-07: research, drafts and weekly rollups already moved to CR-0W/Research, Drafts, Rollups with symlinks back from local/state.
- Vault-side files outside CR-0W (templates, Home pages, vault CLAUDE.md) belong to the vault Claude.

# Acceptance criteria
- AC1: Owner-facing crow files (open threads, status, decisions, rollups, drafts, research, retro reports) live under the vault CR-0W folder and scripts, skills and crow-deck read and write them there.
- AC2: Machine state stays in local/state whatever vault.crow_dir says.
- AC3: Activity lines go to the daily note's CR-0W section through one script; brief, eod, gather and crow-deck read them from there.
- AC4: vault.crow_dir is set and one brief, one eod-style log line and one wrap-style thread edit work end to end.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_006_crow_data_in_the_vault`
- Architecture decision(s): (none yet)

# References
- Obsidian vault note: CR-0W in the Vault - Proposal.md (2026-10-29, owner decisions in section 11)

# Backlog
- `item_127_path_helper_separating_vault_notes_from_machine_state`
- `item_128_activity_lines_in_the_daily_note`
- `item_129_move_owner_notes_into_the_vault_and_switch_on_crow_dir`
- `item_130_vault_side_dashboard_templates_and_folder_conventions`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

