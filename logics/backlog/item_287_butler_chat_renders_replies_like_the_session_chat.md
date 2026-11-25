## item_287_butler_chat_renders_replies_like_the_session_chat - Butler chat renders replies like the session chat
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Small
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The butler rail and sheet print CR-0W replies as raw text; the session chat parses Markdown. Share one reply renderer.
- Keywords: nest, butler, chat, markdown, render, rail, sheet
- Use when: changing how butler or session chat messages render.
- Skip when: changing what the butler says (item_289).

# Problem
- Owner 2026-11-25, screenshot `owner screenshot 20261125_141736.png`: the same reply shows as formatted text in the session chat (bold, numbered list, paragraphs) and as one raw block in the butler panel (`**Where Desk action results go**`, list items run together).
- Cause (read in the code):
  - `ui/ButlerRail.tsx:67` `Msg` renders `<p>{m.text}</p>`. No Markdown parsing, so `**` shows and newlines collapse.
  - The typing bubble (`ButlerRail.tsx:176`, `b.partial`) has the same raw `<p>`.
  - `screens/Chat.tsx:32-45` `Spans` and `Reply` parse the text with `partsOf`/`spansOf` from `data/conversation.ts` (paragraphs, list items, headings, code, bold).
  - Labels differ: session chat says `crow` (`Chat.tsx:271`), the butler says `cr-0w` (`ButlerRail.tsx:66`).
- `ButlerChat` is used by the rail, the butler sheet (`ui/ButlerSheet.tsx:65`) and Home (`screens/Home.tsx:150`), so one fix covers all three.

# Scope
- In:
  - Move `Spans` and `Reply` out of `screens/Chat.tsx` into a shared `ui/Reply.tsx`; Chat imports it.
  - `ButlerRail.tsx` `Msg` renders CR-0W replies and the typing partial with `Reply`. Owner messages stay plain text in both.
  - Note lines keep their `code` spans (`noteText`), styled like the session chat `code`.
  - One speaker label for CR-0W in both chats: `cr-0w` (the project spelling).
  - CSS: the rail `.rl-msg` gets paragraph spacing, `.li` list rows, `b` and `code` rules matching `.msg` in the session chat, sized for the narrower rail.
- Out:
  - Bubble layout and widths of each surface (rail stays compact).
  - Message content (item_289).

# Acceptance criteria
- AC1: A reply with bold, a numbered list and two paragraphs renders the same structure in the butler rail, the butler sheet, Home and the session chat.
- AC2: No literal `**` or run-together list items in the butler panel for the reply in the screenshot (seq 60 of the butler conversation).
- AC3: The streaming partial renders with the same parser, without flicker on unclosed `**`.
- AC4: Both chats label CR-0W `cr-0w`.
- AC5: vitest covers the shared `Reply` (bold, list, code, heading) and a `ButlerChat` render; tsc and vitest pass; shots of rail and session chat side by side at 1440x900 in both themes, zoomed for contrast.

# Links
- Request: `req_069_nest_butler_chat_polish_and_failure_clarity`
- Related: item_176 (butler chat), item_289 (Sir sparingly).

# Launch prompt
```text
Read local/owner-prefs.md and local/prefs/design.md first. Work in a worktree under <repo>.wt/, commit with type[SCOPE]: message, no AI attribution, never push to main without the owner OK. Implement item_287 (butler chat renders replies like the session chat) end to end. Read logics/backlog/item_287_butler_chat_renders_replies_like_the_session_chat.md first, then tools/nest/src/screens/Chat.tsx and tools/nest/src/ui/ButlerRail.tsx. Look at the owner screenshot owner screenshot 20261125_141736.png. Extract the Reply renderer into tools/nest/src/ui/Reply.tsx, use it for CR-0W replies and the typing partial in the butler chat, align the CSS and the speaker label. Run tsc and vitest, take side-by-side shots into the vault CR-0W/Research/nest-chat/ and compare before reporting.
```

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
