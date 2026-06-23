// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { effectScope } from "vue";

const setEnabled = vi.fn();
const parentRef = { value: null };

vi.mock("@formkit/auto-animate/vue", () => ({
	useAutoAnimate: vi.fn(() => [parentRef, setEnabled]),
}));

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

import { useMotionAutoAnimate } from "../useMotionAutoAnimate";

beforeEach(() => setEnabled.mockClear());

describe("useMotionAutoAnimate", () => {
	it("enables list animations when motion is allowed", () => {
		setReducedMotion(false);
		const ref = effectScope().run(() => useMotionAutoAnimate());
		expect(ref).toBe(parentRef);
		expect(setEnabled).toHaveBeenCalledWith(true);
	});

	it("disables list animations under reduced motion", () => {
		setReducedMotion(true);
		effectScope().run(() => useMotionAutoAnimate());
		expect(setEnabled).toHaveBeenCalledWith(false);
	});
});
