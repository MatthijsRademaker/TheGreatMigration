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

export function useHomePagination() {
	const planningWindow = usePlanningWindow();

	const page = ref<number>(1);
	const daysPerPage = ref<number>(4);

	const totalPages = computed<number>(() => {
		const count = planningWindow.planWindowDays.value.length;
		if (count === 0) return 1;
		return Math.max(1, Math.ceil(count / daysPerPage.value));
	});

	/** ISO date string for the first day of the current page window. */
	const start = computed<string>(() => {
		const days = planningWindow.planWindowDays.value;
		if (days.length === 0) return "";
		const idx = (page.value - 1) * daysPerPage.value;
		const day = days[idx] ?? days[0];
		return day.dateString;
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
		start,
		rangeLabel,
		isLoading,
		isError,
		goPrev,
		goNext,
		goToday,
	};
}
