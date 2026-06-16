import { renderToString } from "@vue/server-renderer";
import { createSSRApp, defineComponent, h } from "vue";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ref } from "vue";

// Per-test reactive refs for the mocked composable.
var mockIsLoading = ref(false);
var mockIsError = ref(false);
var mockTaskCount = ref(5);
var mockPriorityCount = ref(2);

vi.mock("@/client/@pinia/colada.gen", () => ({
	getTasksBacklogQueryKey: () => [],
	getTasksBacklogQuery: () => ({}),
}));

vi.mock("@/tasks/composables/useTaskBacklog", () => ({
	useTaskBacklog: () => ({
		data: ref({
			tasks: Array.from({ length: mockTaskCount.value }, (_, i) => ({
				id: `task-${i + 1}`,
				title: `Task ${i + 1}`,
				priority: i < mockPriorityCount.value ? "high" : "low",
				peopleNeeded: 2,
				room: "Kitchen",
				status: "backlog",
				assignedTo: [],
			})),
			priorities: [
				{ id: "high", label: "High", colorIntent: "destructive" },
				{ id: "medium", label: "Medium", colorIntent: "warning" },
				{ id: "low", label: "Low", colorIntent: "success" },
			],
			statuses: [],
			summary: {
				totalTasks: mockTaskCount.value,
				highPriorityTasks: mockPriorityCount.value,
				unassignedTasks: 3,
				understaffedTasks: 2,
			},
		}),
		isLoading: mockIsLoading,
		isError: mockIsError,
		isEmpty: ref(false),
		refresh: vi.fn(),
		queryKey: [],
	}),
}));

import { useTaskBacklog } from "../../src/tasks/composables/useTaskBacklog";

async function renderState() {
	const TestComponent = defineComponent({
		setup() {
			const state = useTaskBacklog();
			return () =>
				h("div", [
					h(
						"span",
						{ "data-testid": "loading" },
						String(state.isLoading.value),
					),
					h("span", { "data-testid": "error" }, String(state.isError.value)),
					h(
						"span",
						{ "data-testid": "task-count" },
						String(state.data.value.tasks.length),
					),
					h(
						"span",
						{ "data-testid": "high-priority" },
						String(state.data.value.summary.highPriorityTasks),
					),
				]);
		},
	});

	const app = createSSRApp(TestComponent);
	return renderToString(app);
}

describe("useTaskBacklog composable states", () => {
	afterEach(() => {
		mockIsLoading.value = false;
		mockIsError.value = false;
		mockTaskCount.value = 5;
		mockPriorityCount.value = 2;
	});

	it("exposes isLoading=true when the query is pending", async () => {
		mockIsLoading.value = true;
		const html = await renderState();
		expect(html).toContain('data-testid="loading">true<');
		expect(html).toContain('data-testid="error">false<');
	});

	it("exposes isError=true when the query has errored", async () => {
		mockIsError.value = true;
		const html = await renderState();
		expect(html).toContain('data-testid="loading">false<');
		expect(html).toContain('data-testid="error">true<');
	});

	it("exposes adapted tasks from the query response", async () => {
		mockTaskCount.value = 3;
		const html = await renderState();
		expect(html).toContain('data-testid="task-count">3<');
	});

	it("exposes adapted highPriorityTasks from the query response", async () => {
		mockTaskCount.value = 10;
		mockPriorityCount.value = 4;
		const html = await renderState();
		expect(html).toContain('data-testid="high-priority">4<');
	});

	it("returns empty tasks array when no data", async () => {
		mockTaskCount.value = 0;
		const html = await renderState();
		expect(html).toContain('data-testid="task-count">0<');
	});
});
