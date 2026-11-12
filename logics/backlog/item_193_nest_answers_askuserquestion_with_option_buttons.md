## item_193_nest_answers_askuserquestion_with_option_buttons - Nest answers a session's AskUserQuestion with option buttons
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 100
> Confidence: 95
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: a session waiting on AskUserQuestion shows each question with its options as buttons in the Nest chat (laptop, Fold, phone); a tap or a free "Other" answer reaches the session.
- Keywords: AskUserQuestion, nest, chat, options, multiSelect, other, crow-prompt, permission hook
- Use when: answering a session's AskUserQuestion from Nest.
- Skip when: tool permission prompts (item_174), free text send (item_170).

# Problem
- The owner puts most decisions in AskUserQuestion. Nest shows the tool call folded and the session as waiting, with no way to pick an option. Today a typed line from the send box lands as the "Other" answer (seen 2026-11-12 in item_187's live test), which loses the options.

# Scope
- In:
  - Source: `bin/crow-prompt hook` already records AskUserQuestion (`tool: AskUserQuestion`, `detail` = the questions JSON with header, options, descriptions, multiSelect) in `local/state/prompts/<pid>.json`. Parse it into `prompt.questions[]` in `crow-runs json`.
  - Key mapping: spike the terminal keys for single select, multiSelect, several questions (Tab or arrows to the next question, submit step) and "Other" with text. Note them in this doc.
  - Core: extend `POST /api/sessions/<pid>/answer` with `{prompt, answers: [[option index...] | {other: text}]}`; 409 when stale, as item_174.
  - Nest: render options as a board first (owner picks): the docked ask in the chat shows header, question, option label and description, recommended marker, one button per option, multi-question stepper, Other field. Then build in Chat.tsx (all layouts).
- Out:
  - Changing how sessions ask. Prompts outside kitty or tmux.

# Acceptance criteria
- AC1: The owner approves a rendered board of the docked question first.
- AC2: A session waiting on AskUserQuestion shows every question with its options and descriptions in Nest within 10 s.
- AC3: Picking options (single, multiSelect, two questions) answers the session with exactly those options; the session carries on.
- AC4: An "Other" text answer arrives as typed.
- AC5: A stale prompt gets 409 and presses nothing.
- AC6: Live test on a hidden tmux session; tsc, vitest, pytest pass.

# Links
- Request: req_057 (Nest on the laptop); builds on item_174 and item_187.

# Priority
- Priority: High
- Rationale: the owner's main way of taking decisions from sessions; raised 2026-11-12.

# Notes
- Raised by the owner during the item_187 live test.

# Key map (Claude Code 2.1, spiked 2026-11-12)
- Single select: the option digit selects and moves on. With one single-select question it submits (no review step).
- Other: digit N+1 (the "Type something." row), the text, Enter.
- MultiSelect: digits toggle, Down to the Next row, Enter. Other in a multiSelect: digit N+1, Down to its row, the text, Down, Enter.
- Several questions or any multiSelect end on a review step: "1" submits. Escape cancels.

# Outcome
- Board: option A (option cards, one question per step) picked by the owner (AC1). Board and build shots in the vault crow folder, Research/nest-design/item_193.
- crow-prompt records `questions`; `crow-prompt answer --answers JSON` presses the keys. Core `POST /api/sessions/<pid>/answer` takes `{prompt, answers}`; allow on a question prompt and answers on a tool prompt are 409; deny cancels.
- Nest: `AskDock` replaces the permission card when the prompt has questions (Chat.tsx), all layouts and both themes.

# AC proofs
- AC1: owner picked A on 2026-11-12.
- AC2: core `/api/state/dispatch` listed the live session's two questions with options right after the hook fired; Nest renders them from that section (fixture shots cover, laptop, both themes).
- AC3: live haiku session in hidden tmux: `[[1],[0,2]]` answered "Color → Blue; Fruits → Apple, Plum"; a single question `[[1]]` answered "Pet → Dog" with no review step. The session carried on each time.
- AC4: through core: `[{"other":"Medium please"},{"other":"Basil","options":[0]}]` arrived as "Size → Medium please; Toppings → Cheese, Basil".
- AC5: core 409 on a stale prompt id and on allow for a question prompt; nothing was pressed (the dialog stayed up).
- AC6: pytest 622 passed, go test ok, tsc clean, vitest 57 passed.
