# Microsoft Graph for the deck Today box

`bin/crow-today` reads today's and tomorrow's calendar, the unread and drafts counts, and picks the mail worth a look through your own Entra app. The deck Today rail shows it; without data the rail hides.

## One-time setup
1. Register a single-tenant app in Entra: "Allow public client flows" on, no secret. Delegated permissions: User.Read, Calendars.Read, Mail.Read, offline_access.
2. If the tenant blocks user consent, ask an admin to grant admin consent. Search your mail for the answer before saying it is still pending: it often comes back within minutes.
3. Copy `docs/graph.example.yaml` to `local/graph.yaml` (tenant, client id, user).
4. `bin/crow-today login` prints a code. Open https://login.microsoft.com/device in the real browser (`bin/crow-open`), enter it, sign in. The refresh token goes to the keyring (`secret-tool`, service `crow-graph`).
   From a Claude session run it in the background, read the code from its log, and wait for the owner to say "signed in".

## Use
- `bin/crow-today json` prints the data and caches it 5 minutes in `local/state/today.json`; `-f` skips the cache.
- Exit 2 means not signed in or the refresh token expired: run `login` again.
- The token request asks for `.default`, so it gets exactly the scopes the admin consented.

## Mail to look at
Inbox mail from the last 3 days, kept when all hold: you are in To, the sender is a colleague (your domain) or someone you wrote to in the last 90 days (Sent Items), it is not a meeting reply or a robot address, and you have not answered the thread since. High importance first, then unread.

With `mark_important: true` in local/graph.yaml the Outlook flag becomes the shared marker: crow flags each picked mail once and adds the category `crow` (existing categories kept), and Today lists every flagged mail (◆ crow's, ● yours). A flag the owner clears stays cleared (`local/state/today-marked.json` remembers what crow flagged). Crow marks its own flag complete once the owner answered the thread; it never touches the owner's flags.

## Mail actions (deck Today pane, and from any session)
- `bin/crow-today done <id>` marks the flag complete; `snooze <id>` flags it due tomorrow 08:00 and Today hides it until then (`until` in the JSON).
- `bin/crow-today thread <id>` prints the thread newest first, plain text, so a session reads mail without the claude.ai connector.
- `bin/crow-today draft <id> [--all] < text` saves a reply draft (text above the native quote) in Outlook Drafts. It never sends.
- Deck: focus Today with `tab`, `j k` walks the mail rows, `o` opens it in Outlook (real browser), `x` done, `s` snooze, `enter` hands off to a hub session that reads the thread and leaves a French draft.

## Gotchas
- `calendarView` returns events that only touch the window edge (an all-day event ending at midnight shows on the next day). The script drops them and duplicates with the same start and subject.
- Times come in the `timezone` of local/graph.yaml (default UTC); it must match the machine zone for the "in N min" countdown.
