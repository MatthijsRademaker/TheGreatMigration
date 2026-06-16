package main

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Request / Response types ----------

// PlanningWindowInput is empty: the endpoint accepts no query parameters.
type PlanningWindowInput struct{}

// PlanningWindowOutput is the response wrapper for GET /api/planning-window.
type PlanningWindowOutput struct {
	Body PlanningWindowBody
}

// UpdatePlanningWindowInput is the request body for PUT /api/planning-window.
type UpdatePlanningWindowInput struct {
	Body struct {
		StartDate string `json:"startDate" required:"true" doc:"Start date (ISO 8601, YYYY-MM-DD)"`
		EndDate   string `json:"endDate" required:"true" doc:"End date (ISO 8601, YYYY-MM-DD)"`
	}
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

	huma.Register(api, huma.Operation{
		OperationID: "put-planning-window",
		Method:      http.MethodPut,
		Path:        "/api/planning-window",
		Summary:     "Update planning window",
		Description: "Updates the planning window start and end dates. Both dates are required and endDate must be >= startDate. Returns the updated planning window with recalculated inclusive day count.",
		Tags:        []string{"Planning"},
		Errors:      []int{422},
	}, func(ctx context.Context, input *UpdatePlanningWindowInput) (*PlanningWindowOutput, error) {
		startDate, err := time.Parse("2006-01-02", input.Body.StartDate)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid startDate format, expected YYYY-MM-DD", err)
		}
		endDate, err := time.Parse("2006-01-02", input.Body.EndDate)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid endDate format, expected YYYY-MM-DD", err)
		}
		if endDate.Before(startDate) {
			return nil, huma.Error422UnprocessableEntity("endDate must be >= startDate")
		}

		body, err := store.UpdatePlanningWindow(ctx, startDate, endDate)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to update planning window", err)
		}

		return &PlanningWindowOutput{
			Body: *body,
		}, nil
	})
}
