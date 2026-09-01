#!/usr/bin/env bash
set -Eeuo pipefail

directorio_script="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
raiz_p6="$(cd "$directorio_script/.." && pwd)"
raiz_repositorio="$(cd "$directorio_script/../.." && pwd)"

confirmar_herramienta() {
  local nombre="$1"
  command -v "$nombre" >/dev/null 2>&1 || { echo "Falta la herramienta: $nombre" >&2; exit 1; }
}

crear_clave_hexadecimal() {
  local cantidad_bytes="${1:-32}"
  openssl rand -hex "$cantidad_bytes"
}

crear_infraestructura() {
  local region="${1:-us-east-1}"
  local nombre_cluster="${2:-sa-p6}"
  local tipo_instancia="${3:-t3.medium}"
  local nombre_rol="AmazonEKS_EBS_CSI_DriverRole-${nombre_cluster}"
  local id_cuenta arn_rol estado_addon
  confirmar_herramienta aws
  confirmar_herramienta eksctl
  confirmar_herramienta kubectl

  aws sts get-caller-identity >/dev/null
  if ! aws eks describe-cluster --name "$nombre_cluster" --region "$region" >/dev/null 2>&1; then
    eksctl create cluster \
      --name "$nombre_cluster" \
      --region "$region" \
      --managed \
      --nodegroup-name nodos-p6 \
      --node-type "$tipo_instancia" \
      --nodes 2 \
      --nodes-min 2 \
      --nodes-max 2 \
      --node-volume-size 30
  fi
  aws eks update-kubeconfig --name "$nombre_cluster" --region "$region"
  eksctl utils associate-iam-oidc-provider --cluster "$nombre_cluster" --region "$region" --approve

  id_cuenta="$(aws sts get-caller-identity --query Account --output text)"
  arn_rol="arn:aws:iam::${id_cuenta}:role/${nombre_rol}"
  if ! aws iam get-role --role-name "$nombre_rol" >/dev/null 2>&1; then
    eksctl create iamserviceaccount \
      --name ebs-csi-controller-sa \
      --namespace kube-system \
      --cluster "$nombre_cluster" \
      --region "$region" \
      --role-name "$nombre_rol" \
      --role-only \
      --attach-policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicyV2 \
      --approve
  fi
  if ! aws eks describe-addon --cluster-name "$nombre_cluster" --addon-name aws-ebs-csi-driver --region "$region" >/dev/null 2>&1; then
    aws eks create-addon \
      --cluster-name "$nombre_cluster" \
      --addon-name aws-ebs-csi-driver \
      --service-account-role-arn "$arn_rol" \
      --resolve-conflicts OVERWRITE \
      --region "$region" >/dev/null
  fi
  for _ in {1..30}; do
    estado_addon="$(aws eks describe-addon --cluster-name "$nombre_cluster" --addon-name aws-ebs-csi-driver --region "$region" --query 'addon.status' --output text)"
    [[ "$estado_addon" == "ACTIVE" ]] && break
    [[ "$estado_addon" == "DEGRADED" || "$estado_addon" == "CREATE_FAILED" ]] && { echo "El complemento EBS CSI fallo: $estado_addon" >&2; exit 1; }
    sleep 10
  done
  [[ "$estado_addon" == "ACTIVE" ]] || { echo "EBS CSI no alcanzo el estado ACTIVE." >&2; exit 1; }
  kubectl get nodes -o wide
  kubectl -n kube-system get pods -l app.kubernetes.io/name=aws-ebs-csi-driver
}

publicar_imagenes() {
  local usuario_dockerhub="${1:-jeremy142004}"
  local etiqueta="${2:-1.0.0}"
  local repositorio_dockerhub="$usuario_dockerhub/sa-p6"
  local herramienta_contenedores
  if command -v podman >/dev/null 2>&1; then herramienta_contenedores=podman; else herramienta_contenedores=docker; fi
  confirmar_herramienta "$herramienta_contenedores"
  if [[ -z "${DOCKERHUB_TOKEN:-}" ]]; then
    read -r -s -p "Token nuevo de Docker Hub: " DOCKERHUB_TOKEN
    echo
  fi
  printf '%s' "$DOCKERHUB_TOKEN" | "$herramienta_contenedores" login --username "$usuario_dockerhub" --password-stdin
  trap '"$herramienta_contenedores" logout >/dev/null 2>&1 || true; unset DOCKERHUB_TOKEN' EXIT

  local nombres=(api-gateway auth-ms catalog-ms orders-ms payments-ms scheduled-jobs)
  local contextos=("P4/api-gateway" "P4/auth-ms" "P4/catalog-ms" "P4/orders-ms" "P4/payments-ms" "P5/cronjobs")
  local indice imagen contexto
  for indice in "${!nombres[@]}"; do
    imagen="$repositorio_dockerhub:${nombres[$indice]}-$etiqueta"
    contexto="$raiz_repositorio/${contextos[$indice]}"
    "$herramienta_contenedores" build --pull -t "$imagen" "$contexto"
    "$herramienta_contenedores" push "$imagen"
  done
  "$herramienta_contenedores" logout
  unset DOCKERHUB_TOKEN
  trap - EXIT
}

preparar_secretos() {
  local ruta_salida="${1:-$raiz_p6/valores-secretos.yaml}"
  [[ ! -e "$ruta_salida" ]] || { echo "El archivo ya existe: $ruta_salida" >&2; exit 1; }
  confirmar_herramienta openssl
  local clave_bd clave_rabbit clave_jwt clave_cifrado usuario_rabbit url_rabbit
  clave_bd="$(crear_clave_hexadecimal 32)"
  clave_rabbit="$(crear_clave_hexadecimal 24)"
  clave_jwt="$(crear_clave_hexadecimal 48)"
  clave_cifrado="$(crear_clave_hexadecimal 48)"
  usuario_rabbit="usuario_p6"
  url_rabbit="amqp://${usuario_rabbit}:${clave_rabbit}@sa-platform-rabbitmq:5672/"
  umask 077
  cat > "$ruta_salida" <<EOF
secrets:
  auth-ms-secret: {DB_PASSWORD: "$clave_bd", JWT_SECRET: "$clave_jwt", ENCRYPTION_KEY: "$clave_cifrado"}
  catalog-ms-secret: {DB_PASSWORD: "$clave_bd"}
  orders-ms-secret: {DB_PASSWORD: "$clave_bd", RABBITMQ_URL: "$url_rabbit"}
  payments-ms-secret: {DB_PASSWORD: "$clave_bd", RABBITMQ_URL: "$url_rabbit"}
  scheduled-jobs-secret: {DB_PASSWORD: "$clave_bd", RABBITMQ_URL: "$url_rabbit"}
rabbitmq:
  auth: {username: "$usuario_rabbit", password: "$clave_rabbit"}
auth-database:
  auth: {password: "$clave_bd"}
catalog-database:
  auth: {password: "$clave_bd"}
orders-database:
  auth: {password: "$clave_bd"}
payments-database:
  auth: {password: "$clave_bd"}
EOF
  echo "Secretos creados en $ruta_salida"
}

desplegar_plataforma() {
  local usuario_dockerhub="${1:-jeremy142004}"
  local etiqueta="${2:-1.0.0}"
  local ruta_chart="$raiz_p6/charts/sa-platform"
  local ruta_secretos="$raiz_p6/valores-secretos.yaml"
  local repositorio_dockerhub="$usuario_dockerhub/sa-p6"
  confirmar_herramienta helm
  confirmar_herramienta kubectl
  [[ -f "$ruta_secretos" ]] || { echo "Ejecute primero: $0 preparar-secretos" >&2; exit 1; }
  helm dependency build "$ruta_chart"
  helm upgrade --install sa-platform "$ruta_chart" --namespace default --values "$ruta_secretos" \
    --set "api-gateway.image.repository=$repositorio_dockerhub" --set "api-gateway.image.tag=api-gateway-$etiqueta" \
    --set "auth-ms.image.repository=$repositorio_dockerhub" --set "auth-ms.image.tag=auth-ms-$etiqueta" \
    --set "catalog-ms.image.repository=$repositorio_dockerhub" --set "catalog-ms.image.tag=catalog-ms-$etiqueta" \
    --set "orders-ms.image.repository=$repositorio_dockerhub" --set "orders-ms.image.tag=orders-ms-$etiqueta" \
    --set "payments-ms.image.repository=$repositorio_dockerhub" --set "payments-ms.image.tag=payments-ms-$etiqueta" \
    --set "scheduled-jobs.image.repository=$repositorio_dockerhub" --set "scheduled-jobs.image.tag=scheduled-jobs-$etiqueta" \
    --set-file "auth-database.primary.initdb.scripts.init\.sql=$raiz_repositorio/P4/auth-ms/database/init.sql" \
    --set-file "catalog-database.primary.initdb.scripts.init\.sql=$raiz_repositorio/P4/catalog-ms/database/init.sql" \
    --set-file "orders-database.primary.initdb.scripts.init\.sql=$raiz_repositorio/P4/orders-ms/database/init.sql" \
    --set-file "payments-database.primary.initdb.scripts.init\.sql=$raiz_repositorio/P4/payments-ms/database/init.sql" \
    --atomic --wait --timeout 15m
  kubectl -n sa-p6 get pods,svc,pvc
}

verificar_despliegue() {
  confirmar_herramienta kubectl
  kubectl -n sa-p6 wait --for=condition=Ready pods --all --timeout=10m
  local direccion_publica=""
  local intento
  for intento in {1..30}; do
    direccion_publica="$(kubectl -n sa-p6 get service api-gateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || true)"
    [[ -n "$direccion_publica" ]] && break
    sleep 10
  done
  [[ -n "$direccion_publica" ]] || { echo "No se obtuvo IP publica." >&2; exit 1; }
  curl --fail --show-error "http://${direccion_publica}:4000/health"
  echo
  kubectl get nodes -o wide
  kubectl -n sa-p6 get pods,svc,pvc,hpa,pdb
}

eliminar_infraestructura() {
  local region="${1:-us-east-1}"
  local nombre_cluster="${2:-sa-p6}"
  helm uninstall sa-platform --namespace default --wait 2>/dev/null || true
  kubectl wait --for=delete service/api-gateway --namespace sa-p6 --timeout=5m 2>/dev/null || true
  aws eks delete-addon --cluster-name "$nombre_cluster" --addon-name aws-ebs-csi-driver --region "$region" >/dev/null 2>&1 || true
  eksctl delete iamserviceaccount \
    --name ebs-csi-controller-sa \
    --namespace kube-system \
    --cluster "$nombre_cluster" \
    --region "$region" \
    --wait 2>/dev/null || true
  eksctl delete cluster --name "$nombre_cluster" --region "$region" --wait
}

mostrar_uso() {
  cat <<'EOF'
Uso:
  ./scripts/p6.sh crear-infraestructura [REGION] [CLUSTER] [TIPO_INSTANCIA]
  ./scripts/p6.sh publicar-imagenes [USUARIO_DOCKERHUB] [ETIQUETA]
  ./scripts/p6.sh preparar-secretos
  ./scripts/p6.sh desplegar [USUARIO_DOCKERHUB] [ETIQUETA]
  ./scripts/p6.sh verificar
  ./scripts/p6.sh eliminar [REGION] [CLUSTER]
EOF
}

accion="${1:-}"
shift || true
case "$accion" in
  crear-infraestructura) crear_infraestructura "$@" ;;
  publicar-imagenes) publicar_imagenes "$@" ;;
  preparar-secretos) preparar_secretos "$@" ;;
  desplegar) desplegar_plataforma "$@" ;;
  verificar) verificar_despliegue "$@" ;;
  eliminar) eliminar_infraestructura "$@" ;;
  *) mostrar_uso; exit 1 ;;
esac
