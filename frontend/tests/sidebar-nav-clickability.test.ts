// @vitest-environment jsdom

// ---------------------------------------------------------------------------
// Must be defined before any module imports that use @vueuse/core.
// ---------------------------------------------------------------------------
import { vi } from "vitest";

// ---------------------------------------------------------------------------
import { flushPromises, mount } from "@vue/test-utils";
import { defineComponent, nextTick } from "vue";
import { createMemoryHistory, createRouter } from "vue-router";
import { beforeAll, describe, expect, it } from "vitest";
import AppSidebar from "../src/shared/layout/app-sidebar/AppSidebar.vue";
import { SidebarProvider } from "../src/shared/ui/sidebar";

const NAV_ITEMS = [
	{ title: "Dashboard", to: "/" },
	{ title: "Tasks", to: "/tasks" },
	{ title: "Schedule", to: "/calendar" },
	{ title: "People", to: "/people" },
	{ title: "Rooms / Areas", to: "/rooms" },
	{ title: "Settings", to: "/settings" },
] as const;

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
describe("sidebar nav-item clickability", () => {
	beforeAll(() => {
		// Stub matchMedia for DESKTOP viewport so the sidebar renders inline
		// (not in a teleported mobile Sheet). This makes nav-item DOM access
		// straightforward without teleportation or Sheet-opening ceremony.
		vi.stubGlobal(
			"matchMedia",
			vi.fn((query: string) => ({
				matches: false,
				media: query,
				onchange: null,
				addEventListener: vi.fn(),
				removeEventListener: vi.fn(),
			})),
		);
	});

	async function mountSidebarDesktop(initialPath = "/") {
		const router = createRouter({
			history: createMemoryHistory(),
			routes: [
				{
					path: "/",
					name: "home",
					component: defineComponent({ template: "<div>Home</div>" }),
				},
				{
					path: "/people",
					name: "people",
					component: defineComponent({ template: "<div>People</div>" }),
				},
				{
					path: "/tasks",
					name: "tasks",
					component: defineComponent({ template: "<div>Tasks</div>" }),
				},
				{
					path: "/calendar",
					name: "calendar",
					component: defineComponent({ template: "<div>Calendar</div>" }),
				},
				{
					path: "/rooms",
					name: "rooms",
					component: defineComponent({ template: "<div>Rooms</div>" }),
				},
				{
					path: "/settings",
					name: "settings",
					component: defineComponent({ template: "<div>Settings</div>" }),
				},
			],
		});
		await router.push(initialPath);
		await router.isReady();

		const TestWrapper = defineComponent({
			components: { SidebarProvider, AppSidebar },
			template: `
				<SidebarProvider>
					<AppSidebar />
				</SidebarProvider>
			`,
		});

		const wrapper = mount(TestWrapper, { global: { plugins: [router] } });
		await nextTick();
		await nextTick();
		await nextTick();

		return { wrapper, router };
	}

	it("renders a native <a> for each nav item with correct href", async () => {
		const { wrapper } = await mountSidebarDesktop();

		for (const item of NAV_ITEMS) {
			// Exclude the header brand link (has aria-label="Open moving dashboard")
			// which shares href="/" with the Dashboard nav item
			const selector = `a[data-sidebar="menu-button"]:not([aria-label])[href="${item.to}"]`;
			const link = wrapper.find(selector);
			expect(link.exists(), `Expected ${selector}`).toBe(true);
		}
	});

	it("has a title attribute on each nav <a> for collapsed-mode tooltip", async () => {
		const { wrapper } = await mountSidebarDesktop();

		for (const item of NAV_ITEMS) {
			const selector = `a[data-sidebar="menu-button"]:not([aria-label])[href="${item.to}"]`;
			const link = wrapper.find(selector);
			expect(link.exists(), `Expected ${selector}`).toBe(true);
			expect(link.attributes("title"), `${item.title}: title attribute`).toBe(
				item.title,
			);
		}
	});

	it("renders native <a> elements with href and title for SPA navigation", async () => {
		const { wrapper, router } = await mountSidebarDesktop();

		for (const item of NAV_ITEMS) {
			const selector = `a[data-sidebar="menu-button"]:not([aria-label])[href="${item.to}"]`;
			const link = wrapper.find(selector);
			expect(link.exists(), `Expected ${selector}`).toBe(true);
			expect(link.attributes("href")).toBe(item.to);
			expect(link.attributes("title")).toBe(item.title);
		}

		// Verify navigation routing works via click on a sidebar nav <a> element.
		const peopleLink = wrapper.find(
			'a[data-sidebar="menu-button"]:not([aria-label])[href="/people"]',
		);
		await peopleLink.trigger("click");
		await flushPromises();
		await nextTick();
		expect(router.currentRoute.value.path).toBe("/people");
	});

	it("does not contain duplicate branding in the footer", async () => {
		const { wrapper } = await mountSidebarDesktop();

		const footer = wrapper.find('[data-slot="sidebar-footer"]');
		expect(footer.exists()).toBe(true);

		expect(footer.text()).not.toContain("The Great Migration");
		expect(footer.text()).not.toContain("House move planner");
	});
});
