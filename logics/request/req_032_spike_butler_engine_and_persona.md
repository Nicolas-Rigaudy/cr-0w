## req_032_spike_butler_engine_and_persona - Spike butler engine and persona
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 65
> Confidence: 55
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. Pick the butler engine and define its tools, context rules, persona and surfaces. One scope-aware butler with a focus mode; full character persona; reachable from nest, the deck pane, voice, one chat bot and the Android assistant.
- Keywords: butler, chat, stream-json, Agent SDK, MCP tools, persona, focus mode, transcript search, tell, brief, voice, Discord, Telegram
- Use when: building the chat, its tools or its persona.
- Skip when: changing mission kinds (req_037).

# Needs
1. **Engine options.**
   - Ultron "Vision": one long-lived `claude -p` process in stream-json mode; tools as HTTP calls to the server (`brief`, chat search, `tell`, dispatch); ~10 line system prompt; no preloaded state.
   - Jarvis: Hermes gateway with MCP tools, persona hooks, web and Discord; preview-then-confirm on costly tools.
   - Claude Agent SDK in the core backend.
   - Criteria: cost per turn, cache use, restart and memory loss, tool latency, streaming to several surfaces, which account pays.
2. **Tools.** Read state (`brief` re-read per question), transcript search (FTS), dispatch, `tell` a mission (comment on a waiting plan, message on a waiting report, resume), Logics context packs, read-only forge. Never `go`, `ok` or mockup approval. Costly or outward tools return a preview first.
3. **Context rules.** Evidence rows not whole docs; long outputs as linked files; shorter answers on the phone; owner prefs loaded.
4. **Scope and focus.** Work data answered only on Claude accounts; focus mode hides the other side in chat and push.
5. **Persona.** Name, voice, tone, sounds, in-character reports. Outward prose still follows owner prefs. One persona file, versioned.
6. **Surfaces.** nest chat panel, deck pane, voice (local speech to text and text to speech), one of Discord or Telegram, Android assistant entry (req_033). Same conversation across surfaces with a surface tag.
7. **Memory.** Transcript mirror and index; no hidden memory beyond owner prefs, Logics and the store.

# Context
- Owner decisions 2026-11-06: one butler, scope-aware, focus mode, full persona, all listed surfaces.
- Ultron dropped Discord and Telegram; Jarvis kept Discord only. One bot channel is enough.

# Acceptance criteria
- AC1: An engine comparison with measured cost and latency for 10 scripted questions.
- AC2: A tool list with input, output and preview rules.
- AC3: A persona file draft the owner approves.
- AC4: A prototype answering "what waits on me" and dispatching a research mission from nest and the bot.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
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
