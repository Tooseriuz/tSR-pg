# Terraform

This directory manages:

- Cloud Run API service: `tooseriuzdotcom-api`
- Artifact Registry Docker repository: `tooseriuzdotcom-ar`
- GCS bucket for API file storage
- GitHub Actions OIDC for deploys from `Tooseriuz/tSR-pg`

Secret Manager is intentionally not enabled yet.

## Apply prod

```sh
cd infra/terraform/envs/prod
terraform init
terraform apply -var="project_id=<gcp-project-id>"
```

The GitHub Actions workflow expects these repository variables:

- `GCP_PROJECT_ID`
- `GCP_PROJECT_NUMBER`
