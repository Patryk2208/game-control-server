resource "google_service_account" "github_deployer" {
  account_id   = "github-deployer"
  display_name = "GitHub Actions Deployment Account"
}

resource "google_project_iam_member" "deployer_roles" {
  for_each = toset([
    "roles/container.developer",
    "roles/storage.admin",
    "roles/artifactregistry.admin"
  ])

  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.github_deployer.email}"
}

resource "google_service_account_iam_member" "workload_identity_control_server" {
  service_account_id = google_service_account.github_deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/${var.github_owner}/${var.github_control_server_repo}"
}

resource "google_service_account_iam_member" "workload_identity_gameserver" {
  service_account_id = google_service_account.github_deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/${var.github_owner}/${var.github_gameserver_repo}"
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

resource "google_iam_workload_identity_pool_provider" "github_ci_provider_control_server" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "ci-provider-control-server"
  display_name                       = "GHA CI for control-server"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = local.attribute_mapping

  attribute_condition = <<-EOT
    assertion.repository == "${var.github_owner}/${var.github_control_server_repo}" &&
    assertion.workflow == "${var.ci_workflow}" &&
    assertion.ref == "${var.github_provider_branch}"
  EOT
}

resource "google_iam_workload_identity_pool_provider" "github_cd_provider_control_server" {
  workload_identity_pool_id = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "cd-provider-control-server"
  display_name = "GHA CD for control-server"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = local.attribute_mapping

  attribute_condition = <<-EOT
    assertion.repository == "${var.github_owner}/${var.github_control_server_repo}" &&
    assertion.workflow == "${var.cd_workflow}" &&
    assertion.ref == "${var.github_provider_branch}"
  EOT
}

resource "google_iam_workload_identity_pool_provider" "github_ci_provider_gameserver" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "ci-provider-gameserver"
  display_name                       = "GHA CI for gameserver"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = local.attribute_mapping

  attribute_condition = <<-EOT
    assertion.repository == "${var.github_owner}/${var.github_gameserver_repo}" &&
    assertion.workflow == "${var.ci_workflow}" &&
    assertion.ref == "${var.github_provider_branch}"
  EOT
}

resource "google_iam_workload_identity_pool_provider" "github_cd_provider_gameserver" {
  workload_identity_pool_id = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "cd-provider-gameserver"
  display_name = "GHA CD for gameserver"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = local.attribute_mapping

  attribute_condition = <<-EOT
    assertion.repository == "${var.github_owner}/${var.github_gameserver_repo}" &&
    assertion.workflow == "${var.cd_workflow}" &&
    assertion.ref == "${var.github_provider_branch}"
  EOT
}


resource "google_artifact_registry_repository" "game_server" {
  location      = var.region
  repository_id = var.gar_repo_name
  description   = "Docker repository for game containers"
  format        = "DOCKER"

  docker_config {
    immutable_tags = false
  }
}