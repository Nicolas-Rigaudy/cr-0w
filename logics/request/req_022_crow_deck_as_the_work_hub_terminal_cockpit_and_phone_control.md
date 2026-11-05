## req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control - Crow deck as the work hub: terminal cockpit and phone control
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Make crow deck the place the owner works from, not a side dashboard: a terminal cockpit (deck pane plus the active session pane in one tmux/kitty layout, with a send box and inline approvals), then phone control through SSH into the same cockpit, ntfy notifications, and a small PWA over `crow serve` only for touch actions.
- Keywords: deck, hub, cockpit, tmux, layout, session, send, approvals, phone, ssh, termux, pwa, ntfy, crow serve
- Use when: planning how the owner drives sessions and crow actions from the deck or from the phone.
- Skip when: deck widget visuals (req_020) or the multi-machine sync design (req_016).

# Needs
- Today the owner works in Claude Code sessions in terminals and keeps crow deck beside them as an information dashboard. The goal (owner, 2026-11-05) is to work from the deck itself.
- Direction decided 2026-11-05: **A first**, the terminal cockpit. No desktop web UI for now.
  - One layout: the deck in a fixed pane, the active session in the pane next to it. Picking a session, MR, run or thread in the deck opens or focuses its session in that pane instead of a new window.
  - A send box in the deck that types into the focused session (`crow-focus --send` exists).
  - A session switcher (live sessions across accounts, with state: working, waiting for input, waiting for approval).
  - Inline approvals: a session waiting on a permission or an AskUserQuestion shows up in the deck, and the owner can jump to it or answer it from there.
- Phone control (owner question 2026-11-05: web UI, Android app or SSH?). Decided order:
  1. SSH into the cockpit: Termux or Termius over Tailscale, `tmux attach` to the same layout. Free once the layout exists; the deck is vertical-first already.
  2. Notifications: ntfy (self-hosted on the home server, req_016) when a run finishes or a session waits on the owner.
  3. A PWA over `crow serve` only for actions that are painful over SSH (approve, dispatch, open an MR, read a run report). Same API as req_016.
  - No native Android app: PWA plus ntfy covers notifications. Chatting with a single Claude session from the phone already works with Claude Code Remote Control.

# Context
- `bin/crow_term.py` (item_031) already has kitty, tmux and print backends; `crow term`, `crow-focus` and `crow-relaunch` go through it.
- The deck already runs actions: handoff (main.go), mail session (today_data.go), MR review and posting, runs, usage.
- req_016 plans Tailscale, `crow serve` (HTTP API over existing commands), the home server as aggregator, ntfy and a PWA. This request uses those pieces and does not redesign them.
- req_020 owns the deck's look and composable widgets; new deck surfaces here follow it.
- Open spike: a Claude Code panel mod (plugin-authoring) could show deck widgets beside the chat. Check whether it makes the split layout unnecessary before building the layout.

# Acceptance criteria
- AC1: One command (`crow deck --cockpit` or similar) opens the layout in tmux or kitty: deck pane plus a session pane; selecting a session in the deck shows it in the session pane.
- AC2: From the deck, the owner can send a prompt to the focused session and see sessions waiting on input or approval, without leaving the deck.
- AC3: From the phone over Tailscale, `ssh` plus `tmux attach` reaches the same cockpit; a runbook covers the phone setup (Termux or Termius, keys, tmux binding).
- AC4: A finished run or a session waiting on the owner sends an ntfy notification to the phone.
- AC5: A list of the touch actions the PWA must cover exists, derived from real phone use over SSH, and is linked to req_016 before any PWA work starts.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics_manager/flow.py`
- `logics_manager/assist.py`
- `tests/python/test_logics_manager_cli.py`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

