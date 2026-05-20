# Application Runbook

## Deploy Frontend

1. Merge a frontend change to `main`.
2. Confirm the frontend workflow passes tests and audit.
3. Confirm the workflow syncs files to S3 and invalidates CloudFront.
4. Open the CloudFront URL from the infrastructure output.

## Deploy Backend

1. Merge a backend change to `main`.
2. Confirm Go tests, `go vet`, and vulnerability checks pass.
3. Confirm the Docker image is pushed to ECR with the commit SHA and `latest`.
4. Confirm the Auto Scaling Group instance refresh completes.
5. Run `scripts/health-check.sh`.

## Roll Back Backend

```bash
AWS_REGION=eu-west-2 \
ECR_REPOSITORY=<repo-url> \
BACKEND_ASG_NAME=<asg-name> \
ROLLBACK_IMAGE_TAG=<previous-sha> \
./scripts/rollback.sh
```

## Troubleshooting

- `frontend shows API offline`: confirm `VITE_API_BASE_URL` points to the ALB URL and the backend `/healthz` endpoint responds.
- `backend readiness returns 503`: check MongoDB Atlas network access, Redis endpoint, security groups, and environment variables.
- `ASG refresh stalls`: check ALB target group health and EC2 system logs.
- `Docker image pull fails`: confirm the EC2 role has ECR read permissions and the image exists in ECR.
