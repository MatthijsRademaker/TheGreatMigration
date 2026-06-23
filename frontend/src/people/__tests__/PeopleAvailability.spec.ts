// @vitest-environment jsdom
import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { Badge } from "@/shared/ui/badge";
import PeopleAvailability from "../PeopleAvailability.vue";

describe("PeopleAvailability", () => {
	it("does not render availability legend footer", () => {
		const wrapper = mount(PeopleAvailability);
		expect(wrapper.findAllComponents(Badge)).toHaveLength(16);
	});
});
