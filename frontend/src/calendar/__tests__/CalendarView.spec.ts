// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { nextTick, computed, ref } from "vue";
import CalendarView from "../CalendarView.vue";

// ── Mock helpers ────────────────────────────────────────────────────────────

function createBacklogMock(
	overrides: Partial<{
		isLoading: boolean;
		isEmpty: boolean;
		tasks: {
			id: string;
			title: string;
			priority: string;
			room: string;
			peopleNeeded: number;
		}[];
	}> = {},
) {
	const state = {
		isLoading: overrides.isLoading ?? false,
		isEmpty: overrides.isEmpty ?? true,
		tasks: overrides.tasks ?? [],
	};
	return {
		data: computed(() => ({
			tasks: state.tasks,
			priorities: [],
			statuses: [],
			summary: {
				totalTasks: state.tasks.length,
				highPriorityTasks: 0,
				unassignedTasks: 0,
				understaffedTasks: 0,
			},
		})),
		isLoading: ref(state.isLoading),
		isEmpty: ref(state.isEmpty),
		refresh: vi.fn(),
		queryKey: ["task-backlog"],
	};
}

let mockBacklog = createBacklogMock();

// ── Module mocks ────────────────────────────────────────────────────────────

vi.mock("@pinia/colada", () => ({
	useQuery: vi.fn(() => ({
		isLoading: ref(false),
		error: ref(null),
		data: ref(null),
		refetch: vi.fn(),
		isPending: ref(false),
	})),
	useMutation: vi.fn(() => ({
		mutateAsync: vi.fn(),
		isLoading: { value: false },
	})),
	useQueryCache: vi.fn(() => ({
		invalidateQueries: vi.fn(),
	})),
}));

vi.mock("@/client/@pinia/colada.gen", () => ({
	getTasksBacklogQuery: vi.fn(() => ({
		key: ["task-backlog"],
		query: vi.fn(),
	})),
	createScheduleCardMutation: vi.fn(() => ({ mutation: vi.fn() })),
	updateScheduleCardMutation: vi.fn(() => ({ mutation: vi.fn() })),
	deleteScheduleCardMutation: vi.fn(() => ({ mutation: vi.fn() })),
}));

vi.mock("@/calendar/composables/useDailySchedule", () => ({
	useDailySchedule: vi.fn(() => ({
		data: computed(() => ({ days: [] })),
		isLoading: ref(false),
		isError: ref(false),
		isEmpty: ref(true),
		queryKey: ["daily-schedule"],
	})),
}));

vi.mock("@/shared/composables/usePeopleAvailability", () => ({
	usePeopleAvailability: vi.fn(() => ({
		data: computed(() => ({ people: [] })),
	})),
}));

vi.mock("@/tasks/composables/useTaskBacklog", () => ({
	useTaskBacklog: vi.fn(() => mockBacklog),
}));

// ── Helpers ─────────────────────────────────────────────────────────────────

async function mountCalendar() {
	const wrapper = mount(CalendarView, { attachTo: document.body });
	await nextTick();
	await nextTick();
	await nextTick();
	return wrapper;
}

async function openAddModal(wrapper: ReturnType<typeof mount>) {
	const addBtn = wrapper.find("button");
	expect(addBtn.exists()).toBe(true);
	await addBtn.trigger("click");
	await nextTick();
	await nextTick();
	await nextTick();
}

// ── Tests ────────────────────────────────────────────────────────────────────

describe("CalendarView – Task Selector", () => {
	beforeEach(() => {
		mockBacklog = createBacklogMock();
	});

	it("renders without crashing and shows empty state", async () => {
		const wrapper = mount(CalendarView);
		expect(wrapper.text()).toContain("Add your first task");
		wrapper.unmount();
	});

	it("shows empty backlog guidance when no tasks exist", async () => {
		mockBacklog = createBacklogMock({ isEmpty: true, tasks: [] });

		const wrapper = await mountCalendar();
		await openAddModal(wrapper);

		const html = document.body.innerHTML;
		expect(html).toContain("No tasks in the backlog yet");
		expect(html).toContain("Task Panel");

		wrapper.unmount();
	});

	it("shows task selector with backlog tasks when tasks exist", async () => {
		const mockTasks = [
			{
				id: "task-1",
				title: "Pack kitchen",
				priority: "high",
				room: "Kitchen",
				peopleNeeded: 2,
			},
			{
				id: "task-2",
				title: "Wrap furniture",
				priority: "medium",
				room: "Living Room",
				peopleNeeded: 3,
			},
		];
		mockBacklog = createBacklogMock({
			isEmpty: false,
			tasks: mockTasks,
		});

		const wrapper = await mountCalendar();
		await openAddModal(wrapper);

		const html = document.body.innerHTML;
		expect(html).toContain("Select a task");
		expect(html).toContain("Search task");

		wrapper.unmount();
	});

	it("shows DatePicker for the scheduled date", async () => {
		mockBacklog = createBacklogMock();

		const wrapper = await mountCalendar();
		await openAddModal(wrapper);

		const html = document.body.innerHTML;
		expect(html).toContain("Select date");

		wrapper.unmount();
	});
});
