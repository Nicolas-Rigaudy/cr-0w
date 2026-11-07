## item_112_skill_and_reply_contract_for_each_mission_kind - Skill and reply contract for each mission kind
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 75%
> Confidence: 65%
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:21:50

# AI Context
- Summary: Give every kind a skill: reuse corpus, review, review-threads and implement-task, and fill the gaps for research, design, repair and generic write.
- Keywords: skill per kind, reply contract, research skill, design skill
- Use when: Before a kind runs as a mission.
- Skip when: Kind defaults (owner review item).

# Problem
- Four kinds have no skill: research, design, repair and a generic write.
- Existing skills end with chat output, not a contract file.

# Scope
- In:
  - Per kind: name the skill in the kind file, add the reply contract through the brief.
  - Write the missing skills chosen in D5.
  - Check each existing skill can end its turn by writing the version file.
- Out:
  - Changing bundled logics-manager skills (upstream).

# Acceptance criteria
- AC1: Every kind file names a skill with status `exists`, local or bundled.
- AC3: Each existing skill reused by a kind (corpus, review, review-threads, implement-task) ends its turn by writing the version file when the brief asks, without edits to bundled skills.
- AC2: A run of each phase 1 kind with its skill writes a valid version file.

# AC Traceability
- request-AC1 -> AC1. Proof: kind files with no missing skill.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `logics/request/req_037_mission_kinds_and_run_contract.md`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Phase 1 kinds need it before their first mission.

# Notes
- D5 decided 2026-11-07: one skill per kind, each in its own phase (research and design first, repair and write later). Ready for phase 1.
- Map in the spec README.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
