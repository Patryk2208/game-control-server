variable "minikube_cpus" {
  description = "Number of CPUs allocated to Minikube VM"
  type        = number
  default     = 4
}

variable "minikube_memory" {
  description = "Amount of memory allocated to Minikube VM"
  type        = string
  default     = "8g"
}

variable "exposed_ports" {
  description = "TCP ports to expose for game servers"
  type        = list(number)
  default     = [7000, 7100, 7200, 7300]  # todo
}