package main

import (
	"context"
	"time"
)

// seedPeopleDefs is the canonical set of people IDs, names, and initials.
// It mirrors the seed data in migrations/002_seed_demo_data.sql.
var seedPeopleDefs = []struct {
	ID       string
	Name     string
	Initials string
	Status   func(int) string
}{
	{ID: "p1", Name: "Sophia Chen", Initials: "SC", Status: always("available")},
	{ID: "p2", Name: "Marcus Rivera", Initials: "MR", Status: always("available")},
	{ID: "p3", Name: "Elena Kowalski", Initials: "EK", Status: always("available")},
	{ID: "p4", Name: "James Okafor", Initials: "JO", Status: always("available")},
	{ID: "p5", Name: "Priya Nair", Initials: "PN", Status: always("available")},
	{ID: "p6", Name: "Thomas Berg", Initials: "TB", Status: always("available")},
	{ID: "p7", Name: "Amara Diallo", Initials: "AD", Status: always("busy")},
	{ID: "p8", Name: "Noah Larsson", Initials: "NL", Status: cycleStatuses},
}

// mockStore implements Store for fast unit tests without a database.
type mockStore struct {
	planningWindow *PlanningWindowBody
	dashboard      *DashboardBody
}

func newMockStore() *mockStore {
	startDate, _ := time.Parse("2006-01-02", "2026-07-05")
	endDate, _ := time.Parse("2006-01-02", "2026-08-13")
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	return &mockStore{
		planningWindow: &PlanningWindowBody{
			StartDate: "2026-07-05",
			EndDate:   "2026-08-13",
			Days:      days,
		},
		dashboard: &DashboardBody{
			Range: Range{
				StartDate:    "2026-07-05",
				EndDate:      "2026-07-08",
				Days:         4,
				SelectedDate: "2026-07-05",
			},
			Summary: Summary{
				AvailableToday: 6,
				TotalPeople:    8,
			},
			People:   buildMockPeople(startDate, 4, seedPeopleDefs),
			Statuses: statusLegend,
		},
	}
}

func (m *mockStore) GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error) {
	return m.planningWindow, nil
}

func (m *mockStore) GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*DashboardBody, error) {
	endDate := startDate.AddDate(0, 0, days-1)
	selectedDate := startDate.Format("2006-01-02")
	people := buildMockPeople(startDate, days, seedPeopleDefs)

	// Compute summary.
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

func buildMockPeople(startDate time.Time, days int, mp []struct {
	ID, Name, Initials string
	Status             func(int) string
}) []Person {
	people := make([]Person, len(mp))
	for i, sp := range mp {
		avail := make([]AvailabilityEntry, days)
		for d := 0; d < days; d++ {
			date := startDate.AddDate(0, 0, d)
			avail[d] = AvailabilityEntry{
				Date:   date.Format("2006-01-02"),
				Status: sp.Status(d),
			}
		}
		people[i] = Person{
			ID:           sp.ID,
			Name:         sp.Name,
			Initials:     sp.Initials,
			Availability: avail,
		}
	}
	return people
}
