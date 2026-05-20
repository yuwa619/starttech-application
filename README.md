# StartTech Application

Full-stack StartTech application used by the CI/CD assessment.

## Structure

- `frontend/`: React application built with Vite and deployed to S3/CloudFront
- `backend/`: Go HTTP API packaged as a Docker image and deployed to EC2
- `.github/workflows/`: frontend and backend CI/CD pipelines
- `scripts/`: local deployment, health check, and rollback scripts

## Frontend

```bash
cd frontend
npm ci
npm test
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

Production deployment is handled by `.github/workflows/frontend-ci-cd.yml`.

## Backend

```bash
cd backend
go test ./...
go run ./cmd/server
```

Environment variables:

- `APP_ENV`
- `PORT`
- `MONGODB_URI`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `REDIS_TLS`
- `SERVICE_NAME`

## GitHub Actions Secrets

Set these secrets in the application repository:

- `AWS_ROLE_TO_ASSUME`
- `AWS_REGION`
- `AWS_S3_BUCKET`
- `CLOUDFRONT_DISTRIBUTION_ID`
- `ECR_REPOSITORY`
- `BACKEND_ASG_NAME`
- `BACKEND_BASE_URL`
- `BACKEND_LOG_GROUP`
- `MONGODB_URI`
- `DEPLOYMENT_WEBHOOK_URL` optional

## MongoDB Atlas

Create a MongoDB Atlas cluster, allow the backend egress IPs or VPC peering path, create an application user, and store the connection string as `MONGODB_URI` in GitHub secrets or the EC2 runtime environment.

The backend deployment workflow writes `MONGODB_URI` into AWS Systems Manager Parameter Store at `/starttech/prod/backend/mongodb_uri`. EC2 instances read that secure parameter during launch.

## Pipelines

The frontend pipeline installs dependencies, runs tests, performs `npm audit`, builds the bundle, syncs to S3, invalidates CloudFront, and sends an optional deployment notification.

The backend pipeline runs Go tests, `go vet`, `govulncheck`, builds a Docker image, scans it with Trivy, pushes it to ECR, starts an Auto Scaling Group rolling refresh, runs smoke tests, and records a CloudWatch deployment marker.
