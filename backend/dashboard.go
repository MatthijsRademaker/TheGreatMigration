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

// ---------- Seed data ----------

var seedPeople = []struct {
	Id       string
	Name     string
	Initials string
	// Status returns the availability status for a given date offset from start.
	Status func(dayOffset int) string
}{
	{Id: "p1", Name: "Sophia Chen", Initials: "SC", Status: always("available")},
	{Id: "p2", Name: "Marcus Rivera", Initials: "MR", Status: always("available")},
	{Id: "p3", Name: "Elena Kowalski", Initials: "EK", Status: always("available")},
	{Id: "p4", Name: "James Okafor", Initials: "JO", Status: always("available")},
	{Id: "p5", Name: "Priya Nair", Initials: "PN", Status: always("available")},
	{Id: "p6", Name: "Thomas Berg", Initials: "TB", Status: always("available")},
	{Id: "p7", Name: "Amara Diallo", Initials: "AD", Status: always("busy")},
	{Id: "p8", Name: "Noah Larsson", Initials: "NL", Status: cycleStatuses},
}

func always(s string) func(int) string {
	return func(_ int) string { return s }
}

var cycleOrder = []string{"off", "partial", "busy", "available"}

func cycleStatuses(dayOffset int) string {
	return cycleOrder[dayOffset%len(cycleOrder)]
}

// ---------- Status legend ----------

var statusLegend = []StatusLegend{
	{ID: "available", Label: "Available", ColorIntent: "success"},
	{ID: "busy", Label: "Busy", ColorIntent: "destructive"},
	{ID: "partial", Label: "Partial", ColorIntent: "warning"},
	{ID: "off", Label: "Off", ColorIntent: "muted"},
}

// ---------- Handler ----------

func registerDashboardPeopleAvailability(api huma.API) {
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

		// Build the range.
		endDate := startDate.AddDate(0, 0, days-1)

		// Build per-person availability.
		people := make([]Person, len(seedPeople))
		for i, sp := range seedPeople {
			avail := make([]AvailabilityEntry, days)
			for d := 0; d < days; d++ {
				date := startDate.AddDate(0, 0, d)
				avail[d] = AvailabilityEntry{
					Date:   date.Format("2006-01-02"),
					Status: sp.Status(d),
				}
			}
			people[i] = Person{
				ID:           sp.Id,
				Name:         sp.Name,
				Initials:     sp.Initials,
				Availability: avail,
			}
		}

		// Compute summary.
		selectedDate := startDate.Format("2006-01-02")
		availableToday := 0
		for _, p := range people {
			for _, e := range p.Availability {
				if e.Date == selectedDate && e.Status == "available" {
					availableToday++
					break
				}
			}
		}

		return &DashboardOutput{
			Body: DashboardBody{
				Range: Range{
					StartDate:    startDate.Format("2006-01-02"),
					EndDate:      endDate.Format("2006-01-02"),
					Days:         days,
					SelectedDate: selectedDate,
				},
				Summary: Summary{
					AvailableToday: availableToday,
					TotalPeople:    len(people),
				},
				People:   people,
				Statuses: statusLegend,
			},
		}, nil
	})

}
