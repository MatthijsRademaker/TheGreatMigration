import type { VariantProps } from "class-variance-authority";
import { cva } from "class-variance-authority";

export { default as Button } from "./Button.vue";

export const buttonVariants = cva(
	"focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 rounded-md border border-transparent bg-clip-padding text-label font-medium focus-visible:ring-2 aria-invalid:ring-2 active:not-aria-[haspopup]:translate-y-px [&_svg:not([class*=size-])]:size-4 group/button inline-flex shrink-0 items-center justify-center whitespace-nowrap transition-all outline-none select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
	{
		variants: {
			variant: {
				default:
					"bg-primary text-primary-foreground shadow-sm hover:bg-primary/90",
				outline:
					"border-border bg-background text-foreground hover:bg-muted/70 hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground",
				secondary:
					"border-border bg-secondary text-foreground hover:bg-muted aria-expanded:bg-muted aria-expanded:text-foreground",
				ghost:
					"text-primary hover:bg-muted/80 hover:text-primary aria-expanded:bg-muted aria-expanded:text-primary",
				destructive:
					"border-destructive/10 bg-destructive-soft text-destructive hover:bg-destructive-soft/80 focus-visible:ring-destructive/20 focus-visible:border-destructive/30",
				link: "text-primary underline-offset-4 hover:underline",
			},
			size: {
				default:
					"h-8 gap-1.5 px-3 has-data-[icon=inline-end]:pr-2.5 has-data-[icon=inline-start]:pl-2.5 [&_svg:not([class*=size-])]:size-3.5",
				xs: "h-6 gap-1 rounded-sm px-2 has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*=size-])]:size-2.5",
				sm: "h-7 gap-1 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 [&_svg:not([class*=size-])]:size-3",
				lg: "h-9 gap-1.5 px-4 has-data-[icon=inline-end]:pr-3 has-data-[icon=inline-start]:pl-3 [&_svg:not([class*=size-])]:size-4",
				icon: "size-8 [&_svg:not([class*=size-])]:size-3.5",
				"icon-xs": "size-6 rounded-sm [&_svg:not([class*=size-])]:size-2.5",
				"icon-sm": "size-7 [&_svg:not([class*=size-])]:size-3",
				"icon-lg": "size-9 [&_svg:not([class*=size-])]:size-4",
			},
		},
		defaultVariants: {
			variant: "default",
			size: "default",
		},
	},
);
export type ButtonVariants = VariantProps<typeof buttonVariants>;
