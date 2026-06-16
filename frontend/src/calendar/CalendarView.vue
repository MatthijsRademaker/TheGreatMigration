<script setup lang="ts">
import DailySchedule from '@/calendar/DailySchedule.vue'
import { useDailySchedule } from '@/calendar/composables/useDailySchedule'

const { data: scheduleData } = useDailySchedule()

// TODO: Wire schedule-card mutation flows (create/update/delete) with query invalidation
// for the /calendar editable workflow. Use useDailySchedule() loading/error/empty for state handling.
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <DailySchedule :days="scheduleData.days.length > 0 ? scheduleData.days : undefined" />
  </section>
</template>
