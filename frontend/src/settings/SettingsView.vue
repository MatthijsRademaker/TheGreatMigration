<script setup lang="ts">
import { computed, shallowRef, ref, watchEffect } from "vue";
import { CalendarDate, type DateValue } from "@internationalized/date";
import { useQuery } from "@pinia/colada";
import { getPlanningWindowQuery } from "@/client/@pinia/colada.gen";
import { useUpdatePlanningWindow } from "@/shared/composables/useUpdatePlanningWindow";
import { DatePicker } from "@/shared/ui/date-picker";
import { Button } from "@/shared/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/shared/ui/card";

const planningQuery = useQuery({
	...getPlanningWindowQuery(),
	ssrCatchError: true,
});
const { mutate, isPending, error: mutationError } = useUpdatePlanningWindow();

const isLoading = computed(() => planningQuery.isPending.value);
const isError = computed(() => !!planningQuery.error.value);

// Local editable state: CalendarDate | null.
const startDate = shallowRef<DateValue | null>(null);
const endDate = shallowRef<DateValue | null>(null);

// Original loaded values as ISO strings.
const originalStart = ref<string>("");
const originalEnd = ref<string>("");

// Populate local state when planning window data loads.
watchEffect(() => {
	const data = planningQuery.data.value;
	if (data && data.startDate && data.endDate) {
		const s = new CalendarDate(
			Number(data.startDate.slice(0, 4)),
			Number(data.startDate.slice(5, 7)),
			Number(data.startDate.slice(8, 10)),
		);
		const e = new CalendarDate(
			Number(data.endDate.slice(0, 4)),
			Number(data.endDate.slice(5, 7)),
			Number(data.endDate.slice(8, 10)),
		);
		startDate.value = s;
		endDate.value = e;
		originalStart.value = data.startDate;
		originalEnd.value = data.endDate;
	}
});

// Derived validation.
const isDirty = computed(() => {
	const sStr = startDate.value?.toString() ?? "";
	const eStr = endDate.value?.toString() ?? "";
	return sStr !== originalStart.value || eStr !== originalEnd.value;
});

const isEndBeforeStart = computed(() => {
	if (!startDate.value || !endDate.value) return false;
	return endDate.value.compare(startDate.value) < 0;
});

const canSave = computed(() => {
	if (isLoading.value || isPending.value) return false;
	if (!startDate.value || !endDate.value) return false;
	if (isEndBeforeStart.value) return false;
	if (!isDirty.value) return false;
	return true;
});

function handleSave() {
	if (!canSave.value) return;
	mutate({
		body: {
			startDate: startDate.value!.toString(),
			endDate: endDate.value!.toString(),
		},
	});
}

function handleReset() {
	if (originalStart.value && originalEnd.value) {
		const s = new CalendarDate(
			Number(originalStart.value.slice(0, 4)),
			Number(originalStart.value.slice(5, 7)),
			Number(originalStart.value.slice(8, 10)),
		);
		const e = new CalendarDate(
			Number(originalEnd.value.slice(0, 4)),
			Number(originalEnd.value.slice(5, 7)),
			Number(originalEnd.value.slice(8, 10)),
		);
		startDate.value = s;
		endDate.value = e;
	}
}
</script>

<template>
	<section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
		<Card>
			<CardHeader>
				<CardTitle>Planning window</CardTitle>
				<CardDescription>
					Set the start and end dates for your move timeline. All schedule and
					availability views adapt to this range. The seeded demo data covers
					2026-07-05 through 2026-08-13.
				</CardDescription>
			</CardHeader>
			<CardContent class="flex flex-col gap-4">
				<div v-if="isLoading" class="text-sm text-muted-foreground">
					Loading planning window&hellip;
				</div>

				<div
					v-else-if="isError"
					class="text-sm text-destructive"
				>
					Failed to load planning window. Please try again later.
				</div>

				<template v-else>
					<div class="flex flex-col gap-2">
						<label class="text-sm font-medium">Start date</label>
						<DatePicker v-model="startDate" placeholder="Select start date" />
					</div>
					<div class="flex flex-col gap-2">
						<label class="text-sm font-medium">End date</label>
						<DatePicker v-model="endDate" placeholder="Select end date" />
					</div>
					<div
						v-if="isEndBeforeStart"
						class="text-sm text-destructive"
					>
						End date must be on or after the start date.
					</div>
					<div
						v-if="mutationError"
						class="text-sm text-destructive"
					>
						Failed to save changes. Please try again.
					</div>
				</template>
			</CardContent>
			<CardFooter class="gap-2">
				<Button
					variant="outline"
					:disabled="!isDirty || isPending"
					@click="handleReset"
				>
					Reset
				</Button>
				<Button :disabled="!canSave" :data-loading="isPending" @click="handleSave">
					{{ isPending ? "Saving&hellip;" : "Save" }}
				</Button>
			</CardFooter>
		</Card>
	</section>
</template>
