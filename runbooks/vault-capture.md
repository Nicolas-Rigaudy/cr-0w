# Vault capture

How the owner's implicit work reaches every recap. Plain logging costs no tokens.

## Channels
| Channel | Use |
|---|---|
| `crow note "text #tag"` | One line, timestamped, into today's daily under Capture. `--section meetings\|eod` for the other sections. Creates the daily from the vault template when missing. |
| Obsidian mobile / `00 Inbox` | Free notes, picked up by `crow notes`. |
| `eod` skill | Asks at most 3 questions on live topics, files the answers. |
| Brief | Asks at most 3 questions on topics with no trace in 3 days. |

## Reading
`bin/crow-notes --since <Nd|date>` prints JSON: per day Capture, Meetings & notes and End of day lines, plus Inbox notes added and `10 Work` notes modified (first lines). Recaps (mgmt-weekly, daily, weekly, brief) run it in their gather step and say they did.

## Tags
`#team #naval #blocage #decision #demo`. Extend by adding a line here.

## Rules
- Writers touch only Capture, Meetings & notes, End of day. Never rewrite existing lines.
- Machine state stays in the state dir, not in the vault notes.
- The vault's own Claude owns triage. Convention text for its CLAUDE.md: see handoff below.

## Vault session handoff
Add to the vault CLAUDE.md: "crow appends timestamped lines under Capture, Meetings & notes and End of day of the daily note (`crow note`, `eod`). Keep those headings and the `- HH:MM text #tag` line shape. Tags: #team #naval #blocage #decision #demo."
