<script setup lang="ts">
import { Badge } from '@/shared/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import KpiCards from './components/KpiCards.vue'

const upcomingWork = [
  'Pack kitchen essentials and label fragile boxes',
  'Confirm van pickup window and parking access',
  'Clear garden tools from the shed',
]
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">

    <KpiCards />

    <div class="grid flex-1 gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <Card>
        <CardHeader>
          <div class="flex items-start justify-between gap-4">
            <div class="flex flex-col gap-1">
              <CardTitle>Today’s plan</CardTitle>
              <CardDescription>A calm command center for the next practical steps.</CardDescription>
            </div>
            <Badge variant="secondary">Draft schedule</Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div class="grid gap-3 md:grid-cols-3">
            <div v-for="item in upcomingWork" :key="item" class="rounded-lg border bg-muted/40 p-4">
              <p class="text-sm font-medium leading-6">{{ item }}</p>
            </div>
          </div>
        </CardContent>
      </Card>

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
