# Runbook de recuperación P9

## 1. Preparar el entorno

1. Tener `gcloud`, `terraform`, `kubectl` y `helm` instalados.
2. Clonar o actualizar el repositorio GitOps en `~/Practicas-SA-B-202300644-gitops`.
3. Copiar `P9/terraform/infra/terraform.tfvars.example` a
   `P9/terraform/infra/terraform.tfvars` y colocar valores propios.
4. No guardar el archivo `terraform.tfvars` en Git.

## 2. Crear la plataforma

Desde la raíz del repositorio:

```bash
chmod +x P9/scripts/bootstrap.sh
P9/scripts/bootstrap.sh
```

El script crea o verifica el estado remoto en GCS, crea el cluster con
Terraform, instala ArgoCD y aplica únicamente la Application raíz. ArgoCD
levanta el resto desde GitOps.

## 3. Verificar el funcionamiento

```bash
kubectl get applications -n argocd
kubectl get pods -n sa-p9
kubectl get pvc -n sa-p9
kubectl get statefulset orders-db -n sa-p9
kubectl get deployment orders-ms -n sa-p9
kubectl get rollout api-gateway -n sa-p9
kubectl get backupstoragelocation -n velero
velero schedule get
```

La tabla `dr_marker` permite verificar que los datos reales de la base fueron
conservados después de una restauración.

## 4. Crear y comprobar un respaldo

```bash
velero backup create p9-manual --include-namespaces sa-p9 --default-volumes-to-fs-backup
velero backup describe p9-manual --details
velero backup get
```

Antes de reconstruir el cluster, confirmar que el respaldo termine en estado
`Completed` y guardar la hora de inicio y de finalización.

## 5. Probar la restauración en un namespace separado

```bash
velero restore create p9-prueba --from-backup p9-manual --namespace-mappings sa-p9:sa-p9-restored
velero restore describe p9-prueba --details
kubectl get pods,pvc -n sa-p9-restored
```

La prueba no reemplaza el namespace activo y permite verificar que el PVC y la
información de PostgreSQL se recuperen sin interrumpir la aplicación.

## 6. Probar pérdida de un nodo

Con dos nodos disponibles, identificar uno y simular su mantenimiento:

```bash
kubectl get nodes
kubectl cordon NOMBRE_DEL_NODO
kubectl drain NOMBRE_DEL_NODO --ignore-daemonsets --delete-emptydir-data --timeout=5m
kubectl get pods -n sa-p9 -o wide
kubectl uncordon NOMBRE_DEL_NODO
```

Registrar que `orders-ms` continúa disponible y que Kubernetes reprograma los
pods respetando el PDB y la anti-afinidad.

## 7. Reconstrucción controlada

La destrucción solo se realiza después de un respaldo `Completed` y con la
autorización del responsable. Se conserva el bucket de Velero y el estado
remoto de Terraform. Después se vuelve a ejecutar `P9/scripts/bootstrap.sh`.

```bash
terraform -chdir=P9/terraform/infra destroy \
  -target=module.base.google_container_node_pool.default \
  -target=module.base.google_container_cluster.p8
P9/scripts/bootstrap.sh
```

El bucket de respaldos usa `force_destroy = false`, por lo que Terraform no
borra los respaldos accidentalmente. Al terminar, verificar ArgoCD, la
aplicación, el StatefulSet, el PVC y el dato de `dr_marker`.

## 8. RTO y RPO

Registrar en el informe:

- Hora del incidente.
- Hora del último respaldo válido.
- Hora en que comienza la reconstrucción.
- Hora en que ArgoCD queda sincronizado.
- Hora en que la aplicación vuelve a estar saludable.
- Diferencia entre el último dato confirmado y el dato recuperado.
