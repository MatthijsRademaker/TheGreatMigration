// @vitest-environment jsdom
import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { nextTick } from "vue";
import DailySchedule from "../DailySchedule.vue";

function days() {
	return [
		{
			date: "2026-07-05",
			label: "5 Jul",
			availablePeopleCount: 3,
			tasks: [
				{
					id: "card-1",
					title: "Pack kitchen",
					priority: "high" as const,
					area: { id: "room-1", name: "Kitchen" },
					assignedPeople: [{ id: "p2", name: "Jordan Lee", initials: "JL" }],
					peopleNeeded: 2,
					assignedCount: 0,
					staffingStatus: "underStaffed" as const,
					completed: false,
					scheduledDate: "2026-07-05",
					taskId: "task-1",
				},
			],
		},
		{
			date: "2026-07-06",
			label: "6 Jul",
			availablePeopleCount: 2,
			tasks: [],
		},
	];
}

const people = [{ id: "p1", name: "Alex Kim" }];

describe("DailySchedule drag-and-drop", () => {
	it("renders a draggable people rail when interactive", () => {
		const wrapper = mount(DailySchedule, { props: { days: days(), people } });
		const rail = wrapper.find("[data-testid='people-rail']");
		expect(rail.exists()).toBe(true);
		expect(wrapper.findAll("[data-testid='rail-person']")).toHaveLength(1);
		wrapper.unmount();
	});

	it("hides entire people rail on mobile", () => {
		const wrapper = mount(DailySchedule, { props: { days: days(), people } });
		const rail = wrapper.find("[data-testid='people-rail']");
		expect(rail.exists()).toBe(true);
		expect(rail.classes()).toContain("hidden");
		expect(rail.classes()).toContain("sm:flex");
		wrapper.unmount();
	});

	it("shows assigned people names on task cards", () => {
		const wrapper = mount(DailySchedule, { props: { days: days(), people } });
		const card = wrapper.find("[data-slot='task-board-card']");
		expect(card.text()).toContain("Jordan Lee");
		expect(card.text()).not.toContain("JL");
		wrapper.unmount();
	});

	it("emits assign-person when a person is dropped on a task card", async () => {
		const wrapper = mount(DailySchedule, { props: { days: days(), people } });
		await wrapper.find("[data-testid='rail-person']").trigger("dragstart");
		await wrapper.find("[data-slot='task-board-card']").trigger("dragover");
		await wrapper.find("[data-slot='task-board-card']").trigger("drop");

		const emitted = wrapper.emitted("assign-person");
		expect(emitted).toBeTruthy();
		expect(emitted![0]).toEqual(["card-1", { id: "p1", name: "Alex Kim" }]);
		wrapper.unmount();
	});

	it("emits reschedule-card when a card is dropped on a different day column", async () => {
		const wrapper = mount(DailySchedule, { props: { days: days(), people } });
		await wrapper.find("[data-slot='task-board-card']").trigger("dragstart");
		const columns = wrapper.findAll("[data-testid='day-column']");
		await columns[1].trigger("dragover");
		await columns[1].trigger("drop");

		const emitted = wrapper.emitted("reschedule-card");
		expect(emitted).toBeTruthy();
		expect(emitted![0]).toEqual(["card-1", "2026-07-06"]);
		wrapper.unmount();
	});

	it("uses compact date label on mobile and hides page indicator there", () => {
		const wrapper = mount(DailySchedule, {
			props: {
				days: days(),
				people,
				page: 1,
				totalPages: 3,
				dateRangeLabel: "5 Jul (Sun) – 6 Jul (Mon)",
			},
		});
		const fullLabel = wrapper.find("[data-testid='date-range-label-full']");
		const compactLabel = wrapper.find(
			"[data-testid='date-range-label-compact']",
		);
		const pageIndicator = wrapper.find("[data-testid='page-indicator']");
		expect(fullLabel.classes()).toContain("hidden");
		expect(fullLabel.classes()).toContain("sm:inline");
		expect(compactLabel.classes()).toContain("sm:hidden");
		expect(compactLabel.text()).toBe("5–6 Jul");
		expect(pageIndicator.classes()).toContain("hidden");
		expect(pageIndicator.classes()).toContain("sm:inline");
		wrapper.unmount();
	});

	it("plays celebration burst, then greys task out after Done click", async () => {
		vi.useFakeTimers();
		const originalMatchMedia = window.matchMedia;
		window.matchMedia = vi.fn().mockImplementation((query: string) => ({
			matches: false,
			media: query,
			onchange: null,
			addEventListener: vi.fn(),
			removeEventListener: vi.fn(),
			addListener: vi.fn(),
			removeListener: vi.fn(),
			dispatchEvent: vi.fn(),
		}));
		const wrapper = mount(DailySchedule, { props: { days: days(), people } });

		await wrapper.find("[data-testid='done-button']").trigger("click");
		await nextTick();
		expect(wrapper.find("[data-testid='celebration']").exists()).toBe(true);
		expect(wrapper.findAll("[data-slot='task-board-card']")).toHaveLength(1);

		await vi.advanceTimersByTimeAsync(120);
		await nextTick();
		expect(wrapper.findAll("[data-slot='task-board-card']")).toHaveLength(0);

		await vi.advanceTimersByTimeAsync(220);
		await nextTick();

		const card = wrapper.find("[data-slot='task-board-card']");
		expect(card.exists()).toBe(true);
		expect(card.attributes("data-done")).toBe("true");
		expect(card.attributes("draggable")).toBeUndefined();
		expect(card.text()).toContain("Revert");

		wrapper.unmount();
		window.matchMedia = originalMatchMedia;
		vi.useRealTimers();
	});

	it("exposes no rail or draggable cards in read-only mode", () => {
		const wrapper = mount(DailySchedule, {
			props: { days: days(), people, readOnly: true },
		});
		expect(wrapper.find("[data-testid='people-rail']").exists()).toBe(false);
		const card = wrapper.find("[data-slot='task-board-card']");
		expect(card.attributes("draggable")).toBeUndefined();
		wrapper.unmount();
	});
});
