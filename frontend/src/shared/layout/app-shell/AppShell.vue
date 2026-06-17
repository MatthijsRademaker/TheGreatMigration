<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppSidebar from '@/shared/layout/app-sidebar/AppSidebar.vue'
import { Button } from '@/shared/ui/button'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/shared/ui/sidebar'
import { usePlanningWindow } from '@/shared/composables/usePlanningWindow'

const route = useRoute()
const { formattedRange, isLoading, isError } = usePlanningWindow()

const routeTitle = computed(() => typeof route.meta.title === 'string' ? route.meta.title : '')
const routeDescription = computed(() => typeof route.meta.description === 'string' ? route.meta.description : '')
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

      <section class="border-b bg-background px-4 py-4 sm:px-6">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <p class="text-xs font-medium uppercase tracking-wide text-muted-foreground">House move planner</p>
            <h1 class="font-heading text-2xl font-semibold tracking-tight">{{ routeTitle }}</h1>
            <p v-if="routeDescription" class="mt-1 text-sm text-muted-foreground">{{ routeDescription }}</p>
          </div>
          <div class="flex items-center gap-2">
            <Button variant="outline" size="sm">Add note</Button>
            <Button variant="outline" size="sm">Help &amp; Support</Button>
          </div>
        </div>
      </section>

      <RouterView />
    </SidebarInset>
  </SidebarProvider>
</template>
