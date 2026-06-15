<script setup lang="ts">
import { computed } from "vue"
import type { AvatarImageEmits, AvatarImageProps, AvatarRootProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { AvatarFallback, AvatarImage, AvatarRoot, useForwardPropsEmits } from "reka-ui"
import { reactiveOmit } from "@vueuse/core"
import { cn } from '@/shared/lib/utils'

interface Props extends AvatarRootProps {
  src?: AvatarImageProps["src"]
  name?: string
  class?: HTMLAttributes["class"]
}

const props = withDefaults(defineProps<Props>(), {
  name: "",
})

const emits = defineEmits<AvatarImageEmits>()

const delegatedProps = reactiveOmit(props, "class", "src", "name")
const forwarded = useForwardPropsEmits(delegatedProps, emits)

const initials = computed(() => props.name.charAt(0).toUpperCase())
</script>

<template>
  <AvatarRoot
    data-slot="avatar"
    :class="cn('relative flex size-10 shrink-0 overflow-hidden rounded-full', props.class)"
    v-bind="forwarded"
  >
    <AvatarImage
      v-if="src"
      :src="src"
      :alt="name"
      class="aspect-square size-full"
    />
    <AvatarFallback
      :class="cn('bg-muted text-label flex size-full items-center justify-center rounded-full text-sm font-medium')"
    >
      {{ initials || '?' }}
    </AvatarFallback>
  </AvatarRoot>
</template>
