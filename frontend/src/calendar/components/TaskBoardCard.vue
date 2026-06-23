<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Motion, AnimatePresence } from 'motion-v'
import { CheckIcon } from '@lucide/vue'
import { Badge } from '@/shared/ui/badge'
import type { BadgeVariants } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import Celebration from '@/shared/motion/Celebration.vue'
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

const props = withDefaults(defineProps<{
  task: BoardTaskCard
  readOnly?: boolean
  /** Whether the card participates in board drag affordances (hover lift). */
  interactive?: boolean
  /** True while a draggable person is hovering this card as a drop target. */
  dropActive?: boolean
  /** Render the card as completed/greyed-out. */
  done?: boolean
}>(), {
  done: false,
})

const emit = defineEmits<{
  done: []
  revert: []
  edit: []
  delete: []
}>()

const { enabled } = useMotionPreference()
const doneCelebrationTrigger = ref(0)
let doneTimer: ReturnType<typeof setTimeout> | undefined

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
  props.interactive && !props.done && enabled.value ? { y: -3, scale: 1.02 } : undefined,
)

function handleDone() {
  if (doneTimer) clearTimeout(doneTimer)

  if (enabled.value) {
    doneCelebrationTrigger.value += 1
    doneTimer = setTimeout(() => {
      emit('done')
      doneTimer = undefined
    }, 120)
    return
  }

  emit('done')
}

onBeforeUnmount(() => {
  clearTimeout(doneTimer)
})
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
      done ? 'border-l-border border-border bg-muted/40 opacity-60 shadow-none grayscale' : priorityAccentMap[task.priority],
      dropActive && !done ? 'ring-2 ring-primary ring-offset-1' : '',
    ]"
    :data-done="done ? 'true' : 'false'"
  >
    <Celebration :trigger="doneCelebrationTrigger" />

    <div class="flex items-start justify-between gap-2 mb-2">
      <span class="text-sm font-medium" :class="done ? 'line-through text-muted-foreground' : ''">
        {{ task.title }}
      </span>
      <Badge :variant="priorityVariantMap[task.priority]">{{ task.priority }}</Badge>
    </div>
    <p
      v-if="task.assignedPeople.length > 0"
      class="text-xs text-muted-foreground mb-2"
      :class="done ? 'line-through' : ''"
    >
      {{ task.assignedPeople.map((person) => person.name).join(', ') }}
    </p>
    <p class="flex items-center gap-1 text-xs text-muted-foreground" :class="done ? 'line-through' : ''">
      {{ task.assignedCount }} / {{ task.peopleNeeded }}
      <span v-if="!done && task.staffingStatus === 'underStaffed'" class="text-destructive">
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
      <template v-if="done">
        <Button variant="ghost" size="xs" data-testid="revert-button" @click.stop="emit('revert')">
          ↩ Revert
        </Button>
      </template>
      <template v-else>
        <Button variant="ghost" size="xs" data-testid="done-button" @click.stop="handleDone">
          <CheckIcon class="size-3" />
          Done
        </Button>
        <Button variant="ghost" size="xs" @click.stop="emit('edit')">Edit</Button>
        <Button variant="ghost" size="xs" @click.stop="emit('delete')">Delete</Button>
      </template>
    </div>
  </Motion>
</template>
