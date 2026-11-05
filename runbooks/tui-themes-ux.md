# TUI themes that feel like different instruments

Scope: how to design themes for a terminal dashboard (Go, lipgloss) so each theme is a different instrument and not a recolour. Generic. Written for the crow deck redesign.

## (a) Diagnosis of the current themes

Source: `tools/crow-deck/theme.go`, `chrome.go`, `frame.go`, and the views.

A theme today is: a frame kind (`round`, `hud`, `cyber`, `soft`), one default accent slot, four role slots (title, key, label, head) and one flag (`caps`).

What a theme changes today:

| Area | lazytf | hud | cyber | soft |
|---|---|---|---|---|
| Box frame | `╭─ Title ─╮` round, `│` sides | `┏━ ◢ TITLE ◣`, `▏` left edge only, `━━┛` corner, panel index on the bottom edge | `┌─ ▌// TITLE // ─┐`, `◢` corner | none, a title row only |
| Title | on the border, magenta | caps tag, accent | caps, `▌//` | sentence case, bold |
| Focus | border accent vs ANSI 8 | tag accent, left edge accent | border accent, `◢` accent | title accent only |
| Brand | `◆ crow deck` | `◆ CROW//DECK ░▒▓` | `▌CROW_DECK // page ░▒▓` | `Crow deck` |
| Tabs | underline on active | `▶` marker | `▌NAME╱` | spaces plus underline |
| Rule under tabs | `─` | `━━━━ ╌╌╌ ░░` then `╌` | `//////////` then `─` | blank line |
| Status bar | plain | caps descriptions | caps plus `//` between hints | plain |
| Table header | `  ` prefix | caps | caps, `//` prefix | plain |
| Layout gaps | 0 | 0 | 0 | 3 columns, 1 row |
| Caps | no | yes | yes | no |

What stays identical in all four themes:

- Selected row: accent plus reverse video. Same fill, same behaviour.
- Gauge: `━` fill in the accent, `━` track in slot 8. Same glyphs everywhere.
- State glyphs: `●`, `○`, `◆`, `‖`. Dots in the header. Same.
- Semantic colours: green ok, yellow warn, red bad. Same.
- Mode chip, live switch, key hints order, `syncing…`, empty text, sparklines.
- Row anatomy, density (one line per item) and column layout.
- Box proportions and page layout.

Why they feel half-baked:

1. Themes only swap the frame. The content inside the frame is the same, so the screen reads the same at a glance.
2. Differences are decoration (`//`, `░▒▓`) added on top of one skeleton, not a different skeleton.
3. Two themes share the accent family (HUD and Soft are both cyan by default).
4. Unfocused borders use slot 8, which nearly vanishes on a translucent terminal.
5. Soft removes all lines. Without fills there is nothing left to delimit a box, so panels melt into each other.
6. No theme owns selection, gauges, density, empty states or the status bar. Those are where an instrument is felt.

## (b) Apps studied and the pattern taken

Sources were fetched on 2026-11-02. A line marked "from memory" was not verified on the page.

| App | What it teaches | Source |
|---|---|---|
| lazygit | `gui.border` is a first-class choice: rounded, single, double, hidden, bold. Active and inactive border colours are separate. Selected line can be bg, bold only, or reverse. Optional Nerd Font icons. Panel jump keys are shown in titles. List footer `5 of 20`. | https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md |
| btop | Graph symbols are a choice per box: braille, block, tty. Rounded corners is a switch. A TTY mode restricts to 16 colours and ANSI symbols. `theme_background = false` gives terminal transparency. Box drawing needs U+2500 to U+259F. Gradients (start, mid, end) colour graphs. | https://raw.githubusercontent.com/aristocratos/btop/main/README.md and https://raw.githubusercontent.com/aristocratos/btop/main/src/btop_theme.cpp |
| k9s | Skins are colour only, but split by region: body, frame border (fg and focus), title (fg, highlight, counter), crumbs (active). Counter colour is a separate role. | https://k9scli.io/topics/skins/ |
| Helix | Theme is scopes with fg, bg, underline style and modifiers. Statusline has normal, insert, select variants and an inactive variant. Underline styles (line, curl, dashed, dotted, double) are used as emphasis. Window separators only get colour. | https://docs.helix-editor.com/themes.html |
| Textual (Posting, Harlequin, Elia) | A theme is about 11 base colours with generated shades and muted variants (70 percent blend). Borders only carry colour. Lesson: generated shades need fills, so they do not apply to a 16 slot, no-fill design. Posting ships 37 themes. Harlequin and Elia internals not verified. | https://textual.textualize.io/guide/design/ and https://github.com/darrenburns/posting |
| zellij | A theme has components: text, ribbon (tabs and mode), table, list, frame_unselected, frame_selected, frame_highlight (focused pane in a non-base mode). Three frame states, not two. | https://zellij.dev/documentation/themes.html |
| bottom (btm) | `border_colour` and `selected_border_colour`. `bar_type`: pipe, block, square. Separate table header and widget title styles. Three built-in themes. | https://clementtsang.github.io/bottom/nightly/configuration/config-file/styling/ |
| yazi | A flavor is a whole pack (theme.toml, syntax theme, preview). User overrides sit on top. Status bar separators and tab style are theme fields (from memory, not on the fetched page). | https://yazi-rs.github.io/docs/flavors/overview |
| superfile | Per-region colour keys: file panel, footer, sidebar, modal, each with a border and an active border. Gradient colour pair. Modal confirm and cancel colours. | https://superfile.dev/configure/custom-theme/ |
| gitui | `selection_fg`, `use_selection_fg`, and a configurable diff `line_break` symbol. Glyphs are partly themable. | https://github.com/gitui-org/gitui/blob/master/THEMES.md |
| gum, lipgloss | Border is a named set: normal, rounded, thick, double, hidden. Padding and margin are theme knobs too. | https://github.com/charmbracelet/gum |
| edex-ui | The sci-fi look comes from named themes (tron, blade) that change fonts, CSS and layout, plus boot animation. Not only colour. The HUD feel in a TUI must come from frame shape, caps and numbering since fonts are fixed. | https://github.com/GitSquared/edex-ui |
| crush, soft-serve, cmatrix | Not fetched. Not used as sources. | none |

Patterns taken:

1. Border form is a separate setting from border colour (lazygit, btop).
2. Frames have three states: unfocused, focused, focused in a special mode (zellij).
3. Regions have their own roles: title, counter, crumbs, header, status (k9s, superfile).
4. Selection style is its own setting: fill, bold only, reverse (lazygit, gitui, Helix).
5. Graph and bar glyph sets are a theme choice (btop, bottom).
6. A theme can ship icons, separators and tab style, not only colours (yazi).
7. A restricted mode of 16 ANSI colours is a legitimate design target (btop TTY).
8. A big look change changes layout and rhythm (edex-ui themes, lazygit screen modes).

### Transparent terminals

- btop: leave the theme background empty, `theme_background = false`.
- Helix: omit `bg` in `ui.*` scopes so the terminal shows through.
- lazygit: `selectedLineBgColor` can be `default` (bold only) or `reverse`.
- Rule for the deck: no background fills except the selected row. Delimit with line glyphs, not tints. Secondary text on slot 7, never 8. Slot 8 only for leaders and gauge tracks that carry no meaning. Distinguish states by glyph shape and weight (round, heavy, double), because colour contrast shifts with the wallpaper.

## (c) Axes a theme should own

A theme is a table with these axes. Everything uses the 16 ANSI slots (0 to 15) and plain glyphs. No fills except an optional reverse-video selection and the mode chip.

1. Border set: round, square, heavy, double, corners only, rules only, spine only.
2. Unfocused border: glyph weight and slot.
3. Focused border: glyph weight, slot, and form change (for example single to double).
4. Corner and edge ornaments: notch, bracket, tick, index number.
5. Panel title style: inline on the border, row above, tab-like, bracketed, with a glyph, case.
6. Title right label and bottom label: position, separator, case.
7. Counter style in titles: `Findings 3`, `(3)`, `[03]`, dim count.
8. Panel spacing: inner padding, gap columns, gap rows between boxes.
9. Header and brand: text, glyph, trailing ornament.
10. Tab bar: active marker (underline, chevron, bar, bracket), inactive style, separators, badge style.
11. Rule under tabs: glyph and slot, or none.
12. Mode switch and live switch: chip, brackets, or bold word.
13. Selection marker: reverse fill, bar, chevron, bold plus colour, underline.
14. Selection on unfocused panes.
15. List density: row height, blank lines between groups, group header style.
16. Table header: case, prefix, underline rule or none.
17. Separators inside panels: rule glyph, dotted, none, spacing only.
18. Gauge: fill glyph, track glyph, cap marker, brackets, percent position.
19. Sparkline and chart glyphs: eighths, braille (only after testing in the font), dots, slashes.
20. State glyph set: ok, wait, bad, idle, running, marker for now.
21. Status bar: key style (bare, bracket, bar), separator, case, position.
22. Accent usage: what the accent colours (title, border, tab, selection, keys) and what it never colours.
23. Secondary slot map: text 7, rule 7 or 8, key slot, label slot, head slot, counter slot.
24. Emphasis: bold, underline, reverse, dim rules.
25. Case convention: sentence, lower, upper for titles, headers, hints, labels.
26. Empty, loading, error states: text, glyph, framing.
27. Loading animation frames.
28. Overlay and palette frame: border set, title style, dimming.
29. Icon set: plain Unicode or Nerd Font, per theme.
30. Narrow-width drop order: which ornament goes first.

Colour rules for all themes: ANSI 1 red bad, 2 green ok, 3 yellow attention, 7 secondary text, 15 primary emphasis, 8 only for non-essential tracks and leaders. Bright slots (9 to 15) may equal the normal ones, so never rely on brightness for focus.

## (d) Four theme proposals

Common to all four: FiraCode Nerd Font, no fills except selection and the mode chip, avoid `▰ ▱ ⚙ ↵ ↻ ⏵ ⏸ ▸`, secondary text 7. Slot 8 is used only for the empty part of a gauge and never for a border. Each theme changes the skeleton, not only the frame.

Identities:

- Lazytf: the dense workbench. Boxes everywhere, one line per row, rounded, magenta titles, blue focus.
- Soft: the paper ledger. Horizontal rules and a left spine delimit panels. Air between panels. Sentence case.
- Cyber: the neon console. Square corners, double lines on focus, slashes everywhere, caps, magenta with cyan.
- HUD: the instrument panel. Corner brackets only, panel numbers, ruler gauges, caps, amber.

### D1. lazytf (dense default)

| Axis | Value |
|---|---|
| Accent slot | 4 blue (focus, keys, selection) |
| Title slot | 5 magenta, bold |
| Border unfocused | `╭ ─ ╮ │ ╰ ╯` in slot 7 |
| Border focused | same glyphs, slot 4 bold |
| Corner ornament | none |
| Panel title | inline on the top border: `╭─ Title ─────── label ─╮`, sentence case |
| Title counter | `Findings 3`, count in slot 7 |
| Right label | on the top border, `label` in slot 4, one `─` gap each side |
| Bottom label | on the bottom border right: `1 of 5` in slot 7 |
| Padding | 1 column inside, gap 0 between boxes |
| Brand | `◆ crow deck`, slot 5 bold |
| Tabs | `1 Home  2 Agents  3 MRs`, active in slot 4 bold underline, inactive slot 7, badge `²` style: yellow 3 bold number after the name |
| Rule under tabs | `─` repeated, slot 7 |
| Mode switch | `all [work] perso m`, active as reverse chip in slot 4, others slot 7, key `m` in slot 4 |
| Selection | reverse video in slot 4 over the whole row (focused pane); unfocused pane: bold slot 4 text with `▏` prefix |
| Density | 1 line per row, no blank lines inside lists, group header is a bold slot 7 row |
| Table header | `  Agent  Account  Doing  Since`, slot 7, no caps, no rule |
| Inner separator | `─` in slot 7 |
| Gauge | fill `━` slot 4, track `─` slot 8, cap marker `│` slot 7, percent right of bar in slot 7: `━━━━━━────  62%` |
| Sparkline | `▁▂▃▅▆▇` slot 4 |
| State glyphs | `●` ok 2, `○` idle 7, `◆` needs you 5, `‖` waiting 3, `●` bad 1 |
| Status bar | `key desc  key desc` with key slot 4, desc slot 7, two spaces apart, no separator glyph |
| Case | sentence case everywhere |
| Empty | `Nothing here.` slot 7 |
| Loading | `syncing…` slot 7 |
| Error | `forge sync failed` slot 1 |
| Overlay frame | round, slot 5 border, title slot 5 |
| Icons | plain Unicode only |
| Narrow drop order | right labels, then bottom labels, then padding |

### D2. soft (paper ledger, delimited)

Delimitation comes from a full-width top rule with the title in it and a left spine on every body row. No right edge, no bottom edge. Panels sit 2 columns apart with no blank row (owner, 2026-11-05).

| Axis | Value |
|---|---|
| Accent slot | 6 cyan (spine when focused, active tab, selection bar) |
| Title slot | 15 bold (focused: slot 6 bold) |
| Border unfocused | top rule `─` in slot 7, spine `▏` in slot 7 |
| Border focused | top rule `━` in slot 6, spine `▌` in slot 6 |
| Corner ornament | top-left `╶` then rule: `╶─ Title ────────── label` |
| Panel title | inside the top rule, one space each side: `╶─ Agents 3 ─────────  → 5 Usage` sentence case |
| Title counter | dim count in slot 7 after the title |
| Right label | at the rule end, slot 6, no brackets |
| Bottom label | none; the `1 of 5` goes at the right end of the last row in slot 7 |
| Padding | 2 columns after the spine, gap 2 columns between panels, no blank row between stacked panels |
| Brand | `Crow deck` slot 6 bold, then a short rule `╶──` slot 7 |
| Tabs | active: name in slot 6 bold underline; inactive slot 7; tabs separated by 3 spaces; badge is a yellow 3 `●` before the name |
| Rule under tabs | one blank row, then nothing |
| Mode switch | `all  work  perso` active in slot 6 bold underline, others slot 7, key `m` in slot 6 |
| Selection | `▍` in slot 6 at the row start, row text bold slot 15, no fill, no reverse. Unfocused pane: `▍` in slot 7 |
| Density | airy: blank row between row groups, group header is slot 7 text followed by `╌` leader to the right edge |
| Table header | sentence case, slot 7, then one `╌` rule row in slot 7 below |
| Inner separator | `╌` in slot 7 |
| Gauge | fill `█` with partial eighth at the end `▏▎▍▌▋▊▉`, slot 6; track `░` slot 8; percent left of bar in slot 7: `62%  ██████▌░░░░` |
| Sparkline | `▁▂▃▄▅▆▇█` slot 6 |
| State glyphs | `●` ok 2, `○` idle 7, `◆` needs you 5, `‖` waiting 3, `✕` bad 1 |
| Status bar | `key desc · key desc` with `·` in slot 7 and key slot 6 bold; preceded by a `╶─` rule row |
| Case | sentence case, titles bold |
| Empty | `All clear.` slot 7 |
| Loading | `syncing  ·  ··  ···` cycling dots, slot 7 |
| Error | `forge sync failed` slot 1, with `✕` |
| Overlay frame | top and left rules only, same as panels, title slot 6 |
| Icons | plain Unicode |
| Narrow drop order | gap columns, then right labels, then padding |

### D3. cyber (neon console)

| Axis | Value |
|---|---|
| Accent slot | 5 magenta (title, focus, selection, brand) |
| Secondary accent | 6 cyan (keys, right labels, headers) |
| Title slot | 5 bold, uppercase |
| Border unfocused | single square `┌ ─ ┐ │ └ ┘` in slot 7 |
| Border focused | double `╔ ═ ╗ ║ ╚ ╝` in slot 5 |
| Corner ornament | bottom-right `◢` in slot 5 on focused, slot 7 otherwise |
| Panel title | `▌// TITLE //` on the top border, caps, `▌//` in slot 5, closing `//` in slot 7 |
| Title counter | `[03]` in slot 6 after the title |
| Right label | `░▒ label` in slot 6, caps |
| Bottom label | `── 1/5 ──` in slot 7 before the corner |
| Padding | 1 column inside, gap 1 column between boxes, 0 rows |
| Brand | `▌CROW_DECK // home` slot 5 bold, then `░▒▓` slot 7 |
| Tabs | `▌1 HOME╱` active: `▌` and `╱` slot 5, name slot 5 bold; inactive: `▌` and `╱` slot 7, name slot 7; badge `!2` in slot 3 bold |
| Rule under tabs | `//////////` slot 5 then `─` slot 7 |
| Mode switch | `[ALL] WORK PERSO` active in reverse chip slot 5, others slot 7, key `m` in slot 6 |
| Selection | `▌` in slot 5 then row text bold slot 15, row prefix `//` replaced by marker; plus underline on the first column. No fill |
| Density | 1 line per row, group header `// NEEDS YOU ////////` slot 6 |
| Table header | caps, prefix `//`, slot 6, rule `─` in slot 7 below |
| Inner separator | `╱╱╱` repeated in slot 7 |
| Gauge | fill `/` slot 5, track `.` slot 7, brackets `[` `]` slot 7, percent inside after bar: `[//////////........] 62%` |
| Sparkline | `▁▂▃▅▆▇` slot 6 |
| State glyphs | `◆` needs you 5, `●` ok 2, `○` idle 7, `‖` waiting 3, `✕` bad 1 |
| Status bar | `[key] DESC // [key] DESC`, key slot 6 in brackets, desc slot 7 caps, `//` slot 7 |
| Case | uppercase titles, headers, hints, labels; data stays as is |
| Empty | `// NO DATA //` slot 7 |
| Loading | `// SYNC / // ///` cycling, slot 6 |
| Error | `// ERR forge sync failed` slot 1 |
| Overlay frame | double border slot 5, title `▌// PALETTE //` |
| Icons | plain Unicode |
| Narrow drop order | `░▒` fades, then `//` closers, then counters |

### D4. hud (instrument panel)

No side lines. Each panel is four corner brackets, a numbered tag and tick marks on the top edge.

| Axis | Value |
|---|---|
| Accent slot | 3 amber (brackets, tag, selection) |
| Secondary accent | 6 cyan (live values, labels) |
| Title slot | 3 bold, uppercase |
| Border unfocused | corners only: `┏━` `━┓` `┗━` `━┛` in slot 7; sides empty |
| Border focused | same corners in slot 3, plus top and bottom edge ticks `╌` in slot 7 between corners |
| Corner ornament | panel index `01` to `05` on the bottom edge: `┗━ 01 ━━━━` in slot 6 |
| Panel title | `◢ TITLE ◣` after the top-left corner, caps; focused slot 3 bold, unfocused slot 15 |
| Title counter | `×03` in slot 6 after the tag |
| Right label | top edge, `╌╌ LABEL` in slot 6, caps |
| Bottom label | `1/5` before the bottom-right corner, slot 7 |
| Padding | 2 columns inside, gap 1 column between panels, 0 rows |
| Brand | `◆ CROW//DECK ░▒▓` slot 3 bold, ornament slot 7 |
| Tabs | `▶ 1 HOME` active: `▶` and name slot 3 bold; inactive: two spaces and name slot 7; badge `●2` in slot 1 |
| Rule under tabs | `━━━━ ╌╌╌ ░░` slot 3, then `╌` in slot 7 |
| Mode switch | `MODE ALL ▏WORK▕ PERSO` active between `▏` `▕` in slot 3 bold, others slot 7, key `m` slot 6 |
| Selection | `▶` in slot 3 at row start, whole row bold slot 3, no fill. Unfocused pane: `›` in slot 7 |
| Density | 1 line per row, blank row between groups, group header `── NEEDS YOU ──` slot 6 |
| Table header | caps, slot 6, `╌` rule below in slot 7 |
| Inner separator | `╌` in slot 7 |
| Gauge | ruler: fill `█` slot 3, track `░` slot 8, tick `│` slot 6 at the policy cap, brackets `▕` `▏`, percent before bar: `62% ▕██████░░░░│░░▏` |
| Sparkline | horizontal eighths `▏▎▍▌▋▊▉█` for single values, `▁▂▃▅▆▇` for series, slot 6 |
| State glyphs | `◆` needs you 3, `●` ok 2, `○` idle 7, `‖` waiting 3, `✕` bad 1 |
| Status bar | `▏key▕ DESC  ▏key▕ DESC`, key slot 3, desc slot 7 caps, two spaces apart |
| Case | uppercase titles, labels, hints, table headers; data as is |
| Empty | `[ NO SIGNAL ]` slot 7 |
| Loading | `[ SYNC ▏▎▍▌▋▊▉█ ]` frames cycling, slot 6 |
| Error | `[ FAULT forge sync failed ]` slot 1 |
| Overlay frame | four corner brackets slot 3, tag `◢ PALETTE ◣` |
| Icons | plain Unicode; Nerd Font glyphs optional later for panel icons |
| Narrow drop order | edge ticks, then panel index, then `░▒▓`, then counters |

Design checks before building:

1. Render each table in a script with the real slots on a translucent terminal. Confirm focus is visible from shape alone (heavy or double or amber).
2. Test `▏▎▍▌▋▊▉` and `░▒▓` in FiraCode Nerd Font at the terminal font size. Test braille separately before adopting it.
3. HUD uses amber 3 as accent. Yellow also means attention. Keep the attention badge as `●n` in slot 1 or a reverse chip so the two do not collide, or let the config override the accent.
4. Keep one skeleton in code (rows, columns, data), and let each theme supply the table above. Selection, gauge, density, status bar, empty states and case are theme fields, not `if frame ==` branches scattered in views.
