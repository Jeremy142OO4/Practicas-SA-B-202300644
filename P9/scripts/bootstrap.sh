#!/usr/bin/env bash
set -Eeuo pipefail

PROJECT_ID="${PROJECT_ID:-sa-202300644}"
REGION="${REGION:-us-central1}"
ZONE="${ZONE:-us-central1-a}"
CLUSTER_NAME="${CLUSTER_NAME:-p9-202300644}"
STATE_BUCKET="${STATE_BUCKET:-sa-p9-202300644-tfstate}"
STATE_LOCATION="${STATE_LOCATION:-us-central1}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TFVARS="${TFVARS:-$ROOT_DIR/P9/terraform/infra/terraform.tfvars}"
GITOPS_DIR="${GITOPS_DIR:-$HOME/Practicas-SA-B-202300644-gitops}"
ROOT_APPLICATION="${ROOT_APPLICATION:-$GITOPS_DIR/p9/argocd/root.yaml}"

for command_name in gcloud terraform kubectl helm; do
  command -v "$command_name" >/dev/null || {
    echo "Falta la herramienta: $command_name" >&2
    exit 1
  }
done

if [ ! -f "$TFVARS" ]; then
  echo "Falta $TFVARS. Copie terraform.tfvars.example y complete sus valores." >&2
  exit 1
fi

gcloud config set project "$PROJECT_ID" >/dev/null

if ! gcloud storage buckets describe "gs://$STATE_BUCKET" >/dev/null 2>&1; then
  gcloud storage buckets create "gs://$STATE_BUCKET" \
    --project="$PROJECT_ID" \
    --location="$STATE_LOCATION" \
    --uniform-bucket-level-access
fi

gcloud storage buckets update "gs://$STATE_BUCKET" --versioning

terraform -chdir="$ROOT_DIR/P9/terraform/infra" init -reconfigure \
  -backend-config="bucket=$STATE_BUCKET" \
  -backend-config="prefix=p9/infra"

# Primera fase: crear únicamente GKE. La segunda fase aplica los recursos
# Kubernetes del módulo cuando el endpoint del clúster ya existe.
terraform -chdir="$ROOT_DIR/P9/terraform/infra" apply \
  -target=module.base.google_container_cluster.p8 \
  -target=module.base.google_container_node_pool.default \
  -var-file="$TFVARS" \
  -auto-approve

terraform -chdir="$ROOT_DIR/P9/terraform/infra" apply \
  -var-file="$TFVARS" \
  -auto-approve

gcloud container clusters get-credentials "$CLUSTER_NAME" \
  --zone "$ZONE" \
  --project "$PROJECT_ID"

terraform -chdir="$ROOT_DIR/P9/terraform/platform" init -reconfigure \
  -backend-config="bucket=$STATE_BUCKET" \
  -backend-config="prefix=p9/platform"

terraform -chdir="$ROOT_DIR/P9/terraform/platform" apply \
  -var="project_id=$PROJECT_ID" \
  -var="region=$REGION" \
  -var="zone=$ZONE" \
  -var="cluster_name=$CLUSTER_NAME" \
  -auto-approve

if [ ! -f "$ROOT_APPLICATION" ]; then
  echo "No se encuentra el Application raíz: $ROOT_APPLICATION" >&2
  echo "Clone o actualice el repositorio GitOps antes de ejecutar el bootstrap." >&2
  exit 1
fi

kubectl apply -f "$ROOT_APPLICATION"

echo "Infraestructura, ArgoCD y Application raíz instalados."
echo "ArgoCD continuará levantando los componentes de P9 desde GitOps."
