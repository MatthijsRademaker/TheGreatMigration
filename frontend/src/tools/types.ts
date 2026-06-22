export interface Tool {
	id: string;
	name: string;
	/** Bringer's person ID, or null when the tool is open. */
	broughtBy: string | null;
}

export interface ToolSummary {
	total: number;
	claimed: number;
	open: number;
}
