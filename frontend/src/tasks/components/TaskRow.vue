<script setup lang="ts">
import { UsersRoundIcon } from '@lucide/vue'
import type { TaskRow } from '../types'
import { getTaskDisplayState, priorityBadgeVariant, priorityLabels } from '../helpers'
import { Badge } from '@/shared/ui/badge'
import { areaColor } from '@/shared/lib/areaColor'

defineProps<{
  task: TaskRow
}>()
</script>

<template>
  <tr class="border-b border-border last:border-b-0">
    <td class="h-12 px-4 py-2 text-sm font-medium">{{ task.title }}</td>
    <td class="h-12 px-4 py-2">
      <Badge :variant="priorityBadgeVariant[task.priority]">
        {{ priorityLabels[task.priority] }}
      </Badge>
    </td>
    <td class="h-12 px-4 py-2">
      <span class="inline-flex items-center gap-1.5 text-sm text-muted-foreground">
        <UsersRoundIcon class="size-3.5" />
        {{ task.peopleNeeded }}
      </span>
    </td>
    <td class="h-12 px-4 py-2 text-sm text-muted-foreground">
      <span class="inline-flex items-center gap-1.5" data-testid="area-chip">
        <span
          v-if="task.area.name"
          class="inline-block size-2 shrink-0 rounded-full"
          :style="{ backgroundColor: areaColor(task.area.id) }"
          aria-hidden="true"
        />
        {{ task.area.name }}
      </span>
    </td>
    <td class="h-12 px-4 py-2">
      <Badge variant="secondary">{{ getTaskDisplayState(task) }}</Badge>
    </td>
    <td v-if="$slots.actions" class="h-12 px-4 py-2">
      <div class="flex items-center gap-1">
        <slot name="actions" />
      </div>
    </td>
  </tr>
</template>
