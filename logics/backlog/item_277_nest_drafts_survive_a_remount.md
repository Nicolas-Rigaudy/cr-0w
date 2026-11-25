## item_277_nest_drafts_survive_a_remount - Nest drafts survive a remount
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Text typed in the butler chat, the session chat or the Desk boxes is lost when the component remounts (window crosses 1200px, Home and back, tab change). Keep drafts per box in sessionStorage.
- Keywords: nest, draft, chat, remount, sessionStorage
- Use when: draft text in Nest inputs.
- Skip when: chat layout.

# Problem
- Draft text lives in component state and dies on unmount:
  - Session chat: `Chat.tsx:91` `text`.
  - Butler chat: `ButlerRail.tsx:142` `text` in `ButlerChat`.
  - Desk: `Desk.tsx:99-100` `text` (Tell it) and `ask`.
- Remounts:
  - Crossing 1200px swaps `DeskLaptop` and `DeskInner` (`Desk.tsx:864`) and the butler rail and sheet (`App.tsx:59,66`). Moving the window to another screen crosses it.
  - The rail is keyed `home`/`page` (`App.tsx:59`): Home and back remounts it.
  - A tab change unmounts Desk.

# Scope
- In:
  - New `data/drafts.ts`: module map mirrored to sessionStorage (`nest.draft.<key>`), guarded like `data/useScope.ts`. Hook `useDraft(key)`.
  - Keys: `butler`, `chat:p<pid>` / `chat:m<mission>` (existing `tkey`, `Chat.tsx:113`), `desk.tell:<pid>`, `desk.ask`.
  - Replace the four `useState('')` drafts. A successful send clears the key; a failed send keeps it.
- Out:
  - /dispatch and /session forms. Cross-device sync.

# Acceptance criteria
- AC1: Type in the session chat, resize across 1200px both ways: the text stays.
- AC2: Type in the butler chat, cross 1200px and go Home and back: the text stays.
- AC3: Session X draft comes back after visiting Y; Y's box is empty.
- AC4: Send clears the draft; a failed send keeps it.
- AC5: vitest covers `drafts.ts` and a remount of `Chat`; tsc and vitest pass.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`. Same remount as item_210.

# Priority
- Priority: High
- Rationale: lost typing makes the owner fall back to the terminal.

# Notes
- Raised by the owner 2026-11-25. Drafted by plan mission 20261125-plan-nest-desk-bugs.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
