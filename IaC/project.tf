resource "google_project_service" "services" {
  for_each = toset([
    "container.googleapis.com",
    "compute.googleapis.com",
    "containerregistry.googleapis.com"
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
    "roles/iam.serviceAccountAdmin",    # Bind service accounts to resources
    "roles/serviceusage.serviceUsageAdmin", # Enable/disable services
    "roles/logging.admin",             # For cluster logging
    "roles/monitoring.admin",          # For cluster monitoring
    "roles/compute.securityAdmin", # For firewall management
    "roles/artifactregistry.admin",
    "roles/iam.workloadIdentityPoolAdmin"
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.tf_runner.email}"
}