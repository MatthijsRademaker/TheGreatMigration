## Why

The only way to run the stack today is the dev `compose.yml`: the frontend is a Vite **dev
server** (HMR, source bind-mounted, unminified) and the backend always seeds demo data. There is
no production-shaped way to stand the app up — optimized static frontend assets, no fictional demo
data, images that don't ship the full toolchain and source.

This change adds a standalone `compose.prod.yml` entrypoint plus production image builds, so
`docker compose -f compose.prod.yml up` brings up a clean, optimized stack.

## Depends on

`extract-seed-data-flag` — that change makes seed data opt-in via `DB_SEED` (default `false`).
This change relies on simply **omitting** `DB_SEED` in prod to get a schema-only database. Without
it, the production stack would still seed demo data. Land that change first.

## What Changes

- **Add `compose.prod.yml`** as a standalone (not override-layered) entrypoint orchestrating the
  same three services in production shape:
  - `db` — unchanged image, but **no host port** published (internal only).
  - `backend` — built from a new production Dockerfile, `DB_AUTO_MIGRATE=true`, `DB_SEED` **unset**
    (schema-only), **no host port** published (reached only via the frontend proxy on the compose
    network).
  - `frontend` — built from a new production Dockerfile that compiles an optimized `vite build`
    and serves the static `dist/` via nginx, published on host port `80`. **No bind mounts, no
    `node_modules` volume, no HMR.**
- **Add `frontend/Dockerfile.prod`** — multi-stage: `node` build stage (`npm run build`) →
  `nginx:alpine` runtime serving `dist/`. nginx config provides SPA history-mode fallback
  (`try_files … /index.html`) and reverse-proxies `/api` to `backend:8080`, so the relative-URL
  API client keeps working same-origin with no rebuild-time base URL.
- **Add `frontend/nginx.conf`** — the static-serve + `/api` proxy + SPA fallback configuration.
- **Add `backend/Dockerfile.prod`** — multi-stage: `golang` builder producing the static
  `CGO_ENABLED=0` binary → minimal `alpine` runtime (keeps `wget` for the existing healthcheck),
  instead of shipping the Go toolchain and source.

The dev `compose.yml`, `frontend/Dockerfile`, and `backend/Dockerfile` are **left untouched**.

## Capabilities

### New Capabilities

- `production-compose`: A standalone production compose entrypoint and the production frontend and
  backend image builds — optimized static frontend served by nginx with an `/api` reverse proxy
  and SPA fallback, a slim backend runtime image, a schema-only database (no demo seed), and
  internal-only `db`/`backend` networking with only the frontend exposed.

## Impact

- **New files**: `compose.prod.yml`, `frontend/Dockerfile.prod`, `frontend/nginx.conf`,
  `backend/Dockerfile.prod`.
- **No changes** to application code, the dev compose, the dev Dockerfiles, the API client
  (relative `/api` URLs already work same-origin), or the router (history-mode fallback handled in
  nginx).
- **Networking**: in prod only the frontend publishes a host port (`80`); `backend:8080` and
  `db:5432` are reachable only inside the compose network.
- **Data**: the `pgdata` named volume is reused; a fresh prod volume yields a schema-only database.
