## item_068_quota_warning_and_account_failover - Quota warning and account failover
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Warn live sessions before the 5h limit, fail background runs over to an allowed account, resume limit-killed runs at reset (lever 15).
- Keywords: usage limit, quota, failover, account, ask-first, resume, reset
- Use when: changing how crow reacts to account limits.
- Skip when: changing per-account launch settings only.

# Problem
- 2026-11-03: the default account ran out at 15:19 UTC with no warning. Run ab76eee8 died. The owner switched to the shared account by hand.

# Scope
- In:
  - UserPromptSubmit hook line (via `crow-runs inbox --hook` or a new `crow-usage --hook`): at 80% and 95% of the current account's 5h window used, one line with the reset time and the best allowed fallback (`cdx status --cached` order, policy `use` first, `ask-first` named as "ask").
  - At 95%: tell the session to /wrap now (handoff note) so the switch costs one reload, not a lost context (lever 15: compact what would cool before reset).
  - Dispatch: `cdx run --failover` restricted to accounts with policy `use`; `ask-first` accounts only after an OK given in the hook prompt or a standing OK in accounts.yaml.
  - A run that dies on `limit`: requeue it with `--continue-run` on a fallback, or schedule it at the reset time. State `limit`, never `failed`, in the inbox hook.
  - `crow rescue <account>`: every session of the exhausted account interrupted in the last window (crow-sessions), each resumed on the fallback with `cdx handoff`/`rescue` and a resume note, after one OK listing them.
  - `crow doctor`: flag when every `use` account is above 90%.
- Out:
  - Using teammates' accounts (policy `never`).
  - Changing cdx itself; file upstream if `--failover` cannot filter by policy.

# Acceptance criteria
- AC1: With the current account at 80% (faked sample), the next prompt shows the warning with reset time and fallback.
- AC2: A run killed by a limit is shown as `limit` and restarts on a fallback or at reset without the owner.
- AC3: `crow rescue <default> --dry-run` lists all sessions the limit stopped (2026-11-03: five), not only the last one.
- AC6: With the 5h reset under 45 min away, crow keeps using the shared account past 50% of the window, up to 100%.
- AC5: Crow stops picking the shared account once its 5h or week use passes 50%.
- AC4: No run ever lands on a `never` account; an `ask-first` account only after an OK.

# AC Traceability
- AC1 -> `crow-relaunch hook` (UserPromptSubmit). Proof: tests/test_crow_quota_failover.py::test_hook_warns_once_per_level (82% then 96%, once each, reset and fallback named); live run with the default account at 100% printed the 95% warning with "Fallback: <shared> (no OK needed)".
- AC2 -> `crow-runs` run_state, inbox, notify/requeue. Proof: test_run_state_reads_limit_from_tail, test_requeue_moves_to_auto_fallback_with_resume_note, test_requeue_waits_for_reset_then_continues, test_requeue_stops_after_max_chain.
- AC3 -> `crow rescue <account> --dry-run`. Proof: test_rescue_lists_every_stopped_and_open_session (three cut-off and one open session listed; stale, relaunched, other-account and run entries left out). On 2026-11-03 the three real cut-off sessions were already relaunched by hand; the other four limit stops were /retro one-shots, which the registry leaves out.
- AC4 -> `crow_lib.pick_accounts` and `crow-runs pick_account`. Proof: test_plain_ask_first_is_never_auto, test_dispatch_refuses_ask_first_without_ok; `never` accounts are dropped by the picker and refused by `crow-accounts require`.
- AC5 -> `crow_lib.cap_block`. Proof: test_shared_caps_block_auto_fallback (5h 60% and week 55% both block).
- AC6 -> `crow_lib.cap_block` free_5h_min. Proof: test_shared_5h_cap_lifts_near_its_reset (70% used, reset in 30 min, still picked).

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Request: req_017_token_levers_quota_switch
- Primary task(s): (none; built directly from this item)

# Priority
- Priority: High
- Rationale: A limit stopped all hub work on 2026-11-03.

# Notes
- Check `cdx run --failover` and `cdx select --require-ready` behaviour first; build on them.
- Decided 2026-11-03: the shared account is an automatic fallback (accounts.yaml `fallback: auto`) with caps: crow uses at most 50% of the shared account's 5h window and 50% of its week, because it is shared. Past a cap, the shared account counts as blocked. Exception (owner, 2026-11-03): when the shared account's 5h window resets within a short time (default 45 min, `fallback_free_5h_min`) and its week use is under the week cap, crow may use up to 100% of the current 5h window, since no one else is likely using it. Policy stays `ask-first` in cdx so existing code keeps working; the caps need code here.
- Built 2026-11-03: picker in `bin/crow_lib.py` (pick_accounts, account_state), accounts.yaml keys `fallback: auto`, `max_5h_pct`, `max_week_pct`, `free_5h_min` (read by crow-accounts), warning hook `crow-relaunch hook`, `crow rescue`, dispatch `--account` and auto fallback, limit requeue in the run notifier (max 3 per task), `crow doctor` quota check. Runbooks: interrupted-sessions.md, dispatch-background-runs.md.
- `cdx run --failover` was not used: it never picks ask-first accounts, so the auto fallback lives in crow. `--continue-run` cannot change account, so a requeue on another account is a fresh run with a RESUME NOTE.
- `cdx rescue` moves only the latest conversation per account, and cdx has no policies set (it treats every account as use). `crow rescue` is the one to use.
- Existing tooling to build on: `crow runs sessions` finds limit-stopped sessions, `crow runs relaunch <id> <account>` hands them off (typed `yes` for ask-first), `crow runs relaunch status` names the target. The gap is the trigger and the warning, not the relaunch itself.
