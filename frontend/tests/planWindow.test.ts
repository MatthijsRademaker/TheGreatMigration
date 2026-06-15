import { PiniaColada } from "@pinia/colada";
import { renderToString } from "@vue/server-renderer";
import { createPinia } from "pinia";
import { createSSRApp, defineComponent, h } from "vue";
import { describe, expect, it, vi } from "vitest";
import { configureApiClient } from "../src/shared/lib/api-client";
import { usePlanningWindow } from "../src/shared/composables/usePlanningWindow";
import { formatPlanDayLabel } from "../src/shared/lib/planWindow";

function renderComposable(mockResponse: unknown) {
	const apiFetch = vi.fn(async () => {
		return new Response(JSON.stringify(mockResponse), {
			status: 200,
			headers: { "Content-Type": "application/json" },
		});
	});

	configureApiClient({ baseUrl: "http://example.test", fetch: apiFetch });

	const TestComponent = defineComponent({
		setup() {
			const state = usePlanningWindow();
			return () =>
				h("div", [
					h(
						"span",
						{ "data-testid": "days" },
						String(state.planWindowDayCount.value),
					),
					h(
						"span",
						{ "data-testid": "loading" },
						String(state.isLoading.value),
					),
					h("span", { "data-testid": "error" }, String(state.isError.value)),
					state.planWindowDays.value.length > 0
						? h(
								"span",
								{ "data-testid": "firstDate" },
								state.planWindowDays.value[0].dateString,
							)
						: null,
					state.planWindowDays.value.length > 0
						? h(
								"span",
								{ "data-testid": "lastDate" },
								state.planWindowDays.value[
									state.planWindowDays.value.length - 1
								].dateString,
							)
						: null,
				]);
		},
	});

	const app = createSSRApp(TestComponent);
	app.use(createPinia());
	app.use(PiniaColada);

	return renderToString(app);
}

describe("usePlanningWindow composable", () => {
	it("derives planWindowDays and planWindowDayCount on success", async () => {
		const html = await renderComposable({
			startDate: "2026-07-05",
			endDate: "2026-08-13",
			days: 40,
		});

		// Success: days count matches response
		expect(html).toContain('data-testid="days">40<');

		// First and last date boundaries
		expect(html).toContain('data-testid="firstDate">2026-07-05<');
		expect(html).toContain('data-testid="lastDate">2026-08-13<');

		// Not loading, not error
		expect(html).toContain('data-testid="loading">false<');
		expect(html).toContain('data-testid="error">false<');
	});
});

describe("formatPlanDayLabel utility", () => {
	it("for 2026-07-05 contains day and month", () => {
		const label = formatPlanDayLabel(new Date("2026-07-05T12:00:00Z"));
		expect(label).toContain("5");
		expect(label).toContain("Jul");
	});

	it("for 2026-08-13 contains day and month", () => {
		const label = formatPlanDayLabel(new Date("2026-08-13T12:00:00Z"));
		expect(label).toContain("13");
		expect(label).toContain("Aug");
	});
});
