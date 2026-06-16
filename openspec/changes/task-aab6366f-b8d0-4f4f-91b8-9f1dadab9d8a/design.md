## Context

The home dashboard (`frontend/src/home/HomeView.vue`) currently renders a top summary row with a hello-world card, two hardcoded summary cards ("Available today: 6", "Under-staffed: 3"), and a planning-window card. The backend already exposes contract-backed endpoints for people availability (`GET /api/dashboard/people-availability` returning `availableToday`/`totalPeople`) and task backlog summaries (`GET /api/tasks/backlog` returning `highPriorityTasks`/`unassignedTasks`). However, the committed `frontend/openapi-snapshot.json` only contains hello, planning-window, and people-availability — the task backlog endpoint is absent from the generated frontend client. A fourth design card ("Rooms completed") has no backend contract, migration table, or API endpoint anywhere in the repository.

The project's design system v2 (`docs/design-system-v2.md`) mandates reuse of semantic tokens and existing primitives (Card, Badge, Button). The repository verification flow (`scripts/precommit-run`) runs through Docker and must pass without a live backend.

## Goals

- Replace the current home-page top summary row with a feature-local KPI card grid component that renders four KPI cards structurally aligned to section 3 of `designs/components.png`.
- Drive the **People available today** card from the existing `getDashboardPeopleAvailabilityQuery`, displaying `availableToday` of `totalPeople`.
- Drive the **High priority tasks** and **Unassigned jobs** cards from a newly generated `getTasksBacklogQuery` exposed through the committed frontend client.
- Isolate the **Rooms completed** card as an explicit placeholder with a documented seam, not deriving values from unrelated backend data.
- Update the SSR route-render smoke test to mock the new queries and assert KPI card labels.
- Preserve the existing lower `Today's plan` and `Move notes` sections.
- Pass the repository verification flow (`scripts/precommit-run`).

## Non-Goals

- Redesigning the full dashboard, sidebar, top bar, people matrix, schedule board, or notes panels.
- Adding new global theme tokens or new shared primitive directories beyond what existing Card/Badge/Button styling already supports.
- Implementing real filtering, task editing, room management, or schedule mutations.
- Inventing a backend room-completion metric or persistence model without a defined contract.
- Adding a new backend endpoint or modifying backend Go code.
- Rendering the planning-window card inside the KPI grid (the design shows it belongs to a different section).

## Decisions

### 1. Codegen sequencing: regenerate snapshot and client first

The task backlog endpoint must be present in the committed generated client before KPI cards can consume `highPriorityTasks` and `unassignedTasks`. The implementation must refresh `frontend/openapi-snapshot.json` from the running backend (or manually append the tasks backlog schema fragment) and regenerate `frontend/src/client/` as a prerequisite step.

### 2. Component location and naming

Create a feature-local component at `frontend/src/home/components/KpiCards.vue`. This keeps the KPI card logic colocated with the home feature rather than in a shared directory, following the dossier's guidance to stay inside the home feature.

### 3. Hello-world card removal

The hello-world card is removed from the top summary row to make space for the four KPI cards shown in the design. The dossier goals state "replace the current home-page summary row." If developer backend-health visibility is needed, that is a separate follow-up concern.

### 4. Fourth card placeholder treatment

The Rooms completed card is rendered as a deliberate placeholder using the same Card primitive layout as the other three cards. It displays `—` for the value and `Rooms completed` as the label, with `data-testid="kpi-placeholder-rooms-completed"` and a code comment documenting it as a placeholder for a future room-progress contract. It must not derive counts from task backlog, schedule, or any other unrelated data.

### 5. People available display format

The card displays `availableToday` of `totalPeople` (e.g., "6 of 8 available") using both fields from the `Summary` response, matching the availability contract's canonical homepage summary definition.

### 6. Styling approach

All cards use existing Card primitives (Card, CardHeader, CardContent, CardTitle, CardDescription), semantic accent classes from the design system v2 token surface, and lucide-vue icons. The implementation does not introduce new global tokens or shared component directories. Decorative details from `designs/components.png` (leaf watermark, specific accent borders) that cannot be confirmed from `docs/design-system-v2.md` are approximated using available semantic utility classes.

### 7. SSR test expansion

The `renderRoute` helper in `frontend/tests/app-routes-render.test.ts` must add conditional mock responses for `/api/dashboard/people-availability` and `/api/tasks/backlog`. The home-route test assertion must verify the four KPI card labels appear in the rendered output.

## Risks

- **Codegen drift (high)**: If the snapshot is refreshed manually, the response shape must exactly match the running backend's `TaskBacklogBody` contract. A mismatch would cause type errors in the generated Pinia Colada queries. Mitigation: verify the backend response shape from `backend/tasks.go` before constructing the snapshot schema.
- **SSR test failure (high)**: Adding two new queries to HomeView without updating the SSR mock will cause the home route test to fail. Mitigation: update the mock before or alongside the HomeView changes.
- **Design fidelity (medium)**: Decorative details from `designs/components.png` (leaf watermark, accent border treatments) are not fully specified in `docs/design-system-v2.md`. Mitigation: stick to existing Card primitives and semantic tokens; accept approximation for decorative elements.
- **Fourth-card scope creep (low)**: Without explicit isolation, the placeholder card could be mistaken for a data-driven card. Mitigation: use `data-testid`, code comments, and a clearly non-numeric placeholder value.
- **Snapshot regeneration surprise diffs (medium)**: Regenerating from a running backend may pull in additional schema changes beyond `/api/tasks/backlog` (e.g., daily-schedule). Mitigation: review the diff for unexpected changes to existing contracts before committing.

## Traceability

- **Task**: `aab6366f-b8d0-4f4f-91b8-9f1dadab9d8a`
- **Dossier**: `2026-06-16T04:54:22.288Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Round outputs**: `round:1:agent:swarm-architect` (session `019ecec7-ec10-709d-b5b1-dad7ceb54a89`), `round:1:agent:swarm-lead-dev` (session `019ecec8-e960-7bc5-8aac-2fa27656e8c7`), `round:1:agent:swarm-reviewer` (session `019ececa-3e35-76bb-b207-b97040bce7b6`)
- **Artifact base**: `initial` snapshot for `task-aab6366f-b8d0-4f4f-91b8-9f1dadab9d8a`

## Conflict Resolution

- **Hello-world card**: The architect recommended keeping the hello-world card as a developer visibility aid unless explicitly removed, while the lead-dev specified removal from the top row to match the 4-card design. **Resolved**: Remove from the top KPI row. The dossier goals state "replace the current home-page summary row," and the design requires 4 KPI cards. A separate developer-health indicator is deferred to follow-up.
- **Fourth card rendering**: The architect suggested a static placeholder with dashed border; the lead-dev specified rendering it with `data-testid` and a documented seam; the reviewer asked for explicit placeholder isolation. **Resolved**: Render the fourth card with the same Card layout as the others, displaying `—` for the value, with `data-testid="kpi-placeholder-rooms-completed"` and a code comment.
- **Component name**: Architect used `KpiCardGrid`, lead-dev used `KpiCards`. **Resolved**: Use `KpiCards.vue` for brevity and to match lead-dev's concrete implementation path.