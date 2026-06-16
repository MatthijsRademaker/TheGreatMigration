<script setup lang="ts">
import { Badge } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import { Avatar } from '@/shared/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/ui/popover'
import type { AvailabilityStatus, CellChangePayload, PeopleAvailabilityProps } from './types'

withDefaults(defineProps<PeopleAvailabilityProps>(), {
  title: 'People availability',
  description: 'Track who is available and where each person can help.',
  days: () => ['Mon', 'Tue', 'Wed', 'Thu'],
  people: () => [
    {
      id: 'alex',
      name: 'Alex',
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
  // Keep in sync with demo people data: counts people with status 'available' on the first day column.
  availableToday: 1,
  totalPeople: 4,
  editable: false,
})

const emit = defineEmits<{
  'update-cell': [payload: CellChangePayload]
  'delete-person': [personId: string]
}>()

const statusOptions: AvailabilityStatus[] = ['available', 'busy', 'partial', 'off']
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
          <caption class="sr-only">
            People availability matrix showing each person's status across the planning days
          </caption>
          <thead>
            <tr>
              <th scope="col" class="px-2 py-1 text-left text-sm font-medium text-muted-foreground">
                Person
              </th>
              <th
                v-for="day in days"
                :key="day"
                scope="col"
                class="px-2 py-1 text-center text-sm font-medium text-muted-foreground"
              >
                {{ day }}
              </th>
              <th
                v-if="editable"
                scope="col"
                class="px-2 py-1 text-center text-sm font-medium text-muted-foreground"
              >
                Actions
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
                v-for="(entry, dayIdx) in person.availability"
                :key="`${person.id}-${entry.date}`"
                class="px-2 py-2 text-center"
              >
                <!-- Editable mode: clickable badge wrapped in Popover -->
                <template v-if="editable">
                  <Popover>
                    <PopoverTrigger as-child>
                      <button
                        type="button"
                        class="cursor-pointer rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1"
                      >
                        <Badge :variant="entry.status" class="pointer-events-none">
                          {{ entry.status.charAt(0).toUpperCase() + entry.status.slice(1) }}
                        </Badge>
                      </button>
                    </PopoverTrigger>
                    <PopoverContent class="w-48 p-2">
                      <div class="flex flex-col gap-1">
                        <p class="mb-1 text-xs font-medium text-muted-foreground">
                          Set status
                        </p>
                        <button
                          v-for="s in statusOptions"
                          :key="s"
                          type="button"
                          class="flex w-full cursor-pointer items-center gap-2 rounded-md px-2 py-1 text-sm hover:bg-muted"
                          @click="emit('update-cell', { personId: person.id, dayIndex: dayIdx, status: s })"
                        >
                          <Badge :variant="s" class="pointer-events-none">
                            {{ s.charAt(0).toUpperCase() + s.slice(1) }}
                          </Badge>
                        </button>
                        <hr class="my-1 border-border">
                        <button
                          type="button"
                          class="flex w-full cursor-pointer items-center gap-2 rounded-md px-2 py-1 text-sm text-muted-foreground hover:bg-muted"
                          @click="emit('update-cell', { personId: person.id, dayIndex: dayIdx, status: null })"
                        >
                          Clear (reset to off)
                        </button>
                      </div>
                    </PopoverContent>
                  </Popover>
                </template>
                <!-- Read-only mode: plain badge -->
                <template v-else>
                  <Badge :variant="entry.status">
                    {{ entry.status.charAt(0).toUpperCase() + entry.status.slice(1) }}
                  </Badge>
                </template>
              </td>
              <td v-if="editable" class="px-2 py-2 text-center">
                <Button
                  variant="destructive"
                  size="sm"
                  @click="emit('delete-person', person.id)"
                >Delete</Button>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Legend -->
        <div class="mt-4 flex flex-wrap items-center gap-3">
          <Badge
            v-for="item in legend"
            :key="item.id"
            :variant="item.id"
          >
            {{ item.label }}
          </Badge>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
