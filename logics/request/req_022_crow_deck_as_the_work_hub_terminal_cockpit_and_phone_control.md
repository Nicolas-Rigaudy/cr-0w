## req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control - Crow deck as the work hub: terminal cockpit and phone control
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85%
> Confidence: 80%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 16:13:51

# AI Context
- Summary: Make crow deck the place the owner works from, not a side dashboard: a terminal cockpit (deck pane plus the active session pane in one tmux/kitty layout, with a send box and inline approvals), then one responsive web app over core for the phone (touch session views, ntfy action notifications) and for the laptop's graphical views (charts, diffs, design picks, reports); SSH as fallback.
- Keywords: deck, hub, cockpit, tmux, layout, session, send, approvals, phone, ssh, termux, pwa, ntfy, nest, core, comlink, cyberdeck
- Use when: planning how the owner drives sessions and crow actions from the deck or from the phone.
- Skip when: deck widget visuals (req_020) or the multi-machine sync design (req_016).

# Status (owner, 2026-11-06)
- Superseded in part by prod_005 and req_028. nest is the main surface and deck gets an inbox pane; the terminal-first cockpit (tmux layout, send box, session switcher) is no longer the target. ntfy and comlink give way to Web Push, ntfy only as a fallback. The 0.7 milestone folds into the Jarvis tracks.
- Kept: the Claude Code panel mod spike, the naming (deck, nest, core), SSH plus `tmux attach` as fallback, and the rule that no surface keeps its own logic.

# Needs
- Today the owner works in Claude Code sessions in terminals and keeps crow deck beside them as an information dashboard. The goal (owner, 2026-11-05) is to work from the deck itself.
- Direction decided 2026-11-05: **A first**, the terminal cockpit. No desktop web UI for now.
  - One layout: the deck in a fixed pane, the active session in the pane next to it. Picking a session, MR, run or thread in the deck opens or focuses its session in that pane instead of a new window.
  - A send box in the deck that types into the focused session (`crow-focus --send` exists).
  - A session switcher (live sessions across accounts, with state: working, waiting for input, waiting for approval).
  - Inline approvals: a session waiting on a permission or an AskUserQuestion shows up in the deck, and the owner can jump to it or answer it from there.
- Phone and graphics (discussed 2026-11-05, replaces the earlier SSH-first and PWA-only drafts):
  - **Claude Code Remote Control is not the phone path.** The owner works only in cdx terminal sessions across several accounts; the Claude mobile app is logged into one account and Remote Control is enabled per session, untested with cdx profiles.
  - **Sessions run in tmux** (item_031 backend), so the cockpit, SSH and the phone all see the same session.
  - **Naming (owner, 2026-11-05):** CR-0W is a droid; its parts are named in a cyberdeck / sci-fi register.
    - **deck**: the terminal interface (`crow deck`), the cyberdeck.
    - **nest**: the graphical interface (`crow nest`), on the laptop in the browser and on the phone as the installed app. The crow's nest is the lookout post where everything is visible.
    - **core**: the always-on engine that feeds deck and nest. A systemd user service started at login; `crow core status` / `restart` only when needed. Replaces the earlier `crow serve` and `crow nest` names.
    - **comlink**: the phone notifications (ntfy topic and the feature name).
  - **core** exposes per session: read (last screen via `tmux capture-pane`, or the conversation via `cdx conversation --json`, see req_027) and send (`crow-focus --send`), plus runs, MRs, threads, today, usage. Every UI reads it; no UI keeps its own logic.
  - **Phone step 1: comlink**, ntfy notifications with action buttons (approve, dispatch, open) when a run finishes or a session waits on the owner.
  - ****nest**: one responsive web app (PWA)** over core:
    - Phone: touch views of sessions (pane view, send box, approve buttons), runs, MRs, threads, today.
    - Laptop: the graphical side the TUI renders poorly: usage and cost charts, run timelines, full diffs, design picks with images, the Logics viewer (item_073), AFK and run reports. It replaces the scattered HTML pages (crow-afk-page, crow-night-page, crow-compare, run reports) over time.
  - **Split:** the TUI deck owns keyboard work (sessions, triage, actions); the web app owns graphics and the phone. A deck key opens the selected item in nest (real browser, xdg-open with the real HOME).
  - **One look:** the web app takes its colour tokens from the active archeotech theme, the same source as kitty's ANSI 16 (req_020).
  - **Fallback:** SSH plus `tmux attach` reaches the full cockpit from the phone.
  - **No native app** (desktop or Android) unless the PWA proves short on notifications or widgets. Desktop presence stays in archeotech-shell.

# Context
- `bin/crow_term.py` (item_031) already has kitty, tmux and print backends; `crow term`, `crow-focus` and `crow-relaunch` go through it.
- The deck already runs actions: handoff (main.go), mail session (today_data.go), MR review and posting, runs, usage.
- req_016 plans Tailscale, core (HTTP API over existing commands), the home server as aggregator, ntfy and a PWA. This request uses those pieces and does not redesign them.
- req_027 owns the conversation feed (cdx conversation) that core serves to deck and nest.
- req_020 owns the deck's look and composable widgets; new deck surfaces here follow it.
- Open spike: a Claude Code panel mod (plugin-authoring) could show deck widgets beside the chat. Check whether it makes the split layout unnecessary before building the layout.

# Acceptance criteria
- AC1: One command (`crow deck --cockpit` or similar) opens the layout in tmux or kitty: deck pane plus a session pane; selecting a session in the deck shows it in the session pane.
- AC2: From the deck, the owner can send a prompt to the focused session and see sessions waiting on input or approval, without leaving the deck.
- AC3: core returns each live session's last screen and accepts a prompt for it, for sessions on every cdx account.
- AC4: A finished run or a session waiting on the owner sends an ntfy notification with action buttons; approve and dispatch work from the notification.
- AC5: On the phone, the installed web app lists sessions waiting on the owner, live and finished runs, and open MRs; a session view shows its screen, takes a prompt and approves with one tap.
- AC6: On the laptop, the same web app shows at least usage charts and full MR diffs, and a deck key opens the selected item there.
- AC7: nest's look is shown as options and picked by the owner before it is built. Direction (owner, 2026-11-05): a cyberdeck, a bit sci-fi, CR-0W as a droid with its own character, inspired by protocol droids and cyberpunk AIs, original rather than copied from any film or game. Colours come from the active archeotech theme.
- AC8: SSH plus `tmux attach` reaches the full cockpit from the phone; a runbook covers the phone setup.

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
