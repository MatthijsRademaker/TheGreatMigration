<script setup lang="ts">
import { Badge } from '@/shared/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

const moveDays = [
  { day: 'Mon', focus: 'Packing', helpers: '4 people' },
  { day: 'Tue', focus: 'Disassembly', helpers: '3 people' },
  { day: 'Wed', focus: 'Loading', helpers: '6 people' },
  { day: 'Thu', focus: 'Travel', helpers: '2 people' },
  { day: 'Fri', focus: 'Unpack essentials', helpers: '5 people' },
]
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <CardTitle>Schedule board foundation</CardTitle>
        <CardDescription>Calendar work will live here as day columns, task blocks, and staffing warnings.</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid gap-3 md:grid-cols-5">
          <div v-for="day in moveDays" :key="day.day" class="flex min-h-36 flex-col gap-3 rounded-lg border bg-muted/40 p-4">
            <div class="flex items-center justify-between gap-2">
              <p class="font-semibold">{{ day.day }}</p>
              <Badge variant="outline">{{ day.helpers }}</Badge>
            </div>
            <p class="text-sm text-muted-foreground">{{ day.focus }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
