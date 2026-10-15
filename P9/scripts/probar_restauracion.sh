#!/usr/bin/env bash
# Prueba real y aislada: respaldo, perdida de un dato de prueba y recuperacion.
# No destruye el cluster ni elimina datos de negocio.
set -euo pipefail

CTX="gke_sa-202300644_us-central1-a_p9-202300644"
NS="sa-p9"
DEST="sa-p9-restored"
BACKUP="p9-verificacion"
RESTORE="p9-verificacion"
MARKER="p9-verificacion"
VALUE="dato-recuperado-$(date -u +%Y%m%dT%H%M%SZ)"

utc() { date -u +%Y-%m-%dT%H:%M:%SZ; }
sql() { kubectl exec -n "$1" orders-db-0 -- psql -v ON_ERROR_STOP=1 -U p9user -d orders_db -tAc "$2"; }
wait_phase() {
  local kind="$1" name="$2" deadline=$((SECONDS + 900)) phase=""
  while (( SECONDS < deadline )); do
    phase=$(kubectl get "$kind" "$name" -n velero -o jsonpath='{.status.phase}')
    printf '%s %s %s: %s\n' "$(utc)" "$kind" "$name" "$phase"
    case "$phase" in
      Completed) return 0 ;;
      Failed|FailedValidation|PartiallyFailed) return 1 ;;
    esac
    sleep 15
  done
  echo "Tiempo agotado; revisar el recurso antes de reintentar." >&2
  return 1
}

[[ $(kubectl config current-context) == "$CTX" ]] || { echo "Contexto incorrecto" >&2; exit 1; }
for entry in "namespace/$DEST" "backup.velero.io/$BACKUP" "restore.velero.io/$RESTORE"; do
  if kubectl get "$entry" -n velero --ignore-not-found -o name | grep -q .; then
    echo "Ya existe $entry. No se sobrescribe una prueba anterior." >&2; exit 1
  fi
done
[[ $(sql "$NS" "select count(*) from dr_marker where id='$MARKER'") == 0 ]] || {
  echo "El marcador ya existe: no se modificara." >&2; exit 1;
}

echo "Inicio UTC: $(utc)"
echo "Dato de prueba: $MARKER|$VALUE"
sql "$NS" "insert into dr_marker(id,value) values ('$MARKER','$VALUE')"
sql "$NS" "select id,value from dr_marker where id='$MARKER'"

kubectl create -f - <<EOF
apiVersion: velero.io/v1
kind: Backup
metadata:
  name: $BACKUP
  namespace: velero
spec:
  includedNamespaces: [$NS]
  defaultVolumesToFsBackup: true
  snapshotVolumes: false
  storageLocation: default
  ttl: 168h
EOF
wait_phase backup "$BACKUP"
kubectl get backup "$BACKUP" -n velero -o json | jq '{name:.metadata.name,phase:.status.phase,start:.status.startTimestamp,end:.status.completionTimestamp,progress:.status.progress}'

# Se elimina solamente la fila recien creada, despues de confirmar el respaldo.
echo "Perdida del dato de prueba UTC: $(utc)"
sql "$NS" "delete from dr_marker where id='$MARKER'"
[[ $(sql "$NS" "select count(*) from dr_marker where id='$MARKER'") == 0 ]]
echo "Marcador ausente en la base activa: confirmado"

kubectl create namespace "$DEST"
# Dependencias de arranque; no se imprimen ni se guardan secretos en archivos.
for item in secret/orders-ms-secrets configmap/orders-db-init; do
  kubectl get "$item" -n "$NS" -o json | jq --arg ns "$DEST" '
    del(.metadata.uid,.metadata.resourceVersion,.metadata.creationTimestamp,
        .metadata.managedFields,.metadata.ownerReferences,.metadata.annotations)
    | .metadata.namespace=$ns' | kubectl create -f -
done

START=$(date +%s)
echo "Inicio restauracion UTC: $(utc)"
kubectl create -f - <<EOF
apiVersion: velero.io/v1
kind: Restore
metadata:
  name: $RESTORE
  namespace: velero
spec:
  backupName: $BACKUP
  includedNamespaces: [$NS]
  namespaceMapping:
    $NS: $DEST
  includedResources: [pods, persistentvolumeclaims, persistentvolumes]
  labelSelector:
    matchLabels:
      app.kubernetes.io/name: orders-db
  restorePVs: true
EOF
wait_phase restore "$RESTORE"
kubectl wait --for=condition=Ready pod/orders-db-0 -n "$DEST" --timeout=300s
RECOVERED=$(sql "$DEST" "select id,value from dr_marker where id='$MARKER'")
[[ "$RECOVERED" == "$MARKER|$VALUE" ]] || { echo "El dato recuperado no coincide" >&2; exit 1; }
SOURCE_PV=$(kubectl get pvc data-orders-db-0 -n "$NS" -o jsonpath='{.spec.volumeName}')
TARGET_PV=$(kubectl get pvc data-orders-db-0 -n "$DEST" -o jsonpath='{.spec.volumeName}')
[[ -n "$TARGET_PV" && "$SOURCE_PV" != "$TARGET_PV" ]]
END=$(date +%s)
echo "Verificacion UTC: $(utc)"
echo "Dato recuperado: $RECOVERED"
echo "RTO de restauracion aislada hasta verificar dato: $((END - START)) segundos"
echo "RPO para el marcador respaldado: 0 datos perdidos"
echo "PV origen: $SOURCE_PV"
echo "PV recuperado: $TARGET_PV"
kubectl get restore "$RESTORE" -n velero -o json | jq '{name:.metadata.name,phase:.status.phase,backup:.spec.backupName,start:.status.startTimestamp,end:.status.completionTimestamp,progress:.status.progress,warnings:.status.warnings,errors:.status.errors}'
kubectl get pods,pvc -n "$DEST"
echo "Prueba completada. Se conserva el Restore en velero."
echo "El namespace aislado puede eliminarse despues de guardar esta evidencia."
