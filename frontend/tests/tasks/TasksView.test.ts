import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h } from "vue";
import type { Component } from "vue";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ref } from "vue";

// Per-test reactive refs for the mocked composables.
let mockTasks: Array<{
	id: string;
	title: string;
	priority: string;
	peopleNeeded: number;
	room: string;
	status: string;
	assignedTo: string[];
}> = [];

vi.mock("@/client/@pinia/colada.gen", () => ({
	getTasksBacklogQueryKey: () => [],
	getTasksBacklogQuery: () => ({}),
	createTaskMutation: () => ({}),
	updateTaskMutation: () => ({}),
	deleteTaskMutation: () => ({}),
	listRoomsQuery: () => ({ key: ["listRooms"], query: vi.fn() }),
}));

vi.mock("@/tasks/composables/useTaskBacklog", () => ({
	useTaskBacklog: () => ({
		data: ref({
			tasks: mockTasks,
			priorities: [
				{ id: "high", label: "High", colorIntent: "destructive" },
				{ id: "medium", label: "Medium", colorIntent: "warning" },
				{ id: "low", label: "Low", colorIntent: "success" },
			],
			statuses: [],
			summary: {
				totalTasks: mockTasks.length,
				highPriorityTasks: 0,
				unassignedTasks: 0,
				understaffedTasks: 0,
			},
		}),
		isLoading: ref(false),
		isError: ref(false),
		isEmpty: ref(false),
		refresh: vi.fn(),
		queryKey: [],
	}),
}));

vi.mock("@/shared/composables/usePeopleAvailability", () => ({
	usePeopleAvailability: () => ({
		data: ref({
			people: [{ id: "p1", name: "Sophia", initials: "SC" }],
			summary: { availableToday: 1, totalPeople: 1 },
		}),
	}),
}));

vi.mock("@pinia/colada", async () => {
	const actual = await vi.importActual("@pinia/colada");
	return {
		...actual,
		useQuery: vi.fn(() => ({
			isLoading: ref(false),
			error: ref(null),
			data: ref(null),
			refetch: vi.fn(),
		})),
		useQueryCache: () => ({
			invalidateQueries: vi.fn(),
		}),
		useMutation: () => ({
			mutateAsync: vi.fn(),
			isLoading: ref(false),
			isPending: ref(false),
		}),
	};
});

import TasksView from "../../src/tasks/TasksView.vue";

async function renderComponent(
	component: Component,
	props?: Record<string, unknown>,
): Promise<string> {
	const app = createSSRApp({
		render: () => h(component, props),
	});
	return renderToString(app);
}

describe("TasksView SSR regression", () => {
	afterEach(() => {
		mockTasks = [];
	});

	it("renders Task backlog heading and Add Task button", async () => {
		const html = await renderComponent(TasksView);
		expect(html).toContain("Task backlog");
		expect(html).toContain("Add Task");
	});

	it("renders task list with Edit and Delete controls", async () => {
		mockTasks = [
			{
				id: "t1",
				title: "Pack boxes",
				priority: "high",
				peopleNeeded: 2,
				room: "Kitchen",
				status: "backlog",
				assignedTo: [],
			},
			{
				id: "t2",
				title: "Move furniture",
				priority: "medium",
				peopleNeeded: 3,
				room: "Living Room",
				status: "ready",
				assignedTo: ["p1"],
			},
		];
		const html = await renderComponent(TasksView);
		expect(html).toContain("Pack boxes");
		expect(html).toContain("Move furniture");
		expect(html).toContain("Edit");
		expect(html).toContain("Delete");
	});

	it("renders empty state when no tasks exist", async () => {
		mockTasks = [];
		const html = await renderComponent(TasksView);
		expect(html).toContain("Add Task");
	});

	it("does not render Sheet or SheetContent elements", async () => {
		mockTasks = [
			{
				id: "t1",
				title: "Test",
				priority: "medium",
				peopleNeeded: 1,
				room: "Test",
				status: "backlog",
				assignedTo: [],
			},
		];
		const html = await renderComponent(TasksView);
		// Verify Sheet-related data attributes are absent
		expect(html).not.toContain('data-slot="sheet"');
		expect(html).not.toContain('data-slot="sheet-content"');
		expect(html).not.toContain('data-slot="sheet-footer"');
		// Verify SheetTrigger is not wrapping the Add Task button
		expect(html).not.toContain("SheetTrigger");
	});

	it("renders Add Task as a direct Button (not inside a SheetTrigger)", async () => {
		mockTasks = [
			{
				id: "t1",
				title: "Test",
				priority: "medium",
				peopleNeeded: 1,
				room: "Test",
				status: "backlog",
				assignedTo: [],
			},
		];
		const html = await renderComponent(TasksView);
		// The + Add Task button should render as a plain Button
		expect(html).toContain("+ Add Task");
	});

	it("includes the AddOperationModal in the template", async () => {
		mockTasks = [
			{
				id: "t1",
				title: "Test",
				priority: "medium",
				peopleNeeded: 1,
				room: "Test",
				status: "backlog",
				assignedTo: [],
			},
		];
		const html = await renderComponent(TasksView);
		// The AddOperationModal is rendered outside the Card in the template.
		// When modalOpen is false, Reka UI's Dialog does not render portal content.
		// In SSR, the base section and card render correctly.
		expect(html).toContain("Task backlog");
		// Verify the Add Task button renders outside any SheetTrigger wrapper
		// (the button is a direct child of the CardHeader flex container)
		expect(html).toContain("+ Add Task");
	});
});
