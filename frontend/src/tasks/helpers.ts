import type { TaskRow, TaskPriority } from "./types";
import type { BadgeVariants } from "@/shared/ui/badge";

/**
 * Centralized display mapping for task row status.
 * Derives the display label from empty `assignedTo` rather than
 * introducing a new canonical status value.
 */
export function getTaskDisplayState(task: TaskRow): string {
	if (task.assignedTo.length === 0) {
		return "Unassigned";
	}
	return task.status;
}

export const priorityLabels: Record<TaskPriority, string> = {
	high: "High",
	medium: "Medium",
	low: "Low",
};

export const priorityBadgeVariant: Record<
	TaskPriority,
	NonNullable<BadgeVariants["variant"]>
> = {
	high: "priorityHigh",
	medium: "priorityMedium",
	low: "priorityLow",
};
