variable "namespace" {
  description = "Namespace administrado por Terraform para P8."
  type        = string
  default     = "sa-p8"
}

variable "kubeconfig_path" {
  description = "Ruta local al kubeconfig. Se proporciona fuera del repositorio."
  type        = string
  default     = null
}

variable "kube_context" {
  description = "Contexto opcional del kubeconfig."
  type        = string
  default     = null
}

variable "workload_service_accounts" {
  description = "Cuentas de servicio sin privilegios de los componentes de la plataforma."
  type        = set(string)
  default = [
    "api-gateway",
    "auth-ms",
    "catalog-ms",
    "orders-ms",
    "payments-ms",
    "scheduled-jobs",
  ]
}

