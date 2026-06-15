import type { VariantProps } from "class-variance-authority";
import { cva } from "class-variance-authority";

export { default as Badge } from "./Badge.vue";

export const badgeVariants = cva(
	"h-5 gap-1 rounded-full border border-transparent px-2.5 py-0.5 text-label font-medium transition-all has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 [&>svg]:size-2.5! group/badge inline-flex w-fit shrink-0 items-center justify-center overflow-hidden whitespace-nowrap focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none",
	{
		variants: {
			variant: {
				default: "bg-primary text-primary-foreground [a]:hover:bg-primary/90",
				secondary: "bg-secondary text-secondary-foreground [a]:hover:bg-muted",
				destructive:
					"border-destructive/10 bg-destructive-soft text-destructive [a]:hover:bg-destructive-soft/80 focus-visible:ring-destructive/20",
				outline:
					"border-border bg-background text-foreground [a]:hover:bg-muted [a]:hover:text-foreground",
				ghost: "hover:bg-muted hover:text-foreground",
				link: "text-primary underline-offset-4 hover:underline",
				priorityHigh:
					"border-destructive/10 bg-destructive-soft text-destructive",
				priorityMedium: "border-warning/20 bg-warning-soft text-warning",
				priorityLow: "border-success/20 bg-success-soft text-success",
				available: "border-success/20 bg-success-soft text-success",
				busy: "border-destructive/10 bg-destructive-soft text-destructive",
				partial: "border-warning/20 bg-warning-soft text-warning",
				off: "border-border bg-secondary text-muted-foreground",
			},
		},
		defaultVariants: {
			variant: "default",
		},
	},
);
export type BadgeVariants = VariantProps<typeof badgeVariants>;
