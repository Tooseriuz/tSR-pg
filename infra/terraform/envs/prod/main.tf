terraform {
  required_version = ">= 1.6.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}

variable "project_id" {
  type        = string
  description = "Google Cloud project id."
}

variable "project_number" {
  type        = string
  description = "Google Cloud project number."
}

variable "region" {
  type        = string
  description = "Google Cloud region."
  default     = "asia-southeast1"
}

variable "github_repository" {
  type        = string
  description = "GitHub repository in owner/name format."
  default     = "Tooseriuz/tSR-pg"
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "required" {
  for_each = toset([
    "artifactregistry.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "run.googleapis.com",
    "sts.googleapis.com",
  ])

  project = var.project_id
  service = each.key
}

module "cloud_run_api" {
  source       = "../../modules/cloud_run_api"
  service_name = "tooseriuzdotcom-api"
  location     = var.region

  depends_on = [google_project_service.required]
}

module "artifact_registry" {
  source        = "../../modules/artifact_registry"
  repository_id = "tooseriuzdotcom-ar"
  location      = var.region

  depends_on = [google_project_service.required]
}

module "github_oidc" {
  source                          = "../../modules/github_oidc"
  pool_id                         = "github-actions"
  provider_id                     = "github"
  project_id                      = var.project_id
  project_number                  = var.project_number
  github_repository               = var.github_repository
  service_account_id              = "github-actions-deployer"
  cloud_run_service_account_email = module.cloud_run_api.service_account_email

  depends_on = [google_project_service.required]
}

output "artifact_registry_repository_id" {
  value = module.artifact_registry.repository_id
}

output "cloud_run_service_name" {
  value = module.cloud_run_api.service_name
}

output "cloud_run_uri" {
  value = module.cloud_run_api.uri
}

output "github_actions_service_account" {
  value = module.github_oidc.service_account_email
}

output "workload_identity_provider" {
  value = module.github_oidc.provider_name
}
