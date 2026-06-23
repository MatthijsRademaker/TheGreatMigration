import { computed, ref, watch, type Ref } from "vue";
import { useQuery } from "@pinia/colada";
import {
	getDashboardPeopleAvailabilityQuery,
	getDashboardPeopleAvailabilityQueryKey,
} from "@/client/@pinia/colada.gen";
import { formatPlanDayLabel } from "@/shared/lib/planWindow";
import { usePlanningWindow } from "@/shared/composables/usePlanningWindow";
import type { DashboardBody } from "@/client/types.gen";
import type {
	PeopleAvailabilityProps,
	PersonAvailability,
	PersonAvailabilityEntry,
} from "@/people/types";

const CANONICAL_STATUSES = new Set(["available", "busy", "partial", "off"]);

function isoDateToDayLabel(isoDate: string): string {
	const date = new Date(isoDate);
	return formatPlanDayLabel(date);
}

interface AdaptedResult {
	props: PeopleAvailabilityProps;
	daysISO: string[];
}

/**
 * Adapt a generated DashboardBody into the component-local PeopleAvailabilityProps.
 *
 * Handles three gaps between the API contract and the component contract:
 * 1. Nullable API arrays → empty arrays
 * 2. ISO date strings → human-readable day labels
 * 3. Loose `string` status → narrowed AvailabilityStatus union (canonical filtering)
 *
 * Also produces a `daysISO` array (parallel to `days`) of ISO 8601 date strings.
 */
function adaptToComponentProps(data: DashboardBody | undefined): AdaptedResult {
	if (!data) {
		return {
			props: {
				people: [],
				days: [],
			},
			daysISO: [],
		};
	}

	const peopleRaw = data.people ?? [];
	const days: string[] = [];
	const daysISO: string[] = [];
	const people: PersonAvailability[] = [];

	// Build day labels and ISO dates from the range.
	if (data.range) {
		const startDate = new Date(data.range.startDate);
		const dayCount = data.range.days;
		for (let i = 0; i < dayCount; i++) {
			const cursor = new Date(startDate);
			cursor.setDate(cursor.getDate() + i);
			const isoDate = cursor.toISOString().slice(0, 10);
			days.push(isoDateToDayLabel(isoDate));
			daysISO.push(isoDate);
		}
	}

	for (const rawPerson of peopleRaw) {
		const availRaw = rawPerson.availability ?? [];
		const availability: PersonAvailabilityEntry[] = [];

		for (const entry of availRaw) {
			if (!CANONICAL_STATUSES.has(entry.status)) {
				continue;
			}
			const dayLabel = isoDateToDayLabel(entry.date);
			availability.push({
				date: dayLabel,
				status: entry.status as PersonAvailabilityEntry["status"],
			});
		}

		people.push({
			id: rawPerson.id,
			name: rawPerson.name,
			availability,
		});
	}

	return {
		props: {
			title: "People availability",
			days,
			people,
		},
		daysISO,
	};
}

interface UsePeopleAvailabilityOptions {
	/** Explicit start date (YYYY-MM-DD). When omitted, the planning window start date is used. */
	start?: string;
	/** 1-indexed page for day pagination. Default 1. */
	page?: number;
	/** Number of days per page. Default 7. */
	daysPerPage?: number;
	/** Number of people to skip for person-level pagination. Default 0. */
	offset?: number;
	/** Maximum people to return. 0 means no limit. Default 0. */
	limit?: number;
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

export function usePeopleAvailability(options?: UsePeopleAvailabilityOptions) {
	const page = options?.pageRef ?? ref(options?.page ?? 1);
	const daysPerPage = options?.daysPerPageRef ?? ref(options?.daysPerPage ?? 7);
	const offset = ref(options?.offset ?? 0);
	const limit = ref(options?.limit ?? 0);

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

	// Defer the dashboard query until the planning window resolves (if needed).
	const queryEnabled = computed<boolean>(() => {
		if (options?.start) return true;
		return !planningWindow?.isLoading.value && startParam.value != null;
	});

	const query = useQuery(() => ({
		...getDashboardPeopleAvailabilityQuery(
			startParam.value
				? {
						query: {
							start: startParam.value,
							days: daysPerPage.value,
							...(offset.value > 0 ? { offset: offset.value } : {}),
							...(limit.value > 0 ? { limit: limit.value } : {}),
						},
					}
				: undefined,
		),
		enabled: queryEnabled.value,
	}));

	const adapted = computed<AdaptedResult>(() =>
		adaptToComponentProps(query.data.value),
	);

	/** The adapted PeopleAvailabilityProps for the component. */
	const data = computed<PeopleAvailabilityProps>(() => adapted.value.props);

	/** ISO date strings parallel to `data.value.days`. */
	const daysISO = computed<string[]>(() => adapted.value.daysISO);

	/** Raw DashboardBody for write operations that need ISO dates (legacy, prefer daysISO). */
	const rawData = computed<DashboardBody | undefined>(() => query.data.value);

	const isLoading = computed<boolean>(() => query.isPending.value);

	const isError = computed<boolean>(() => query.error.value != null);

	const isEmpty = computed<boolean>(
		() =>
			!isLoading.value &&
			!isError.value &&
			(data.value.people ?? []).length === 0,
	);

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
		daysISO,
		rawData,
		isLoading,
		isError,
		isEmpty,
		page,
		totalPages,
		daysPerPage,
		totalDays,
		goToPrevPage,
		goToNextPage,
		/** Refresh the underlying query. */
		refresh: () => query.refetch(),
		queryKey: getDashboardPeopleAvailabilityQueryKey(),
	};
}
