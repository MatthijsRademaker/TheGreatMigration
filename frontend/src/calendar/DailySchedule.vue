<script setup lang="ts">
import { computed } from 'vue'
import { Avatar } from '@/shared/ui/avatar'
import { Badge } from '@/shared/ui/badge'
import type { BadgeVariants } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'

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
  /** Navigate to the previous page. */
  goToPrevPage?: () => void
  /** Navigate to the next page. */
  goToNextPage?: () => void
}

const props = withDefaults(defineProps<DailyScheduleProps>(), {
  readOnly: false,
  page: 0,
  totalPages: 0,
  dateRangeLabel: '',
  goToPrevPage: () => {},
  goToNextPage: () => {},
})

const hasPagination = computed(() => props.page > 0 && props.totalPages > 0)

const emit = defineEmits<{
  "add-task": [date?: string]
  "edit-task": [task: TaskCard]
  "delete-task": [taskId: string]
}>()

const scheduleDays = computed(() => props.days ?? [])
</script>

<template>
  <Card class="relative">
    <!-- Pagination bar -- sits above the schedule content, visually attached to the top of the card -->
    <div
      v-if="hasPagination"
      class="flex items-center justify-between border-b border-border px-4 py-2"
    >
      <span class="text-sm text-muted-foreground">
        {{ dateRangeLabel || '—' }}
      </span>
      <div class="flex items-center gap-2">
        <span class="text-sm text-muted-foreground">
          Page {{ page }} of {{ totalPages }}
        </span>
        <Button
          variant="outline"
          size="sm"
          :disabled="page <= 1"
          @click="goToPrevPage"
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          :disabled="page >= totalPages"
          @click="goToNextPage"
        >
          Next
        </Button>
      </div>
    </div>

    <CardHeader>
      <CardTitle>Daily Schedule</CardTitle>
    </CardHeader>
    <CardContent>
      <!-- Header controls -->
      <div class="flex items-center justify-between mb-4">
        <Button variant="link" size="sm">View by: Day</Button>
        <Button v-if="!readOnly" variant="link" size="sm" @click="emit('add-task')">Add task</Button>
      </div>

      <!-- Day columns -->
      <div class="overflow-x-auto">
        <div class="flex gap-4">
          <div
            v-for="day in scheduleDays"
            :key="day.date"
            class="min-w-[280px] shrink-0"
          >
            <!-- Day header -->
            <div class="mb-3">
              <h3 class="text-sm font-semibold">{{ day.label }}</h3>
              <p class="text-xs text-muted-foreground">
                {{ day.availablePeopleCount }} available
              </p>
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
                  <div class="flex items-center gap-1.5">
                    <Badge v-if="task.taskId" variant="outline" class="text-[10px] px-1.5 py-0">
                      From backlog
                    </Badge>
                    <Badge :variant="priorityVariantMap[task.priority]">
                      {{ task.priority }}
                    </Badge>
                  </div>
                </div>
                <div class="flex flex-wrap items-center gap-2 mb-2">
                  <span
                    v-for="person in task.assignedPeople"
                    :key="person.id"
                    class="inline-flex items-center gap-1"
                  >
                    <Avatar
                      :name="person.name"
                      :initials="person.initials"
                      class="size-5 text-xs"
                    />
                    <span class="text-xs text-muted-foreground">
                      {{ person.name }}
                    </span>
                  </span>
                </div>
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

              <!-- Task creation placeholder -->
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
    </CardContent>
  </Card>
</template>
