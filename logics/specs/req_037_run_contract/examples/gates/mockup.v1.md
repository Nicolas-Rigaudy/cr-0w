# Mockup v1: mission card in the inbox

Gate shape `[mockup]` (kind design). Same scenario framed the same way in each option. Answer `pick A`, `change ...` or `approve`.

```crow
contract: 1
mission: 20261109-design-mission-card
kind: design
phase: mockup
version: 1
summary: Three card layouts for a mission waiting on a gate.
confidence: 60
verdict: needs-input
artifacts:
  - {path: outbox/a-gate.png, kind: image, option: A, scenario: waiting-on-gate}
  - {path: outbox/b-gate.png, kind: image, option: B, scenario: waiting-on-gate}
  - {path: outbox/c-gate.png, kind: image, option: C, scenario: waiting-on-gate}
questions:
  - id: q1
    text: "Three layouts for a mission card that waits on your answer. Which one goes forward?"
    blocking: true
    choices:
      - {id: a, label: "A: one line", preview: outbox/a-gate.png}
      - {id: b, label: "B: two lines with the verdict", preview: outbox/b-gate.png}
      - {id: c, label: "C: card with buttons", preview: outbox/c-gate.png}
publish: []
```
