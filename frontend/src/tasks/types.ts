export type TaskPriority = "high" | "medium" | "low";

export interface TaskRow {
	id: string;
	title: string;
	priority: TaskPriority;
	peopleNeeded: number;
	room: string;
	status: string;
	assignedTo: string[];
}
