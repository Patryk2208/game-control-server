resource "google_project_service" "services" {
  for_each = toset([
    "container.googleapis.com",   # GKE
    "compute.googleapis.com",     # Networking
    //"gameservices.googleapis.com",# Game servers
  ])
  service            = each.key
  disable_on_destroy = false
}

resource "google_service_account" "tf_runner" {
  account_id   = "tf-runner"
  display_name = "Terraform Deployment Account"
}

resource "google_project_iam_member" "sa_roles" {
  for_each = toset([
    "roles/resourcemanager.projectIamAdmin",
    "roles/container.admin",        # GKE management
    "roles/compute.networkAdmin",   # VPC/firewall rules
    "roles/iam.serviceAccountUser",    # Bind service accounts to resources
    "roles/serviceusage.serviceUsageAdmin", # Enable/disable services
    "roles/logging.admin",             # For cluster logging
    "roles/monitoring.admin",          # For cluster monitoring
    "roles/compute.securityAdmin" # For firewall management
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.tf_runner.email}"
}