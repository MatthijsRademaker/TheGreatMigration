<script lang="ts">
export interface AssignedPerson {
  id: string
  name: string
  initials: string
}

export interface TaskCard {
  id: string
  title: string
  priority: 'high' | 'medium' | 'low'
  roomArea: string
  assignedPeople: AssignedPerson[]
  peopleNeeded: number
  assignedCount: number
  staffingStatus: 'fullyStaffed' | 'underStaffed'
}

export interface ScheduleDay {
  date: string
  label: string
  availablePeopleCount: number
  tasks: TaskCard[]
}

const defaultDays: ScheduleDay[] = [
  {
    date: '2026-07-02',
    label: '2 Jul (Tue)',
    availablePeopleCount: 6,
    tasks: [
      {
        id: 't1',
        title: 'Painting hall',
        priority: 'high',
        roomArea: 'Hallway',
        assignedPeople: [
          { id: 'p1', name: 'Alex', initials: 'A' },
          { id: 'p2', name: 'Morgan', initials: 'M' },
        ],
        peopleNeeded: 2,
        assignedCount: 2,
        staffingStatus: 'fullyStaffed',
      },
      {
        id: 't2',
        title: 'Steam walls',
        priority: 'medium',
        roomArea: 'Bathroom',
        assignedPeople: [{ id: 'p3', name: 'Sam', initials: 'S' }],
        peopleNeeded: 1,
        assignedCount: 1,
        staffingStatus: 'fullyStaffed',
      },
      {
        id: 't3',
        title: 'Clean up',
        priority: 'low',
        roomArea: 'Kitchen',
        assignedPeople: [{ id: 'p4', name: 'Riley', initials: 'R' }],
        peopleNeeded: 2,
        assignedCount: 1,
        staffingStatus: 'underStaffed',
      },
    ],
  },
  {
    date: '2026-07-03',
    label: '3 Jul (Wed)',
    availablePeopleCount: 7,
    tasks: [
      {
        id: 't4',
        title: 'Sanding',
        priority: 'high',
        roomArea: 'Living Room',
        assignedPeople: [
          { id: 'p1', name: 'Alex', initials: 'A' },
          { id: 'p5', name: 'Jordan', initials: 'J' },
        ],
        peopleNeeded: 3,
        assignedCount: 2,
        staffingStatus: 'underStaffed',
      },
      {
        id: 't5',
        title: 'Bedroom painting',
        priority: 'medium',
        roomArea: 'Bedroom',
        assignedPeople: [{ id: 'p3', name: 'Sam', initials: 'S' }],
        peopleNeeded: 1,
        assignedCount: 1,
        staffingStatus: 'fullyStaffed',
      },
    ],
  },
  {
    date: '2026-07-04',
    label: '4 Jul (Thu)',
    availablePeopleCount: 7,
    tasks: [
      {
        id: 't6',
        title: 'Touch up woodwork',
        priority: 'low',
        roomArea: 'Dining Room',
        assignedPeople: [{ id: 'p2', name: 'Morgan', initials: 'M' }],
        peopleNeeded: 1,
        assignedCount: 1,
        staffingStatus: 'fullyStaffed',
      },
      {
        id: 't7',
        title: 'Living room finishing',
        priority: 'high',
        roomArea: 'Living Room',
        assignedPeople: [
          { id: 'p4', name: 'Riley', initials: 'R' },
          { id: 'p1', name: 'Alex', initials: 'A' },
        ],
        peopleNeeded: 2,
        assignedCount: 2,
        staffingStatus: 'fullyStaffed',
      },
      {
        id: 't8',
        title: '2nd floor walls',
        priority: 'medium',
        roomArea: 'Upstairs',
        assignedPeople: [],
        peopleNeeded: 2,
        assignedCount: 0,
        staffingStatus: 'underStaffed',
      },
    ],
  },
  {
    date: '2026-07-05',
    label: '5 Jul (Fri)',
    availablePeopleCount: 5,
    tasks: [
      {
        id: 't9',
        title: 'Kitchen painting',
        priority: 'medium',
        roomArea: 'Kitchen',
        assignedPeople: [{ id: 'p3', name: 'Sam', initials: 'S' }],
        peopleNeeded: 1,
        assignedCount: 1,
        staffingStatus: 'fullyStaffed',
      },
      {
        id: 't10',
        title: 'Final clean',
        priority: 'low',
        roomArea: 'Whole House',
        assignedPeople: [{ id: 'p5', name: 'Jordan', initials: 'J' }],
        peopleNeeded: 2,
        assignedCount: 1,
        staffingStatus: 'underStaffed',
      },
    ],
  },
]

const priorityAccentMap: Record<TaskCard['priority'], string> = {
  high: 'border-l-destructive',
  medium: 'border-l-warning',
  low: 'border-l-success',
}
</script>

<script setup lang="ts">
import { Avatar } from '@/shared/ui/avatar'
import { Badge } from '@/shared/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'

export interface DailyScheduleProps {
  days?: ScheduleDay[]
}

withDefaults(defineProps<DailyScheduleProps>(), {
  days: () => defaultDays,
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Daily Schedule</CardTitle>
    </CardHeader>
    <CardContent>
      <!-- Header controls -->
      <div class="flex items-center justify-between mb-4">
        <span class="text-sm font-medium">View by: Day</span>
        <span class="text-sm font-medium">Add task</span>
      </div>

      <!-- Day columns -->
      <div class="overflow-x-auto">
        <div class="flex gap-4">
          <div
            v-for="day in days"
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
                  <Badge
                    :variant="`priority${task.priority.charAt(0).toUpperCase() + task.priority.slice(1)}` as any"
                  >
                    {{ task.priority }}
                  </Badge>
                </div>
                <div class="flex flex-wrap items-center gap-2 mb-2">
                  <span
                    v-for="person in task.assignedPeople"
                    :key="person.id"
                    class="inline-flex items-center gap-1"
                  >
                    <Avatar
                      :name="person.name"
                      class="size-5 text-xs"
                    />
                    <span class="text-xs text-muted-foreground">
                      {{ person.name }}
                    </span>
                  </span>
                </div>
                <p class="text-xs text-muted-foreground">
                  {{ task.assignedCount }} / {{ task.peopleNeeded }}
                </p>
              </div>

              <!-- Add task placeholder -->
              <div
                class="rounded-lg border-2 border-dashed border-muted-foreground/25 p-3 text-center"
              >
                <span class="text-xs text-muted-foreground">+ Add task</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
