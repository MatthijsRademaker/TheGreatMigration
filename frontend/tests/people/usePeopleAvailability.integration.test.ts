import { renderToString } from "@vue/server-renderer";
import { createSSRApp, defineComponent, h, ref, computed, nextTick } from "vue";
import { describe, expect, it, vi, beforeEach } from "vitest";
import type { DashboardBody } from "@/client/types.gen";

// --- Controlled test data ---
const PLANNING_WINDOW_START = "2026-07-05";
const EXPLICIT_START = "2024-07-02";

function makeDashboardBody(startDate: string, days: number): DashboardBody {
	const people: {
		id: string;
		name: string;
		initials: string;
		availability: { date: string; status: string }[];
	}[] = [];
	for (let i = 0; i < days; i++) {
		const cursor = new Date(startDate);
		cursor.setDate(cursor.getDate() + i);
		const isoDate = cursor.toISOString().slice(0, 10);
		people.push({
			id: "p1",
			name: "Test",
			initials: "TT",
			availability: [{ date: isoDate, status: "available" }],
		});
	}
	return {
		range: { startDate, endDate: "", days, selectedDate: startDate },
		summary: { availableToday: 1, totalPeople: 1 },
		people,
		pagination: { totalPeople: 1, page: 1, perPage: 1 },
		statuses: [
			{ id: "available", label: "Available", colorIntent: "green" },
			{ id: "busy", label: "Busy", colorIntent: "red" },
			{ id: "partial", label: "Partial", colorIntent: "yellow" },
			{ id: "off", label: "Off", colorIntent: "gray" },
		],
	};
}

// --- Mock planning window composable ---
const mockPlanningWindow = ref<{
	planWindowDays: { dateString: string; label: string; date: Date }[];
	isLoading: boolean;
	isError: boolean;
	formattedRange: string | null;
	queryKey: string[];
}>({
	planWindowDays: [
		{
			dateString: PLANNING_WINDOW_START,
			label: "Sun 5 Jul",
			date: new Date(PLANNING_WINDOW_START),
		},
	],
	isLoading: false,
	isError: false,
	formattedRange: "Jul 5 – Aug 13",
	queryKey: ["planning-window"],
});

vi.mock("@/shared/composables/usePlanningWindow", () => ({
	usePlanningWindow: () => ({
		get planWindowDays() {
			return computed(() => mockPlanningWindow.value.planWindowDays);
		},
		get isLoading() {
			return computed(() => mockPlanningWindow.value.isLoading);
		},
		get isError() {
			return computed(() => mockPlanningWindow.value.isError);
		},
		get formattedRange() {
			return computed(() => mockPlanningWindow.value.formattedRange);
		},
		queryKey: computed(() => mockPlanningWindow.value.queryKey),
	}),
}));

// --- Mock generated client ---
const capturedOptions = ref<{ query?: { start?: string } } | undefined>();

vi.mock("@/client/@pinia/colada.gen", () => ({
	getDashboardPeopleAvailabilityQueryKey: (options?: {
		query?: { start?: string };
	}) => ["getDashboardPeopleAvailability", options?.query?.start],
	getDashboardPeopleAvailabilityQuery: (options?: {
		query?: { start?: string };
	}) => {
		capturedOptions.value = options;
		return {
			key: ["getDashboardPeopleAvailability", options?.query?.start],
			query: () => {
				const start = options?.query?.start ?? PLANNING_WINDOW_START;
				return makeDashboardBody(start, 4);
			},
		};
	},
}));

// Also mock the actual query implementation useQuery will call
vi.mock("@pinia/colada", async (importOriginal) => {
	const actual = await importOriginal<typeof import("@pinia/colada")>();
	return {
		...actual,
		useQuery: vi.fn((optionsOrGetter: unknown) => {
			// Resolve options — could be a function or an object
			const opts =
				typeof optionsOrGetter === "function"
					? (
							optionsOrGetter as () => {
								key: unknown[];
								query: () => DashboardBody | Promise<DashboardBody>;
								enabled: boolean;
							}
						)()
					: (optionsOrGetter as {
							key: unknown[];
							query: () => DashboardBody | Promise<DashboardBody>;
							enabled: boolean;
						});

			const key = opts.key as string[];
			const shouldFetch = opts.enabled ?? true;
			const queryResult = shouldFetch ? opts.query() : undefined;
			const data = ref<DashboardBody | undefined>(
				queryResult instanceof Promise ? undefined : queryResult,
			);
			const isPending = ref(false);

			// Handle async query results
			if (queryResult instanceof Promise) {
				isPending.value = true;
				queryResult.then((result) => {
					data.value = result;
					isPending.value = false;
				});
			}

			return {
				data,
				isPending,
				isLoading: computed(() => isPending.value),
				error: ref(null),
				isFetching: ref(shouldFetch),
				isRefetching: ref(false),
				status: computed(() => (isPending.value ? "pending" : "success")),
				state: computed(() => ({
					status: isPending.value ? ("pending" as const) : ("success" as const),
					data: data.value,
					error: null,
				})),
				asyncStatus: computed(() =>
					isPending.value ? ("loading" as const) : ("idle" as const),
				),
				refresh: vi.fn(),
				refetch: vi.fn(),
				queryKey: key,
			};
		}),
	};
});

// --- Test helper ---
import { usePeopleAvailability } from "../../src/shared/composables/usePeopleAvailability";

async function mountComposable(opts?: { start?: string }) {
	const TestComponent = defineComponent({
		setup() {
			const state = usePeopleAvailability(opts);
			return () =>
				h("div", [
					h(
						"span",
						{ "data-testid": "days-count" },
						String(state.data.value?.days?.length ?? 0),
					),
					h(
						"span",
						{ "data-testid": "days-first" },
						state.data.value?.days?.[0] ?? "",
					),
					h(
						"span",
						{ "data-testid": "daysISO-count" },
						String(state.daysISO.value.length),
					),
					h(
						"span",
						{ "data-testid": "daysISO-first" },
						state.daysISO.value[0] ?? "",
					),
					h(
						"span",
						{ "data-testid": "daysISO-last" },
						state.daysISO.value[state.daysISO.value.length - 1] ?? "",
					),
					h(
						"span",
						{ "data-testid": "total" },
						String(state.data.value.totalPeople),
					),
				]);
		},
	});

	const app = createSSRApp(TestComponent);
	return renderToString(app);
}

describe("usePeopleAvailability planning-window integration", () => {
	beforeEach(() => {
		capturedOptions.value = undefined;
		mockPlanningWindow.value = {
			planWindowDays: [
				{
					dateString: PLANNING_WINDOW_START,
					label: "Sun 5 Jul",
					date: new Date(PLANNING_WINDOW_START),
				},
			],
			isLoading: false,
			isError: false,
			formattedRange: "Jul 5 – Aug 13",
			queryKey: ["planning-window"],
		};
	});

	it("passes planning window startDate as query param when no start given", async () => {
		await mountComposable();

		// The captured options should include start=PLANNING_WINDOW_START
		expect(capturedOptions.value).toBeDefined();
		expect(capturedOptions.value?.query?.start).toBe(PLANNING_WINDOW_START);
	});

	it("uses explicit start when provided, bypassing planning window", async () => {
		await mountComposable({ start: EXPLICIT_START });

		// The captured options should include start=EXPLICIT_START
		expect(capturedOptions.value).toBeDefined();
		expect(capturedOptions.value?.query?.start).toBe(EXPLICIT_START);
	});

	it("exposes daysISO array matching days label length", async () => {
		const html = await mountComposable();

		// Days and daysISO should have same length (4 days from makeDashboardBody with days=4)
		const daysCount = html.match(/data-testid="days-count">(\d+)</);
		const daysISOCount = html.match(/data-testid="daysISO-count">(\d+)</);
		expect(daysCount).toBeDefined();
		expect(daysISOCount).toBeDefined();
		expect(daysCount![1]).toBe(daysISOCount![1]);
	});

	it("exposes daysISO in correct order matching day labels", async () => {
		const html = await mountComposable();

		// First day label should correspond to first ISO date
		const daysFirst = html.match(/data-testid="days-first">([^<]+)</);
		const daysISOFirst = html.match(/data-testid="daysISO-first">([^<]+)</);
		expect(daysFirst).toBeDefined();
		expect(daysISOFirst).toBeDefined();

		// The first ISO date string should be the planning window start
		expect(daysISOFirst![1]).toBe(PLANNING_WINDOW_START);
	});

	it("defers dashboard query when planning window is loading", async () => {
		// Set planning window as loading with empty planWindowDays (data not yet loaded)
		mockPlanningWindow.value = {
			planWindowDays: [],
			isLoading: true,
			isError: false,
			formattedRange: null,
			queryKey: ["planning-window"],
		};

		await mountComposable();

		// When planning window is loading with empty planWindowDays,
		// startParam should be undefined (no planWindowDays to derive from)
		expect(capturedOptions.value?.query?.start).toBeUndefined();
	});
});

// --- Client-render tests for page navigation and page-reset ---
// @vitest-environment jsdom

function makeMultiDayPlanWindow(
	count: number,
	startDateStr: string = PLANNING_WINDOW_START,
) {
	const days: { dateString: string; label: string; date: Date }[] = [];
	for (let i = 0; i < count; i++) {
		const d = new Date(startDateStr);
		d.setUTCDate(d.getUTCDate() + i);
		const isoDate = d.toISOString().slice(0, 10);
		days.push({
			dateString: isoDate,
			label: "Day " + (i + 1),
			date: new Date(d),
		});
	}
	return days;
}

describe("usePeopleAvailability page navigation", () => {
	beforeEach(() => {
		capturedOptions.value = undefined;
		// Set up 21 planning window days => 3 pages with default 7 daysPerPage.
		mockPlanningWindow.value = {
			planWindowDays: makeMultiDayPlanWindow(21),
			isLoading: false,
			isError: false,
			formattedRange: "Jul 5 – Jul 25",
			queryKey: ["planning-window"],
		};
	});

	async function createHarness() {
		const { mount } = await import("@vue/test-utils");
		let state: ReturnType<typeof usePeopleAvailability> | undefined;
		const wrapper = mount(
			defineComponent({
				setup() {
					state = usePeopleAvailability();
					return () => h("div");
				},
			}),
			{ attachTo: document.body },
		);
		await nextTick();
		return {
			wrapper,
			get state() {
				return state!;
			},
		};
	}

	it("initializes page at 1", async () => {
		const { state, wrapper } = await createHarness();
		expect(state.page.value).toBe(1);
		wrapper.unmount();
	});

	it("goToNextPage advances page by 1", async () => {
		const { state, wrapper } = await createHarness();
		state.goToNextPage();
		await nextTick();
		expect(state.page.value).toBe(2);
		wrapper.unmount();
	});

	it("goToNextPage advances multiple times", async () => {
		const { state, wrapper } = await createHarness();
		state.goToNextPage();
		state.goToNextPage();
		await nextTick();
		expect(state.page.value).toBe(3);
		wrapper.unmount();
	});

	it("goToPrevPage retreats page by 1", async () => {
		const { state, wrapper } = await createHarness();
		state.goToNextPage(); // page 2
		state.goToNextPage(); // page 3
		state.goToPrevPage(); // page 2
		await nextTick();
		expect(state.page.value).toBe(2);
		wrapper.unmount();
	});

	it("goToPrevPage does not go below 1", async () => {
		const { state, wrapper } = await createHarness();
		state.goToPrevPage();
		await nextTick();
		expect(state.page.value).toBe(1);
		wrapper.unmount();
	});

	it("goToNextPage does not go beyond totalPages", async () => {
		const { state, wrapper } = await createHarness();
		state.goToNextPage(); // 2
		state.goToNextPage(); // 3
		state.goToNextPage(); // should stay at 3 (21/7 = 3 pages)
		await nextTick();
		expect(state.page.value).toBe(3);
		wrapper.unmount();
	});

	it("totalPages reflects planning window days / daysPerPage", async () => {
		const { state, wrapper } = await createHarness();
		expect(state.totalPages.value).toBe(3); // ceil(21/7)
		wrapper.unmount();
	});

	it("totalPages is at least 1", async () => {
		const { state, wrapper } = await createHarness();
		expect(state.totalPages.value).toBeGreaterThanOrEqual(1);
		wrapper.unmount();
	});

	it("resets page to 1 when planning window changes and page > 1", async () => {
		const { state, wrapper } = await createHarness();

		// Navigate to page 2.
		state.goToNextPage();
		await nextTick();
		expect(state.page.value).toBe(2);

		// Change the planning window (simulate admin updating the window).
		mockPlanningWindow.value = {
			planWindowDays: makeMultiDayPlanWindow(14, "2026-08-01"),
			isLoading: false,
			isError: false,
			formattedRange: "Aug 1 – Aug 14",
			queryKey: ["planning-window"],
		};
		await nextTick();

		// Page should reset to 1.
		expect(state.page.value).toBe(1);
		wrapper.unmount();
	});

	it("does not reset page on initial planning window load", async () => {
		// Start with loading state (empty days).
		mockPlanningWindow.value = {
			planWindowDays: [],
			isLoading: true,
			isError: false,
			formattedRange: null,
			queryKey: ["planning-window"],
		};

		const { mount } = await import("@vue/test-utils");
		let state: ReturnType<typeof usePeopleAvailability> | undefined;
		const wrapper = mount(
			defineComponent({
				setup() {
					state = usePeopleAvailability({ page: 2 });
					return () => h("div");
				},
			}),
			{ attachTo: document.body },
		);
		await nextTick();

		// Page should stay at 2 (no reset on initial load).
		expect(state!.page.value).toBe(2);

		// Now simulate planning window data loading.
		mockPlanningWindow.value = {
			planWindowDays: makeMultiDayPlanWindow(21),
			isLoading: false,
			isError: false,
			formattedRange: "Jul 5 – Jul 25",
			queryKey: ["planning-window"],
		};
		await nextTick();

		// Page should not reset when transitioning from loading to loaded.
		expect(state!.page.value).toBe(2);
		wrapper.unmount();
	});
});
