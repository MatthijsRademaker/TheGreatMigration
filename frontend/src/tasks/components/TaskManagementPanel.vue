<script setup lang="ts">
import { ListFilterIcon, PlusIcon } from '@lucide/vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Button } from '@/shared/ui/button'
import { Badge } from '@/shared/ui/badge'
import { useTaskBacklog } from '../composables/useTaskBacklog'
import { priorityLabels, priorityBadgeVariant } from '../helpers'
import type { TaskPriority } from '../types'
import TaskRow from './TaskRow.vue'

withDefaults(defineProps<{
  readOnly?: boolean
}>(), {
  readOnly: false,
})

defineEmits<{
  filter: []
  'add-task': []
  'edit-task': [task: import('../types').TaskRow]
  'delete-task': [taskId: string]
}>()

const { data, isLoading, isError, isEmpty } = useTaskBacklog()

const priorities: TaskPriority[] = ['high', 'medium', 'low']
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-center justify-between gap-4">
        <div class="flex items-center gap-2">
          <CardTitle>Tasks Backlog</CardTitle>
          <Badge v-if="data" variant="secondary">{{ data.summary.totalTasks }} tasks</Badge>
        </div>
        <div v-if="!readOnly" class="flex items-center gap-2">
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
      <!-- Loading state -->
      <div v-if="isLoading" class="py-8 text-center text-sm text-muted-foreground">
        Loading tasks&hellip;
      </div>

      <!-- Error state -->
      <div v-else-if="isError" class="py-8 text-center text-sm text-destructive">
        Could not load tasks. Please try again.
      </div>

      <!-- Empty state -->
      <div v-else-if="isEmpty" class="py-8 text-center text-sm text-muted-foreground">
        No tasks yet.
      </div>

      <!-- Data state -->
      <template v-else>
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
              <TaskRow v-for="task in data.tasks" :key="task.id" :task="task">
                <template v-if="!readOnly" #actions>
                  <Button variant="outline" size="sm" @click="$emit('edit-task', task)">Edit</Button>
                  <Button variant="outline" size="sm" @click="$emit('delete-task', task.id)">Delete</Button>
                </template>
              </TaskRow>
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
      </template>
    </CardContent>
  </Card>
</template>
