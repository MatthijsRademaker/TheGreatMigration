## Why

This is a pipeline validation task — the first end-to-end test of the swarm development pipeline (OpenSpec → Refinement Room → task execution → implementation) inside this repository. Success is defined by the pipeline producing a correct, mergeable implementation, not by the feature itself.

## What Changes

Add a visible greeting banner at the top of the home page in `src/home/HomeView.vue` to prove the pipeline can generate and execute a minimal implementation task. The banner is a plain `<div>` element rendered above the existing summary cards grid, containing the text "👋 Hello world — swarm pipeline test" and styled with Tailwind CSS v4 utility classes.

**Files changed:**
- `src/home/HomeView.vue` — ~8–12 lines added to the `<template>` block (no script or style changes)

## Impact

- **Zero functional impact** — the banner is purely presentational and adds no data fetching, routing, state management, or user interaction.
- **No new dependencies** — uses only the existing Vue 3 + Tailwind v4 + shadcn-vue stack already in the project.
- **No architecture changes** — the app shell, router, and all existing components are untouched.
- **Verification:** the change passes `scripts/precommit-run` and `npm run build` (vue-tsc -b && vite build).