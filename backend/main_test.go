package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	registerTasksBacklog(api)

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
