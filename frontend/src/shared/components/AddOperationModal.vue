<script setup lang="ts">
import { computed } from "vue"
import { Button } from "@/shared/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/dialog"

export interface AddOperationModalProps {
  open: boolean
  title?: string
  description?: string
  submitLabel?: string
  cancelLabel?: string
  disabled?: boolean
  submitting?: boolean
}

const props = withDefaults(defineProps<AddOperationModalProps>(), {
  title: "",
  description: "",
  submitLabel: "Save",
  cancelLabel: "Cancel",
  disabled: false,
  submitting: false,
})

const emit = defineEmits<{
  "update:open": [value: boolean]
  submit: []
  cancel: []
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit("update:open", value),
})

function handleSubmit() {
  emit("submit")
}

function handleCancel() {
  emit("cancel")
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent class="flex flex-col">
      <DialogHeader>
        <DialogTitle v-if="title">{{ title }}</DialogTitle>
        <DialogDescription v-if="description">
          {{ description }}
        </DialogDescription>
      </DialogHeader>

      <!-- Body region – scrollable for long forms -->
      <div class="flex-1 overflow-y-auto px-6 py-2">
        <slot />
      </div>

      <!-- Footer: caller can override via named slot, otherwise default actions -->
      <DialogFooter v-if="$slots.footer">
        <slot name="footer" />
      </DialogFooter>
      <DialogFooter v-else>
        <slot name="footer-actions">
          <Button
            type="submit"
            :disabled="disabled || submitting"
            @click="handleSubmit"
          >
            {{ submitLabel }}
          </Button>
          <Button
            variant="outline"
            type="button"
            :disabled="submitting"
            @click="handleCancel"
          >
            {{ cancelLabel }}
          </Button>
        </slot>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
