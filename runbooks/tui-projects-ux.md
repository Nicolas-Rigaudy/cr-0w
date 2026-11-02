# Crow deck: Projects tab UX research

Date: 2026-11-02. Scope: tab 4 Projects. Inputs: the audit (P1 to P5), the design spec 5.4, req_011, `tools/crow-deck/projects*.go`.
Research limits: several vendor pages returned 404 or only a landing page (Linear updates, taskwarrior docs, Notion rollups, gh-dash docs). Claims below come only from pages that loaded. Where a page was thin, the pattern is marked "known from use".

## 1. What the tab is for

Jobs to be done, in the owner's terms:

1. Answer "what is the state of each repo" in one glance: running, needs me, blocked, ready, quiet.
2. See the next real action per project and who acts: me, an agent, or someone else.
3. Know if the prod launch (Dec 1) is on time, and which project threatens it.
4. See roadmap progress as items done of total, the current milestone, its exit signal, the target date, at-risk flag.
5. See what is open around a project: MRs, review threads, Jira tickets, git state (dirty, ahead, behind, stale worktrees).
6. See dated milestones on a timeline and what comes before the launch.
7. Act from the row: continue, review, research, edit roadmap, open session.
8. Filter by scope (work or personal) without work rows leaking personal ones.

Data today (projects_data.go): per project name, release, scope, state, now, next, block, done/of, target, at-risk, current milestone, exit, open MRs/threads. Wrong today: NEXT comes from the open-threads log, no dates, lambda has no release entry.

## 2. Apps studied

- Linear projects. Projects are grouped by status, lead, target date or activity, with list, board and timeline as display options. Start and target dates can be coarse (quarter, month) or exact. Take: coarse-or-exact dates, group-by switch, a health word per project. https://linear.app/docs/projects
- GitHub Projects roadmap. Items sit on a time axis from start and target date fields, zoom month/quarter/year, vertical markers for iterations and milestones, group by a field. Take: zoom levels, markers, grouping rows by a field. https://docs.github.com/en/issues/planning-and-tracking-with-projects/customizing-views-in-your-project/customizing-the-roadmap-layout
- Jira Advanced Roadmaps (plans). One plan merges several boards, shows potential dependencies, schedules by capacity, and keeps changes in the plan until saved. Take: cross-project dependency lines and a visible "this threatens the date" warning. https://support.atlassian.com/jira-software-cloud/docs/what-is-advanced-roadmaps/
- Notion databases. One dataset, many views (table default, list, board, calendar, gallery), with filter, sort and group per view. Take: the same rows as table or lanes by a keypress. https://www.notion.com/help/intro-to-databases
- taskwarrior-tui. Keyboard-first list, report switching, fast filter, contexts, details pane. Take: filter-as-you-type and named reports (views) switched by one key. https://kdheepak.com/taskwarrior-tui/ and https://github.com/GothenburgBitFactory/taskwarrior-tui
- gh-dash. User-defined sections per repo, configurable columns, vim keys, a preview pane. Take: sections as saved filters, tab to switch, preview of the selected row. https://github.com/dlvhdr/gh-dash and https://www.gh-dash.dev/getting-started/
- lazygit. List panel plus main view, tab moves between panes, screen modes cycle to enlarge a pane. Take: list left, detail right, one key to grow the detail. https://github.com/jesseduffield/lazygit
- Mermaid gantt. Sections, `milestone` points, `done/active/crit` status tags, a today marker, a configurable axis. Take: crit tag for the path to the launch, one glyph per status, today marker. https://mermaid.js.org/syntax/gantt.html
- Portfolio RAG dashboards (known from use, Linear health, Jira portfolio, Asana portfolios). One coloured status word per project (on track, at risk, off track) plus a one-line reason. Take: RAG as colour slots green/yellow/red plus a text word, never colour alone.
- Taskwarrior burndown and ghistory (known from use, page not reachable). Text bars per period with stacked glyphs. Take: a per-project mini sparkline of items done per week as a velocity hint.

## 3. Principles

1. One row per project, full repo name, never a card that spans five lines.
2. Every row names who acts (`you`, `agent`, `wait`, `none`) in its own column, coloured.
3. Health is a word plus a colour (`on time`, `at risk`, `late`, `no date`) with a one-line reason on the selected row.
4. Progress is `n/m` plus a bar, with the unit in the header (`items`). Never a bare `6/8`.
5. Dates are first-class. A project with a roadmap and no date shows `no date` in yellow, which is a to-do for the owner.
6. The launch gets a countdown (`D-29`) and the projects on its critical path are marked.
7. Detail pane for the selected row holds the long text (milestone, exit criteria, next action in full). The list holds short columns.
8. Leftover height becomes more content: quiet repos get rows, then the timeline grows, then recent project events. No empty boxes.
9. Same data, three views on one key (`v`: table, runway, lanes), like Notion views.

## 4. Design directions

Shared glyphs: `●` running, `◐` needs you, `■` blocked, `○` ready, `·` quiet, `◆` milestone, `▼` today, bars `█▓░` and `━ ─`. Colours by role: green on time, yellow at risk or needs you, red late or blocked, cyan agent, magenta box titles, dim for quiet.

### Direction A: Ledger (table-first)

One idea: a dense sortable table of every repo, a detail pane for the selected one, a thin timeline strip under it.

Layout at 235x62:
- Header row of the tab: filter (`/`), scope chip (`work`, `all`), launch countdown `prod launch Dec 1 · D-29 · 2 projects at risk`. 1 row.
- Left region, `Projects 11` table, 140 cols x about 34 rows. Columns: state glyph (2), repo (22, full name), scope (4), health (9), road bar and `n/m items` (18), milestone (24), target (7), next action in full (34, wraps to 2 rows when long), who (6), MRs/threads (9), Jira (5), git (8: `dirty 3`, `+2 -0`, `clean`). Group rows under three headers: needs you, running or ready, quiet. Quiet repos are one row each, dimmed.
- Right region, `Detail: <repo>` 95 cols x 34 rows: now (agent, since), next action full text, block reasons, current milestone and exit criteria, open MRs (id, title, who must act), open threads, Jira tickets, git (branch, worktrees, stale), action row `c continue  r review  s research  e roadmap  o session`.
- Bottom region, `Runway` 235 cols x 22 rows: month axis Nov to Feb, one line per project that has a date, `▼ today`, `◆` milestones, `crit` rows in red. It grows with leftover height. Below the timeline, a 4-row `Recent` list: last events per project (merge, push, release).

At 118x60:
- Table on top (60 rows share: 30), columns drop to state, repo, health, `n/m`, target, who, next (truncated last, never repo). MRs/Jira/git move to the detail pane.
- Detail pane under the table, 14 rows, full width.
- Runway under that, months compress to weeks of the next 6 only (`z` zoom).

Data:
- Existing: registry (crow-registry json: name, scope, road done/total, current, exit, release target), agents (state, now, since), mrs (open, blockers), `mrCounts`.
- New: `@actor` and type in open-threads (audit gap 1) for next and who; `target:` and milestone dates in roadmap frontmatter; lambda `release` entry; `bin/crow-jira json` for the Jira column; `crow-activity` for git state; Recent from crow-activity.

Inspiration: taskwarrior-tui (dense list plus details, filter), gh-dash (columns, preview), Linear list (health word and target date), lazygit (list plus main view).

### Direction B: Runway (timeline-first)

One idea: the tab is a gantt of the next 90 days. Rows are projects, the axis is time, and the Dec 1 launch is the right-hand wall. Text columns hug the left.

Layout at 235x62:
- Top band, 3 rows: launch countdown bar from today to Dec 1 (`━━━━▼───────◆ Dec 1  D-29`), critical path in one line (`lambda init prod seed ← web tokens ← lint flow`).
- Left gutter, 60 cols: state glyph, repo (full), health word, `n/m`, who. One row per project, 2 rows each when it has a block reason (second row dim, shows the reason).
- Main area, 170 cols: time axis (weeks, `Nov 02 ... Jan 25`, 1 col per half-day at the default zoom), bar per project from start to target: `█` done share, `░` remaining, `◆` milestones with a 1-line label above on the selected row, `▼` today column drawn as `│` through all rows, red `crit` bars, `!` where the bar passes its target. Projects without dates draw a dotted `· · ·` stub at the left with `no date` so the gap is visible.
- Bottom, 14 rows: `Selected` strip split in three: milestone and exit, MRs and threads with who acts, next action and action keys.
- Grouping rows by scope (work, personal) with `g`, or by health.

At 118x60:
- Gutter shrinks to glyph, repo, health, who (34 cols). Axis 80 cols, 1 col per day, the window slides to start 7 days before today. `h l` pan, `z` zoom week/month/quarter.
- Bottom strip stacks vertically under the axis, 16 rows.

Data:
- Existing: registry progress and release target, agents, mrs.
- New and required: start date and target date per roadmap, milestone dates (the whole direction fails without them), dependency edges between projects (frontmatter `after:`), a `crit` flag derived from dependencies on the launch.
- Fallback when dates are missing: compute the bar from done/total and a velocity of milestones closed per week from git history (labelled `est`).

Inspiration: GitHub Projects roadmap (zoom, markers), Mermaid gantt (milestone, crit, today marker), Jira Advanced Roadmaps (dependencies and date warnings), taskwarrior burndown (text bars).

### Direction C: Who-acts lanes (attention-first)

One idea: no table. Four lanes by who must act, so the tab answers "what do I do next" first and "state of repos" second.

Layout at 235x62:
- Four columns of about 58 cols each, 50 rows high: `You` (needs your decision or review, sorted by launch impact), `Agents` (running now, with since and what), `Waiting` (on a reviewer, CI, Jira, a person, with age), `Ready or quiet` (nothing blocks, idle).
- Each project is one tile of 3 rows inside its lane: row 1 `● repo  health  n/m bar`, row 2 next action in full, wrapped, row 3 dim: `!214 your review · 2 threads · Jira 3 · dirty 3`. A project appears in exactly one lane, the one that owns the next move. Lane titles carry counts (`You 3`).
- Bottom, 10 rows: a compact runway, a single line per dated project, same glyphs as Direction B, plus the launch countdown. The lane tiles and the runway share colours.
- `v` swaps lane grouping to health (late, at risk, on time, no date), `scope` filter on `m`.

At 118x60:
- Lanes stack vertically as one list with lane headers (`You 3`, `Agents 2`, `Waiting 2`, `Quiet 4`), same tiles at full width, 3 rows each. Runway shrinks to a 8-row strip at the bottom. Tab key jumps between lane headers.

Data:
- Existing: agents groups (gNeeds, gRunning), mrBlockers (who must act), registry progress.
- New: actor on next actions (`@you|@agent|@wait`), Jira waiting-on-me flag, CI state per MR for the Waiting lane, dates for the strip.

Inspiration: Notion board view (same rows grouped by a property), Linear group by status or lead, gh-dash sections, portfolio RAG dashboards (a worst-first reading order).

### Choosing

- A is the safest and the closest to the audit proposal. It works with today's data and adds columns as sources arrive.
- B is the best answer to the Dec 1 question but needs dates and dependencies first. Build the data work (frontmatter dates, lambda release) in any case, since A and C also use the timeline strip.
- C is the best for daily use when the owner asks "what do I do", but it hides the full portfolio order. It overlaps with the Home Next list, so check duplication before building it.

Recommended: A as the default view, B as the `v` alternate, C only if Home Next does not already cover it.
