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
	{
		path: "/tools",
		name: "tools",
		component: () => import("@/tools/ToolsView.vue"),
		meta: {
			title: "Tools",
			description:
				"Track required tools and who is bringing each item on move day.",
		},
	},
	{
		path: "/rooms",
		name: "rooms",
		component: () => import("@/rooms/RoomsView.vue"),
		meta: {
			title: "Rooms / Areas",
			description:
				"Organize and label rooms, floors, and zones for a clear move-day plan.",
		},
	},
	{
		path: "/settings",
		name: "settings",
		component: () => import("@/settings/SettingsView.vue"),
		meta: {
			title: "Settings",
			description:
				"Configure your move preferences, notification defaults, and account details.",
		},
	},
];
