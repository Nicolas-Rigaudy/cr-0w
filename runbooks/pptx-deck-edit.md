# Edit an image-heavy pptx deck (texts, icons, new slides)

For decks exported from a design tool: each slide is one full-bleed background PNG plus editable text boxes. Icons and bullets are baked into the PNG.

## Setup
- Work dir: `local/state/shots/<topic>/`. Copy the source pptx to `src.pptx`; never touch the original download.
- Dump texts once: shape id, position, runs, font sizes per slide (python-pptx). Slide numbers in scripts stay the ORIGINAL ones, even after inserts and deletes.

## Text edits: `apply.py src.pptx out.pptx`
- Helpers: `lines(n, sid, [..])` replaces paragraphs keeping the first run's format; `runs(n, sid, [..])` keeps per-run colours; `para(n, sid, i, t)` sets one paragraph.
- Text that overflows: shorten the wording first, then widen the box, then lower the size.
- New slides: `add_slide` with an existing layout, native shapes, the deck's colours and font. Insert/delete by moving `sldIdLst` entries at the end of the script.

## Icons baked in the background
- Extract each background (`slide.part.related_part(blip rId).blob`) to `bg/s{n}.png`.
- `patch_icons.py` reads `bg/s{n}.png`, writes `bg/p{n}.png`: erase the old glyph with the sampled background colour, paste a new one.
  - Simple glyphs: Lucide SVGs (unpkg `lucide-static`), recoloured, rendered with `rsvg-convert`.
  - Illustrated icons: a small hand-written SVG in the slide's palette (outline colour, fills).
- Measure positions with a colour bbox scan in PIL, not by eye. Text boxes are positioned in EMU: px = emu / slide_width_emu * image_width.
- `apply.py` swaps the blob: `slides[n-1].part.related_part(rid)._blob = open('bg/p{n}.png','rb').read()`.
- A stray character that is not in any text box is baked in the background: erase it there.

## Harmonise text
- Install the deck font first (also into the sandboxed `$HOME/.local/share/fonts`, then `fc-cache`); `fc-match <font>` must resolve or renders fall back and hide overflow.
- Last pass before `prs.save`, on final slide order: title = shape with the largest run; one weight, two sizes (hero / long one-line), two colours (dark + accent); dark banners their own size.
- Bullets baked in the image: move the text box left edge after the bullet, same gap in every column.
- If the text size changes, baked bullets stop matching the lines: erase them (row blend of the pixels left and right of the box) and add real bullets (`buClr`/`buSzPct`/`buFont`/`buChar` after `lnSpc`/`spcBef`/`spcAft` in `a:pPr`, `marL` = -`indent`).
- Erasing an old glyph: pad each side separately (left, top, right, bottom); a single pad leaves remnants on the side where the glyph is bigger.
- Re-render the whole deck and zoom any title that wraps; widen its box rather than shrinking the font.

- One type scale for the whole deck (e.g. 12/14/16/18/22/26/32): snap every run to the nearest step at the end of the script, keep big numerals, scale explicit line spacing by the same ratio.
- After any size change, re-check every title and label: grown text wraps or runs under images.
- `wrap="none"` with `spAutoFit` still wraps in LibreOffice: remove `spAutoFit` from no-wrap boxes.
- Replacing text in bullet paragraphs: keep the bullet run and the text run separate, or the whole line takes the bullet colour.
- Delete slides as the very last step and keep the script on the old numbering (say so in its docstring).

## Check
- `soffice --headless --convert-to pdf out.pptx`, `pdftoppm -r 50` for a grid, zoomed crops (`magick -crop`) for each patched icon.
- Zoom straight from the pdf: `pdftoppm -r 150 -f N -l N -x X -y Y -W w -H h -png deck.pdf z` writes `z-NN.png`. Coordinates are in pixels at that dpi (13.33 in wide slide = 2000 px at 150).
- Small text boxes stacked on an icon (labels like "TVS / Taxe") overflow in LibreOffice: drop them when the slide text already says the same.
- Before/after for the owner: render the same slides from `src.pdf` and the new pdf into the vault research folder and build a `bin/crow-compare` page.
- Copy the deck to the vault next to the meeting notes.

## AI photo edits
- `codex exec --skip-git-repo-check --sandbox workspace-write -i in.png -- "prompt"`: the `--` is required because `-i` takes several values.
- Needs codex logged in under the sandboxed `$HOME/.codex` (401 "Missing bearer" otherwise).
- Edit a crop, not the full background, then composite it back with a feathered mask. Show before/after before integrating.

## Screenshot slides
- Tall screenshots are height-limited: grow the slot height first. Widening the slot only helps wide ones.
- For a very wide screenshot, trim an edge (no key content) so it fills the same height as its neighbours.
- Check the logo tab on slides whose background already has one baked in: its size can differ from the standard one. Paint it out and place the standard tab.
