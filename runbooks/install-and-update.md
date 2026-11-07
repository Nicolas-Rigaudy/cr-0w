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

## Install
`sh install.sh [--version X.Y.Z] [--prefix DIR] [--bin-dir DIR]`
- Downloads `crow-<version>.tar.gz` and `SHA256SUMS` from the GitHub release (latest when no `--version`).
  A checksum mismatch aborts before anything is written.
- Engine: `<prefix>/versions/<version>`, `<prefix>/current` points at it (prefix default `~/.local/share/crow`).
  `<bin-dir>/crow` (default `~/.local/bin`) links to `current/bin/crow`. `.crow-install` in the engine records the method.
- Data: `local/` in the data dir gets the `docs/*.example.yaml` files that are missing. Existing files stay as they are.
- Release assets it expects (item_035 builds them): `crow-<version>.tar.gz` with a top dir `crow-<version>/`, no `local/`,
  and `SHA256SUMS` in `sha256sum` format.

## Update
`crow update [--check] [--to X.Y.Z]`
- Package install (`method=package` in `.crow-install`): prints the package manager to use.
- Script install: runs the current engine's `install.sh --engine-only` for the target, then prints the changelog entries
  between the two versions. Old versions stay under `<prefix>/versions/` for a manual rollback (`ln -sfn versions/<old> <prefix>/current`).
- Git checkout: refuses uncommitted changes, fast-forwards the branch to the tag (or checks out the tag on a detached HEAD),
  prints the changelog entries in between. Local commits on the branch stop it.
- `--check` reads the release and changes nothing. No method ever writes to the data dir.

## Restore on a new machine
1. `sh install.sh --restore <git url of the private data repo>` clones it into `<data>/local`. It refuses a non-empty `local/`.
2. Run the steps it prints. Git does not carry them:
   - keyring tokens: `crow secret set <service>` for each entry of `tokens.yaml` (runbooks/service-tokens.md);
   - cdx accounts: `cdx add <name> && cdx login <name>` for each account in `accounts.yaml`, then `cdx set` options;
   - the vault: clone or sync it, then check `vault.path` in `projects.yaml`.
3. `crow doctor` checks the rest, including the data repo (remote, uncommitted changes, unpushed commits).
- Keep the data repo light: no generated state or large files (`state/` caches, renders). Commit and push it regularly.

## Test
- `python3 -m pytest tests/test_install_update.py`: install, checksum mismatch, update per method, data byte for byte, restore.
- `tests/smoke/install-in-docker.sh`: the same round trip in a clean container (default `python:3.13-alpine`).
