# Plan v1: Nest desk bugs

Three items, drafts in `outbox/items.md`. Order A, B, C (all touch `screens/Desk.tsx`).

Checked: `screens/Desk.tsx`, `screens/Chat.tsx`, `App.tsx`, `router.tsx`, `ui/Shell.tsx`, `ui/ButlerRail.tsx`, `ui/ButlerSheet.tsx` (header), `ui/useMedia.ts`, `data/useScope.ts`, `data/work.ts`, storage use across `tools/nest/src`, item_210 and item_246.

- A (bug 1): Desk unmounts on every tab change and its pick is plain state; the default puts `runs[0]` first (`Desk.tsx:108`). Fix: remembered pick in localStorage, fallback to the first session of the work feed (already most urgent first), never a run. item_210 planned the same pick and was never built.
- B (bug 2): four drafts live in component state. Crossing 1200px swaps Desk layouts and rail/sheet, and the rail remounts on Home. Fix: one `useDraft` store in sessionStorage keyed per target, the `useScope` pattern.
- C (bug 3): Desk has its own CR-0W ask box and log (Ask CR-0W on Fold inner, Launch on laptop, AskBox on the cover) next to the butler side panel, and its free text already goes to the butler. Fix: drop the Desk copy; where action results go is q1.

Not done: the mission worktree `cr-0w.wt/20261125-plan-nest-desk-bugs` does not exist, so the items are not created in `logics/` and nothing is committed. They get created once the answers are in.

```crow
contract: 1
mission: 20261125-plan-nest-desk-bugs
kind: plan
phase: plan
version: 1
summary: Three Nest desk items drafted (remembered pick, drafts store, single CR-0W chat); three choices wait on the owner
confidence: 75
verdict: needs-input
questions:
  - id: q1
    text: "Bug 3: Desk shows its own CR-0W box (Ask CR-0W panel on Fold, Launch panel on laptop, ask box on phone) next to the butler side panel. After removing it, where do Desk action results go (Pause, Stop, Allow, Tell it sent or not delivered)?"
    blocking: true
    choices:
      - {id: a, label: "Into the butler conversation as note lines; /dispatch and /session forms stay on Desk as a plain form panel (recommended)"}
      - {id: b, label: "A short toast on Desk; forms stay on Desk"}
      - {id: c, label: "Keep a small Desk log under the forms, with no input box"}
  - id: q2
    text: "Bug 1: when the last session you picked is a background run, should Desk reopen it?"
    blocking: true
    choices:
      - {id: a, label: "Yes, an explicit pick is restored; only the automatic default skips runs (recommended)"}
      - {id: b, label: "No, Desk never opens a background run on its own, even one you picked"}
  - id: q3
    text: "item_210 (Ready, never built) planned the same last-pick fix plus a layout that item_246 has since shipped. What happens to it?"
    blocking: true
    choices:
      - {id: a, label: "Close item_210 as superseded by item A and item_246 (recommended)"}
      - {id: b, label: "Reuse item_210 for bug 1 instead of a new item"}
artifacts:
  - {path: outbox/items.md}
```
