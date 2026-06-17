<script setup lang="ts">
import { useQuery } from '@pinia/colada'
import { computed } from 'vue'
import type { Component } from 'vue'
import {
  Building2Icon,
  HammerIcon,
  TriangleAlertIcon,
  UsersRoundIcon,
} from '@lucide/vue'
import { getDashboardPeopleAvailabilityQuery, getTasksBacklogQuery } from '@/client/@pinia/colada.gen'
import { Card } from '@/shared/ui/card'

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
const peopleStatus = computed(() => {
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

interface KpiCardConfig {
  id: string
  label: string
  description: string
  icon: Component
  accentClass: string
  status: 'loading' | 'error' | 'ready' | 'empty'
  value: number
}

/** Unified config driving all four KPI cards from one computed. */
const cardConfigs = computed<KpiCardConfig[]>(() => [
  {
    id: 'people',
    label: 'People available today',
    description: 'Helpers with confirmed availability',
    icon: UsersRoundIcon,
    accentClass: 'bg-info-soft text-info',
    status: peopleStatus.value,
    value: displayAvailableToday.value,
  },
  {
    id: 'high-priority',
    label: 'High priority tasks',
    description: 'Tasks marked as high priority',
    icon: TriangleAlertIcon,
    accentClass: 'bg-destructive-soft text-destructive',
    status: backlogStatus.value,
    value: highPriorityTasks.value,
  },
  {
    id: 'unassigned',
    label: 'Unassigned jobs',
    description: 'Tasks with no one assigned yet',
    icon: HammerIcon,
    accentClass: 'bg-warning-soft text-warning',
    status: backlogStatus.value,
    value: unassignedTasks.value,
  },
  // Rooms completed — placeholder for future room-progress contract
  {
    id: 'rooms',
    label: 'Rooms completed',
    description: 'Rooms fully packed and cleared',
    icon: Building2Icon,
    accentClass: 'bg-success-soft text-success',
    status: 'empty',
    value: 0,
  },
])
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
    <Card
      v-for="card in cardConfigs"
      :key="card.id"
      class="!py-0 !gap-0 relative"
      :data-testid="card.id === 'rooms' ? 'kpi-placeholder-rooms-completed' : undefined"
    >
      <!-- !py-0 !gap-0 overrides Card.vue&#39;s default padding so the left accent column is flush with card edges -->
      <div class="flex flex-row">
        <!-- Left accent column: full-height semantic color band -->
        <div
          class="flex w-[72px] shrink-0 items-center justify-center self-stretch rounded-l-lg"
          :class="card.accentClass"
        >
          <component :is="card.icon" class="size-8" />
        </div>
        <!-- Right content column: title, KPI value, subtitle stacked -->
        <div class="flex flex-1 flex-col gap-1 p-panel">
          <span class="[font-size:var(--text-caption)] [line-height:var(--text-caption--line-height)] text-muted-foreground">{{ card.label }}</span>
          <span
            class="text-3xl font-semibold"
            :class="{
              'text-muted-foreground': card.status === 'loading',
              'text-destructive': card.status === 'error',
            }"
          >
            <template v-if="card.status === 'loading'">Loading…</template>
            <template v-else-if="card.status === 'error'">Backend unavailable</template>
            <template v-else-if="card.status === 'empty'">—</template>
            <template v-else-if="card.id === 'people' && card.status === 'ready'">
              {{ card.value }}<span class="text-xl text-muted-foreground"> of {{ totalPeople }}</span>&nbsp;<span class="text-xl text-muted-foreground">available</span>
            </template>
            <template v-else>{{ card.value }}</template>
          </span>
          <span class="text-sm text-muted-foreground">{{ card.description }}</span>
        </div>
      </div>
      <img src="/images/leaf.png" alt="" class="absolute bottom-0 right-0 h-14 w-auto pointer-events-none" />
    </Card>
  </div>
</template>
