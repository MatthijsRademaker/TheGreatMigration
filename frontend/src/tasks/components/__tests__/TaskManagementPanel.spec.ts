// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { computed, ref } from "vue";
import TaskManagementPanel from "../TaskManagementPanel.vue";

// ── Default state factory ──────────────────────────────────────────────────

function defaultBacklogState() {
	return {
		tasks: [] as any[],
		summary: {
			totalTasks: 0,
			highPriorityTasks: 0,
			unassignedTasks: 0,
			understaffedTasks: 0,
		},
		priorities: [] as any[],
		statuses: [] as any[],
	};
}

// ── Mock state (reassigned in beforeEach) ──────────────────────────────────

let mockData = computed(() => defaultBacklogState());
let mockIsLoading = ref(false);
let mockIsError = ref(false);
let mockIsEmpty = ref(false);

vi.mock("@/tasks/composables/useTaskBacklog", () => ({
	useTaskBacklog: vi.fn(() => ({
		data: mockData,
		isLoading: mockIsLoading,
		isError: mockIsError,
		isEmpty: mockIsEmpty,
		queryKey: ["tasks-backlog"],
		refresh: vi.fn(),
	})),
}));

// ── Helpers ─────────────────────────────────────────────────────────────────

function makeTask(
	overrides: Partial<{
		id: string;
		title: string;
		priority: "high" | "medium" | "low";
		peopleNeeded: number;
		room: string;
		status: "backlog" | "ready" | "assigned";
		assignedTo: string[];
	}> = {},
) {
	return {
		id: overrides.id ?? "task-1",
		title: overrides.title ?? "Pack kitchen boxes",
		priority: overrides.priority ?? "high",
		peopleNeeded: overrides.peopleNeeded ?? 3,
		room: overrides.room ?? "Kitchen",
		status: overrides.status ?? "ready",
		assignedTo: overrides.assignedTo ?? [],
	};
}

function setBacklogState(state: {
	tasks?: any[];
	summary?: Partial<{
		totalTasks: number;
		highPriorityTasks: number;
		unassignedTasks: number;
		understaffedTasks: number;
	}>;
	isLoading?: boolean;
	isError?: boolean;
	isEmpty?: boolean;
}) {
	const summary = {
		totalTasks: state.summary?.totalTasks ?? state.tasks?.length ?? 0,
		highPriorityTasks: state.summary?.highPriorityTasks ?? 0,
		unassignedTasks: state.summary?.unassignedTasks ?? 0,
		understaffedTasks: state.summary?.understaffedTasks ?? 0,
	};
	mockData = computed(() => ({
		tasks: state.tasks ?? [],
		summary,
		priorities: [],
		statuses: [],
	}));
	mockIsLoading = ref(state.isLoading ?? false);
	mockIsError = ref(state.isError ?? false);
	mockIsEmpty = ref(state.isEmpty ?? false);
}

// ── Tests ───────────────────────────────────────────────────────────────────

describe("TaskManagementPanel", () => {
	beforeEach(() => {
		setBacklogState({});
	});

	// 6.2
	it("renders heading 'Tasks Backlog'", () => {
		setBacklogState({
			tasks: [makeTask()],
			summary: { totalTasks: 1 },
		});
		const wrapper = mount(TaskManagementPanel);
		expect(wrapper.text()).toContain("Tasks Backlog");
	});

	// 6.3
	it("renders count chip with totalTasks value", () => {
		setBacklogState({
			tasks: [makeTask(), makeTask({ id: "task-2" })],
			summary: { totalTasks: 5 },
		});
		const wrapper = mount(TaskManagementPanel);
		expect(wrapper.text()).toContain("5 tasks");
	});

	// 6.4
	it("hides Filter and Add Task toolbar buttons in readOnly mode", () => {
		setBacklogState({
			tasks: [makeTask()],
			summary: { totalTasks: 1 },
		});
		const wrapper = mount(TaskManagementPanel, {
			props: { readOnly: true },
		});
		expect(wrapper.text()).not.toContain("Filter");
		expect(wrapper.text()).not.toContain("Add Task");
	});

	// 6.5
	it("shows Filter and Add Task toolbar buttons in non-readOnly mode", () => {
		setBacklogState({
			tasks: [makeTask()],
			summary: { totalTasks: 1 },
		});
		const wrapper = mount(TaskManagementPanel, {
			props: { readOnly: false },
		});
		expect(wrapper.text()).toContain("Filter");
		expect(wrapper.text()).toContain("Add Task");
	});

	// 6.6
	it("hides per-row Edit and Delete buttons in readOnly mode", () => {
		setBacklogState({
			tasks: [makeTask()],
			summary: { totalTasks: 1 },
		});
		const wrapper = mount(TaskManagementPanel, {
			props: { readOnly: true },
		});
		expect(wrapper.text()).not.toContain("Edit");
		expect(wrapper.text()).not.toContain("Delete");
	});

	// 6.7
	it("shows per-row Edit and Delete buttons in non-readOnly mode", () => {
		setBacklogState({
			tasks: [makeTask()],
			summary: { totalTasks: 1 },
		});
		const wrapper = mount(TaskManagementPanel, {
			props: { readOnly: false },
		});
		expect(wrapper.text()).toContain("Edit");
		expect(wrapper.text()).toContain("Delete");
	});

	// 6.8
	it("renders a shimmer skeleton while loading", () => {
		setBacklogState({ isLoading: true });
		const wrapper = mount(TaskManagementPanel);
		expect(wrapper.find("[data-testid='skeleton-rows']").exists()).toBe(true);
		// Loading state should not render the table or footer
		expect(wrapper.text()).not.toContain("Priority:");
	});

	it("renders error state text", () => {
		setBacklogState({ isError: true });
		const wrapper = mount(TaskManagementPanel);
		expect(wrapper.text()).toContain("Could not load tasks. Please try again.");
		// Error state should not render the table or footer
		expect(wrapper.text()).not.toContain("Priority:");
	});

	it("renders empty state text", () => {
		setBacklogState({ isEmpty: true, tasks: [] });
		const wrapper = mount(TaskManagementPanel);
		expect(wrapper.text()).toContain("No tasks yet.");
		// Empty state should not render the table or footer
		expect(wrapper.text()).not.toContain("Priority:");
	});

	// 6.9
	it("renders five column headers: Task, Priority, People Needed, Room / Area, Status", () => {
		setBacklogState({
			tasks: [makeTask()],
			summary: { totalTasks: 1 },
		});
		const wrapper = mount(TaskManagementPanel);
		const headers = wrapper.findAll("th");
		expect(headers).toHaveLength(5);
		expect(headers[0].text()).toBe("Task");
		expect(headers[1].text()).toBe("Priority");
		expect(headers[2].text()).toBe("People Needed");
		expect(headers[3].text()).toBe("Room / Area");
		expect(headers[4].text()).toBe("Status");
	});

	// 6.10
	it("does not render priority legend footer", () => {
		setBacklogState({
			tasks: [makeTask()],
			summary: { totalTasks: 1 },
		});
		const wrapper = mount(TaskManagementPanel);
		expect(wrapper.text()).not.toContain("Priority:");
	});
});
