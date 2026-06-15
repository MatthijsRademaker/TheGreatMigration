## Implementation Tasks

### 1. Add greeting banner to HomeView.vue

- [x] Insert a greeting `<div>` as the first child of the `<section>` element in `src/home/HomeView.vue`, above the summary cards grid
- [x] Use Tailwind CSS v4 utility classes: `bg-muted/40 rounded-lg border p-4 border-l-4 border-primary text-sm`
- [x] Set the text content to: `👋 Hello world — swarm pipeline test`
- [x] Run `scripts/precommit-run` to verify pre-commit hooks pass
- [x] Run `npm run build` to verify the app builds successfully (vue-tsc -b && vite build)
