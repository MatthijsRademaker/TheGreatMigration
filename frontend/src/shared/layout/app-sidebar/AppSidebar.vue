<script setup lang="ts">
import type { Component } from 'vue'
import { CalendarDaysIcon, ClipboardListIcon, HomeIcon, NotebookTabsIcon, UsersRoundIcon } from '@lucide/vue'
import { RouterLink, useRoute } from 'vue-router'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from '@/shared/ui/sidebar'

type NavigationItem = {
  title: string
  to: string
  icon: Component
  badge?: string
}

const route = useRoute()

const primaryNavigation: NavigationItem[] = [
  { title: 'Dashboard', to: '/', icon: HomeIcon },
  { title: 'Tasks', to: '/tasks', icon: ClipboardListIcon, badge: '12' },
  { title: 'Schedule', to: '/calendar', icon: CalendarDaysIcon },
  { title: 'People', to: '/people', icon: UsersRoundIcon, badge: '6' },
]

function isActive(path: string) {
  return route.path === path
}
</script>

<template>
  <Sidebar collapsible="icon" variant="inset">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <RouterLink to="/" aria-label="Open moving dashboard">
              <div class="flex aspect-square size-9 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <NotebookTabsIcon />
              </div>
              <div class="grid flex-1 text-left leading-tight">
                <span class="truncate text-sm font-semibold">The Great Migration</span>
                <span class="truncate text-xs text-muted-foreground">House move planner</span>
              </div>
            </RouterLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>Plan</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in primaryNavigation" :key="item.to">
              <SidebarMenuButton :is-active="isActive(item.to)" :tooltip="item.title" as-child>
                <RouterLink :to="item.to">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </RouterLink>
              </SidebarMenuButton>
              <SidebarMenuBadge v-if="item.badge">{{ item.badge }}</SidebarMenuBadge>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

    </SidebarContent>

    <SidebarFooter>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton>
            <span class="flex size-6 items-center justify-center rounded-md bg-sidebar-accent text-xs font-semibold">GM</span>
            <span>The Great Migration</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
