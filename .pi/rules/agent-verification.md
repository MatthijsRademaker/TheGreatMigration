# Agent Verification Workflow

You are running inside a Docker container without Go, Node.js, Python, or other SDKs installed.
All verification (build, test, lint) must run through the Docker-backed scripts in `scripts/`.

For container and runtime execution details, see `runtime-environment.md`.

## Recommended workflow

After making code changes, run verification through Docker:

```bash
scripts/precommit-run
```

This runs all pre-commit hooks (`go vet`, `go build`, `go test`, dashboard checks, proto generation)
inside Docker containers. It is the authoritative verification path used by CI and worker agents.

## Individual scripts

To run specific checks in isolation, use the individual scripts instead of the full pre-commit run:

```bash
scripts/check-go          # go vet
scripts/test-go           # go test -v ./...
scripts/build-go          # go build -v ./...
scripts/check-dashboard   # dashboard type-check
scripts/build-dashboard   # dashboard production build
scripts/generate-proto    # buf generate
```

## Do not use

- `go build`, `go test`, `go vet` directly from host (no Go SDK)
- `make lint`, `make test`, `make build` from `src/` (host SDK convenience, not available to agents)
- Any Node.js or Python tooling directly from host
