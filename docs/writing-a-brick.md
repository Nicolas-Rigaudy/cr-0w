# Writing a brick

A brick is a set of scripts, skills and runbooks that crow turns on or off as one unit. The engine ships its own bricks in `modules/<brick>.yaml`. You can add your own without patching the engine.

## Where it lives

```
<CROW_HOME>/modules/<name>/        # local/modules/<name>/ when CROW_HOME is the engine checkout
  module.yaml                      # the manifest (required)
  bin/crow-<name>                  # scripts, executable
  skills/<skill>/SKILL.md          # Claude skills
  runbooks/<runbook>.md
```

`crow modules` lists engine and user bricks. User bricks are marked `[user: <dir>]`, with the pinned ref when installed from git.

## Manifest

The schema is the same as the engine manifests (`modules/forge.yaml` is a good example). Every field is optional.

```yaml
summary: One line shown by crow modules.
default: true            # on when crow.yaml does not say; user bricks default to on
requires: [git]          # bricks it needs; it stays off while one of them is off
scripts: [crow-hello]    # files in bin/; `crow hello` runs bin/crow-hello
skills: [hello]          # skills/<name>, linked into every allowed account profile
hub_skills: []           # skills/<name>, linked into the hub session only
runbooks: []
gather: {}               # <routine>: [keys] this brick owns in crow-gather output
doctor:
  needs: [jq]            # tools that must be on PATH (needs_any: one of them)
config:
  greeting: hi           # defaults; crow.yaml can override each key
```

`name` defaults to the folder name. When it is set, it must match the folder name.

## Switch and config

```yaml
# crow.yaml
version: 1
modules:
  hello: false             # off
  # or keep it on and override config:
  # hello: {greeting: salut}
```

After a change, run `crow modules link` to relink skills.

## Script contract

- `crow <cmd>` runs `bin/<cmd>` or `bin/crow-<cmd>` from an on user brick when no engine command has that name. Engine commands always win.
- The script gets `CROW_HOME` and `CROW_ENGINE` in its environment. Read config with `crow modules json`. Each brick row has `config`, `on` and `dir`.
- Exit 0 on success. Exit 2 when a precondition is missing. Errors go to stderr.

## Install from git

```
crow modules add <git url> [--ref <tag|sha>] [--name <name>]   # clone, check, pin, turn on, link skills
crow modules update [<name>] [--ref <tag|sha>]                 # newest tag (or the ref), lock rewritten
crow modules remove <name>                                     # clone, lock entry and crow.yaml key gone
```

- Without `--ref`, `add` and `update` take the newest tag, or the default branch when there is no tag. Tag your releases (`v1`, `v2`...).
- The pin goes to `modules.lock`, next to the modules dir. crow writes this file, so don't edit it.
- An update with a broken manifest is refused, and the clone goes back to the locked commit.
- `crow doctor` fails when a locked clone is missing, has a HEAD that differs from the lock, or has local changes.

## Rules

- A brick named after an engine brick (`forge`, `vault`...) is refused. Rename its folder.
- A folder in the modules dir without `module.yaml` is a `crow doctor` finding.
- If a user skill has the same name as an engine skill, the engine skill is linked.
