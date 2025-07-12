resource "minikube_cluster" "game_dev" {
  cluster_name = "game-platform"
  driver       = "docker"  # Docker driver for Linux/macOS
  cpus         = var.minikube_cpus
  memory       = var.minikube_memory
  addons       = ["ingress", "metrics-server", "registry"]

  # Expose game server ports (TCP 7000-8000 range)
  extra_config = {
    "ports" = join(",", [for port in var.exposed_ports : "${port}:${port}"])
  }

  provisioner "local-exec" {
    command = <<-EOT
      helm repo add agones https://agones.dev/chart/stable
      helm install agones agones/agones \
        --namespace agones-system \
        --create-namespace \
        --set "agones.ping.http.enabled=false" \
        --set "agones.ping.udp.enabled=false"
    EOT
  }
}

/*resource "kubernetes_storage_class" "local_storage" {
  metadata {
    name = "local-storage"
  }
  storage_provisioner = "kubernetes.io/no-provisioner"
  volume_binding_mode = "WaitForFirstConsumer"
}*/

resource "kubernetes_persistent_volume" "postgres_pv" {
  metadata {
    name = "postgres-pv"
  }
  spec {
    capacity = {
      storage = "5Gi"
    }
    access_modes = ["ReadWriteOnce"]
    persistent_volume_source {
      host_path {
        path = "/data/postgres"
      }
    }
    storage_class_name = kubernetes_storage_class.local_storage.metadata[0].name
  }
}

resource "kubernetes_namespace" "game" {
  metadata {
    name = "game"
  }
}

resource "kubernetes_storage_class" "standard" {
  metadata {
    name = "standard"
  }
  storage_provisioner = "k8s.io/minikube-hostpath"
  reclaim_policy      = "Delete"
}


module "local_env" {
  source = "./environments/local"
}

# Deploy database
module "database" {
  source = "./modules/database"
  namespace   = kubernetes_namespace.game.metadata[0].name
  db_username = "patryk"       # Replace with actual values
  db_password = "sql"   # Replace with actual values
}

# Deploy matchmaker
module "matchmaker" {
  source = "./modules/matchmaker"
  namespace         = kubernetes_namespace.game.metadata[0].name
  matchmaker_image  = "test-game-server:local"
  db_secret_name    = module.database.db_secret_name
  depends_on        = [module.database]
}

# Deploy game servers
module "agones" {
  source = "./modules/agones"
  namespace       = kubernetes_namespace.game.metadata[0].name
  game_server_image = "test-rpg-gameserver:local"
  depends_on      = [module.local_env]
}