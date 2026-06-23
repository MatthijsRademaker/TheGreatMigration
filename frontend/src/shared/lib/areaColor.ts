/**
 * Deterministic color for a room/area, derived by hashing the area id into a
 * fixed palette. The same id always yields the same color, and the color does
 * not depend on the area name — renaming an area preserves its color.
 *
 * Color is a purely presentational concern: it is computed here in the frontend
 * and never persisted or returned by the API.
 */

// Fixed palette of accent colors that read well as small chips/dots in both
// light and dark themes.
const PALETTE = [
	"#3b82f6", // blue
	"#f97316", // orange
	"#10b981", // emerald
	"#a855f7", // purple
	"#ef4444", // red
	"#eab308", // yellow
	"#06b6d4", // cyan
	"#ec4899", // pink
];

/** Returns a stable palette color for the given area id (djb2 hash). */
export function areaColor(id: string): string {
	let hash = 5381;
	for (let i = 0; i < id.length; i++) {
		hash = ((hash << 5) + hash) ^ id.charCodeAt(i);
	}
	return PALETTE[Math.abs(hash) % PALETTE.length];
}
