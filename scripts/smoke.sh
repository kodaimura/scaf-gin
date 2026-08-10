#!/bin/sh
set -eu

base_url="${1:-http://localhost:8000}"

curl -fsS "${base_url}/health" >/dev/null
printf 'smoke ok: %s/health\n' "${base_url}"
