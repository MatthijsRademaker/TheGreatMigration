<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { cn } from '@/shared/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <div
    data-slot="card-title"
    :class="cn('font-heading text-card-title font-semibold text-card-foreground', props.class)"
  >
    <slot />
  </div>
</template>
