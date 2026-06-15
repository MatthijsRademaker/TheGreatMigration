<script setup lang="ts">
import type { CheckboxRootEmits, CheckboxRootProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import type { CheckboxVariants } from "."
import { Check } from "@lucide/vue"
import { CheckboxIndicator, CheckboxRoot, useForwardPropsEmits } from "reka-ui"
import { reactiveOmit } from "@vueuse/core"
import { cn } from '@/shared/lib/utils'
import { checkboxVariants } from "."

interface Props extends CheckboxRootProps {
  variant?: CheckboxVariants["variant"]
  size?: CheckboxVariants["size"]
  class?: HTMLAttributes["class"]
}

const props = withDefaults(defineProps<Props>(), {
  variant: "default",
  size: "default",
})

const emits = defineEmits<CheckboxRootEmits>()

const delegatedProps = reactiveOmit(props, "class", "variant", "size")
const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <CheckboxRoot
    data-slot="checkbox"
    :data-variant="variant"
    :data-size="size"
    :class="cn(checkboxVariants({ variant, size }), props.class)"
    v-bind="forwarded"
  >
    <CheckboxIndicator
      class="flex size-full items-center justify-center text-current"
    >
      <Check class="size-3.5" />
    </CheckboxIndicator>
  </CheckboxRoot>
</template>
