# Crow deck Home tab and themes: UX research

Scope: the Home tab of `crow deck` and the theme system. Written 2026-11-02 for req_011. Inputs: the audit (H1-H5, X2-X8), the design spec, `tools/crow-deck/home.go`, `data.go`, `chrome.go`, `theme.go`, `frame.go`.
Fetch notes: some vendor pages were thin (k9s pulses, Things, Todoist Today). Claims below are limited to what the pages said. Anything else is marked as a design choice, not a fact about the app.

## 1. What the tab is for

Home is the one screen read at the start of the day and re-read between tasks. It answers five questions in this order.

1. What needs me now, and who must act? Next to-dos, blockers (MRs, Jira), MRs waiting. Every row names its actor: me, an agent, someone else.
2. What do the agents do right now? State, project, account, what they do, since when, and which one is asking me something.
3. What does today look like? Calendar, mail counts, drafts due, shared mailbox.
4. Where do the projects stand? Roadmap progress with a real target date. The prod launch must always be visible.
5. How much budget is left? Per account: 5h and weekly gauges, policy, burn, runs x/2. Gauges say "left" or "used" in words.

Jobs, in the owner's terms: unblock others first (reviews), keep the launch on track, spot a stuck or asking agent in under three seconds, never open another tab to find out what to do next.
Home is read-mostly. Its actions are small: open, hand off, review, snooze, relaunch. It must not become a second copy of the other tabs.

Data today: `fixtures` (Gauges, Agents, Blockers, Today, Mail, Next, Projects, repo counters), fed by cdx, `crow runs`, forges, open-threads, registry. Known gaps from the audit: Next has no actor or age (req_012), Blockers has no Jira (`bin/crow-jira json`), Today has no source (`bin/crow-today`), Projects lacks lambda and dates.

## 2. Apps studied

- **btop**: separate boxes for cpu, mem, net, proc, each with a numbered hotkey on its title, boxes can be hidden or moved by config, meters or graphs per value. Take: a header strip of gauges that costs 3 to 4 rows and numbered boxes. Take the TTY theme idea: a theme that only uses the 16 terminal colors. https://github.com/aristocratos/btop
- **htop**: header meters (cpu bars, memory, tasks, load) above one dense table. Take: gauge rows in the header, each with a label, bar and number in one line, then the main table gets all remaining height. https://htop.dev
- **k9s**: pulses view gives a cluster overview (`:pulses`), and skins have separate colors for border and focused border, plus title, counter and filter colors. The pulses page itself was too thin to cite details. Take: counters in titles and a skin with explicit border and focus roles. https://k9scli.io/topics/commands/ and https://k9scli.io/topics/skins/
- **glance**: pages as tabs, each page split into columns, each column `small` or `full`, widgets stacked in a column, `group` puts several widgets in one slot. Take: a narrow summary rail beside a wide main column, and a group slot (our Blockers and Today share one). https://github.com/glanceapp/glance
- **gh-dash**: user-defined sections of PRs and issues, one table per section, configurable columns, vim keys, a theme with text, background and border colors (primary, secondary, faint). Take: sections as filtered tables with an "awaiting my review" section first, and the three-level border palette. https://github.com/dlvhdr/gh-dash and https://dlvhdr.github.io/gh-dash/configuration/theme
- **wtfutil**: module grid with `columns` and `rows` as sizes and each module placed by top, left, height, width; three border colors (`focusable`, `focused`, `normal`). Take: the three-state border (can focus, has focus, static) so read-only boxes look different from selectable ones. https://wtfutil.com/configuration/files/
- **sampler**: components runchart, sparkline, barchart, gauge, textbox, each with a title; conditional triggers raise an alert when a value passes a limit. Take: sparkline plus gauge pairs for burn, and a trigger idea (a gauge below 15% turns the row yellow and adds a badge). https://github.com/sqshq/sampler
- **Linear inbox**: a Priority tab for what needs attention and an Other tab for the rest, `J K` to move, snooze on `H`, mark read on `U`, reminders resurface items. Take: a split between "needs you" and "FYI", and snooze with a return time. https://linear.app/docs/inbox
- **GitHub dashboard and "My issues"** style views: not fetched in detail. Design choice only: group by "assigned to you", "review requested", "waiting on others". Used here as the model for the actor column. https://github.com
- **Things and Todoist Today**: Things has a This Evening list and headings inside a project (URL scheme page only). Todoist keeps overdue, all-day and time-blocked tasks in a Plan sidebar next to the calendar, and lets you push a task to tomorrow. Take: overdue first, calendar and tasks in one time-ordered list, and a one-key postpone. https://culturedcode.com/things/support/articles/2803573/ and https://www.todoist.com/help/articles/plan-your-day-with-the-today-view-UVUXaiSs
- **lazygit**: theme keys `activeBorderColor`, `inactiveBorderColor`, `selectedLineBgColor`, `optionsTextColor`, border style `rounded`, `single`, `double`, `hidden`, `bold`. Take: focus is carried by border color plus weight, and `reverse` for the selected line. https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md
- **helix**: themes use the 16 ANSI names by default, `inherits` builds one theme on another, `ui.window` styles the lines that separate splits, `ui.statusline.inactive` dims unfocused panes. Take: a theme as a small override of a base, and a separate role for separators. https://docs.helix-editor.com/themes.html

## 3. Principles

1. One list answers "what do I do next". It merges to-dos, blockers and waiting MRs, sorted by who must act then by age. Everything else supports it.
2. Every row carries an actor glyph and word: `@me`, `@agent`, `@wait <name or system>`. Colour follows actor: yellow me, green agent, ANSI 7 wait.
3. Show age on every item that waits. Age turns yellow past one day and red past three.
4. Gauges say what they measure: `left` or `used`, with the reset time. One line per account.
5. Fixed columns only for short fields (glyph, project, age). The title column takes all remaining width and wraps to a second row before it truncates.
6. Leftover height adds rows to the lists in priority order. No box is taller than its content, no empty boxes. A source that is not live is not drawn.
7. Counts in titles (`Needs you 7 · 3 me`), destination in the right label (`→ 3 MRs`), cursor position bottom right.
8. Home links out, it does not duplicate. A row opens its tab filtered on that item.

## 4. Design directions

Notation: `W` is columns, `H` is rows. Chrome (header, tabs, rule, status bar) takes 4 rows, so the body is 58 rows at 62 and 56 at 60 in portrait. All boxes include their border rows.

### A. Triage table ("one list to act on")

Idea: one wide unified table of everything that needs action, grouped by actor, with a rail of context on the right.

Layout at 235x62:
- Row band 1, full width, 4 rows (2 content): Pulse. One line per account (cdx accounts, usually 2): `acct 5h ━━━━━━━━ 72% left, resets 17:00  wk ━━━ 41% left  ask-first · burn 3%/h · runs 1/2`. Right end: a one-line health text.
- Row band 2, left column 160 cols, 36 rows: **Needs you** table. Columns: `◆` glyph 1, actor 7 (`@me`, `@agent`, `@wait`), project 22 (full repo name), kind 6 (`MR`, `JIRA`, `TODO`, `ASK`), title flexible (full, wraps to row 2), who or system 14, age 5. Group header rows: `you decide 5`, `agents work 3`, `waiting on others 4`, each foldable with `space`. A selected row shows a one-line preview under it: the question an agent asks, the MR title, the Jira status.
- Row band 2, right column 75 cols: three stacked boxes. **Right now** (agents, about 24 rows: glyph, project, account, doing, since). **Today** (calendar then a mail counts line, about 30 rows, time-ordered, next event on the accent).
- Row band 3, left 160 cols, rest of height (about 14 rows): **Projects** table: repo name, scope, state, road bar `6/8`, target date, days left, at risk flag, MRs, threads, git dirty/ahead.
- Right column keeps Right now and Today through band 3, so they get the full height and show more rows.

At 118x60: one column. Pulse (4 rows), Needs you (about 24 rows, project column drops to 14, kind column merges into the glyph), Right now (10), Today (8), Projects (rest, columns: name, road bar, target, flags). Order is the same as landscape reading order.

Data per region:
- Pulse: existing `fx.Gauges` (crow usage, accounts policy, runs).
- Needs you: existing Blockers (forges), plus new actor field from the to-do model (req_012: `@me|@agent|@wait`), plus Jira from `bin/crow-jira json`, plus agents in state `needs` from cdx.
- Right now: existing `fx.Agents`.
- Today: new `bin/crow-today` (Graph), shown only when its file exists.
- Projects: existing registry plus new `release` for lambda and `target:` dates in roadmap frontmatter.

Inspiration: gh-dash (section tables, awaiting-review first), Linear inbox (Priority versus Other, snooze), htop (header meters, one dense table), glance (rail beside main column).

### B. Day timeline ("the day as one ruler")

Idea: time is the spine. A horizontal day ruler shows meetings, agent activity and deadlines; the action list sits under the moment it relates to.

Layout at 235x62:
- Band 1, full width, 12 rows: **Today** timeline. Ruler 08 to 19 with `▼ now`. Lanes: `calendar` (event blocks as `▕ standup ▏` ranges with titles inside, full width of the span), one lane per active agent (`●` working stretches, `‖` paused, `◆` asking), `deadlines` lane (`◆` prod launch countdown, MR due dates, Jira due). The lane for each account carries its gauge at the left edge: `acct 72% ━━━━` in a 24 col gutter. Mail counts at the right end of the calendar lane label.
- Band 2, 3 columns of 78/78/79 cols, 30 rows: **Needs me** (the unified actor table, no wait group), **Waiting on others** (MRs awaiting review by someone, Jira in review, support mailbox, with age and a nudge key `n`), **Agents** (full list: state, project, doing, since, ask text on a second row).
- Band 3, full width, rest (about 16 rows): **Projects** as a mini timeline. One row per project: name, bar from start to target with `▼ today`, `◆` milestones, `6/8`, and days left in the accent when under 30.

At 118x60: ruler band shortened to 08 to 19 at 1 col per 15 min (44 cols) beside lane labels, 10 rows. Needs me, Waiting, Agents stack as three boxes of equal height (about 12 each). Projects mini timeline drops to bar plus date columns. The `deadlines` lane becomes a single line under the ruler.

Data per region:
- Timeline calendar lane: new `bin/crow-today` events.
- Agent lanes: existing "Who worked when" data (`whoRow` with `On`, `Off` spans from sessions) already used in the Agents tab.
- Deadlines: roadmap targets (new), Jira due dates (`crow-jira`), MR age from forges.
- Needs me and Waiting: Blockers plus to-do model plus Jira; split by actor field.
- Projects timeline: registry roads and targets (new `target:`).
- Gauges in the lane gutter: existing `fx.Gauges`.

Inspiration: Things and Todoist Today (calendar and tasks in one time-ordered view, postpone), sampler (sparkline and gauge components with alert thresholds), the existing Agents "Who worked when" ruler, wtfutil (grid by sizes).

Risk: if the Today source is not wired, band 1 loses its calendar lane and shrinks to agents and deadlines. It must collapse, not show a placeholder.

### C. Project cards ("projects as the unit, with a pulse header")

Idea: the owner thinks in projects first. Each project is a compact card that carries its own needs, agents, MRs and date. A thin header holds budgets and the day.

Layout at 235x62:
- Band 1, full width, 5 rows: **Pulse** in btop style. Left 120 cols: per account gauge lines (5h and week, left and reset, policy, burn sparkline `▁▂▃▅▆▇` for the last 6 h, runs x/2). Right 115 cols: **Today** compact: next three events on one row each, mail counts and drafts due on the last row.
- Band 2, full width, 3 rows: **Prod launch strip**: one line for the top-priority release: name, target date, days left, road bar, at-risk reason, and the blocking items count. Stays on screen always, accent coloured.
- Band 3, rest (about 48 rows): a grid of project cards, 3 across by 2 or 3 down, each 77 cols by 15 to 16 rows. Card rows: title row (full repo name, scope, state glyph), `NOW` (what the active agent does, since), `NEXT` (the top to-do with actor and age), `BLOCK` (MR, Jira, thread with who and age), `ROAD` bar with `6/8 · Dec 1 · 29d`, `OPEN` counters (MRs mine, awaiting me, threads, dirty, ahead). Sorted by urgency: cards with an `@me` item first. The selected card shows `h handoff  r review  c continue` as key text in its bottom border. Cards for repos with nothing open collapse into one line list at the end.
- No separate Next or Blockers box. A `filter` key (`!`) shows only cards that need me.

At 118x60: one column of 2 by N cards is too narrow, so cards go one per row at 118 cols and 7 rows (title, NOW, NEXT, BLOCK, ROAD plus OPEN merged on one row), about 6 cards visible; Pulse and launch strip stay on top; Today drops into the Pulse as one line.

Data per region:
- Pulse and Today: as in A.
- Prod launch strip: registry release (new for lambda), `target:`, at-risk from `crow-registry json`.
- Cards: registry (names, scope, git counters from `crow-activity`), forges (MRs, threads), Jira (new), to-do model (NEXT with actor), cdx (NOW). The card is a join on the project tag, which needs the unregistered tags (`system`, `scenefx`, profile repos) mapped in the registry (audit X6).

Inspiration: glance (widget groups per column, small and full), k9s (counters in titles and skin roles), project dashboards in management tools (card per project with status and progress), btop (header meters and sparklines). The Projects tab keeps its table, so Home cards must stay summaries.

Which direction fits best: A is the cheapest and the most direct answer to H1 and X2, so it is the safest default. B wins if Today (calendar) becomes reliable. C wins if the owner scans by project more than by action. A and B can share the Pulse band and the unified actor table, so choosing A first does not block B.

## 5. Themes

### How other TUIs delimit panels without fills

- lazygit: `activeBorderColor` (bold green default) versus `inactiveBorderColor` (default fg), border styles `rounded`, `single`, `double`, `bold`, `hidden`. Focus uses colour and weight. The selected line uses a background or `reverse`.
- k9s skins: `frame.border.fgColor` and `focusColor`, title has its own `fgColor`, `highlightColor`, `counterColor`, `filterColor`. Title text is coloured, the counter has its own colour.
- btop: every box has an outline colour per box, the title sits on the top border with a hotkey digit highlighted. TTY theme uses the 16 colours only.
- helix: `ui.window` is a separate role for the separator lines, `ui.statusline.inactive` dims the unfocused pane. Separators exist as a role, they are not fills.
- wtfutil: three border colours: focusable, focused, normal.
- gh-dash: three border levels (primary, secondary, faint). Primary lines sit under tabs and between sections.

Rules distilled for a translucent terminal with no fills:
1. Delimit by lines, never by background. A line glyph in ANSI 8 may vanish on some themes, so the unfocused line uses ANSI 7 when the theme has low-contrast 8, or a config switch `dim: 7|8`.
2. Three states: static (not selectable), selectable unfocused, focused. Focus changes weight (`─` to `━`, single to double or heavy corners) and colour, never only colour. Bright slots equal normal ones in some themes, so weight carries the signal.
3. The title sits on the border, coloured, with its counter in a second colour. The focused title adds a marker glyph (`▶` or `◆`) and the hotkey digit.
4. The selected row is reverse video on the accent. Only this row and the mode chip are filled.
5. Between boxes use one blank cell column or a shared border, never a mixed gap. Columns inside a box use `│` in ANSI 8 or 7.
6. The status bar is a line separated by a rule, not a band.

### Four theme directions

1. **Round (lazytf, default)**. Rounded borders `╭─╮╰╯`, title on the top border, magenta title, blue accent. Unfocused border ANSI 8 (fallback 7), focused border accent plus `━`-weight top edge replaced by heavy rounded look: use `┏━┓` corners in the accent for the focused box only, `╭─╮` for the rest. Static boxes (Pulse) use dotted `┄` top and bottom. Fix over today: the focused box is currently distinguished by colour only. Weight change is added.
2. **Soft, fixed ("Rules")**. The current Soft had no lines, so boxes were not visible. Keep the airy look but add three delimiters: a title rule `Needs you 7 ───────────────────────` full width of the box in ANSI 8 (this is the box top), a one-column gutter `│` in ANSI 8 between side-by-side boxes, and a left rail `▌` in the accent on the focused box (all its rows). No box bottom, no side lines. Between stacked boxes one blank row stays, which gives the airy feel. Sentence-case titles, cyan accent, underline on the active tab. The title rule is what restores delimitation: every box starts with a visible line and the rail shows focus.
3. **HUD (brackets)**. Corner brackets only: `┏` and `┛` plus short ticks `╸` at the other two corners, title as `◢ TITLE ◣` on the top edge, panel index `01` on the edge. Edges between corners drawn in ANSI 8 with `╌`. Focused panel: edges become `━` in the cyan accent and the index is bold. Cyan accent. Drops first when narrow: the index, then the `◢ ◣` tags.
4. **Ledger (Cyber, double focus)**. Single rules `┌─┐`, title as `▌// TITLE //` in magenta, table headers prefixed `//`. Unfocused: single line ANSI 8. Focused: double line `╔═╗` in the accent magenta, the only double-line box on screen, so focus is readable at a glance without brights. Corner notch `◢` on the top right. Drops first when narrow: `//` marks, then the notch.

Config: one accent slot per theme (any of 16), an optional `dim` slot (7 or 8) for unfocused lines, a `weight` flag to force heavy or double focus. Themes defined as an override of a base (helix `inherits` idea) keeps the four small. Check each theme on Home and one dense tab (AC2) with a translucent terminal screenshot before choosing.
