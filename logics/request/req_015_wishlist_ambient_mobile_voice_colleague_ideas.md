## req_015_wishlist_ambient_mobile_voice_colleague_ideas - Wishlist: ambient screensaver, mobile control, voice, colleague ideas
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 60
> Confidence: 50
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 17:29:18

# AI Context
- Summary: Owner wishlist for crow from 2026-11-03, inspired by two colleagues' systems (Ultron, Jarvis). Four tracks: ambient screensaver of live runs, full mobile control, voice and visual agents, automated report pipelines. Capture only; nothing is scoped or committed.
- Keywords: wishlist, screensaver, ambient display, mobile, remote control, voice, drones, browser, reports, jarvis, ultron
- Use when: picking the next crow track after the deck work, or ingesting a colleague's design doc.
- Skip when: implementing an existing deck or release task.

# Needs
1. **Ambient screensaver.** When the machine idles, show a full-screen screensaver. Reference: Antonin's Ultron screensaver. Decided 2026-11-03:
   - The engine lives in **archeotech-shell**: picker in Settings, preview, shuffle, screensaver sources as plugins (terminal animation, QML scene, web page). Easy to add new ones.
   - Source ideas: cmatrix/unimatrix, pipes.sh, asciiquarium, cbonsai, tty-clock, ASCII Star Wars (asciimation, local copy since the telnet server is unreliable), plus original QML scenes.
   - ascii-magic (ascii-magic.com/ascii, added 2026-11-05): 734 MIT-licensed animated ASCII pieces, playable in a terminal with `npx ascii-magic-art` or exported as GIF/video; also has an image/video to ASCII converter for original scenes.
   - **archeotech-dotfiles**: swayidle calls the shell's screensaver before the lock; package list gains the animation tools.
   - **crow** ships a screensaver plugin (like `integrations/archeotech/crow-agents`): a "CR-0W is working" card with run bars. Setting with three modes (overlay, take over, off), default overlay on top of the chosen animation.
   - Order: shell engine scoped in its repo first, then the dotfiles hook, then the crow plugin.
2. **Mobile control.** Split out to req_016 (crow everywhere: multi-machine and mobile control).
3. **Voice and visual agents.** Voice commands and spoken replies; "drones" in the UI that visibly carry out actions; browser automation (claude-browser). Reference: Ultron.
4. **Automated report pipelines.** Weekly slides, daily reports and similar produced with little owner effort. Reference: Alex's Jarvis. The Jarvis design doc arrived 2026-11-06 and is ingested in prod_005; reports now live in req_034, screensaver in req_036, voice in req_032.

# Context
- Existing building blocks:
  - `crow json` / `bin/crow-runs` already expose slots, agents and finished runs (data source for screensaver and mobile).
  - `bin/crow-afk-page` already renders a run status page (HTML), a possible base for both a screensaver and a phone view.
  - `bin/crow-focus --pid P --send TEXT` sends text into a live session (kitty remote control): the basis for "chat with a session" outside the terminal.
  - Desktop: Hyprland with swayidle/swaylock, so a screensaver can hook the idle timeout before the lock.
  - Claude Code Remote Control already lets the Claude mobile app drive individual Claude sessions; it does not cover crow actions (dispatch, runs, deck).
  - Report skills already exist: `daily`, `weekly`, `mgmt-weekly`, `brief`, `eod`. Track 4 is about automating and scheduling them, not new skills from scratch.
- Related: req_007 / req_011 / req_014 (deck), item_031 (terminal backend so kitty is optional; mobile chat depends on it or on kitty), req_003 (R&D priorities).
- Priority order from CLAUDE.md: these are side-project R&D, after prod launch, reviews and sprint work.

# Open questions (owner)
- Voice: local STT/TTS (whisper.cpp, piper) vs a cloud service.
- Ask Antonin how Ultron does the screensaver and drones before designing (reuse over reinvent).

# Acceptance criteria
- AC1: The four tracks are captured with their references and existing building blocks.
- AC2: Alex's Jarvis design doc is ingested and track 4 lists the features worth adopting.
- AC3: Each track is split into its own request or backlog items with scope, AC and order before any build starts.

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
- `bin/crow-afk-page`
- `bin/crow-focus`
- `logics/backlog/item_031_put_terminal_actions_behind_a_backend_so_kitty_is_optional.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
