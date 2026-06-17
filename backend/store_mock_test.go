package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/user/the-great-migration/backend/api"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------- Seed data (test-only, used by mockStore) ----------

// seedPeople is the canonical set of people for test mocks.
// Mirrors the seed data in migrations/002_seed_demo_data.sql.
var seedPeople = []struct {
	Id       string
	Name     string
	Initials string
	// Status returns the availability status for a given date offset from start.
	Status func(dayOffset int) string
}{
	{Id: "p1", Name: "Sophia Chen", Initials: "SC", Status: always("available")},
	{Id: "p2", Name: "Marcus Rivera", Initials: "MR", Status: always("available")},
	{Id: "p3", Name: "Elena Kowalski", Initials: "EK", Status: always("available")},
	{Id: "p4", Name: "James Okafor", Initials: "JO", Status: always("available")},
	{Id: "p5", Name: "Priya Nair", Initials: "PN", Status: always("available")},
	{Id: "p6", Name: "Thomas Berg", Initials: "TB", Status: always("available")},
	{Id: "p7", Name: "Amara Diallo", Initials: "AD", Status: always("busy")},
	{Id: "p8", Name: "Noah Larsson", Initials: "NL", Status: cycleStatuses},
}

func always(s string) func(int) string {
	return func(_ int) string { return s }
}

var cycleOrder = []string{"off", "partial", "busy", "available"}

func cycleStatuses(dayOffset int) string {
	return cycleOrder[dayOffset%len(cycleOrder)]
}

// seedTasks is the canonical set of tasks for test mocks.
var seedTasks = []api.TaskRow{
	{ID: "task-1", Title: "Disconnect kitchen appliances", Priority: "high", PeopleNeeded: 3, Room: "Kitchen", Status: "backlog", AssignedTo: []string{}},
	{ID: "task-2", Title: "Wrap living room furniture", Priority: "high", PeopleNeeded: 2, Room: "Living Room", Status: "ready", AssignedTo: []string{}},
	{ID: "task-3", Title: "Pack kitchen fragile items", Priority: "high", PeopleNeeded: 2, Room: "Kitchen", Status: "assigned", AssignedTo: []string{"p1"}},
	{ID: "task-4", Title: "Disassemble bedroom furniture", Priority: "high", PeopleNeeded: 2, Room: "Bedroom 1", Status: "assigned", AssignedTo: []string{"p2", "p3"}},
	{ID: "task-5", Title: "Sort and label moving boxes", Priority: "medium", PeopleNeeded: 3, Room: "Living Room", Status: "backlog", AssignedTo: []string{"p4"}},
	{ID: "task-6", Title: "Clear garage shelving", Priority: "medium", PeopleNeeded: 1, Room: "Garage", Status: "ready", AssignedTo: []string{}},
	{ID: "task-7", Title: "Move bedroom wardrobe", Priority: "medium", PeopleNeeded: 3, Room: "Bedroom 2", Status: "assigned", AssignedTo: []string{"p5", "p6", "p7"}},
	{ID: "task-8", Title: "Sweep garage floor", Priority: "low", PeopleNeeded: 1, Room: "Garage", Status: "backlog", AssignedTo: []string{"p8"}},
	{ID: "task-9", Title: "Dust living room shelves", Priority: "low", PeopleNeeded: 2, Room: "Living Room", Status: "ready", AssignedTo: []string{}},
	{ID: "task-10", Title: "Wipe down kitchen counters", Priority: "low", PeopleNeeded: 2, Room: "Kitchen", Status: "assigned", AssignedTo: []string{"p1", "p2"}},
	{ID: "task-11", Title: "Inventory bedroom closet", Priority: "medium", PeopleNeeded: 3, Room: "Bedroom 1", Status: "ready", AssignedTo: []string{"p3"}},
}

// Seed daily-schedule helpers.

type seedTaskTemplate struct {
	title        string
	priority     string
	roomArea     string
	peopleNeeded int
	assigneeIds  []string
}

func seedTasksForDay(dayOffset int) []seedTaskTemplate {
	dayGroups := [][]seedTaskTemplate{
		{
			{title: "Kitchen deep clean", priority: "high", roomArea: "Kitchen", peopleNeeded: 2, assigneeIds: []string{"p1", "p2"}},
			{title: "Window washing", priority: "medium", roomArea: "Living Room", peopleNeeded: 2, assigneeIds: []string{"p3"}},
		},
		{
			{title: "Furniture assembly", priority: "high", roomArea: "Bedroom", peopleNeeded: 1, assigneeIds: []string{"p4"}},
			{title: "Packing supplies inventory", priority: "low", roomArea: "Storage", peopleNeeded: 3, assigneeIds: []string{"p5"}},
		},
		{
			{title: "Electronics setup", priority: "medium", roomArea: "Office", peopleNeeded: 2, assigneeIds: []string{"p6", "p7"}},
			{title: "Declutter garage", priority: "low", roomArea: "Garage", peopleNeeded: 1, assigneeIds: []string{}},
		},
		{
			{title: "Box labeling", priority: "medium", roomArea: "Storage", peopleNeeded: 2, assigneeIds: []string{"p8"}},
			{title: "Curtain removal", priority: "high", roomArea: "Living Room", peopleNeeded: 2, assigneeIds: []string{"p1", "p3"}},
		},
	}
	return dayGroups[dayOffset%len(dayGroups)]
}

func findPersonByID(id string) (api.AssignedPerson, bool) {
	for _, sp := range seedPeople {
		if sp.Id == id {
			return api.AssignedPerson{ID: sp.Id, Name: sp.Name, Initials: sp.Initials}, true
		}
	}
	return api.AssignedPerson{}, false
}

func countAvailableForDay(dayOffset int) int {
	count := 0
	for _, sp := range seedPeople {
		if sp.Status(dayOffset) == "available" {
			count++
		}
	}
	return count
}

func buildMockPeople(startDate time.Time, days int, mp []struct {
	Id, Name, Initials string
	Status             func(int) string
}) []api.Person {
	people := make([]api.Person, len(mp))
	for i, sp := range mp {
		avail := make([]api.AvailabilityEntry, days)
		for d := 0; d < days; d++ {
			date := startDate.AddDate(0, 0, d)
			avail[d] = api.AvailabilityEntry{
				Date:   date.Format("2006-01-02"),
				Status: sp.Status(d),
			}
		}
		people[i] = api.Person{
			ID:           sp.Id,
			Name:         sp.Name,
			Initials:     sp.Initials,
			Availability: avail,
		}
	}
	return people
}

// ---------- mockStore ----------

// mockStore implements Store for fast unit tests without a database.
// Parameterless methods return pre-computed struct fields.
// Parameterised methods compute from seed data on each call.
type mockStore struct {
	planningWindow   *api.PlanningWindowBody
	tasks            map[string]api.TaskRow
	nextTaskID       int
	rooms            map[string]api.Room
	nextRoomID       int
	scheduleCards    map[string]api.TaskCard
	scheduleDates    map[string]string // card ID -> scheduledDate string
	scheduleTaskRefs map[string]string // card ID -> referenced task ID
	nextScheduleID   int
}

func newMockStore() *mockStore {
	startDate, _ := time.Parse("2006-01-02", "2026-07-05")
	endDate, _ := time.Parse("2006-01-02", "2026-08-13")
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	// Build mutable task map from seed tasks.
	tasks := make(map[string]api.TaskRow, len(seedTasks))
	for _, t := range seedTasks {
		tasks[t.ID] = t
	}

	return &mockStore{
		planningWindow: &api.PlanningWindowBody{
			StartDate: "2026-07-05",
			EndDate:   "2026-08-13",
			Days:      days,
		},
		tasks:      tasks,
		nextTaskID: 12,
		rooms: map[string]api.Room{
			"room-1": {ID: "room-1", Name: "Kitchen", Type: "room", CreatedAt: "2026-01-01T00:00:00Z", UpdatedAt: "2026-01-01T00:00:00Z"},
			"room-2": {ID: "room-2", Name: "Living Room", Type: "room", CreatedAt: "2026-01-01T00:00:00Z", UpdatedAt: "2026-01-01T00:00:00Z"},
		},
		nextRoomID:       3,
		scheduleCards:    make(map[string]api.TaskCard),
		scheduleDates:    make(map[string]string),
		scheduleTaskRefs: make(map[string]string),
		nextScheduleID:   1,
	}
}

func (m *mockStore) GetPlanningWindow(ctx context.Context) (*api.PlanningWindowBody, error) {
	return m.planningWindow, nil
}

func (m *mockStore) UpdatePlanningWindow(ctx context.Context, startDate, endDate time.Time) (*api.PlanningWindowBody, error) {
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	m.planningWindow = &api.PlanningWindowBody{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		Days:      days,
	}
	return m.planningWindow, nil
}

func (m *mockStore) GetTaskBacklog(ctx context.Context) (*api.TaskBacklogBody, error) {
	// Build the backlog payload dynamically from the mutable task map so
	// that writes via CreateTask / UpdateTask / DeleteTask are reflected.
	tasks := make([]api.TaskRow, 0, len(m.tasks))
	total := 0
	highPriority := 0
	unassigned := 0
	understaffed := 0
	for _, t := range m.tasks {
		tasks = append(tasks, t)
		total++
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

func (m *mockStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int, offset int, limit int) (*api.DashboardBody, error) {
	endDate := startDate.AddDate(0, 0, days-1)
	selectedDate := startDate.Format("2006-01-02")
	allPeople := buildMockPeople(startDate, days, seedPeople)

	total := len(allPeople)

	// Apply pagination slicing.
	var people []api.Person
	if limit > 0 {
		start := offset
		if start > total {
			start = total
		}
		end := offset + limit
		if end > total {
			end = total
		}
		people = allPeople[start:end]
	} else {
		people = allPeople
	}

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
			TotalPeople:    total,
		},
		People:   people,
		Statuses: api.StatusLegendData,
	}, nil
}

func (m *mockStore) GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*api.DailyScheduleBody, error) {
	endDate := startDate.AddDate(0, 0, days-1)

	scheduleDays := make([]api.ScheduleDay, days)
	for d := 0; d < days; d++ {
		date := startDate.AddDate(0, 0, d)
		dateStr := date.Format("2006-01-02")

		availableCount := countAvailableForDay(d)

		templates := seedTasksForDay(d)
		tasks := make([]api.TaskCard, len(templates))
		for ti, tmpl := range templates {
			assignees := make([]api.AssignedPerson, 0, len(tmpl.assigneeIds))
			for _, pid := range tmpl.assigneeIds {
				if p, ok := findPersonByID(pid); ok {
					assignees = append(assignees, p)
				}
			}
			assignedCount := len(assignees)

			staffingStatus := "underStaffed"
			if assignedCount == tmpl.peopleNeeded {
				staffingStatus = "fullyStaffed"
			}

			tasks[ti] = api.TaskCard{
				ID:             fmt.Sprintf("sched-d%d-%d", d, ti),
				Title:          tmpl.title,
				Priority:       tmpl.priority,
				RoomArea:       tmpl.roomArea,
				AssignedPeople: assignees,
				PeopleNeeded:   tmpl.peopleNeeded,
				AssignedCount:  assignedCount,
				StaffingStatus: staffingStatus,
				TaskId:         nil,
			}
		}

		scheduleDays[d] = api.ScheduleDay{
			Date:                 dateStr,
			Label:                api.FormatDayLabel(date),
			AvailablePeopleCount: availableCount,
			Tasks:                tasks,
		}
	}

	// Merge cards from the mutable scheduleCards store by scheduled date.
	for d := 0; d < days; d++ {
		dateStr := startDate.AddDate(0, 0, d).Format("2006-01-02")
		for cardID, card := range m.scheduleCards {
			if m.scheduleDates[cardID] == dateStr {
				scheduleDays[d].Tasks = append(scheduleDays[d].Tasks, card)
			}
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

// Person CRUD mocks (no-op for existing tests).

func (m *mockStore) CreatePerson(ctx context.Context, name, initials string) (string, error) {
	return "p-mock-id", nil
}

func (m *mockStore) UpdatePerson(ctx context.Context, id, name, initials string) error {
	return nil
}

func (m *mockStore) DeletePerson(ctx context.Context, id string) error {
	return nil
}

func (m *mockStore) PersonExists(ctx context.Context, id string) (bool, error) {
	for _, sp := range seedPeople {
		if sp.Id == id {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockStore) PersonHasReferences(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func (m *mockStore) UpsertAvailability(ctx context.Context, personID string, date pgtype.Date, status string) error {
	return nil
}

func (m *mockStore) DeleteAvailability(ctx context.Context, personID string, date pgtype.Date) error {
	return nil
}

// ---------- Task CRUD (mockStore) ----------

func (m *mockStore) CreateTask(ctx context.Context, input api.CreateTaskInput) (*api.TaskRow, error) {
	id := fmt.Sprintf("task-%d", m.nextTaskID)
	m.nextTaskID++

	assigned := input.AssignedTo
	if assigned == nil {
		assigned = []string{}
	}

	task := api.TaskRow{
		ID:           id,
		Title:        input.Title,
		Priority:     input.Priority,
		PeopleNeeded: input.PeopleNeeded,
		Room:         input.Room,
		Status:       input.Status,
		AssignedTo:   assigned,
	}
	m.tasks[id] = task
	return &task, nil
}

func (m *mockStore) UpdateTask(ctx context.Context, id string, input api.UpdateTaskInput) (*api.TaskRow, error) {
	_, ok := m.tasks[id]
	if !ok {
		return nil, api.ErrTaskNotFound
	}

	assigned := input.AssignedTo
	if assigned == nil {
		assigned = []string{}
	}

	task := api.TaskRow{
		ID:           id,
		Title:        input.Title,
		Priority:     input.Priority,
		PeopleNeeded: input.PeopleNeeded,
		Room:         input.Room,
		Status:       input.Status,
		AssignedTo:   assigned,
	}
	m.tasks[id] = task
	return &task, nil
}

func (m *mockStore) DeleteTask(ctx context.Context, id string) error {
	if _, ok := m.tasks[id]; !ok {
		return api.ErrTaskNotFound
	}
	// Check for referencing schedule cards before deleting.
	hasCards, err := m.TaskHasScheduleCards(ctx, id)
	if err != nil {
		return err
	}
	if hasCards {
		return api.ErrTaskHasScheduleCards
	}
	delete(m.tasks, id)
	return nil
}

// ---------- Room CRUD (mockStore) ----------

func (m *mockStore) ListRooms(ctx context.Context) ([]api.Room, error) {
	rooms := make([]api.Room, 0, len(m.rooms))
	for _, r := range m.rooms {
		rooms = append(rooms, r)
	}
	sort.Slice(rooms, func(i, j int) bool { return rooms[i].Name < rooms[j].Name })
	return rooms, nil
}

func (m *mockStore) CreateRoom(ctx context.Context, input api.CreateRoomInput) (*api.Room, error) {
	id := fmt.Sprintf("room-%d", m.nextRoomID)
	m.nextRoomID++
	now := time.Now().UTC().Format(time.RFC3339)
	r := api.Room{
		ID:        id,
		Name:      input.Name,
		Type:      input.Type,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.rooms[id] = r
	return &r, nil
}

func (m *mockStore) UpdateRoom(ctx context.Context, id string, input api.UpdateRoomInput) (*api.Room, error) {
	r, ok := m.rooms[id]
	if !ok {
		return nil, api.ErrRoomNotFound
	}
	r.Name = input.Name
	r.Type = input.Type
	m.rooms[id] = r
	return &r, nil
}

func (m *mockStore) DeleteRoom(ctx context.Context, id string) error {
	if _, ok := m.rooms[id]; !ok {
		return api.ErrRoomNotFound
	}
	delete(m.rooms, id)
	return nil
}

// ---------- Schedule-card CRUD (mockStore) ----------

func (m *mockStore) CreateScheduleCard(ctx context.Context, input api.CreateScheduleCardInput) (*api.TaskCard, error) {
	id := fmt.Sprintf("sched-%d", m.nextScheduleID)
	m.nextScheduleID++

	// Resolve inherited fields from referenced backlog task.
	title := input.Title
	priority := input.Priority
	roomArea := input.RoomArea
	peopleNeeded := input.PeopleNeeded

	if input.TaskId != "" {
		refTask, ok := m.tasks[input.TaskId]
		if !ok {
			return nil, fmt.Errorf("referenced task '%s' not found", input.TaskId)
		}
		if title == "" {
			title = refTask.Title
		}
		if priority == "" {
			priority = refTask.Priority
		}
		if roomArea == "" {
			roomArea = refTask.Room
		}
		if peopleNeeded < 1 {
			peopleNeeded = refTask.PeopleNeeded
		}
	}

	assignees := make([]api.AssignedPerson, 0, len(input.AssignedTo))
	for _, pid := range input.AssignedTo {
		if p, ok := findPersonByID(pid); ok {
			assignees = append(assignees, p)
		}
	}
	assignedCount := len(assignees)
	staffingStatus := "underStaffed"
	if assignedCount == peopleNeeded {
		staffingStatus = "fullyStaffed"
	}

	var taskIDPtr *string
	if input.TaskId != "" {
		taskIDPtr = &input.TaskId
	}

	card := api.TaskCard{
		ID:             id,
		Title:          title,
		Priority:       priority,
		RoomArea:       roomArea,
		AssignedPeople: assignees,
		PeopleNeeded:   peopleNeeded,
		AssignedCount:  assignedCount,
		StaffingStatus: staffingStatus,
		TaskId:         taskIDPtr,
	}
	m.scheduleCards[id] = card
	m.scheduleDates[id] = input.ScheduledDate
	m.scheduleTaskRefs[id] = input.TaskId
	return &card, nil
}

func (m *mockStore) UpdateScheduleCard(ctx context.Context, idStr string, input api.CreateScheduleCardInput) (*api.TaskCard, error) {
	existing, ok := m.scheduleCards[idStr]
	if !ok {
		return nil, api.ErrScheduleCardNotFound
	}

	// Determine effective taskId: use input if provided, otherwise preserve the existing reference.
	effectiveTaskID := input.TaskId
	if effectiveTaskID == "" && existing.TaskId != nil {
		effectiveTaskID = *existing.TaskId
	}

	// Resolve inherited fields from referenced backlog task.
	title := input.Title
	priority := input.Priority
	roomArea := input.RoomArea
	peopleNeeded := input.PeopleNeeded

	if effectiveTaskID != "" {
		refTask, ok := m.tasks[effectiveTaskID]
		if !ok {
			return nil, fmt.Errorf("referenced task '%s' not found", effectiveTaskID)
		}
		if title == "" {
			title = refTask.Title
		}
		if priority == "" {
			priority = refTask.Priority
		}
		if roomArea == "" {
			roomArea = refTask.Room
		}
		if peopleNeeded < 1 {
			peopleNeeded = refTask.PeopleNeeded
		}
	}

	var taskIDPtr *string
	if effectiveTaskID != "" {
		taskIDPtr = &effectiveTaskID
	}

	assignees := make([]api.AssignedPerson, 0, len(input.AssignedTo))
	for _, pid := range input.AssignedTo {
		if p, ok := findPersonByID(pid); ok {
			assignees = append(assignees, p)
		}
	}
	assignedCount := len(assignees)
	staffingStatus := "underStaffed"
	if assignedCount == peopleNeeded {
		staffingStatus = "fullyStaffed"
	}

	card := api.TaskCard{
		ID:             idStr,
		Title:          title,
		Priority:       priority,
		RoomArea:       roomArea,
		AssignedPeople: assignees,
		PeopleNeeded:   peopleNeeded,
		AssignedCount:  assignedCount,
		StaffingStatus: staffingStatus,
		TaskId:         taskIDPtr,
	}
	m.scheduleCards[idStr] = card
	m.scheduleDates[idStr] = input.ScheduledDate
	m.scheduleTaskRefs[idStr] = effectiveTaskID
	return &card, nil
}

func (m *mockStore) DeleteScheduleCard(ctx context.Context, idStr string) error {
	if _, ok := m.scheduleCards[idStr]; !ok {
		return api.ErrScheduleCardNotFound
	}
	delete(m.scheduleCards, idStr)
	delete(m.scheduleDates, idStr)
	delete(m.scheduleTaskRefs, idStr)
	return nil
}

func (m *mockStore) TaskExists(ctx context.Context, id string) (bool, error) {
	_, ok := m.tasks[id]
	return ok, nil
}

func (m *mockStore) TaskHasScheduleCards(ctx context.Context, id string) (bool, error) {
	for _, refID := range m.scheduleTaskRefs {
		if refID == id {
			return true, nil
		}
	}
	return false, nil
}
