variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
  default     = "game"
}

variable "game_server_image" {
  description = "Game server container image"
  type        = string
  default     = "test-rpg-gameserver:local"
}

variable "replicas" {
  description = "Number of game server replicas"
  type        = number
  default     = 1
}