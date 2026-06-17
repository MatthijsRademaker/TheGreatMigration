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

/**
 * Clamp availableToday so it never exceeds totalPeople (defensive guard against
 * stale / corrupted data). When totalPeople is zero we display a placeholder.
 */
const displayAvailableToday = computed(() => Math.min(rawAvailableToday.value, totalPeople.value))

/** Consolidated display status for the People available today card. */
const availabilityStatus = computed(() => {
  if (availabilityQuery.isPending.value) return 'loading'
  if (availabilityQuery.error.value != null) return 'error'
  if (totalPeople.value === 0) return 'empty'
  return 'ready'
})

const highPriorityTasks = computed(() => tasksBacklogQuery.data.value?.summary.highPriorityTasks ?? 0)
const unassignedTasks = computed(() => tasksBacklogQuery.data.value?.summary.unassignedTasks ?? 0)

/** Consolidated display status for the two task-backlog KPI cards. */
const backlogStatus = computed(() => {
  if (tasksBacklogQuery.isPending.value) return 'loading'
  if (tasksBacklogQuery.error.value != null) return 'error'
  return 'ready'
})

/** Shared config for the two cards that consume getTasksBacklogQuery. */
interface BacklogCardConfig {
  label: string
  value: number
  description: string
  icon: typeof TriangleAlertIcon
  bgClass: string
}

const backlogCards = computed<BacklogCardConfig[]>(() => [
  {
    label: 'High priority tasks',
    value: highPriorityTasks.value,
    description: 'Tasks marked as high priority',
    icon: TriangleAlertIcon,
    bgClass: 'bg-destructive-soft text-destructive',
  },
  {
    label: 'Unassigned jobs',
    value: unassignedTasks.value,
    description: 'Tasks with no one assigned yet',
    icon: HammerIcon,
    bgClass: 'bg-warning-soft text-warning',
  },
])
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
    <!-- People available today -->
    <Card class="relative">
      <CardHeader class="flex flex-row items-start justify-between gap-3">
        <div class="flex flex-col gap-1">
          <CardDescription>People available today</CardDescription>
          <CardTitle class="text-3xl">
            <span v-if="availabilityStatus === 'loading'" class="text-muted-foreground">Loading…</span>
            <span v-else-if="availabilityStatus === 'error'" class="text-destructive">Backend unavailable</span>
            <span v-else-if="availabilityStatus === 'empty'">—</span>
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
      <img src="/images/leaf.png" alt="" class="absolute bottom-0 right-0 h-14 w-auto pointer-events-none" />
    </Card>

    <!-- Task-backlog cards (shared loading / error / value layout) -->
    <Card v-for="card in backlogCards" :key="card.label" class="relative">
      <CardHeader class="flex flex-row items-start justify-between gap-3">
        <div class="flex flex-col gap-1">
          <CardDescription>{{ card.label }}</CardDescription>
          <CardTitle class="text-3xl">
            <span v-if="backlogStatus === 'loading'" class="text-muted-foreground">Loading…</span>
            <span v-else-if="backlogStatus === 'error'" class="text-destructive">Backend unavailable</span>
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
      <img src="/images/leaf.png" alt="" class="absolute bottom-0 right-0 h-14 w-auto pointer-events-none" />
    </Card>

    <!-- Rooms completed (placeholder for future room-progress contract) -->
    <Card data-testid="kpi-placeholder-rooms-completed" class="relative">
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
      <img src="/images/leaf.png" alt="" class="absolute bottom-0 right-0 h-14 w-auto pointer-events-none" />
    </Card>
  </div>
</template>
