<script setup lang="ts">
import { XIcon } from '@lucide/vue';

import type { DialogContentEmits, DialogContentProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import {
  DialogClose,
  DialogContent,
  DialogPortal,
  useForwardPropsEmits,
} from "reka-ui"
import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui/button'
import DialogOverlay from "./DialogOverlay.vue"

interface DialogContentPropsExtended extends DialogContentProps {
  class?: HTMLAttributes["class"]
  showCloseButton?: boolean
}

defineOptions({
  inheritAttrs: false,
})

const props = withDefaults(defineProps<DialogContentPropsExtended>(), {
  showCloseButton: true,
})
const emits = defineEmits<DialogContentEmits>()

const delegatedProps = reactiveOmit(props, "class", "showCloseButton")

const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <DialogPortal>
    <DialogOverlay />
    <DialogContent
      data-slot="dialog-content"
      :class="cn(
        'bg-popover text-popover-foreground fixed top-1/2 left-1/2 z-50 -translate-x-1/2 -translate-y-1/2 flex flex-col rounded-xl shadow-lg border border-border bg-clip-padding text-xs/relaxed',
        'w-[calc(100%-2rem)] h-[calc(100%-2rem)] max-w-[90vw] max-h-[85vh]',
        'sm:w-[60vw] sm:h-[60vh] sm:max-w-[90vw] sm:max-h-[85vh]',
        'data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95',
        'data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95',
        'duration-200',
        props.class
      )"
      v-bind="{ ...$attrs, ...forwarded }"
    >
      <slot />

      <DialogClose
        v-if="showCloseButton"
        data-slot="dialog-close"
        as-child
      >
        <Button variant="ghost" class="absolute top-4 right-4" size="icon-sm">
          <XIcon />
          <span class="sr-only">Close</span>
        </Button>
      </DialogClose>
    </DialogContent>
  </DialogPortal>
</template>
