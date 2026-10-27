## req_004_portfolio_readmes_profile_readme_and_personal_page - Portfolio READMEs, profile README and personal page
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Finish the public portfolio: showcase READMEs for lazytf, lorenzini-track and obsidian-stacks, then the GitHub profile README and a personal page linking every public project.
- Keywords: portfolio, readme, profile, personal page, showcase
- Use when: working on any README, the profile repo or the personal site.
- Skip when: the project is cr-0w, archeotech-dotfiles or archeotech-shell (done 2026-10-26) or a school project.

# Needs
- One repo session per project, method in `runbooks/readme-showcase.md` (genre first, real fixture screenshots, no cloned skeleton).
- Profile README (`Nicolas-Rigaudy/Nicolas-Rigaudy`) with pinned repos, then a personal page hosting a resume and links to the public projects (domain to decide).
- Final side-by-side of all READMEs and logos (`bin/crow-compare`): one design family or distinct styles, owner decides.
- Out: Aquarium-Monitoring-System, PixabayImageSearcher (school), SkaterXLProjects (no marketing, not pinned).

# Context
- Started 2026-10-22; cut by the weekly limit before its wrap. Done: cr-0w, archeotech-dotfiles, archeotech-shell on main.
- Left on the shell/dotfiles side: showcase video (owner records it from `docs/showcase-video.md`, then the placeholder screenshot is replaced).
- Open question: obsidian-stacks README has 117 lines and no visuals; lazytf has a newcomer README from an AFK run, lorenzini-track has one from the demo build.
- Order: profile README first (decided 2026-10-27), then lazytf, lorenzini-track, obsidian-stacks, personal page, compare.

# Acceptance criteria
- AC1: lazytf, lorenzini-track and obsidian-stacks each have a showcase README with real screenshots, on main.
- AC2: The profile README exists and pins the public projects.
- AC3: A personal page is live or a decision with options is recorded.
- AC4: The final compare is shown and the owner's pick is applied.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics_manager/flow.py`
- `logics_manager/assist.py`
- `tests/python/test_logics_manager_cli.py`

# Backlog
- none
