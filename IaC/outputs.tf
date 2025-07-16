/*output "tf_service_account" {
  value = google_service_account.tf_runner.email
  description = "Service account for automation"
}

output "enabled_apis" {
  value = [for s in google_project_service.services : s.service]
  description = "List of enabled APIs"
}*/