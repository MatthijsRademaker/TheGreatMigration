## 1. Remove greeting banner from HomeView.vue

- [ ] 1.1 Delete the greeting banner `<div>` block (4 lines starting with `<div class="bg-muted/40...">` and ending with `</div>`) from the `<template>` section in `src/home/HomeView.vue`
- [ ] 1.2 Verify the removal leaves the `<section>` element with no orphaned whitespace or empty lines above the summary cards grid

## 2. Remove active greeting-banner spec

- [ ] 2.1 Delete the `openspec/specs/greeting-banner/` directory and all its contents

## 3. Verify project integrity

- [ ] 3.1 Run `scripts/precommit-run` to confirm pre-commit hooks pass
- [ ] 3.2 Run `npm run build` (vue-tsc -b && vite build) to confirm the app builds without errors
