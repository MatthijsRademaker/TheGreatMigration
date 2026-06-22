import { computed } from "vue";
import { useQuery, useMutation, useQueryCache } from "@pinia/colada";
import {
	getToolsQuery,
	getToolsQueryKey,
	createToolMutation,
	deleteToolMutation,
	setToolBringerMutation,
	clearToolBringerMutation,
} from "@/client/@pinia/colada.gen";
import type { ToolsBody } from "@/client/types.gen";
import type { Tool, ToolSummary } from "../types";

export interface ToolsState {
	tools: Tool[];
	summary: ToolSummary;
}

/**
 * Maps the generated API response shape into the component contract, with
 * explicit loading / error / empty states. Modeled on useTaskBacklog.
 */
function adaptToState(data: ToolsBody | undefined): ToolsState {
	if (!data) {
		return {
			tools: [],
			summary: { total: 0, claimed: 0, open: 0 },
		};
	}
	return {
		tools: (data.tools ?? []).map((tool) => ({
			id: tool.id,
			name: tool.name,
			broughtBy: tool.broughtBy ?? null,
		})),
		summary: data.summary,
	};
}

export function useTools() {
	const query = useQuery(getToolsQuery());
	const queryKey = getToolsQueryKey();
	const queryCache = useQueryCache();

	const invalidate = () => queryCache.invalidateQueries({ key: queryKey });

	const createMut = useMutation({
		...createToolMutation(),
		onSuccess: invalidate,
	});
	const deleteMut = useMutation({
		...deleteToolMutation(),
		onSuccess: invalidate,
	});
	const claimMut = useMutation({
		...setToolBringerMutation(),
		onSuccess: invalidate,
	});
	const unclaimMut = useMutation({
		...clearToolBringerMutation(),
		onSuccess: invalidate,
	});

	const data = computed<ToolsState>(() => adaptToState(query.data.value));

	const isLoading = computed<boolean>(() => query.isPending.value);
	const isError = computed<boolean>(() => query.error.value != null);
	const isEmpty = computed<boolean>(
		() => !isLoading.value && !isError.value && data.value.tools.length === 0,
	);

	return {
		data,
		isLoading,
		isError,
		isEmpty,
		/** Create a tool by name. */
		createTool: (name: string) => createMut.mutateAsync({ body: { name } }),
		/** Delete a tool by ID. */
		deleteTool: (id: string) => deleteMut.mutateAsync({ path: { id } }),
		/** Claim a tool: set its bringer. */
		claimTool: (id: string, personId: string) =>
			claimMut.mutateAsync({ path: { id }, body: { personId } }),
		/** Unclaim a tool: clear its bringer. */
		unclaimTool: (id: string) => unclaimMut.mutateAsync({ path: { id } }),
		/** Refresh the underlying tools query. */
		refresh: () => query.refetch(),
		queryKey,
	};
}
