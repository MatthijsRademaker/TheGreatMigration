<script setup lang="ts">
import { ref } from 'vue'
import { useMutation, useQueryCache } from '@pinia/colada'
import { Badge } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import PeopleAvailability from './PeopleAvailability.vue'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
import {
  createPersonMutation,
  deletePersonMutation,
  deletePersonAvailabilityMutation,
  getDashboardPeopleAvailabilityQueryKey,
  upsertPersonAvailabilityMutation,
} from '@/client/@pinia/colada.gen'
import { getHttpErrorStatus } from '@/shared/lib/errorStatus'
import type { AvailabilityStatus } from './types'

const queryCache = useQueryCache()

const {
  data: availabilityData,
  rawData,
  isLoading,
  isError,
  isEmpty,
} = usePeopleAvailability()

// --- Mutation state ---
const createMutation = useMutation({
  ...createPersonMutation(),
  onSettled: () => {
    queryCache.invalidateQueries({ key: getDashboardPeopleAvailabilityQueryKey() })
  },
})
const deleteMutation = useMutation({
  ...deletePersonMutation(),
  onSettled: () => {
    queryCache.invalidateQueries({ key: getDashboardPeopleAvailabilityQueryKey() })
  },
})
const upsertMutation = useMutation({
  ...upsertPersonAvailabilityMutation(),
  onSettled: () => {
    queryCache.invalidateQueries({ key: getDashboardPeopleAvailabilityQueryKey() })
  },
})
const deleteAvailabilityMutation = useMutation({
  ...deletePersonAvailabilityMutation(),
  onSettled: () => {
    queryCache.invalidateQueries({ key: getDashboardPeopleAvailabilityQueryKey() })
  },
})

// --- Error state clearing ---
function clearErrors() {
  createError.value = ''
  deleteError.value = ''
  statusError.value = ''
  clearAvailabilityError.value = ''
}

// --- Create person form ---
const newId = ref('')
const newName = ref('')
const newInitials = ref('')
const createError = ref('')

async function handleCreate() {
  clearErrors()
  if (!newId.value.trim() || !newName.value.trim() || !newInitials.value.trim()) {
    createError.value = 'All fields are required.'
    return
  }
  try {
    await createMutation.mutateAsync({
      body: {
        id: newId.value.trim(),
        name: newName.value.trim(),
        initials: newInitials.value.trim(),
      },
    })
    newId.value = ''
    newName.value = ''
    newInitials.value = ''
  } catch (err: unknown) {
    const status = getHttpErrorStatus(err)
    if (status === 409) {
      createError.value = 'A person with this ID already exists.'
    } else {
      const msg = err instanceof Error ? err.message : String(err)
      createError.value = `Failed to create person: ${msg}`
    }
  }
}

// --- Delete person ---
const deleteError = ref('')
const deletingId = ref<string | null>(null)

async function handleDelete(id: string) {
  clearErrors()
  deletingId.value = id
  try {
    await deleteMutation.mutateAsync({ path: { id } })
    deletingId.value = null
  } catch (err: unknown) {
    deletingId.value = null
    const status = getHttpErrorStatus(err)
    if (status === 409) {
      deleteError.value = `Cannot delete this person: they are referenced by existing assignments.`
    } else if (status === 404) {
      deleteError.value = 'Person not found.'
    } else {
      const msg = err instanceof Error ? err.message : String(err)
      deleteError.value = `Failed to delete person: ${msg}`
    }
  }
}

// --- Status update ---
const editingCell = ref<{ personId: string; date: string } | null>(null)
const statusError = ref('')

async function handleStatusUpdate(status: string) {
  if (!editingCell.value) return
  clearErrors()
  const { personId, date } = editingCell.value
  try {
    await upsertMutation.mutateAsync({
      path: { id: personId, date },
      body: { status: status as AvailabilityStatus },
    })
    editingCell.value = null
  } catch (err: unknown) {
    const status = getHttpErrorStatus(err)
    if (status === 400) {
      statusError.value = 'Invalid status or date.'
    } else if (status === 404) {
      statusError.value = 'Person not found.'
    } else {
      const msg = err instanceof Error ? err.message : String(err)
      statusError.value = `Failed to update status: ${msg}`
    }
    editingCell.value = null
  }
}

// --- Clear availability (delete) ---
const clearAvailabilityError = ref('')

async function handleClearAvailability() {
  if (!editingCell.value) return
  clearErrors()
  const { personId, date } = editingCell.value
  try {
    await deleteAvailabilityMutation.mutateAsync({
      path: { id: personId, date },
    })
    editingCell.value = null
  } catch (err: unknown) {
    const status = getHttpErrorStatus(err)
    if (status === 404) {
      clearAvailabilityError.value = 'Person not found.'
    } else {
      const msg = err instanceof Error ? err.message : String(err)
      clearAvailabilityError.value = `Failed to clear availability: ${msg}`
    }
    editingCell.value = null
  }
}

// --- Derive ISO date from day index ---
function getISODate(dayIndex: number): string {
  if (!rawData.value?.range) return ''
  const start = new Date(rawData.value.range.startDate)
  start.setDate(start.getDate() + dayIndex)
  return start.toISOString().slice(0, 10)
}

const statusOptions: AvailabilityStatus[] = ['available', 'busy', 'partial', 'off']

// Helper to safely access mutation loading state (SSR-safe).
function isMutationLoading(mutation: typeof createMutation | typeof deleteMutation | typeof upsertMutation | typeof deleteAvailabilityMutation): boolean {
  return mutation?.isLoading?.value ?? false
}
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <!-- Loading state -->
    <Card v-if="isLoading">
      <CardHeader>
        <CardTitle>People availability</CardTitle>
        <CardDescription>Loading availability data…</CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">Fetching people and their statuses from the backend.</p>
      </CardContent>
    </Card>

    <!-- Error state -->
    <Card v-else-if="isError">
      <CardHeader>
        <CardTitle>People availability</CardTitle>
        <CardDescription>Backend unavailable</CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-destructive">
          Could not load availability data. The backend may be unreachable or experiencing issues.
        </p>
      </CardContent>
    </Card>

    <!-- Empty state -->
    <Card v-else-if="isEmpty">
      <CardHeader>
        <CardTitle>People availability</CardTitle>
        <CardDescription>No people found</CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">
          There are currently no people with availability data.
        </p>
      </CardContent>
    </Card>

    <!-- Success: people matrix + management controls -->
    <template v-else>
      <!-- Create person form -->
      <Card>
        <CardHeader>
          <CardTitle>Add a person</CardTitle>
          <CardDescription>Create a new person with a short unique ID.</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="flex flex-wrap items-end gap-3">
            <div class="flex flex-col gap-1">
              <label for="new-id" class="text-xs font-medium text-muted-foreground">ID</label>
              <Input
                id="new-id"
                v-model="newId"
                placeholder="e.g. p9"
                class="w-28"
                :disabled="isMutationLoading(createMutation)"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label for="new-name" class="text-xs font-medium text-muted-foreground">Name</label>
              <Input
                id="new-name"
                v-model="newName"
                placeholder="Full name"
                class="w-48"
                :disabled="isMutationLoading(createMutation)"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label for="new-initials" class="text-xs font-medium text-muted-foreground">Initials</label>
              <Input
                id="new-initials"
                v-model="newInitials"
                placeholder="e.g. JS"
                class="w-20"
                maxlength="10"
                :disabled="isMutationLoading(createMutation)"
              />
            </div>
            <Button
              variant="default"
              :disabled="isMutationLoading(createMutation)"
              @click="handleCreate"
            >
              {{ isMutationLoading(createMutation) ? 'Creating…' : 'Create' }}
            </Button>
          </div>
          <p v-if="createError" class="mt-2 text-sm text-destructive">{{ createError }}</p>
        </CardContent>
      </Card>

      <!-- People availability matrix -->
      <PeopleAvailability v-bind="availabilityData" />

      <!-- Per-person management controls -->
      <Card>
        <CardHeader>
          <CardTitle>Manage people</CardTitle>
          <CardDescription>Update statuses or remove people no longer needed.</CardDescription>
        </CardHeader>
        <CardContent>
          <!-- Status editing -->
          <div v-if="editingCell" class="mb-4 rounded-md border p-3">
            <p class="mb-2 text-sm">
              Change status for <strong>{{ editingCell.personId }}</strong> on {{ editingCell.date }}:
            </p>
            <div class="flex flex-wrap gap-2">
              <Badge
                v-for="s in statusOptions"
                :key="s"
                :variant="s"
                class="cursor-pointer hover:ring-2 hover:ring-ring"
                @click="handleStatusUpdate(s)"
              >
                {{ s.charAt(0).toUpperCase() + s.slice(1) }}
              </Badge>
            </div>
            <div class="mt-2 flex flex-wrap gap-2">
              <Button
                variant="outline"
                size="sm"
                :disabled="isMutationLoading(deleteAvailabilityMutation)"
                @click="handleClearAvailability"
              >
                {{ isMutationLoading(deleteAvailabilityMutation) ? 'Clearing…' : 'Clear (reset to off)' }}
              </Button>
              <Button variant="ghost" size="sm" @click="editingCell = null">
                Cancel
              </Button>
            </div>
          </div>
          <p v-if="clearAvailabilityError" class="mb-3 text-sm text-destructive">{{ clearAvailabilityError }}</p>
          <p v-if="statusError" class="mb-3 text-sm text-destructive">{{ statusError }}</p>
          <p v-if="deleteError" class="mb-3 text-sm text-destructive">{{ deleteError }}</p>

          <!-- Person list with actions -->
          <div class="space-y-2">
            <div
              v-for="person in availabilityData.people"
              :key="person.id"
              class="flex flex-wrap items-center gap-2 rounded-md border px-3 py-2"
            >
              <span class="min-w-[120px] text-sm font-medium">{{ person.name }}</span>

              <!-- Day cells as clickable badges -->
              <div class="flex flex-wrap gap-1">
                <Badge
                  v-for="(entry, dayIdx) in person.availability"
                  :key="`${person.id}-${entry.date}`"
                  :variant="entry.status"
                  class="cursor-pointer hover:ring-2 hover:ring-ring"
                  @click="clearErrors(); editingCell = { personId: person.id, date: getISODate(dayIdx) }"
                >
                  {{ entry.status.charAt(0).toUpperCase() + entry.status.slice(1) }}
                </Badge>
              </div>

              <div class="ml-auto">
                <Button
                  variant="destructive"
                  size="sm"
                  :disabled="isMutationLoading(deleteMutation) && deletingId === person.id"
                  @click="handleDelete(person.id)"
                >
                  {{ isMutationLoading(deleteMutation) && deletingId === person.id ? 'Deleting…' : 'Delete' }}
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </template>
  </section>
</template>
