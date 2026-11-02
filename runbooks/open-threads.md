# Open threads

One file in the state dir (`vault.crow_dir` or `local/state/`): `open-threads.md`. It lists open work by who acts. It stays a plain checkbox list that Obsidian Tasks can read.

## Format
```
## Me
- [ ] 2026-10-22 [project] text → next step
## Crow
- [ ] 2026-11-02 [project] item_047: text → next step
## Waiting
- [ ] 2026-11-02 [project] Person: text → next step
```
- Line shape: `- [ ] YYYY-MM-DD [project] text → next step`. The date is the last touch.
- `## Me`: the owner acts.
- `## Crow`: the hub or a project session acts. Point to the Logics item (`item_047`) when one exists instead of restating it.
- `## Waiting`: someone else acts. The text starts with the person's name, then a colon.
- Moving an item to another actor means moving the line and re-dating it.

## Rules
- Done means gone. A ticked item leaves the file in the same write and goes to `activity/YYYY-MM-DD.md` as `- HH:MM [project] done: text`.
- Never write a status note here. Project state goes to `status.md`, constraints and decisions to `decisions.md`.
- Stale: no touch for 14 days (date older than today minus 14). Readers list stale lines and ask keep, re-date or drop. Drop means the line leaves the file.
- One line per item. Update a line in place, never duplicate it.

## Who reads and writes
| Skill | Access |
|---|---|
| `wrap`, `eod` | Write: add, update, move, remove done lines |
| `brief` | Read: Me first, then Waiting with age in days, then stale lines |
| `afk` | Read: the Crow section only |

Any new writer follows this runbook and does not add other sections.
