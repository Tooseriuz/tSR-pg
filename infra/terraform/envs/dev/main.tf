terraform {
  required_version = ">= 1.6.0"
}

module "cloud_run_api" {
  source       = "../../modules/cloud_run_api"
  service_name = "tsr-pg-api-dev"
}

module "artifact_registry" {
  source        = "../../modules/artifact_registry"
  repository_id = "tsr-pg-dev"
}

module "github_oidc" {
  source  = "../../modules/github_oidc"
  pool_id = "github-actions-dev"
}

module "secret_manager" {
  source     = "../../modules/secret_manager"
  secret_ids = ["api-config-dev"]
}
