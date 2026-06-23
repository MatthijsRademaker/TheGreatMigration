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

export interface CellChangePayload {
	personId: string;
	dayIndex: number;
	status: AvailabilityStatus | null;
}

export interface PeopleAvailabilityProps {
	title?: string;
	days?: string[];
	people?: PersonAvailability[];
	availableToday?: number;
	totalPeople?: number;
	editable?: boolean;
	/** When non-null, the Delete button for this person is disabled and shows "Deleting…". */
	deletingPersonId?: string | null;
	/** When true, status popover triggers are disabled to prevent concurrent edits. */
	updating?: boolean;
}
