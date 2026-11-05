## item_087_work_prose_ai_writing_checklist - Work prose AI-writing checklist from humanizer
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Turn the AI-writing signs from blader/humanizer (based on Wikipedia's Signs of AI writing) into a short checklist the hub applies to work prose, without installing the skill.
- Keywords: work, prose, writing, checklist, humanizer
- Use when: Editing owner-prefs prose rules or the review gate for MR, Jira, mail, Teams and weekly text.
- Skip when: Code review patterns or UI design.

# Problem
- Work prose (MR descriptions, Jira, mails, weekly cards) must read as written by the owner. owner-prefs.md already bans em-dashes and "X, so Y" justifications, but there is no fuller list of AI tells.
- blader/humanizer (MIT, 32 KB SKILL.md) lists them. Installing it costs context each time it loads; a distilled checklist costs almost nothing.

# Scope
- In:
  - Read humanizer's SKILL.md once, keep the tells that apply to French and English work prose.
  - Add them as a compact checklist to local/owner-prefs.md (prose section) and to the review ledger used by the mr gate, daily, weekly and mgmt-weekly drafts.
- Out:
  - Installing the humanizer skill.
  - Rewriting past artifacts.

# Acceptance criteria
- AC1: The checklist exists in owner-prefs.md, under 25 lines, each tell with a one-line fix.
- AC2: The mr gate and the draft skills (daily, weekly, mgmt-weekly) point to it.
- AC3: One real draft run through the checklist shows it catches at least one tell the current rules miss.

# AC Traceability
- AC1 -> this item. Proof: local/owner-prefs.md, Writing, "AI-writing tells": 18 tells plus 2 header lines, each with a one-line fix, French and English examples (local file, gitignored).
- AC2 -> this item. Proof: skills/mr/SKILL.md (Prepare), .claude/skills/daily, weekly and mgmt-weekly SKILL.md, and the AI-sounding row of local/reviews/common.md point to the checklist.
- AC3 -> this item. Proof: run on the owner's open lambda MR descriptions (2026-11-05). The Sonar MR gives every list item a bold label ("**Type hints:**", "**Python:**", ...), a formatting tell no earlier rule covered. The booking schedules MR has a hidden actor ("the conflict was only logged"), weak alone.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (none, from the vault note "Github repos and websites", owner pick 2026-11-05)
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Cheap, and every work artifact benefits.

# Notes
- Source: https://github.com/blader/humanizer (MIT). Owner picked this on 2026-11-05 from the vault link review.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
