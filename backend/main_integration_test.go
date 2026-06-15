//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

// TestDBBackedEndpoints runs endpoint contract tests against a real Postgres database.
// The caller (e.g., scripts/test-integration) is responsible for starting a Postgres
// sidecar via verification_start_postgres_sidecar() and setting DATABASE_URL.
func TestDBBackedEndpoints(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL must be set (run via scripts/test-integration)")
	}

	ctx := context.Background()

	// Create connection pool.
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse DATABASE_URL: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	// Run migrations via goose using the embedded migration files.
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open database for migrations: %v", err)
	}
	defer sqlDB.Close()

	goose.SetBaseFS(migrationsFS)
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	store := NewPgStore(pool)
	router, api := newTestAPI(store)

	// Test planning window endpoint.
	t.Run("PlanningWindow", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/planning-window", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
		}

		var body PlanningWindowBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
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
	})

	// Test dashboard endpoint with default params.
	t.Run("DashboardDefault", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/dashboard/people-availability?start=2026-07-05&days=4", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
		}

		var body DashboardBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if body.Range.StartDate != "2026-07-05" {
			t.Fatalf("expected startDate '2026-07-05', got %q", body.Range.StartDate)
		}
		if body.Range.Days != 4 {
			t.Fatalf("expected days=4, got %d", body.Range.Days)
		}
		if len(body.People) != 8 {
			t.Fatalf("expected 8 people, got %d", len(body.People))
		}
		if len(body.Statuses) != 4 {
			t.Fatalf("expected 4 statuses, got %d", len(body.Statuses))
		}

		// Verify p7 is always busy.
		for _, p := range body.People {
			if p.ID == "p7" {
				for _, e := range p.Availability {
					if e.Status != "busy" {
						t.Fatalf("p7 has status %q on date %s, expected 'busy'", e.Status, e.Date)
					}
				}
			}
		}

		// Verify p8 cycles through off/partial/busy/available.
		for _, p := range body.People {
			if p.ID == "p8" {
				if len(p.Availability) != 4 {
					t.Fatalf("p8 has %d availability entries, expected 4", len(p.Availability))
				}
				expected := []string{"off", "partial", "busy", "available"}
				for i, e := range p.Availability {
					if e.Status != expected[i] {
						t.Fatalf("p8 on day %d: expected %q, got %q", i, expected[i], e.Status)
					}
				}
			}
		}
	})

	// Test tasks backlog from Postgres-backed store.
	t.Run("TasksBacklog", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tasks/backlog", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
		}

		var body TaskBacklogBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		// Verify 11 seeded tasks.
		if len(body.Tasks) != 11 {
			t.Fatalf("expected 11 tasks, got %d", len(body.Tasks))
		}

		// Verify summary consistency.
		if body.Summary.TotalTasks != len(body.Tasks) {
			t.Fatalf("summary.totalTasks=%d != len(tasks)=%d", body.Summary.TotalTasks, len(body.Tasks))
		}

		// Verify canonical vocabularies.
		allPriorities := map[string]bool{}
		allStatuses := map[string]bool{}
		hasEmptyAssigned := false
		hasPartialAssigned := false
		for _, task := range body.Tasks {
			allPriorities[task.Priority] = true
			allStatuses[task.Status] = true
			if len(task.AssignedTo) == 0 {
				hasEmptyAssigned = true
			} else if len(task.AssignedTo) < task.PeopleNeeded {
				hasPartialAssigned = true
			}
		}
		for _, p := range []string{"high", "medium", "low"} {
			if !allPriorities[p] {
				t.Fatalf("task backlog missing priority %q", p)
			}
		}
		for _, s := range []string{"backlog", "ready", "assigned"} {
			if !allStatuses[s] {
				t.Fatalf("task backlog missing status %q", s)
			}
		}
		if !hasEmptyAssigned {
			t.Fatal("expected at least one task with empty assignedTo")
		}
		if !hasPartialAssigned {
			t.Fatal("expected at least one task with partial assignedTo")
		}
	})

	// Test daily schedule from Postgres-backed store.
	t.Run("DailySchedule", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
		}

		var body DailyScheduleBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		// Default window: start from planning window (2026-07-05), 4 days.
		if body.Range.StartDate != "2026-07-05" {
			t.Fatalf("expected startDate '2026-07-05', got %q", body.Range.StartDate)
		}
		if body.Range.EndDate != "2026-07-08" {
			t.Fatalf("expected endDate '2026-07-08', got %q", body.Range.EndDate)
		}
		if body.Range.Days != 4 {
			t.Fatalf("expected days=4, got %d", body.Range.Days)
		}
		if len(body.Days) != 4 {
			t.Fatalf("expected 4 day entries, got %d", len(body.Days))
		}

		// Verify each day has tasks and availablePeopleCount is populated.
		allPriorities := map[string]bool{}
		allStaffing := map[string]bool{}
		hasFullyStaffed2x2 := false
		hasFullyStaffed1x1 := false
		hasUnderstaffed := false
		for _, d := range body.Days {
			if d.AvailablePeopleCount < 1 {
				t.Fatalf("day %s: availablePeopleCount=%d, expected >= 1", d.Date, d.AvailablePeopleCount)
			}
			if len(d.Tasks) < 1 {
				t.Fatalf("day %s: expected at least 1 task, got %d", d.Date, len(d.Tasks))
			}
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
				// Verify assigned people identities match people table.
				for _, ap := range task.AssignedPeople {
					if ap.ID == "" || ap.Name == "" || ap.Initials == "" {
						t.Fatalf("task %s has incomplete assigned person: id=%q name=%q initials=%q", task.ID, ap.ID, ap.Name, ap.Initials)
					}
				}
				// Derived-field invariants.
				if task.AssignedCount != len(task.AssignedPeople) {
					t.Fatalf("task %s: assignedCount=%d != len(assignedPeople)=%d", task.ID, task.AssignedCount, len(task.AssignedPeople))
				}
				if task.PeopleNeeded < 1 {
					t.Fatalf("task %s: peopleNeeded=%d < 1", task.ID, task.PeopleNeeded)
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
	})

	// Test daily schedule with explicit params.
	t.Run("DailyScheduleExplicit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/dashboard/daily-schedule?start=2026-07-10&days=3", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
		}

		var body DailyScheduleBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
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
	})

	// Test hello endpoint.
	t.Run("Hello", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
	})

	// Verify OpenAPI includes all endpoints.
	t.Run("OpenAPI", func(t *testing.T) {
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
		expectedPaths := []string{"/api/hello", "/api/planning-window", "/api/dashboard/people-availability", "/api/tasks/backlog", "/api/dashboard/daily-schedule"}
		for _, p := range expectedPaths {
			if _, exists := paths[p]; !exists {
				t.Fatalf("OpenAPI paths missing %q", p)
			}
		}
	})
}
