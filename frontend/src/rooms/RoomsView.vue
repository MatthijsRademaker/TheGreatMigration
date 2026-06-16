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
import {
  listRoomsQuery,
  createRoomMutation,
  updateRoomMutation,
  deleteRoomMutation,
} from '@/client/@pinia/colada.gen'

// ---- State ----
const editingId = ref<string | null>(null)
const formName = ref('')
const formType = ref<'room' | 'area'>('room')

// ---- Queries ----
const roomsQuery = useQuery(listRoomsQuery())

// ---- Mutations ----
const queryCache = useQueryCache()
const createMut = useMutation({
  ...createRoomMutation(),
  onSettled: () => queryCache.invalidateQueries({ key: ['listRooms'] }),
})
const updateMut = useMutation({
  ...updateRoomMutation(),
  onSettled: () => queryCache.invalidateQueries({ key: ['listRooms'] }),
})
const deleteMut = useMutation({
  ...deleteRoomMutation(),
  onSettled: () => queryCache.invalidateQueries({ key: ['listRooms'] }),
})

// ---- Actions ----
function startEdit(room: { id: string; name: string; type: string }) {
  editingId.value = room.id
  formName.value = room.name
  formType.value = room.type as 'room' | 'area'
}

function cancelEdit() {
  editingId.value = null
  formName.value = ''
  formType.value = 'room'
}

async function handleSubmit() {
  if (!formName.value.trim()) return

  if (editingId.value) {
    await updateMut.mutateAsync({
      path: { id: editingId.value },
      body: { name: formName.value.trim(), type: formType.value },
    })
    cancelEdit()
  } else {
    await createMut.mutateAsync({
      body: { name: formName.value.trim(), type: formType.value },
    })
    formName.value = ''
    formType.value = 'room'
  }
}

async function handleDelete(id: string) {
  await deleteMut.mutateAsync({ path: { id } })
}

// Reset form when a pending room is deleted while editing.
watch(() => roomsQuery.data.value?.rooms, (rooms) => {
  if (editingId.value && rooms && !rooms.find((r) => r.id === editingId.value)) {
    cancelEdit()
  }
})
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <!-- Create / Edit form -->
    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? 'Edit' : 'Add' }} Room / Area</CardTitle>
        <CardDescription>
          {{ editingId ? 'Update the room name and type.' : 'Define a new room or area for your move plan.' }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          class="flex flex-col gap-3 sm:flex-row sm:items-end"
          @submit.prevent="handleSubmit"
        >
          <div class="flex-1">
            <label class="mb-1 block text-xs font-medium text-muted-foreground">Name</label>
            <Input v-model="formName" placeholder="e.g. Basement" />
          </div>
          <div class="w-full sm:w-44">
            <label class="mb-1 block text-xs font-medium text-muted-foreground">Type</label>
            <Select v-model="formType">
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="room">Room</SelectItem>
                <SelectItem value="area">Area</SelectItem>
              </SelectContent>
            </Select>
          </div>
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

    <!-- Room list -->
    <Card>
      <CardHeader>
        <CardTitle>Rooms &amp; Areas</CardTitle>
        <CardDescription>
          {{ roomsQuery.data.value?.rooms?.length ?? 0 }} record{{ roomsQuery.data.value?.rooms?.length === 1 ? '' : 's' }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <!-- Loading -->
        <p v-if="roomsQuery.isLoading.value" class="text-sm text-muted-foreground">
          Loading&hellip;
        </p>

        <!-- Error -->
        <p v-else-if="roomsQuery.error.value" class="text-sm text-destructive">
          Could not load rooms. Please try again.
        </p>

        <!-- Empty -->
        <p
          v-else-if="!roomsQuery.data.value?.rooms?.length"
          class="text-sm text-muted-foreground"
        >
          No rooms or areas yet. Add one above.
        </p>

        <!-- List -->
        <ul v-else class="divide-y">
          <li
            v-for="room in roomsQuery.data.value.rooms"
            :key="room.id"
            class="flex items-center justify-between py-2 first:pt-0 last:pb-0"
          >
            <div>
              <p class="text-sm font-medium">
                {{ room.name }}
                <span class="ml-1 text-xs text-muted-foreground">
                  ({{ room.type }})
                </span>
              </p>
            </div>
            <div class="flex gap-1">
              <Button
                variant="outline"
                size="sm"
                @click="startEdit(room)"
              >
                Edit
              </Button>
              <Button
                variant="outline"
                size="sm"
                :disabled="deleteMut.isLoading.value"
                @click="handleDelete(room.id)"
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
