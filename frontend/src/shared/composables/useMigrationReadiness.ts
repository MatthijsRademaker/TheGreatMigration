import { computed, type ComputedRef } from "vue";
import { useTools } from "@/tools/composables/useTools";
import type { ScheduleDay } from "@/calendar/composables/useDailySchedule";

interface Dimension {
	/** 0..1 coverage for this dimension. */
	ratio: number;
	/** Whether the dimension has a denominator (otherwise excluded from the aggregate). */
	applicable: boolean;
}

/**
 * Derived, read-only readiness for the Move-Day Readiness Journey. Aggregates
 * purely from existing query data — per-task staffing completeness from the
 * passed schedule days and tool coverage from the tools query — with no new
 * persisted state. As those queries update (including after a drag-and-drop
 * mutation invalidates the schedule), the readiness recomputes reactively.
 *
 * Rooms coverage has no backend contract yet (the "Rooms completed" KPI is a
 * placeholder), so it is intentionally excluded until that data exists.
 */
export function useMigrationReadiness(days: ComputedRef<ScheduleDay[]>) {
	const { data: tools } = useTools();

	/** Staffing coverage: assigned helpers vs helpers needed across all cards. */
	const staffing = computed<Dimension>(() => {
		let needed = 0;
		let assigned = 0;
		for (const day of days.value) {
			for (const task of day.tasks) {
				needed += task.peopleNeeded;
				assigned += Math.min(task.assignedCount, task.peopleNeeded);
			}
		}
		return needed > 0
			? { ratio: assigned / needed, applicable: true }
			: { ratio: 0, applicable: false };
	});

	/** Tool coverage: claimed tools vs total tools. */
	const toolCoverage = computed<Dimension>(() => {
		const total = tools.value.summary.total;
		return total > 0
			? { ratio: tools.value.summary.claimed / total, applicable: true }
			: { ratio: 0, applicable: false };
	});

	/** Aggregate readiness (0..1): mean of the applicable dimensions. */
	const readiness = computed<number>(() => {
		const dims = [staffing.value, toolCoverage.value].filter((d) => d.applicable);
		if (dims.length === 0) return 0;
		return dims.reduce((sum, d) => sum + d.ratio, 0) / dims.length;
	});

	/** Readiness as a whole-number percent (0..100). */
	const percent = computed<number>(() => Math.round(readiness.value * 100));

	/** True only when there is real data and every applicable dimension is fully covered. */
	const isComplete = computed<boolean>(
		() =>
			(staffing.value.applicable || toolCoverage.value.applicable) &&
			readiness.value >= 1,
	);

	return { readiness, percent, isComplete, staffing, toolCoverage };
}
