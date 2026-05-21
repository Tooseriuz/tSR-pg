variable "pool_id" {
  type        = string
  description = "Workload identity pool id."
}

output "pool_id" {
  value       = var.pool_id
  description = "Placeholder workload identity pool id."
}
