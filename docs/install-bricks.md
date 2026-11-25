# Install the bricks

Crow is one core plus optional bricks (adr_002). This page installs the two that serve a
browser: `core-api` (the back, crow-core) and `nest` (the front, a static build that core
serves). Each one installs on its own. Nothing here needs a forge or Tailscale.

## Needs

- Linux or WSL2, with `git`, `go` 1.27+, `node` 20+ and `npm`.
- The engine checkout (`CROW_ENGINE`) and a home folder (`CROW_HOME`, which can be the same
  checkout). `install.sh` sets both and writes a starter `crow.yaml`.

## Pick bricks

`install.sh` asks which bricks to turn on for a new data dir. `--modules a,b,c` answers
without a prompt. The default set is `core,git,core-api,nest,deck`. Every other brick is off.
An existing data dir keeps its `crow.yaml`, or every default when it has none.

| Brick | What it adds | Default |
|---|---|---|
| core | Dispatcher, runs, sessions, accounts, doctor, notes, brief and wrap. Always on. | on |
| git | Commit and push hooks, git activity, worktrees, identity check. | on |
| forge | Merge requests and reviews on GitLab or GitHub. Needs git. | off |
| vault | Obsidian vault notes, tasks and the eod skill. | off |
| mail | Calendar and mail counts, Today box, reply drafts. | off |
| team | Sprint, 1:1s and mail to answer. | off |
| afk | AFK and night runs: page, preflight, keep-awake. | off |
| retro | Session retros, digest, cost audit, tool update watch. | off |
| releases | Changelog, release proposals, releases view. | off |
| core-api | crow-core: store, inbox, missions and the API. | on |
| nest | Nest web front. Needs core-api. | on |
| deck | crow-deck and crow-watch terminal front. | on |
| design | Screenshot, crop and compare tools, design notes. | off |
| office | Mail text, tracked-changes DOCX, simple A4 DOCX. | off |

Optional bricks from their own repo install with `crow modules add <git url>`
(see [writing-a-brick.md](writing-a-brick.md)):

| Brick | What it adds | Install |
|---|---|---|
| cr-0w-archeotech | Crow agents and usage panels for the archeotech desktop shell. | `crow modules add https://github.com/Nicolas-Rigaudy/cr-0w-archeotech`, then `crow archeotech link` |

The Default column is the install starter set. Without a `crow.yaml`, the manifest default
applies (`crow modules` shows both).

Switch one later with `crow modules on <brick>` or `crow modules off <brick>`. It edits
`crow.yaml`, keeps the other keys and config, and relinks the skills. It refuses to turn a
brick on while its requires are off, and to turn off a brick that an on brick requires.
An off brick's commands print `crow: module <name> is off (crow.yaml)` and exit 2.

## crow.yaml

Turn off the bricks you do not use. Their sections, pages and boxes disappear from Nest and
the deck, and their core routes answer 404 `module off`.

```yaml
# $CROW_HOME/crow.yaml
version: 1
modules:
  forge: false   # no GitLab or GitHub
  mail: false    # no mail or calendar
```

`crow modules` shows the result. Core reads the file again when it changes, and Nest follows
within a few seconds. No restart is needed.

The Nest look lives under `modules.nest`: `theme`, `hero` and `pages` (nav order). Set them
with `crow modules set nest theme cyber` (`pages` takes a comma list, `""` for every page).
It refuses a name Nest cannot build and edits only the `modules.nest` lines. Core serves the
same writer at `GET` and `POST /api/settings/nest`. A hero or pages change needs
`crow nest build`. A save through core runs it and Nest reloads when it succeeds. A failed
build names the brick and keeps the last good build served. The theme switches live.

## core-api (back)

```bash
crow core serve        # builds crow-core when the Go sources changed, then serves
```

- It listens on `127.0.0.1:8787` and on the unix socket `local/state/core.sock`.
- To change the port, write `local/core.yaml` with `serve: {listen: 127.0.0.1:9000}`.
- To keep it running, use `crow core install`. This writes a systemd user unit and starts it.
  On WSL, systemd must be on: add `[boot]` and `systemd=true` to `/etc/wsl.conf`, then run
  `wsl --shutdown` from Windows. Without systemd, run `crow core serve` in a terminal you keep
  open.
- `crow core install` also tries `tailscale serve`. Without Tailscale it prints that this
  failed. You can ignore that line.

### Login

Core lets a browser in with its token. To get a link:

```bash
crow core url          # http://127.0.0.1:8787/?token=<token>
```

Open the link once. Core sets an HttpOnly cookie and drops the token from the address bar.
The cookie lasts about a year. To revoke every browser, delete `local/state/core.token` and
restart core: the next `crow core url` prints a new token.

Behind `tailscale serve`, the owner's tailnet login is enough and no token is needed.

## nest (front)

```bash
cd "$CROW_ENGINE/tools/nest"
npm ci
npm run build          # writes dist/
```

Core serves `tools/nest/dist` at `/`. To serve a build from somewhere else, set
`serve: {nest: <path>}` in `local/core.yaml`. Rebuild after each engine update. The build
does not depend on your config, so the same `dist/` runs against any core.

Nest reads `/api/modules`. A page shows when at least one of its sections is on. A brick's
generic list sections sit on the page they name, after its panels.

## WSL from Windows

WSL2 forwards `127.0.0.1` to Windows. Open the `crow core url` link in a Windows browser.
If it does not load, `wsl hostname -I` gives the WSL address. Then set
`serve: {listen: 0.0.0.0:8787}` in `local/core.yaml` and open
`http://<that address>:8787/?token=<token>`. Only do this on a network you trust: the
listener is plain HTTP.

## Check

```bash
crow doctor                                   # identity, accounts, bricks
curl -s --unix-socket "$CROW_HOME/local/state/core.sock" http://core/api/modules | head
```

The second command lists every brick with `on`, its core sections and its Nest sections.
