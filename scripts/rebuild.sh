#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "🧹 Tearing down containers and volumes..."
docker compose down -v

echo "🏗️  Rebuilding images from scratch..."
docker compose build --no-cache

echo "🚀 Starting services..."
docker compose up -d

echo "✅ Done! Services are starting up."
