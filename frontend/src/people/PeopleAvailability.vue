<script setup lang="ts">
import { Badge } from '@/shared/ui/badge'
import { Avatar } from '@/shared/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

/** Badge variants used for availability status legend entries. */
type StatusVariant = 'available' | 'busy' | 'partial' | 'off'

export interface PersonAvailabilityEntry {
  date: string
  status: 'available' | 'busy' | 'partial' | 'off'
}

export interface PersonAvailability {
  id: string
  name: string
  initials: string
  availability: PersonAvailabilityEntry[]
}

export interface StatusLegendItem {
  id: string
  label: string
}

export interface PeopleAvailabilityProps {
  title?: string
  description?: string
  days?: string[]
  people?: PersonAvailability[]
  legend?: StatusLegendItem[]
  availableToday?: number
  totalPeople?: number
}

withDefaults(defineProps<PeopleAvailabilityProps>(), {
  title: 'People availability',
  description: 'Track who is available and where each person can help.',
  days: () => ['Mon', 'Tue', 'Wed', 'Thu'],
  people: () => [
    {
      id: 'alex',
      name: 'Alex',
      initials: 'A',
      availability: [
        { date: 'Mon', status: 'available' },
        { date: 'Tue', status: 'available' },
        { date: 'Wed', status: 'busy' },
        { date: 'Thu', status: 'available' },
      ],
    },
    {
      id: 'morgan',
      name: 'Morgan',
      initials: 'M',
      availability: [
        { date: 'Mon', status: 'busy' },
        { date: 'Tue', status: 'partial' },
        { date: 'Wed', status: 'available' },
        { date: 'Thu', status: 'off' },
      ],
    },
    {
      id: 'sam',
      name: 'Sam',
      initials: 'S',
      availability: [
        { date: 'Mon', status: 'off' },
        { date: 'Tue', status: 'available' },
        { date: 'Wed', status: 'partial' },
        { date: 'Thu', status: 'busy' },
      ],
    },
    {
      id: 'riley',
      name: 'Riley',
      initials: 'R',
      availability: [
        { date: 'Mon', status: 'partial' },
        { date: 'Tue', status: 'busy' },
        { date: 'Wed', status: 'available' },
        { date: 'Thu', status: 'available' },
      ],
    },
  ],
  legend: () => [
    { id: 'available', label: 'Available' },
    { id: 'busy', label: 'Busy' },
    { id: 'partial', label: 'Partial' },
    { id: 'off', label: 'Off' },
  ],
  availableToday: 2,
  totalPeople: 4,
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>{{ title }}</CardTitle>
      <CardDescription>{{ description }}</CardDescription>
    </CardHeader>
    <CardContent>
      <div class="overflow-x-auto">
        <!-- Summary row -->
        <p class="mb-4 text-sm text-muted-foreground">
          {{ availableToday }} of {{ totalPeople }} available today
        </p>

        <!-- Matrix -->
        <table class="w-full border-collapse">
          <thead>
            <tr>
              <th class="px-2 py-1 text-left text-sm font-medium text-muted-foreground">
                Person
              </th>
              <th
                v-for="day in days"
                :key="day"
                class="px-2 py-1 text-center text-sm font-medium text-muted-foreground"
              >
                {{ day }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="person in people"
              :key="person.id"
              class="border-t border-border"
            >
              <td class="px-2 py-2">
                <div class="flex items-center gap-2">
                  <Avatar :name="person.name" class="size-7 text-xs" />
                  <span class="text-sm font-medium">{{ person.name }}</span>
                </div>
              </td>
              <td
                v-for="entry in person.availability"
                :key="entry.date"
                class="px-2 py-2 text-center"
              >
                <Badge :variant="entry.status">
                  {{ entry.status.charAt(0).toUpperCase() + entry.status.slice(1) }}
                </Badge>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Legend -->
        <div class="mt-4 flex flex-wrap items-center gap-3">
          <Badge
            v-for="item in legend"
            :key="item.id"
            :variant="item.id as StatusVariant"
          >
            {{ item.label }}
          </Badge>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
