output "argocd_namespace" {
  value = var.argocd_namespace
}

output "argocd_release" {
  value = helm_release.argocd.name
}

