# Reply drafts in the mailbox (crow-mail-draft)

`bin/crow-mail-draft` saves one reply draft per input file, reads it back and checks it. It never sends: no backend has a send path.

## Setup
1. Copy `docs/mail.example.yaml` to `local/mail.yaml`. Pick `draft_backend`:
   - `graph`: Microsoft Graph with the `crow-today` sign-in (`runbooks/graph-today.md`). Drafts need the `Mail.ReadWrite` scope.
   - `com`: classic Outlook on Windows, driven from WSL. Outlook must be running. Nothing to sign in.
   - `connector`: no local mailbox access. The script queues a request for a hub session that has the M365 connector.
2. Fill `internal_domains` and `signatures` (one list of lines per kind). A line can be `{image: PATH, width: PX, alt: TEXT}`: the image goes in as an inline attachment (`cid:`), PATH relative to the crow home, 3 MB max. Graph cannot read Outlook's own signature, so keep a PNG copy and replace it when the Outlook signature changes.
3. `com` only: `bin/crow-mail-draft probe` must print `"connected": true`.

## One draft
Write a draft file (YAML or JSON):
```yaml
message_id: <id of the mail to answer>
mode: all              # all (reply all) or sender
signature: auto        # auto | internal | external | none
text: |
  Hello,
  The reply.
attachments: [report.pdf]   # relative to the draft file
cc: []
```
`html:` replaces `text:` for formatted bodies. `bin/crow-mail-draft save draft.yaml` writes the draft and prints the checks; `--json` for the payload. Exit 1 means a check failed: open the draft in Outlook before anything else.

## Checks after the write
| Check | Fails when |
| --- | --- |
| body_visible | the reply text is not in the draft |
| signature_once | the configured signature (text and each image) appears zero or several times in the reply part (`#crow-reply`); without that wrapper, outside the quoted thread |
| no_native_signature | Outlook's automatic signature (`_MailAutoSig`) is still in the draft |
| attachments | an attachment named in the draft file is missing |
| still_draft / not_sent | the message left Drafts or carries the sent flag |

Classic Outlook `ReplyAll()` inserts the owner's automatic signature. The COM bridge removes it before inserting the reply; the two signature checks catch a draft where that did not happen.

## Connector backend
`save` writes `local/state/mail/connector/<stamp>.json` with the HTML, the tool to call and the readback path. A hub session creates the draft with the connector, adds `inline_images` as inline attachments with their `cid`, writes the readback JSON (`draft_id, html, body_text, attachments, is_draft, sent`) and runs `bin/crow-mail-draft check <draft file> --readback <readback>`. `signature: auto` needs a `recipients:` list in the draft file here.

## Rules
- A write whose answer is lost is reported `uncertain: true` and never retried. Look in Drafts before a new attempt.
- `delete-draft <id>` refuses anything outside the Drafts folder.
- Proofs land in `local/state/mail/drafts/` (expected values, readback, checks).
