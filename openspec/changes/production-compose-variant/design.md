# Design — Production compose variant

## Context

Dev `compose.yml`: frontend = Vite dev server on `:5173` (HMR, `./frontend` bind-mounted,
`/app/node_modules` anonymous volume), proxying `/api` to the backend via Vite's `server.proxy`.
Backend = `golang:1.26.2-alpine` image carrying the toolchain + source, always seeding. The API
client calls **relative `/api/...`** (no base URL configured anywhere) and the router uses
`createWebHistory()`.

Production needs: optimized static assets, no demo data, slim images, minimal host exposure.

## Key facts that drive the design

1. **Relative API URLs + same-origin serving.** Because the client uses relative `/api`, the
   simplest correct production topology is to serve the static assets and proxy `/api` from the
   **same origin**. nginx serving `dist/` and `proxy_pass`-ing `/api` to `backend:8080` satisfies
   this with **no build-time `VITE_API_BASE_URL`** and no client changes.
2. **History-mode routing.** `createWebHistory()` means deep links (`/tools`, `/people`) must fall
   back to `index.html`; nginx `try_files $uri $uri/ /index.html` handles it.
3. **Seed is already opt-in** (after `extract-seed-data-flag`): omit `DB_SEED` → schema-only DB.

## Decision 1 — standalone `compose.prod.yml`, not override layering

Two ways to express a prod variant:

- **Override layering** (`-f compose.yml -f compose.prod.yml`): prod states only diffs.
- **Standalone** (`-f compose.prod.yml`): a complete file.

Chosen: **standalone**. Reasons:
- Matches the requested "`compose.prod.yml` as entrypoint" and a one-flag invocation.
- Compose **cannot remove** an inherited `volumes:` entry, so an override could not cleanly drop
  the dev frontend's bind mount + `node_modules` volume — it would fight the base file.
- The prod frontend uses a different image, command, port, and has no mounts; it shares almost
  nothing with the dev service. Layering would be mostly overrides anyway.

Cost: the `db` service definition is duplicated across the two files. Accepted — it is small and
stable, and duplication is cheaper than the wrong abstraction here.

## Decision 2 — frontend prod image (multi-stage → nginx)

```
# frontend/Dockerfile.prod
FROM node:24-alpine AS build
WORKDIR /app
COPY package.json skills-lock.json* ./
RUN npm install --no-audit --no-fund --package-lock=false --legacy-peer-deps
COPY . .
RUN npm run build                      # vue-tsc -b && vite build → /app/dist

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

```
# frontend/nginx.conf (sketch)
server {
  listen 80;
  root /usr/share/nginx/html;
  location /api/ { proxy_pass http://backend:8080; proxy_set_header Host $host; }
  location /     { try_files $uri $uri/ /index.html; }   # SPA history fallback
}
```

nginx replaces Vite's dev proxy. A separate `Dockerfile.prod` (rather than a `target` in the dev
Dockerfile) keeps the dev image untouched — surgical, and dev never runs `npm run build`.

## Decision 3 — backend prod image (multi-stage → slim runtime)

The dev `backend/Dockerfile` ships `golang:1.26.2-alpine` with the full toolchain and source. For
prod, build the static binary in a builder stage and copy it into a minimal `alpine` runtime:

```
# backend/Dockerfile.prod
FROM golang:1.26.2-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o backend .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/backend .
EXPOSE 8080
CMD ["./backend"]
```

`alpine` (not `scratch`) is chosen so the existing compose **healthcheck** (`wget … /api/hello`)
keeps working without bundling extra tooling. Migrations are embedded in the binary via
`//go:embed`, so nothing else needs copying.

## Decision 4 — networking / exposure

Only the `frontend` publishes a host port (`80`). `backend` and `db` are **not** published — the
frontend nginx reaches `backend:8080` over the compose network, and the backend reaches `db:5432`.
This is the production posture: one ingress, internals private. `depends_on` healthchecks are
retained so startup ordering (db → backend → frontend) is unchanged.

## Out of scope (explicitly)

- TLS / certificates — assume a TLS-terminating ingress in front, or a follow-up change.
- Secrets management — `db` keeps the simple `app/app` credentials for now; productionizing
  credentials is a separate concern.
- Backend horizontal scaling, logging/observability — not part of "a production variant" as asked.
