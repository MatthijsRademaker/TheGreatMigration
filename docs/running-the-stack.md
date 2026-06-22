# Running the Stack

The Great Migration ships two Docker Compose entrypoints: a **development** stack with hot reload
and seeded demo data, and a **production** stack with optimized assets and a clean database.

Both define the same three services — `db` (Postgres 16), `backend` (Go), `frontend` (Vue) — on a
shared `app-network`, backed by the `pgdata` named volume.

## Development — `compose.yml`

```bash
docker compose up
```

- **Frontend:** Vite dev server with HMR on `http://localhost:5173`. The source is bind-mounted
  (`./frontend`), so edits reload live. Vite proxies `/api` to the backend.
- **Backend:** built from `backend/Dockerfile` (full Go toolchain + source), published on
  `http://localhost:8080`.
- **Database:** published on `localhost:5432`, runs migrations **and seeds demo data**
  (`DB_SEED: "true"`) — people, tasks, rooms, schedule cards.

Use this for day-to-day development.

## Production — `compose.prod.yml`

```bash
docker compose -f compose.prod.yml up
```

This is a **standalone** file — do not layer it over `compose.yml`. It brings up a
production-shaped stack:

- **Frontend:** `frontend/Dockerfile.prod` runs an optimized `npm run build` and serves the static
  `dist/` via **nginx** on `http://localhost:80`. nginx reverse-proxies `/api` to the backend and
  provides SPA history-mode fallback, so deep links (e.g. `/tools`) resolve to the app. No dev
  server, no bind mount, no HMR.
- **Backend:** `backend/Dockerfile.prod` builds a static `CGO_ENABLED=0` binary into a slim
  `alpine` runtime (no Go toolchain or source shipped). Runs migrations
  (`DB_AUTO_MIGRATE=true`) but **does not seed** demo data — a fresh `pgdata` volume yields a
  schema-only database.
- **Networking:** only the frontend publishes a host port (`80`). `backend:8080` and `db:5432` are
  reachable **only** inside the compose network.

Open `http://localhost/`.

### Notes

- **Fresh database:** to start prod from an empty DB (e.g. after dev seeding), reset the shared
  volume first: `docker compose -f compose.prod.yml down -v`.
- **Rebuild images** after code changes: `docker compose -f compose.prod.yml build`.
- **Out of scope** for this variant: TLS (assume a TLS-terminating ingress in front) and secrets
  management (the DB still uses the simple `app/app` credentials).
