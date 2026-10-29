# Runbooks

Short, tested procedures for things the hub does often. Read the matching runbook
before exploring; update it when reality differs. One file per procedure.

- `runbooks/`: generic, publishable (no names, companies or personal paths).
- `local/runbooks/`: private specifics (accounts, repos, people), indexed in `local/runbooks/README.md`.

Promote a runbook to a skill (`skills/` or `.claude/skills/`) once the user would
want to trigger it by name.

## Index
- [dispatch-background-runs](dispatch-background-runs.md): start, track and recover cdx background runs in project repos
- [cost-audit](cost-audit.md): measure token spend by context size, model, project and skill; re-measure after a change
- [archeotech-crow-agents](archeotech-crow-agents.md): test the crow agents bar/panel module with fixtures and headless renders
- [archeotech-crow-usage](archeotech-crow-usage.md): test the crow usage bar/panel module (fixtures, render script, sandbox gotchas)
- [sub-project-work](sub-project-work.md): hand repo work off with `crow handoff`, or `crow work` for lookups; the hub guard hook
- [open-project-session](open-project-session.md): open a project session in a new kitty window from the hub (HOME + inherited Claude markers)
- [session-retros](session-retros.md): automatic retrospective after every session (SessionEnd hook → queue → /retro on Sonnet)
- [agent-tooling](agent-tooling.md): third-party agent tools installed (code graph, graphify, ast-grep, claude-hud), pinned versions, verdicts, undo
- [ci-secrets](ci-secrets.md): prod-scoped CI variables and secrets (owner runs the generator script; scope/protect rules)
- [interrupted-sessions](interrupted-sessions.md): list sessions cut off by a usage limit and relaunch them on another allowed account (always by id)
- [memory](memory.md): memory lives in local/memory, linked into allowed cdx profiles (never copied)
- [design-canvas](design-canvas.md): add or change boards in a Design canvas artifact without losing the owner's layout (live index merge, one file per publish)
- [tui-render](tui-render.md): render a TUI (or ANSI mockup) to PNG in the kitty theme at real monitor sizes; fallback-glyph and contrast checks
- [afk-runs](afk-runs.md): unattended AFK runs (night or day): queue, guardrails, pre-push denylist, return report
- [repo-tidy](repo-tidy.md): daily repo hygiene via `bin/crow-tidy` (safe auto-cleanup in the brief) and the by-hand cleanup steps
- [own-routines](own-routines.md): write your own daily/weekly-style skills in local/hub-skills from the generic examples
- [readme-showcase](readme-showcase.md): method for a marketing README per repo (visuals, header options, structure, audit, preview)
- [claude-and-codex](claude-and-codex.md): what changes on a Codex account (instructions, skills, permission, models) and what the hub guard covers
- [long-sessions](long-sessions.md): context cap at 150k, warning at 120k, when to wrap and how to split one session into scoped handoffs
