resource "google_service_account" "github_deployer" {
  account_id   = "github-deployer"
  display_name = "GitHub Actions Deployment Account"
}

resource "google_project_iam_member" "deployer_roles" {
  for_each = toset([
    "roles/container.developer",
    "roles/storage.admin",
    "roles/artifactregistry.writer"
  ])

  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.github_deployer.email}"
}

resource "google_service_account_iam_member" "workload_identity" {
  service_account_id = google_service_account.github_deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/${var.github_owner}/${var.github_repo}"
}

locals {
  attribute_mapping = {
    "google.subject"      = "assertion.sub"
    "attribute.actor"     = "assertion.actor"
    "attribute.repository"= "assertion.repository"
    "attribute.ref"       = "assertion.ref"
    "attribute.workflow"  = "assertion.workflow"
  }
}

resource "google_iam_workload_identity_pool" "github_pool" {
  workload_identity_pool_id = "github-pool"
  display_name              = "GitHub Actions Pool"
  project = var.project_id
}

resource "google_iam_workload_identity_pool_provider" "github_ci_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "ci-provider"
  display_name                       = "GH Actions CI Provider"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = local.attribute_mapping

  attribute_condition = <<-EOT
    assertion.repository == "${var.github_owner}/${var.github_repo}" &&
    assertion.workflow == "${var.ci_workflow_file}" &&
    assertion.ref == "${var.github_provider_branch}"
  EOT
}

resource "google_iam_workload_identity_pool_provider" "github_cd_provider" {
  workload_identity_pool_id = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "cd-provider"
  display_name = "GH Actions CD Provider"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = local.attribute_mapping

  attribute_condition = <<-EOT
    assertion.repository == "${var.github_owner}/${var.github_repo}" &&
    assertion.workflow == "${var.cd_workflow_file}" &&
    assertion.ref == "${var.github_provider_branch}"
  EOT
}

resource "google_artifact_registry_repository" "game_server" {
  location      = var.region
  repository_id = "game-repo"
  description   = "Docker repository for game containers"
  format        = "DOCKER"

  docker_config {
    immutable_tags = true
  }
}