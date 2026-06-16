package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Request / Response types ----------

// DailyScheduleInput holds query parameters for GET /api/dashboard/daily-schedule.
type DailyScheduleInput struct {
	Start string `query:"start" doc:"Start date in ISO 8601 format (YYYY-MM-DD). Defaults to the planning-window start date 2026-07-05."`
	Days  int    `query:"days" default:"4" minimum:"1" doc:"Number of days inclusive of start date."`
}

// DailyScheduleOutput is the response wrapper for GET /api/dashboard/daily-schedule.
type DailyScheduleOutput struct {
	Body DailyScheduleBody
}

// DailyScheduleBody is the top-level response body for the daily schedule board.
type DailyScheduleBody struct {
	Range ScheduleRange `json:"range" doc:"Date range metadata"`
	Days  []ScheduleDay `json:"days" doc:"One entry per date in the requested window"`
}

// ScheduleRange holds date-range metadata for the daily schedule.
type ScheduleRange struct {
	StartDate string `json:"startDate" doc:"Start date of the window (ISO 8601)"`
	EndDate   string `json:"endDate" doc:"End date of the window (ISO 8601)"`
	Days      int    `json:"days" doc:"Number of dates in the window (inclusive of startDate)"`
}

// ScheduleDay represents one day column in the daily schedule board.
type ScheduleDay struct {
	Date                 string     `json:"date" doc:"Date in ISO 8601 format (YYYY-MM-DD)"`
	Label                string     `json:"label" doc:"Human-readable day header"`
	AvailablePeopleCount int        `json:"availablePeopleCount" doc:"Number of people available on this date"`
	Tasks                []TaskCard `json:"tasks" doc:"Ordered schedule task cards"`
}

// TaskCard is a single schedule task card displayed on the dashboard.
type TaskCard struct {
	ID             string           `json:"id" doc:"Stable task identifier"`
	Title          string           `json:"title" doc:"Task title"`
	Priority       string           `json:"priority" doc:"One of: high, medium, low"`
	RoomArea       string           `json:"roomArea" doc:"Room or area name"`
	AssignedPeople []AssignedPerson `json:"assignedPeople" doc:"People assigned to this task"`
	PeopleNeeded   int              `json:"peopleNeeded" doc:"Number of people needed for the task (>=1)"`
	AssignedCount  int              `json:"assignedCount" doc:"Number of people currently assigned (derived from assignedPeople)"`
	StaffingStatus string           `json:"staffingStatus" doc:"One of: fullyStaffed, underStaffed"`
}

// AssignedPerson represents a person assigned to a task card.
type AssignedPerson struct {
	ID       string `json:"id" doc:"Stable person key"`
	Name     string `json:"name" doc:"Full name"`
	Initials string `json:"initials" doc:"Initials"`
}

// ---------- Helpers ----------

var dayLabels = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
var monthLabels = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// FormatDayLabel returns a human-readable day label (e.g., "Mon, Jan 5").
func FormatDayLabel(t time.Time) string {
	return dayLabels[t.Weekday()] + ", " + monthLabels[t.Month()-1] + " " + t.Format("2")
}

// ---------- Handler ----------

func registerDailySchedule(api huma.API, store Store) {
	huma.Register(api, huma.Operation{
		OperationID: "get-dashboard-daily-schedule",
		Method:      http.MethodGet,
		Path:        "/api/dashboard/daily-schedule",
		Summary:     "Daily schedule for the dashboard",
		Description: "Returns a dashboard-ready daily schedule read model with date columns, " +
			"available-helper counts, and ordered schedule task cards. " +
			"The start parameter defaults to the planning-window start date (2026-07-05); " +
			"clients should pass start explicitly for a different window.",
		Tags: []string{"Dashboard"},
	}, func(ctx context.Context, input *DailyScheduleInput) (*DailyScheduleOutput, error) {
		// Parse or default start date.
		var startDate time.Time
		if input.Start != "" {
			var err error
			startDate, err = time.Parse("2006-01-02", input.Start)
			if err != nil {
				return nil, huma.Error400BadRequest("start must be a valid ISO 8601 date (YYYY-MM-DD)")
			}
		} else {
			pw, err := store.GetPlanningWindow(ctx)
			if err != nil {
				return nil, huma.Error500InternalServerError("failed to retrieve planning window for default start", err)
			}
			startDate, _ = time.Parse("2006-01-02", pw.StartDate)
		}

		days := input.Days

		body, err := store.GetDailySchedule(ctx, startDate, days)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to retrieve daily schedule", err)
		}

		return &DailyScheduleOutput{
			Body: *body,
		}, nil
	})
}
