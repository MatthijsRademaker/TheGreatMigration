// @vitest-environment jsdom
import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { computed, nextTick } from "vue";
import ToolsView from "../ToolsView.vue";
import type { Tool } from "../types";

// ── Controllable useTools mock ────────────────────────────────────────────────

interface ToolsMockState {
	tools: Tool[];
	summary: { total: number; claimed: number; open: number };
	isLoading: boolean;
	isError: boolean;
	isEmpty: boolean;
}

let toolsState: ToolsMockState;

const createTool = vi.fn();
const deleteTool = vi.fn();
const claimTool = vi.fn();
const unclaimTool = vi.fn();

function resetToolsState(overrides: Partial<ToolsMockState> = {}) {
	toolsState = {
		tools: overrides.tools ?? [],
		summary: overrides.summary ?? { total: 0, claimed: 0, open: 0 },
		isLoading: overrides.isLoading ?? false,
		isError: overrides.isError ?? false,
		isEmpty: overrides.isEmpty ?? false,
	};
}

vi.mock("@/tools/composables/useTools", () => ({
	useTools: vi.fn(() => ({
		data: computed(() => ({
			tools: toolsState.tools,
			summary: toolsState.summary,
		})),
		isLoading: computed(() => toolsState.isLoading),
		isError: computed(() => toolsState.isError),
		isEmpty: computed(() => toolsState.isEmpty),
		createTool,
		deleteTool,
		claimTool,
		unclaimTool,
		refresh: vi.fn(),
		queryKey: ["getTools"],
	})),
}));

vi.mock("@/shared/composables/usePeopleAvailability", () => ({
	usePeopleAvailability: vi.fn(() => ({
		data: computed(() => ({
			people: [
				{ id: "p1", name: "Sophia Chen", availability: [] },
				{ id: "p2", name: "Marcus Rivera", availability: [] },
			],
		})),
	})),
}));

// ── Tests ─────────────────────────────────────────────────────────────────────

describe("ToolsView", () => {
	beforeEach(() => {
		resetToolsState();
		createTool.mockReset();
		deleteTool.mockReset();
		claimTool.mockReset();
		unclaimTool.mockReset();
	});

	it("shows the empty state when there are no tools", () => {
		resetToolsState({ isEmpty: true });
		const wrapper = mount(ToolsView);
		expect(wrapper.text()).toContain("No tools yet.");
		wrapper.unmount();
	});

	it("shows a shimmer skeleton while loading", () => {
		resetToolsState({ isLoading: true });
		const wrapper = mount(ToolsView);
		expect(wrapper.find("[data-testid='skeleton-rows']").exists()).toBe(true);
		wrapper.unmount();
	});

	it("shows the error state", () => {
		resetToolsState({ isError: true });
		const wrapper = mount(ToolsView);
		expect(wrapper.text()).toContain("Could not load tools");
		wrapper.unmount();
	});

	it("renders open and crossed-off tools, distinguishing them and showing the bringer", () => {
		resetToolsState({
			tools: [
				{ id: "tool-1", name: "Ladder", broughtBy: null },
				{ id: "tool-2", name: "Power drill", broughtBy: "p1" },
			],
			summary: { total: 2, claimed: 1, open: 1 },
		});
		const wrapper = mount(ToolsView);

		expect(wrapper.text()).toContain("Ladder");
		expect(wrapper.text()).toContain("Power drill");
		// Crossed-off tool shows its bringer's name.
		expect(wrapper.text()).toContain("Sophia Chen");

		const claimedRow = wrapper.find("[data-testid='tool-tool-2']");
		expect(claimedRow.attributes("data-claimed")).toBe("true");
		// Crossed-off name uses a line-through.
		expect(claimedRow.find(".line-through").exists()).toBe(true);

		const openRow = wrapper.find("[data-testid='tool-tool-1']");
		expect(openRow.attributes("data-claimed")).toBe("false");
		// Coverage badge.
		expect(wrapper.text()).toContain("1 / 2 covered");
		wrapper.unmount();
	});

	it("calls createTool when a new tool name is submitted", async () => {
		resetToolsState({ isEmpty: true });
		const wrapper = mount(ToolsView);

		const input = wrapper.find("input");
		await input.setValue("Crowbar");
		await wrapper.find("form").trigger("submit.prevent");

		expect(createTool).toHaveBeenCalledWith("Crowbar");
		wrapper.unmount();
	});

	it("calls unclaimTool when the Unclaim button on a crossed-off tool is clicked", async () => {
		resetToolsState({
			tools: [{ id: "tool-2", name: "Power drill", broughtBy: "p1" }],
			summary: { total: 1, claimed: 1, open: 0 },
		});
		const wrapper = mount(ToolsView);

		const unclaimBtn = wrapper
			.findAll("button")
			.find((b) => b.text() === "Unclaim");
		expect(unclaimBtn).toBeTruthy();
		await unclaimBtn!.trigger("click");

		expect(unclaimTool).toHaveBeenCalledWith("tool-2");
		wrapper.unmount();
	});

	it("calls deleteTool when the remove control is clicked", async () => {
		resetToolsState({
			tools: [{ id: "tool-1", name: "Ladder", broughtBy: null }],
			summary: { total: 1, claimed: 0, open: 1 },
		});
		const wrapper = mount(ToolsView);
		await nextTick();

		const removeBtn = wrapper.find("[aria-label='Remove Ladder']");
		expect(removeBtn.exists()).toBe(true);
		await removeBtn.trigger("click");

		expect(deleteTool).toHaveBeenCalledWith("tool-1");
		wrapper.unmount();
	});
});
