package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Request / Response types ----------

// PlanningWindowInput is empty: the endpoint accepts no query parameters.
type PlanningWindowInput struct{}

// PlanningWindowOutput is the response wrapper for GET /api/planning-window.
type PlanningWindowOutput struct {
	Body PlanningWindowBody
}

// PlanningWindowBody is the response body for the planning-window endpoint.
type PlanningWindowBody struct {
	StartDate string `json:"startDate" doc:"Start date of the planning window (ISO 8601)"`
	EndDate   string `json:"endDate" doc:"End date of the planning window (ISO 8601)"`
	Days      int    `json:"days" doc:"Inclusive day count between startDate and endDate"`
}

// ---------- Handler ----------

func registerPlanningWindow(api huma.API, store Store) {
	huma.Register(api, huma.Operation{
		OperationID: "get-planning-window",
		Method:      http.MethodGet,
		Path:        "/api/planning-window",
		Summary:     "Global planning window",
		Description: "Returns the global move range with startDate, endDate, and inclusive day count. " +
			"The planning window is the canonical source of truth for the move timeline. " +
			"All date-dependent views derive their rendered content from this contract.",
		Tags: []string{"Planning"},
	}, func(ctx context.Context, input *PlanningWindowInput) (*PlanningWindowOutput, error) {
		body, err := store.GetPlanningWindow(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to retrieve planning window", err)
		}

		return &PlanningWindowOutput{
			Body: *body,
		}, nil
	})
}
