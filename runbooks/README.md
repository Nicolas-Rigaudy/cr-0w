# Runbooks

Short, tested procedures for things the hub does often. Read the matching runbook
before exploring; update it when reality differs. One file per procedure.

- `runbooks/`: generic, publishable (no names, companies or personal paths).
- `local/runbooks/`: private specifics (accounts, repos, people), indexed in `local/runbooks/README.md`.

Promote a runbook to a skill (`skills/` or `.claude/skills/`) once the user would
want to trigger it by name.

## Index
- [install-and-update](install-and-update.md): engine vs data dir, install, `crow update`, restore on a new machine
- [dispatch-background-runs](dispatch-background-runs.md): start, track and recover cdx background runs in project repos
- [review-findings](review-findings.md): read-only review run of one MR, its findings JSON schema, resume for the apply phase
- [cost-audit](cost-audit.md): measure token spend by context size, model, project and skill; re-measure after a change
- [archeotech-crow-agents](archeotech-crow-agents.md): test the crow agents bar/panel module with fixtures and headless renders
- [archeotech-crow-usage](archeotech-crow-usage.md): test the crow usage bar/panel module (fixtures, render script, sandbox gotchas)
- [sub-project-work](sub-project-work.md): hand repo work off with `crow handoff`, or `crow work` for lookups; the hub guard hook
- [open-project-session](open-project-session.md): open a project session in a new kitty window from the hub (HOME + inherited Claude markers)
- [session-retros](session-retros.md): daily retro batch (Mon–Fri timer → no-model pre-filter → one /retro run → dated report with skipped list)
- [agent-tooling](agent-tooling.md): third-party agent tools installed (code graph, graphify, ast-grep, claude-hud), pinned versions, verdicts, undo
- [service-tokens](service-tokens.md): personal API tokens in the keyring via `crow secret` (store, read without printing, rotate, revoke), registry in local/tokens.yaml
- [ci-secrets](ci-secrets.md): prod-scoped CI variables and secrets (owner runs the generator script; scope/protect rules)
- [interrupted-sessions](interrupted-sessions.md): quota warning before a usage limit, `crow rescue` for every session it stopped, relaunch one by id
- [new-account](new-account.md): bring a fresh cdx profile up to crow parity (memory links, env, hooks, skills, MCP); `crow doctor` checks it
- [memory](memory.md): memory lives in local/memory, linked into allowed cdx profiles (never copied)
- [design-canvas](design-canvas.md): add or change boards in a Design canvas artifact without losing the owner's layout (live index merge, one file per publish)
- [visual-concepts](visual-concepts.md): purpose, content, flat sketches, then image-model prompts before any render run for a visual world
- [tui-render](tui-render.md): render a TUI (or ANSI mockup) to PNG in the kitty theme at real monitor sizes; fallback-glyph and contrast checks
- [afk-runs](afk-runs.md): unattended AFK runs (night or day): queue, guardrails, pre-push denylist, return report
- [repo-tidy](repo-tidy.md): daily repo hygiene via `bin/crow-tidy` (safe auto-cleanup in the brief) and the by-hand cleanup steps
- [own-routines](own-routines.md): write your own daily/weekly-style skills in local/hub-skills from the generic examples
- [readme-showcase](readme-showcase.md): method for a marketing README per repo (visuals, header options, structure, audit, preview)
- [claude-and-codex](claude-and-codex.md): what changes on a Codex account (instructions, skills, permission, models) and what the hub guard covers
- [long-sessions](long-sessions.md): context cap at 150k, warning at 120k, when to wrap and how to split one session into scoped handoffs
- [vault-capture](vault-capture.md): `crow note`, `crow notes` JSON, the `eod` skill and the capture tags that feed every recap
- [tui-review-ux](tui-review-ux.md): review-tool UX research for the deck MR page: stage rail, whose-move rule, badge options, TUI constraints
- [tui-home-ux](tui-home-ux.md): home dashboard UX research: triage table, day timeline, project cards
- [tui-agents-ux](tui-agents-ux.md): agent monitor UX research: process table, day lanes, inbox-first
- [tui-projects-ux](tui-projects-ux.md): project portfolio UX research: ledger, runway gantt, who-acts lanes
- [tui-usage-ux](tui-usage-ux.md): token and cost usage UX research: ledger, runway timeline, account cards
- [tui-themes-ux](tui-themes-ux.md): TUI theme research: axes a theme can change, four skeleton-level theme specs
- [open-threads](open-threads.md): the three-section format (me, crow, waiting) of `open-threads.md`, the done and stale rules, who reads and writes it
- [graph-today](graph-today.md): Entra app + device-code login for `bin/crow-today`, the Graph calendar and mail feed of the deck Today box
- [releases](releases.md): VERSION, per-release changelogs, the aggregate and the version-bump gate, `crow release propose`
