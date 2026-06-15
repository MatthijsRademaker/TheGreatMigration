/**
 * Planning Window — canonical source of truth for the move-plan date range.
 *
 * All date-dependent views (CalendarView, HomeView, and future consumers)
 * derive their rendered content from this module's constants and pure helpers.
 *
 * ## Inclusive contract
 * Both PLAN_WINDOW_START and PLAN_WINDOW_END are included in every derived
 * value (planWindowDays, planWindowDayCount). Changing either constant
 * propagates automatically to all consumers.
 *
 * ## Default range
 * 2026-07-05 through 2026-08-13 (40 inclusive days; 27 July + 13 August).
 *
 * ## Consumer pattern
 * Import named exports directly:
 * ```
 * import { planWindowDays } from "@/shared/lib/planWindow";
 * // Iterate planWindowDays in templates, derive counts, format labels, etc.
 * ```
 * For future reactive plan-window state, a Vue composable can wrap this module
 * without changing consumer code.
 *
 * ## Evaluation strategy
 * `planWindowDays` and `planWindowDayCount` are evaluated eagerly at module
 * import time. This keeps consumer code simple (no optional chaining, no
 * async), but means changing the window requires a rebuild. If runtime
 * reconfiguration becomes necessary, the eagerly-evaluated constants can be
 * replaced with lazy getters without affecting consumer code.
 *
 * ## Future PeopleView integration
 * PeopleView currently uses weekday availability strings ("Mon, Wed, Fri", …)
 * that represent availability *patterns*, not specific dates. When real
 * availability CRUD is built (CRUD helpers with per-date availability), the
 * PeopleView will consume planWindowDays to bind helpers to actual calendar
 * dates. The exports below are designed with that future consumer in mind.
 */

export const PLAN_WINDOW_START = "2026-07-05";
export const PLAN_WINDOW_END = "2026-08-13";

/**
 * Format a Date as a compact human-readable label.
 *
 * Uses `Intl.DateTimeFormat` with `en-US` locale, producing labels like
 * `"Sun 5 Jul"` (abbreviated weekday, numeric day, abbreviated month).
 */
export function formatPlanDayLabel(date: Date): string {
	return new Intl.DateTimeFormat("en-US", {
		weekday: "short",
		day: "numeric",
		month: "short",
	}).format(date);
}

export interface PlanWindowDay {
	date: Date;
	label: string;
	dateString: string;
}

function generatePlanWindowDays(): PlanWindowDay[] {
	// Date-only ISO strings parse as UTC midnight per ECMAScript spec,
	// ensuring deterministic dateString values regardless of runtime timezone.
	const start = new Date(PLAN_WINDOW_START);
	const end = new Date(PLAN_WINDOW_END);

	if (start >= end) {
		throw new Error(
			`planWindow: PLAN_WINDOW_START (${PLAN_WINDOW_START}) must precede PLAN_WINDOW_END (${PLAN_WINDOW_END})`,
		);
	}
	const days: PlanWindowDay[] = [];

	const cursor = new Date(start);
	while (cursor <= end) {
		days.push({
			date: new Date(cursor),
			label: formatPlanDayLabel(cursor),
			dateString: cursor.toISOString().slice(0, 10),
		});
		cursor.setDate(cursor.getDate() + 1);
	}

	return days;
}

/** Array of inclusive day objects between PLAN_WINDOW_START and PLAN_WINDOW_END. */
export const planWindowDays: PlanWindowDay[] = generatePlanWindowDays();

/** Inclusive count of days between PLAN_WINDOW_START and PLAN_WINDOW_END. */
export const planWindowDayCount: number = planWindowDays.length;
