resource "google_service_account" "gke_node_sa" {
  account_id   = "gke-node-sa"
  display_name = "GKE Node Service Account"
}

resource "google_project_iam_member" "node_roles" {
  for_each = toset([
    "roles/compute.storageAdmin",
    "roles/container.defaultNodeServiceAccount"
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.gke_node_sa.email}"
}

resource "google_service_account_iam_member" "allow_tf_runner_to_use_gke_node_sa" {
  service_account_id = "projects/${var.project_id}/serviceAccounts/${google_service_account.gke_node_sa.email}"
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.tf_runner.email}"
}


resource "google_container_cluster" "game_cluster" {
  name               = var.cluster_name
  location           = var.zone
  initial_node_count = 1

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  network    = google_compute_network.game_vpc.name
  subnetwork = google_compute_subnetwork.game_subnet.name

  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  private_cluster_config {
    enable_private_nodes    = false
    enable_private_endpoint = false
  }

  node_config {
    service_account = google_service_account.gke_node_sa.email
    machine_type = "e2-standard-4"
    disk_size_gb = 40
    disk_type    = "pd-ssd"

    tags = ["game-server", "matchmaking", "database"]

    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }

  cluster_autoscaling {
    enabled = false
  }

  logging_service    = "none"
  monitoring_service = "none"

  addons_config {
    http_load_balancing {
      disabled = false
    }
    horizontal_pod_autoscaling {
      disabled = true
    }
  }
}