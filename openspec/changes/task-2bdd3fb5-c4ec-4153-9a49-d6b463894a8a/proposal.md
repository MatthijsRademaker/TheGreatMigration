## Why

The frontend now has a visual homepage target (`designs/home-page.png`), a component inventory (`designs/components.png`), shared UI primitives, sidebar/navigation, and backend dashboard contracts for task backlog, people availability, and daily schedule. However, `frontend/src/home/HomeView.vue` is still a placeholder composition — it renders a Hello card, two hardcoded summary cards, a Move days card, a static Today's plan, and Move notes rather than the full dashboard composition from the design artifacts.

A text-only LLM cannot reliably infer the intended homepage layout from a PNG or from scattered route/spec files. This change produces a monospaced ASCII blueprint that translates the homepage composition into text while staying aligned with existing component names, dashboard sections, and data contracts.

## What Changes

- Create a non-canonical ASCII homepage mockup at `openspec/changes/task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a/home-page-ascii.md`.
- The artifact contains: (1) a legend mapping each dashboard region to current repo components, OpenSpec specs, and data contracts; (2) a desktop-width monospaced wireframe showing the 4-row layout (top toolbar, 4 KPI cards, Tasks Backlog + People Availability, Daily Schedule + Move Notes) using existing repo vocabulary; (3) a collapsed-sidebar variant showing content-area reflow; (4) a gap analysis comparing the current `HomeView.vue` placeholder against the target composition.
- No canonical OpenSpec specs, Vue component files, or backend files are modified.

## Impact

- **Affected specs**: None — this is a non-canonical change-local artifact.
- **Affected code**: None — no Vue, Go, CSS, or configuration files are modified.
- **New artifact**: `openspec/changes/task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a/home-page-ascii.md` — a companion reference for future homepage composition work.
- **Downstream**: A future task implementing the `HomeView` composition can read this ASCII blueprint without opening PNG design files.