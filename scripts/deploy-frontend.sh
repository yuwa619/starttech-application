#!/usr/bin/env bash
set -euo pipefail

: "${AWS_S3_BUCKET:?AWS_S3_BUCKET is required}"
: "${CLOUDFRONT_DISTRIBUTION_ID:?CLOUDFRONT_DISTRIBUTION_ID is required}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$(cd "${SCRIPT_DIR}/../frontend" && pwd)"

npm --prefix "${FRONTEND_DIR}" ci
npm --prefix "${FRONTEND_DIR}" test
npm --prefix "${FRONTEND_DIR}" run build

aws s3 sync "${FRONTEND_DIR}/dist" "s3://${AWS_S3_BUCKET}" \
  --delete \
  --cache-control "public,max-age=31536000,immutable" \
  --exclude "index.html"

aws s3 cp "${FRONTEND_DIR}/dist/index.html" "s3://${AWS_S3_BUCKET}/index.html" \
  --cache-control "no-cache,no-store,must-revalidate"

aws cloudfront create-invalidation \
  --distribution-id "${CLOUDFRONT_DISTRIBUTION_ID}" \
  --paths "/*"
