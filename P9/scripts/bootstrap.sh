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

for command_name in gcloud terraform kubectl helm curl; do
  command -v "$command_name" >/dev/null || {
    echo "Falta la herramienta: $command_name" >&2
    exit 1
  }
done

if [ ! -f "$TFVARS" ]; then
  echo "Falta $TFVARS. Copie terraform.tfvars.example y complete sus valores." >&2
  exit 1
fi

export CLOUDSDK_CORE_PROJECT="$PROJECT_ID"
export CLOUDSDK_CORE_DISABLE_PROMPTS=1

# La creación del estado usa la API de GCS directamente. Así se evita que una
# espera del cliente gcloud deje el bootstrap detenido antes de Terraform.
gcs_token() {
  gcloud auth print-access-token
}

gcs_request() {
  local method="$1"
  local url="$2"
  shift 2
  local token
  token="$(gcs_token)"
  curl --fail --silent --show-error -X "$method" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    "$@" "$url"
}

bucket_status="$(curl --silent --show-error \
  -o /dev/null -w '%{http_code}' \
  -H "Authorization: Bearer $(gcs_token)" \
  "https://storage.googleapis.com/storage/v1/b/$STATE_BUCKET")"

if [ "$bucket_status" = "404" ]; then
  gcs_request POST "https://storage.googleapis.com/storage/v1/b?project=$PROJECT_ID" \
    --data "{\"name\":\"$STATE_BUCKET\",\"location\":\"$STATE_LOCATION\",\"iamConfiguration\":{\"uniformBucketLevelAccess\":{\"enabled\":true}}}" \
    >/dev/null
elif [ "$bucket_status" != "200" ]; then
  echo "No se pudo consultar el bucket de estado (HTTP $bucket_status)." >&2
  exit 1
fi

gcs_request PATCH "https://storage.googleapis.com/storage/v1/b/$STATE_BUCKET" \
  --data '{"versioning":{"enabled":true}}' >/dev/null

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
