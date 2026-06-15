## Implementation Tasks

### 1. Add greeting banner to HomeView.vue

- [ ] Insert a greeting `<div>` as the first child of the `<section>` element in `src/home/HomeView.vue`, above the summary cards grid
- [ ] Use Tailwind CSS v4 utility classes: `bg-muted/40 rounded-lg border p-4 border-l-4 border-primary text-sm`
- [ ] Set the text content to: `👋 Hello world v2 — swarm pipeline test`
- [ ] Run `scripts/precommit-run` to verify pre-commit hooks pass
- [ ] Run `npm run build` to verify the app builds successfully (vue-tsc -b && vite build)