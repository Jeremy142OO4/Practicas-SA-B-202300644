# Informe de recuperación ante desastre

## 1. Incidente

Se simuló la pérdida del cluster de la Práctica 9 después de contar con un
respaldo Velero completado en el bucket externo de GCS.

## 2. Detección

Se comprobó que el cluster o sus nodos ya no estaban disponibles mediante
`kubectl get nodes` y se registró la hora de inicio.

## 3. Acciones tomadas

1. Se confirmó que el último backup estaba en estado `Completed`.
2. Se conservó el bucket de Velero y el estado remoto de Terraform.
3. Se ejecutó el bootstrap de P9.
4. ArgoCD volvió a instalar los componentes definidos en GitOps.
5. Se ejecutó o verificó la restauración de la aplicación y sus datos.

## 4. Resultado

Se verificó el estado `Synced/Healthy` de ArgoCD, los pods de `sa-p9`, el
StatefulSet de PostgreSQL, el PVC y el dato de prueba persistente.

## 5. RTO y RPO

| Medición | Hora | Resultado |
|---|---|---|
| Incidente detectado | Pendiente de medición |  |
| Último backup válido | Pendiente de medición |  |
| Inicio de reconstrucción | Pendiente de medición |  |
| ArgoCD sincronizado | Pendiente de medición |  |
| Aplicación saludable | Pendiente de medición |  |
| RTO total | Pendiente de medición |  |
| RPO observado | Pendiente de medición |  |

## 6. Lecciones aprendidas

El estado remoto, el bucket externo de Velero y los secretos externos reducen
la dependencia del cluster original. La recuperación debe ejecutarse primero
en un entorno controlado y con un backup confirmado antes de destruir recursos.
