# Practica 9: continuidad operativa y recuperacion

Esta practica reutiliza la plataforma de P8 sin copiar sus microservicios,
charts ni manifiestos. P9 agrega la reconstruccion con Terraform, el estado
remoto, los backups con Velero, la restauracion de datos, los secretos
externos y la resiliencia ante perdida de nodo.

## Flujo

```text
bootstrap.sh
  -> backend GCS con bloqueo
  -> Terraform crea el cluster GKE
  -> Terraform instala ArgoCD
  -> Application raiz en GitOps
  -> Rollouts, secretos, politicas y Velero
  -> backup y restauracion de la aplicacion
```

## Reutilizacion

- La infraestructura base se referencia desde `P8/terraform` como modulo.
- Los servicios y la imagen de la aplicacion se reutilizan desde P8 y GitOps.
- P9 agrega solamente los recursos de continuidad operativa.

## Archivos importantes

| Archivo | Uso |
|---|---|
| `scripts/bootstrap.sh` | Crea la infraestructura e instala ArgoCD. |
| `terraform/infra` | GKE, bucket de Velero, cuentas GCP y Secret Manager. |
| `terraform/platform` | Instalacion de ArgoCD con Terraform. |
| `docs/DIAGRAMA_BOOTSTRAP.md` | Diagrama del flujo completo. |
| `docs/RUNBOOK.md` | Pasos de creacion, backup, restore y reconstruccion. |
| `docs/INFORME_DR.md` | Incidente, resultado, RTO y RPO. |

## Repositorios

- Codigo: [Practicas-SA-B-202300644](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644)
- GitOps: [Practicas-SA-B-202300644-gitops](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644-gitops)
- Application raiz: `p9/argocd/root.yaml` en GitOps.

## Resultado comprobado

- ArgoCD: todas las aplicaciones `Synced/Healthy`.
- Backup `p9-manual`: `Completed`, 208/208 recursos.
- Restore `p9-prueba`: `Completed`, 86/86 recursos.
- PVC restaurado: `Bound`.
- Dato restaurado: `p9-inicial|dato-real-de-prueba`.
- Cluster final: 2 nodos y la aplicacion principal `Running/Ready`.

## Video demostrativo

- Enlace: [**Video demostrativo de la P9**](https://drive.google.com/file/d/1zS4KpNyHk5F95h6WW3VNfJK2ZDrV8X81/view?usp=sharing)
- Duracion aproximada: **5 minutos con 36 segundos**.
- 00:00-00:45: Presentacion del README y estructura de la carpeta P9.
- 00:45-01:10: Informe de recuperacion y runbook.
- 01:10-01:50: Script `bootstrap.sh` y flujo de creacion automatica.
- 01:50-02:20: Cluster GKE P9 y conexion desde Cloud Shell.
- 02:20-03:20: Verificacion de nodos, pods, PVC, servicios y rollout.
- 03:20-03:50: Aplicaciones de ArgoCD en estado `Synced/Healthy`.
- 03:50-04:15: Backup Storage Location, Schedule y backup de Velero.
- 04:15-04:30: Detalle de la restauracion con Velero.
- 04:30-05:05: Resultados de la restauracion y datos recuperados.
- 05:05-05:36: Informe final con RTO, RPO y lecciones aprendidas.
