package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/user/the-great-migration/backend/db"
)

// Store is the data access interface for the backend.
type Store interface {
	GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error)
	GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*DashboardBody, error)
}

// PgStore implements Store backed by a pgx connection pool and sqlc queries.
type PgStore struct {
	queries *db.Queries
}

// NewPgStore creates a new PgStore from a pgx connection pool.
func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{
		queries: db.New(pool),
	}
}

// GetPlanningWindow returns the singleton planning window row from the database.
func (s *PgStore) GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error) {
	row, err := s.queries.GetPlanningWindow(ctx)
	if err != nil {
		return nil, err
	}

	startDate := pgDateToTime(row.StartDate)
	endDate := pgDateToTime(row.EndDate)
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	return &PlanningWindowBody{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		Days:      days,
	}, nil
}

// GetPeopleAvailability returns people and availability data for the given date range.
func (s *PgStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*DashboardBody, error) {
	// Fetch all people.
	peopleRows, err := s.queries.GetAllPeople(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch availability for the date range.
	startPG := pgtype.Date{Time: startDate, Valid: true}
	availRows, err := s.queries.GetAvailabilityByDateRange(ctx, db.GetAvailabilityByDateRangeParams{
		StartDate: startPG,
		Days:      int32(days),
	})
	if err != nil {
		return nil, err
	}

	// Index availability by person_id -> date -> status.
	availMap := make(map[string]map[string]string)
	for _, a := range availRows {
		if !a.Date.Valid {
			continue
		}
		dateStr := a.Date.Time.Format("2006-01-02")
		if availMap[a.PersonID] == nil {
			availMap[a.PersonID] = make(map[string]string)
		}
		availMap[a.PersonID][dateStr] = a.Status
	}

	// Build end date.
	endDate := startDate.AddDate(0, 0, days-1)

	// Build per-person availability.
	people := make([]Person, 0, len(peopleRows))
	for _, pr := range peopleRows {
		avail := make([]AvailabilityEntry, days)
		for d := 0; d < days; d++ {
			date := startDate.AddDate(0, 0, d)
			dateStr := date.Format("2006-01-02")
			status := "off" // default if no availability row
			if m, ok := availMap[pr.ID]; ok {
				if s, ok := m[dateStr]; ok {
					status = s
				}
			}
			avail[d] = AvailabilityEntry{
				Date:   dateStr,
				Status: status,
			}
		}
		people = append(people, Person{
			ID:           pr.ID,
			Name:         pr.Name,
			Initials:     pr.Initials,
			Availability: avail,
		})
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

	return &DashboardBody{
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
	}, nil
}

// pgDateToTime converts a pgtype.Date to time.Time.
// Returns zero time if the date is not valid.
func pgDateToTime(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}
