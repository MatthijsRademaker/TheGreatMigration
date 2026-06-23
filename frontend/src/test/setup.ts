import { vi } from "vitest";

// jsdom does not implement ResizeObserver / IntersectionObserver, which motion-v
// uses on client mount (e.g. for `layout` measurement). Without these, mounting
// a <Motion layout> component throws and aborts the surrounding render. Provide
// inert stubs so motion-driven components mount cleanly under test. Real
// browsers supply these natively.
class StubObserver {
	observe() {}
	unobserve() {}
	disconnect() {}
	takeRecords() {
		return [];
	}
}
if (typeof globalThis.ResizeObserver === "undefined") {
	globalThis.ResizeObserver = StubObserver as unknown as typeof ResizeObserver;
}
if (typeof globalThis.IntersectionObserver === "undefined") {
	globalThis.IntersectionObserver =
		StubObserver as unknown as typeof IntersectionObserver;
}

// jsdom does not implement matchMedia. Provide a stub that reports
// `prefers-reduced-motion: reduce` so motion is gated OFF by default in the
// test environment — component output is then deterministic (final state, no
// rAF-driven animation that never advances under vitest). Tests that exercise
// the motion-enabled path override matchMedia explicitly.
if (typeof window !== "undefined" && !window.matchMedia) {
	window.matchMedia = vi.fn().mockImplementation((query: string) => ({
		matches: query.includes("prefers-reduced-motion"),
		media: query,
		onchange: null,
		addEventListener: vi.fn(),
		removeEventListener: vi.fn(),
		addListener: vi.fn(),
		removeListener: vi.fn(),
		dispatchEvent: vi.fn(),
	}));
}
