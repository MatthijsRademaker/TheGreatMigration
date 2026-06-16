import { describe, expect, it } from "vitest";
import { getHttpErrorStatus } from "../../src/shared/lib/errorStatus";

describe("getHttpErrorStatus", () => {
	it("returns the status from an ErrorModel-shaped object", () => {
		expect(getHttpErrorStatus({ status: 409 })).toBe(409);
		expect(getHttpErrorStatus({ status: 404 })).toBe(404);
		expect(getHttpErrorStatus({ status: 400 })).toBe(400);
		expect(getHttpErrorStatus({ status: 500 })).toBe(500);
	});

	it("returns undefined for non-ErrorModel errors", () => {
		expect(getHttpErrorStatus(new Error("network error"))).toBeUndefined();
		expect(getHttpErrorStatus("plain string")).toBeUndefined();
		expect(getHttpErrorStatus(null)).toBeUndefined();
		expect(getHttpErrorStatus(undefined)).toBeUndefined();
		expect(getHttpErrorStatus(42)).toBeUndefined();
	});

	it("returns undefined when status is not a number", () => {
		expect(getHttpErrorStatus({ status: "409" })).toBeUndefined();
		expect(getHttpErrorStatus({ status: null })).toBeUndefined();
	});

	it("returns undefined for empty objects", () => {
		expect(getHttpErrorStatus({})).toBeUndefined();
	});
});
