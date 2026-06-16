/**
 * Planning Window — pure format utilities and the `PlanWindowDay` interface.
 *
 * The canonical source of truth for the move timeline is now
 * `GET /api/planning-window` (backend). Consumers should use
 * `usePlanningWindow()` from `@/shared/composables/usePlanningWindow`
 * to fetch, derive, and reactively share the planning-window data.
 *
 * ## Consumer pattern
 * ```
 * import { usePlanningWindow } from "@/shared/composables/usePlanningWindow";
 * const { planWindowDays, planWindowDayCount, isLoading, isError } = usePlanningWindow();
 * ```
 *
 * Pure utilities below (formatPlanDayLabel, PlanWindowDay interface)
 * remain available for use alongside the composable.
 *
 * ## Date handling
 * `formatPlanDayLabel` uses UTC to ensure deterministic labels regardless
 * of runtime timezone. Dates derived from the backend (`startDate` / `endDate`
 * ISO 8601 strings) produce a stable `dateString` via `toISOString().slice(0, 10)`.
 */

/**
 * Format a Date as a compact human-readable label.
 *
 * Uses `Intl.DateTimeFormat` with `en-US` locale and explicit `UTC`
 * timezone. Produces labels like `"Sun 5 Jul"` (abbreviated weekday,
 * numeric day, abbreviated month) deterministically regardless of
 * runtime timezone.
 */
export function formatPlanDayLabel(date: Date): string {
	return new Intl.DateTimeFormat("en-US", {
		weekday: "short",
		day: "numeric",
		month: "short",
		timeZone: "UTC",
	}).format(date);
}

/**
 * Format a planning window range as a compact human-readable string.
 *
 * Accepts ISO 8601 date strings from the raw API response and an inclusive
 * day count. Uses `Intl.DateTimeFormat` with `en-US` locale and explicit
 * `UTC` timezone for deterministic labels regardless of runtime timezone,
 * mirroring the existing `formatPlanDayLabel` pattern.
 *
 * Produces format like: `"5 Jul – 13 Aug 2026 · 40 days"`
 */
export function formatPlanWindowRange(
	startDate: string,
	endDate: string,
	days: number,
): string {
	const monthFmt = new Intl.DateTimeFormat("en-US", {
		month: "short",
		timeZone: "UTC",
	});

	const start = new Date(startDate);
	const end = new Date(endDate);

	const startMonth = monthFmt.format(start);
	const startDay = start.getUTCDate();
	const endMonth = monthFmt.format(end);
	const endDay = end.getUTCDate();
	const year = start.getUTCFullYear();

	return `${startDay} ${startMonth} – ${endDay} ${endMonth} ${year} · ${days} days`;
}

export interface PlanWindowDay {
	date: Date;
	label: string;
	dateString: string;
}
