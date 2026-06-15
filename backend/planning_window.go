package main

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Seed data ----------

const (
	planWindowStart = "2026-07-05"
	planWindowEnd   = "2026-08-13"
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

func registerPlanningWindow(api huma.API) {
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
		// planWindowStart and planWindowEnd are compile-time constants in the exact
		// "2006-01-02" layout, so time.Parse always succeeds. Errors are discarded
		// safely — the parsed values are only used to compute the inclusive day count.
		startDate, _ := time.Parse("2006-01-02", planWindowStart)
		endDate, _ := time.Parse("2006-01-02", planWindowEnd)
		days := int(endDate.Sub(startDate).Hours()/24) + 1

		return &PlanningWindowOutput{
			Body: PlanningWindowBody{
				StartDate: planWindowStart,
				EndDate:   planWindowEnd,
				Days:      days,
			},
		}, nil
	})
}
