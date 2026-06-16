<script setup lang="ts">
import { useQuery } from '@pinia/colada'
import { computed } from 'vue'
import {
  Building2Icon,
  HammerIcon,
  TriangleAlertIcon,
  UsersRoundIcon,
} from '@lucide/vue'
import { getDashboardPeopleAvailabilityQuery, getTasksBacklogQuery } from '@/client/@pinia/colada.gen'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

const availabilityQuery = useQuery(getDashboardPeopleAvailabilityQuery())
const tasksBacklogQuery = useQuery(getTasksBacklogQuery())

const rawAvailableToday = computed(() => availabilityQuery.data.value?.summary.availableToday ?? 0)
const totalPeople = computed(() => availabilityQuery.data.value?.summary.totalPeople ?? 0)
const availabilityLoading = computed(() => availabilityQuery.isPending.value)
const availabilityError = computed(() => availabilityQuery.error.value != null)

/**
 * Clamp availableToday so it never exceeds totalPeople (defensive guard against
 * stale / corrupted data). When totalPeople is zero we display a placeholder.
 */
const displayAvailableToday = computed(() => Math.min(rawAvailableToday.value, totalPeople.value))

const highPriorityTasks = computed(() => tasksBacklogQuery.data.value?.summary.highPriorityTasks ?? 0)
const unassignedTasks = computed(() => tasksBacklogQuery.data.value?.summary.unassignedTasks ?? 0)
const backlogLoading = computed(() => tasksBacklogQuery.isPending.value)
const backlogError = computed(() => tasksBacklogQuery.error.value != null)

/** Shared config for the two cards that consume getTasksBacklogQuery. */
interface BacklogCardConfig {
  label: string
  value: number
  description: string
  icon: typeof TriangleAlertIcon
  bgClass: string
  textClass: string
}

const backlogCards = computed<BacklogCardConfig[]>(() => [
  {
    label: 'High priority tasks',
    value: highPriorityTasks.value,
    description: 'Tasks marked as high priority',
    icon: TriangleAlertIcon,
    bgClass: 'bg-destructive-soft text-destructive',
    textClass: '',
  },
  {
    label: 'Unassigned jobs',
    value: unassignedTasks.value,
    description: 'Tasks with no one assigned yet',
    icon: HammerIcon,
    bgClass: 'bg-warning-soft text-warning',
    textClass: '',
  },
])
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
    <!-- People available today -->
    <Card>
      <CardHeader class="flex flex-row items-start justify-between gap-3">
        <div class="flex flex-col gap-1">
          <CardDescription>People available today</CardDescription>
          <CardTitle class="text-3xl">
            <span v-if="availabilityLoading" class="text-muted-foreground">Loading…</span>
            <span v-else-if="availabilityError" class="text-destructive">Backend unavailable</span>
            <span v-else-if="totalPeople === 0">—</span>
            <span v-else>{{ displayAvailableToday }}<span class="text-xl text-muted-foreground"> of {{ totalPeople }}</span>&nbsp;<span class="text-xl text-muted-foreground">available</span></span>
          </CardTitle>
        </div>
        <div class="flex size-10 items-center justify-center rounded-full bg-secondary text-secondary-foreground">
          <UsersRoundIcon />
        </div>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">Helpers with confirmed availability</p>
      </CardContent>
    </Card>

    <!-- Task-backlog cards (shared loading / error / value layout) -->
    <Card v-for="card in backlogCards" :key="card.label">
      <CardHeader class="flex flex-row items-start justify-between gap-3">
        <div class="flex flex-col gap-1">
          <CardDescription>{{ card.label }}</CardDescription>
          <CardTitle class="text-3xl">
            <span v-if="backlogLoading" class="text-muted-foreground">Loading…</span>
            <span v-else-if="backlogError" class="text-destructive">Backend unavailable</span>
            <span v-else>{{ card.value }}</span>
          </CardTitle>
        </div>
        <div :class="['flex size-10 items-center justify-center rounded-full', card.bgClass]">
          <component :is="card.icon" />
        </div>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">{{ card.description }}</p>
      </CardContent>
    </Card>

    <!-- Rooms completed (placeholder for future room-progress contract) -->
    <Card data-testid="kpi-placeholder-rooms-completed">
      <CardHeader class="flex flex-row items-start justify-between gap-3">
        <div class="flex flex-col gap-1">
          <CardDescription>Rooms completed</CardDescription>
          <CardTitle class="text-3xl">—</CardTitle>
        </div>
        <div class="flex size-10 items-center justify-center rounded-full bg-secondary text-secondary-foreground">
          <Building2Icon />
        </div>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">Rooms fully packed and cleared</p>
      </CardContent>
    </Card>
  </div>
</template>
