## item_136_read_cdx_profiles_through_the_cli - Read cdx profiles through the cdx CLI: profile dirs, relaunch by conversation, live sessions
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub portability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Replace crow's hand-built cdx profile paths with the cdx JSON CLI in three slices: profile and transcript dirs, relaunch by conversation id, live sessions.
- Keywords: cdx sessions, profile_dir, transcripts_dir, can-resume, conversation, stdout_path, /proc
- Use when: building req_042.
- Skip when: the read has no cdx equivalent (see Out).

# Problem
- About 15 call sites glob `<profile>/claude-home/...` directly. That layout is a cdx internal and Claude-only.

# Scope
- In, in this order (one commit each, tests green after each):
  - S1 profile and transcript dirs. Add one cached helper in crow_lib, `cdx_profiles()` from `cdx sessions --all --json` (name, provider, profile_dir, transcripts_dir), with a fallback to today's layout when cdx is missing or older than 0.23. Route through it: `crow_lib.claude_transcripts` (:95), `crow-runs` `run_transcript` (:170) and `report_from_transcript` (:447), `crow-usage` transcript scan (:749), `crow-memory` (:24), `crow-cost-audit` (:83), `crow-relaunch` fallback scan (:142). `crow-usage` `dispatch_titles` (:588): take the session id from `cdx run-report --json` / `cdx runs --json` (`conversation_id`, `artifacts.stdout_path`) instead of globbing `cdx-run-<id>.stdout.log`.
  - S2 relaunch by conversation. `crow-relaunch` checks `cdx can-resume <account> --conversation <id> --json` before launching and resumes with `--conversation <id>`; `conversation_not_found` is a clear error, no fallback to `--last`. Works for Codex accounts too.
  - S3 live sessions. `crow-usage` `live_session_ids` (:616) and `crow-runs` `collect` (:628) take cdx-launched sessions from `cdx sessions --json` (pid, provider, cwd, conversation_id); `agent_procs` (/proc) stays for sessions started outside cdx, deduplicated by pid.
- Out:
  - `crow-usage` `rtk_homes` (rtk data, not cdx), `crow-ctx` settings.json, `crow-doctor` profile checks, `crow_lib` setup-token file: no cdx CLI equivalent, keep as is.

# Acceptance criteria
- AC1: S1 done: `grep -rn 'claude-home' bin` only hits crow_lib's fallback and the Out list.
- AC2: S2 done: relaunch of a real Claude and a real Codex conversation id works; an unknown id fails with conversation_not_found.
- AC3: S3 done: `crow runs` shows the same live sessions as before on a machine with cdx sessions running, with no duplicates.
- AC4: tests fake `cdx sessions` / `can-resume` output; full suite green. Closing comment on #6 drafted for the owner's OK.

# AC Traceability
- request-AC1 -> S1. Proof: AC1 grep.
- request-AC2 -> S2. Proof: AC2 live check.
- request-AC3 -> S3. Proof: AC3 live check.
- request-AC4 -> all slices. Proof: AC4.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_042_crow_reads_cdx_through_its_cli`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: no user-facing bug today; removes silent breakage on a cdx layout change and is a prerequisite for Codex parity. Below task_016.

# Notes
- Line numbers checked 2026-11-09 against main ac1d396.
- `cdx sessions --all --json` checked live on cdx 0.26.11: rows carry session, provider, active, pid, cwd, conversation_id, profile_dir, transcripts_dir.
- Fresh session: `implement-task` on this item, one slice per commit.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
