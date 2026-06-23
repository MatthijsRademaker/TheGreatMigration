<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Motion, AnimatePresence } from 'motion-v'
import { CheckIcon } from '@lucide/vue'
import { Badge } from '@/shared/ui/badge'
import type { BadgeVariants } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import { useMotionPreference } from '@/shared/composables/useMotionPreference'
import { springs } from '@/shared/motion/tokens'

interface AssignedPerson {
  id: string
  name: string
  initials: string
}

interface BoardTaskCard {
  id: string
  title: string
  priority: 'high' | 'medium' | 'low'
  assignedPeople: AssignedPerson[]
  peopleNeeded: number
  assignedCount: number
  staffingStatus: 'fullyStaffed' | 'underStaffed'
}

const props = defineProps<{
  task: BoardTaskCard
  readOnly?: boolean
  /** Whether the card participates in board drag affordances (hover lift). */
  interactive?: boolean
  /** True while a draggable person is hovering this card as a drop target. */
  dropActive?: boolean
}>()

const emit = defineEmits<{
  edit: []
  delete: []
}>()

const { enabled } = useMotionPreference()

const priorityAccentMap: Record<BoardTaskCard['priority'], string> = {
  high: 'border-l-destructive',
  medium: 'border-l-warning',
  low: 'border-l-success',
}

const priorityVariantMap: Record<BoardTaskCard['priority'], BadgeVariants['variant']> = {
  high: 'priorityHigh',
  medium: 'priorityMedium',
  low: 'priorityLow',
}

// Reward: a one-shot pop when staffing reaches fullyStaffed. The animate prop
// briefly drives a scale keyframe; under reduced motion it stays at rest.
const popState = ref<{ scale: number | number[] }>({ scale: 1 })

watch(
  () => props.task.staffingStatus,
  (status, previous) => {
    if (enabled.value && status === 'fullyStaffed' && previous === 'underStaffed') {
      popState.value = { scale: [1, 1.06, 1] }
    }
  },
)

// Spring hover lift, only for interactive (board) cards with motion enabled.
const whileHover = computed(() =>
  props.interactive && enabled.value ? { y: -3, scale: 1.02 } : undefined,
)
</script>

<template>
  <Motion
    as="div"
    layout
    :animate="popState"
    :while-hover="whileHover"
    :transition="springs.bouncy"
    data-slot="task-board-card"
    class="relative rounded-lg border bg-card shadow-sm p-3 border-l-4 transition-colors"
    :class="[
      priorityAccentMap[task.priority],
      dropActive ? 'ring-2 ring-primary ring-offset-1' : '',
    ]"
  >
    <div class="flex items-start justify-between gap-2 mb-2">
      <span class="text-sm font-medium">{{ task.title }}</span>
      <Badge :variant="priorityVariantMap[task.priority]">{{ task.priority }}</Badge>
    </div>
    <p
      v-if="task.assignedPeople.length > 0"
      class="text-xs text-muted-foreground mb-2"
    >
      {{ task.assignedPeople.map((person) => person.name).join(', ') }}
    </p>
    <p class="flex items-center gap-1 text-xs text-muted-foreground">
      {{ task.assignedCount }} / {{ task.peopleNeeded }}
      <span v-if="task.staffingStatus === 'underStaffed'" class="text-destructive">
        &nbsp;— needs help
      </span>
      <AnimatePresence>
        <Motion
          v-if="task.staffingStatus === 'fullyStaffed'"
          as="span"
          data-testid="fully-staffed-check"
          class="inline-flex size-4 items-center justify-center rounded-full bg-success text-success-foreground"
          :initial="enabled ? { scale: 0, rotate: -45 } : false"
          :animate="{ scale: 1, rotate: 0 }"
          :exit="enabled ? { scale: 0 } : undefined"
          :transition="springs.bouncy"
        >
          <CheckIcon class="size-3" />
        </Motion>
      </AnimatePresence>
    </p>
    <div v-if="!readOnly" class="flex items-center gap-1 mt-2">
      <Button variant="ghost" size="xs" @click.stop="emit('edit')">Edit</Button>
      <Button variant="ghost" size="xs" @click.stop="emit('delete')">Delete</Button>
    </div>
  </Motion>
</template>
