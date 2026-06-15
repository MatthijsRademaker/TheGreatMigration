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

export interface PlanWindowDay {
	date: Date;
	label: string;
	dateString: string;
}
