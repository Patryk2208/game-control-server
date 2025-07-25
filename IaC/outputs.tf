/*output "tf_service_account" {
  value = google_service_account.tf_runner.email
  description = "Service account for automation"
}

output "enabled_apis" {
  value = [for s in google_project_service.services : s.service]
  description = "List of enabled APIs"
}*/

/*output "cluster_name" {
  value = var.cluster_name
  description = "Cluster name"
}

output "zone" {
  value = var.zone
  description = "Zone"
}

output "gcp_github_service_account_email" {
  value       = google_service_account.github_deployer.email
  description = "Email of the CI Service Account"
}*/

/*output "gcp_wid_ci_provider_cs" {
  value = google_iam_workload_identity_pool_provider.github_ci_provider_control_server.name
  description = "CS CI"
}

output "gcp_wid_cd_provider_cs" {
  value = google_iam_workload_identity_pool_provider.github_cd_provider_control_server.name
  description = "CS CD"
}

output "gcp_wid_ci_provider_gs" {
  value = google_iam_workload_identity_pool_provider.github_ci_provider_gameserver.name
  description = "GS CI"
}

output "gcp_wid_cd_provider_gs" {
  value = google_iam_workload_identity_pool_provider.github_cd_provider_gameserver.name
  description = "GS CD"
}*/