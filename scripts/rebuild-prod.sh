#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

RUN_SEED=false

usage() {
	cat <<'EOF'
Usage: scripts/rebuild-prod.sh [--seed]

Options:
  --seed   Run prod seed job after services start
  --help   Show this help
EOF
}

for arg in "$@"; do
	case "$arg" in
	--seed)
		RUN_SEED=true
		;;
	--help)
		usage
		exit 0
		;;
	*)
		echo "Unknown argument: $arg" >&2
		usage >&2
		exit 1
		;;
	esac
done

echo "🧹 Tearing down containers (keeping volumes)..."
docker compose -f compose.prod.yml down

echo "🏗️  Rebuilding images from scratch..."
docker compose -f compose.prod.yml build --no-cache

echo "🚀 Starting services (volumes preserved)..."
docker compose -f compose.prod.yml up -d

if [[ "$RUN_SEED" == "true" ]]; then
	echo "🌱 Running seed job..."
	docker compose -f compose.prod.yml --profile ops run --rm seed
	echo "✅ Done! Prod services rebuilt. Volumes preserved. Seed ran."
else
	echo "✅ Done! Prod services rebuilt. Volumes preserved. Seed skipped."
fi
