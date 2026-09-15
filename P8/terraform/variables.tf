variable "project_id" {
  description = "Proyecto de Google Cloud donde se crea el clúster de P8."
  type        = string
  default     = "sa-202300644"
}

variable "region" {
  description = "Región de Google Cloud para recursos regionales."
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "Zona de Google Cloud para el clúster zonal."
  type        = string
  default     = "us-central1-a"
}

variable "cluster_name" {
  description = "Nombre del clúster GKE administrado por Terraform."
  type        = string
  default     = "p7-202300644"
}

variable "network" {
  description = "Red VPC del clúster."
  type        = string
  default     = "default"
}

variable "subnetwork" {
  description = "Subred del clúster."
  type        = string
  default     = "default"
}

variable "cluster_ipv4_cidr" {
  description = "CIDR secundario para pods."
  type        = string
  default     = "10.6.0.0/17"
}

variable "services_ipv4_cidr" {
  description = "CIDR secundario para servicios."
  type        = string
  default     = "34.118.224.0/20"
}

variable "node_count" {
  description = "Cantidad de nodos del pool principal."
  type        = number
  default     = 2
}

variable "jwt_secret" {
  description = "Secreto JWT para el servicio de autenticación; solo local."
  type        = string
  sensitive   = true
}

variable "encryption_key" {
  description = "Clave de cifrado para el servicio de autenticación; solo local."
  type        = string
  sensitive   = true
}

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
