## item_125_workshop_concept_art_round_2_original_cast_states_phone_and_night_views - Workshop concept art round 2: original cast, states, phone and night views
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 40%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 13:52:33

# AI Context
- Summary: Second image-prompt round on the picked droid workshop world: an original droid cast per theme, one image per droid state, a phone view and a night (AFK return) view. Prompts only, run by the owner in Gemini or ChatGPT; no render runs.
- Keywords: nest, workshop, droid cast, concept art, image prompt, Gemini, states, phone, night
- Use when: continuing the nest visual design after the workshop pick.
- Skip when: building the nest app (item_126 and after).

# Problem
- The workshop concept (refs/v4 g1 cyber, g3 modern) is picked, but its droids copy film droids (WALL-E, astromechs) and nothing shows the states up close, the phone layout or the night return view.

# Scope
- In:
  - Prompt file v5 in the vault, `CR-0W/Research/nest-design/Nest Prompts v5.md`: cast sheet per theme (6 to 8 original droids, distinct silhouettes), state close-ups (working, needs you, failed, sleeping, queued), phone portrait top-down, night view.
  - Same rules as v4: "in-game screenshot of an isometric <genre> game", stations numbered, amber/red/green/blue reserved per state, no franchise names (Gemini blocks them), no UI words.
  - Compare page from the owner's images with `bin/crow-compare`, decisions asked with AskUserQuestion.
- Out:
  - Blender or code renders (rejected three times, see memory design-renders-need-wow).
  - Hangar, grid city and war table (kept in mind in req_040, not this round).

# Acceptance criteria
- AC1: v5 prompts written for cast (both themes), states, phone and night, each pasteable as one block.
- AC2: The owner's images are filed under the vault `CR-0W/Research/nest-design/v5/` and shown side by side in one compare page.
- AC3: The owner picks a cast per theme and approves the state, phone and night reads; picks recorded in req_040.

# AC Traceability
- request-AC1 -> AC1, AC2, AC3. Proof: vault folder nest-design/v5 and req_040 picks section.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_040_nest_deck_two_themes_3d_world`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: R&D, after reviews and sprint work. Unblocks item_126.

# Notes
- 2026-11-07 owner ran prompt 1 (cyber cast), 7 takes in the vault `v5/`, compare page `local/state/research/nest-design/compare-cast-v5.html`.
- 2026-11-09 cyber cast picked as 8 families (req_040 "Cast pick, cyber"). State colours dropped: state through pose and props. Modern deferred to a later round, so AC1 to AC3 now cover cyber only.
- Next: v5b prompts (merged cast sheet, states, phone, night) in `Nest Prompts v5.md`, owner runs them, compare, approve.
- Images v4 and the compare page: `local/state/research/nest-design/refs/v4/`, `compare-concepts-v4.html`.
- 2026-11-09 paused by owner: too much effort on droid art for the need. Cast picks kept in req_040; 1c prompt ready in the vault. Step back on what nest needs to design, and a Jarvis-style 3D hero is being considered as an alternative.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
