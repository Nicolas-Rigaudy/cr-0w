# README that sells (showcase a project)

Method used on the hub's own README. Reuse it per repo, one repo session each (`crow handoff`).

1. **Read first.** The repo's README, docs, tests and what it really does. Never claim a feature that is not built.
2. **Look at 3 or 4 popular repos of the same kind** (raw README from GitHub). Note: header, first screen, install placement, depth.
3. **Visuals before words.** Produce real screenshots from fixtures, never real data. TUI: `runbooks/tui-render.md`. GUI or shell: the repo's own screenshot script. Compress (`magick in.png -colors 128 PNG8:out.png`).
4. **Header options side by side.** Render 3 to 5 variants on a dark and a light background with `bin/crow-compare`; the owner picks. Use block letters or paths, not fonts, so it renders the same everywhere.
5. **Structure** (about 120 to 160 lines): header and tagline, hero screenshot, one tight paragraph, a day-in-the-life list, a pain-to-benefit table, a short tour with one image per feature, honest setup steps, collapsible detail, license.
6. **Move depth to docs/** and link it. Keep the README skimmable.
7. **Audit pass.** Cut repeats of the same pitch, jargon (script names), and anything that overpromises (quick start that needs setup).
8. **Preview as GitHub renders it** (marked + github-markdown-css + mermaid in a local HTML) before committing.
9. **Prose rules:** short and declarative, no em-dashes, no "X, so Y", no "not X but Y", no AI mention. English for personal projects.
10. Commit as `chg[MD/SVG/PNG]: ...`; push only on the owner's OK (public repos).

Portfolio level: a profile README (`<user>/<user>` repo) and pinned repos link the projects together. See the owner's thread on the personal site.
