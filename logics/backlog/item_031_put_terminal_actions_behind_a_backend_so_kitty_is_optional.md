## item_031_put_terminal_actions_behind_a_backend_so_kitty_is_optional - Put terminal actions behind a backend so kitty is optional
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: release
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Kitty becomes one terminal backend; focus, open, relaunch and screenshots work without it.
- Keywords: kitty terminal backend focus relaunch ansi-shot portable
- Use when: Making crow run outside kitty, or adding a terminal backend.
- Skip when: Release plumbing.

# Problem
- kitty is hard-wired in four scripts: `bin/crow` (`term`, `focus`), `bin/crow-focus` (kitty remote control sockets), `bin/crow-relaunch`, `bin/ansi-shot` (reads the kitty theme file).
- The owner wants crow usable on macOS and, through WSL, Windows.

# Scope
- In:
  - List every kitty-specific call (above) and the real behaviour each provides: open a window, focus a window by pid, resume in a window, take a theme.
  - A small backend interface (python module in bin/, or a `CROW_TERMINAL` env/config key): kitty, plus a fallback that prints the command to run. Add a second real backend only if cheap (tmux is the likely one, it also works in WSL).
  - `ansi-shot` takes `--theme` for any kitty-format file and ships a default theme, no dependency on `~/.config/kitty`.
  - Docs state which features need which backend.
- Out:
  - Native Windows terminals.
  - A GUI for choosing backends.

# Acceptance criteria
- AC1: With no kitty installed, `crow term`, `crow focus` and `crow-relaunch` do not crash: they print the command or use the configured backend.
- AC2: With kitty, behaviour is unchanged (existing tests still pass, focus tests cover the backend).
- AC3: `ansi-shot` renders with the bundled default theme when no kitty theme exists.
- AC4: Feature-by-platform table written for item_037.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: bounded slice with testable criteria above.
- request-AC4 -> The slice's criteria feed the first tagged release and the install/update round trip.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Needed so other Linux users, macOS and WSL are not blocked.

# Dependencies
- Independent of item_030. Lets item_032 run the macOS job green.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
