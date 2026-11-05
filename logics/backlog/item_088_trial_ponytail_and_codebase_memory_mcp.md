## item_088_trial_ponytail_and_codebase_memory_mcp - Trial ponytail and codebase-memory-mcp in one repo session
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
- Summary: Trial two token levers in one repo session each: the ponytail skill (pushes the agent toward less code) and codebase-memory-mcp (local code graph that replaces exploration), measured against a baseline.
- Keywords: trial, ponytail, codebase, memory, mcp, repo, session
- Use when: Running the trial or deciding whether either tool joins a repo's tool set.
- Skip when: Hub sessions: neither tool should load there (req_021).

# Problem
- Repo sessions spend many tokens exploring and writing code. ponytail (MIT, ~3 KB core SKILL.md) claims -45% tokens on its own benchmark; codebase-memory-mcp (MIT, native C binary, 17 MCP tools, no telemetry) indexes code into a local graph.
- Both claims are self-reported. The vault note's star counts look inflated, so measure before adopting.

# Scope
- In:
  - Read both installers before running (codebase-memory installs via curl | bash and may add hooks).
  - Enable each only in one personal repo via its project .claude/settings.json (req_021 AC3), never profile-wide.
  - Compare tokens and result quality on two comparable tasks with and without, using bin/crow-usage and bin/crow-cost-audit.
- Out:
  - Work repos until a personal trial shows a gain.
  - Always-on memory plugins (claude-mem, Trail-Code): rejected 2026-11-05 for start-context cost.

# Acceptance criteria
- AC1: Installer scripts read and any hooks they add are listed before install.
- AC2: Each tool runs in one repo only; hub start context unchanged (crow-cost-audit).
- AC3: A short result table (tokens, turns, diff size, quality notes) decides keep or drop for each tool.

# AC Traceability
- AC1 -> this item. Proof: to record at closeout.
- AC2 -> this item. Proof: to record at closeout.
- AC3 -> this item. Proof: to record at closeout.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_021_load_plugins_and_mcp_per_session_kind (per-repo loading); also feeds req_008 cost pass
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: R&D after req_021 per-repo loading exists.

# Notes
- Sources: https://github.com/DietrichGebert/ponytail, https://github.com/DeusData/codebase-memory-mcp. Owner picked this on 2026-11-05 from the vault link review.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

