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

## Tabla de entrega - seccion 4.1

| Item | Enlace o dato |
|---|---|
| Repositorio GitOps | [Practicas-SA-B-202300644-gitops](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644-gitops), publico |
| Aplicacion raiz | `p9-root`, namespace `argocd`; [manifiesto](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644-gitops/blob/main/p9/argocd/root.yaml) |
| Bootstrap | [P9/scripts/bootstrap.sh](scripts/bootstrap.sh), ejecutar desde la raiz del repositorio |
| Backend remoto | GCS `gs://sa-p9-202300644-tfstate`, prefijos `p9/infra` y `p9/platform`; bloqueo integrado de GCS |
| Schedule de Velero | `p9-velero-p9-daily`, cada seis horas, retencion 168h; bucket `sa-202300644-p9-202300644-velero` |
| Reconstruccion cronometrada | Pendiente: la prueba publicada es una restauracion aislada, no destruccion y reconstruccion total del cluster. [Alcance del informe](docs/INFORME_DR.md) |
| Restauracion de datos | [Registro real](evidencias/restauracion.txt); Backup y Restore `p9-verificacion`, ambos Completed |
| Prueba de perdida de nodo | [Registro de reprogramacion selectiva](evidencias/reprogramacion.txt): 30/30 consultas correctas; [script](scripts/probar_nodo.py). No equivale a apagar una maquina completa |
| RTO y RPO | RTO medido de restauracion aislada: 50 s; RPO del marcador respaldado: cero datos perdidos. Los objetivos para futuras pruebas se explican en el [informe](docs/INFORME_DR.md) |
| Video | [Video y minutaje](#video-demostrativo); muestra el ensayo anterior, no reemplaza el registro nuevo |

## Resultado comprobado

- ArgoCD: plataforma y aplicaciones de P10 `Synced/Healthy`; Kyverno
  sincronizado, con su controlador de reportes reiniciandose en la ultima revision.
- Backup `p9-verificacion`: `Completed`, 141/141 recursos.
- Restore `p9-verificacion`: `Completed`, 3/3 recursos seleccionados.
- PVC restaurado: `Bound`.
- Dato recuperado: `p9-verificacion|dato-recuperado-20261015T230448Z`.
- La fila se elimino de la base activa despues del backup y se recupero en
  otro volumen. No estaba en el script de inicializacion de PostgreSQL.
- Cluster actual compartido con P10: cuatro nodos; gateway y pedidos tienen
  dos replicas. PostgreSQL y RabbitMQ siguen con una instancia cada uno.

Los datos del verificador estan en [p9.conf](p9.conf). Los scripts del auxiliar
no se modifican. La documentacion no afirma una reconstruccion completa ni
alta disponibilidad de PostgreSQL o RabbitMQ.

El [verificador original](evidencias/verificador-original.txt), sin modificaciones,
dio 43/60. Sus comprobaciones de Velero y Restore fallan al procesar JSON
([explicacion y comandos independientes](docs/RUNBOOK.md#9-verificador-original)).
La penalizacion por cargas con una sola replica si corresponde a la configuracion
actual. Estos resultados no se presentan como una calificacion aprobada.

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
