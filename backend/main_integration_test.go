//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	backendapi "github.com/user/the-great-migration/backend/api"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	// Apply the demo seed dataset explicitly (mirrors main.go's DB_SEED pass)
	// under its own version table so the seeded-count assertions below hold.
	goose.SetBaseFS(seedFS)
	goose.SetTableName("goose_seed_version")
	if err := goose.Up(sqlDB, "seed"); err != nil {
		t.Fatalf("failed to run seed dataset: %v", err)
	}
	goose.SetTableName("goose_db_version")

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

		var body backendapi.PlanningWindowBody
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

	// Test planning window update round-trip.
	t.Run("PlanningWindowUpdate", func(t *testing.T) {
		bodyJSON := `{"startDate": "2026-07-15", "endDate": "2026-07-25"}`
		req := httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
		}

		var body backendapi.PlanningWindowBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if body.StartDate != "2026-07-15" {
			t.Fatalf("expected startDate '2026-07-15', got %q", body.StartDate)
		}
		if body.EndDate != "2026-07-25" {
			t.Fatalf("expected endDate '2026-07-25', got %q", body.EndDate)
		}
		if body.Days != 11 {
			t.Fatalf("expected days=11, got %d", body.Days)
		}

		// Verify persistence: GET should return updated values.
		req = httptest.NewRequest(http.MethodGet, "/api/planning-window", nil)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
		}

		var getBody backendapi.PlanningWindowBody
		if err := json.Unmarshal(rec.Body.Bytes(), &getBody); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if getBody.StartDate != "2026-07-15" {
			t.Fatalf("expected startDate '2026-07-15', got %q", getBody.StartDate)
		}
		if getBody.EndDate != "2026-07-25" {
			t.Fatalf("expected endDate '2026-07-25', got %q", getBody.EndDate)
		}
		if getBody.Days != 11 {
			t.Fatalf("expected days=11, got %d", getBody.Days)
		}

		// Restore the seeded planning window so later subtests that rely on the
		// default window (e.g. DailySchedule) see the seeded 2026-07-05 start.
		restoreJSON := `{"startDate": "2026-07-05", "endDate": "2026-08-13"}`
		restoreReq := httptest.NewRequest(http.MethodPut, "/api/planning-window", strings.NewReader(restoreJSON))
		restoreReq.Header.Set("Content-Type", "application/json")
		restoreRec := httptest.NewRecorder()
		router.ServeHTTP(restoreRec, restoreReq)
		if restoreRec.Code != http.StatusOK {
			t.Fatalf("failed to restore planning window: status %d", restoreRec.Code)
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

		var body backendapi.DashboardBody
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

		var body backendapi.TaskBacklogBody
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

	// Test task CRUD lifecycle against Postgres.
	t.Run("TasksCRUD", func(t *testing.T) {
		// Create: add a new task with assignments.
		createBody := `{"title":"Integration test task","priority":"high","peopleNeeded":2,"areaId":"room-1","status":"backlog","assignedTo":["p1","p2"]}`
		createReq := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createRec := httptest.NewRecorder()
		router.ServeHTTP(createRec, createReq)

		if createRec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", createRec.Code, createRec.Body.String())
		}

		var created backendapi.TaskRow
		if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to unmarshal created task: %v", err)
		}
		if created.Title != "Integration test task" {
			t.Fatalf("unexpected title: %q", created.Title)
		}
		if created.Priority != "high" || created.Status != "backlog" || created.PeopleNeeded != 2 || created.Area.ID != "room-1" || created.Area.Name != "Kitchen" {
			t.Fatalf("unexpected created fields: prio=%q status=%q ppl=%d area=%+v", created.Priority, created.Status, created.PeopleNeeded, created.Area)
		}
		if created.ID == "" || created.ID[:5] != "task-" {
			t.Fatalf("created task has invalid ID: %q", created.ID)
		}
		if len(created.AssignedTo) != 2 || created.AssignedTo[0] != "p1" || created.AssignedTo[1] != "p2" {
			t.Fatalf("unexpected assignedTo: %v", created.AssignedTo)
		}

		// Verify created task appears in backlog.
		backlogReq := httptest.NewRequest(http.MethodGet, "/api/tasks/backlog", nil)
		backlogRec := httptest.NewRecorder()
		router.ServeHTTP(backlogRec, backlogReq)

		if backlogRec.Code != http.StatusOK {
			t.Fatalf("backlog GET failed: %d", backlogRec.Code)
		}

		var backlog backendapi.TaskBacklogBody
		if err := json.Unmarshal(backlogRec.Body.Bytes(), &backlog); err != nil {
			t.Fatalf("failed to unmarshal backlog: %v", err)
		}

		found := false
		for _, task := range backlog.Tasks {
			if task.ID == created.ID {
				found = true
				if task.Title != "Integration test task" {
					t.Fatalf("backlog task has wrong title: %q", task.Title)
				}
				break
			}
		}
		if !found {
			t.Fatalf("created task %s not found in backlog", created.ID)
		}

		// Update: change title and assignments.
		updateBody := `{"title":"Updated integration task","priority":"medium","peopleNeeded":1,"areaId":"room-5","status":"ready","assignedTo":["p3"]}`
		updateReq := httptest.NewRequest(http.MethodPut, "/api/tasks/"+created.ID, strings.NewReader(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateRec := httptest.NewRecorder()
		router.ServeHTTP(updateRec, updateReq)

		if updateRec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", updateRec.Code, updateRec.Body.String())
		}

		var updated backendapi.TaskRow
		if err := json.Unmarshal(updateRec.Body.Bytes(), &updated); err != nil {
			t.Fatalf("failed to unmarshal updated task: %v", err)
		}
		if updated.ID != created.ID {
			t.Fatalf("updated task id changed: %q -> %q", created.ID, updated.ID)
		}
		if updated.Title != "Updated integration task" || updated.Priority != "medium" || updated.Status != "ready" {
			t.Fatalf("unexpected updated fields: title=%q prio=%q status=%q", updated.Title, updated.Priority, updated.Status)
		}
		if updated.PeopleNeeded != 1 || updated.Area.ID != "room-5" || updated.Area.Name != "Garage" {
			t.Fatalf("unexpected updated fields: ppl=%d area=%+v", updated.PeopleNeeded, updated.Area)
		}
		if len(updated.AssignedTo) != 1 || updated.AssignedTo[0] != "p3" {
			t.Fatalf("unexpected updated assignedTo: %v", updated.AssignedTo)
		}

		// Verify update appears in backlog.
		backlogReq2 := httptest.NewRequest(http.MethodGet, "/api/tasks/backlog", nil)
		backlogRec2 := httptest.NewRecorder()
		router.ServeHTTP(backlogRec2, backlogReq2)

		var backlog2 backendapi.TaskBacklogBody
		if err := json.Unmarshal(backlogRec2.Body.Bytes(), &backlog2); err != nil {
			t.Fatalf("failed to unmarshal backlog: %v", err)
		}
		foundUpdated := false
		for _, task := range backlog2.Tasks {
			if task.ID == created.ID {
				foundUpdated = true
				if task.Title != "Updated integration task" {
					t.Fatalf("backlog shows stale title: %q", task.Title)
				}
				if len(task.AssignedTo) != 1 || task.AssignedTo[0] != "p3" {
					t.Fatalf("backlog shows stale assignments: %v", task.AssignedTo)
				}
				break
			}
		}
		if !foundUpdated {
			t.Fatalf("updated task %s not found in backlog", created.ID)
		}

		// Delete the task.
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/tasks/"+created.ID, nil)
		deleteRec := httptest.NewRecorder()
		router.ServeHTTP(deleteRec, deleteReq)

		if deleteRec.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d\nbody: %s", deleteRec.Code, deleteRec.Body.String())
		}

		// Delete again: should 404.
		deleteReq2 := httptest.NewRequest(http.MethodDelete, "/api/tasks/"+created.ID, nil)
		deleteRec2 := httptest.NewRecorder()
		router.ServeHTTP(deleteRec2, deleteReq2)

		if deleteRec2.Code != http.StatusNotFound {
			t.Fatalf("expected status 404 on re-delete, got %d", deleteRec2.Code)
		}

		// Verify task is gone from backlog.
		backlogReq3 := httptest.NewRequest(http.MethodGet, "/api/tasks/backlog", nil)
		backlogRec3 := httptest.NewRecorder()
		router.ServeHTTP(backlogRec3, backlogReq3)

		var backlog3 backendapi.TaskBacklogBody
		if err := json.Unmarshal(backlogRec3.Body.Bytes(), &backlog3); err != nil {
			t.Fatalf("failed to unmarshal backlog: %v", err)
		}
		for _, task := range backlog3.Tasks {
			if task.ID == created.ID {
				t.Fatalf("deleted task %s still present in backlog", created.ID)
			}
		}

		// Update non-existent task: should 404.
		updateReq2 := httptest.NewRequest(http.MethodPut, "/api/tasks/"+created.ID, strings.NewReader(updateBody))
		updateReq2.Header.Set("Content-Type", "application/json")
		updateRec2 := httptest.NewRecorder()
		router.ServeHTTP(updateRec2, updateReq2)

		if updateRec2.Code != http.StatusNotFound {
			t.Fatalf("expected status 404 for update on deleted task, got %d", updateRec2.Code)
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

		var body backendapi.DailyScheduleBody
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

		var body backendapi.DailyScheduleBody
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

	// Test room CRUD lifecycle against Postgres.
	t.Run("RoomCRUD", func(t *testing.T) {
		// List: seed data has 8 rooms.
		listReq := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
		listRec := httptest.NewRecorder()
		router.ServeHTTP(listRec, listReq)

		if listRec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", listRec.Code, listRec.Body.String())
		}

		var listBody struct {
			Rooms []backendapi.Room `json:"rooms"`
		}
		if err := json.Unmarshal(listRec.Body.Bytes(), &listBody); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(listBody.Rooms) != 8 {
			t.Fatalf("expected 8 seed rooms, got %d", len(listBody.Rooms))
		}
		for _, r := range listBody.Rooms {
			if r.ID == "" || r.Name == "" || (r.Type != "room" && r.Type != "area") {
				t.Fatalf("invalid room: id=%q name=%q type=%q", r.ID, r.Name, r.Type)
			}
			if r.CreatedAt == "" || r.UpdatedAt == "" {
				t.Fatalf("room %s missing timestamps", r.ID)
			}
		}

		// Create: add a new room.
		createBody := `{"name":"Integration Test Room","type":"room"}`
		createReq := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createRec := httptest.NewRecorder()
		router.ServeHTTP(createRec, createReq)

		if createRec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d\nbody: %s", createRec.Code, createRec.Body.String())
		}

		var created backendapi.Room
		if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if created.Name != "Integration Test Room" || created.Type != "room" {
			t.Fatalf("unexpected created room: %+v", created)
		}
		if created.ID == "" {
			t.Fatal("created room has empty id")
		}

		// List again: now 9 rooms.
		listReq2 := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
		listRec2 := httptest.NewRecorder()
		router.ServeHTTP(listRec2, listReq2)

		var listBody2 struct {
			Rooms []backendapi.Room `json:"rooms"`
		}
		if err := json.Unmarshal(listRec2.Body.Bytes(), &listBody2); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(listBody2.Rooms) != 9 {
			t.Fatalf("expected 9 rooms after create, got %d", len(listBody2.Rooms))
		}

		// Update: rename the created room.
		updateBody := `{"name":"Updated Integration Room","type":"area"}`
		updateReq := httptest.NewRequest(http.MethodPut, "/api/rooms/"+created.ID, strings.NewReader(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateRec := httptest.NewRecorder()
		router.ServeHTTP(updateRec, updateReq)

		if updateRec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d\nbody: %s", updateRec.Code, updateRec.Body.String())
		}

		var updated backendapi.Room
		if err := json.Unmarshal(updateRec.Body.Bytes(), &updated); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if updated.Name != "Updated Integration Room" || updated.Type != "area" {
			t.Fatalf("unexpected updated room: %+v", updated)
		}
		if updated.ID != created.ID {
			t.Fatalf("updated room id changed: %q -> %q", created.ID, updated.ID)
		}

		// Delete the room.
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/rooms/"+created.ID, nil)
		deleteRec := httptest.NewRecorder()
		router.ServeHTTP(deleteRec, deleteReq)

		if deleteRec.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d\nbody: %s", deleteRec.Code, deleteRec.Body.String())
		}

		// Delete again: should 404.
		deleteReq2 := httptest.NewRequest(http.MethodDelete, "/api/rooms/"+created.ID, nil)
		deleteRec2 := httptest.NewRecorder()
		router.ServeHTTP(deleteRec2, deleteReq2)

		if deleteRec2.Code != http.StatusNotFound {
			t.Fatalf("expected status 404 on re-delete, got %d\nbody: %s", deleteRec2.Code, deleteRec2.Body.String())
		}

		// List: back to 8 rooms.
		listReq3 := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
		listRec3 := httptest.NewRecorder()
		router.ServeHTTP(listRec3, listReq3)

		var listBody3 struct {
			Rooms []backendapi.Room `json:"rooms"`
		}
		if err := json.Unmarshal(listRec3.Body.Bytes(), &listBody3); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(listBody3.Rooms) != 8 {
			t.Fatalf("expected 8 rooms after delete, got %d", len(listBody3.Rooms))
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
		expectedPaths := []string{"/api/hello", "/api/planning-window", "/api/dashboard/people-availability", "/api/tasks", "/api/tasks/{id}", "/api/tasks/backlog", "/api/dashboard/daily-schedule", "/api/rooms", "/api/rooms/{id}"}
		for _, p := range expectedPaths {
			if _, exists := paths[p]; !exists {
				t.Fatalf("OpenAPI paths missing %q", p)
			}
		}
	})

	// Write the OpenAPI spec to the frontend snapshot file for offline
	// code generation. The path is relative to the workspace root where
	// the frontend directory lives.
	t.Run("WriteSnapshot", func(t *testing.T) {
		snapshotPath := os.Getenv("OPENAPI_SNAPSHOT_PATH")
		if snapshotPath == "" {
			t.Skip("OPENAPI_SNAPSHOT_PATH not set; skipping snapshot write")
		}
		openapiBytes, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal OpenAPI: %v", err)
		}
		openapiBytes = append(openapiBytes, '\n')
		if err := os.WriteFile(snapshotPath, openapiBytes, 0o644); err != nil {
			t.Fatalf("failed to write OpenAPI snapshot to %s: %v", snapshotPath, err)
		}
		t.Logf("OpenAPI snapshot written to %s (%d bytes)", snapshotPath, len(openapiBytes))
	})
}

// TestSchemaOnlyStartup verifies that applying the schema migrations WITHOUT the
// seed dataset leaves all domain tables empty and the server still serves HTTP.
func TestSchemaOnlyStartup(t *testing.T) {
	baseDSN := os.Getenv("DATABASE_URL")
	if baseDSN == "" {
		t.Fatal("DATABASE_URL must be set (run via scripts/test-integration)")
	}

	ctx := context.Background()
	dsn := provisionFreshDB(t, baseDSN)
	applySchema(t, dsn) // schema only — no seed pass

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse DSN: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	// Domain tables must be empty when no seed is applied.
	for _, table := range []string{"people", "backlog_tasks", "rooms_areas", "schedule_task_cards"} {
		var count int
		if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
			t.Fatalf("failed to count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("expected %s to be empty without seed, got %d rows", table, count)
		}
	}

	// The server starts and serves HTTP against the schema-only database.
	store := NewPgStore(pool)
	router, _ := newTestAPI(store)
	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected hello 200 on schema-only DB, got %d", rec.Code)
	}
}

// TestSeedAdvancesSequences verifies the seed dataset advances the ID sequences
// past the seeded maximum so newly created entities get non-colliding IDs.
func TestSeedAdvancesSequences(t *testing.T) {
	baseDSN := os.Getenv("DATABASE_URL")
	if baseDSN == "" {
		t.Fatal("DATABASE_URL must be set (run via scripts/test-integration)")
	}

	ctx := context.Background()
	dsn := provisionFreshDB(t, baseDSN)
	applySchema(t, dsn)
	applySeed(t, dsn)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse DSN: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()
	store := NewPgStore(pool)

	personID, err := store.CreatePerson(ctx, "New Person", "NP")
	if err != nil {
		t.Fatalf("CreatePerson failed: %v", err)
	}
	if personID != "p9" {
		t.Fatalf("expected new person id 'p9', got %q", personID)
	}

	room, err := store.CreateRoom(ctx, backendapi.CreateRoomInput{Name: "New Room", Type: "room"})
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}
	if room.ID != "room-9" {
		t.Fatalf("expected new room id 'room-9', got %q", room.ID)
	}

	task, err := store.CreateTask(ctx, backendapi.CreateTaskInput{
		Title:        "New Task",
		Priority:     "high",
		PeopleNeeded: 1,
		AreaID:       "room-1",
		Status:       "backlog",
	})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if task.ID != "task-12" {
		t.Fatalf("expected new task id 'task-12', got %q", task.ID)
	}
}

// provisionFreshDB creates a brand-new empty database on the same Postgres
// server as baseDSN and returns a DSN pointing at it. The database is dropped
// via t.Cleanup. Tests use this when they need an isolated empty database, since
// the shared sidecar DB is mutated by TestDBBackedEndpoints.
func provisionFreshDB(t *testing.T, baseDSN string) string {
	t.Helper()

	u, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatalf("failed to parse DATABASE_URL: %v", err)
	}
	dbName := "tgm_test_" + sanitizeDBName(t.Name())

	admin, err := sql.Open("pgx", baseDSN)
	if err != nil {
		t.Fatalf("failed to open admin connection: %v", err)
	}
	defer admin.Close()

	ctx := context.Background()
	if _, err := admin.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName+" WITH (FORCE)"); err != nil {
		t.Fatalf("failed to drop pre-existing test database %s: %v", dbName, err)
	}
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE "+dbName); err != nil {
		t.Fatalf("failed to create test database %s: %v", dbName, err)
	}
	t.Cleanup(func() {
		cleanup, err := sql.Open("pgx", baseDSN)
		if err != nil {
			return
		}
		defer cleanup.Close()
		_, _ = cleanup.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+dbName+" WITH (FORCE)")
	})

	u.Path = "/" + dbName
	return u.String()
}

// sanitizeDBName turns a test name into a valid lowercase Postgres identifier.
func sanitizeDBName(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}

// applySchema runs the schema migrations (only) against the given DSN.
func applySchema(t *testing.T, dsn string) {
	t.Helper()
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open database for schema migrations: %v", err)
	}
	defer sqlDB.Close()
	goose.SetBaseFS(migrationsFS)
	goose.SetTableName("goose_db_version")
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		t.Fatalf("failed to run schema migrations: %v", err)
	}
}

// applySeed runs the demo seed dataset against the given DSN under its own
// version table, then restores the default table name.
func applySeed(t *testing.T, dsn string) {
	t.Helper()
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open database for seed dataset: %v", err)
	}
	defer sqlDB.Close()
	goose.SetBaseFS(seedFS)
	goose.SetTableName("goose_seed_version")
	if err := goose.Up(sqlDB, "seed"); err != nil {
		t.Fatalf("failed to run seed dataset: %v", err)
	}
	goose.SetTableName("goose_db_version")
}

// TestAreaBackfillMigration exercises migration 014's self-healing backfill:
// name-match linking, deterministic MIN(id) on duplicate names, and orphan
// auto-creation — by applying migrations up to v13, seeding legacy free-text
// rooms, then applying v14 and asserting the resulting area_id foreign keys.
func TestAreaBackfillMigration(t *testing.T) {
	baseDSN := os.Getenv("DATABASE_URL")
	if baseDSN == "" {
		t.Fatal("DATABASE_URL must be set (run via scripts/test-integration)")
	}
	dsn := provisionFreshDB(t, baseDSN)
	ctx := context.Background()

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer sqlDB.Close()

	// Apply schema up to v13 (before the area_id FK migration).
	goose.SetBaseFS(migrationsFS)
	goose.SetTableName("goose_db_version")
	if err := goose.UpTo(sqlDB, "migrations", 13); err != nil {
		t.Fatalf("failed to migrate to v13: %v", err)
	}

	// Seed legacy catalog: a normal room, plus a duplicate-named pair to test
	// deterministic MIN(id) resolution.
	if _, err := sqlDB.ExecContext(ctx, `
		INSERT INTO rooms_areas (id, name, type) VALUES
			('room-1', 'Kitchen', 'room'),
			('room-50', 'Attic', 'area'),
			('room-60', 'Attic', 'area');
		SELECT setval('rooms_areas_id_seq', 60);
	`); err != nil {
		t.Fatalf("failed to seed legacy rooms_areas: %v", err)
	}

	// Legacy free-text rows: a clean match, a duplicate-name match, and an
	// orphan string present in BOTH tables (should auto-create exactly one row).
	if _, err := sqlDB.ExecContext(ctx, `
		INSERT INTO backlog_tasks (id, title, priority, people_needed, room, status, sort_order) VALUES
			('task-a', 'Match',      'high',   1, 'Kitchen',   'backlog', 1),
			('task-b', 'DupName',    'medium', 1, 'Attic',     'backlog', 2),
			('task-c', 'Orphan',     'low',    1, 'Ghostroom', 'backlog', 3);
		INSERT INTO schedule_task_cards (title, priority, room_area, people_needed, scheduled_date, sort_order) VALUES
			('Card orphan', 'low', 'Ghostroom', 1, '2026-07-05', 0);
	`); err != nil {
		t.Fatalf("failed to seed legacy free-text rows: %v", err)
	}

	// Apply the area_id FK migration.
	if err := goose.UpTo(sqlDB, "migrations", 14); err != nil {
		t.Fatalf("failed to migrate to v14: %v", err)
	}

	// Clean match links to the catalog id.
	var areaID string
	if err := sqlDB.QueryRowContext(ctx, `SELECT area_id FROM backlog_tasks WHERE id = 'task-a'`).Scan(&areaID); err != nil {
		t.Fatalf("query task-a: %v", err)
	}
	if areaID != "room-1" {
		t.Fatalf("expected task-a area_id 'room-1', got %q", areaID)
	}

	// Duplicate name resolves to the numerically lowest id.
	if err := sqlDB.QueryRowContext(ctx, `SELECT area_id FROM backlog_tasks WHERE id = 'task-b'`).Scan(&areaID); err != nil {
		t.Fatalf("query task-b: %v", err)
	}
	if areaID != "room-50" {
		t.Fatalf("expected task-b area_id 'room-50' (MIN id), got %q", areaID)
	}

	// Orphan string auto-created a single 'area' row, shared by both tables.
	var taskOrphanArea, cardOrphanArea, orphanName, orphanType string
	if err := sqlDB.QueryRowContext(ctx, `SELECT area_id FROM backlog_tasks WHERE id = 'task-c'`).Scan(&taskOrphanArea); err != nil {
		t.Fatalf("query task-c: %v", err)
	}
	if err := sqlDB.QueryRowContext(ctx, `SELECT area_id FROM schedule_task_cards WHERE title = 'Card orphan'`).Scan(&cardOrphanArea); err != nil {
		t.Fatalf("query card orphan: %v", err)
	}
	if taskOrphanArea != cardOrphanArea {
		t.Fatalf("expected orphan rows to share one auto-created area, got %q and %q", taskOrphanArea, cardOrphanArea)
	}
	if err := sqlDB.QueryRowContext(ctx, `SELECT name, type FROM rooms_areas WHERE id = $1`, taskOrphanArea).Scan(&orphanName, &orphanType); err != nil {
		t.Fatalf("query auto-created area: %v", err)
	}
	if orphanName != "Ghostroom" || orphanType != "area" {
		t.Fatalf("expected auto-created area Ghostroom/area, got %q/%q", orphanName, orphanType)
	}

	// No row is left without an area_id.
	var nullCount int
	if err := sqlDB.QueryRowContext(ctx, `
		SELECT (SELECT COUNT(*) FROM backlog_tasks WHERE area_id IS NULL)
		     + (SELECT COUNT(*) FROM schedule_task_cards WHERE area_id IS NULL)
	`).Scan(&nullCount); err != nil {
		t.Fatalf("query null area_id count: %v", err)
	}
	if nullCount != 0 {
		t.Fatalf("expected 0 rows with null area_id, got %d", nullCount)
	}

	// The foreign key is enforced: an unknown area_id is rejected.
	if _, err := sqlDB.ExecContext(ctx, `
		INSERT INTO backlog_tasks (id, title, priority, people_needed, area_id, status, sort_order)
		VALUES ('task-bad', 'Bad', 'low', 1, 'room-does-not-exist', 'backlog', 99)
	`); err == nil {
		t.Fatal("expected FK violation inserting unknown area_id, got nil error")
	}
}
