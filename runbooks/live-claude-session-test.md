# Live-test a hook or key press on a real Claude session

Use when a feature reacts to a live Claude Code session (hooks, permission prompts, key presses) and unit tests are not enough.

1. Make a scratch folder in the session scratchpad (never /tmp, never a repo).
2. Start a hidden session there: `tmux new-session -d -s <name> -c <dir> 'claude --model haiku --permission-mode default'`. No visible window.
3. Accept the folder trust dialog: `tmux send-keys -t <name> Down Enter`, then wait for `for shortcuts` in `tmux capture-pane -p -t <name>`.
4. Type the prompt with `tmux send-keys -t <name> -l '<text>'` then `tmux send-keys -t <name> Enter`. Ask for one exact harmless command (`date +%s > stamp.txt`).
5. Wait on a condition, not a fixed sleep: `timeout 40 bash -c 'until <check>; do sleep 1; done'`.
6. Check the result through the feature's own CLI, then on disk (the stamp file).
7. Scripts that pick a terminal backend default to kitty: set `CROW_TERMINAL=tmux`.
8. Clean up: `tmux kill-session -t <name>` and remove state files the test wrote.

Notes
- Permission menu: 1 Yes, 2 always allow, 3 auto mode, 4 No. `1` allows once; Escape denies and ends the turn.
- Claude's own status file `<HOME>/.claude/sessions/<pid>.json` shows `waiting` / `permission prompt` while the dialog is up.
- Never `pkill -f` a pattern that matches your own command line.
