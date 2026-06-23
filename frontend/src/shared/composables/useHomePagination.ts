import { computed, ref } from "vue";
import { usePlanningWindow } from "@/shared/composables/usePlanningWindow";

/**
 * Format a date range label from two ISO date strings.
 *
 * Produces format like "2 Jul – 5 Jul, 2024" — day, abbreviated month, year.
 * Uses Intl.DateTimeFormat with en-US locale and UTC timezone for deterministic labels.
 */
function formatPageRangeLabel(startISO: string, endISO: string): string {
	const dayMonthFmt = new Intl.DateTimeFormat("en-US", {
		day: "numeric",
		month: "short",
		timeZone: "UTC",
	});

	const yearFmt = new Intl.DateTimeFormat("en-US", {
		year: "numeric",
		timeZone: "UTC",
	});

	const start = new Date(startISO);
	const end = new Date(endISO);

	const startPart = dayMonthFmt.format(start);
	const endDayMonth = dayMonthFmt.format(end);
	const endYear = yearFmt.format(end);

	return `${startPart} – ${endDayMonth}, ${endYear}`;
}

/** Compact label for narrow viewports, e.g. "5–8 Jul" or "30 Jul–2 Aug". */
function formatCompactRangeLabel(startISO: string, endISO: string): string {
	const monthFmt = new Intl.DateTimeFormat("en-US", { month: "short", timeZone: "UTC" });
	const start = new Date(startISO);
	const end = new Date(endISO);
	const startDay = start.getUTCDate();
	const endDay = end.getUTCDate();
	const startMonth = monthFmt.format(start);
	const endMonth = monthFmt.format(end);
	if (startMonth === endMonth) {
		return startDay === endDay ? `${startDay} ${endMonth}` : `${startDay}–${endDay} ${endMonth}`;
	}
	return `${startDay} ${startMonth}–${endDay} ${endMonth}`;
}

/** Module-scoped shared state so all callers share the same page/daysPerPage refs. */
const page = ref<number>(1);
const daysPerPage = ref<number>(4);

export function useHomePagination() {
	const planningWindow = usePlanningWindow();

	const totalPages = computed<number>(() => {
		const count = planningWindow.planWindowDays.value.length;
		if (count === 0) return 1;
		return Math.max(1, Math.ceil(count / daysPerPage.value));
	});

	/** Human-readable date range label for the current page, e.g. "2 Jul – 5 Jul, 2024". */
	const rangeLabel = computed<string>(() => {
		const days = planningWindow.planWindowDays.value;
		if (days.length === 0) return "—";
		const startIdx = (page.value - 1) * daysPerPage.value;
		const endIdx = Math.min(startIdx + daysPerPage.value - 1, days.length - 1);
		const startDate = days[startIdx]?.dateString ?? "";
		const endDate = days[endIdx]?.dateString ?? "";
		if (!startDate || !endDate) return "—";
		return formatPageRangeLabel(startDate, endDate);
	});

	/** Compact label for narrow viewports, e.g. "5–8 Jul". */
	const compactRangeLabel = computed<string>(() => {
		const days = planningWindow.planWindowDays.value;
		if (days.length === 0) return "—";
		const startIdx = (page.value - 1) * daysPerPage.value;
		const endIdx = Math.min(startIdx + daysPerPage.value - 1, days.length - 1);
		const startDate = days[startIdx]?.dateString ?? "";
		const endDate = days[endIdx]?.dateString ?? "";
		if (!startDate || !endDate) return "—";
		return formatCompactRangeLabel(startDate, endDate);
	});

	const isLoading = computed<boolean>(() => planningWindow.isLoading.value);
	const isError = computed<boolean>(() => planningWindow.isError.value);

	function goPrev(): void {
		if (page.value > 1) {
			page.value--;
		}
	}

	function goNext(): void {
		if (page.value < totalPages.value) {
			page.value++;
		}
	}

	function goToday(): void {
		const days = planningWindow.planWindowDays.value;
		if (days.length === 0) return;

		// Get today's date as a UTC date string (YYYY-MM-DD).
		const now = new Date();
		const todayISO =
			now.getUTCFullYear() +
			"-" +
			String(now.getUTCMonth() + 1).padStart(2, "0") +
			"-" +
			String(now.getUTCDate()).padStart(2, "0");

		// Find today's index in the planning window.
		const todayIdx = days.findIndex((d) => d.dateString === todayISO);

		if (todayIdx >= 0) {
			// Compute the 1-indexed page containing today.
			page.value = Math.floor(todayIdx / daysPerPage.value) + 1;
		} else {
			// Today is outside the planning window — fall back to page 1.
			page.value = 1;
		}
	}

	return {
		page,
		daysPerPage,
		totalPages,
		rangeLabel,
		compactRangeLabel,
		isLoading,
		isError,
		goPrev,
		goNext,
		goToday,
	};
}
