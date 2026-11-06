# Plan v1: offline mode for the mobile app

Gate shape `[plan]` (kind plan). `go` commits the request and backlog docs on the mission branch.

## Items
1. Cache the last trip list. 2. Queue actions while offline. 3. Sync on reconnect.

```crow
contract: 1
mission: 20261109-plan-offline-mode
kind: plan
phase: plan
version: 1
summary: One request, three backlog items, none Ready (two wait on a product choice).
confidence: 75
verdict: ready
artifacts:
  - {path: outbox/req_draft.md, kind: doc}
publish:
  - action: commit-docs
    target: "<repo> mission/20261109-plan-offline-mode"
    title: "new[MD]: plan offline mode"
next: Commit the docs after go.
```
