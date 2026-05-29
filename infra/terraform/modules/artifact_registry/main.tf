variable "repository_id" {
  type        = string
  description = "Artifact Registry repository id."
}

variable "location" {
  type        = string
  description = "Artifact Registry location."
}

variable "description" {
  type        = string
  description = "Repository description."
  default     = "Docker images."
}

resource "google_artifact_registry_repository" "docker" {
  location      = var.location
  repository_id = var.repository_id
  description   = var.description
  format        = "DOCKER"
}

output "repository_id" {
  value       = google_artifact_registry_repository.docker.repository_id
  description = "Artifact Registry repository id."
}

output "repository_name" {
  value       = google_artifact_registry_repository.docker.name
  description = "Artifact Registry repository resource name."
}
