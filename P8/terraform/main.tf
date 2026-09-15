resource "kubernetes_namespace" "p8" {
  metadata {
    name = var.namespace

    labels = {
      "app.kubernetes.io/part-of"          = "practica-8"
      "app.kubernetes.io/managed-by"       = "terraform"
      "pod-security.kubernetes.io/enforce" = "restricted"
    }
  }
}

resource "kubernetes_resource_quota" "p8" {
  metadata {
    name      = "sa-p8-quota"
    namespace = kubernetes_namespace.p8.metadata[0].name
  }

  spec {
    hard = {
      pods                   = "60"
      "requests.cpu"         = "4"
      "requests.memory"      = "6Gi"
      "limits.cpu"           = "12"
      "limits.memory"        = "16Gi"
      "requests.storage"     = "60Gi"
      persistentvolumeclaims = "10"
    }
  }
}

resource "kubernetes_limit_range" "p8" {
  metadata {
    name      = "sa-p8-limits"
    namespace = kubernetes_namespace.p8.metadata[0].name
  }

  spec {
    limit {
      type = "Container"

      default_request = {
        cpu    = "50m"
        memory = "64Mi"
      }

      default = {
        cpu    = "500m"
        memory = "512Mi"
      }

      min = {
        cpu    = "10m"
        memory = "32Mi"
      }

      max = {
        cpu    = "2"
        memory = "2Gi"
      }
    }
  }
}

resource "kubernetes_service_account" "workload" {
  for_each = var.workload_service_accounts

  metadata {
    name      = each.value
    namespace = kubernetes_namespace.p8.metadata[0].name
    labels = {
      "app.kubernetes.io/part-of"    = "practica-8"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_secret" "auth_source" {
  metadata {
    name      = "auth-ms-secret"
    namespace = kubernetes_namespace.p8.metadata[0].name
    labels = {
      "app.kubernetes.io/managed-by" = "terraform"
      "app.kubernetes.io/part-of"    = "practica-8"
    }
  }

  type = "Opaque"

  data = {
    JWT_SECRET     = var.jwt_secret
    ENCRYPTION_KEY = var.encryption_key
  }
}

resource "kubernetes_role" "workload_read" {
  metadata {
    name      = "p8-workload-read"
    namespace = kubernetes_namespace.p8.metadata[0].name
  }

  rule {
    api_groups = [""]
    resources  = ["configmaps"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_role_binding" "workload_read" {
  for_each = kubernetes_service_account.workload

  metadata {
    name      = "${each.key}-read"
    namespace = kubernetes_namespace.p8.metadata[0].name
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = kubernetes_role.workload_read.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = each.value.metadata[0].name
    namespace = kubernetes_namespace.p8.metadata[0].name
  }
}
