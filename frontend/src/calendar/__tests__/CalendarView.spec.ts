// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { nextTick, computed, ref } from "vue";
import CalendarView from "../CalendarView.vue";

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

// ── Helpers ─────────────────────────────────────────────────────────────────

async function mountCalendar() {
	const wrapper = mount(CalendarView, { attachTo: document.body });
	// Wait for initial render (portal, reka-ui)
	await nextTick();
	await nextTick();
	await nextTick();
	return wrapper;
}

async function openAddModal(wrapper: ReturnType<typeof mount>) {
	// The empty state shows "Add your first task" button
	const addBtn = wrapper.find("button");
	expect(addBtn.exists()).toBe(true);
	await addBtn.trigger("click");
	await nextTick();
	await nextTick();
	await nextTick();
}

// ── Tests ────────────────────────────────────────────────────────────────────

describe("CalendarView – Room Select", () => {
	beforeEach(() => {
		mockRoomsQuery = createRoomsQuery();
	});

	it("renders without crashing and shows empty state", async () => {
		const wrapper = mount(CalendarView);
		// The component renders the empty state with the add button
		expect(wrapper.text()).toContain("Add your first task");
		wrapper.unmount();
	});

	it("renders room Select with room names when rooms query succeeds", async () => {
		const mockRooms = [
			{ id: "r1", name: "Kitchen", type: "room" },
			{ id: "r2", name: "Living Room", type: "room" },
			{ id: "r3", name: "Garage", type: "room" },
		];
		mockRoomsQuery = createRoomsQuery({ rooms: mockRooms });

		const wrapper = await mountCalendar();
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

		const wrapper = await mountCalendar();
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

		const wrapper = await mountCalendar();
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

describe("CalendarView – DatePicker", () => {
	beforeEach(() => {
		mockRoomsQuery = createRoomsQuery();
	});

	it("renders DatePicker instead of a free-form date Input", async () => {
		const wrapper = await mountCalendar();
		await openAddModal(wrapper);

		// The DatePicker renders a button with "Select date" placeholder and a CalendarIcon
		const html = document.body.innerHTML;
		expect(html).toContain("Select date");

		// A plain text Input for the date should NOT be present in the modal body
		expect(html).not.toContain('id="form-date"');

		wrapper.unmount();
	});
});
