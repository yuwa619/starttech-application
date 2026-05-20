#!/usr/bin/env bash
set -euo pipefail

: "${AWS_REGION:?AWS_REGION is required}"
: "${ECR_REPOSITORY:?ECR_REPOSITORY is required}"
: "${BACKEND_ASG_NAME:?BACKEND_ASG_NAME is required}"
: "${ROLLBACK_IMAGE_TAG:?ROLLBACK_IMAGE_TAG is required}"

aws ecr get-login-password --region "${AWS_REGION}" \
  | docker login --username AWS --password-stdin "${ECR_REPOSITORY%/*}"

docker pull "${ECR_REPOSITORY}:${ROLLBACK_IMAGE_TAG}"
docker tag "${ECR_REPOSITORY}:${ROLLBACK_IMAGE_TAG}" "${ECR_REPOSITORY}:latest"
docker push "${ECR_REPOSITORY}:latest"

aws autoscaling start-instance-refresh \
  --auto-scaling-group-name "${BACKEND_ASG_NAME}" \
  --strategy Rolling \
  --preferences MinHealthyPercentage=50,InstanceWarmup=120
