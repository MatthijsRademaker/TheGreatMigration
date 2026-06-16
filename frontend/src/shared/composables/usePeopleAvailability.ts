import { computed } from "vue";
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
				availableToday: 0,
				totalPeople: 0,
				legend: [],
			},
			daysISO: [],
		};
	}

	const peopleRaw = data.people ?? [];
	const legendRaw = data.statuses ?? [];

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

	const legend = legendRaw
		.filter((s) => CANONICAL_STATUSES.has(s.id))
		.map((s) => ({
			id: s.id as PersonAvailabilityEntry["status"],
			label: s.label,
		}));

	return {
		props: {
			title: "People availability",
			description: "Track who is available and where each person can help.",
			days,
			people,
			legend,
			availableToday: data.summary?.availableToday ?? 0,
			totalPeople: data.summary?.totalPeople ?? 0,
		},
		daysISO,
	};
}

interface UsePeopleAvailabilityOptions {
	/** Explicit start date (YYYY-MM-DD). When omitted, the planning window start date is used. */
	start?: string;
}

export function usePeopleAvailability(options?: UsePeopleAvailabilityOptions) {
	// Resolve the start date: explicit or from the planning window.
	const planningWindow = options?.start ? null : usePlanningWindow();

	const startParam = computed<string | undefined>(() => {
		if (options?.start) return options.start;
		if (planningWindow?.planWindowDays.value.length) {
			return planningWindow.planWindowDays.value[0].dateString;
		}
		return undefined;
	});

	// Defer the dashboard query until the planning window resolves (if needed).
	const queryEnabled = computed<boolean>(() => {
		if (options?.start) return true;
		return !planningWindow?.isLoading.value && startParam.value != null;
	});

	const query = useQuery({
		...getDashboardPeopleAvailabilityQuery(
			startParam.value ? { query: { start: startParam.value } } : undefined,
		),
		enabled: queryEnabled,
	});

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
		() => !isLoading.value && !isError.value && data.value.totalPeople === 0,
	);

	return {
		data,
		daysISO,
		rawData,
		isLoading,
		isError,
		isEmpty,
		/** Refresh the underlying query. */
		refresh: () => query.refetch(),
		queryKey: getDashboardPeopleAvailabilityQueryKey(),
	};
}
