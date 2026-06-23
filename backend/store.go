package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/user/the-great-migration/backend/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/user/the-great-migration/backend/db"
)

// PgStore implements api.Store backed by a pgx connection pool and sqlc queries.
type PgStore struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// NewPgStore creates a new PgStore from a pgx connection pool.
func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{
		queries: db.New(pool),
		pool:    pool,
	}
}

// GetPlanningWindow returns the singleton planning window row from the database.
func (s *PgStore) GetPlanningWindow(ctx context.Context) (*api.PlanningWindowBody, error) {
	row, err := s.queries.GetPlanningWindow(ctx)
	if err != nil {
		return nil, err
	}

	startDate := pgDateToTime(row.StartDate)
	endDate := pgDateToTime(row.EndDate)
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	return &api.PlanningWindowBody{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		Days:      days,
	}, nil
}

// UpdatePlanningWindow persists the singleton planning-window row via UPSERT and returns the updated body.
func (s *PgStore) UpdatePlanningWindow(ctx context.Context, startDate, endDate time.Time) (*api.PlanningWindowBody, error) {
	row, err := s.queries.UpsertPlanningWindow(ctx, db.UpsertPlanningWindowParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	start := pgDateToTime(row.StartDate)
	end := pgDateToTime(row.EndDate)
	days := int(end.Sub(start).Hours()/24) + 1

	return &api.PlanningWindowBody{
		StartDate: start.Format("2006-01-02"),
		EndDate:   end.Format("2006-01-02"),
		Days:      days,
	}, nil
}

// GetPeopleAvailability returns people and availability data for the given date range,
// with optional pagination. When limit <= 0, all people are returned.
func (s *PgStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int, offset int, limit int) (*api.DashboardBody, error) {
	// Get total count regardless of pagination.
	totalCount, err := s.queries.CountPeople(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch people — paginated or all.
	var peopleRows []db.Person
	if limit > 0 {
		peopleRows, err = s.queries.GetPeoplePaginated(ctx, db.GetPeoplePaginatedParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	} else {
		peopleRows, err = s.queries.GetAllPeople(ctx)
	}
	if err != nil {
		return nil, err
	}

	// Fetch availability for the date range.
	startPG := pgtype.Date{Time: startDate, Valid: true}
	availRows, err := s.queries.GetAvailabilityByDateRange(ctx, db.GetAvailabilityByDateRangeParams{
		StartDate: startPG,
		Days:      int32(days),
	})
	if err != nil {
		return nil, err
	}

	// Index availability by person_id -> date -> status.
	availMap := make(map[string]map[string]string)
	for _, a := range availRows {
		if !a.Date.Valid {
			continue
		}
		dateStr := a.Date.Time.Format("2006-01-02")
		if availMap[a.PersonID] == nil {
			availMap[a.PersonID] = make(map[string]string)
		}
		availMap[a.PersonID][dateStr] = a.Status
	}

	// Build end date.
	endDate := startDate.AddDate(0, 0, days-1)

	// Build per-person availability.
	people := make([]api.Person, 0, len(peopleRows))
	for _, pr := range peopleRows {
		avail := make([]api.AvailabilityEntry, days)
		for d := 0; d < days; d++ {
			date := startDate.AddDate(0, 0, d)
			dateStr := date.Format("2006-01-02")
			status := "off" // default if no availability row
			if m, ok := availMap[pr.ID]; ok {
				if s, ok := m[dateStr]; ok {
					status = s
				}
			}
			avail[d] = api.AvailabilityEntry{
				Date:   dateStr,
				Status: status,
			}
		}
		people = append(people, api.Person{
			ID:           pr.ID,
			Name:         pr.Name,
			Initials:     pr.Initials,
			Availability: avail,
		})
	}

	// Compute summary.
	// AvailableToday is computed globally from availMap (all people, not just paginated subset).
	selectedDate := startDate.Format("2006-01-02")
	availableToday := 0
	for _, dayMap := range availMap {
		if dayMap[selectedDate] == "available" {
			availableToday++
		}
	}

	return &api.DashboardBody{
		Range: api.Range{
			StartDate:    startDate.Format("2006-01-02"),
			EndDate:      endDate.Format("2006-01-02"),
			Days:         days,
			SelectedDate: selectedDate,
		},
		Summary: api.Summary{
			AvailableToday: availableToday,
			TotalPeople:    int(totalCount),
		},
		// Pagination is computed by the handler from offset/limit/totalCount.
		People:   people,
		Statuses: api.StatusLegendData,
	}, nil
}

// GetTaskBacklog returns the full task backlog response from the database.
func (s *PgStore) GetTaskBacklog(ctx context.Context) (*api.TaskBacklogBody, error) {
	taskRows, err := s.queries.GetTaskBacklog(ctx)
	if err != nil {
		return nil, err
	}

	assignmentRows, err := s.queries.GetTaskBacklogAssignments(ctx)
	if err != nil {
		return nil, err
	}

	assignMap := make(map[string][]string)
	for _, a := range assignmentRows {
		assignMap[a.TaskID] = append(assignMap[a.TaskID], a.PersonID)
	}

	tasks := make([]api.TaskRow, len(taskRows))
	for i, tr := range taskRows {
		assigned := assignMap[tr.ID]
		if assigned == nil {
			assigned = []string{}
		}
		tasks[i] = api.TaskRow{
			ID:           tr.ID,
			Title:        tr.Title,
			Priority:     tr.Priority,
			PeopleNeeded: int(tr.PeopleNeeded),
			Area:         api.Area{ID: tr.AreaID, Name: tr.AreaName},
			Status:       tr.Status,
			AssignedTo:   assigned,
		}
	}

	total := len(tasks)
	highPriority := 0
	unassigned := 0
	understaffed := 0
	for _, t := range tasks {
		if t.Priority == "high" {
			highPriority++
		}
		if len(t.AssignedTo) == 0 {
			unassigned++
		} else if len(t.AssignedTo) < t.PeopleNeeded {
			understaffed++
		}
	}

	return &api.TaskBacklogBody{
		Summary: api.TaskSummary{
			TotalTasks:        total,
			HighPriorityTasks: highPriority,
			UnassignedTasks:   unassigned,
			UnderstaffedTasks: understaffed,
		},
		Tasks:      tasks,
		Priorities: api.PriorityLegendData,
		Statuses:   api.TaskStatusLegendData,
	}, nil
}

// GetDailySchedule returns the full daily schedule response from the database.
func (s *PgStore) GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*api.DailyScheduleBody, error) {
	peopleRows, err := s.queries.GetAllPeople(ctx)
	if err != nil {
		return nil, err
	}
	peopleMap := make(map[string]db.Person)
	for _, p := range peopleRows {
		peopleMap[p.ID] = p
	}

	cardRows, err := s.queries.GetDailyScheduleTaskCards(ctx, db.GetDailyScheduleTaskCardsParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		Days:      int32(days),
	})
	if err != nil {
		return nil, err
	}

	assignmentRows, err := s.queries.GetDailyScheduleAssignments(ctx)
	if err != nil {
		return nil, err
	}

	assignMap := make(map[int32][]string)
	for _, a := range assignmentRows {
		assignMap[a.TaskCardID] = append(assignMap[a.TaskCardID], a.PersonID)
	}

	// Fetch available counts per date.
	availCountRows, err := s.queries.GetAvailableCountByDate(ctx, db.GetAvailableCountByDateParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		Days:      int32(days),
	})
	if err != nil {
		return nil, err
	}
	availCountMap := make(map[string]int)
	for _, ac := range availCountRows {
		if ac.Date.Valid {
			availCountMap[ac.Date.Time.Format("2006-01-02")] = int(ac.AvailableCount)
		}
	}

	// Group cards by scheduled_date for date-based lookup.
	dateCards := make(map[string][]db.GetDailyScheduleTaskCardsRow)
	for _, cr := range cardRows {
		if cr.ScheduledDate.Valid {
			dateStr := cr.ScheduledDate.Time.Format("2006-01-02")
			dateCards[dateStr] = append(dateCards[dateStr], cr)
		}
	}

	endDate := startDate.AddDate(0, 0, days-1)

	scheduleDays := make([]api.ScheduleDay, days)
	for d := 0; d < days; d++ {
		date := startDate.AddDate(0, 0, d)
		dateStr := date.Format("2006-01-02")

		cards := dateCards[dateStr]

		tasks := make([]api.TaskCard, 0, len(cards))
		for _, cr := range cards {
			assigneeIDs := assignMap[cr.ID]
			assignedPeople := make([]api.AssignedPerson, 0, len(assigneeIDs))
			for _, pid := range assigneeIDs {
				if p, ok := peopleMap[pid]; ok {
					assignedPeople = append(assignedPeople, api.AssignedPerson{
						ID:       p.ID,
						Name:     p.Name,
						Initials: p.Initials,
					})
				}
			}
			assignedCount := len(assignedPeople)
			staffingStatus := "underStaffed"
			if assignedCount == int(cr.PeopleNeeded) {
				staffingStatus = "fullyStaffed"
			}
			var taskIDPtr *string
			if cr.TaskID.Valid {
				taskIDPtr = &cr.TaskID.String
			}
			tasks = append(tasks, api.TaskCard{
				ID:             fmt.Sprintf("sched-%d", cr.ID),
				Title:          cr.Title,
				Priority:       cr.Priority,
				Area:           api.Area{ID: cr.AreaID, Name: cr.AreaName},
				AssignedPeople: assignedPeople,
				PeopleNeeded:   int(cr.PeopleNeeded),
				AssignedCount:  assignedCount,
				StaffingStatus: staffingStatus,
				TaskId:         taskIDPtr,
			})
		}

		availableCount := availCountMap[dateStr]

		scheduleDays[d] = api.ScheduleDay{
			Date:                 dateStr,
			Label:                api.FormatDayLabel(date),
			AvailablePeopleCount: availableCount,
			Tasks:                tasks,
		}
	}

	return &api.DailyScheduleBody{
		Range: api.ScheduleRange{
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   endDate.Format("2006-01-02"),
			Days:      days,
		},
		Days: scheduleDays,
	}, nil
}

// CreatePerson inserts a new person row and returns the generated ID.
func (s *PgStore) CreatePerson(ctx context.Context, name, initials string) (string, error) {
	return s.queries.CreatePerson(ctx, db.CreatePersonParams{
		Name:     name,
		Initials: initials,
	})
}

// UpdatePerson updates an existing person's name and initials.
func (s *PgStore) UpdatePerson(ctx context.Context, id, name, initials string) error {
	return s.queries.UpdatePerson(ctx, db.UpdatePersonParams{
		ID:       id,
		Name:     name,
		Initials: initials,
	})
}

// DeletePerson removes a person by id.
func (s *PgStore) DeletePerson(ctx context.Context, id string) error {
	return s.queries.DeletePerson(ctx, id)
}

// PersonExists checks whether a person with the given id exists.
func (s *PgStore) PersonExists(ctx context.Context, id string) (bool, error) {
	return s.queries.PersonExists(ctx, id)
}

// AreaExists checks whether a room/area with the given id exists.
func (s *PgStore) AreaExists(ctx context.Context, id string) (bool, error) {
	return s.queries.AreaExists(ctx, id)
}

// PersonHasReferences checks whether a person is referenced by backlog or schedule assignments.
func (s *PgStore) PersonHasReferences(ctx context.Context, id string) (bool, error) {
	backlogRefs, err := s.queries.CheckPersonBacklogReferences(ctx, id)
	if err != nil {
		return false, err
	}
	if backlogRefs {
		return true, nil
	}
	return s.queries.CheckPersonScheduleReferences(ctx, id)
}

// UpsertAvailability inserts or updates an availability row for a person on a given date.
func (s *PgStore) UpsertAvailability(ctx context.Context, personID string, date pgtype.Date, status string) error {
	return s.queries.UpsertAvailability(ctx, db.UpsertAvailabilityParams{
		PersonID: personID,
		Date:     date,
		Status:   status,
	})
}

// DeleteAvailability removes an availability row for a person on a given date.
func (s *PgStore) DeleteAvailability(ctx context.Context, personID string, date pgtype.Date) error {
	return s.queries.DeleteAvailability(ctx, db.DeleteAvailabilityParams{
		PersonID: personID,
		Date:     date,
	})
}

// pgDateToTime converts a pgtype.Date to time.Time.
// Returns zero time if the date is not valid.
func pgDateToTime(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}

// ---------- Room CRUD ----------

// ListRooms returns all room/area records ordered by name.
func (s *PgStore) ListRooms(ctx context.Context) ([]api.Room, error) {
	rows, err := s.queries.ListRooms(ctx)
	if err != nil {
		return nil, err
	}
	rooms := make([]api.Room, len(rows))
	for i, r := range rows {
		rooms[i] = dbRoomToAPI(r)
	}
	return rooms, nil
}

// CreateRoom creates a new room/area record and returns it.
func (s *PgStore) CreateRoom(ctx context.Context, input api.CreateRoomInput) (*api.Room, error) {
	r, err := s.queries.CreateRoom(ctx, db.CreateRoomParams{
		Name: input.Name,
		Type: input.Type,
	})
	if err != nil {
		return nil, err
	}
	room := dbRoomToAPI(r)
	return &room, nil
}

// UpdateRoom updates an existing room/area record by ID.
func (s *PgStore) UpdateRoom(ctx context.Context, id string, input api.UpdateRoomInput) (*api.Room, error) {
	r, err := s.queries.UpdateRoom(ctx, db.UpdateRoomParams{
		ID:   id,
		Name: input.Name,
		Type: input.Type,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrRoomNotFound
		}
		return nil, err
	}
	room := dbRoomToAPI(r)
	return &room, nil
}

// DeleteRoom deletes a room/area record by ID.
func (s *PgStore) DeleteRoom(ctx context.Context, id string) error {
	_, err := s.queries.DeleteRoom(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.ErrRoomNotFound
		}
		return err
	}
	return nil
}

// ---------- Task CRUD ----------

// CreateTask creates a new backlog task with server-assigned ID and next sort_order,
// then inserts assignment rows — all in a single transaction so that a partial
// assignment failure does not leave an orphaned task row.
func (s *PgStore) CreateTask(ctx context.Context, input api.CreateTaskInput) (*api.TaskRow, error) {
	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer pgxTx.Rollback(ctx) //nolint:errcheck

	tx := s.queries.WithTx(pgxTx)

	maxSort, err := tx.GetMaxSortOrder(ctx)
	if err != nil {
		return nil, fmt.Errorf("get max sort order: %w", err)
	}

	taskRow, err := tx.CreateTask(ctx, db.CreateTaskParams{
		Title:        input.Title,
		Priority:     input.Priority,
		PeopleNeeded: int32(input.PeopleNeeded),
		AreaID:       input.AreaID,
		Status:       input.Status,
		SortOrder:    maxSort + 1,
	})
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	for i, pid := range input.AssignedTo {
		if err := tx.CreateTaskAssignment(ctx, db.CreateTaskAssignmentParams{
			TaskID:    taskRow.ID,
			PersonID:  pid,
			SortOrder: int32(i),
		}); err != nil {
			return nil, fmt.Errorf("create task assignment: %w", err)
		}
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	assigned := input.AssignedTo
	if assigned == nil {
		assigned = []string{}
	}

	return &api.TaskRow{
		ID:           taskRow.ID,
		Title:        taskRow.Title,
		Priority:     taskRow.Priority,
		PeopleNeeded: int(taskRow.PeopleNeeded),
		Area:         api.Area{ID: taskRow.AreaID, Name: taskRow.AreaName},
		Status:       taskRow.Status,
		AssignedTo:   assigned,
	}, nil
}

// UpdateTask updates a backlog task and replaces assignments transactionally.
func (s *PgStore) UpdateTask(ctx context.Context, id string, input api.UpdateTaskInput) (*api.TaskRow, error) {
	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer pgxTx.Rollback(ctx) //nolint:errcheck

	tx := s.queries.WithTx(pgxTx)

	// Check task exists inside the transaction to avoid a TOCTOU race
	// where a concurrent delete removes the row between check and update.
	_, err = tx.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrTaskNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}

	if err := tx.DeleteTaskAssignments(ctx, id); err != nil {
		return nil, fmt.Errorf("delete task assignments: %w", err)
	}

	for i, pid := range input.AssignedTo {
		if err := tx.CreateTaskAssignment(ctx, db.CreateTaskAssignmentParams{
			TaskID:    id,
			PersonID:  pid,
			SortOrder: int32(i),
		}); err != nil {
			return nil, fmt.Errorf("create task assignment: %w", err)
		}
	}

	taskRow, err := tx.UpdateTask(ctx, db.UpdateTaskParams{
		ID:           id,
		Title:        input.Title,
		Priority:     input.Priority,
		PeopleNeeded: int32(input.PeopleNeeded),
		AreaID:       input.AreaID,
		Status:       input.Status,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrTaskNotFound
		}
		return nil, fmt.Errorf("update task: %w", err)
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	assigned := input.AssignedTo
	if assigned == nil {
		assigned = []string{}
	}

	return &api.TaskRow{
		ID:           taskRow.ID,
		Title:        taskRow.Title,
		Priority:     taskRow.Priority,
		PeopleNeeded: int(taskRow.PeopleNeeded),
		Area:         api.Area{ID: taskRow.AreaID, Name: taskRow.AreaName},
		Status:       taskRow.Status,
		AssignedTo:   assigned,
	}, nil
}

// DeleteTask removes a backlog task and its assignments in a single transaction
// so that assignment rows are never orphaned if the task-delete query fails.
// Returns an error if the task has referencing schedule cards.
func (s *PgStore) DeleteTask(ctx context.Context, id string) error {
	// Check for referencing schedule cards before starting the transaction.
	hasCards, err := s.queries.TaskHasScheduleCards(ctx, pgtype.Text{String: id, Valid: true})
	if err != nil {
		return fmt.Errorf("check task schedule references: %w", err)
	}
	if hasCards {
		return fmt.Errorf("%w: task '%s' has referencing schedule cards. Remove them first", api.ErrTaskHasScheduleCards, id)
	}

	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer pgxTx.Rollback(ctx) //nolint:errcheck

	tx := s.queries.WithTx(pgxTx)

	_, err = tx.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.ErrTaskNotFound
		}
		return fmt.Errorf("get task by id: %w", err)
	}

	if err := tx.DeleteTaskAssignments(ctx, id); err != nil {
		return fmt.Errorf("delete task assignments: %w", err)
	}

	if err := tx.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// ---------- Schedule-card CRUD ----------

// CreateScheduleCard creates a new schedule card with assignments in a single transaction.
// If input.TaskId is set, title/priority/roomArea/peopleNeeded inherit from the referenced
// backlog task unless explicit (non-empty/non-zero) values are supplied.
// Inheritance resolution runs inside the transaction to avoid TOCTOU races.
func (s *PgStore) CreateScheduleCard(ctx context.Context, input api.CreateScheduleCardInput) (*api.TaskCard, error) {
	scheduledDate, err := time.Parse("2006-01-02", input.ScheduledDate)
	if err != nil {
		return nil, fmt.Errorf("parse scheduledDate: %w", err)
	}

	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer pgxTx.Rollback(ctx) //nolint:errcheck

	tx := s.queries.WithTx(pgxTx)

	// Resolve inherited fields from referenced backlog task (inside the transaction).
	title := input.Title
	priority := input.Priority
	areaID := input.AreaId
	peopleNeeded := input.PeopleNeeded

	if input.TaskId != "" {
		refTask, err := tx.GetTaskByIDForRef(ctx, input.TaskId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("referenced task '%s' not found", input.TaskId)
			}
			return nil, fmt.Errorf("get referenced task: %w", err)
		}

		// Inherit only when the caller did not provide an explicit value.
		if title == "" {
			title = refTask.Title
		}
		if priority == "" {
			priority = refTask.Priority
		}
		if areaID == "" {
			areaID = refTask.AreaID
		}
		if peopleNeeded < 1 {
			peopleNeeded = int(refTask.PeopleNeeded)
		}
	}

	maxSort, err := tx.GetMaxScheduleSortOrder(ctx)
	if err != nil {
		return nil, fmt.Errorf("get max sort order: %w", err)
	}

	cardRow, err := tx.CreateScheduleCard(ctx, db.CreateScheduleCardParams{
		Title:         title,
		Priority:      priority,
		AreaID:        areaID,
		PeopleNeeded:  int32(peopleNeeded),
		ScheduledDate: pgtype.Date{Time: scheduledDate, Valid: true},
		SortOrder:     maxSort + 1,
		TaskID:        pgtype.Text{String: input.TaskId, Valid: input.TaskId != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("create schedule card: %w", err)
	}

	for i, pid := range input.AssignedTo {
		if err := tx.CreateScheduleAssignment(ctx, db.CreateScheduleAssignmentParams{
			TaskCardID: cardRow.ID,
			PersonID:   pid,
			SortOrder:  int32(i),
		}); err != nil {
			return nil, fmt.Errorf("create schedule assignment: %w", err)
		}
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return s.buildTaskCardResponse(cardRow.ID, title, priority, cardRow.AreaID, cardRow.AreaName, peopleNeeded, input.AssignedTo, cardRow.TaskID, cardRow.Completed, ctx)
}

// UpdateScheduleCard updates a schedule card and replaces assignments transactionally.
// The referenced backlog task's fields are resolved inside the transaction to avoid TOCTOU races.
// When input.TaskId is empty, the existing card's taskId is preserved to prevent silent reference clearing.
func (s *PgStore) UpdateScheduleCard(ctx context.Context, idStr string, input api.CreateScheduleCardInput) (*api.TaskCard, error) {
	// Parse the sched-{id} identifier.
	var id int32
	if _, err := fmt.Sscanf(idStr, "sched-%d", &id); err != nil || id <= 0 {
		return nil, api.ErrScheduleCardNotFound
	}

	scheduledDate, err := time.Parse("2006-01-02", input.ScheduledDate)
	if err != nil {
		return nil, fmt.Errorf("parse scheduledDate: %w", err)
	}

	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer pgxTx.Rollback(ctx) //nolint:errcheck

	tx := s.queries.WithTx(pgxTx)

	// Check card exists and preserve its sort order and existing taskId.
	existing, err := tx.GetScheduleCardByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrScheduleCardNotFound
		}
		return nil, fmt.Errorf("get schedule card by id: %w", err)
	}

	// Determine effective taskId: use input if provided, otherwise preserve the existing reference.
	effectiveTaskID := input.TaskId
	if effectiveTaskID == "" && existing.TaskID.Valid {
		effectiveTaskID = existing.TaskID.String
	}

	// Resolve inherited fields from referenced backlog task (inside the transaction).
	title := input.Title
	priority := input.Priority
	areaID := input.AreaId
	peopleNeeded := input.PeopleNeeded

	if effectiveTaskID != "" {
		refTask, err := tx.GetTaskByIDForRef(ctx, effectiveTaskID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("referenced task '%s' not found", effectiveTaskID)
			}
			return nil, fmt.Errorf("get referenced task: %w", err)
		}

		// Inherit only when the caller did not provide an explicit value.
		if title == "" {
			title = refTask.Title
		}
		if priority == "" {
			priority = refTask.Priority
		}
		if areaID == "" {
			areaID = refTask.AreaID
		}
		if peopleNeeded < 1 {
			peopleNeeded = int(refTask.PeopleNeeded)
		}
	}

	// Replace assignment rows.
	if err := tx.DeleteScheduleAssignments(ctx, id); err != nil {
		return nil, fmt.Errorf("delete schedule assignments: %w", err)
	}

	for i, pid := range input.AssignedTo {
		if err := tx.CreateScheduleAssignment(ctx, db.CreateScheduleAssignmentParams{
			TaskCardID: id,
			PersonID:   pid,
			SortOrder:  int32(i),
		}); err != nil {
			return nil, fmt.Errorf("create schedule assignment: %w", err)
		}
	}

	cardRow, err := tx.UpdateScheduleCard(ctx, db.UpdateScheduleCardParams{
		ID:            id,
		Title:         title,
		Priority:      priority,
		AreaID:        areaID,
		PeopleNeeded:  int32(peopleNeeded),
		ScheduledDate: pgtype.Date{Time: scheduledDate, Valid: true},
		SortOrder:     existing.SortOrder, // Preserve existing sort order.
		TaskID:        pgtype.Text{String: effectiveTaskID, Valid: effectiveTaskID != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrScheduleCardNotFound
		}
		return nil, fmt.Errorf("update schedule card: %w", err)
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return s.buildTaskCardResponse(cardRow.ID, cardRow.Title, cardRow.Priority, cardRow.AreaID, cardRow.AreaName, int(cardRow.PeopleNeeded), input.AssignedTo, cardRow.TaskID, cardRow.Completed, ctx)
}

// DeleteScheduleCard removes a schedule card and its assignments in a single transaction.
func (s *PgStore) DeleteScheduleCard(ctx context.Context, idStr string) error {
	var id int32
	if _, err := fmt.Sscanf(idStr, "sched-%d", &id); err != nil || id <= 0 {
		return api.ErrScheduleCardNotFound
	}

	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer pgxTx.Rollback(ctx) //nolint:errcheck

	tx := s.queries.WithTx(pgxTx)

	_, err = tx.GetScheduleCardByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.ErrScheduleCardNotFound
		}
		return fmt.Errorf("get schedule card by id: %w", err)
	}

	if err := tx.DeleteScheduleAssignments(ctx, id); err != nil {
		return fmt.Errorf("delete schedule assignments: %w", err)
	}

	if err := tx.DeleteScheduleCard(ctx, id); err != nil {
		return fmt.Errorf("delete schedule card: %w", err)
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// SetScheduleCardCompleted sets the completed status of a schedule card.
func (s *PgStore) SetScheduleCardCompleted(ctx context.Context, idStr string, completed bool) error {
	var id int32
	if _, err := fmt.Sscanf(idStr, "sched-%d", &id); err != nil || id <= 0 {
		return api.ErrScheduleCardNotFound
	}

	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer pgxTx.Rollback(ctx) //nolint:errcheck

	tx := s.queries.WithTx(pgxTx)

	// Verify the card exists.
	_, err = tx.GetScheduleCardByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.ErrScheduleCardNotFound
		}
		return fmt.Errorf("get schedule card by id: %w", err)
	}

	if err := tx.SetScheduleCardCompleted(ctx, db.SetScheduleCardCompletedParams{
		ID:        id,
		Completed: completed,
	}); err != nil {
		return fmt.Errorf("set schedule card completed: %w", err)
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// buildTaskCardResponse resolves assignee identities and constructs the API TaskCard.
func (s *PgStore) buildTaskCardResponse(id int32, title, priority, areaID, areaName string, peopleNeeded int, assigneeIDs []string, taskID pgtype.Text, completed bool, ctx context.Context) (*api.TaskCard, error) {
	// Fetch all people once and build a lookup map.
	peopleRows, err := s.queries.GetAllPeople(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all people: %w", err)
	}
	peopleByID := make(map[string]db.Person, len(peopleRows))
	for _, p := range peopleRows {
		peopleByID[p.ID] = p
	}

	assignees := make([]api.AssignedPerson, 0, len(assigneeIDs))
	for _, pid := range assigneeIDs {
		if p, ok := peopleByID[pid]; ok {
			assignees = append(assignees, api.AssignedPerson{
				ID:       p.ID,
				Name:     p.Name,
				Initials: p.Initials,
			})
		}
	}
	assignedCount := len(assignees)
	staffingStatus := "underStaffed"
	if assignedCount == peopleNeeded {
		staffingStatus = "fullyStaffed"
	}

	var taskIDPtr *string
	if taskID.Valid {
		taskIDPtr = &taskID.String
	}

	return &api.TaskCard{
		ID:             fmt.Sprintf("sched-%d", id),
		Title:          title,
		Priority:       priority,
		Area:           api.Area{ID: areaID, Name: areaName},
		AssignedPeople: assignees,
		PeopleNeeded:   peopleNeeded,
		AssignedCount:  assignedCount,
		StaffingStatus: staffingStatus,
		Completed:      completed,
		TaskId:         taskIDPtr,
	}, nil
}

// ---------- Task reference checks ----------

// TaskExists checks whether a backlog task with the given id exists.
func (s *PgStore) TaskExists(ctx context.Context, id string) (bool, error) {
	return s.queries.TaskExists(ctx, id)
}

// TaskHasScheduleCards checks whether a backlog task is referenced by any schedule cards.
func (s *PgStore) TaskHasScheduleCards(ctx context.Context, id string) (bool, error) {
	return s.queries.TaskHasScheduleCards(ctx, pgtype.Text{String: id, Valid: true})
}

// ---------- Tool CRUD ----------

// GetTools returns all tools ordered by sort order plus a derived coverage summary.
func (s *PgStore) GetTools(ctx context.Context) (*api.ToolsBody, error) {
	rows, err := s.queries.GetTools(ctx)
	if err != nil {
		return nil, err
	}

	tools := make([]api.Tool, len(rows))
	claimed := 0
	for i, r := range rows {
		tools[i] = dbToolToAPI(r)
		if tools[i].BroughtBy != nil {
			claimed++
		}
	}

	total := len(tools)
	return &api.ToolsBody{
		Summary: api.ToolSummary{
			Total:   total,
			Claimed: claimed,
			Open:    total - claimed,
		},
		Tools: tools,
	}, nil
}

// CreateTool creates a new tool with a server-assigned ID and append sort order, no bringer.
func (s *PgStore) CreateTool(ctx context.Context, input api.CreateToolInput) (*api.Tool, error) {
	maxSort, err := s.queries.GetMaxToolSortOrder(ctx)
	if err != nil {
		return nil, fmt.Errorf("get max tool sort order: %w", err)
	}
	row, err := s.queries.CreateTool(ctx, db.CreateToolParams{
		Name:      input.Name,
		SortOrder: maxSort + 1,
	})
	if err != nil {
		return nil, fmt.Errorf("create tool: %w", err)
	}
	tool := dbToolToAPI(row)
	return &tool, nil
}

// UpdateTool updates a tool's name and sort order.
func (s *PgStore) UpdateTool(ctx context.Context, id string, input api.UpdateToolInput) (*api.Tool, error) {
	row, err := s.queries.UpdateTool(ctx, db.UpdateToolParams{
		ID:        id,
		Name:      input.Name,
		SortOrder: int32(input.SortOrder),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrToolNotFound
		}
		return nil, fmt.Errorf("update tool: %w", err)
	}
	tool := dbToolToAPI(row)
	return &tool, nil
}

// DeleteTool removes a tool by ID.
func (s *PgStore) DeleteTool(ctx context.Context, id string) error {
	_, err := s.queries.DeleteTool(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.ErrToolNotFound
		}
		return fmt.Errorf("delete tool: %w", err)
	}
	return nil
}

// ToolExists checks whether a tool with the given ID exists.
func (s *PgStore) ToolExists(ctx context.Context, id string) (bool, error) {
	_, err := s.queries.GetToolByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SetToolBringer sets a tool's bringer, replacing any existing one.
func (s *PgStore) SetToolBringer(ctx context.Context, id, personID string) (*api.Tool, error) {
	row, err := s.queries.SetToolBringer(ctx, db.SetToolBringerParams{
		ID:        id,
		BroughtBy: pgtype.Text{String: personID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrToolNotFound
		}
		return nil, fmt.Errorf("set tool bringer: %w", err)
	}
	tool := dbToolToAPI(row)
	return &tool, nil
}

// ClearToolBringer clears a tool's bringer, returning it to open. Idempotent
// with respect to an already-open tool, but returns ErrToolNotFound when the
// tool does not exist.
func (s *PgStore) ClearToolBringer(ctx context.Context, id string) error {
	_, err := s.queries.ClearToolBringer(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.ErrToolNotFound
		}
		return fmt.Errorf("clear tool bringer: %w", err)
	}
	return nil
}

// dbToolToAPI converts a db.Tool to the API-facing Tool.
func dbToolToAPI(t db.Tool) api.Tool {
	var broughtBy *string
	if t.BroughtBy.Valid {
		broughtBy = &t.BroughtBy.String
	}
	return api.Tool{
		ID:        t.ID,
		Name:      t.Name,
		BroughtBy: broughtBy,
	}
}

// dbRoomToAPI converts a db.RoomsArea to the API-facing Room.
func dbRoomToAPI(r db.RoomsArea) api.Room {
	createdAt := ""
	if r.CreatedAt.Valid {
		createdAt = r.CreatedAt.Time.Format(time.RFC3339)
	}
	updatedAt := ""
	if r.UpdatedAt.Valid {
		updatedAt = r.UpdatedAt.Time.Format(time.RFC3339)
	}
	return api.Room{
		ID:        r.ID,
		Name:      r.Name,
		Type:      r.Type,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
