variable "pool_id" {
  type        = string
  description = "Workload identity pool id."
}

variable "provider_id" {
  type        = string
  description = "Workload identity provider id."
}

variable "project_id" {
  type        = string
  description = "Google Cloud project id."
}

variable "project_number" {
  type        = string
  description = "Google Cloud project number."
}

variable "github_repository" {
  type        = string
  description = "GitHub repository in owner/name format."
}

variable "service_account_id" {
  type        = string
  description = "Service account id for GitHub Actions deployments."
}

variable "cloud_run_service_account_email" {
  type        = string
  description = "Runtime service account email used by Cloud Run."
}

resource "google_iam_workload_identity_pool" "github" {
  workload_identity_pool_id = var.pool_id
  display_name              = "GitHub Actions"
}

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = var.provider_id
  display_name                       = "GitHub"

  attribute_mapping = {
    "google.subject"             = "assertion.sub"
    "attribute.actor"            = "assertion.actor"
    "attribute.repository"       = "assertion.repository"
    "attribute.repository_owner" = "assertion.repository_owner"
  }

  attribute_condition = "assertion.repository == '${var.github_repository}'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account" "github_actions" {
  account_id   = var.service_account_id
  display_name = "GitHub Actions deployer"
}

resource "google_service_account_iam_member" "workload_identity_user" {
  service_account_id = google_service_account.github_actions.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/projects/${var.project_number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.github.workload_identity_pool_id}/attribute.repository/${var.github_repository}"
}

resource "google_project_iam_member" "artifact_registry_writer" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

resource "google_project_iam_member" "cloud_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

resource "google_service_account_iam_member" "cloud_run_service_account_user" {
  service_account_id = "projects/${var.project_id}/serviceAccounts/${var.cloud_run_service_account_email}"
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_actions.email}"
}

output "pool_id" {
  value       = google_iam_workload_identity_pool.github.workload_identity_pool_id
  description = "Workload identity pool id."
}

output "provider_name" {
  value       = google_iam_workload_identity_pool_provider.github.name
  description = "Workload identity provider resource name."
}

output "service_account_email" {
  value       = google_service_account.github_actions.email
  description = "GitHub Actions deployer service account email."
}
