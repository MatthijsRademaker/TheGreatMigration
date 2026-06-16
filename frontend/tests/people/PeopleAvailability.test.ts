import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h } from "vue";
import { describe, expect, it } from "vitest";
import PeopleAvailability from "../../src/people/PeopleAvailability.vue";
import type { PeopleAvailabilityProps } from "../../src/people/types";

describe("PeopleAvailability", () => {
	it("renders the panel title", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).toContain("People availability");
	});

	it("renders sample people", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).toContain("Alex");
		expect(html).toContain("Morgan");
		expect(html).toContain("Sam");
		expect(html).toContain("Riley");
	});

	it("renders four day labels", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).toContain(">Mon<");
		expect(html).toContain(">Tue<");
		expect(html).toContain(">Wed<");
		expect(html).toContain(">Thu<");
	});

	it("renders all four status labels in cells", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).toContain(">Available<");
		expect(html).toContain(">Busy<");
		expect(html).toContain(">Partial<");
		expect(html).toContain(">Off<");
	});

	it("renders the legend with all four status badges", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		// Legend Badges with `data-variant` attributes
		expect(html).toContain('data-variant="available"');
		expect(html).toContain('data-variant="busy"');
		expect(html).toContain('data-variant="partial"');
		expect(html).toContain('data-variant="off"');
	});

	it("renders the summary row with availability count", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).toContain("1 of 4 available today");
	});

	it("renders person avatars with initials fallback", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).toContain(">A<");
		expect(html).toContain(">M<");
		expect(html).toContain(">S<");
		expect(html).toContain(">R<");
	});

	it("renders a Card shell with card-title", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).toContain('data-slot="card"');
		expect(html).toContain('data-slot="card-title"');
	});

	it("renders no PersonChip references", async () => {
		const app = createSSRApp({
			render: () => h(PeopleAvailability),
		});
		const html = await renderToString(app);
		expect(html).not.toContain("person-chip");
	});

	it("accepts custom props and renders them", async () => {
		const app = createSSRApp({
			render: () =>
				h(PeopleAvailability, {
					title: "Custom Title",
					days: ["Day 1", "Day 2"],
					people: [
						{
							id: "test",
							name: "Test",
							availability: [
								{ date: "Day 1", status: "available" },
								{ date: "Day 2", status: "off" },
							],
						},
					],
					availableToday: 1,
					totalPeople: 1,
				}),
		});
		const html = await renderToString(app);
		expect(html).toContain("Custom Title");
		expect(html).toContain(">Day 1<");
		expect(html).toContain(">Day 2<");
		expect(html).toContain("Test");
		expect(html).toContain("1 of 1 available today");
	});

	// --- Backend-shaped props tests ---

	it("renders backend-shaped props with ISO-to-label conversion", async () => {
		const props: PeopleAvailabilityProps = {
			title: "People availability",
			description: "Test description",
			days: ["Sun 5 Jul", "Mon 6 Jul"],
			people: [
				{
					id: "p1",
					name: "Sophia Chen",
					availability: [
						{ date: "Sun 5 Jul", status: "available" },
						{ date: "Mon 6 Jul", status: "busy" },
					],
				},
				{
					id: "p7",
					name: "Amara Diallo",
					availability: [
						{ date: "Sun 5 Jul", status: "busy" },
						{ date: "Mon 6 Jul", status: "off" },
					],
				},
			],
			legend: [
				{ id: "available", label: "Available" },
				{ id: "busy", label: "Busy" },
				{ id: "partial", label: "Partial" },
				{ id: "off", label: "Off" },
			],
			availableToday: 1,
			totalPeople: 2,
		};

		const app = createSSRApp({
			render: () => h(PeopleAvailability, props),
		});
		const html = await renderToString(app);
		expect(html).toContain("Sophia Chen");
		expect(html).toContain("Amara Diallo");
		expect(html).toContain("1 of 2 available today");
		expect(html).toContain('data-variant="available"');
		expect(html).toContain('data-variant="busy"');
		expect(html).toContain('data-variant="off"');
	});

	it("renders empty people list without errors", async () => {
		const props: PeopleAvailabilityProps = {
			days: [],
			people: [],
			legend: [],
			availableToday: 0,
			totalPeople: 0,
		};

		const app = createSSRApp({
			render: () => h(PeopleAvailability, props),
		});
		const html = await renderToString(app);
		expect(html).toContain("0 of 0 available today");
		expect(html).toContain("People availability");
	});

	it("renders canonical statuses only", async () => {
		const props: PeopleAvailabilityProps = {
			days: ["Day 1"],
			people: [
				{
					id: "p1",
					name: "Test",
					availability: [{ date: "Day 1", status: "available" }],
				},
			],
			legend: [
				{ id: "available", label: "Available" },
				{ id: "busy", label: "Busy" },
				{ id: "partial", label: "Partial" },
				{ id: "off", label: "Off" },
			],
			availableToday: 1,
			totalPeople: 1,
		};

		const app = createSSRApp({
			render: () => h(PeopleAvailability, props),
		});
		const html = await renderToString(app);
		expect(html).toContain('data-variant="available"');
		expect(html).not.toContain('data-variant="unknown"');
	});
});
