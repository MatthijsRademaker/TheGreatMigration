<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { ScheduleSkeleton } from '@/shared/ui/skeleton'
import KpiCards from './components/KpiCards.vue'
import MigrationJourney from './components/MigrationJourney.vue'
import TaskManagementPanel from '@/tasks/components/TaskManagementPanel.vue'
import PeopleAvailability from '@/people/PeopleAvailability.vue'
import DailySchedule from '@/calendar/DailySchedule.vue'
import ToolsPanel from '@/tools/components/ToolsPanel.vue'
import NavArrow from '@/shared/ui/nav-arrow/NavArrow.vue'
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

    <MigrationJourney :days="scheduleData.days" class="hidden sm:block" />

    <KpiCards />

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <!-- Daily Schedule: state-driven rendering, read-only on home -->
      <div class="relative min-w-0">
        <Card v-if="scheduleLoading">
          <CardHeader>
            <CardTitle>Daily Schedule</CardTitle>
            <CardDescription>Loading schedule data…</CardDescription>
          </CardHeader>
          <CardContent>
            <ScheduleSkeleton />
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

        <NavArrow to="/calendar" class="absolute top-2 right-2" />
      </div>

      <div class="relative min-w-0">
        <ToolsPanel />
        <NavArrow to="/tools" class="absolute top-2 right-2" />
      </div>
    </div>

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <div class="relative min-w-0">
        <TaskManagementPanel read-only />
        <NavArrow to="/tasks" class="absolute top-2 right-2" />
      </div>
      <div class="relative min-w-0">
        <PeopleAvailability v-bind="availabilityData" />
        <NavArrow to="/people" class="absolute top-2 right-2" />
      </div>
    </div>
  </section>
</template>
