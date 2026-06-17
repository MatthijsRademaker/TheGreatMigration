<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import KpiCards from './components/KpiCards.vue'
import TaskManagementPanel from '@/tasks/components/TaskManagementPanel.vue'
import PeopleAvailability from '@/people/PeopleAvailability.vue'
import DailySchedule from '@/calendar/DailySchedule.vue'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
import { useDailySchedule } from '@/calendar/composables/useDailySchedule'


const { data: availabilityData } = usePeopleAvailability()
const {
  data: scheduleData,
  isLoading: scheduleLoading,
  isError: scheduleError,
  isEmpty: scheduleEmpty,
  page,
  totalPages,
  goToPrevPage,
  goToNextPage,
  dateRangeLabel,
} = useDailySchedule()
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
          :page="page"
          :total-pages="totalPages"
          :date-range-label="dateRangeLabel"
          @prev-page="goToPrevPage"
          @next-page="goToNextPage"
          read-only
        />
      </template>

      <Card>
        <CardHeader>
          <CardTitle>Move notes</CardTitle>
          <CardDescription>Keep reminders visible without turning the app into heavy project management.</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="flex flex-col gap-3 text-sm text-muted-foreground">
            <p>Check building access times before assigning early morning jobs.</p>
            <p>Keep tea, chargers, tape, markers, and bin bags in the first-day essentials box.</p>
          </div>
        </CardContent>
      </Card>
    </div>
  </section>
</template>
