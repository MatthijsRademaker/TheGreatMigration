<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { cn } from '@/shared/lib/utils'

const props = defineProps<{ class?: HTMLAttributes["class"] }>()
</script>

<template>
  <div
    data-slot="dialog-footer"
    :class="cn('gap-2 p-6 pt-4 flex flex-col sm:flex-row sm:justify-end', props.class)"
  >
    <slot />
  </div>
</template>
