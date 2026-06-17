<script setup lang="ts">
import { ref } from 'vue'
import { useMutation, useQueryCache } from '@pinia/colada'
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
import type { CellChangePayload } from './types'

const queryCache = useQueryCache()

const {
  data: availabilityData,
  daysISO,
  isLoading,
  isError,
  isEmpty,
  page,
  totalPages,
  goToPrevPage,
  goToNextPage,
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

// --- Error state ---
const createError = ref('')
const deleteError = ref('')
const updateError = ref('')

function clearErrors() {
  createError.value = ''
  deleteError.value = ''
  updateError.value = ''
}

// --- Create person form ---
const newName = ref('')
const newInitials = ref('')

async function handleCreate() {
  clearErrors()
  if (!newName.value.trim() || !newInitials.value.trim()) {
    createError.value = 'Name and initials are required.'
    return
  }
  try {
    await createMutation.mutateAsync({
      body: {
        name: newName.value.trim(),
        initials: newInitials.value.trim(),
      },
    })
    newName.value = ''
    newInitials.value = ''
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err)
    createError.value = `Failed to create person: ${msg}`
  }
}

// --- Delete person ---
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

// --- Handle cell update from editable matrix ---
async function handleCellUpdate(payload: CellChangePayload) {
  clearErrors()
  const { personId, dayIndex, status } = payload

  // Look up the ISO date from the adapted daysISO array (parallel to days labels).
  const date = daysISO.value[dayIndex]
  if (!date) {
    updateError.value = 'Selected cell cannot be mapped to a date. Try refreshing the page.'
    return
  }

  try {
    if (status === null) {
      await deleteAvailabilityMutation.mutateAsync({
        path: { id: personId, date },
      })
    } else {
      await upsertMutation.mutateAsync({
        path: { id: personId, date },
        body: { status },
      })
    }
  } catch (err: unknown) {
    const httpStatus = getHttpErrorStatus(err)
    if (httpStatus === 400) {
      // Inspect the Huma error detail field for distinct 400 subtypes.
      const detail = ((err as any)?.cause?.body?.detail as string | undefined) ?? ''
      if (detail.includes('outside the planning window')) {
        updateError.value = 'This date is outside the planning window.'
      } else if (detail.includes('must be a valid ISO 8601 date')) {
        updateError.value = 'Invalid date format.'
      } else if (detail.includes('status must be one of')) {
        updateError.value = 'Invalid status value.'
      } else {
        updateError.value = `Failed to update: ${detail}`
      }
    } else if (httpStatus === 404) {
      updateError.value = 'Person not found.'
    } else {
      const msg = err instanceof Error ? err.message : String(err)
      updateError.value = `Failed to update availability: ${msg}`
    }
  }
}

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

    <!-- Success: create form + editable matrix -->
    <template v-else>
      <!-- Create person form -->
      <Card>
        <CardHeader>
          <CardTitle>Add a person</CardTitle>
          <CardDescription>Create a new person. The ID is assigned server-side.</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="flex flex-wrap items-end gap-3">
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

      <!-- Mutation error display -->
      <p v-if="deleteError" class="text-sm text-destructive">{{ deleteError }}</p>
      <p v-if="updateError" class="text-sm text-destructive">{{ updateError }}</p>

      <!-- Day pagination navigation -->
      <div
        class="flex items-center justify-between rounded-lg border border-border bg-card px-4 py-2"
      >
        <span class="text-sm text-muted-foreground">
          {{ availabilityData.days?.[0] ?? '—' }}
          –
          {{ availabilityData.days?.[availabilityData.days.length - 1] ?? '—' }}
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

      <!-- People availability matrix (editable) -->
      <PeopleAvailability
        v-bind="availabilityData"
        :editable="true"
        :deleting-person-id="deletingId"
        :updating="isMutationLoading(upsertMutation) || isMutationLoading(deleteAvailabilityMutation)"
        @update-cell="handleCellUpdate"
        @delete-person="handleDelete"
      />
    </template>
  </section>
</template>
