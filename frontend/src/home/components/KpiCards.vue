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

const availableToday = computed(() => availabilityQuery.data.value?.summary.availableToday ?? 0)
const totalPeople = computed(() => availabilityQuery.data.value?.summary.totalPeople ?? 0)
const availabilityLoading = computed(() => availabilityQuery.isPending.value)
const availabilityError = computed(() => availabilityQuery.error.value != null)

const highPriorityTasks = computed(() => tasksBacklogQuery.data.value?.summary.highPriorityTasks ?? 0)
const unassignedTasks = computed(() => tasksBacklogQuery.data.value?.summary.unassignedTasks ?? 0)
const backlogLoading = computed(() => tasksBacklogQuery.isPending.value)
const backlogError = computed(() => tasksBacklogQuery.error.value != null)
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
            <span v-else>{{ availableToday }}<span class="text-xl text-muted-foreground"> of {{ totalPeople }}</span> available</span>
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

    <!-- High priority tasks -->
    <Card>
      <CardHeader class="flex flex-row items-start justify-between gap-3">
        <div class="flex flex-col gap-1">
          <CardDescription>High priority tasks</CardDescription>
          <CardTitle class="text-3xl">
            <span v-if="backlogLoading" class="text-muted-foreground">Loading…</span>
            <span v-else-if="backlogError" class="text-destructive">Backend unavailable</span>
            <span v-else>{{ highPriorityTasks }}</span>
          </CardTitle>
        </div>
        <div class="flex size-10 items-center justify-center rounded-full bg-destructive-soft text-destructive">
          <TriangleAlertIcon />
        </div>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">Tasks marked as high priority</p>
      </CardContent>
    </Card>

    <!-- Unassigned jobs -->
    <Card>
      <CardHeader class="flex flex-row items-start justify-between gap-3">
        <div class="flex flex-col gap-1">
          <CardDescription>Unassigned jobs</CardDescription>
          <CardTitle class="text-3xl">
            <span v-if="backlogLoading" class="text-muted-foreground">Loading…</span>
            <span v-else-if="backlogError" class="text-destructive">Backend unavailable</span>
            <span v-else>{{ unassignedTasks }}</span>
          </CardTitle>
        </div>
        <div class="flex size-10 items-center justify-center rounded-full bg-warning-soft text-warning">
          <HammerIcon />
        </div>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">Tasks with no one assigned yet</p>
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
