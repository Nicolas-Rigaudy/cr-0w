## req_032_spike_butler_engine_and_persona - Spike butler engine and persona
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85
> Confidence: 75
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 13:07:27

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
6. **Surfaces.** nest chat panel, deck pane, voice (local speech to text and text to speech), a Discord bot (free-text chat, personal data only; work data never goes through it), Android assistant entry (req_033). Same conversation across surfaces with a surface tag.
7. **Memory.** Transcript mirror and index; no hidden memory beyond owner prefs, Logics and the store.

# Context
- Owner decisions 2026-11-06: one butler, scope-aware, focus mode, full persona, all listed surfaces.
- Owner decisions 2026-11-07: see "Decisions" below.
- Ultron dropped Discord and Telegram; Jarvis kept Discord only. One bot channel is enough.
- Owner 2026-11-06: Discord chosen (no Telegram account). The butler reaches core through the core MCP server (adr_001). Built after core.

# Spike findings (2026-11-07)

Spec draft: `logics/specs/req_032_butler/` (engine flags, tool list, context, scope, surfaces, memory, persona skeleton).

## Sources checked
- Claude Code headless docs, [code.claude.com/docs/en/headless](https://code.claude.com/docs/en/headless), checked 2026-11-07 against Claude Code 2.1.319: stream-json in and out, `--include-partial-messages`, `--resume <id>`, `--permission-prompts none` (2.1.259+), `-p` starts in Manual mode when nothing sets a mode (2.1.302+), SIGTERM leaves the running turn unfinished, `--bare` skips OAuth and needs an API key.
- Agent SDK overview, [code.claude.com/docs/en/agent-sdk/overview](https://code.claude.com/docs/en/agent-sdk/overview), checked 2026-11-07: the SDK runs the Claude Code binary; "Anthropic does not allow third party developers to offer claude.ai login or rate limits for their products, including agents built on the Claude Agent SDK".
- Agent SDK sessions, [code.claude.com/docs/en/agent-sdk/sessions](https://code.claude.com/docs/en/agent-sdk/sessions), checked 2026-11-07: `ClaudeSDKClient` (Python) for multi-turn, `resume`, `fork_session`, `SessionStore` adapter to mirror transcripts; sessions stored as `.jsonl` per project folder.
- Agent SDK cost tracking, [code.claude.com/docs/en/agent-sdk/cost-tracking](https://code.claude.com/docs/en/agent-sdk/cost-tracking), checked 2026-11-07: caching is automatic; `cache_creation_input_tokens` and `cache_read_input_tokens` per request; 5-minute TTL with an API key, 1-hour TTL on a Claude subscription within included usage (5 minutes once on usage credits unless `promptCacheTtl` is `1h`); `total_cost_usd` is a client-side estimate; resumed sessions report the session total (2.1.277+).
- Python package [pypi.org/project/claude-agent-sdk](https://pypi.org/project/claude-agent-sdk/): 0.2.179, released 2026-11-06, bundles the CLI, Python 3.10+.
- Claude Code memory docs, [code.claude.com/docs/en/memory](https://code.claude.com/docs/en/memory), checked 2026-11-07: `autoMemoryEnabled: false` or `CLAUDE_CODE_DISABLE_AUTO_MEMORY=1`; `claudeMdExcludes` skips ancestor `CLAUDE.md` files.
- Hermes Agent, [github.com/NousResearch/hermes-agent](https://github.com/NousResearch/hermes-agent), page read 2026-11-07 (no tagged release shown): MIT, one gateway process for Telegram, Discord, Slack, WhatsApp, Signal, email; MCP client; built-in cron; agent-curated persistent memory and self-written skills; providers named are Nous Portal, OpenRouter, OpenAI and custom endpoints.
- Subscription OAuth outside Claude Code, Feb 2026 coverage: [The Register, 2026-02-20](https://www.theregister.com/2026/02/20/anthropic_clarifies_ban_third_party_claude_access), [alternativeto.net](https://alternativeto.net/news/2026/2/anthropic-officially-bans-using-subscription-authentication-for-third-party-claude-use). Subscription tokens are for Claude Code and claude.ai; other tools need an API key. Secondary sources quote an Anthropic engineer encouraging local, personal use of the Agent SDK and `claude -p`. The primary legal page was not fetched in this spike.

## What crow does today
- No butler, no core, no MCP server, no transcript index. `crow brief`, `crow daily` and `crow weekly` are one-shot runs; `crow dispatch` starts tracked background runs; `crow-accounts pick` chooses an allowed account; `crow-digest` condenses a transcript; `crow-usage` already calls `claude -p --model <haiku> --no-session-persistence --tools '' --strict-mcp-config` for cheap titles.
- The butler therefore needs core's chat bridge, core's MCP server (adr_001) and the transcript index before it can run for real.

## Engine comparison (Need 1)

| Criterion | A. Long-lived `claude -p` stream-json + core MCP tools | B. Gateway (Hermes style) with MCP tools | C. Agent SDK inside core |
|---|---|---|---|
| Who pays | The Claude subscription of the process (work or personal account), inside its 5-hour and weekly windows | An API key billed per token (subscription OAuth is not allowed outside Claude Code), or a non-Claude provider. Work turns cannot use a non-Claude provider | Same binary as A. With a subscription login it sits in the grey zone of the Feb 2026 rule; with an API key, billed per token |
| Cost per turn | Probe: Haiku, 1 turn, 1 tool call, about 4.2k tokens written to cache, 4.0k read, 67 out, estimated 0.0009 USD (counts against the window, not billed) | Same tokens priced at API list price, plus the gateway's own memory and persona hooks per call | Same as A |
| Cache use | Automatic. Probe shows 1-hour TTL writes on the subscription and a full prefix read on the second request of the turn | API key: 5-minute TTL by default, 1-hour at a higher write price | Same as A; `promptCacheTtl` setting available |
| Restart and memory loss | Conversation on disk; `--resume <id>` restores it. The prefix stays cached if the restart is within the TTL and the prompt files are unchanged. A SIGTERM mid-turn drops that turn | Gateway keeps its own sessions and memory store (a second memory to audit) | Same as A plus `SessionStore` to mirror transcripts into core |
| Tool latency | Probe: MCP stdio call answered in about 30 ms; process start to init 0.8 s; first text 3.0 s with one tool call; turn done 3.4 s | One more hop (gateway to MCP server); unmeasured | In-process SDK MCP tools skip the stdio hop; model time dominates either way |
| Multi-surface streaming | core bridge reads the stdout lines and fans them out (SSE to nest and deck, a Discord bot as one more client) | Built in for chat apps (Discord); nest and deck need its web API or a bridge anyway | Typed messages in core, same fan-out as A |
| Fit with crow rules | Hard limits in flags: `--tools ""`, `--strict-mcp-config`, `dontAsk`; tools only from core | Brings auto-memory and self-written skills that Need 7 forbids; must be turned off | Adds `canUseTool` and hooks in core code; needs core written in Python or TypeScript |

Notes:
- B adds a third-party runtime and an API bill for every Claude turn, and duplicates core's scheduler and memory. Its one strength, a ready Discord gateway, is a small bot process in A or C.
- C is A with a typed wrapper. Moving from A to C later changes the bridge only: the stream-json protocol, the MCP tools and the persona file stay.
- AC1 (10 scripted questions, measured) needs core's MCP server; this spike ran the single allowed probe. The script is ready below.

## Probe (2026-11-07)
- One stream-json turn on a personal subscription account, Haiku, system prompt replaced, built-in tools off, one stub MCP server with a `brief` tool, permission mode `dontAsk`.
- Result: correct one-sentence answer from the stub rows; `system/init` lists the MCP server as connected; usage and timings as in the table. The 4k prefix includes ancestor `CLAUDE.md` files of the probe folder; a butler folder outside every repo, or `claudeMdExcludes`, removes them.

## AC1 question script (to run once core MCP exists)
1. What waits on me?
2. What finished overnight?
3. Is any run stuck or out of quota?
4. Which account has room for a long run?
5. Summarise the open review threads on my personal repos.
6. What did I decide about push notifications? (transcript search)
7. Show the plan of the newest research mission. (mission_get)
8. Tell the parked design mission to use the second option. (tell, comment)
9. Start a research mission on local speech-to-text engines. (dispatch, preview then confirm)
10. What is on the roadmap for the next milestone? (logics_context)

Measure per question: wall time to first text, total time, tool calls, tokens (input, cache write, cache read, output) and the window share used. Run on A first, then C if the owner keeps it as an option.

## Tools, context, scope, surfaces, memory (Needs 2 to 4, 6, 7)
See `logics/specs/req_032_butler/README.md`. Summary: 11 core MCP tools in three exposed classes (read, write-local, costly with a single-use confirm token); outward actions never exposed; evidence rows capped at 20; surface tag on every message; core filters scope on every tool result; private personal data never reaches the model; Discord always personal; transcripts mirrored and indexed in core with auto-memory off.

## Persona (Need 5)
Skeleton and three candidates written as sample replies to the same three questions: `logics/specs/req_032_butler/persona.template.md`. Owner choice (AC3).

## Decisions (owner, 2026-11-07)
1. **Engine.** A long-lived `claude -p` process in stream-json mode now, behind a thin bridge. The Agent SDK later.
2. **Process routing.** Split by data scope. `work` and `both` turns run on a work-allowed account. `personal` turns and Discord run on the personal account.
3. **Butler model.** Sonnet everywhere.
4. **Persona.** A Steward base with a droid flavour: formal and polite. The butler is CR-0W, the droid the project is named after; it addresses the owner as "Sir" (owner, 2026-11-09). Written up in item_118.
5. **Persona at work.** Neutral while work focus is on.

Still open (owner): the persona name and the address form.

# Acceptance criteria
- AC1: An engine comparison with measured cost and latency for 10 scripted questions.
- AC2: A tool list with input, output and preview rules.
- AC3: A persona file draft the owner approves (Steward base with a droid flavour, item_118; name and address form still open).
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
- `item_118_butler_persona_file_steward_with_a_droid_flavour`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
