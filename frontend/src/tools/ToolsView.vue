<script setup lang="ts">
import { computed, ref } from 'vue'
import { PlusIcon, WrenchIcon, XIcon } from '@lucide/vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Button } from '@/shared/ui/button'
import { Badge } from '@/shared/ui/badge'
import { Input } from '@/shared/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { useTools } from '@/tools/composables/useTools'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'

const {
  data,
  isLoading,
  isError,
  isEmpty,
  createTool,
  deleteTool,
  claimTool,
  unclaimTool,
} = useTools()

// The bringer picker is fed by the existing people query — the same people
// shown on the availability screen, with no availability filtering.
const { data: peopleData } = usePeopleAvailability()

const peopleNameById = computed<Record<string, string>>(() => {
  const map: Record<string, string> = {}
  for (const person of peopleData.value.people ?? []) {
    map[person.id] = person.name
  }
  return map
})

function bringerName(personId: string): string {
  return peopleNameById.value[personId] ?? personId
}

// ---- Add tool ----
const newToolName = ref('')

async function handleAddTool() {
  const name = newToolName.value.trim()
  if (!name) return
  await createTool(name)
  newToolName.value = ''
}

async function handleClaim(toolId: string, personId: string) {
  if (!personId) return
  await claimTool(toolId, personId)
}
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between gap-4">
          <div class="flex items-center gap-2">
            <CardTitle>Tools</CardTitle>
            <Badge v-if="!isLoading && !isError" variant="secondary">
              {{ data.summary.claimed }} / {{ data.summary.total }} covered
            </Badge>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <!-- Add tool -->
        <form class="mb-4 flex items-center gap-2" @submit.prevent="handleAddTool">
          <Input
            v-model="newToolName"
            placeholder="Add a tool, e.g. Ladder"
            class="max-w-xs"
          />
          <Button type="submit" :disabled="!newToolName.trim()">
            <PlusIcon />
            Add Tool
          </Button>
        </form>

        <!-- Loading state -->
        <div v-if="isLoading" class="py-8 text-center text-sm text-muted-foreground">
          Loading tools&hellip;
        </div>

        <!-- Error state -->
        <div v-else-if="isError" class="py-8 text-center text-sm text-destructive">
          Could not load tools. Please try again.
        </div>

        <!-- Empty state -->
        <div v-else-if="isEmpty" class="py-8 text-center text-sm text-muted-foreground">
          No tools yet.
        </div>

        <!-- Data state -->
        <ul v-else class="flex flex-col divide-y divide-border">
          <li
            v-for="tool in data.tools"
            :key="tool.id"
            class="flex items-center justify-between gap-4 py-3"
            :data-testid="`tool-${tool.id}`"
            :data-claimed="tool.broughtBy ? 'true' : 'false'"
          >
            <div class="flex items-center gap-2">
              <WrenchIcon class="size-4 text-muted-foreground" />
              <span
                class="text-sm font-medium"
                :class="tool.broughtBy ? 'text-muted-foreground line-through' : ''"
              >
                {{ tool.name }}
              </span>
            </div>

            <div class="flex items-center gap-2">
              <!-- Crossed-off tool: show bringer + unclaim -->
              <template v-if="tool.broughtBy">
                <Badge variant="secondary">{{ bringerName(tool.broughtBy) }}</Badge>
                <Button variant="outline" size="sm" @click="unclaimTool(tool.id)">
                  Unclaim
                </Button>
              </template>

              <!-- Open tool: bringer picker -->
              <template v-else>
                <Select
                  :model-value="''"
                  @update:model-value="(v) => handleClaim(tool.id, String(v))"
                >
                  <SelectTrigger class="w-44">
                    <SelectValue placeholder="Claim — pick a person" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="person in peopleData.people ?? []"
                      :key="person.id"
                      :value="person.id"
                    >
                      {{ person.name }}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </template>

              <!-- Organizer remove control -->
              <Button
                variant="ghost"
                size="sm"
                :aria-label="`Remove ${tool.name}`"
                @click="deleteTool(tool.id)"
              >
                <XIcon />
              </Button>
            </div>
          </li>
        </ul>
      </CardContent>
    </Card>
  </section>
</template>
