#!/usr/bin/env bash
set -euo pipefail

: "${BACKEND_BASE_URL:?BACKEND_BASE_URL is required}"

curl -fsS "${BACKEND_BASE_URL}/healthz"
curl -fsS "${BACKEND_BASE_URL}/readyz"
