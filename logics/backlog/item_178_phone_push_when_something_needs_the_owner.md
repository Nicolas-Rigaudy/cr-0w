## item_178_phone_push_when_something_needs_the_owner - Phone push when something needs the owner
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: web push from core to the installed Nest PWA when an inbox item, permission prompt or approval waits, or a mission ends; tap opens the item.
- Keywords: web push, vapid, service worker, notification, inbox, approval, pwa
- Use when: building notifications to the phone.
- Skip when: in-app live updates (item_146).

# Problem
- Without push, the owner has to keep Nest open to see that a mission waits. req_031 researched web push; nothing sends it.

# Scope
- In:
  - VAPID keys in local state (0600); `POST /api/push/subscribe`; subscriptions stored per device.
  - Core sends on: new inbox item (ask, approval), session permission prompt (item_174), mission done or failed. Payload under 4096 bytes: kind, title, item id, no work text.
  - Service worker shows it and opens Nest on the item. Action buttons only for safe answers (open, snooze); never approve from the notification.
  - Quiet hours and a per-kind mute in Nest settings.
- Out:
  - Native app. SMS or mail fallback.

# Acceptance criteria
- AC1: A mission asking a question buzzes the phone within 15 s; the tap opens its answer sheet.
- AC2: An approval notification has no approve button and no work text in the payload.
- AC3: Quiet hours hold back non-urgent pushes and send one digest after.

# AC Traceability
- request-AC2 -> AC1. Proof: the phone gets a push when a run needs the owner.

# Decision framing
- Product framing: Not needed
- Architecture framing: Consider. Push library (webpush-go) and key storage, decide during the build.

# Links
- Request: `req_055_a_work_day_from_nest_without_a_terminal`, `req_031_spike_nest_stack_and_app_shape`
- Depends on: `item_175_work_approvals_given_from_nest_on_the_exact_text`, `item_174_nest_answers_a_session_s_permission_prompt_with_buttons`
- Primary task(s): `task_021_orchestrate_a_work_day_from_nest`

# Priority
- Priority: High
- Rationale: missions only run unattended if the owner hears when they stop.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
