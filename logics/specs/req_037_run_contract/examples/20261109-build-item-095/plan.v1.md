# Plan v1: scope-aware account pick

The dispatcher picks the first allowed account today. This plan reads the scope class of the target repo and keeps only accounts allowed for it.

## Steps
1. Read the scope from the project registry in `crow_lib`.
2. Filter the account pool by scope before the existing failover order.
3. Add tests for work, personal and crow scopes.

## Open point
A repo with no scope in the registry. Two ways: refuse the dispatch, or treat it as personal.

```crow
contract: 1
mission: 20261109-build-item-095
kind: build
phase: plan
version: 1
summary: Filter the account pool by repo scope before failover.
confidence: 70
verdict: needs-input
questions:
  - id: q1
    text: "A repo has no scope in the registry. The dispatcher cannot tell which accounts may run it. What should it do?"
    blocking: true
    choices:
      - {id: refuse, label: Refuse the dispatch, note: Safe. The owner adds the scope then retries.}
      - {id: personal, label: Treat it as personal, note: Runs at once. Risk of using a personal account on work code.}
    default: refuse
    recommended: refuse
next: Implement the 3 steps on the mission branch, run the checks, write the report.
publish: []
```
