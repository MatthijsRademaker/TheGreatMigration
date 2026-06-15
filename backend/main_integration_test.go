//go:build integration

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// startPostgresSidecar starts a disposable Postgres container using Docker
// with a published port. Returns the container name and DSN (via localhost).
func startPostgresSidecar(t *testing.T) (containerName, dsn string) {
	t.Helper()

	containerName = fmt.Sprintf("integ-postgres-%d", time.Now().UnixNano())

	// Start Postgres container with a random published port.
	cmd := exec.Command("docker", "run", "-d",
		"--name", containerName,
		"-p", "0:5432",
		"-e", "POSTGRES_DB=testdb",
		"-e", "POSTGRES_USER=test",
		"-e", "POSTGRES_PASSWORD=test",
		"postgres:16-alpine")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to start postgres container: %v\noutput: %s", err, out)
	}

	// Get the published port.
	cmd = exec.Command("docker", "port", containerName, "5432")
	out, err = cmd.CombinedOutput()
	if err != nil {
		cleanupPostgresSidecar(containerName)
		t.Fatalf("failed to get port mapping: %v\noutput: %s", err, out)
	}
	portMapping := strings.TrimSpace(string(out))
	// portMapping is like "0.0.0.0:32768" or "[::]:32768"
	parts := strings.Split(portMapping, ":")
	publishedPort := parts[len(parts)-1]

	// Wait for Postgres to be ready.
	dsn = fmt.Sprintf("postgres://test:test@localhost:%s/testdb?sslmode=disable", publishedPort)
	deadline := time.Now().Add(30 * time.Second)
	for {
		cmd = exec.Command("docker", "exec", containerName, "pg_isready", "-U", "test", "-d", "testdb")
		if err := cmd.Run(); err == nil {
			break
		}
		if time.Now().After(deadline) {
			cleanupPostgresSidecar(containerName)
			t.Fatalf("postgres did not become ready within 30s")
		}
		time.Sleep(1 * time.Second)
	}

	t.Logf("started postgres sidecar: %s (port: %s)", containerName, publishedPort)
	return containerName, dsn
}

// cleanupPostgresSidecar stops and removes the Postgres container.
func cleanupPostgresSidecar(containerName string) {
	if containerName != "" {
		exec.Command("docker", "stop", containerName).Run()
		exec.Command("docker", "rm", "-f", containerName).Run()
	}
}

// TestDBBackedEndpoints runs endpoint contract tests against a real Postgres database.
func TestDBBackedEndpoints(t *testing.T) {
	containerName, dsn := startPostgresSidecar(t)
	t.Cleanup(func() {
		cleanupPostgresSidecar(containerName)
	})

	// Create connection pool.
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse DSN: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	// Run migrations manually using raw SQL.
	migrations := []string{
		`CREATE TABLE planning_windows (id SERIAL PRIMARY KEY, start_date DATE NOT NULL, end_date DATE NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
		`CREATE TABLE people (id TEXT PRIMARY KEY, name TEXT NOT NULL, initials TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
		`CREATE TABLE availability (id SERIAL PRIMARY KEY, person_id TEXT NOT NULL REFERENCES people(id), date DATE NOT NULL, status TEXT NOT NULL CHECK (status IN ('available','busy','partial','off')), created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE(person_id, date))`,
	}
	for _, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			t.Fatalf("migration failed: %v\nsql: %s", err, m)
		}
	}

	// Seed data.
	seed := `
		INSERT INTO planning_windows (start_date, end_date) VALUES ('2026-07-05', '2026-08-13');
		INSERT INTO people (id, name, initials) VALUES
			('p1', 'Sophia Chen', 'SC'),
			('p2', 'Marcus Rivera', 'MR'),
			('p3', 'Elena Kowalski', 'EK'),
			('p4', 'James Okafor', 'JO'),
			('p5', 'Priya Nair', 'PN'),
			('p6', 'Thomas Berg', 'TB'),
			('p7', 'Amara Diallo', 'AD'),
			('p8', 'Noah Larsson', 'NL');
		INSERT INTO availability (person_id, date, status)
		SELECT
			p.id,
			d.date,
			CASE
				WHEN p.id IN ('p1','p2','p3','p4','p5','p6') THEN 'available'
				WHEN p.id = 'p7' THEN 'busy'
				WHEN p.id = 'p8' THEN
					CASE (d.date::date - '2026-07-05'::date) % 4
						WHEN 0 THEN 'off'
						WHEN 1 THEN 'partial'
						WHEN 2 THEN 'busy'
						WHEN 3 THEN 'available'
					END
			END AS status
		FROM people p
		CROSS JOIN generate_series('2026-07-05'::date, '2026-08-13'::date, '1 day'::interval) AS d(date);
	`
	if _, err := pool.Exec(ctx, seed); err != nil {
		t.Fatalf("seed failed: %v", err)
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
