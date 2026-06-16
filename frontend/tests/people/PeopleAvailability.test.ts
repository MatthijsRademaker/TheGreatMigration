import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h, nextTick } from "vue";
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

// --- Editable mode tests (client-render) ---
// @vitest-environment jsdom

describe("PeopleAvailability editable mode", () => {
	it("renders Popover triggers in editable mode", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleAvailability, {
			props: {
				editable: true,
			},
			attachTo: document.body,
		});

		// Popover triggers should exist
		const triggers = wrapper.findAll('[data-slot="popover-trigger"]');
		// 4 people × 4 days = 16 triggers
		expect(triggers.length).toBe(16);

		// Actions column header should be present
		expect(wrapper.text()).toContain("Actions");

		// Delete buttons should be present (one per person)
		const deleteButtons = wrapper
			.findAll("button")
			.filter((b) => b.text() === "Delete");
		expect(deleteButtons.length).toBe(4);

		wrapper.unmount();
	});

	it("emits update-cell when a status is selected in the popover", async () => {
		const { mount } = await import("@vue/test-utils");
		const props = {
			editable: true,
			days: ["Mon", "Tue"],
			people: [
				{
					id: "p1",
					name: "Test Person",
					availability: [
						{ date: "Mon", status: "available" as const },
						{ date: "Tue", status: "busy" as const },
					],
				},
			],
			availableToday: 0,
			totalPeople: 1,
		};

		const wrapper = mount(PeopleAvailability, {
			props,
			attachTo: document.body,
		});

		// Click the first popover trigger (person p1, day 0 = Mon, status="available")
		const firstTrigger = wrapper.find('[data-slot="popover-trigger"]');
		expect(firstTrigger.exists()).toBe(true);
		await firstTrigger.trigger("click");
		await nextTick();

		// The popover content is teleported to document.body.
		const popoverContent = document.body.querySelector(
			'[data-slot="popover-content"]',
		);
		expect(popoverContent).not.toBeNull();

		// Find all option buttons inside popover content (excluding the Clear button which is last)
		const optionButtons = Array.from(
			popoverContent?.querySelectorAll("button") ?? [],
		).filter((button) => {
			const text = button.textContent?.trim();
			return (
				text === "Available" ||
				text === "Busy" ||
				text === "Partial" ||
				text === "Off"
			);
		});
		expect(optionButtons.length).toBeGreaterThanOrEqual(2);

		// Click the "Busy" option
		const busyButton = optionButtons.find(
			(button) => button.textContent?.trim() === "Busy",
		);
		expect(busyButton).toBeDefined();
		busyButton?.click();
		await nextTick();

		const emitted = wrapper.emitted("update-cell");
		expect(emitted).toBeTruthy();
		if (emitted) {
			const payload = emitted[0][0] as {
				personId: string;
				dayIndex: number;
				status: string | null;
			};
			expect(payload.personId).toBe("p1");
			expect(payload.dayIndex).toBe(0);
			expect(payload.status).toBe("busy");
		}

		wrapper.unmount();
	});

	it("emits update-cell with status null when Clear is selected", async () => {
		const { mount } = await import("@vue/test-utils");
		const props = {
			editable: true,
			days: ["Mon"],
			people: [
				{
					id: "p1",
					name: "Test",
					availability: [{ date: "Mon", status: "available" as const }],
				},
			],
			availableToday: 1,
			totalPeople: 1,
		};

		const wrapper = mount(PeopleAvailability, {
			props,
			attachTo: document.body,
		});

		// Click the trigger to open popover
		const trigger = wrapper.find('[data-slot="popover-trigger"]');
		expect(trigger.exists()).toBe(true);
		await trigger.trigger("click");
		await nextTick();

		// Find the Clear button by text. The popover content is teleported to document.body.
		const popoverContent = document.body.querySelector(
			'[data-slot="popover-content"]',
		);
		expect(popoverContent).not.toBeNull();

		const clearButton = Array.from(
			popoverContent?.querySelectorAll("button") ?? [],
		).filter((button) => button.textContent?.includes("Clear"));
		expect(clearButton.length).toBe(1);
		clearButton[0]?.click();
		await nextTick();

		const emitted = wrapper.emitted("update-cell");
		expect(emitted).toBeTruthy();
		if (emitted) {
			const payload = emitted[0][0] as {
				personId: string;
				dayIndex: number;
				status: string | null;
			};
			expect(payload.personId).toBe("p1");
			expect(payload.dayIndex).toBe(0);
			expect(payload.status).toBeNull();
		}

		wrapper.unmount();
	});

	it("emits delete-person when delete button is clicked", async () => {
		const { mount } = await import("@vue/test-utils");
		const props = {
			editable: true,
			days: ["Mon"],
			people: [
				{
					id: "p1",
					name: "Test",
					availability: [{ date: "Mon", status: "available" as const }],
				},
			],
			availableToday: 1,
			totalPeople: 1,
		};

		const wrapper = mount(PeopleAvailability, {
			props,
			attachTo: document.body,
		});

		const deleteButton = wrapper
			.findAll("button")
			.filter((b) => b.text() === "Delete");
		expect(deleteButton.length).toBe(1);
		await deleteButton[0].trigger("click");

		const emitted = wrapper.emitted("delete-person");
		expect(emitted).toBeTruthy();
		if (emitted) {
			expect(emitted[0][0]).toBe("p1");
		}

		wrapper.unmount();
	});

	it("does not render Popover triggers without editable prop", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleAvailability, {
			props: {},
			attachTo: document.body,
		});

		const triggers = wrapper.findAll('[data-slot="popover-trigger"]');
		expect(triggers.length).toBe(0);

		// No Actions column header
		expect(wrapper.text()).not.toContain("Actions");

		// No Delete buttons (read-only)
		const deleteButtons = wrapper
			.findAll("button")
			.filter((b) => b.text() === "Delete");
		expect(deleteButtons.length).toBe(0);

		wrapper.unmount();
	});

	it("does not render Popover triggers when editable is false", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleAvailability, {
			props: { editable: false },
			attachTo: document.body,
		});

		const triggers = wrapper.findAll('[data-slot="popover-trigger"]');
		expect(triggers.length).toBe(0);

		wrapper.unmount();
	});
});
