// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { computed, effectScope, ref } from "vue";
import type { ScheduleDay } from "../useDailySchedule";

// ── Mocks ─────────────────────────────────────────────────────────────────
let mutateAsync: ReturnType<typeof vi.fn>;
const invalidateQueries = vi.fn();

vi.mock("@pinia/colada", () => ({
	useMutation: vi.fn(() => ({ mutateAsync })),
	useQueryCache: vi.fn(() => ({ invalidateQueries })),
}));

vi.mock("@/client/@pinia/colada.gen", () => ({
	updateScheduleCardMutation: vi.fn(() => ({ mutation: vi.fn() })),
}));

import { useScheduleBoardDnd } from "../useScheduleBoardDnd";

// ── Fixtures ────────────────────────────────────────────────────────────────
function makeDays(): ScheduleDay[] {
	return [
		{
			date: "2026-07-05",
			label: "5 Jul",
			availablePeopleCount: 3,
			tasks: [
				{
					id: "card-1",
					title: "Pack kitchen",
					priority: "high",
					area: { id: "room-1", name: "Kitchen" },
					assignedPeople: [],
					peopleNeeded: 2,
					assignedCount: 0,
					staffingStatus: "underStaffed",
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
			tasks: [
				{
					id: "card-low",
					title: "Sweep floor",
					priority: "low",
					area: { id: "room-1", name: "Kitchen" },
					assignedPeople: [],
					peopleNeeded: 1,
					assignedCount: 0,
					staffingStatus: "underStaffed",
					completed: false,
					scheduledDate: "2026-07-06",
					taskId: null,
				},
			],
		},
	];
}

function setup(days: ScheduleDay[]) {
	const source = ref(days);
	const scope = effectScope();
	const api = scope.run(() =>
		useScheduleBoardDnd({
			source: computed(() => source.value),
			queryKey: ["daily-schedule"],
		}),
	)!;
	return { api, source };
}

beforeEach(() => {
	mutateAsync = vi.fn();
	invalidateQueries.mockClear();
});

describe("useScheduleBoardDnd", () => {
	it("optimistically assigns a person, sends the full body, reconciles, and invalidates", async () => {
		mutateAsync.mockResolvedValue({
			id: "card-1",
			assignedPeople: [{ id: "p1", name: "Alex Kim", initials: "AK" }],
			assignedCount: 1,
			peopleNeeded: 2,
			staffingStatus: "underStaffed",
			completed: false,
			priority: "high",
			area: { id: "room-1", name: "Kitchen" },
			title: "Pack kitchen",
			taskId: "task-1",
		});
		const { api } = setup(makeDays());

		await api.assignPerson("card-1", { id: "p1", name: "Alex Kim" });

		// Body is the full, valid taskId-linked update with the new assignment.
		expect(mutateAsync).toHaveBeenCalledTimes(1);
		const body = mutateAsync.mock.calls[0][0].body;
		expect(body).toMatchObject({
			taskId: "task-1",
			scheduledDate: "2026-07-05",
			assignedTo: ["p1"],
		});
		// Reconciled to the server response.
		expect(api.board.value[0].tasks[0].assignedCount).toBe(1);
		expect(api.board.value[0].tasks[0].assignedPeople[0].id).toBe("p1");
		expect(invalidateQueries).toHaveBeenCalledTimes(1);
		expect(api.error.value).toBeNull();
	});

	it("treats a duplicate assignment as a no-op (no duplicate, no mutation)", async () => {
		const days = makeDays();
		days[0].tasks[0].assignedPeople = [
			{ id: "p1", name: "Alex Kim", initials: "AK" },
		];
		days[0].tasks[0].assignedCount = 1;
		const { api } = setup(days);

		await api.assignPerson("card-1", { id: "p1", name: "Alex Kim" });

		expect(mutateAsync).not.toHaveBeenCalled();
		expect(api.board.value[0].tasks[0].assignedPeople).toHaveLength(1);
	});

	it("rolls back to the untouched data and surfaces an error when assignment fails", async () => {
		mutateAsync.mockRejectedValue(
			new Error("assignedTo count must not exceed peopleNeeded"),
		);
		const { api } = setup(makeDays());

		await api.assignPerson("card-1", { id: "p1", name: "Alex Kim" });

		// Reverted: no optimistic person remains.
		expect(api.board.value[0].tasks[0].assignedPeople).toHaveLength(0);
		expect(api.board.value[0].tasks[0].assignedCount).toBe(0);
		expect(api.error.value).toContain("exceed peopleNeeded");
	});

	it("optimistically reschedules a card to another day and persists the new date", async () => {
		mutateAsync.mockResolvedValue({
			id: "card-1",
			assignedPeople: [],
			assignedCount: 0,
			peopleNeeded: 2,
			staffingStatus: "underStaffed",
			completed: false,
			priority: "high",
			area: { id: "room-1", name: "Kitchen" },
			title: "Pack kitchen",
			taskId: "task-1",
		});
		const { api } = setup(makeDays());

		await api.rescheduleCard("card-1", "2026-07-06");

		expect(api.board.value[0].tasks).toHaveLength(0); // left source day
		expect(api.board.value[1].tasks.map((task) => task.id)).toEqual([
			"card-1",
			"card-low",
		]);
		expect(mutateAsync.mock.calls[0][0].body.scheduledDate).toBe("2026-07-06");
	});

	it("reverts a failed reschedule to the original column", async () => {
		mutateAsync.mockRejectedValue(
			new Error("scheduledDate must be within the planning window"),
		);
		const { api } = setup(makeDays());

		await api.rescheduleCard("card-1", "2026-07-06");

		expect(api.board.value[0].tasks[0].id).toBe("card-1"); // back in source day
		expect(api.board.value[1].tasks.map((task) => task.id)).toEqual([
			"card-low",
		]);
		expect(api.error.value).toContain("planning window");
	});
});
