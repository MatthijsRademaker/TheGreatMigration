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

function createDailyScheduleMock(
	overrides: Partial<{
		isLoading: boolean;
		isError: boolean;
		isEmpty: boolean;
		days: unknown[];
		page: number;
		totalPages: number;
		dateRangeLabel: string;
	}> = {},
) {
	const state = {
		isLoading: overrides.isLoading ?? false,
		isError: overrides.isError ?? false,
		isEmpty: overrides.isEmpty ?? true,
		days: overrides.days ?? [],
		page: overrides.page ?? 0,
		totalPages: overrides.totalPages ?? 0,
		dateRangeLabel: overrides.dateRangeLabel ?? "",
	};
	return {
		data: computed(() => ({ days: state.days })),
		isLoading: ref(state.isLoading),
		isError: ref(state.isError),
		isEmpty: ref(state.isEmpty),
		page: ref(state.page),
		totalPages: ref(state.totalPages),
		dateRangeLabel: computed(() => state.dateRangeLabel),
		goToPrevPage: vi.fn(),
		goToNextPage: vi.fn(),
		queryKey: ["daily-schedule"],
	};
}

let mockBacklog = createBacklogMock();
let mockDailySchedule = createDailyScheduleMock();

// Capturing mutation spies — recreated before each test so assertions are clean.
let createMutateAsync: ReturnType<typeof vi.fn>;
let updateMutateAsync: ReturnType<typeof vi.fn>;
let deleteMutateAsync: ReturnType<typeof vi.fn>;

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
	useDailySchedule: vi.fn(() => mockDailySchedule),
}));

vi.mock("@/shared/composables/usePeopleAvailability", () => ({
	usePeopleAvailability: vi.fn(() => ({
		data: computed(() => ({ people: [] })),
	})),
}));

vi.mock("@/tasks/composables/useTaskBacklog", () => ({
	useTaskBacklog: vi.fn(() => mockBacklog),
}));

// We need distinct spies per mutation. Patch the useMutation mock after imports.
import { useMutation } from "@pinia/colada";

function setupMutationSpies() {
	createMutateAsync = vi.fn().mockResolvedValue(undefined);
	updateMutateAsync = vi.fn().mockResolvedValue(undefined);
	deleteMutateAsync = vi.fn().mockResolvedValue(undefined);
	// The component creates three useMutation calls in order: createMut, updateMut, deleteMut.
	const callOrder: (() => {
		mutateAsync: typeof createMutateAsync;
		isLoading: { value: boolean };
	})[] = [
		() => ({ mutateAsync: createMutateAsync, isLoading: { value: false } }),
		() => ({ mutateAsync: updateMutateAsync, isLoading: { value: false } }),
		() => ({ mutateAsync: deleteMutateAsync, isLoading: { value: false } }),
	];
	let callIndex = 0;
	(useMutation as ReturnType<typeof vi.fn>).mockImplementation(() => {
		const result = callOrder[callIndex % callOrder.length]();
		callIndex++;
		return result;
	});
}

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
		mockDailySchedule = createDailyScheduleMock();
		setupMutationSpies();
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

	it("sends taskId in create mutation body when a task is selected", async () => {
		mockBacklog = createBacklogMock({
			isEmpty: false,
			tasks: [
				{
					id: "task-1",
					title: "Pack kitchen",
					priority: "high",
					room: "Kitchen",
					peopleNeeded: 2,
				},
			],
		});

		const wrapper = await mountCalendar();
		await openAddModal(wrapper);

		// Set formTaskId directly via the component's exposed refs.
		// wrapper.vm auto-unwraps Vue refs, so we set the value directly.
		(wrapper.vm as unknown as Record<string, unknown>).formTaskId = "task-1";
		await nextTick();

		// Click Save
		const saveBtn = Array.from(document.body.querySelectorAll("button")).find(
			(b) => b.textContent?.trim() === "Save",
		);
		expect(saveBtn).toBeTruthy();
		(saveBtn as HTMLElement).click();
		await nextTick();
		await nextTick();

		expect(createMutateAsync).toHaveBeenCalledTimes(1);
		const callArg = createMutateAsync.mock.calls[0][0];
		expect(callArg.body.taskId).toBe("task-1");

		wrapper.unmount();
	});
});

describe("CalendarView – Edit Modal", () => {
	beforeEach(() => {
		mockBacklog = createBacklogMock();
		mockDailySchedule = createDailyScheduleMock();
		setupMutationSpies();
	});

	it("shows read-only task info when editing a referenced card (has taskId)", async () => {
		mockBacklog = createBacklogMock({
			isEmpty: false,
			tasks: [
				{
					id: "task-1",
					title: "Pack kitchen",
					priority: "high",
					room: "Kitchen",
					peopleNeeded: 2,
				},
			],
		});
		mockDailySchedule = createDailyScheduleMock({
			isEmpty: false,
			days: [
				{
					date: "2026-07-10",
					cards: [
						{
							id: "card-1",
							title: "Pack kitchen",
							priority: "high",
							roomArea: "Kitchen",
							peopleNeeded: 2,
							scheduledDate: "2026-07-10",
							taskId: "task-1",
							staffingStatus: "underStaffed",
							assignedCount: 0,
							assignedPeople: [],
						},
					],
				},
			],
		});

		const wrapper = await mountCalendar();

		// Open the edit modal programmatically by calling the component's openEdit.
		const anyVM = wrapper.vm as unknown as Record<string, unknown>;
		const openEditFn = anyVM["openEdit"] as
			| ((card: Record<string, unknown>) => void)
			| undefined;
		if (openEditFn) {
			openEditFn({
				id: "card-1",
				title: "Pack kitchen",
				priority: "high",
				roomArea: "Kitchen",
				peopleNeeded: 2,
				scheduledDate: "2026-07-10",
				assignedPeople: [],
				taskId: "task-1",
			});
		} else {
			// Fallback: trigger edit via DailySchedule by finding an action button
			const allBtns = wrapper.findAll("button");
			// Click the first non-"Add" button
			for (const btn of allBtns) {
				const text = btn.text();
				if (text !== "Add your first task" && text !== "Save") {
					await btn.trigger("click");
					break;
				}
			}
		}
		await nextTick();
		await nextTick();
		await nextTick();

		const html = document.body.innerHTML;
		expect(html).toContain("From backlog");
		expect(html).toContain("Pack kitchen");

		wrapper.unmount();
	});

	it("shows free-form fields when editing an unreferenced card (no taskId)", async () => {
		mockBacklog = createBacklogMock({
			isEmpty: false,
			tasks: [],
		});
		mockDailySchedule = createDailyScheduleMock({
			isEmpty: false,
			days: [
				{
					date: "2026-07-10",
					cards: [
						{
							id: "card-2",
							title: "Legacy card",
							priority: "low",
							roomArea: "Garage",
							peopleNeeded: 3,
							scheduledDate: "2026-07-10",
							taskId: null,
							staffingStatus: "underStaffed",
							assignedCount: 0,
							assignedPeople: [],
						},
					],
				},
			],
		});

		const wrapper = await mountCalendar();

		// Call openEdit programmatically
		const anyVM = wrapper.vm as unknown as Record<string, unknown>;
		const openEditFn = anyVM["openEdit"] as
			| ((card: Record<string, unknown>) => void)
			| undefined;
		if (openEditFn) {
			openEditFn({
				id: "card-2",
				title: "Legacy card",
				priority: "low",
				roomArea: "Garage",
				peopleNeeded: 3,
				scheduledDate: "2026-07-10",
				assignedPeople: [],
				taskId: null,
			});
		}
		await nextTick();
		await nextTick();
		await nextTick();

		const html = document.body.innerHTML;
		expect(html).toContain("form-title");
		expect(html).toContain("form-priority");
		expect(html).toContain("form-room");
		expect(html).toContain("form-people");
		expect(html).not.toContain("From backlog");

		wrapper.unmount();
	});
});

describe("CalendarView – Pagination Integration", () => {
	beforeEach(() => {
		mockBacklog = createBacklogMock();
		mockDailySchedule = createDailyScheduleMock();
		setupMutationSpies();
	});

	it("passes page and totalPages to DailySchedule", async () => {
		mockDailySchedule = createDailyScheduleMock({
			isEmpty: false,
			days: [
				{
					date: "2026-07-10",
					label: "10 Jul (Fri)",
					availablePeopleCount: 4,
					tasks: [],
				},
			],
			page: 2,
			totalPages: 5,
			dateRangeLabel: "10 Jul (Fri) – 13 Jul (Mon)",
		});

		const wrapper = await mountCalendar();

		const html = document.body.innerHTML;
		expect(html).toContain("Page 2 of 5");
		expect(html).toContain("Previous");
		expect(html).toContain("Next");
		expect(html).toContain("10 Jul (Fri) – 13 Jul (Mon)");

		wrapper.unmount();
	});

	it("does not show pagination bar when page is 0", async () => {
		mockDailySchedule = createDailyScheduleMock({
			isEmpty: false,
			days: [
				{
					date: "2026-07-10",
					label: "10 Jul (Fri)",
					availablePeopleCount: 4,
					tasks: [],
				},
			],
			page: 0,
			totalPages: 0,
		});

		const wrapper = await mountCalendar();

		const html = document.body.innerHTML;
		expect(html).not.toContain("Previous");
		expect(html).not.toContain("Next");

		wrapper.unmount();
	});

	it("disables Previous button on page 1", async () => {
		mockDailySchedule = createDailyScheduleMock({
			isEmpty: false,
			days: [
				{
					date: "2026-07-10",
					label: "10 Jul (Fri)",
					availablePeopleCount: 4,
					tasks: [],
				},
			],
			page: 1,
			totalPages: 3,
		});

		const wrapper = await mountCalendar();

		const html = document.body.innerHTML;
		const allButtons = html.match(/<button[\s\S]*?<\/button>/g) ?? [];
		const prevButton = allButtons.find((b) => b.includes("Previous"));
		expect(prevButton).toBeDefined();
		expect(prevButton).toContain("disabled");

		wrapper.unmount();
	});

	it("disables Next button on last page", async () => {
		mockDailySchedule = createDailyScheduleMock({
			isEmpty: false,
			days: [
				{
					date: "2026-07-10",
					label: "10 Jul (Fri)",
					availablePeopleCount: 4,
					tasks: [],
				},
			],
			page: 3,
			totalPages: 3,
		});

		const wrapper = await mountCalendar();

		const html = document.body.innerHTML;
		const allButtons = html.match(/<button[\s\S]*?<\/button>/g) ?? [];
		const nextButton = allButtons.find((b) => b.includes("Next"));
		expect(nextButton).toBeDefined();
		expect(nextButton).toContain("disabled");

		wrapper.unmount();
	});
});
