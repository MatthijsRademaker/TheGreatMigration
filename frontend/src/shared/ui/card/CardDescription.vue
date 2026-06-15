<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { cn } from '@/shared/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <div
    data-slot="card-description"
    :class="cn('[font-size:var(--text-caption)] [line-height:var(--text-caption--line-height)] text-muted-foreground', props.class)"
  >
    <slot />
  </div>
</template>
