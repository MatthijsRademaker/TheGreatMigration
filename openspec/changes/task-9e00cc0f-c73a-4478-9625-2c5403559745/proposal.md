## Why

The LikeC4 architecture diagram at `.devagent/architecture/workspace.c4` models the Dashboard feature area as a single opaque block with generic descriptions, but the frontend now has four distinct KPI summary cards — each backed by a different data source. The architecture model should reflect this decomposition so the diagram accurately communicates the system's component structure and data flow to readers.

## What Changes

- Decompose the `Dashboard` feature area in the C4 model into four explicit KPI card sub-elements within the dashboard, each with its own description, data dependencies, and visual identity.
- Add a dedicated scoped C4 component view of the Dashboard that renders the KPI cards and their data-flow relationships.
- Update the existing `feature-map` and `system-context` views if needed to surface the refined dashboard structure.
- Keep the change strictly to the LikeC4 architecture model — no frontend Vue code, backend handlers, tests, or application behavior changes.

## Capabilities

### New Capabilities

- `kpi-summary-cards`: Architectural model for the four KPI summary cards on the home dashboard, covering their element definitions, data-source relationships, and diagram views in LikeC4 DSL.

### Modified Capabilities

- (none — this is a documentation/architecture-only change, no application requirements are modified)

## Impact

- `.devagent/architecture/workspace.c4` — the single LikeC4 source file, extended with new element definitions, relationships, and a scoped dashboard view.
- `designs/components.png` and `designs/home-page.png` — no change (these are the design reference the architecture should align with).
- All frontend and backend source files — no change.
