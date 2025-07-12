variable "project_id" {
  description = "Project ID"
  type = string
}

variable "region" {
  description = "GCP Region"
  type = string
  default = "europe-central2"
}

variable "db_password" {
  description = "DB password"
  type = string
  sensitive = true
}