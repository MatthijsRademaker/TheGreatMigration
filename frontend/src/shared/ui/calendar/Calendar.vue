<script setup lang="ts">
import type { CalendarRootEmits, CalendarRootProps } from "reka-ui"
import { ChevronLeft, ChevronRight } from "@lucide/vue"
import {
  CalendarCell,
  CalendarCellTrigger,
  CalendarGrid,
  CalendarGridBody,
  CalendarGridHead,
  CalendarGridRow,
  CalendarHeadCell,
  CalendarHeader,
  CalendarHeading,
  CalendarNext,
  CalendarPrev,
  CalendarRoot,
  useForwardPropsEmits,
} from "reka-ui"
import { cn } from '@/shared/lib/utils'

interface Props extends CalendarRootProps {
  class?: string
}

const props = defineProps<Props>()
const emits = defineEmits<CalendarRootEmits>()

const forwarded = useForwardPropsEmits(props, emits)

function formatWeekday(day: string) {
  return day.slice(0, 2)
}
</script>

<template>
  <CalendarRoot
    v-slot="{ date, grid, weekDays }"
    v-bind="forwarded"
    :class="cn('p-3', props.class)"
  >
    <CalendarHeader>
      <CalendarPrev
        class="hover:bg-muted focus-visible:ring-ring/30 inline-flex size-7 items-center justify-center rounded-md transition-colors focus-visible:ring-2 focus-visible:outline-none"
      >
        <ChevronLeft class="size-4" />
      </CalendarPrev>
      <CalendarHeading class="text-body flex-1 text-center text-sm font-medium" />
      <CalendarNext
        class="hover:bg-muted focus-visible:ring-ring/30 inline-flex size-7 items-center justify-center rounded-md transition-colors focus-visible:ring-2 focus-visible:outline-none"
      >
        <ChevronRight class="size-4" />
      </CalendarNext>
    </CalendarHeader>

    <div class="mt-3 flex flex-col gap-y-2">
      <CalendarGrid v-for="month in grid" :key="month.value.toString()">
        <CalendarGridHead>
          <CalendarGridRow>
            <CalendarHeadCell
              v-for="day in weekDays"
              :key="day"
              class="text-label w-7 rounded-md text-center text-xs font-normal"
            >
              {{ formatWeekday(day) }}
            </CalendarHeadCell>
          </CalendarGridRow>
        </CalendarGridHead>
        <CalendarGridBody>
          <CalendarGridRow v-for="(week, wi) in month.rows" :key="`${wi}`" class="mt-1">
            <CalendarCell
              v-for="cellDate in week"
              :key="cellDate.toString()"
              :date="cellDate"
            >
              <CalendarCellTrigger
                :day="cellDate"
                :month="date"
                class="data-[selected]:bg-primary data-[selected]:text-primary-foreground hover:bg-muted focus-visible:ring-ring/30 relative flex size-7 items-center justify-center rounded-md text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-offset-1 data-[today]:border data-[today]:border-border data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[unavailable]:pointer-events-none data-[unavailable]:text-destructive/30 data-[unavailable]:line-through"
              />
            </CalendarCell>
          </CalendarGridRow>
        </CalendarGridBody>
      </CalendarGrid>
    </div>
  </CalendarRoot>
</template>
