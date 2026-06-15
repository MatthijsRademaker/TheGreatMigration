export const routes = [
	{
		path: "/",
		name: "home",
		component: () => import("@/home/HomeView.vue"),
		meta: {
			title: "Moving dashboard",
			description:
				"Today’s priorities, staffing gaps, and move notes at a glance.",
		},
	},
	{
		path: "/tasks",
		name: "tasks",
		component: () => import("@/tasks/TasksView.vue"),
		meta: {
			title: "Task backlog",
			description:
				"Capture jobs, priorities, staffing needs, and planning status.",
		},
	},
	{
		path: "/calendar",
		name: "calendar",
		component: () => import("@/calendar/CalendarView.vue"),
		meta: {
			title: "Schedule board",
			description: "Plan work across move days and balance available helpers.",
		},
	},
	{
		path: "/people",
		name: "people",
		component: () => import("@/people/PeopleView.vue"),
		meta: {
			title: "People availability",
			description: "Track who is available and where each person can help.",
		},
	},
];
