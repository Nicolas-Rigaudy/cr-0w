# Visual concepts (worlds, scenes, themed boards)

Use before any render run for a visual world or themed view of a system. Renders made before the purpose is agreed get rejected whole.

## Steps
1. **Purpose.** Ask the owner: when is it looked at, what must it answer at a glance, watch or act, why a picture instead of a list. Record the answers in the request doc.
2. **Content inventory.** List every real signal it shows (runs, states, decisions waiting, quota, finished work) and what each click opens. Agree what stays out.
3. **Flat sketches.** Draw 3 to 4 worlds as annotated 2D SVGs of the same moment, each element labelled with the data it shows. The owner picks meanings and worlds, not looks.
4. **Image prompts.** One prompt per picked world and theme, run by the owner in an image model:
   - frame it as "in-game screenshot of an isometric <genre> game";
   - number the stations in flow order;
   - reserve one colour per state (working, needs you, failed, done, queued);
   - no UI words (dashboard, panel, chart, widget, screen, HUD, app);
   - no franchise or character names (image models refuse them); describe original designs instead.
5. **Compare.** Put the images in a side-by-side compare page and ask the pick with a self-contained question.
6. **Only then** plan the build: data map, rendering method (at most 3 options), click targets.

## Keep
- Prompts versioned next to the research (`PROMPT.md`, rejected versions kept with a suffix).
- Worlds not picked noted in the request as later candidates.
