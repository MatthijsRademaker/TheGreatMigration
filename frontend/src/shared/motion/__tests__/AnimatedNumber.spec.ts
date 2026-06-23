// @vitest-environment jsdom
import { describe, it, expect, vi, afterEach } from "vitest";
import { mount } from "@vue/test-utils";
import { nextTick } from "vue";
import AnimatedNumber from "../AnimatedNumber.vue";

function setReducedMotion(reduce: boolean) {
	window.matchMedia = vi.fn().mockImplementation((query: string) => ({
		matches: query.includes("prefers-reduced-motion") ? reduce : false,
		media: query,
		onchange: null,
		addEventListener: vi.fn(),
		removeEventListener: vi.fn(),
		addListener: vi.fn(),
		removeListener: vi.fn(),
		dispatchEvent: vi.fn(),
	}));
}

afterEach(() => {
	setReducedMotion(true);
	vi.useRealTimers();
});

describe("AnimatedNumber", () => {
	it("shows the exact value immediately under reduced motion (no count-up)", () => {
		setReducedMotion(true);
		const wrapper = mount(AnimatedNumber, { props: { value: 42 } });
		expect(wrapper.text()).toBe("42");
		wrapper.unmount();
	});

	it("reflects value changes under reduced motion", async () => {
		setReducedMotion(true);
		const wrapper = mount(AnimatedNumber, { props: { value: 3 } });
		expect(wrapper.text()).toBe("3");
		await wrapper.setProps({ value: 9 });
		expect(wrapper.text()).toBe("9");
		wrapper.unmount();
	});

	it("starts the count-up from a baseline when motion is enabled", async () => {
		setReducedMotion(false);
		const wrapper = mount(AnimatedNumber, { props: { value: 50 } });
		await nextTick();
		// With motion enabled the intro count-up runs from 0, so the value shown
		// just after mount is below the target (it has not snapped to the final
		// value the way reduced motion does).
		expect(Number(wrapper.text())).toBeLessThan(50);
		wrapper.unmount();
	});
});
