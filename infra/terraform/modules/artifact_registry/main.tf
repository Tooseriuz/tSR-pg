variable "repository_id" {
  type        = string
  description = "Artifact Registry repository id."
}

output "repository_id" {
  value       = var.repository_id
  description = "Placeholder Artifact Registry repository id."
}
