<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import KpiCards from './components/KpiCards.vue'
import TaskManagementPanel from '@/tasks/components/TaskManagementPanel.vue'
import PeopleAvailability from '@/people/PeopleAvailability.vue'
import DailySchedule from '@/calendar/DailySchedule.vue'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
import { useDailySchedule } from '@/calendar/composables/useDailySchedule'

const { data: availabilityData } = usePeopleAvailability()
const { data: scheduleData } = useDailySchedule()
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">

    <KpiCards />

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <TaskManagementPanel read-only />
      <PeopleAvailability v-bind="availabilityData" />
    </div>

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <!-- Daily Schedule: shows backend data on success, demo data otherwise (read-only) -->
      <DailySchedule :days="scheduleData.days.length > 0 ? scheduleData.days : undefined" />

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
