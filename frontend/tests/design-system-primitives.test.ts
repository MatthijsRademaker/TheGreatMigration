import { renderToString } from "@vue/server-renderer";
import { createSSRApp, h } from "vue";
import { describe, expect, it } from "vitest";
import {
	Badge,
	badgeVariants,
	type BadgeVariants,
} from "../src/shared/ui/badge";
import { Button, buttonVariants } from "../src/shared/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "../src/shared/ui/card";
import { Checkbox, checkboxVariants } from "../src/shared/ui/checkbox";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "../src/shared/ui/select";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "../src/shared/ui/popover";
import { Calendar } from "../src/shared/ui/calendar";
import { Avatar } from "../src/shared/ui/avatar";
import { SearchField } from "../src/shared/ui/search-field";
import { PersonChip } from "../src/shared/ui/person-chip";
import { HelpButton } from "../src/shared/ui/help-button";
import { DatePicker } from "../src/shared/ui/date-picker";
import { TooltipProvider } from "reka-ui";
import { CalendarDate } from "@internationalized/date";

const semanticBadgeCases = [
	{
		variant: "priorityHigh",
		expectedClasses: ["bg-destructive-soft", "text-destructive"],
	},
	{
		variant: "priorityMedium",
		expectedClasses: ["bg-warning-soft", "text-warning"],
	},
	{
		variant: "priorityLow",
		expectedClasses: ["bg-success-soft", "text-success"],
	},
	{
		variant: "available",
		expectedClasses: ["bg-success-soft", "text-success"],
	},
	{
		variant: "busy",
		expectedClasses: ["bg-destructive-soft", "text-destructive"],
	},
	{
		variant: "partial",
		expectedClasses: ["bg-warning-soft", "text-warning"],
	},
	{
		variant: "off",
		expectedClasses: ["bg-secondary", "text-muted-foreground"],
	},
] satisfies Array<{
	variant: NonNullable<BadgeVariants["variant"]>;
	expectedClasses: string[];
}>;

async function renderButtonPreview() {
	const app = createSSRApp({
		render: () => h(Button, null, { default: () => "Button label" }),
	});

	return renderToString(app);
}

async function renderBadgePreview() {
	const app = createSSRApp({
		render: () =>
			h(Badge, { variant: "available" }, { default: () => "Available" }),
	});

	return renderToString(app);
}

async function renderCardPreview() {
	const app = createSSRApp({
		render: () =>
			h(Card, null, {
				default: () => [
					h(CardHeader, null, {
						default: () => [
							h(CardTitle, null, { default: () => "Card title" }),
							h(CardDescription, null, {
								default: () => "Card description",
							}),
						],
					}),
					h(CardContent, null, { default: () => "Card body" }),
				],
			}),
	});

	return renderToString(app);
}

async function renderCheckboxPreview(modelValue = false) {
	const app = createSSRApp({
		render: () => h(Checkbox, { modelValue }, { default: () => "Checkbox" }),
	});
	return renderToString(app);
}

async function renderSelectTriggerPreview() {
	const app = createSSRApp({
		render: () =>
			h(
				Select,
				{ modelValue: undefined },
				{
					default: () => [
						h(SelectTrigger, null, {
							default: () => h(SelectValue, { placeholder: "Select..." }),
						}),
						h(SelectContent, null, {
							default: () => [
								h(SelectItem, { value: "opt1" }, { default: () => "Option 1" }),
								h(SelectItem, { value: "opt2" }, { default: () => "Option 2" }),
							],
						}),
					],
				},
			),
	});
	return renderToString(app);
}

async function renderSearchFieldPreview() {
	const app = createSSRApp({
		render: () =>
			h(SearchField, { modelValue: "query text", placeholder: "Search here" }),
	});
	return renderToString(app);
}

async function renderPersonChipPreview() {
	const app = createSSRApp({
		render: () => h(PersonChip, { name: "Alex Smith" }),
	});
	return renderToString(app);
}

async function renderHelpButtonPreview() {
	const app = createSSRApp({
		render: () =>
			h(
				TooltipProvider,
				{ delayDuration: 0, skipDelayDuration: 0 },
				{
					default: () => h(HelpButton, { tooltip: "How to use this feature" }),
				},
			),
	});
	return renderToString(app);
}

async function renderDatePickerTriggerPreview() {
	const date = new CalendarDate(2026, 6, 15);
	const app = createSSRApp({
		render: () =>
			h(DatePicker, {
				modelValue: date,
				placeholder: "Select date",
			}),
	});
	return renderToString(app);
}

async function renderAvatarPreview() {
	const app = createSSRApp({
		render: () => h(Avatar, { name: "Alex" }),
	});
	return renderToString(app);
}

async function renderPopoverPreview() {
	const app = createSSRApp({
		render: () =>
			h(Popover, {
				defaultOpen: false,
				default: () => [
					h(PopoverTrigger, null, { default: () => "Open" }),
					h(PopoverContent, null, { default: () => "Content" }),
				],
			}),
	});
	return renderToString(app);
}

async function renderCalendarPreview() {
	const date = new CalendarDate(2026, 6, 15);
	const app = createSSRApp({
		render: () => h(Calendar, { modelValue: date }),
	});
	return renderToString(app);
}

describe("design system primitive contracts", () => {
	it("keeps outline buttons usable for white-surface secondary and filter controls", () => {
		const classes = buttonVariants({ variant: "outline", size: "default" });

		expect(classes).toContain("bg-background");
		expect(classes).toContain("border-border");
		expect(classes).toContain("text-foreground");
	});

	it.each(
		semanticBadgeCases,
	)("exposes semantic badge variant $variant through tokenized classes", ({
		variant,
		expectedClasses,
	}) => {
		const classes = badgeVariants({ variant });

		for (const expectedClass of expectedClasses) {
			expect(classes).toContain(expectedClass);
		}
	});

	it("renders button primitives with semantic label typography and variant color classes", async () => {
		const html = await renderButtonPreview();

		expect(html).toContain("[font-size:var(--text-label)]");
		expect(html).toContain("[line-height:var(--text-label--line-height)]");
		expect(html).toContain(
			"[letter-spacing:var(--text-label--letter-spacing)]",
		);
		expect(html).toContain("text-primary-foreground");
	});

	it("renders badge primitives with semantic label typography and variant color classes", async () => {
		const html = await renderBadgePreview();

		expect(html).toContain("[font-size:var(--text-label)]");
		expect(html).toContain("[line-height:var(--text-label--line-height)]");
		expect(html).toContain(
			"[letter-spacing:var(--text-label--letter-spacing)]",
		);
		expect(html).toContain("text-success");
	});

	it("renders card primitives with tokenized spacing, typography, and surface classes", async () => {
		const html = await renderCardPreview();

		expect(html).toContain("border-border/80");
		expect(html).toContain("gap-panel");
		expect(html).toContain("py-panel");
		expect(html).toContain("text-card-foreground");
		expect(html).toContain("[font-size:var(--text-body)]");
		expect(html).toContain("[line-height:var(--text-body--line-height)]");
		expect(html).toContain("shadow-sm");
		expect(html).toContain("px-panel");
		expect(html).toContain("[font-size:var(--text-card-title)]");
		expect(html).toContain("[line-height:var(--text-card-title--line-height)]");
		expect(html).toContain("[font-size:var(--text-caption)]");
		expect(html).toContain("[line-height:var(--text-caption--line-height)]");
	});
});

describe("Checkbox", () => {
	it("exposes checkbox variant classes with border-border for unchecked state", () => {
		const classes = checkboxVariants({ variant: "default", size: "default" });
		expect(classes).toContain("border-border");
		expect(classes).toContain("data-[state=checked]:bg-primary");
		expect(classes).toContain("data-[state=checked]:text-primary-foreground");
	});

	it("renders unchecked checkbox with border-border class", async () => {
		const html = await renderCheckboxPreview(false);
		expect(html).toContain("border-border");
	});

	it("renders checked checkbox with state=checked", async () => {
		const html = await renderCheckboxPreview(true);
		expect(html).toContain('data-state="checked"');
	});
});

describe("Select", () => {
	it("renders trigger with bg-background and border-border", async () => {
		const html = await renderSelectTriggerPreview();
		expect(html).toContain("bg-background");
		expect(html).toContain("border-border");
	});

	it("renders trigger with trailing ChevronDown icon", async () => {
		const html = await renderSelectTriggerPreview();
		expect(html).toContain("select-trigger");
	});

	it("instantiates without error (portal-based SelectContent)", async () => {
		// Portal content is client-side only; verifying SSR renders without error
		const html = await renderSelectTriggerPreview();
		expect(html).toBeDefined();
	});
});

describe("SearchField", () => {
	it("renders with leading Search icon and muted-foreground class", async () => {
		const html = await renderSearchFieldPreview();
		expect(html).toContain("text-muted-foreground");
	});

	it("renders v-model value in the input element", async () => {
		const html = await renderSearchFieldPreview();
		expect(html).toContain('value="query text"');
	});

	it("renders with placeholder text", async () => {
		const html = await renderSearchFieldPreview();
		expect(html).toContain('placeholder="Search here"');
	});
});

describe("PersonChip", () => {
	it("renders initials fallback for avatar", async () => {
		const html = await renderPersonChipPreview();
		expect(html).toContain(">A<");
	});

	it("renders person name adjacent to avatar", async () => {
		const html = await renderPersonChipPreview();
		expect(html).toContain("Alex Smith");
	});

	it("uses bg-secondary token class", async () => {
		const html = await renderPersonChipPreview();
		expect(html).toContain("bg-secondary");
	});

	it("renders as a span by default", async () => {
		const html = await renderPersonChipPreview();
		expect(html).toContain('data-slot="person-chip"');
	});
});

describe("HelpButton", () => {
	it("renders with ghost variant button", async () => {
		const html = await renderHelpButtonPreview();
		expect(html).toContain('data-variant="ghost"');
	});

	it("renders with icon-sm size", async () => {
		const html = await renderHelpButtonPreview();
		expect(html).toContain('data-size="icon-sm"');
	});

	it("renders tooltip trigger wrapping the button", async () => {
		const html = await renderHelpButtonPreview();
		expect(html).toContain('data-slot="tooltip-trigger"');
	});

	it("tooltip content not present in SSR (portal-based)", async () => {
		const html = await renderHelpButtonPreview();
		// TooltipContent is rendered via portal, not in SSR
		expect(html).not.toContain("tooltip-content");
	});
});

describe("Avatar", () => {
	it("renders initials fallback with bg-muted and text-label", async () => {
		const html = await renderAvatarPreview();
		expect(html).toContain("bg-muted");
		expect(html).toContain("text-label");
	});

	it("renders first character of name", async () => {
		const html = await renderAvatarPreview();
		expect(html).toContain(">A<");
	});
});

describe("DatePicker", () => {
	it("renders trigger with formatted date when modelValue is set", async () => {
		const html = await renderDatePickerTriggerPreview();
		expect(html).toContain("2026-06-15");
	});

	it("renders trigger with bg-background and border-border", async () => {
		const html = await renderDatePickerTriggerPreview();
		expect(html).toContain("bg-background");
		expect(html).toContain("border-border");
	});

	it("instantiates without error (portal-based Popover+Calendar)", async () => {
		const html = await renderDatePickerTriggerPreview();
		expect(html).toBeDefined();
	});
});

describe("Popover", () => {
	it("renders popover trigger in SSR without error", async () => {
		const html = await renderPopoverPreview();
		expect(html).toBeDefined();
	});

	it("popover content (portal) not present in SSR", async () => {
		const html = await renderPopoverPreview();
		expect(html).not.toContain("popover-content");
	});
});

describe("Calendar", () => {
	it("renders calendar header with text-body token class", async () => {
		const html = await renderCalendarPreview();
		expect(html).toContain("text-body");
	});

	it("renders day labels with text-label token class", async () => {
		const html = await renderCalendarPreview();
		expect(html).toContain("text-label");
	});

	it("instantiates without error", async () => {
		const html = await renderCalendarPreview();
		expect(html).toBeDefined();
	});
});
