import { PiniaColada } from "@pinia/colada";
import { renderToString } from "@vue/server-renderer";
import { createPinia } from "pinia";
import { createSSRApp } from "vue";
import { createMemoryHistory, createRouter } from "vue-router";
import { describe, expect, it, vi } from "vitest";
import App from "../src/app/App.vue";
import { configureApiClient } from "../src/shared/lib/api-client";
import { routes } from "../src/app/routes";

async function renderRoute(path: string) {
	const router = createRouter({
		history: createMemoryHistory(),
		routes,
	});

	await router.push(path);
	await router.isReady();

	configureApiClient({
		baseUrl: "http://example.test",
		fetch: vi.fn(async (input: RequestInfo | URL) => {
			const url = input instanceof Request ? input.url : input.toString();
			if (url.includes("/api/dashboard/people-availability")) {
				return new Response(
					JSON.stringify({
						range: {
							startDate: "2026-07-05",
							endDate: "2026-08-13",
							days: 40,
							selectedDate: "2026-07-05",
						},
						summary: {
							availableToday: 6,
							totalPeople: 8,
						},
						people: [],
						statuses: [],
					}),
					{
						status: 200,
						headers: { "Content-Type": "application/json" },
					},
				);
			}
			if (url.includes("/api/tasks/backlog")) {
				return new Response(
					JSON.stringify({
						summary: {
							totalTasks: 11,
							highPriorityTasks: 4,
							unassignedTasks: 3,
							understaffedTasks: 2,
						},
						tasks: [],
						priorities: [],
						statuses: [],
					}),
					{
						status: 200,
						headers: { "Content-Type": "application/json" },
					},
				);
			}
			if (url.includes("/api/planning-window")) {
				return new Response(
					JSON.stringify({
						startDate: "2026-07-05",
						endDate: "2026-08-13",
						days: 40,
					}),
					{
						status: 200,
						headers: { "Content-Type": "application/json" },
					},
				);
			}
			return new Response(
				JSON.stringify({ message: "Hello from the backend!" }),
				{
					status: 200,
					headers: { "Content-Type": "application/json" },
				},
			);
		}),
	});

	const app = createSSRApp(App);
	app.use(createPinia());
	app.use(PiniaColada);
	app.use(router);

	return renderToString(app);
}

describe("application route rendering", () => {
	it("renders the shared shell and sidebar chrome around the home route", async () => {
		const html = await renderRoute("/");

		expect(html).toContain("The Great Migration");
		expect(html).toContain("House move planner");
		expect(html).toContain("Planning mode");
		expect(html).toContain("Dashboard");
		expect(html).toContain("Tasks");
		expect(html).toContain("Schedule");
		expect(html).toContain("People");
		expect(html).toContain("Rooms / Areas");
		expect(html).toContain("Settings");
		expect(html).toContain("Moving dashboard");
		expect(html).toContain("Today’s plan");
		expect(html).toContain("Add note");
		expect(html).toContain("Help &amp; Support");
	});

	const routeCases = [
		{
			path: "/",
			title: "Moving dashboard",
			description:
				"Today’s priorities, staffing gaps, and move notes at a glance.",
			content: "Today’s plan",
		},
		{
			path: "/tasks",
			title: "Task backlog",
			description:
				"Capture jobs, priorities, staffing needs, and planning status.",
			content: "Task Management",
		},
		{
			path: "/calendar",
			title: "Schedule board",
			description: "Plan work across move days and balance available helpers.",
			content: "Daily Schedule",
		},
		{
			path: "/people",
			title: "People availability",
			description: "Track who is available and where each person can help.",
			content: "1 of 4 available today",
		},
		{
			path: "/rooms",
			title: "Rooms / Areas",
			description:
				"Organize and label rooms, floors, and zones for a clear move-day plan.",
			content: "Feature coming soon",
		},
		{
			path: "/settings",
			title: "Settings",
			description:
				"Configure your move preferences, notification defaults, and account details.",
			content: "Feature coming soon",
		},
	] as const;

	for (const { path, title, description, content } of routeCases) {
		it(`renders ${path} with route metadata and view content`, async () => {
			const html = await renderRoute(path);

			expect(html).toContain(title);
			expect(html).toContain(description);
			expect(html).toContain(content);

			if (path === "/calendar") {
				expect(html).toContain("Painting hall");
				expect(html).toContain("2 / 2");
				expect(html).toContain('data-variant="priorityHigh"');
				expect(html).toContain("2 Jul (Tue)");
				expect(html).toContain("5 Jul (Fri)");
				expect(html).toContain("+ Add task");
				expect(html).not.toContain("plan-day-column");
			}

			if (path === "/") {
				expect(html).toContain("People available today");
				expect(html).toContain("High priority tasks");
				expect(html).toContain("Unassigned jobs");
				expect(html).toContain("Rooms completed");

				// Verify rendered KPI values from mock data
				expect(html).toContain("6");
				expect(html).toContain("of 8");
				expect(html).toContain("available");
				expect(html).toContain("4");
				expect(html).toContain("3");
			}

			if (path === "/tasks") {
				expect(html).toContain("Painting hall");
				expect(html).toContain("People Needed");
				expect(html).toContain("Room / Area");
				expect(html).toContain("Unassigned");
				expect(html).toContain("Filter");
				expect(html).toContain("Add Task");
			}
		});
	}
});
