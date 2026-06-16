package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Domain input types ----------

// CreateScheduleCardInput is the domain-layer input for creating/updating a schedule card.
type CreateScheduleCardInput struct {
	Title         string
	Priority      string
	RoomArea      string
	PeopleNeeded  int
	ScheduledDate string
	AssignedTo    []string
}

// ErrScheduleCardNotFound is returned when a schedule card ID is not found.
var ErrScheduleCardNotFound = errors.New("schedule card not found")

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

// ---------- Schedule-card CRUD request / response types ----------

// CreateScheduleCardRequestBody holds the fields for creating/updating a schedule card.
type CreateScheduleCardRequestBody struct {
	Title         string   `json:"title" required:"true" doc:"Task title"`
	Priority      string   `json:"priority" required:"true" enum:"high,medium,low" doc:"One of: high, medium, low"`
	RoomArea      string   `json:"roomArea" required:"true" doc:"Room or area name"`
	PeopleNeeded  int      `json:"peopleNeeded" required:"true" minimum:"1" doc:"Number of people needed for the task (>=1)"`
	ScheduledDate string   `json:"scheduledDate" required:"true" doc:"ISO 8601 date (YYYY-MM-DD) the card is scheduled for"`
	AssignedTo    []string `json:"assignedTo" doc:"Person-ID strings for assigned helpers, may be empty"`
}

// CreateScheduleCardInputHuma is the Huma input for POST /api/schedule/cards.
type CreateScheduleCardInputHuma struct {
	Body CreateScheduleCardRequestBody
}

// CreateScheduleCardOutput is the response wrapper for POST /api/schedule/cards.
type CreateScheduleCardOutput struct {
	Body TaskCard
}

// UpdateScheduleCardInputHuma is the Huma input for PUT /api/schedule/cards/{id}.
type UpdateScheduleCardInputHuma struct {
	ID   string                            `path:"id" doc:"Schedule card identifier (e.g., sched-1)"`
	Body CreateScheduleCardRequestBody
}

// UpdateScheduleCardOutput is the response wrapper for PUT /api/schedule/cards/{id}.
type UpdateScheduleCardOutput struct {
	Body TaskCard
}

// DeleteScheduleCardInput is the Huma input for DELETE /api/schedule/cards/{id}.
type DeleteScheduleCardInput struct {
	ID string `path:"id" doc:"Schedule card identifier (e.g., sched-1)"`
}

// DeleteScheduleCardOutput is the empty response for DELETE /api/schedule/cards/{id}.
type DeleteScheduleCardOutput struct{}

// ---------- Schedule-card CRUD validation ----------

// validateScheduleCardInput checks domain-level constraints and returns a Huma error or nil.
func validateScheduleCardInput(body CreateScheduleCardRequestBody, store Store, ctx context.Context) error {
	if body.Title == "" {
		return huma.Error400BadRequest("title is required")
	}
	if body.Priority != "high" && body.Priority != "medium" && body.Priority != "low" {
		return huma.Error400BadRequest("priority must be one of: high, medium, low")
	}
	if body.PeopleNeeded < 1 {
		return huma.Error400BadRequest("peopleNeeded must be at least 1")
	}
	if body.RoomArea == "" {
		return huma.Error400BadRequest("roomArea is required")
	}

	// Validate scheduled date is parseable.
	scheduledDate, err := time.Parse("2006-01-02", body.ScheduledDate)
	if err != nil {
		return huma.Error400BadRequest("scheduledDate must be a valid ISO 8601 date (YYYY-MM-DD)")
	}

	// Validate scheduled date is within the planning window.
	pw, err := store.GetPlanningWindow(ctx)
	if err != nil {
		return huma.Error500InternalServerError("failed to validate planning window", err)
	}
	pwStart, _ := time.Parse("2006-01-02", pw.StartDate)
	pwEnd, _ := time.Parse("2006-01-02", pw.EndDate)
	if scheduledDate.Before(pwStart) || scheduledDate.After(pwEnd) {
		return huma.Error400BadRequest(fmt.Sprintf("scheduledDate must be within the planning window (%s to %s)", pw.StartDate, pw.EndDate))
	}

	// Validate assigned people exist.
	for _, pid := range body.AssignedTo {
		exists, err := store.PersonExists(ctx, pid)
		if err != nil {
			return huma.Error500InternalServerError("failed to validate assigned person", err)
		}
		if !exists {
			return huma.Error400BadRequest("assigned person '" + pid + "' not found")
		}
	}

	return nil
}

// ---------- Schedule-card CRUD handlers ----------

func registerScheduleCardEndpoints(api huma.API, store Store) {
	// POST /api/schedule/cards — create a new schedule card.
	huma.Register(api, huma.Operation{
		OperationID: "create-schedule-card",
		Method:      http.MethodPost,
		Path:        "/api/schedule/cards",
		Summary:     "Create a schedule card",
		Description: "Creates a new schedule card with server-assigned sched-* ID. Returns 400 for invalid input or out-of-window date.",
		Tags:        []string{"Schedule"},
	}, func(ctx context.Context, input *CreateScheduleCardInputHuma) (*CreateScheduleCardOutput, error) {
		if err := validateScheduleCardInput(input.Body, store, ctx); err != nil {
			return nil, err
		}

		card, err := store.CreateScheduleCard(ctx, CreateScheduleCardInput{
			Title:         input.Body.Title,
			Priority:      input.Body.Priority,
			RoomArea:      input.Body.RoomArea,
			PeopleNeeded:  input.Body.PeopleNeeded,
			ScheduledDate: input.Body.ScheduledDate,
			AssignedTo:    input.Body.AssignedTo,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create schedule card", err)
		}

		return &CreateScheduleCardOutput{Body: *card}, nil
	})

	// PUT /api/schedule/cards/{id} — update an existing schedule card.
	huma.Register(api, huma.Operation{
		OperationID: "update-schedule-card",
		Method:      http.MethodPut,
		Path:        "/api/schedule/cards/{id}",
		Summary:     "Update a schedule card",
		Description: "Updates a schedule card and replaces assignments transactionally. Returns 400 for validation errors, 404 if the card ID is unknown.",
		Tags:        []string{"Schedule"},
	}, func(ctx context.Context, input *UpdateScheduleCardInputHuma) (*UpdateScheduleCardOutput, error) {
		if err := validateScheduleCardInput(input.Body, store, ctx); err != nil {
			return nil, err
		}

		card, err := store.UpdateScheduleCard(ctx, input.ID, CreateScheduleCardInput{
			Title:         input.Body.Title,
			Priority:      input.Body.Priority,
			RoomArea:      input.Body.RoomArea,
			PeopleNeeded:  input.Body.PeopleNeeded,
			ScheduledDate: input.Body.ScheduledDate,
			AssignedTo:    input.Body.AssignedTo,
		})
		if err != nil {
			if errors.Is(err, ErrScheduleCardNotFound) {
				return nil, huma.Error404NotFound("schedule card not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to update schedule card", err)
		}

		return &UpdateScheduleCardOutput{Body: *card}, nil
	})

	// DELETE /api/schedule/cards/{id} — delete a schedule card.
	huma.Register(api, huma.Operation{
		OperationID:  "delete-schedule-card",
		Method:       http.MethodDelete,
		Path:         "/api/schedule/cards/{id}",
		Summary:      "Delete a schedule card",
		Description:  "Deletes a schedule card and its assignments transactionally. Returns 404 if the card ID is unknown.",
		Tags:         []string{"Schedule"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *DeleteScheduleCardInput) (*DeleteScheduleCardOutput, error) {
		err := store.DeleteScheduleCard(ctx, input.ID)
		if err != nil {
			if errors.Is(err, ErrScheduleCardNotFound) {
				return nil, huma.Error404NotFound("schedule card not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to delete schedule card", err)
		}
		return &DeleteScheduleCardOutput{}, nil
	})
}
