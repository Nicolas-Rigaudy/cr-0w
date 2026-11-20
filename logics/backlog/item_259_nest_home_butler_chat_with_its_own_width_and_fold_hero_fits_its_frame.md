## item_259_nest_home_butler_chat_with_its_own_width_and_fold_hero_fits_its_frame - Nest Home butler chat with its own width and fold, hero fits its frame
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 100%
> Confidence: 95%
> Progress: 100%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-19 17:32:08

# AI Context
- Summary: Nest Home in two parts: the hero takes about 2/3 of the screen, the right third is the butler chat with the old left column folded into a header above it. Home keeps its own chat width and fold state, and the hero fits its frame without clipping or a janky resize.
- Keywords: nest, home, butler, chat, rail, fold, width, hero, holo table, resize
- Use when: changing the butler rail on Home or the hero sizing.
- Skip when: the rail on other pages (item_256 done there).

# Problem
- Owner 2026-11-19: folding the chat on Home widens the hero and the holo table clips at the edges.
- The hero resize animates badly while the rail width changes.
- Home should be the most comfortable place to talk to the butler; today it is the same rail as every page.
- One fold key (`nest.rail`) and one width key (`nest.railw`) are shared by every page.
- Owner 2026-11-19: the Home left column is often quite empty. Merge it into a header above the chat and give the hero about 2/3 of the screen.

# Scope
- In:
  - Home layout: hero about 2/3 wide; right third = left-column facts as a compact header, chat below it.
  - Mockups first (laptop and Fold inner, both themes): owner picks the header content and density before the build.
  - Fold inner (about 707x691): redesign Home so it uses the panel space (owner 2026-11-19: the inner display wastes space). The mockups decide whether it takes the laptop split or its own layout.
  - Home keeps its own fold and width keys, separate from the other pages.
  - Home opens the chat wider by default (about 40% of the window, clamped like item_256).
  - The hero fits its frame at every width: the holo table never clips. No width tween on the hero; it refits once when the fold or drag ends.
  - Shots at rail open, folded and dragged, both themes, laptop sizes.
- Out:
  - Other pages' layout (they keep the item_256 rail).
  - Phone cover Home (keeps its own layout; item_228). Owner likes the cover hero size as is.

# Acceptance criteria
- AC0: Owner approved the Home mockup (hero 2/3, header plus chat) before the build.
- AC1: Folding or resizing on Home leaves the other pages' rail as it was, and the reverse.
- AC2: First visit to Home shows the wider chat.
- AC3: No clipping of the holo table in shots at 1280 and 1920, rail open and folded, both themes.
- AC5: Fold inner Home uses its space per the approved mockup; shots in both themes, hero not clipped.
- AC4: Fold and drag on Home show no stretched or jumping frames (recorded check or frame shots).

# AC Traceability
- `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits` AC2 -> This backlog slice. Proof: Home chat works without clipping.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Builds on: `item_256_resizable_butler_chat_rail_in_nest`

# Priority
- Priority: Medium
- Rationale: Owner-visible glitch on the most used page; small build.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Notes
- 2026-11-19, table fit done ahead of the layout: the table scale width is `min(w, 1.05 * h)` (`FIT` in `src/hero/scene.ts`). Wide, short heroes (Fold inner 683x420, laptop with the chat folded 1010x816, 1920 folded) keep the table and planets in the frame. Cover and laptop with the chat open are unchanged (width under 1.05x height). Shots: `scripts/hero-fit.mjs`, `local/state/shots/hero-fit/compare.png`.
- 2026-11-20, AC0 met: owner approved the mockup. Design notes and boards: vault `Research/nest-design/item_259/design-notes.md`.
  - Laptop Home: hero left 2/3, right third (480px default, own fold and width keys) = mood line, HUD card, chat transcript, input. Build reference `sketch.html?v=H3`.
  - HUD card (H3 tiles): header "Needs you" + count + Inbox link; each Needs-you item is its own card with an icon tile, the top ask highlighted with the only filled button, others ghost buttons; below, Next (big time, title, countdown bar) and Running (live count, slots pips) as two tinted cards.
  - Fold inner: F1 laptop split (hero left, about 300px header + chat right), header = H3 compacted (top ask card + one summary line). Reference `sketch.html?v=F1`.
  - Rejected: tiles strip (A), tabs (C), wide-hero fold (F2), lists + Ask bar fold (F3), styles glass/console/minimal.
- 2026-11-20, built.
  - Home rail: own keys `nest.rail.home` (opens by default) and `nest.railw.home` (480px default); the rail is keyed apart in App.tsx so state never leaks. Header = mood line + `HomeHud` (H3 tiles, Next and Running cards). Rail modifier class is `home-rail` (`.home` is an old grid rule).
  - Hero: refit debounced 160ms and held while `data-rail-drag` is set; the canvases sit in a `.hero-stage` kept at the last fitted size and centred, so nothing stretches mid-drag.
  - Fold inner: `HomeFold` split, hero left, 300px side (mood line, `HomeHudCompact`, chat). No header badge on inner Home since the chat is in the side. The side keeps the account picker (the board left it out).
  - Proof (vault `Research/nest-design/item_259/build/`, script `build-shots.mjs`, both themes):
    - AC1: after a drag Home stores 680, other pages keep 320 (Inbox rail 320px).
    - AC2: 1280 open stage 800x744 (rail 480).
    - AC3: no planet, rim or card outside the hero box at 1280/1920 open and folded.
    - AC4: during an 8-step drag the stage stays 800x744 / 1440x1024, then refits once (600x744 / 1240x1024).
    - AC5: `fold-inner-*.png` against `shots/F1-*.png`; stage 407x587, not clipped.
  - Tests: Home width key and 480 default in `railWidth.test.ts`; 247 tests pass, typecheck and build clean.
