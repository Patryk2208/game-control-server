# Service Account
resource "kubernetes_service_account" "matchmaker" {
  metadata {
    name = "matchmaker-service-account"
    namespace = var.namespace
  }
}

# Role
resource "kubernetes_role" "agones_management" {
  metadata {
    name = "agones-management-role"
    namespace = var.namespace
  }
  rule {
    api_groups = ["allocation.agones.dev"]
    resources  = ["gameserverallocations"]
    verbs      = ["create"]
  }
  rule {
    api_groups = ["agones.dev"]
    resources  = ["gameservers"]
    verbs      = ["get", "watch", "list"]
  }
}

# Role Binding
resource "kubernetes_role_binding" "agones_management" {
  metadata {
    name = "agones-management-bind"
    namespace = var.namespace
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.matchmaker.metadata[0].name
    namespace = var.namespace
  }
  role_ref {
    kind      = "Role"
    name      = kubernetes_role.agones_management.metadata[0].name
    api_group = "rbac.authorization.k8s.io"
  }
}

# Deployment
resource "kubernetes_deployment" "matchmaker" {
  metadata {
    name = "server-pod"
    namespace = var.namespace
    labels = {
      app = "server"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "server-pod"
      }
    }
    template {
      metadata {
        labels = {
          app = "server-pod"
        }
      }
      spec {
        service_account_name = kubernetes_service_account.matchmaker.metadata[0].name
        container {
          name  = "server-container"
          image = var.server_image
          image_pull_policy = "Never"
          port {
            container_port = 8080
            protocol       = "TCP"
          }
          env {
            name  = "DB_HOST"
            value = "postgres-service.game.svc.cluster.local"
          }
          env {
            name  = "DB_PORT"
            value = "5432"
          }
          env {
            name  = "DB_NAME"
            value = "users"
          }
          env {
            name = "DB_USER"
            value_from {
              secret_key_ref {
                name = var.db_secret_name
                key  = "username"
              }
            }
          }
          env {
            name = "DB_PASSWORD"
            value_from {
              secret_key_ref {
                name = var.db_secret_name
                key  = "password"
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "matchmaker" {
  metadata {
    name = "server-service"
    namespace = var.namespace
  }
  spec {
    selector = {
      app = "server-pod"
    }
    port {
      protocol    = "TCP"
      port        = 8080
      target_port = 8080
    }
    type = "LoadBalancer"
  }
}