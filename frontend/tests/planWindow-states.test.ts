import { ref } from "vue";
import { renderToString } from "@vue/server-renderer";
import { createSSRApp, defineComponent, h } from "vue";
import { afterEach, describe, expect, it, vi } from "vitest";

// Per-test reactive refs that the mocked composable returns.
// eslint-disable-next-line no-var
var mockIsLoading = ref(false);
// eslint-disable-next-line no-var
var mockIsError = ref(false);

vi.mock("@/shared/composables/usePlanningWindow", () => ({
	usePlanningWindow: () => ({
		planWindowDays: ref([]),
		planWindowDayCount: ref(0),
		isLoading: mockIsLoading,
		isError: mockIsError,
		queryKey: [],
	}),
}));

import { usePlanningWindow } from "../src/shared/composables/usePlanningWindow";

/**
 * Render the composable via a minimal SSR component that mirrors the
 * testids used by the success-path test in planWindow.test.ts.
 */
async function renderState() {
	const TestComponent = defineComponent({
		setup() {
			const state = usePlanningWindow();
			return () =>
				h("div", [
					h(
						"span",
						{ "data-testid": "loading" },
						String(state.isLoading.value),
					),
					h("span", { "data-testid": "error" }, String(state.isError.value)),
				]);
		},
	});

	const app = createSSRApp(TestComponent);
	return renderToString(app);
}

describe("usePlanningWindow loading and error states", () => {
	afterEach(() => {
		mockIsLoading.value = false;
		mockIsError.value = false;
	});

	it("exposes isLoading=true when the query is pending", async () => {
		mockIsLoading.value = true;

		const html = await renderState();

		expect(html).toContain('data-testid="loading">true<');
		expect(html).toContain('data-testid="error">false<');
	});

	it("exposes isError=true when the query has errored", async () => {
		mockIsError.value = true;

		const html = await renderState();

		expect(html).toContain('data-testid="loading">false<');
		expect(html).toContain('data-testid="error">true<');
	});
});
