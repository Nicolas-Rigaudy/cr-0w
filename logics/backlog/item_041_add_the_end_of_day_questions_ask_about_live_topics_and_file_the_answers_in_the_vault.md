## item_041_add_the_end_of_day_questions_ask_about_live_topics_and_file_the_answers_in_the_vault - Add the end-of-day questions: ask about live topics and file the answers in the vault
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Add the end-of-day questions: ask about live topics and file the answers.
- Keywords: add, end, day, questions, ask, about, live, topics, file, answers, vault
- Use when: building this slice of req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap.
- Skip when: another slice of the same request.

# Problem
- The owner wants crow to ask where important topics stand and write the answers where they belong.

# Scope
- In:
  - a skill `eod` (user-invoked: end of day), also called from `/wrap`
  - picks live topics without a model: open threads, Jira In Progress and Merge Review items, hub active focus, vault `hot.md`, topics tagged in the last days
  - asks at most three questions with AskUserQuestion, self-contained, with choices (done, in progress, blocked by X, dropped) and a free text option
  - writes the answers: End of day (Shipped, Blocked / carry over, Note for tomorrow), Meetings & notes, and a dated line in the relevant project note
  - shows what it wrote
- Out:
  - editing notes outside the daily sections and the project note line it appends
  - asking about topics with a trace today

# Acceptance criteria
- AC1: with no input the skill asks only about topics without a trace today
- AC2: each answer lands in exactly one named section, shown back to the owner
- AC3: a skipped question writes nothing

# AC Traceability
- request-AC4 -> This backlog slice. Proof: the ACs above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: depends on: item_039, item_040.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

