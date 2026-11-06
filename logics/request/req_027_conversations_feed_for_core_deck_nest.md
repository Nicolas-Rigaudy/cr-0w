## req_027_conversations_feed_for_core_deck_nest - Conversations as the assistant's live feed for core, deck and nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 75%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Make `cdx conversation` (and `cdx sessions` / `cdx history --interactive`) the one source core uses to read what every session said, on every account. deck and nest show live conversation views from it, and the assistant layer (briefs, waiting-on-owner detection, digests, the personal butler) reads conversations instead of parsing transcripts or scraping panes.
- Keywords: cdx conversation, transcript, cursor, tail, sessions, history, core, deck, nest, comlink, assistant, butler, digest
- Use when: building core's session read API, deck/nest session views, or any feature that needs what a session said.
- Skip when: account routing (req_025); sending prompts (stays `crow-focus --send`).

# Needs
- See any session's conversation (any account, any project) from deck and nest, live, without attaching to its pane.
- One read path: core polls `cdx conversation <account> --cursor` per live session and serves normalized entries (role, text, tool_name, tool_summary, at, id). No crow code parses Claude .jsonl itself.
- Build the assistant on top of it:
  - waiting-on-owner: the last entry is an assistant question or a tool waiting for approval → comlink notification with the excerpt.
  - session digest: what a session did since a cursor (for /brief, /eod, /wrap, retro, weekly) without re-reading the transcript.
  - cross-session search and "what did I decide about X" over recent conversations.
  - personal butler (req_026): the rk1 conversations are its memory of requests and follow-ups.

# Context
- cdx 0.26 (cdx-manager README "Read an interactive conversation", docs/cli.md "Find interactive conversations"):
  - `cdx conversation <name> [--conversation ID] [--cwd PATH] [--since ISO] [--cursor C | --tail N | --before C] [--limit N] --json`: read-only pages. Envelope: status, reason, provider, session, entries, next_cursor, has_more, before_cursor, resolution_method (conversation_id | recency | null).
  - Entries carry role, text, tool_name, tool_summary, at, id; tool args/results and reasoning are omitted; common secrets are redacted (not all sensitive prose).
  - `--tail N` starts near the present; `next_cursor` polls for new entries; `before_cursor` loads older pages; `reset_required` on rewrite.
  - Without an ID it falls back to the latest transcript for `--cwd` (recency, unverified); `--since` narrows that fallback.
  - `cdx sessions [--all] --json` gives conversation_id, profile_dir, transcripts_dir per live session; `cdx history --interactive --cwd P --json` finds past interactive launches with conversation_id and first_user_prompt.
- req_022 currently plans core's read as `tmux capture-pane` or "the transcript tail from the session .jsonl". This request replaces the .jsonl part with `cdx conversation`; capture-pane stays for the raw screen (approval prompts are not in the transcript).
- Existing readers to migrate: crow-retro (transcript digests), crow-ctx (context size), deck session list, review-threads/afk reports where they read transcripts. Check each before replacing; keep what cdx does not return (token counts, tool args).

# Acceptance criteria
- AC1: core exposes `GET /sessions/<id>/conversation?tail=N|cursor=C|before=C` backed by `cdx conversation --json`, for live sessions on every cdx account, resolved by exact conversation_id from `cdx sessions --json` (recency fallback flagged in the response).
- AC2: core keeps one cursor per live session and polls; new entries reach deck and nest by push (SSE or websocket), no client polls cdx.
- AC3: deck shows a conversation pane for the selected session (tail, scroll back with before_cursor), next to or instead of the pane capture.
- AC4: nest (phone and laptop) shows the same conversation view with a send box; the look follows req_022 AC7.
- AC5: waiting-on-owner detection from conversation entries feeds comlink; a test covers question, approval wait and finished turn.
- AC6: a `crow conversations digest <session> --since C` command returns what changed since a cursor; /brief and /wrap use it.
- AC7: no crow code parses provider transcripts for content that `cdx conversation` returns; remaining direct readers are listed with the reason.
- AC8: conversation text never leaves the machine or tailnet without the owner's action (redaction is partial); nest over Tailscale only.

# Out of scope
- Writing to conversations (send stays crow-focus/tmux), changes in cdx-manager (request them through `crow handoff cdx-manager` if a gap shows up, e.g. token counts or approval state in entries).

# Risks
- Recency fallback can show the wrong conversation: always prefer conversation_id; mark recency in the UI.
- Polling many sessions: one poller in core with cursors, bounded limit, back-off for idle sessions.
- Sensitive text on the phone: Tailscale only, no public exposure.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics/request/req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control.md` (core, deck, nest, comlink)
- `logics/request/req_016_crow_everywhere_multi_machine_and_mobile_control.md`
- `logics/request/req_026_personal_butler.md`
- cdx-manager README.md, docs/cli.md

# Backlog
- (none yet: split after req_022's core item is groomed)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
