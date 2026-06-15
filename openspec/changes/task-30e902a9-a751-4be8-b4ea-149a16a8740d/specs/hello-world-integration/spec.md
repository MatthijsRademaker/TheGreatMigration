## ADDED Requirements

### Requirement: The HomeView SHALL display a hello-world message fetched from the backend at runtime
The change SHALL replace the first static summary card in `frontend/src/home/HomeView.vue` (currently "High priority / 4") with a live fetch to `GET /api/hello`. The card SHALL display the returned `message` field from the response. The remaining three summary cards and the "Today's plan" / "Move notes" sections SHALL be preserved unchanged.

#### Scenario: Hello message replaces static card content
- **WHEN** the frontend loads with a reachable backend
- **THEN** the first summary card displays the message from `GET /api/hello` instead of the static "High priority / 4" content

#### Scenario: Loading state is shown while fetching
- **WHEN** the HomeView mounts and the hello fetch is in-flight
- **THEN** the first summary card shows a loading indicator rather than placeholder text or an error

#### Scenario: Error state is shown when backend is unreachable
- **WHEN** the HomeView mounts and `GET /api/hello` fails (backend down, network error, or CORS rejection)
- **THEN** the first summary card shows a graceful error state indicating the backend is unavailable

#### Scenario: Other HomeView sections are preserved
- **WHEN** the HomeView renders after the hello-world integration
- **THEN** the three remaining summary cards ("Available today", "Under-staffed", "Move days"), the "Today's plan" section, and the "Move notes" section continue to render with their existing static content