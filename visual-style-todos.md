# Visual Style — Future Work & Out-of-Scope Items

This file captures ideas and deferred items from the visual style work on Curious Ape.

## Current Theme State

**Direction**: Earthy natural color scheme. Less dark and stark than the initial Oxide-inspired dark technical theme. Warm olive/forest tones. More grounded, organic, and approachable.

**Primary colors** (as specified by user):
- `background`: `#111c18`
- `foreground`: `#C1C497`

**Key palette** (defined in `assets/css/main.css` `:root`):

- `--bg-base`: #111c18 (main page background)
- `--bg-surface`: #1a2a24
- `--bg-elevated`: #22372f (used for cards, day rows, log entries, and surfaces to increase separation/contrast)
- `--bg-header`: #0f1916
- `--fg-base`: #C1C497 (warm muted olive — primary text)
- `--fg-muted`: #8a9578
- `--fg-inverse`: #111c18
- `--accent`: #7d8f5e (earthy sage green — primary action/highlight color)
- `--accent-hover`: #96a97a
- `--accent-dim`: #374a33
- `--accent-weak`: #2c3c2a
- `--border`: #32433c
- `--border-strong`: #45544d
- `--success`: #6b8a58 (moss green)
- `--warning`: #9c8a5e (warm earth)
- `--danger`: #8b6355 (clay / terracotta)

**Score states** (Bujo / `.day-score`, earthy progression from low to high):
- `.score-0`: #c48d7e (warm clay)
- `.score-1`: #b88f72 (clay)
- `.score-2`: #b5a16e (warm earth)
- `.score-3`: #8fa36d (sage)
- `.score-4`: #8C5FA3 (regal purple — top tier / "king")

**Key visual treatments implemented**:
- Bujo day rows (`.day`) and Deadlines/log entries (`.deadline-item`, `.log-entry`, `.surface`) use `--bg-elevated` + `--border-strong` for stronger contrast against the main background.
- `.day` rows feature a prominent left accent bar using `--accent` (instead of dim).
- Habit grid, buttons, forms, nav, and integrations all use the unified earthy tokens.
- Stronger borders and elevated surfaces throughout for better definition on the dark base.
- Focus/hover states consistently use accent green.
- Habit states (done/not done) and error states use the semantic tokens.
- Score colors chosen for readability on elevated cards while staying in the natural palette.
- Mono font used for scores, data, and grid headers.

**Typography & structure**: System sans-serif base + monospace for data/numbers. No custom fonts. All styling is CSS-variables driven + targeted classes. Changes strictly followed "CSS + class/ID additions only" — no layout, grid, or logic modifications.

**Status**: Theme is cohesive and earthy. Specific contrast improvements were made for Bujo and Deadlines pages. Score-4 uses a distinct regal purple as the highest achievement state.

## Out of Scope for Initial Overhaul (but valuable later)

- **Deadline urgency accents**: Color or icon treatment based on `DaysLeft` (soon = stronger accent).
- **Dark / light or user theme toggle**: Full theme switcher (persisted in session or localStorage via datastar). Requires more CSS vars + a small toggle component.
- Skeleton or loading states for datastar in-flight requests (integrations, sync).

## Polish Ideas (CSS-only or tiny class additions)

- Add a very subtle top green accent line under the header h1.
- Experiment with 0 or 2px border-radius for a harder "machine" Oxide-rack feel vs current 3-4px.
- Stronger use of `--ff-mono` on dates, durations, percentages, IDs.
- Hover lift or underline style on `.day` rows and log entries.
- Consistent icon + text sizing in nav and day rows.

## Notes

- All future visual work should continue the rule: **CSS + class/ID additions only**. No layout structure or business logic changes.
- Palette lives in `assets/css/main.css` `:root`. Update tokens there first.
- When adding new classes, prefer documenting them in `pkg/ui/ui.go` constants where reusable.
- After changes: `mage build`, hard-refresh dev server, and manual check of all views + partial re-renders (habit flip, day sync).

See also the original plan in the session notes for context on the first pass.
