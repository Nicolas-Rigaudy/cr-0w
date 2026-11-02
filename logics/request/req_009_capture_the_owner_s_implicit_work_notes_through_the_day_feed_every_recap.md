## req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap - Capture the owner's implicit work: notes through the day feed every recap
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Hub context
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: The owner's work that is not in git, Jira or mail (calls, decisions, topics in motion) is invisible to recaps. Capture it cheaply during the day in the vault daily note and make every recap read it.
- Keywords: capture, owner, implicit, work, notes, through, day, feed, recap
- Use when: planning or building note capture, end-of-day questions, or the sources of mgmt-weekly, daily, weekly and brief.
- Skip when: vault folder structure for crow data (see the vault proposal) or Confluence posting.

# Needs
- Recaps (Monday management weekly, daily, Friday rollup, brief) only see git, Jira and mail. Calls, oral decisions, vendor talks and topics the owner works on without a ticket stay invisible (2026-11-02: the weekly draft missed the Naval document and the Google ID wait).
- Writing a note must cost under five seconds, or the Capture section of the daily stays empty (it was empty on 2026-10-29).
- Crow should ask the owner where live topics stand and file the answers itself, at the brief and at end of day.

# Context
- The vault already has the places: daily note `40 Journal/Daily/` with `📥 Capture`, `📆 Meetings & notes`, `🌙 End of day` (Shipped, Blocked / carry over, Note for tomorrow), and `00 Inbox`. The vault proposal puts crow's own lines under `## CR-0W` in the daily.
- The vault is owned by the owner and the vault Claude (vault CLAUDE.md contract: groom freely, propose big changes). Crow only appends to the daily sections above and to `## CR-0W`.
- Plain logging costs no tokens: a script appends, the model only asks and summarizes.
- Owner decisions (2026-11-02): channels are `crow note` CLI, question at /wrap, Obsidian mobile and Inbox, plus end-of-day and brief questions via AskUserQuestion about the topics that matter now, filled in by crow.
- Tags in Capture lines: #team #naval #blocage #decision #demo (free to extend).
- Skill mgmt-weekly (local/hub-skills/mgmt-weekly) lists its sources to read first; this request adds the vault to that list.

# Acceptance criteria
- AC1: `crow note "text" #tag` appends a timestamped line to today's daily Capture, creating the daily from the template when missing, with no model call.
- AC2: one gather script returns, for the last N days, the Capture, Meetings and End of day lines, the Inbox notes and the project notes modified, as compact JSON.
- AC3: mgmt-weekly, daily, weekly and brief read that JSON and say they did.
- AC4: an end-of-day step asks about the live topics with AskUserQuestion, writes the answers into End of day / Meetings / the project note, and never writes outside the agreed sections.
- AC5: the brief asks at most three questions about topics with no recent trace (no commit, ticket move, mail or note in 3 days).
- AC6: the convention (tags, sections, who writes where) is documented and the vault session is told to add it to the vault CLAUDE.md.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics_manager/flow.py`
- `logics_manager/assist.py`
- `tests/python/test_logics_manager_cli.py`

# Backlog
- item_039 -> item_040 -> item_041 -> item_042 (this order)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

