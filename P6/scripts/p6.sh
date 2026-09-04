#!/usr/bin/env bash
set -Eeuo pipefail

directorio_script="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
raiz_p6="$(cd "$directorio_script/.." && pwd)"
raiz_repositorio="$(cd "$directorio_script/../.." && pwd)"
usuario_dockerhub="jeremy142004"
repositorio_dockerhub="${usuario_dockerhub}/sa-p6"

confirmar_herramienta() {
  local nombre="$1"
  command -v "$nombre" >/dev/null 2>&1 || { echo "Falta la herramienta: $nombre" >&2; exit 1; }
}

crear_clave_hexadecimal() {
  local cantidad_bytes="${1:-32}"
  openssl rand -hex "$cantidad_bytes"
}

leer_variable_env() {
  local ruta_archivo="$1"
  local nombre_variable="$2"
  local linea valor
  [[ -f "$ruta_archivo" ]] || { echo "No existe $ruta_archivo. Copie primero su .env.example." >&2; return 1; }
  linea="$(grep -m1 -E "^[[:space:]]*${nombre_variable}=" "$ruta_archivo" || true)"
  [[ -n "$linea" ]] || { echo "Falta ${nombre_variable} en $ruta_archivo." >&2; return 1; }
  valor="${linea#*=}"
  valor="${valor%$'\r'}"
  if [[ "$valor" == \"*\" && "$valor" == *\" ]]; then valor="${valor:1:${#valor}-2}"; fi
  if [[ "$valor" == \'*\' && "$valor" == *\' ]]; then valor="${valor:1:${#valor}-2}"; fi
  printf '%s' "$valor"
}

obtener_proyecto_gcp() {
  local id_proyecto
  id_proyecto="$(gcloud config get-value project 2>/dev/null)"
  [[ -n "$id_proyecto" && "$id_proyecto" != "(unset)" ]] || {
    echo "Configure primero el proyecto: gcloud config set project ID_PROYECTO" >&2
    return 1
  }
  printf '%s' "$id_proyecto"
}

ejecutar_psql() {
  local url_base_datos="$1"
  shift
  if [[ "$herramienta_contenedores" == "psql" ]]; then
    psql "$url_base_datos" "$@"
  else
    "$herramienta_contenedores" run --rm -i docker.io/library/postgres:18-alpine psql "$url_base_datos" "$@"
  fi
}

crear_infraestructura() {
  local zona="${1:-us-central1-a}"
  local nombre_cluster="${2:-sa-p6}"
  local tipo_instancia="${3:-e2-standard-2}"
  local id_proyecto
  confirmar_herramienta gcloud
  confirmar_herramienta kubectl
  id_proyecto="$(obtener_proyecto_gcp)"

  gcloud services enable container.googleapis.com compute.googleapis.com --project "$id_proyecto"
  if ! gcloud container clusters describe "$nombre_cluster" --zone "$zona" --project "$id_proyecto" >/dev/null 2>&1; then
    gcloud container clusters create "$nombre_cluster" \
      --project "$id_proyecto" \
      --zone "$zona" \
      --release-channel regular \
      --machine-type "$tipo_instancia" \
      --num-nodes 1 \
      --disk-type pd-balanced \
      --disk-size 30 \
      --enable-ip-alias \
      --addons GcePersistentDiskCsiDriver
  fi
  gcloud container clusters get-credentials "$nombre_cluster" --zone "$zona" --project "$id_proyecto"
  kubectl get nodes -o wide
  kubectl get csidriver pd.csi.storage.gke.io
}

publicar_imagenes() {
  local etiqueta="${1:-1.0.0}"
  local herramienta_contenedores
  if command -v podman >/dev/null 2>&1; then herramienta_contenedores=podman
  elif command -v docker >/dev/null 2>&1; then herramienta_contenedores=docker
  else herramienta_contenedores=docker
  fi
  confirmar_herramienta "$herramienta_contenedores"

  echo "El repositorio ${repositorio_dockerhub} debe existir y estar configurado como privado."
  "$herramienta_contenedores" login -u "$usuario_dockerhub"

  local nombres=(api-gateway auth-ms catalog-ms orders-ms payments-ms scheduled-jobs)
  local contextos=("P4/api-gateway" "P4/auth-ms" "P4/catalog-ms" "P4/orders-ms" "P4/payments-ms" "P5/cronjobs")
  local indice imagen contexto
  for indice in "${!nombres[@]}"; do
    imagen="$repositorio_dockerhub:${nombres[$indice]}-$etiqueta"
    contexto="$raiz_repositorio/${contextos[$indice]}"
    "$herramienta_contenedores" build --pull -f "$contexto/Dockerfile.prod" -t "$imagen" "$contexto"
    "$herramienta_contenedores" push "$imagen"
  done
  "$herramienta_contenedores" logout
}

preparar_secretos() {
  local ruta_salida="${1:-$raiz_p6/valores-secretos.yaml}"
  [[ ! -e "$ruta_salida" ]] || { echo "El archivo ya existe: $ruta_salida" >&2; exit 1; }
  confirmar_herramienta openssl
  local clave_rabbit clave_jwt clave_cifrado usuario_rabbit url_rabbit
  local url_auth url_catalog url_orders url_payments url_cronjobs url_neon herramienta_contenedores
  if command -v podman >/dev/null 2>&1; then herramienta_contenedores=podman; else herramienta_contenedores=docker; fi
  confirmar_herramienta "$herramienta_contenedores"
  url_auth="$(leer_variable_env "$raiz_repositorio/P4/auth-ms/.env" DATABASE_URL)"
  url_catalog="$(leer_variable_env "$raiz_repositorio/P4/catalog-ms/.env" DATABASE_URL)"
  url_orders="$(leer_variable_env "$raiz_repositorio/P4/orders-ms/.env" DATABASE_URL)"
  url_payments="$(leer_variable_env "$raiz_repositorio/P4/payments-ms/.env" DATABASE_URL)"
  url_cronjobs="$(leer_variable_env "$raiz_repositorio/P5/cronjobs/.env" DATABASE_URL)"
  for url_neon in "$url_auth" "$url_catalog" "$url_orders" "$url_payments" "$url_cronjobs"; do
    [[ "$url_neon" == *"sslmode=require"* ]] || { echo "Cada URL de Neon debe incluir sslmode=require." >&2; exit 1; }
  done
  clave_rabbit="$(crear_clave_hexadecimal 24)"
  clave_jwt="$(crear_clave_hexadecimal 48)"
  clave_cifrado="$(crear_clave_hexadecimal 48)"
  usuario_rabbit="usuario_p6"
  url_rabbit="amqp://${usuario_rabbit}:${clave_rabbit}@sa-platform-rabbitmq:5672/"
  umask 077
  cat > "$ruta_salida" <<EOF
secrets:
  auth-ms-secret:
    DATABASE_URL: >-
      $url_auth
    JWT_SECRET: "$clave_jwt"
    ENCRYPTION_KEY: "$clave_cifrado"
  catalog-ms-secret:
    DATABASE_URL: >-
      $url_catalog
  orders-ms-secret:
    DATABASE_URL: >-
      $url_orders
    RABBITMQ_URL: "$url_rabbit"
  payments-ms-secret:
    DATABASE_URL: >-
      $url_payments
    RABBITMQ_URL: "$url_rabbit"
  scheduled-jobs-secret:
    DATABASE_URL: >-
      $url_cronjobs
    RABBITMQ_URL: "$url_rabbit"
rabbitmq:
  auth: {username: "$usuario_rabbit", password: "$clave_rabbit"}
EOF
  ejecutar_psql "$url_auth" < "$raiz_repositorio/P4/auth-ms/database/init.sql"
  ejecutar_psql "$url_catalog" < "$raiz_repositorio/P4/catalog-ms/database/init.sql"
  ejecutar_psql "$url_orders" < "$raiz_repositorio/P4/orders-ms/database/init.sql"
  ejecutar_psql "$url_payments" < "$raiz_repositorio/P4/payments-ms/database/init.sql"
  ejecutar_psql "$url_cronjobs" -c '
    CREATE SCHEMA IF NOT EXISTS scheduled_jobs;
    CREATE TABLE IF NOT EXISTS scheduled_jobs.ejecuciones (
      id_ejecucion BIGSERIAL PRIMARY KEY,
      fecha_hora TIMESTAMPTZ NOT NULL,
      zona_horaria VARCHAR(20) NOT NULL,
      carne VARCHAR(20) NOT NULL
    );
    CREATE TABLE IF NOT EXISTS scheduled_jobs.resumenes_hora (
      id_resumen BIGSERIAL PRIMARY KEY,
      event_id UUID NOT NULL,
      hora TIMESTAMP NOT NULL,
      cantidad INTEGER NOT NULL,
      recibido_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
      CONSTRAINT unique_evento_hora UNIQUE (event_id, hora)
    );'
  unset url_auth url_catalog url_orders url_payments url_cronjobs
  echo "Neon inicializado y secretos creados en $ruta_salida"
}

crear_secreto_dockerhub() {
  local token_dockerhub
  read -r -s -p "Token nuevo de Docker Hub: " token_dockerhub
  echo
  [[ -n "$token_dockerhub" ]] || { echo "El token no puede estar vacio." >&2; return 1; }
  kubectl create namespace sa-p6 --dry-run=client -o yaml | kubectl apply -f -
  kubectl -n sa-p6 create secret docker-registry dockerhub-cred \
    --docker-server=https://index.docker.io/v1/ \
    --docker-username="$usuario_dockerhub" \
    --docker-password="$token_dockerhub" \
    --dry-run=client -o yaml | kubectl apply -f -
  unset token_dockerhub
}

desplegar_plataforma() {
  local etiqueta="${1:-1.0.0}"
  local ruta_chart="$raiz_p6/charts/sa-platform"
  local ruta_secretos="$raiz_p6/valores-secretos.yaml"
  confirmar_herramienta helm
  confirmar_herramienta kubectl
  [[ -f "$ruta_secretos" ]] || { echo "Ejecute primero: $0 preparar-secretos" >&2; exit 1; }
  crear_secreto_dockerhub
  helm dependency build "$ruta_chart"
  helm upgrade --install sa-platform "$ruta_chart" --namespace default --values "$ruta_secretos" \
    --set "api-gateway.image.repository=$repositorio_dockerhub" --set "api-gateway.image.tag=api-gateway-$etiqueta" \
    --set "auth-ms.image.repository=$repositorio_dockerhub" --set "auth-ms.image.tag=auth-ms-$etiqueta" \
    --set "catalog-ms.image.repository=$repositorio_dockerhub" --set "catalog-ms.image.tag=catalog-ms-$etiqueta" \
    --set "orders-ms.image.repository=$repositorio_dockerhub" --set "orders-ms.image.tag=orders-ms-$etiqueta" \
    --set "payments-ms.image.repository=$repositorio_dockerhub" --set "payments-ms.image.tag=payments-ms-$etiqueta" \
    --set "scheduled-jobs.image.repository=$repositorio_dockerhub" --set "scheduled-jobs.image.tag=scheduled-jobs-$etiqueta" \
    --atomic --wait --timeout 15m
  kubectl -n sa-p6 get pods,svc,pvc
}

verificar_despliegue() {
  confirmar_herramienta kubectl
  kubectl -n sa-p6 wait --for=condition=Ready pods --all --timeout=10m
  local direccion_publica="" intento
  for intento in {1..30}; do
    direccion_publica="$(kubectl -n sa-p6 get service api-gateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || true)"
    [[ -n "$direccion_publica" ]] && break
    sleep 10
  done
  [[ -n "$direccion_publica" ]] || { echo "No se obtuvo IP publica." >&2; exit 1; }
  curl --fail --show-error "http://${direccion_publica}:4000/health"
  echo
  kubectl get nodes -o wide
  kubectl -n sa-p6 get pods,svc,pvc
}

eliminar_infraestructura() {
  local zona="${1:-us-central1-a}"
  local nombre_cluster="${2:-sa-p6}"
  local id_proyecto
  confirmar_herramienta gcloud
  id_proyecto="$(obtener_proyecto_gcp)"
  helm uninstall sa-platform --namespace default --wait 2>/dev/null || true
  kubectl wait --for=delete service/api-gateway --namespace sa-p6 --timeout=5m 2>/dev/null || true
  gcloud container clusters delete "$nombre_cluster" --zone "$zona" --project "$id_proyecto" --quiet
}

mostrar_uso() {
  cat <<'EOF'
Uso:
  ./scripts/p6.sh crear-infraestructura [ZONA] [CLUSTER] [TIPO_INSTANCIA]
  ./scripts/p6.sh publicar-imagenes [ETIQUETA]
  ./scripts/p6.sh preparar-secretos
  ./scripts/p6.sh desplegar [ETIQUETA]
  ./scripts/p6.sh verificar
  ./scripts/p6.sh eliminar [ZONA] [CLUSTER]
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
