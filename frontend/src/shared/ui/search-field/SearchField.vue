<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { Search } from "@lucide/vue"
import { useVModel } from "@vueuse/core"
import { cn } from '@/shared/lib/utils'
import { Input } from "@/shared/ui/input"

defineOptions({
  inheritAttrs: false,
})

const props = withDefaults(defineProps<{
  modelValue?: string | number
  defaultValue?: string | number
  placeholder?: string
  disabled?: boolean
  class?: HTMLAttributes["class"]
}>(), {
  placeholder: "Search...",
})

const emits = defineEmits<{
  (e: "update:modelValue", payload: string | number): void
}>()

const modelValue = useVModel(props, "modelValue", emits, {
  passive: true,
  defaultValue: props.defaultValue,
})
</script>

<template>
  <div :class="cn('relative', props.class)">
    <Search class="text-muted-foreground pointer-events-none absolute left-2 top-1/2 size-3.5 -translate-y-1/2" />
    <Input
      v-model="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      class="pl-7"
      v-bind="$attrs"
    />
  </div>
</template>
