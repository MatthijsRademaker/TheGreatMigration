<script setup lang="ts">
import { ChevronLeft, ChevronRight } from "@lucide/vue"
import { Button } from "@/shared/ui/button"

interface ToolbarProps {
  /** Formatted date range label, e.g. "2 Jul – 5 Jul, 2024". */
  rangeLabel: string
  /** Compact label for narrow viewports, e.g. "5–8 Jul". */
  compactRangeLabel: string
  /** Whether the Previous chevron should be enabled. */
  canGoPrev: boolean
  /** Whether the Next chevron should be enabled. */
  canGoNext: boolean
}

defineProps<ToolbarProps>()

interface ToolbarEmits {
  today: []
  prev: []
  next: []
}

const emit = defineEmits<ToolbarEmits>()
</script>

<template>
  <div class="flex items-center gap-2">
    <span class="hidden sm:inline text-sm text-muted-foreground whitespace-nowrap">{{ rangeLabel }}</span>
    <span class="sm:hidden text-sm text-muted-foreground whitespace-nowrap">{{ compactRangeLabel }}</span>
    <Button variant="outline" size="sm" @click="emit('today')">
      Today
    </Button>
    <Button
      variant="ghost"
      size="sm"
      :disabled="!canGoPrev"
      @click="emit('prev')"
    >
      <ChevronLeft class="size-4" />
    </Button>
    <Button
      variant="ghost"
      size="sm"
      :disabled="!canGoNext"
      @click="emit('next')"
    >
      <ChevronRight class="size-4" />
    </Button>
  </div>
</template>
