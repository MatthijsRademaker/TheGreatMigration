<script setup lang="ts">
import { usePlanningWindow } from '@/shared/composables/usePlanningWindow'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

const { planWindowDays, isLoading, isError } = usePlanningWindow()
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <CardTitle>Schedule board foundation</CardTitle>
        <CardDescription>Calendar work will live here as day columns, task blocks, and staffing warnings.</CardDescription>
      </CardHeader>
      <CardContent>
        <!-- Loading skeleton -->
        <div v-if="isLoading" class="grid gap-3 grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7">
          <div
            v-for="n in 7"
            :key="n"
            class="min-h-36 rounded-lg border bg-muted/40 p-4 animate-pulse"
          />
        </div>

        <!-- Error state -->
        <div v-else-if="isError" class="text-destructive text-sm">
          Planning window unavailable — check backend
        </div>

        <!-- Success: day columns -->
        <div v-else class="grid gap-3 grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7">
          <div
            v-for="planDay in planWindowDays"
            :key="planDay.dateString"
            data-testid="plan-day-column"
            class="flex min-h-36 flex-col gap-3 rounded-lg border bg-muted/40 p-4"
          >
            <div class="flex items-center justify-between gap-2">
              <p class="font-semibold">{{ planDay.label }}</p>
            </div>
            <p class="text-sm text-muted-foreground">Planned move day</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
