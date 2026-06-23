/**
 * Shared motion tokens — the single source of truth for the app's playful
 * motion language. Spring physics and durations are consumed by `motion-v`
 * configs and `@vueuse/core`'s `useTransition`; the matching CSS custom
 * properties (durations/easings, the shimmer keyframe) live in
 * `src/app/styles.css`. Tune the feel here and every animated surface follows.
 */

/** `motion-v` spring transition presets. */
export const springs = {
	/** Snappy, minimal overshoot — hover/press and most UI feedback. */
	default: { type: "spring", stiffness: 400, damping: 30 } as const,
	/** Lively overshoot — count-ups, reward pops, celebratory moments. */
	bouncy: { type: "spring", stiffness: 520, damping: 17 } as const,
	/** Smooth, no overshoot — layout/position transitions (FLIP, reschedule). */
	smooth: { type: "spring", stiffness: 300, damping: 34 } as const,
} as const;

/** Durations in seconds (for `motion-v`); multiply by 1000 for ms APIs. */
export const durations = {
	fast: 0.15,
	base: 0.25,
	slow: 0.4,
} as const;

/** Cubic-bézier easing curves. */
export const easings = {
	/** Material standard ease. */
	standard: [0.4, 0, 0.2, 1] as const,
	/** Overshoots then settles — count-up "spring" feel for `useTransition`. */
	overshoot: [0.34, 1.56, 0.64, 1] as const,
} as const;

/**
 * Count-up settings for `@vueuse/core` `useTransition` (KPI numbers).
 * `duration` is milliseconds; `transition` is the easing curve.
 */
export const countUp = {
	duration: 650,
	transition: easings.overshoot,
} as const;

/** `motion-v` transition that applies the final state with no animation. */
export const INSTANT = { duration: 0 } as const;
