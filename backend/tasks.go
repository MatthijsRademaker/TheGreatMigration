package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Request / Response types ----------

// TaskBacklogInput is empty: the endpoint accepts no query parameters.
type TaskBacklogInput struct{}

// TaskBacklogOutput is the combined payload for the task backlog endpoint.
type TaskBacklogOutput struct {
	Body TaskBacklogBody
}

// TaskBacklogBody is the top-level response body.
type TaskBacklogBody struct {
	Summary    TaskSummary       `json:"summary" doc:"Derived summary counts"`
	Tasks      []TaskRow         `json:"tasks" doc:"Backlog task rows"`
	Priorities []PriorityLegend  `json:"priorities" doc:"Canonical priority legend"`
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

var priorityLegend = []PriorityLegend{
	{ID: "high", Label: "High", ColorIntent: "destructive"},
	{ID: "medium", Label: "Medium", ColorIntent: "warning"},
	{ID: "low", Label: "Low", ColorIntent: "success"},
}

var taskStatusLegend = []TaskStatusLegend{
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
