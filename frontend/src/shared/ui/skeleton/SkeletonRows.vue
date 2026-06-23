<script setup lang="ts">
import Skeleton from './Skeleton.vue'

/** Shimmer rows that approximate a loading list/table panel. */
const props = withDefaults(defineProps<{ rows?: number }>(), { rows: 5 })
</script>

<template>
  <div class="flex flex-col gap-3 py-2" data-testid="skeleton-rows">
    <div
      v-for="n in props.rows"
      :key="n"
      class="flex items-center gap-3"
    >
      <Skeleton class="size-8 shrink-0 rounded-lg" />
      <Skeleton class="h-4 flex-1" />
      <Skeleton class="h-4 w-16 shrink-0" />
    </div>
  </div>
</template>
