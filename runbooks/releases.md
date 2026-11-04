# Releases: version, changelogs, bump check

`VERSION` at the root is the single source (semver, from 0.1.0). `bin/crow-check` (group `version`) fails when:
- VERSION is not semver;
- another carrier (`package.json`, `pyproject.toml`, `tools/*/package.json`, `tools/*/VERSION`) names a different version;
- `changelogs/CHANGELOGS_X_Y_Z.md` for VERSION is missing;
- `CHANGELOG.md` differs from what `bin/crow-changelog` builds.

## Cutting a release (the owner approves each step)
1. `bin/crow release propose`: reads commits since the last tag, offline, prints the version and notes by `new`/`chg`/`fix`. `new` bumps the minor, `chg`/`fix` the patch.
2. After approval: write `changelogs/CHANGELOGS_X_Y_Z.md` in plain words, set `VERSION`, run `bin/crow-changelog`, run `bin/crow-check`.
3. Commit. Tagging and publishing are separate and need an explicit OK (item_035).
