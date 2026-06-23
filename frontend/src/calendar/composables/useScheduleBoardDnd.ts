import { computed, ref, watch, type ComputedRef } from "vue";
import { useMutation, useQueryCache } from "@pinia/colada";
import { updateScheduleCardMutation } from "@/client/@pinia/colada.gen";
import type {
	CreateScheduleCardRequestBodyWritable,
	TaskCard as ApiTaskCard,
} from "@/client/types.gen";
import {
	sortScheduleTasksByPriority,
	type ScheduleDay,
	type ScheduleTaskCard,
} from "./useDailySchedule";

/** Derive uppercase initials from a name as a client-side optimistic stand-in. */
function deriveInitials(name: string): string {
	return name
		.split(/\s+/)
		.filter(Boolean)
		.slice(0, 2)
		.map((part) => part[0]?.toUpperCase() ?? "")
		.join("");
}

/**
 * Deep clone the day list so optimistic edits never mutate the query data.
 * Uses JSON round-tripping: the day/card data is plain and serializable, and
 * this avoids `structuredClone` throwing on Vue's reactive proxies.
 */
function cloneDays(days: ScheduleDay[]): ScheduleDay[] {
	return JSON.parse(JSON.stringify(days)) as ScheduleDay[];
}

/** Build the full, valid update body the API requires (taskId-linked or free-form). */
function buildUpdateBody(
	card: ScheduleTaskCard,
	scheduledDate: string,
): CreateScheduleCardRequestBodyWritable {
	const assignedTo = card.assignedPeople.map((p) => p.id);
	if (card.taskId) {
		return { taskId: card.taskId, scheduledDate, assignedTo };
	}
	return {
		title: card.title,
		priority: card.priority,
		roomArea: card.roomArea,
		peopleNeeded: card.peopleNeeded,
		scheduledDate,
		assignedTo,
	};
}

interface CardLocation {
	dayIndex: number;
	cardIndex: number;
	card: ScheduleTaskCard;
}

function locateCard(days: ScheduleDay[], cardId: string): CardLocation | null {
	for (let d = 0; d < days.length; d++) {
		const cardIndex = days[d].tasks.findIndex((t) => t.id === cardId);
		if (cardIndex !== -1) {
			return { dayIndex: d, cardIndex, card: days[d].tasks[cardIndex] };
		}
	}
	return null;
}

/** Apply the server-returned card onto the optimistic copy (reconcile). */
function reconcileCard(card: ScheduleTaskCard, updated: ApiTaskCard): void {
	card.assignedPeople = updated.assignedPeople ?? [];
	card.assignedCount = updated.assignedCount;
	card.staffingStatus =
		updated.staffingStatus as ScheduleTaskCard["staffingStatus"];
	card.peopleNeeded = updated.peopleNeeded;
}

export interface UseScheduleBoardDndOptions {
	/** Adapted day list from the schedule query (source of truth). */
	source: ComputedRef<ScheduleDay[]>;
	/** Schedule query key for background invalidation after a successful write. */
	queryKey: unknown;
}

/**
 * Drives optimistic drag-and-drop assignment and rescheduling on the schedule
 * board. While a write is in flight an optimistic override of the day list is
 * shown; the change is persisted via `updateScheduleCard`, reconciled with the
 * server response, and the query is invalidated so the board converges with the
 * server. On error the override is dropped, reverting to the untouched query
 * data, and a non-blocking error is surfaced — the board is never left in a
 * fabricated state.
 *
 * `board` is a computed that falls back to the query data whenever no optimistic
 * write is pending, which keeps it reactive (and correct under SSR).
 */
export function useScheduleBoardDnd(options: UseScheduleBoardDndOptions) {
	const { source, queryKey } = options;
	const queryCache = useQueryCache();
	const updateMut = useMutation(updateScheduleCardMutation());

	/** Optimistic override; null means "follow the query". */
	const override = ref<ScheduleDay[] | null>(null);
	/** Non-blocking error message for the most recent failed mutation. */
	const error = ref<string | null>(null);
	/** In-flight optimistic writes; while > 0 the override owns the board. */
	const inFlight = ref(0);

	/** What the board renders: the optimistic override, else the live query. */
	const board = computed<ScheduleDay[]>(() => override.value ?? source.value);

	// Once the query catches up (e.g. after invalidation) and nothing is in
	// flight, drop the override so the board follows server truth again.
	watch(source, () => {
		if (inFlight.value === 0) override.value = null;
	});

	async function persist(
		cardId: string,
		card: ScheduleTaskCard,
		scheduledDate: string,
		next: ScheduleDay[],
		fallbackMessage: string,
	): Promise<void> {
		override.value = next;
		error.value = null;
		inFlight.value += 1;
		try {
			const updated = await updateMut.mutateAsync({
				path: { id: cardId },
				body: buildUpdateBody(card, scheduledDate),
			});
			const after = locateCard(override.value ?? [], cardId);
			if (after) reconcileCard(after.card, updated);
			queryCache.invalidateQueries({ key: queryKey as never });
		} catch (e) {
			override.value = null; // revert to untouched query data
			error.value = e instanceof Error ? e.message : fallbackMessage;
		} finally {
			inFlight.value = Math.max(0, inFlight.value - 1);
		}
	}

	/**
	 * Assign a person to a card by appending their ID and persisting. Duplicate
	 * assignment is a no-op (no duplicate entry, no redundant mutation).
	 */
	async function assignPerson(
		cardId: string,
		person: { id: string; name: string },
	): Promise<void> {
		const loc = locateCard(board.value, cardId);
		if (!loc) return;
		if (loc.card.assignedPeople.some((p) => p.id === person.id)) {
			return; // duplicate — no-op
		}

		const next = cloneDays(board.value);
		const target = locateCard(next, cardId);
		if (!target) return;
		target.card.assignedPeople.push({
			id: person.id,
			name: person.name,
			initials: deriveInitials(person.name),
		});
		target.card.assignedCount = target.card.assignedPeople.length;

		await persist(
			cardId,
			target.card,
			target.card.scheduledDate,
			next,
			"Failed to assign helper",
		);
	}

	/**
	 * Reschedule a card to a different day by moving it between columns and
	 * persisting the new `scheduledDate`. Reverts to the original column on error.
	 */
	async function rescheduleCard(
		cardId: string,
		targetDate: string,
	): Promise<void> {
		const loc = locateCard(board.value, cardId);
		if (!loc) return;
		if (loc.card.scheduledDate === targetDate) return; // same column — no-op

		const next = cloneDays(board.value);
		const targetDayIndex = next.findIndex((d) => d.date === targetDate);
		const source = locateCard(next, cardId);
		if (targetDayIndex === -1 || !source) return;

		const [moved] = next[source.dayIndex].tasks.splice(source.cardIndex, 1);
		moved.scheduledDate = targetDate;
		next[targetDayIndex].tasks.push(moved);
		sortScheduleTasksByPriority(next[targetDayIndex].tasks);

		await persist(cardId, moved, targetDate, next, "Failed to reschedule card");
	}

	function dismissError(): void {
		error.value = null;
	}

	return { board, error, assignPerson, rescheduleCard, dismissError };
}
