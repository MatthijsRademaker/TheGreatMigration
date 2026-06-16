<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppSidebar from '@/shared/layout/app-sidebar/AppSidebar.vue'
import { SidebarInset, SidebarProvider } from '@/shared/ui/sidebar'
import { usePlanningWindow } from '@/shared/composables/usePlanningWindow'

const route = useRoute()

const pageTitle = computed(() => String(route.meta.title ?? 'The Great Migration'))
const pageDescription = computed(() => String(route.meta.description ?? 'A shared plan for the practical work of moving house.'))

const { formattedRange, isLoading, isError } = usePlanningWindow()
</script>

<template>
  <SidebarProvider class="bg-sidebar">
    <AppSidebar />
    <SidebarInset class="overflow-hidden">
      <header class="flex min-h-16 shrink-0 items-center gap-3 border-b bg-background/90 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/70">
        <div class="flex min-w-0 flex-1 flex-col gap-1 py-3">
          <div class="flex items-center gap-2">
            <div class="flex h-6 items-center">
              <span v-if="formattedRange" class="text-sm text-muted-foreground">{{ formattedRange }}</span>
              <div v-else-if="isLoading" class="h-4 w-48 animate-pulse rounded bg-muted" />
              <span v-else-if="isError" class="text-sm text-muted-foreground">&mdash;</span>
            </div>
            <h1 class="truncate text-base font-semibold tracking-tight sm:text-lg">{{ pageTitle }}</h1>
          </div>
          <p class="truncate text-sm text-muted-foreground">{{ pageDescription }}</p>
        </div>
      </header>

      <RouterView />
    </SidebarInset>
  </SidebarProvider>
</template>
