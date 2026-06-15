import { computed } from "vue";
import { useQuery } from "@pinia/colada";
import {
	getPlanningWindowQuery,
	getPlanningWindowQueryKey,
} from "@/client/@pinia/colada.gen";
import type { PlanWindowDay } from "@/shared/lib/planWindow";
import { formatPlanDayLabel } from "@/shared/lib/planWindow";

export function usePlanningWindow() {
	const query = useQuery(getPlanningWindowQuery());

	const planWindowDays = computed<PlanWindowDay[]>(() => {
		const data = query.data.value;
		if (!data) return [];

		const start = new Date(data.startDate);
		const end = new Date(data.endDate);
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
	});

	const planWindowDayCount = computed<number>(
		() => planWindowDays.value.length,
	);

	const isLoading = computed<boolean>(() => query.isPending.value);

	const isError = computed<boolean>(() => query.error.value != null);

	return {
		planWindowDays,
		planWindowDayCount,
		isLoading,
		isError,
		queryKey: getPlanningWindowQueryKey(),
	};
}
