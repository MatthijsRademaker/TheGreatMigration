<script setup lang="ts">
import type { PrimitiveProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { Primitive } from "reka-ui"
import { cn } from '@/shared/lib/utils'
import { Avatar } from "@/shared/ui/avatar"

interface Props extends PrimitiveProps {
  name: string
  src?: string
  class?: HTMLAttributes["class"]
}

const props = withDefaults(defineProps<Props>(), {
  as: "span",
})
</script>

<template>
  <Primitive
    data-slot="person-chip"
    :as="as"
    :as-child="asChild"
    :class="cn('bg-secondary inline-flex items-center gap-2 rounded-full py-1 pl-1 pr-3', props.class)"
  >
    <Avatar
      :name="name"
      :src="src"
      class="size-6"
    />
    <span class="text-body truncate text-sm font-medium">
      {{ name }}
    </span>
  </Primitive>
</template>
