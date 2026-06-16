<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
</script>

<template>
  <section class="flex flex-1 flex-col gap-6 p-4 sm:p-6">
    <Card>
      <CardHeader>
        <CardTitle>Settings</CardTitle>
        <CardDescription>
          Configure your move preferences, notification defaults, and account details. Feature coming soon.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground">
          Settings will include planning window dates, team visibility controls, notification preferences, and account management.
        </p>
      </CardContent>
    </Card>
  </section>
</template>
