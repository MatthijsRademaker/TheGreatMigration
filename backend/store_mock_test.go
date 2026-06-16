package main

import (
	"context"
	"fmt"
	"sort"
	"time"
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
var seedTasks = []TaskRow{
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

func findPersonByID(id string) (AssignedPerson, bool) {
	for _, sp := range seedPeople {
		if sp.Id == id {
			return AssignedPerson{ID: sp.Id, Name: sp.Name, Initials: sp.Initials}, true
		}
	}
	return AssignedPerson{}, false
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
}) []Person {
	people := make([]Person, len(mp))
	for i, sp := range mp {
		avail := make([]AvailabilityEntry, days)
		for d := 0; d < days; d++ {
			date := startDate.AddDate(0, 0, d)
			avail[d] = AvailabilityEntry{
				Date:   date.Format("2006-01-02"),
				Status: sp.Status(d),
			}
		}
		people[i] = Person{
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
	planningWindow *PlanningWindowBody
	taskBacklog    *TaskBacklogBody
	rooms          map[string]Room
	nextRoomID     int
}

func newMockStore() *mockStore {
	startDate, _ := time.Parse("2006-01-02", "2026-07-05")
	endDate, _ := time.Parse("2006-01-02", "2026-08-13")
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	// Pre-compute task backlog summary from seed tasks.
	total := len(seedTasks)
	highPriority := 0
	unassigned := 0
	understaffed := 0
	for _, t := range seedTasks {
		if t.Priority == "high" {
			highPriority++
		}
		if len(t.AssignedTo) == 0 {
			unassigned++
		} else if len(t.AssignedTo) < t.PeopleNeeded {
			understaffed++
		}
	}

	return &mockStore{
		planningWindow: &PlanningWindowBody{
			StartDate: "2026-07-05",
			EndDate:   "2026-08-13",
			Days:      days,
		},
		taskBacklog: &TaskBacklogBody{
			Summary: TaskSummary{
				TotalTasks:        total,
				HighPriorityTasks: highPriority,
				UnassignedTasks:   unassigned,
				UnderstaffedTasks: understaffed,
			},
			Tasks:      seedTasks,
			Priorities: priorityLegend,
			Statuses:   taskStatusLegend,
		},
		rooms: map[string]Room{
			"room-1": {ID: "room-1", Name: "Kitchen", Type: "room", CreatedAt: "2026-01-01T00:00:00Z", UpdatedAt: "2026-01-01T00:00:00Z"},
			"room-2": {ID: "room-2", Name: "Living Room", Type: "room", CreatedAt: "2026-01-01T00:00:00Z", UpdatedAt: "2026-01-01T00:00:00Z"},
		},
		nextRoomID: 3,
	}
}

func (m *mockStore) GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error) {
	return m.planningWindow, nil
}

func (m *mockStore) GetTaskBacklog(ctx context.Context) (*TaskBacklogBody, error) {
	return m.taskBacklog, nil
}

func (m *mockStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*DashboardBody, error) {
	endDate := startDate.AddDate(0, 0, days-1)
	selectedDate := startDate.Format("2006-01-02")
	people := buildMockPeople(startDate, days, seedPeople)

	availableToday := 0
	for _, p := range people {
		for _, e := range p.Availability {
			if e.Date == selectedDate && e.Status == "available" {
				availableToday++
				break
			}
		}
	}

	return &DashboardBody{
		Range: Range{
			StartDate:    startDate.Format("2006-01-02"),
			EndDate:      endDate.Format("2006-01-02"),
			Days:         days,
			SelectedDate: selectedDate,
		},
		Summary: Summary{
			AvailableToday: availableToday,
			TotalPeople:    len(people),
		},
		People:   people,
		Statuses: statusLegend,
	}, nil
}

func (m *mockStore) GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*DailyScheduleBody, error) {
	endDate := startDate.AddDate(0, 0, days-1)

	scheduleDays := make([]ScheduleDay, days)
	for d := 0; d < days; d++ {
		date := startDate.AddDate(0, 0, d)
		dateStr := date.Format("2006-01-02")

		availableCount := countAvailableForDay(d)

		templates := seedTasksForDay(d)
		tasks := make([]TaskCard, len(templates))
		for ti, tmpl := range templates {
			assignees := make([]AssignedPerson, 0, len(tmpl.assigneeIds))
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

			tasks[ti] = TaskCard{
				ID:             fmt.Sprintf("task-d%d-%d", d, ti),
				Title:          tmpl.title,
				Priority:       tmpl.priority,
				RoomArea:       tmpl.roomArea,
				AssignedPeople: assignees,
				PeopleNeeded:   tmpl.peopleNeeded,
				AssignedCount:  assignedCount,
				StaffingStatus: staffingStatus,
			}
		}

		scheduleDays[d] = ScheduleDay{
			Date:                 dateStr,
			Label:                formatDayLabel(date),
			AvailablePeopleCount: availableCount,
			Tasks:                tasks,
		}
	}

	return &DailyScheduleBody{
		Range: ScheduleRange{
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   endDate.Format("2006-01-02"),
			Days:      days,
		},
		Days: scheduleDays,
	}, nil
}

// ---------- Room CRUD (mockStore) ----------

func (m *mockStore) ListRooms(ctx context.Context) ([]Room, error) {
	rooms := make([]Room, 0, len(m.rooms))
	for _, r := range m.rooms {
		rooms = append(rooms, r)
	}
	sort.Slice(rooms, func(i, j int) bool { return rooms[i].Name < rooms[j].Name })
	return rooms, nil
}

func (m *mockStore) CreateRoom(ctx context.Context, input CreateRoomInput) (*Room, error) {
	id := fmt.Sprintf("room-%d", m.nextRoomID)
	m.nextRoomID++
	now := time.Now().UTC().Format(time.RFC3339)
	r := Room{
		ID:        id,
		Name:      input.Name,
		Type:      input.Type,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.rooms[id] = r
	return &r, nil
}

func (m *mockStore) UpdateRoom(ctx context.Context, id string, input UpdateRoomInput) (*Room, error) {
	r, ok := m.rooms[id]
	if !ok {
		return nil, ErrRoomNotFound
	}
	r.Name = input.Name
	r.Type = input.Type
	m.rooms[id] = r
	return &r, nil
}

func (m *mockStore) DeleteRoom(ctx context.Context, id string) error {
	if _, ok := m.rooms[id]; !ok {
		return ErrRoomNotFound
	}
	delete(m.rooms, id)
	return nil
}
