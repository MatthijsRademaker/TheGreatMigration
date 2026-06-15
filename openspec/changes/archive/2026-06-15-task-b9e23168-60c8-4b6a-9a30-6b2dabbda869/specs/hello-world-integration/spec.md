# hello-world-integration Specification (Delta)

## Purpose
Amend the canonical `openspec/specs/hello-world-integration/spec.md` to accommodate the planning-window change, which makes the "Move days" summary card value dynamic.

## MODIFIED Requirements

### Requirement: The HomeView SHALL display a hello-world message fetched from the backend at runtime

The change SHALL replace the first static summary card in `frontend/src/home/HomeView.vue` (currently "High priority / 4") with a live fetch to `GET /api/hello`. The card SHALL display the returned `message` field from the response. The remaining summary cards ("Available today" and "Under-staffed") and the "Today's plan" / "Move notes" sections SHALL be preserved unchanged. The "Move days" summary card SHALL derive its value from the planning window rather than a static constant.

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
- **THEN** the two static summary cards ("Available today" with value `"6"`, "Under-staffed" with value `"3"`), the "Today's plan" section, and the "Move notes" section continue to render with their existing static content

#### Scenario: Move days card derives from the planning window
- **WHEN** the HomeView renders after the planning-window integration
- **THEN** the "Move days" summary card displays a value derived from the shared planning window's `planWindowDayCount` rather than a static `"5"`
