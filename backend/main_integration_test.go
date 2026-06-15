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

	// Test tasks backlog still works.
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
		if len(body.Tasks) < 10 {
			t.Fatalf("expected at least 10 tasks, got %d", len(body.Tasks))
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
		expectedPaths := []string{"/api/hello", "/api/planning-window", "/api/dashboard/people-availability", "/api/tasks/backlog"}
		for _, p := range expectedPaths {
			if _, exists := paths[p]; !exists {
				t.Fatalf("OpenAPI paths missing %q", p)
			}
		}
	})
}
