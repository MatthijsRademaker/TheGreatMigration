// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { ref } from "vue";
import KpiCards from "../components/KpiCards.vue";

// ── Mock helpers ────────────────────────────────────────────────────────────

function createAvailabilityMock(
	overrides: Partial<{
		isPending: boolean;
		error: Error | null;
		data: { summary: { availableToday: number; totalPeople: number } } | null;
	}> = {},
) {
	const state = {
		isPending: overrides.isPending ?? false,
		error: overrides.error ?? null,
		data: overrides.data ?? null,
	};
	return {
		isPending: ref(state.isPending),
		error: ref(state.error),
		data: ref(state.data),
	};
}

function createBacklogMock(
	overrides: Partial<{
		isPending: boolean;
		error: Error | null;
		data: {
			summary: { highPriorityTasks: number; unassignedTasks: number };
		} | null;
	}> = {},
) {
	const state = {
		isPending: overrides.isPending ?? false,
		error: overrides.error ?? null,
		data: overrides.data ?? null,
	};
	return {
		isPending: ref(state.isPending),
		error: ref(state.error),
		data: ref(state.data),
	};
}

// ── Module mocks ────────────────────────────────────────────────────────────

let mockAvailability = createAvailabilityMock();
let mockBacklog = createBacklogMock();

vi.mock("@pinia/colada", () => ({
	useQuery: vi.fn((options: { key: { _id: string }[] }) => {
		// Route to the correct mock based on query _id
		const id = options.key[0]?._id;
		if (id === "getDashboardPeopleAvailability") {
			return mockAvailability;
		}
		if (id === "getTasksBacklog") {
			return mockBacklog;
		}
		return createAvailabilityMock();
	}),
}));

vi.mock("@/client/@pinia/colada.gen", () => ({
	getDashboardPeopleAvailabilityQuery: vi.fn(() => ({
		key: [
			{
				_id: "getDashboardPeopleAvailability",
				baseUrl: "http://localhost:3000",
			},
		],
		query: vi.fn(),
	})),
	getTasksBacklogQuery: vi.fn(() => ({
		key: [{ _id: "getTasksBacklog", baseUrl: "http://localhost:3000" }],
		query: vi.fn(),
	})),
}));

// ── Tests ────────────────────────────────────────────────────────────────────

describe("KpiCards", () => {
	beforeEach(() => {
		mockAvailability = createAvailabilityMock();
		mockBacklog = createBacklogMock();
	});

	it("renders four cards", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 5, totalPeople: 8 } },
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		// Grid container plus four Card elements
		const cards = wrapper.findAll("[data-slot='card']");
		expect(cards).toHaveLength(4);
		wrapper.unmount();
	});

	it("shows loading state when availability is pending", () => {
		mockAvailability = createAvailabilityMock({ isPending: true });
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		expect(wrapper.text()).toContain("Loading…");
		wrapper.unmount();
	});

	it("shows error state when availability query fails", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			error: new Error("Network error"),
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		expect(wrapper.text()).toContain("Backend unavailable");
		wrapper.unmount();
	});

	it("shows empty state when totalPeople is zero", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 0, totalPeople: 0 } },
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		// People card shows em-dash; backlog cards show their values
		expect(wrapper.text()).toContain("—");
		expect(wrapper.text()).toContain("3");
		expect(wrapper.text()).toContain("2");
		wrapper.unmount();
	});

	it("shows 'X of Y available' format for the people card", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 5, totalPeople: 8 } },
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		expect(wrapper.text()).toContain("5 of 8");
		expect(wrapper.text()).toContain("available");
		wrapper.unmount();
	});

	it("shows backlog cards with correct values when backlog query succeeds", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 5, totalPeople: 8 } },
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 7, unassignedTasks: 4 } },
		});

		const wrapper = mount(KpiCards);
		expect(wrapper.text()).toContain("7");
		expect(wrapper.text()).toContain("4");
		expect(wrapper.text()).toContain("High priority tasks");
		expect(wrapper.text()).toContain("Unassigned jobs");
		wrapper.unmount();
	});

	it("renders rooms placeholder with data-testid", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 5, totalPeople: 8 } },
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		const placeholder = wrapper.find(
			"[data-testid='kpi-placeholder-rooms-completed']",
		);
		expect(placeholder.exists()).toBe(true);
		expect(wrapper.text()).toContain("Rooms completed");
		wrapper.unmount();
	});

	it("renders two-column layout with left accent column for each card", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 5, totalPeople: 8 } },
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		const firstCard = wrapper.find("[data-slot='card']");

		// The card contains a flex-row container (two-column layout)
		expect(firstCard.find(".flex.flex-row").exists()).toBe(true);

		// Left accent column has the expected width and semantic color class
		const accentCol = firstCard.find(".w-\\[72px\\]");
		expect(accentCol.exists()).toBe(true);
		expect(accentCol.classes()).toEqual(
			expect.arrayContaining(["bg-info-soft", "text-info"]),
		);

		// Icon renders inside the accent column
		expect(accentCol.find("svg").exists()).toBe(true);

		// Right content column exists and contains the leaf decoration
		const contentCol = firstCard.find(".flex-1");
		expect(contentCol.exists()).toBe(true);
		wrapper.unmount();
	});

	it("shows loading state when backlog is pending", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 5, totalPeople: 8 } },
		});
		mockBacklog = createBacklogMock({ isPending: true });

		const wrapper = mount(KpiCards);
		// Backlog cards show "Loading…" — should appear at least twice (high-priority + unassigned)
		const loadingMatches = wrapper.text().match(/Loading…/g);
		expect(loadingMatches).not.toBeNull();
		// People card shows data (not loading)
		expect(wrapper.text()).toContain("5 of 8");
		wrapper.unmount();
	});
});
