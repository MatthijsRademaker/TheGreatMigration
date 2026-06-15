import { createRouter, createWebHistory } from "vue-router";

export const router = createRouter({
	history: createWebHistory(),
	routes: [
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
				description:
					"Plan work across move days and balance available helpers.",
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
			path: "/showcase",
			name: "showcase",
			component: () => import("@/showcase/ShowcaseView.vue"),
			meta: {
				title: "Component showcase",
				description:
					"Developer preview of reusable UI components and patterns.",
			},
		},
	],
});
