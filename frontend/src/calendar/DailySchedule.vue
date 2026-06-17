<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/shared/ui/badge'
import type { BadgeVariants } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import { Card } from '@/shared/ui/card'

interface AssignedPerson {
  id: string
  name: string
  initials: string
}

interface TaskCard {
  id: string
  title: string
  priority: 'high' | 'medium' | 'low'
  roomArea: string
  assignedPeople: AssignedPerson[]
  peopleNeeded: number
  assignedCount: number
  staffingStatus: 'fullyStaffed' | 'underStaffed'
  scheduledDate: string
  taskId: string | null
}

interface ScheduleDay {
  date: string
  label: string
  availablePeopleCount: number
  tasks: TaskCard[]
}

const priorityAccentMap: Record<TaskCard['priority'], string> = {
  high: 'border-l-destructive',
  medium: 'border-l-warning',
  low: 'border-l-success',
}

const priorityVariantMap: Record<TaskCard['priority'], BadgeVariants['variant']> = {
  high: 'priorityHigh',
  medium: 'priorityMedium',
  low: 'priorityLow',
}

interface DailyScheduleProps {
  days?: ScheduleDay[]
  readOnly?: boolean
  /** Page number for pagination (1-indexed). Renders a pagination bar when > 0 together with totalPages. */
  page?: number
  /** Total number of pages. */
  totalPages?: number
  /** Date range label shown in the pagination bar. */
  dateRangeLabel?: string
}

interface DailyScheduleEmits {
  "add-task": [date?: string]
  "edit-task": [task: TaskCard]
  "delete-task": [taskId: string]
  "prev-page": []
  "next-page": []
}

const props = withDefaults(defineProps<DailyScheduleProps>(), {
  readOnly: false,
  page: 0,
  totalPages: 0,
  dateRangeLabel: '',
})

const hasPagination = computed(() => props.page > 0 && props.totalPages > 0)

const emit = defineEmits<DailyScheduleEmits>()

const scheduleDays = computed(() => props.days ?? [])
</script>

<template>
  <Card class="relative">
    <!-- Compact header row: title left, controls right -->
    <div class="flex items-center justify-between border-b border-border px-4 py-3">
      <h2 class="text-lg font-semibold">Daily Schedule</h2>
      <div class="flex items-center gap-3">
        <template v-if="hasPagination">
          <span class="text-sm text-muted-foreground">{{ dateRangeLabel || '—' }}</span>
          <span class="text-sm text-muted-foreground">Page {{ page }} of {{ totalPages }}</span>
          <Button
            variant="outline"
            size="sm"
            :disabled="page <= 1"
            @click="emit('prev-page')"
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            :disabled="page >= totalPages"
            @click="emit('next-page')"
          >
            Next
          </Button>
        </template>
        <Button variant="link" size="sm">View by: Day</Button>
        <Button v-if="!readOnly" variant="link" size="sm" @click="emit('add-task')">Add task</Button>
      </div>
    </div>

    <!-- Day columns -->
    <div class="px-4 py-3">
      <div class="overflow-x-auto">
        <div class="flex gap-3">
          <div
            v-for="day in scheduleDays"
            :key="day.date"
            class="min-w-[260px] shrink-0"
          >
            <!-- Day header: compact horizontal layout -->
            <div class="flex items-baseline gap-2 mb-3">
              <span class="text-sm font-semibold">{{ day.label }}</span>
              <span class="text-xs text-muted-foreground">{{ day.availablePeopleCount }} available</span>
            </div>

            <!-- Task cards -->
            <div class="flex flex-col gap-3">
              <div
                v-for="task in day.tasks"
                :key="task.id"
                class="rounded-lg border bg-card shadow-sm p-3 border-l-4"
                :class="priorityAccentMap[task.priority]"
              >
                <div class="flex items-start justify-between gap-2 mb-2">
                  <span class="text-sm font-medium">{{ task.title }}</span>
                  <Badge :variant="priorityVariantMap[task.priority]">
                    {{ task.priority }}
                  </Badge>
                </div>
                <p
                  v-if="task.assignedPeople.length > 0"
                  class="text-xs text-muted-foreground mb-2"
                >
                  {{ task.assignedPeople.map(p => p.initials).join(', ') }}
                </p>
                <p class="text-xs text-muted-foreground">
                  {{ task.assignedCount }} / {{ task.peopleNeeded }}
                  <span
                    v-if="task.staffingStatus === 'underStaffed'"
                    class="text-destructive"
                  >&nbsp;— needs help</span>
                </p>
                <div v-if="!readOnly" class="flex items-center gap-1 mt-2">
                  <Button variant="ghost" size="xs" @click.stop="emit('edit-task', task)">Edit</Button>
                  <Button variant="ghost" size="xs" @click.stop="emit('delete-task', task.id)">Delete</Button>
                </div>
              </div>

              <!-- Task creation placeholder: only when not readOnly -->
              <div
                v-if="!readOnly"
                class="rounded-lg border-2 border-dashed border-muted-foreground/25 p-3 text-center"
              >
                <Button variant="ghost" size="xs" @click="emit('add-task', day.date)">+ Add task</Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Card>
</template>
