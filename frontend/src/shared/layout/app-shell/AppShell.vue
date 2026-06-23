<script setup lang="ts">
import { computed } from 'vue';
import { RouterView, useRoute } from 'vue-router';
import AppSidebar from '@/shared/layout/app-sidebar/AppSidebar.vue';
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/shared/ui/sidebar';
import { usePlanningWindow } from '@/shared/composables/usePlanningWindow';
import { useHomePagination } from '@/shared/composables/useHomePagination';
import { useMotionPreference } from '@/shared/composables/useMotionPreference';
import AppShellTimelineToolbar from './AppShellTimelineToolbar.vue';

const route = useRoute();
const isHomeRoute = computed(() => route.path === '/');

const { enabled: motionEnabled } = useMotionPreference();
// Empty transition name → no transition classes → instant route change under
// reduced motion. The route-fade CSS additionally short-circuits via media query.
const routeTransition = computed(() =>
  motionEnabled.value ? 'route-fade' : '',
);

const { formattedRange, isLoading, isError } = usePlanningWindow();
const homePagination = useHomePagination();
</script>

<template>
  <SidebarProvider class="bg-sidebar">
    <AppSidebar />
    <SidebarInset class="overflow-y-auto">
      <header
        class="sticky top-0 z-10 flex h-10 shrink-0 items-center gap-2 border-b bg-background/90 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/70"
      >
        <SidebarTrigger class="md:hidden" />
        <span
          class="hidden sm:inline-block text-sm pr-1 font-semibold tracking-tight italic bg-gradient-to-r from-primary to-chart-3 bg-clip-text text-transparent"
          >The Great Migration</span
        >
        <img src="/images/birds.png" alt="" class="h-5 w-auto" />

        <div class="ml-auto flex items-center gap-2">
          <!-- Homepage: centralized timeline toolbar -->
          <template v-if="isHomeRoute">
            <div
              v-if="homePagination.isLoading.value"
              class="h-3 w-32 animate-pulse rounded bg-muted"
            />
            <span
              v-else-if="homePagination.isError.value"
              class="text-sm text-muted-foreground"
              >&mdash;</span
            >
            <AppShellTimelineToolbar
              v-else
              :range-label="homePagination.rangeLabel.value"
              :compact-range-label="homePagination.compactRangeLabel.value"
              :can-go-prev="homePagination.page.value > 1"
              :can-go-next="
                homePagination.page.value < homePagination.totalPages.value
              "
              @today="homePagination.goToday"
              @prev="homePagination.goPrev"
              @next="homePagination.goNext"
            />
          </template>

          <!-- Other routes: static planning window range -->
          <template v-else>
            <span
              v-if="formattedRange"
              class="text-sm text-muted-foreground whitespace-nowrap"
              >{{ formattedRange }}</span
            >
            <div
              v-else-if="isLoading"
              class="h-3 w-32 animate-pulse rounded bg-muted"
            />
            <span v-else-if="isError" class="text-sm text-muted-foreground"
              >&mdash;</span
            >
          </template>
        </div>
      </header>

      <RouterView v-slot="{ Component }">
        <Transition :name="routeTransition" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </SidebarInset>
  </SidebarProvider>
</template>
