# Runbook de recuperacion P9

## 1. Preparar el entorno

1. Tener `gcloud`, `terraform`, `kubectl` y `helm` instalados.
2. Clonar o actualizar el repositorio GitOps.
3. Copiar `P9/terraform/infra/terraform.tfvars.example` a
   `P9/terraform/infra/terraform.tfvars` y completar los valores.
4. No guardar `terraform.tfvars` en Git.

## 2. Crear la plataforma

Desde la raiz del repositorio:

```bash
chmod +x P9/scripts/bootstrap.sh
P9/scripts/bootstrap.sh
```

El script crea el estado remoto, levanta GKE con Terraform, instala ArgoCD y
aplica solamente la Application raiz. ArgoCD levanta el resto desde GitOps.

## 3. Verificar el funcionamiento

```bash
kubectl get applications -n argocd
kubectl get pods -n sa-p9
kubectl get pvc -n sa-p9
kubectl get statefulset orders-db -n sa-p9
kubectl get deployment orders-ms -n sa-p9
kubectl get rollout api-gateway -n sa-p9
kubectl get backupstoragelocation -n velero
kubectl get schedule -n velero
```

El dato `dr_marker` sirve para comprobar que PostgreSQL conservo la
informacion despues de una restauracion.

## 4. Crear y comprobar un backup

Si esta instalado el cliente de Velero:

```bash
velero backup create p9-manual \
  --include-namespaces sa-p9 \
  --default-volumes-to-fs-backup \
  --wait
velero backup get
```

Alternativa usando solamente `kubectl`:

```bash
kubectl apply -f - <<'EOF'
apiVersion: velero.io/v1
kind: Backup
metadata:
  name: p9-manual
  namespace: velero
spec:
  includedNamespaces: [sa-p9]
  defaultVolumesToFsBackup: true
  storageLocation: default
  ttl: 168h
EOF
kubectl get backup p9-manual -n velero
```

Antes de reconstruir, el estado debe ser `Completed`.

## 5. Probar la restauracion en otro namespace

La prueba automatizada usada en la verificacion reciente esta en:

```bash
bash P9/scripts/probar_restauracion.sh
```

Se ejecuta desde la raiz del repositorio. Crea un marcador nuevo, completa un
backup, elimina solamente ese marcador de la base activa y lo recupera en
`sa-p9-restored`. Compara el contenido y verifica que se use otro volumen.
No sobrescribe recursos existentes ni destruye el cluster. La [salida real](../evidencias/restauracion.txt)
documenta el resultado. Si la prueba ya existe, el script se detiene para
conservarla; no borrar backups ni repetir sobre un namespace con datos.

El resto de esta seccion describe la alternativa manual historica:

```bash
velero restore create p9-prueba \
  --from-backup p9-manual \
  --namespace-mappings sa-p9:sa-p9-restored \
  --wait
kubectl get pods,pvc -n sa-p9-restored
```

Alternativa con `kubectl`:

```bash
kubectl apply -f - <<'EOF'
apiVersion: velero.io/v1
kind: Restore
metadata:
  name: p9-prueba
  namespace: velero
spec:
  backupName: p9-manual
  namespaceMapping:
    sa-p9: sa-p9-restored
  restorePVs: true
EOF
kubectl get restore p9-prueba -n velero
```

La prueba no reemplaza el namespace activo. Para validar el dato restaurado:

```bash
kubectl exec -n sa-p9-restored orders-db-0 -- \
  psql -U p9user -d orders_db -tAc \
  'select id,value from dr_marker'
```

Al terminar la demostracion se puede limpiar el namespace temporal:

```bash
kubectl delete namespace sa-p9-restored
```

## 6. Probar perdida de un nodo

Con P10 en uso, la prueba segura de reprogramacion de pedidos es:

```bash
python3 P9/scripts/probar_nodo.py
```

El script comprueba que haya dos replicas Ready, cordona un nodo, drena solo
los pods de pedidos y consulta una operacion de lectura a traves del gateway.
Al terminar, incluso si detecta errores, vuelve a habilitar el nodo.
No desalojar Tempo, Loki o Prometheus para esta prueba: usan almacenamiento
temporal y podrian perder el historial de P10.

Pedidos espera 15 segundos en preStop y dispone de 45 segundos de terminacion:
esto da margen para retirar el endpoint antes de cerrar el servicio y su pool
de conexiones. El PDB y la anti-afinidad mantienen otra replica disponible.
Esto no convierte a PostgreSQL y RabbitMQ en servicios de alta disponibilidad.

El drenaje completo siguiente es mas disruptivo y requiere preparar una
ventana de prueba; no ejecutarlo durante la grabacion o calificacion de P10.

Con dos nodos disponibles:

```bash
kubectl get nodes
kubectl cordon NOMBRE_DEL_NODO
kubectl drain NOMBRE_DEL_NODO \
  --ignore-daemonsets \
  --delete-emptydir-data \
  --timeout=5m
kubectl get pods -n sa-p9 -o wide
kubectl uncordon NOMBRE_DEL_NODO
```

Registrar que `orders-ms` continua disponible y que Kubernetes reprograma los
pods respetando el PDB y la anti-afinidad.

El verificador compara las replicas de todas las cargas, incluso las bases y
colas. No aumentar PostgreSQL o RabbitMQ a dos sin replicacion real: crearia
instancias independientes. Revisar ese criterio con el auxiliar.

## 7. Reconstruccion controlada

Solo destruir despues de un backup `Completed` y con autorizacion:

```bash
terraform -chdir=P9/terraform/infra destroy \
  -target=module.base.google_container_node_pool.default \
  -target=module.base.google_container_cluster.p8
P9/scripts/bootstrap.sh
```

El bucket usa `force_destroy = false`, por lo que Terraform no borra los
backups accidentalmente. Al terminar, revisar ArgoCD, pods, PVC y `dr_marker`.

## 8. RTO y RPO

Registrar la hora del incidente, del ultimo backup valido, del inicio de la
reconstruccion, de la sincronizacion de ArgoCD y de la aplicacion saludable.
El informe actual registra la restauracion aislada comprobada en P9.

## 9. Verificador original

Usar [p9.conf](../p9.conf) con `APP_ROOT=p9-root`, `NS=sa-p9` y
`PVC=data-orders-db-0`. Ejecutar el script original desde la carpeta que
contenga esa configuracion. No contiene credenciales.

La version recibida usa expresiones como `${SCHED:-{}}`: con un JSON presente
agregan una llave extra y jq termina con error. Por eso puede informar que no
hay backups o schedules aunque existan. Las comprobaciones reales son:

```bash
kubectl get schedules,backups,restores -n velero
kubectl get backup p9-verificacion -n velero -o json | jq '.status'
kubectl get restore p9-verificacion -n velero -o json | jq '.status'
velero backup get
velero restore get
```

El script y su hash se conservan. Si marca esos falsos negativos, presentar
los recursos y el registro original al auxiliar, sin editar el verificador.
