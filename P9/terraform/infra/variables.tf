variable "project_id" {
  type    = string
  default = "sa-202300644"
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "zone" {
  type    = string
  default = "us-central1-a"
}

variable "cluster_name" {
  type    = string
  default = "p9-202300644"
}

variable "network" {
  type    = string
  default = "default"
}

variable "subnetwork" {
  type    = string
  default = "default"
}

variable "cluster_ipv4_cidr" {
  type    = string
  default = "10.7.0.0/17"
}

variable "services_ipv4_cidr" {
  type    = string
  default = "34.118.240.0/20"
}

variable "master_authorized_network_cidr" {
  type    = string
  default = "0.0.0.0/0"
}

variable "node_count" {
  type    = number
  default = 2
}

variable "namespace" {
  type    = string
  default = "sa-p9"
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "encryption_key" {
  type      = string
  sensitive = true
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "orders_database_url" {
  type      = string
  sensitive = true
}

variable "rabbitmq_url" {
  type      = string
  sensitive = true
}

variable "rabbitmq_password" {
  type      = string
  sensitive = true
}

variable "velero_bucket_name" {
  type    = string
  default = null
}
