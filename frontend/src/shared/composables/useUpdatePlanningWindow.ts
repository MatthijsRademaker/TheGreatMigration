import { useMutation, useQueryCache } from "@pinia/colada";
import { computed } from "vue";
import {
	getPlanningWindowQueryKey,
	putPlanningWindowMutation,
} from "@/client/@pinia/colada.gen";

export function useUpdatePlanningWindow() {
	const queryCache = useQueryCache();

	const { mutate, status, error } = useMutation({
		...putPlanningWindowMutation(),
		onSuccess: () => {
			queryCache.invalidateQueries({
				key: getPlanningWindowQueryKey(),
				exact: true,
			});
		},
	});

	const isPending = computed(() => status.value === "pending");

	return {
		mutate,
		isPending,
		error,
	};
}
