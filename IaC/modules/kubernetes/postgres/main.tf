resource "kubernetes_secret" "postgres_basic_auth" {
  metadata {
    name = "postgres-basic-auth"
    namespace = var.namespace
  }
  type = "kubernetes.io/basic-auth"
  data = {
    username = var.db_username
    password = var.db_password
  }
}

resource "kubernetes_stateful_set" "postgres" {
  metadata {
    name = "postgres"
    namespace = var.namespace
    labels = {
      app = "postgres"
    }
  }
  spec {
    service_name = "postgres-service"
    replicas = 1
    selector {
      match_labels = {
        app = "postgres"
      }
    }
    template {
      metadata {
        labels = {
          app = "postgres"
        }
      }
      spec {
        container {
          name  = "postgres"
          image = "postgres:15"
          image_pull_policy = "IfNotPresent"
          port {
            container_port = 5432
          }
          env {
            name = "POSTGRES_DB"
            value = var.db_name
          }
          env {
            name = "POSTGRES_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_basic_auth.metadata[0].name
                key  = "username"
              }
            }
          }
          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_basic_auth.metadata[0].name
                key  = "password"
              }
            }
          }
          volume_mount {
            name       = "postgres-data"
            mount_path = "/var/lib/postgresql/data"
          }
          resources {
            requests = {
              cpu    = "100m"
              memory = "256Mi"
            }
          }
        }
        volume {
          name = "postgres-data"
          persistent_volume_claim {
            claim_name = "test-postgres-claim"
          }
        }
      }
    }
    volume_claim_template {
      metadata {
        name = "test-postgres-claim"
      }
      spec {
        access_modes       = ["ReadWriteOnce"]
        storage_class_name = "standard"
        resources {
          requests = {
            storage = "1Gi"
          }
        }
      }
    }
  }
}

# Service
resource "kubernetes_service" "postgres_service" {
  metadata {
    name = "postgres-service"
    namespace = var.namespace
  }
  spec {
    selector = {
      app = "postgres"
    }
    port {
      protocol    = "TCP"
      port        = 5432
      target_port = 5432
    }
    type = "NodePort"
  }
}