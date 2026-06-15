package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// HelloInput is the input for the hello endpoint.
type HelloInput struct{}

// HelloOutput is the response for the hello endpoint.
type HelloOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello from the backend!"`
	}
}

func main() {
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

	fmt.Println("Backend listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
	}
}
