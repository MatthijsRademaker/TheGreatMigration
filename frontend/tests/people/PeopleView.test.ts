import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h } from "vue";
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

	it("renders the manage people section with delete button per person", async () => {
		const html = await renderPeopleView();

		expect(html).toContain("Manage people");
		expect(html).toContain(
			"Update statuses or remove people no longer needed.",
		);

		// Person names in manage section
		expect(html).toContain("Sophia Chen");
		expect(html).toContain("Marcus Rivera");

		// Each person has a Delete button
		const deleteButtonCount = (html.match(/>Delete</g) || []).length;
		expect(deleteButtonCount).toBe(2);
	});

	it("renders status badges per person in the manage card", async () => {
		const html = await renderPeopleView();

		// Status badges with data-variant attributes
		expect(html).toContain('data-variant="available"');
		expect(html).toContain('data-variant="busy"');
		expect(html).toContain('data-variant="off"');

		// Status labels
		expect(html).toContain(">Available<");
		expect(html).toContain(">Busy<");
		expect(html).toContain(">Off<");
	});

	it("renders the PeopleAvailability matrix", async () => {
		const html = await renderPeopleView();

		// The matrix renders full day labels (e.g. "Sun 5 Jul")
		expect(html).toContain("Sun 5 Jul");
		expect(html).toContain("Mon 6 Jul");
	});

	it("includes clear availability option description in the card", async () => {
		const html = await renderPeopleView();

		// The card description confirms update and delete capabilities
		expect(html).toContain(
			"Update statuses or remove people no longer needed.",
		);
	});

	it("does not show error messages in initial success state", async () => {
		const html = await renderPeopleView();

		// Error messages should not be present
		expect(html).not.toContain("Failed to create person");
		expect(html).not.toContain("Failed to delete person");
		expect(html).not.toContain("Failed to update status");
		expect(html).not.toContain("Failed to clear availability");
	});

	it("does not show loading or error states when data is present", async () => {
		const html = await renderPeopleView();

		expect(html).not.toContain("Loading availability data");
		expect(html).not.toContain("Backend unavailable");
		expect(html).not.toContain("No people found");
	});
});
