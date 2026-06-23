// @vitest-environment jsdom
import { describe, it, expect, vi } from "vitest";
import { effectScope, ref } from "vue";

const queryData = ref<any>(null);

vi.mock("@pinia/colada", () => ({
	defineQueryOptions: vi.fn((options: unknown) => options),
	useQuery: vi.fn(() => ({
		data: queryData,
		isPending: ref(false),
		error: ref(null),
		refetch: vi.fn(),
	})),
}));

import { useDailySchedule } from "../useDailySchedule";

function run<T>(fn: () => T): T {
	return effectScope().run(fn) as T;
}

describe("useDailySchedule", () => {
	it("sorts tasks by priority inside each day", () => {
		queryData.value = {
			days: [
				{
					date: "2026-07-05",
					label: "5 Jul",
					availablePeopleCount: 3,
					tasks: [
						{
							id: "low",
							title: "Low",
							priority: "low",
							roomArea: "Kitchen",
							assignedPeople: [],
							peopleNeeded: 1,
							assignedCount: 0,
							staffingStatus: "underStaffed",
							scheduledDate: "2026-07-05",
							taskId: null,
						},
						{
							id: "high",
							title: "High",
							priority: "high",
							roomArea: "Kitchen",
							assignedPeople: [],
							peopleNeeded: 1,
							assignedCount: 0,
							staffingStatus: "underStaffed",
							scheduledDate: "2026-07-05",
							taskId: null,
						},
						{
							id: "medium",
							title: "Medium",
							priority: "medium",
							roomArea: "Kitchen",
							assignedPeople: [],
							peopleNeeded: 1,
							assignedCount: 0,
							staffingStatus: "underStaffed",
							scheduledDate: "2026-07-05",
							taskId: null,
						},
					],
				},
			],
		};

		const { data } = run(() => useDailySchedule({ start: "2026-07-05" }));

		expect(data.value.days[0].tasks.map((task) => task.id)).toEqual([
			"high",
			"medium",
			"low",
		]);
	});
});
