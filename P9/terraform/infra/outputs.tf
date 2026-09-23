output "cluster_name" {
  value = module.base.cluster_name
}

output "cluster_location" {
  value = module.base.cluster_location
}

output "namespace" {
  value = module.base.namespace
}

output "velero_bucket_name" {
  value = google_storage_bucket.velero.name
}

output "velero_service_account" {
  value = google_service_account.velero.email
}

output "external_secrets_service_account" {
  value = google_service_account.external_secrets.email
}
