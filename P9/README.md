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
