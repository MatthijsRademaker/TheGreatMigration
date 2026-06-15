# greeting-banner Specification

## Purpose
TBD - created by archiving change task-200fb581-9ace-4ca8-ac95-39ce94d09626. Update Purpose after archive.
## Requirements
### Requirement: REQ-GREET-01 Banner Visibility

**Description:** A greeting banner MUST be rendered at the top of the home page as the first child of the `<section>` element in `src/home/HomeView.vue`, above the summary cards grid.

**Rationale:** The home page at the `/` route is the app's entry point. Placing the banner at the top of the view ensures immediate visibility for anyone validating the swarm pipeline.

**Keywords:** SHALL

#### Scenario: Banner renders on home page load
- **Given** the app is running
- **When** a user navigates to the `/` route
- **Then** a banner with text "👋 Hello world — swarm pipeline test" is visible as the first child of the `<section>` element in HomeView.vue, above the summary cards grid

### Requirement: REQ-GREET-02 Consistent Styling

**Description:** The banner SHALL use Tailwind CSS v4 utility classes consistent with the project's design system.

**Rationale:** The `bg-muted/40 rounded-lg border p-4` pattern is established in `tasks/TasksView.vue` and `shared/ui/input/Input.vue`. The `border-l-4 border-primary` provides a subtle accent consistent with the green/nature-inspired theme. This avoids introducing new shadcn-vue components.

**Keywords:** SHALL

#### Scenario: Banner uses existing themed tokens
- **Given** the home page is rendered
- **Then** the greeting banner has the utility classes `bg-muted/40 rounded-lg border p-4 border-l-4 border-primary text-sm`
- **And** no new CSS imports, design tokens, or shadcn-vue component additions are required

### Requirement: REQ-GREET-03 No Side Effects

**Description:** The banner MUST NOT introduce any new script imports, composables, data fetching, router configuration, or layout changes.

**Rationale:** The change is purely presentational. Any side effects would violate the non-goals and add unnecessary scope to a pipeline validation task.

**Keywords:** MUST NOT

#### Scenario: No script or data layer changes
- **Given** the current HomeView.vue `<script setup>` block
- **When** the greeting banner is added
- **Then** no new imports, composables, store access, or props are added to the `<script>` block
- **And** no router configuration, app shell, or layout components are modified

### Requirement: REQ-GREET-04 Build Verification

**Description:** The change SHALL pass the project's verification gates after implementation.

**Rationale:** Acceptance criteria require pre-commit hooks and production build to succeed, confirming the change does not break existing functionality.

**Keywords:** SHALL

#### Scenario: Verification passes
- **Given** the greeting banner has been added to HomeView.vue
- **When** `scripts/precommit-run` is executed
- **Then** all pre-commit hooks pass without errors
- **And** when `npm run build` is executed, the app builds successfully (vue-tsc -b && vite build)
