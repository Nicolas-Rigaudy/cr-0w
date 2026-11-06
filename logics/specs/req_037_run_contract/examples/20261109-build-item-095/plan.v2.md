# Plan v2: scope-aware account pick

## Responses
1. q1: the dispatch refuses a repo with no scope and names the fix.
2. `crow doctor` lists repos with no scope. Added as step 4.

## Steps
1. Read the scope from the project registry in `crow_lib`.
2. Filter the account pool by scope before the existing failover order.
3. Refuse a repo with no scope, with the registry line to add.
4. List repos with no scope in `crow doctor`.
5. Add tests for work, personal, crow and missing scopes.

```crow
contract: 1
mission: 20261109-build-item-095
kind: build
phase: plan
version: 2
summary: Filter accounts by repo scope, refuse unscoped repos, flag them in doctor.
confidence: 85
verdict: ready
responses:
  - {ref: a1, outcome: done, note: Refuse with the registry line to add.}
  - {ref: a2, outcome: done, note: Step 4.}
next: Implement the 5 steps on the mission branch, run the checks, write the report.
publish: []
```
