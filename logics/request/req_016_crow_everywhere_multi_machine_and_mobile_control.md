## req_016_crow_everywhere_multi_machine_and_mobile_control - crow everywhere: multi-machine and mobile control
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 55
> Confidence: 50
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Run crow on several machines (work laptop, home Windows tower via WSL, a self-built home server) and control all of them from the phone: see runs, launch sessions, run any crow action, chat with sessions. Split out of req_015 track 2.
- Keywords: mobile, remote, multi-machine, wsl, tailscale, home server, pwa, tmux, sync, self-host
- Use when: scoping remote or multi-machine crow.
- Skip when: working on the deck or screensaver.

# Needs
- Same crow on three machines: work laptop (Arch, today's install), home tower (Windows, crow in WSL), home server (old Lenovo Yoga 2 Pro, Linux, always on, also self-hosts other tools).
- From the phone: list runs and sessions on every machine, launch a session or dispatch, run any `crow` action, chat with a live session, get notified when a run needs the owner.

# Decided 2026-11-03
- The home server is the owner's private cloud: crow, the Obsidian vault and other personal data all pass through it (sync hub, crow aggregator, notifications, self-hosted services).
- Sync: Syncthing everywhere, server as the always-on peer with staggered versioning (replaces the git option below). Phone runs Syncthing-Fork with battery optimization off; LiveSync for the phone only if conflicts remain.
- Devices keep working offline when the server is down (Syncthing peers still sync directly).
- Backups are mandatory once everything goes through one laptop SSD: restic (or similar) nightly to the tower disk plus an off-site copy. Scoped in the homelab repo.

# Proposed architecture (to confirm)
- **Private network:** Tailscale tailnet joining the three machines and the phone. Nothing exposed publicly. On the tower, Tailscale on Windows with WSL mirrored networking.
- **Sessions in tmux (item_031):** a tmux terminal backend makes sessions attachable from anywhere and works on WSL and a headless server, where kitty is absent. Chat = send keys / read pane through crow, or Claude Remote Control for a single session.
- **`crow serve` on each machine:** small HTTP API over the tailnet wrapping existing commands (`crow json`, runs, dispatch, handoff, focus/send). Read only first, actions second.
- **Home server as the aggregator:** always on, polls each machine's `crow serve`, hosts the phone web app (PWA) and push notifications (self-hosted ntfy). If the server is off, the phone talks to one machine directly.
- **Shared crow state across machines (proposal 2026-11-03):** the phone is a thin client (web app over the tailnet), never a sync peer, which removes the mobile background-sync conflicts the owner sees with Syncthing today. Split data by kind:
  - Live state (runs, sessions now): not synced, read live from each machine's `crow serve`.
  - Append-only logs (activity, run history, retros): per-host files (`activity/YYYY-MM-DD.<host>.md`), merged at read time. No conflicts by construction.
  - Shared mutable docs (prefs, memory, open threads, review ledgers): git with automatic commit, pull --rebase and push on a timer and at wrap, remote on the home server (Forgejo or a bare repo). Conflicts become merges with history instead of conflict copies.
  - Alternative under discussion (2026-11-03): Syncthing for everything, home server as the always-on peer with staggered file versioning (history without git). Fewer moving parts and shared with the vault. Excludes: SQLite databases, `.git` dirs, live lock files.
  - Claude session transcripts (append-only, one writer each) sync safely either way; resuming on another machine needs the same project path there.
  - Secrets and account tokens: never synced, per machine.
- **Vault (owned by the vault session, to propose there):** Obsidian Self-hosted LiveSync with CouchDB on the home server handles mobile well; Syncthing can stay between desktops.

# Order (draft)
1. tmux backend (item_031, already Ready).
2. `crow serve` read only + phone view on the tailnet, one machine.
3. Actions and session chat from the phone.
4. crow on the tower (WSL) and shared state sync.
5. Home server build: OS, Tailscale, aggregator, ntfy. The rest of the homelab (DNS ad blocker, media, photos, local AI) is scoped in its own future `homelab` repo, not in crow.

# Open questions (owner)
- Decided 2026-11-03: the work laptop joins the personal tailnet like the other machines.
- Own tailnet, or share nodes with Alex's tailnet? (Tailscale node sharing allows both.)
- Yoga 2 Pro specs (RAM 4 or 8 GB, soldered). Owner's homelab wishes (2026-11-03): local AI, media server, network ad blocker, photo storage; goal is to start homelabbing and move to a small rack later.

# Acceptance criteria
- AC1: Network, sync and serve choices are decided by the owner.
- AC2: Each order step becomes a backlog item with scope and AC before build.
- AC3: Security boundary stated: tailnet only, no public port, actions authenticated.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-runs`
- `bin/crow-focus`
- `logics/request/req_015_wishlist_ambient_mobile_voice_colleague_ideas.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

