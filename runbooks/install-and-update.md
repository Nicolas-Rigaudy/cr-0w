# Install, update and data dir

## Engine and data
- Engine: the scripts, skills and runbooks. `bin/crow` finds it from its own location (`CROW_ENGINE`).
- Data: the owner's `local/` (registry, accounts, state, memory). Its parent is `CROW_HOME`.
- `CROW_HOME` resolution: the variable when set; else the checkout when it has `local/`; else
  `$XDG_CONFIG_HOME/crow` (default `~/.config/crow`), with the data in `~/.config/crow/local/`.
- `crow version` prints the version, the engine path and the data dir.

## Migrating a checkout install
Nothing moves by itself. A checkout with `local/` keeps working as before.
To keep the engine and the data apart:
1. `mkdir -p ~/.config/crow && mv <checkout>/local ~/.config/crow/local`
2. `crow version` shows `data: ~/.config/crow`. Set `CROW_HOME` instead when the data lives elsewhere.
3. Hooks that read `local/` from the checkout (`crow-commit-msg-check`, `crow-prepush-check`,
   `crow-precommit-check`) still expect it there: keep a symlink `<checkout>/local -> ~/.config/crow/local`.
