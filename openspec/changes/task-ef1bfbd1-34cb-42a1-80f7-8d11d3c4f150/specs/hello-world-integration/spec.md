## MODIFIED Requirements

### Requirement: The HomeView SHALL display a hello-world message fetched from the backend through generated query artifacts

The change SHALL replace the manual hello fetch in `frontend/src/home/HomeView.vue` with the generated Hey API SDK and Pinia Colada query artifacts for `GET /api/hello`. The first summary card SHALL preserve the current loading, error, and success states while displaying the backend `message` field. The remaining three summary cards and the "Today's plan" / "Move notes" sections SHALL remain unchanged.

#### Scenario: Hello message replaces manual fetch logic

- **WHEN** `frontend/src/home/HomeView.vue` is updated for the generated client integration
- **THEN** it consumes the generated query artifacts for `GET /api/hello`
- **AND** it no longer performs its own manual `onMounted` fetch for that card

#### Scenario: Loading state is shown while the query is pending

- **WHEN** the HomeView renders and the hello query is in-flight
- **THEN** the first summary card shows a loading indicator rather than placeholder text or an error

#### Scenario: Error state is shown when backend access fails

- **WHEN** the HomeView renders and `GET /api/hello` fails
- **THEN** the first summary card shows a graceful error state indicating the backend is unavailable

#### Scenario: Success state shows the backend message

- **WHEN** the hello query resolves successfully
- **THEN** the first summary card displays the returned backend `message` value

#### Scenario: Other HomeView sections are preserved

- **WHEN** the HomeView renders after the generated-client integration
- **THEN** the three remaining summary cards ("Available today", "Under-staffed", "Move days"), the "Today's plan" section, and the "Move notes" section continue to render with their existing static content