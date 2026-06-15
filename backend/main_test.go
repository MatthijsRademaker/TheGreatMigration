package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func newTestAPI() (chi.Router, huma.API) {
	router := chi.NewMux()
	config := huma.DefaultConfig("Test API", "1.0.0")
	api := humachi.New(router, config)

	huma.Register(api, huma.Operation{
		OperationID: "get-hello",
		Method:      http.MethodGet,
		Path:        "/api/hello",
	}, func(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
		resp := &HelloOutput{}
		resp.Body.Message = "Hello from the backend!"
		return resp, nil
	})

	registerDashboardPeopleAvailability(api)

	return router, api
}

func TestHelloEndpoint(t *testing.T) {
	router, _ := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	if body["message"] != "Hello from the backend!" {
		t.Fatalf("expected message 'Hello from the backend!', got %q", body["message"])
	}
}

func TestDashboardPeopleAvailability(t *testing.T) {
	router, _ := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/people-availability", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Happy path: 200 OK and JSON content-type.
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d\nbody: %s", rec.Code, rec.Body.String())
	}
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body DashboardBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, rec.Body.String())
	}

	// Response shape assertions.
	if body.Range.StartDate == "" {
		t.Fatal("range.startDate is empty")
	}
	if body.Range.EndDate == "" {
		t.Fatal("range.endDate is empty")
	}
	if body.Range.Days != 4 {
		t.Fatalf("expected range.days=4 (default), got %d", body.Range.Days)
	}
	if body.Range.SelectedDate != body.Range.StartDate {
		t.Fatalf("expected selectedDate=%s to equal startDate", body.Range.StartDate)
	}

	if body.Summary.TotalPeople == 0 {
		t.Fatal("summary.totalPeople is 0")
	}
	if len(body.People) != body.Summary.TotalPeople {
		t.Fatalf("len(people)=%d != summary.totalPeople=%d", len(body.People), body.Summary.TotalPeople)
	}

	if len(body.Statuses) != 4 {
		t.Fatalf("expected 4 statuses in legend, got %d", len(body.Statuses))
	}

	// Validate all statuses in the legend.
	canonical := map[string]bool{}
	for _, s := range body.Statuses {
		canonical[s.ID] = true
	}
	expectedStatuses := []string{"available", "busy", "partial", "off"}
	for _, s := range expectedStatuses {
		if !canonical[s] {
			t.Fatalf("statuses legend missing %q", s)
		}
	}

	// Validate all availability statuses are canonical.
	for _, p := range body.People {
		for _, e := range p.Availability {
			if !canonical[e.Status] {
				t.Fatalf("person %s has non-canonical status %q on date %s", p.ID, e.Status, e.Date)
			}
		}
		if p.ID == "" {
			t.Fatal("person has empty id")
		}
		if p.Name == "" {
			t.Fatalf("person %s has empty name", p.ID)
		}
		if p.Initials == "" {
			t.Fatalf("person %s has empty initials", p.ID)
		}
		// Each person should have exactly D availability entries.
		if len(p.Availability) != body.Range.Days {
			t.Fatalf("person %s has %d availability entries, expected %d", p.ID, len(p.Availability), body.Range.Days)
		}
	}

	// Verify summary.availableToday matches actual count on selectedDate.
	actualAvailable := 0
	for _, p := range body.People {
		for _, e := range p.Availability {
			if e.Date == body.Range.SelectedDate && e.Status == "available" {
				actualAvailable++
				break
			}
		}
	}
	if body.Summary.AvailableToday != actualAvailable {
		t.Fatalf("summary.availableToday=%d does not match actual count=%d", body.Summary.AvailableToday, actualAvailable)
	}

	// Verify at least 8 people.
	if len(body.People) < 8 {
		t.Fatalf("expected at least 8 people, got %d", len(body.People))
	}
}
