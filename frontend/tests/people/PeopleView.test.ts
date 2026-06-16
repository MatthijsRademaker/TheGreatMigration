import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h, nextTick } from "vue";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ref, computed } from "vue";

// --- Shared mutation mock functions (controllable from tests) ---
const mutationMockState = {
	mutateAsync: vi
		.fn<(...args: unknown[]) => Promise<unknown>>()
		.mockResolvedValue(undefined),
};

// --- Mock Pinia Colada ---
vi.mock("@pinia/colada", () => ({
	useMutation: () => ({
		mutateAsync: mutationMockState.mutateAsync,
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
		daysISO: ref(["2026-07-05", "2026-07-06"]),
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
		expect(html).toContain(
			"Create a new person. The ID is assigned server-side.",
		);
		expect(html).toContain(">Create<");
		// Input fields with placeholders
		expect(html).toContain("Full name");
		expect(html).toContain("e.g. JS");
		// ID input should not be present
		expect(html).not.toContain("e.g. p9");
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

// --- Client-render tests for handleCellUpdate ---
// @vitest-environment jsdom

describe("PeopleView handleCellUpdate", () => {
	beforeEach(() => {
		mutationMockState.mutateAsync.mockReset();
		mutationMockState.mutateAsync.mockResolvedValue(undefined);
	});

	it("derives correct date from daysISO by index", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 0,
			status: "busy",
		});
		await nextTick();

		// The mutation should be called with the ISO date from daysISO[0] ("2026-07-05")
		expect(mutationMockState.mutateAsync).toHaveBeenCalledWith(
			expect.objectContaining({
				path: expect.objectContaining({ id: "p1", date: "2026-07-05" }),
				body: expect.objectContaining({ status: "busy" }),
			}),
		);

		wrapper.unmount();
	});

	it("derives date from daysISO for non-zero index", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p2",
			dayIndex: 1,
			status: "available",
		});
		await nextTick();

		// The mutation should be called with daysISO[1] ("2026-07-06")
		expect(mutationMockState.mutateAsync).toHaveBeenCalledWith(
			expect.objectContaining({
				path: expect.objectContaining({ id: "p2", date: "2026-07-06" }),
			}),
		);

		wrapper.unmount();
	});

	it("shows error without API call when dayIndex is out of range", async () => {
		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 10,
			status: "busy",
		});
		await nextTick();

		// Error message without API call
		expect(wrapper.text()).toContain(
			"Selected cell cannot be mapped to a date",
		);
		expect(mutationMockState.mutateAsync).not.toHaveBeenCalled();

		wrapper.unmount();
	});

	it("shows planning window error for 400 outside planning window", async () => {
		mutationMockState.mutateAsync.mockRejectedValueOnce({
			status: 400,
			cause: {
				body: { detail: "date 2026-07-05 is outside the planning window" },
			},
		});

		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 0,
			status: "busy",
		});

		// Wait for the async mutation rejection to propagate through Vue's reactivity
		await new Promise((r) => setTimeout(r, 0));
		await nextTick();

		expect(wrapper.text()).toContain(
			"This date is outside the planning window.",
		);

		wrapper.unmount();
	});

	it("shows invalid date format error for 400 with invalid ISO date", async () => {
		mutationMockState.mutateAsync.mockRejectedValueOnce({
			status: 400,
			cause: {
				body: {
					detail:
						"must be a valid ISO 8601 date: invalid-date' does not match pattern '^\\d{4}-\\d{2}-\\d{2}$'",
				},
			},
		});

		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 0,
			status: "busy",
		});

		// Wait for the async mutation rejection to propagate
		await new Promise((r) => setTimeout(r, 0));
		await nextTick();

		expect(wrapper.text()).toContain("Invalid date format.");

		wrapper.unmount();
	});

	it("shows invalid status error for 400 with invalid status value", async () => {
		mutationMockState.mutateAsync.mockRejectedValueOnce({
			status: 400,
			cause: {
				body: {
					detail: "status must be one of [available, busy, partial, off]",
				},
			},
		});

		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 0,
			status: "busy",
		});

		// Wait for the async mutation rejection to propagate
		await new Promise((r) => setTimeout(r, 0));
		await nextTick();

		expect(wrapper.text()).toContain("Invalid status value.");

		wrapper.unmount();
	});

	it("shows person not found for 404 error", async () => {
		mutationMockState.mutateAsync.mockRejectedValueOnce({
			status: 404,
		});

		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 0,
			status: "busy",
		});

		// Wait for the async mutation rejection to propagate
		await new Promise((r) => setTimeout(r, 0));
		await nextTick();

		expect(wrapper.text()).toContain("Person not found.");

		wrapper.unmount();
	});

	it("shows fallback for generic 400 error without known detail", async () => {
		mutationMockState.mutateAsync.mockRejectedValueOnce({
			status: 400,
			cause: {
				body: { detail: "some unexpected validation error" },
			},
		});

		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 0,
			status: "busy",
		});

		// Wait for the async mutation rejection to propagate
		await new Promise((r) => setTimeout(r, 0));
		await nextTick();

		expect(wrapper.text()).toContain(
			"Failed to update: some unexpected validation error",
		);

		wrapper.unmount();
	});

	it("shows generic error for non-HTTP errors", async () => {
		mutationMockState.mutateAsync.mockRejectedValueOnce(
			new Error("Network failure"),
		);

		const { mount } = await import("@vue/test-utils");
		const wrapper = mount(PeopleView, {
			attachTo: document.body,
		});

		await nextTick();

		const matrix = wrapper.findComponent({ name: "PeopleAvailability" });
		matrix.vm.$emit("update-cell", {
			personId: "p1",
			dayIndex: 0,
			status: "busy",
		});

		// Wait for the async mutation rejection to propagate
		await new Promise((r) => setTimeout(r, 0));
		await nextTick();

		expect(wrapper.text()).toContain(
			"Failed to update availability: Network failure",
		);

		wrapper.unmount();
	});
});
