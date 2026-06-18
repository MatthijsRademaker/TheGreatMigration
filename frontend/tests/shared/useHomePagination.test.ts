// @vitest-environment jsdom

import { defineComponent, h, nextTick, ref, computed } from "vue";
import { describe, expect, it, vi, beforeEach } from "vitest";

// ── Mock planning window composable with reactive state ──────────────
const mockPlanWindowDays = ref<
	{ dateString: string; label: string; date: Date }[]
>([]);
const mockPlIsLoading = ref(false);
const mockPlIsError = ref(false);
const mockPlFormattedRange = ref<string | null>(null);

// Factory must be a plain function (not a ref-based reactive) for vi.mock hoisting.
vi.mock("@/shared/composables/usePlanningWindow", () => ({
	usePlanningWindow: () => ({
		planWindowDays: computed(() => mockPlanWindowDays.value),
		isLoading: computed(() => mockPlIsLoading.value),
		isError: computed(() => mockPlIsError.value),
		formattedRange: computed(() => mockPlFormattedRange.value),
		queryKey: ["planning-window"],
	}),
}));

// ── Helpers ──────────────────────────────────────────────────────────

/** Create a synthetic plan window of `count` days starting from `startISO` (YYYY-MM-DD). */
function makePlanWindow(count: number, startISO = "2026-07-05") {
	const days: { dateString: string; label: string; date: Date }[] = [];
	for (let i = 0; i < count; i++) {
		const d = new Date(startISO);
		d.setUTCDate(d.getUTCDate() + i);
		const iso = d.toISOString().slice(0, 10);
		days.push({ dateString: iso, label: iso, date: new Date(d) });
	}
	return days;
}

async function createHarness() {
	const { mount } = await import("@vue/test-utils");
	const { useHomePagination } = await import(
		"../../src/shared/composables/useHomePagination"
	);

	let state: ReturnType<typeof useHomePagination> | undefined;
	const wrapper = mount(
		defineComponent({
			setup() {
				state = useHomePagination();
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

// ── Tests ────────────────────────────────────────────────────────────

describe("useHomePagination", () => {
	beforeEach(async () => {
		mockPlanWindowDays.value = [];
		mockPlIsLoading.value = false;
		mockPlIsError.value = false;
		mockPlFormattedRange.value = null;
		// Reset the useHomePagination module to get fresh module-scoped refs.
		await vi.resetModules();
	});

	describe("page and totalPages", () => {
		it("initializes page at 1", async () => {
			mockPlanWindowDays.value = makePlanWindow(12);
			const { state, wrapper } = await createHarness();
			expect(state.page.value).toBe(1);
			wrapper.unmount();
		});

		it("computes totalPages as ceil(dayCount / daysPerPage)", async () => {
			// 12 days / 4 daysPerPage = 3 pages
			mockPlanWindowDays.value = makePlanWindow(12);
			const { state, wrapper } = await createHarness();
			expect(state.totalPages.value).toBe(3);
			wrapper.unmount();
		});

		it("computes totalPages with fractional remainder", async () => {
			// 13 days / 4 = ceil(3.25) = 4 pages
			mockPlanWindowDays.value = makePlanWindow(13);
			const { state, wrapper } = await createHarness();
			expect(state.totalPages.value).toBe(4);
			wrapper.unmount();
		});

		it("returns totalPages of 1 when planning window is empty", async () => {
			mockPlanWindowDays.value = [];
			const { state, wrapper } = await createHarness();
			expect(state.totalPages.value).toBe(1);
			wrapper.unmount();
		});

		it("returns totalPages of 1 for a single day", async () => {
			mockPlanWindowDays.value = makePlanWindow(1);
			const { state, wrapper } = await createHarness();
			expect(state.totalPages.value).toBe(1);
			wrapper.unmount();
		});
	});

	describe("goPrev / goNext clamping", () => {
		it("goNext advances page by 1", async () => {
			mockPlanWindowDays.value = makePlanWindow(20);
			const { state, wrapper } = await createHarness();
			state.goNext();
			await nextTick();
			expect(state.page.value).toBe(2);
			wrapper.unmount();
		});

		it("goPrev does not go below 1", async () => {
			mockPlanWindowDays.value = makePlanWindow(20);
			const { state, wrapper } = await createHarness();
			state.goPrev();
			await nextTick();
			expect(state.page.value).toBe(1);
			wrapper.unmount();
		});

		it("goNext does not exceed totalPages", async () => {
			mockPlanWindowDays.value = makePlanWindow(8);
			const { state, wrapper } = await createHarness();
			state.goNext(); // 2
			state.goNext(); // 3 (totalPages = 8/4 = 2, clamped to 2)
			await nextTick();
			// With 8 days and 4 daysPerPage, totalPages = 2.
			// goNext already clamps at 2.
			expect(state.page.value).toBe(2);
			wrapper.unmount();
		});

		it("goPrev retreats from higher page", async () => {
			mockPlanWindowDays.value = makePlanWindow(20);
			const { state, wrapper } = await createHarness();
			state.goNext(); // 2
			state.goNext(); // 3
			state.goPrev(); // 2
			await nextTick();
			expect(state.page.value).toBe(2);
			wrapper.unmount();
		});
	});

	describe("goToday", () => {
		it("navigates to the page containing today when today is within planning window", async () => {
			// Build 12 days starting from 2 days ago so today is included.
			const start = new Date();
			start.setUTCDate(start.getUTCDate() - 2);
			const startISO = start.toISOString().slice(0, 10);

			mockPlanWindowDays.value = makePlanWindow(12, startISO);
			const { state, wrapper } = await createHarness();

			// Navigate away first
			state.goNext();
			state.goNext();
			await nextTick();
			expect(state.page.value).toBeGreaterThan(1);

			state.goToday();
			await nextTick();

			// Today should be within the first few days; page should be 1.
			expect(state.page.value).toBe(1);
			wrapper.unmount();
		});

		it("falls back to page 1 when today is outside planning window", async () => {
			// Planning window entirely in the past
			const farPast = "2020-01-05";
			mockPlanWindowDays.value = makePlanWindow(8, farPast);
			const { state, wrapper } = await createHarness();

			// Navigate to a different page first
			state.goNext();
			await nextTick();
			expect(state.page.value).toBe(2);

			// Today (2026) is not in the 2020 window
			state.goToday();
			await nextTick();
			expect(state.page.value).toBe(1);
			wrapper.unmount();
		});

		it("does nothing when planning window is empty", async () => {
			mockPlanWindowDays.value = [];
			const { state, wrapper } = await createHarness();
			state.goToday();
			await nextTick();
			// Page should remain at 1 (default)
			expect(state.page.value).toBe(1);
			wrapper.unmount();
		});
	});

	describe("rangeLabel", () => {
		it("returns formatted range for days within page", async () => {
			mockPlanWindowDays.value = makePlanWindow(8, "2026-07-05");
			const { state, wrapper } = await createHarness();
			expect(state.rangeLabel.value).toContain("Jul");
			expect(state.rangeLabel.value).toContain("2026");
			wrapper.unmount();
		});

		it("returns em-dash when planning window is empty", async () => {
			mockPlanWindowDays.value = [];
			const { state, wrapper } = await createHarness();
			expect(state.rangeLabel.value).toBe("—");
			wrapper.unmount();
		});

		it("updates rangeLabel when page changes", async () => {
			mockPlanWindowDays.value = makePlanWindow(40, "2026-07-05");
			const { state, wrapper } = await createHarness();
			const page1Label = state.rangeLabel.value;

			state.goNext();
			await nextTick();
			const page2Label = state.rangeLabel.value;

			expect(page2Label).not.toBe(page1Label);
			wrapper.unmount();
		});
	});

	describe("isLoading / isError passthrough", () => {
		it("reflects planning window isLoading", async () => {
			mockPlIsLoading.value = true;
			const { state, wrapper } = await createHarness();
			expect(state.isLoading.value).toBe(true);
			mockPlIsLoading.value = false;
			await nextTick();
			expect(state.isLoading.value).toBe(false);
			wrapper.unmount();
		});

		it("reflects planning window isError", async () => {
			mockPlIsError.value = true;
			const { state, wrapper } = await createHarness();
			expect(state.isError.value).toBe(true);
			mockPlIsError.value = false;
			await nextTick();
			expect(state.isError.value).toBe(false);
			wrapper.unmount();
		});
	});
});
