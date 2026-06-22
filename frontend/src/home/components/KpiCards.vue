<script setup lang="ts">
import { useQuery } from '@pinia/colada'
import { computed } from 'vue'
import type { Component } from 'vue'
import {
  BriefcaseIcon,
  CheckCircleIcon,
  FlagIcon,
  UsersRoundIcon,
  WrenchIcon,
} from '@lucide/vue'
import { getDashboardPeopleAvailabilityQuery, getTasksBacklogQuery, getToolsQuery } from '@/client/@pinia/colada.gen'
import { Card, CardHeader, CardContent } from '@/shared/ui/card'

const availabilityQuery = useQuery(getDashboardPeopleAvailabilityQuery())
const tasksBacklogQuery = useQuery(getTasksBacklogQuery())
const toolsQuery = useQuery(getToolsQuery())

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

const toolsClaimed = computed(() => toolsQuery.data.value?.summary.claimed ?? 0)
const toolsTotal = computed(() => toolsQuery.data.value?.summary.total ?? 0)

/** Consolidated display status for the Tools covered KPI card. */
const toolsStatus = computed(() => {
  if (toolsQuery.isPending.value) return 'loading'
  if (toolsQuery.error.value != null) return 'error'
  return 'ready'
})

/** Destructure range from the availability response to access selectedDate. */
const range = computed(() => availabilityQuery.data.value?.range)

/** Format range.selectedDate as "MMM D" (e.g. "Jul 5"). Returns undefined if missing. */
const formattedSelectedDate = computed(() => {
  const dateStr = range.value?.selectedDate
  if (!dateStr) return undefined
  const [year, month, day] = dateStr.split('-').map(Number)
  const date = new Date(Date.UTC(year, month - 1, day))
  return new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', timeZone: 'UTC' }).format(date)
})

interface KpiCardConfig {
  id: string
  label: string
  description: string
  icon: Component
  iconBgClass: string
  borderClass: string
  status: 'loading' | 'error' | 'ready' | 'empty'
  value: number
}

/**
 * Unified config driving all five KPI cards. Card order follows the design:
 * High priority tasks → People available today → Unassigned jobs → Rooms completed → Tools covered.
 */
const cardConfigs = computed<KpiCardConfig[]>(() => [
  {
    id: 'high-priority',
    label: 'High priority tasks',
    description: 'high priority tasks need attention',
    icon: FlagIcon,
    iconBgClass: 'bg-destructive-soft text-destructive',
    borderClass: 'border-destructive',
    status: backlogStatus.value,
    value: highPriorityTasks.value,
  },
  {
    id: 'people',
    label: 'People available today',
    description: formattedSelectedDate.value
      ? `available on ${formattedSelectedDate.value}`
      : 'available today',
    icon: UsersRoundIcon,
    iconBgClass: 'bg-info-soft text-info',
    borderClass: 'border-info',
    status: peopleStatus.value,
    value: displayAvailableToday.value,
  },
  {
    id: 'unassigned',
    label: 'Unassigned jobs',
    description: 'jobs that need assignment',
    icon: BriefcaseIcon,
    iconBgClass: 'bg-warning-soft text-warning',
    borderClass: 'border-warning',
    status: backlogStatus.value,
    value: unassignedTasks.value,
  },
  // Rooms completed — placeholder for future room-progress contract
  {
    id: 'rooms',
    label: 'Rooms completed',
    description: 'rooms fully packed and cleared',
    icon: CheckCircleIcon,
    iconBgClass: 'bg-success-soft text-success',
    borderClass: 'border-success',
    status: 'empty',
    value: 0,
  },
  {
    id: 'tools',
    label: 'Tools covered',
    description: 'tools claimed by a helper',
    icon: WrenchIcon,
    iconBgClass: 'bg-info-soft text-info',
    borderClass: 'border-info',
    status: toolsStatus.value,
    value: toolsClaimed.value,
  },
])
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
    <Card
      v-for="card in cardConfigs"
      :key="card.id"
      class="border-l-4 relative"
      :class="card.borderClass"
      :data-testid="card.id === 'rooms' ? 'kpi-placeholder-rooms-completed' : undefined"
    >
      <CardHeader class="!pb-0">
        <div class="flex items-center gap-2">
          <div
            class="size-8 rounded-lg flex items-center justify-center shrink-0"
            :class="card.iconBgClass"
          >
            <component :is="card.icon" class="size-4" />
          </div>
          <span class="[font-size:var(--text-caption)] [line-height:var(--text-caption--line-height)] text-muted-foreground">
            {{ card.label }}
          </span>
        </div>
      </CardHeader>
      <CardContent class="!pt-0">
        <div class="flex flex-col gap-1">
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
              {{ card.value }} / {{ totalPeople }}
            </template>
            <template v-else-if="card.id === 'tools' && card.status === 'ready'">
              {{ toolsClaimed }} / {{ toolsTotal }}
            </template>
            <template v-else>{{ card.value }}</template>
          </span>
          <span class="text-sm text-muted-foreground">{{ card.description }}</span>
        </div>
      </CardContent>
      <img src="/images/leaf.png" alt="" class="absolute bottom-0 right-0 h-14 w-auto pointer-events-none" />
    </Card>
  </div>
</template>
