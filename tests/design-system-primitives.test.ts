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
