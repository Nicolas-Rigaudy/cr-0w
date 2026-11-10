# Design references in planning items

Use when a Logics item, plan or handoff cites a design image, board or mockup as the thing to build or match.

## Rule
- Every design pass folder gets a `Decision.md` once the owner picks: approved files, rejected files, date. Write it in the same session as the OK.
- An item cites the `Decision.md` and the approved file named in it. A file that only exists in a folder, or that a compare page shows, is an option, not a reference.
- A folder with no `Decision.md` holds no approved reference. Ask the owner before citing anything from it.
- When a corpus or plan copies a "Sources" line from a parent item, re-check each file against its `Decision.md`. Parent items go stale when a later pass changes the pick.
- An item title such as "from the approved concept art" does not approve anything. Name the file.

## When the owner rejects a reference later
1. Add it under "Rejected" in the folder's `Decision.md`.
2. `grep -rn <file> logics/` and fix every item that cites it, the parent plan item first.
