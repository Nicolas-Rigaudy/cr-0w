# Terminal backends and platforms

Crow opens, focuses and types into terminal windows through one module, `bin/crow_term.py`.
`CROW_TERMINAL` picks the backend. Unset, crow takes the first that fits: `kitty` when it is on PATH,
else `tmux` when running inside tmux, else `print`.

| Backend | Opens a window (`crow term`, `crow handoff`, resume in `crow focus`, `crow runs relaunch`) | Focuses a session (`crow focus`, deck jump) | Types an answer (deck reply, `crow focus --send`) |
|---|---|---|---|
| `kitty` | new kitty window (`setsid -f` where it exists) | remote control on a socket matching `$XDG_RUNTIME_DIR/kitty-*` (kitty.conf: `allow_remote_control socket-only`, `listen_on unix:${XDG_RUNTIME_DIR}/kitty-{kitty_pid}`) | yes |
| `tmux` | new tmux window in the current server | selects the pane whose shell started the session | yes (`send-keys`) |
| `print` | prints the command to run in a new terminal | no | no (exits 1) |

`crow runs relaunch --here` runs the handoff in the current terminal and needs no backend.

## Feature by platform

Linux is the tested platform. macOS and Windows through WSL are not covered by CI yet; the cells
below say what the code needs, not what was run there.

| Feature | Needs | Linux | macOS | Windows (WSL) |
|---|---|---|---|---|
| Hub, `crow open`, briefs, dispatch, runs, usage | Python 3 + PyYAML, `jq`, cdx | yes | yes | yes |
| `crow term`, `crow handoff`, `crow runs relaunch` | a terminal backend | kitty, tmux or print | kitty, tmux or print | tmux or print |
| `crow focus`, deck jump and reply | kitty or tmux backend | yes | kitty or tmux | tmux |
| Live agents (`crow json`, watch and deck Agents) | `/proc` | yes | no | yes |
| `ansi-shot` | Pillow, fontTools, a monospace font (`--font`) | yes | yes, with `--font` | yes |
| Run notifications (`crow runs` wait, `crow-ctx`) | `notify-send` | yes | no | with a notification daemon |
| `crow-open` (links with the real HOME) | `xdg-open` | yes | no | with wslu |
| `crow secret` | `secret-tool` (libsecret keyring) | yes | no | with a keyring daemon |
| Copy an id in `crow watch` | `wl-copy` (Wayland) | Wayland only | no | no |
