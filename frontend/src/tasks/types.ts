export type TaskPriority = "high" | "medium" | "low";

export interface Area {
	id: string;
	name: string;
}

export interface TaskRow {
	id: string;
	title: string;
	priority: TaskPriority;
	peopleNeeded: number;
	area: Area;
	status: "backlog" | "ready" | "assigned";
	assignedTo: string[];
}
