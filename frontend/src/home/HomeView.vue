<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { CalendarDaysIcon, ClipboardCheckIcon, MapPinnedIcon, UsersRoundIcon } from '@lucide/vue'
import { Badge } from '@/shared/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

type SummaryCard = {
  label: string
  value: string
  description: string
  icon: typeof ClipboardCheckIcon
}

const summaries: SummaryCard[] = [
  { label: 'Available today', value: '6', description: 'Helpers with confirmed availability', icon: UsersRoundIcon },
  { label: 'Under-staffed', value: '3', description: 'Tasks needing more people assigned', icon: MapPinnedIcon },
  { label: 'Move days', value: '5', description: 'Scheduled working days in the plan', icon: CalendarDaysIcon },
]

const helloMessage = ref('')
const helloLoading = ref(true)
const helloError = ref(false)

onMounted(async () => {
  try {
    const baseUrl = import.meta.env.VITE_API_BASE_URL || ''
    const res = await fetch(`${baseUrl}/api/hello`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    helloMessage.value = data.message
  } catch {
    helloError.value = true
  } finally {
    helloLoading.value = false
  }
})

const upcomingWork = [
  'Pack kitchen essentials and label fragile boxes',
  'Confirm van pickup window and parking access',
  'Clear garden tools from the shed',
]
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">

    <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <Card>
        <CardHeader class="flex flex-row items-start justify-between gap-3">
          <div class="flex flex-col gap-1">
            <CardDescription>Hello world</CardDescription>
            <CardTitle class="text-3xl">
              <span v-if="helloLoading" class="text-muted-foreground">Loading…</span>
              <span v-else-if="helloError" class="text-destructive">Backend unavailable</span>
              <span v-else>{{ helloMessage }}</span>
            </CardTitle>
          </div>
          <div class="flex size-10 items-center justify-center rounded-full bg-secondary text-secondary-foreground">
            <component :is="ClipboardCheckIcon" />
          </div>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">Live message from the backend</p>
        </CardContent>
      </Card>

      <Card v-for="summary in summaries" :key="summary.label">
        <CardHeader class="flex flex-row items-start justify-between gap-3">
          <div class="flex flex-col gap-1">
            <CardDescription>{{ summary.label }}</CardDescription>
            <CardTitle class="text-3xl">{{ summary.value }}</CardTitle>
          </div>
          <div class="flex size-10 items-center justify-center rounded-full bg-secondary text-secondary-foreground">
            <component :is="summary.icon" />
          </div>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">{{ summary.description }}</p>
        </CardContent>
      </Card>
    </div>

    <div class="grid flex-1 gap-4 xl:grid-cols-[1.4fr_0.9fr]">
      <Card>
        <CardHeader>
          <div class="flex items-start justify-between gap-4">
            <div class="flex flex-col gap-1">
              <CardTitle>Today’s plan</CardTitle>
              <CardDescription>A calm command center for the next practical steps.</CardDescription>
            </div>
            <Badge variant="secondary">Draft schedule</Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div class="grid gap-3 md:grid-cols-3">
            <div v-for="item in upcomingWork" :key="item" class="rounded-lg border bg-muted/40 p-4">
              <p class="text-sm font-medium leading-6">{{ item }}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Move notes</CardTitle>
          <CardDescription>Keep reminders visible without turning the app into heavy project management.</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="flex flex-col gap-3 text-sm text-muted-foreground">
            <p>Check building access times before assigning early morning jobs.</p>
            <p>Keep tea, chargers, tape, markers, and bin bags in the first-day essentials box.</p>
          </div>
        </CardContent>
      </Card>
    </div>
  </section>
</template>
