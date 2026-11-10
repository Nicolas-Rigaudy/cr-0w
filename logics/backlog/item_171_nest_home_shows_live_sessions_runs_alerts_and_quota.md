## item_171_nest_home_shows_live_sessions_runs_alerts_and_quota - Nest Home shows live sessions, runs, alerts and quota
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 95%
> Confidence: 90%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 18:11:26

# AI Context
- Summary: Nest Home and header read the item_169 sections: live sessions and dispatch runs as one-liners beside missions, alerts from checks, quota per account and the live count in the header, today's events in Agenda.
- Keywords: nest home, sessions, dispatch runs, alerts, quota gauges, header, phone, fold 6
- Use when: wiring Home or the header to deck sources.
- Skip when: Work and System screens (item_172), hero motion (item_158).

# Problem
- Home lists missions and inbox only. With item_169 the data exists, but no screen shows sessions, dispatch runs, checks or quota (req_054).

# Scope
- In:
  - Mockup first (owner prefs: options side by side, then OK). Base: the approved Home boards (item_158 and prod_007 sources), not new art.
  - Running work: missions, live sessions and dispatch runs in one list of one-liners (project, title or label, state, age, account). Tap: session row opens its send box (item_170); run row shows its report path and commits.
  - Alerts: `checks` items that need the owner, above running work.
  - Header chips keep need/working counts. Quota and live counts move to the HUD row above the holotable (approved design below).
  - Quota per account from `usage` and `accounts`; live count = running missions + working sessions + running dispatch runs. Stale sources show their age.
  - Hero mood reads the same live count (input to item_158's open mood order question, not a change to it).
  - Both Fold 6 screens checked (req_038).
- Out:
  - Work and System screens (item_172). Session conversation view (req_027).

# Acceptance criteria
- AC1: With the owner's real day running, Home on the phone lists the same live sessions and running dispatch runs as deck's agents panel.
- AC2: Header gauges match `crow-usage` per account; the live count matches deck.
- AC3: A stale source shows its age on its block; the rest stays live.
- AC4: Cover and inner screens pass the req_038 checks; bundle under 90 kB gzip.

# AC Traceability
- request-AC1 -> AC1, AC2. Proof: Home matches deck with an empty mission store.
- request-AC3 -> AC3. Proof: stale shown per source.

# Decision framing
- Product framing: Not needed (prod_005 Home and header)
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`, `prod_007_nest_home_design_pass`
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`
- Depends on: `item_169_core_serves_every_deck_source_as_read_only_sections`
- Related: `item_158_nest_hero_live_data_state_motions_and_actions`, `item_165_telemetry_strip_on_the_laptop_home`
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: High
- Rationale: Home is what the owner opens on the phone.

# Notes
- Owner 2026-11-10: most runs will be missions. Mission rows and actions come first; sessions and dispatch runs follow.
- Design gate: mockup and owner OK before code. If it waits on the owner, build item_170 meanwhile (task_020 order rule).

# Approved design (owner 2026-11-10, pass 3b option D)
- Source: vault `CR-0W/Research/nest-design/item_171/` — `D-{cover,inner}-{modern,cyber}.png`, `mockup.html?v=D`, `compare.html`. Build from these only (not pass 1 or 2, not A/B/C/E).
- Hero: the canvas holotable `Hero` (item_162, as on HomeLaptop) is the centerpiece on both screens. Drop the static `.hero.tall` from phone Home.
- HUD row above the holotable, two panels with the same border, height and type:
  - Quota HUD: per account a name, two thin bars (5h and week left), `5h/wk` numbers; footer `5h / week left · usage N min old` (stale age).
  - Status HUD: need you ◆N, working N, slots N/2; footer `live` (or the source age when stale).
- Butler line under the HUD row: `CR-0W › N need you, Sir. N working.` with cursor.
- Cover (368): HUD row, holotable (~400px), Needs you first, Running, Agenda; tab bar.
- Inner (707): left column 330px sticky = HUD row, holotable (fills), Agenda (non-interactive). Right column = Needs you first, Running.
- Running: one list with a kind icon (mission / session / run), no subheads. Tap: session = send box (item_170), run = report path and commits, mission = mission view.
- Owner on kinds: sessions = interactive windows you answer, runs = headless `crow dispatch` jobs with a report, missions = multi-step tracked jobs. Same list, icon only.
- Butler chat lives on Home (decision 2026-11-10, built by item_176): a text box under the holotable, reply types out in the butler line, the droid reacts while it thinks/talks; tapping the line opens the transcript on Desk. item_171 leaves room under the holotable for it but does not build it.
- Build session: fresh, on an account with week quota left.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
