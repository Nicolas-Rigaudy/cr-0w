## item_098_monthly_budget_from_bank_exports - Monthly budget from bank CSV exports in the personal vault
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Import the bank's CSV export into the personal vault, categorise operations, and produce a monthly budget vs actual review with alerts.
- Keywords: butler, budget, bank, csv, import, categories, monthly review
- Use when: building this butler slice.
- Skip when: work rituals on the work account.

# Problem
- No budget view exists. The bank offers a CSV export of operations; the owner wants crow to do the sorting and the monthly review.

# Scope
- In:
  - `bin/crow-budget import <csv>`: parse the bank export, deduplicate against earlier imports, store normalised rows next to the personal vault.
  - Categorisation by rules (payee patterns) kept in a rules file, with Claude suggesting rules for uncategorised rows and the owner approving them.
  - Monthly review note in the personal vault: budget vs actual per category, top changes, recurring payments, alerts on overspend.
  - Source-agnostic input layer so a DSP2 aggregator can replace the manual export later.
- Out:
  - Payments or transfers of any kind.
  - Stored bank credentials or scraping the bank site.
  - The aggregator itself (later, if manual export becomes a chore).

# Acceptance criteria
- AC1: Importing the same export twice adds no duplicate rows.
- AC2: Every imported row gets a category from the rules or lands in an uncategorised list; new rules are only added after the owner approves them.
- AC3: The monthly review note shows budget vs actual per category and flags categories over budget.
- AC4: Budget alerts feed the personal brief as a short line without amounts outside the personal vault.
- AC5: Parsers are tested on a redacted sample export; no real statement is committed.

# AC Traceability
- AC1 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC2 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC3 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC4 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC5 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_026_personal_butler`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Most valuable butler slice long term, but needs the export check on the owner's account first.

# Notes
- Grooming check before build: confirm the export button and columns on the owner's account page (bank specifics in local/context.md).
- Generic parser code in bin/; bank-specific column mapping in local/.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
