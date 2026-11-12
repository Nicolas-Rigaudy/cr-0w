## item_193_nest_answers_askuserquestion_with_option_buttons - Nest answers a session's AskUserQuestion with option buttons
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80
> Confidence: 70
> Progress: 0%
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
