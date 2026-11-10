## item_184_crow_helps_manage_game_mods_on_the_tower - crow helps manage game mods on the tower
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hosts
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: crow, helps, manage, game, mods, tower
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- The owner wants help with game mods on Windows. Games and mod managers live on the Windows side.

# Scope
- In:
  - Discover first: which games are modded, which mod managers (Vortex, Mod Organizer 2, r2modman, manual), where mods and profiles live.
  - Read access from WSL through /mnt/c (or the game drive) and powershell.exe for Windows-only commands.
  - Backup of a game's mod folder or profile before any change, and a rollback.
  - A tower runbook section, and a script in bin/ once a step repeats.
- Out:
  - Store inventories, installs, cleanup, save backups: only if a need shows.

# Acceptance criteria
- AC1: crow lists one game's installed mods and load order.
- AC2: One mod change is backed up first and rolls back cleanly.

# AC Traceability
- request-AC2 -> This backlog slice. Proof: AC1: crow lists one game's installed mods and load order.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_008_crow_on_the_tower`
- Architecture decision(s): (none yet)
- Request: `req_056_crow_on_the_windows_tower`
- Primary task(s): `task_022_orchestrate_crow_on_the_windows_tower`

# Priority
- Priority: Medium
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

