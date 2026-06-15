import { describe, expect, it } from "vitest";
import { badgeVariants } from "../src/shared/ui/badge";
import { buttonVariants } from "../src/shared/ui/button";

describe("design system primitive contracts", () => {
	it("keeps outline buttons usable for white-surface secondary and filter controls", () => {
		const classes = buttonVariants({ variant: "outline", size: "default" });

		expect(classes).toContain("bg-background");
		expect(classes).toContain("border-border");
		expect(classes).toContain("text-foreground");
	});

	it("exposes semantic badge variants for priority and availability chips", () => {
		const highPriorityClasses = badgeVariants({
			variant: "priorityHigh" as never,
		});
		const availableClasses = badgeVariants({ variant: "available" as never });

		expect(highPriorityClasses).toContain("bg-destructive-soft");
		expect(highPriorityClasses).toContain("text-destructive");
		expect(availableClasses).toContain("bg-success-soft");
		expect(availableClasses).toContain("text-success");
	});
});
