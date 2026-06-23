<script setup lang="ts">
import { ref, watch } from 'vue'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import { Input } from '@/shared/ui/input'
import { Button } from '@/shared/ui/button'
import TaskManagementPanel from '@/tasks/components/TaskManagementPanel.vue'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { Checkbox } from '@/shared/ui/checkbox'
import AddOperationModal from '@/shared/components/AddOperationModal.vue'
import {
  createTaskMutation,
  updateTaskMutation,
  deleteTaskMutation,
  listRoomsQuery,
} from '@/client/@pinia/colada.gen'
import { useTaskBacklog } from '@/tasks/composables/useTaskBacklog'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
import type { TaskRow } from '@/tasks/types'

// ---- State ----
const modalOpen = ref(false)
const editingId = ref<string | null>(null)
const deletingId = ref<string | null>(null)
const mutationError = ref<string | null>(null)
const formTitle = ref('')
const formPriority = ref<'high' | 'medium' | 'low'>('medium')
const formPeopleNeeded = ref(2)
const formRoom = ref('')
const formStatus = ref<'backlog' | 'ready' | 'assigned'>('backlog')
const formAssignedTo = ref<string[]>([])

// ---- Queries ----
const { data, queryKey } = useTaskBacklog()
const { data: peopleData } = usePeopleAvailability()
const roomsQuery = useQuery(listRoomsQuery())

// ---- Mutations ----
const queryCache = useQueryCache()
const createMut = useMutation({
  ...createTaskMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: queryKey }),
})
const updateMut = useMutation({
  ...updateTaskMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: queryKey }),
})
const deleteMut = useMutation({
  ...deleteTaskMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: queryKey }),
})

// ---- Actions ----
function startNewTask() {
  editingId.value = null
  formTitle.value = ''
  formPriority.value = 'medium'
  formPeopleNeeded.value = 2
  formRoom.value = ''
  formStatus.value = 'backlog'
  formAssignedTo.value = []
  modalOpen.value = true
}

function startEdit(task: TaskRow) {
  editingId.value = task.id
  formTitle.value = task.title
  formPriority.value = task.priority as 'high' | 'medium' | 'low'
  formPeopleNeeded.value = task.peopleNeeded
  formRoom.value = task.room
  formStatus.value = task.status as 'backlog' | 'ready' | 'assigned'
  formAssignedTo.value = task.assignedTo ? [...task.assignedTo] : []
  modalOpen.value = true
}

function cancelEdit() {
  editingId.value = null
  formTitle.value = ''
  formPriority.value = 'medium'
  formPeopleNeeded.value = 2
  formRoom.value = ''
  formStatus.value = 'backlog'
  formAssignedTo.value = []
  modalOpen.value = false
}

function toggleAssignment(personId: string) {
  const idx = formAssignedTo.value.indexOf(personId)
  if (idx === -1) {
    formAssignedTo.value.push(personId)
  } else {
    formAssignedTo.value.splice(idx, 1)
  }
}

async function handleSubmit() {
  if (!formTitle.value.trim() || !formRoom.value.trim()) return
  mutationError.value = null

  try {
    if (editingId.value) {
      await updateMut.mutateAsync({
        path: { id: editingId.value },
        body: {
          title: formTitle.value.trim(),
          priority: formPriority.value,
          peopleNeeded: formPeopleNeeded.value,
          room: formRoom.value.trim(),
          status: formStatus.value,
          assignedTo: [...formAssignedTo.value],
        },
      })
      cancelEdit()
    } else {
      await createMut.mutateAsync({
        body: {
          title: formTitle.value.trim(),
          priority: formPriority.value,
          peopleNeeded: formPeopleNeeded.value,
          room: formRoom.value.trim(),
          status: formStatus.value,
          assignedTo: [...formAssignedTo.value],
        },
      })
      cancelEdit()
    }
  } catch (e: any) {
    mutationError.value = e?.message || 'Something went wrong. Please try again.'
  }
}

async function handleDelete(id: string) {
  deletingId.value = id
  mutationError.value = null
  try {
    await deleteMut.mutateAsync({ path: { id } })
  } catch (e: any) {
    mutationError.value = e?.message || 'Could not delete. Please try again.'
  } finally {
    deletingId.value = null
  }
}

// Reset form when a pending task is deleted while editing.
watch(() => data.value.tasks, (tasks) => {
  if (editingId.value && tasks && !tasks.find((t) => t.id === editingId.value)) {
    cancelEdit()
  }
})
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <TaskManagementPanel
      @add-task="startNewTask"
      @edit-task="(task: TaskRow) => startEdit(task)"
      @delete-task="(id: string) => handleDelete(id)"
    />

    <!-- Add/Edit Modal -->
    <AddOperationModal
    v-model:open="modalOpen"
    :title="editingId ? 'Edit Task' : 'Add Task'"
    :description="editingId ? 'Update the task details.' : 'Add a new task to the backlog.'"
    :submit-label="editingId ? 'Save' : 'Add'"
    :disabled="createMut.isLoading.value || updateMut.isLoading.value"
    :submitting="createMut.isLoading.value || updateMut.isLoading.value"
    @submit="handleSubmit"
    @cancel="cancelEdit"
  >
    <form
      class="flex flex-col gap-4"
      @submit.prevent="handleSubmit"
    >
      <!-- Mutation error inside modal body -->
      <p
        v-if="mutationError"
        class="rounded border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive"
      >
        {{ mutationError }}
      </p>

      <div class="grid gap-3">
        <div>
          <label class="mb-1 block text-xs font-medium text-muted-foreground">Title</label>
          <Input v-model="formTitle" placeholder="e.g. Pack kitchen boxes" />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-muted-foreground">Room / Area</label>
          <Select v-if="roomsQuery.isLoading.value" disabled>
            <SelectTrigger>
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
            v-model="formRoom"
          >
            <SelectTrigger>
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
        <div>
          <label class="mb-1 block text-xs font-medium text-muted-foreground">Priority</label>
          <Select v-model="formPriority">
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="high">High</SelectItem>
              <SelectItem value="medium">Medium</SelectItem>
              <SelectItem value="low">Low</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-muted-foreground">Status</label>
          <Select v-model="formStatus">
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="backlog">Backlog</SelectItem>
              <SelectItem value="ready">Ready</SelectItem>
              <SelectItem value="assigned">Assigned</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-muted-foreground">People Needed</label>
          <Input v-model.number="formPeopleNeeded" type="number" min="1" />
        </div>
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
    </form>
  </AddOperationModal>
  </section>
</template>
