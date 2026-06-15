<script setup lang="ts">
import type { SelectItemEmits, SelectItemProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { SelectItemIndicator, SelectItem, SelectItemText } from "reka-ui"
import { Check } from "@lucide/vue"
import { cn } from '@/shared/lib/utils'

defineOptions({
  inheritAttrs: false,
})

const props = defineProps<SelectItemProps & { class?: HTMLAttributes["class"] }>()
const emits = defineEmits<SelectItemEmits>()
</script>

<template>
  <SelectItem
    :class="cn(
      'focus:bg-accent focus:text-accent-foreground [&_svg]:size-4 relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
      props.class,
    )"
    v-bind="props"
    @select="emits('select', $event)"
  >
    <span class="absolute left-2 flex size-3.5 items-center justify-center">
      <SelectItemIndicator>
        <Check class="size-4" />
      </SelectItemIndicator>
    </span>
    <SelectItemText>
      <slot />
    </SelectItemText>
  </SelectItem>
</template>
