// @vitest-environment jsdom

import { mount } from "@vue/test-utils";
import { createPinia } from "pinia";
import { PiniaColada } from "@pinia/colada";
import { nextTick, computed, ref, type Ref } from "vue";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { CalendarDate, type DateValue } from "@internationalized/date";
import SettingsView from "../src/settings/SettingsView.vue";

// ---------------------------------------------------------------------------
// Controlled reactive state — declared before vi.mock so the hoisted
// mock factory closures capture the module-level bindings.
// ---------------------------------------------------------------------------
let controlledState: {
	startDate: DateValue | null;
	endDate: DateValue | null;
	originalStart: string;
	originalEnd: string;
	isLoading: Ref<boolean>;
	isError: Ref<boolean>;
};
let mutateSpy: ReturnType<typeof vi.fn>;
let controlledMutationError: Ref<Error | null>;
let controlledIsPending: Ref<boolean>;

// ---------------------------------------------------------------------------
// Mock the generated Pinia Colada query module.
// ---------------------------------------------------------------------------
vi.mock("@/client/@pinia/colada.gen", () => ({
	getPlanningWindowQueryKey: () => ["planning-window"],
	getPlanningWindowQuery: () => ({
		key: ["planning-window"],
		query: async () => ({
			startDate: "2026-07-05",
			endDate: "2026-08-13",
			days: 40,
		}),
	}),
	putPlanningWindowMutation: () => ({
		mutation: async () => ({
			startDate: "2026-07-05",
			endDate: "2026-08-13",
			days: 40,
		}),
	}),
}));

// ---------------------------------------------------------------------------
// Mock @pinia/colada's useQuery to return controlled data.
// ---------------------------------------------------------------------------
vi.mock("@pinia/colada", async () => {
	const actual = await vi.importActual("@pinia/colada");
	return {
		...actual,
		useQuery: () => ({
			data: computed(() =>
				controlledState.isError.value
					? null
					: controlledState.isLoading.value
						? undefined
						: {
								startDate: controlledState.originalStart,
								endDate: controlledState.originalEnd,
								days: 40,
							},
			),
			isPending: computed(() => controlledState.isLoading.value),
			error: computed(() =>
				controlledState.isError.value ? new Error("fail") : null,
			),
			status: computed(() =>
				controlledState.isLoading.value
					? "pending"
					: controlledState.isError.value
						? "error"
						: "success",
			),
		}),
	};
});

// ---------------------------------------------------------------------------
// Mock useUpdatePlanningWindow composable.
// ---------------------------------------------------------------------------
vi.mock("@/shared/composables/useUpdatePlanningWindow", () => ({
	useUpdatePlanningWindow: () => ({
		mutate: mutateSpy,
		isPending: computed(() => controlledIsPending.value),
		error: computed(() => controlledMutationError.value),
	}),
}));

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
interface SettingsViewVM {
	startDate: DateValue | null;
	endDate: DateValue | null;
	originalStart: string;
	originalEnd: string;
	isEndBeforeStart: boolean;
	isDirty: boolean;
	canSave: boolean;
}

function sv(wrapper: ReturnType<typeof mount>): SettingsViewVM {
	return wrapper.vm as unknown as SettingsViewVM;
}

async function mountSettingsView() {
	const wrapper = mount(SettingsView, {
		global: {
			plugins: [createPinia(), PiniaColada],
		},
	});

	await nextTick();
	await nextTick();

	return wrapper;
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
describe("SettingsView interactive", () => {
	beforeEach(() => {
		controlledState = {
			startDate: new CalendarDate(2026, 7, 5),
			endDate: new CalendarDate(2026, 8, 13),
			originalStart: "2026-07-05",
			originalEnd: "2026-08-13",
			isLoading: ref(false),
			isError: ref(false),
		};
		mutateSpy = vi.fn();
		controlledMutationError = ref(null);
		controlledIsPending = ref(false);
	});

	it("renders the planning window card with prefilled dates", async () => {
		const wrapper = await mountSettingsView();

		expect(wrapper.html()).toContain("Planning window");
		expect(wrapper.html()).toContain("Start date");
		expect(wrapper.html()).toContain("End date");
	});

	it("shows validation error when endDate is before startDate", async () => {
		const wrapper = await mountSettingsView();

		sv(wrapper).startDate = new CalendarDate(2026, 8, 15);
		sv(wrapper).endDate = new CalendarDate(2026, 7, 10);
		await nextTick();

		expect(wrapper.html()).toContain(
			"End date must be on or after the start date",
		);
	});

	it("disables Save when endDate is before startDate", async () => {
		const wrapper = await mountSettingsView();

		sv(wrapper).startDate = new CalendarDate(2026, 8, 15);
		sv(wrapper).endDate = new CalendarDate(2026, 7, 10);
		await nextTick();

		expect(sv(wrapper).canSave).toBe(false);
	});

	it("enables Save when dates are valid and dirty", async () => {
		const wrapper = await mountSettingsView();

		expect(sv(wrapper).originalStart).toBe("2026-07-05");
		expect(sv(wrapper).originalEnd).toBe("2026-08-13");

		sv(wrapper).startDate = new CalendarDate(2026, 7, 6);
		await nextTick();

		expect(sv(wrapper).isDirty).toBe(true);
		expect(sv(wrapper).canSave).toBe(true);
	});

	it("disables Save when form is clean", async () => {
		const wrapper = await mountSettingsView();

		expect(sv(wrapper).isDirty).toBe(false);
		expect(sv(wrapper).canSave).toBe(false);
	});

	it("calls mutate when Save is clicked", async () => {
		const wrapper = await mountSettingsView();

		sv(wrapper).startDate = new CalendarDate(2026, 7, 10);
		sv(wrapper).endDate = new CalendarDate(2026, 7, 20);
		await nextTick();

		expect(sv(wrapper).canSave).toBe(true);

		const buttons = wrapper.findAll("button");
		const saveBtn = buttons.find(
			(b) => b.text() === "Save" || b.text().includes("Saving"),
		);
		expect(saveBtn).toBeTruthy();
		await saveBtn!.trigger("click");

		expect(mutateSpy).toHaveBeenCalledWith({
			body: {
				startDate: "2026-07-10",
				endDate: "2026-07-20",
			},
		});
	});

	it("displays error message when mutation fails", async () => {
		controlledMutationError.value = new Error("Server error");
		const wrapper = await mountSettingsView();

		expect(wrapper.html()).toContain("Failed to save changes");
	});

	it("Reset restores original values after editing", async () => {
		const wrapper = await mountSettingsView();

		expect(sv(wrapper).originalStart).toBe("2026-07-05");
		expect(sv(wrapper).originalEnd).toBe("2026-08-13");

		sv(wrapper).startDate = new CalendarDate(2026, 8, 1);
		sv(wrapper).endDate = new CalendarDate(2026, 8, 5);
		await nextTick();

		expect(sv(wrapper).isDirty).toBe(true);

		const buttons = wrapper.findAll("button");
		const resetBtn = buttons.find((b) => b.text() === "Reset");
		expect(resetBtn).toBeTruthy();

		await resetBtn!.trigger("click");
		await nextTick();

		expect(sv(wrapper).startDate!.toString()).toBe("2026-07-05");
		expect(sv(wrapper).endDate!.toString()).toBe("2026-08-13");
	});

	it("disables buttons while mutation is pending", async () => {
		const wrapper = await mountSettingsView();

		sv(wrapper).startDate = new CalendarDate(2026, 7, 10);
		sv(wrapper).endDate = new CalendarDate(2026, 7, 20);
		await nextTick();

		controlledIsPending.value = true;
		await nextTick();

		const buttons = wrapper.findAll("button");
		const resetBtn = buttons.find((b) => b.text() === "Reset");
		const saveBtn = buttons.find(
			(b) => b.text().includes("Saving") || b.text() === "Save",
		);

		expect(resetBtn).toBeTruthy();
		expect(saveBtn).toBeTruthy();
		expect(resetBtn!.attributes("disabled")).toBeDefined();
		expect(saveBtn!.attributes("disabled")).toBeDefined();
	});

	it("shows saving text on button while mutation is pending", async () => {
		const wrapper = await mountSettingsView();

		sv(wrapper).startDate = new CalendarDate(2026, 7, 10);
		sv(wrapper).endDate = new CalendarDate(2026, 7, 20);
		await nextTick();

		controlledIsPending.value = true;
		await nextTick();

		const buttons = wrapper.findAll("button");
		const saveBtn = buttons.find((b) => b.text().includes("Saving"));
		expect(saveBtn).toBeTruthy();
	});

	it("shows loading state when query is pending", async () => {
		controlledState.isLoading.value = true;
		const wrapper = await mountSettingsView();

		expect(wrapper.html()).toContain("Loading planning window");
	});

	it("shows error state when query fails", async () => {
		controlledState.isError.value = true;
		const wrapper = await mountSettingsView();

		expect(wrapper.html()).toContain("Failed to load planning window");
	});
});
