package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/user/the-great-migration/backend/api"
)

func main() {
	router := chi.NewMux()
	config := huma.DefaultConfig("The Great Migration API", "1.0.0")
	humaAPI := humachi.New(router, config)

	// Register all endpoints with nil store.
	// Huma's OpenAPI() only uses registered type metadata — handlers are never invoked.
	api.RegisterAll(humaAPI, nil)

	spec := humaAPI.OpenAPI()
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(spec); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding OpenAPI spec: %v\n", err)
		os.Exit(1)
	}
}
