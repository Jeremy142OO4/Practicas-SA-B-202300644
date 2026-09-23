# Informe de recuperacion ante desastre

## 1. Incidente

Se probo la recuperacion de P9 en un namespace separado. No se destruyo el
namespace activo `sa-p9`; de esta forma la prueba no interrumpio la aplicacion.

## 2. Deteccion

Se comprobo el estado del cluster con `kubectl get nodes` y el estado de las
aplicaciones con `kubectl get applications -n argocd`.

## 3. Acciones tomadas

1. Se creo el respaldo Velero `p9-manual` en el bucket externo de GCS.
2. El respaldo termino en `Completed` con 208 de 208 recursos.
3. Se creo `p9-prueba` desde ese respaldo y se uso el mapeo
   `sa-p9:sa-p9-restored`.
4. Velero recupero el PVC de PostgreSQL mediante Kopia.
5. Se verifico el dato `p9-inicial|dato-real-de-prueba` en la base restaurada.
6. Se elimino el namespace temporal y el cluster regreso a 2 nodos.

## 4. Resultado

Al finalizar:

- ArgoCD: todas las aplicaciones `Synced` y `Healthy`.
- Aplicacion: API Gateway, RabbitMQ, PostgreSQL y las 2 replicas de
  `orders-ms` en `Running/Ready`.
- Backup: `p9-manual`, `Completed`, 208/208 recursos.
- Restore: `p9-prueba`, `Completed`, 86/86 recursos y PVC `Bound`.
- Dato validado: `p9-inicial|dato-real-de-prueba`.

## 5. RTO y RPO observados

| Medicion | Resultado |
|---|---|
| Inicio del restore aislado | 2026-09-23 17:20:44 UTC |
| Restore completado | 2026-09-23 17:32:59 UTC |
| RTO de la prueba aislada | 12 min 15 s |
| Backup valido | 2026-09-23 17:18:02 UTC |
| RPO observado | 0 perdida del dato de prueba |

Este RTO corresponde a la restauracion aislada. No se presenta como el RTO
de una destruccion completa del cluster porque esa medicion no se ejecuto.

## 6. Lecciones aprendidas

El estado remoto de Terraform, el bucket externo de Velero y Secret Manager
permiten reconstruir la plataforma sin depender de los recursos temporales del
cluster. Antes de destruir un cluster se debe confirmar un backup `Completed`.
