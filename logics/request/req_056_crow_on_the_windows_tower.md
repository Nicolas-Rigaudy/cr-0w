## req_056_crow_on_the_windows_tower - crow on the Windows tower
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Hosts
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: crow, windows, tower
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Needs
- crow runs on the Windows tower and helps the owner manage game mods there.
- Every tool crow needs works on the tower (WSL Arch, Tailscale, git, gh, claude, cdx, tmux, rtk, logics-manager).

# Context
- Owner 2026-11-10: focus for 2026-11-11 is the tower and the homelab.
- Tower: Windows. crow is bash and Python, cdx sandboxes with bubblewrap, pasta and nftables, Linux only. crow runs in WSL2. Distro: Arch, same as the laptop (owner 2026-11-10).
- Tower network: Tailscale on Windows with WSL mirrored networking (req_016).
- Games: crow helps manage mods. Other game features wait until a need shows (owner 2026-11-10).
- The Yoga homelab is planned in its own local homelab repo (owner 2026-11-10: local, not pushed).

# Acceptance criteria
- AC1: crow doctor passes on the tower and the tower is on the tailnet.
- AC2: From a tower crow session, the owner lists one game's mods and changes one with a backup taken first.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_008_crow_on_the_tower`
- Architecture decision(s): (none yet)

# References
- logics/request/req_016_crow_everywhere_multi_machine_and_mobile_control.md
- logics/specs/req_029_hosts/README.md
- logics/backlog/item_120_tower_as_a_mission_host_for_claude_missions.md

# Backlog
- `item_183_crow_runs_on_the_windows_tower_in_wsl_arch`
- `item_184_crow_helps_manage_game_mods_on_the_tower`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

