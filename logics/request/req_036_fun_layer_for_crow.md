## req_036_fun_layer_for_crow - Fun layer for crow
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 65
> Confidence: 60
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Feature track. The fun layer: ambient screensaver, animated droid avatar with sound, stats and achievements, drones view of agents. Reads events only; never on the critical path.
- Keywords: screensaver, ambient display, avatar, droid, sounds, achievements, streaks, stats, drones, cyberdeck
- Use when: building delight features.
- Skip when: anything that gates work.

# Needs
1. **Ambient screensaver.** Live missions, asks waiting, quota, machines as a calm animated scene; laptop idle screen and an old phone as an always-on cyberdeck display (with req_033). Ultron shipped a macOS screen saver over its API.
2. **Droid avatar.** Animated character tied to the butler persona; states (idle, working, waiting on you, alarm); short sound chirps; mute with quiet hours.
3. **Stats and achievements.** Missions shipped, tokens saved, asks answered fast, streaks, first auto-merge, incidents turned into automations. Weekly in-character report.
4. **Drones view.** Agents as drones moving between projects, coloured by state; click opens the mission.
5. Same tokens as nest and deck (req_038); respects reduced motion.

# Context
- Owner choices 2026-11-06: ambient screensaver, droid avatar, stats and achievements, drones view; full character persona.
- Lowest priority track (phase 4).

# Acceptance criteria
- AC1: The screensaver runs on the laptop and the old phone from the core API.
- AC2: The avatar reacts to at least four mission states with sound and mute.
- AC3: Five achievements computed from the event log.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_015_wishlist_ambient_mobile_voice_colleague_ideas.md`
- `logics/request/req_032_spike_butler_engine_and_persona.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
