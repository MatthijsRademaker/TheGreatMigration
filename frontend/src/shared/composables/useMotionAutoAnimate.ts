import { watch } from "vue";
import { useAutoAnimate } from "@formkit/auto-animate/vue";
import { useMotionPreference } from "./useMotionPreference";

/**
 * AutoAnimate FLIP transitions for a list container, gated by the centralized
 * motion preference. Bind the returned ref to the element whose direct children
 * are the list items (e.g. `<tbody>`, `<ul>`): added/removed/reordered children
 * animate to their new positions instead of teleporting. When the user prefers
 * reduced motion the animations are disabled and items snap into place.
 *
 * Usage: `const list = useMotionAutoAnimate()` then `<ul ref="list">`.
 */
export function useMotionAutoAnimate() {
	const { enabled } = useMotionPreference();
	const [parent, setEnabled] = useAutoAnimate();

	watch(enabled, (value) => setEnabled(value), { immediate: true });

	return parent;
}
