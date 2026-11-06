## req_030_spike_mission_safety_and_sandbox - Spike mission safety and sandbox
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90
> Confidence: 75
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 21:04:50

# AI Context
- Summary: Research spike. Pick the layered safety model for unattended missions from the two colleague designs and current Claude Code features: guard hook and validation records (Ultron), OS sandbox with egress rules and an engine that holds every key (Jarvis), native sandbox and permission modes, item_077 deny-by-default runs. Define the personal auto-merge gate.
- Keywords: safety, sandbox, bubblewrap, pasta, nftables, guard hook, Stop hook, allowlist, egress, secrets, engine publishes, auto-merge, secret scan
- Use when: deciding what a mission agent can touch and who publishes.
- Skip when: tuning one allowlist for one run kind (item_077).

# Needs
1. **Options to compare.**
   - Ultron: about 170 allow rules in run settings; a PreToolUse guard hook that blocks push, MR creation and destructive git unless a `validated/KEY` marker exists; a Stop hook that checks the phase output exists before letting a run end.
   - Jarvis: bubblewrap with a temp home, project dir bound read-write, no git or forge credentials inside; network through pasta with nftables blocking the LAN and the tailnet; the engine pushes, opens PRs and merges.
   - Claude Code native sandbox and permission modes; cdx `--permission` per account.
   - crow item_077: deny by default, whitelist per run kind, scrub-env.
2. **Threats.** Destructive command on the host, secret read and leak, LAN or tailnet access, unwanted push or post, prompt injection from a ticket, MR or mail, runaway cost.
3. **Layered design.** For each threat: what stops it in run config, in cdx, in core, and at the OS. Who holds which credential.
4. **Auto-merge gate (personal).** Green CI, secret scan over the commit range, diff size limit with fallback to asking, sensitive paths always ask (CI, sandbox, units, policy, merge gate), full-suite receipt before main, deploy hook after merge.
5. **Work rules.** Push, MR, merge, Jira and Confluence writes, mail and Teams only through core after a validation record for the exact text.
6. **Escape hatch.** Live terminal "take control" and how it is logged as a manual gesture.

# Context
- Owner decision 2026-11-06: research what is doable and recommended, steal from both.
- Memory: scrub-env broke headless runs until `--allowedTools` and socat; MCP inside sandboxes untested.
- Builds on req_019 / item_077 and adr_001 (engine holds the keys).

# Spike findings (2026-11-06)
Research run, no product code. Settled calls from prod_005, adr_001 and the owner's roadmap review are not reopened: core alone publishes, work repos never merge, personal repos auto-merge on green CI and a secret scan, sensitive paths always ask, one "pause all publishing" switch.

## Sources checked
| Source | Version or date | Used for |
| --- | --- | --- |
| [Claude Code sandboxing](https://code.claude.com/docs/en/sandboxing) | Claude Code 2.1.319, read 2026-11-06 | Native Bash sandbox, its scope and limits |
| [Claude Code permissions](https://code.claude.com/docs/en/permissions) | Claude Code 2.1.319, read 2026-11-06 | Modes, rule precedence, hooks vs rules |
| cdx-manager `docs/cli.md`, `docs/operations.md`, `src/workspace_sandbox.py`, changelogs 0.21.6 to 0.26 | cdx 0.26.6 installed, source at 2026-11-04 | `--scrub-env`, `--isolate-workspace`, `--network-access`, `--continue-run` limits |
| [GitHub push protection for public repos](https://github.blog/changelog/2023-05-09-secret-scannings-push-protection-is-available-on-public-repositories-for-free), [default on for user accounts](https://github.blog/news-insights/product-news/keeping-secrets-out-of-public-repositories/) | 2023-05 and 2024-02 posts, read 2026-11-06 | Server-side secret block on personal public repos |
| [GitHub auto-merge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request) | read 2026-11-06 | Repo setting plus required checks |
| [gitleaks `--log-opts` range scans](https://developer.harness.io/docs/security-testing-orchestration/sto-techref-category/gitleaks-scanner-reference) | read 2026-11-06, gitleaks not installed here | Secret scan over a commit range |
| crow code: `bin/crow_lib.py` (`RUN_KINDS`, `run_flags`, `permission`), `bin/crow-runs`, `bin/crow-hub-guard`, `bin/crow-prepush-check` | origin/main 8e8c3b4 | What crow does today |
| Private design notes of both reference systems (`local/reference/`) | 2026-11-06 | Option descriptions |

## What crow does today
- Task and one-shot runs launch with cdx `--permission auto` (`crow_lib.permission`); reviews use `--permission review` (plan mode). Every run carries a budget, a timeout, `--failover` and two headless env vars (`crow_lib.run_flags`). The notifier kills a run whose transcript stalls 20 min.
- No run uses `--scrub-env`, `--strict-mcp-config` or `--isolate-workspace` yet. item_077 adds the first two (deny by default, per-kind allowlists).
- `crow-hub-guard` asks before commit, push, rebase, merge or MR in a registered repo from a hub session. `crow-prepush-check` refuses a push whose new commits match a deny list. Both are host-side, neither runs inside background runs.

## Option facts
- **Allow rules + guard hook + Stop hook** (Ultron). Rules and hooks live in Claude Code, inside the agent's process. Rules are evaluated deny, then ask, then allow; a PreToolUse hook cannot override a deny rule ([permissions](https://code.claude.com/docs/en/permissions)). A Bash rule does not match the same program by path or inside `sh -c`, and Read deny rules do not stop a script that opens files itself (same page). So the guard is a policy layer for an honest agent, not a boundary against an injected one.
- **OS sandbox + engine publishes** (Jarvis). bubblewrap with a temp home, the worktree bound read-write, read-only git metadata, an env allowlist, pasta for a private network namespace and an nftables chain that rejects private, CGNAT (tailnet) and loopback ranges. No forge credentials inside; the engine pushes, opens PRs and merges.
- **cdx already ships the Jarvis boundary.** `cdx run --isolate-workspace` (cdx 0.24.5, pasta egress since 0.24.6) builds exactly that sandbox: allowlisted read-only toolchains, `/home/agent` volatile home, the checkout bound read-write (read-only for `--permission review`), a read-only git view with no remotes, hooks or credential helpers, an env allowlist, and the pasta + nftables chain above. Online runs refuse to start without pasta ("no shared-network fallback"). Isolated Claude needs a cdx OAuth setup token (`CLAUDE_CODE_OAUTH_TOKEN`) and sets `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=1` inside. Limits: `--continue-run` cannot be combined with `--isolate-workspace` (0.24.9); local runs keep the public internet open (the inference-only allowlist applies to remote Claude runs); git metadata is read-only, so the agent cannot commit ("committing and pushing remain the supervising mission flow's job", `docs/operations.md`).
- **Claude Code native sandbox.** Off by default. On Linux: bubblewrap for Bash and its children, socat to a host proxy that filters by hostname (`sandbox.network.allowedDomains`), writes limited to the working dir, reads open to most of the machine unless `filesystem.denyRead` is set. Read, Edit, Write, WebFetch, MCP servers and hooks run outside it. Domain fronting can reach hosts outside the allowlist because the proxy does not inspect TLS. `sandbox.failIfUnavailable` and `allowUnsandboxedCommands: false` make it fail closed. With `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB` on Linux every command runs sandboxed, which explains why the 2026-10-29 scrub-env probe needed socat.
- **cdx run config.** `--scrub-env` strips credential-looking env vars from what the agent starts; GitHub tokens and proxy settings stay (`docs/operations.md`), and it hides no files. It forces Claude's default permission mode, so headless runs need `--allowedTools` through `--extra-args`. `--budget`, `--timeout-seconds` and `--max-parallel` cap cost. `permission_denials` come back as run warnings.
- **item_077** (being built tonight by another run, not edited here): default permission mode, per-kind allowlist next to `RUN_KINDS`, scrub-env, strict MCP config, denials surfaced in `crow runs`. It is the run-config layer of the design below.

## Probe (bwrap on this machine)
Throwaway script and raw output kept outside the repo (path in the private note). Fixture repo on a `mission/probe` branch, a bare local remote, an HTTPS remote, a synthetic `FAKE_GH_TOKEN` in the launcher's env. bubblewrap 0.13.1, no pasta, no tailscale.

| Check inside the sandbox | Net off (`--unshare-net`) | Net on, shared host network |
| --- | --- | --- |
| HOME is a temp dir, owner home holds only the mount path to the repo | yes | yes |
| Forge config (`gh` hosts file), `~/.ssh`, `~/.aws` readable | no | no |
| Any `TOKEN`/`KEY`/`SECRET` env var | no | no |
| Write and commit in the bound repo | yes | yes |
| Push to the local remote (outside the bind) | no | no |
| Push to the HTTPS remote (no credentials) | no | no |
| Public HTTPS | no | yes |
| Home LAN gateway | no | **yes** |
| A web UI the host serves on 127.0.0.1 | no | **yes** |
| Write outside the repo | no | no |

- After the run, the host pushed the agent's commit on `mission/probe` to the remote: core can publish what the sandbox produced.
- With `.git` bound read-only (cdx layout) the agent edits files but `git commit` fails: core must commit, or cdx must offer writable git metadata.
- Verdict: credentials and push are contained by bwrap alone. Network must be off, or go through pasta + nftables: a shared network namespace exposes the LAN and host loopback services. AC3 is met for "cannot push" and "cannot read credentials", not for "online but no LAN" until pasta is installed. No agent ran inside the probe: cdx refuses isolated Claude without pasta.

## Threats by option (AC1)
Verdict per cell: **stops** (holds against an injected agent), **partial** (holds for an honest agent or for part of the threat), **no**.

| Threat | Allow rules + guard (Ultron) | Sandbox + engine publishes (Jarvis, = cdx isolate) | Claude native sandbox | cdx scrub-env + budget | item_077 deny by default |
| --- | --- | --- | --- | --- | --- |
| Destructive command on the host | partial: rules miss `sh -c`, scripts | stops: only the worktree is writable | partial: Bash only; Edit/Write, MCP, hooks outside | no | partial: allowlist, no boundary |
| Secret read (files, env) | no: owner home and env visible | stops: temp home, env allowlist (probe) | partial: reads open unless `denyRead`; Read tool outside | partial: env only, GitHub tokens kept, files visible | partial: as scrub-env |
| Secret or code leak to the internet | no | partial: public internet open; nothing secret inside | partial: domain allowlist, fronting possible | no | no |
| LAN or tailnet access | no | stops: nftables rejects private and CGNAT ranges | partial: Bash only | no | no |
| Unwanted push or post | partial: guard parses commands, needs a marker | stops: no remotes, read-only git, no tokens | partial: network allowlist can drop forge hosts | partial: env tokens stripped, files remain | partial: push not allowlisted |
| Prompt injection (ticket, MR, mail) | partial: blast radius = what rules allow | partial: blast radius = the worktree and public net | partial | no | partial |
| Runaway cost | partial: Stop hook rounds | partial: 2 repair rounds | no | stops: `--budget`, `--timeout-seconds`, `--max-parallel` | no (crow watchdog already does) |

No layer stops prompt injection itself. It is contained by what the injected run can reach (rows above) and by core: work publish needs the owner's OK on the exact payload, personal merge passes the gate in Need 4, sensitive paths always ask.

## Recommended layered design (AC2)
Steal the boundary from Jarvis (through cdx, which already has it), the validation file and the phase Stop check from Ultron, and keep item_077 as the run-config layer.

| Threat | Run config (kind file, item_077) | cdx | core | OS |
| --- | --- | --- | --- | --- |
| Destructive command | Per-kind allowlist, default mode | `--isolate-workspace` | Worktree per mission, reconciler | bwrap: only the worktree writable |
| Secret read | No MCP unless listed | volatile home, env allowlist, setup token only | Holds every forge, tracker, mail and cloud credential, files 0600 outside repos | Owner home not mounted |
| Leak | WebFetch off by default | Inference-only egress where possible; `--network-access` only for kinds that install packages | `data_class` routes private work to local or work accounts | pasta + nftables |
| LAN / tailnet | none | pasta + nftables | API bound to localhost and tailnet, owner header | nftables rejects private and CGNAT |
| Push / post | Push, MR and post tools absent | No remotes, read-only git | Publishes after the gate (Needs 4 and 5); pause switch | No credentials in the namespace |
| Injection | Untrusted text passed as files, never inline | same | Validation on exact payload hash; message after `ok` drops it | Blast radius limited above |
| Cost | Kind budget and timeout | `--budget`, `--timeout-seconds`, `--max-parallel` | Resume cap, watchdog, quota probe | none |

**Who holds which credential**

| Credential | Holder | Inside a mission run |
| --- | --- | --- |
| Provider login (Claude setup token, Codex auth) | cdx, per account | Yes, only the provider binary's own login |
| Forge tokens (`gh`, `glab`), SSH keys | core (host) | Never |
| Tracker, wiki, mail and chat tokens | core (host), owner-approved scopes | Never; MCP servers that need them run in core, read tools only |
| Cloud credentials (SSO) | Owner, interactive only | Never |
| Web Push keys, core API identity | core | Never |
| Owner validation | Event log (owner `ok` events) | Never writable by a run; the MCP server exposes no `ok`, `go` or approve tools (adr_001) |

Ultron's guard hook stays as a second layer for runs that cannot be isolated (continued runs, interactive take-control sessions, Codex hub sessions): deny push, MR creation, merge, tracker and wiki writes and gate answers unless core published them. The Stop hook check (a turn may end only on its phase output) belongs to the run contract (req_037).

## Personal auto-merge gate (Need 4, AC4 as a checklist)
core runs these in order on a mission branch and stops at the first failure, which turns into an inbox item. Not yet tested on a repo: that is the gate backlog item.
1. `pause all publishing` is off; mission scope is personal from the registry; kind is `build` or `repair`.
2. Branch is the mission branch; it rebases on the default branch without conflicts.
3. Commit range = merge base..head. Every subject passes `crow-commit-msg-check`.
4. Secret scan over the range: the `crow-prepush-check` deny list plus a secret scanner (gitleaks over the range when installed). GitHub push protection (default on for public repos of user accounts) is the server-side backstop, not the gate.
5. Sensitive paths always ask: CI workflows, sandbox and egress config, systemd units, kind and policy files, the merge gate itself, `.claude/` settings and hooks, `.mcp.json`, env files, git hooks and deny lists. Global list in core, extra globs per repo in the registry.
6. Diff size limit; above it the mission asks instead of merging (limits are a decision below).
7. Full-suite receipt: core runs the registry's test command on the exact head SHA in a fresh isolated workspace with network off, and stores `{sha, command, exit, duration}`. A push to the default branch without a receipt for that SHA is refused.
8. Green CI on the head SHA when the repo has CI (`gh pr checks`); a repo without CI relies on the receipt.
9. core merges (squash through the forge, or fast-forward push when the repo has no PR flow) and logs `publish`.
10. Deploy hook from the registry runs after the merge; a failure logs an incident and opens a `repair` mission (2 rounds max, then ask).

## Work rules through core (Need 5)
- A publish action is one of: push, MR create or update, tracker write, wiki draft or publish, mail draft, chat post. Work merge is refused outright.
- core shows the exact payload (branch and head SHA, MR title and body, ticket fields, page diff, message text). The owner's `ok` writes a validation record `{mission, version, action, target, payload_sha256, owner_event, ts}` into the event log, chained by hash so an old record cannot be replayed (pattern from the roadmap review).
- At publish time core recomputes the payload hash; a mismatch or a newer owner message cancels the record. Mail stays a draft whatever the record says (owner rule).
- Today's bridge until core exists: `crow-hub-guard` in hub sessions, the autonomy table, and item_077 runs without push in their allowlist.

## Escape hatch (Need 6)
- Isolated runs are headless (`claude -p` inside bwrap): there is no terminal to attach to. "Take control" therefore means: core stops the run, the mission enters a `manual` sub-state, and the owner gets an interactive session in the same worktree (the `crow open` / `handoff` path), outside the sandbox under the owner's own permission mode.
- core logs `manual_gesture` events at start and end with the worktree head before and after. Commits made in that window are tagged as owner work in the mission history.
- Returning the mission to auto starts a fresh isolated run with the brief and the owner's notes (the RESUME NOTE pattern), since `--continue-run` cannot be isolated.
- Open: whether owner commits made in take-control still pass the personal gate (decision below).

## Open decisions (for the owner)
1. **Network for build missions.** Inference-only egress (needs a cdx change for local runs), public internet minus LAN (cdx today), or network off with a host-side dependency cache. Recommendation: public minus LAN now, inference-only later for work-scope runs.
2. **Who commits in isolated runs.** core commits the worktree diff after the run (cdx as is), or ask cdx for an opt-in writable git index. Recommendation: core commits; it already owns publishing and the message check.
3. **Resume of isolated missions.** Fresh isolated run from brief + answers, or ask cdx to allow `--continue-run` with isolation. Recommendation: fresh run now; raise the cdx ask only if resume quality suffers.
4. **Install pasta (passt).** Needed for any online isolated run. A system package the owner installs. Recommendation: install it, then rerun the probe with pasta.
5. **Diff size limit for auto-merge.** 200, 400 or 800 changed lines (and a file count). Recommendation: 400 lines, 20 files, tune from the first month.
6. **Take-control commits and the gate.** Gate still runs on owner commits, or owner commits merge as manual work. Recommendation: the gate still runs; sensitive paths then ask as usual.

# Acceptance criteria
- AC1: A threat × option table with a verdict per cell.
- AC2: A recommended design split into run config, cdx, core and OS, with credential ownership.
- AC3: A prototype run proves an agent cannot push, read the owner's credentials or reach the LAN, while core publishes its branch.
- AC4: The auto-merge gate written as a checklist core can run, tested on one personal repo.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_019_release_leverage_cdx_0_25_logics_2_24.md`
- `logics/architecture/adr_001_core_missions_and_inbox_architecture.md`

# Backlog
- Core publishes and holds every forge credential: `item_099_core_publishes_and_holds_every_forge_credential` (Draft)
- Isolated mission runs through cdx workspace isolation: `item_100_isolated_mission_runs_through_cdx_workspace_isolation` (Draft, needs pasta)
- Personal auto-merge gate for mission branches: `item_101_personal_auto_merge_gate_for_mission_branches` (Draft)
- Mission guard hook for runs that cannot be isolated: `item_102_mission_guard_hook_for_runs_that_cannot_be_isolated` (Draft)
- Take control of a mission as a logged manual gesture: `item_103_take_control_of_a_mission_as_a_logged_manual_gesture` (Draft)
- Deny-by-default background runs: `item_077` stays the run-config layer, unchanged by this spike.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
