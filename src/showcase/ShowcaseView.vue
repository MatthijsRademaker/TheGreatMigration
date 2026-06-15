<script setup lang="ts">
import { Badge } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import { Separator } from '@/shared/ui/separator'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/shared/ui/sheet'
import { Skeleton } from '@/shared/ui/skeleton'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/shared/ui/tooltip'

// --- Catalog data structure (2.2) ---

type ButtonVariant = 'default' | 'outline' | 'secondary' | 'ghost' | 'destructive' | 'link'
type ButtonSize = 'default' | 'sm' | 'lg'

interface ButtonExample {
  variant: ButtonVariant
  size: ButtonSize
}

const buttonVariants: ButtonVariant[] = ['default', 'outline', 'secondary', 'ghost', 'destructive', 'link']
const buttonSizes: ButtonSize[] = ['default', 'sm', 'lg']
const buttonExamples: ButtonExample[] = buttonVariants.flatMap(variant =>
  buttonSizes.map(size => ({ variant, size }))
)

const badgeVariants: ButtonVariant[] = ['default', 'secondary', 'destructive', 'outline', 'ghost', 'link']

// --- Template helpers ---

function variantLabel(v: string): string {
  return v.charAt(0).toUpperCase() + v.slice(1)
}
</script>

<template>
  <section class="flex flex-1 flex-col gap-8 p-4 sm:p-6">

    <!-- 2.3 Actions: Button variant × size matrix -->
    <Card>
      <CardHeader>
        <CardTitle>Actions</CardTitle>
        <CardDescription>
          Buttons trigger actions. Each variant conveys a different semantic weight, and
          sizes adapt to layout density.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex flex-wrap items-start gap-3">
          <div
            v-for="ex in buttonExamples"
            :key="`${ex.variant}-${ex.size}`"
            class="flex flex-col items-center gap-1"
          >
            <Button :variant="ex.variant" :size="ex.size">
              {{ variantLabel(ex.variant) }}
            </Button>
            <span class="text-[0.625rem] text-muted-foreground">
              {{ variantLabel(ex.variant) }} &middot; {{ ex.size }}
            </span>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- 2.4 Inputs: Input -->
    <Card>
      <CardHeader>
        <CardTitle>Inputs</CardTitle>
        <CardDescription>
          A text-entry field for forms and search.
        </CardDescription>
      </CardHeader>
      <CardContent class="flex flex-col gap-3">
        <Input placeholder="Type here…" class="max-w-sm" />
        <p class="text-xs text-muted-foreground">
          Usage: binds with <code class="rounded bg-muted px-1 text-[0.625rem]">v-model</code> for form state.
        </p>
      </CardContent>
    </Card>

    <!-- 2.5 Feedback: Badge + Skeleton -->
    <Card>
      <CardHeader>
        <CardTitle>Feedback</CardTitle>
        <CardDescription>
          Badges label status or counts. Skeletons reserve space while content loads.
        </CardDescription>
      </CardHeader>
      <CardContent class="flex flex-col gap-6">
        <div class="flex flex-col gap-2">
          <p class="text-xs font-medium">Badge variants</p>
          <div class="flex flex-wrap items-center gap-2">
            <div
              v-for="v in badgeVariants"
              :key="v"
              class="flex flex-col items-center gap-1"
            >
              <Badge :variant="v">{{ variantLabel(v) }}</Badge>
              <span class="text-[0.625rem] text-muted-foreground">{{ variantLabel(v) }}</span>
            </div>
          </div>
        </div>

        <Separator />

        <div class="flex flex-col gap-2">
          <p class="text-xs font-medium">Skeleton placeholder</p>
          <Skeleton class="h-5 w-64 rounded" />
          <p class="text-xs text-muted-foreground">
            The skeleton pulses to indicate a loading placeholder. Use it to reserve
            layout space while async content resolves.
          </p>
        </div>
      </CardContent>
    </Card>

    <!-- 2.6 Layout: Card + Separator -->
    <Card>
      <CardHeader>
        <CardTitle>Layout</CardTitle>
        <CardDescription>
          Cards structure content into discrete sections. Separators divide content
          horizontally or vertically.
        </CardDescription>
      </CardHeader>
      <CardContent class="flex flex-col gap-6">
        <div class="flex flex-col gap-2">
          <p class="text-xs font-medium">Composed card</p>
          <Card class="max-w-md">
            <CardHeader>
              <CardTitle>Example card</CardTitle>
              <CardDescription>
                This card combines Card, CardHeader, CardTitle, CardDescription, and
                CardContent.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p class="text-sm text-muted-foreground">
                Card subcomponents make it easy to compose consistent layouts.
              </p>
            </CardContent>
          </Card>
        </div>

        <Separator />

        <div class="flex flex-col gap-2">
          <p class="text-xs font-medium">Separator orientations</p>
          <div class="flex flex-col gap-3">
            <div class="flex items-center gap-2">
              <span class="text-[0.625rem] text-muted-foreground w-24">Horizontal</span>
              <Separator class="flex-1" orientation="horizontal" />
            </div>
            <div class="flex items-center gap-2 h-16">
              <span class="text-[0.625rem] text-muted-foreground w-24">Vertical</span>
              <Separator orientation="vertical" />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- 2.7 Overlays: Sheet + Tooltip -->
    <Card>
      <CardHeader>
        <CardTitle>Overlays</CardTitle>
        <CardDescription>
          Sheets slide in from the edge for focused tasks. Tooltips show contextual
          help on hover.
        </CardDescription>
      </CardHeader>
      <CardContent class="flex flex-col gap-6">
        <div class="flex flex-col gap-2">
          <p class="text-xs font-medium">Sheet demo</p>
          <Sheet>
            <SheetTrigger as-child>
              <Button variant="outline" class="w-fit">Open sheet</Button>
            </SheetTrigger>
            <SheetContent>
              <SheetHeader>
                <SheetTitle>Example sheet</SheetTitle>
                <SheetDescription>
                  This sheet demonstrates the Sheet, SheetTrigger, SheetContent,
                  SheetHeader, SheetTitle, and SheetDescription components.
                </SheetDescription>
              </SheetHeader>
            </SheetContent>
          </Sheet>
          <p class="text-xs text-muted-foreground">
            Click the trigger button to open the sheet panel.
          </p>
        </div>

        <Separator />

        <div class="flex flex-col gap-2">
          <p class="text-xs font-medium">Tooltip demo</p>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="outline" class="w-fit">Hover me</Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>This is a tooltip with contextual help.</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <p class="text-xs text-muted-foreground">
            Hover the button to see the tooltip. Touch devices may require a long-press.
          </p>
        </div>
      </CardContent>
    </Card>

    <!-- 2.8 Sidebar family prose entry -->
    <Card>
      <CardHeader>
        <CardTitle>Sidebar</CardTitle>
        <CardDescription>
          Layout-coupled component family — not rendered standalone.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <p class="text-sm text-muted-foreground leading-relaxed">
          The Sidebar component family (<code class="rounded bg-muted px-1 text-[0.625rem]">Sidebar</code>,
          <code class="rounded bg-muted px-1 text-[0.625rem]">SidebarProvider</code>,
          <code class="rounded bg-muted px-1 text-[0.625rem]">SidebarMenu</code>, etc.) requires a
          <code class="rounded bg-muted px-1 text-[0.625rem]">SidebarProvider</code> context and is
          designed for layout-level composition. It cannot be rendered standalone inside a
          catalog card. The canonical usage example is the live
          <code class="rounded bg-muted px-1 text-[0.625rem]">AppSidebar</code> component
          visible in the sidebar of this application — it uses collapsible
          <code class="rounded bg-muted px-1 text-[0.625rem]">"icon"</code> mode, inset variant,
          grouped navigation menus, and a footer, all backed by
          <code class="rounded bg-muted px-1 text-[0.625rem]">SidebarProvider</code> in
          <code class="rounded bg-muted px-1 text-[0.625rem]">AppShell</code>.
        </p>
      </CardContent>
    </Card>

  </section>
</template>
