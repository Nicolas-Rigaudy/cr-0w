# Agents tab: UX research and design directions

Scope: tab `2 Agents` of the crow deck. Inputs: audit findings A1 to A3 and X1 to X8, spec section 5.2 and 5.6, req_011. Code read: `tools/crow-deck/agents.go`, `history.go`, `history_data.go`, `live.go`, `fixtures.go`.

## 1. What the tab is for

Jobs to be done, in the owner's terms.

1. See at once which agents need me, which work, which are paused, which finished today.
2. Answer a waiting agent without leaving the deck (`a x e c`, `enter`).
3. Read the tail of any agent to judge if it is on track.
4. Know who worked when today: which project, which account, how long, with the full title.
5. Find an old session by words, project, account, date or model, and resume it, hand it off or retro it.
6. Keep the background runs budget visible (`runs 1/2`) and relaunch a failed one.

Data now: `agentRow` (state, project, account, title, doing, since, ask, tail, ctx), `whoRow` (name, account, on and off hour spans), `histSession` (id, project, account, model, task, start, took, cost, tokens, messages, branch, end state, touched, commits, linked, last messages). Sources: cdx sessions, `crow runs`, `crow-usage --json`.

Known gaps from the audit: list labels cut at 17 chars next to a wide empty area (A2), who-worked labels cut and no project column, 10 hour ruler with 2 to 6 cell bars (A1), about 30 blank rows (X3), no statement of who must act (X2).

## 2. Apps studied

- claude-squad: one list of sessions on the left, preview and diff as two tabs on the right, one-key menu at the bottom, pause and resume. Take: preview/diff tab switch on the selected row, always-visible key menu. https://github.com/smtg-ai/claude-squad
- ccmanager: three explicit states per session, Busy, Waiting, Idle, and per project counts `[active/busy/waiting]`, recent projects first, vi-style `/` filter. Take: the state vocabulary and the per project counters in group headers. https://github.com/kbwo/ccmanager
- k9s: `:resource` switch, `/` regex filter with invert `/!`, `l` for logs, columns with age. Take: filter syntax with negation, one key to the log of the row. https://k9scli.io/topics/commands/
- lazydocker: left stack of panels, main pane shows logs or stats of the selected item, ASCII graphs, logs limited to a recent window by default. Take: main pane switches content (logs, stats) with tabs, short default window. https://github.com/jesseduffield/lazydocker
- btop: process list with sort cycle, tree view with collapsed parents, filter, a detail box for the selected process with a small graph, list can follow the selected process, pause. Take: sort key, tree for run to child sessions, detail box with a sparkline per agent. https://github.com/aristocratos/btop
- htop: header row of columns with the sort column highlighted, colored state letter, tag with space, F-keys shown in a bottom bar. Take: sticky column header, sort marker, tagging to act on several rows. https://htop.dev/
- sesh: merges several sources (live, configured, recent dirs) into one picker, a glyph per source, `#digit` to jump to row, exact alias wins over fuzzy. Take: one history picker over live plus past sessions with a source glyph, number jump. https://github.com/joshmedeski/sesh
- Conductor: one workspace per agent, status at a glance, then review and merge from the same app. Take: row ends in the next action (review, merge) not only a state. https://www.conductor.build/
- GitHub Actions: run list with a status icon (running, success, failure, cancelled), job graph, failed steps auto-expanded, log search box. Take: auto-expand the part that failed or asks, log search. https://docs.github.com/en/actions/monitoring-and-troubleshooting-workflows/monitoring-workflows/using-workflow-run-logs
- Grafana state timeline: one swimlane per entity, regions of equal state merged into one bar, legend for state colors, hover for details. Take: merge adjacent identical states, one lane per agent, state color per segment (working, waiting, paused). https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/state-timeline/
- tig: split view, log on top, diff below, `/` search inside the view, Enter opens the detail in a split. Take: list on top, preview below in portrait. https://www.mankier.com/7/tigmanual
- Timewarrior: `day` and `week` reports draw an hour ruler with one bar per tracked interval, tags on the line, totals per row and per day. Take: label with tags, total hours at the line end, day total at the bottom. https://timewarrior.net/docs/

## 3. Principles

1. State first, then actor. Every row starts with one state glyph and ends with who must act (`you`, `agent`, `wait`).
2. Group by who must act, not by technical state. `needs you` on top, always visible, never scrolled away.
3. Full titles. The title column takes all the width left after fixed columns (project, account, age). Cut only the last column, never the title.
4. One list, one preview. The tail is the main pane. The question card is part of the tail, with the answer keys under it.
5. Timeline lanes carry the same label as the list (project and full title) so the eye needs no lookup. Ruler fits the day's real span, not a fixed 10 hours.
6. Merge equal neighbours in lanes (Grafana), show the total hours at the lane end (timewarrior).
7. No empty area. Leftover rows go to more history, a longer tail or the day total, never padding.
8. Filter language is shared: `/` free text, `project:`, `account:`, `model:`, `!` negation (k9s). Live and history use the same bar.

## 4. Design directions

### A. Table-first: "Process table"

One line: a wide htop style table of all agents, with a bottom split of tail and a day timeline in the same grid.

Layout at 235x62.
- Region 1, table, full width, rows 3 to 26 (24 rows). Sticky header: `S  who  project  title  account  model  ctx  tokens  since  next`. Group header rows inside: `needs you 2`, `running 3`, `background runs 1/2`, `finished today 6`. Title column flexible (about 110 cols), the rest fixed. Last column is the actor and action (`you: answer`, `agent`, `wait: MR review`). Sort marker on the header, `s` cycles sort.
- Region 2, bottom left, cols 1 to 140, rows 27 to 61: tail of the selected agent. Header line with title, account, ctx bar. Question card at the bottom with `a x e c g` keys.
- Region 3, bottom right, cols 141 to 235: detail box (btop style): task, branch, touched files, linked MR, tokens and a sparkline of tokens per 10 minutes, commits. Then a compact day strip of this agent only.
- No who-worked-when box. The strip lives in the detail box and in `w` full screen overlay.

At 118x60: table on top (rows 3 to 24, `next` and `model` columns dropped, title keeps the width), tail below (rows 25 to 59), detail box becomes tab `d` of the tail pane (claude-squad style tab). `w` overlay for the timeline.

Data: table from `agentRow` (all fields exist; `Next` actor column is new, derived from state and Ask). Model and tokens from `crow-usage --json` joined on session id. Sparkline: new, bucket the transcript timestamps. Detail: `histSession.Touched/Commits/Linked` exist for history rows, needs the same join for live rows.

Inspiration: htop, btop detail box, claude-squad tabs, ccmanager state names.

### B. Timeline-first: "Day lanes"

One line: the day is the page. Each agent is a swimlane on an hour ruler, the list and the label are the same thing.

Layout at 235x62.
- Region 1, lanes, full width, rows 3 to 30 (up to 28 lanes, grouped by project with a header row each). Lane = `S project  full title  acct` (about 70 cols, fixed to the longest title up to 60) then the ruler area (about 150 cols). Ruler spans first activity to now plus 1 hour, one tick per hour, `▼` now. Segments: `█` working in green, `▓` waiting for you in yellow, `░` paused in ANSI 8, `╌` idle gaps. Adjacent equal states merged (Grafana). At the lane end: total hours (`2h40`). Selected lane highlighted as the one filled row. A waiting lane ends with a blinking-free `?` in yellow at the now marker.
- Region 2, rows 31 to 33: totals row per account (`work 5h used 62%`, `perso 4h 18%`) and day total hours, summed with `unionHours`.
- Region 3, rows 34 to 61: tail of the selected lane, left 60 percent, question card under it. Right 40 percent: history of the same project today (previous sessions, titles, took), so the lane has its context.
- Cursor on a lane moves a vertical guide on the ruler and shows the time under it in the box footer.

At 118x60: lanes keep the label (50 cols) and shrink the ruler to 60 cols (hour ticks every 2 hours). Tail below, project history hidden behind `h`. Under 100 cols the lane list and the tail are two panes switched by `tab`.

Data: lanes from `whoRow` (On and Off spans exist). Needs project and full title joined from `crow-usage` sessions (the audit says the data is there). Waiting spans (state `needs`) are new: derive from the transcript when the last assistant message asked and no user message followed. Account totals from the Usage source.

Inspiration: Grafana state timeline, timewarrior day report, GitHub Actions job graph, tig log plus detail.

### C. Inbox-first: "Needs you and the rest"

One line: a triage queue. Left is a short ranked list of what needs the owner with the question inline, right is the tail, and the history is a first-class column instead of a hidden mode.

Layout at 235x62.
- Region 1, left column, cols 1 to 95, rows 3 to 30: "Needs you N" cards, 3 rows each: line 1 state glyph, project, full title, age; line 2 the agent's question in full (wrapped to 2 lines at 90 cols); line 3 keys `a keep  x drop  e edit  c chat`. Selected card expands to show the last 6 tail lines.
- Region 2, left column, rows 31 to 61: "Working and waiting on others" compact one line rows (`S project title account since actor`), then `background runs 1/2` with the run command label, then `finished today` rows with result and cost. Group headers carry ccmanager style counters.
- Region 3, right column, cols 96 to 235, rows 3 to 36: tail of the selected agent, tabs `tail | diff | detail` (claude-squad).
- Region 4, right column, rows 37 to 61: history, always on: search bar `/` on one line, then sessions grouped by day (today, yesterday, older) as one line each (`time project full title took cost`). Selecting one fills the tail pane with its preview (task, end state, touched, last messages), `enter` resumes. No separate mode and no `live / history` switch.
- Day timeline: one compact line per account under the header (`work ██▓▓░░██  09 to 14`), 2 rows total, replacing the big box.

At 118x60: one column. Needs you cards (rows 3 to 20), working list (21 to 34), tail with tabs (35 to 50), history one line rows (51 to 59). `tab` rotates focus between the four boxes. Hidden: nothing, only lines per box shrink.

Data: cards from `agentRow.Ask` and state. Actor column new but trivial. History from `histSession` (exists, same fetcher as today). Account timeline from `whoRow` grouped by account (new grouping only).

Inspiration: ccmanager Waiting state, Conductor (row ends in the next action), sesh (one picker over live and past), k9s filter syntax, GitHub Actions auto-expanded failed step.

### How they differ

| | A Process table | B Day lanes | C Inbox |
|---|---|---|---|
| Main object | the agent row | the time span | the question to answer |
| Best at | scanning many agents, sorting | seeing who worked when | clearing what blocks the owner |
| Weakest at | time view (overlay only) | many short sessions (lane count) | overview of 20+ agents |
| New data | actor, sparkline | waiting spans, joined titles | actor only |
| History | `/` mode | detail box and project panel | always-on box |

Suggested pick to compare first: B for the timeline complaint (A1), C for the "who must act" rule (X2). A is the safest if the agent count stays under 15.
