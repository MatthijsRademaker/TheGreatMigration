import { beforeEach, describe, expect, it, vi } from "vitest";
import { client } from "../src/client/client.gen";
import {
	configureApiClient,
	resolveApiBaseUrl,
} from "../src/shared/lib/api-client";

describe("api client configuration", () => {
	beforeEach(() => {
		client.setConfig({
			baseUrl: "",
			fetch: undefined,
		});
	});

	it("defaults to same-origin API requests when no base URL is provided", () => {
		expect(resolveApiBaseUrl(undefined)).toBe("");

		configureApiClient({ baseUrl: undefined });

		expect(client.getConfig().baseUrl).toBe("");
	});

	it("uses an explicit runtime API base URL when provided", () => {
		configureApiClient({ baseUrl: "http://example.test:8080" });

		expect(client.getConfig().baseUrl).toBe("http://example.test:8080");
	});

	it("allows SSR and tests to install a custom fetch implementation", () => {
		const fetchMock = vi.fn<typeof fetch>();

		configureApiClient({
			baseUrl: "",
			fetch: fetchMock,
		});

		expect(client.getConfig().fetch).toBe(fetchMock);
	});
});
