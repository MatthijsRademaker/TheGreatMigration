<script setup lang="ts">
import type { DateValue } from "@internationalized/date"
import type { HTMLAttributes } from "vue"
import { Calendar as CalendarIcon } from "@lucide/vue"
import { cn } from '@/shared/lib/utils'
import { Button } from "@/shared/ui/button"
import { Calendar } from "@/shared/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/shared/ui/popover"

interface Props {
  modelValue?: DateValue | null
  placeholder?: string
  class?: HTMLAttributes["class"]
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: "Select date",
})

const emits = defineEmits<{
  (e: "update:modelValue", value: DateValue | undefined): void
}>()

function formatDate(value: DateValue | null | undefined): string {
  if (!value) return props.placeholder
  return `${value.year}-${String(value.month).padStart(2, "0")}-${String(value.day).padStart(2, "0")}`
}

function onDateChange(value: DateValue | undefined) {
  emits("update:modelValue", value)
}
</script>

<template>
  <Popover>
    <PopoverTrigger as-child>
      <Button
        variant="outline"
        :class="cn(
          'justify-start gap-2 font-normal',
          !modelValue && 'text-muted-foreground',
          props.class,
        )"
        :disabled="disabled"
      >
        <CalendarIcon class="size-4" />
        {{ formatDate(modelValue) }}
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-auto p-0">
      <Calendar
        :model-value="modelValue ?? undefined"
        @update:model-value="onDateChange"
      />
    </PopoverContent>
  </Popover>
</template>
