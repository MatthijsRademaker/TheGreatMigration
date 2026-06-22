# The Great Migration

A full-stack app for planning a household move — people, tasks, rooms, tools, and a daily schedule.
Vue 3 frontend, Go backend, Postgres, orchestrated with Docker Compose.

## Quick start

**Development** (Vite HMR + seeded demo data):

```bash
docker compose up
```

→ app on `http://localhost:5173`, backend on `http://localhost:8080`.

**Production** (optimized static build, no demo data, single ingress):

```bash
docker compose -f compose.prod.yml up
```

→ app on `http://localhost/` (nginx serves the built frontend and proxies `/api` to the backend;
only port `80` is exposed).

See [`docs/running-the-stack.md`](docs/running-the-stack.md) for the full breakdown of both stacks.

## Layout

- `frontend/` — Vue 3 + Vite SPA. `Dockerfile` (dev), `Dockerfile.prod` + `nginx.conf` (prod).
- `backend/` — Go API with embedded migrations. `Dockerfile` (dev), `Dockerfile.prod` (prod).
- `compose.yml` — development stack. `compose.prod.yml` — production stack.
- `openspec/` — spec-driven change proposals and the live capability specs.
- `docs/` — design system and operational docs.
