variable "service_name" {
  type        = string
  description = "Cloud Run service name."
}

output "service_name" {
  value       = var.service_name
  description = "Placeholder Cloud Run service name."
}
