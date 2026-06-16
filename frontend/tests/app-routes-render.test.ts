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
							endDate: "2026-07-08",
							days: 4,
							selectedDate: "2026-07-05",
						},
						summary: {
							availableToday: 6,
							totalPeople: 8,
						},
						people: [
							{
								id: "p1",
								name: "Sophia Chen",
								initials: "SC",
								availability: [
									{ date: "2026-07-05", status: "available" },
									{ date: "2026-07-06", status: "available" },
									{ date: "2026-07-07", status: "available" },
									{ date: "2026-07-08", status: "available" },
								],
							},
							{
								id: "p2",
								name: "Marcus Rivera",
								initials: "MR",
								availability: [
									{ date: "2026-07-05", status: "available" },
									{ date: "2026-07-06", status: "available" },
									{ date: "2026-07-07", status: "available" },
									{ date: "2026-07-08", status: "available" },
								],
							},
							{
								id: "p3",
								name: "Elena Kowalski",
								initials: "EK",
								availability: [
									{ date: "2026-07-05", status: "available" },
									{ date: "2026-07-06", status: "available" },
									{ date: "2026-07-07", status: "available" },
									{ date: "2026-07-08", status: "available" },
								],
							},
							{
								id: "p4",
								name: "James Okafor",
								initials: "JO",
								availability: [
									{ date: "2026-07-05", status: "busy" },
									{ date: "2026-07-06", status: "busy" },
									{ date: "2026-07-07", status: "busy" },
									{ date: "2026-07-08", status: "busy" },
								],
							},
						],
						statuses: [
							{ id: "available", label: "Available", colorIntent: "success" },
							{ id: "busy", label: "Busy", colorIntent: "destructive" },
							{ id: "partial", label: "Partial", colorIntent: "warning" },
							{ id: "off", label: "Off", colorIntent: "muted" },
						],
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
							totalTasks: 3,
							highPriorityTasks: 1,
							unassignedTasks: 1,
							understaffedTasks: 1,
						},
						tasks: [
							{
								id: "task-1",
								title: "Pack kitchen boxes",
								priority: "high",
								peopleNeeded: 3,
								room: "Kitchen",
								status: "ready",
								assignedTo: [],
							},
							{
								id: "task-2",
								title: "Disassemble bed frames",
								priority: "medium",
								peopleNeeded: 2,
								room: "Bedroom",
								status: "ready",
								assignedTo: [],
							},
							{
								id: "task-3",
								title: "Move living room furniture",
								priority: "medium",
								peopleNeeded: 4,
								room: "Living Room",
								status: "assigned",
								assignedTo: ["p1"],
							},
						],
						priorities: [
							{ id: "high", label: "High", colorIntent: "destructive" },
							{ id: "medium", label: "Medium", colorIntent: "warning" },
							{ id: "low", label: "Low", colorIntent: "success" },
						],
						statuses: [
							{ id: "backlog", label: "Backlog", colorIntent: "muted" },
							{ id: "ready", label: "Ready", colorIntent: "info" },
							{ id: "assigned", label: "Assigned", colorIntent: "success" },
						],
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
			if (url.includes("/api/rooms")) {
				return new Response(
					JSON.stringify({
						rooms: [],
					}),
					{
						status: 200,
						headers: { "Content-Type": "application/json" },
					},
				);
			}
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
		expect(html).toContain("Jul");
		expect(html).toContain("2026");
		expect(html).toContain("40 days");
		expect(html).not.toContain("Planning mode");
		expect(html).toContain("Dashboard");
		expect(html).toContain("Tasks");
		expect(html).toContain("Schedule");
		expect(html).toContain("People");
		expect(html).toContain("Rooms / Areas");
		expect(html).toContain("Settings");
		expect(html).toContain("Moving dashboard");
		expect(html).toContain("Add note");
		expect(html).toContain("Help &amp; Support");
	});

	const routeCases = [
		{
			path: "/",
			title: "Moving dashboard",
			description:
				"Today’s priorities, staffing gaps, and move notes at a glance.",
			content: "Move notes",
		},
		{
			path: "/tasks",
			title: "Task backlog",
			description:
				"Capture jobs, priorities, staffing needs, and planning status.",
			content: "Add Task",
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
			content: "6 of 8 available today",
		},
		{
			path: "/rooms",
			title: "Rooms / Areas",
			description:
				"Organize and label rooms, floors, and zones for a clear move-day plan.",
			content: "Add Room / Area",
		},
		{
			path: "/settings",
			title: "Settings",
			description:
				"Configure your move preferences, notification defaults, and account details.",
			content: "Planning window",
		},
	] as const;

	for (const { path, title, description, content } of routeCases) {
		it(`renders ${path} with route metadata and view content`, async () => {
			const html = await renderRoute(path);

			expect(html).toContain(title);
			expect(html).toContain(description);
			expect(html).toContain(content);

			if (path === "/calendar") {
				expect(html).toContain("API-driven task");
				expect(html).toContain("1 / 2");
				expect(html).toContain('data-variant="priorityHigh"');
				expect(html).toContain("5 Jul (Sun)");
				expect(html).toContain("+ Add task");
				expect(html).not.toContain("plan-day-column");
			}

			if (path === "/rooms") {
				expect(html).toContain("No rooms or areas yet");
				expect(html).not.toContain("Feature coming soon");
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
				expect(html).toContain("1");

				// Home renders TaskManagementPanel in readOnly mode
				expect(html).toContain("Task Management");
				// readOnly mode: no Filter or Add Task buttons
				expect(html).not.toContain("Filter");
				expect(html).not.toContain("Add Task");
				// Backlog-backed task rows
				expect(html).toContain("Pack kitchen boxes");
				expect(html).toContain("Disassemble bed frames");

				expect(html).toContain("People availability");
				expect(html).toContain("Daily Schedule");
				expect(html).toContain("Move notes");
			}

			if (path === "/people") {
				expect(html).toContain("Add a person");
				expect(html).toContain("Manage people");
				expect(html).toContain('data-variant="available"');
				expect(html).toContain('data-variant="busy"');
				expect(html).toContain("Sophia Chen");
				expect(html).toContain("Marcus Rivera");
				expect(html).toContain("Elena Kowalski");
				expect(html).toContain("James Okafor");
			}

			if (path === "/tasks") {
				// TasksView renders its own management controls and backlog-backed task list
				expect(html).toContain("Add Task");
				expect(html).toContain("Task backlog");
				// Backlog-backed task rows from the mock
				expect(html).toContain("Pack kitchen boxes");
				expect(html).toContain("Disassemble bed frames");
				expect(html).toContain("Move living room furniture");
				// Management controls
				expect(html).toContain("Edit");
				expect(html).toContain("Delete");
				// No fixture content
				expect(html).not.toContain("Painting hall");
			}
		});
	}
});
