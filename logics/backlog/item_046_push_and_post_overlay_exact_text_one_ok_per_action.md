## item_046_push_and_post_overlay_exact_text_one_ok_per_action - Push and post overlay: exact text, one OK per action
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: One confirm overlay per push, reply or resolve.
- Keywords: push, post, overlay, exact, text, per, action
- Use when: building the deck review flow (req_010).
- Skip when: the work is outside the deck review flow.

# Problem
- One confirm overlay per push, reply or resolve. Part of req_010. Needs item_045. Same overlay serves item_020 actions.

# Scope
- In:
  - The postOverlay becomes real: it shows the exact text and the command, enter runs it, esc cancels. One OK covers one action. Work repos keep the explicit OK rule.
- Out:
  - Replaces the disabled post overlay in overlay.go.

# Acceptance criteria
- AC1: Each action shows its exact text.
- AC2: Enter runs only that action.
- AC3: Replies are in French, code in English.
- AC4: Tested with simulated commands.

# AC Traceability
- req_010 AC1-AC5 -> this slice. Proof: its own acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Default until groomed.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

