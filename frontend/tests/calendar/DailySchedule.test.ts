import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h } from "vue";
import type { Component } from "vue";
import { describe, expect, it } from "vitest";
import DailySchedule from "../../src/calendar/DailySchedule.vue";

interface AssignedPerson {
	id: string;
	name: string;
	initials: string;
}

interface TaskCard {
	id: string;
	title: string;
	priority: "high" | "medium" | "low";
	roomArea: string;
	assignedPeople: AssignedPerson[];
	peopleNeeded: number;
	assignedCount: number;
	staffingStatus: "fullyStaffed" | "underStaffed";
	scheduledDate: string;
}

interface ScheduleDay {
	date: string;
	label: string;
	availablePeopleCount: number;
	tasks: TaskCard[];
}

const sampleDays: ScheduleDay[] = [
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
				scheduledDate: "2026-08-01",
			},
			{
				id: "c2",
				title: "Urgent fix",
				priority: "high",
				roomArea: "Main Hall",
				assignedPeople: [
					{ id: "x2", name: "Alex", initials: "A" },
					{ id: "x3", name: "Morgan", initials: "M" },
				],
				peopleNeeded: 3,
				assignedCount: 2,
				staffingStatus: "underStaffed",
				scheduledDate: "2026-08-01",
			},
		],
	},
	{
		date: "2026-08-02",
		label: "2 Aug (Sun)",
		availablePeopleCount: 5,
		tasks: [
			{
				id: "c3",
				title: "Medium work",
				priority: "medium",
				roomArea: "Office",
				assignedPeople: [{ id: "x4", name: "Sam", initials: "S" }],
				peopleNeeded: 1,
				assignedCount: 1,
				staffingStatus: "fullyStaffed",
				scheduledDate: "2026-08-02",
			},
		],
	},
];

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
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("Daily Schedule");
	});

	it("renders header controls", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("View by: Day");
		expect(html).toContain("Add task");
	});

	it("renders day labels from props", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("1 Aug (Sat)");
		expect(html).toContain("2 Aug (Sun)");
	});

	it("renders availability counts", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("3 available");
		expect(html).toContain("5 available");
	});

	it("renders task titles from props", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("Custom task");
		expect(html).toContain("Urgent fix");
		expect(html).toContain("Medium work");
	});

	it("renders priority badges with correct variants", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain('data-variant="priorityHigh"');
		expect(html).toContain('data-variant="priorityMedium"');
		expect(html).toContain('data-variant="priorityLow"');
	});

	it("renders staffing counts", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("1 / 1");
		expect(html).toContain("2 / 3");
	});

	it("renders under-staffed indicator for applicable tasks", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("— needs help");
	});

	it("renders per-column Add task placeholders", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		const matches = html.match(/\+ Add task/g);
		expect(matches?.length).toBe(2);
	});

	it("renders assignee metadata", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain("Taylor");
		expect(html).toContain("Alex");
		expect(html).toContain("Morgan");
		expect(html).toContain("Sam");
	});

	it("renders a Card shell", async () => {
		const html = await renderComponent(DailySchedule, { days: sampleDays });
		expect(html).toContain('data-slot="card"');
		expect(html).toContain('data-slot="card-title"');
	});

	it("renders correctly with empty days array (no crash)", async () => {
		const html = await renderComponent(DailySchedule, { days: [] });
		expect(html).toContain("Daily Schedule");
		expect(html).toContain("View by: Day");
		// Should render but with no day columns or task cards
	});

	it("renders correctly with undefined days (graceful empty)", async () => {
		const html = await renderComponent(DailySchedule, {});
		expect(html).toContain("Daily Schedule");
		// No crash, renders shell with no days
	});

	it("hides edit, delete, and add-task controls when readonly is true", async () => {
		const html = await renderComponent(DailySchedule, {
			days: sampleDays,
			readOnly: true,
		});
		expect(html).toContain("Daily Schedule");
		expect(html).toContain("Custom task");
		expect(html).not.toContain(">Edit<");
		expect(html).not.toContain(">Delete<");
		expect(html).not.toContain("+ Add task");
	});

	it("shows edit, delete, and add-task controls when readonly is false (default)", async () => {
		const html = await renderComponent(DailySchedule, {
			days: sampleDays,
			readOnly: false,
		});
		expect(html).toContain(">Edit<");
		expect(html).toContain(">Delete<");
		expect(html).toContain("+ Add task");
		expect(html).toContain("Add task");
	});

	describe("pagination bar", () => {
		it("renders pagination bar when page and totalPages are provided", async () => {
			const html = await renderComponent(DailySchedule, {
				days: sampleDays,
				page: 2,
				totalPages: 5,
			});
			expect(html).toContain("Page 2 of 5");
			expect(html).toContain("Previous");
			expect(html).toContain("Next");
		});

		it("does not render pagination bar when pagination props are absent", async () => {
			const html = await renderComponent(DailySchedule, {
				days: sampleDays,
			});
			expect(html).not.toContain("Page 1 of");
			expect(html).not.toContain("Previous");
			expect(html).not.toContain("Next");
		});

		it("renders date range label when provided", async () => {
			const html = await renderComponent(DailySchedule, {
				days: sampleDays,
				page: 1,
				totalPages: 3,
				dateRangeLabel: "1 Aug (Sat) – 4 Aug (Tue)",
			});
			expect(html).toContain("1 Aug (Sat) – 4 Aug (Tue)");
		});

		it("disables Previous button on page 1", async () => {
			const html = await renderComponent(DailySchedule, {
				days: sampleDays,
				page: 1,
				totalPages: 3,
			});
			const allButtons = html.match(/<button[\s\S]*?<\/button>/g) ?? [];
			const prevButton = allButtons.find((b) => b.includes("Previous"));
			expect(prevButton).toBeDefined();
			expect(prevButton).toContain("disabled");
		});

		it("disables Next button on last page", async () => {
			const html = await renderComponent(DailySchedule, {
				days: sampleDays,
				page: 3,
				totalPages: 3,
			});
			const allButtons = html.match(/<button[\s\S]*?<\/button>/g) ?? [];
			const nextButton = allButtons.find((b) => b.includes("Next"));
			expect(nextButton).toBeDefined();
			expect(nextButton).toContain("disabled");
		});

		it("renders pagination bar even without callbacks (they are optional emits)", async () => {
			const html = await renderComponent(DailySchedule, {
				days: sampleDays,
				page: 2,
				totalPages: 5,
			});
			expect(html).toContain("Previous");
			expect(html).toContain("Next");
			expect(html).toContain("Page 2 of 5");
		});

		it("does not render pagination bar when page is 0 (no pagination state)", async () => {
			const html = await renderComponent(DailySchedule, {
				days: sampleDays,
				page: 0,
				totalPages: 0,
			});
			expect(html).not.toContain("Previous");
			expect(html).not.toContain("Next");
		});
	});
});
