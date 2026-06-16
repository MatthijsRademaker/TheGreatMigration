import { computed } from "vue";
import { useQuery } from "@pinia/colada";
import {
	getDashboardPeopleAvailabilityQuery,
	getDashboardPeopleAvailabilityQueryKey,
} from "@/client/@pinia/colada.gen";
import { formatPlanDayLabel } from "@/shared/lib/planWindow";
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

/**
 * Adapt a generated DashboardBody into the component-local PeopleAvailabilityProps.
 *
 * Handles three gaps between the API contract and the component contract:
 * 1. Nullable API arrays → empty arrays
 * 2. ISO date strings → human-readable day labels
 * 3. Loose `string` status → narrowed AvailabilityStatus union (canonical filtering)
 */
function adaptToComponentProps(
	data: DashboardBody | undefined,
): PeopleAvailabilityProps {
	if (!data) {
		return {
			people: [],
			days: [],
			availableToday: 0,
			totalPeople: 0,
			legend: [],
		};
	}

	const peopleRaw = data.people ?? [];
	const legendRaw = data.statuses ?? [];

	const days: string[] = [];
	const people: PersonAvailability[] = [];

	// Build day labels from the range.
	if (data.range) {
		const startDate = new Date(data.range.startDate);
		const dayCount = data.range.days;
		for (let i = 0; i < dayCount; i++) {
			const cursor = new Date(startDate);
			cursor.setDate(cursor.getDate() + i);
			days.push(isoDateToDayLabel(cursor.toISOString().slice(0, 10)));
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
		title: "People availability",
		description: "Track who is available and where each person can help.",
		days,
		people,
		legend,
		availableToday: data.summary?.availableToday ?? 0,
		totalPeople: data.summary?.totalPeople ?? 0,
	};
}

export function usePeopleAvailability() {
	const query = useQuery(getDashboardPeopleAvailabilityQuery());

	const adapted = computed<PeopleAvailabilityProps>(() =>
		adaptToComponentProps(query.data.value),
	);

	/** Raw DashboardBody for write operations that need ISO dates. */
	const rawData = computed<DashboardBody | undefined>(() => query.data.value);

	const isLoading = computed<boolean>(() => query.isPending.value);

	const isError = computed<boolean>(() => query.error.value != null);

	const isEmpty = computed<boolean>(
		() => !isLoading.value && !isError.value && adapted.value.totalPeople === 0,
	);

	return {
		data: adapted,
		rawData,
		isLoading,
		isError,
		isEmpty,
		/** Refresh the underlying query. */
		refresh: () => query.refetch(),
		queryKey: getDashboardPeopleAvailabilityQueryKey(),
	};
}
