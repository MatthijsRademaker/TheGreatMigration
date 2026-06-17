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
		data: {
			summary: { availableToday: number; totalPeople: number };
			range?: { selectedDate?: string };
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

	it("renders four cards in correct order: High priority, People, Unassigned, Rooms", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: {
				summary: { availableToday: 5, totalPeople: 8 },
				range: { selectedDate: "2026-07-05" },
			},
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		const cards = wrapper.findAll("[data-slot='card']");
		expect(cards).toHaveLength(4);

		// Verify order: High priority → People → Unassigned → Rooms
		const cardTexts = cards.map((c) => c.text());
		expect(cardTexts[0]).toContain("High priority tasks");
		expect(cardTexts[1]).toContain("People available today");
		expect(cardTexts[2]).toContain("Unassigned jobs");
		expect(cardTexts[3]).toContain("Rooms completed");
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
		// Loading text appears inside CardContent
		expect(
			wrapper.findAll("[data-slot='card-content']").length,
		).toBeGreaterThan(0);
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
		// People card shows em-dash (empty); backlog cards show their values
		expect(wrapper.text()).toContain("—");
		expect(wrapper.text()).toContain("3");
		expect(wrapper.text()).toContain("2");
		// When totalPeople is zero and no selectedDate, subtitle falls back
		expect(wrapper.text()).toContain("available today");
		wrapper.unmount();
	});

	it("shows 'X / Y' fraction format with date subtitle for the people card", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: {
				summary: { availableToday: 5, totalPeople: 8 },
				range: { selectedDate: "2026-07-05" },
			},
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		expect(wrapper.text()).toContain("5 / 8");
		expect(wrapper.text()).toContain("available on Jul 5");
		wrapper.unmount();
	});

	it("shows 'available today' fallback subtitle when selectedDate is missing", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: { summary: { availableToday: 5, totalPeople: 8 } },
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);
		expect(wrapper.text()).toContain("5 / 8");
		expect(wrapper.text()).toContain("available today");
		wrapper.unmount();
	});

	it("shows backlog cards with correct values and new subtitle copy", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: {
				summary: { availableToday: 5, totalPeople: 8 },
				range: { selectedDate: "2026-07-05" },
			},
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
		// New outcome-focused subtitles
		expect(wrapper.text()).toContain("high priority tasks need attention");
		expect(wrapper.text()).toContain("jobs that need assignment");
		wrapper.unmount();
	});

	it("renders rooms placeholder with data-testid and border-success + icon chip", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: {
				summary: { availableToday: 5, totalPeople: 8 },
				range: { selectedDate: "2026-07-05" },
			},
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
		expect(placeholder.text()).toContain("Rooms completed");
		expect(placeholder.text()).toContain("—");
		expect(placeholder.text()).toContain("rooms fully packed and cleared");
		// Verify border-success left accent
		expect(placeholder.classes()).toContain("border-success");
		wrapper.unmount();
	});

	it("renders thin accent border and compact icon chip layout for each card", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: {
				summary: { availableToday: 5, totalPeople: 8 },
				range: { selectedDate: "2026-07-05" },
			},
		});
		mockBacklog = createBacklogMock({
			isPending: false,
			data: { summary: { highPriorityTasks: 3, unassignedTasks: 2 } },
		});

		const wrapper = mount(KpiCards);

		// No w-[72px] accent column should exist
		expect(wrapper.find(".w-\\[72px\\]").exists()).toBe(false);

		// First card (high priority) should have border-l-4 and border-destructive
		const cards = wrapper.findAll("[data-slot='card']");
		expect(cards[0].classes()).toContain("border-l-4");
		expect(cards[0].classes()).toContain("border-destructive");

		// Second card (people) should have border-info
		expect(cards[1].classes()).toContain("border-l-4");
		expect(cards[1].classes()).toContain("border-info");

		// Check compact icon chip (size-8 rounded-lg) exists with semantic background
		const iconChip = cards[0].find(".size-8.rounded-lg");
		expect(iconChip.exists()).toBe(true);
		expect(iconChip.classes()).toEqual(
			expect.arrayContaining(["bg-destructive-soft", "text-destructive"]),
		);

		// CardHeader and CardContent primitives are present
		expect(cards[0].find("[data-slot='card-header']").exists()).toBe(true);
		expect(cards[0].find("[data-slot='card-content']").exists()).toBe(true);

		// No flex-row two-column wrapper
		expect(cards[0].find(".flex.flex-row").exists()).toBe(false);
		wrapper.unmount();
	});

	it("shows loading state when backlog is pending", () => {
		mockAvailability = createAvailabilityMock({
			isPending: false,
			data: {
				summary: { availableToday: 5, totalPeople: 8 },
				range: { selectedDate: "2026-07-05" },
			},
		});
		mockBacklog = createBacklogMock({ isPending: true });

		const wrapper = mount(KpiCards);
		// Backlog cards show "Loading…" inside CardContent — appears at least twice (high-priority + unassigned)
		const loadingMatches = wrapper.text().match(/Loading…/g);
		expect(loadingMatches).not.toBeNull();
		expect(loadingMatches!.length).toBeGreaterThanOrEqual(2);
		// People card shows data (not loading) in its CardContent
		expect(wrapper.text()).toContain("5 / 8");
		expect(wrapper.text()).toContain("available on Jul 5");
		wrapper.unmount();
	});
});
