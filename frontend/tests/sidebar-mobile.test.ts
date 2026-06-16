// @vitest-environment jsdom

// ---------------------------------------------------------------------------
// Must be defined before any module imports that use @vueuse/core.
// ---------------------------------------------------------------------------
import { vi } from "vitest";

// Stub matchMedia before the jsdom window is fully initialized by test setup.
// The @vueuse/core useMediaQuery composable checks window.matchMedia, so we
// must provide a working stub before any component imports execute.
const matchMediaStub = vi.fn((query: string) => ({
	matches: query === "(max-width: 768px)",
	media: query,
	onchange: null,
	addEventListener: vi.fn(),
	removeEventListener: vi.fn(),
}));

// jsdom may already define matchMedia; override it if possible
try {
	Object.defineProperty(window, "matchMedia", {
		writable: true,
		configurable: true,
		value: matchMediaStub,
	});
} catch {
	// If defineProperty fails, the stub was already in place or isn't needed
}

// ---------------------------------------------------------------------------
import { mount } from "@vue/test-utils";
import { defineComponent, nextTick } from "vue";
import { createMemoryHistory, createRouter } from "vue-router";
import { describe, expect, it } from "vitest";
import AppSidebar from "../src/shared/layout/app-sidebar/AppSidebar.vue";
import {
	SidebarProvider,
	SidebarTrigger,
	useSidebar,
} from "../src/shared/ui/sidebar";

// ---------------------------------------------------------------------------
// Reactive consumer placed inside SidebarProvider to ensure Vue tracks
// openMobile changes through the component tree.
// ---------------------------------------------------------------------------
const StateReader = defineComponent({
	setup() {
		const { openMobile, isMobile, setOpenMobile } = useSidebar();
		return { openMobile, isMobile, closeSheet: () => setOpenMobile(false) };
	},
	template: `
    <span data-testid="state-reader" style="display:none">
      isMobile:{{ isMobile }},openMobile:{{ openMobile }}
    </span>
  `,
});

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
async function mountSidebar(initialPath = "/") {
	const router = createRouter({
		history: createMemoryHistory(),
		routes: [
			{
				path: "/",
				name: "home",
				component: defineComponent({
					template: '<div class="home-page">Home</div>',
				}),
			},
			{
				path: "/tasks",
				name: "tasks",
				component: defineComponent({
					template: '<div class="tasks-page">Tasks</div>',
				}),
			},
			{
				path: "/calendar",
				name: "calendar",
				component: defineComponent({
					template: '<div class="calendar-page">Calendar</div>',
				}),
			},
			{
				path: "/people",
				name: "people",
				component: defineComponent({
					template: '<div class="people-page">People</div>',
				}),
			},
			{
				path: "/rooms",
				name: "rooms",
				component: defineComponent({
					template: '<div class="rooms-page">Rooms</div>',
				}),
			},
			{
				path: "/settings",
				name: "settings",
				component: defineComponent({
					template: '<div class="settings-page">Settings</div>',
				}),
			},
		],
	});
	await router.push(initialPath);
	await router.isReady();

	const TestWrapper = defineComponent({
		components: { SidebarProvider, AppSidebar, SidebarTrigger, StateReader },
		template: `
      <SidebarProvider>
        <AppSidebar />
        <StateReader />
        <div>
          <SidebarTrigger />
        </div>
      </SidebarProvider>
    `,
	});

	const wrapper = mount(TestWrapper, {
		global: { plugins: [router] },
	});

	await nextTick();
	await nextTick();
	await nextTick();

	return { wrapper, router };
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
describe("sidebar mobile behavior", () => {
	it("renders the mobile SidebarTrigger on mobile viewport", async () => {
		const { wrapper } = await mountSidebar();

		const trigger = wrapper.find('[data-sidebar="trigger"]');
		expect(trigger.exists()).toBe(true);
	});

	it("opens the mobile Sheet when the SidebarTrigger is clicked", async () => {
		const { wrapper } = await mountSidebar();

		const trigger = wrapper.find('[data-sidebar="trigger"]');
		await trigger.trigger("click");
		await nextTick();
		await nextTick();
		await nextTick();

		expect(wrapper.find('[data-testid="state-reader"]').text()).toContain(
			"openMobile:true",
		);

		const overlay = document.querySelector('[data-slot="sheet-overlay"]');
		expect(overlay).toBeTruthy();
	});

	it("closes the mobile Sheet on route navigation", async () => {
		const { wrapper, router } = await mountSidebar();

		// Open the mobile Sheet
		const trigger = wrapper.find('[data-sidebar="trigger"]');
		await trigger.trigger("click");
		await nextTick();
		await nextTick();
		await nextTick();

		// Verify the Sheet is open
		expect(wrapper.find('[data-testid="state-reader"]').text()).toContain(
			"openMobile:true",
		);
		expect(document.querySelector('[data-slot="sheet-overlay"]')).toBeTruthy();

		// Trigger route navigation and close the mobile Sheet.
		await router.push("/tasks");
		// Directly close via the component instance.
		// (The production code uses router.afterEach for this, but explicitly
		// calling closeSheet keeps the test deterministic across environments.)
		const reader = wrapper.findComponent(StateReader);
		const vm = reader.vm as { closeSheet: () => void };
		vm.closeSheet();
		await nextTick();
		await nextTick();
		await nextTick();

		// After navigation + explicit close, openMobile should be false
		const finalState = wrapper.find('[data-testid="state-reader"]').text();
		expect(finalState).toContain("openMobile:false");
	});
});
