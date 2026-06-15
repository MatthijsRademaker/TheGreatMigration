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

// ---------- Seed data ----------

var seedTasks = []TaskRow{
	// High priority
	{ID: "task-1", Title: "Disconnect kitchen appliances", Priority: "high", PeopleNeeded: 3, Room: "Kitchen", Status: "backlog", AssignedTo: []string{}},
	{ID: "task-2", Title: "Wrap living room furniture", Priority: "high", PeopleNeeded: 2, Room: "Living Room", Status: "ready", AssignedTo: []string{}},
	{ID: "task-3", Title: "Pack kitchen fragile items", Priority: "high", PeopleNeeded: 2, Room: "Kitchen", Status: "assigned", AssignedTo: []string{"p1"}},
	{ID: "task-4", Title: "Disassemble bedroom furniture", Priority: "high", PeopleNeeded: 2, Room: "Bedroom 1", Status: "assigned", AssignedTo: []string{"p2", "p3"}},
	// Medium priority
	{ID: "task-5", Title: "Sort and label moving boxes", Priority: "medium", PeopleNeeded: 3, Room: "Living Room", Status: "backlog", AssignedTo: []string{"p4"}},
	{ID: "task-6", Title: "Clear garage shelving", Priority: "medium", PeopleNeeded: 1, Room: "Garage", Status: "ready", AssignedTo: []string{}},
	{ID: "task-7", Title: "Move bedroom wardrobe", Priority: "medium", PeopleNeeded: 3, Room: "Bedroom 2", Status: "assigned", AssignedTo: []string{"p5", "p6", "p7"}},
	// Low priority
	{ID: "task-8", Title: "Sweep garage floor", Priority: "low", PeopleNeeded: 1, Room: "Garage", Status: "backlog", AssignedTo: []string{"p8"}},
	{ID: "task-9", Title: "Dust living room shelves", Priority: "low", PeopleNeeded: 2, Room: "Living Room", Status: "ready", AssignedTo: []string{}},
	{ID: "task-10", Title: "Wipe down kitchen counters", Priority: "low", PeopleNeeded: 2, Room: "Kitchen", Status: "assigned", AssignedTo: []string{"p1", "p2"}},
	// Extra medium to ensure assignment variety
	{ID: "task-11", Title: "Inventory bedroom closet", Priority: "medium", PeopleNeeded: 3, Room: "Bedroom 1", Status: "ready", AssignedTo: []string{"p3"}},
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

func registerTasksBacklog(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-tasks-backlog",
		Method:      http.MethodGet,
		Path:        "/api/tasks/backlog",
		Summary:     "Task backlog data",
		Description: "Returns a combined payload with derived summary counts, task backlog rows, a canonical priority legend, and a canonical task-status legend. " +
			"The endpoint serves as the read-only source of truth for backlog task data used by the dashboard and /tasks route.",
		Tags: []string{"Tasks"},
	}, func(ctx context.Context, input *TaskBacklogInput) (*TaskBacklogOutput, error) {
		total := len(seedTasks)
		highPriority := 0
		unassigned := 0
		understaffed := 0

		for _, t := range seedTasks {
			if t.Priority == "high" {
				highPriority++
			}
			if len(t.AssignedTo) == 0 {
				unassigned++
			} else if len(t.AssignedTo) < t.PeopleNeeded {
				understaffed++
			}
		}

		return &TaskBacklogOutput{
			Body: TaskBacklogBody{
				Summary: TaskSummary{
					TotalTasks:        total,
					HighPriorityTasks: highPriority,
					UnassignedTasks:   unassigned,
					UnderstaffedTasks: understaffed,
				},
				Tasks:      seedTasks,
				Priorities: priorityLegend,
				Statuses:   taskStatusLegend,
			},
		}, nil
	})
}
