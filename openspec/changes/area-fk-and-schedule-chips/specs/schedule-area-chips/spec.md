## ADDED Requirements

### Requirement: Area color SHALL be derived deterministically from the area id

The frontend SHALL provide a pure `areaColor(id)` utility that maps an area id to a color from a fixed palette by hashing the id. The same id SHALL always yield the same color; the color SHALL NOT depend on the area name, so renaming an area preserves its color. No color value SHALL be persisted in the backend or read from the API.

#### Scenario: Same area id yields the same color

- **WHEN** `areaColor` is called twice with the same area id
- **THEN** it returns the same palette color both times

#### Scenario: Color is independent of the area name

- **WHEN** an area is renamed but keeps its id
- **THEN** `areaColor(id)` returns the same color as before the rename

### Requirement: Schedule cards and task rows SHALL render a color-coded area chip

`TaskBoardCard.vue` SHALL render the card's area name as a chip tinted with `areaColor(area.id)`, displayed distinctly from the priority accent (the chip SHALL NOT reuse the priority left-border). The task-board row SHALL render the same area chip. The chip SHALL make two cards with identical titles but different areas visually distinguishable.

#### Scenario: A schedule card shows its area name and color

- **WHEN** a schedule card with an area is rendered
- **THEN** the card displays the area name
- **AND** the chip is tinted with the color returned by `areaColor(area.id)`

#### Scenario: Two cards with the same title but different areas are distinguishable

- **WHEN** two schedule cards both titled "Paint" reference different areas
- **THEN** each card shows its own area name
- **AND** the two chips render in different colors

#### Scenario: The area chip does not collide with the priority accent

- **WHEN** a card renders both its priority accent and its area chip
- **THEN** the area color is shown on the chip, not on the priority left-border
- **AND** the priority accent remains driven by priority, not by area
