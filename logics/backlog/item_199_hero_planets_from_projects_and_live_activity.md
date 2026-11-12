## item_199_hero_planets_from_projects_and_live_activity - Hero planets from projects and live activity
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Build hero planets from the project registry plus live sessions and dispatch agents, with missions and asks on top.
- Keywords: hero, planets, projects, dispatch, sessions, model.ts
- Use when: building req_059 step 1.
- Skip when: deciding how many planets fit (item_200).

# Problem
- `heroModel` only sees missions. A live session in a repo with no running mission leaves its project off the table or shows it done.

# Scope
- In:
  - Feed `projects`, `dispatch` and `sessions` into the hero (they are already polled for other screens; reuse that, check `src/data` first).
  - `heroModel(missions, inbox, live, now)`: one `HeroProject` per registered project name; live agent match by cwd under `path`, else by `projects[].name` with the largest share.
  - State precedence: failed > needs > working > queued > done > sleeping (no activity).
  - A live agent planet carries `runFrom` (agent start or oldest recent action) so the label shows the run time.
  - Unit tests and fixture updates (`fixtures/states/live.json`), reel still plays.
- Out:
  - Cap and layout of many planets (item_200).

# Acceptance criteria
- AC1: Live core during a hub session: cr-0w planet is working, label shows a run time.
- AC2: A failed mission on a project with a live agent shows failed.
- AC3: An inbox ask on a project shows needs with the wait time.
- AC4: `npm test` and `tsc` pass; `scripts/hero-reel.mjs` renders all beats.

# AC Traceability
- request-AC1 -> AC1. Proof: live screenshot.
- request-AC2 -> AC2, AC3. Proof: model unit tests.
- request-AC4 -> AC4. Proof: tests and reel.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_059_nest_hero_planets_from_projects
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: The hero misreports live work; side project.

# Notes
- Build in a fresh session from this doc. Read `local/owner-prefs.md` and `runbooks/nest-screenshots.md` first.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

