<script setup lang="ts">
import { ref, watch } from 'vue'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
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
import {
  getTasksBacklogQuery,
  createTaskMutation,
  updateTaskMutation,
  deleteTaskMutation,
  getTasksBacklogQueryKey,
} from '@/client/@pinia/colada.gen'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'

// ---- State ----
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
const backlogQuery = useQuery(getTasksBacklogQuery())
const { data: peopleData } = usePeopleAvailability()

// ---- Mutations ----
const queryCache = useQueryCache()
const createMut = useMutation({
  ...createTaskMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: getTasksBacklogQueryKey() }),
})
const updateMut = useMutation({
  ...updateTaskMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: getTasksBacklogQueryKey() }),
})
const deleteMut = useMutation({
  ...deleteTaskMutation(),
  onSuccess: () => queryCache.invalidateQueries({ key: getTasksBacklogQueryKey() }),
})

// ---- Actions ----
function startEdit(task: { id: string; title: string; priority: string; peopleNeeded: number; room: string; status: string; assignedTo?: string[] | null }) {
  editingId.value = task.id
  formTitle.value = task.title
  formPriority.value = task.priority as 'high' | 'medium' | 'low'
  formPeopleNeeded.value = task.peopleNeeded
  formRoom.value = task.room
  formStatus.value = task.status as 'backlog' | 'ready' | 'assigned'
  formAssignedTo.value = task.assignedTo ? [...task.assignedTo] : []
}

function cancelEdit() {
  editingId.value = null
  formTitle.value = ''
  formPriority.value = 'medium'
  formPeopleNeeded.value = 2
  formRoom.value = ''
  formStatus.value = 'backlog'
  formAssignedTo.value = []
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
      formTitle.value = ''
      formPriority.value = 'medium'
      formPeopleNeeded.value = 2
      formRoom.value = ''
      formStatus.value = 'backlog'
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
watch(() => backlogQuery.data.value?.tasks, (tasks) => {
  if (editingId.value && tasks && !tasks.find((t) => t.id === editingId.value)) {
    cancelEdit()
  }
})
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <!-- Mutation error -->
    <p
      v-if="mutationError"
      class="rounded border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive"
    >
      {{ mutationError }}
    </p>

    <!-- Create / Edit form -->
    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? 'Edit' : 'Add' }} Task</CardTitle>
        <CardDescription>
          {{ editingId ? 'Update the task details.' : 'Add a new task to the backlog.' }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          class="flex flex-col gap-3"
          @submit.prevent="handleSubmit"
        >
          <div class="grid gap-3 sm:grid-cols-2">
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">Title</label>
              <Input v-model="formTitle" placeholder="e.g. Pack kitchen boxes" />
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">Room / Area</label>
              <Input v-model="formRoom" placeholder="e.g. Kitchen" />
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

          <div class="flex gap-2">
            <Button type="submit" :disabled="createMut.isLoading.value || updateMut.isLoading.value">
              {{ editingId ? 'Save' : 'Add' }}
            </Button>
            <Button v-if="editingId" variant="outline" type="button" @click="cancelEdit">
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <!-- Task panel with edit/delete actions -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between gap-4">
          <div>
            <CardTitle>Task backlog</CardTitle>
            <CardDescription>
              Capture jobs, priorities, staffing needs, and planning status.
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <!-- Loading -->
        <p v-if="backlogQuery.isLoading.value" class="text-sm text-muted-foreground">
          Loading&hellip;
        </p>

        <!-- Error -->
        <p v-else-if="backlogQuery.error.value" class="text-sm text-destructive">
          Could not load tasks. Please try again.
        </p>

        <!-- Empty -->
        <p
          v-else-if="!backlogQuery.data.value?.tasks?.length"
          class="text-sm text-muted-foreground"
        >
          No tasks yet. Add one above.
        </p>

        <!-- List with management controls -->
        <ul v-else class="divide-y">
          <li
            v-for="task in backlogQuery.data.value.tasks"
            :key="task.id"
            class="flex items-center justify-between py-2 first:pt-0 last:pb-0"
          >
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium">{{ task.title }}</p>
              <p class="text-xs text-muted-foreground">
                {{ task.room }} &middot; {{ task.priority }} &middot; {{ task.status }} &middot; {{ task.peopleNeeded }} needed
              </p>
            </div>
            <div class="ml-4 flex shrink-0 gap-1">
              <Button
                variant="outline"
                size="sm"
                @click="startEdit(task)"
              >
                Edit
              </Button>
              <Button
                variant="outline"
                size="sm"
                :disabled="deletingId === task.id"
                @click="handleDelete(task.id)"
              >
                Delete
              </Button>
            </div>
          </li>
        </ul>
      </CardContent>
    </Card>
  </section>
</template>
