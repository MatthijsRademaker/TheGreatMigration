package main

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Request / Response types ----------

// DashboardInput holds query parameters for GET /api/dashboard/people-availability.
type DashboardInput struct {
	Start string `query:"start" doc:"Start date in ISO 8601 format (YYYY-MM-DD). Defaults to the server-local current date."`
	Days  int    `query:"days" default:"4" minimum:"1" doc:"Number of days inclusive of start date."`
}

// DashboardOutput is the combined payload for the dashboard people-availability endpoint.
type DashboardOutput struct {
	Body DashboardBody
}

// DashboardBody is the top-level response body.
type DashboardBody struct {
	Range    Range           `json:"range" doc:"Date range metadata"`
	Summary  Summary         `json:"summary" doc:"Summary counts"`
	People   []Person        `json:"people" doc:"People with daily availability"`
	Statuses []StatusLegend  `json:"statuses" doc:"Canonical status legend"`
}

// Range holds date-range metadata.
type Range struct {
	StartDate    string `json:"startDate" doc:"Start date of the window (ISO 8601)"`
	EndDate      string `json:"endDate" doc:"End date of the window (ISO 8601)"`
	Days         int    `json:"days" doc:"Number of dates in the window (inclusive of startDate)"`
	SelectedDate string `json:"selectedDate" doc:"Reference date for summary counts (equals startDate)"`
}

// Summary holds the summary counts for the dashboard card.
type Summary struct {
	AvailableToday int `json:"availableToday" doc:"Number of people with status 'available' on selectedDate"`
	TotalPeople    int `json:"totalPeople" doc:"Total number of people in the response"`
}

// Person represents one person with their daily availability over the requested window.
type Person struct {
	ID           string             `json:"id" doc:"Stable person key"`
	Name         string             `json:"name" doc:"Full name"`
	Initials     string             `json:"initials" doc:"Initials"`
	Availability []AvailabilityEntry `json:"availability" doc:"One entry per date in the range"`
}

// AvailabilityEntry is a single date+status pair.
type AvailabilityEntry struct {
	Date   string `json:"date" doc:"Date in ISO 8601 format (YYYY-MM-DD)"`
	Status string `json:"status" doc:"One of: available, busy, partial, off"`
}

// StatusLegend is a canonical status definition from the design system.
type StatusLegend struct {
	ID          string `json:"id" doc:"Status identifier"`
	Label       string `json:"label" doc:"Human-readable label"`
	ColorIntent string `json:"colorIntent" doc:"Design system color intent"`
}

// ---------- Status legend ----------

var statusLegend = []StatusLegend{
	{ID: "available", Label: "Available", ColorIntent: "success"},
	{ID: "busy", Label: "Busy", ColorIntent: "destructive"},
	{ID: "partial", Label: "Partial", ColorIntent: "warning"},
	{ID: "off", Label: "Off", ColorIntent: "muted"},
}

// ---------- Handler ----------

func registerDashboardPeopleAvailability(api huma.API, store Store) {
	huma.Register(api, huma.Operation{
		OperationID: "get-dashboard-people-availability",
		Method:      http.MethodGet,
		Path:        "/api/dashboard/people-availability",
		Summary:     "People availability for the dashboard",
		Description: "Returns a combined payload with date-range metadata, summary counts, per-person daily availability, and a status legend. " +
			"The start parameter defaults to the server-local current date; clients should pass start explicitly for timezone-correct results. " +
			"availableToday counts only people whose status on selectedDate equals 'available'.",
		Tags: []string{"Dashboard"},
	}, func(ctx context.Context, input *DashboardInput) (*DashboardOutput, error) {
		// Parse or default start date.
		var startDate time.Time
		if input.Start != "" {
			var err error
			startDate, err = time.Parse("2006-01-02", input.Start)
			if err != nil {
				return nil, huma.Error400BadRequest("start must be a valid ISO 8601 date (YYYY-MM-DD)")
			}
		} else {
			now := time.Now()
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		}

		days := input.Days

		body, err := store.GetPeopleAvailability(ctx, startDate, days)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to retrieve availability data", err)
		}

		return &DashboardOutput{
			Body: *body,
		}, nil
	})

}
