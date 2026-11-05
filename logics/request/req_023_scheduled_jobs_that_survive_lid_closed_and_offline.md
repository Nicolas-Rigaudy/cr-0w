## req_023_scheduled_jobs_that_survive_lid_closed_and_offline - Scheduled jobs that survive lid closed and offline
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Make every crow scheduled job robust to the laptop being suspended (lid closed) or offline, keep the laptop awake only while a crow run is in flight, and decide what moves to an always-on node (the old phone as cyberdeck, or a server).
- Keywords: timer, systemd, suspend, lid, inhibit, offline, catch-up, hotspot, cyberdeck, phone, always-on, termux, tailscale
- Use when: adding or changing any scheduled or background crow job, or planning where crow runs when the laptop is closed.
- Skip when: deck/nest UI work (req_022).

# Needs
- The owner closes the lid to carry the laptop (commute, evenings); the network can be missing at the scheduled time. Raised 2026-11-05 on the 19:00 retro batch.
- **Catch-up contract for every job**: a stored "done up to" point, a slot or window rather than a clock time, idempotent no-op once done, a network probe before any model call, a capped retry count. The retro batch already follows it (hourly ticks, offline no-op, 3 failures per slot). Audit the others against it: update-dev-tools.timer, AFK/night queues, dispatch runs, brief pre-gather, release/MR watchers.
- **Lid closed while a run is in flight**: wrap crow runs (dispatch, retro batch, AFK jobs) in `systemd-inhibit --what=handle-lid-switch:sleep` so a closing lid waits for the run to finish, then suspends. No always-on lid override: a running laptop in a bag overheats and drains.
  - Optional per-session switch (`crow keepawake on|off`) for a trip where the owner wants work to continue, with a battery floor and a temperature guard that release the inhibit.
- **Network on the move**: the laptop auto-joins the main phone's hotspot (saved Wi-Fi, priority below home). The probe treats a metered hotspot as online; heavy jobs (big digests, updates) may wait for unmetered.
- **Always-on node (owner idea, 2026-11-05)**: the old Note 20 Ultra as a portable control centre (cyberdeck), on the main phone's hotspot or home Wi-Fi, no extra SIM.
  - Role first: comlink and nest client for req_022 (ntfy, PWA, SSH/mosh to the laptop over Tailscale).
  - Role second, to evaluate: run light scheduled jobs itself in Termux (mail/Jira/calendar gathers, ntfy pushes, a queue the laptop drains on wake). Jobs that read local transcripts or repos (retro, git activity) stay on the laptop.
  - Check: Termux:boot + wake lock survival under Android battery optimisation, power draw on charger, whether a cdx/Claude CLI runs there at all.

# Context
- Retro catch-up shipped in `bin/crow-retro` (slot + hourly ticks), runbook `runbooks/session-retros.md`.
- systemd user timers with `Persistent=true` fire a tick missed while suspended on wake (seen: update-dev-tools due 00:03 ran 10:26).
- logind default here: `HandleLidSwitch=suspend`.

# Acceptance criteria
- AC1: a runbook states the catch-up contract and every scheduled crow job is listed against it (compliant or with a backlog slice).
- AC2: closing the lid during a crow run lets the run finish before suspend; tested once by hand.
- AC3: a written decision on the phone node's role (control only, or light jobs too) after a Termux spike.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-retro`
- `runbooks/session-retros.md`
- `logics/request/req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

