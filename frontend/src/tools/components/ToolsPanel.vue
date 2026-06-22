<script setup lang="ts">
import { WrenchIcon } from '@lucide/vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Badge } from '@/shared/ui/badge'
import { useTools } from '@/tools/composables/useTools'
import { usePeopleAvailability } from '@/shared/composables/usePeopleAvailability'
import { computed } from 'vue'

const { data, isLoading, isError, isEmpty } = useTools()
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
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-center gap-2">
        <CardTitle>Tools</CardTitle>
        <Badge v-if="!isLoading && !isError" variant="secondary">
          {{ data.summary.claimed }} / {{ data.summary.total }} covered
        </Badge>
      </div>
    </CardHeader>
    <CardContent>
      <div v-if="isLoading" class="py-4 text-center text-sm text-muted-foreground">
        Loading tools&hellip;
      </div>
      <div v-else-if="isError" class="py-4 text-center text-sm text-destructive">
        Could not load tools.
      </div>
      <div v-else-if="isEmpty" class="py-4 text-center text-sm text-muted-foreground">
        No tools yet.
      </div>
      <ul v-else class="flex flex-col divide-y divide-border">
        <li
          v-for="tool in data.tools"
          :key="tool.id"
          class="flex items-center justify-between gap-2 py-2"
        >
          <div class="flex items-center gap-2 min-w-0">
            <WrenchIcon class="size-4 shrink-0 text-muted-foreground" />
            <span
              class="text-sm truncate"
              :class="tool.broughtBy ? 'line-through text-muted-foreground' : 'font-medium'"
            >
              {{ tool.name }}
            </span>
          </div>
          <Badge v-if="tool.broughtBy" variant="secondary" class="shrink-0">
            {{ bringerName(tool.broughtBy) }}
          </Badge>
        </li>
      </ul>
    </CardContent>
  </Card>
</template>
