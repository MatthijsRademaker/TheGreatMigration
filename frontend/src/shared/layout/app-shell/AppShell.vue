<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppSidebar from '@/shared/layout/app-sidebar/AppSidebar.vue'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/shared/ui/sidebar'
import { usePlanningWindow } from '@/shared/composables/usePlanningWindow'

const { formattedRange, isLoading, isError } = usePlanningWindow()
</script>

<template>
  <SidebarProvider class="bg-sidebar">
    <AppSidebar />
    <SidebarInset class="overflow-y-auto">
      <header class="sticky top-0 z-10 flex h-10 shrink-0 items-center gap-2 border-b bg-background/90 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/70">
        <SidebarTrigger class="md:hidden" />
        <span class="text-sm font-semibold tracking-tight">The Great Migration</span>
        <img src="/images/birds.png" alt="" class="h-5 w-auto" />
        <span v-if="formattedRange" class="text-sm text-muted-foreground">{{ formattedRange }}</span>
        <div v-else-if="isLoading" class="h-3 w-32 animate-pulse rounded bg-muted" />
        <span v-else-if="isError" class="text-sm text-muted-foreground">&mdash;</span>
      </header>

      <RouterView />
    </SidebarInset>
  </SidebarProvider>
</template>
