import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h, nextTick } from "vue";
import { describe, expect, it, vi } from "vitest";
import { ref, computed } from "vue";

// --- Mock Pinia Colada ---
vi.mock("@pinia/colada", () => ({
	useMutation: () => ({
		mutateAsync: vi.fn(),
		mutate: vi.fn(),
		isLoading: ref(false),
		error: ref(null),
		data: ref(undefined),
		status: ref("idle"),
		asyncStatus: ref("idle"),
		state: computed(() => ({ status: "idle" as const })),
		reset: vi.fn(),
	}),
	useQueryCache: () => ({
		invalidateQueries: vi.fn(),
	}),
}));

// --- Mock generated client ---
vi.mock("@/client/@pinia/colada.gen", () => ({
	getDashboardPeopleAvailabilityQueryKey: () => [],
	getDashboardPeopleAvailabilityQuery: () => ({}),
	createPersonMutation: () => ({ mutation: vi.fn() }),
	deletePersonMutation: () => ({ mutation: vi.fn() }),
	upsertPersonAvailabilityMutation: () => ({ mutation: vi.fn() }),
	deletePersonAvailabilityMutation: () => ({ mutation: vi.fn() }),
}));

// --- Mock composable ---
const mockPeople = [
	{
		id: "p1",
		name: "Sophia Chen",
		availability: [
			{ date: "Sun 5 Jul", status: "available" as const },
			{ date: "Mon 6 Jul", status: "busy" as const },
		],
	},
	{
		id: "p2",
		name: "Marcus Rivera",
		availability: [
			{ date: "Sun 5 Jul", status: "off" as const },
			{ date: "Mon 6 Jul", status: "available" as const },
		],
	},
];

vi.mock("@/shared/composables/usePeopleAvailability", () => ({
	usePeopleAvailability: () => ({
		data: ref({
			title: "People availability",
			description: "Test desc",
			days: ["Sun 5 Jul", "Mon 6 Jul"],
			people: mockPeople,
			legend: [
				{ id: "available", label: "Available" },
				{ id: "busy", label: "Busy" },
				{ id: "partial", label: "Partial" },
				{ id: "off", label: "Off" },
			],
			availableToday: 1,
			totalPeople: 2,
		}),
		rawData: ref({
			range: {
				startDate: "2026-07-05",
				endDate: "2026-07-06",
				days: 2,
				selectedDate: "2026-07-05",
			},
			summary: { availableToday: 1, totalPeople: 2 },
			people: [],
			statuses: [],
		}),
		isLoading: ref(false),
		isError: ref(false),
		isEmpty: ref(false),
		refresh: vi.fn(),
		queryKey: [],
	}),
}));

import PeopleView from "../../src/people/PeopleView.vue";

async function renderPeopleView() {
	const app = createSSRApp({
		render: () => h(PeopleView),
	});
	return renderToString(app);
}

describe("PeopleView management controls", () => {
	it("renders the create person form with input fields and submit button", async () => {
		const html = await renderPeopleView();

		expect(html).toContain("Add a person");
		expect(html).toContain("Create a new person with a short unique ID.");
		expect(html).toContain(">Create<");
		// Input fields with placeholders
		expect(html).toContain("e.g. p9");
		expect(html).toContain("Full name");
		expect(html).toContain("e.g. JS");
	});

	it("renders the editable PeopleAvailability matrix with person names and Delete buttons", async () => {
		const html = await renderPeopleView();

		// Person names are rendered in the matrix
		expect(html).toContain("Sophia Chen");
		expect(html).toContain("Marcus Rivera");

		// The matrix renders full day labels
		expect(html).toContain("Sun 5 Jul");
		expect(html).toContain("Mon 6 Jul");

		// Actions column is present (editable mode)
		expect(html).toContain("Actions");

		// Each person has a Delete button (2 people = 2 Delete buttons)
		const deleteButtonCount = (html.match(/>Delete</g) || []).length;
		expect(deleteButtonCount).toBe(2);
	});

	it("renders status badges in the editable matrix", async () => {
		const html = await renderPeopleView();

		// Status badges with data-variant attributes (now from the matrix, not the manage card)
		expect(html).toContain('data-variant="available"');
		expect(html).toContain('data-variant="busy"');
		expect(html).toContain('data-variant="off"');

		// Status labels
		expect(html).toContain(">Available<");
		expect(html).toContain(">Busy<");
		expect(html).toContain(">Off<");
	});

	it("does not render the old Manage people card", async () => {
		const html = await renderPeopleView();

		// The old manage card heading and description should not appear
		expect(html).not.toContain("Manage people");
	});

	it("does not show error messages in initial success state", async () => {
		const html = await renderPeopleView();

		// Error messages should not be present
		expect(html).not.toContain("Failed to create person");
		expect(html).not.toContain("Failed to delete person");
		expect(html).not.toContain("Failed to update status");
		expect(html).not.toContain("Failed to update availability");
		expect(html).not.toContain("Failed to clear availability");
	});
	it("does not show loading or error states when data is present", async () => {
		const html = await renderPeopleView();

		expect(html).not.toContain("Loading availability data");
		expect(html).not.toContain("Backend unavailable");
		expect(html).not.toContain("No people found");
	});
});

// --- Client-render tests for loading state passing ---
// @vitest-environment jsdom

describe("PeopleView loading state wiring", () => {
	it("passes deletingPersonId from handleDelete to the matrix", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		// The matrix should have a PeopleAvailability component
		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		expect(matrix.exists()).toBe(true);

		// deletingPersonId prop should be undefined initially (no delete in progress)
		expect(matrix.props("deletingPersonId")).toBeNull();

		// updating prop should be falsy (no upsert in progress)
		expect(matrix.props("updating")).toBe(false);

		wrapper.unmount();
	});

	it("passes editable=true to the matrix on success", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		expect(matrix.exists()).toBe(true);
		expect(matrix.props("editable")).toBe(true);

		wrapper.unmount();
	});
});
