<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import { Motion } from 'motion-v'
import { useMotionPreference } from '@/shared/composables/useMotionPreference'

/**
 * A one-shot particle burst overlay. Increment `trigger` to replay the burst
 * (e.g. when a day reaches 100% readiness or the move arrives). Renders nothing
 * under reduced motion. Position it inside a `relative` container.
 */
const props = defineProps<{ trigger: number }>()

const { enabled } = useMotionPreference()

const PARTICLE_COUNT = 12
const DURATION_MS = 800
const colors = ['var(--success)', 'var(--primary)', 'var(--warning)', 'var(--info)']

const particles = Array.from({ length: PARTICLE_COUNT }, (_, i) => {
  const angle = (i / PARTICLE_COUNT) * Math.PI * 2
  return {
    id: i,
    x: Math.cos(angle) * 44,
    y: Math.sin(angle) * 44,
    color: colors[i % colors.length],
  }
})

const active = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined

watch(
  () => props.trigger,
  (next, prev) => {
    if (!enabled.value || next === prev) return
    active.value = true
    clearTimeout(timer)
    timer = setTimeout(() => {
      active.value = false
    }, DURATION_MS)
  },
)

onBeforeUnmount(() => clearTimeout(timer))
</script>

<template>
  <div
    v-if="active"
    data-testid="celebration"
    class="pointer-events-none absolute inset-0 z-10 overflow-visible"
  >
    <Motion
      v-for="p in particles"
      :key="p.id"
      as="span"
      class="absolute left-1/2 top-1/2 size-1.5 rounded-full"
      :style="{ backgroundColor: p.color }"
      :initial="{ opacity: 1, x: 0, y: 0, scale: 1 }"
      :animate="{ opacity: 0, x: p.x, y: p.y, scale: 0.5 }"
      :transition="{ duration: 0.7 }"
    />
  </div>
</template>
