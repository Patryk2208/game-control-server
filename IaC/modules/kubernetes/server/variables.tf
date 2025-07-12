variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
  default     = "game"
}

variable "server_image" {
  description = "Matchmaker container image"
  type        = string
  default     = "test-game-server:local"
}

variable "db_secret_name" {
  description = "Database secret name"
  type        = string
  default     = "postgres-basic-auth"
}