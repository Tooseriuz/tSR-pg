variable "secret_ids" {
  type        = list(string)
  description = "Secret Manager secret ids."
  default     = []
}

output "secret_ids" {
  value       = var.secret_ids
  description = "Placeholder Secret Manager secret ids."
}
