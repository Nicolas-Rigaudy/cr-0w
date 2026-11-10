## req_055_a_work_day_from_nest_without_a_terminal - A work day from Nest without a terminal
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 75%
> Complexity: High
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: after req_054, close the gaps that still send the owner to a terminal: butler chat, chat inside a mission or session, phone push, mail and calendar actions, rituals, a core that stays up when the laptop sleeps, and retros run as missions.
- Keywords: nest, no terminal, butler chat, mission chat, web push, mail, calendar, rituals, retro, always on
- Use when: grooming or building anything that lets the owner run a work day from Nest.
- Skip when: deck parity (req_054), Nest visual design passes.

# Needs
- Owner, 2026-11-10: "nearly not have to open a terminal anymore and just talk in nest to the orchestrator butler chat or in missions/sessions/runs themselves, handle my calendar, mails". Also: run retros on crow easily, maybe as missions.
- Owner, 2026-11-10: "basically most runs will be missions". Missions are the main run path. Plain sessions and `crow dispatch` stay reachable but come second in every list and action.
- The owner picked all four gap groups (2026-11-10): butler chat; session and mission chat plus push; mail and calendar; rituals plus always-on core. Retros: retro mission kind.

# Context
- Butler engine decided in req_032 (logics/specs/req_032_butler/README.md, owner 2026-11-07): one long-lived `claude -p` stream-json process per conversation, core MCP server as its only tools, a bridge in core fanning the output to every surface. Persona in item_118 (Draft).
- Conversation feed: req_027, item_148 (run conversation poller with cursors), item_146 (SSE). Web push research: req_031 (VAPID, `userVisibleOnly`, webpush-go, 4096-byte payload, notification actions).
- Mail and calendar: Microsoft Graph through `bin/crow-today` (local/graph.yaml, refresh token from `crow-today login`). Sending mail on work always goes through an approval on the exact text (item_175).
- Rituals today are skills run in a terminal session: brief, daily, eod, wrap, weekly.
- Retros: `bin/crow-retro` runs from `crow-retro.timer`, one batch per weekday 19:00 slot, reports in the vault Retros folder. The 2026-11-06 batch failed 3/3 with cdx `budget_unverifiable` ("Measured usage cannot be priced; USD budget enforcement is unavailable."). Mission kinds live in logics/specs/req_037_run_contract/kinds/.
- Hosts: req_029 spike, item_120 (tower as mission host, Draft).

# Acceptance criteria
- AC1: From Nest on the phone the owner asks the butler to start, steer or check work, and the butler acts through core actions only, with the same confirmations and approvals as a tap.
- AC2: The owner reads and writes inside any running mission or session from Nest, and the phone gets a push when one needs them.
- AC3: The owner reads flagged mail and today's meetings, approves a drafted reply, and accepts or declines a meeting from Nest. Nothing is sent without an approval on the exact text.
- AC4: brief, daily, eod, wrap, weekly and retro start from Nest as missions and their output reads in Nest.
- AC5: Nest and core answer while the laptop sleeps, or the gap is a decided, written limit.
- AC6: A retro is a mission; each proposed crow fix is an inbox item; an approved fix runs as a build mission on cr-0w.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `req_054_nest_sees_and_pilots_everything_deck_does` (prerequisite: sections, actions, approvals)
- `req_027_conversations_feed_for_core_deck_nest`, `req_029_spike_hosts_model_sources_and_local_ai_benchmarks`, `req_031_spike_nest_stack_and_app_shape`, `req_032_spike_butler_engine_and_persona`, `req_037_mission_kinds_and_run_contract`
- `bin/crow-today`, `bin/crow-retro`, `.claude/skills/retro`

# Backlog
- `item_176_butler_chat_in_nest_on_the_core_mcp_tools`
- `item_177_chat_inside_any_session_or_mission_run_from_nest`
- `item_178_phone_push_when_something_needs_the_owner`
- `item_179_mail_and_calendar_actions_from_nest`
- `item_180_rituals_as_mission_kinds_started_from_nest`
- `item_181_core_reachable_when_the_laptop_sleeps`
- `item_182_retro_mission_kind_that_improves_crow`

# Out of scope
- Voice. Teams chat reading. Jira and Confluence editing from Nest (req_034). Agents holding work credentials.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
