package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ErrToolNotFound is returned when a tool ID is not found.
var ErrToolNotFound = errors.New("tool not found")

// ---------- Domain input types ----------

// CreateToolInput is the domain-layer input for creating a tool.
type CreateToolInput struct {
	Name string
}

// UpdateToolInput is the domain-layer input for updating a tool.
type UpdateToolInput struct {
	Name      string
	SortOrder int
}

// ---------- Response types ----------

// Tool represents a single tool in the list. A null BroughtBy means the tool
// is open; a non-null value is the bringer's person ID (crossed off).
type Tool struct {
	ID        string  `json:"id" doc:"Stable tool identifier, prefixed 'tool-'"`
	Name      string  `json:"name" doc:"Human-readable tool name"`
	BroughtBy *string `json:"broughtBy" doc:"Bringer's person ID, or null when the tool is open"`
}

// ToolSummary holds derived coverage counts for the tool list.
type ToolSummary struct {
	Total   int `json:"total" doc:"Total number of tools"`
	Claimed int `json:"claimed" doc:"Count of tools with a bringer (crossed off)"`
	Open    int `json:"open" doc:"Count of tools without a bringer"`
}

// ToolsBody is the top-level response body for the tools read endpoint.
type ToolsBody struct {
	Summary ToolSummary `json:"summary" doc:"Derived coverage summary"`
	Tools   []Tool      `json:"tools" doc:"Tool rows ordered by sort order"`
}

// ---------- Request / Response wrappers ----------

// ToolsInput is empty: the read endpoint accepts no query parameters.
type ToolsInput struct{}

// ToolsOutput is the combined payload for GET /api/tools.
type ToolsOutput struct {
	Body ToolsBody
}

// CreateToolRequestBody holds the fields for creating a tool.
type CreateToolRequestBody struct {
	Name string `json:"name" required:"true" doc:"Human-readable tool name"`
}

// CreateToolInputHuma is the Huma input for POST /api/tools.
type CreateToolInputHuma struct {
	Body CreateToolRequestBody
}

// CreateToolOutput is the response wrapper for POST /api/tools.
type CreateToolOutput struct {
	Body Tool
}

// UpdateToolRequestBody holds the fields for updating a tool.
type UpdateToolRequestBody struct {
	Name      string `json:"name" required:"true" doc:"Updated tool name"`
	SortOrder int    `json:"sortOrder" required:"true" doc:"Stable sort order"`
}

// UpdateToolInputHuma is the Huma input for PUT /api/tools/{id}.
type UpdateToolInputHuma struct {
	ID   string `path:"id" doc:"Tool identifier"`
	Body UpdateToolRequestBody
}

// UpdateToolOutput is the response wrapper for PUT /api/tools/{id}.
type UpdateToolOutput struct {
	Body Tool
}

// DeleteToolInput is the Huma input for DELETE /api/tools/{id}.
type DeleteToolInput struct {
	ID string `path:"id" doc:"Tool identifier"`
}

// DeleteToolOutput is the empty response for DELETE /api/tools/{id}.
type DeleteToolOutput struct{}

// SetToolBringerInput is the Huma input for PUT /api/tools/{id}/bringer.
type SetToolBringerInput struct {
	ID   string `path:"id" doc:"Tool identifier"`
	Body struct {
		PersonID string `json:"personId" required:"true" doc:"ID of the person bringing the tool"`
	}
}

// SetToolBringerOutput is the response wrapper for PUT /api/tools/{id}/bringer.
type SetToolBringerOutput struct {
	Body Tool
}

// ClearToolBringerInput is the Huma input for DELETE /api/tools/{id}/bringer.
type ClearToolBringerInput struct {
	ID string `path:"id" doc:"Tool identifier"`
}

// ClearToolBringerOutput is the empty response for DELETE /api/tools/{id}/bringer.
type ClearToolBringerOutput struct{}

// ---------- Handlers ----------

func registerToolsEndpoints(api huma.API, store Store) {
	// GET /api/tools — read all tools plus a coverage summary.
	huma.Register(api, huma.Operation{
		OperationID: "get-tools",
		Method:      http.MethodGet,
		Path:        "/api/tools",
		Summary:     "Tool list with coverage summary",
		Description: "Returns all tools ordered by sort order plus a derived coverage summary (total, claimed, open). " +
			"A tool is claimed (crossed off) when it has a bringer and open when it has none.",
		Tags: []string{"Tools"},
	}, func(ctx context.Context, input *ToolsInput) (*ToolsOutput, error) {
		body, err := store.GetTools(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to retrieve tools", err)
		}
		return &ToolsOutput{Body: *body}, nil
	})

	// POST /api/tools — create a new tool.
	huma.Register(api, huma.Operation{
		OperationID:   "create-tool",
		Method:        http.MethodPost,
		Path:          "/api/tools",
		Summary:       "Create a tool",
		Description:   "Creates a new tool with a server-assigned tool-* ID and append sort order, with no bringer. Returns 400 for an empty name.",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Tools"},
	}, func(ctx context.Context, input *CreateToolInputHuma) (*CreateToolOutput, error) {
		if input.Body.Name == "" {
			return nil, huma.Error400BadRequest("name is required")
		}
		tool, err := store.CreateTool(ctx, CreateToolInput{Name: input.Body.Name})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create tool", err)
		}
		return &CreateToolOutput{Body: *tool}, nil
	})

	// PUT /api/tools/{id} — update a tool's name and sort order.
	huma.Register(api, huma.Operation{
		OperationID: "update-tool",
		Method:      http.MethodPut,
		Path:        "/api/tools/{id}",
		Summary:     "Update a tool",
		Description: "Updates a tool's name and sort order. Returns 400 for an empty name, 404 if the tool ID is unknown.",
		Tags:        []string{"Tools"},
	}, func(ctx context.Context, input *UpdateToolInputHuma) (*UpdateToolOutput, error) {
		if input.Body.Name == "" {
			return nil, huma.Error400BadRequest("name is required")
		}
		tool, err := store.UpdateTool(ctx, input.ID, UpdateToolInput{
			Name:      input.Body.Name,
			SortOrder: input.Body.SortOrder,
		})
		if err != nil {
			if errors.Is(err, ErrToolNotFound) {
				return nil, huma.Error404NotFound("tool not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to update tool", err)
		}
		return &UpdateToolOutput{Body: *tool}, nil
	})

	// DELETE /api/tools/{id} — delete a tool.
	huma.Register(api, huma.Operation{
		OperationID:   "delete-tool",
		Method:        http.MethodDelete,
		Path:          "/api/tools/{id}",
		Summary:       "Delete a tool",
		Description:   "Deletes a tool. Returns 404 if the tool ID is unknown.",
		Tags:          []string{"Tools"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *DeleteToolInput) (*DeleteToolOutput, error) {
		if err := store.DeleteTool(ctx, input.ID); err != nil {
			if errors.Is(err, ErrToolNotFound) {
				return nil, huma.Error404NotFound("tool not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to delete tool", err)
		}
		return &DeleteToolOutput{}, nil
	})

	// PUT /api/tools/{id}/bringer — claim a tool by setting its bringer.
	huma.Register(api, huma.Operation{
		OperationID: "set-tool-bringer",
		Method:      http.MethodPut,
		Path:        "/api/tools/{id}/bringer",
		Summary:     "Claim a tool",
		Description: "Sets the tool's bringer, replacing any existing one. Claiming does not depend on the person's availability. " +
			"Returns 404 for an unknown tool and 400 for an unknown person.",
		Tags: []string{"Tools"},
	}, func(ctx context.Context, input *SetToolBringerInput) (*SetToolBringerOutput, error) {
		toolExists, err := store.ToolExists(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check tool existence", err)
		}
		if !toolExists {
			return nil, huma.Error404NotFound("tool not found")
		}

		personExists, err := store.PersonExists(ctx, input.Body.PersonID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check person existence", err)
		}
		if !personExists {
			return nil, huma.Error400BadRequest("person '" + input.Body.PersonID + "' not found")
		}

		tool, err := store.SetToolBringer(ctx, input.ID, input.Body.PersonID)
		if err != nil {
			if errors.Is(err, ErrToolNotFound) {
				return nil, huma.Error404NotFound("tool not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to set tool bringer", err)
		}
		return &SetToolBringerOutput{Body: *tool}, nil
	})

	// DELETE /api/tools/{id}/bringer — unclaim a tool. Idempotent.
	huma.Register(api, huma.Operation{
		OperationID:   "clear-tool-bringer",
		Method:        http.MethodDelete,
		Path:          "/api/tools/{id}/bringer",
		Summary:       "Unclaim a tool",
		Description:   "Clears the tool's bringer, returning it to open. Idempotent: clearing an already-open tool succeeds. Returns 404 for an unknown tool.",
		Tags:          []string{"Tools"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *ClearToolBringerInput) (*ClearToolBringerOutput, error) {
		if err := store.ClearToolBringer(ctx, input.ID); err != nil {
			if errors.Is(err, ErrToolNotFound) {
				return nil, huma.Error404NotFound("tool not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to clear tool bringer", err)
		}
		return &ClearToolBringerOutput{}, nil
	})
}
