// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { computed, ref } from "vue";
import TasksView from "../TasksView.vue";

// ── Mock helpers ────────────────────────────────────────────────────────────

function createRoomsQuery(
	overrides: Partial<{
		isLoading: boolean;
		error: Error | null;
		rooms: { id: string; name: string; type: string }[] | null;
	}> = {},
) {
	const state = {
		isLoading: overrides.isLoading ?? false,
		error: overrides.error ?? null,
		rooms: overrides.rooms,
	};
	return {
		isLoading: ref(state.isLoading),
		error: ref(state.error),
		data: ref(state.rooms ? { rooms: state.rooms } : null),
		refetch: vi.fn(),
	};
}

// ── Tasks backlog mock ──────────────────────────────────────────────────────

const mockBacklogData = computed(() => ({
	tasks: [
		{
			id: "task-1",
			title: "Pack kitchen",
			priority: "high",
			room: "Kitchen",
			peopleNeeded: 2,
			status: "backlog",
			assignedTo: [],
		},
	],
	summary: {
		totalTasks: 1,
		highPriorityTasks: 1,
		unassignedTasks: 1,
		understaffedTasks: 1,
	},
	priorities: [],
	statuses: [],
}));

// ── Module mocks ────────────────────────────────────────────────────────────

let mockRoomsQuery = createRoomsQuery();

vi.mock("@pinia/colada", () => ({
	useQuery: vi.fn(() => mockRoomsQuery),
	useMutation: vi.fn(() => ({
		mutateAsync: vi.fn(),
		isLoading: { value: false },
	})),
	useQueryCache: vi.fn(() => ({
		invalidateQueries: vi.fn(),
	})),
}));

vi.mock("@/client/@pinia/colada.gen", () => ({
	listRoomsQuery: vi.fn(() => ({ key: ["listRooms"], query: vi.fn() })),
	createTaskMutation: vi.fn(() => ({ mutation: vi.fn() })),
	updateTaskMutation: vi.fn(() => ({ mutation: vi.fn() })),
	deleteTaskMutation: vi.fn(() => ({ mutation: vi.fn() })),
}));

vi.mock("@/tasks/composables/useTaskBacklog", () => ({
	useTaskBacklog: vi.fn(() => ({
		data: mockBacklogData,
		isLoading: ref(false),
		isError: ref(false),
		isEmpty: ref(false),
		queryKey: ["tasks-backlog"],
	})),
}));

vi.mock("@/shared/composables/usePeopleAvailability", () => ({
	usePeopleAvailability: vi.fn(() => ({
		data: computed(() => ({ people: [] })),
	})),
}));

// ── Tests ────────────────────────────────────────────────────────────────────

describe("TasksView – Room Select", () => {
	beforeEach(() => {
		mockRoomsQuery = createRoomsQuery();
	});

	it("renders without crashing and shows backlog", async () => {
		const wrapper = mount(TasksView);
		expect(wrapper.text()).toContain("Task backlog");
		expect(wrapper.text()).toContain("Pack kitchen");
		wrapper.unmount();
	});
});
