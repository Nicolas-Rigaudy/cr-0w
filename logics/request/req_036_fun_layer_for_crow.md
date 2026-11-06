## req_036_fun_layer_for_crow - Fun layer for crow
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 55
> Confidence: 55
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: The fun layer: ambient screensaver, droid avatar and sounds, stats and achievements, drones view of agents at work.
- Keywords: screensaver, ambient, avatar, droid, sounds, achievements, stats, drones, cyberdeck
- Use when: building delight features on top of core events.
- Skip when: anything on the critical path of missions.

# Needs
1. **Ambient screensaver.** Missions moving through the system, waiting items, clock, next meeting. Runs on the laptop lock screen and on an old phone as an always-on cyberdeck display.
2. **Droid avatar.** Animated crow character reacting to states (working, waiting on the owner, error, idle) with short sound chirps; shown in nest and deck.
3. **Stats and achievements.** Missions shipped, tokens saved, streaks, per-agent board, a weekly system report in character.
4. **Drones view.** Agents drawn as drones moving between projects on a map while they work.

# Context
- Reads only core events; never blocks a mission.
- Persona from req_032; display shape from req_031 and req_033; ties to req_015 wishlist.

# Acceptance criteria
- AC1: The screensaver shows live missions from core events on the laptop and the old phone.
- AC2: The avatar changes state within 5 seconds of a mission event.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
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
