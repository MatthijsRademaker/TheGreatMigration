## 1. Frontend production image

- [x] 1.1 Add `frontend/Dockerfile.prod`: multi-stage `node:24-alpine` build stage running the
      existing install step + `npm run build`, then an `nginx:alpine` runtime copying `dist/` into
      `/usr/share/nginx/html`.
- [x] 1.2 Add `frontend/nginx.conf`: `listen 80`, static root, `location /api/ { proxy_pass
      http://backend:8080; }`, and `location / { try_files $uri $uri/ /index.html; }` for
      history-mode SPA fallback. Copy it to `/etc/nginx/conf.d/default.conf` in the Dockerfile.
- [x] 1.3 Confirm the production bundle builds: run `npm run build` (`vue-tsc -b && vite build`)
      in the pinned `node:24-alpine` image, or `docker compose -f compose.prod.yml build frontend`.
      Note: no existing script runs the prod bundle — precommit only runs `npm run check`
      (`vue-tsc --noEmit`), so this path is otherwise unverified.

## 2. Backend production image

- [x] 2.1 Add `backend/Dockerfile.prod`: `golang:1.26.2-alpine` builder producing the
      `CGO_ENABLED=0` static binary, then an `alpine:latest` runtime copying only the binary and
      exposing `8080`.
- [x] 2.2 Confirm the backend builds via `scripts/build-go`.

## 3. Production compose file

- [x] 3.1 Add `compose.prod.yml` (standalone) with `db`, `backend`, `frontend` on a shared network
      and the `pgdata` named volume.
- [x] 3.2 `db`: same `postgres:16-alpine` config and healthcheck as dev, but **no** `ports:` entry.
- [x] 3.3 `backend`: `build` from `backend/Dockerfile.prod`, `DATABASE_URL` as in dev,
      `DB_AUTO_MIGRATE: "true"`, **no `DB_SEED`** (schema-only), `depends_on: db service_healthy`,
      retain the `wget /api/hello` healthcheck, **no** host `ports:`.
- [x] 3.4 `frontend`: `build` from `frontend/Dockerfile.prod`, `ports: ["80:80"]`,
      `depends_on: backend service_healthy`, **no** bind mounts and **no** `node_modules` volume.

## 4. Verification

- [x] 4.1 `docker compose -f compose.prod.yml up` brings all three services to healthy.
- [x] 4.2 The app loads on `http://localhost/`, a deep link (e.g. `/tools`) served directly
      returns the SPA (not a 404), and API calls to `/api/...` succeed through the nginx proxy.
- [x] 4.3 With a fresh `pgdata` volume, the domain tables are **empty** (no demo people/tasks),
      confirming `DB_SEED` is off in prod.
- [x] 4.4 `backend:8080` and `db:5432` are **not** reachable from the host; only `:80` is.
- [x] 4.5 Dev `docker compose up` (the unchanged `compose.yml`) still works with HMR and seeded data.
