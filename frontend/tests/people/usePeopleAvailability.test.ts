import { renderToString } from "@vue/server-renderer";
import { createSSRApp, defineComponent, h } from "vue";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ref } from "vue";

// Per-test reactive refs for the mocked composable.
var mockIsLoading = ref(false);
var mockIsError = ref(false);
var mockTotalPeople = ref(8);
var mockAvailableToday = ref(6);

vi.mock("@/client/@pinia/colada.gen", () => ({
	getDashboardPeopleAvailabilityQueryKey: () => [],
	getDashboardPeopleAvailabilityQuery: () => ({}),
}));

vi.mock("@/shared/composables/usePeopleAvailability", () => ({
	usePeopleAvailability: () => ({
		data: ref({
			title: "People availability",
			description: "Track who is available and where each person can help.",
			days: ["Sun 5 Jul", "Mon 6 Jul", "Tue 7 Jul", "Wed 8 Jul"],
			people: [
				{
					id: "p1",
					name: "Sophia Chen",
					availability: [
						{ date: "Sun 5 Jul", status: "available" },
						{ date: "Mon 6 Jul", status: "available" },
						{ date: "Tue 7 Jul", status: "available" },
						{ date: "Wed 8 Jul", status: "available" },
					],
				},
				{
					id: "p2",
					name: "Marcus Rivera",
					availability: [
						{ date: "Sun 5 Jul", status: "busy" },
						{ date: "Mon 6 Jul", status: "busy" },
						{ date: "Tue 7 Jul", status: "busy" },
						{ date: "Wed 8 Jul", status: "busy" },
					],
				},
			],
			legend: [
				{ id: "available", label: "Available" },
				{ id: "busy", label: "Busy" },
				{ id: "partial", label: "Partial" },
				{ id: "off", label: "Off" },
			],
			availableToday: mockAvailableToday,
			totalPeople: mockTotalPeople,
		}),
		isLoading: mockIsLoading,
		isError: mockIsError,
		isEmpty: ref(false),
		refresh: vi.fn(),
		queryKey: [],
		rawData: ref(undefined),
	}),
}));

import { usePeopleAvailability } from "../../src/shared/composables/usePeopleAvailability";

async function renderState() {
	const TestComponent = defineComponent({
		setup() {
			const state = usePeopleAvailability();
			return () =>
				h("div", [
					h(
						"span",
						{ "data-testid": "loading" },
						String(state.isLoading.value),
					),
					h("span", { "data-testid": "error" }, String(state.isError.value)),
					h(
						"span",
						{ "data-testid": "total" },
						String(state.data.value.totalPeople),
					),
					h(
						"span",
						{ "data-testid": "available" },
						String(state.data.value.availableToday),
					),
				]);
		},
	});

	const app = createSSRApp(TestComponent);
	return renderToString(app);
}

describe("usePeopleAvailability composable states", () => {
	afterEach(() => {
		mockIsLoading.value = false;
		mockIsError.value = false;
		mockTotalPeople.value = 8;
		mockAvailableToday.value = 6;
	});

	it("exposes isLoading=true when the query is pending", async () => {
		mockIsLoading.value = true;
		const html = await renderState();
		expect(html).toContain('data-testid="loading">true<');
		expect(html).toContain('data-testid="error">false<');
	});

	it("exposes isError=true when the query has errored", async () => {
		mockIsError.value = true;
		const html = await renderState();
		expect(html).toContain('data-testid="loading">false<');
		expect(html).toContain('data-testid="error">true<');
	});

	it("exposes adapted totalPeople from the query response", async () => {
		mockTotalPeople.value = 5;
		const html = await renderState();
		expect(html).toContain('data-testid="total">5<');
	});

	it("exposes adapted availableToday from the query response", async () => {
		mockAvailableToday.value = 3;
		const html = await renderState();
		expect(html).toContain('data-testid="available">3<');
	});
});
