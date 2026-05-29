variable "service_name" {
  type        = string
  description = "Cloud Run service name."
}

variable "location" {
  type        = string
  description = "Cloud Run location."
}

variable "image" {
  type        = string
  description = "Initial container image. CI owns updates after service creation."
  default     = "us-docker.pkg.dev/cloudrun/container/hello"
}

variable "allow_unauthenticated" {
  type        = bool
  description = "Allow unauthenticated requests to the Cloud Run service."
  default     = true
}

resource "google_service_account" "runtime" {
  account_id   = var.service_name
  display_name = "${var.service_name} runtime"
}

resource "google_cloud_run_v2_service" "api" {
  name     = var.service_name
  location = var.location

  template {
    service_account = google_service_account.runtime.email

    containers {
      image = var.image

      ports {
        container_port = 8080
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
    ]
  }
}

resource "google_cloud_run_v2_service_iam_member" "public_invoker" {
  count    = var.allow_unauthenticated ? 1 : 0
  name     = google_cloud_run_v2_service.api.name
  location = google_cloud_run_v2_service.api.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "service_name" {
  value       = google_cloud_run_v2_service.api.name
  description = "Cloud Run service name."
}

output "service_account_email" {
  value       = google_service_account.runtime.email
  description = "Cloud Run runtime service account email."
}

output "uri" {
  value       = google_cloud_run_v2_service.api.uri
  description = "Cloud Run service URI."
}
