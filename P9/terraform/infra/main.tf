module "base" {
  # Reutiliza la infraestructura de P8 sin copiar el código.
  source = "../../../P8/terraform"

  project_id                     = var.project_id
  region                         = var.region
  zone                           = var.zone
  cluster_name                   = var.cluster_name
  network                        = var.network
  subnetwork                     = var.subnetwork
  cluster_ipv4_cidr              = var.cluster_ipv4_cidr
  services_ipv4_cidr             = var.services_ipv4_cidr
  master_authorized_network_cidr = var.master_authorized_network_cidr
  node_count                     = var.node_count
  namespace                      = var.namespace
  jwt_secret                     = var.jwt_secret
  encryption_key                 = var.encryption_key
}
