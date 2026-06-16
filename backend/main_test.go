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

	backendapi "github.com/user/the-great-migration/backend/api"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func newTestAPI(store backendapi.Store) (chi.Router, huma.API) {
	router := chi.NewMux()
	config := huma.DefaultConfig("Test API", "1.0.0")
	api := humachi.New(router, config)

	backendapi.RegisterAll(api, store)

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

	var body backendapi.DashboardBody
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

	var body backendapi.PlanningWindowBody
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

func TestUpdatePlanningWindowSuccess(t *testing.T) {
	store := newMockStore()
	router, _ := newTestAPI(store)

	bodyJSON := `{"startDate": "2026-07-10", "endDate": "2026-07-20"}`
	req := httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body backendapi.PlanningWindowBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	if body.StartDate != "2026-07-10" {
		t.Fatalf("expected startDate '2026-07-10', got %q", body.StartDate)
	}
	if body.EndDate != "2026-07-20" {
		t.Fatalf("expected endDate '2026-07-20', got %q", body.EndDate)
	}
	if body.Days != 11 {
		t.Fatalf("expected days=11, got %d", body.Days)
	}

	// Verify store was updated.
	pw, err := store.GetPlanningWindow(context.Background())
	if err != nil {
		t.Fatalf("failed to get planning window: %v", err)
	}
	if pw.StartDate != "2026-07-10" || pw.EndDate != "2026-07-20" || pw.Days != 11 {
		t.Fatalf("store was not updated: startDate=%q, endDate=%q, days=%d", pw.StartDate, pw.EndDate, pw.Days)
	}
}

func TestUpdatePlanningWindowValidationFailure(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	// endDate < startDate.
	bodyJSON := `{"startDate": "2026-08-01", "endDate": "2026-07-01"}`
	req := httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 422 {
		t.Fatalf("expected status 422, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	// Malformed dates.
	bodyJSON = `{"startDate": "not-a-date", "endDate": "2026-07-01"}`
	req = httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 422 {
		t.Fatalf("expected status 422, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdatePlanningWindowMaxRange(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	// Range exactly at the max (365 days) should succeed.
	bodyJSON := `{"startDate": "2026-01-01", "endDate": "2026-12-31"}`
	req := httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for exactly 365 days, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	// Range exceeding max (366 days) should be rejected.
	bodyJSON = `{"startDate": "2026-01-01", "endDate": "2027-01-01"}`
	req = httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 422 {
		t.Fatalf("expected status 422 for range exceeding max, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdatePlanningWindowStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	bodyJSON := `{"startDate": "2026-07-10", "endDate": "2026-07-20"}`
	req := httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
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

	var body backendapi.TaskBacklogBody
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

	var body backendapi.DailyScheduleBody
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

	var body backendapi.DailyScheduleBody
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

func (f *failingStore) GetPlanningWindow(ctx context.Context) (*backendapi.PlanningWindowBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) UpdatePlanningWindow(ctx context.Context, startDate, endDate time.Time) (*backendapi.PlanningWindowBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*backendapi.DashboardBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) GetTaskBacklog(ctx context.Context) (*backendapi.TaskBacklogBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*backendapi.DailyScheduleBody, error) {
	return nil, errTestFailure
}

func (f *failingStore) CreatePerson(ctx context.Context, id, name, initials string) error {
	return errTestFailure
}

func (f *failingStore) UpdatePerson(ctx context.Context, id, name, initials string) error {
	return errTestFailure
}

func (f *failingStore) DeletePerson(ctx context.Context, id string) error {
	return errTestFailure
}

func (f *failingStore) PersonExists(ctx context.Context, id string) (bool, error) {
	return false, errTestFailure
}

func (f *failingStore) PersonHasReferences(ctx context.Context, id string) (bool, error) {
	return false, errTestFailure
}

func (f *failingStore) UpsertAvailability(ctx context.Context, personID string, date pgtype.Date, status string) error {
	return errTestFailure
}

func (f *failingStore) DeleteAvailability(ctx context.Context, personID string, date pgtype.Date) error {
	return errTestFailure
}

func (f *failingStore) ListRooms(ctx context.Context) ([]backendapi.Room, error) {
	return nil, errTestFailure
}

func (f *failingStore) CreateRoom(ctx context.Context, input backendapi.CreateRoomInput) (*backendapi.Room, error) {
	return nil, errTestFailure
}

func (f *failingStore) UpdateRoom(ctx context.Context, id string, input backendapi.UpdateRoomInput) (*backendapi.Room, error) {
	return nil, errTestFailure
}

func (f *failingStore) DeleteRoom(ctx context.Context, id string) error {
	return errTestFailure
}

func (f *failingStore) CreateTask(ctx context.Context, input backendapi.CreateTaskInput) (*backendapi.TaskRow, error) {
	return nil, errTestFailure
}

func (f *failingStore) UpdateTask(ctx context.Context, id string, input backendapi.UpdateTaskInput) (*backendapi.TaskRow, error) {
	return nil, errTestFailure
}

func (f *failingStore) DeleteTask(ctx context.Context, id string) error {
	return errTestFailure
}

func (f *failingStore) CreateScheduleCard(ctx context.Context, input backendapi.CreateScheduleCardInput) (*backendapi.TaskCard, error) {
	return nil, errTestFailure
}

func (f *failingStore) UpdateScheduleCard(ctx context.Context, id string, input backendapi.CreateScheduleCardInput) (*backendapi.TaskCard, error) {
	return nil, errTestFailure
}

func (f *failingStore) DeleteScheduleCard(ctx context.Context, id string) error {
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

func (f *partialFailingStore) GetPlanningWindow(ctx context.Context) (*backendapi.PlanningWindowBody, error) {
	return &backendapi.PlanningWindowBody{
		StartDate: "2026-07-05",
		EndDate:   "2026-08-13",
		Days:      40,
	}, nil
}

func (f *partialFailingStore) UpdatePlanningWindow(ctx context.Context, startDate, endDate time.Time) (*backendapi.PlanningWindowBody, error) {
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	return &backendapi.PlanningWindowBody{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		Days:      days,
	}, nil
}

func (f *partialFailingStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*backendapi.DashboardBody, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) GetTaskBacklog(ctx context.Context) (*backendapi.TaskBacklogBody, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*backendapi.DailyScheduleBody, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) CreatePerson(ctx context.Context, id, name, initials string) error {
	return errTestFailure
}

func (f *partialFailingStore) UpdatePerson(ctx context.Context, id, name, initials string) error {
	return errTestFailure
}

func (f *partialFailingStore) DeletePerson(ctx context.Context, id string) error {
	return errTestFailure
}

func (f *partialFailingStore) PersonExists(ctx context.Context, id string) (bool, error) {
	return false, errTestFailure
}

func (f *partialFailingStore) PersonHasReferences(ctx context.Context, id string) (bool, error) {
	return false, errTestFailure
}

func (f *partialFailingStore) UpsertAvailability(ctx context.Context, personID string, date pgtype.Date, status string) error {
	return errTestFailure
}

func (f *partialFailingStore) DeleteAvailability(ctx context.Context, personID string, date pgtype.Date) error {
	return errTestFailure
}

func (f *partialFailingStore) ListRooms(ctx context.Context) ([]backendapi.Room, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) CreateRoom(ctx context.Context, input backendapi.CreateRoomInput) (*backendapi.Room, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) UpdateRoom(ctx context.Context, id string, input backendapi.UpdateRoomInput) (*backendapi.Room, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) DeleteRoom(ctx context.Context, id string) error {
	return errTestFailure
}

func (f *partialFailingStore) CreateTask(ctx context.Context, input backendapi.CreateTaskInput) (*backendapi.TaskRow, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) UpdateTask(ctx context.Context, id string, input backendapi.UpdateTaskInput) (*backendapi.TaskRow, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) DeleteTask(ctx context.Context, id string) error {
	return errTestFailure
}

func (f *partialFailingStore) CreateScheduleCard(ctx context.Context, input backendapi.CreateScheduleCardInput) (*backendapi.TaskCard, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) UpdateScheduleCard(ctx context.Context, id string, input backendapi.CreateScheduleCardInput) (*backendapi.TaskCard, error) {
	return nil, errTestFailure
}

func (f *partialFailingStore) DeleteScheduleCard(ctx context.Context, id string) error {
	return errTestFailure
}

// ---------- Task CRUD tests ----------

func TestCreateTask(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	body := `{"title":"Pack kitchen boxes","priority":"high","peopleNeeded":3,"room":"Kitchen","status":"backlog","assignedTo":["p1","p2"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var resp backendapi.TaskRow
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Title != "Pack kitchen boxes" {
		t.Fatalf("expected title 'Pack kitchen boxes', got %q", resp.Title)
	}
	if resp.Priority != "high" {
		t.Fatalf("expected priority 'high', got %q", resp.Priority)
	}
	if resp.PeopleNeeded != 3 {
		t.Fatalf("expected peopleNeeded 3, got %d", resp.PeopleNeeded)
	}
	if resp.Room != "Kitchen" {
		t.Fatalf("expected room 'Kitchen', got %q", resp.Room)
	}
	if resp.Status != "backlog" {
		t.Fatalf("expected status 'backlog', got %q", resp.Status)
	}
	if len(resp.AssignedTo) != 2 || resp.AssignedTo[0] != "p1" || resp.AssignedTo[1] != "p2" {
		t.Fatalf("expected assignedTo [p1 p2], got %v", resp.AssignedTo)
	}
	if resp.ID == "" || resp.ID[:5] != "task-" {
		t.Fatalf("expected task ID prefixed 'task-', got %q", resp.ID)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "empty title",
			body:           `{"title":"","priority":"medium","peopleNeeded":2,"room":"Kitchen","status":"backlog","assignedTo":[]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid priority",
			body:           `{"title":"Task","priority":"urgent","peopleNeeded":2,"room":"Kitchen","status":"backlog","assignedTo":[]}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "invalid status",
			body:           `{"title":"Task","priority":"medium","peopleNeeded":2,"room":"Kitchen","status":"done","assignedTo":[]}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "peopleNeeded < 1",
			body:           `{"title":"Task","priority":"medium","peopleNeeded":0,"room":"Kitchen","status":"backlog","assignedTo":[]}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty room",
			body:           `{"title":"Task","priority":"medium","peopleNeeded":2,"room":"","status":"backlog","assignedTo":[]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unknown assigned person",
			body:           `{"title":"Task","priority":"medium","peopleNeeded":2,"room":"Kitchen","status":"backlog","assignedTo":["p99"]}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d\nbody: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	store := newMockStore()
	router, _ := newTestAPI(store)

	// Update an existing task (task-1 from seed data).
	body := `{"title":"Updated kitchen task","priority":"low","peopleNeeded":1,"room":"Kitchen","status":"ready","assignedTo":["p3"]}`
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/task-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var resp backendapi.TaskRow
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Title != "Updated kitchen task" {
		t.Fatalf("expected title 'Updated kitchen task', got %q", resp.Title)
	}
	if resp.Priority != "low" {
		t.Fatalf("expected priority 'low', got %q", resp.Priority)
	}
	if resp.PeopleNeeded != 1 {
		t.Fatalf("expected peopleNeeded 1, got %d", resp.PeopleNeeded)
	}
	if resp.Status != "ready" {
		t.Fatalf("expected status 'ready', got %q", resp.Status)
	}
	if len(resp.AssignedTo) != 1 || resp.AssignedTo[0] != "p3" {
		t.Fatalf("expected assignedTo [p3], got %v", resp.AssignedTo)
	}

	// Verify the store reflects the update.
	task := store.tasks["task-1"]
	if task.Title != "Updated kitchen task" {
		t.Fatalf("store task title not updated: got %q", task.Title)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	body := `{"title":"Nobody","priority":"low","peopleNeeded":1,"room":"Nowhere","status":"backlog","assignedTo":[]}`
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/task-nonexistent", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteTask(t *testing.T) {
	store := newMockStore()
	router, _ := newTestAPI(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/task-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	// Verify the task was removed from the store.
	if _, exists := store.tasks["task-1"]; exists {
		t.Fatal("task should have been deleted")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/task-nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestTaskWriteReflectsInBacklog(t *testing.T) {
	store := newMockStore()
	router, _ := newTestAPI(store)

	// Create a task.
	createBody := `{"title":"New reflected task","priority":"high","peopleNeeded":2,"room":"Garage","status":"backlog","assignedTo":["p1"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create failed: status %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var created backendapi.TaskRow
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to unmarshal created task: %v", err)
	}

	// GET the backlog and verify the new task appears.
	getReq := httptest.NewRequest(http.MethodGet, "/api/tasks/backlog", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("backlog GET failed: status %d\nbody: %s", getRec.Code, getRec.Body.String())
	}

	var backlog backendapi.TaskBacklogBody
	if err := json.Unmarshal(getRec.Body.Bytes(), &backlog); err != nil {
		t.Fatalf("failed to unmarshal backlog: %v", err)
	}

	found := false
	for _, task := range backlog.Tasks {
		if task.ID == created.ID {
			found = true
			if task.Title != "New reflected task" {
				t.Fatalf("reflected task has wrong title: %q", task.Title)
			}
			if task.Priority != "high" {
				t.Fatalf("reflected task has wrong priority: %q", task.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatalf("created task %s not found in backlog response", created.ID)
	}
}

// ---------- People CRUD test store ----------

// peopleTestStore wraps a mockStore with in-memory CRUD support for testing the write surface.
type peopleTestStore struct {
	*mockStore
	people                 map[string]testPerson
	availability           map[string]map[string]string // personID -> date -> status
	deleteShouldFailWithFK bool
}

type testPerson struct {
	Name     string
	Initials string
}

func newPeopleTestStore() *peopleTestStore {
	return &peopleTestStore{
		mockStore:    newMockStore(),
		people:       make(map[string]testPerson),
		availability: make(map[string]map[string]string),
	}
}

func (s *peopleTestStore) CreatePerson(ctx context.Context, id, name, initials string) error {
	if _, exists := s.people[id]; exists {
		return errors.New("duplicate key")
	}
	s.people[id] = testPerson{Name: name, Initials: initials}
	return nil
}

func (s *peopleTestStore) UpdatePerson(ctx context.Context, id, name, initials string) error {
	if _, exists := s.people[id]; !exists {
		return errors.New("not found")
	}
	s.people[id] = testPerson{Name: name, Initials: initials}
	return nil
}

func (s *peopleTestStore) DeletePerson(ctx context.Context, id string) error {
	if s.deleteShouldFailWithFK {
		return &pgconn.PgError{Code: "23503", Message: "violates foreign key constraint"}
	}
	delete(s.people, id)
	return nil
}

func (s *peopleTestStore) PersonExists(ctx context.Context, id string) (bool, error) {
	_, exists := s.people[id]
	return exists, nil
}

func (s *peopleTestStore) PersonHasReferences(ctx context.Context, id string) (bool, error) {
	// Simulate: p1, p2, p3 are referenced by backlog/schedule assignments.
	return id == "p1" || id == "p2" || id == "p3", nil
}

func (s *peopleTestStore) UpsertAvailability(ctx context.Context, personID string, date pgtype.Date, status string) error {
	if s.availability[personID] == nil {
		s.availability[personID] = make(map[string]string)
	}
	dateStr := date.Time.Format("2006-01-02")
	s.availability[personID][dateStr] = status
	return nil
}

func (s *peopleTestStore) DeleteAvailability(ctx context.Context, personID string, date pgtype.Date) error {
	if m, ok := s.availability[personID]; ok {
		delete(m, date.Time.Format("2006-01-02"))
	}
	return nil
}

// Override GetPeopleAvailability to use the CRUD-backed data.
func (s *peopleTestStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*backendapi.DashboardBody, error) {
	// Build from in-memory people + availability.
	people := make([]backendapi.Person, 0, len(s.people))
	endDate := startDate.AddDate(0, 0, days-1)
	selectedDate := startDate.Format("2006-01-02")

	for id, tp := range s.people {
		avail := make([]backendapi.AvailabilityEntry, days)
		for d := 0; d < days; d++ {
			date := startDate.AddDate(0, 0, d)
			dateStr := date.Format("2006-01-02")
			status := "off" // default
			if m, ok := s.availability[id]; ok {
				if st, ok := m[dateStr]; ok {
					status = st
				}
			}
			avail[d] = backendapi.AvailabilityEntry{
				Date:   dateStr,
				Status: status,
			}
		}
		people = append(people, backendapi.Person{
			ID:           id,
			Name:         tp.Name,
			Initials:     tp.Initials,
			Availability: avail,
		})
	}

	// Compute availableToday: count people who are "available" on the selected date.
	availableToday := 0
	for _, p := range people {
		for _, e := range p.Availability {
			if e.Date == selectedDate && e.Status == "available" {
				availableToday++
				break
			}
		}
	}

	return &backendapi.DashboardBody{
		Range: backendapi.Range{
			StartDate:    startDate.Format("2006-01-02"),
			EndDate:      endDate.Format("2006-01-02"),
			Days:         days,
			SelectedDate: selectedDate,
		},
		Summary: backendapi.Summary{
			AvailableToday: availableToday,
			TotalPeople:    len(people),
		},
		People:   people,
		Statuses: backendapi.StatusLegendData,
	}, nil
}

// ---------- People CRUD tests ----------

func TestCreatePerson(t *testing.T) {
	store := newPeopleTestStore()
	router, _ := newTestAPI(store)

	body := `{"id":"p9","name":"Test User","initials":"TU"}`
	req := httptest.NewRequest(http.MethodPost, "/api/people", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var resp backendapi.Person
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.ID != "p9" {
		t.Fatalf("expected id p9, got %q", resp.ID)
	}
	if resp.Name != "Test User" {
		t.Fatalf("expected name 'Test User', got %q", resp.Name)
	}
	if resp.Initials != "TU" {
		t.Fatalf("expected initials 'TU', got %q", resp.Initials)
	}
}

func TestCreatePersonDuplicate(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p9", "Existing", "EX"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	body := `{"id":"p9","name":"Test User","initials":"TU"}`
	req := httptest.NewRequest(http.MethodPost, "/api/people", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409 for duplicate, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestCreatePersonMissingFields(t *testing.T) {
	router, _ := newTestAPI(newPeopleTestStore())

	body := `{"id":"","name":"","initials":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/people", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422 for empty fields, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdatePerson(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p9", "Original", "OR"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	body := `{"name":"Updated Name","initials":"UN"}`
	req := httptest.NewRequest(http.MethodPut, "/api/people/p9", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var resp backendapi.Person
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Name != "Updated Name" {
		t.Fatalf("expected name 'Updated Name', got %q", resp.Name)
	}
	if resp.Initials != "UN" {
		t.Fatalf("expected initials 'UN', got %q", resp.Initials)
	}
}

func TestUpdatePersonNotFound(t *testing.T) {
	router, _ := newTestAPI(newPeopleTestStore())

	body := `{"name":"Nobody","initials":"NB"}`
	req := httptest.NewRequest(http.MethodPut, "/api/people/nonexistent", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeletePerson(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p99", "To Delete", "TD"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/people/p99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	exists, _ := store.PersonExists(context.Background(), "p99")
	if exists {
		t.Fatal("person should have been deleted")
	}
}

func TestDeletePersonNotFound(t *testing.T) {
	router, _ := newTestAPI(newPeopleTestStore())

	req := httptest.NewRequest(http.MethodDelete, "/api/people/nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for nonexistent person, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeletePersonConflict(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p1", "Sophia Chen", "SC"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/people/p1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeletePersonForeignKeyViolation(t *testing.T) {
	store := newPeopleTestStore()
	// p9 has no references (PersonHasReferences returns false), but a concurrent
	// request could insert a reference between the check and delete.
	if err := store.CreatePerson(context.Background(), "p9", "Raced Person", "RP"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}
	store.deleteShouldFailWithFK = true

	router, _ := newTestAPI(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/people/p9", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409 for FK violation, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpsertAvailability(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p9", "Test User", "TU"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	body := `{"status":"available"}`
	req := httptest.NewRequest(http.MethodPut, "/api/people/p9/availability/2026-07-10", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
	if got := store.availability["p9"]["2026-07-10"]; got != "available" {
		t.Fatalf("expected persisted availability 'available', got %q", got)
	}
}

func TestUpsertAvailabilityInvalidStatus(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p9", "Test User", "TU"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	body := `{"status":"unknown"}`
	req := httptest.NewRequest(http.MethodPut, "/api/people/p9/availability/2026-07-10", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid status, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpsertAvailabilityPersonNotFound(t *testing.T) {
	router, _ := newTestAPI(newPeopleTestStore())

	body := `{"status":"available"}`
	req := httptest.NewRequest(http.MethodPut, "/api/people/nonexistent/availability/2026-07-10", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpsertAvailabilityOutOfWindow(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p9", "Test User", "TU"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	body := `{"status":"available"}`
	req := httptest.NewRequest(http.MethodPut, "/api/people/p9/availability/2025-01-01", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for out-of-window date, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpsertAvailabilityMalformedDate(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p9", "Test User", "TU"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}

	router, _ := newTestAPI(store)

	body := `{"status":"available"}`
	req := httptest.NewRequest(http.MethodPut, "/api/people/p9/availability/not-a-date", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for malformed date, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteAvailability(t *testing.T) {
	store := newPeopleTestStore()
	if err := store.CreatePerson(context.Background(), "p9", "Test User", "TU"); err != nil {
		t.Fatalf("failed to seed person: %v", err)
	}
	date := pgtype.Date{Time: time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC), Valid: true}
	if err := store.UpsertAvailability(context.Background(), "p9", date, "available"); err != nil {
		t.Fatalf("failed to seed availability: %v", err)
	}

	router, _ := newTestAPI(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/people/p9/availability/2026-07-10", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
	if _, ok := store.availability["p9"]["2026-07-10"]; ok {
		t.Fatal("availability should have been deleted")
	}
}

func TestDeleteAvailabilityPersonNotFound(t *testing.T) {
	router, _ := newTestAPI(newPeopleTestStore())

	req := httptest.NewRequest(http.MethodDelete, "/api/people/nonexistent/availability/2026-07-10", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
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

	var body struct {
		Rooms []backendapi.Room `json:"rooms"`
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

	var room backendapi.Room
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

	var room backendapi.Room
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

func TestUpdateRoomStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	bodyJSON := `{"name":"Updated","type":"room"}`
	req := httptest.NewRequest(http.MethodPut, "/api/rooms/room-1", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteRoomStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	req := httptest.NewRequest(http.MethodDelete, "/api/rooms/room-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestOpenAPIIncludesMergedEndpoints(t *testing.T) {
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

	expectedPaths := []string{
		"/api/hello",
		"/api/dashboard/people-availability",
		"/api/dashboard/daily-schedule",
		"/api/planning-window",
		"/api/tasks",
		"/api/tasks/backlog",
		"/api/tasks/{id}",
		"/api/people",
		"/api/people/{id}",
		"/api/people/{id}/availability/{date}",
		"/api/rooms",
		"/api/rooms/{id}",
		"/api/schedule/cards",
		"/api/schedule/cards/{id}",
	}
	for _, path := range expectedPaths {
		if _, exists := paths[path]; !exists {
			t.Fatalf("OpenAPI spec does not include %s", path)
		}
	}
}

// ---------- Schedule-card CRUD tests ----------

func TestCreateScheduleCard(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	body := `{"title":"Kitchen painting","priority":"high","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"2026-07-05","assignedTo":["p1","p2"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/schedule/cards", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var resp backendapi.TaskCard
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Title != "Kitchen painting" {
		t.Fatalf("expected title 'Kitchen painting', got %q", resp.Title)
	}
	if resp.Priority != "high" {
		t.Fatalf("expected priority 'high', got %q", resp.Priority)
	}
	if resp.PeopleNeeded != 2 {
		t.Fatalf("expected peopleNeeded 2, got %d", resp.PeopleNeeded)
	}
	if resp.RoomArea != "Kitchen" {
		t.Fatalf("expected roomArea 'Kitchen', got %q", resp.RoomArea)
	}
	if len(resp.AssignedPeople) != 2 {
		t.Fatalf("expected 2 assigned people, got %d", len(resp.AssignedPeople))
	}
	if resp.AssignedCount != 2 {
		t.Fatalf("expected assignedCount 2, got %d", resp.AssignedCount)
	}
	if resp.StaffingStatus != "fullyStaffed" {
		t.Fatalf("expected staffingStatus 'fullyStaffed', got %q", resp.StaffingStatus)
	}
	if resp.ID == "" || resp.ID[:6] != "sched-" {
		t.Fatalf("expected id prefixed 'sched-', got %q", resp.ID)
	}
}

func TestCreateScheduleCardValidation(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "empty title",
			body:           `{"title":"","priority":"medium","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"2026-07-05","assignedTo":[]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid priority",
			body:           `{"title":"Task","priority":"urgent","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"2026-07-05","assignedTo":[]}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "peopleNeeded < 1",
			body:           `{"title":"Task","priority":"medium","roomArea":"Kitchen","peopleNeeded":0,"scheduledDate":"2026-07-05","assignedTo":[]}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty roomArea",
			body:           `{"title":"Task","priority":"medium","roomArea":"","peopleNeeded":2,"scheduledDate":"2026-07-05","assignedTo":[]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unknown assigned person",
			body:           `{"title":"Task","priority":"medium","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"2026-07-05","assignedTo":["p99"]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "out-of-window date",
			body:           `{"title":"Task","priority":"medium","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"2025-01-01","assignedTo":[]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "assignedTo exceeds peopleNeeded",
			body:           `{"title":"Task","priority":"medium","roomArea":"Kitchen","peopleNeeded":1,"scheduledDate":"2026-07-05","assignedTo":["p1","p2"]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "malformed date",
			body:           `{"title":"Task","priority":"medium","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"not-a-date","assignedTo":[]}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/schedule/cards", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d\nbody: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestUpdateScheduleCard(t *testing.T) {
	store := newMockStore()

	// Seed a schedule card.
	card, err := store.CreateScheduleCard(context.Background(), backendapi.CreateScheduleCardInput{
		Title:         "Original task",
		Priority:      "low",
		RoomArea:      "Garage",
		PeopleNeeded:  1,
		ScheduledDate: "2026-07-05",
		AssignedTo:    []string{"p1"},
	})
	if err != nil {
		t.Fatalf("failed to seed schedule card: %v", err)
	}

	router, _ := newTestAPI(store)

	body := `{"title":"Updated task","priority":"high","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"2026-07-06","assignedTo":["p2","p3"]}`
	req := httptest.NewRequest(http.MethodPut, "/api/schedule/cards/"+card.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	var resp backendapi.TaskCard
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Title != "Updated task" {
		t.Fatalf("expected title 'Updated task', got %q", resp.Title)
	}
	if resp.Priority != "high" {
		t.Fatalf("expected priority 'high', got %q", resp.Priority)
	}
	if resp.PeopleNeeded != 2 {
		t.Fatalf("expected peopleNeeded 2, got %d", resp.PeopleNeeded)
	}
	if len(resp.AssignedPeople) != 2 {
		t.Fatalf("expected 2 assigned people, got %d", len(resp.AssignedPeople))
	}

	// Verify store reflects the update.
	stored, ok := store.scheduleCards[card.ID]
	if !ok {
		t.Fatal("schedule card should still exist")
	}
	if stored.Title != "Updated task" {
		t.Fatalf("store card title not updated: got %q", stored.Title)
	}
}

func TestUpdateScheduleCardNotFound(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	body := `{"title":"Nowhere","priority":"low","roomArea":"Garage","peopleNeeded":1,"scheduledDate":"2026-07-05","assignedTo":[]}`
	req := httptest.NewRequest(http.MethodPut, "/api/schedule/cards/sched-999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteScheduleCard(t *testing.T) {
	store := newMockStore()

	// Seed a schedule card.
	card, err := store.CreateScheduleCard(context.Background(), backendapi.CreateScheduleCardInput{
		Title:         "To be deleted",
		Priority:      "medium",
		RoomArea:      "Storage",
		PeopleNeeded:  1,
		ScheduledDate: "2026-07-05",
		AssignedTo:    []string{},
	})
	if err != nil {
		t.Fatalf("failed to seed schedule card: %v", err)
	}

	router, _ := newTestAPI(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/schedule/cards/"+card.ID, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	// Verify card was removed.
	if _, exists := store.scheduleCards[card.ID]; exists {
		t.Fatal("schedule card should have been deleted")
	}
}

func TestDeleteScheduleCardNotFound(t *testing.T) {
	router, _ := newTestAPI(newMockStore())

	req := httptest.NewRequest(http.MethodDelete, "/api/schedule/cards/sched-999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestScheduleCardWriteReflectsInDailySchedule(t *testing.T) {
	store := newMockStore()
	router, _ := newTestAPI(store)

	// Create a schedule card on 2026-07-05.
	createBody := `{"title":"Reflected scheduled task","priority":"high","roomArea":"Kitchen","peopleNeeded":2,"scheduledDate":"2026-07-05","assignedTo":["p1","p2"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/schedule/cards", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create failed: status %d\nbody: %s", rec.Code, rec.Body.String())
	}

	// GET the daily schedule and verify the new card appears.
	getReq := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule?start=2026-07-05&days=4", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("daily schedule GET failed: status %d\nbody: %s", getRec.Code, getRec.Body.String())
	}

	var ds backendapi.DailyScheduleBody
	if err := json.Unmarshal(getRec.Body.Bytes(), &ds); err != nil {
		t.Fatalf("failed to unmarshal daily schedule: %v", err)
	}

	// The first day (2026-07-05) should contain the new card.
	if len(ds.Days) == 0 {
		t.Fatal("no days in daily schedule")
	}
	found := false
	for _, task := range ds.Days[0].Tasks {
		if task.Title == "Reflected scheduled task" {
			found = true
			if task.Priority != "high" {
				t.Fatalf("reflected task has wrong priority: %q", task.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatal("created schedule card not found in daily schedule response")
	}
}

func TestCreateScheduleCardStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	body := `{"title":"Test","priority":"medium","roomArea":"Kitchen","peopleNeeded":1,"scheduledDate":"2026-07-05","assignedTo":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/schedule/cards", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateScheduleCardStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	body := `{"title":"Test","priority":"medium","roomArea":"Kitchen","peopleNeeded":1,"scheduledDate":"2026-07-05","assignedTo":[]}`
	req := httptest.NewRequest(http.MethodPut, "/api/schedule/cards/sched-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteScheduleCardStoreFailure(t *testing.T) {
	router, _ := newTestAPI(&failingStore{})

	req := httptest.NewRequest(http.MethodDelete, "/api/schedule/cards/sched-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
}
