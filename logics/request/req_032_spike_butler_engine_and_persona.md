## req_032_spike_butler_engine_and_persona - Spike butler engine and persona
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 55
> Confidence: 45
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. Pick how the butler runs and define its persona: one scope-aware butler with a focus mode, reachable from nest, deck, voice, Discord or Telegram and the Android assistant.
- Keywords: butler, chat, persona, focus mode, agent sdk, stream-json, tools, transcript search, voice, discord, telegram
- Use when: designing the butler chat engine or its character.
- Skip when: building a single surface client.

# Needs
1. Engine options: a long-lived `claude -p` stream-json process with crow tools (Ultron), a gateway with MCP tools and persona hooks (Jarvis), the Claude Agent SDK. Compare cost, context size, restarts, tool surface.
2. Tool surface: missions, inbox, runs, summaries, dispatch, transcript search (SQLite FTS over every session). The butler can tell a mission something but never gives a go on its behalf.
3. Scope: one butler aware of work and personal; a focus mode (work only, personal only) hides the other side.
4. Persona: full character. Name, voice, tone rules, sounds, in-character reports; owner prose rules still apply to anything sent outward.
5. Surfaces: nest, deck pane, voice (local speech to text, text to speech), Discord or Telegram bot, Android assistant.

# Context
- Owner decisions 2026-11-06: research the engine first; one butler plus focus mode; full persona; surfaces as listed.
- Work data stays on Claude accounts; private personal data stays local.

# Acceptance criteria
- AC1: An engine pick with a cost estimate per day of use.
- AC2: A persona sheet (name, tone, do and don't, sample lines).
- AC3: A prototype answers "what is running and what waits on me" from the core API in under 10 seconds.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_028_jarvis_layer_missions_inbox_butler_chat_model_router.md`
- `logics/request/req_026_personal_butler.md`
- `logics/request/req_027_conversations_feed_for_core_deck_nest.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
