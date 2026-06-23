import { computed, ref, watch, type Ref } from "vue";
import { useQuery } from "@pinia/colada";
import {
	getDashboardDailyScheduleQuery,
	getDashboardDailyScheduleQueryKey,
} from "@/client/@pinia/colada.gen";
import { usePlanningWindow } from "@/shared/composables/usePlanningWindow";
import type {
	DailyScheduleBody,
	ScheduleDay as ApiScheduleDay,
	TaskCard as ApiTaskCard,
	AssignedPerson,
} from "@/client/types.gen";

// ── Local types (narrowed from generated) ──────────────────────────────────

export interface ScheduleTaskCard {
	id: string;
	title: string;
	priority: "high" | "medium" | "low";
	roomArea: string;
	assignedPeople: AssignedPerson[];
	peopleNeeded: number;
	assignedCount: number;
	staffingStatus: "fullyStaffed" | "underStaffed";
	completed: boolean;
	scheduledDate: string;
	taskId: string | null;
}

export interface ScheduleDay {
	date: string;
	label: string;
	availablePeopleCount: number;
	tasks: ScheduleTaskCard[];
}

const CANONICAL_PRIORITIES = new Set(["high", "medium", "low"]);
const CANONICAL_STAFFING = new Set(["fullyStaffed", "underStaffed"]);
const PRIORITY_RANK: Record<ScheduleTaskCard["priority"], number> = {
	high: 0,
	medium: 1,
	low: 2,
};

export function sortScheduleTasksByPriority(
	tasks: ScheduleTaskCard[],
): ScheduleTaskCard[] {
	return tasks.sort(
		(left, right) =>
			PRIORITY_RANK[left.priority] - PRIORITY_RANK[right.priority],
	);
}

function isCanonicalTask(card: ApiTaskCard): boolean {
	return (
		CANONICAL_PRIORITIES.has(card.priority) &&
		CANONICAL_STAFFING.has(card.staffingStatus)
	);
}

function adaptTaskCard(
	card: ApiTaskCard,
	scheduledDate: string,
): ScheduleTaskCard {
	return {
		id: card.id,
		title: card.title,
		priority: card.priority as ScheduleTaskCard["priority"],
		roomArea: card.roomArea,
		assignedPeople: card.assignedPeople ?? [],
		peopleNeeded: card.peopleNeeded,
		assignedCount: card.assignedCount,
		staffingStatus: card.staffingStatus as ScheduleTaskCard["staffingStatus"],
		completed: card.completed ?? false,
		scheduledDate,
		taskId: card.taskId ?? null,
	};
}

function adaptDay(day: ApiScheduleDay): ScheduleDay {
	const rawTasks = day.tasks ?? [];
	const tasks: ScheduleTaskCard[] = [];
	for (const card of rawTasks) {
		if (isCanonicalTask(card)) {
			tasks.push(adaptTaskCard(card, day.date));
		}
	}
	sortScheduleTasksByPriority(tasks);

	return {
		date: day.date,
		label: day.label,
		availablePeopleCount: day.availablePeopleCount,
		tasks,
	};
}

function adaptToComponentDays(
	data: DailyScheduleBody | undefined,
): ScheduleDay[] {
	if (!data || !data.days) {
		return [];
	}
	return data.days.map(adaptDay);
}

export interface DailyScheduleState {
	days: ScheduleDay[];
	range: {
		startDate: string;
		endDate: string;
		days: number;
	} | null;
}

export interface UseDailyScheduleOptions {
	/** Explicit start date (YYYY-MM-DD). When omitted, the planning window start date is used. */
	start?: string;
	/** 1-indexed page for day pagination. Default 1. */
	page?: number;
	/** Number of days per page. Default 4. */
	daysPerPage?: number;
	/**
	 * Reactive ref for page, owned externally. When provided, the composable uses
	 * this ref directly instead of creating its own internal ref.
	 */
	pageRef?: Ref<number>;
	/**
	 * Reactive ref for daysPerPage, owned externally. When provided, the composable uses
	 * this ref directly instead of creating its own internal ref.
	 */
	daysPerPageRef?: Ref<number>;
}

export function useDailySchedule(options?: UseDailyScheduleOptions) {
	const page = options?.pageRef ?? ref(options?.page ?? 1);
	const daysPerPage = options?.daysPerPageRef ?? ref(options?.daysPerPage ?? 4);

	// Resolve the planning window.
	const planningWindow = options?.start ? null : usePlanningWindow();

	// Compute total days from planning window.
	const totalDays = computed<number>(() => {
		if (planningWindow?.planWindowDays.value.length) {
			return planningWindow.planWindowDays.value.length;
		}
		// Fall back to daysPerPage when planning window is unavailable.
		return daysPerPage.value;
	});

	// Compute total pages.
	const totalPages = computed<number>(() => {
		return Math.max(1, Math.ceil(totalDays.value / daysPerPage.value));
	});

	// Watch for planning window changes and reset page to 1 when the window actually changes
	// (not on initial load, not on loading→loaded transition).
	// Skip when page is externally owned (pageRef provided) — the external owner handles resets.
	if (planningWindow && !options?.pageRef) {
		watch(
			() => planningWindow.planWindowDays.value.map((d) => d.dateString),
			(_newDates, oldDates) => {
				if (oldDates && oldDates.length > 0 && page.value > 1) {
					page.value = 1;
				}
			},
		);
	}

	// Compute the effective start date: planning window start + (page-1) * daysPerPage.
	const startParam = computed<string | undefined>(() => {
		if (options?.start) return options.start;
		if (planningWindow?.planWindowDays.value.length) {
			const baseDate = planningWindow.planWindowDays.value[0].dateString;
			if (page.value > 1) {
				const d = new Date(baseDate);
				d.setDate(d.getDate() + (page.value - 1) * daysPerPage.value);
				return d.toISOString().slice(0, 10);
			}
			return baseDate;
		}
		return undefined;
	});

	// Defer the dashboard query until the planning window resolves (when no explicit start).
	const queryEnabled = computed<boolean>(() => {
		if (options?.start) return true;
		return !planningWindow?.isLoading.value && startParam.value != null;
	});

	const query = useQuery(() => ({
		...getDashboardDailyScheduleQuery(
			startParam.value
				? {
						query: {
							start: startParam.value,
							days: daysPerPage.value,
						},
					}
				: undefined,
		),
		enabled: queryEnabled.value,
	}));

	const data = computed<DailyScheduleState>(() => {
		const raw = query.data.value;
		return {
			days: adaptToComponentDays(raw),
			range: raw?.range
				? {
						startDate: raw.range.startDate,
						endDate: raw.range.endDate,
						days: raw.range.days,
					}
				: null,
		};
	});

	const isLoading = computed<boolean>(() => query.isPending.value);

	const isError = computed<boolean>(() => query.error.value != null);

	const isEmpty = computed<boolean>(
		() => !isLoading.value && !isError.value && data.value.days.length === 0,
	);

	/** Human-readable label for the currently visible date range, e.g. "1 Aug (Sat) – 4 Aug (Tue)". */
	const dateRangeLabel = computed<string>(() => {
		const days = data.value.days;
		if (!days || days.length === 0) return "—";
		const first = days[0]?.label ?? "—";
		const last = days[days.length - 1]?.label ?? "—";
		return `${first} – ${last}`;
	});

	/** Navigate to the previous page. */
	function goToPrevPage() {
		if (page.value > 1) {
			page.value--;
		}
	}

	/** Navigate to the next page. */
	function goToNextPage() {
		if (page.value < totalPages.value) {
			page.value++;
		}
	}

	return {
		data,
		isLoading,
		isError,
		isEmpty,
		page,
		totalPages,
		daysPerPage,
		totalDays,
		goToPrevPage,
		goToNextPage,
		/** Human-readable label for the currently visible date range. */
		dateRangeLabel,
		/** Refresh the underlying daily-schedule query. */
		refresh: () => query.refetch(),
		queryKey: getDashboardDailyScheduleQueryKey(),
	};
}
