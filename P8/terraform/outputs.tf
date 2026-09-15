output "cluster_name" {
  description = "Nombre del clúster GKE creado por Terraform."
  value       = google_container_cluster.p8.name
}

output "cluster_location" {
  description = "Zona del clúster GKE creado por Terraform."
  value       = google_container_cluster.p8.location
}

output "namespace" {
  description = "Namespace administrado para P8."
  value       = kubernetes_namespace.p8.metadata[0].name
}

output "service_accounts" {
  description = "Cuentas de servicio creadas para los workloads."
  value       = sort(keys(kubernetes_service_account.workload))
}
