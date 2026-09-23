# Práctica 9: continuidad operativa y recuperación

Esta práctica reutiliza la plataforma de la P8 sin copiar sus microservicios,
charts ni manifiestos. La P9 agrega la reconstrucción desde cero, el estado
remoto de Terraform, los respaldos con Velero, la restauración de datos, la
resiliencia ante pérdida de nodo y la medición de RTO y RPO.

## Flujo previsto

```text
bootstrap.sh
  -> backend GCS con bloqueo
  -> Terraform reutiliza P8/terraform
  -> cluster GKE p9-202300644
  -> Terraform instala ArgoCD
  -> Application raíz en GitOps
  -> aplicaciones, Velero y políticas
  -> prueba de respaldo, destrucción y restauración
```

## Reutilización

- La infraestructura base se referencia desde `../P8/terraform` como módulo.
- Los servicios y la imagen de la aplicación se reutilizan desde P8 y GitOps.
- P9 contiene únicamente la automatización y los recursos nuevos de
  continuidad operativa.

## Archivos importantes

| Archivo | Uso |
|---|---|
| `scripts/bootstrap.sh` | Crea la infraestructura, instala ArgoCD y aplica la Application raíz. |
| `terraform/infra` | Cluster GKE, bucket de Velero, cuentas GCP y Secret Manager. |
| `terraform/platform` | Instalación de ArgoCD con Terraform. |
| `docs/DIAGRAMA_BOOTSTRAP.md` | Diagrama del flujo completo. |
| `docs/RUNBOOK.md` | Pasos de creación, backup, restauración y reconstrucción. |
| `docs/INFORME_DR.md` | Informe de incidente, RTO y RPO. |

## Repositorios y evidencias

- Código: [Practicas-SA-B-202300644](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644)
- GitOps: [Practicas-SA-B-202300644-gitops](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644-gitops)
- Application raíz: `p9/argocd/root.yaml` en el repositorio GitOps.
- Video: pendiente de agregar después de grabar la demostración.
- Capturas y resultados: `P9/evidencias/`.

## Estado

La implementación está preparada para validarse por fases. No se debe
destruir el cluster hasta contar con un respaldo Velero en estado `Completed`,
conservar el bucket externo y registrar las horas necesarias para calcular RTO
y RPO.
