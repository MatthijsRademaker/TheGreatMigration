package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func newTestAPI(store Store) (chi.Router, huma.API) {
	router := chi.NewMux()
	config := huma.DefaultConfig("Test API", "1.0.0")
	api := humachi.New(router, config)

	huma.Register(api, huma.Operation{
		OperationID: "get-hello",
		Method:      http.MethodGet,
		Path:        "/api/hello",
	}, func(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
		resp := &HelloOutput{}
		resp.Body.Message = "Hello from the backend!"
		return resp, nil
	})

	registerDashboardPeopleAvailability(api, store)
	registerPlanningWindow(api, store)
	registerTasksBacklog(api, store)
	registerDailySchedule(api, store)
	registerRoomsAreas(api, store)

	return router, api
}

func TestHelloEndpoint(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	if body["message"] != "Hello from the backend!" {
		t.Fatalf("expected message 'Hello from the backend!', got %q", body["message"])
	}
}

func TestDashboardPeopleAvailability(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/people-availability", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Happy path: 200 OK and JSON content-type.
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body DashboardBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	// Response shape assertions.
	if body.Range.StartDate == "" {
		t.Fatal("range.startDate is empty")
	}
	if body.Range.EndDate == "" {
		t.Fatal("range.endDate is empty")
	}
	if body.Range.Days != 4 {
		t.Fatalf("expected range.days=4 (default), got %d", body.Range.Days)
	}
	if body.Range.SelectedDate != body.Range.StartDate {
		t.Fatalf("expected selectedDate=%s to equal startDate", body.Range.StartDate)
	}

	if body.Summary.TotalPeople == 0 {
		t.Fatal("summary.totalPeople is 0")
	}
	if len(body.People) != body.Summary.TotalPeople {
		t.Fatalf("len(people)=%d != summary.totalPeople=%d", len(body.People), body.Summary.TotalPeople)
	}

	if len(body.Statuses) != 4 {
		t.Fatalf("expected 4 statuses in legend, got %d", len(body.Statuses))
	}

	// Validate all statuses in the legend.
	canonical := map[string]bool{}
	for _, s := range body.Statuses {
		canonical[s.ID] = true
	}
	expectedStatuses := []string{"available", "busy", "partial", "off"}
	for _, s := range expectedStatuses {
		if !canonical[s] {
			t.Fatalf("statuses legend missing %q", s)
		}
	}

	// Validate all availability statuses are canonical.
	for _, p := range body.People {
		for _, e := range p.Availability {
			if !canonical[e.Status] {
				t.Fatalf("person %s has non-canonical status %q on date %s", p.ID, e.Status, e.Date)
			}
		}
		if p.ID == "" {
			t.Fatal("person has empty id")
		}
		if p.Name == "" {
			t.Fatalf("person %s has empty name", p.ID)
		}
		if p.Initials == "" {
			t.Fatalf("person %s has empty initials", p.ID)
		}
		// Each person should have exactly D availability entries.
		if len(p.Availability) != body.Range.Days {
			t.Fatalf("person %s has %d availability entries, expected %d", p.ID, len(p.Availability), body.Range.Days)
		}
	}

	// Verify summary.availableToday matches actual count on selectedDate.
	actualAvailable := 0
	for _, p := range body.People {
		for _, e := range p.Availability {
			if e.Date == body.Range.SelectedDate && e.Status == "available" {
				actualAvailable++
				break
			}
		}
	}
	if body.Summary.AvailableToday != actualAvailable {
		t.Fatalf("summary.availableToday=%d does not match actual count=%d", body.Summary.AvailableToday, actualAvailable)
	}

	// Verify at least 8 people.
	if len(body.People) < 8 {
		t.Fatalf("expected at least 8 people, got %d", len(body.People))
	}
}

func TestPlanningWindowEndpoint(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/planning-window", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body PlanningWindowBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	if body.StartDate != "2026-07-05" {
		t.Fatalf("expected startDate '2026-07-05', got %q", body.StartDate)
	}
	if body.EndDate != "2026-08-13" {
		t.Fatalf("expected endDate '2026-08-13', got %q", body.EndDate)
	}
	if body.Days != 40 {
		t.Fatalf("expected days=40, got %d", body.Days)
	}

	// Verify startDate lexicographically precedes endDate.
	if body.StartDate >= body.EndDate {
		t.Fatalf("expected startDate < endDate, got startDate=%q endDate=%q", body.StartDate, body.EndDate)
	}
}

func TestTaskBacklog(t *testing.T) {
	router, api := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/backlog", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Happy path: 200 OK and JSON content-type.
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body TaskBacklogBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	// Top-level fields present.
	if len(body.Tasks) == 0 {
		t.Fatal("tasks array is empty")
	}
	if len(body.Priorities) == 0 {
		t.Fatal("priorities array is empty")
	}
	if len(body.Statuses) == 0 {
		t.Fatal("statuses array is empty")
	}

	// Summary count consistency.
	if body.Summary.TotalTasks != len(body.Tasks) {
		t.Fatalf("summary.totalTasks=%d != len(tasks)=%d", body.Summary.TotalTasks, len(body.Tasks))
	}

	highCount := 0
	unassignedCount := 0
	understaffedCount := 0
	for _, t := range body.Tasks {
		if t.Priority == "high" {
			highCount++
		}
		if len(t.AssignedTo) == 0 {
			unassignedCount++
		} else if len(t.AssignedTo) < t.PeopleNeeded {
			understaffedCount++
		}
	}
	if body.Summary.HighPriorityTasks != highCount {
		t.Fatalf("summary.highPriorityTasks=%d != actual high count=%d", body.Summary.HighPriorityTasks, highCount)
	}
	if body.Summary.UnassignedTasks != unassignedCount {
		t.Fatalf("summary.unassignedTasks=%d != actual unassigned count=%d", body.Summary.UnassignedTasks, unassignedCount)
	}
	if body.Summary.UnderstaffedTasks != understaffedCount {
		t.Fatalf("summary.understaffedTasks=%d != actual understaffed count=%d", body.Summary.UnderstaffedTasks, understaffedCount)
	}

	// Canonical priority values in legend.
	priorityIDs := map[string]bool{}
	for _, p := range body.Priorities {
		priorityIDs[p.ID] = true
	}
	expectedPriorities := []string{"high", "medium", "low"}
	for _, id := range expectedPriorities {
		if !priorityIDs[id] {
			t.Fatalf("priorities legend missing %q", id)
		}
	}

	// Canonical status values in legend.
	statusIDs := map[string]bool{}
	for _, s := range body.Statuses {
		statusIDs[s.ID] = true
	}
	expectedStatuses := []string{"backlog", "ready", "assigned"}
	for _, id := range expectedStatuses {
		if !statusIDs[id] {
			t.Fatalf("statuses legend missing %q", id)
		}
	}

	// Validate all priority and status values across tasks are canonical.
	for _, task := range body.Tasks {
		if !priorityIDs[task.Priority] {
			t.Fatalf("task %s has non-canonical priority %q", task.ID, task.Priority)
		}
		if !statusIDs[task.Status] {
			t.Fatalf("task %s has non-canonical status %q", task.ID, task.Status)
		}
		if task.ID == "" {
			t.Fatal("task has empty id")
		}
		if task.Title == "" {
			t.Fatalf("task %s has empty title", task.ID)
		}
		if task.Room == "" {
			t.Fatalf("task %s has empty room", task.ID)
		}
		if task.PeopleNeeded < 1 {
			t.Fatalf("task %s has peopleNeeded=%d, expected >= 1", task.ID, task.PeopleNeeded)
		}
	}

	// Verify at least 10 tasks.
	if len(body.Tasks) < 10 {
		t.Fatalf("expected at least 10 tasks, got %d", len(body.Tasks))
	}

	// Verify OpenAPI includes the new path.
	openapiBytes, err := json.Marshal(api.OpenAPI())
	if err != nil {
		t.Fatalf("failed to marshal OpenAPI: %v", err)
	}
	var openapi map[string]interface{}
	if err := json.Unmarshal(openapiBytes, &openapi); err != nil {
		t.Fatalf("failed to unmarshal OpenAPI: %v", err)
	}
	paths, ok := openapi["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("OpenAPI spec missing paths")
	}
	if _, exists := paths["/api/tasks/backlog"]; !exists {
		t.Fatal("OpenAPI paths missing /api/tasks/backlog")
	}
}

func TestDailyScheduleHappyPath(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body DailyScheduleBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	// Default-window range assertions.
	if body.Range.StartDate != "2026-07-05" {
		t.Fatalf("expected range.startDate '2026-07-05', got %q", body.Range.StartDate)
	}
	if body.Range.EndDate != "2026-07-08" {
		t.Fatalf("expected range.endDate '2026-07-08', got %q", body.Range.EndDate)
	}
	if body.Range.Days != 4 {
		t.Fatalf("expected range.days=4, got %d", body.Range.Days)
	}

	// One day per requested date.
	if len(body.Days) != body.Range.Days {
		t.Fatalf("expected %d day entries, got %d", body.Range.Days, len(body.Days))
	}

	// Day date uniqueness.
	seenDates := map[string]bool{}
	for _, d := range body.Days {
		if d.Date == "" {
			t.Fatal("day.Date is empty")
		}
		if seenDates[d.Date] {
			t.Fatalf("duplicate day date %q", d.Date)
		}
		seenDates[d.Date] = true

		if d.Label == "" {
			t.Fatalf("day %s has empty label", d.Date)
		}

		// Validate task cards.
		for _, task := range d.Tasks {
			if task.ID == "" {
				t.Fatalf("task on day %s has empty id", d.Date)
			}
			if task.Title == "" {
				t.Fatalf("task %s on day %s has empty title", task.ID, d.Date)
			}

			// Validate priority enum.
			switch task.Priority {
			case "high", "medium", "low":
			default:
				t.Fatalf("task %s on day %s has invalid priority %q", task.ID, d.Date, task.Priority)
			}

			// Validate staffingStatus enum.
			switch task.StaffingStatus {
			case "fullyStaffed", "underStaffed":
			default:
				t.Fatalf("task %s on day %s has invalid staffingStatus %q", task.ID, d.Date, task.StaffingStatus)
			}

			// Derived-field invariants.
			if task.AssignedCount != len(task.AssignedPeople) {
				t.Fatalf("task %s: assignedCount=%d != len(assignedPeople)=%d", task.ID, task.AssignedCount, len(task.AssignedPeople))
			}
			if task.AssignedCount > task.PeopleNeeded {
				t.Fatalf("task %s: assignedCount=%d > peopleNeeded=%d", task.ID, task.AssignedCount, task.PeopleNeeded)
			}
			if task.AssignedCount == task.PeopleNeeded && task.StaffingStatus != "fullyStaffed" {
				t.Fatalf("task %s: assigned=%d needed=%d but staffingStatus=%q (expected fullyStaffed)", task.ID, task.AssignedCount, task.PeopleNeeded, task.StaffingStatus)
			}
			if task.AssignedCount < task.PeopleNeeded && task.StaffingStatus != "underStaffed" {
				t.Fatalf("task %s: assigned=%d needed=%d but staffingStatus=%q (expected underStaffed)", task.ID, task.AssignedCount, task.PeopleNeeded, task.StaffingStatus)
			}

			if task.PeopleNeeded < 1 {
				t.Fatalf("task %s: peopleNeeded=%d (must be >= 1)", task.ID, task.PeopleNeeded)
			}

			// Validate assigned people reuse seed identities.
			for _, ap := range task.AssignedPeople {
				if ap.ID == "" {
					t.Fatalf("task %s has assigned person with empty id", task.ID)
				}
				if ap.Name == "" {
					t.Fatalf("task %s assigned person %s has empty name", task.ID, ap.ID)
				}
				if ap.Initials == "" {
					t.Fatalf("task %s assigned person %s has empty initials", task.ID, ap.ID)
				}
				// Verify the person id is from seedPeople.
				found := false
				for _, sp := range seedPeople {
					if sp.Id == ap.ID {
						found = true
						if sp.Name != ap.Name {
							t.Fatalf("task %s assigned person %s name mismatch: seed=%q response=%q", task.ID, ap.ID, sp.Name, ap.Name)
						}
						if sp.Initials != ap.Initials {
							t.Fatalf("task %s assigned person %s initials mismatch: seed=%q response=%q", task.ID, ap.ID, sp.Initials, ap.Initials)
						}
						break
					}
				}
				if !found {
					t.Fatalf("task %s assigned person %s not found in seedPeople", task.ID, ap.ID)
				}
			}
		}
	}

	// Default window seed variety: priorities and staffing states.
	allPriorities := map[string]bool{}
	allStaffing := map[string]bool{}
	hasFullyStaffed2x2 := false
	hasFullyStaffed1x1 := false
	hasUnderstaffed := false
	for _, d := range body.Days {
		for _, task := range d.Tasks {
			allPriorities[task.Priority] = true
			allStaffing[task.StaffingStatus] = true
			if task.PeopleNeeded == 2 && task.AssignedCount == 2 && task.StaffingStatus == "fullyStaffed" {
				hasFullyStaffed2x2 = true
			}
			if task.PeopleNeeded == 1 && task.AssignedCount == 1 && task.StaffingStatus == "fullyStaffed" {
				hasFullyStaffed1x1 = true
			}
			if task.StaffingStatus == "underStaffed" {
				hasUnderstaffed = true
			}
		}
	}

	if !allPriorities["high"] || !allPriorities["medium"] || !allPriorities["low"] {
		t.Fatalf("expected high/medium/low priorities across default window, got %v", allPriorities)
	}
	if !hasFullyStaffed2x2 {
		t.Fatal("expected at least one fully staffed 2/2 task card")
	}
	if !hasFullyStaffed1x1 {
		t.Fatal("expected at least one fully staffed 1/1 task card")
	}
	if !hasUnderstaffed {
		t.Fatal("expected at least one understaffed task card")
	}

	// Verify availablePeopleCount matches seedPeople availability.
	for _, d := range body.Days {
		// Parse the date to compute day index relative to start.
		date, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			t.Fatalf("failed to parse day date %q: %v", d.Date, err)
		}
		startDate, _ := time.Parse("2006-01-02", body.Range.StartDate)
		dayOffset := int(date.Sub(startDate).Hours() / 24)
		expectedAvailable := countAvailableForDay(dayOffset)
		if d.AvailablePeopleCount != expectedAvailable {
			t.Fatalf("day %s: availablePeopleCount=%d != expected=%d", d.Date, d.AvailablePeopleCount, expectedAvailable)
		}
	}
}

func TestDailyScheduleExplicitParams(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule?start=2026-07-10&days=3", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var body DailyScheduleBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if body.Range.StartDate != "2026-07-10" {
		t.Fatalf("expected startDate '2026-07-10', got %q", body.Range.StartDate)
	}
	if body.Range.EndDate != "2026-07-12" {
		t.Fatalf("expected endDate '2026-07-12', got %q", body.Range.EndDate)
	}
	if body.Range.Days != 3 {
		t.Fatalf("expected days=3, got %d", body.Range.Days)
	}
	if len(body.Days) != 3 {
		t.Fatalf("expected 3 day entries, got %d", len(body.Days))
	}
}

func TestDailyScheduleMalformedStart(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule?start=2026-13-99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDailyScheduleOpenAPIInclusion(t *testing.T) {
	_, api := newTestAPI(newMockStore())

	openAPIBytes, err := json.Marshal(api.OpenAPI())
	if err != nil {
		t.Fatalf("failed to marshal OpenAPI: %v", err)
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(openAPIBytes, &spec); err != nil {
		t.Fatalf("failed to unmarshal OpenAPI spec: %v", err)
	}

	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("OpenAPI spec missing 'paths' key")
	}

	if _, exists := paths["/api/dashboard/daily-schedule"]; !exists {
		t.Fatal("OpenAPI spec does not include /api/dashboard/daily-schedule")
	}
}

func TestDailyScheduleDeterministic(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	// Two identical requests should produce the same response.
	req1 := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule?start=2026-07-05&days=4", nil)
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule?start=2026-07-05&days=4", nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	if rec1.Body.String() != rec2.Body.String() {
		t.Fatal("identical requests produced different responses")
	}
}

// ---------- Failure-mode tests ----------

// failingStore always returns an error for every Store method.
type failingStore struct{}

func (f *failingStore) GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*DashboardBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) GetTaskBacklog(ctx context.Context) (*TaskBacklogBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*DailyScheduleBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) ListRooms(ctx context.Context) ([]Room, error) {
	return nil, errTestFailure
}

func (f *failingStore) CreateRoom(ctx context.Context, input CreateRoomInput) (*Room, error) {
	return nil, errTestFailure
}

func (f *failingStore) UpdateRoom(ctx context.Context, id string, input UpdateRoomInput) (*Room, error) {
	return nil, errTestFailure
}

func (f *failingStore) DeleteRoom(ctx context.Context, id string) error {
	return errTestFailure
}

// errTestFailure is a sentinel error used by failingStore.
var errTestFailure = errors.New("test-induced store failure")

func TestTaskBacklogStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/backlog", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDailyScheduleStoreGetPlanningWindowFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	// No start param -> handler calls store.GetPlanningWindow which fails.
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDailyScheduleStoreGetDailyScheduleFailure(t *testing.T) {
	// Create a store that returns a planning window but fails on GetDailySchedule.
	store := &partialFailingStore{}
	router, _ := newTestAPI(store)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule?start=2026-07-05&days=4", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

// partialFailingStore succeeds on GetPlanningWindow but fails on GetDailySchedule.
type partialFailingStore struct{}

func (f *partialFailingStore) GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error) {
	return &PlanningWindowBody{
		StartDate: "2026-07-05",
		EndDate:   "2026-08-13",
		Days:      40,
	}, nil
}

func (f *partialFailingStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*DashboardBody, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) GetTaskBacklog(ctx context.Context) (*TaskBacklogBody, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*DailyScheduleBody, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) ListRooms(ctx context.Context) ([]Room, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) CreateRoom(ctx context.Context, input CreateRoomInput) (*Room, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) UpdateRoom(ctx context.Context, id string, input UpdateRoomInput) (*Room, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) DeleteRoom(ctx context.Context, id string) error {
	return errTestFailure
}

// ---------- Room CRUD tests ----------

func TestListRooms(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body struct {
		Rooms []Room `json:"rooms"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	if len(body.Rooms) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(body.Rooms))
	}

	for _, room := range body.Rooms {
		if room.ID == "" {
			t.Fatal("room has empty id")
		}
		if room.Name == "" {
			t.Fatalf("room %s has empty name", room.ID)
		}
		if room.Type != "room" && room.Type != "area" {
			t.Fatalf("room %s has invalid type %q", room.ID, room.Type)
		}
	}
}

func TestCreateRoom(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	bodyJSON := `{"name":"Basement","type":"area"}`
	req := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var room Room
	if err := json.Unmarshal(rec.Body.Bytes(), &room); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	if room.Name != "Basement" {
		t.Fatalf("expected name 'Basement', got %q", room.Name)
	}
	if room.Type != "area" {
		t.Fatalf("expected type 'area', got %q", room.Type)
	}
	if room.ID == "" {
		t.Fatal("room has empty id")
	}
}

func TestCreateRoomInvalidType(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	bodyJSON := `{"name":"Lobby","type":"lobby"}`
	req := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateRoom(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	bodyJSON := `{"name":"Updated Kitchen","type":"room"}`
	req := httptest.NewRequest(http.MethodPut, "/api/rooms/room-1", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var room Room
	if err := json.Unmarshal(rec.Body.Bytes(), &room); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	if room.ID != "room-1" {
		t.Fatalf("expected id 'room-1', got %q", room.ID)
	}
	if room.Name != "Updated Kitchen" {
		t.Fatalf("expected name 'Updated Kitchen', got %q", room.Name)
	}
}

func TestUpdateRoomNotFound(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	bodyJSON := `{"name":"Ghost","type":"room"}`
	req := httptest.NewRequest(http.MethodPut, "/api/rooms/nonexistent", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteRoom(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodDelete, "/api/rooms/room-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteRoomNotFound(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodDelete, "/api/rooms/nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestListRoomsStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	req := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateRoomStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	bodyJSON := `{"name":"Basement","type":"area"}`
	req := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestRoomsOpenAPIInclusion(t *testing.T) {
	_, api := newTestAPI(newMockStore())

	openAPIBytes, err := json.Marshal(api.OpenAPI())
	if err != nil {
		t.Fatalf("failed to marshal OpenAPI: %v", err)
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(openAPIBytes, &spec); err != nil {
		t.Fatalf("failed to unmarshal OpenAPI spec: %v", err)
	}

	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("OpenAPI spec missing 'paths' key")
	}

	if _, exists := paths["/api/rooms"]; !exists {
		t.Fatal("OpenAPI spec does not include /api/rooms")
	}
	if _, exists := paths["/api/rooms/{id}"]; !exists {
		t.Fatal("OpenAPI spec does not include /api/rooms/{id}")
	}
}
