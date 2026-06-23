// @vitest-environment jsdom
import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
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
					roomArea: "Kitchen",
					assignedPeople: [{ id: "p2", name: "Jordan Lee", initials: "JL" }],
					peopleNeeded: 2,
					assignedCount: 0,
					staffingStatus: "underStaffed" as const,
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
