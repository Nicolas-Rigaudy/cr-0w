## item_183_crow_runs_on_the_windows_tower_in_wsl_arch - crow runs on the Windows tower in WSL Arch
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
- Keywords: crow, runs, windows, tower, wsl, arch
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- crow and cdx need Linux. The tower runs Windows.

# Scope
- In:
  - WSL2 with the Arch image (check the official archlinux WSL image first), systemd on, mirrored networking.
  - Tailscale on Windows; the WSL side reachable over the tailnet.
  - Tools: git, gh, glab, python, node, tmux, claude, cdx, rtk, logics-manager, fish if wanted. Check bubblewrap, pasta and nftables work under WSL2 for the cdx sandbox.
  - crow clone and private config (CROW_HOME), git identity per scope, per-machine secrets never synced.
  - Note GPU model and VRAM for the later local AI item.
  - Write a tower runbook (generic part in runbooks/, private part in local/runbooks/).
- Out:
  - Tower as mission host (item_120), Ollama, state sync.

# Acceptance criteria
- AC1: crow doctor passes on the tower.
- AC2: A cdx run starts on the tower with its sandbox.
- AC3: The laptop reaches the tower WSL over ssh on the tailnet.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: crow doctor passes on the tower.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_008_crow_on_the_tower`
- Architecture decision(s): (none yet)
- Request: `req_056_crow_on_the_windows_tower`
- Primary task(s): `task_022_orchestrate_crow_on_the_windows_tower`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

