## item_099_core_publishes_and_holds_every_forge_credential - Core publishes and holds every forge credential
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Safety
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: core is the only process that pushes, opens MRs or PRs, merges and posts; it holds every forge, tracker and chat credential and publishes only after a hash-chained validation record or the personal merge gate.
- Keywords: core, publishes, holds, forge, credential
- Use when: Building the publish path of core or moving a credential out of agent reach.
- Skip when: Work on the run sandbox itself (isolated runs item).

# Problem
- Agents today run with the owner's environment and could push or post with the owner's tokens.
- adr_001 settles that core publishes; the spike defines the validation record and the credential split.

# Scope
- In:
  - Publish actions: push, MR or PR create and update, tracker write, wiki draft or publish, mail draft, chat post.
  - Validation record `{mission, version, action, target, payload_sha256, owner_event, ts}` in the event log, chained by hash.
  - Payload hash recheck at publish time; a newer owner message cancels the record.
  - Credential table from the spike: forge, tracker, wiki, mail and chat tokens live with core only.
  - `pause all publishing` switch honoured by every action.
- Out:
  - Work merges (always refused).
  - The personal auto-merge gate checks (own item).

# Acceptance criteria
- AC1: Every publish action goes through one core function that refuses without a matching validation record or gate pass.
- AC2: A changed payload or a newer owner message after `ok` cancels the record; a test covers both.
- AC3: A replayed or edited record fails the hash chain check.
- AC4: The pause switch blocks every publish action and logs the refusal.
- AC5: No mission run environment or mount contains a forge, tracker, wiki, mail or chat credential (probe test).

# AC Traceability
- request-AC2 -> credential ownership and core publish path.
- request-AC3 -> core publishes what the sandboxed run produced.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_030_spike_mission_safety_and_sandbox.md`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Every other safety layer assumes core holds the keys.

# Notes
- Proposed by the req_030 spike (2026-11-06); Draft until the owner grooms it.
- Depends on the core event log (adr_001).
- Mail stays a draft whatever the record says.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

