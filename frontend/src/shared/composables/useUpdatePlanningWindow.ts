import { useMutation, useQueryCache } from "@pinia/colada";
import {
	getPlanningWindowQueryKey,
	putPlanningWindowMutation,
} from "@/client/@pinia/colada.gen";

export function useUpdatePlanningWindow() {
	const queryCache = useQueryCache();

	const {
		mutate,
		isLoading: isPending,
		error,
	} = useMutation({
		...putPlanningWindowMutation(),
		onSuccess: () => {
			queryCache.invalidateQueries({
				key: getPlanningWindowQueryKey(),
				exact: true,
			});
		},
	});

	return {
		mutate,
		isPending,
		error,
	};
}
