## ADDED Requirements

### Requirement: Assign a helper by dragging

The schedule board SHALL allow a user to drag a person from the people rail onto a task card to assign that person. On drop, the system SHALL append the person's ID to the card's `assignedTo` and persist via `updateScheduleCard`.

#### Scenario: Successful assignment

- **WHEN** the user drops a person onto a task card that has capacity
- **THEN** the person is added to the card's assigned helpers optimistically
- **AND** `updateScheduleCard` is called with the updated `assignedTo`
- **AND** the card reconciles to the server-returned `assignedCount` and `staffingStatus`

#### Scenario: Assignment completes staffing

- **WHEN** an assignment causes the server to return `staffingStatus: fullyStaffed`
- **THEN** the card reflects fully-staffed state and the fully-staffed reward animation plays (subject to motion preference)

#### Scenario: Assignment rejected by server

- **WHEN** `updateScheduleCard` returns an error for the assignment
- **THEN** the optimistic change is rolled back to the pre-drop state
- **AND** a non-blocking error is surfaced to the user
- **AND** no fully-staffed reward is played

#### Scenario: Duplicate assignment

- **WHEN** the user drops a person already assigned to that card
- **THEN** no duplicate is added and no redundant mutation is sent

### Requirement: Reschedule a card by dragging

The schedule board SHALL allow a user to drag a task card from one day column to another to reschedule it. On drop, the system SHALL set the card's `scheduledDate` to the target day and persist via `updateScheduleCard`.

#### Scenario: Successful reschedule

- **WHEN** the user drops a card onto a different day column
- **THEN** the card moves to the target column optimistically with a FLIP animation
- **AND** `updateScheduleCard` is called with the new `scheduledDate`
- **AND** the board reconciles with the server response

#### Scenario: Reschedule rejected by server

- **WHEN** `updateScheduleCard` returns an error for the reschedule
- **THEN** the card animates back to its original column
- **AND** a non-blocking error is surfaced to the user

### Requirement: Drag affordances

The board SHALL provide clear visual affordances during drag, gated by motion preference for the non-essential parts.

#### Scenario: Picking up a draggable

- **WHEN** the user begins dragging a person or card and motion is enabled
- **THEN** the dragged element lifts (elevation/scale/tilt) and valid drop targets are visually highlighted

#### Scenario: Reduced motion drag

- **WHEN** the user drags with motion disabled
- **THEN** drop targets are still highlighted and the operation is fully usable without the lift/tilt flourish

### Requirement: Read-only schedule is not draggable

The home page's read-only `DailySchedule` SHALL NOT expose drag interactions.

#### Scenario: Home schedule roll-up

- **WHEN** the schedule is rendered in read-only mode on the home page
- **THEN** no drag handles or drop targets are active and the data is display-only
