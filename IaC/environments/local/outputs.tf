output "minikube_ip" {
  description = "IP address of Minikube cluster"
  value       = minikube_cluster.game_dev.ip_address
}

output "kubeconfig_path" {
  description = "Path to Kubernetes config file"
  value       = minikube_cluster.game_dev.kubeconfig
}

output "registry_endpoint" {
  description = "Local Docker registry endpoint"
  value       = "localhost:${minikube_cluster.game_dev.registry_port}"
}