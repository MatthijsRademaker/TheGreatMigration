package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// HelloInput is the input for the hello endpoint.
type HelloInput struct{}

// HelloOutput is the response for the hello endpoint.
type HelloOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello from the backend!"`
	}
}

func main() {
	// Initialize database connection pool.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintf(os.Stderr, "DATABASE_URL environment variable is required\n")
		os.Exit(1)
	}

	// Establish the connection pool first so that pool-creation failures
	// do not leave a partially migrated database.
	pool, err := connectDB(databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Run goose migrations at startup (gated by DB_AUTO_MIGRATE).
	if shouldAutoMigrate() {
		sqlDB, err := sql.Open("pgx", databaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open database for migrations: %v\n", err)
			os.Exit(1)
		}
		defer sqlDB.Close()

		goose.SetBaseFS(migrationsFS)
		if err := goose.Up(sqlDB, "migrations"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database migrations applied successfully")
	}

	store := NewPgStore(pool)

	router := chi.NewMux()

	// Configure CORS to allow Vite dev server and compose frontend origins.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://frontend:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	// Create the Huma v2 API with the chi adapter.
	config := huma.DefaultConfig("The Great Migration API", "1.0.0")
	config.DocsPath = "/docs"
	config.OpenAPIPath = "/openapi.json"
	api := humachi.New(router, config)

	// Serve the raw OpenAPI specification at /openapi.json.
	router.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(api.OpenAPI())
		w.Write(b)
	})

	// Register GET /api/hello.
	huma.Register(api, huma.Operation{
		OperationID: "get-hello",
		Method:      http.MethodGet,
		Path:        "/api/hello",
		Summary:     "Hello world",
		Description: "Returns a hello-world message confirming the backend is reachable.",
	}, func(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
		resp := &HelloOutput{}
		resp.Body.Message = "Hello from the backend!"
		return resp, nil
	})

	// Register GET /api/dashboard/people-availability.
	registerDashboardPeopleAvailability(api, store)

	// Register GET /api/planning-window.
	registerPlanningWindow(api, store)

	// Register GET /api/tasks/backlog.
	registerTasksBacklog(api, store)

	// Register GET /api/dashboard/daily-schedule.
	registerDailySchedule(api, store)

	fmt.Println("Backend listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
	}
}

// connectDB creates a pgx connection pool with retry logic.
func connectDB(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	deadline := time.Now().Add(30 * time.Second)
	for {
		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
		}

		err = pool.Ping(ctx)
		if err == nil {
			return pool, nil
		}
		pool.Close()

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out connecting to database after 30s: %w", err)
		}

		time.Sleep(1 * time.Second)
	}
}

// shouldAutoMigrate returns true if DB_AUTO_MIGRATE is not explicitly set to "false".
func shouldAutoMigrate() bool {
	v := strings.ToLower(os.Getenv("DB_AUTO_MIGRATE"))
	return v != "false"
}
