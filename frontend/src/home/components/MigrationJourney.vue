<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Motion } from 'motion-v'
import { HomeIcon, FlagIcon, TruckIcon } from '@lucide/vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import Celebration from '@/shared/motion/Celebration.vue'
import { useMotionPreference } from '@/shared/composables/useMotionPreference'
import { useMigrationReadiness } from '@/shared/composables/useMigrationReadiness'
import { springs } from '@/shared/motion/tokens'
import type { ScheduleDay } from '@/calendar/composables/useDailySchedule'

/**
 * The signature "Move-Day Readiness Journey": an on-theme progress path showing
 * the move travelling from the old place → in transit → the new place as cards
 * get staffed and tools get claimed. Read-only; derived from existing queries.
 */
const props = defineProps<{ days: ScheduleDay[] }>()

const daysRef = computed(() => props.days)
const { percent, isComplete, staffing, toolCoverage } = useMigrationReadiness(daysRef)

const { transition } = useMotionPreference()

// Spring the traveller toward its readiness position; snap instantly when reduced.
const travellerTransition = computed(() => transition(springs.smooth))

// Celebrate the arrival once, on the transition into completeness.
const arrivalTrigger = ref(0)
watch(isComplete, (complete, was) => {
  if (complete && !was) arrivalTrigger.value += 1
})

const statusLabel = computed(() => {
  if (isComplete.value) return 'Arrived — the move is ready!'
  if (percent.value === 0) return 'Not started'
  if (percent.value < 50) return 'Getting packed'
  if (percent.value < 100) return 'On the way'
  return 'Almost there'
})
</script>

<template>
  <Card class="relative overflow-hidden">
    <CardHeader>
      <div class="flex items-center justify-between gap-2">
        <CardTitle>Move-Day Readiness Journey</CardTitle>
        <span
          class="text-sm font-semibold"
          :class="isComplete ? 'text-success' : 'text-muted-foreground'"
        >
          {{ percent }}%
        </span>
      </div>
    </CardHeader>
    <CardContent>
      <Celebration :trigger="arrivalTrigger" />

      <!-- Journey path: old place → in transit → new place -->
      <div class="relative mx-1 mt-2 mb-6 h-12">
        <!-- Track -->
        <div class="absolute left-0 right-0 top-1/2 h-1.5 -translate-y-1/2 rounded-full bg-muted" />
        <!-- Progress fill -->
        <Motion
          class="absolute left-0 top-1/2 h-1.5 -translate-y-1/2 rounded-full bg-success"
          :animate="{ width: `${percent}%` }"
          :transition="travellerTransition"
          :style="{ width: `${percent}%` }"
          data-testid="journey-progress"
        />

        <!-- Old place marker -->
        <div class="absolute left-0 top-1/2 flex size-7 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full border bg-card text-muted-foreground">
          <HomeIcon class="size-4" />
        </div>
        <!-- New place marker -->
        <div
          class="absolute right-0 top-1/2 flex size-7 translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full border bg-card"
          :class="isComplete ? 'border-success text-success' : 'text-muted-foreground'"
        >
          <FlagIcon class="size-4" />
        </div>

        <!-- Traveller -->
        <Motion
          class="absolute top-1/2 z-[1] flex size-8 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-md"
          :animate="{ left: `${percent}%` }"
          :transition="travellerTransition"
          :style="{ left: `${percent}%` }"
          data-testid="journey-traveller"
        >
          <TruckIcon class="size-4" />
        </Motion>
      </div>

      <!-- Status + dimension breakdown -->
      <div class="flex flex-wrap items-center justify-between gap-2 text-sm">
        <span :class="isComplete ? 'font-medium text-success' : 'text-muted-foreground'">
          {{ statusLabel }}
        </span>
        <div class="flex items-center gap-4 text-xs text-muted-foreground">
          <span v-if="staffing.applicable">
            Staffing {{ Math.round(staffing.ratio * 100) }}%
          </span>
          <span v-if="toolCoverage.applicable">
            Tools {{ Math.round(toolCoverage.ratio * 100) }}%
          </span>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
