# Why

This is a repeat pipeline validation exercise ("hello world v2") for the swarm development pipeline inside this repository. The original hello world task (200FB581) added a greeting banner to HomeView.vue to prove end-to-end pipeline correctness — it was archived after completion, then cleaned up via a removal task (50d758b7) that deleted the banner and its spec. The project is now clean, and this v2 test re-validates that the OpenSpec → Refinement Room → task execution → implementation pipeline works correctly after pipeline changes/improvements. Success is defined by pipeline correctness and mergeable output, not by banner utility.

# What Changes

Add a visible greeting banner at the top of the home page in `src/home/HomeView.vue` to re-validate the swarm pipeline. The banner is a plain `<div>` element rendered as the first child of the `<section>` element, above the existing summary cards grid, containing the text "👋 Hello world v2 — swarm pipeline test" and styled with Tailwind CSS v4 utility classes.

**Files changed:**
- `src/home/HomeView.vue` — ~8–12 lines added to the `<template>` block (no script or style changes)
- `openspec/changes/task-5b60ec70-b3e1-4995-9aef-57257cf3e226/specs/greeting-banner/spec.md` — new spec matching the spec-driven contract
- `openspec/changes/task-5b60ec70-b3e1-4995-9aef-57257cf3e226/design.md` — new design document following v1 pattern

# Impact

- **Zero functional impact** — the banner is purely presentational and adds no data fetching, routing, state management, or user interaction.
- **No new dependencies** — uses only the existing Vue 3 + Tailwind v4 + shadcn-vue stack already in the project.
- **No architecture changes** — the app shell, router, and all existing components are untouched.
- **v2 traceability** — banner text includes "v2" to clearly distinguish this iteration from the archived v1 banner.
- **Verification:** the change passes `scripts/precommit-run` and `npm run build` (vue-tsc -b && vite build).
