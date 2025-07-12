variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
  default     = "game"
}

variable "db_username" {
  description = "Database username"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type = string
  default = "users"
}