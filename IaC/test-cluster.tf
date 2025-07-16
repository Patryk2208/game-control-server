resource "google_container_cluster" "test_cluster" {
  name     = "network-test-cluster"
  location = var.zone
  deletion_protection = false

  initial_node_count = 1
  node_config {
    machine_type = "e2-micro"
    disk_size_gb = 20

    tags = ["test-node", "ssh"]

    service_account = google_service_account.gke_node.email
  }

  network    = google_compute_network.game_vpc.name
  subnetwork = google_compute_subnetwork.game_subnet.name

  addons_config {
    http_load_balancing {
      disabled = true
    }
    horizontal_pod_autoscaling {
      disabled = true
    }
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }
}

resource "google_service_account" "gke_node" {
  account_id   = "gke-node-sa"
  display_name = "GKE Node Service Account"
}

resource "google_project_iam_member" "node_roles" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/container.nodeServiceAccount"
  ])
  project = var.project_id
  role   = each.key
  member = "serviceAccount:${google_service_account.gke_node.email}"
}