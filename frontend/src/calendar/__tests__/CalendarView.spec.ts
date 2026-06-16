// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { computed, ref } from "vue";
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
});
