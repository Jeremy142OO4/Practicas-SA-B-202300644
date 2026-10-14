# Practica 10: Observabilidad y operacion

P10 agrega observabilidad a la plataforma reutilizada de P9. No copia los
microservicios ni el cluster: el codigo de negocio continua en `P4/`, la
infraestructura base continua en `P9/` y los manifiestos de observabilidad se
administran en el repositorio GitOps.

## Objetivo

El flujo de la practica es:

```text
Solicitud -> api-gateway -> orders-ms -> catalog-ms
     |             |             |             |
     +---- trace_id y contexto W3C en los tres servicios
     +---- metricas HTTP y de negocio para Prometheus
     +---- logs JSON con trace_id

OpenTelemetry Collector -> Tempo
Prometheus -> reglas SLO y AnalysisTemplate de Argo Rollouts
Prometheus + Loki + Tempo -> Grafana
```

## Servicios instrumentados

| Servicio | Ruta del codigo | Señales |
|---|---|---|
| `api-gateway` | `P4/api-gateway` | trazas, metricas HTTP y logs JSON |
| `catalog-ms` | `P4/catalog-ms` | trazas, metricas HTTP y logs JSON |
| `orders-ms` | `P4/orders-ms` | trazas, metricas HTTP y logs JSON |

Las librerias se inicializan desde el codigo existente. El contexto W3C se
propaga entre llamadas HTTP y todos los logs de solicitud incluyen `trace_id`.
Ademas, `orders-ms` publica el contador de negocio `business_events_total`
para distinguir resultados de pago exitosos y fallidos.

## SLO definidos

| SLO | Objetivo | Ventana | Consulta |
|---|---:|---|---|
| Disponibilidad del API | 99% | 30 dias | `slo:api_availability:ratio` |
| Latencia p95 del API | 95% menor a 500 ms | 30 dias | `slo:api_latency:ratio` |

Los error budgets son 1% para disponibilidad y 5% para latencia. Las alertas
`ApiAvailabilityErrorBudgetBurn` y `ApiLatencyErrorBudgetBurn` avisan cuando
se consume mas presupuesto del permitido.

Las reglas, alertas y el AnalysisTemplate se encuentran en el repositorio
GitOps. La promocion canary consulta Prometheus mediante PromQL; no usa un Job
de prueba como criterio de promocion.

## Archivos de esta practica

| Archivo | Uso |
|---|---|
| `docs/FLUJO.md` | Explicacion del flujo y de las cuatro señales doradas |
| `docs/RUNBOOK.md` | Pasos para investigar una alerta |
| `docs/INFORME_OPERACION.md` | Informe de incidente de dos paginas |
| `docs/DIAGRAMA.md` | Diagrama de flujo de telemetria |
| `dashboards/sa-platform.json` | Dashboard de Grafana |
| `evidencias/` | Capturas y resultados de la demostracion |

## Repositorios

| Elemento | Enlace |
|---|---|
| Codigo | [Practicas-SA-B-202300644](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644) |
| GitOps | [Practicas-SA-B-202300644-gitops](https://github.com/Jeremy142OO4/Practicas-SA-B-202300644-gitops) |

## Tabla de entrega

| Requisito | Ubicacion |
|---|---|
| Servicios instrumentados | `P4/api-gateway`, `P4/catalog-ms`, `P4/orders-ms` |
| AnalysisTemplate con Prometheus | GitOps `p10/platform/analysis-templates.yaml` |
| Consultas PromQL | GitOps `p10/observability/prometheus-rules.yaml` |
| SLO y error budget | GitOps `p10/observability/prometheus-rules.yaml` |
| PrometheusRule | GitOps `p10/observability/prometheus-rules.yaml` |
| Dashboard Grafana | `dashboards/sa-platform.json` |
| AnalysisRun fallido y Rollout | `evidencias/` (se agregan durante la prueba) |
| Ruta sintoma a causa | `docs/INFORME_OPERACION.md` |
| Video y minutaje | Pendiente de colocar al grabar |

## Video

- Enlace: **PENDIENTE**
- Minutaje: **PENDIENTE**

## Estado actual

La preparacion local de P10 no levanta el cluster ni genera costos. Antes de
validar en Kubernetes se deben encender temporalmente los nodos del cluster P9,
sin modificar la configuracion base de Terraform.
