# Service Account
resource "kubernetes_service_account" "agones_sdk" {
  metadata {
    name = "agones-sdk"
    namespace = var.namespace
  }
}

# Role Binding
resource "kubernetes_role_binding" "agones_sdk" {
  metadata {
    name = "agones-sdk-bind"
    namespace = var.namespace
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.agones_sdk.metadata[0].name
    namespace = var.namespace
  }
  role_ref {
    kind      = "ClusterRole"
    name      = "agones-sdk"
    api_group = "rbac.authorization.k8s.io"
  }
}

# Fleet
resource "kubernetes_manifest" "game_server_fleet" {
  manifest = {
    apiVersion = "agones.dev/v1"
    kind       = "Fleet"
    metadata = {
      name      = "game-server-fleet"
      namespace = var.namespace
    }
    spec = {
      replicas = var.replicas
      template = {
        spec = {
          ports = [
            {
              name          = "game"
              containerPort = 7777
              portPolicy    = "Dynamic"
              protocol      = "TCP"
            },
            {
              name       = "agones-sdk"
              containerPort = 9357
              protocol   = "TCP"
            }
          ]
          template = {
            spec = {
              containers = [
                {
                  name  = "game-server"
                  image = var.game_server_image
                  resources = {
                    limits = {
                      cpu    = "4"
                      memory = "1Gi"
                    }
                  }
                }
              ]
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "agones-autoscaler" {
  manifest = {
    apiVersion = "autoscaling.agones.dev/v1"
    kind = "FleetAutoscaler"
    metadata = {
      name = "game-autoscaler"
      namespace = "game"
    }
    spec = {
      fleetName = "game-server-fleet"
    }
    policy = {
      type = "Buffer"
      buffer = {
        bufferSize = "20%"
        maxReplicas = 10
        minReplicas = 3
      }
    }
  }
}