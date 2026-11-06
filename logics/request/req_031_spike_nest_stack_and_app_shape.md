## req_031_spike_nest_stack_and_app_shape - Spike nest stack and app shape
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 70
> Confidence: 60
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 18:12:20

# AI Context
- Summary: Research spike. Pick the nest stack and app shape: React front end (owner decision), Go or Python backend, live update method, push mechanism, terminal attach, and PWA vs web app vs TWA or native shell on Android.
- Keywords: React, Vite, Go, Python, SSE, polling, ETag, Web Push, service worker, VAPID, PWA, TWA, xterm, pty, SQLite
- Use when: starting the nest code or the core API.
- Skip when: writing design rules (req_038) or Android extras (req_033).

# Needs
1. **Reference stacks.**
   - Ultron: React 19 + Vite, xterm.js, OKLCH tokens, dark only; Node `node:http` server, ws, node-pty, `node:sqlite`; polling every 3 s, long-poll for chat; Declarative Web Push without a service worker (iOS).
   - Jarvis: vanilla ES modules with hash routes; Python stdlib `ThreadingHTTPServer`, ETags, CSP; polling 30 s, 2.5 s during actions; Web Push with VAPID; no build step.
2. **Backend.** Go (single binary, matches deck, good SSE and pty libraries) vs Python (matches `bin/` scripts and logics-manager, faster to write). Criteria: deploy simplicity, pty and SSE support, SQLite access, sharing code with deck and scripts, the owner's pleasure.
3. **Front end.** React with Vite; routing; state fetching with ETags; component kit or hand-made over req_038 tokens; bundle size on the phone.
4. **Live updates.** SSE vs polling with ETags vs WebSocket. Battery on the phone, behaviour behind Tailscale, reconnects.
5. **Push.** Service worker with action buttons, grouping and badges on Android; VAPID key storage; payload limits; what happens without the app open.
6. **Live terminal.** Read-only attach to a cdx run (xterm.js + pty or tmux), "take control" flow, logging.
7. **App shape.** Installable PWA first; what a TWA or native shell adds (widgets, tiles, assistant hooks from req_033) and costs.
8. **Galaxy Z Fold 6.** The owner's phone. Check Viewport Segments API and fold posture support in Chrome and Samsung Internet, state kept across fold and unfold, the share target (req_039) and Web Push over the tailnet (HTTPS through `tailscale serve`).
9. **Auth.** Tailnet identity header, CSRF header, localhost bind; nothing public.

# Context
- Owner decisions 2026-11-06: React (work team stack, strong ecosystem); Go or Python backend; app shape by research.
- Jarvis lessons: data contract before screens, cache from day one. Ultron lesson: one event log.
- deck is Go and Bubble Tea; scripts and logics-manager are Python.

# Acceptance criteria
- AC1: A comparison table for backend, live updates and app shape with a recommendation each.
- AC2: A hello-world slice: core API section with ETag, a React page reading it, one Web Push with an action button received on the owner's Android phone over the tailnet.
- AC3: A read-only terminal attach to a running cdx run in the browser.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control.md`
- `logics/request/req_020_terminal_design_system_and_design_review_workflow.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
