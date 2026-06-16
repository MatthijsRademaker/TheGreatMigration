import { computed } from "vue";
import { useQuery } from "@pinia/colada";
import {
	getDashboardDailyScheduleQuery,
	getDashboardDailyScheduleQueryKey,
} from "@/client/@pinia/colada.gen";
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
}

export interface ScheduleDay {
	date: string;
	label: string;
	availablePeopleCount: number;
	tasks: ScheduleTaskCard[];
}

const CANONICAL_PRIORITIES = new Set(["high", "medium", "low"]);
const CANONICAL_STAFFING = new Set(["fullyStaffed", "underStaffed"]);

function isCanonicalTask(card: ApiTaskCard): boolean {
	return (
		CANONICAL_PRIORITIES.has(card.priority) &&
		CANONICAL_STAFFING.has(card.staffingStatus)
	);
}

function adaptTaskCard(card: ApiTaskCard): ScheduleTaskCard {
	return {
		id: card.id,
		title: card.title,
		priority: card.priority as ScheduleTaskCard["priority"],
		roomArea: card.roomArea,
		assignedPeople: card.assignedPeople ?? [],
		peopleNeeded: card.peopleNeeded,
		assignedCount: card.assignedCount,
		staffingStatus: card.staffingStatus as ScheduleTaskCard["staffingStatus"],
	};
}

function adaptDay(day: ApiScheduleDay): ScheduleDay {
	const rawTasks = day.tasks ?? [];
	const tasks: ScheduleTaskCard[] = [];
	for (const card of rawTasks) {
		if (isCanonicalTask(card)) {
			tasks.push(adaptTaskCard(card));
		}
	}
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

export function useDailySchedule() {
	const query = useQuery(getDashboardDailyScheduleQuery());

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

	return {
		data,
		isLoading,
		isError,
		isEmpty,
		/** Refresh the underlying daily-schedule query. */
		refresh: () => query.refetch(),
		queryKey: getDashboardDailyScheduleQueryKey(),
	};
}
