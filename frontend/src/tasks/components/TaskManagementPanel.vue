<script setup lang="ts">
import { ListFilterIcon, PlusIcon } from '@lucide/vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Button } from '@/shared/ui/button'
import { Badge } from '@/shared/ui/badge'
import { taskFixtures } from '../fixtures'
import { priorityLabels, priorityBadgeVariant } from '../helpers'
import type { TaskPriority } from '../types'
import TaskRow from './TaskRow.vue'

defineEmits<{
  filter: []
  'add-task': []
}>()

const priorities: TaskPriority[] = ['high', 'medium', 'low']
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-center justify-between gap-4">
        <CardTitle>Task Management</CardTitle>
        <div class="flex items-center gap-2">
          <Button variant="outline" @click="$emit('filter')">
            <ListFilterIcon />
            Filter
          </Button>
          <Button @click="$emit('add-task')">
            <PlusIcon />
            Add Task
          </Button>
        </div>
      </div>
    </CardHeader>
    <CardContent>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-border">
              <th class="h-10 px-4 text-left text-sm font-medium text-muted-foreground">Task</th>
              <th class="h-10 px-4 text-left text-sm font-medium text-muted-foreground">Priority</th>
              <th class="h-10 px-4 text-left text-sm font-medium text-muted-foreground">People Needed</th>
              <th class="h-10 px-4 text-left text-sm font-medium text-muted-foreground">Room / Area</th>
              <th class="h-10 px-4 text-left text-sm font-medium text-muted-foreground">Status</th>
            </tr>
          </thead>
          <tbody>
            <TaskRow v-for="task in taskFixtures" :key="task.id" :task="task" />
          </tbody>
        </table>
      </div>

      <div class="mt-4 flex items-center gap-3 border-t border-border pt-4">
        <span class="text-sm text-muted-foreground">Priority:</span>
        <div class="flex items-center gap-2">
          <Badge v-for="p in priorities" :key="p" :variant="priorityBadgeVariant[p]">
            {{ priorityLabels[p] }}
          </Badge>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
