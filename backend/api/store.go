package api

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Store is the data access interface for the backend.
type Store interface {
	GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error)
	UpdatePlanningWindow(ctx context.Context, startDate, endDate time.Time) (*PlanningWindowBody, error)
	GetPeopleAvailability(ctx context.Context, startDate time.Time, days int, offset int, limit int) (*DashboardBody, error)
	GetTaskBacklog(ctx context.Context) (*TaskBacklogBody, error)
	GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*DailyScheduleBody, error)
	// Person CRUD
	CreatePerson(ctx context.Context, name, initials string) (string, error)
	UpdatePerson(ctx context.Context, id, name, initials string) error
	DeletePerson(ctx context.Context, id string) error
	PersonExists(ctx context.Context, id string) (bool, error)
	PersonHasReferences(ctx context.Context, id string) (bool, error)

	// Availability write
	UpsertAvailability(ctx context.Context, personID string, date pgtype.Date, status string) error
	DeleteAvailability(ctx context.Context, personID string, date pgtype.Date) error

	// Task CRUD
	CreateTask(ctx context.Context, input CreateTaskInput) (*TaskRow, error)
	UpdateTask(ctx context.Context, id string, input UpdateTaskInput) (*TaskRow, error)
	DeleteTask(ctx context.Context, id string) error
	TaskExists(ctx context.Context, id string) (bool, error)
	TaskHasScheduleCards(ctx context.Context, id string) (bool, error)

	// Schedule-card CRUD
	CreateScheduleCard(ctx context.Context, input CreateScheduleCardInput) (*TaskCard, error)
	UpdateScheduleCard(ctx context.Context, id string, input CreateScheduleCardInput) (*TaskCard, error)
	DeleteScheduleCard(ctx context.Context, id string) error

	ListRooms(ctx context.Context) ([]Room, error)
	CreateRoom(ctx context.Context, input CreateRoomInput) (*Room, error)
	UpdateRoom(ctx context.Context, id string, input UpdateRoomInput) (*Room, error)
	DeleteRoom(ctx context.Context, id string) error
}
