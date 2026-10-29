# Design canvas: add boards without losing the owner's layout

A "Design" artifact (canvas) holds boards as `project/<name>.dc.html`, indexed by `project/canvas.json`.
The owner moves boards in the browser, so the published index is the truth, never the local copy.

## Index format (v3)

- `boards`: `{"<name>.dc.html": {x, y, w, h, title}}`. Keys have no `project/` prefix.
- `order`: list of board keys.
- `notes`: `{id: {fill, w, x, y, text}}`, `noteOrder`: list of ids.

Each board keeps the `<script src="./support.js">` line and the `data-dc-script` block with `$preview` width and height.

## Steps

1. Generate boards with a script (copy a base board, swap CSS or colors with string replaces, set `<title>`). Run Python with `-I`. Keep generators in the scratchpad, not in the repo.
2. Count `<div` and `</div>` after any structural edit; an unbalanced tag breaks the board.
3. Publish one board per call: Artifact publish with `url`, `root` = the local canvas folder, `file_path` = absolute path of the board, `files {"project/X.dc.html": "project/X.dc.html"}`. Several calls can go in parallel.
4. Right before touching the index, read the live one: Artifact read with `path: project/canvas.json` and an `out_dir` (e.g. `<canvas>/live`).
5. Merge only your boards and notes into that live file (check for overlaps with existing boxes), write it to `project/canvas.json`, publish it the same way.
6. Never overwrite an original board: new ideas are copies with a new name, placed next to what they change.
7. Do not verify with screenshots. Open the link with `bin/crow-open <url>` (real browser).

## Showing choices

Draw every visual option as its own board beside the original, add a blue note naming the boards, then ask with AskUserQuestion. Say when the data is invented.
