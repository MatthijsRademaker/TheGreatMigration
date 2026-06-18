// @vitest-environment jsdom

// ---------------------------------------------------------------------------
// Must be defined before any module imports that use @vueuse/core.
// ---------------------------------------------------------------------------
import { vi } from "vitest";

// Stub matchMedia for all breakpoints.
// Must execute before module imports because @vueuse/core checks
// window.matchMedia synchronously during module initialization.
const matchMediaStub = vi.fn((query: string) => ({
	matches: query === "(max-width: 768px)",
	media: query,
	onchange: null,
	addEventListener: vi.fn(),
	removeEventListener: vi.fn(),
}));

try {
	Object.defineProperty(window, "matchMedia", {
		writable: true,
		configurable: true,
		value: matchMediaStub,
	});
} catch {
	// Global already stubbed
}

// ---------------------------------------------------------------------------
import { mount } from "@vue/test-utils";
import { defineComponent, nextTick } from "vue";
import { createMemoryHistory, createRouter } from "vue-router";
import { describe, expect, it } from "vitest";
import AppShell from "../src/shared/layout/app-shell/AppShell.vue";

// Stub the planning window composable with a known value
vi.mock("@/shared/composables/usePlanningWindow", () => ({
	usePlanningWindow: () => ({
		formattedRange: "Jul 5–8",
		isLoading: false,
		isError: false,
	}),
}));

// Stub the home pagination composable
vi.mock("@/shared/composables/useHomePagination", () => ({
	useHomePagination: () => ({
		page: { value: 1 },
		daysPerPage: { value: 4 },
		totalPages: { value: 10 },
		start: { value: "2026-07-05" },
		rangeLabel: { value: "5 Jul – 8 Jul, 2026" },
		isLoading: { value: false },
		isError: { value: false },
		goPrev: vi.fn(),
		goNext: vi.fn(),
		goToday: vi.fn(),
	}),
}));

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
async function mountAppShell() {
	const router = createRouter({
		history: createMemoryHistory(),
		routes: [
			{
				path: "/",
				name: "dashboard",
				component: defineComponent({
					template: '<div class="dashboard-page">Dashboard</div>',
				}),
			},
		],
	});
	await router.push("/");
	await router.isReady();

	const wrapper = mount(AppShell, {
		global: { plugins: [router] },
	});

	await nextTick();
	await nextTick();
	await nextTick();

	return { wrapper, router };
}

// ---------------------------------------------------------------------------
// Tests — scrolling and sticky-header layout verification
// ---------------------------------------------------------------------------
describe("AppShell scroll and sticky header layout", () => {
	it("renders the SidebarInset with overflow-y-auto class for vertical scrolling", async () => {
		const { wrapper } = await mountAppShell();

		const sidebarInset = wrapper.findComponent({ name: "SidebarInset" });
		expect(sidebarInset.exists()).toBe(true);

		const classes = sidebarInset.classes();
		expect(classes).toContain("overflow-y-auto");
		expect(classes).not.toContain("overflow-hidden");
	});

	it("renders the header with sticky position classes (sticky, top-0, z-10)", async () => {
		const { wrapper } = await mountAppShell();

		const header = wrapper.find("header");
		expect(header.exists()).toBe(true);

		const classes = header.classes();
		expect(classes).toContain("sticky");
		expect(classes).toContain("top-0");
		expect(classes).toContain("z-10");
	});

	it("renders the header with visual separation classes (border-b, backdrop-blur, bg-background)", async () => {
		const { wrapper } = await mountAppShell();

		const header = wrapper.find("header");
		expect(header.exists()).toBe(true);

		const classes = header.classes();
		expect(classes).toContain("border-b");
		expect(classes).toContain("backdrop-blur");
		expect(classes).toContain("bg-background/90");
	});

	it("renders the mobile SidebarTrigger inside the header", async () => {
		const { wrapper } = await mountAppShell();

		const header = wrapper.find("header");
		const trigger = header.find('[data-sidebar="trigger"]');
		expect(trigger.exists()).toBe(true);
	});

	it("renders the timeline toolbar on the home route", async () => {
		const { wrapper } = await mountAppShell();

		const header = wrapper.find("header");
		expect(header.text()).toContain("Today");
		expect(header.text()).toContain("5 Jul – 8 Jul, 2026");
	});

	it("renders the static planning window range on non-home routes", async () => {
		const router = createRouter({
			history: createMemoryHistory(),
			routes: [
				{
					path: "/",
					name: "dashboard",
					component: defineComponent({
						template: '<div class="dashboard-page">Dashboard</div>',
					}),
				},
				{
					path: "/calendar",
					name: "calendar",
					component: defineComponent({
						template: '<div class="calendar-page">Calendar</div>',
					}),
				},
			],
		});
		await router.push("/calendar");
		await router.isReady();

		const wrapper = mount(AppShell, {
			global: { plugins: [router] },
		});

		await nextTick();
		await nextTick();
		await nextTick();

		const header = wrapper.find("header");
		// On non-home routes, the static planning-window range is shown
		expect(header.text()).toContain("Jul 5–8");
		// The Today button should not be present on non-home routes
		expect(header.text()).not.toContain("Today");
	});
});
