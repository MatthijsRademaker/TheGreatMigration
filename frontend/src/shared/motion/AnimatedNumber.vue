<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useTransition } from '@vueuse/core'
import { useMotionPreference } from '@/shared/composables/useMotionPreference'
import { countUp } from '@/shared/motion/tokens'

/**
 * Renders an integer that springs up to its value with a slight overshoot when
 * motion is enabled, and shows the exact value immediately under reduced motion.
 * The intro count-up runs from a baseline of 0 on mount; later value changes
 * animate from the previous value.
 */
const props = defineProps<{ value: number }>()

const { enabled } = useMotionPreference()

// Baseline 0 so the client intro animation counts up from zero.
const source = ref(0)
const output = useTransition(source, {
  duration: countUp.duration,
  transition: [...countUp.transition],
})

// False during SSR and the first client render; the count-up only runs after
// mount. Until then (and under reduced motion) the exact value is shown.
const mounted = ref(false)

onMounted(() => {
  mounted.value = true
  source.value = props.value
})

watch(
  () => props.value,
  (next) => {
    source.value = next
  },
)

const display = computed(() =>
  mounted.value && enabled.value
    ? Math.round(output.value)
    : Math.round(props.value),
)
</script>

<template>{{ display }}</template>
