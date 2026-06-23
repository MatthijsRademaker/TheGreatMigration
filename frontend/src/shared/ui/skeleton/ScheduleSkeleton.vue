<script setup lang="ts">
import Skeleton from './Skeleton.vue'

/** Shimmer placeholder approximating the day-columns schedule board layout. */
withDefaults(defineProps<{ columns?: number; cardsPerColumn?: number }>(), {
  columns: 4,
  cardsPerColumn: 2,
})
</script>

<template>
  <div class="overflow-x-auto" data-testid="skeleton-schedule">
    <div class="flex gap-3">
      <div v-for="c in columns" :key="c" class="min-w-[260px] shrink-0">
        <div class="mb-3 flex items-baseline gap-2">
          <Skeleton class="h-4 w-20" />
          <Skeleton class="h-3 w-16" />
        </div>
        <div class="flex flex-col gap-3">
          <Skeleton
            v-for="n in cardsPerColumn"
            :key="n"
            class="h-24 w-full rounded-lg"
          />
        </div>
      </div>
    </div>
  </div>
</template>
