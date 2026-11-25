## item_289_butler_says_sir_sparingly - Butler says Sir sparingly
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
- Summary: The persona tells the butler to address the owner as Sir, and it does so in every reply; Nest UI strings add more. Keep Sir for greetings and rare formal moments.
- Keywords: nest, butler, persona, sir, tone, copy
- Use when: changing the butler voice or Nest copy that addresses the owner.
- Skip when: rendering (item_287).

# Problem
- Owner 2026-11-25: "sir isn't needed on all messages, it gets annoying".
- Cause:
  - `tools/crow-core/internal/butler/persona.md:1` says `Address the owner as "Sir".` with no frequency limit. Every reply in the 2026-11-25 conversation (seq 57 to 71) opens or closes with Sir.
  - Nest UI copy repeats it on the same screen: `format.ts:32` greeting, `hero/model.ts:187,193` hero line, `Home.tsx:71,119-120`, `HomeLaptop` HUD (`ui/HomeHud.tsx:65`), `WorkList.tsx:78` and `Desk.tsx:689` empty states, `data/butler.ts:112` canned note.
- The persona memory `butler-is-cr-0w` (personal mode uses Sir, work focus none) is kept; this item sets frequency.

# Scope
- In:
  - persona.md: use "Sir" at most once per conversation turn group: the first reply of a conversation or after a long pause, a refusal, or bad news. Never in two replies in a row. Short acks carry no address.
  - Nest copy: keep Sir in the time greeting only. Hero lines, empty states and the canned note drop it.
  - Update tests and fixtures that assert Sir (`scope.test.ts`, `work.test.ts`, `hero/model.test.ts`, `fixtures.ts`, core `butler_test.go` strings stay valid).
- Out:
  - Persona voice otherwise (formal, calm, brief).
  - Work focus mode rules (already no address).

# Acceptance criteria
- AC1: In a scripted five-message butler exchange, Sir appears in at most one reply and never in two consecutive ones.
- AC2: Home, Desk and the hero show Sir once on screen at most (the greeting).
- AC3: The butler memory and runbooks that describe the persona say "sparingly".
- AC4: go test, tsc and vitest pass.

# Links
- Request: `req_069_nest_butler_chat_polish_and_failure_clarity`
- Related: item_287 (render), memory `butler-is-cr-0w`.

# Launch prompt
```text
Read local/owner-prefs.md and local/prefs/design.md first. Work in a worktree under <repo>.wt/, commit with type[SCOPE]: message, no AI attribution, never push to main without the owner OK. Implement item_289 (butler says Sir sparingly) end to end. Read logics/backlog/item_289_butler_says_sir_sparingly.md first. Edit tools/crow-core/internal/butler/persona.md for the frequency rule, drop Sir from Nest copy except the greeting (files listed in the item), update tests and fixtures, and update the butler persona memory and any runbook that says "addresses the owner as Sir" to say sparingly. Check with a short live exchange on the butler that Sir shows at most once.
```

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
