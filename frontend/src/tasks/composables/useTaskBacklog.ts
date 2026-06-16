import { computed } from "vue";
import { useQuery } from "@pinia/colada";
import {
	getTasksBacklogQuery,
	getTasksBacklogQueryKey,
} from "@/client/@pinia/colada.gen";
import type { TaskBacklogBody, TaskRow } from "@/client/types.gen";

export interface TaskBacklogState {
	tasks: TaskRow[];
	priorities: TaskBacklogBody["priorities"];
	statuses: TaskBacklogBody["statuses"];
	summary: TaskBacklogBody["summary"];
}

/**
 * Maps the generated API response shape into the panel contract consumed by
 * TaskManagementPanel, with explicit loading / error / empty states.
 */
function adaptToPanelState(
	data: TaskBacklogBody | undefined,
): TaskBacklogState {
	if (!data) {
		return {
			tasks: [],
			priorities: [],
			statuses: [],
			summary: {
				totalTasks: 0,
				highPriorityTasks: 0,
				unassignedTasks: 0,
				understaffedTasks: 0,
			},
		};
	}

	return {
		tasks: data.tasks ?? [],
		priorities: data.priorities ?? [],
		statuses: data.statuses ?? [],
		summary: data.summary,
	};
}

export function useTaskBacklog() {
	const query = useQuery(getTasksBacklogQuery());

	const data = computed<TaskBacklogState>(() =>
		adaptToPanelState(query.data.value),
	);

	const isLoading = computed<boolean>(() => query.isPending.value);

	const isError = computed<boolean>(() => query.error.value != null);

	const isEmpty = computed<boolean>(
		() => !isLoading.value && !isError.value && data.value.tasks.length === 0,
	);

	return {
		data,
		isLoading,
		isError,
		isEmpty,
		/** Refresh the underlying backlog query. */
		refresh: () => query.refetch(),
		queryKey: getTasksBacklogQueryKey(),
	};
}
