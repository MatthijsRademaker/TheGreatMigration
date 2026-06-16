package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Canonical task vocabularies ----------

var canonicalTaskPriorities = map[string]bool{
	"high":   true,
	"medium": true,
	"low":    true,
}

var canonicalTaskStatuses = map[string]bool{
	"backlog":  true,
	"ready":    true,
	"assigned": true,
}

// ErrTaskNotFound is returned when a task ID is not found.
var ErrTaskNotFound = errors.New("task not found")

// ---------- Domain input types ----------

// CreateTaskInput is the domain-layer input for creating a task.
type CreateTaskInput struct {
	Title        string
	Priority     string
	PeopleNeeded int
	Room         string
	Status       string
	AssignedTo   []string
}

// UpdateTaskInput is the domain-layer input for updating a task.
type UpdateTaskInput struct {
	Title        string
	Priority     string
	PeopleNeeded int
	Room         string
	Status       string
	AssignedTo   []string
}

// ---------- Request / Response types ----------

// TaskBacklogInput is empty: the endpoint accepts no query parameters.
type TaskBacklogInput struct{}

// TaskBacklogOutput is the combined payload for the task backlog endpoint.
type TaskBacklogOutput struct {
	Body TaskBacklogBody
}

// TaskBacklogBody is the top-level response body.
type TaskBacklogBody struct {
	Summary    TaskSummary        `json:"summary" doc:"Derived summary counts"`
	Tasks      []TaskRow          `json:"tasks" doc:"Backlog task rows"`
	Priorities []PriorityLegend   `json:"priorities" doc:"Canonical priority legend"`
	Statuses   []TaskStatusLegend `json:"statuses" doc:"Canonical task-status legend"`
}

// TaskSummary holds derived counts for the task backlog.
type TaskSummary struct {
	TotalTasks        int `json:"totalTasks" doc:"Total number of tasks in the backlog"`
	HighPriorityTasks int `json:"highPriorityTasks" doc:"Count of tasks with priority 'high'"`
	UnassignedTasks   int `json:"unassignedTasks" doc:"Count of tasks with empty assignedTo"`
	UnderstaffedTasks int `json:"understaffedTasks" doc:"Count of tasks with non-empty assignedTo where len(assignedTo) < peopleNeeded"`
}

// TaskRow represents a single task in the backlog.
type TaskRow struct {
	ID           string   `json:"id" doc:"Stable task identifier, prefixed 'task-'"`
	Title        string   `json:"title" doc:"Human-readable task description"`
	Priority     string   `json:"priority" doc:"One of: high, medium, low"`
	PeopleNeeded int      `json:"peopleNeeded" doc:"Number of people required for the task, minimum 1"`
	Room         string   `json:"room" doc:"Room or area the task belongs to"`
	Status       string   `json:"status" doc:"One of: backlog, ready, assigned"`
	AssignedTo   []string `json:"assignedTo" doc:"Person-ID strings for assigned helpers, may be empty"`
}

// PriorityLegend is a canonical priority definition from the design system.
type PriorityLegend struct {
	ID          string `json:"id" doc:"Priority identifier"`
	Label       string `json:"label" doc:"Human-readable label"`
	ColorIntent string `json:"colorIntent" doc:"Design system color intent"`
}

// TaskStatusLegend is a canonical task-status definition from the design system.
type TaskStatusLegend struct {
	ID          string `json:"id" doc:"Status identifier"`
	Label       string `json:"label" doc:"Human-readable label"`
	ColorIntent string `json:"colorIntent" doc:"Design system color intent"`
}

// ---------- Legends ----------

// PriorityLegendData is the canonical priority legend used by task backlog responses.
var PriorityLegendData = []PriorityLegend{
	{ID: "high", Label: "High", ColorIntent: "destructive"},
	{ID: "medium", Label: "Medium", ColorIntent: "warning"},
	{ID: "low", Label: "Low", ColorIntent: "success"},
}

// TaskStatusLegendData is the canonical task-status legend used by backlog responses.
var TaskStatusLegendData = []TaskStatusLegend{
	{ID: "backlog", Label: "Backlog", ColorIntent: "muted"},
	{ID: "ready", Label: "Ready", ColorIntent: "info"},
	{ID: "assigned", Label: "Assigned", ColorIntent: "success"},
}

// ---------- Handler ----------

func registerTasksBacklog(api huma.API, store Store) {
	huma.Register(api, huma.Operation{
		OperationID: "get-tasks-backlog",
		Method:      http.MethodGet,
		Path:        "/api/tasks/backlog",
		Summary:     "Task backlog data",
		Description: "Returns a combined payload with derived summary counts, task backlog rows, a canonical priority legend, and a canonical task-status legend. " +
			"The endpoint serves as the read-only source of truth for backlog task data used by the dashboard and /tasks route.",
		Tags: []string{"Tasks"},
	}, func(ctx context.Context, input *TaskBacklogInput) (*TaskBacklogOutput, error) {
		body, err := store.GetTaskBacklog(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to retrieve task backlog", err)
		}

		return &TaskBacklogOutput{
			Body: *body,
		}, nil
	})
}

// ---------- Task CRUD request / response types ----------

// CreateTaskRequestBody holds the fields for creating a task.
type CreateTaskRequestBody struct {
	Title        string   `json:"title" required:"true" doc:"Human-readable task description"`
	Priority     string   `json:"priority" required:"true" enum:"high,medium,low" doc:"One of: high, medium, low"`
	PeopleNeeded int      `json:"peopleNeeded" required:"true" minimum:"1" doc:"Number of people required for the task, minimum 1"`
	Room         string   `json:"room" required:"true" doc:"Room or area the task belongs to"`
	Status       string   `json:"status" required:"true" enum:"backlog,ready,assigned" doc:"One of: backlog, ready, assigned"`
	AssignedTo   []string `json:"assignedTo" doc:"Person-ID strings for assigned helpers, may be empty"`
}

// CreateTaskInputHuma is the Huma input for POST /api/tasks.
type CreateTaskInputHuma struct {
	Body CreateTaskRequestBody
}

// CreateTaskOutput is the response wrapper for POST /api/tasks.
type CreateTaskOutput struct {
	Body TaskRow
}

// UpdateTaskInputHuma is the Huma input for PUT /api/tasks/{id}.
type UpdateTaskInputHuma struct {
	ID   string                `path:"id" doc:"Task identifier"`
	Body CreateTaskRequestBody
}

// UpdateTaskOutput is the response wrapper for PUT /api/tasks/{id}.
type UpdateTaskOutput struct {
	Body TaskRow
}

// DeleteTaskInput is the Huma input for DELETE /api/tasks/{id}.
type DeleteTaskInput struct {
	ID string `path:"id" doc:"Task identifier"`
}

// DeleteTaskOutput is the empty response for DELETE /api/tasks/{id}.
type DeleteTaskOutput struct{}

// ---------- Task CRUD validation ----------

// validateTaskInput checks domain-level constraints and returns a Huma error or nil.
func validateTaskInput(body CreateTaskRequestBody, store Store, ctx context.Context) error {
	if body.Title == "" {
		return huma.Error400BadRequest("title is required")
	}
	if !canonicalTaskPriorities[body.Priority] {
		return huma.Error400BadRequest("priority must be one of: high, medium, low")
	}
	if !canonicalTaskStatuses[body.Status] {
		return huma.Error400BadRequest("status must be one of: backlog, ready, assigned")
	}
	if body.PeopleNeeded < 1 {
		return huma.Error400BadRequest("peopleNeeded must be at least 1")
	}
	if body.Room == "" {
		return huma.Error400BadRequest("room is required")
	}
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

// ---------- Task CRUD handlers ----------

func registerTasksEndpoints(api huma.API, store Store) {
	// POST /api/tasks — create a new task.
	huma.Register(api, huma.Operation{
		OperationID: "create-task",
		Method:      http.MethodPost,
		Path:        "/api/tasks",
		Summary:     "Create a backlog task",
		Description: "Creates a new backlog task with server-assigned task-* ID and append sort order. Returns 400 for invalid priority, status, empty title, peopleNeeded < 1, missing room, or unknown assigned person IDs.",
		Tags:        []string{"Tasks"},
	}, func(ctx context.Context, input *CreateTaskInputHuma) (*CreateTaskOutput, error) {
		if err := validateTaskInput(input.Body, store, ctx); err != nil {
			return nil, err
		}

		task, err := store.CreateTask(ctx, CreateTaskInput{
			Title:        input.Body.Title,
			Priority:     input.Body.Priority,
			PeopleNeeded: input.Body.PeopleNeeded,
			Room:         input.Body.Room,
			Status:       input.Body.Status,
			AssignedTo:   input.Body.AssignedTo,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create task", err)
		}

		return &CreateTaskOutput{Body: *task}, nil
	})

	// PUT /api/tasks/{id} — update an existing task.
	huma.Register(api, huma.Operation{
		OperationID: "update-task",
		Method:      http.MethodPut,
		Path:        "/api/tasks/{id}",
		Summary:     "Update a backlog task",
		Description: "Updates a backlog task and replaces assignments transactionally. Returns 400 for validation errors, 404 if the task ID is unknown.",
		Tags:        []string{"Tasks"},
	}, func(ctx context.Context, input *UpdateTaskInputHuma) (*UpdateTaskOutput, error) {
		if err := validateTaskInput(input.Body, store, ctx); err != nil {
			return nil, err
		}

		task, err := store.UpdateTask(ctx, input.ID, UpdateTaskInput{
			Title:        input.Body.Title,
			Priority:     input.Body.Priority,
			PeopleNeeded: input.Body.PeopleNeeded,
			Room:         input.Body.Room,
			Status:       input.Body.Status,
			AssignedTo:   input.Body.AssignedTo,
		})
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				return nil, huma.Error404NotFound("task not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to update task", err)
		}

		return &UpdateTaskOutput{Body: *task}, nil
	})

	// DELETE /api/tasks/{id} — delete a task.
	huma.Register(api, huma.Operation{
		OperationID: "delete-task",
		Method:      http.MethodDelete,
		Path:        "/api/tasks/{id}",
		Summary:     "Delete a backlog task",
		Description: "Deletes a backlog task and its assignments transactionally. Returns 404 if the task ID is unknown.",
		Tags:        []string{"Tasks"},
	}, func(ctx context.Context, input *DeleteTaskInput) (*DeleteTaskOutput, error) {
		err := store.DeleteTask(ctx, input.ID)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				return nil, huma.Error404NotFound("task not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to delete task", err)
		}
		return &DeleteTaskOutput{}, nil
	})
}
