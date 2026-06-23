import { computed } from "vue";
import { usePreferredReducedMotion } from "@vueuse/core";
import { INSTANT, countUp, durations, springs } from "@/shared/motion/tokens";

/**
 * Centralized reduced-motion gate. Every animated surface SHALL consume this
 * composable instead of checking `prefers-reduced-motion` independently, so the
 * playful motion layer is consistent and centrally controllable.
 *
 * When the user prefers reduced motion, decorative motion (count-ups, springs,
 * confetti, route slides) is skipped while the underlying state change still
 * applies instantly. The boolean is reactive: toggling the OS/browser setting
 * updates every surface without a reload.
 */
export function useMotionPreference() {
	const preference = usePreferredReducedMotion();

	/** True when non-essential (decorative) motion should play. */
	const enabled = computed(() => preference.value !== "reduce");

	/**
	 * Resolve a `motion-v` transition config: returns it as-is when motion is
	 * enabled, or an instant (duration 0) transition when reduced — so the
	 * element snaps to its final state instead of animating.
	 */
	function transition<T extends object>(config: T): T | typeof INSTANT {
		return enabled.value ? config : INSTANT;
	}

	/** Count-up duration in ms, collapsed to 0 under reduced motion. */
	const countUpDuration = computed(() => (enabled.value ? countUp.duration : 0));

	return {
		enabled,
		transition,
		countUpDuration,
		/** Raw tokens, re-exported for ergonomic access at call sites. */
		springs,
		durations,
		countUp,
	};
}
