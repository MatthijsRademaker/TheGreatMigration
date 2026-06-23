<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/shared/lib/utils'

/**
 * A shimmer placeholder block. The sweep is CSS-driven (`.motion-shimmer`) and
 * gated by `prefers-reduced-motion`: under reduced motion it renders as a
 * static muted block with no animation. Compose several to approximate the
 * eventual content layout while a panel's query is in flight.
 */
const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <div data-slot="skeleton" :class="cn('motion-shimmer rounded-md', props.class)" />
</template>
