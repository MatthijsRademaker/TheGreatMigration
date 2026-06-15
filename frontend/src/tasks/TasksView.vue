<script setup lang="ts">
import { Badge } from '@/shared/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

const taskStatuses = [
  { label: 'Backlog', count: 8, tone: 'Needs planning' },
  { label: 'Ready', count: 5, tone: 'Can be scheduled' },
  { label: 'Assigned', count: 7, tone: 'People attached' },
]
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <CardTitle>Task foundation</CardTitle>
        <CardDescription>Tasks will own priority, staffing needs, status, room, and assignment behavior.</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid gap-3 md:grid-cols-3">
          <div v-for="status in taskStatuses" :key="status.label" class="rounded-lg border bg-muted/40 p-4">
            <div class="flex items-center justify-between gap-3">
              <p class="font-medium">{{ status.label }}</p>
              <Badge variant="secondary">{{ status.count }}</Badge>
            </div>
            <p class="mt-3 text-sm text-muted-foreground">{{ status.tone }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
