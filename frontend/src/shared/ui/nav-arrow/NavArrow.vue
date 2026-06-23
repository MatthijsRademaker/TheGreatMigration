<script setup lang="ts">
import { ArrowUpRightIcon } from '@lucide/vue'
import { useRouter } from 'vue-router'
import { cn } from '@/shared/lib/utils'

const props = withDefaults(defineProps<{
  to: string
  class?: string
}>(), {})

const router = useRouter()
</script>

<template>
  <button
    :class="cn(
      'flex items-center justify-center size-6 rounded-full',
      'border border-border/40 bg-background/80 shadow-xs',
      'text-muted-foreground/50',
      'hover:text-muted-foreground hover:bg-accent hover:border-border',
      'cursor-pointer',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
      props.class,
    )"
    title="Navigate to full page"
    @click.stop="router.push(to)"
  >
    <ArrowUpRightIcon class="size-3.5" />
  </button>
</template>
