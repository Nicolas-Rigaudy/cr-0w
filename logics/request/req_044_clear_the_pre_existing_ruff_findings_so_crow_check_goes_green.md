## req_044_clear_the_pre_existing_ruff_findings_so_crow_check_goes_green - Clear the pre-existing ruff findings so crow-check goes green
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 100%
> Confidence: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Fix the pre-existing `ruff check .` findings in bin/ and tests/ so crow-check passes again.
- Keywords: clear, pre, existing, ruff, findings, crow, check, goes, green
- Use when: crow-check is red on ruff, or before a release.
- Skip when: Go code (tools/), Logics docs.

# Needs
- `ruff check .` reports about 90 findings across 25 files (seen 2026-11-09): PLW1510 subprocess without check (bin/crow_lib.py x12, many tests), F811 redefinitions (tests/test_crow_log.py x7), RUF059, SIM115, FURB167 (bin/crow_vault.py, bin/crow_lib.py), I001, F401, re.I aliases.
- crow-check stays red, so a new finding goes unnoticed.

# Context
- Found while closing item_131; left out of that work as a separate chore.
- Fixes must keep behaviour: PLW1510 needs an explicit `check=` per call, chosen from what each call expects.

# Acceptance criteria
- AC1: `ruff check .` reports no findings.
- AC2: The Python test suite passes unchanged.
- AC3: No rule is silenced in config to get there, except a noqa with a reason where the code is right.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow_lib.py`, `bin/crow_vault.py`, `tests/test_crow_log.py`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |


# Outcome
- 2026-11-09: `ruff check .` clean, 502 tests pass, crow-check green.
- Safe and unsafe autofixes (imports, UTC alias, noqa, regex flags, unused unpacked names, dict literals).
- `vault` fixture and `TPL` moved from tests/test_crow_note.py to tests/conftest.py (fixes F811 in tests/test_crow_log.py).
- Bare `open().read()` in tests replaced with `Path.read_text`/`write_text`.
- Every flagged `subprocess.run` got `check=False`: each site already reads the return code or tolerates failure, so behaviour is unchanged.
- No rule silenced in config.
