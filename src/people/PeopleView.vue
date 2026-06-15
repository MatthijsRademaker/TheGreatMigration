<script setup lang="ts">
import { Badge } from '@/shared/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

const helpers = [
  { name: 'Alex', availability: 'Mon, Wed, Fri' },
  { name: 'Morgan', availability: 'Tue, Wed' },
  { name: 'Sam', availability: 'Wed, Fri' },
  { name: 'Riley', availability: 'Mon, Thu' },
]
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <CardTitle>People availability foundation</CardTitle>
        <CardDescription>People are simple free-form helpers with day-level availability and assignment visibility.</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
          <div v-for="helper in helpers" :key="helper.name" class="rounded-lg border bg-muted/40 p-4">
            <div class="flex items-center justify-between gap-3">
              <p class="font-medium">{{ helper.name }}</p>
              <Badge variant="secondary">Available</Badge>
            </div>
            <p class="mt-3 text-sm text-muted-foreground">{{ helper.availability }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
