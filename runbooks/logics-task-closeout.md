# Close a Logics orchestration task

For a task whose items are all Done but whose closeout preflight fails.

1. `logics-manager flow validate-closeout <task>`: each finding names its fix.
2. Check every listed item is `Status: Done` before running `flow repair gates`. The repair ticks all Plan/DoD boxes without checking.
3. AC traceability on a slice task: the validator wants `- request-ACn -> This task. Proof: ...` for every AC of every request linked through its items. The lines are keyed by AC number only, so one line serves all linked requests. Write the lines by hand. For each request, state what this task proved and name the owner of the open part (later task or item). Never claim an AC this slice did not meet. `flow repair ac-traceability` skips ACs that already carry a "deferred" line.
4. `flow repair links <task>` for missing back-links.
5. Validation: write concrete lines in `# Validation` first (`<command> passed on <date>: <result>`). `flow closeout --validation "..."` runs its preflight before writing the text and rolls back.
6. `logics-manager flow closeout <task> --index --lint --audit`.
7. Review `git diff` before committing. The repairs can rewrite a product brief header: they shrink `Related request` to one request and set `Status: Settled`. Restore both when the brief covers more than this task.
8. Linked requests stay Draft while any AC is open.
