// @vitest-environment jsdom
import { describe, it, expect, afterEach, vi } from "vitest";
import { effectScope } from "vue";
import { useMotionPreference } from "../useMotionPreference";
import { INSTANT, springs } from "@/shared/motion/tokens";

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

function run<T>(fn: () => T): T {
	const scope = effectScope();
	const result = scope.run(fn) as T;
	return result;
}

afterEach(() => {
	// Restore the test-default (reduced) matchMedia.
	setReducedMotion(true);
});

describe("useMotionPreference", () => {
	it("reports motion disabled and instant transitions under reduced motion", () => {
		setReducedMotion(true);
		const { enabled, transition } = run(() => useMotionPreference());
		expect(enabled.value).toBe(false);
		expect(transition(springs.default)).toBe(INSTANT);
		expect(transition(springs.bouncy)).toEqual({ duration: 0 });
	});

	it("reports motion enabled and passes transitions through when not reduced", () => {
		setReducedMotion(false);
		const { enabled, transition, countUpDuration } = run(() =>
			useMotionPreference(),
		);
		expect(enabled.value).toBe(true);
		expect(transition(springs.default)).toBe(springs.default);
		expect(countUpDuration.value).toBeGreaterThan(0);
	});
});
