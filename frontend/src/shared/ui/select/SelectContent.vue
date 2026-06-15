<script setup lang="ts">
import type { SelectContentEmits, SelectContentProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { SelectContent, SelectPortal, SelectViewport, useForwardPropsEmits } from "reka-ui"
import { cn } from '@/shared/lib/utils'

defineOptions({
  inheritAttrs: false,
})

const props = withDefaults(defineProps<SelectContentProps & { class?: HTMLAttributes["class"] }>(), {
  position: "popper",
  sideOffset: 4,
})

const emits = defineEmits<SelectContentEmits>()

const delegatedProps = reactiveOmit(props, "class")
const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <SelectPortal>
    <SelectContent
      data-slot="select-content"
      v-bind="{ ...forwarded, ...$attrs }"
      :class="cn(
        'bg-popover text-popover-foreground border-border shadow-md z-50 min-w-[8rem] overflow-hidden rounded-md border',
        'data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95',
        'data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2',
        'data-[position=popper]:data-[side=bottom]:translate-y-1 data-[position=popper]:data-[side=left]:-translate-x-1 data-[position=popper]:data-[side=right]:translate-x-1 data-[position=popper]:data-[side=top]:-translate-y-1',
        props.class,
      )"
    >
      <SelectViewport
        class="p-1"
      >
        <slot />
      </SelectViewport>
    </SelectContent>
  </SelectPortal>
</template>
