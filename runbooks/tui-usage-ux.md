# Usage tab UX research (crow deck, tab 5)

Date: 2026-11-02. Scope: design directions only. Nothing here is built.
Inputs: the deck audit (U1 to U6), the deck spec (section 5.5), `bin/crow-usage --json`, `tools/crow-deck/usage*.go`, the old `tools/crow-watch/usage.go` page.

## 1. What the tab is for

Jobs to be done, in the owner's terms.

1. Can I keep working? Per account: 5h and weekly share left, reset time, burn rate, time until empty.
2. Which account takes the next session? Show the routing policy (use, ask-first, never) next to the gauge, and what each account allows.
3. What costs money? Cost today, week, month. Per session: project, full title, model, tokens, cache read and write, cost, context size warning. Per project and per model.
4. Where did the week go? Daily spend over 7 and 14 days, with the 5h windows that were hit.
5. What happens tonight? Forecast of the empty time per account, planned overnight runs per account, whether they fit.
6. Is it worth it? Savings: cache reads, routing to a smaller model, rtk output filter, scripts that replaced model calls.
7. Who must act? A limit near empty, a session over the context cap, an account that is `never` but still used, a planned run that will not fit. Every flagged row says what to do.

Costs are API-equivalent dollars from the price table, not the bill. The tab must say so once, in a small label.

## 2. Data available today

`crow-usage --json` keys: `limits` (per account: 5h and week percent left, reset times, `burn_5h_pct_per_h`, `empty_at`, `active`), `today`, `week`, `month` (cost, tokens, cache read/write, cache and routing savings, runs, routed runs, compactions, by project, by model, tools), `days` (up to 30 daily rollups), `history`, `sessions` (account, project, title, branch, kind interactive or run, models, tokens split, cache percent, cost, `context_tokens`, `context_level` ok/warn/over, compactions, `limit_share` of 5h and week, `by_day`, per-day `activity` slots, multi-project shares), `rtk` (saved tokens today and all time, percent, commands, failures), `retros`.

Not shown in the deck today: context level, cache split, compactions, tools, routed runs, limit share, month, tokens versus cost, per-account sessions.

New sources needed:
- Planned overnight runs per account: `local/state/afk/` job files. A small `bin/crow-afk-plan --json` that lists job, account, estimated share.
- Per-account policy: `local/accounts.yaml` (exists), edit through `crow-accounts set`.
- Scripts saving: no source. Count script invocations that replace a model call (needs a log line per call in `bin/` helpers).
- Hourly burn history per account for a real sparkline: `crow-usage` returns only the current rate. Add a `burn_series` (hourly percent used, last 24h).

## 3. Apps studied

- ccusage: reports by day, week, month, session and 5h block, a compact table mode under 100 columns, per-model breakdown, project grouping, `--since/--until`, cache creation and cache read as separate columns. Take: one table schema reused for every period, and a compact column set at narrow width. https://github.com/ryoppippi/ccusage
- Claude-Code-Usage-Monitor (claude-monitor): realtime view with progress bars for the session, burn rate from the last hour, predictions of when tokens run out and when the session resets, daily and monthly table views, a model distribution bar, labels that say whether a number is official, estimated or unknown. Take: burn rate next to every gauge, "runs out at" time, and provenance labels on numbers. https://github.com/Maciek-roboblog/Claude-Code-Usage-Monitor
- better-ccflare: accounts with health indicators, priority-based selection, fallback chains ("combos") per model family, rate-limit reason and time recorded, analytics filtered by account, model and status. Take: show the routing order as a visible chain and say why an account is skipped (limited until, policy never). https://github.com/tombii/better-ccflare
- Claude Console Usage and Cost pages (and the Usage and Cost API): usage in time buckets (minute, hour, day), group by model, workspace, key, service tier, context window; token types split into uncached input, cached input, cache creation, output; cost page is daily only. Take: one grouping selector over the same data (by model, project, account), and the four-way token split. https://platform.claude.com/docs/en/manage-claude/usage-cost-api
- OpenAI usage dashboard: cost and activity views, filter and group by project, model, key, with export. Take: same idea, confirms group-by as the primary control. https://developers.openai.com/api/docs/guides/usage-dashboard (page could not be read in full, pattern taken from its title and the Console equivalent)
- Cursor usage page: a meter for the included allowance against on-demand spend, plus an event table. Take: separate "included" from "extra" so a gauge reads as allowance left. https://cursor.com/docs (specific page not retrievable, pattern from the product, treat as unverified)
- AWS Cost Explorer: group by and filter, daily or monthly granularity, a forecast over a chosen future range drawn as a continuation of the history, saved reports. Take: forecast as a dashed continuation of the daily series, and group-by as a single key press. https://aws.amazon.com/aws-cost-management/aws-cost-explorer/
- Vercel usage: an allotment indicator per metric showing the share used in the cycle and the projected cost, a per-project breakdown, view modes (count, per project, ratio, average), alerts and spend management that can act at a threshold. Take: "used so far, projected at reset" on one line, and a ratio view (cache read against total). https://vercel.com/docs/pricing/manage-and-optimize-usage
- btop network box: auto-scaled graph with a fixed minimum, upload and download shared scale option, three symbol sets (braille, block, tty) because font support varies. Take: auto-scale the chart to the visible window, pick the block symbol set (verified in FiraCode), keep a scale label. https://github.com/aristocratos/btop
- ncdu: sorted list with a bar and percent per row, drill down with enter or right and back with left, sort keys that flip on a second press (`s` size, `n` name), toggle keys for columns. Take: drill-down by size for project then session, bar plus percent column, sort key flips. https://dev.yorhel.nl/ncdu/man
- k9s: column sort shortcuts (shift plus letter), warn and critical thresholds that color a usage column, skins. Take: threshold colors for context size and limit share, one-key column sort. https://github.com/derailed/k9s
- asciigraph: line graph from box characters with height, caption, series legend, value thresholds that change color. Take: multi-series chart with a caption and a threshold line, if a 6 to 8 row chart is wanted. Sparklines stay on `▁▂▃▅▆▇`. https://github.com/guptarohit/asciigraph

## 4. Principles

1. Gauges first, always on screen. Left share, reset time, burn and empty time on one row per account. Word "left" printed.
2. A number that is a guess says so. Empty time and forecast carry "est". Costs carry "API-equivalent" once.
3. One table schema. Session, project, model and account rows use the same columns, so grouping changes the rows and not the layout.
4. Group-by and period are one key each. Period `p` cycles today, week, month (never `m`, the global mode key), group `g` (project, model, account, day). Sort flips on a second press.
5. Drill down by size. `enter` on a project opens its sessions, `enter` on a session opens it in Agents history. `esc` goes back.
6. Full titles and project names. The title column takes all free width, the numeric columns are fixed and right aligned.
7. Warnings name the action. `ctx 182k over: wrap or compact`, `empty 16:26: move runs to account X`. Color marks the flag, text says the action.
8. Leftover height is more rows. The ledger grows to the bottom. Charts stay small (6 to 8 rows or one row per day).
9. Forecast is drawn as a continuation, not a separate box: the daily series and the 5h timeline both continue to the right of "now" with `╌`.

## 5. Design directions

### Direction A: Ledger (table-first)

One line: gauges on top, one large foldable ledger below, everything else is a column or a one-line summary.

Inspiration: ccusage tables, ncdu drill-down and size bars, k9s thresholds, Console group-by.

Layout at 235x62:
- Row band 1, "Accounts", rows 1 to 8 (4 accounts, 2 rows each is not needed, one row each plus header and a policy sub-line count). Columns: account, policy chip, 5h gauge (24 cols) with percent left and reset, week gauge (24 cols) with percent left and reset, burn in percent per hour, empty at (est), runs planned tonight, flag text with the action. Full width.
- Row band 2, "Spend", rows 9 to 16, three side by side regions: left 90 cols period totals (today, week, month: cost, tokens, cache percent, sessions, compactions) as a 3 row mini table; middle 80 cols daily bars, one row per day for 8 days (`Mon ███████▏ $84`), `╌` bar for today's forecast to end of day; right 60 cols savings, one line each: cache, routing, rtk, scripts, total, each with the period word.
- Row band 3, "Ledger", rows 17 to 60 (44 rows): group header per project (name, sessions, tokens, cost, bar, share percent), folded sessions inside. Session row columns: state glyph, account, full title (flex), model, tokens, cache read, cache write, cost, 5h share, week share, context size with `ok warn over` color. Group-by key `g` swaps project for model, account or day.
- Status bar: `p period  g group  s sort  enter drill  / filter  P policy`.

At 118x60: gauges stack to two rows per account (5h row, week row), spend band becomes totals line plus 8 day rows with bars, savings one line. Ledger drops cache write, week share and model to a second line under the title (`opus  r 2.5M w 224k  wk 0.8%`). Same data, taller.

Data per region: Accounts from `limits` plus `accounts.yaml` plus afk plan (new). Spend from `today/week/month/days`. Savings from `rtk`, `cache_saved_usd`, `routing_saved_usd`, scripts (new). Ledger from `sessions` (all fields exist).

Strength: fixes U1 and U5 directly, lowest build cost. Weakness: time is not visible, forecast is a column.

### Direction B: Runway (timeline-first)

One line: a shared time axis. Each account is a lane of 5h windows, the daily chart continues into the forecast, and overnight runs are blocks on the lanes.

Inspiration: claude-monitor session and prediction view, ccusage 5h blocks, btop graph scaling, Cost Explorer forecast continuation, asciigraph threshold line.

Layout at 235x62:
- Region 1, "Runway", rows 1 to 20, full width. Axis: last 24h to next 24h in 4 min per column (about 210 columns of plot after a 22 col label gutter), hour ruler on top, `▼ now` marker. One lane per account, 3 rows each: row one is the 5h window block (`━` for used, `░` for left, ends at reset), row two is burn sparkline for the past (`▁▂▃▅▆▇`) continuing as `╌` to the predicted empty time, row three is the weekly gauge as a text bar with reset tick. Planned afk runs are `▓` blocks on the lane of the account they use, labeled with the job name. Empty-time and conflicts are red text on the lane. Lane gutter: account, policy chip, 5h left, week left.
- Region 2, "Daily spend", rows 21 to 32, left 120 cols: 14 daily columns as vertical bars (8 rows high) with `$` labels, today's bar split solid and `╌` forecast, a threshold line at the weekly budget if set. Right 110 cols: period totals and savings as a 2 column key-value list (today, week, month, cache, routing, rtk, scripts).
- Region 3, "Sessions today", rows 33 to 60: ledger as in direction A but scoped to the period, with an hour sparkline per session (from `activity`) after the cost column so rows line up under the Runway axis. Context warning column.
- Keys: `h l` pan the axis by 6h, `d` 7 or 14 days, `enter` on a lane block opens the run or session.

At 118x60: lane label gutter shrinks to 14 columns and the axis shows 12h back and 12h ahead at 5 min per column. Daily chart drops to 7 bars. Sessions lose the sparkline and fold to two lines.

Data: lanes need `burn_series` per account and afk plan (both new), plus `limits`. Daily from `days`. Sessions from `sessions` with `activity`.

Strength: answers "when does a limit hit, do tonight's runs fit" at a glance, which no table does. Weakness: the most new data and the most drawing code. Needs a fallback when `burn_series` is missing (use the current rate as a straight `╌` line).

### Direction C: Account cards (grid, policy-first)

One line: one card per account holding its gauges, policy, routing chain and top sessions. A strip on top says the total.

Inspiration: better-ccflare account health and fallback chains, Vercel allotment with projected cost, Cursor included versus extra meter, claude-monitor provenance labels.

Layout at 235x62:
- Region 1, "Total" strip, 3 rows, full width: cost today, week, month with `API-equivalent`, tokens, cache percent, savings sum, rtk percent, and the routing chain (`A use > B ask-first > C never > D never`) with the reason an account is skipped (`limited until 20:30`).
- Region 2, four cards in a row at 58 cols each, rows 4 to 36 (33 rows). Card content: title line (account, policy chip, active dot), 5h gauge with left percent, reset, burn and empty time; week gauge likewise; "used so far, projected at reset" line (Vercel style); cost today, week; sparkline of 7 days; runs tonight with job names; flag line with the action. Below, "top sessions on this account", 12 rows, full title wrapped to two lines when long, cost and context flag. A card for an account that is `never` stays visible and says `not used by crow`.
- Region 3, rows 37 to 60: left 120 cols "By project" drill list (ncdu bars with percent), right 110 cols "By model" and "Savings" boxes with the same bar style. `enter` on a project opens its sessions inline.
- Keys: `h l` card, `p` cycle policy on the card, `< >` priority (no AltGr on AZERTY), `enter` drill in the card's list.

At 118x60: cards become a vertical stack of 4 compact blocks of 8 rows (two gauge rows, a projection line, a sessions line count, one flag line), a single "Sessions" ledger replaces per-card session lists, project and model bars share a row at 59 cols each.

Data: limits and policy exist. Projection needs `limit_share` sums (exists) and the burn rate. Runs tonight needs the afk plan (new). Per-account sessions come from `sessions.account`.

Strength: policy edits and routing live where the gauge is, which suits an owner who manages several accounts and wants to see who gets work. Weakness: duplicated chrome per card, session detail is split across cards, so a global "what costs most" needs the lower region.

## 6. Recommendation

Build A first (all data exists, closes U1 and U5, raises information density). Add B's runway as a second view on key `t` when `burn_series` and the afk plan exist. C's routing chain and "used so far, projected" line are cheap and can be lifted into A's Accounts band without the card grid.

## 7. Open choices for the owner

- Show a budget line (weekly dollars) on the daily chart, or keep gauges as the only limit.
- Whether `never` accounts appear in the gauges at all or only as a dimmed row.
- Whether the scripts saving is worth a logging change in `bin/`.
