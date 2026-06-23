<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Button } from '@/shared/ui/button'
import { Card } from '@/shared/ui/card'
import { PersonChip } from '@/shared/ui/person-chip'
import Celebration from '@/shared/motion/Celebration.vue'
import TaskBoardCard from './components/TaskBoardCard.vue'

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

/** Minimal person shape for the draggable people rail. */
interface RailPerson {
  id: string
  name: string
}

interface DailyScheduleProps {
  days?: ScheduleDay[]
  readOnly?: boolean
  /** People available to drag onto task cards. Only shown when not read-only. */
  people?: RailPerson[]
  /** Page number for pagination (1-indexed). Renders a pagination bar when > 0 together with totalPages. */
  page?: number
  /** Total number of pages. */
  totalPages?: number
  /** Date range label shown in the pagination bar. */
  dateRangeLabel?: string
  /** When true, suppress the pagination bar regardless of page/totalPages values. */
  hidePagination?: boolean
}

interface DailyScheduleEmits {
  "add-task": [date?: string]
  "edit-task": [task: TaskCard]
  "delete-task": [taskId: string]
  "prev-page": []
  "next-page": []
  /** A person was dropped onto a task card to be assigned. */
  "assign-person": [cardId: string, person: RailPerson]
  /** A card was dropped onto a different day column to be rescheduled. */
  "reschedule-card": [cardId: string, date: string]
}

const props = withDefaults(defineProps<DailyScheduleProps>(), {
  readOnly: false,
  people: () => [],
  page: 0,
  totalPages: 0,
  dateRangeLabel: '',
  hidePagination: false,
})

const hasPagination = computed(() => !props.hidePagination && props.page > 0 && props.totalPages > 0)

const emit = defineEmits<DailyScheduleEmits>()

const scheduleDays = computed(() => props.days ?? [])
const interactive = computed(() => !props.readOnly)

function formatCompactRangeLabel(days: ScheduleDay[]): string {
  if (days.length === 0) return '—'

  const monthFmt = new Intl.DateTimeFormat('en-US', { month: 'short', timeZone: 'UTC' })
  const start = new Date(days[0].date)
  const end = new Date(days[days.length - 1].date)
  const startDay = start.getUTCDate()
  const endDay = end.getUTCDate()
  const startMonth = monthFmt.format(start)
  const endMonth = monthFmt.format(end)

  if (startMonth === endMonth) {
    return startDay === endDay ? `${startDay} ${endMonth}` : `${startDay}–${endDay} ${endMonth}`
  }

  return `${startDay} ${startMonth}–${endDay} ${endMonth}`
}

const compactDateRangeLabel = computed(() => formatCompactRangeLabel(scheduleDays.value))

// ── Done state ──────────────────────────────────────────────────────────────
// "Done" is currently a frontend-only affordance: the card briefly leaves the
// column using the same transition as deletion, then reappears greyed out.
const TASK_DONE_TRANSITION_MS = 220
const doneTaskIds = ref<Record<string, true>>({})
const hidingTaskIds = ref<Record<string, true>>({})
const doneTimers = new Map<string, ReturnType<typeof setTimeout>>()

function isTaskDone(taskId: string): boolean {
  return doneTaskIds.value[taskId] === true
}

function isTaskHidden(taskId: string): boolean {
  return hidingTaskIds.value[taskId] === true
}

const visibleTasksByDay = computed<Record<string, TaskCard[]>>(() => {
  const visible: Record<string, TaskCard[]> = {}
  for (const day of scheduleDays.value) {
    visible[day.date] = day.tasks.filter((task) => !isTaskHidden(task.id))
  }
  return visible
})

function clearDoneTimer(taskId: string) {
  const timer = doneTimers.get(taskId)
  if (!timer) return
  clearTimeout(timer)
  doneTimers.delete(taskId)
}

function markTaskDone(taskId: string) {
  if (props.readOnly || isTaskDone(taskId) || isTaskHidden(taskId)) return

  hidingTaskIds.value = {
    ...hidingTaskIds.value,
    [taskId]: true,
  }

  clearDoneTimer(taskId)
  doneTimers.set(taskId, setTimeout(() => {
    const nextHiding = { ...hidingTaskIds.value }
    delete nextHiding[taskId]
    hidingTaskIds.value = nextHiding
    doneTaskIds.value = {
      ...doneTaskIds.value,
      [taskId]: true,
    }
    doneTimers.delete(taskId)
  }, TASK_DONE_TRANSITION_MS))
}

watch(
  scheduleDays,
  (days) => {
    const liveTaskIds = new Set(days.flatMap((day) => day.tasks.map((task) => task.id)))
    const nextDone = Object.fromEntries(
      Object.entries(doneTaskIds.value).filter(([taskId]) => liveTaskIds.has(taskId)),
    ) as Record<string, true>
    const nextHiding = Object.fromEntries(
      Object.entries(hidingTaskIds.value).filter(([taskId]) => liveTaskIds.has(taskId)),
    ) as Record<string, true>

    for (const taskId of doneTimers.keys()) {
      if (!liveTaskIds.has(taskId)) {
        clearDoneTimer(taskId)
      }
    }

    doneTaskIds.value = nextDone
    hidingTaskIds.value = nextHiding
  },
  { deep: true },
)

onBeforeUnmount(() => {
  for (const taskId of doneTimers.keys()) {
    clearDoneTimer(taskId)
  }
})

// ── Day-complete celebration ────────────────────────────────────────────────
// A day is "complete" when it has tasks and every task is fully staffed. We
// bump a per-day trigger when a day transitions into completeness (not on
// initial load) so the burst only plays as a reward for the change.
function isDayComplete(day: ScheduleDay): boolean {
  return day.tasks.length > 0 && day.tasks.every((t) => t.staffingStatus === 'fullyStaffed')
}

const dayCelebrationTrigger = ref<Record<string, number>>({})

watch(
  scheduleDays,
  (days, previous) => {
    const prevComplete = new Map((previous ?? []).map((d) => [d.date, isDayComplete(d)]))
    for (const day of days) {
      // Only celebrate days that existed before (skip the first data load).
      if (!prevComplete.has(day.date)) continue
      if (isDayComplete(day) && !prevComplete.get(day.date)) {
        dayCelebrationTrigger.value[day.date] = (dayCelebrationTrigger.value[day.date] ?? 0) + 1
      }
    }
  },
  { deep: true },
)

// ── Drag-and-drop state ─────────────────────────────────────────────────────
type DragKind = 'person' | 'card' | null
const dragKind = ref<DragKind>(null)
const draggedPerson = ref<RailPerson | null>(null)
const draggedCardId = ref<string | null>(null)
const dropTargetCardId = ref<string | null>(null)
const dropTargetDayDate = ref<string | null>(null)

function resetDrag() {
  dragKind.value = null
  draggedPerson.value = null
  draggedCardId.value = null
  dropTargetCardId.value = null
  dropTargetDayDate.value = null
}

function onPersonDragStart(event: DragEvent, person: RailPerson) {
  dragKind.value = 'person'
  draggedPerson.value = person
  if (event.dataTransfer) {
    event.dataTransfer.setData('text/plain', person.id)
    event.dataTransfer.effectAllowed = 'copy'
  }
}

function onCardDragStart(event: DragEvent, cardId: string) {
  dragKind.value = 'card'
  draggedCardId.value = cardId
  if (event.dataTransfer) {
    event.dataTransfer.setData('text/plain', cardId)
    event.dataTransfer.effectAllowed = 'move'
  }
}

// Assignment: a person dropped onto a task card.
function onCardDragOver(event: DragEvent, cardId: string) {
  if (dragKind.value !== 'person') return
  event.preventDefault()
  dropTargetCardId.value = cardId
}

function onCardDrop(event: DragEvent, cardId: string) {
  if (dragKind.value !== 'person' || !draggedPerson.value) return
  event.preventDefault()
  emit('assign-person', cardId, draggedPerson.value)
  resetDrag()
}

// Reschedule: a card dropped onto a day column.
function onDayDragOver(event: DragEvent, date: string) {
  if (dragKind.value !== 'card') return
  event.preventDefault()
  dropTargetDayDate.value = date
}

function onDayDrop(event: DragEvent, date: string) {
  if (dragKind.value !== 'card' || !draggedCardId.value) return
  event.preventDefault()
  emit('reschedule-card', draggedCardId.value, date)
  resetDrag()
}
</script>

<template>
  <Card class="!gap-0 relative">
    <!-- Compact header row: title left, controls right -->
    <div class="flex flex-col gap-3 border-b border-border px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
      <h2 class="text-lg font-semibold">Daily Schedule</h2>
      <div class="flex flex-wrap items-center gap-2 sm:justify-end sm:gap-3">
        <template v-if="hasPagination">
          <span
            data-testid="date-range-label-full"
            class="hidden text-sm text-muted-foreground whitespace-nowrap sm:inline"
          >
            {{ dateRangeLabel || '—' }}
          </span>
          <span
            data-testid="date-range-label-compact"
            class="text-sm text-muted-foreground whitespace-nowrap sm:hidden"
          >
            {{ compactDateRangeLabel }}
          </span>
          <span
            data-testid="page-indicator"
            class="hidden text-sm text-muted-foreground whitespace-nowrap sm:inline"
          >
            Page {{ page }} of {{ totalPages }}
          </span>
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

    <!-- People rail: drag a helper onto a task card to assign them -->
    <div
      v-if="interactive && people.length"
      data-testid="people-rail"
      class="hidden flex-wrap items-center gap-2 border-b border-border px-4 py-3 sm:flex"
    >
      <span
        data-testid="people-rail-helper-text"
        class="text-xs font-medium text-muted-foreground"
      >
        Drag a helper onto a task:
      </span>
      <PersonChip
        v-for="person in people"
        :key="person.id"
        :name="person.name"
        draggable="true"
        data-testid="rail-person"
        class="cursor-grab select-none transition-transform active:scale-105 active:cursor-grabbing"
        @dragstart="onPersonDragStart($event, person)"
        @dragend="resetDrag"
      />
    </div>

    <!-- Day columns -->
    <div class="px-4 py-3">
      <div class="overflow-x-auto">
        <div class="flex gap-3">
          <div
            v-for="day in scheduleDays"
            :key="day.date"
            data-testid="day-column"
            class="relative min-w-[260px] shrink-0 rounded-lg p-1 transition-colors"
            :class="interactive && dragKind === 'card' && dropTargetDayDate === day.date
              ? 'bg-primary/5 ring-2 ring-primary/40'
              : ''"
            @dragover="onDayDragOver($event, day.date)"
            @drop="onDayDrop($event, day.date)"
          >
            <Celebration :trigger="dayCelebrationTrigger[day.date] ?? 0" />

            <!-- Day header: compact horizontal layout -->
            <div class="flex items-baseline gap-2 mb-3">
              <span class="text-sm font-semibold">{{ day.label }}</span>
              <span class="text-xs text-muted-foreground">{{ day.availablePeopleCount }} available</span>
            </div>

            <!-- Task cards -->
            <div class="flex flex-col gap-3">
              <TransitionGroup name="schedule-card" tag="div" class="flex flex-col gap-3">
                <TaskBoardCard
                  v-for="task in visibleTasksByDay[day.date] ?? []"
                  :key="task.id"
                  :task="task"
                  :done="isTaskDone(task.id)"
                  :read-only="readOnly"
                  :interactive="interactive && !isTaskDone(task.id)"
                  :drop-active="interactive && !isTaskDone(task.id) && dragKind === 'person' && dropTargetCardId === task.id"
                  :draggable="interactive && !isTaskDone(task.id) ? 'true' : undefined"
                  :class="interactive && !isTaskDone(task.id) ? 'cursor-grab active:cursor-grabbing' : ''"
                  @dragstart="onCardDragStart($event, task.id)"
                  @dragend="resetDrag"
                  @dragover="onCardDragOver($event, task.id)"
                  @drop="onCardDrop($event, task.id)"
                  @done="markTaskDone(task.id)"
                  @edit="emit('edit-task', task)"
                  @delete="emit('delete-task', task.id)"
                />
              </TransitionGroup>

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

<style scoped>
.schedule-card-enter-active,
.schedule-card-leave-active {
  transition: opacity 220ms ease, transform 220ms ease;
}

.schedule-card-enter-from,
.schedule-card-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.96);
}
</style>
