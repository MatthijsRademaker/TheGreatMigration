<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <CardTitle>Rooms / Areas</CardTitle>
        <CardDescription>
          Organize and label rooms, floors, and zones for a clear move-day plan. Feature coming soon.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">
          Room management will let you define the floor plan, group areas by type, and attach packing checklists.
        </p>
      </CardContent>
    </Card>
  </section>
</template>
