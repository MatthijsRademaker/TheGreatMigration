<script setup lang="ts">
import { planWindowDays } from '@/shared/lib/planWindow'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <CardTitle>Schedule board foundation</CardTitle>
        <CardDescription>Calendar work will live here as day columns, task blocks, and staffing warnings.</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid gap-3 grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7">
          <div v-for="planDay in planWindowDays" :key="planDay.dateString" class="flex min-h-36 flex-col gap-3 rounded-lg border bg-muted/40 p-4">
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
