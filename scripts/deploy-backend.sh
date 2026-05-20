#!/usr/bin/env bash
set -euo pipefail

: "${AWS_REGION:?AWS_REGION is required}"
: "${ECR_REPOSITORY:?ECR_REPOSITORY is required}"
: "${BACKEND_ASG_NAME:?BACKEND_ASG_NAME is required}"

IMAGE_TAG="${IMAGE_TAG:-$(git rev-parse --short HEAD)}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/../backend" && pwd)"

(cd "${BACKEND_DIR}" && go test ./...)

aws ecr get-login-password --region "${AWS_REGION}" \
  | docker login --username AWS --password-stdin "${ECR_REPOSITORY%/*}"

docker build -t "${ECR_REPOSITORY}:${IMAGE_TAG}" -t "${ECR_REPOSITORY}:latest" "${BACKEND_DIR}"
docker push "${ECR_REPOSITORY}:${IMAGE_TAG}"
docker push "${ECR_REPOSITORY}:latest"

aws autoscaling start-instance-refresh \
  --auto-scaling-group-name "${BACKEND_ASG_NAME}" \
  --strategy Rolling \
  --preferences MinHealthyPercentage=50,InstanceWarmup=120
