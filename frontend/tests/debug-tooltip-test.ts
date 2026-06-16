// @vitest-environment jsdom
import { vi } from "vitest";
import { mount } from "@vue/test-utils";
import { defineComponent, nextTick } from "vue";
import { createMemoryHistory, createRouter } from "vue-router";
import { describe, expect, it } from "vitest";
import AppSidebar from "../src/shared/layout/app-sidebar/AppSidebar.vue";
import { SidebarProvider } from "../src/shared/ui/sidebar";

describe("tooltip debug", () => {
	it("checks tooltip", async () => {
		vi.stubGlobal(
			"matchMedia",
			vi.fn(() => ({
				matches: false,
				media: "",
				onchange: null,
				addEventListener: vi.fn(),
				removeEventListener: vi.fn(),
			})),
		);
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
			],
		});
		await router.push("/");
		await router.isReady();
		const TestWrapper = defineComponent({
			components: { SidebarProvider, AppSidebar },
			template: `<SidebarProvider><AppSidebar /></SidebarProvider>`,
		});
		const wrapper = mount(TestWrapper, { global: { plugins: [router] } });
		await nextTick();
		await nextTick();
		await nextTick();
		expect(wrapper.findAll('[data-slot="tooltip"]').length).toBe(1);
	});
});
