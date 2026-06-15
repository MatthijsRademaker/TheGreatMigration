import { describe, expect, it } from "vitest";
import {
	PLAN_WINDOW_START,
	PLAN_WINDOW_END,
	formatPlanDayLabel,
	planWindowDayCount,
	planWindowDays,
} from "../src/shared/lib/planWindow";

describe("planWindow module", () => {
	it("exports the default ISO date constants", () => {
		expect(PLAN_WINDOW_START).toBe("2026-07-05");
		expect(PLAN_WINDOW_END).toBe("2026-08-13");
	});

	it("planWindowDayCount equals 40 for the default inclusive range", () => {
		// 27 July days + 13 August days, both inclusive
		expect(planWindowDayCount).toBe(40);
	});

	it("planWindowDays length equals planWindowDayCount", () => {
		expect(planWindowDays.length).toBe(planWindowDayCount);
	});

	it("planWindowDays first and last entries match the constants", () => {
		expect(planWindowDays[0].dateString).toBe(PLAN_WINDOW_START);
		expect(planWindowDays[planWindowDays.length - 1].dateString).toBe(
			PLAN_WINDOW_END,
		);
	});

	it("formatPlanDayLabel for 2026-07-05 contains day and month", () => {
		const label = formatPlanDayLabel(new Date("2026-07-05T12:00:00Z"));
		expect(label).toContain("5");
		expect(label).toContain("Jul");
	});

	it("formatPlanDayLabel for 2026-08-13 contains day and month", () => {
		const label = formatPlanDayLabel(new Date("2026-08-13T12:00:00Z"));
		expect(label).toContain("13");
		expect(label).toContain("Aug");
	});
});
