<script setup lang="ts">
import { ref, computed } from 'vue'
import { parseDate } from '@internationalized/date'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import type { DateValue } from '@internationalized/date'
import {
  createScheduleCardMutation,
  updateScheduleCardMutation,
  deleteScheduleCardMutation,
  listRoomsQuery,
} from '@/client/@pinia/colada.gen'
import { useDailySchedule } from '@/calendar/composables/useDailySchedule'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
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
const roomsQuery = useQuery(listRoomsQuery())

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
const mutationError = ref<string | null>(null)
const formTitle = ref('')
const formPriority = ref<'high' | 'medium' | 'low'>('medium')
const formPeopleNeeded = ref(2)
const formRoomArea = ref('')
const formScheduledDate = ref('')
const formAssignedTo = ref<string[]>([])
const mutationLoading = ref(false)

// ---- DatePicker model bridge (DateValue ↔ YYYY-MM-DD string) ----
const scheduledDateModel = computed({
  get: (): DateValue | undefined =>
    formScheduledDate.value ? parseDate(formScheduledDate.value) : undefined,
  set: (val: DateValue | undefined) => {
    formScheduledDate.value = val ? val.toString() : ''
  },
})

// ---- Helpers ----
function resetForm() {
  editingId.value = null
  formTitle.value = ''
  formPriority.value = 'medium'
  formPeopleNeeded.value = 2
  formRoomArea.value = ''
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

function openEdit(card: { id: string; title: string; priority: string; roomArea: string; peopleNeeded: number; scheduledDate: string; assignedPeople?: { id: string }[] }) {
  resetForm()
  editingId.value = card.id
  formTitle.value = card.title
  formPriority.value = card.priority as 'high' | 'medium' | 'low'
  formPeopleNeeded.value = card.peopleNeeded
  formRoomArea.value = card.roomArea
  formScheduledDate.value = card.scheduledDate
  formAssignedTo.value = card.assignedPeople?.map(p => p.id) ?? []
  modalOpen.value = true
}

async function handleSubmit() {
  mutationError.value = null
  mutationLoading.value = true
  try {
    const body = {
      title: formTitle.value,
      priority: formPriority.value,
      roomArea: formRoomArea.value,
      peopleNeeded: formPeopleNeeded.value,
      scheduledDate: formScheduledDate.value,
      assignedTo: [...formAssignedTo.value],
    }
    if (editingId.value) {
      await updateMut.mutateAsync({
        path: { id: editingId.value },
        body,
      })
    } else {
      await createMut.mutateAsync({ body })
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
      :title="editingId ? 'Edit task' : 'Add task'"
      :description="editingId ? 'Update an existing schedule card.' : 'Create a new schedule card.'"
      submit-label="Save"
      :disabled="mutationLoading"
      :submitting="mutationLoading"
      @submit="handleSubmit"
      @cancel="handleCancel"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1.5">
          <label for="form-title" class="text-xs font-medium text-muted-foreground">Title</label>
          <Input id="form-title" v-model="formTitle" placeholder="Task title" />
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="form-priority" class="text-xs font-medium text-muted-foreground">Priority</label>
          <Select v-model="formPriority">
            <SelectTrigger id="form-priority">
              <SelectValue placeholder="Select priority" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="high">High</SelectItem>
              <SelectItem value="medium">Medium</SelectItem>
              <SelectItem value="low">Low</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="form-room" class="text-xs font-medium text-muted-foreground">Room / Area</label>
          <Select v-if="roomsQuery.isLoading.value" disabled>
            <SelectTrigger id="form-room">
              <SelectValue placeholder="Loading rooms…" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="loading" disabled>Loading…</SelectItem>
            </SelectContent>
          </Select>
          <div v-else-if="roomsQuery.error.value" class="flex items-center gap-2">
            <span class="text-xs text-destructive">Could not load rooms.</span>
            <Button variant="outline" size="sm" @click="roomsQuery.refetch()">Retry</Button>
          </div>
          <Select
            v-else
            v-model="formRoomArea"
          >
            <SelectTrigger id="form-room">
              <SelectValue placeholder="Select a room…" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="room in roomsQuery.data.value?.rooms ?? []"
                :key="room.id"
                :value="room.name"
              >
                {{ room.name }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="form-people" class="text-xs font-medium text-muted-foreground">People needed</label>
          <Input id="form-people" v-model.number="formPeopleNeeded" type="number" min="1" />
        </div>

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
