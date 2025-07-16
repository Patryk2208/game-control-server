resource "google_container_cluster" "game_cluster" {
  name     = "game-platform-cluster"
  location = var.region


  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "game_nodes" {
  name       = "game-node-pool"
  cluster    = google_container_cluster.game_cluster.name
  location   = var.region
  node_count = 1

  node_config {
    machine_type = "e2-micro"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring"
    ]
  }
}

resource "google_sql_database_instance" "postgres" {
  name             = "game-db-instance"
  database_version = "POSTGRES_15"
  region           = var.region
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = false
      private_network = "default"
    }
  }
}

resource "google_sql_database" "users_db" {
  name     = "users"
  instance = google_sql_database_instance.postgres.name
}

resource "google_sql_user" "dev_user" {
  name     = "devuser"
  instance = google_sql_database_instance.postgres.name
  password = var.db_password
}


resource "helm_release" "agones" {
  name       = "agones"
  repository = "https://agones.dev/chart/stable"
  chart      = "agones"
  version    = "1.32.0"
  namespace  = "agones-system"
  create_namespace = true

  set {
    name  = "agones.ping.http.enabled"
    value = "false"
  }
  set {
    name = "agones.crds.install"
    value = "true"
  }
  set {
    name = "agones.featureGates=PlayerTracking"
    value = "true"
  }

  depends_on = [google_container_node_pool.game_nodes]
}

provider "kubernetes" {
  host  = "https://${google_container_cluster.game_cluster.endpoint}"
  token = data.google_client_config.current.access_token
  cluster_ca_certificate = base64decode(
    google_container_cluster.game_cluster.master_auth[0].cluster_ca_certificate
  )
}

data "google_client_config" "current" {}