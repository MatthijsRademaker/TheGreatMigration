<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { CircleHelp } from "@lucide/vue"
import { cn } from '@/shared/lib/utils'
import { Button } from "@/shared/ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/shared/ui/tooltip"

const props = defineProps<{
  tooltip: string
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <Tooltip>
    <TooltipTrigger as-child>
      <Button
        variant="ghost"
        size="icon-sm"
        :class="cn('text-muted-foreground hover:text-foreground', props.class)"
      >
        <CircleHelp class="size-4" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      {{ tooltip }}
    </TooltipContent>
  </Tooltip>
</template>
