resource "kubernetes_namespace" "game" {
  metadata {
    name = "game"
  }
}

# Database secret module
module "database_secret" {
  source = "./modules/postgres"

  namespace   = kubernetes_namespace.game.metadata[0].name
  db_username = google_sql_user.game_user.name
  db_password = google_sql_user.game_user.password
}

module "matchmaker" {
  source = "./modules/server"

  namespace        = kubernetes_namespace.game.metadata[0].name
  matchmaker_image = "gcr.io/${var.project_id}/server:latest"
  db_secret_name   = module.database_secret.secret_name
  db_host          = google_sql_database_instance.postgres.private_ip_address
}

module "agones_fleet" {
  source = "./modules/agones"

  namespace         = kubernetes_namespace.game.metadata[0].name
  game_server_image = "gcr.io/${var.project_id}/game-server:latest"
}

/*# (Optional) Network exposure
resource "kubernetes_ingress_v1" "matchmaker_ingress" {
  metadata {
    name      = "matchmaker-ingress"
    namespace = kubernetes_namespace.game.metadata[0].name
    annotations = {
      "kubernetes.io/ingress.class" = "gce"
    }
  }
  spec {
    rule {
      http {
        path {
          path = "/*"
          backend {
            service {
              name = module.matchmaker.service_name
              port {
                number = 8080
              }
            }
          }
        }
      }
    }
  }
}*/