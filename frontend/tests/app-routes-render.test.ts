import { renderToString } from "@vue/server-renderer";
import { createSSRApp } from "vue";
import { createMemoryHistory, createRouter } from "vue-router";
import { describe, expect, it } from "vitest";
import App from "../src/app/App.vue";
import { routes } from "../src/app/routes";

async function renderRoute(path: string) {
	const router = createRouter({
		history: createMemoryHistory(),
		routes,
	});

	await router.push(path);
	await router.isReady();

	const app = createSSRApp(App);
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
		expect(html).toContain("Moving dashboard");
		expect(html).toContain("Today’s plan");
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
			content: "Task foundation",
		},
		{
			path: "/calendar",
			title: "Schedule board",
			description: "Plan work across move days and balance available helpers.",
			content: "Schedule board foundation",
		},
		{
			path: "/people",
			title: "People availability",
			description: "Track who is available and where each person can help.",
			content: "People availability foundation",
		},
	] as const;

	for (const { path, title, description, content } of routeCases) {
		it(`renders ${path} with route metadata and view content`, async () => {
			const html = await renderRoute(path);

			expect(html).toContain(title);
			expect(html).toContain(description);
			expect(html).toContain(content);
		});
	}
});
