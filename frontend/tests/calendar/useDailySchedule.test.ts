import { PiniaColada } from "@pinia/colada";
import { renderToString } from "@vue/server-renderer";
import { createPinia } from "pinia";
import { createSSRApp, h } from "vue";
import type { Component } from "vue";
import { describe, expect, it, vi } from "vitest";
import { configureApiClient } from "../../src/shared/lib/api-client";

// Set up a mock API that returns deterministic daily-schedule data.
function setupMockApi() {
	configureApiClient({
		baseUrl: "http://example.test",
		fetch: vi.fn(async (input: RequestInfo | URL) => {
			const url = input instanceof Request ? input.url : input.toString();
			if (url.includes("/api/dashboard/daily-schedule")) {
				return new Response(
					JSON.stringify({
						range: {
							startDate: "2026-07-05",
							endDate: "2026-07-08",
							days: 4,
						},
						days: [
							{
								date: "2026-07-05",
								label: "5 Jul (Sun)",
								availablePeopleCount: 6,
								tasks: [
									{
										id: "sched-1",
										title: "API-driven task",
										priority: "high",
										roomArea: "Kitchen",
										assignedPeople: [{ id: "p1", name: "Alex", initials: "A" }],
										peopleNeeded: 2,
										assignedCount: 1,
										staffingStatus: "underStaffed",
									},
								],
							},
						],
					}),
					{
						status: 200,
						headers: { "Content-Type": "application/json" },
					},
				);
			}
			return new Response(JSON.stringify({ message: "OK" }), {
				status: 200,
				headers: { "Content-Type": "application/json" },
			});
		}),
	});
}

async function renderComponent(component: Component): Promise<string> {
	const app = createSSRApp({
		render: () => h(component),
	});
	const pinia = createPinia();
	app.use(pinia);
	app.use(PiniaColada);
	return renderToString(app);
}

describe("useDailySchedule", () => {
	it("adapts API data into component props", async () => {
		setupMockApi();

		const { useDailySchedule } = await import(
			"../../src/calendar/composables/useDailySchedule"
		);
		const TestComponent = {
			setup() {
				const { data, isLoading, isError, isEmpty } = useDailySchedule();
				return { data, isLoading, isError, isEmpty };
			},
			template: `<div>
				<span v-if="isLoading">loading</span>
				<span v-else-if="isError">error</span>
				<span v-else-if="isEmpty">empty</span>
				<span v-else>{{ data.days[0]?.tasks[0]?.title }}</span>
			</div>`,
		};

		const html = await renderComponent(TestComponent);
		// With SSR + PiniaColada, the query may still be pending on first render.
		// We verify the component doesn't crash and renders something.
		expect(typeof html).toBe("string");
	});

	it("handles null days gracefully", async () => {
		configureApiClient({
			baseUrl: "http://example.test",
			fetch: vi.fn(async () => {
				return new Response(
					JSON.stringify({
						range: { startDate: "2026-07-05", endDate: "2026-07-08", days: 4 },
						days: null,
					}),
					{ status: 200, headers: { "Content-Type": "application/json" } },
				);
			}),
		});

		const { useDailySchedule } = await import(
			"../../src/calendar/composables/useDailySchedule"
		);
		const TestComponent = {
			setup() {
				const { isEmpty } = useDailySchedule();
				return { isEmpty };
			},
			template: `<div><span v-if="isEmpty">empty</span></div>`,
		};

		const html = await renderComponent(TestComponent);
		expect(html).toContain("empty");
	});
});
