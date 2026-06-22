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

	backendapi "github.com/user/the-great-migration/backend/api"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

//go:embed seed/*.sql
var seedFS embed.FS

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

		// Always bootstrap the singleton planning-window row if missing.
		// This is essential config, not demo data — the app cannot function
		// without it. Uses a 30-day default window starting today.
		if err := bootstrapPlanningWindow(sqlDB); err != nil {
			fmt.Fprintf(os.Stderr, "failed to bootstrap planning window: %v\n", err)
			os.Exit(1)
		}

		// Conditionally apply the demo seed dataset (gated by DB_SEED).
		// It is a separate goose dataset tracked in its own version table so
		// it stays independent of the schema chain and idempotent on restart.
		if shouldSeed() {
			goose.SetBaseFS(seedFS)
			goose.SetTableName("goose_seed_version")
			if err := goose.Up(sqlDB, "seed"); err != nil {
				fmt.Fprintf(os.Stderr, "failed to run seed dataset: %v\n", err)
				os.Exit(1)
			}
			goose.SetTableName("goose_db_version")
			fmt.Println("Demo seed dataset applied successfully")
		}
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

	// Register all API endpoints (from backend/api package).
	backendapi.RegisterAll(api, store)

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

// bootstrapPlanningWindow inserts a default 30-day planning window (today → today+29)
// only when the planning_windows table is empty. This ensures the app never starts
// without a valid planning window, regardless of whether demo seed was applied.
func bootstrapPlanningWindow(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT count(*) FROM planning_windows").Scan(&count); err != nil {
		return fmt.Errorf("checking planning_windows: %w", err)
	}
	if count > 0 {
		return nil
	}

	today := time.Now().Truncate(24 * time.Hour)
	end := today.AddDate(0, 0, 29)
	_, err := db.Exec(
		"INSERT INTO planning_windows (id, start_date, end_date) VALUES (1, $1, $2)",
		today.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)
	if err != nil {
		return fmt.Errorf("inserting default planning window: %w", err)
	}

	fmt.Printf("Bootstrapped default planning window: %s → %s\n",
		today.Format("2006-01-02"), end.Format("2006-01-02"))
	return nil
}

// shouldAutoMigrate returns true if DB_AUTO_MIGRATE is not explicitly set to "false".
func shouldAutoMigrate() bool {
	v := strings.ToLower(os.Getenv("DB_AUTO_MIGRATE"))
	return v != "false"
}

// shouldSeed returns true only if DB_SEED is explicitly set to "true".
// Seed data is opt-in (defaults to false), the inverse of auto-migrate.
func shouldSeed() bool {
	v := strings.ToLower(os.Getenv("DB_SEED"))
	return v == "true"
}
