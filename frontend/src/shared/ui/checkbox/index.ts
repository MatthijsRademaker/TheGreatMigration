import type { VariantProps } from "class-variance-authority";
import { cva } from "class-variance-authority";

export { default as Checkbox } from "./Checkbox.vue";

export const checkboxVariants = cva(
	"peer shrink-0 rounded-sm border border-border ring-offset-background focus-visible:border-ring focus-visible:ring-ring/30 focus-visible:ring-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground",
	{
		variants: {
			size: {
				default: "size-5",
				sm: "size-4",
				lg: "size-6",
			},
		},
		defaultVariants: {
			size: "default",
		},
	},
);
export type CheckboxVariants = VariantProps<typeof checkboxVariants>;
