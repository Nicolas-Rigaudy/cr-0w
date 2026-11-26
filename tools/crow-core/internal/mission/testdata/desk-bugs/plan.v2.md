# Plan v2: Nest desk bugs

## Responses
- a1: "i don't understand the question". Question 1 is rewritten below as q4, with what you see on screen today and what changes.

## What changed since v1
- Bug 1 is already being built on item_210 (Desk opens on the last selected session) in the worktree `cr-0w.wt/item-210-desk-pick`. Its `data/deskPick.ts` does what item A described: remembered pick per device, else the first session of the work list, never a background run by default, and an explicit pick of a run is restored. Item A is dropped; v1 questions q2 and q3 are settled by that build and withdrawn.
- Two items remain, drafts in `outbox/items.md`: B (typed text survives a screen move) and C (one CR-0W chat on Desk). Order: B, then C, after item_210 lands, since all three edit `screens/Desk.tsx`.
- The mission worktree `cr-0w.wt/20261125-plan-nest-desk-bugs` still does not exist, so the items are not created in `logics/` and nothing is committed.

## Question 4, in plain terms
Today on Desk there are two places to talk to CR-0W: the side panel, and a second box on the Desk page itself (called "Ask CR-0W" on a medium window, "Launch" on the laptop). Bug 3 removes the second box.

That second box also shows short replies after you press a button on Desk, for example "Pause my-session: paused", "Stop my-run: stopped", or "sending… / not delivered, Retry" after Tell it. Once the box is gone, those replies need a new place.

```crow
contract: 1
mission: 20261125-plan-nest-desk-bugs
kind: plan
phase: plan
version: 2
summary: Bug 1 is already built on item_210; two items drafted (drafts survive a screen move, one CR-0W chat on Desk); one question rewritten
confidence: 80
verdict: needs-input
questions:
  - id: q4
    text: "When you press Pause, Stop, Allow or Send on Desk, a short reply shows today in the second CR-0W box on the Desk page (e.g. 'Stop my-run: stopped'). Bug 3 removes that box. Where should those replies show instead?"
    blocking: true
    choices:
      - {id: a, label: "In the side panel chat, as a small grey line (recommended)"}
      - {id: b, label: "A pop-up message on Desk that fades after a few seconds"}
      - {id: c, label: "A small list of replies on Desk with no text box"}
responses:
  - {ref: a1, outcome: done, note: "\"i don't understand the question\": rewritten as q4 with an example"}
artifacts:
  - {path: outbox/items.md}
```
