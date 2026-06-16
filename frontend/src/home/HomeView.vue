<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import KpiCards from './components/KpiCards.vue'
import TaskManagementPanel from '@/tasks/components/TaskManagementPanel.vue'
import PeopleAvailability from '@/people/PeopleAvailability.vue'
import DailySchedule from '@/calendar/DailySchedule.vue'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'

const { data: availabilityData } = usePeopleAvailability()
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">

    <KpiCards />

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <TaskManagementPanel />
      <PeopleAvailability v-bind="availabilityData" />
    </div>

    <div class="grid gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <DailySchedule />

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
