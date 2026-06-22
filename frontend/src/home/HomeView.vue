<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import KpiCards from './components/KpiCards.vue'
import TaskManagementPanel from '@/tasks/components/TaskManagementPanel.vue'
import PeopleAvailability from '@/people/PeopleAvailability.vue'
import DailySchedule from '@/calendar/DailySchedule.vue'
import ToolsPanel from '@/tools/components/ToolsPanel.vue'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
import { useDailySchedule } from '@/calendar/composables/useDailySchedule'
import { useHomePagination } from '@/shared/composables/useHomePagination'

const homePagination = useHomePagination()

const { data: availabilityData } = usePeopleAvailability({
  pageRef: homePagination.page,
  daysPerPageRef: homePagination.daysPerPage,
})
const {
  data: scheduleData,
  isLoading: scheduleLoading,
  isError: scheduleError,
  isEmpty: scheduleEmpty,
} = useDailySchedule({
  pageRef: homePagination.page,
  daysPerPageRef: homePagination.daysPerPage,
})
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">

    <KpiCards />

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <TaskManagementPanel read-only />
      <PeopleAvailability v-bind="availabilityData" />
    </div>

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <!-- Daily Schedule: state-driven rendering, read-only on home -->
      <Card v-if="scheduleLoading">
        <CardHeader>
          <CardTitle>Daily Schedule</CardTitle>
          <CardDescription>Loading schedule data…</CardDescription>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">Fetching task cards from the backend.</p>
        </CardContent>
      </Card>

      <Card v-else-if="scheduleError">
        <CardHeader>
          <CardTitle>Daily Schedule</CardTitle>
          <CardDescription>Backend unavailable</CardDescription>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-destructive">
            Could not load daily schedule data. The backend may be unreachable or experiencing issues.
          </p>
        </CardContent>
      </Card>

      <Card v-else-if="scheduleEmpty">
        <CardHeader>
          <CardTitle>Daily Schedule</CardTitle>
          <CardDescription>No tasks scheduled</CardDescription>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">
            There are currently no task cards scheduled. Add tasks from the calendar page.
          </p>
        </CardContent>
      </Card>

      <template v-else>
        <DailySchedule
          :days="scheduleData.days"
          hide-pagination
          read-only
        />
      </template>

      <ToolsPanel />
    </div>
  </section>
</template>
