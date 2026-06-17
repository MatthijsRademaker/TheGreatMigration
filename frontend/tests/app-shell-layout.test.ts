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
});
