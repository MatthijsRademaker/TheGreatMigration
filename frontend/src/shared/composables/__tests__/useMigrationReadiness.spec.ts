// @vitest-environment jsdom
import { describe, it, expect, vi } from "vitest";
import { computed, effectScope, ref } from "vue";
import type { ScheduleDay } from "@/calendar/composables/useDailySchedule";

const toolsSummary = ref({ total: 0, claimed: 0 });

vi.mock("@/tools/composables/useTools", () => ({
	useTools: vi.fn(() => ({
		data: computed(() => ({ summary: toolsSummary.value })),
	})),
}));

import { useMigrationReadiness } from "../useMigrationReadiness";

function dayWith(assigned: number, needed: number): ScheduleDay {
	return {
		date: "2026-07-05",
		label: "5 Jul",
		availablePeopleCount: 0,
		tasks: [
			{
				id: "c1",
				title: "Pack",
				priority: "high",
				roomArea: "Kitchen",
				assignedPeople: [],
				peopleNeeded: needed,
				assignedCount: assigned,
				staffingStatus: assigned >= needed ? "fullyStaffed" : "underStaffed",
				completed: false,
				scheduledDate: "2026-07-05",
				taskId: null,
			},
		],
	};
}

function run<T>(fn: () => T): T {
	return effectScope().run(fn) as T;
}

describe("useMigrationReadiness", () => {
	it("averages staffing and tool coverage into a percent", () => {
		toolsSummary.value = { total: 4, claimed: 2 }; // 50%
		const days = ref<ScheduleDay[]>([dayWith(1, 2)]); // staffing 50%
		const { percent, isComplete } = run(() =>
			useMigrationReadiness(computed(() => days.value)),
		);
		expect(percent.value).toBe(50);
		expect(isComplete.value).toBe(false);
	});

	it("reports complete only when every applicable dimension is fully covered", () => {
		toolsSummary.value = { total: 3, claimed: 3 }; // 100%
		const days = ref<ScheduleDay[]>([dayWith(2, 2)]); // staffing 100%
		const { percent, isComplete } = run(() =>
			useMigrationReadiness(computed(() => days.value)),
		);
		expect(percent.value).toBe(100);
		expect(isComplete.value).toBe(true);
	});

	it("excludes dimensions with no data and reports 0% when nothing applies", () => {
		toolsSummary.value = { total: 0, claimed: 0 };
		const days = ref<ScheduleDay[]>([]);
		const { percent, isComplete } = run(() =>
			useMigrationReadiness(computed(() => days.value)),
		);
		expect(percent.value).toBe(0);
		expect(isComplete.value).toBe(false);
	});
});
