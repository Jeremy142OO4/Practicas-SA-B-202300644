locals {
  velero_bucket_name = coalesce(
    var.velero_bucket_name,
    "${var.project_id}-${var.cluster_name}-velero",
  )
  external_secrets_gsa = "external-secrets-${var.cluster_name}"
  velero_gsa           = "velero-${var.cluster_name}"
}

resource "google_storage_bucket" "velero" {
  name                        = local.velero_bucket_name
  location                    = var.region
  project                     = var.project_id
  uniform_bucket_level_access = true
  force_destroy               = false

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      age = 30
    }

    action {
      type = "Delete"
    }
  }
}

resource "google_service_account" "velero" {
  account_id   = local.velero_gsa
  display_name = "Velero P9 backup service account"
  project      = var.project_id
}

resource "google_storage_bucket_iam_member" "velero" {
  bucket = google_storage_bucket.velero.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.velero.email}"
}

resource "google_service_account_iam_member" "velero_workload_identity" {
  service_account_id = google_service_account.velero.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[velero/velero-server]"
}

resource "google_service_account" "external_secrets" {
  account_id   = local.external_secrets_gsa
  display_name = "External Secrets P9 service account"
  project      = var.project_id
}

resource "google_project_iam_member" "external_secrets_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.external_secrets.email}"
}

resource "google_service_account_iam_member" "external_secrets_workload_identity" {
  service_account_id = google_service_account.external_secrets.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[external-secrets/external-secrets]"
}

resource "google_secret_manager_secret" "db_password" {
  secret_id = "${var.cluster_name}-db-password"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = var.db_password
}

resource "google_secret_manager_secret" "jwt_secret" {
  secret_id = "${var.cluster_name}-jwt-secret"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "jwt_secret" {
  secret      = google_secret_manager_secret.jwt_secret.id
  secret_data = var.jwt_secret
}

resource "google_secret_manager_secret" "encryption_key" {
  secret_id = "${var.cluster_name}-encryption-key"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "encryption_key" {
  secret      = google_secret_manager_secret.encryption_key.id
  secret_data = var.encryption_key
}

resource "google_secret_manager_secret" "orders_database_url" {
  secret_id = "${var.cluster_name}-orders-database-url"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "orders_database_url" {
  secret      = google_secret_manager_secret.orders_database_url.id
  secret_data = var.orders_database_url
}

resource "google_secret_manager_secret" "rabbitmq_url" {
  secret_id = "${var.cluster_name}-rabbitmq-url"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "rabbitmq_url" {
  secret      = google_secret_manager_secret.rabbitmq_url.id
  secret_data = var.rabbitmq_url
}

resource "google_secret_manager_secret" "rabbitmq_password" {
  secret_id = "${var.cluster_name}-rabbitmq-password"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "rabbitmq_password" {
  secret      = google_secret_manager_secret.rabbitmq_password.id
  secret_data = var.rabbitmq_password
}
