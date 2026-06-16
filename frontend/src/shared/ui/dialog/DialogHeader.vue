<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { cn } from '@/shared/lib/utils'

const props = defineProps<{ class?: HTMLAttributes["class"] }>()
</script>

<template>
  <div
    data-slot="dialog-header"
    :class="cn('gap-1.5 p-6 pb-4 flex flex-col', props.class)"
  >
    <slot />
  </div>
</template>
