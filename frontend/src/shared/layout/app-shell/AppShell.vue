<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppSidebar from '@/shared/layout/app-sidebar/AppSidebar.vue'
import { Badge } from '@/shared/ui/badge'
import { Separator } from '@/shared/ui/separator'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/shared/ui/sidebar'

const route = useRoute()

const pageTitle = computed(() => String(route.meta.title ?? 'The Great Migration'))
const pageDescription = computed(() => String(route.meta.description ?? 'A shared plan for the practical work of moving house.'))
</script>

<template>
  <SidebarProvider class="bg-sidebar">
    <AppSidebar />
    <SidebarInset class="overflow-hidden">
      <header class="flex min-h-16 shrink-0 items-center gap-3 border-b bg-background/90 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/70">
        <SidebarTrigger class="-ml-1" />
        <Separator class="min-h-6" orientation="vertical" />
        <div class="flex min-w-0 flex-1 flex-col gap-1 py-3">
          <div class="flex items-center gap-2">
            <Badge class="hidden sm:inline-flex" variant="secondary">Planning mode</Badge>
            <h1 class="truncate text-base font-semibold tracking-tight sm:text-lg">{{ pageTitle }}</h1>
          </div>
          <p class="truncate text-sm text-muted-foreground">{{ pageDescription }}</p>
        </div>
      </header>

      <RouterView />
    </SidebarInset>
  </SidebarProvider>
</template>
