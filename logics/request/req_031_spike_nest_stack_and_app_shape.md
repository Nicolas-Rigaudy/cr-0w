## req_031_spike_nest_stack_and_app_shape - Spike nest stack and app shape
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90
> Confidence: 80
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

# Spike findings (2026-11-07)
Research run, no product code. Settled calls are not reopened: React front end, core owns the API and Web Push (adr_001), ntfy only as a fallback, nest is the main surface, Fold 6 handled by breakpoints and segments only (req_038).

## Sources checked
| Source | Version or date | Used for |
| --- | --- | --- |
| crow code: `tools/crow-deck/data.go`, `live.go`, `bin/` imports, `.github/workflows/ci.yml`, `bin/crow_term.py` | origin/main 526d175 | What deck and scripts share today |
| cdx-manager `docs/cli.md` (`run --detach`, `run-tail`), `src/provider_runtime.py` | cdx 0.26.6 installed, source at c3e861f (2026-11-04) | What a running cdx run exposes |
| Throwaway probes (Go `net/http`, Python `http.server`), raw numbers in the private note | Go 1.27.3, Python 3.14.9, run 2026-11-07 | Idle CPU, memory, lines of code |
| [MDN browser-compat-data `api/Viewport.json`](https://github.com/mdn/browser-compat-data/blob/main/api/Viewport.json), [`api/DevicePosture.json`](https://github.com/mdn/browser-compat-data/blob/main/api/DevicePosture.json) | main branch, read 2026-11-07 | Fold support per browser |
| [Chrome: Viewport Segments API shipped](https://developer.chrome.com/blog/viewport-segments-api-shipped), [web-features: viewport-segments](https://web-platform-dx.github.io/web-features-explorer/features/viewport-segments) | 2025-06-09 post, read 2026-11-07 | Chrome 138, `window.viewport.segments` |
| [Device Posture intent to ship](https://groups.google.com/a/chromium.org/g/blink-dev/c/Xbhdn0u6Z9A), [interop issue 1439](https://github.com/web-platform-tests/interop/issues/1439), [W3C Device Posture](https://www.w3.org/TR/device-posture/) | issue opened 2026-09-22, CRD 2026-05-20 | Posture status |
| [Posture and segments guide](https://azukiazusa.dev/en/blog/device-posture-viewport-segments/) | 2026-10-26 | CSS media features and `env(viewport-segment-*)` names |
| [Tailscale Serve](https://tailscale.com/kb/1312/serve), [Tailscale HTTPS](https://tailscale.com/kb/1153/enabling-https) | read 2026-11-07, identity headers noted as v1.92+ | Identity headers, certs, CT logs |
| [Chrome notification actions](https://developer.chrome.com/blog/notification-actions), [MDN `Notification.maxActions`](https://developer.mozilla.org/en-US/docs/Web/API/Notification/maxActions), [Chrome 50 notifications (tag, renotify)](https://developer.chrome.com/blog/notifications) | read 2026-11-07 | Buttons, grouping |
| [Chrome Badging API](https://developer.chrome.com/docs/capabilities/web-apis/badging-api) | read 2026-11-07 | Badges on Android |
| [Chromium web push sender](https://chromium.googlesource.com/chromium/src/+/913eaefa55beab765aa1d99da170db474d8c63ca/components/sharing_message/web_push/web_push_sender.cc), [RFC 8030 section 7.2](https://www.rfc-editor.org/rfc/rfc8030#section-7.2) | read 2026-11-07 | VAPID audience, 4096-byte floor (3993 bytes of plaintext with aes128gcm), FCM 400 on oversize |
| [web.dev: subscribing a user](https://web.dev/articles/push-notifications-subscribing-a-user) | read 2026-11-07 | `userVisibleOnly` required in Chrome |
| [pywebpush 1.14.1](https://pypi.org/project/pywebpush/1.14.1/), [webpush-go](https://socket.dev/go/package/github.com/SherClockHolmes/webpush-go) | 1.14.1; v1.4.0 (2025-01-02) | Push libraries |
| [dontkillmyapp: Samsung](https://dontkillmyapp.com/samsung) | read 2026-11-07 | One UI sleeping apps |
| [web.dev WebAPKs](https://web.dev/articles/webapks), [web.dev PWA installation](https://web.dev/learn/pwa/installation) | read 2026-11-07 | WebAPK from Chrome with Google services and from Samsung Internet |
| [TWA quick start](https://developer.android.com/develop/ui/views/layout/webapps/guide-trusted-web-activities-version2), [TWA asset links on private hosts](https://groups.google.com/a/chromium.org/g/chromium-dev/c/Ki3TZIGfmTY) | read 2026-11-07, thread from 2018 | TWA cost and verification |
| [xterm.js 6.0.0](https://unpkg.com/@xterm/xterm@6.0.0/README.md), [addon-attach](https://xtermjs.org/docs/api/addons/attach/) | 6.0.0 | Browser terminal |
| [Vite 8 with Rolldown](https://alternativeto.net/news/2026/3/vite-8-debuts-with-unified-and-faster-rolldown-bundler-integrated-devtools-and-wasm-ssr), [Rolldown 1.0](https://voidzero.dev/posts/announcing-rolldown-1-0), [React releases](https://www.versio.io/it/product-release-end-of-life-eol-facebook-react.html) | Vite 8 (2026-03), React 19.3.0 (2026-09-09) | Front end versions |
| [SSE closed in hidden tabs](https://data-star.dev/how_tos/prevent_sse_connections_closing), [stream lost after phone sleep](https://github.com/anomalyco/opencode/pull/17770) | read 2026-11-07 | SSE on a phone |

## What crow does today
- deck (Go, about 17k lines) holds no data logic of its own: it execs the Python scripts' JSON (`crow-runs json`, `crow-sessions list --json`, `crow-usage --json`, `crow-accounts json`, `crow-registry json`) and maps it to views. Run kinds, permissions, registry and accounts live in `bin/crow_lib.py` and the scripts (about 11k lines of Python).
- Python already pulls PyYAML and markdown; CI sets up Python 3.12 (`pytest pyyaml markdown`) and Go from `tools/crow-deck/go.mod` and runs both test suites.
- cdx runs Claude headless as `claude --print --output-format json` with a fixed `--session-id` (`src/provider_runtime.py`). The run's stdout holds one JSON result at the end; `cdx run-tail` returns raw stdout lines. A running cdx run has no terminal. Its live trace is the provider transcript (JSONL), which the crow notifier already watches for stalls.
- `crow_term.py` opens interactive sessions in kitty or tmux (tmux is not installed on the laptop today).

## Need 1: reference stacks
- Kept from Ultron: React 19 with Vite, xterm.js for the terminal, one event log. Not kept: Node server and node-pty (no Node in crow's runtime), Declarative Web Push without a service worker (aimed at iOS; nest targets Chrome on Android and needs a service worker for action buttons and the share target).
- Kept from Jarvis: Python stdlib `ThreadingHTTPServer`, ETags, CSP, Web Push with VAPID, data contract before screens. Not kept: vanilla modules with hash routes (React is decided), 30 s polling as the only live channel.
- The probe below reproduces the Jarvis server shape and a Go equivalent.

## Need 2: backend
Probe: one ETag section endpoint and one SSE stream, stdlib only, 3 SSE clients plus one conditional poller every 3 s.

| | Go `net/http` | Python `ThreadingHTTPServer` |
| --- | --- | --- |
| Idle CPU, 120 s | 30 ms (0.025 %) | 60 ms (0.05 %) |
| One change per second, 60 s | 50 ms | 80 ms |
| Memory (RSS) | 12 to 14 MB | 23 MB |
| Lines | 85 | 56 |
| Per open stream | goroutine | thread |
| Artifact | 10 MB static binary | script, needs the interpreter |

| Criterion | Go | Python |
| --- | --- | --- |
| Deploy | One binary, systemd unit | Script plus PyYAML, already how crow installs |
| SSE | Goroutines, `http.Flusher` | Works in stdlib (probe); one thread per stream is fine for one owner |
| pty | `creack/pty` (third party) | stdlib `pty`, `os.openpty` |
| SQLite | cgo driver or a pure-Go module | stdlib `sqlite3` |
| Web Push crypto | stdlib ECDH, AES-GCM, HKDF, or webpush-go | needs `cryptography` (pywebpush) |
| Sharing | Would re-implement `crow_lib`, registry, kinds and policy, or exec the scripts like deck | Imports `crow_lib`, registry, kinds files and logics-manager directly |
| Idle cost | Negligible | Negligible |

Recommendation: Python core. The domain code it needs (kinds, accounts, registry, runs, Logics docs) is Python today and deck is a thin client of it. Performance is not a criterion at one owner (both probes idle near zero). Go keeps deck. Decision 1.

## Need 3: front end
- React 19 (19.3.0 latest, 2026-09-09) with Vite 8 (Rolldown bundler), TypeScript.
- Routing: URL routes for every screen and overlay state (mission id, open sheet), so a fold or unfold, a reload and a notification deep link land on the same view (req_038 Need 4 and 5). A small router is enough; no server rendering.
- Fetching with ETags: core answers `Cache-Control: no-cache` plus `ETag`. The browser revalidates by itself and a 304 reaches the code as a 200 from cache, so the client needs no ETag logic. Each section is one query keyed by name; an SSE `section` event invalidates that key. Every section carries `updated_at` and the screen shows its age (req_038 Need 8).
- Component kit: none. req_038 tokens derive from the deck ANSI theme; a styled kit would fight them. Hand-made components over CSS custom properties, with headless primitives only where focus handling is hard (sheets, palette). Decision 3.
- Bundle: `react-dom` client 19.2.8 is 95 kB gzip unminified (measured from an existing local install); the minified figure and a budget come from the hello-world slice. xterm.js loads only on the live terminal route.

## Need 4: live updates
| | SSE + ETag refetch | Polling + ETag | WebSocket |
| --- | --- | --- | --- |
| Latency | Instant | Poll interval (3 to 30 s) | Instant |
| Phone battery | One idle socket while visible, 15 s heartbeat | One request per interval, 304 is a few hundred bytes | One idle socket |
| Hidden tab or phone asleep | Closed by the browser; reconnects on resume (sources above) | Timers stop | Same as SSE |
| Reconnect | Built into `EventSource` (`retry`, `Last-Event-ID`) | Nothing to reconnect | Hand-written |
| Behind `tailscale serve` | Streaming not verified | Plain requests | Not verified |
| Server cost | Probe above | Lowest | Library in both stacks |

Recommendation: SSE for "section X changed, version N" only; the client refetches that section with its ETag. Close the stream on `visibilitychange` hidden, reopen and refetch all sections when visible. Polling with ETags every 5 s is the fallback when the stream fails twice. Web Push covers the app closed or the phone asleep. WebSocket only for the live terminal (Need 6). Decision 2.

## Need 5: push
- Chrome requires `userVisibleOnly: true`: every push shows a notification, no silent pushes.
- Delivery goes through the browser's push service (FCM for Chrome). core needs outbound internet; the phone needs none of the tailnet to receive. The payload is encrypted end to end. Plaintext limit 3993 bytes with `aes128gcm`; FCM answers 400 when it is too big. Send the mission id, version, title and one line; the click fetches the rest.
- Buttons: Chrome shows 2 action buttons (Chrome blog; read `Notification.maxActions` on the phone). `notificationclick` gets `event.action`; the service worker POSTs the verb with the mission version and the CSRF header, then closes or opens the mission. core refuses an answer for a stale version. A tap on an action works without nest open; the phone must reach the tailnet at that moment.
- Grouping: one `tag` per mission replaces the previous notification silently; `renotify: true` only on a new question.
- Badges: the Badging API is not supported on Android. Android badges the installed app icon from unread notifications by itself.
- VAPID: core generates one P-256 key pair at first start, private key in core's config dir outside every repo (0600), public key served to the subscribe page. Subscriptions are stored per device with their endpoint and removed on 404 or 410. Matches the req_030 credential table (core holds the Web Push keys).
- Risk: One UI puts unused apps to sleep. Chrome is used daily; listing Chrome under "never sleeping apps" is the remedy if pushes arrive late.

## Need 6: live terminal
- A running cdx run has no terminal to attach to (headless `--print`). The live view of a run is its transcript: core tails the JSONL for the run's session id and streams new entries over SSE, rendered as a timeline (text, tool calls, results). Same trace for isolated runs if cdx exposes it (req_030: isolated runs are headless too).
- A real terminal exists only in take control. core stops the run, starts the interactive session in the mission worktree under a pty it owns (or tmux), and nest attaches with xterm.js 6.0 plus `@xterm/addon-attach` over a WebSocket. Read-only attach drops input server side. core records the pty output to the mission dir and logs `manual_gesture` at start and end (req_030 escape hatch).
- AC3 asks for a terminal attach to a running cdx run; that object does not exist. Decision 4.

## Need 7: app shape
| | Installed PWA | TWA (Bubblewrap) | Native shell |
| --- | --- | --- | --- |
| Install | Chrome install, WebAPK on devices with Google services | APK, sideloaded or Play | APK |
| Push, share target, shortcuts | Yes (manifest) | Yes, through Chrome | Yes, own code |
| Home widgets, quick tiles, assistant hooks (req_033) | No | Only with native Android code added to the project | Yes |
| Fold APIs | Chrome | Chrome | Lost in a WebView (Device Posture is off in WebView) |
| Cost | Manifest and a service worker | Android project, signing key, Digital Asset Links on the tailnet origin, Android toolchain | Full app |
| Updates | Every deploy | Web part every deploy; shell rebuilt on manifest change | Rebuild per change |

Recommendation: installed PWA first. TWA when req_033 extras are groomed and need native widgets or tiles; no WebView shell. Two unknowns need the phone: whether Chrome mints a WebAPK for a tailnet-only origin (else a shortcut, and the share target of req_039 may not register), and whether TWA asset-link verification works on a private origin (it needed a public origin in 2018; the thread says it later moved on device). Decision 5.

## Need 8: Galaxy Z Fold 6
| API | Chrome Android | Samsung Internet | Android WebView |
| --- | --- | --- | --- |
| Viewport Segments (`window.viewport.segments`, `horizontal-viewport-segments`, `vertical-viewport-segments`, `env(viewport-segment-*)`) | 138 (2025-06) | Mirrors Chrome 138 in MDN data, not checked on a release | Mirrors Chrome in MDN data |
| Device Posture (`navigator.devicePosture`, `device-posture: folded or continuous`) | 132 in MDN data (intent to ship said 131) | 16.2 | No |
| Firefox, Safari | No | | |

- Segments describe the split, posture the fold angle; a flat device reports one segment and `continuous`. Whether the Fold 6 reports two segments in Flex mode in Chrome is not documented: hand check.
- State across fold and unfold: the screen and the open overlay live in the URL, drafts in `sessionStorage` keyed by mission and version, layout from container breakpoints and segments only. Chrome resizes the page on a fold; the hand check confirms no reload (req_038 AC4).
- Share target (req_039) needs the installed app; see the WebAPK unknown above.
- Web Push over the tailnet works through `tailscale serve` HTTPS (a valid public certificate on the tailnet name gives the secure context push needs). Receiving goes through FCM, not the tailnet.

## Need 9: auth
- `tailscale serve` on the control host proxies HTTPS to core on `127.0.0.1`. Serve adds `Tailscale-User-Login` and strips any copy sent by the client. core accepts a request only when that header equals the owner's login from config. Tagged devices get no identity headers.
- Gap: every local process also reaches `127.0.0.1` and can forge that header. core therefore listens for Serve on one loopback port that requires the header, and serves local clients (deck, scripts, MCP) on a unix socket with file mode 0600. No loopback request is trusted without one of the two.
- CSRF: every POST needs a custom header (for example `X-Crow: 1`) and a matching `Origin`. No CORS headers are ever sent, so a cross-site page cannot add the header without a preflight that fails.
- Certificates land in public Certificate Transparency logs with the machine name: the serving machine gets a neutral name before `tailscale cert`.
- Nothing public: no Funnel (Funnel traffic carries no identity headers).

## Recommended stack
- core: Python, stdlib HTTP server (threaded), stdlib `sqlite3` store rebuilt from the JSONL event log, `cryptography` for Web Push, stdlib `pty` for take control, systemd user unit.
- nest: React 19, Vite 8, TypeScript, URL routes, section queries over ETag with SSE invalidation, hand-made components over req_038 tokens, installed PWA with a service worker for push and share target.
- Live: SSE invalidation, ETag refetch, 5 s polling fallback; WebSocket only for take control.
- Terminal: transcript timeline for runs; xterm.js over a core-owned pty for take control.
- Auth: `tailscale serve` identity header on a dedicated loopback port, unix socket for local clients, CSRF header and Origin check.
- deck stays Go and reads core over the unix socket once core exists.

## Open decisions (for the owner)
1. **Backend language.** Python (shares `crow_lib`, registry, kinds, logics-manager; stdlib sqlite and pty) or Go (single binary, one language with deck, would re-implement or exec the Python domain code). Recommendation: Python.
2. **Live updates.** SSE invalidation plus ETag refetch with polling fallback, polling with ETags only (simplest, Jarvis), or WebSocket for everything. Recommendation: SSE plus ETag with polling fallback.
3. **Component kit.** Hand-made over req_038 tokens with headless primitives for sheets and palette, a headless library everywhere, or a styled kit themed to the tokens. Recommendation: hand-made plus headless primitives where needed.
4. **Live terminal scope (AC3).** Rewrite AC3 as "live transcript timeline of a running run, plus a read-only xterm attach in take control", ask cdx for a pty-backed run mode, or trial an existing mobile client for live sessions first (the roadmap review lists one). Recommendation: rewrite AC3 as the first option.
5. **App shape.** Installed PWA now, TWA later for req_033 widgets and tiles; TWA from day one; or PWA only and drop native extras. Recommendation: PWA now, TWA later.

## Hand checks for the owner (AC2 and Fold rows)
- On the Fold 6 in Chrome: `Notification.maxActions`, one push with two buttons from a tailnet origin with the app closed, install result (app or shortcut), share target listed, two segments in Flex mode, no reload across fold and unfold. Samsung Internet: same segment and posture checks.

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
