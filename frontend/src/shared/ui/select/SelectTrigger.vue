<script setup lang="ts">
import type { SelectTriggerProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { ChevronDown } from "@lucide/vue"
import { SelectIcon, SelectTrigger } from "reka-ui"
import { cn } from '@/shared/lib/utils'

const props = defineProps<SelectTriggerProps & {
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <SelectTrigger
    data-slot="select-trigger"
    v-bind="props"
    :class="cn(
      'bg-background border-border text-foreground h-7 rounded-md border px-2 py-0.5 text-sm transition-colors focus-visible:border-ring focus-visible:ring-ring/30 focus-visible:ring-2 focus-visible:outline-none disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 outline-none flex w-full items-center justify-between gap-1 whitespace-nowrap [&>span]:truncate',
      props.class,
    )"
  >
    <slot />
    <SelectIcon>
      <ChevronDown class="size-4 shrink-0 text-muted-foreground" />
    </SelectIcon>
  </SelectTrigger>
</template>
