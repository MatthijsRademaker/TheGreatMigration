import { describe, it, expect } from "vitest";
import { areaColor } from "../areaColor";

describe("areaColor", () => {
	it("is deterministic: same id yields the same color", () => {
		expect(areaColor("room-1")).toBe(areaColor("room-1"));
		expect(areaColor("room-42")).toBe(areaColor("room-42"));
	});

	it("returns a hex color from the palette", () => {
		expect(areaColor("room-1")).toMatch(/^#[0-9a-f]{6}$/i);
	});

	it("depends on the id, not the name (different ids can differ)", () => {
		// Across a spread of ids we expect more than one distinct color.
		const colors = new Set(
			Array.from({ length: 8 }, (_, i) => areaColor(`room-${i + 1}`)),
		);
		expect(colors.size).toBeGreaterThan(1);
	});

	it("handles empty and arbitrary strings without throwing", () => {
		expect(areaColor("")).toMatch(/^#[0-9a-f]{6}$/i);
		expect(areaColor("Living Room")).toMatch(/^#[0-9a-f]{6}$/i);
	});
});
