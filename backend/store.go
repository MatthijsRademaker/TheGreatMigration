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

// GetPeopleAvailability returns people and availability data for the given date range.
func (s *PgStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*api.DashboardBody, error) {
	// Fetch all people.
	peopleRows, err := s.queries.GetAllPeople(ctx)
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
	selectedDate := startDate.Format("2006-01-02")
	availableToday := 0
	for _, p := range people {
		for _, e := range p.Availability {
			if e.Date == selectedDate && e.Status == "available" {
				availableToday++
				break
			}
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
			TotalPeople:    len(people),
		},
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
			Room:         tr.Room,
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

	cardRows, err := s.queries.GetDailyScheduleTaskCards(ctx)
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

	// Group cards by day_group for modulo-based lookup.
	groupCards := make(map[int32][]db.ScheduleTaskCard)
	for _, cr := range cardRows {
		groupCards[cr.DayGroup] = append(groupCards[cr.DayGroup], cr)
	}

	endDate := startDate.AddDate(0, 0, days-1)

	scheduleDays := make([]api.ScheduleDay, days)
	for d := 0; d < days; d++ {
		date := startDate.AddDate(0, 0, d)
		dateStr := date.Format("2006-01-02")

		dayGroup := int32(d % len(groupCards))
		cards := groupCards[dayGroup]

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
			tasks = append(tasks, api.TaskCard{
				ID:             fmt.Sprintf("task-d%d-%d", d, cr.SortOrder),
				Title:          cr.Title,
				Priority:       cr.Priority,
				RoomArea:       cr.RoomArea,
				AssignedPeople: assignedPeople,
				PeopleNeeded:   int(cr.PeopleNeeded),
				AssignedCount:  assignedCount,
				StaffingStatus: staffingStatus,
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

// CreatePerson inserts a new person row.
func (s *PgStore) CreatePerson(ctx context.Context, id, name, initials string) error {
	return s.queries.CreatePerson(ctx, db.CreatePersonParams{
		ID:       id,
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
		Room:         input.Room,
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
		Room:         taskRow.Room,
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
		Room:         input.Room,
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
		Room:         taskRow.Room,
		Status:       taskRow.Status,
		AssignedTo:   assigned,
	}, nil
}

// DeleteTask removes a backlog task and its assignments in a single transaction
// so that assignment rows are never orphaned if the task-delete query fails.
func (s *PgStore) DeleteTask(ctx context.Context, id string) error {
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
