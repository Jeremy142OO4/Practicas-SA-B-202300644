output "namespace" {
  description = "Namespace administrado para P8."
  value       = kubernetes_namespace.p8.metadata[0].name
}

output "service_accounts" {
  description = "Cuentas de servicio creadas para los workloads."
  value       = sort(keys(kubernetes_service_account.workload))
}

