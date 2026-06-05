variable "bucket_name" {
  type        = string
  description = "GCS bucket name."
}

variable "location" {
  type        = string
  description = "GCS bucket location."
}

variable "cloud_run_service_account_email" {
  type        = string
  description = "Cloud Run runtime service account email."
}

variable "force_destroy" {
  type        = bool
  description = "Delete objects when destroying the bucket."
  default     = false
}

resource "google_storage_bucket" "files" {
  name                        = var.bucket_name
  location                    = var.location
  force_destroy               = var.force_destroy
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
}

resource "google_storage_bucket_iam_member" "cloud_run_object_user" {
  bucket = google_storage_bucket.files.name
  role   = "roles/storage.objectUser"
  member = "serviceAccount:${var.cloud_run_service_account_email}"
}

output "bucket_name" {
  value       = google_storage_bucket.files.name
  description = "GCS bucket name."
}

output "bucket_url" {
  value       = google_storage_bucket.files.url
  description = "GCS bucket URL."
}
