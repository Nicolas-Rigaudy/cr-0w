# Butler spec (req_032)

Draft from the 2026-11-07 spike. Nothing here is wired into code. The butler is built after core (roadmap J tracks) and reaches core only through the core MCP server (adr_001).

## Files

| Path | What |
|---|---|
| `README.md` | Engine, tools, context rules, scope and focus, surfaces, memory |
| `persona.template.md` | Persona file skeleton and three candidate personas with sample replies (owner picks) |

## Engine (Need 1)

Recommended shape, pending the owner's decision: one long-lived `claude -p` process per conversation in stream-json mode, with the core MCP server as its only tool source. A small bridge in core owns the process, tags each message with its surface and fans the output stream out to every surface. Details and the comparison sit in the request (`req_032`, Spike findings).

Process flags (Claude Code 2.1.319, checked 2026-11-07):

```
claude -p --input-format stream-json --output-format stream-json --verbose --include-partial-messages \
  --model <butler model> --tools "" --strict-mcp-config --mcp-config <core-mcp.json> \
  --allowedTools "mcp__core__*" --permission-mode dontAsk --permission-prompts none \
  --system-prompt-file <butler-system.md> --append-system-prompt-file <persona.md> \
  [--resume <session_id>]
```

- `--tools ""` removes the built-in file, shell and web tools. The butler sees only core tools.
- `--strict-mcp-config` ignores every other MCP config on the machine.
- `--permission-mode dontAsk` plus `--permission-prompts none` denies anything not listed, without waiting on a prompt.
- The working directory is a dedicated empty folder outside every repo, so no repo `CLAUDE.md`, hooks or `.mcp.json` load.
- The bridge stops the process after 30 minutes idle and resumes it with `--resume <session_id>` on the next message.
- One process per account side (work, personal). core routes each turn by scope (see Scope and focus).

## Tools (Need 2, AC2)

All tools are MCP tools on the core server. Classes:
- **read**: no side effect, answers at once.
- **write-local**: changes core state only (a note on a mission, a draft). Runs at once, logged in the event log.
- **costly**: starts model work or spends quota. Two calls: the first returns a preview and a `confirm_token`; the second, with the token, runs it. The token is single use, bound to the preview hash, and expires after 10 minutes.
- **outward**: anything that leaves the machine. Never exposed to the butler. core publishes after an owner OK from nest.

| Tool | Class | Input | Output | Preview rule |
|---|---|---|---|---|
| `brief` | read | `section?` (inbox, missions, runs, accounts, jobs), `scope?` | Evidence rows, max 20, each with id, one line and a nest link; `updated_at` | none; re-read before any state answer |
| `mission_get` | read | `mission_id` | Status, kind, last ask, last 5 events, links to artifacts | none |
| `search_transcripts` | read | `query`, `scope?`, `since?`, `limit` (default 10) | Matching turns: date, session, surface, 2-line snippet, link | none; tool output is not indexed |
| `logics_context` | read | `ref` (request, item, task id) or `query` | Context pack rows (max 20 rows, 10 docs), never whole docs | none |
| `forge_read` | read | `project`, `what` (mrs, pr, pipeline, threads), `id?` | Rows with state, author, age, link | none; read-only token |
| `usage` | read | none | Room per account (5-hour and weekly), as rows | none |
| `memo_add` | write-local | `text`, `scope` | Memo id | none; append-only journal |
| `tell` | write-local | `mission_id`, `text`, `files?` | Channel used (plan comment, report message, resume) and event id | none for a comment on a parked mission; costly class when it resumes a run (preview names the account and model) |
| `dispatch` | costly | `project`, `kind`, `goal`, `scope` | Preview: kind, gates, account, model, effort, cap; then mission id | always previews |
| `job_run` | costly | `job` (a scheduled core job) | Preview: what it reads, model, account; then run id | always previews |
| `draft_write` | write-local | `target` (mail, ticket, doc), `text` | Draft id in core, link | none; publishing stays an owner action in nest |

Never exposed: `go`, `ok`, `no` on a gate, mockup approval, merge, push, send, post, delete. The core MCP server does not register them for the butler client, so a prompt injection cannot call them.

A message to a mission after its OK reopens it and drops the validation (core rule, adr_001). `tell` is refused while the run is mid-turn; the butler says so and offers to queue it.

## Context rules (Need 3)

- System prompt about 10 lines: role, the tool list in one line each, the rules below. No preloaded state.
- Owner preferences: core serves a condensed digest of `local/owner-prefs.md` (writing and acting rules, under 1,500 tokens) appended once per process start. The full file stays a `logics_context`-style lookup.
- State only from tools, re-read per question. Never answer state from memory of earlier turns.
- Tool outputs are evidence rows (max 20 rows per call). Long outputs (digests, reports, diffs) are files in core with a link; the butler links them.
- Each user message starts with a surface tag: `[nest]`, `[deck]`, `[phone]`, `[cover]`, `[voice]`, `[discord]`. Phone, cover screen and voice answers: 2 sentences or 3 bullets max, links instead of tables. Voice answers: no Markdown, no links, numbers spelled for speech.
- Persona voice in chat only. Anything that becomes an artifact (draft, ticket, MR text, report) uses the neutral voice and the owner's prose rules.
- Context budget: compaction at 150k as for every session; the bridge starts a fresh session after a compaction or after 24 h, seeded with the last 10 turns' summaries from the transcript index.

## Scope and focus (Need 4)

- Each turn carries a scope: `work`, `personal` or `both` (the default outside focus). core sets it from the surface and the focus mode, never the model.
- Work data (work projects, work mail, work forge, work missions) is served only to a butler process on a Claude account allowed for work. A personal-side process gets those rows filtered out by core.
- The Discord surface is always personal: core filters every tool result to personal scope for Discord turns, and refuses work-scope tools there, whatever the focus.
- Focus mode `work` or `personal` is a core setting (nest toggle, deck key, schedule). In focus, core filters tool results and pushes to that side. The butler does not mention hidden items; `brief` returns a count line "N items hidden by focus".
- A cross-scope question in focus gets "that is on the other side, switch focus in nest" with no content.
- Private personal data (bank, mail bodies, health) never reaches the butler model. Tools return a count and a nest link, or a summary a local-model core job already wrote.
- Process routing (owner decision pending): `work` and `both` turns run on the process of a work-allowed Claude account; `personal` turns (personal focus, Discord) run on the personal account's process.
- Scope filtering happens in core on every tool result. The prompt states the rule too, but core is the gate.

## Surfaces (Need 6)

| Surface | Transport | Notes |
|---|---|---|
| nest chat panel | core bridge over SSE (req_031 picks SSE or long-poll); text deltas from `stream_event` lines | Tool calls shown as one-line summaries from `assistant` tool_use blocks |
| deck pane | Same endpoint, polled or SSE from Go | Read and `tell` only in the first version |
| voice | Local speech to text into the nest answer box; local text to speech of the final text | No server voice stack; `[voice]` tag shortens the answer |
| Discord bot | A small bot process on the homelab talks to the core chat endpoint as surface `discord` | Personal scope only, see above; allowlisted owner user id; DM only |
| Android assistant | req_033 | Opens nest chat with `[phone]` |

One conversation per scope side, shared across surfaces. The surface tag rides on each message. The bridge serialises turns: a message sent mid-turn is queued and sent as the next user message (stream-json input accepts it; the bridge does not interleave).

## Memory (Need 7)

- Transcript mirror: core copies each butler session `.jsonl` and every mission transcript into its store, and indexes typed messages and replies (SQLite FTS5). Tool output is excluded.
- No hidden memory: Claude Code auto-memory is off for the butler folder (`autoMemoryEnabled: false` in its settings, or `CLAUDE_CODE_DISABLE_AUTO_MEMORY=1`), and `claudeMdExcludes` skips any ancestor `CLAUDE.md`. Durable facts go to owner prefs (owner edits), Logics or the core store.
- The session id per scope side is kept in core, so a restart resumes the same conversation.
