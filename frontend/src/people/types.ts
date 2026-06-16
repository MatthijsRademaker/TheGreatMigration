import type { BadgeVariants } from "@/shared/ui/badge";

/** Availability status values kept in sync with Badge variants. */
export type AvailabilityStatus = Extract<
	NonNullable<BadgeVariants["variant"]>,
	"available" | "busy" | "partial" | "off"
>;

export interface PersonAvailabilityEntry {
	date: string;
	status: AvailabilityStatus;
}

export interface PersonAvailability {
	id: string;
	name: string;
	availability: PersonAvailabilityEntry[];
}

export interface StatusLegendItem {
	id: AvailabilityStatus;
	label: string;
}

export interface PeopleAvailabilityProps {
	title?: string;
	description?: string;
	days?: string[];
	people?: PersonAvailability[];
	legend?: StatusLegendItem[];
	availableToday?: number;
	totalPeople?: number;
}
