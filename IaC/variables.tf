variable "project_id" {
  description = "Project ID"
  type = string
  sensitive = true
}

variable "region" {
  description = "GCP Region"
  type = string
  default = "europe-central2"
}

variable "zone" {
  description = "GCP zone"
  type        = string
  default     = "europe-central2-a"
}

variable "service_account_key_path" {
  description = "Path to terraform-init service account key file"
  type        = string
  sensitive   = true
}

variable "cluster_name" {
  description = "Gke cluster name"
  type = string
  default = "game-cluster"
}

variable "github_owner" {
  description = "Github repo owner name"
  type = string
  default = "Patryk2208"
}

variable "github_repo" {
  description = "Github project repo"
  type = string
  default = "game-control-server"
}

variable "ci_workflow_file" {
  description = "Path to CI.yml"
  type = string
  default = "CI.yml"
}

variable "cd_workflow_file" {
  description = "Path to CD.yml"
  type = string
  default = "CD.yml"
}

variable "github_provider_branch" {
  description = "branch for providers"
  type = string
  default = "refs/heads/main"
}