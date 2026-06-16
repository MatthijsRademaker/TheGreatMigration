import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h } from "vue";
import type { Component } from "vue";
import { describe, expect, it } from "vitest";
import DailySchedule from "../../src/calendar/DailySchedule.vue";

async function renderComponent(
	component: Component,
	props?: Record<string, unknown>,
): Promise<string> {
	const app = createSSRApp({
		render: () => h(component, props),
	});
	return renderToString(app);
}

describe("DailySchedule", () => {
	it("renders the panel title", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("Daily Schedule");
	});

	it("renders header controls", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("View by: Day");
		expect(html).toContain("Add task");
	});

	it("renders four day labels", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("2 Jul (Tue)");
		expect(html).toContain("3 Jul (Wed)");
		expect(html).toContain("4 Jul (Thu)");
		expect(html).toContain("5 Jul (Fri)");
	});

	it("renders availability counts", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("6 available");
		expect(html).toContain("7 available");
		expect(html).toContain("5 available");
	});

	it("renders representative task titles", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("Painting hall");
		expect(html).toContain("Steam walls");
		expect(html).toContain("Clean up");
		expect(html).toContain("Sanding");
		expect(html).toContain("Bedroom painting");
		expect(html).toContain("Touch up woodwork");
		expect(html).toContain("Living room finishing");
		expect(html).toContain("2nd floor walls");
		expect(html).toContain("Kitchen painting");
		expect(html).toContain("Final clean");
	});

	it("renders priority badges with correct variants", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain('data-variant="priorityHigh"');
		expect(html).toContain('data-variant="priorityMedium"');
		expect(html).toContain('data-variant="priorityLow"');
	});

	it("renders staffing counts", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("2 / 2");
		expect(html).toContain("1 / 1");
		expect(html).toContain("1 / 2");
	});

	it("renders under-staffed indicator for applicable tasks", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("— needs help");
		expect(html.match(/— needs help/g)?.length).toBe(4);
	});

	it("renders per-column Add task placeholders", async () => {
		const html = await renderComponent(DailySchedule);
		// Four day columns × one placeholder each
		const matches = html.match(/\+ Add task/g);
		expect(matches?.length).toBe(4);
	});

	it("renders assignee metadata", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain("Alex");
		expect(html).toContain("Morgan");
		expect(html).toContain("Sam");
		expect(html).toContain("Riley");
	});

	it("renders a Card shell", async () => {
		const html = await renderComponent(DailySchedule);
		expect(html).toContain('data-slot="card"');
		expect(html).toContain('data-slot="card-title"');
	});

	it("accepts custom props and renders them", async () => {
		const html = await renderComponent(DailySchedule, {
			days: [
				{
					date: "2026-08-01",
					label: "1 Aug (Sat)",
					availablePeopleCount: 3,
					tasks: [
						{
							id: "c1",
							title: "Custom task",
							priority: "low",
							roomArea: "Test Room",
							assignedPeople: [{ id: "x1", name: "Taylor", initials: "T" }],
							peopleNeeded: 1,
							assignedCount: 1,
							staffingStatus: "fullyStaffed",
						},
					],
				},
			],
		});
		expect(html).toContain("1 Aug (Sat)");
		expect(html).toContain("3 available");
		expect(html).toContain("Custom task");
		expect(html).toContain("Taylor");
	});
});
