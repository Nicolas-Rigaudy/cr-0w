## item_118_butler_persona_file_steward_with_a_droid_flavour - Butler persona file Steward with a droid flavour
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 95%
> Confidence: 90%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:22:07

# AI Context
- Summary: Write the butler persona file: a Steward base with a droid flavour. Named CR-0W, formal and polite, addresses the owner as "Sir". Neutral while work focus is on.
- Keywords: butler, persona, steward, droid, work focus
- Use when: Giving the butler its voice.
- Skip when: Butler engine and tools.

# Problem
- The persona skeleton has three candidates (Steward, Ship AI, Familiar).
- Owner decision 2026-11-07: Steward base with a droid flavour. Neutral at work while work focus is on.

# Scope
- In:
  - `persona.md` from `persona.template.md`: voice rules, sample replies to the three spec questions, forbidden tics.
  - A neutral variant used whenever work focus is on.
  - Name CR-0W and address form "Sir" (owner, 2026-11-09).
- Out:
  - Sounds and voice output.

# Acceptance criteria
- AC1: The file holds the Steward voice with the droid flavour and three sample replies.
- AC2: A work-focus turn uses the neutral variant: same structure, no tics.
- AC3: The file names the butler CR-0W and uses "Sir" in personal mode only; the work variant uses no form of address.
- AC4: The owner approves the draft (req_032 AC3).

# AC Traceability
- request-AC3 -> AC1 to AC4. Proof: the persona file the owner approves.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_032_spike_butler_engine_and_persona`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Req_032 AC3 needs an approved persona file.

# Notes
- 2026-11-07 owner decisions 20 and 21.
- 2026-11-09 owner: the butler is CR-0W (the droid behind the project name), address form "Sir". Formal and polite; not fussy, not anxious (the 07/11 "fussy, anxious" wording was never the owner's).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
