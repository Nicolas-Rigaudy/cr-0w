# Report v3: scope-aware account pick

## Responses
No answers since plan v2 (go).

## Done
- The dispatcher filters accounts by repo scope before failover.
- A repo with no scope is refused with the registry line to add.
- `crow doctor` lists repos with no scope.
- 9 new tests. Full checks green (`outbox/checks.log`).

## To publish after ok
Push the mission branch and open one merge request.

```crow
contract: 1
mission: 20261109-build-item-095
kind: build
phase: report
version: 3
summary: Scope-aware account pick done, checks green, ready to push.
confidence: 90
verdict: ready
responses: []
artifacts:
  - {path: outbox/checks.log, kind: data, caption: Full check run}
proof:
  - {what: command, cmd: bin/crow-check, exit: 0}
  - {what: file, path: outbox/checks.log}
publish:
  - action: push
    target: "<repo> mission/20261109-build-item-095"
    commits: [3f2a9c1, 8b41d07]
  - action: mr
    target: "<repo> mission/20261109-build-item-095 -> main"
    title: "new[PY]: pick the dispatch account from the repo scope"
    body: "Filters the account pool by repo scope before failover. Refuses repos with no scope. Lists them in crow doctor."
```
