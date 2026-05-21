terraform {
  required_version = ">= 1.6.0"
}

module "cloud_run_api" {
  source       = "../../modules/cloud_run_api"
  service_name = "tsr-pg-api-prod"
}

module "artifact_registry" {
  source        = "../../modules/artifact_registry"
  repository_id = "tsr-pg-prod"
}

module "github_oidc" {
  source  = "../../modules/github_oidc"
  pool_id = "github-actions-prod"
}

module "secret_manager" {
  source     = "../../modules/secret_manager"
  secret_ids = ["api-config-prod"]
}
