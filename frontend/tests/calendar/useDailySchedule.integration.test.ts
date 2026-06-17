import { createSSRApp, defineComponent, h, ref, computed, nextTick } from "vue";
import { describe, expect, it, vi, beforeEach } from "vitest";
import type { DailyScheduleBody, ScheduleDay } from "@/client/types.gen";

// --- Controlled test data ---
const PLANNING_WINDOW_START = "2026-07-05";

function makeScheduleDay(isoDate: string): ScheduleDay {
	return {
		date: isoDate,
		label: isoDate,
		availablePeopleCount: 6,
		tasks: [],
	};
}

function makeDailyScheduleBody(
	startDate: string,
	days: number,
): DailyScheduleBody {
	const endDate = new Date(startDate);
	endDate.setDate(endDate.getDate() + days - 1);

	const scheduleDays: ScheduleDay[] = [];
	for (let i = 0; i < days; i++) {
		const cursor = new Date(startDate);
		cursor.setDate(cursor.getDate() + i);
		scheduleDays.push(makeScheduleDay(cursor.toISOString().slice(0, 10)));
	}

	return {
		range: {
			startDate,
			endDate: endDate.toISOString().slice(0, 10),
			days,
		},
		days: scheduleDays,
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
const capturedOptions = ref<
	{ query?: { start?: string; days?: number } } | undefined
>();

vi.mock("@/client/@pinia/colada.gen", () => ({
	getDashboardDailyScheduleQueryKey: (options?: {
		query?: { start?: string; days?: number };
	}) => ["getDashboardDailySchedule", options?.query?.start],
	getDashboardDailyScheduleQuery: (options?: {
		query?: { start?: string; days?: number };
	}) => {
		capturedOptions.value = options;
		return {
			key: ["getDashboardDailySchedule", options?.query?.start],
			query: () => {
				const start = options?.query?.start ?? PLANNING_WINDOW_START;
				const days = options?.query?.days ?? 4;
				return makeDailyScheduleBody(start, days);
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
								query: () => DailyScheduleBody | Promise<DailyScheduleBody>;
								enabled: boolean;
							}
						)()
					: (optionsOrGetter as {
							key: unknown[];
							query: () => DailyScheduleBody | Promise<DailyScheduleBody>;
							enabled: boolean;
						});

			const key = opts.key as string[];
			const shouldFetch = opts.enabled ?? true;
			const queryResult = shouldFetch ? opts.query() : undefined;
			const data = ref<DailyScheduleBody | undefined>(
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

// --- Test helpers ---
import { useDailySchedule } from "../../src/calendar/composables/useDailySchedule";

async function mountComposable(opts?: { start?: string }) {
	const TestComponent = defineComponent({
		setup() {
			const state = useDailySchedule(opts);
			return () =>
				h("div", [
					h(
						"span",
						{ "data-testid": "days-count" },
						String(state.data.value.days.length),
					),
					h(
						"span",
						{ "data-testid": "days-first" },
						state.data.value.days[0]?.label ?? "",
					),
				]);
		},
	});

	const app = createSSRApp(TestComponent);
	return renderToString(app);
}

// Need renderToString for the SSR tests above
import { renderToString } from "@vue/server-renderer";

describe("useDailySchedule planning-window integration", () => {
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

	it("passes planning window startDate and daysPerPage as query params when no start given", async () => {
		await mountComposable();

		expect(capturedOptions.value).toBeDefined();
		expect(capturedOptions.value?.query?.start).toBe(PLANNING_WINDOW_START);
		expect(capturedOptions.value?.query?.days).toBe(4);
	});

	it("uses explicit start when provided, bypassing planning window", async () => {
		await mountComposable({ start: "2024-07-02" });

		expect(capturedOptions.value).toBeDefined();
		expect(capturedOptions.value?.query?.start).toBe("2024-07-02");
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

describe("useDailySchedule page navigation", () => {
	beforeEach(() => {
		capturedOptions.value = undefined;
		// Set up 20 planning window days => 5 pages with default 4 daysPerPage.
		mockPlanningWindow.value = {
			planWindowDays: makeMultiDayPlanWindow(20),
			isLoading: false,
			isError: false,
			formattedRange: "Jul 5 – Jul 24",
			queryKey: ["planning-window"],
		};
	});

	async function createHarness() {
		const { mount } = await import("@vue/test-utils");
		let state: ReturnType<typeof useDailySchedule> | undefined;
		const wrapper = mount(
			defineComponent({
				setup() {
					state = useDailySchedule();
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
		state.goToNextPage(); // 4
		state.goToNextPage(); // 5
		state.goToNextPage(); // should stay at 5 (20/4 = 5 pages)
		await nextTick();
		expect(state.page.value).toBe(5);
		wrapper.unmount();
	});

	it("totalPages reflects planning window days / daysPerPage", async () => {
		const { state, wrapper } = await createHarness();
		expect(state.totalPages.value).toBe(5); // ceil(20/4)
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
		let state: ReturnType<typeof useDailySchedule> | undefined;
		const wrapper = mount(
			defineComponent({
				setup() {
					state = useDailySchedule({ page: 2 });
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
			planWindowDays: makeMultiDayPlanWindow(20),
			isLoading: false,
			isError: false,
			formattedRange: "Jul 5 – Jul 24",
			queryKey: ["planning-window"],
		};
		await nextTick();

		// Page should not reset when transitioning from loading to loaded.
		expect(state!.page.value).toBe(2);
		wrapper.unmount();
	});
});
