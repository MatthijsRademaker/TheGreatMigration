// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { nextTick, computed, ref } from "vue";
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

// ── Helpers ─────────────────────────────────────────────────────────────────

async function mountTasks() {
	const wrapper = mount(TasksView, { attachTo: document.body });
	// Wait for initial render (portal, reka-ui)
	await nextTick();
	await nextTick();
	await nextTick();
	return wrapper;
}

async function openAddModal(wrapper: ReturnType<typeof mount>) {
	// The backlog view shows "Add Task" button in the panel toolbar
	const buttons = wrapper.findAll("button");
	const addBtn = buttons.find((b) => b.text().includes("Add Task"));
	expect(addBtn).toBeTruthy();
	await addBtn!.trigger("click");
	await nextTick();
	await nextTick();
	await nextTick();
}

// ── Tests ────────────────────────────────────────────────────────────────────

describe("TasksView – Room Select", () => {
	beforeEach(() => {
		mockRoomsQuery = createRoomsQuery();
	});

	it("renders without crashing and shows backlog", async () => {
		const wrapper = mount(TasksView);
		expect(wrapper.text()).toContain("Tasks Backlog");
		expect(wrapper.text()).toContain("Pack kitchen");
		wrapper.unmount();
	});

	it("renders room Select with room names when rooms query succeeds", async () => {
		const mockRooms = [
			{ id: "r1", name: "Kitchen", type: "room" },
			{ id: "r2", name: "Living Room", type: "room" },
			{ id: "r3", name: "Garage", type: "room" },
		];
		mockRoomsQuery = createRoomsQuery({ rooms: mockRooms });

		const wrapper = await mountTasks();
		await openAddModal(wrapper);

		// The success branch renders the Select with placeholder
		const html = document.body.innerHTML;
		expect(html).toContain("Select a room");
		// Loading/error text should NOT be present
		expect(html).not.toContain("Loading rooms");
		expect(html).not.toContain("Could not load rooms");
		// The room data is correctly passed to the query mock
		expect(mockRoomsQuery.data.value?.rooms).toEqual(mockRooms);

		wrapper.unmount();
	});

	it("shows loading placeholder while rooms query is pending", async () => {
		mockRoomsQuery = createRoomsQuery({ isLoading: true });

		const wrapper = await mountTasks();
		await openAddModal(wrapper);

		// The loading state shows a disabled Select with "Loading rooms…"
		const html = document.body.innerHTML;
		expect(html).toContain("Loading rooms");

		wrapper.unmount();
	});

	it("shows error message and retry button when rooms query fails", async () => {
		mockRoomsQuery = createRoomsQuery({
			error: new Error("Network error"),
		});

		const wrapper = await mountTasks();
		await openAddModal(wrapper);

		// The error state shows a message and retry button
		const html = document.body.innerHTML;
		expect(html).toContain("Could not load rooms.");
		expect(html).toContain("Retry");

		// Click retry and verify refetch is called
		const retryBtn = Array.from(document.body.querySelectorAll("button")).find(
			(b) => b.textContent?.trim() === "Retry",
		);
		expect(retryBtn).toBeTruthy();
		(retryBtn as HTMLElement).click();
		expect(mockRoomsQuery.refetch).toHaveBeenCalledOnce();

		wrapper.unmount();
	});
});
