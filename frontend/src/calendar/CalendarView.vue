<script setup lang="ts">
import { ref, computed } from 'vue'
import { parseDate } from '@internationalized/date'
import { useMutation, useQueryCache } from '@pinia/colada'
import type { DateValue } from '@internationalized/date'
import {
  createScheduleCardMutation,
  updateScheduleCardMutation,
  deleteScheduleCardMutation,
} from '@/client/@pinia/colada.gen'
import { useDailySchedule } from '@/calendar/composables/useDailySchedule'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
import { useTaskBacklog } from '@/tasks/composables/useTaskBacklog'
import DailySchedule from '@/calendar/DailySchedule.vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import { Button } from '@/shared/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { Checkbox } from '@/shared/ui/checkbox'
import { DatePicker } from '@/shared/ui/date-picker'
import AddOperationModal from '@/shared/components/AddOperationModal.vue'

// ---- Queries ----
const { data: scheduleData, isLoading, isError, isEmpty, queryKey } = useDailySchedule()
const { data: peopleData } = usePeopleAvailability()
const { data: backlog, isLoading: backlogLoading, isEmpty: backlogEmpty } = useTaskBacklog()

// ---- Mutations ----
const queryCache = useQueryCache()
const createMut = useMutation({
  ...createScheduleCardMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: queryKey }),
})
const updateMut = useMutation({
  ...updateScheduleCardMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: queryKey }),
})
const deleteMut = useMutation({
  ...deleteScheduleCardMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: queryKey }),
})

// ---- Modal / form state ----
const modalOpen = ref(false)
const editingId = ref<string | null>(null)
const isEditMode = computed(() => editingId.value !== null)
const mutationError = ref<string | null>(null)
const formTaskId = ref('')
const formScheduledDate = ref('')
const formAssignedTo = ref<string[]>([])
const mutationLoading = ref(false)

// ---- Selected task read-only display ----
const selectedTask = computed(() => {
  if (!formTaskId.value) return null
  return backlog.value.tasks.find(t => t.id === formTaskId.value) ?? null
})

// ---- DatePicker model bridge (DateValue ↔ YYYY-MM-DD string) ----
const scheduledDateModel = computed({
  get: (): DateValue | undefined =>
    formScheduledDate.value ? parseDate(formScheduledDate.value) : undefined,
  set: (val: DateValue | undefined) => {
    formScheduledDate.value = val ? val.toString() : ''
  },
})

// ---- Search filter for tasks ----
const taskSearch = ref('')
const filteredTasks = computed(() => {
  if (!taskSearch.value) return backlog.value.tasks
  const q = taskSearch.value.toLowerCase()
  return backlog.value.tasks.filter(t => t.title.toLowerCase().includes(q))
})

// ---- Helpers ----
function resetForm() {
  editingId.value = null
  formTaskId.value = ''
  taskSearch.value = ''
  formScheduledDate.value = ''
  formAssignedTo.value = []
  mutationError.value = null
}

function openCreate(date?: string) {
  resetForm()
  if (date) {
    formScheduledDate.value = date
  }
  modalOpen.value = true
}

function openEdit(card: { id: string; title: string; priority: string; roomArea: string; peopleNeeded: number; scheduledDate: string; assignedPeople?: { id: string }[]; taskId?: string | null }) {
  resetForm()
  editingId.value = card.id
  // For cards linked to a backlog task, show task reference as read-only
  if (card.taskId) {
    formTaskId.value = card.taskId
  }
  formScheduledDate.value = card.scheduledDate
  formAssignedTo.value = card.assignedPeople?.map(p => p.id) ?? []
  modalOpen.value = true
}

async function handleSubmit() {
  mutationError.value = null
  mutationLoading.value = true
  try {
    const hasTaskId = !!formTaskId.value
    const body: Record<string, unknown> = {
      scheduledDate: formScheduledDate.value,
      assignedTo: [...formAssignedTo.value],
    }
    if (hasTaskId) {
      // When referencing a backlog task, send taskId; omit title/priority/roomArea/peopleNeeded
      body.taskId = formTaskId.value
    } else {
      // Fall back to free-form fields for cards with no backlog reference
      body.title = ''
      body.priority = 'medium'
      body.roomArea = ''
      body.peopleNeeded = 2
    }
    if (editingId.value) {
      await updateMut.mutateAsync({
        path: { id: editingId.value },
        body: body as Parameters<typeof updateMut.mutateAsync>[0]['body'],
      })
    } else {
      await createMut.mutateAsync({
        body: body as Parameters<typeof createMut.mutateAsync>[0]['body'],
      })
    }
    modalOpen.value = false
  } catch (e: unknown) {
    mutationError.value = e instanceof Error ? e.message : 'An unexpected error occurred'
  } finally {
    mutationLoading.value = false
  }
}

async function handleDelete(id: string) {
  mutationError.value = null
  try {
    await deleteMut.mutateAsync({ path: { id } })
  } catch (e: unknown) {
    mutationError.value = e instanceof Error ? e.message : 'Failed to delete'
  }
}

function toggleAssignment(personId: string) {
  const idx = formAssignedTo.value.indexOf(personId)
  if (idx === -1) {
    formAssignedTo.value.push(personId)
  } else {
    formAssignedTo.value.splice(idx, 1)
  }
}

function handleCancel() {
  modalOpen.value = false
}


</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <!-- Loading state -->
    <Card v-if="isLoading">
      <CardHeader>
        <CardTitle>Calendar</CardTitle>
        <CardDescription>Loading schedule data…</CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">Fetching task cards from the backend.</p>
      </CardContent>
    </Card>

    <!-- Error state -->
    <Card v-else-if="isError">
      <CardHeader>
        <CardTitle>Calendar</CardTitle>
        <CardDescription>Backend unavailable</CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-destructive">
          Could not load daily schedule data. The backend may be unreachable or experiencing issues.
        </p>
      </CardContent>
    </Card>

    <!-- Empty state -->
    <Card v-else-if="isEmpty">
      <CardHeader>
        <CardTitle>Calendar</CardTitle>
        <CardDescription>No tasks scheduled</CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">
          There are currently no task cards scheduled.
        </p>
        <Button variant="outline" class="mt-3" @click="openCreate()">Add your first task</Button>
      </CardContent>
    </Card>

    <!-- Success -->
    <template v-else>
      <DailySchedule
        :days="scheduleData.days"
        @add-task="openCreate"
        @edit-task="openEdit"
        @delete-task="handleDelete"
      />
    </template>

    <!-- Add / Edit Modal -->
    <AddOperationModal
      v-model:open="modalOpen"
      :title="isEditMode && selectedTask ? 'Edit scheduled task' : isEditMode ? 'Edit schedule card' : 'Add task to schedule'"
      :description="isEditMode ? 'Adjust the scheduled date or assigned people.' : 'Select an existing backlog task to schedule. Tasks are defined in the Task Panel.'"
      submit-label="Save"
      :disabled="mutationLoading"
      :submitting="mutationLoading"
      @submit="handleSubmit"
      @cancel="handleCancel"
    >
      <div class="flex flex-col gap-4">
        <!-- Task selector (create mode, or edit mode for cards without taskId) -->
        <template v-if="!isEditMode || (!selectedTask && editingId)">
          <!-- Empty backlog state -->
          <div v-if="backlogEmpty && !backlogLoading" class="rounded border border-dashed border-muted-foreground/30 p-4 text-center">
            <p class="text-sm text-muted-foreground mb-2">
              No tasks in the backlog yet.
            </p>
            <p class="text-xs text-muted-foreground">
              Create tasks in the Task Panel first, then schedule them here.
            </p>
          </div>

          <!-- Loading state -->
          <div v-else-if="backlogLoading" class="py-4 text-center">
            <p class="text-sm text-muted-foreground">Loading tasks…</p>
          </div>

          <!-- Task selector with search -->
          <template v-else>
            <div class="flex flex-col gap-1.5">
              <label for="form-task-search" class="text-xs font-medium text-muted-foreground">Search task</label>
              <Input
                id="form-task-search"
                v-model="taskSearch"
                placeholder="Type to filter backlog tasks…"
              />
            </div>

            <div class="flex flex-col gap-1.5">
              <label for="form-task-select" class="text-xs font-medium text-muted-foreground">Select a task</label>
              <Select v-model="formTaskId">
                <SelectTrigger id="form-task-select">
                  <SelectValue :placeholder="filteredTasks.length === 0 ? 'No matching tasks' : 'Choose a task…'" />
                </SelectTrigger>
                <SelectContent class="max-h-[240px]">
                  <SelectItem
                    v-for="task in filteredTasks"
                    :key="task.id"
                    :value="task.id"
                  >
                    {{ task.title }}
                  </SelectItem>
                  <div v-if="filteredTasks.length === 0" class="px-2 py-4 text-center text-xs text-muted-foreground">
                    No tasks match your search.
                  </div>
                </SelectContent>
              </Select>
            </div>
          </template>
        </template>

        <!-- Edit mode: show read-only task info when card has taskId -->
        <template v-if="selectedTask">
          <div class="rounded border border-border bg-muted/30 p-3 space-y-2">
            <div class="flex items-center justify-between">
              <span class="text-xs font-medium text-muted-foreground">Task</span>
              <span class="text-xs text-muted-foreground">From backlog</span>
            </div>
            <p class="text-sm font-medium">{{ selectedTask.title }}</p>
            <div class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
              <span>Priority: <strong>{{ selectedTask.priority }}</strong></span>
              <span>Room: <strong>{{ selectedTask.room }}</strong></span>
              <span>People needed: <strong>{{ selectedTask.peopleNeeded }}</strong></span>
            </div>
          </div>
        </template>

        <!-- Scheduled date -->
        <div class="flex flex-col gap-1.5">
          <label for="form-date" class="text-xs font-medium text-muted-foreground">Scheduled date</label>
          <DatePicker id="form-date" v-model="scheduledDateModel" />
        </div>

        <!-- Assignment -->
        <fieldset v-if="peopleData.people?.length" class="rounded border border-border p-3">
          <legend class="px-1 text-xs font-medium text-muted-foreground">Assign People</legend>
          <div class="flex flex-wrap gap-3">
            <label
              v-for="person in peopleData.people"
              :key="person.id"
              class="flex items-center gap-1.5 text-sm"
            >
              <Checkbox
                :model-value="formAssignedTo.includes(person.id)"
                size="sm"
                @update:model-value="toggleAssignment(person.id)"
              />
              {{ person.name }}
            </label>
          </div>
        </fieldset>

        <p v-if="mutationError" class="text-sm text-destructive">{{ mutationError }}</p>
      </div>
    </AddOperationModal>
  </section>
</template>
