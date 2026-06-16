<script setup lang="ts">
import { type Component } from 'vue'
import {
  Building2Icon,
  CalendarDaysIcon,
  CircleHelpIcon,
  ClipboardListIcon,
  HomeIcon,
  NotebookTabsIcon,
  PlusIcon,
  SettingsIcon,
  UsersRoundIcon,
} from '@lucide/vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarSeparator,
  SidebarTrigger,
  useSidebar,
} from '@/shared/ui/sidebar'

type NavigationItem = {
  title: string
  to: string
  icon: Component
}

const route = useRoute()
const router = useRouter()
const { isMobile, openMobile, setOpenMobile } = useSidebar()

// Auto-close the mobile sidebar Sheet when navigating to a new route.
router.afterEach(() => {
  if (isMobile.value && openMobile.value) {
    setOpenMobile(false)
  }
})

// TODO: Re-add badge property to navigation items when real data subscriptions exist.
// Badge counts should be driven by live backend queries (e.g. open task count,
// assigned helper count) rather than hardcoded values.
const planNavigation: NavigationItem[] = [
  { title: 'Dashboard', to: '/', icon: HomeIcon },
  { title: 'Tasks', to: '/tasks', icon: ClipboardListIcon },
  { title: 'Schedule', to: '/calendar', icon: CalendarDaysIcon },
  { title: 'People', to: '/people', icon: UsersRoundIcon },
]

const organizationNavigation: NavigationItem[] = [
  { title: 'Rooms / Areas', to: '/rooms', icon: Building2Icon },
  { title: 'Settings', to: '/settings', icon: SettingsIcon },
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
      <SidebarMenu>
        <SidebarMenuItem v-if="!isMobile">
          <SidebarMenuButton as-child variant="ghost">
            <SidebarTrigger />
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>Plan</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in planNavigation" :key="item.to">
              <SidebarMenuButton :is-active="isActive(item.to)" :tooltip="item.title" as-child>
                <RouterLink :to="item.to">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </RouterLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <SidebarSeparator />

      <SidebarGroup>
        <SidebarGroupLabel>Organization</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in organizationNavigation" :key="item.to">
              <SidebarMenuButton :is-active="isActive(item.to)" :tooltip="item.title" as-child>
                <RouterLink :to="item.to">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </RouterLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter>
      <!-- TODO: Make utility actions interactive when backend wiring is available for
           creating notes and opening help/support flows. -->
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg">
            <div class="flex aspect-square size-9 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
              <NotebookTabsIcon />
            </div>
            <div class="grid flex-1 text-left leading-tight">
              <span class="truncate text-sm font-semibold">The Great Migration</span>
              <span class="truncate text-xs text-muted-foreground">House move planner</span>
            </div>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton>
            <PlusIcon />
            <span>Add note</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton>
            <CircleHelpIcon />
            <span>Help &amp; Support</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
