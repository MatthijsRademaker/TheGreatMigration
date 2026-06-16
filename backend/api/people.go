package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---------- Canonical statuses ----------

var canonicalStatuses = map[string]bool{
	"available": true,
	"busy":      true,
	"partial":   true,
	"off":       true,
}

// ---------- Request / Response types ----------

// CreatePersonInput holds the body for POST /api/people.
type CreatePersonInput struct {
	Body struct {
		ID       string `json:"id" doc:"Client-supplied stable person key (slug)" minLength:"1" maxLength:"100"`
		Name     string `json:"name" doc:"Full name" minLength:"1" maxLength:"200"`
		Initials string `json:"initials" doc:"Initials" minLength:"1" maxLength:"10"`
	}
}

// CreatePersonOutput is the response for POST /api/people.
type CreatePersonOutput struct {
	Body Person
}

// UpdatePersonInput holds the body and path param for PUT /api/people/{id}.
type UpdatePersonInput struct {
	ID   string `path:"id" doc:"Stable person key"`
	Body struct {
		Name     string `json:"name" doc:"Full name" minLength:"1" maxLength:"200"`
		Initials string `json:"initials" doc:"Initials" minLength:"1" maxLength:"10"`
	}
}

// UpdatePersonOutput is the response for PUT /api/people/{id}.
type UpdatePersonOutput struct {
	Body Person
}

// DeletePersonInput holds the path param for DELETE /api/people/{id}.
type DeletePersonInput struct {
	ID string `path:"id" doc:"Stable person key"`
}

// DeletePersonOutput is the response for DELETE /api/people/{id}.
type DeletePersonOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// UpsertAvailabilityInput holds the path params and body for PUT /api/people/{id}/availability/{date}.
type UpsertAvailabilityInput struct {
	ID   string `path:"id" doc:"Stable person key"`
	Date string `path:"date" doc:"Date in ISO 8601 format (YYYY-MM-DD)"`
	Body struct {
		Status string `json:"status" doc:"One of: available, busy, partial, off"`
	}
}

// UpsertAvailabilityOutput is the response for PUT /api/people/{id}/availability/{date}.
type UpsertAvailabilityOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// DeleteAvailabilityInput holds the path params for DELETE /api/people/{id}/availability/{date}.
type DeleteAvailabilityInput struct {
	ID   string `path:"id" doc:"Stable person key"`
	Date string `path:"date" doc:"Date in ISO 8601 format (YYYY-MM-DD)"`
}

// DeleteAvailabilityOutput is the response for DELETE /api/people/{id}/availability/{date}.
type DeleteAvailabilityOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// ---------- Helpers ----------

// isWithinPlanningWindow checks whether the given date falls within the planning window.
func isWithinPlanningWindow(dateStr string, pw *PlanningWindowBody) bool {
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	start, err := time.Parse("2006-01-02", pw.StartDate)
	if err != nil {
		return false
	}
	end, err := time.Parse("2006-01-02", pw.EndDate)
	if err != nil {
		return false
	}
	return !d.Before(start) && !d.After(end)
}

// isoToDate parses an ISO date string and returns a time.Time.
func isoToDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// isForeignKeyViolation returns true if err is a PostgreSQL foreign-key violation (SQLSTATE 23503).
func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return true
	}
	return false
}

// ---------- Handlers ----------

// registerPeopleEndpoints registers all people management endpoints.
func registerPeopleEndpoints(api huma.API, store Store) {
	// POST /api/people
	huma.Register(api, huma.Operation{
		OperationID: "create-person",
		Method:      http.MethodPost,
		Path:        "/api/people",
		Summary:     "Create a person",
		Description: "Creates a new person with a client-supplied stable ID suitable for name-derived slugs.",
		Tags:        []string{"People"},
	}, func(ctx context.Context, input *CreatePersonInput) (*CreatePersonOutput, error) {
		// Validate fields are non-empty (Huma's minLength handles this, but double-check).
		if input.Body.ID == "" || input.Body.Name == "" || input.Body.Initials == "" {
			return nil, huma.Error400BadRequest("id, name, and initials are required")
		}

		// Check for duplicate ID.
		exists, err := store.PersonExists(ctx, input.Body.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check person existence", err)
		}
		if exists {
			return nil, huma.Error409Conflict("a person with this id already exists")
		}

		if err := store.CreatePerson(ctx, input.Body.ID, input.Body.Name, input.Body.Initials); err != nil {
			return nil, huma.Error500InternalServerError("failed to create person", err)
		}

		return &CreatePersonOutput{
			Body: Person{
				ID:           input.Body.ID,
				Name:         input.Body.Name,
				Initials:     input.Body.Initials,
				Availability: []AvailabilityEntry{},
			},
		}, nil
	})

	// PUT /api/people/{id}
	huma.Register(api, huma.Operation{
		OperationID: "update-person",
		Method:      http.MethodPut,
		Path:        "/api/people/{id}",
		Summary:     "Update a person",
		Description: "Updates the name and initials of an existing person.",
		Tags:        []string{"People"},
	}, func(ctx context.Context, input *UpdatePersonInput) (*UpdatePersonOutput, error) {
		if input.Body.Name == "" || input.Body.Initials == "" {
			return nil, huma.Error400BadRequest("name and initials are required")
		}

		exists, err := store.PersonExists(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check person existence", err)
		}
		if !exists {
			return nil, huma.Error404NotFound("person not found")
		}

		if err := store.UpdatePerson(ctx, input.ID, input.Body.Name, input.Body.Initials); err != nil {
			return nil, huma.Error500InternalServerError("failed to update person", err)
		}

		return &UpdatePersonOutput{
			Body: Person{
				ID:           input.ID,
				Name:         input.Body.Name,
				Initials:     input.Body.Initials,
				Availability: []AvailabilityEntry{},
			},
		}, nil
	})

	// DELETE /api/people/{id}
	huma.Register(api, huma.Operation{
		OperationID: "delete-person",
		Method:      http.MethodDelete,
		Path:        "/api/people/{id}",
		Summary:     "Delete a person",
		Description: "Deletes a person. Returns 409 Conflict if the person is referenced by backlog or schedule assignments.",
		Tags:        []string{"People"},
	}, func(ctx context.Context, input *DeletePersonInput) (*DeletePersonOutput, error) {
		exists, err := store.PersonExists(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check person existence", err)
		}
		if !exists {
			return nil, huma.Error404NotFound("person not found")
		}

		hasRefs, err := store.PersonHasReferences(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check person references", err)
		}
		if hasRefs {
			return nil, huma.Error409Conflict("person is referenced by backlog or schedule assignments and cannot be deleted")
		}

		if err := store.DeletePerson(ctx, input.ID); err != nil {
			if isForeignKeyViolation(err) {
				return nil, huma.Error409Conflict("person is referenced by backlog or schedule assignments and cannot be deleted")
			}
			return nil, huma.Error500InternalServerError("failed to delete person", err)
		}

		return &DeletePersonOutput{
			Body: struct {
				Message string `json:"message"`
			}{Message: "person deleted"},
		}, nil
	})

	// PUT /api/people/{id}/availability/{date}
	huma.Register(api, huma.Operation{
		OperationID: "upsert-person-availability",
		Method:      http.MethodPut,
		Path:        "/api/people/{id}/availability/{date}",
		Summary:     "Upsert person availability for a date",
		Description: "Creates or updates a single availability entry for a person on a specific date. The status must be one of: available, busy, partial, off. The date must be within the planning window.",
		Tags:        []string{"People"},
	}, func(ctx context.Context, input *UpsertAvailabilityInput) (*UpsertAvailabilityOutput, error) {
		// Validate status is canonical.
		if !canonicalStatuses[input.Body.Status] {
			return nil, huma.Error400BadRequest("status must be one of: available, busy, partial, off")
		}

		// Validate date format.
		date, err := isoToDate(input.Date)
		if err != nil {
			return nil, huma.Error400BadRequest("date must be a valid ISO 8601 date (YYYY-MM-DD)")
		}

		// Check person exists.
		exists, err := store.PersonExists(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check person existence", err)
		}
		if !exists {
			return nil, huma.Error404NotFound("person not found")
		}

		// Validate within planning window.
		pw, err := store.GetPlanningWindow(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to retrieve planning window", err)
		}
		if !isWithinPlanningWindow(input.Date, pw) {
			return nil, huma.Error400BadRequest("date is outside the planning window")
		}

		if err := store.UpsertAvailability(ctx, input.ID, pgtype.Date{Time: date, Valid: true}, input.Body.Status); err != nil {
			return nil, huma.Error500InternalServerError("failed to upsert availability", err)
		}

		return &UpsertAvailabilityOutput{
			Body: struct {
				Message string `json:"message"`
			}{Message: "availability updated"},
		}, nil
	})

	// DELETE /api/people/{id}/availability/{date}
	huma.Register(api, huma.Operation{
		OperationID: "delete-person-availability",
		Method:      http.MethodDelete,
		Path:        "/api/people/{id}/availability/{date}",
		Summary:     "Delete person availability for a date",
		Description: "Deletes a single availability entry for a person on a specific date. Idempotent: deleting an already-absent entry succeeds without error. Subsequent dashboard reads will fall back to the default-missing-availability behavior.",
		Tags:        []string{"People"},
	}, func(ctx context.Context, input *DeleteAvailabilityInput) (*DeleteAvailabilityOutput, error) {
		// Validate date format.
		date, err := isoToDate(input.Date)
		if err != nil {
			return nil, huma.Error400BadRequest("date must be a valid ISO 8601 date (YYYY-MM-DD)")
		}

		// Check person exists.
		exists, err := store.PersonExists(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to check person existence", err)
		}
		if !exists {
			return nil, huma.Error404NotFound("person not found")
		}

		// Validate within planning window.
		pw, err := store.GetPlanningWindow(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to retrieve planning window", err)
		}
		if !isWithinPlanningWindow(input.Date, pw) {
			return nil, huma.Error400BadRequest("date is outside the planning window")
		}

		if err := store.DeleteAvailability(ctx, input.ID, pgtype.Date{Time: date, Valid: true}); err != nil {
			return nil, huma.Error500InternalServerError("failed to delete availability", err)
		}

		return &DeleteAvailabilityOutput{
			Body: struct {
				Message string `json:"message"`
			}{Message: "availability deleted"},
		}, nil
	})
}
