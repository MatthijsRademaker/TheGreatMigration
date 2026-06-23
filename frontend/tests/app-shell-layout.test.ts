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
import { computed, defineComponent, nextTick, ref } from "vue";
import { createMemoryHistory, createRouter } from "vue-router";
import { describe, expect, it, beforeEach } from "vitest";
import AppShell from "../src/shared/layout/app-shell/AppShell.vue";

// Stub the planning window composable with a known value
vi.mock("@/shared/composables/usePlanningWindow", () => ({
	usePlanningWindow: () => ({
		formattedRange: "Jul 5–8",
		isLoading: false,
		isError: false,
	}),
}));

// Mutable mock state for useHomePagination — tests can flip these refs to
// exercise loading, error, and normal states.
const mockHomePage = ref(1);
const mockHomeDaysPerPage = ref(4);
const mockHomeTotalPages = ref(10);
const mockHomeRangeLabel = ref("5 Jul – 8 Jul, 2026");
const mockHomeCompactRangeLabel = ref("5–8 Jul");
const mockHomeIsLoading = ref(false);
const mockHomeIsError = ref(false);

const mockHomeGoPrev = vi.fn();
const mockHomeGoNext = vi.fn();
const mockHomeGoToday = vi.fn();

vi.mock("@/shared/composables/useHomePagination", () => ({
	useHomePagination: () => ({
		page: mockHomePage,
		daysPerPage: mockHomeDaysPerPage,
		totalPages: computed(() => mockHomeTotalPages.value),
		rangeLabel: computed(() => mockHomeRangeLabel.value),
		compactRangeLabel: computed(() => mockHomeCompactRangeLabel.value),
		isLoading: computed(() => mockHomeIsLoading.value),
		isError: computed(() => mockHomeIsError.value),
		goPrev: mockHomeGoPrev,
		goNext: mockHomeGoNext,
		goToday: mockHomeGoToday,
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
	beforeEach(() => {
		mockHomePage.value = 1;
		mockHomeDaysPerPage.value = 4;
		mockHomeTotalPages.value = 10;
		mockHomeRangeLabel.value = "5 Jul – 8 Jul, 2026";
		mockHomeCompactRangeLabel.value = "5–8 Jul";
		mockHomeIsLoading.value = false;
		mockHomeIsError.value = false;
		mockHomeGoPrev.mockClear();
		mockHomeGoNext.mockClear();
		mockHomeGoToday.mockClear();
	});

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

	it("renders loading skeleton on home route when pagination isLoading is true", async () => {
		mockHomeIsLoading.value = true;

		const { wrapper } = await mountAppShell();

		const header = wrapper.find("header");
		expect(header.find(".animate-pulse").exists()).toBe(true);
		expect(header.text()).not.toContain("Today");
	});

	it("renders dash fallback on home route when pagination isError is true", async () => {
		mockHomeIsError.value = true;

		const { wrapper } = await mountAppShell();

		const header = wrapper.find("header");
		expect(header.text()).toContain("—");
		expect(header.text()).not.toContain("Today");
	});
});
