<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { cn } from '@/shared/lib/utils'

const props = withDefaults(defineProps<{
  class?: HTMLAttributes["class"]
  size?: "default" | "sm"
}>(), {
  size: "default",
})
</script>

<template>
  <div
    data-slot="card"
    :data-size="size"
    :class="cn('border-border/80 bg-card text-card-foreground gap-panel overflow-hidden rounded-lg border py-panel [font-size:var(--text-body)] [line-height:var(--text-body--line-height)] shadow-sm has-[>img:first-child]:pt-0 data-[size=sm]:gap-comfy data-[size=sm]:py-comfy *:[img:first-child]:rounded-t-lg *:[img:last-child]:rounded-b-lg group/card flex flex-col', props.class)"
  >
    <slot />
  </div>
</template>
